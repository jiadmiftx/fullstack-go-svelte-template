package permission_controller

import (
	"net/http"
	"reflect"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sekolah-madrasah/app/use_case/permission_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"
)

type permissionController struct {
	permissionUseCase permission_use_case.PermissionUseCase
}

func NewPermissionController(permissionUseCase permission_use_case.PermissionUseCase) PermissionController {
	return &permissionController{permissionUseCase: permissionUseCase}
}

func (ctrl *permissionController) toPermissionResponse(p permission_use_case.Permission) Permission {
	return Permission{
		Id:          p.Id,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

// GetPermission godoc
// @Summary Get permission by ID
// @Description Retrieves a single permission by its UUID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=Permission}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/permissions/{id} [get]
func (ctrl *permissionController) GetPermission(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid permission id"})
		return
	}

	permission, code, err := ctrl.permissionUseCase.GetPermission(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    ctrl.toPermissionResponse(permission),
	})
}

// GetPermissions godoc
// @Summary List all permissions
// @Description Retrieves a paginated list of permissions with optional filtering
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param resource query string false "Filter by resource"
// @Param action query string false "Filter by action"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]Permission}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/permissions [get]
func (ctrl *permissionController) GetPermissions(c *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := permission_use_case.PermissionFilter{}

	if resource := c.Query("resource"); resource != "" {
		filter.Resource = &resource
	}
	if action := c.Query("action"); action != "" {
		filter.Action = &action
	}

	permissions, code, err := ctrl.permissionUseCase.GetPermissions(c.Request.Context(), filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]Permission, len(permissions))
	for i, p := range permissions {
		result[i] = ctrl.toPermissionResponse(p)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Creates a new permission with resource, action, and description
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePermissionRequest true "Permission creation data"
// @Success 201 {object} gin_utils.DataResponse{data=Permission}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/permissions [post]
func (ctrl *permissionController) CreatePermission(c *gin.Context) {
	roles := map_validator.BuildRoles().
		SetRule("resource", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(50),
		}).
		SetRule("action", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(20),
		}).
		SetRule("description", map_validator.Rules{
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

	var req CreatePermissionRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	permission, code, err := ctrl.permissionUseCase.CreatePermission(c.Request.Context(), permission_use_case.CreatePermissionRequest{
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "permission created successfully",
		Data:    ctrl.toPermissionResponse(permission),
	})
}

// DeletePermission godoc
// @Summary Delete permission
// @Description Deletes a permission by its ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/permissions/{id} [delete]
func (ctrl *permissionController) DeletePermission(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid permission id"})
		return
	}

	code, err := ctrl.permissionUseCase.DeletePermission(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "permission deleted successfully"})
}
