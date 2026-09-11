package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"oss.kftd.co.id/v2/main/shared/middleware"
	"oss.kftd.co.id/v2/user-management/handler"
	"oss.kftd.co.id/v2/user-management/repository"
	"oss.kftd.co.id/v2/user-management/service"
)

func Setup(
	router *gin.Engine,
	db *gorm.DB,
) {

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(
		userRepo,
	)

	authHandler := handler.NewAuthHandler(
		authService,
	)

	// User
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Role
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	// Permission
	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// Role Permission
	rolePermissionRepo := repository.NewRolePermissionRepository(db)
	rolePermissionService := service.NewRolePermissionService(
		roleRepo,
		permissionRepo,
		rolePermissionRepo,
	)
	rolePermissionHandler := handler.NewRolePermissionHandler(rolePermissionService)

	// User Role
	userRoleRepo := repository.NewUserRoleRepository(db)
	userRoleService := service.NewUserRoleService(
		userRepo,
		roleRepo,
		userRoleRepo,
	)
	userRoleHandler := handler.NewUserRoleHandler(userRoleService)

	route := router.Group("/v1")
	{
		auth := route.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
			auth.POST("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		users := route.Group("/users")
		{
			users.Use(middleware.AuthMiddleware())

			users.GET("", middleware.RequirePermission(permissionService, "user.read"), userHandler.GetUsers)
			users.GET("/:uuid", middleware.RequirePermission(permissionService, "user.read"), userHandler.GetUser)
			users.POST("", middleware.RequirePermission(permissionService, "user.create"), userHandler.CreateUser)
			users.PUT("/:uuid", middleware.RequirePermission(permissionService, "user.update"), userHandler.UpdateUser)
			users.DELETE("/:uuid", middleware.RequirePermission(permissionService, "user.delete"), userHandler.DeleteUser)
			users.GET(
				"/:uuid/roles",
				middleware.RequirePermission(permissionService, "user.read"),
				userRoleHandler.GetUserRoles,
			)

			users.PUT(
				"/:uuid/roles",
				middleware.RequirePermission(permissionService, "user.update"),
				userRoleHandler.UpdateUserRoles,
			)

			roles := users.Group("/roles")
			{
				roles.Use(middleware.AuthMiddleware())

				roles.GET("", middleware.RequirePermission(permissionService, "role.read"), roleHandler.GetRoles)
				roles.GET("/:uuid", middleware.RequirePermission(permissionService, "role.read"), roleHandler.GetRole)
				roles.POST("", middleware.RequirePermission(permissionService, "role.create"), roleHandler.CreateRole)
				roles.PUT("/:uuid", middleware.RequirePermission(permissionService, "role.update"), roleHandler.UpdateRole)
				roles.DELETE("/:uuid", middleware.RequirePermission(permissionService, "role.delete"), roleHandler.DeleteRole)
				roles.PUT(
					"/:uuid/permissions",
					middleware.RequirePermission(permissionService, "role.update"),
					rolePermissionHandler.UpdateRolePermissions,
				)
				roles.GET(
					"/:uuid/permissions",
					middleware.RequirePermission(permissionService, "role.read"),
					rolePermissionHandler.GetRolePermissions,
				)
			}

			permissions := users.Group("/permissions")
			{
				permissions.Use(middleware.AuthMiddleware())

				permissions.GET("", middleware.RequirePermission(permissionService, "permission.read"), permissionHandler.GetPermissions)
				permissions.GET("/:uuid", middleware.RequirePermission(permissionService, "permission.read"), permissionHandler.GetPermission)
				permissions.POST("", middleware.RequirePermission(permissionService, "permission.create"), permissionHandler.CreatePermission)
				permissions.PUT("/:uuid", middleware.RequirePermission(permissionService, "permission.update"), permissionHandler.UpdatePermission)
				permissions.DELETE("/:uuid", middleware.RequirePermission(permissionService, "permission.delete"), permissionHandler.DeletePermission)
			}
		}
	}
}
