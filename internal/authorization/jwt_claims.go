package authorization

import (
	errs "LeotureWeb/internal/errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

// Token类型
const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

// Token状态
const (
	StatusValid   = 1 // 活跃
	StatusRevoked = 2 // 主动注销
	StatusKicked  = 3 // 已被替换(已有新登录)
	StatusLocked  = 4 // 风控/封禁
)

// JWTClaims 自定义 JWT 载荷
// 嵌入 jwt.RegisteredClaims 以支持标准字段
// iss 签发者 		推荐 	"leoture-web"
// sub 主体 			必须 	user_id
// aud 受众 			可选 	"leoture-api"
// exp 过期时间 		必须 	access 15m / Refresh 7d
// nbf 生效时间 		可选 	 基本不用
// iat 签发时间 		推荐 	 防重放
// jti Token唯一ID 	强烈推荐  黑名单 / 防复用
type JWTClaims struct {
	Type      string `json:"type"`     // "access" or "refresh"
	Username  string `json:"username"` // 用户名称
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

// ValidateClaims 验证JWT声明的有效性
func (s *JWTService) ValidateClaims(claims *JWTClaims) *errs.BizError {
	// 校验Token类型
	if claims.Type != AccessToken {
		slog.Error("jwt: Claim token type verification failed",
			slog.String("claims.TokenType", claims.Type))
		return errs.ErrUnauthenticated

	}

	// 检查必需字段
	if claims.ID == "" || claims.Subject == "" {
		slog.Error("jwt: Claim main field verification failed",
			slog.String("claims.ID", claims.ID),
			slog.String("claims.Subject", claims.Subject))
		return errs.ErrUnauthenticated
	}
	if claims.Username == "" {
		slog.Error("jwt: Claim business field verification failed",
			slog.Any("claims.Username", claims.Username))
		return errs.ErrUnauthenticated
	}

	// 可以添加更多业务验证逻辑，如检查用户是否被禁用等
	// 这里需要集成用户服务来查询用户状态

	return nil
}
