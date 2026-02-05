package user_controller

import "github.com/gin-gonic/gin"

type UserController interface {
	GetUser(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	GetUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetMyMemberships(c *gin.Context)
}
