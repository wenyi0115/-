package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xunlvji/xunlvji/pkg/logger"
	"go.uber.org/zap"
)

// RequestIDKey 请求ID在上下文中的键
const RequestIDKey = "request_id"

// RequestIDHeader 请求ID的HTTP头名
const RequestIDHeader = "X-Request-Id"

// RequestID 为每个请求注入唯一 request_id
// 优先使用上游传入的 X-Request-Id，否则生成新的 UUID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(RequestIDHeader)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(RequestIDKey, rid)
		c.Writer.Header().Set(RequestIDHeader, rid)
		c.Next()
	}
}

// GetRequestID 从 gin.Context 获取 request_id
func GetRequestID(c *gin.Context) string {
	if v, exists := c.Get(RequestIDKey); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Logger 请求日志中间件（访问日志）
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		rid := GetRequestID(c)
		logger.L().Info("request",
			zap.String("request_id", rid),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
		)
	}
}
