package routes

import (
	"interview-user/handlers"
	"interview-user/middleware"
	"interview-user/models"
	"interview-user/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, jwtService *utils.JWTService) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, jwtService)
	saleOrderHandler := handlers.NewSaleOrderHandler(db)
	userHandler := handlers.NewUserHandler(db)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		utils.OKResponse(c, "Service is healthy", nil)
	})

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// Logout (requires auth)
		protected.POST("/auth/logout", authHandler.Logout)

		// Sale Orders - accessible by both cashier and owner
		saleOrders := protected.Group("/sale-orders")
		saleOrders.Use(middleware.RBACMiddleware(models.RoleCashier, models.RoleOwner))
		{
			saleOrders.GET("", saleOrderHandler.GetAll)
			saleOrders.GET("/:id", saleOrderHandler.GetByID)
			saleOrders.POST("", saleOrderHandler.Create)
			saleOrders.PATCH("/:id", saleOrderHandler.Update)
			saleOrders.DELETE("/:id", saleOrderHandler.Delete)
		}

		// User Cashier management - owner only
		users := protected.Group("/users/cashier")
		users.Use(middleware.RBACMiddleware(models.RoleOwner))
		{
			users.GET("", userHandler.GetAllCashiers)
			users.GET("/:id", userHandler.GetCashierByID)
			users.POST("", userHandler.CreateCashier)
			users.PATCH("/:id", userHandler.UpdateCashier)
			users.DELETE("/:id", userHandler.DeleteCashier)
		}
	}
}
