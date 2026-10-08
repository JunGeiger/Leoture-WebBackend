package router

import (
	"LeotureWeb/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(r *gin.RouterGroup, h *handler.Role) {
	g := r.Group("/role")
	{
		g.POST("/create", h.Create)
		g.POST("/update", h.Update)
		g.DELETE("/delete/:id", h.Delete)
		g.GET("/:id", h.GetByID)
		g.GET("/code/:code", h.GetByCode)
		g.GET("/list", h.List)
		g.GET("/page", h.ListByPage)
	}
}
