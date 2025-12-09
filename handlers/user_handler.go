package handlers

import (
	"strconv"

	"interview-user/models"
	"interview-user/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,max=255"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=100"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Name     string `json:"name" binding:"omitempty,max=255"`
	IsActive *bool  `json:"is_active"`
}

type CashierResponse struct {
	ID        uint        `json:"id"`
	Username  string      `json:"username"`
	Name      string      `json:"name"`
	Role      models.Role `json:"role"`
	IsActive  bool        `json:"is_active"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// GetAllCashiers returns all cashier users with pagination
func (h *UserHandler) GetAllCashiers(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var total int64
	var users []models.User

	// Count total cashiers
	if err := h.DB.Model(&models.User{}).Where("role = ?", models.RoleCashier).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count users")
		return
	}

	// Get paginated data
	if err := h.DB.Where("role = ?", models.RoleCashier).
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&users).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch users")
		return
	}

	// Convert to response format
	var cashiers []CashierResponse
	for _, u := range users {
		cashiers = append(cashiers, CashierResponse{
			ID:        u.ID,
			Username:  u.Username,
			Name:      u.Name,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	utils.OKResponse(c, "Cashiers retrieved successfully", utils.PaginatedResponse{
		Items:      cashiers,
		TotalItems: total,
		TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		Page:       pagination.Page,
		Limit:      pagination.Limit,
	})
}

// GetCashierByID returns a cashier user by ID
func (h *UserHandler) GetCashierByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND role = ?", id, models.RoleCashier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Cashier not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch user")
		return
	}

	utils.OKResponse(c, "Cashier retrieved successfully", CashierResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// CreateCashier creates a new cashier user
func (h *UserHandler) CreateCashier(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequestResponse(c, "Username already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to hash password")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     models.RoleCashier,
		IsActive: true,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create user")
		return
	}

	utils.CreatedResponse(c, "Cashier created successfully", CashierResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateCashier updates a cashier user
func (h *UserHandler) UpdateCashier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND role = ?", id, models.RoleCashier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Cashier not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch user")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update fields if provided
	if req.Username != "" {
		// Check if username already exists (for another user)
		var existingUser models.User
		if err := h.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			utils.BadRequestResponse(c, "Username already exists")
			return
		}
		user.Username = req.Username
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.InternalServerErrorResponse(c, "Failed to hash password")
			return
		}
		user.Password = string(hashedPassword)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.DB.Save(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user")
		return
	}

	utils.OKResponse(c, "Cashier updated successfully", CashierResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// DeleteCashier soft deletes a cashier user
func (h *UserHandler) DeleteCashier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND role = ?", id, models.RoleCashier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Cashier not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch user")
		return
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete user")
		return
	}

	utils.OKResponse(c, "Cashier deleted successfully", nil)
}
