package router

import (
	"admin-api/handler"
	"admin-api/middleware"

	"github.com/labstack/echo/v4"
)

func VersionOne(v1 *echo.Group) {
	adminGroup := v1.Group("/admins")
	{
		adminHandler := handler.NewAdminHandler()
		adminGroup.POST("", adminHandler.CreateAdminUser, middleware.CheckApiKey)
		adminGroup.POST("/sign-in", adminHandler.SignIn)
		adminGroup.POST("/access-token", adminHandler.RefreshToken)
		adminGroup.GET("", adminHandler.GetAllAdminEmail, middleware.AuthenticateToken)
	}

}
