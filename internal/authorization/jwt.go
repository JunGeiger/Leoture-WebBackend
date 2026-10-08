package authorization

import (
	"LeotureWeb/internal/config"
	"LeotureWeb/internal/utils"
	"crypto"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 加密算法白名单
const (
	AlgHS256 = "HS256" // HMAC-SHA256 对称签名
	AlgES256 = "ES256" // ECDSA P-256 非对称签名
	AlgEdDSA = "EdDSA" // Ed25519 非对称签名
)

// validAlgorithms 签名算法白名单，解析 Token 时仅接受这些算法
var allowedAlgorithms = []string{AlgHS256, AlgES256, AlgEdDSA}

var jwtService *JWTService

// JWTService JWT服务，封装签发与解析逻辑
type JWTService struct {
	issuer        string            // 令牌签发者
	audience      jwt.ClaimStrings  // 令牌接收着
	algorithm     string            // 签名算法白名单，仅支持: HS256 / ES256 / EdDSA
	signingMethod jwt.SigningMethod // 签名算法，仅支持: HS256 / ES256 / EdDSA
	signingKey    any               // HS256: []byte; ES256: *ecdsa.PrivateKey; EdDSA: ed25519.PrivateKey
	verifyKey     any               // HS256: []byte; ES256: *ecdsa.PublicKey;  EdDSA: ed25519.PublicKey
	accessExpiry  time.Duration     // 访问令牌过期时间
	refreshExpiry time.Duration     // 刷新令牌过期时间 (默认7天)
}

func SetupJWT(cfg config.JWT) (*JWTService, error) {
	jwtService = &JWTService{}

	// 校验算法是否在白名单内
	if !isAllowedAlgorithm(cfg.Algorithm) {
		return nil, fmt.Errorf("jwt: unsupported algorithm: %s, allowed: %v", cfg.Algorithm, allowedAlgorithms)
	}
	jwtService.issuer = cfg.Issuer
	jwtService.audience = cfg.Audience
	jwtService.algorithm = cfg.Algorithm

	// 解析令牌过期时间
	jwtService.accessExpiry = cfg.AccessTokenExpire
	if jwtService.accessExpiry <= 0 {
		return nil, fmt.Errorf("jwt: access_token_ttl must be positive")
	}
	jwtService.refreshExpiry = cfg.RefreshTokenExpire
	if jwtService.refreshExpiry <= 0 {
		return nil, fmt.Errorf("jwt: refresh_token_ttl must be positive")
	}

	// 根据算法加载对应的签名/验签密钥
	switch cfg.Algorithm {
	case AlgHS256:
		jwtService.signingMethod = jwt.SigningMethodHS256
		if len(cfg.Secret) < 32 {
			return nil, fmt.Errorf("jwt: HS256 secret must be at least 32 bytes")
		}
		key := []byte(cfg.Secret)
		jwtService.signingKey = key
		jwtService.verifyKey = key // 对称算法，签名和验签使用同一密钥

	case AlgES256:
		jwtService.signingMethod = jwt.SigningMethodES256
		privateKeyPEM, publicKeyPEM, err := loadECDSAKeys(cfg.PrivateKey, cfg.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("jwt: ES256 keys failed to load: %w", err)
		}
		jwtService.signingKey = privateKeyPEM
		jwtService.verifyKey = publicKeyPEM

	case AlgEdDSA:
		jwtService.signingMethod = jwt.SigningMethodEdDSA
		privateKeyPEM, publicKeyPEM, err := loadEdDSAKeys(cfg.PrivateKey, cfg.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("jwt: EdDSA keys failed to load: %w", err)
		}
		jwtService.signingKey = privateKeyPEM
		jwtService.verifyKey = publicKeyPEM
	}

	return jwtService, nil
}

// isAllowedAlgorithm 检查算法是否在白名单中
func isAllowedAlgorithm(alg string) bool {
	for _, a := range allowedAlgorithms {
		if a == alg {
			return true
		}
	}
	return false
}

// loadECDSAKeys 加载ES256密钥对
func loadECDSAKeys(privateKeyPEM string, publicKeyPEM string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKeyStr, err := resolvePEMContent(privateKeyPEM, "private_key")
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: parse ec private key failed: %w", err)
	}

	publicKeyStr, err := resolvePEMContent(publicKeyPEM, "public_key")
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: parse ec public key failed: %w", err)
	}

	return privateKey, publicKey, nil
}

// loadEdDSAKeys 加载EdDSA密钥对
func loadEdDSAKeys(privateKeyPEM string, publicKeyPEM string) (crypto.PrivateKey, crypto.PublicKey, error) {
	privateKeyStr, err := resolvePEMContent(privateKeyPEM, "private_key")
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM(privateKeyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: parse ed private key failed: %w", err)
	}

	publicKeyStr, err := resolvePEMContent(publicKeyPEM, "public_key")
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := jwt.ParseEdPublicKeyFromPEM(publicKeyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: parse ed public key failed: %w", err)
	}

	return privateKey, publicKey, nil
}

// resolvePEMContent 统一解析PEM内容
// name 仅用于错误提示，如 "private_key"、"public_key"
func resolvePEMContent(content, keyName string) ([]byte, error) {
	// 优先使用内联字符串（去除首尾空白，避免YAML多行语法引入多余换行）
	if trimmed := strings.TrimSpace(content); trimmed != "" {
		return []byte(trimmed), nil
	}
	return nil, fmt.Errorf("jwt: secret key content is empty, key name: %s", keyName)
}

// ParseToken 解析并验证Token，强制使用白名单算法防止算法混淆攻击
func (s *JWTService) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (any, error) {
			// 校验签名方法类型是否匹配（双保险，与 WithValidMethods 配合）
			if token.Method.Alg() != s.algorithm {
				return nil, fmt.Errorf("jwt: unexpected signing method: %v", token.Method.Alg())
			}
			return s.verifyKey, nil
		},
		// 【安全关键】白名单校验：只允许指定的算法通过验证
		jwt.WithValidMethods(allowedAlgorithms),
		// 校验Issuer
		jwt.WithIssuer(s.issuer),
		// 校验Audience
		jwt.WithAudience(s.audience...),
	)

	if err != nil {
		return nil, err
	}

	// 类型断言提取 Claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("jwt: token claims invalid")
	}

	return claims, nil
}

// GenerateAccessToken 签发 Access Token
func (s *JWTService) GenerateAccessToken(userID string, sessionID string, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessExpiry)
	claims := JWTClaims{
		Username:  username,
		Type:      AccessToken,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        utils.GenUUIDv7().String(),
			Issuer:    s.issuer,
			Audience:  s.audience,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token, err := s.signToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GenerateRefreshToken 签发 Refresh Token
func (s *JWTService) GenerateRefreshToken(userID string, sessionID string) (string, uuid.UUID, error) {
	id := utils.GenUUIDv7()
	now := time.Now()
	claims := JWTClaims{
		Type:      RefreshToken,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Issuer:    s.issuer,
			Audience:  s.audience,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
		},
	}
	token, err := s.signToken(claims)
	if err != nil {
		return "", uuid.Nil, err
	}
	return token, id, nil
}

// signToken 使用配置的签名方法和密钥对 claims 进行签名，返回 token 字符串
func (s *JWTService) signToken(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(s.signingMethod, claims)
	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", fmt.Errorf("jwt: token signature failed: %w", err)
	}
	return tokenString, nil
}
