package middleware

import (
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// Trace 注入 request_id
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. request_id
		requestID := c.GetHeader(request.RequestIDHeaderKey)
		if requestID == "" {
			requestID = utils.GenUUIDv7().String()
		}

		// 2. 注入 Gin Context
		c.Set(request.RequestIDKey, requestID)
		// 记录请求开始时间
		c.Set(request.RequestTimeKey, time.Now().String())

		// 3. 注入 HTTP Header（方便下游）
		c.Writer.Header().Set(request.RequestIDHeaderKey, requestID)

		// 4. 注入 request.Context（关键）
		ctx := context.WithValue(
			c.Request.Context(),
			request.RequestIDKey,
			requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
