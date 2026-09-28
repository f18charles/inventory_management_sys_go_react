package handlers

import (
	"net/http"
	"strconv"

	"i_m_s/internal/models"
	"i_m_s/internal/services"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler exposes user management endpoints for administrators.
type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type createUserRequest struct {
	FirstName string       `json:"first_name" binding:"required,min=1"`
	LastName  string       `json:"last_name" binding:"required,min=1"`
	Username  string       `json:"username" binding:"required,min=3"`
	Email     string       `json:"email" binding:"required,email"`
	Password  string       `json:"password" binding:"required,min=8"`
	Role      models.Roles `json:"role" binding:"required,oneof=admin manager staff"`
}

type updateRoleRequest struct {
	Role models.Roles `json:"role" binding:"required,oneof=admin manager staff"`
}

// Create godoc
// @Summary      Create a user
// @Description  Creates a new user account. Admin only.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body createUserRequest true "New user details"
// @Success      201 {object} object{data=models.User}
// @Failure      400 {object} object{error=object{code=string,message=string}}
// @Failure      401 {object} object{error=object{code=string,message=string}}
// @Failure      403 {object} object{error=object{code=string,message=string}}
// @Failure      409 {object} object{error=object{code=string,message=string}}
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if !BindAndValidate(c, &req) {
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), services.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
	})
	if err != nil {
		RespondError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// List godoc
// @Summary      List users
// @Description  Returns a paginated list of users. Admin/Manager only.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Items per page" default(20)
// @Success      200 {object} object{data=[]models.User,meta=response.PaginationMeta}
// @Failure      401 {object} object{error=object{code=string,message=string}}
// @Failure      403 {object} object{error=object{code=string,message=string}}
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 20)

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		RespondError(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, users, response.NewPaginationMeta(page, pageSize, total))
}

// UpdateRole godoc
// @Summary      Update a user's role
// @Description  Changes the role assigned to a user. Admin only.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        request body updateRoleRequest true "New role"
// @Success      200 {object} object{data=models.User}
// @Failure      400 {object} object{error=object{code=string,message=string}}
// @Failure      401 {object} object{error=object{code=string,message=string}}
// @Failure      403 {object} object{error=object{code=string,message=string}}
// @Failure      404 {object} object{error=object{code=string,message=string}}
// @Router       /users/{id}/role [patch]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid user id")
		return
	}

	var req updateRoleRequest
	if !BindAndValidate(c, &req) {
		return
	}

	user, err := h.userService.UpdateUserRole(c.Request.Context(), id, req.Role)
	if err != nil {
		RespondError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

func queryInt(c *gin.Context, key string, fallback int) int {
	if raw := c.Query(key); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			return value
		}
	}
	return fallback
}
