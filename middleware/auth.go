package middleware

import (
	"strings"
	"sync"

	"interview-user/models"
	"interview-user/utils"

	"github.com/gin-gonic/gin"
)

// TokenBlacklist stores invalidated tokens
var TokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]bool
}{tokens: make(map[string]bool)}

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check if token is blacklisted
		TokenBlacklist.RLock()
		if TokenBlacklist.tokens[tokenString] {
			TokenBlacklist.RUnlock()
			utils.UnauthorizedResponse(c, "Token has been invalidated")
			c.Abort()
			return
		}
		TokenBlacklist.RUnlock()

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		c.Next()
	}
}

// RBACMiddleware checks if user has required role
func RBACMiddleware(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			utils.UnauthorizedResponse(c, "User role not found")
			c.Abort()
			return
		}

		userRole := roleValue.(models.Role)

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c, "You don't have permission to access this resource")
		c.Abort()
	}
}

// BlacklistToken adds a token to the blacklist
func BlacklistToken(token string) {
	TokenBlacklist.Lock()
	TokenBlacklist.tokens[token] = true
	TokenBlacklist.Unlock()
}
