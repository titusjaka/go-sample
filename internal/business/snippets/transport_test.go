package snippets_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

type snippetResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func TestTransport_MakeSnippetsHandler(t *testing.T) {
	t.Run("List snippets", func(t *testing.T) {
		t.Run("Successfully list snippets", func(t *testing.T) {
			t.Run("Empty list", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedPagination := service.Pagination{
					Limit:       100,
					Offset:      0,
					Total:       0,
					TotalPages:  1,
					CurrentPage: 1,
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().List(gomock.Any(), uint(0), uint(0)).Return(nil, expectedPagination, nil)

				client := server.Client()
				resp, err := client.Get(server.URL + "/snippets")
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetsResponse struct {
					Snippets   []snippetResponse  `json:"snippets"`
					Pagination service.Pagination `json:"pagination"`
				}

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetsResponse)
				require.NoError(t, err)
				assert.Empty(t, actualListSnippetsResponse.Snippets)
				assert.Equal(t, expectedPagination, actualListSnippetsResponse.Pagination)
			})
			t.Run("Without pagination", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
				fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
				expectedSnippets := []snippets.Snippet{
					{
						ID:        100,
						Title:     "Snippet #100",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated,
						UpdatedAt: fakeTimeCreated,
						ExpiresAt: fakeTimeExpires,
					},
					{
						ID:        1000,
						Title:     "Snippet #1000",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated.Add(time.Hour),
						UpdatedAt: fakeTimeCreated.Add(time.Hour),
						ExpiresAt: fakeTimeExpires.Add(time.Hour),
					},
				}
				limit := uint(0)
				offset := uint(0)
				expectedPagination := service.Pagination{
					Limit:       100,
					Offset:      0,
					Total:       2,
					TotalPages:  1,
					CurrentPage: 1,
				}
				expectedSnippetsResponse := []snippetResponse{
					{
						ID:        100,
						Title:     "Snippet #100",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated,
						ExpiresAt: fakeTimeExpires,
					},
					{
						ID:        1000,
						Title:     "Snippet #1000",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated.Add(time.Hour),
						ExpiresAt: fakeTimeExpires.Add(time.Hour),
					},
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedSnippets, expectedPagination, nil)

				client := server.Client()
				resp, err := client.Get(server.URL + "/snippets")
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetsResponse struct {
					Snippets   []snippetResponse  `json:"snippets"`
					Pagination service.Pagination `json:"pagination"`
				}

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetsResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, actualListSnippetsResponse.Snippets)
				assert.Equal(t, len(actualListSnippetsResponse.Snippets), len(expectedSnippetsResponse))
				assert.Equal(t, expectedPagination, actualListSnippetsResponse.Pagination)
				assert.EqualValues(t, expectedSnippetsResponse, actualListSnippetsResponse.Snippets)
			})
			t.Run("With pagination", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
				fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
				expectedSnippets := []snippets.Snippet{
					{
						ID:        100,
						Title:     "Snippet #100",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated,
						UpdatedAt: fakeTimeCreated,
						ExpiresAt: fakeTimeExpires,
					},
					{
						ID:        1000,
						Title:     "Snippet #1000",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated.Add(time.Hour),
						UpdatedAt: fakeTimeCreated.Add(time.Hour),
						ExpiresAt: fakeTimeExpires.Add(time.Hour),
					},
				}
				limit := uint(100)
				offset := uint(0)
				expectedPagination := service.Pagination{
					Limit:       100,
					Offset:      0,
					Total:       2,
					TotalPages:  1,
					CurrentPage: 1,
				}
				expectedSnippetsResponse := []snippetResponse{
					{
						ID:        100,
						Title:     "Snippet #100",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated,
						ExpiresAt: fakeTimeExpires,
					},
					{
						ID:        1000,
						Title:     "Snippet #1000",
						Content:   "Very important text",
						CreatedAt: fakeTimeCreated.Add(time.Hour),
						ExpiresAt: fakeTimeExpires.Add(time.Hour),
					},
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedSnippets, expectedPagination, nil)

				client := server.Client()
				resp, err := client.Get(server.URL + fmt.Sprintf("/snippets?limit=%d&offset=%d", limit, offset))
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetsResponse struct {
					Snippets   []snippetResponse  `json:"snippets"`
					Pagination service.Pagination `json:"pagination"`
				}

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetsResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, actualListSnippetsResponse.Snippets)
				assert.Equal(t, len(actualListSnippetsResponse.Snippets), len(expectedSnippetsResponse))
				assert.Equal(t, expectedPagination, actualListSnippetsResponse.Pagination)
				assert.EqualValues(t, expectedSnippetsResponse, actualListSnippetsResponse.Snippets)
			})
		})
		t.Run("Failed to list snippets", func(t *testing.T) {
			t.Run("Error error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedPagination := service.Pagination{
					Limit:       100,
					Offset:      0,
					Total:       0,
					TotalPages:  1,
					CurrentPage: 1,
				}
				expectedErrorMsg := "internal error"
				expectedSvcErr := &service.Error{
					Type: service.InternalError,
					Base: errors.New(expectedErrorMsg),
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().List(gomock.Any(), uint(0), uint(0)).Return(nil, expectedPagination, expectedSvcErr)

				client := server.Client()
				resp, err := client.Get(server.URL + "/snippets")
				require.NoError(t, err)
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListResponse.Error, expectedErrorMsg)
			})
			t.Run("Wrong limit", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				client := server.Client()
				resp, err := client.Get(server.URL + "/snippets?limit=-1")
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, actualListResponse.Error)
			})
		})
	})

	t.Run("Get snippet by ID", func(t *testing.T) {
		t.Run("Successfully get a snippet", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockService(ctrl)

			fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			expectedID := uint(100)
			expectedSnippet := snippets.Snippet{
				ID:        expectedID,
				Title:     "Snippet #100",
				Content:   "Very important text",
				CreatedAt: fakeTimeCreated,
				UpdatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}
			expectedSnippetResponse := snippetResponse{
				ID:        expectedID,
				Title:     "Snippet #100",
				Content:   "Very important text",
				CreatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}

			handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
			router := chi.NewRouter()
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Mount("/snippets", handler)

			server := httptest.NewServer(router)
			defer server.Close()

			mockService.EXPECT().Get(gomock.Any(), expectedID).Return(expectedSnippet, nil)

			client := server.Client()
			resp, err := client.Get(server.URL + fmt.Sprintf("/snippets/%d", expectedID))
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			defer func() {
				require.NoError(t, resp.Body.Close())
			}()

			var actualListSnippetResponse snippetResponse

			err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
			require.NoError(t, err)
			assert.EqualValues(t, expectedSnippetResponse, actualListSnippetResponse)
		})
		t.Run("Failed to get a snippet", func(t *testing.T) {
			t.Run("Error error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedID := uint(100)
				expectedErrorMsg := "internal error"
				expectedSvcErr := &service.Error{
					Type: service.InternalError,
					Base: errors.New(expectedErrorMsg),
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().Get(gomock.Any(), expectedID).Return(snippets.Snippet{}, expectedSvcErr)

				client := server.Client()
				resp, err := client.Get(server.URL + fmt.Sprintf("/snippets/%d", expectedID))
				require.NoError(t, err)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListSnippetResponse.Error, expectedErrorMsg)
			})

			t.Run("Not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedID := uint(100)
				expectedErrorMsg := snippets.ErrNotFound.Error()
				expectedSvcErr := &service.Error{
					Type: service.NotFound,
					Base: snippets.ErrNotFound,
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().Get(gomock.Any(), expectedID).Return(snippets.Snippet{}, expectedSvcErr)

				client := server.Client()
				resp, err := client.Get(server.URL + fmt.Sprintf("/snippets/%d", expectedID))
				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListSnippetResponse.Error, expectedErrorMsg)
			})

			t.Run("Wrong id param", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				client := server.Client()
				resp, err := client.Get(server.URL + "/snippets/-1")
				require.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, actualListSnippetResponse.Error)
			})
		})
	})

	t.Run("Create a new snippet", func(t *testing.T) {
		t.Run("Successfully create a snippet", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockService(ctrl)

			fakeTimeCreated := time.Now().UTC()
			fakeTimeExpires := fakeTimeCreated.Add(time.Hour * 24 * 120) // 120 days after
			fakeTimeString := fakeTimeExpires.Format(time.RFC3339)
			expectedID := uint(100)
			expectedSnippet := snippets.Snippet{
				ID:        expectedID,
				Title:     "Snippet #100",
				Content:   "Very important text",
				CreatedAt: fakeTimeCreated,
				UpdatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}
			expectedSnippetResponse := snippetResponse{
				ID:        expectedID,
				Title:     "Snippet #100",
				Content:   "Very important text",
				CreatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}

			handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
			router := chi.NewRouter()
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Mount("/snippets", handler)

			server := httptest.NewServer(router)
			defer server.Close()

			mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedSnippet, nil)

			client := server.Client()

			// language=JSON
			snippetJSON := fmt.Sprintf(`{"title": "Snippet #100", "content": "Very important text", "expires_at": "%s"}`,
				fakeTimeString,
			)

			resp, err := client.Post(
				server.URL+"/snippets",
				"application/json",
				bytes.NewReader([]byte(snippetJSON)),
			)
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			defer func() {
				require.NoError(t, resp.Body.Close())
			}()

			var actualListSnippetResponse snippetResponse

			err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
			require.NoError(t, err)
			assert.EqualValues(t, expectedSnippetResponse, actualListSnippetResponse)
		})

		t.Run("Failed to create a snippet", func(t *testing.T) {
			t.Run("Error error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				fakeTimeCreated := time.Now().UTC()
				fakeTimeExpires := fakeTimeCreated.Add(time.Hour * 24 * 120) // 120 days after
				fakeTimeString := fakeTimeExpires.Format(time.RFC3339)

				expectedErrorMsg := "internal error"
				expectedSvcError := &service.Error{
					Type: service.InternalError,
					Base: errors.New(expectedErrorMsg),
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(snippets.Snippet{}, expectedSvcError)

				client := server.Client()

				// language=JSON
				snippetJSON := fmt.Sprintf(`{"title": "Snippet #100", "content": "Very important text", "expires_at": "%s"}`,
					fakeTimeString,
				)

				resp, err := client.Post(
					server.URL+"/snippets",
					"application/json",
					bytes.NewReader([]byte(snippetJSON)),
				)
				require.NoError(t, err)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListResponse.Error, expectedErrorMsg)
			})
			t.Run("Bad request", func(t *testing.T) {
				t.Run("Wrong JSON", func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					mockService := NewMockService(ctrl)

					fakeTimeCreated := time.Now().UTC()
					fakeTimeExpires := fakeTimeCreated.Add(time.Hour * 24 * 120) // 120 days after
					fakeTimeString := fakeTimeExpires.Format(time.RFC3339)

					handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
					router := chi.NewRouter()
					router.Use(render.SetContentType(render.ContentTypeJSON))
					router.Mount("/snippets", handler)

					server := httptest.NewServer(router)
					defer server.Close()

					client := server.Client()

					// language=JSON
					snippetJSON := fmt.Sprintf(`{"title": 100, "content": "Very important text", "expires_at": "%s"}`,
						fakeTimeString,
					)

					resp, err := client.Post(
						server.URL+"/snippets",
						"application/json",
						bytes.NewReader([]byte(snippetJSON)),
					)
					require.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

					defer func() {
						require.NoError(t, resp.Body.Close())
					}()

					var actualListSnippetResponse api.ErrResponse

					err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
					require.NoError(t, err)
					assert.NotEmpty(t, actualListSnippetResponse.Error)
				})
				t.Run("Wrong time format", func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					mockService := NewMockService(ctrl)

					fakeTimeCreated := time.Now().UTC()
					fakeTimeExpires := fakeTimeCreated.Add(time.Hour * 24 * 120) // 120 days after
					fakeTimeString := fakeTimeExpires.Format(time.RFC1123Z)

					handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
					router := chi.NewRouter()
					router.Use(render.SetContentType(render.ContentTypeJSON))
					router.Mount("/snippets", handler)

					server := httptest.NewServer(router)
					defer server.Close()

					client := server.Client()

					// language=JSON
					snippetJSON := fmt.Sprintf(`{"title": "Snippet #100", "content": "Very important text", "expires_at": "%s"}`,
						fakeTimeString,
					)

					resp, err := client.Post(
						server.URL+"/snippets",
						"application/json",
						bytes.NewReader([]byte(snippetJSON)),
					)
					require.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

					defer func() {
						require.NoError(t, resp.Body.Close())
					}()

					var actualListSnippetResponse api.ErrResponse

					err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
					require.NoError(t, err)
					assert.NotEmpty(t, actualListSnippetResponse.Error)
				})
			})
		})
	})

	t.Run("Delete an existing snippet", func(t *testing.T) {
		t.Run("Successfully delete a snippet", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockService(ctrl)

			expectedID := uint(100)

			handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
			router := chi.NewRouter()
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Mount("/snippets", handler)

			server := httptest.NewServer(router)
			defer server.Close()

			mockService.EXPECT().SoftDelete(gomock.Any(), expectedID).Return(nil)

			client := server.Client()
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/snippets/%d", server.URL, expectedID),
				nil,
			)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)

			defer func() {
				require.NoError(t, resp.Body.Close())
			}()
		})
		t.Run("Failed to delete a snippet", func(t *testing.T) {
			t.Run("Error error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedID := uint(100)
				expectedErrorMsg := "internal error"
				expectedSvcErr := &service.Error{
					Type: service.InternalError,
					Base: errors.New(expectedErrorMsg),
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().SoftDelete(gomock.Any(), expectedID).Return(expectedSvcErr)

				client := server.Client()
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("%s/snippets/%d", server.URL, expectedID),
					nil,
				)
				require.NoError(t, err)

				resp, err := client.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListSnippetResponse.Error, expectedErrorMsg)
			})

			t.Run("Bad request", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				client := server.Client()
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("%s/snippets/wrong", server.URL),
					nil,
				)
				require.NoError(t, err)

				resp, err := client.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, actualListSnippetResponse.Error)
			})

			t.Run("Not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockService := NewMockService(ctrl)

				expectedID := uint(100)
				expectedErrorMsg := snippets.ErrNotFound.Error()
				expectedSvcErr := &service.Error{
					Type: service.NotFound,
					Base: snippets.ErrNotFound,
				}

				handler := snippets.MakeSnippetsHandler(mockService, log.NopLogger{})
				router := chi.NewRouter()
				router.Use(render.SetContentType(render.ContentTypeJSON))
				router.Mount("/snippets", handler)

				server := httptest.NewServer(router)
				defer server.Close()

				mockService.EXPECT().SoftDelete(gomock.Any(), expectedID).Return(expectedSvcErr)

				client := server.Client()
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("%s/snippets/%d", server.URL, expectedID),
					nil,
				)
				require.NoError(t, err)

				resp, err := client.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)

				defer func() {
					require.NoError(t, resp.Body.Close())
				}()

				var actualListSnippetResponse api.ErrResponse

				err = json.NewDecoder(resp.Body).Decode(&actualListSnippetResponse)
				require.NoError(t, err)
				assert.Equal(t, actualListSnippetResponse.Error, expectedErrorMsg)
			})
		})
	})
}
