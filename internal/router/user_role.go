package router

import (
	"LeotureWeb/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoleRoutes(r *gin.RouterGroup, h *handler.UserRole) {
	g := r.Group("/user-role")
	{
		g.POST("/update", h.Update)
		g.GET("/roles/:userID", h.GetRolesByUserID)
	}
}
