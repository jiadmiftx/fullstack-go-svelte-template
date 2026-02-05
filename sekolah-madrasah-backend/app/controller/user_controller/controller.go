package user_controller

import (
	"net/http"
	"reflect"

	"sekolah-madrasah/app/service/membership_service"
	"sekolah-madrasah/app/use_case/user_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userController struct {
	userUseCase       user_use_case.UserUseCase
	membershipService membership_service.MembershipService
}

func NewUserController(userUseCase user_use_case.UserUseCase, membershipService membership_service.MembershipService) UserController {
	return &userController{
		userUseCase:       userUseCase,
		membershipService: membershipService,
	}
}

func (ctrl *userController) toUserResponse(u user_use_case.User) User {
	return User{
		Id:              u.Id,
		Email:           u.Email,
		FullName:        u.FullName,
		Phone:           u.Phone,
		Avatar:          u.Avatar,
		IsSuperAdmin:    u.IsSuperAdmin,
		IsActive:        u.IsActive,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Returns the profile of the currently authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gin_utils.DataResponse{data=User}
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users/me [get]
func (ctrl *userController) GetCurrentUser(c *gin.Context) {
	userIdVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "unauthorized"})
		return
	}

	// Handle both uuid.UUID and string types
	var userId uuid.UUID
	switch v := userIdVal.(type) {
	case uuid.UUID:
		userId = v
	case string:
		var err error
		userId, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "invalid user id"})
			return
		}
	default:
		c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "invalid user id type"})
		return
	}

	user, code, err := ctrl.userUseCase.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "success",
		Data:    ctrl.toUserResponse(user),
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieves a single user by their UUID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} gin_utils.DataResponse{data=User}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users/{id} [get]
func (ctrl *userController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid user id"})
		return
	}

	user, code, err := ctrl.userUseCase.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    ctrl.toUserResponse(user),
	})
}

// GetUsers godoc
// @Summary List all users
// @Description Retrieves a paginated list of users with optional filtering
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param email query string false "Filter by email"
// @Param is_active query boolean false "Filter by active status"
// @Success 200 {object} gin_utils.DataWithPaginateResponse{data=[]User}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users [get]
func (ctrl *userController) GetUsers(c *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := user_use_case.UserFilter{}

	if email := c.Query("email"); email != "" {
		filter.Email = &email
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	if platformOnly := c.Query("platform_only"); platformOnly != "" {
		isPlatformOnly := platformOnly == "true"
		filter.PlatformOnly = &isPlatformOnly
	}

	users, code, err := ctrl.userUseCase.GetUsers(c.Request.Context(), filter, paginate)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result := make([]User, len(users))
	for i, u := range users {
		result[i] = ctrl.toUserResponse(u)
	}

	c.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "success",
			Data:    result,
		},
		Paginate: paginate,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with the provided details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} gin_utils.DataResponse{data=User}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users [post]
func (ctrl *userController) CreateUser(c *gin.Context) {
	roles := map_validator.BuildRoles().
		SetRule("email", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(255),
		}).
		SetRule("password", map_validator.Rules{
			Type: reflect.String,
			Min:  map_validator.SetTotal(8),
		}).
		SetRule("full_name", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(100),
		}).
		SetRule("phone", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(20),
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

	var req CreateUserRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	user, code, err := ctrl.userUseCase.CreateUser(c.Request.Context(), user_use_case.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "user created successfully",
		Data:    ctrl.toUserResponse(user),
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates an existing user's information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} gin_utils.DataResponse{data=User}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users/{id} [put]
func (ctrl *userController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid user id"})
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("email", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(255),
			Null: true,
		}).
		SetRule("full_name", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(100),
			Null: true,
		}).
		SetRule("phone", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(20),
			Null: true,
		}).
		SetRule("avatar", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(500),
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

	var req UpdateUserRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	user, code, err := ctrl.userUseCase.UpdateUser(c.Request.Context(), id, user_use_case.UpdateUserRequest{
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "user updated successfully",
		Data:    ctrl.toUserResponse(user),
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft deletes a user by their ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} gin_utils.MessageResponse
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 404 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users/{id} [delete]
func (ctrl *userController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid user id"})
		return
	}

	code, err := ctrl.userUseCase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.MessageResponse{Message: "user deleted successfully"})
}

// GetMyMemberships godoc
// @Summary Get current user's memberships
// @Description Returns the authenticated user's organization and perumahan memberships
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gin_utils.DataResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/users/me/memberships [get]
func (ctrl *userController) GetMyMemberships(c *gin.Context) {
	// Get user ID from JWT context (set by auth middleware)
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "unauthorized"})
		return
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		// Try parsing from string
		userIdStr, ok := userIdRaw.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "invalid user id in token"})
			return
		}
		var err error
		userId, err = uuid.Parse(userIdStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "invalid user id format"})
			return
		}
	}

	memberships, code, err := ctrl.membershipService.GetUserMemberships(c.Request.Context(), userId)
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "success",
		Data:    memberships,
	})
}
