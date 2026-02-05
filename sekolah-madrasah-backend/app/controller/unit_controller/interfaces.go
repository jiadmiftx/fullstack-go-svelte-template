package unit_controller

import "github.com/gin-gonic/gin"

type GinUnitController interface {
	GetUnit(c *gin.Context)
	GetUnits(c *gin.Context)
	CreateUnit(c *gin.Context)
	UpdateUnit(c *gin.Context)
	DeleteUnit(c *gin.Context)
}
