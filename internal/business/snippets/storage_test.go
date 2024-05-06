package snippets_test

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titusjaka/go-sample/v2/internal/business/snippets"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres/pgtest"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/service"
)

const envFile = "../../../.env"

//nolint:gosec // this is a fake DSN, no password or sensitive data here
const fakePostgresDSN = "postgres://fake:fake@localhost:5432/fake"

func TestPGStorage_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test due to 'short' flag")
	}
	t.Parallel()

	pgConn := pgtest.InitTestDatabase(
		t,
		pgtest.WithConfigFiles(envFile),
	)

	ctx := context.Background()
	pgStorage := snippets.NewPGStorage(pgConn)

	t.Run("Successfully create a new snippet", func(t *testing.T) {
		t.Run("With id 1", func(t *testing.T) {
			fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			snippet := snippets.Snippet{
				Title:     "Snippet title #1",
				Content:   "Very important content",
				CreatedAt: fakeTimeCreated,
				UpdatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}

			id, err := pgStorage.Create(ctx, snippet)
			require.NoError(t, err)
			assert.EqualValues(t, 1, id)
		})

		t.Run("With id 2", func(t *testing.T) {
			fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			snippet := snippets.Snippet{
				Title:     "Snippet title #2",
				Content:   "Very important content",
				CreatedAt: fakeTimeCreated,
				UpdatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}

			id, err := pgStorage.Create(ctx, snippet)
			require.NoError(t, err)
			assert.EqualValues(t, 2, id)
		})
	})

	t.Run("Handle errors", func(t *testing.T) {
		t.Run("Create a snippet on not initialized DB", func(t *testing.T) {
			db, err := sql.Open("pgx/v5", fakePostgresDSN)
			require.NoError(t, err)

			fakePG := snippets.NewPGStorage(db)

			_, err = fakePG.Create(context.Background(), snippets.Snippet{})
			require.Error(t, err)
		})

		t.Run("Context timeout", func(t *testing.T) {
			expiredCtx, cancel := context.WithTimeout(ctx, time.Nanosecond)
			defer cancel()

			fakeTimeCreated := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			snippet := snippets.Snippet{
				Title:     "Snippet title #1",
				Content:   "Very important content",
				CreatedAt: fakeTimeCreated,
				UpdatedAt: fakeTimeCreated,
				ExpiresAt: fakeTimeExpires,
			}

			id, err := pgStorage.Create(expiredCtx, snippet)
			require.Error(t, err)
			assert.Empty(t, id)
		})
	})
}

func TestPGStorage_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test due to 'short' flag")
	}
	t.Parallel()

	pgConn := pgtest.InitTestDatabase(
		t,
		pgtest.WithConfigFiles(envFile),
	)

	ctx := context.Background()
	pgStorage := snippets.NewPGStorage(pgConn)

	t.Run("Successfully get a snippet", func(t *testing.T) {
		// Initialize snippets
		fakeTimeCreated1 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
		fakeTimeExpires1 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
		snippet1 := snippets.Snippet{
			ID:        1,
			Title:     "Snippet title #1",
			Content:   "Very important content",
			CreatedAt: fakeTimeCreated1,
			UpdatedAt: fakeTimeCreated1,
			ExpiresAt: fakeTimeExpires1,
		}

		fakeTimeCreated2 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
		fakeTimeExpires2 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
		snippet2 := snippets.Snippet{
			ID:        2,
			Title:     "Snippet title #2",
			Content:   "Very important content",
			CreatedAt: fakeTimeCreated2,
			UpdatedAt: fakeTimeCreated2,
			ExpiresAt: fakeTimeExpires2,
		}

		t.Run("Create snippets", func(t *testing.T) {
			id, err := pgStorage.Create(ctx, snippet1)
			require.NoError(t, err)
			assert.Equal(t, uint(1), id)

			id, err = pgStorage.Create(ctx, snippet2)
			require.NoError(t, err)
			assert.Equal(t, uint(2), id)
		})

		t.Run("Get snippet #1", func(t *testing.T) {
			actualSnippet, err := pgStorage.Get(ctx, 1)
			require.NoError(t, err)
			assert.Equal(t, snippet1, actualSnippet)
		})

		t.Run("Get snippet #2", func(t *testing.T) {
			actualSnippet, err := pgStorage.Get(ctx, 2)
			require.NoError(t, err)
			assert.Equal(t, snippet2, actualSnippet)
		})
	})

	t.Run("Handle errors", func(t *testing.T) {
		t.Run("Get snippet on not initialized DB", func(t *testing.T) {
			db, err := sql.Open("pgx/v5", fakePostgresDSN)
			require.NoError(t, err)

			fakePG := snippets.NewPGStorage(db)

			_, err = fakePG.Get(context.Background(), 0)
			require.Error(t, err)
		})

		t.Run("Context timeout", func(t *testing.T) {
			expiredCtx, cancel := context.WithTimeout(ctx, time.Nanosecond)
			defer cancel()

			snippet, err := pgStorage.Get(expiredCtx, 1)
			require.Error(t, err)
			assert.Empty(t, snippet)
		})

		t.Run("Not Found", func(t *testing.T) {
			snippet, err := pgStorage.Get(ctx, 3)
			require.ErrorIs(t, err, snippets.ErrNotFound)
			assert.Empty(t, snippet)
		})
	})
}

func TestPGStorage_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test due to 'short' flag")
	}
	t.Parallel()

	pgConn := pgtest.InitTestDatabase(
		t,
		pgtest.WithConfigFiles(envFile),
	)

	ctx := context.Background()
	pgStorage := snippets.NewPGStorage(pgConn)

	t.Run("Successfully list snippets", func(t *testing.T) {
		t.Run("Empty list on empty DB", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  100,
				Offset: 0,
			}

			actualSnippets, err := pgStorage.List(ctx, expectedPagination)
			require.NoError(t, err)
			assert.Empty(t, actualSnippets)
		})

		t.Run("Fill the DB and return snippets", func(t *testing.T) {
			// Initialize snippets
			createdSnippets := make([]snippets.Snippet, 0, 10)

			startTime := time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC)
			startExpiresTime := time.Date(2500, 1, 1, 1, 1, 1, 0, time.UTC)

			for i := 0; i < 10; i++ {
				createdSnippets = append(createdSnippets, snippets.Snippet{
					ID:        uint(i + 1),
					Title:     fmt.Sprintf("Very important snippet #%d", i+1),
					Content:   "Some kind of content",
					CreatedAt: startTime.Add(time.Hour * -time.Duration(i)),
					UpdatedAt: startTime.Add(time.Hour * -time.Duration(i)),
					ExpiresAt: startExpiresTime.Add(time.Hour * -time.Duration(i)),
				})
			}

			t.Run("Create snippets", func(t *testing.T) {
				for i := range createdSnippets {
					id, err := pgStorage.Create(ctx, createdSnippets[i])
					require.NoError(t, err)
					assert.Equal(t, uint(i+1), id)
				}
			})

			t.Run("List snippets with the default pagination", func(t *testing.T) {
				expectedPagination := service.Pagination{
					Limit:  100,
					Offset: 0,
				}

				actualSnippets, err := pgStorage.List(ctx, expectedPagination)
				require.NoError(t, err)

				expectedSnippets := slices.Clone(createdSnippets)
				assert.Equal(t, expectedSnippets, actualSnippets)
			})

			t.Run("List snippets with a certain limit", func(t *testing.T) {
				expectedPagination := service.Pagination{
					Limit:  5,
					Offset: 0,
				}

				actualSnippets, err := pgStorage.List(ctx, expectedPagination)
				require.NoError(t, err)

				expectedSnippets := slices.Clone(createdSnippets[:5])
				assert.Equal(t, expectedSnippets, actualSnippets)
			})

			t.Run("List snippets with certain limit and offset", func(t *testing.T) {
				expectedPagination := service.Pagination{
					Limit:  5,
					Offset: 5,
				}

				actualSnippets, err := pgStorage.List(ctx, expectedPagination)
				require.NoError(t, err)

				expectedSnippets := slices.Clone(createdSnippets[5:])
				assert.Equal(t, expectedSnippets, actualSnippets)
			})
		})
	})

	t.Run("Handle errors", func(t *testing.T) {
		t.Run("List on not initialized DB", func(t *testing.T) {
			db, err := sql.Open("pgx/v5", fakePostgresDSN)
			require.NoError(t, err)

			fakePG := snippets.NewPGStorage(db)

			_, err = fakePG.List(context.Background(), service.Pagination{})
			require.Error(t, err)
		})

		t.Run("Context timeout", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  100,
				Offset: 0,
			}

			expiredCtx, cancel := context.WithTimeout(ctx, time.Nanosecond)
			defer cancel()

			actualSnippets, err := pgStorage.List(expiredCtx, expectedPagination)
			require.Error(t, err)
			assert.Empty(t, actualSnippets)
		})
	})
}

func TestPGStorage_SoftDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test due to 'short' flag")
	}
	t.Parallel()

	pgConn := pgtest.InitTestDatabase(
		t,
		pgtest.WithConfigFiles(envFile),
	)

	ctx := context.Background()
	pgStorage := snippets.NewPGStorage(pgConn)

	t.Run("Successfully soft delete a snippet", func(t *testing.T) {
		// Initialize snippets
		fakeTimeCreated1 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
		fakeTimeExpires1 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
		snippet1 := snippets.Snippet{
			ID:        1,
			Title:     "Snippet title #1",
			Content:   "Very important content",
			CreatedAt: fakeTimeCreated1,
			UpdatedAt: fakeTimeCreated1,
			ExpiresAt: fakeTimeExpires1,
		}

		fakeTimeCreated2 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
		fakeTimeExpires2 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
		snippet2 := snippets.Snippet{
			ID:        2,
			Title:     "Snippet title #2",
			Content:   "Very important content",
			CreatedAt: fakeTimeCreated2,
			UpdatedAt: fakeTimeCreated2,
			ExpiresAt: fakeTimeExpires2,
		}

		t.Run("Create snippets", func(t *testing.T) {
			id, err := pgStorage.Create(ctx, snippet1)
			require.NoError(t, err)
			assert.EqualValues(t, 1, id)

			id, err = pgStorage.Create(ctx, snippet2)
			require.NoError(t, err)
			assert.EqualValues(t, 2, id)
		})

		t.Run("Soft delete snippet #1", func(t *testing.T) {
			err := pgStorage.SoftDelete(ctx, 1)
			require.NoError(t, err)
		})
	})

	t.Run("Handle errors", func(t *testing.T) {
		t.Run("Soft delete a snippet on not initialized DB", func(t *testing.T) {
			db, err := sql.Open("pgx/v5", fakePostgresDSN)
			require.NoError(t, err)

			fakePG := snippets.NewPGStorage(db)

			err = fakePG.SoftDelete(context.Background(), 0)
			require.Error(t, err)
		})

		t.Run("Context timeout", func(t *testing.T) {
			expiredCtx, cancel := context.WithTimeout(ctx, time.Nanosecond)
			defer cancel()

			err := pgStorage.SoftDelete(expiredCtx, 2)
			require.Error(t, err)
		})

		t.Run("Not Found", func(t *testing.T) {
			err := pgStorage.SoftDelete(ctx, 3)
			require.Equal(t, snippets.ErrNotFound, err)
		})
	})
}

func TestPGStorage_Total(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test due to 'short' flag")
	}
	t.Parallel()

	pgConn := pgtest.InitTestDatabase(
		t,
		pgtest.WithConfigFiles(envFile),
	)

	ctx := context.Background()
	pgStorage := snippets.NewPGStorage(pgConn)

	t.Run("Successfully count snippets", func(t *testing.T) {
		t.Run("Empty database", func(t *testing.T) {
			count, err := pgStorage.Total(ctx)
			require.NoError(t, err)
			require.Zero(t, count)
		})

		t.Run("Fill the DB and count snippets", func(t *testing.T) {
			// Initialize snippets
			fakeTimeCreated1 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires1 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			snippet1 := snippets.Snippet{
				ID:        1,
				Title:     "Snippet title #1",
				Content:   "Very important content",
				CreatedAt: fakeTimeCreated1,
				UpdatedAt: fakeTimeCreated1,
				ExpiresAt: fakeTimeExpires1,
			}

			fakeTimeCreated2 := time.Date(2020, 10, 7, 12, 0, 0, 0, time.UTC)
			fakeTimeExpires2 := time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)
			snippet2 := snippets.Snippet{
				ID:        2,
				Title:     "Snippet title #2",
				Content:   "Very important content",
				CreatedAt: fakeTimeCreated2,
				UpdatedAt: fakeTimeCreated2,
				ExpiresAt: fakeTimeExpires2,
			}

			t.Run("Create snippets", func(t *testing.T) {
				id, err := pgStorage.Create(ctx, snippet1)
				require.NoError(t, err)
				assert.EqualValues(t, 1, id)

				id, err = pgStorage.Create(ctx, snippet2)
				require.NoError(t, err)
				assert.EqualValues(t, 2, id)
			})

			t.Run("Count snippets", func(t *testing.T) {
				count, err := pgStorage.Total(ctx)
				require.NoError(t, err)
				assert.EqualValues(t, 2, count)
			})
		})
	})

	t.Run("Handle errors", func(t *testing.T) {
		t.Run("Count snippets on not initialized DB", func(t *testing.T) {
			db, err := sql.Open("pgx/v5", fakePostgresDSN)
			require.NoError(t, err)

			fakePG := snippets.NewPGStorage(db)

			_, err = fakePG.Total(context.Background())
			require.Error(t, err)
		})

		t.Run("Context timeout", func(t *testing.T) {
			expiredCtx, cancel := context.WithTimeout(ctx, time.Nanosecond)
			defer cancel()

			count, err := pgStorage.Total(expiredCtx)
			require.Error(t, err)
			assert.Zero(t, count)
		})
	})
}
