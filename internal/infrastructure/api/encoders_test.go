package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

type testResponse struct {
	Data string `json:"data"`
	err  *service.Error
}

func (t *testResponse) ServiceError() *service.Error {
	return t.err
}

func TestEncoders_EncodeResponse(t *testing.T) {
	t.Run("Pass response without error", func(t *testing.T) {
		response := map[string]interface{}{
			"key1": 1.0,
			"key2": "2",
		}

		recorder := httptest.NewRecorder()

		err := api.EncodeResponse(context.Background(), recorder, response)
		require.NoError(t, err)

		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		var actualResponse map[string]interface{}
		err = json.NewDecoder(result.Body).Decode(&actualResponse)
		require.NoError(t, err)
		assert.Equal(t, response, actualResponse)
	})

	t.Run("Pass ServiceErrorer with empty error", func(t *testing.T) {
		response := testResponse{
			Data: "some string",
			err:  nil,
		}

		recorder := httptest.NewRecorder()

		err := api.EncodeResponse(context.Background(), recorder, response)
		require.NoError(t, err)

		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		var actualResponse testResponse
		err = json.NewDecoder(result.Body).Decode(&actualResponse)
		require.NoError(t, err)
		assert.Equal(t, response, actualResponse)
	})

	t.Run("Pass ServiceErrorer with error", func(t *testing.T) {
		response := &testResponse{
			Data: "some string",
			err: &service.Error{
				Type: service.BadRequest,
				Base: errors.New("bad request"),
			},
		}

		recorder := httptest.NewRecorder()

		err := api.EncodeResponse(context.Background(), recorder, response)
		require.NoError(t, err)

		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)

		var actualResponse api.ErrResponse
		err = json.NewDecoder(result.Body).Decode(&actualResponse)
		require.NoError(t, err)
		assert.Equal(t, "bad request", actualResponse.Error)
	})
}

func TestEncoders_EncodeError(t *testing.T) {
	t.Run("Not a service error", func(t *testing.T) {
		err := errors.New("bad romance")

		recorder := httptest.NewRecorder()

		api.EncodeError(context.Background(), err, recorder)

		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

		var actualResponse api.ErrResponse
		err = json.NewDecoder(result.Body).Decode(&actualResponse)
		require.NoError(t, err)
		assert.Equal(t, "bad romance", actualResponse.Error)
	})

	t.Run("A service error", func(t *testing.T) {
		svcErr := &service.Error{
			Type: service.BadRequest,
			Base: errors.New("bad romance"),
		}

		recorder := httptest.NewRecorder()

		api.EncodeError(context.Background(), svcErr, recorder)

		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)

		var actualResponse api.ErrResponse
		err := json.NewDecoder(result.Body).Decode(&actualResponse)
		require.NoError(t, err)
		assert.Equal(t, "bad romance", actualResponse.Error)
	})
}
