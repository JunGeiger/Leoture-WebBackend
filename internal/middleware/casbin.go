package middleware

import (
	"LeotureWeb/internal/authorization"
	errs "LeotureWeb/internal/errors"
	"LeotureWeb/internal/types/request"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Casbin Casbin RBAC 鉴权中间件
// 必须在 JWT 认证中间件之后使用（依赖 ctx 中的 claims 信息）
//
// 鉴权逻辑：
//  1. 从 JWT claims 中获取用户ID
//  2. 通过 CasbinService 查询用户角色并校验权限
//  3. 权限不足返回 403 Forbidden
func Casbin(c *gin.Context, casbinSvc *authorization.CasbinService) bool {
	// 从 gin.Context 获取 JWT claims（由 JWT 中间件设置）
	claims := request.GetClaims(c)
	if claims == nil {
		// 未登录或 claims 解析失败，由 JWT 中间件处理，此处兜底
		abortWithBizError(c, errs.ErrUnauthenticated)
		return false
	}

	// 构建鉴权参数
	sub := claims.Username    // 用户名称
	obj := c.Request.URL.Path // 请求路径
	act := c.Request.Method
	// 以用户维度进行权限校验（自动查询用户角色并逐一匹配）
	ok, err := casbinSvc.Enforcer.Enforce(sub, obj, act)

	// 获取所有匹配的策略
	if err != nil {
		slog.Error("casbin: enforced error",
			slog.String("username", sub),
			slog.String("path", obj),
			slog.String("method", act),
		)
		abortWithBizError(c, errs.ErrForbidden)
		return false
	}

	if !ok {
		slog.Warn("casbin: access denied",
			slog.String("roleCode", sub),
			slog.String("path", obj),
			slog.String("method", act),
		)
		// 返回 403 权限不足
		abortWithBizError(c, errs.ErrForbidden)
		return false
	}
	return true
}
