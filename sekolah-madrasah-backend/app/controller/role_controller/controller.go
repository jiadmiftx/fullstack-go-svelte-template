package role_controller

import (
	"log"
	"net/http"

	"sekolah-madrasah/app/use_case/role_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type roleController struct {
	roleUseCase role_use_case.RoleUseCase
}

func NewRoleController(roleUseCase role_use_case.RoleUseCase) RoleController {
	return &roleController{roleUseCase: roleUseCase}
}

func (ctrl *roleController) toRoleResponse(r role_use_case.Role) Role {
	perms := make([]Permission, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = Permission{
			Id:          p.Id,
			Name:        p.Name,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
		}
	}

	return Role{
		Id:             r.Id,
		OrganizationId: r.OrganizationId,
		Name:           r.Name,
		DisplayName:    r.DisplayName,
		Type:           r.Type,
		Level:          r.Level,
		Description:    r.Description,
		IsDefault:      r.IsDefault,
		Permissions:    perms,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

// GetRole godoc
// @Summary Get role by ID
// @Description Retrieves a single role by its UUID including assigned permissions
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=Role}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/roles/{id} [get]
func (ctrl *roleController) GetRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid role id"})
		return
	}

	role, code, err := ctrl.roleUseCase.GetRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    ctrl.toRoleResponse(role),
	})
}

// GetRoles godoc
// @Summary List all roles
// @Description Retrieves a paginated list of roles with optional filtering
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param name query string false "Filter by role name"
// @Param type query string false "Filter by role type"
// @Param organization_id query string false "Filter by organization ID"
// @Param is_global query boolean false "Filter global roles only"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]Role}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/roles [get]
func (ctrl *roleController) GetRoles(c *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := role_use_case.RoleFilter{}

	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}
	if roleType := c.Query("type"); roleType != "" {
		filter.Type = &roleType
	}
	if orgId := c.Query("organization_id"); orgId != "" {
		if parsedId, err := uuid.Parse(orgId); err == nil {
			filter.OrganizationId = &parsedId
		}
	}
	if isGlobal := c.Query("is_global"); isGlobal == "true" {
		global := true
		filter.IsGlobal = &global
	}

	roles, code, err := ctrl.roleUseCase.GetRoles(c.Request.Context(), filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]Role, len(roles))
	for i, r := range roles {
		result[i] = ctrl.toRoleResponse(r)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// CreateRole godoc
// @Summary Create a new role
// @Description Creates a new role with optional permission assignments
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRoleRequest true "Role creation data"
// @Success 201 {object} gin_utils.DataResponse{data=Role}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/roles [post]
func (ctrl *roleController) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	// Debug log
	log.Printf("CreateRole: Received permission_ids: %v (count: %d)", req.PermissionIds, len(req.PermissionIds))

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "name is required"})
		return
	}

	role, code, err := ctrl.roleUseCase.CreateRole(c.Request.Context(), role_use_case.CreateRoleRequest{
		OrganizationId: req.OrganizationId,
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Type:           req.Type,
		Level:          req.Level,
		Description:    req.Description,
		IsDefault:      req.IsDefault,
		PermissionIds:  req.PermissionIds,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "role created successfully",
		Data:    ctrl.toRoleResponse(role),
	})
}

// UpdateRole godoc
// @Summary Update role
// @Description Updates an existing role's information and permissions
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID (UUID)"
// @Param request body UpdateRoleRequest true "Role update data"
// @Success 200 {object} gin_utils.DataResponse{data=Role}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/roles/{id} [put]
func (ctrl *roleController) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid role id"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	// Debug log
	log.Printf("UpdateRole: Received permission_ids: %v (count: %d)", req.PermissionIds, len(req.PermissionIds))

	role, code, err := ctrl.roleUseCase.UpdateRole(c.Request.Context(), id, role_use_case.UpdateRoleRequest{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Level:         req.Level,
		PermissionIds: req.PermissionIds,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "role updated successfully",
		Data:    ctrl.toRoleResponse(role),
	})
}

// DeleteRole godoc
// @Summary Delete role
// @Description Soft deletes a role by its ID
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/roles/{id} [delete]
func (ctrl *roleController) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid role id"})
		return
	}

	code, err := ctrl.roleUseCase.DeleteRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "role deleted successfully"})
}
