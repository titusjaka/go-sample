package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type authenticatorKey int

const (
	// AuthorizationHeaderKey is the value used to place
	// the authorization header content in context.Context.
	AuthorizationHeaderKey authenticatorKey = iota
)

// AuthorizationHeader tries to retrieve the token string from the "Authorization"
// request header and adds it to the context.
func AuthorizationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(
				r.Context(),
				AuthorizationHeaderKey,
				tokenFromHeader(r),
			)),
		)
	})
}

// InternalCommunication performs bearer authentication with provided token
// for internal service communication purposes.
func InternalCommunication(token string, logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logError := func() {
				logger.Info("Unauthorized request",
					slog.String("method", r.Method),
					slog.String("host", r.Host),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("uri", r.URL.RequestURI()),
					slog.Int("status", http.StatusUnauthorized),
				)
			}

			if v, ok := r.Context().Value(AuthorizationHeaderKey).(string); !ok || v != token {
				logError()

				_ = render.Render(w, r, ErrUnauthorized())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")

	if len(bearer) > 7 && strings.EqualFold(bearer[0:6], "bearer") {
		return bearer[7:]
	}

	return ""
}
