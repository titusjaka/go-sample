package api

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// ServiceErrorer is used to pass business-logic errors through service -> endpoint -> transport layers.
// If the endpoint can return a business-logic error, then you need to implement the ServiceError() method.
type ServiceErrorer interface {
	ServiceError() *service.Error
}

// EncodeResponse is a default response encoder. It uses go-kit JSON Encoder.
// If response implements ServiceErrorer interface and response.ServiceError() is not nil,
// then it's handled as ErrResponse.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if svcErrorer, ok := response.(ServiceErrorer); ok && svcErrorer.ServiceError() != nil {
		// Not a Go kit transport error, but a business-logic error.
		return kithttp.EncodeJSONResponse(ctx, w, NewErrResponse(svcErrorer.ServiceError()))
	}
	return kithttp.EncodeJSONResponse(ctx, w, response)
}

// EncodeError encodes errors occurred on transport layer
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	var errResp *ErrResponse
	switch err := err.(type) {
	case *service.Error:
		errResp = NewErrResponse(err)
	default:
		errResp = ErrInternal(err)
	}
	_ = kithttp.EncodeJSONResponse(ctx, w, errResp)
}
