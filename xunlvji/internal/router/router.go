package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xunlvji/xunlvji/internal/middleware"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/response"
	"gorm.io/gorm"
)

// Dependencies 路由依赖
type Dependencies struct {
	DB      *gorm.DB
	RDB     *redis.Client
	JWTMgr  *jwt.Manager
}

// New 创建并注册路由
func New(deps Dependencies) *gin.Engine {
	r := gin.New()

	// 全局中间件（顺序：Recovery → RequestID → Logger → CORS）
	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(),
	)

	// 健康检查（无需鉴权）
	r.GET("/health", healthCheck(deps))
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	// 全局限流：每IP每秒60次
	v1.Use(middleware.RateLimitByIP(deps.RDB, 60, time.Second))

	// --- 公开接口（无需鉴权）---
	pub := v1.Group("")
	{
		// 健康检查
		pub.GET("/status", healthCheck(deps))
		// 后续接入：微信登录(code2session)、短信验证码发送等
	}

	// --- 可选鉴权接口（登录态可选，用于个性化推荐） ---
	optional := v1.Group("")
	optional.Use(middleware.OptionalAuth(deps.JWTMgr))
	{
		// 后续接入：首页Feed流、推荐、攻略列表、路线列表、地点详情等
	}

	// --- 鉴权接口（必须登录）---
	auth := v1.Group("")
	auth.Use(middleware.Auth(deps.JWTMgr))
	{
		// 后续接入：发布内容、打卡、行程管理、关注/点赞/收藏/评论等
	}

	return r
}

// healthCheck 健康检查处理
func healthCheck(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"app":     "xunlvji",
			"version": "v2.1.7",
			"status":  "ok",
		}
		if deps.DB != nil {
			if sqlDB, err := deps.DB.DB(); err == nil {
				if err := sqlDB.Ping(); err != nil {
					status["mysql"] = "down"
				} else {
					status["mysql"] = "up"
				}
			}
		}
		if deps.RDB != nil {
			if err := deps.RDB.Ping(c).Err(); err != nil {
				status["redis"] = "down"
			} else {
				status["redis"] = "up"
			}
		}
		response.Success(c, status)
	}
}
