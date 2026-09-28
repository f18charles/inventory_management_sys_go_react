package handlers

import (
	"net/http"

	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	"i_m_s/internal/services"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler exposes authentication endpoints. It only binds requests, calls
// the auth service, and shapes HTTP responses.
type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required,min=1"`
}

// LoginResponse is the payload returned by a successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Verifies username/email and password, then returns a JWT and the user profile.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login credentials"
// @Success      200 {object} object{data=handlers.LoginResponse}
// @Failure      400 {object} object{error=object{code=string,message=string}}
// @Failure      401 {object} object{error=object{code=string,message=string}}
// @Failure      403 {object} object{error=object{code=string,message=string}}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if !BindAndValidate(c, &req) {
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.UsernameOrEmail, req.Password)
	if err != nil {
		RespondError(c, err)
		return
	}

	response.Success(c, http.StatusOK, LoginResponse{Token: token, User: user})
}

// Me godoc
// @Summary      Get the current user
// @Description  Returns the profile of the authenticated user.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} object{data=models.User}
// @Failure      401 {object} object{error=object{code=string,message=string}}
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// currentUserID reads the authenticated user id placed on the context by the
// JWT middleware.
func currentUserID(c *gin.Context) (uuid.UUID, error) {
	raw := c.GetString(middleware.ContextUserID)
	if raw == "" {
		return uuid.Nil, models.ErrUnauthorized
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, models.ErrUnauthorized
	}
	return id, nil
}
