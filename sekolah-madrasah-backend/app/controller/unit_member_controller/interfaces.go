package unit_member_controller

import "github.com/gin-gonic/gin"

type GinUnitMemberController interface {
	GetMembers(c *gin.Context)
	GetMember(c *gin.Context)
	AddMember(c *gin.Context)
	UpdateMember(c *gin.Context)
	RemoveMember(c *gin.Context)
}
