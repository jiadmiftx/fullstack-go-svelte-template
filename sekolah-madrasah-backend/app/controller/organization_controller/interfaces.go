package organization_controller

import "github.com/gin-gonic/gin"

type OrganizationController interface {
	GetOrganization(c *gin.Context)
	GetOrganizations(c *gin.Context)
	CreateOrganization(c *gin.Context)
	UpdateOrganization(c *gin.Context)
	DeleteOrganization(c *gin.Context)

	GetMembers(c *gin.Context)
	AddMember(c *gin.Context)
	UpdateMember(c *gin.Context)
	RemoveMember(c *gin.Context)
}
