package snippets_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"go.uber.org/mock/gomock"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/nopslog"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

func TestTransport_listSnippets(t *testing.T) {
	t.Parallel()

	t.Run("Successfully list snippets", func(t *testing.T) {
		t.Parallel()

		t.Run("Empty list", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			pagination := service.Pagination{
				Limit:       100,
				Offset:      0,
				Total:       0,
				TotalPages:  1,
				CurrentPage: 1,
			}

			// ================================================
			// Describe mock calls
			mockService.EXPECT().List(
				gomock.Any(),
				uint(0),
				uint(0),
			).Return(nil, pagination, nil)

			// ================================================
			// Run test
			expected := map[string]any{
				"pagination": pagination,
			}

			response := expect.GET("/").
				Expect()

			response.
				Status(http.StatusOK).
				JSON().Object().IsEqual(expected)
		})

		t.Run("Return list of snippets with pagination", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			limit := uint(0)
			offset := uint(0)

			pagination := service.Pagination{
				Limit:       100,
				Offset:      0,
				Total:       2,
				TotalPages:  1,
				CurrentPage: 1,
			}

			fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)

			listOfSnippets := []snippets.Snippet{
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

			// ================================================
			// Describe mock calls
			mockService.EXPECT().List(
				gomock.Any(),
				limit,
				offset,
			).Return(listOfSnippets, pagination, nil)

			// ================================================
			// Run test
			expectedSnippetsResponse := []snippets.SnippetResponse{
				{
					ID:        listOfSnippets[0].ID,
					Title:     listOfSnippets[0].Title,
					Content:   listOfSnippets[0].Content,
					CreatedAt: listOfSnippets[0].CreatedAt,
					ExpiresAt: listOfSnippets[0].ExpiresAt,
				},
				{
					ID:        listOfSnippets[1].ID,
					Title:     listOfSnippets[1].Title,
					Content:   listOfSnippets[1].Content,
					CreatedAt: listOfSnippets[1].CreatedAt,
					ExpiresAt: listOfSnippets[1].ExpiresAt,
				},
			}

			expected := map[string]any{
				"snippets":   expectedSnippetsResponse,
				"pagination": pagination,
			}

			response := expect.GET("/").
				Expect()

			response.
				Status(http.StatusOK).
				JSON().Object().IsEqual(expected)
		})
	})

	t.Run("Failed to list snippets", func(t *testing.T) {
		t.Parallel()

		t.Run("Service error", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			svcErr := &service.Error{
				Type: service.InternalError,
				Base: errors.New("internal error"),
			}

			// ================================================
			// Describe mock calls
			mockService.EXPECT().List(
				gomock.Any(),
				uint(0),
				uint(0),
			).Return(nil, service.Pagination{}, svcErr)

			// ================================================
			// Run test
			expected := map[string]any{
				"error": "internal error",
			}

			response := expect.GET("/").
				Expect()

			response.
				Status(http.StatusInternalServerError).
				JSON().Object().IsEqual(expected)
		})

		t.Run("Bad request", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Run test
			expected := map[string]any{
				"error": `schema: error converting value for "limit"`,
			}

			response := expect.GET("/").
				WithQuery("limit", "-1").
				Expect()

			response.
				Status(http.StatusBadRequest).
				JSON().Object().IsEqual(expected)
		})
	})
}

func TestTransport_getSnippet(t *testing.T) {
	t.Parallel()

	t.Run("Successfully get snippet", func(t *testing.T) {
		t.Parallel()

		// ================================================
		// Init mocks and service
		ctrl := gomock.NewController(t)

		mockService := NewMockService(ctrl)
		transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
		handler := transport.Routes()

		// ================================================
		// Create httpexpect instance
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})

		// ================================================
		// Init test data
		id := uint(100)

		fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
		fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)

		snippet := snippets.Snippet{
			ID:        id,
			Title:     "Snippet #100",
			Content:   "Very important text",
			CreatedAt: fakeTimeCreated,
			UpdatedAt: fakeTimeCreated,
			ExpiresAt: fakeTimeExpires,
		}

		// ================================================
		// Describe mock calls
		mockService.EXPECT().Get(
			gomock.Any(),
			id,
		).Return(snippet, nil)

		// ================================================
		// Run test
		expectedSnippetResponse := snippets.SnippetResponse{
			ID:        id,
			Title:     "Snippet #100",
			Content:   "Very important text",
			CreatedAt: fakeTimeCreated,
			ExpiresAt: fakeTimeExpires,
		}

		response := expect.GET("/{id}", id).
			Expect()

		response.
			Status(http.StatusOK).
			JSON().Object().IsEqual(expectedSnippetResponse)
	})

	t.Run("Failed to get snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("Service error", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			svcErr := &service.Error{
				Type: service.InternalError,
				Base: errors.New("internal error"),
			}

			// ================================================
			// Describe mock calls
			mockService.EXPECT().Get(
				gomock.Any(),
				uint(1),
			).Return(snippets.Snippet{}, svcErr)

			// ================================================
			// Run test
			expected := map[string]any{
				"error": "internal error",
			}

			response := expect.GET("/{id}", 1).
				Expect()

			response.
				Status(http.StatusInternalServerError).
				JSON().Object().IsEqual(expected)
		})

		t.Run("Bad request", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Run test
			expected := map[string]any{
				"error": `invalid id param: -1`,
			}

			response := expect.GET("/{id}", -1).
				Expect()

			response.
				Status(http.StatusBadRequest).
				JSON().Object().IsEqual(expected)
		})
	})
}

func TestTransport_createSnippet(t *testing.T) {
	t.Parallel()

	t.Run("Successfully create snippet", func(t *testing.T) {
		t.Parallel()

		// ================================================
		// Init mocks and service
		ctrl := gomock.NewController(t)

		mockService := NewMockService(ctrl)
		transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
		handler := transport.Routes()

		// ================================================
		// Create httpexpect instance
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})

		// ================================================
		// Init test data
		createSnippetRequest := snippets.CreateSnippetRequest{
			Title:     "Snippet #100",
			Content:   "Very important text",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 120).Truncate(time.Second),
		}

		createdSnippet := snippets.Snippet{
			ID:        100,
			Title:     createSnippetRequest.Title,
			Content:   createSnippetRequest.Content,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			ExpiresAt: createSnippetRequest.ExpiresAt,
		}

		// ================================================
		// Describe mock calls
		mockService.EXPECT().Create(
			gomock.Any(),
			snippets.Snippet{
				Title:     createSnippetRequest.Title,
				Content:   createSnippetRequest.Content,
				ExpiresAt: createSnippetRequest.ExpiresAt,
			},
		).Return(createdSnippet, nil)

		// ================================================
		// Run test
		expectedSnippetResponse := snippets.SnippetResponse{
			ID:        createdSnippet.ID,
			Title:     createdSnippet.Title,
			Content:   createdSnippet.Content,
			CreatedAt: createdSnippet.CreatedAt,
			ExpiresAt: createdSnippet.ExpiresAt,
		}

		response := expect.POST("/").
			WithJSON(createSnippetRequest).
			Expect()

		response.
			Status(http.StatusOK).
			JSON().Object().IsEqual(expectedSnippetResponse)
	})

	t.Run("Failed to create snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("Service error", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			createSnippetRequest := snippets.CreateSnippetRequest{
				Title:     "Snippet #100",
				Content:   "Very important text",
				ExpiresAt: time.Now().Add(time.Hour * 24 * 120).Truncate(time.Second),
			}

			svcErr := &service.Error{
				Type: service.InternalError,
				Base: errors.New("internal error"),
			}

			// ================================================
			// Describe mock calls
			mockService.EXPECT().Create(
				gomock.Any(),
				snippets.Snippet{
					Title:     createSnippetRequest.Title,
					Content:   createSnippetRequest.Content,
					ExpiresAt: createSnippetRequest.ExpiresAt,
				},
			).Return(snippets.Snippet{}, svcErr)

			// ================================================
			// Run test
			expected := map[string]any{
				"error": "internal error",
			}

			response := expect.POST("/").
				WithJSON(createSnippetRequest).
				Expect()

			response.
				Status(http.StatusInternalServerError).
				JSON().Object().IsEqual(expected)
		})

		t.Run("Bad request", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			createSnippetRequest := snippets.CreateSnippetRequest{
				Title:     "",
				Content:   "Very important text",
				ExpiresAt: time.Now().Add(time.Hour * 24 * 120).Truncate(time.Second),
			}

			// ================================================
			// Run test
			expected := map[string]any{
				"error": `title: cannot be blank.`,
			}

			response := expect.POST("/").
				WithJSON(createSnippetRequest).
				Expect()

			response.
				Status(http.StatusBadRequest).
				JSON().Object().IsEqual(expected)
		})
	})
}

func TestTransport_deleteSnippet(t *testing.T) {
	t.Parallel()

	t.Run("Successfully delete snippet", func(t *testing.T) {
		t.Parallel()

		// ================================================
		// Init mocks and service
		ctrl := gomock.NewController(t)

		mockService := NewMockService(ctrl)
		transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
		handler := transport.Routes()

		// ================================================
		// Create httpexpect instance
		expect := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})

		// ================================================
		// Init test data
		id := uint(100)

		// ================================================
		// Describe mock calls
		mockService.EXPECT().SoftDelete(
			gomock.Any(),
			id,
		).Return(nil)

		// ================================================
		// Run test
		response := expect.DELETE("/{id}", id).
			Expect()

		response.
			Status(http.StatusNoContent)
	})

	t.Run("Failed to delete snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("Service error", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Init test data
			svcErr := &service.Error{
				Type: service.InternalError,
				Base: errors.New("internal error"),
			}

			// ================================================
			// Describe mock calls
			mockService.EXPECT().SoftDelete(
				gomock.Any(),
				uint(1),
			).Return(svcErr)

			// ================================================
			// Run test
			expected := map[string]any{
				"error": "internal error",
			}

			response := expect.DELETE("/{id}", 1).
				Expect()

			response.
				Status(http.StatusInternalServerError).
				JSON().Object().IsEqual(expected)
		})

		t.Run("Bad request", func(t *testing.T) {
			t.Parallel()

			// ================================================
			// Init mocks and service
			ctrl := gomock.NewController(t)

			mockService := NewMockService(ctrl)
			transport := snippets.NewTransport(mockService, nopslog.NewNoplogger())
			handler := transport.Routes()

			// ================================================
			// Create httpexpect instance
			expect := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(handler),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			// ================================================
			// Run test
			expected := map[string]any{
				"error": `invalid id param: -1`,
			}

			response := expect.GET("/{id}", -1).
				Expect()

			response.
				Status(http.StatusBadRequest).
				JSON().Object().IsEqual(expected)
		})
	})
}
