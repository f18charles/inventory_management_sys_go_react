package repositories_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"i_m_s/internal/config"
	"i_m_s/internal/database"
	"i_m_s/internal/models"
	"i_m_s/internal/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// testDB is the shared connection used by every case in this package. When it
// is nil (no database reachable) each test skips rather than fails.
var testDB *gorm.DB

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err == nil {
		db, dbErr := database.Connect(cfg.DSN(), false)
		if dbErr == nil {
			if migErr := ensureUsersTable(db); migErr == nil {
				testDB = db
			} else {
				fmt.Printf("repository tests: failed to prepare schema: %v\n", migErr)
			}
		} else {
			fmt.Printf("repository tests: database unavailable, skipping: %v\n", dbErr)
		}
	} else {
		fmt.Printf("repository tests: config load failed: %v\n", err)
	}
	os.Exit(m.Run())
}

// ensureUsersTable applies the users migration when the table is absent so the
// integration tests are self-contained in CI. It never drops an existing table.
func ensureUsersTable(db *gorm.DB) error {
	if db.Migrator().HasTable("users") {
		return nil
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", "000001_create_users.up.sql"))
	if err != nil {
		return err
	}
	for _, stmt := range strings.Split(string(raw), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}

// newTestTx hands each test an isolated transaction that is always rolled back,
// keeping the shared table clean without destructive truncation.
func newTestTx(t *testing.T) *gorm.DB {
	t.Helper()
	if testDB == nil {
		t.Skip("test database unavailable; set DB_HOST/DB_USER/DB_PASSWORD/DB_NAME")
	}
	tx := testDB.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })
	return tx
}

func newUser(username, email string) *models.User {
	return &models.User{
		FirstName: "Test",
		LastName:  "User",
		Username:  username,
		Email:     email,
		PassHash:  "hashed-secret",
		Role:      models.Staff,
		IsActive:  true,
	}
}

func TestUserRepository_CreateAndLookups(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	user := newUser("ada", "ada@example.com")
	require.NoError(t, repo.Create(ctx, tx, user))
	assert.NotEqual(t, uuid.Nil, user.ID)

	byID, err := repo.GetByID(ctx, tx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "ada", byID.Username)

	byEmail, err := repo.GetByEmail(ctx, tx, "ada@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, byEmail.ID)

	byUsername, err := repo.GetByUsername(ctx, tx, "ada")
	require.NoError(t, err)
	assert.Equal(t, user.ID, byUsername.ID)
}

func TestUserRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	_, err := repo.GetByID(ctx, tx, uuid.New())
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetByEmail(ctx, tx, "nobody@example.com")
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetByUsername(ctx, tx, "nobody")
	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestUserRepository_Update(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	user := newUser("grace", "grace@example.com")
	require.NoError(t, repo.Create(ctx, tx, user))

	user.Role = models.Manager
	user.IsActive = false
	require.NoError(t, repo.Update(ctx, tx, user))

	got, err := repo.GetByID(ctx, tx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, models.Manager, got.Role)
	assert.False(t, got.IsActive)
}

func TestUserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	user := newUser("linus", "linus@example.com")
	require.NoError(t, repo.Create(ctx, tx, user))

	require.NoError(t, repo.Delete(ctx, tx, user.ID))

	_, err := repo.GetByID(ctx, tx, user.ID)
	assert.ErrorIs(t, err, models.ErrNotFound)

	// Deleting an already soft-deleted row reports not found.
	assert.ErrorIs(t, repo.Delete(ctx, tx, user.ID), models.ErrNotFound)
}

func TestUserRepository_ListPagination(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	for _, name := range []string{"u1", "u2", "u3"} {
		require.NoError(t, repo.Create(ctx, tx, newUser(name, name+"@example.com")))
	}

	firstPage, total, err := repo.List(ctx, tx, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, firstPage, 2)

	secondPage, total, err := repo.List(ctx, tx, 2, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, secondPage, 1)
}

func TestUserRepository_DuplicateUsernameRejected(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := repositories.NewUserRepository()

	require.NoError(t, repo.Create(ctx, tx, newUser("dup", "dup1@example.com")))

	err := repo.Create(ctx, tx, newUser("dup", "dup2@example.com"))
	assert.ErrorIs(t, err, models.ErrConflict)
}
