package swagger

import (
	"LeotureWeb/internal/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 配置 Swagger 路由s
func SetupSwagger(r *gin.Engine, cfg config.Config) {
	// 根据配置决定是否启用，而且只在非生产环境启用 Swagger
	if !cfg.Swagger.Enabled || gin.Mode() == gin.ReleaseMode {
		return
	}

	SwaggerInfo.Title = cfg.Swagger.Title
	SwaggerInfo.Description = cfg.Swagger.Desc
	SwaggerInfo.Version = cfg.App.Version

	//
	// 方式一：使用默认配置
	r.GET(cfg.Swagger.Path+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 方式二：自定义配置（可选）
	// url := ginSwagger.URL("/swagger/doc.json") // 指定访问的 json 文件
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
