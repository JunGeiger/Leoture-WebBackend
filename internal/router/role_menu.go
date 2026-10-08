package router

import (
	"LeotureWeb/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoleMenuRoutes(r *gin.RouterGroup, h *handler.RoleMenu) {
	g := r.Group("/role-menu")
	{
		g.POST("/update", h.Update)
		g.GET("/roles/:roleID", h.GetMenusByRoleID)
	}
}
