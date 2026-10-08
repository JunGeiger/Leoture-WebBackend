package router

import (
	"LeotureWeb/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerMenuRoutes(r *gin.RouterGroup, h *handler.Menu) {
	g := r.Group("/menu")
	{
		g.POST("/create", h.Create)
		g.POST("/update", h.Update)
		g.DELETE("/delete/:id", h.Delete)
		g.GET("/:id", h.GetByID)
		g.GET("/list", h.List)
		g.GET("/page", h.ListByPage)
		g.GET("/tree", h.Tree)
	}
}
