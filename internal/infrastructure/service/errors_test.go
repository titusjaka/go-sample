package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/go-sample/v2/internal/infrastructure/service"
)

func TestError_Error(t *testing.T) {
	t.Run("Empty error", func(t *testing.T) {
		var svcError *service.Error
		assert.Empty(t, svcError.Error())
	})

	t.Run("Embed error", func(t *testing.T) {
		svcError := &service.Error{
			Type: service.InternalError,
			Base: errors.New("internal error"),
		}
		assert.Equal(t, "internal error", svcError.Error())
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Run("Empty error", func(t *testing.T) {
		var svcError *service.Error
		assert.Nil(t, svcError.Unwrap())
	})

	t.Run("Embed error", func(t *testing.T) {
		internalErr := errors.New("internal error")
		svcError := &service.Error{
			Type: service.InternalError,
			Base: internalErr,
		}
		assert.True(t, errors.Is(svcError, internalErr))
	})
}
