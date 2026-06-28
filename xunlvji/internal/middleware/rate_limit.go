package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xunlvji/xunlvji/pkg/errcode"
	"github.com/xunlvji/xunlvji/pkg/response"
)

// RateLimiter 基于Redis的固定窗口限流器
type RateLimiter struct {
	rdb    *redis.Client
	window time.Duration
	limit  int
}

// NewRateLimiter 创建限流器
// window: 时间窗口（如 1*time.Second 表示每秒）
// limit: 窗口内最大请求数
func NewRateLimiter(rdb *redis.Client, window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{rdb: rdb, window: window, limit: limit}
}

// Allow 判断是否允许请求
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	if rl.rdb == nil {
		return true, nil // Redis 不可用时放行
	}
	now := time.Now().UnixNano()
	windowStart := now - (now % int64(rl.window))
	redisKey := fmt.Sprintf("rate_limit:%s:%d", key, windowStart)

	pipe := rl.rdb.TxPipeline()
	incr := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, rl.window+time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return true, nil // Redis 出错时降级放行
	}
	return incr.Val() <= int64(rl.limit), nil
}

// RateLimit 限流中间件
// limit: 窗口内最大请求数
// window: 时间窗口
// keyFn: 限流键生成函数（如按IP或用户ID）
func RateLimit(rdb *redis.Client, limit int, window time.Duration, keyFn func(c *gin.Context) string) gin.HandlerFunc {
	rl := NewRateLimiter(rdb, window, limit)
	return func(c *gin.Context) {
		key := keyFn(c)
		if key == "" {
			c.Next()
			return
		}
		allowed, err := rl.Allow(c.Request.Context(), key)
		if err != nil || !allowed {
			response.Error(c, errcode.ErrTooManyRequests, errcode.GetMsg(errcode.ErrTooManyRequests))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitByIP 按IP限流
func RateLimitByIP(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return RateLimit(rdb, limit, window, func(c *gin.Context) string {
		return "ip:" + c.ClientIP()
	})
}

// RateLimitByUser 按用户ID限流（未登录用户回退到IP）
func RateLimitByUser(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return RateLimit(rdb, limit, window, func(c *gin.Context) string {
		uid := GetUserID(c)
		if uid > 0 {
			return fmt.Sprintf("user:%d", uid)
		}
		return "ip:" + c.ClientIP()
	})
}
