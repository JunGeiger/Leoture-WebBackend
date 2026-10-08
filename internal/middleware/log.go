package middleware

import (
	"LeotureWeb/internal/types/request"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

// ReqLog HTTP请求日志打印
func ReqLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set(request.RequestTimeKey, start)

		slog.Info("HTTP Request",
			slog.String("request_id", c.GetString(request.RequestIDKey)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_time", start.String()))

		c.Next()
	}
}

// RespLog HTTP返回日志打印
func RespLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var latency time.Duration
		s, ok := c.Get(request.RequestTimeKey)
		if ok {
			start := s.(time.Time)
			latency = time.Since(start)
		}

		slog.Info("HTTP response",
			slog.String("request_id", c.GetString(request.RequestIDKey)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("latency", latency.String()))
	}
}
