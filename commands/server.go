package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"golang.org/x/sync/errgroup"

	"github.com/titusjaka/go-sample/commands/flags"
	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/kongflag"
	"github.com/titusjaka/go-sample/internal/infrastructure/postgres"
)

// ServerCmd implements kong.Command for the main server command.
type ServerCmd struct {
	Postgres postgres.Flags `kong:"embed"`
	Logger   flags.Logger   `kong:"embed"`

	Listen string `kong:"optional,default=':4040',group='HTTP Server',env=HTTP_LISTEN,help='HTTP network address'"`
	Token  string `kong:"optional,env=API_TOKEN,group='HTTP Server',help='authentication token used for inter-service communication'"`
}

// Run (ServerCmd) runs the main server command.
func (c ServerCmd) Run(kVars kong.Vars) error {
	// =========================================================================
	// Init Context that listens for stop signals
	notifyCtx, notifyCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer notifyCancel()

	// =========================================================================
	// Init Error Group
	gr, ctx := errgroup.WithContext(notifyCtx)

	// =========================================================================
	// Init Logger
	logger := c.Logger.Init()

	logger.Info(
		"⏳ starting service…",
		slog.Group("service_info",
			slog.String("version", kVars[kongflag.ServiceVersion]),
			slog.String("git_commit_sha", kVars[kongflag.GitCommitSHA]),
			slog.String("git_branch", kVars[kongflag.GitBranch]),
			slog.String("go_version", kVars[kongflag.GoVersion]),
		),
		slog.Any("config", c),
	)

	// =========================================================================
	// Init PostgreSQL Connection
	db, err := c.Postgres.OpenStdSQLDB()
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error(
				"unable to close database connection",
				slog.Any("err", closeErr),
			)
		}
	}()

	// ================================================
	// Apply SQL Migrations
	migrationCmd := UpCmd{
		Postgres: c.Postgres,
		Logger:   c.Logger,
	}

	if err = migrationCmd.Run(); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	// =========================================================================
	// Start Private API Server
	gr.Go(func() error {
		return c.runHTTPServer(
			ctx,
			logger.With(
				slog.String("module", "http-server"),
			),
			db,
		)
	})

	// =========================================================================
	// Wait for stop signal
	if err = gr.Wait(); err != nil {
		return fmt.Errorf("error during wait: %w", err)
	}

	logger.Info("💀 ➡ web server gracefully stopped")
	return nil
}

// runHTTPServer starts the HTTP server.
func (c ServerCmd) runHTTPServer(ctx context.Context, logger *slog.Logger, db *sql.DB) error {
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
		logger.With(slog.String("service", "snippets")),
		func() time.Time { return time.Now().UTC() },
	)
	snippetTransport := snippets.NewTransport(snippetService, logger)

	// =========================================================================
	// Mount API Routes

	r.Route("/v1", func(r chi.Router) {
		r.Use(api.AuthorizationHeader)
		r.Use(api.InternalCommunication(c.Token, logger))
		r.Mount("/snippets", snippetTransport.Routes())
	})

	// =========================================================================
	// Start HTTP Server

	server := http.Server{
		Addr:              c.Listen,
		Handler:           r,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	logger.Info("👨‍💻 ➡ web server started", slog.String(
		"listen-addr",
		c.Listen,
	))

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
