package unit_member_controller

import (
	"net/http"

	"sekolah-madrasah/app/use_case/unit_member_use_case"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type perumahanMemberController struct {
	memberUseCase unit_member_use_case.UnitMemberUseCase
}

func NewUnitMemberController(memberUseCase unit_member_use_case.UnitMemberUseCase) GinUnitMemberController {
	return &perumahanMemberController{memberUseCase: memberUseCase}
}

// GetMembers godoc
// @Summary List perumahan members
// @Description Retrieves a paginated list of members for a perumahan
// @Tags UnitMember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path string true "Perumahan ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role (admin, pengurus, warga, parent, staff)"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]UnitMember}
// @Router /api/v1/companies/{unit_id}/members [get]
func (ctrl *perumahanMemberController) GetMembers(c *gin.Context) {
	perumahanIdStr := c.Param("id")
	perumahanId, err := uuid.Parse(perumahanIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
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

	filter := unit_member_use_case.UnitMemberFilter{}
	if roleStr := c.Query("role"); roleStr != "" {
		role := schemas.UnitMemberRole(roleStr)
		filter.Role = &role
	}

	members, code, err := ctrl.memberUseCase.GetMembers(c.Request.Context(), perumahanId, filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]UnitMember, len(members))
	for i, m := range members {
		result[i] = toMemberResponse(m)
	}

	c.JSON(http.StatusOK, gin_utils.MakeDataPaginateResponse(gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	}))
}

// GetMember godoc
// @Summary Get perumahan member by ID
// @Description Retrieves a single perumahan member
// @Tags UnitMember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path string true "Perumahan ID (UUID)"
// @Param member_id path string true "Member ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=UnitMember}
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{unit_id}/members/{member_id} [get]
func (ctrl *perumahanMemberController) GetMember(c *gin.Context) {
	perumahanIdStr := c.Param("id")
	perumahanId, err := uuid.Parse(perumahanIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	memberIdStr := c.Param("memberId")
	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid member id"})
		return
	}

	member, code, err := ctrl.memberUseCase.GetMember(c.Request.Context(), perumahanId, memberId)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "success",
		Data:    toMemberResponse(member),
	})
}

// AddMember godoc
// @Summary Add member to perumahan
// @Description Adds a user as a member of the perumahan with a specified role
// @Tags UnitMember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path string true "Perumahan ID (UUID)"
// @Param body body AddMemberRequest true "Member data"
// @Success 201 {object} gin_utils.DataResponse{data=UnitMember}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{unit_id}/members [post]
func (ctrl *perumahanMemberController) AddMember(c *gin.Context) {
	perumahanIdStr := c.Param("id")
	perumahanId, err := uuid.Parse(perumahanIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	var body AddMemberRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid request body"})
		return
	}

	if body.UserId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "user_id is required"})
		return
	}

	if body.Role == "" {
		body.Role = schemas.UnitMemberRoleStaff
	}

	req := unit_member_use_case.AddMemberRequest{
		UserId: body.UserId,
		Role:   body.Role,
	}

	// Get current user ID from context if available
	var invitedBy *uuid.UUID
	if userIdVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIdVal.(uuid.UUID); ok {
			invitedBy = &uid
		}
	}

	member, code, err := ctrl.memberUseCase.AddMember(c.Request.Context(), perumahanId, req, invitedBy)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin_utils.DataResponse{
		Message: "member added successfully",
		Data:    toMemberResponse(member),
	})
}

// UpdateMember godoc
// @Summary Update perumahan member
// @Description Updates a perumahan member's role or status
// @Tags UnitMember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path string true "Perumahan ID (UUID)"
// @Param member_id path string true "Member ID (UUID)"
// @Param body body UpdateMemberRequest true "Update data"
// @Success 200 {object} gin_utils.DataResponse{data=UnitMember}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{unit_id}/members/{member_id} [put]
func (ctrl *perumahanMemberController) UpdateMember(c *gin.Context) {
	perumahanIdStr := c.Param("id")
	perumahanId, err := uuid.Parse(perumahanIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	memberIdStr := c.Param("memberId")
	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid member id"})
		return
	}

	var body UpdateMemberRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid request body"})
		return
	}

	req := unit_member_use_case.UpdateMemberRequest{
		Role:     body.Role,
		IsActive: body.IsActive,
	}

	member, code, err := ctrl.memberUseCase.UpdateMember(c.Request.Context(), perumahanId, memberId, req)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "member updated successfully",
		Data:    toMemberResponse(member),
	})
}

// RemoveMember godoc
// @Summary Remove member from perumahan
// @Description Removes a user from a perumahan
// @Tags UnitMember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path string true "Perumahan ID (UUID)"
// @Param member_id path string true "Member ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Router /api/v1/companies/{unit_id}/members/{member_id} [delete]
func (ctrl *perumahanMemberController) RemoveMember(c *gin.Context) {
	perumahanIdStr := c.Param("id")
	perumahanId, err := uuid.Parse(perumahanIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid perumahan id"})
		return
	}

	memberIdStr := c.Param("memberId")
	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid member id"})
		return
	}

	code, err := ctrl.memberUseCase.RemoveMember(c.Request.Context(), perumahanId, memberId)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "member removed successfully"})
}
