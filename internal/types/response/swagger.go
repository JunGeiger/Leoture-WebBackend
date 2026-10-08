package response

import "LeotureWeb/internal/errors"

// SwaggerResponse 仅用作Swagger文档注释读取
type SwaggerResponse struct {
	Code      errors.ResponseCode `json:"code"`                // 业务状态码
	Message   string              `json:"message"`             // 面向用户的提示信息
	Data      any                 `json:"data,omitempty"`      // 业务数据，错误时省略
	RequestID string              `json:"requestId,omitempty"` // 请求唯一标识
}

// SwaggerPaginatedData 仅用作Swagger文档注释读取
type SwaggerPaginatedData struct {
	List       any   `json:"list,omitempty"` // 当前页数据列表（空时返回 [] 而非 null）
	Total      int64 `json:"total"`          // 总记录数
	Page       int   `json:"page"`           // 当前页码（从 1 开始）
	PageSize   int   `json:"pageSize"`       // 每页条数
	TotalPages int   `json:"totalPages"`     // 总页数（自动计算）
}
