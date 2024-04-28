package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
)

func TestApi_NotFoundHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.NotFound(api.NewNotFoundHandler(log.NopLogger{}))
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	request, err := http.NewRequest(http.MethodGet, "/not_existing_URL", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()
	defer func() {
		require.NoError(t, result.Body.Close())
	}()

	assert.Equal(t, http.StatusNotFound, result.StatusCode)

	var actualResponse api.ErrResponse
	err = json.NewDecoder(result.Body).Decode(&actualResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, actualResponse.Error)
}

func TestApi_MethodNotAllowedHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.MethodNotAllowed(api.NewMethodNotAllowedHandler(log.NopLogger{}))
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	request, err := http.NewRequest(http.MethodPost, "/", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()
	defer func() {
		err = result.Body.Close()
		require.NoError(t, err)
	}()

	assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)

	var actualResponse api.ErrResponse
	err = json.NewDecoder(result.Body).Decode(&actualResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, actualResponse.Error)
}
