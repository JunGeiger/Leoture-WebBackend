package request

import (
	"LeotureWeb/internal/authorization"
	"time"

	"github.com/gin-gonic/gin"
)

// Gin请求上下文中携带的参数key值
const (
	RequestIDKey         = "request_id"   // 请求ID在gin上下文中的key值
	RequestTimeKey       = "request_time" // 中间件注入的请求时间key值
	RequestIDHeaderKey   = "x-request-id" // 前端携带的请求ID在请求头中的key值
	ContextKeyUserClaims = "jwt_claims"   // Gin上下文中存储Claims的key
)

// 默认分页参数
const (
	defaultPageSize = 10  // 默认分页条数
	maxPageSize     = 100 // 最大分页条数
)

// Header HTTP请求需要绑定的请求头参数
type Header struct {
	Authorization string `header:"Authorization" binding:"required"`
}

// DataID HTTP请求业务通用ID数据
type DataID struct {
	ID string `json:"id" form:"id" uri:"id" binding:"required"`
}

// DateRange HTTP请求业务数据通用日期范围
type DateRange struct {
	DateStart *time.Time `json:"dateStart" form:"dateStart" binding:"omitempty" time_format:"2006-01-02"`
	DateEnd   *time.Time `json:"dataEnd" form:"dataEnd" binding:"omitempty" time_format:"2006-01-02"`
}

// PageParams HTTP请求业务数据通用分页参数
type PageParams struct {
	CurrentPage int `json:"currentPage" form:"currentPage" binding:"omitempty"` // 当前页码（从 1 开始）
	PageSize    int `json:"pageSize" form:"pageSize" binding:"omitempty"`       // 每页条数
}

// InitPagination 计算分页数据
func (p *PageParams) InitPagination() {
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	} else if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

// GetClaims 从 gin.Context 中提取 JWT Claims。
// 在 Handler 层中使用，获取当前登录用户信息。
//
// 使用方式：
//
//	claims := middleware.GetClaims(c)
//	userID := claims.UserID
func GetClaims(c *gin.Context) *authorization.JWTClaims {
	val, exists := c.Get(ContextKeyUserClaims)
	if !exists {
		return nil
	}
	claims, ok := val.(*authorization.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
