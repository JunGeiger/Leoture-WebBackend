package router

import (
	"LeotureWeb/internal/authorization"
	"LeotureWeb/internal/config"
	"LeotureWeb/internal/handler"
	"LeotureWeb/internal/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupGinRouter 初始化 Gin Engine
func SetupGinRouter(cfg config.Config, handlerSet *handler.Set,
	jwtSvc *authorization.JWTService, casbinSvc *authorization.CasbinService) *gin.Engine {
	// 设置 Gin 运行模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局错误处理
	r.Use(middleware.ErrorHandler())
	// CORS中间件 - 处理跨域请求
	if cfg.Server.EnabledCors {
		r.Use(middleware.Cors(*cfg.Cors))
	}
	// 注入参数到请求
	r.Use(middleware.Trace())
	// Request日志处理
	r.Use(middleware.ReqLog())
	// Response日志处理
	r.Use(middleware.RespLog())
	r.Use(authMiddleware(jwtSvc, casbinSvc))

	apiV1 := r.Group("/api/v1")
	registerRoutes(apiV1, handlerSet)
	return r
}

// authMiddleware 使用权限中间件全量拦截，设置通行白名单
func authMiddleware(jwtSvc *authorization.JWTService, casbinSvc *authorization.CasbinService) gin.HandlerFunc {
	// public白名单
	publicPrefixes := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refreshToken",
	}
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		for _, prefix := range publicPrefixes {
			if strings.HasPrefix(p, prefix) {
				c.Next()
				return
			}
		}
		if !middleware.JWT(c, jwtSvc) {
			return
		}
		if !middleware.Casbin(c, casbinSvc) {
			return
		}
		c.Next()
	}
}

func registerRoutes(apiV1 *gin.RouterGroup, set *handler.Set) {
	registerAuthRoutes(apiV1, set.AuthHandler)
	registerUserRoutes(apiV1, set.UserHandler)
	RegisterRoleRoutes(apiV1, set.RoleHandler)
	registerMenuRoutes(apiV1, set.MenuHandler)
}
