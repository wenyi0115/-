package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xunlvji/xunlvji/pkg/errcode"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/response"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// ContextUserID 用户ID
	ContextUserID ContextKey = "user_id"
	// ContextOpenID 微信OpenID
	ContextOpenID ContextKey = "open_id"
	// ContextNickname 昵称
	ContextNickname ContextKey = "nickname"
	// ContextAvatar 头像
	ContextAvatar ContextKey = "avatar"
)

// AuthorizationHeader 鉴权请求头
const AuthorizationHeader = "Authorization"

// tokenPrefix Bearer 前缀
const tokenPrefix = "Bearer "

// Auth JWT 鉴权中间件（强制校验，失败返回未登录）
func Auth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Error(c, errcode.ErrNotLoggedIn, errcode.GetMsg(errcode.ErrNotLoggedIn))
			c.Abort()
			return
		}
		claims, err := mgr.ParseAccessToken(token)
		if err != nil {
			// 区分过期与无效
			if err == jwt.ErrTokenExpired {
				response.Error(c, errcode.ErrTokenExpired, errcode.GetMsg(errcode.ErrTokenExpired))
			} else {
				response.Error(c, errcode.ErrTokenInvalid, errcode.GetMsg(errcode.ErrTokenInvalid))
			}
			c.Abort()
			return
		}
		c.Set(string(ContextUserID), claims.UserID)
		c.Set(string(ContextOpenID), claims.OpenID)
		c.Set(string(ContextNickname), claims.Nickname)
		c.Set(string(ContextAvatar), claims.Avatar)
		c.Next()
	}
}

// OptionalAuth 可选鉴权中间件（已登录则注入用户信息，未登录放行）
func OptionalAuth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}
		claims, err := mgr.ParseAccessToken(token)
		if err == nil {
			c.Set(string(ContextUserID), claims.UserID)
			c.Set(string(ContextOpenID), claims.OpenID)
			c.Set(string(ContextNickname), claims.Nickname)
			c.Set(string(ContextAvatar), claims.Avatar)
		}
		c.Next()
	}
}

// extractToken 从请求头提取 Token
func extractToken(c *gin.Context) string {
	auth := c.GetHeader(AuthorizationHeader)
	if auth == "" {
		return ""
	}
	if strings.HasPrefix(auth, tokenPrefix) {
		return strings.TrimPrefix(auth, tokenPrefix)
	}
	return auth
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	if v, exists := c.Get(string(ContextUserID)); exists {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// GetOpenID 从上下文获取OpenID
func GetOpenID(c *gin.Context) string {
	if v, exists := c.Get(string(ContextOpenID)); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
