package router

import (
	"LeotureWeb/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.RouterGroup, h *handler.Auth) {
	r.POST("/auth/login", h.LoginByUser)
	r.POST("/auth/refreshToken", h.LoginByRefreshToken)
	r.GET("/auth/getUserInfo", h.GetUserInfo)
	r.GET("/route/getConstantRoutes", h.GetConstantRoutes)
	r.GET("/route/getUserRoutes", h.GetUserRoutes)
	r.GET("/route/isRouteExist", h.IsRouteExist)
}
