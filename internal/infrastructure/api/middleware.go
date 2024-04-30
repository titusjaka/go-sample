package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// NewNotFoundHandler returns http.HandlerFunc, that handles default 404 behavior
func NewNotFoundHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("resource not found",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("uri", r.URL.RequestURI()),
				slog.Int("status", http.StatusNotFound),
			)
		}

		logError()
		_ = render.Render(w, r, ErrNotFound(errors.New("resource not found")))
	}
}

// NewMethodNotAllowedHandler returns http.HandlerFunc, that handles default 405 behavior
func NewMethodNotAllowedHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("method not allowed",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("uri", r.URL.RequestURI()),
				slog.Int("status", http.StatusMethodNotAllowed),
			)
		}

		logError()
		_ = render.Render(w, r, ErrMethodNotAllowed(errors.New("method not allowed")))
	}
}
