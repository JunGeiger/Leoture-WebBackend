package router

import (
	"LeotureWeb/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.RouterGroup, h *handler.User) {
	g := r.Group("/user")
	{
		g.POST("/create", h.Create)
		g.POST("/update", h.Update)
		g.DELETE("/delete/:id", h.Delete)
		g.GET("/:id", h.GetByID)
		g.GET("/username/:username", h.GetByUsername)
		g.GET("/list", h.List)
		g.GET("/page", h.ListByPage)
	}
}
