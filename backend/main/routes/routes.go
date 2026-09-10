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
	}
}
