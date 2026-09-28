package user_test

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
	"i_m_s/internal/user"

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
		db, dbErr := database.Connect(cfg.TestDSN(), false)
		if dbErr == nil {
			if migErr := ensureUsersTable(db); migErr == nil {
				testDB = db
			} else {
				fmt.Printf("user repository tests: failed to prepare schema: %v\n", migErr)
			}
		} else {
			fmt.Printf("user repository tests: test database unavailable, skipping: %v\n", dbErr)
		}
	} else {
		fmt.Printf("user repository tests: config load failed: %v\n", err)
	}
	os.Exit(m.Run())
}

// ensureUsersTable applies the users migration when the table is absent so the
// integration tests are self-contained. It never drops an existing table.
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
		t.Skip("test database unavailable; set TEST_DATABASE_URL")
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

func TestRepository_CreateAndLookups(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

	account := newUser("ada", "ada@example.com")
	require.NoError(t, repo.Create(ctx, tx, account))
	assert.NotEqual(t, uuid.Nil, account.ID)

	byID, err := repo.GetByID(ctx, tx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, "ada", byID.Username)

	byEmail, err := repo.GetByEmail(ctx, tx, "ada@example.com")
	require.NoError(t, err)
	assert.Equal(t, account.ID, byEmail.ID)

	byUsername, err := repo.GetByUsername(ctx, tx, "ada")
	require.NoError(t, err)
	assert.Equal(t, account.ID, byUsername.ID)
}

func TestRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

	_, err := repo.GetByID(ctx, tx, uuid.New())
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetByEmail(ctx, tx, "nobody@example.com")
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetByUsername(ctx, tx, "nobody")
	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestRepository_Update(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

	account := newUser("grace", "grace@example.com")
	require.NoError(t, repo.Create(ctx, tx, account))

	account.Role = models.Manager
	account.IsActive = false
	require.NoError(t, repo.Update(ctx, tx, account))

	got, err := repo.GetByID(ctx, tx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, models.Manager, got.Role)
	assert.False(t, got.IsActive)
}

func TestRepository_Delete(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

	account := newUser("linus", "linus@example.com")
	require.NoError(t, repo.Create(ctx, tx, account))

	require.NoError(t, repo.Delete(ctx, tx, account.ID))

	_, err := repo.GetByID(ctx, tx, account.ID)
	assert.ErrorIs(t, err, models.ErrNotFound)

	// Deleting an already soft-deleted row reports not found.
	assert.ErrorIs(t, repo.Delete(ctx, tx, account.ID), models.ErrNotFound)
}

func TestRepository_ListPagination(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

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

func TestRepository_DuplicateUsernameRejected(t *testing.T) {
	ctx := context.Background()
	tx := newTestTx(t)
	repo := user.NewRepository()

	require.NoError(t, repo.Create(ctx, tx, newUser("dup", "dup1@example.com")))

	err := repo.Create(ctx, tx, newUser("dup", "dup2@example.com"))
	assert.ErrorIs(t, err, models.ErrConflict)
}
