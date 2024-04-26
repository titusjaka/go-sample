package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"golang.org/x/sync/errgroup"

	"github.com/titusjaka/go-sample/commands/flags"
	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/kongflag"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
)

// ErrStopped is returned when the server is gracefully stopped.
var ErrStopped = errors.New("stopped")

// ServerCmd implements kong.Command for the main server command.
type ServerCmd struct {
	Postgres flags.PostgreSQL `kong:"embed"`

	Listen string `kong:"optional,default=':4040',group='HTTP Server',env=HTTP_LISTEN,help='HTTP network address'"`
	Token  string `kong:"optional,env=API_TOKEN,group='HTTP Server',help='authentication token used for inter-service communication'"`
}

// Run (ServerCmd) runs the main server command.
func (c ServerCmd) Run(kVars kong.Vars) error {
	gr, ctx := errgroup.WithContext(context.Background())

	logger := log.New()

	logger.Info(
		"⏳ starting service…",
		log.Field("version", kVars[kongflag.ServiceVersion]),
		log.Field("git_commit_sha", kVars[kongflag.GitCommitSHA]),
		log.Field("git_branch", kVars[kongflag.GitBranch]),
		log.Field("go_version", kVars[kongflag.GoVersion]),
		log.Field("config", c),
	)

	// =========================================================================
	// Init PostgreSQL Connection

	db, err := c.Postgres.OpenStdSQLDB()
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("unable to close database connection", log.Field("err", closeErr))
		}
	}()

	// ================================================
	// Apply SQL Migrations

	migrationCmd := UpCmd{
		Postgres: c.Postgres,
	}

	if err = migrationCmd.Run(); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	// =========================================================================
	// Start Private API Server

	gr.Go(func() error {
		return c.runHTTPServer(
			ctx,
			log.With(logger, log.Field("system", "http-server")),
			db,
		)
	})

	// =========================================================================
	// Wait for stop signal

	gr.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stop)

		select {
		case <-stop:
			logger.Info("🛑 ➡ received stop signal")
			return ErrStopped
		case <-ctx.Done():
			return nil
		}
	})

	if err = gr.Wait(); err != nil && !errors.Is(err, ErrStopped) {
		return fmt.Errorf("error during wait: %w", err)
	}

	logger.Info("💀 ➡ web server gracefully stopped")
	return nil
}

// runHTTPServer starts the HTTP server.
func (c ServerCmd) runHTTPServer(ctx context.Context, logger log.Logger, db *sql.DB) error {
	// =========================================================================
	// Init Chi Router

	r := chi.NewRouter()

	corsOpts := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Accept",
			"Authorization",
			"Content-Type",
		},
	})
	r.Use(corsOpts.Handler)
	r.Use(middleware.SetHeader("X-Frame-Options", "deny"))
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.NotFound(api.NewNotFoundHandler(logger))
	r.MethodNotAllowed(api.NewMethodNotAllowedHandler(logger))

	// =========================================================================
	// Init Snippets Module

	snippetStorage := snippets.NewPGStorage(db)
	snippetService := snippets.NewService(
		snippetStorage,
		log.With(logger, log.Field("service", "snippets")),
	)
	snippetRouter := snippets.MakeSnippetsHandler(snippetService, logger)

	// =========================================================================
	// Mount API Routes

	r.Route("/v1", func(r chi.Router) {
		r.Use(api.AuthorizationHeader)
		r.Use(api.InternalCommunication(c.Token, logger))
		r.Mount("/snippets", snippetRouter)
	})

	// =========================================================================
	// Start HTTP Server

	errLogger, err := log.NewStdLogger(
		log.With(logger, log.Field("sub-system", "http-error-logger")),
		log.Error,
	)
	if err != nil {
		return fmt.Errorf("failed to create std logger: %w", err)
	}

	server := http.Server{
		Addr:              c.Listen,
		Handler:           r,
		ErrorLog:          errLogger,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	logger.Info("👨‍💻 ➡ web server started", log.Field("listen-addr", c.Listen))

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		//nolint:contextcheck // If we inherit context, new ctx will be finished from the very beginning.
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
