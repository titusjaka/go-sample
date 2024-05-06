package snippets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/titusjaka/go-sample/v2/internal/business/snippets"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/nopslog"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/service"
)

func TestSnippetService_Create(t *testing.T) {
	t.Parallel()

	t.Run("Successfully create a new snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("All times in UTC", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			fakeNow := time.Now().UTC()

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return fakeNow },
			)

			// ===============================================
			// Init test data
			expiresAt := fakeNow.Add(time.Hour * 24)

			snippetToCreate := snippets.Snippet{
				Title:     "Best snippet ever",
				Content:   "Some text here…",
				ExpiresAt: expiresAt,
			}

			snippetPassedToStorage := snippets.Snippet{
				Title:     snippetToCreate.Title,
				Content:   snippetToCreate.Content,
				CreatedAt: fakeNow,
				UpdatedAt: fakeNow,
				ExpiresAt: expiresAt,
			}

			snippetID := uint(200)

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().Create(ctx, snippetPassedToStorage).Return(snippetID, nil)

			// ===============================================
			// Run Test
			expectedSnippet := snippetPassedToStorage
			expectedSnippet.ID = snippetID

			actual, svcErr := snippetService.Create(ctx, snippetToCreate)

			require.Nil(t, svcErr)
			assert.Equal(t, expectedSnippet, actual)
		})

		t.Run("Times in different TZ", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			fakeNow := time.Now().In(time.FixedZone("CET", 60*60))

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return fakeNow },
			)

			// ===============================================
			// Init test data
			expiresAt := fakeNow.Add(time.Hour * 24).In(time.FixedZone("EET", 60*60*2))

			snippetToCreate := snippets.Snippet{
				Title:     "Best snippet ever",
				Content:   "Some text here…",
				ExpiresAt: expiresAt,
			}

			snippetPassedToStorage := snippets.Snippet{
				Title:     snippetToCreate.Title,
				Content:   snippetToCreate.Content,
				CreatedAt: fakeNow,
				UpdatedAt: fakeNow,
				ExpiresAt: expiresAt.UTC(),
			}

			snippetID := uint(200)

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().Create(ctx, snippetPassedToStorage).Return(snippetID, nil)

			// ===============================================
			// Run Test
			expectedSnippet := snippetPassedToStorage
			expectedSnippet.ID = snippetID

			actual, svcErr := snippetService.Create(ctx, snippetToCreate)

			require.Nil(t, svcErr)
			assert.Equal(t, expectedSnippet, actual)
		})
	})

	t.Run("Failed to create snippet", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)

		// ===============================================
		// Init Mocks and Service
		mockStorage := NewMockStorage(ctrl)

		fakeNow := time.Now().UTC()

		snippetService := snippets.NewService(
			mockStorage,
			nopslog.NewNoplogger(),
			func() time.Time { return fakeNow },
		)

		// ===============================================
		// Init test data
		expectedErr := errors.New("something wrong happen 😱")

		// ===============================================
		// Describe Mock Calls
		mockStorage.EXPECT().Create(ctx, gomock.Any()).Return(0, expectedErr)

		// ===============================================
		// Run Test
		_, svcErr := snippetService.Create(ctx, snippets.Snippet{})

		require.NotNil(t, svcErr)
		assert.Equal(t, service.InternalError, svcErr.Type)
		assert.ErrorIs(t, svcErr, expectedErr)
	})
}

func TestSnippetService_Get(t *testing.T) {
	t.Parallel()

	t.Run("Successfully get a snippet", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)

		// ===============================================
		// Init Mocks and Service
		mockStorage := NewMockStorage(ctrl)

		snippetService := snippets.NewService(
			mockStorage,
			nopslog.NewNoplogger(),
			func() time.Time { return time.Now().UTC() },
		)

		// ===============================================
		// Init test data
		snippet := snippets.Snippet{
			ID:        200,
			Title:     "Best snippet ever",
			Content:   "Some text here…",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
		}

		// ===============================================
		// Describe Mock Calls
		mockStorage.EXPECT().Get(ctx, snippet.ID).Return(snippet, nil)

		// ===============================================
		// Run Test
		actual, svcErr := snippetService.Get(ctx, 200)

		require.Nil(t, svcErr)
		assert.Equal(t, snippet, actual)
	})

	t.Run("Failed to get a snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("Internal error", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			expectedErr := errors.New("failed to get")

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().Get(ctx, uint(200)).Return(snippets.Snippet{}, expectedErr)

			actual, svcErr := snippetService.Get(ctx, 200)
			require.NotNil(t, svcErr)
			assert.Empty(t, actual)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.ErrorIs(t, svcErr, expectedErr)
		})

		t.Run("Not found", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().Get(ctx, uint(200)).Return(snippets.Snippet{}, snippets.ErrNotFound)

			// ===============================================
			// Run Test
			actual, svcErr := snippetService.Get(ctx, 200)

			require.NotNil(t, svcErr)
			assert.Empty(t, actual)
			assert.Equal(t, service.NotFound, svcErr.Type)
			assert.ErrorIs(t, svcErr, snippets.ErrNotFound)
		})
	})
}

func TestSnippetService_List(t *testing.T) {
	t.Parallel()

	t.Run("Successfully list snippets", func(t *testing.T) {
		t.Parallel()

		t.Run("Total is less than limit", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			listOfSnippets := []snippets.Snippet{
				{
					ID:        1,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
				},
				{
					ID:        2,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: time.Now().UTC().Add(time.Hour),
					UpdatedAt: time.Now().UTC().Add(time.Hour),
					ExpiresAt: time.Now().UTC().Add(time.Hour * 48),
				},
			}
			limit := uint(10)
			offset := uint(0)

			total := uint(2)

			pagination := snippets.NewPagination(limit, offset, total)

			// ===============================================
			// Describe Mock Calls
			gomock.InOrder(
				mockStorage.EXPECT().Total(ctx).Return(total, nil),
				mockStorage.EXPECT().List(ctx, pagination).Return(listOfSnippets, nil),
			)

			// ===============================================
			// Run Test
			actualSnippets, actualPagination, svcErr := snippetService.List(ctx, limit, offset)

			require.Nil(t, svcErr)
			assert.Equal(t, listOfSnippets, actualSnippets)
			assert.Equal(t, pagination, actualPagination)
		})

		t.Run("Total is greater than limit", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			listOfSnippets := []snippets.Snippet{
				{
					ID:        1,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
				},
				{
					ID:        2,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: time.Now().UTC().Add(time.Hour),
					UpdatedAt: time.Now().UTC().Add(time.Hour),
					ExpiresAt: time.Now().UTC().Add(time.Hour * 48),
				},
			}
			limit := uint(2)
			offset := uint(0)

			total := uint(20)

			pagination := snippets.NewPagination(limit, offset, total)

			// ===============================================
			// Describe Mock Calls
			gomock.InOrder(
				mockStorage.EXPECT().Total(ctx).Return(total, nil),
				mockStorage.EXPECT().List(ctx, pagination).Return(listOfSnippets, nil),
			)

			// ===============================================
			// Run Test
			actualSnippets, actualPagination, svcErr := snippetService.List(ctx, limit, offset)

			require.Nil(t, svcErr)
			assert.Equal(t, listOfSnippets, actualSnippets)
			assert.Equal(t, pagination, actualPagination)
		})
	})

	t.Run("Failed to list snippets", func(t *testing.T) {
		t.Parallel()

		t.Run("List returned error", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			limit := uint(10)
			offset := uint(0)
			total := uint(10)

			expectedErr := errors.New("OMG!!! VERY BAD 🤯")

			pagination := snippets.NewPagination(limit, offset, total)

			// ===============================================
			// Describe Mock Calls
			gomock.InOrder(
				mockStorage.EXPECT().Total(ctx).Return(total, nil),
				mockStorage.EXPECT().List(ctx, pagination).Return([]snippets.Snippet{}, expectedErr),
			)

			// ===============================================
			// Run Test
			actualSnippets, _, svcErr := snippetService.List(ctx, limit, offset)

			assert.Empty(t, actualSnippets)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.ErrorIs(t, svcErr, expectedErr)
		})

		t.Run("Total returned error", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			limit := uint(10)
			offset := uint(0)

			expectedErr := errors.New("OMG!!! VERY BAD 🤯")

			// ===============================================
			// Describe Mock Calls
			gomock.InOrder(
				mockStorage.EXPECT().Total(ctx).Return(0, expectedErr),
			)

			// ===============================================
			// Run Test
			actualSnippets, _, svcErr := snippetService.List(ctx, limit, offset)

			assert.Empty(t, actualSnippets)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.ErrorIs(t, svcErr, expectedErr)
		})
	})
}

func TestSnippetService_SoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("Successfully soft-delete snippet", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)

		// ===============================================
		// Init Mocks and Service
		mockStorage := NewMockStorage(ctrl)

		snippetService := snippets.NewService(
			mockStorage,
			nopslog.NewNoplogger(),
			func() time.Time { return time.Now().UTC() },
		)

		// ===============================================
		// Init test data
		id := uint(200)

		// ===============================================
		// Describe Mock Calls
		mockStorage.EXPECT().SoftDelete(ctx, id).Return(nil)

		// ===============================================
		// Run Test
		svcErr := snippetService.SoftDelete(ctx, id)

		require.Nil(t, svcErr)
	})

	t.Run("Failed to soft-delete snippet", func(t *testing.T) {
		t.Parallel()

		t.Run("Not found", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			id := uint(200)
			expectedErr := snippets.ErrNotFound

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().SoftDelete(ctx, id).Return(expectedErr)

			// ===============================================
			// Run Test
			svcErr := snippetService.SoftDelete(ctx, id)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.NotFound, svcErr.Type)
			assert.ErrorIs(t, svcErr, expectedErr)
		})

		t.Run("Internal error", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			ctrl := gomock.NewController(t)

			// ===============================================
			// Init Mocks and Service
			mockStorage := NewMockStorage(ctrl)

			snippetService := snippets.NewService(
				mockStorage,
				nopslog.NewNoplogger(),
				func() time.Time { return time.Now().UTC() },
			)

			// ===============================================
			// Init test data
			id := uint(200)
			expectedErr := errors.New("this is not suppose to happen 🚑")

			// ===============================================
			// Describe Mock Calls
			mockStorage.EXPECT().SoftDelete(ctx, id).Return(expectedErr)

			// ===============================================
			// Run Test
			svcErr := snippetService.SoftDelete(ctx, id)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.ErrorIs(t, svcErr, expectedErr)
		})
	})
}
