package middleware

import (
	"LeotureWeb/internal/authorization"
	errs "LeotureWeb/internal/errors"
	"LeotureWeb/internal/types/request"
	"errors"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 返回 JWT 身份认证中间件。
// 从 Authorization 请求头中提取 Bearer Token，验证后将 Claims 存入 gin.Context。
//
// 使用方式：
//
//	router.Use(middleware.Auth(jwtManager))
//	router.GET("/profile", handler.Profile)
func JWT(c *gin.Context, jwtSvc *authorization.JWTService) bool {
	// 1. 从请求头解析Token
	var header request.Header
	err := c.ShouldBindHeader(&header)
	if err != nil {
		abortWithBizError(c, errs.ErrUnauthenticated.WithErr(err))
		return false
	}
	tokenString, err := extractBearerToken(header.Authorization)
	if err != nil {
		abortWithBizError(c, errs.ErrUnauthenticated.WithErr(err))
		return false
	}

	// 2. 解析Token
	claims, err := jwtSvc.ParseToken(tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			abortWithBizError(c, errs.ErrTokenExpired.WithErr(err))
		} else {
			abortWithBizError(c, errs.ErrUnauthenticated.WithErr(err))
		}
		return false
	}

	// 3. 验证声明有效性
	if err := jwtSvc.ValidateClaims(claims); err != nil {
		abortWithBizError(c, err)
		return false
	}

	// 添加请求日志
	slog.Info("jwt: User authenticated",
		slog.Any("userID", claims.Subject),
		slog.Any("username", claims.Username),
		slog.String("tokenType", claims.Type),
	)

	// 4. 将Claims注入Context，后续Handler可直接获取
	c.Set(request.ContextKeyUserClaims, claims)
	return true
}

// extractBearerToken 从 Authorization 请求头中提取 Bearer Token。
// 期望格式: "Bearer <token>"
func extractBearerToken(authHeader string) (string, error) {
	// 校验 Bearer 前缀（不区分大小写）
	const prefix = "Bearer "

	if len(authHeader) < len(prefix) || !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("invalid Authorization header format, expected Bearer token")
	}

	token := strings.TrimSpace(authHeader[len(prefix):])
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}
