package auth_controller

import (
	"net/http"
	"reflect"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/gin-gonic/gin"
	"sekolah-madrasah/app/use_case/auth_use_case"
	"sekolah-madrasah/pkg/gin_utils"
)

type authController struct {
	authUseCase auth_use_case.AuthUseCase
}

func NewAuthController(authUseCase auth_use_case.AuthUseCase) AuthController {
	return &authController{authUseCase: authUseCase}
}

// Login godoc
// @Summary User login
// @Description Authenticates user credentials and returns JWT access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} gin_utils.DataResponse{data=LoginResponse}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/auth/login [post]
func (ctrl *authController) Login(c *gin.Context) {
	roles := map_validator.BuildRoles().
		SetRule("email", map_validator.Rules{
			Type: reflect.String,
			Max:  map_validator.SetTotal(255),
		}).
		SetRule("password", map_validator.Rules{
			Type: reflect.String,
			Min:  map_validator.SetTotal(1),
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

	var req LoginRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result, code, err := ctrl.authUseCase.Login(c.Request.Context(), auth_use_case.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "login successful",
		Data: LoginResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    result.ExpiresAt,
			User: UserInfo{
				Id:           result.User.Id,
				Email:        result.User.Email,
				FullName:     result.User.FullName,
				IsSuperAdmin: result.User.IsSuperAdmin,
				IsActive:     result.User.IsActive,
				LastLoginAt:  result.User.LastLoginAt,
			},
		},
	})
}

// Register godoc
// @Summary User registration
// @Description Creates a new user account with email, password, and full name
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} gin_utils.DataResponse{data=RegisterResponse}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 409 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/auth/register [post]
func (ctrl *authController) Register(c *gin.Context) {
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

	var req RegisterRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result, code, err := ctrl.authUseCase.Register(c.Request.Context(), auth_use_case.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "registration successful",
		Data: RegisterResponse{
			User: UserInfo{
				Id:           result.Id,
				Email:        result.Email,
				FullName:     result.FullName,
				IsSuperAdmin: result.IsSuperAdmin,
				IsActive:     result.IsActive,
			},
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates new access and refresh tokens using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} gin_utils.DataResponse{data=LoginResponse}
// @Failure 400 {object} gin_utils.MessageResponse
// @Failure 401 {object} gin_utils.MessageResponse
// @Failure 500 {object} gin_utils.MessageResponse
// @Router /api/v1/auth/refresh [post]
func (ctrl *authController) RefreshToken(c *gin.Context) {
	roles := map_validator.BuildRoles().
		SetRule("refresh_token", map_validator.Rules{
			Type: reflect.String,
			Min:  map_validator.SetTotal(1),
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

	var req RefreshTokenRequest
	if err := jsonData.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	result, code, err := ctrl.authUseCase.RefreshToken(c.Request.Context(), auth_use_case.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(code, gin_utils.DataResponse{
		Message: "token refreshed",
		Data: LoginResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    result.ExpiresAt,
			User: UserInfo{
				Id:           result.User.Id,
				Email:        result.User.Email,
				FullName:     result.User.FullName,
				IsSuperAdmin: result.User.IsSuperAdmin,
				IsActive:     result.User.IsActive,
				LastLoginAt:  result.User.LastLoginAt,
			},
		},
	})
}
