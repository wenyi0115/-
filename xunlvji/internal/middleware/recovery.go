package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/xunlvji/xunlvji/pkg/errcode"
	"github.com/xunlvji/xunlvji/pkg/logger"
	"github.com/xunlvji/xunlvji/pkg/response"
	"go.uber.org/zap"
)

// Recovery Panic 恢复中间件，捕获 panic 防止进程崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				rid := GetRequestID(c)
				logger.L().Error("panic recovered",
					zap.Any("error", rec),
					zap.String("request_id", rid),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.ByteString("stack", debug.Stack()),
				)
				if !c.Writer.Written() {
					response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}

// NotFound 404 处理
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Error(c, errcode.ErrNotFound, "接口不存在")
		c.AbortWithStatus(http.StatusNotFound)
	}
}
