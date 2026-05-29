package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

// Auth creates a JWT authentication middleware.
func Auth(db *gorm.DB, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Load user from database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// InternalKey validates the X-Internal-Key header against INTERNAL_API_KEY env var.
func InternalKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Key")
		if key == "" || key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// BusinessAccess validates the X-Business-ID header and checks user membership.
func BusinessAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		businessIDStr := c.GetHeader("X-Business-ID")
		if businessIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Business-ID header is required"})
			c.Abort()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Check user has a role in this business
		var role models.UserBusinessRole
		if err := db.Where("user_id = ? AND business_id = ?", userID, businessIDStr).First(&role).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this business"})
			c.Abort()
			return
		}

		c.Set("business_id", role.BusinessID)
		c.Set("user_role", role.Role)
		c.Set("user_branch_id", role.BranchID)
		c.Next()
	}
}

// RoleRequired checks if the user has one of the allowed roles.
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
		c.Abort()
	}
}
