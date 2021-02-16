package snippets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

func TestSnippetService_Create(t *testing.T) {
	t.Run("Successfully create a new snippet", func(t *testing.T) {
		t.Run("All times in UTC", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			now := time.Now().UTC()
			expiresAt := now.Add(time.Hour * 24)

			expectedID := uint(200)
			expectedResult := snippets.Snippet{
				ID:        expectedID,
				Title:     "Best snippet ever",
				Content:   "Some text here…",
				CreatedAt: now,
				UpdatedAt: now,
				ExpiresAt: expiresAt,
			}

			mockStorage.EXPECT().Create(ctx, gomock.Any()).Return(expectedID, nil)

			gotSnippet, svcErr := snippetService.Create(ctx, snippets.Snippet{
				Title:     expectedResult.Title,
				Content:   expectedResult.Content,
				ExpiresAt: expectedResult.ExpiresAt,
			})

			require.Nil(t, svcErr)
			assert.Equal(t, expectedID, gotSnippet.ID)
			assert.Equal(t, expectedResult.Title, gotSnippet.Title)
			assert.Equal(t, expectedResult.Content, gotSnippet.Content)
			assert.Equal(t, expectedResult.ExpiresAt, gotSnippet.ExpiresAt)
			assert.True(t, gotSnippet.CreatedAt.After(now))
			assert.True(t, gotSnippet.CreatedAt.Before(now.Add(time.Minute)))
		})

		t.Run("Times in different TZ", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			now := time.Now().In(time.FixedZone("CET", 60*60))
			expiresAt := now.Add(time.Hour * 24).In(time.FixedZone("EET", 60*60*2))

			expectedID := uint(200)
			expectedResult := snippets.Snippet{
				ID:        expectedID,
				Title:     "Best snippet ever",
				Content:   "Some text here…",
				CreatedAt: now,
				UpdatedAt: now,
				ExpiresAt: expiresAt.UTC(),
			}

			mockStorage.EXPECT().Create(ctx, gomock.Any()).Return(expectedID, nil)

			gotSnippet, svcErr := snippetService.Create(ctx, snippets.Snippet{
				Title:     expectedResult.Title,
				Content:   expectedResult.Content,
				ExpiresAt: expectedResult.ExpiresAt,
			})

			require.Nil(t, svcErr)
			assert.Equal(t, expectedID, gotSnippet.ID)
			assert.Equal(t, expectedResult.Title, gotSnippet.Title)
			assert.Equal(t, expectedResult.Content, gotSnippet.Content)
			assert.Equal(t, expectedResult.ExpiresAt, gotSnippet.ExpiresAt)
			assert.True(t, gotSnippet.CreatedAt.After(now))
			assert.True(t, gotSnippet.CreatedAt.Before(now.Add(time.Minute)))
		})
	})

	t.Run("Error during creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := NewMockStorage(ctrl)
		mockLogger := log.NewMockLogger(ctrl)

		snippetService := snippets.NewService(mockStorage, mockLogger)

		ctx := context.Background()

		expectedID := uint(0)
		expectedErr := errors.New("something wrong happen 😱")

		mockStorage.EXPECT().Create(ctx, gomock.Any()).Return(expectedID, expectedErr)
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

		_, svcErr := snippetService.Create(ctx, snippets.Snippet{})

		require.NotNil(t, svcErr)
		assert.Equal(t, service.InternalError, svcErr.Type)
		assert.True(t, errors.Is(svcErr, expectedErr))
	})
}

func TestSnippetService_Get(t *testing.T) {
	t.Run("Successfully get a snippet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := NewMockStorage(ctrl)
		mockLogger := log.NewMockLogger(ctrl)

		snippetService := snippets.NewService(mockStorage, mockLogger)

		ctx := context.Background()
		now := time.Now().UTC()
		expiresAt := now.Add(time.Hour * 24)

		expectedID := uint(200)
		expectedResult := snippets.Snippet{
			ID:        expectedID,
			Title:     "Best snippet ever",
			Content:   "Some text here…",
			CreatedAt: now,
			UpdatedAt: now,
			ExpiresAt: expiresAt,
		}

		mockStorage.EXPECT().Get(ctx, expectedID).Return(expectedResult, nil)

		gotSnippet, svcErr := snippetService.Get(ctx, 200)

		require.Nil(t, svcErr)
		assert.Equal(t, expectedID, gotSnippet.ID)
		assert.Equal(t, expectedResult.Title, gotSnippet.Title)
		assert.Equal(t, expectedResult.Content, gotSnippet.Content)
		assert.Equal(t, expectedResult.ExpiresAt, gotSnippet.ExpiresAt)
		assert.Equal(t, expectedResult.CreatedAt, gotSnippet.CreatedAt)
	})

	t.Run("Failed to get a snippet", func(t *testing.T) {
		t.Run("Internal error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			expectedErr := errors.New("failed to get")

			mockStorage.EXPECT().Get(ctx, uint(200)).Return(snippets.Snippet{}, expectedErr)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

			gotSnippet, svcErr := snippetService.Get(ctx, 200)
			require.NotNil(t, svcErr)
			assert.Equal(t, snippets.Snippet{}, gotSnippet)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.True(t, errors.Is(svcErr, expectedErr))
		})

		t.Run("Not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()

			mockStorage.EXPECT().Get(ctx, uint(200)).Return(snippets.Snippet{}, snippets.ErrNotFound)

			gotSnippet, svcErr := snippetService.Get(ctx, 200)
			require.NotNil(t, svcErr)
			assert.Equal(t, snippets.Snippet{}, gotSnippet)
			assert.Equal(t, service.NotFound, svcErr.Type)
			assert.True(t, errors.Is(svcErr, snippets.ErrNotFound))
		})
	})
}

func TestSnippetService_List(t *testing.T) {
	t.Run("Successfully list snippets", func(t *testing.T) {
		t.Run("Total is less than limit", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			now1 := time.Now().UTC()
			now2 := now1.Add(time.Hour)
			expiresAt1 := now1.Add(time.Hour * 24)
			expiresAt2 := expiresAt1.Add(time.Hour * 48)

			expectedResult := []snippets.Snippet{
				{
					ID:        1,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: now1,
					UpdatedAt: now1,
					ExpiresAt: expiresAt1,
				},
				{
					ID:        2,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: now2,
					UpdatedAt: now2,
					ExpiresAt: expiresAt2,
				},
			}
			limit := uint(10)
			offset := uint(0)

			expectedPagination := service.Pagination{
				Limit:       limit,
				Offset:      offset,
				Total:       2,
				TotalPages:  1,
				CurrentPage: 1,
			}

			mockStorage.EXPECT().Total(ctx).Return(expectedPagination.Total, nil)
			mockStorage.EXPECT().List(ctx, expectedPagination).Return(expectedResult, nil)

			actualSnippets, pagination, svcErr := snippetService.List(ctx, limit, offset)
			require.Nil(t, svcErr)
			assert.EqualValues(t, expectedResult, actualSnippets)
			assert.Equal(t, expectedPagination, pagination)
		})

		t.Run("Total is greater than limit", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			now1 := time.Now().UTC()
			now2 := now1.Add(time.Hour)
			expiresAt1 := now1.Add(time.Hour * 24)
			expiresAt2 := expiresAt1.Add(time.Hour * 48)

			expectedResult := []snippets.Snippet{
				{
					ID:        1,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: now1,
					UpdatedAt: now1,
					ExpiresAt: expiresAt1,
				},
				{
					ID:        2,
					Title:     "Best snippet ever",
					Content:   "Some text here…",
					CreatedAt: now2,
					UpdatedAt: now2,
					ExpiresAt: expiresAt2,
				},
			}
			limit := uint(2)
			offset := uint(4)

			expectedPagination := service.Pagination{
				Limit:       limit,
				Offset:      offset,
				Total:       20,
				TotalPages:  10,
				CurrentPage: 3,
			}

			mockStorage.EXPECT().Total(ctx).Return(expectedPagination.Total, nil)
			mockStorage.EXPECT().List(ctx, expectedPagination).Return(expectedResult, nil)

			actualSnippets, pagination, svcErr := snippetService.List(ctx, limit, offset)
			require.Nil(t, svcErr)
			assert.EqualValues(t, expectedResult, actualSnippets)
			assert.Equal(t, expectedPagination, pagination)
		})

		t.Run("Empty list", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()

			var expectedResult []snippets.Snippet
			limit := uint(2)
			offset := uint(40)

			expectedPagination := service.Pagination{
				Limit:       limit,
				Offset:      offset,
				Total:       20,
				TotalPages:  10,
				CurrentPage: 21,
			}

			mockStorage.EXPECT().Total(ctx).Return(expectedPagination.Total, nil)
			mockStorage.EXPECT().List(ctx, expectedPagination).Return(expectedResult, nil)

			actualSnippets, pagination, svcErr := snippetService.List(ctx, limit, offset)
			require.Nil(t, svcErr)
			assert.EqualValues(t, expectedResult, actualSnippets)
			assert.Equal(t, expectedPagination, pagination)
		})
	})

	t.Run("Failed to list snippets", func(t *testing.T) {
		t.Run("List returned error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()

			limit := uint(10)
			offset := uint(0)

			expectedErr := errors.New("OMG!!! VERY BAD 🤯")

			expectedPagination := service.Pagination{
				Limit:       limit,
				Offset:      offset,
				CurrentPage: 1,
			}

			mockStorage.EXPECT().Total(ctx).Return(uint(0), nil)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
			mockStorage.EXPECT().List(ctx, expectedPagination).Return([]snippets.Snippet{}, expectedErr)

			actualSnippets, pagination, svcErr := snippetService.List(ctx, limit, offset)
			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.True(t, errors.Is(svcErr, expectedErr))
			assert.Nil(t, actualSnippets)
			assert.Equal(t, expectedPagination, pagination)
		})

		t.Run("Total returned error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()

			limit := uint(10)
			offset := uint(0)

			expectedErr := errors.New("OMG!!! VERY BAD 🤯")

			expectedPagination := service.Pagination{
				Limit:  0,
				Offset: 0,
			}

			mockStorage.EXPECT().Total(ctx).Return(uint(0), expectedErr)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

			actualSnippets, pagination, svcErr := snippetService.List(ctx, limit, offset)
			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.True(t, errors.Is(svcErr, expectedErr))
			assert.Nil(t, actualSnippets)
			assert.Equal(t, expectedPagination, pagination)
		})
	})
}

func TestSnippetService_SoftDelete(t *testing.T) {
	t.Run("Successfully soft-delete snippet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := NewMockStorage(ctrl)
		mockLogger := log.NewMockLogger(ctrl)

		snippetService := snippets.NewService(mockStorage, mockLogger)

		ctx := context.Background()
		expectedID := uint(200)

		mockStorage.EXPECT().SoftDelete(ctx, expectedID).Return(nil)

		svcErr := snippetService.SoftDelete(ctx, expectedID)

		require.Nil(t, svcErr)
	})

	t.Run("Failed to soft-delete snippet", func(t *testing.T) {
		t.Run("Not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			expectedID := uint(200)
			expectedErr := snippets.ErrNotFound

			mockStorage.EXPECT().SoftDelete(ctx, expectedID).Return(expectedErr)

			svcErr := snippetService.SoftDelete(ctx, expectedID)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.NotFound, svcErr.Type)
			assert.Equal(t, expectedErr, svcErr.Base)
		})

		t.Run("Internal error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockLogger := log.NewMockLogger(ctrl)

			snippetService := snippets.NewService(mockStorage, mockLogger)

			ctx := context.Background()
			expectedID := uint(200)
			expectedErr := errors.New("this is not suppose to happen 🚑")

			mockStorage.EXPECT().SoftDelete(ctx, expectedID).Return(expectedErr)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

			svcErr := snippetService.SoftDelete(ctx, expectedID)

			require.NotNil(t, svcErr)
			assert.Equal(t, service.InternalError, svcErr.Type)
			assert.True(t, errors.Is(svcErr, expectedErr))
		})
	})
}
