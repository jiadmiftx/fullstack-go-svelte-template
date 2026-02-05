package permission_controller

import "github.com/gin-gonic/gin"

type PermissionController interface {
	GetPermission(c *gin.Context)
	GetPermissions(c *gin.Context)
	CreatePermission(c *gin.Context)
	DeletePermission(c *gin.Context)
}
