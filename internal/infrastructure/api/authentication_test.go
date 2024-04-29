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
	"github.com/titusjaka/go-sample/internal/infrastructure/nopslog"
)

func TestApi_InternalCommunication(t *testing.T) {
	t.Run("Successfully authorized", func(t *testing.T) {
		expectedToken := "12345"

		router := chi.NewRouter()
		router.Use(render.SetContentType(render.ContentTypeJSON))
		router.Use(api.AuthorizationHeader)
		router.Use(api.InternalCommunication(expectedToken, nopslog.NewNoplogger()))

		router.Get("/", func(w http.ResponseWriter, req *http.Request) {
			actualToken, ok := req.Context().Value(api.AuthorizationHeaderKey).(string)
			assert.True(t, ok)
			assert.Equal(t, expectedToken, actualToken)
			_, _ = w.Write([]byte("OK"))
		})

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		request.Header.Add("Authorization", "Bearer "+expectedToken)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer func() {
			require.NoError(t, result.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Run("Empty auth header", func(t *testing.T) {
			internalToken := "12345"

			router := chi.NewRouter()
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Use(api.AuthorizationHeader)
			router.Use(api.InternalCommunication(internalToken, nopslog.NewNoplogger()))

			router.Get("/", func(w http.ResponseWriter, req *http.Request) {
				_, _ = w.Write([]byte("OK"))
			})

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer func() {
				require.NoError(t, result.Body.Close())
			}()

			assert.Equal(t, http.StatusUnauthorized, result.StatusCode)

			var actualResponse api.ErrResponse
			err = json.NewDecoder(result.Body).Decode(&actualResponse)
			require.NoError(t, err)
			assert.NotEmpty(t, actualResponse.Error)
		})

		t.Run("Wrong auth header", func(t *testing.T) {
			internalToken := "12345"
			expectedToken := "wrong token"

			router := chi.NewRouter()
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Use(api.AuthorizationHeader)
			router.Use(api.InternalCommunication(internalToken, nopslog.NewNoplogger()))

			router.Get("/", func(w http.ResponseWriter, req *http.Request) {
				actualToken, ok := req.Context().Value(api.AuthorizationHeaderKey).(string)
				assert.True(t, ok)
				assert.Equal(t, expectedToken, actualToken)
				_, _ = w.Write([]byte("OK"))
			})

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			request.Header.Add("Authorization", "Bearer "+expectedToken)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer func() {
				require.NoError(t, result.Body.Close())
			}()

			assert.Equal(t, http.StatusUnauthorized, result.StatusCode)

			var actualResponse api.ErrResponse
			err = json.NewDecoder(result.Body).Decode(&actualResponse)
			require.NoError(t, err)
			assert.NotEmpty(t, actualResponse.Error)
		})
	})
}
