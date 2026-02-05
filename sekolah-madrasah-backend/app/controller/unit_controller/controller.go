package unit_controller

import (
	"net/http"
	"reflect"

	"sekolah-madrasah/app/use_case/unit_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type unitController struct {
	unitUseCase unit_use_case.UnitUseCase
}

func NewUnitController(unitUseCase unit_use_case.UnitUseCase) GinUnitController {
	return &unitController{unitUseCase: unitUseCase}
}

// GetUnit godoc
// @Summary Get perumahan by ID
// @Description Retrieves a single perumahan by its UUID
// @Tags Unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Unit ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=Unit}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{id} [get]
func (ctrl *unitController) GetUnit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	perumahan, code, err := ctrl.unitUseCase.GetUnit(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    toUnitResponse(perumahan),
	})
}

// GetUnits godoc
// @Summary List all perumahans
// @Description Retrieves a paginated list of perumahans
// @Tags Unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param type query string false "Filter by perumahan type"
// @Param organization_id query string false "Filter by organization ID"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]Unit}
// @Router /api/v1/companies [get]
func (ctrl *unitController) GetUnits(c *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := unit_use_case.UnitFilter{}
	if perumahanType := c.Query("type"); perumahanType != "" {
		filter.Type = &perumahanType
	}
	if orgId := c.Query("organization_id"); orgId != "" {
		if parsedId, err := uuid.Parse(orgId); err == nil {
			filter.OrganizationId = &parsedId
		}
	}

	perumahans, code, err := ctrl.unitUseCase.GetUnits(c.Request.Context(), filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]Unit, len(perumahans))
	for i, s := range perumahans {
		result[i] = toUnitResponse(s)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// CreateUnitRequest with organization_id
type CreateUnitRequestWithOrg struct {
	OrganizationId string `json:"organization_id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Email          string `json:"email,omitempty"`
	Logo           string `json:"logo,omitempty"`
}

// CreateUnit godoc
// @Summary Create a new perumahan
// @Description Creates a new perumahan under an organization
// @Tags Unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUnitRequestWithOrg true "Unit creation data"
// @Success 201 {object} gin_utils.DataResponse{data=Unit}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Router /api/v1/companies [post]
func (ctrl *unitController) CreateUnit(c *gin.Context) {
	var req CreateUnitRequestWithOrg
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	if req.Name == "" || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "name and code are required"})
		return
	}

	if req.OrganizationId == "" {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "organization_id is required"})
		return
	}

	orgId, err := uuid.Parse(req.OrganizationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid organization_id"})
		return
	}

	perumahan, code, err := ctrl.unitUseCase.CreateUnit(c.Request.Context(), orgId, unit_use_case.CreateUnitRequest{
		Name:    req.Name,
		Code:    req.Code,
		Type:    req.Type,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
		Logo:    req.Logo,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "perumahan created successfully",
		Data:    toUnitResponse(perumahan),
	})
}

// UpdateUnit godoc
// @Summary Update perumahan
// @Description Updates an existing perumahan
// @Tags Unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Unit ID (UUID)"
// @Param request body UpdateUnitRequest true "Unit update data"
// @Success 200 {object} gin_utils.DataResponse{data=Unit}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{id} [put]
func (ctrl *unitController) UpdateUnit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("address", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("phone", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("email", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("logo", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("settings", map_validator.Rules{Type: reflect.String, Null: true})

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

	var req UpdateUnitRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	perumahan, code, err := ctrl.unitUseCase.UpdateUnit(c.Request.Context(), id, unit_use_case.UpdateUnitRequest{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		Email:    req.Email,
		Logo:     req.Logo,
		Settings: req.Settings,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "perumahan updated successfully",
		Data:    toUnitResponse(perumahan),
	})
}

// DeleteUnit godoc
// @Summary Delete perumahan
// @Description Soft deletes a perumahan
// @Tags Unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Unit ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{id} [delete]
func (ctrl *unitController) DeleteUnit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	code, err := ctrl.unitUseCase.DeleteUnit(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "perumahan deleted successfully"})
}
