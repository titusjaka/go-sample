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

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq" // import pg driver
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
	"github.com/titusjaka/go-sample/migrate"
)

var errStopped = errors.New("stopped")

// NewServerCmd creates a new server CLI sub-command
func NewServerCmd() *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Runs HTTP server",
		Description: "HTTP server with a specimen route “snippets”.",
		Action:      runServer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dsn",
				Usage:    "Data Source Name for PostgreSQL database server",
				EnvVars:  []string{"DATABASE_DSN", "POSTGRES_DSN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "listen",
				Usage:   "HTTP network address",
				EnvVars: []string{"LISTEN", "HTTP_LISTEN"},
				Value:   ":4000",
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "authentication token used for inter-service communication",
				EnvVars: []string{"API_TOKEN", "INTER_SERVICE_TOKEN"},
				Value:   "",
			},
		},
	}
}

func runServer(c *cli.Context) (err error) {
	gr, ctx := errgroup.WithContext(context.Background())

	logger := log.New()

	db, err := sql.Open("postgres", c.String("dsn"))
	if err != nil {
		return fmt.Errorf("unable to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("can't ping database connection: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = fmt.Errorf("unable to close database connection: %w", closeErr)
		}
	}()

	applied, err := applyMigrations(ctx, db)
	if err != nil {
		return fmt.Errorf("unable to apply database migrations: %w", err)
	}

	logger.Info("👟 ➡ migration(-s) applied successfully", log.Field("applied", applied))

	gr.Go(func() error {
		return setupServer(
			ctx,
			c.String("listen"),
			c.String("token"),
			log.With(logger, log.Field("system", "http-server")),
			db,
		)
	})

	gr.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stop)

		select {
		case <-stop:
			return errStopped
		case <-ctx.Done():
			return nil
		}
	})

	if err = gr.Wait(); err != nil && err != errStopped {
		return fmt.Errorf("error during wait: %w", err)
	}

	logger.Info("💀 ➡ web server gracefully stopped")
	return nil
}

func setupServer(
	ctx context.Context,
	listen string,
	internalToken string,
	logger log.Logger,
	db *sql.DB,
) error {
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

	snippetStorage := snippets.NewPGStorage(db)
	snippetService := snippets.NewService(
		snippetStorage,
		log.With(logger, log.Field("service", "snippets")),
	)
	snippetRouter := snippets.MakeSnippetsHandler(snippetService, logger)
	r.Route("/v1", func(r chi.Router) {
		r.Use(api.AuthorizationHeader)
		r.Use(api.InternalCommunication(internalToken, logger))
		r.Mount("/snippets", snippetRouter)
	})

	errLogger, err := log.NewStdLogger(
		log.With(logger, log.Field("sub-system", "http-error-logger")),
		log.Error,
	)
	if err != nil {
		return fmt.Errorf("failed to create std logger: %w", err)
	}

	server := http.Server{
		Addr:     listen,
		Handler:  r,
		ErrorLog: errLogger,
	}
	errCh := make(chan error, 1)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	logger.Info("👨‍💻 ➡ web server started", log.Field("listen-addr", listen))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
	defer shutdownCancel()

	select {
	case <-ctx.Done():
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func applyMigrations(ctx context.Context, db *sql.DB) (int, error) {
	source := migrate.NewEmbeddedSource()
	migrator := migrate.NewMigrator(db, source)
	return migrator.Up(ctx)
}
