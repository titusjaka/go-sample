package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/titusjaka/go-sample/internal/infrastructure/log"
)

// LogErrorHandler is a transport error handler implementation which logs an error.
type LogErrorHandler struct {
	logger log.Logger
}

// NewLogErrorHandler returns a new log middleware for go-kit transport
func NewLogErrorHandler(logger log.Logger) transport.ErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

// Handle implements go-kit ErrorHandler interface
func (h *LogErrorHandler) Handle(_ context.Context, err error) {
	h.logger.Error("error occurred", log.Field("err", err))
}

// NewNotFoundHandler returns http.HandlerFunc, that handles default 404 behavior
func NewNotFoundHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("Resource not found",
				log.Field("method", r.Method),
				log.Field("host", r.Host),
				log.Field("uri", r.URL.RequestURI()),
				log.Field("status", http.StatusNotFound),
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
func NewMethodNotAllowedHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logError := func() {
			logger.Info("Method not allowed",
				log.Field("method", r.Method),
				log.Field("host", r.Host),
				log.Field("uri", r.URL.RequestURI()),
				log.Field("status", http.StatusMethodNotAllowed),
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
