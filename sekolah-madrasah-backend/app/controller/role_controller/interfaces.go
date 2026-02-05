package role_controller

import "github.com/gin-gonic/gin"

type RoleController interface {
	GetRole(c *gin.Context)
	GetRoles(c *gin.Context)
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}
