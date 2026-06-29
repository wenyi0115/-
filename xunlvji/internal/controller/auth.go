// Package controller 控制器层（HTTP 入口）
package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/xunlvji/xunlvji/internal/middleware"
	"github.com/xunlvji/xunlvji/internal/service"
	"github.com/xunlvji/xunlvji/pkg/errcode"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/response"
	"github.com/xunlvji/xunlvji/pkg/wechat"
	"go.uber.org/zap"
)

// AuthController 认证控制器
type AuthController struct {
	authSvc *service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authSvc *service.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

// wechatLoginRequest 微信登录请求
type wechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Platform string `json:"platform" binding:"required"` // miniprogram / app
	IV       string `json:"iv,omitempty"`
	EncryptedData string `json:"encrypted_data,omitempty"`
}

// WechatLogin 微信登录
// POST /api/v1/auth/login/wechat
func (ctrl *AuthController) WechatLogin(c *gin.Context) {
	var req wechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "缺少必要参数: code 和 platform")
		return
	}

	// 平台校验（当前仅支持小程序，App 端后续扩展）
	if req.Platform != "miniprogram" && req.Platform != "app" {
		response.Error(c, errcode.ErrInvalidParam, "platform 必须为 miniprogram 或 app")
		return
	}
	if req.Platform == "app" {
		// App 端登录需要 iv + encrypted_data，P0 阶段暂不支持
		response.Error(c, errcode.ErrInvalidParam, "暂不支持 App 端微信登录")
		return
	}

	result, err := ctrl.authSvc.WechatLogin(c.Request.Context(), req.Code)
	if err != nil {
		ctrl.handleAuthError(c, err)
		return
	}

	response.Success(c, result)
}

// refreshTokenRequest 刷新 Token 请求
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新 Token
// POST /api/v1/auth/refresh-token
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "缺少必要参数: refresh_token")
		return
	}

	result, err := ctrl.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenInvalid) {
			response.Error(c, errcode.ErrRefreshTokenExpired, errcode.GetMsg(errcode.ErrRefreshTokenExpired))
			return
		}
		zap.L().Error("刷新 token 失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
		return
	}

	response.Success(c, result)
}

// handleAuthError 处理认证服务错误
func (ctrl *AuthController) handleAuthError(c *gin.Context, err error) {
	rid := middleware.GetRequestID(c)
	switch {
	case errors.Is(err, service.ErrWechatLoginFailed):
		zap.L().Warn("微信登录失败",
			zap.String("request_id", rid),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrOAuthFailed, errcode.GetMsg(errcode.ErrOAuthFailed))
	case errors.Is(err, service.ErrAccountBanned):
		response.Error(c, errcode.ErrAccountBanned, errcode.GetMsg(errcode.ErrAccountBanned))
	case errors.Is(err, service.ErrAccountCanceled):
		response.Error(c, errcode.ErrAccountDeleted, errcode.GetMsg(errcode.ErrAccountDeleted))
	case errors.Is(err, service.ErrRefreshTokenInvalid):
		response.Error(c, errcode.ErrRefreshTokenExpired, errcode.GetMsg(errcode.ErrRefreshTokenExpired))
	case errors.Is(err, jwt.ErrTokenExpired):
		response.Error(c, errcode.ErrTokenExpired, errcode.GetMsg(errcode.ErrTokenExpired))
	case errors.Is(err, jwt.ErrTokenInvalid):
		response.Error(c, errcode.ErrTokenInvalid, errcode.GetMsg(errcode.ErrTokenInvalid))
	case errors.Is(err, wechat.ErrInvalidCode):
		response.Error(c, errcode.ErrInvalidParam, "无效的微信 code")
	case errors.Is(err, wechat.ErrInvalidAppID):
		zap.L().Error("微信 AppID 未配置",
			zap.String("request_id", rid),
		)
		response.Error(c, errcode.ErrInternal, "微信登录未配置")
	default:
		zap.L().Error("认证服务异常",
			zap.String("request_id", rid),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
	}
}
