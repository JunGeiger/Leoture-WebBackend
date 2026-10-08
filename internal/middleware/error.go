package middleware

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// ErrorHandler 是 Gin 全局错误处理中间件。
//
// 职责:
//  1. 捕获 handler 链中未被处理的 panic，防止进程崩溃
//  2. 统一处理通过 c.Error() 推入 c.Errors 的业务错误
//  3. 根据错误类型（4xx/5xx）分级记录日志，5xx 记录原始错误链，4xx 仅记录摘要
//  4. 向前端返回标准化的 JSON 错误响应
//
// 注册位置:
//
//	必须作为第一个中间件注册，确保能捕获后续所有中间件和 handler 的 panic
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Panic 恢复
		defer func() {
			if r := recover(); r != nil {
				// panic 属于严重的服务端异常，必须记录完整堆栈
				stack := string(debug.Stack())
				slog.Error(
					"Panic recovered",
					slog.Any("panic", r),
					slog.String("stack", stack),
					slog.String("request_id", c.GetString(request.RequestIDKey)),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
				)

				// 检查是否已经写过响应（极端情况：panic 发生在写响应之后）
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError,
						response.Failure(c, errors.CodeInternal, "Internal server error"))
				}
			}
		}()

		// 执行后续 handler 链
		c.Next()

		// 常规错误处理
		// 如果没有错误，而且handler 已经写过响应（如某些特殊接口自行处理了响应），不再重复写入，直接返回
		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		// 提取最后一个错误（Gin 的 c.Error() 会将错误追加到 c.Errors 切片）
		lastErr := c.Errors.Last().Err
		if lastErr == nil {
			return
		}
		// 统一转换为 BizError
		bizErr := errors.FromError(lastErr)

		// 获取 HTTP 状态码
		statusCode := bizErr.HTTPStatus()
		logError(c, bizErr, statusCode)
		c.JSON(statusCode, response.Failure(c, bizErr.Code, bizErr.Message))
	}
}

// logError 根据 HTTP 状态码分级记录错误日志。
// 根据错误类型（4xx/5xx）分级记录日志，5xx 记录原始错误链，4xx 仅记录摘要
func logError(c *gin.Context, bizErr *errors.BizError, statusCode int) {
	// 构建基础日志字段
	attrs := []slog.Attr{
		slog.Int("code", int(bizErr.Code)),
		slog.String("message", bizErr.Message),
		slog.String("request_id", c.GetString(request.RequestIDKey)),
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
	}

	if len(bizErr.Meta) > 0 {
		attrs = append(attrs, slog.Any("meta", bizErr.Meta))
	}

	if statusCode >= http.StatusInternalServerError {
		// 5xx: 记录原始错误链，而非当前 goroutine 的调用栈
		if bizErr.Err != nil {
			attrs = append(attrs, slog.String("root_error", bizErr.Err.Error()))
		}
		slog.LogAttrs(c.Request.Context(), slog.LevelError, "Server error", attrs...)
	} else {
		// 4xx: Warn 级别，不抓取堆栈
		slog.LogAttrs(c.Request.Context(), slog.LevelWarn, "Client error", attrs...)
	}
}

// abortWithBizError 是中间件专用的错误中断辅助函数。
//
// 与 c.AbortWithError 的区别:
//   - 无需手动传入 HTTP 状态码，自动从 BizError 中提取
//   - 语义明确，调用方只需关注"发生了什么业务错误"
//   - 与 Handler 层的 c.Error(err) + return 保持一致的错误对象模型
//
// 注意: 此函数仅应在中间件中使用。Handler 层应使用 c.Error(err) + return，
// 由 ErrorHandler 全局中间件统一处理响应。
func abortWithBizError(c *gin.Context, bizErr *errors.BizError) {
	err := c.Error(bizErr) // 注入错误到 c.Errors，供 ErrorHandler 读取
	if err != nil {
		slog.Warn("gin: Abort with bizError，but an error has occurred", slog.Any("bizErr", bizErr))
	}
	c.Abort() // 中断后续 handler/中间件执行
}
