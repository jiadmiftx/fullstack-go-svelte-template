package organization_controller

import (
	"net/http"
	"reflect"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sekolah-madrasah/app/use_case/organization_use_case"
	"sekolah-madrasah/pkg/auth_utils"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"
)

type organizationController struct {
	orgUseCase organization_use_case.OrganizationUseCase
}

func NewOrganizationController(orgUseCase organization_use_case.OrganizationUseCase) OrganizationController {
	return &organizationController{orgUseCase: orgUseCase}
}

func (ctrl *organizationController) toOrgResponse(o organization_use_case.Organization) Organization {
	return Organization{
		Id:          o.Id,
		OwnerId:     o.OwnerId,
		Name:        o.Name,
		Code:        o.Code,
		Type:        o.Type,
		Description: o.Description,
		Address:     o.Address,
		Logo:        o.Logo,
		IsActive:    o.IsActive,
		MemberCount: o.MemberCount,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func (ctrl *organizationController) toMemberResponse(m organization_use_case.OrganizationMember) OrganizationMember {
	return OrganizationMember{
		Id:             m.Id,
		UserId:         m.UserId,
		OrganizationId: m.OrganizationId,
		RoleId:         m.RoleId,
		IsActive:       m.IsActive,
		JoinedAt:       m.JoinedAt,
		InvitedBy:      m.InvitedBy,
	}
}

// GetOrganization godoc
// @Summary Get organization by ID
// @Description Retrieves a single organization by its UUID
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=Organization}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id} [get]
func (ctrl *organizationController) GetOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	org, code, err := ctrl.orgUseCase.GetOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    ctrl.toOrgResponse(org),
	})
}

// GetOrganizations godoc
// @Summary List all organizations
// @Description Retrieves a paginated list of organizations with optional filtering
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param type query string false "Filter by organization type"
// @Param owner_id query string false "Filter by owner ID"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]Organization}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations [get]
func (ctrl *organizationController) GetOrganizations(c *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := organization_use_case.OrganizationFilter{}

	if orgType := c.Query("type"); orgType != "" {
		filter.Type = &orgType
	}
	if ownerId := c.Query("owner_id"); ownerId != "" {
		if parsedId, err := uuid.Parse(ownerId); err == nil {
			filter.OwnerId = &parsedId
		}
	}

	orgs, code, err := ctrl.orgUseCase.GetOrganizations(c.Request.Context(), filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]Organization, len(orgs))
	for i, o := range orgs {
		result[i] = ctrl.toOrgResponse(o)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Creates a new organization with the authenticated user as owner
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrganizationRequest true "Organization creation data"
// @Success 201 {object} gin_utils.DataResponse{data=Organization}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations [post]
func (ctrl *organizationController) CreateOrganization(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "unauthorized"})
		return
	}
	claims, ok := authData.(*auth_utils.AuthClaim)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "unauthorized"})
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(100),
		}).
		SetRule("code", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(50),
		}).
		SetRule("type", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(50),
			Null: true,
		}).
		SetRule("description", map_validator.Rules{
			Type: reflect.String,
			Null: true,
		}).
		SetRule("address", map_validator.Rules{
			Type: reflect.String,
			Null: true,
		})

	jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
	jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	jsonData, err := jsonDataValidate.RunValidate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var req CreateOrganizationRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	org, code, err := ctrl.orgUseCase.CreateOrganization(c.Request.Context(), claims.UserID, organization_use_case.CreateOrganizationRequest{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Description: req.Description,
		Address:     req.Address,
		Logo:        req.Logo,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "organization created successfully",
		Data:    ctrl.toOrgResponse(org),
	})
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Updates an existing organization's information
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Param request body UpdateOrganizationRequest true "Organization update data"
// @Success 200 {object} gin_utils.DataResponse{data=Organization}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id} [put]
func (ctrl *organizationController) UpdateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(100),
			Null: true,
		}).
		SetRule("description", map_validator.Rules{
			Type: reflect.String,
			Null: true,
		}).
		SetRule("address", map_validator.Rules{
			Type: reflect.String,
			Null: true,
		})

	jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
	jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	jsonData, err := jsonDataValidate.RunValidate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var req UpdateOrganizationRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	org, code, err := ctrl.orgUseCase.UpdateOrganization(c.Request.Context(), id, organization_use_case.UpdateOrganizationRequest{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Logo:        req.Logo,
		Settings:    req.Settings,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "organization updated successfully",
		Data:    ctrl.toOrgResponse(org),
	})
}

// DeleteOrganization godoc
// @Summary Delete organization
// @Description Soft deletes an organization by its ID
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id} [delete]
func (ctrl *organizationController) DeleteOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	code, err := ctrl.orgUseCase.DeleteOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "organization deleted successfully"})
}

// GetMembers godoc
// @Summary Get organization members
// @Description Retrieves all members of an organization with pagination
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]OrganizationMember}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id}/members [get]
func (ctrl *organizationController) GetMembers(c *gin.Context) {
	orgIdParam := c.Param("id")
	orgId, err := uuid.Parse(orgIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	members, code, err := ctrl.orgUseCase.GetMembers(c.Request.Context(), orgId, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]OrganizationMember, len(members))
	for i, m := range members {
		result[i] = ctrl.toMemberResponse(m)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// AddMember godoc
// @Summary Add member to organization
// @Description Adds a user as a member to an organization with a specific role
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Param request body AddMemberRequest true "Member data"
// @Success 201 {object} gin_utils.DataResponse{data=OrganizationMember}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id}/members [post]
func (ctrl *organizationController) AddMember(c *gin.Context) {
	orgIdParam := c.Param("id")
	orgId, err := uuid.Parse(orgIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	authData, _ := c.Get("auth")

	roles := map_validator.BuildRoles().
		SetRule("user_id", map_validator.Rules{
			Type: reflect.String,
		}).
		SetRule("role_id", map_validator.Rules{
			Type: reflect.String,
		})

	jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
	jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	jsonData, err := jsonDataValidate.RunValidate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var req AddMemberRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var invitedBy *uuid.UUID
	if claims, ok := authData.(*auth_utils.AuthClaim); ok {
		invitedBy = &claims.UserID
	}

	member, code, err := ctrl.orgUseCase.AddMember(c.Request.Context(), orgId, organization_use_case.AddMemberRequest{
		UserId:    req.UserId,
		RoleId:    req.RoleId,
		InvitedBy: invitedBy,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "member added successfully",
		Data:    ctrl.toMemberResponse(member),
	})
}

// UpdateMember godoc
// @Summary Update organization member
// @Description Updates a member's role or active status in an organization
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Param userId path string true "User ID (UUID)"
// @Param request body UpdateMemberRequest true "Member update data"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id}/members/{userId} [put]
func (ctrl *organizationController) UpdateMember(c *gin.Context) {
	orgIdParam := c.Param("id")
	orgId, err := uuid.Parse(orgIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	userIdParam := c.Param("userId")
	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid user id"})
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("role_id", map_validator.Rules{
			Type: reflect.String,
			Null: true,
		}).
		SetRule("is_active", map_validator.Rules{
			Type: reflect.Bool,
			Null: true,
		})

	jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
	jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	jsonData, err := jsonDataValidate.RunValidate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var req UpdateMemberRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	code, err := ctrl.orgUseCase.UpdateMember(c.Request.Context(), orgId, userId, organization_use_case.UpdateMemberRequest{
		RoleId:   req.RoleId,
		IsActive: req.IsActive,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "member updated successfully"})
}

// RemoveMember godoc
// @Summary Remove member from organization
// @Description Removes a user from an organization
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID (UUID)"
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/organizations/{id}/members/{userId} [delete]
func (ctrl *organizationController) RemoveMember(c *gin.Context) {
	orgIdParam := c.Param("id")
	orgId, err := uuid.Parse(orgIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization id"})
		return
	}

	userIdParam := c.Param("userId")
	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid user id"})
		return
	}

	code, err := ctrl.orgUseCase.RemoveMember(c.Request.Context(), orgId, userId)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "member removed successfully"})
}
