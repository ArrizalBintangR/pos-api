package handlers

import (
	"interview-user/middleware"
	"interview-user/models"
	"interview-user/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB         *gorm.DB
	JWTService *utils.JWTService
}

func NewAuthHandler(db *gorm.DB, jwtService *utils.JWTService) *AuthHandler {
	return &AuthHandler{
		DB:         db,
		JWTService: jwtService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint        `json:"id"`
	Username string      `json:"username"`
	Name     string      `json:"name"`
	Role     models.Role `json:"role"`
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	var user models.User
	if err := h.DB.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		utils.UnauthorizedResponse(c, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.UnauthorizedResponse(c, "Invalid username or password")
		return
	}

	token, err := h.JWTService.GenerateToken(&user)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token")
		return
	}

	utils.OKResponse(c, "Login successful", LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
		},
	})
}

// Logout handles user logout by blacklisting the token
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		utils.BadRequestResponse(c, "Token not found")
		return
	}

	middleware.BlacklistToken(token.(string))
	utils.OKResponse(c, "Logout successful", nil)
}
