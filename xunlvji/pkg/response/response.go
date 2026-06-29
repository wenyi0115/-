package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response 统一返回格式（见 18_错误码规范.md）
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:      code,
		Message:   msg,
		Data:      nil,
		RequestID: getRequestID(c),
	})
}

// ErrorWithData 错误响应（带数据，用于部分场景）
func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:      code,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// getRequestID 获取请求ID（从中间件注入，或生成新的）
func getRequestID(c *gin.Context) string {
	if rid, exists := c.Get("request_id"); exists {
		if s, ok := rid.(string); ok && s != "" {
			return s
		}
	}
	return uuid.New().String()
}

// getHTTPStatus 根据错误码映射HTTP状态码（见 18_错误码规范.md）
func getHTTPStatus(code int) int {
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 1000 && code < 2000: // 通用错误
		switch code {
		case 1001, 1002, 1003:
			return http.StatusBadRequest
		case 1004:
			return http.StatusNotFound
		case 1005:
			return http.StatusConflict
		case 1006:
			return http.StatusTooManyRequests
		case 1007:
			return http.StatusGone
		case 1008:
			return http.StatusUpgradeRequired
		case 1009:
			return http.StatusServiceUnavailable
		default:
			return http.StatusBadRequest
		}
	case code >= 2000 && code < 3000: // 认证与授权
		switch code {
		case 2001, 2002, 2003, 2012, 2015:
			return http.StatusUnauthorized
		case 2004, 2005, 2013, 2014:
			return http.StatusForbidden
		case 2006, 2007, 2011:
			return http.StatusBadRequest
		case 2008:
			return http.StatusTooManyRequests
		case 2009:
			return http.StatusConflict
		case 2010:
			return http.StatusNotFound
		default:
			return http.StatusUnauthorized
		}
	case code >= 3000 && code < 4000: // 用户相关
		switch code {
		case 3001:
			return http.StatusNotFound
		case 3002:
			return http.StatusConflict
		case 3007, 3008:
			return http.StatusForbidden
		default:
			return http.StatusBadRequest
		}
	case code >= 4000 && code < 5000: // 内容相关
		switch code {
		case 4001, 4002, 4003, 4004:
			return http.StatusNotFound
		case 4005, 4009:
			return http.StatusConflict
		case 4006:
			return http.StatusAccepted
		case 4010:
			return http.StatusTooManyRequests
		case 4011:
			return http.StatusGone
		default:
			return http.StatusBadRequest
		}
	case code >= 5000 && code < 6000: // 行程相关
		switch code {
		case 5001, 5002:
			return http.StatusNotFound
		default:
			return http.StatusBadRequest
		}
	case code >= 6000 && code < 7000: // 社交相关
		switch code {
		case 6001, 6014:
			return http.StatusNotFound
		case 6002:
			return http.StatusGone
		case 6006, 6013:
			return http.StatusTooManyRequests
		case 6007, 6010, 6011, 6012:
			return http.StatusConflict
		default:
			return http.StatusBadRequest
		}
	case code >= 7000 && code < 8000: // 文件上传
		return http.StatusBadRequest
	case code >= 8000 && code < 9000: // AI服务
		switch code {
		case 8001, 8014:
			return http.StatusServiceUnavailable
		case 8002, 8007, 8012:
			return http.StatusGatewayTimeout
		case 8003:
			return http.StatusOK
		case 8005:
			return http.StatusAccepted
		case 8010:
			return http.StatusBadRequest
		case 8011:
			return http.StatusTooManyRequests
		default:
			return http.StatusInternalServerError
		}
	case code >= 9000 && code < 10000: // 系统错误
		switch code {
		case 9007:
			return http.StatusTooManyRequests
		case 9008:
			return http.StatusServiceUnavailable
		default:
			return http.StatusInternalServerError
		}
	case code >= 15000 && code < 16000: // CPS分销
		switch code {
		case 15003:
			return http.StatusGone
		case 15004:
			return http.StatusNotFound
		case 15005:
			return http.StatusForbidden
		case 15002:
			return http.StatusGatewayTimeout
		default:
			return http.StatusInternalServerError
		}
	case code >= 16000 && code < 17000: // 定位
		switch code {
		case 16001:
			return http.StatusForbidden
		case 16002:
			return http.StatusGatewayTimeout
		case 16003:
			return http.StatusServiceUnavailable
		case 16004:
			return http.StatusInternalServerError
		case 16005:
			return http.StatusNotFound
		default:
			return http.StatusBadRequest
		}
	case code >= 18000 && code < 19000: // 行程工具
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
