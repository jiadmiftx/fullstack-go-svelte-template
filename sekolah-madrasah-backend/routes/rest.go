package routes

import (
	"fmt"
	"log"

	"sekolah-madrasah/app/controller/auth_controller"
	"sekolah-madrasah/app/controller/organization_controller"
	"sekolah-madrasah/app/controller/permission_controller"
	"sekolah-madrasah/app/controller/post_controller"
	"sekolah-madrasah/app/controller/role_controller"
	"sekolah-madrasah/app/controller/unit_controller"
	"sekolah-madrasah/app/controller/unit_member_controller"
	"sekolah-madrasah/app/controller/unit_settings_controller"
	"sekolah-madrasah/app/controller/user_controller"
	"sekolah-madrasah/app/repository/org_member_repository"
	"sekolah-madrasah/app/repository/organization_repository"
	"sekolah-madrasah/app/repository/permission_repository"
	"sekolah-madrasah/app/repository/post_repository"
	"sekolah-madrasah/app/repository/role_repository"
	"sekolah-madrasah/app/repository/unit_member_repository"
	"sekolah-madrasah/app/repository/unit_repository"
	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/app/use_case/auth_use_case"
	"sekolah-madrasah/app/use_case/organization_use_case"
	"sekolah-madrasah/app/use_case/permission_use_case"
	"sekolah-madrasah/app/use_case/post_use_case"
	"sekolah-madrasah/app/use_case/role_use_case"
	"sekolah-madrasah/app/use_case/unit_member_use_case"
	"sekolah-madrasah/app/use_case/unit_use_case"
	"sekolah-madrasah/app/use_case/user_use_case"
	"sekolah-madrasah/config"
	"sekolah-madrasah/database"
	"sekolah-madrasah/pkg/http_middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Container struct {
	AuthController         auth_controller.AuthController
	UserController         user_controller.UserController
	RoleController         role_controller.RoleController
	PermissionController   permission_controller.PermissionController
	OrganizationController organization_controller.OrganizationController
	UnitController         unit_controller.GinUnitController
	UnitMemberController   unit_member_controller.GinUnitMemberController
	UnitSettingsController unit_settings_controller.GinUnitSettingsController
	PostController         *post_controller.PostController
}

func NewContainer(db *gorm.DB) *Container {
	userRepo := user_repository.NewUserRepository(db)
	roleRepo := role_repository.NewRoleRepository(db)
	permissionRepo := permission_repository.NewPermissionRepository(db)
	orgRepo := organization_repository.NewOrganizationRepository(db)
	orgMemberRepo := org_member_repository.NewOrgMemberRepository(db)
	unitRepo := unit_repository.NewUnitRepository(db)
	unitMemberRepo := unit_member_repository.NewUnitMemberRepository(db)
	postRepo := post_repository.NewPostRepository(db)

	authUseCase := auth_use_case.NewAuthUseCase(userRepo)
	userUseCase := user_use_case.NewUserUseCase(userRepo)
	roleUseCase := role_use_case.NewRoleUseCase(roleRepo, permissionRepo)
	permissionUseCase := permission_use_case.NewPermissionUseCase(permissionRepo)
	orgUseCase := organization_use_case.NewOrganizationUseCase(orgRepo, orgMemberRepo)
	unitUseCase := unit_use_case.NewUnitUseCase(unitRepo)
	unitMemberUseCase := unit_member_use_case.NewUnitMemberUseCase(unitMemberRepo)
	postUseCase := post_use_case.NewPostUseCase(postRepo, userRepo)

	authController := auth_controller.NewAuthController(authUseCase)
	userController := user_controller.NewUserController(userUseCase, nil)
	roleController := role_controller.NewRoleController(roleUseCase)
	permissionController := permission_controller.NewPermissionController(permissionUseCase)
	orgController := organization_controller.NewOrganizationController(orgUseCase)
	unitController := unit_controller.NewUnitController(unitUseCase)
	unitMemberController := unit_member_controller.NewUnitMemberController(unitMemberUseCase)
	unitSettingsController := unit_settings_controller.NewUnitSettingsController(db)
	postCtrl := post_controller.NewPostController(postUseCase)

	return &Container{
		AuthController:         authController,
		UserController:         userController,
		RoleController:         roleController,
		PermissionController:   permissionController,
		OrganizationController: orgController,
		UnitController:         unitController,
		UnitMemberController:   unitMemberController,
		UnitSettingsController: unitSettingsController,
		PostController:         postCtrl,
	}
}

func InitRoutes() *gin.Engine {
	mainDB, err := database.NewDBConnection(config.APP.MainDB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to main database: %v", err)
	}
	log.Println("✅ Main database connected")

	router := gin.Default()

	router.Use(http_middleware.CORS)

	container := NewContainer(mainDB)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/login", container.AuthController.Login)
			auth.POST("/register", container.AuthController.Register)
			auth.POST("/refresh", container.AuthController.RefreshToken)
		}

		users := v1.Group("/users")
		users.Use(http_middleware.JWTAuthentication)
		{
			users.GET("", container.UserController.GetUsers)
			users.GET("/me", container.UserController.GetCurrentUser)
			users.GET("/:id", container.UserController.GetUser)
			users.POST("", container.UserController.CreateUser)
			users.PUT("/:id", container.UserController.UpdateUser)
			users.DELETE("/:id", container.UserController.DeleteUser)
		}

		roles := v1.Group("/roles")
		roles.Use(http_middleware.JWTAuthentication)
		{
			roles.GET("", container.RoleController.GetRoles)
			roles.GET("/:id", container.RoleController.GetRole)
			roles.POST("", container.RoleController.CreateRole)
			roles.PUT("/:id", container.RoleController.UpdateRole)
			roles.DELETE("/:id", container.RoleController.DeleteRole)
		}

		permissions := v1.Group("/permissions")
		permissions.Use(http_middleware.JWTAuthentication)
		{
			permissions.GET("", container.PermissionController.GetPermissions)
			permissions.GET("/:id", container.PermissionController.GetPermission)
			permissions.POST("", container.PermissionController.CreatePermission)
			permissions.DELETE("/:id", container.PermissionController.DeletePermission)
		}

		organizations := v1.Group("/organizations")
		organizations.Use(http_middleware.JWTAuthentication)
		{
			organizations.GET("", container.OrganizationController.GetOrganizations)
			organizations.GET("/:id", container.OrganizationController.GetOrganization)
			organizations.POST("", container.OrganizationController.CreateOrganization)
			organizations.PUT("/:id", container.OrganizationController.UpdateOrganization)
			organizations.DELETE("/:id", container.OrganizationController.DeleteOrganization)

			organizations.GET("/:id/members", container.OrganizationController.GetMembers)
			organizations.POST("/:id/members", container.OrganizationController.AddMember)
			organizations.PUT("/:id/members/:userId", container.OrganizationController.UpdateMember)
			organizations.DELETE("/:id/members/:userId", container.OrganizationController.RemoveMember)
		}

		units := v1.Group("/units")
		units.Use(http_middleware.JWTAuthentication)
		{
			units.GET("", container.UnitController.GetUnits)
			units.GET("/:id", container.UnitController.GetUnit)
			units.POST("", container.UnitController.CreateUnit)
			units.PUT("/:id", container.UnitController.UpdateUnit)
			units.DELETE("/:id", container.UnitController.DeleteUnit)

			units.GET("/:id/members", container.UnitMemberController.GetMembers)
			units.GET("/:id/members/:memberId", container.UnitMemberController.GetMember)
			units.POST("/:id/members", container.UnitMemberController.AddMember)
			units.PUT("/:id/members/:memberId", container.UnitMemberController.UpdateMember)
			units.DELETE("/:id/members/:memberId", container.UnitMemberController.RemoveMember)

			units.GET("/:id/settings", container.UnitSettingsController.GetSettings)
			units.PUT("/:id/settings", container.UnitSettingsController.UpdateSettings)
		}

		posts := v1.Group("/posts")
		posts.Use(http_middleware.JWTAuthentication)
		{
			posts.GET("", container.PostController.GetPosts)
			posts.GET("/:id", container.PostController.GetPost)
			posts.POST("", container.PostController.CreatePost)
			posts.PUT("/:id", container.PostController.UpdatePost)
			posts.DELETE("/:id", container.PostController.DeletePost)
			posts.GET("/:id/comments", container.PostController.GetComments)
			posts.POST("/:id/comments", container.PostController.CreateComment)
			posts.DELETE("/:id/comments/:commentId", container.PostController.DeleteComment)
			posts.POST("/:id/vote", container.PostController.VotePoll)
		}
	}

	log.Println("✅ All routes configured")
	return router
}

func StartServer(engine *gin.Engine) {
	address := fmt.Sprintf("0.0.0.0:%d", config.APP.Rest.Port)
	log.Printf("🚀 Server starting on %s", address)

	if err := engine.Run(address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
