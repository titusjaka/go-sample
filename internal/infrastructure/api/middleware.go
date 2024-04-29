package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

// LogErrorHandler is a transport error handler implementation which logs an error.
type LogErrorHandler struct {
	logger *slog.Logger
}

// NewLogErrorHandler returns a new log middleware for go-kit transport
func NewLogErrorHandler(logger *slog.Logger) transport.ErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

// Handle implements go-kit ErrorHandler interface
func (h *LogErrorHandler) Handle(_ context.Context, err error) {
	h.logger.Error("error occurred", slog.Any("err", err))
}

// NewNotFoundHandler returns http.HandlerFunc, that handles default 404 behavior
func NewNotFoundHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("Resource not found",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("uri", r.URL.RequestURI()),
				slog.Int("status", http.StatusNotFound),
			)
		}

		respond := func() {
			errNotFound := errors.New("not found")
			_ = kithttp.EncodeJSONResponse(context.Background(), w, ErrNotFound(errNotFound))
		}

		logError()
		respond()
	}
}

// NewMethodNotAllowedHandler returns http.HandlerFunc, that handles default 405 behavior
func NewMethodNotAllowedHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("Method not allowed",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("uri", r.URL.RequestURI()),
				slog.Int("status", http.StatusMethodNotAllowed),
			)
		}

		respond := func() {
			errMethodNotAllowed := errors.New("method not allowed")
			_ = kithttp.EncodeJSONResponse(context.Background(), w, ErrMethodNotAllowed(errMethodNotAllowed))
		}

		logError()
		respond()
	}
}
