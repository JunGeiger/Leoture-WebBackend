package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ---------------------------------------------------------------------------
// BizError 业务错误类型
// ---------------------------------------------------------------------------

// BizError 是项目统一的业务错误类型，携带错误码、消息、原始错误和元数据。
//
// 设计要点:
//   - 所有链式方法 (WithErr / WithMsg / WithMeta) 返回 **新实例**，
//     不修改原对象，确保包级哨兵错误 (ErrXxx) 在并发场景下安全不可变
//   - 实现 Is / As / Unwrap 完整适配 Go 标准库 errors 包
type BizError struct {
	Code    ResponseCode   `json:"code"`           // 业务错误码
	Message string         `json:"message"`        // 面向用户的错误消息
	Err     error          `json:"-"`              // 原始错误（不序列化，仅日志可见）
	Meta    map[string]any `json:"meta,omitempty"` // 附加元数据（如字段校验详情）
}

// Error 实现 error 接口。
// 输出格式: "code=4004: resource not found" 或 "code=4004: resource not found: <原始错误>"
func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d: %s", e.Code, e.Message)
}

// Unwrap 返回原始错误，使 errors.Is / errors.As 可沿错误链向上追溯。
func (e *BizError) Unwrap() error {
	return e.Err
}

// Is 实现 errors.Is 自定义匹配逻辑：只要业务码相同即视为同类错误。
//
// 这使得以下用法成立:
//
//	err := ErrNotFound.WithErr(sql.ErrNoRows)
//	errors.Is(err, ErrNotFound) // → true（按 Code 匹配，而非指针比较）
func (e *BizError) Is(target error) bool {
	var t *BizError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// HTTPStatus 返回该业务错误对应的 HTTP 状态码。
func (e *BizError) HTTPStatus() int {
	return e.Code.HTTPStatus()
}

// copy 浅拷贝结构体，并深拷贝 Meta map，防止新旧实例共享 map 引用。
func (e *BizError) copy() *BizError {
	ne := *e
	if e.Meta != nil {
		ne.Meta = make(map[string]any, len(e.Meta))
		for k, v := range e.Meta {
			ne.Meta[k] = v
		}
	}
	return &ne
}

// WithErr 附加原始错误，返回新实例。
//
// 用法:
//
//	return errors.ErrNotFound.WithErr(err)
func (e *BizError) WithErr(err error) *BizError {
	ne := e.copy()
	ne.Err = err
	return ne
}

// WithMsg 覆盖错误消息，返回新实例。
//
// 用法:
//
//	return errors.ErrNotFound.WithMsg("用户不存在")
func (e *BizError) WithMsg(msg string) *BizError {
	ne := e.copy()
	ne.Message = msg
	return ne
}

// WithMeta 添加元数据键值对，返回新实例。支持链式调用。
//
// 用法:
//
//	return errors.ErrInvalidParam.
//	    WithMeta("field", "email").
//	    WithMeta("reason", "invalid format")
func (e *BizError) WithMeta(key string, value any) *BizError {
	ne := e.copy()
	if ne.Meta == nil {
		ne.Meta = make(map[string]any)
	}
	ne.Meta[key] = value
	return ne
}

// New 创建一个不带原始错误的 BizError。
//
// 用法:
//
//	return errors.New(errors.CodeNotFound, "用户不存在")
func New(code ResponseCode, msg string) *BizError {
	return &BizError{Code: code, Message: msg}
}

// Wrap 创建一个携带原始错误的 BizError。
// 如果 err 为 nil，等价于 New。
//
// 用法:
//
//	user, err := repo.FindByID(id)
//	if err != nil {
//	    return nil, errors.Wrap(errors.CodeDatabaseError, "查询用户失败", err)
//	}
func Wrap(code ResponseCode, msg string, err error) *BizError {
	return &BizError{Code: code, Message: msg, Err: err}
}

// ---------------------------------------------------------------------------
// 错误提取
// ---------------------------------------------------------------------------

// FromError 从任意 error 中提取或构造 *BizError。
//
// 匹配规则（按优先级）:
//  1. err 本身或错误链中包含 *BizError → 直接返回
//  2. err 为 nil → 返回 nil
//  3. 其他情况 → 包装为 CodeInternal 错误
//
// 用法:
//
//	bizErr := errors.FromError(err)
//	c.JSON(bizErr.HTTPStatus(), bizErr)
func FromError(err error) *BizError {
	if err == nil {
		return nil
	}
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr
	}
	return Wrap(CodeInternal, "Internal server error", err)
}

// ---------------------------------------------------------------------------
// 哨兵错误（Sentinel Errors）
// ---------------------------------------------------------------------------
//
// 使用规范:
//   - 哨兵错误是 **不可变模板**，禁止直接修改其字段
//   - 必须通过 WithErr / WithMsg / WithMeta 派生新实例后使用
//   - 在 handler / service / repository 中统一使用派生实例返回
//
// 示例:
//
//	// 正确: 派生新实例
//	return errors.ErrNotFound.WithMsg("用户不存在").WithErr(err)
//
//	// 错误: 直接修改哨兵（并发不安全，且污染全局状态）
//	errors.ErrNotFound.Err = err
//	errors.ErrNotFound.Message = "用户不存在"
var (
	// ---- 4xxx 客户端错误 ----

	// ErrInvalidParam 请求参数校验失败
	// 场景: binding.Validate 失败、自定义校验不通过、JSON 反序列化错误
	ErrInvalidParam = New(CodeInvalidParam, "Parameter verification failed")

	// ErrUnauthenticated 未认证
	// 场景: 请求未携带 Token、Token 格式非法、签名校验失败
	ErrUnauthenticated = New(CodeUnauthenticated, "Authentication failed")

	// ErrTokenExpired 令牌过期
	// 场景: JWT Access Token 过期，前端应使用 Refresh Token 刷新
	ErrTokenExpired = New(CodeTokenExpired, "Token expired")

	// ErrForbidden 无权限
	// 场景: Casbin 策略判定当前用户无权执行目标操作
	ErrForbidden = New(CodeForbidden, "Permission denied")

	// ErrNotFound 资源不存在
	// 场景: 按 ID 查询用户/角色/菜单等实体未命中
	ErrNotFound = New(CodeNotFound, "Resource not found")

	// ErrAlreadyExists 资源已存在
	// 场景: 创建时违反唯一约束（用户名重复、角色名冲突等）
	ErrAlreadyExists = New(CodeAlreadyExists, "Resource already exists")

	// ErrTooManyRequests 请求过于频繁
	// 场景: 登录失败次数过多、接口调用频率超限
	ErrTooManyRequests = New(CodeTooManyRequests, "Too many requests")

	// ---- 5xxx 服务端错误 ----

	// ErrInternal 服务器内部错误（兜底）
	// 场景: 未归类的服务端异常，生产环境不向客户端暴露内部细节
	ErrInternal = New(CodeInternal, "Internal server error")

	// ErrServiceUnavailable 服务暂不可用
	// 场景: 系统维护、依赖组件未就绪、服务降级
	ErrServiceUnavailable = New(CodeServiceUnavailable, "Service unavailable")

	// ErrTimeout 请求超时
	// 场景: 数据库查询超时、外部 API 调用超时
	ErrTimeout = New(CodeTimeout, "Request timeout")

	// ErrDatabase 数据库操作异常
	// 场景: 连接池耗尽、SQL 执行失败、事务提交异常
	ErrDatabase = New(CodeDatabaseError, "Database error")

	// ErrCache 缓存操作异常
	// 场景: Redis 连接失败、序列化异常
	ErrCache = New(CodeCacheError, "Cache error")
)

// ---------------------------------------------------------------------------
// 外部 HTTP 响应转换
// ---------------------------------------------------------------------------

// FromHTTPStatus 将外部依赖返回的 HTTP 状态码转换为内部 BizError。
//
// 主要用于调用第三方 API 或外部微服务时，将外部 HTTP 错误统一映射为内部业务错误，
// 起到"防腐层"（Anti-Corruption Layer）的作用，使上层业务无需关心外部协议细节。
//
// 参数:
//   - statusCode: 外部 HTTP 响应状态码
//   - body: 外部响应体（可选，用于提取外部错误详情写入 Meta）
//   - err: 原始错误（可选，如网络超时、DNS 解析失败等底层错误）
func FromHTTPStatus(statusCode int, body []byte, err error) *BizError {
	// 网络层错误（非 HTTP 响应错误）: 超时、DNS 失败、连接拒绝等
	if statusCode == 0 {
		bizErr := ErrServiceUnavailable.WithMsg("External service unreachable")
		if err != nil {
			bizErr = bizErr.WithErr(err)
		}
		return bizErr
	}

	// 根据 HTTP 状态码映射为内部业务错误码
	var bizErr *BizError
	switch {
	case statusCode == http.StatusBadRequest:
		bizErr = ErrInvalidParam.WithMsg("External service rejected request")

	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		bizErr = ErrUnauthenticated.WithMsg("External service authentication failed")

	case statusCode == http.StatusNotFound:
		bizErr = ErrNotFound.WithMsg("External resource not found")

	case statusCode == http.StatusConflict:
		bizErr = ErrAlreadyExists.WithMsg("External resource conflict")

	case statusCode == http.StatusTooManyRequests:
		bizErr = ErrTooManyRequests.WithMsg("External service rate limited")

	case statusCode == http.StatusRequestTimeout || statusCode == http.StatusGatewayTimeout:
		bizErr = ErrTimeout.WithMsg("External service timeout")

	case statusCode >= 500:
		bizErr = ErrServiceUnavailable.WithMsg("External service unavailable")

	default:
		bizErr = ErrInternal.WithMsg(fmt.Sprintf("External service error: %d %s",
			statusCode, http.StatusText(statusCode)))
	}

	// 附加原始错误
	if err != nil {
		bizErr = bizErr.WithErr(err)
	}

	// 将外部响应体写入 Meta，便于日志排查（截取前 500 字节防止日志膨胀）
	if len(body) > 0 {
		truncated := body
		if len(truncated) > 500 {
			truncated = truncated[:500]
		}
		bizErr = bizErr.WithMeta("external_response", string(truncated))
	}

	// 记录外部状态码，便于监控和告警
	bizErr = bizErr.WithMeta("external_status", statusCode)

	return bizErr
}
