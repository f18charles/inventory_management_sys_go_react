# Phase 3 Execution Plan: Authentication & User Management

This document details the step-by-step technical plan for implementing Phase 3 of the Inventory Management System backend.

---

## Overview & Scope

Phase 3 implements authentication, user management, identity verification, password hashing, JWT token generation/validation, authentication middleware, role-based access control, and user management API endpoints.

All implementations strictly adhere to the layered architecture defined in `AGENTS.md` and the `inventory-backend` skill guidelines:

```
Gin Handler → Service → Repository → GORM → PostgreSQL
```

---

## Phase 3 Implementation Steps

### Step 1: User Repository (`internal/repositories`)
- **Interface & GORM Implementation**:
  - `Create(ctx context.Context, user *models.User) error`
  - `GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)`
  - `GetByEmail(ctx context.Context, email string) (*models.User, error)`
  - `GetByUsername(ctx context.Context, username string) (*models.User, error)`
  - `Update(ctx context.Context, user *models.User) error`
  - `Delete(ctx context.Context, id uuid.UUID) error`
  - `List(ctx context.Context, page, pageSize int) ([]models.User, int64, error)`
- **Repository Tests**: Unit tests executing CRUD operations against a test database instance.

---

### Step 2: Password Security & JWT Utilities (`internal/utils/auth`)
- **Password Hashing**: Secure hashing and password verification routines utilizing `golang.org/x/crypto/bcrypt`.
- **JWT Token Management**: Token generation, signing, and parsing utilizing `github.com/golang-jwt/jwt/v5` containing `user_id`, `role`, and expiration claims.

---

### Step 3: Service Layer (`internal/services`)
- **`AuthService` & `UserService`**:
  - `Login(ctx context.Context, usernameOrEmail, password string)`: Validates active status and credentials, issuing a JWT token upon success.
  - `CreateUser(ctx context.Context, input CreateUserInput)`: Enforces unique email/username constraints and hashes passwords before passing to the repository layer.
  - `GetUserByID(ctx context.Context, id uuid.UUID)` & `ListUsers(ctx context.Context, page, pageSize int)`.
  - `UpdateUserRole(ctx context.Context, id uuid.UUID, role models.Roles)`.
- **Service Tests**: Testify unit tests mocking the repository interface to verify business rules, inactive account rejections, and invalid credential handling.

---

### Step 4: Middleware Layer (`internal/middleware`)
- **`JWTAuthMiddleware(jwtSecret string)`**: Extracts `Authorization: Bearer <token>`, verifies signature/expiration, and attaches `user_id` and `user_role` to the Gin context.
- **`RequireRole(allowedRoles ...models.Roles)`**: Enforces role-based authorization rules (`Admin`, `Manager`, `Staff`).

---

### Step 5: Handler Layer & Endpoints (`internal/handlers`)
- **`AuthHandler` & `UserHandler`**:
  - `POST /api/v1/auth/login`: User login endpoint returning `{ "data": { "token": "...", "user": ... } }`.
  - `GET /api/v1/auth/me`: Retrieves current authenticated user profile.
  - `GET /api/v1/users`: List users (Admin/Manager).
  - `POST /api/v1/users`: Create user (Admin).
  - `PATCH /api/v1/users/:id/role`: Update user role (Admin).
- **Handler Tests**: `httptest` + Gin test mode asserting HTTP status codes and JSON response envelopes.

---

### Step 6: Documentation & Validation
- Add Swaggo OpenAPI annotations above all new endpoints and regenerate Swagger specs.
- Update `docs/Roadmap.md` checklist for Phase 3.

---

After phase three is completed and just before the branch is pushed to the remote repository, This document is to be deleted
