package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/go-sample/v2/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/service"
)

func TestNewErrResponse(t *testing.T) {
	tests := []struct {
		name               string
		svcError           *service.Error
		expectedStatusCode int
	}{
		{
			name:               "empty error",
			svcError:           nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "default error",
			svcError:           &service.Error{},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "empty error type",
			svcError:           &service.Error{Base: errors.New("some text")},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "BadRequest",
			svcError: &service.Error{
				Type: service.BadRequest,
				Base: errors.New("badRequest"),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Forbidden",
			svcError: &service.Error{
				Type: service.Forbidden,
				Base: errors.New("forbidden"),
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name: "InternalError",
			svcError: &service.Error{
				Type: service.InternalError,
				Base: errors.New("internalError"),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Unauthorized",
			svcError: &service.Error{
				Type: service.Unauthorized,
				Base: errors.New("unauthorized"),
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "NotFound",
			svcError: &service.Error{
				Type: service.NotFound,
				Base: errors.New("notFound"),
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := api.NewErrResponse(tt.svcError)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode())
		})
	}
}
