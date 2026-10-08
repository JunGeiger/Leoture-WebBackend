package response

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/types/request"

	"github.com/gin-gonic/gin"
)

// Response HTTP接口统一响应信封
type Response[T any] struct {
	Code      errors.ResponseCode `json:"code"`                // 业务状态码
	Message   string              `json:"message"`             // 提示信息
	Data      *T                  `json:"data,omitempty"`      // 业务数据
	RequestID string              `json:"requestId,omitempty"` // 请求唯一标识
}

// Success 返回标准的成功响应（HTTP 200）
func Success[T any](c *gin.Context, data T) *Response[T] {
	return &Response[T]{
		Code:      errors.CodeSuccess,
		Message:   "success",
		Data:      &data,
		RequestID: c.GetString(request.RequestIDKey),
	}
}

// Failure 组装统一错误响应
func Failure(c *gin.Context, code errors.ResponseCode, msg string) Response[any] {
	return Response[any]{
		Code:      code,
		Message:   msg,
		Data:      nil,
		RequestID: c.GetString(request.RequestIDKey),
	}
}

// Paginated 分页数据结构
// 如果分页数据未传，则需要手动初始化为第一页参数
type Paginated[T any] struct {
	List       []*T  `json:"list,omitempty"` // 当前页数据列表（空时返回 [] 而非 null）
	Total      int64 `json:"total"`          // 总记录数
	Page       int   `json:"page"`           // 当前页码（从 1 开始）
	PageSize   int   `json:"pageSize"`       // 每页条数
	TotalPages int   `json:"totalPages"`     // 总页数（自动计算）
}

// PagingData 返回分页成功响应（HTTP 200）。
// 自动计算总页数，并确保 List 为空时返回 [] 而非 null（前端无需额外判空）。
func PagingData[T any](list []*T, total int64, page int, pageSize int) *Paginated[T] {
	// 确保空列表返回 [] 而非 null，避免前端额外判空
	if list == nil {
		list = []*T{}
	}
	var totalPages int64
	if total == 0 {
		totalPages = 0
	} else {
		totalPages = (total + int64(pageSize) - 1) / int64(pageSize)
	}
	return &Paginated[T]{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}
}
