package router

import (
	"LeotureWeb/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerDictRoutes(r *gin.RouterGroup, h *handler.Dict) {
	g := r.Group("/dict")
	{
		g.POST("/create", h.Create)
		g.POST("/update", h.Update)
		g.DELETE("/delete/:id", h.Delete)
		g.GET("/:id", h.GetByID)
		g.GET("/key/:category/:key", h.GetByKey)
		g.GET("/list", h.List)
		g.GET("/page", h.ListByPage)
	}
}
