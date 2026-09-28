// internal/tests/user/service_test.go
//
// Service tests mock the repository INTERFACE (the payoff of DIP) and assert
// business outcomes, not call choreography. Shared mocks live in
// internal/tests/mocks; repository tests that need real SQL live in
// internal/tests/<domain>/repository_test.go and run against the local test DB.
//
// Note: a method that opens a db.Transaction cannot run with a nil *gorm.DB.
// Cover transaction-owning flows with a repository/integration test against the
// local PostgreSQL test database (TEST_DATABASE_URL) rather than this mock style.

package user_test

import (
	"context"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"
	"i_m_s/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateUserRejectsDuplicateUsername(t *testing.T) {
	// Arrange: the repository reports the username already exists.
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(&models.User{}, nil)

	// Act
	_, err := user.NewService(nil, repo).CreateUser(context.Background(), user.CreateUserInput{
		FirstName: "Grace",
		LastName:  "Hopper",
		Username:  "grace",
		Email:     "grace@example.com",
		Password:  "plaintext-password",
		Role:      models.Staff,
	})

	// Assert: the business rule was enforced, and no write happened.
	assert.ErrorIs(t, err, models.ErrConflict)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

// Table-driven pattern for multiple cases on the same behavior:
//
//	func TestService_CreateUser(t *testing.T) {
//		cases := []struct {
//			name      string
//			userExists bool
//			wantErr   error
//		}{
//			{"unique username succeeds", false, nil},
//			{"duplicate username conflicts", true, models.ErrConflict},
//		}
//		for _, tc := range cases {
//			t.Run(tc.name, func(t *testing.T) {
//				// set up mocks per-case using tc, then assert errors.Is(err, tc.wantErr)
//				// or require.NoError(t, err).
//			})
//		}
//	}
