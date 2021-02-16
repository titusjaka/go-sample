package api

import (
	"errors"
	"net/http"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// ErrResponse renderer type for handling all sorts of business errors.
type ErrResponse struct {
	Error      string `json:"error,omitempty"`
	statusCode int
}

// StatusCode returns an HTTP status code.
// It implements kithttp.StatusCoder interface
func (e *ErrResponse) StatusCode() int {
	return e.statusCode
}

// ErrBadRequest handler returns the pre-defined 400 schema.
func ErrBadRequest(err error) *ErrResponse {
	return &ErrResponse{
		Error:      err.Error(),
		statusCode: http.StatusBadRequest,
	}
}

// ErrMethodNotAllowed handler returns the pre-defined 405 schema.
func ErrMethodNotAllowed(err error) *ErrResponse {
	return &ErrResponse{
		Error:      err.Error(),
		statusCode: http.StatusMethodNotAllowed,
	}
}

// ErrInternal handler returns the pre-defined 500 schema.
func ErrInternal(err error) *ErrResponse {
	return &ErrResponse{
		Error:      err.Error(),
		statusCode: http.StatusInternalServerError,
	}
}

// ErrForbidden handler returns the pre-defined 403 schema.
func ErrForbidden() *ErrResponse {
	return &ErrResponse{
		Error:      http.StatusText(http.StatusForbidden),
		statusCode: http.StatusForbidden,
	}
}

// ErrNotFound handler returns the pre-defined 404 schema.
func ErrNotFound(err error) *ErrResponse {
	return &ErrResponse{
		Error:      err.Error(),
		statusCode: http.StatusNotFound,
	}
}

// ErrUnauthorized handler returns the pre-defined 401 schema.
func ErrUnauthorized() *ErrResponse {
	return &ErrResponse{
		Error:      http.StatusText(http.StatusUnauthorized),
		statusCode: http.StatusUnauthorized,
	}
}

// NewErrResponse wraps Error into HTTP error response
func NewErrResponse(err *service.Error) *ErrResponse {
	internalErr := errors.New("internal error")
	if err == nil {
		return ErrInternal(internalErr)
	}

	switch err.Type {
	case service.BadRequest:
		return ErrBadRequest(err.Base)
	case service.Forbidden:
		return ErrForbidden()
	case service.InternalError:
		return ErrInternal(err.Base)
	case service.Unauthorized:
		return ErrUnauthorized()
	case service.NotFound:
		return ErrNotFound(err.Base)
	default:
		if err.Base != nil {
			return ErrInternal(err.Base)
		}
		return ErrInternal(internalErr)
	}
}
