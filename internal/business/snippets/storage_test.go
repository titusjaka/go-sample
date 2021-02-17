package snippets_test

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/lib/pq" // import pg driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
	"github.com/titusjaka/go-sample/migrate"
)

type stopFunc func() error

const (
	defaultUsername     = "test"
	defaultPassword     = "test"
	defaultPort         = 9876
	defaultStartTimeout = 30 * time.Second
)

type connectionParams struct {
	username     string
	password     string
	databaseName string
	port         uint32
	timeout      time.Duration
}

var (
	defaultParams = connectionParams{
		username: defaultUsername,
		password: defaultPassword,
		port:     defaultPort,
		timeout:  defaultStartTimeout,
	}
)

func TestPGStorage_Create(t *testing.T) {
	params := defaultParams
	params.databaseName = fmt.Sprintf("snippets_test_%d", time.Now().UnixNano())

	pgStorage, stop, err := initPGStorage(params)
	require.NoError(t, err)

	defer func() {
		err := stop()
		require.NoError(t, err)
	}()

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

			id, err := pgStorage.Create(context.Background(), snippet)
			require.NoError(t, err)
			assert.Equal(t, uint(1), id)
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

			id, err := pgStorage.Create(context.Background(), snippet)
			require.NoError(t, err)
			assert.Equal(t, uint(2), id)
		})
	})

	t.Run("Failed to create a new snippet", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
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

		id, err := pgStorage.Create(ctx, snippet)
		require.Error(t, err)
		assert.Zero(t, id)
	})
}

func TestPGStorage_Get(t *testing.T) {
	params := defaultParams
	params.databaseName = fmt.Sprintf("snippets_test_%d", time.Now().UnixNano())

	pgStorage, stop, err := initPGStorage(params)
	require.NoError(t, err)

	defer func() {
		err := stop()
		require.NoError(t, err)
	}()

	t.Run("Successfully get a snippet", func(t *testing.T) {
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
			id, err := pgStorage.Create(context.Background(), snippet1)
			require.NoError(t, err)
			assert.Equal(t, uint(1), id)

			id, err = pgStorage.Create(context.Background(), snippet2)
			require.NoError(t, err)
			assert.Equal(t, uint(2), id)
		})

		t.Run("Get snippet #1", func(t *testing.T) {
			actualSnippet, err := pgStorage.Get(context.Background(), 1)
			require.NoError(t, err)
			assert.Equal(t, snippet1.ID, actualSnippet.ID)
			assert.Equal(t, snippet1.Title, actualSnippet.Title)
			assert.Equal(t, snippet1.Content, actualSnippet.Content)
			assert.Equal(t, snippet1.CreatedAt, actualSnippet.CreatedAt.UTC())
			assert.Equal(t, snippet1.UpdatedAt, actualSnippet.UpdatedAt.UTC())
			assert.Equal(t, snippet1.ExpiresAt, actualSnippet.ExpiresAt.UTC())
		})

		t.Run("Get snippet #2", func(t *testing.T) {
			actualSnippet, err := pgStorage.Get(context.Background(), 2)
			require.NoError(t, err)
			assert.Equal(t, snippet2.ID, actualSnippet.ID)
			assert.Equal(t, snippet2.Title, actualSnippet.Title)
			assert.Equal(t, snippet2.Content, actualSnippet.Content)
			assert.Equal(t, snippet2.CreatedAt, actualSnippet.CreatedAt.UTC())
			assert.Equal(t, snippet2.UpdatedAt, actualSnippet.UpdatedAt.UTC())
			assert.Equal(t, snippet2.ExpiresAt, actualSnippet.ExpiresAt.UTC())
		})
	})

	t.Run("Failed to get a new snippet", func(t *testing.T) {
		t.Run("Error occurred", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			defer cancel()

			snippet, err := pgStorage.Get(ctx, 1)
			require.Error(t, err)
			assert.Zero(t, snippet)
		})

		t.Run("Not Found", func(t *testing.T) {
			snippet, err := pgStorage.Get(context.Background(), 3)
			require.Equal(t, snippets.ErrNotFound, err)
			assert.Zero(t, snippet)
		})
	})
}

func TestPGStorage_List(t *testing.T) {
	params := defaultParams
	params.databaseName = fmt.Sprintf("snippets_test_%d", time.Now().UnixNano())

	pgStorage, stop, err := initPGStorage(params)
	require.NoError(t, err)

	defer func() {
		err := stop()
		require.NoError(t, err)
	}()

	t.Run("Successfully list snippets", func(t *testing.T) {
		t.Run("Empty database", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  100,
				Offset: 0,
			}

			actualSnippets, err := pgStorage.List(context.Background(), expectedPagination)
			require.NoError(t, err)
			assert.Empty(t, actualSnippets)
		})

		expectedSnippets := make([]snippets.Snippet, 0, 10)

		startTime := time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC)
		startExpiresTime := time.Date(2500, 1, 1, 1, 1, 1, 0, time.UTC)
		for i := 0; i < 10; i++ {
			expectedSnippets = append(expectedSnippets, snippets.Snippet{
				ID:        uint(i + 1),
				Title:     fmt.Sprintf("Very important snippet #%d", i+1),
				Content:   "Some kind of content",
				CreatedAt: startTime.Add(time.Hour * time.Duration(i)),
				UpdatedAt: startTime.Add(time.Hour * time.Duration(i)),
				ExpiresAt: startExpiresTime.Add(time.Hour * time.Duration(i)),
			})
		}
		sortedByCreatedDateDesc := make([]snippets.Snippet, 10)
		copy(sortedByCreatedDateDesc, expectedSnippets)

		sort.Slice(sortedByCreatedDateDesc, func(i, j int) bool {
			return sortedByCreatedDateDesc[i].CreatedAt.After(sortedByCreatedDateDesc[j].CreatedAt)
		})

		t.Run("Create snippets", func(t *testing.T) {
			for i := range expectedSnippets {
				id, err := pgStorage.Create(context.Background(), expectedSnippets[i])
				require.NoError(t, err)
				assert.Equal(t, uint(i+1), id)
			}
		})

		t.Run("List snippets with default params", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  100,
				Offset: 0,
			}

			actualSnippets, err := pgStorage.List(context.Background(), expectedPagination)
			require.NoError(t, err)
			for i := range actualSnippets {
				assert.Equal(t, sortedByCreatedDateDesc[i].ID, actualSnippets[i].ID)
				assert.Equal(t, sortedByCreatedDateDesc[i].Title, actualSnippets[i].Title)
				assert.Equal(t, sortedByCreatedDateDesc[i].Content, actualSnippets[i].Content)
				assert.Equal(t, sortedByCreatedDateDesc[i].CreatedAt, actualSnippets[i].CreatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i].UpdatedAt, actualSnippets[i].UpdatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i].ExpiresAt, actualSnippets[i].ExpiresAt.UTC())
			}
		})

		t.Run("List snippets with certain limit", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  5,
				Offset: 0,
			}

			actualSnippets, err := pgStorage.List(context.Background(), expectedPagination)
			require.NoError(t, err)
			for i := range actualSnippets {
				assert.Equal(t, sortedByCreatedDateDesc[i].ID, actualSnippets[i].ID)
				assert.Equal(t, sortedByCreatedDateDesc[i].Title, actualSnippets[i].Title)
				assert.Equal(t, sortedByCreatedDateDesc[i].Content, actualSnippets[i].Content)
				assert.Equal(t, sortedByCreatedDateDesc[i].CreatedAt, actualSnippets[i].CreatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i].UpdatedAt, actualSnippets[i].UpdatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i].ExpiresAt, actualSnippets[i].ExpiresAt.UTC())
			}
		})

		t.Run("List snippets with certain limit and offset", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  5,
				Offset: 5,
			}

			actualSnippets, err := pgStorage.List(context.Background(), expectedPagination)
			require.NoError(t, err)
			for i := range actualSnippets {
				assert.Equal(t, sortedByCreatedDateDesc[i+5].ID, actualSnippets[i].ID)
				assert.Equal(t, sortedByCreatedDateDesc[i+5].Title, actualSnippets[i].Title)
				assert.Equal(t, sortedByCreatedDateDesc[i+5].Content, actualSnippets[i].Content)
				assert.Equal(t, sortedByCreatedDateDesc[i+5].CreatedAt, actualSnippets[i].CreatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i+5].UpdatedAt, actualSnippets[i].UpdatedAt.UTC())
				assert.Equal(t, sortedByCreatedDateDesc[i+5].ExpiresAt, actualSnippets[i].ExpiresAt.UTC())
			}
		})
	})

	t.Run("Failed to list snippets", func(t *testing.T) {
		t.Run("Error occurred", func(t *testing.T) {
			expectedPagination := service.Pagination{
				Limit:  100,
				Offset: 0,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			defer cancel()

			actualSnippets, err := pgStorage.List(ctx, expectedPagination)
			require.Error(t, err)
			assert.Empty(t, actualSnippets)
		})
	})
}

func TestPGStorage_SoftDelete(t *testing.T) {
	params := defaultParams
	params.databaseName = fmt.Sprintf("snippets_test_%d", time.Now().UnixNano())

	pgStorage, stop, err := initPGStorage(params)
	require.NoError(t, err)

	defer func() {
		err := stop()
		require.NoError(t, err)
	}()

	t.Run("Successfully soft delete a snippet", func(t *testing.T) {
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
			id, err := pgStorage.Create(context.Background(), snippet1)
			require.NoError(t, err)
			assert.Equal(t, uint(1), id)

			id, err = pgStorage.Create(context.Background(), snippet2)
			require.NoError(t, err)
			assert.Equal(t, uint(2), id)
		})

		t.Run("Soft delete snippet #1", func(t *testing.T) {
			err := pgStorage.SoftDelete(context.Background(), 1)
			require.NoError(t, err)
		})
	})

	t.Run("Failed to soft delete a snippet", func(t *testing.T) {
		t.Run("Error occurred", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			defer cancel()

			err := pgStorage.SoftDelete(ctx, 2)
			require.Error(t, err)
		})

		t.Run("Not Found", func(t *testing.T) {
			err := pgStorage.SoftDelete(context.Background(), 3)
			require.Equal(t, snippets.ErrNotFound, err)
		})
	})
}

func TestPGStorage_Total(t *testing.T) {
	params := defaultParams
	params.databaseName = fmt.Sprintf("snippets_test_%d", time.Now().UnixNano())

	pgStorage, stop, err := initPGStorage(params)
	require.NoError(t, err)

	defer func() {
		err := stop()
		require.NoError(t, err)
	}()

	t.Run("Successfully count snippets", func(t *testing.T) {
		t.Run("Empty database", func(t *testing.T) {
			count, err := pgStorage.Total(context.Background())
			require.NoError(t, err)
			require.Zero(t, count)
		})

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
			id, err := pgStorage.Create(context.Background(), snippet1)
			require.NoError(t, err)
			assert.Equal(t, uint(1), id)

			id, err = pgStorage.Create(context.Background(), snippet2)
			require.NoError(t, err)
			assert.Equal(t, uint(2), id)
		})

		t.Run("Count snippets", func(t *testing.T) {
			count, err := pgStorage.Total(context.Background())
			require.NoError(t, err)
			assert.Equal(t, uint(2), count)
		})
	})

	t.Run("Failed to count snippets", func(t *testing.T) {
		t.Run("Error occurred", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			defer cancel()

			count, err := pgStorage.Total(ctx)
			require.Error(t, err)
			assert.Zero(t, count)
		})
	})
}

func TestPGStorage_Errors(t *testing.T) {
	t.Run("List on not initialized DB", func(t *testing.T) {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:12345/test_db?sslmode=disable")
		require.NoError(t, err)

		pgStorage := snippets.NewPGStorage(db)

		_, err = pgStorage.List(context.Background(), service.Pagination{})
		require.Error(t, err)
	})

	t.Run("Get snippet on not initialized DB", func(t *testing.T) {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:12345/test_db?sslmode=disable")
		require.NoError(t, err)

		pgStorage := snippets.NewPGStorage(db)

		_, err = pgStorage.Get(context.Background(), 0)
		require.Error(t, err)
	})

	t.Run("Create a snippet on not initialized DB", func(t *testing.T) {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:12345/test_db?sslmode=disable")
		require.NoError(t, err)

		pgStorage := snippets.NewPGStorage(db)

		_, err = pgStorage.Create(context.Background(), snippets.Snippet{})
		require.Error(t, err)
	})

	t.Run("Soft delete a snippet on not initialized DB", func(t *testing.T) {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:12345/test_db?sslmode=disable")
		require.NoError(t, err)

		pgStorage := snippets.NewPGStorage(db)

		err = pgStorage.SoftDelete(context.Background(), 0)
		require.Error(t, err)
	})

	t.Run("Count snippets on not initialized DB", func(t *testing.T) {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:12345/test_db?sslmode=disable")
		require.NoError(t, err)

		pgStorage := snippets.NewPGStorage(db)

		_, err = pgStorage.Total(context.Background())
		require.Error(t, err)
	})
}

func initPGStorage(params connectionParams) (storage *snippets.PGStorage, stop stopFunc, err error) {
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(params.username).
		Password(params.password).
		Database(params.databaseName).
		Version(embeddedpostgres.V13).
		Port(params.port).
		StartTimeout(params.timeout))

	if err = postgres.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start PostgreSQL: %w", err)
	}

	defer func() {
		// If any error happens, must stop PostgreSQL service
		if err != nil {
			stopErr := postgres.Stop()
			if stopErr != nil {
				err = fmt.Errorf("failed to stop PostgreSQL: %w. Previous error: %s", stopErr, err.Error())
			}
		}
	}()

	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%d/%s?sslmode=disable",
		params.username,
		params.password,
		params.port,
		params.databaseName,
	)

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	_, err = applyMigrations(context.Background(), db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	stop = func() (err error) {
		if closeErr := db.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close db connection: %w", closeErr)
		}
		if stopErr := postgres.Stop(); stopErr != nil {
			var previousErrMessage string
			if err != nil {
				previousErrMessage = ": " + err.Error()
			}
			err = fmt.Errorf("failed to stop PostgreSQL: %w%s", stopErr, previousErrMessage)
		}

		return err
	}

	return snippets.NewPGStorage(db), stop, nil
}

func applyMigrations(ctx context.Context, db *sql.DB) (int, error) {
	source := migrate.NewEmbeddedSource()
	migrator := migrate.NewMigrator(db, source)
	return migrator.Up(ctx)
}
