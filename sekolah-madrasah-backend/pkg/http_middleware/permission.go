package http_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sekolah-madrasah/pkg/auth_utils"
)

type PermissionChecker interface {
	HasPermission(userID string, permission string) bool
	GetUserPermissions(userID string) []string
}

var permissionChecker PermissionChecker

func SetPermissionChecker(checker PermissionChecker) {
	permissionChecker = checker
}

func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authData, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		claims, ok := authData.(*auth_utils.AuthClaim)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth claims"})
			c.Abort()
			return
		}

		if permissionChecker == nil {
			c.Next()
			return
		}

		for _, perm := range permissions {
			if permissionChecker.HasPermission(claims.UserID.String(), perm) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "insufficient permissions",
			"required": permissions,
		})
		c.Abort()
	}
}

func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return RequirePermission(permissions...)
}

func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authData, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		claims, ok := authData.(*auth_utils.AuthClaim)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth claims"})
			c.Abort()
			return
		}

		if permissionChecker == nil {
			c.Next()
			return
		}

		for _, perm := range permissions {
			if !permissionChecker.HasPermission(claims.UserID.String(), perm) {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "insufficient permissions",
					"missing": perm,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

type RoleLevel int

const (
	RoleLevelGuest RoleLevel = iota
	RoleLevelUser
	RoleLevelModerator
	RoleLevelAdmin
	RoleLevelSuperAdmin
)

func RequireRoleLevel(minLevel RoleLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		authData, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		claims, ok := authData.(*auth_utils.AuthClaim)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth claims"})
			c.Abort()
			return
		}

		userLevel := getRoleLevelFromClaim(claims)
		if userLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "insufficient role level",
				"required":      minLevel,
				"current_level": userLevel,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getRoleLevelFromClaim(claims *auth_utils.AuthClaim) RoleLevel {
	switch claims.Role {
	case "super_admin", "superadmin":
		return RoleLevelSuperAdmin
	case "admin":
		return RoleLevelAdmin
	case "moderator":
		return RoleLevelModerator
	case "user", "member":
		return RoleLevelUser
	default:
		return RoleLevelGuest
	}
}
