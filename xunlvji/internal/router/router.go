package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xunlvji/xunlvji/internal/config"
	"github.com/xunlvji/xunlvji/internal/controller"
	"github.com/xunlvji/xunlvji/internal/middleware"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/internal/service"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/response"
	"github.com/xunlvji/xunlvji/pkg/wechat"
	"gorm.io/gorm"
)

// Dependencies 路由依赖
type Dependencies struct {
	DB      *gorm.DB
	RDB     *redis.Client
	JWTMgr  *jwt.Manager
	Config  *config.Config
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

	// 开发预览页面（debug 模式注册）
	if deps.Config.App.Mode == "debug" {
		RegisterDebugRoutes(r)
	}

	// === 依赖注入：仓库层 ===
	// debug 模式且 DB 不可用时，使用 Mock 仓库返回预设荆州数据（开发预览）
	var userRepo repository.UserRepository
	var oauthRepo repository.OAuthRepository
	var placeRepo repository.PlaceRepository
	var routeRepo repository.RouteRepository
	var guideRepo repository.GuideRepository
	var checkinRepo repository.CheckinRepository

	useMock := deps.DB == nil && deps.Config.App.Mode == "debug"
	if useMock {
		// Mock 模式：仅内容仓库用 Mock（用户/认证仓库无 Mock，登录走真实逻辑会失败）
		placeRepo = repository.NewMockPlaceRepository()
		routeRepo = repository.NewMockRouteRepository()
		guideRepo = repository.NewMockGuideRepository()
		checkinRepo = repository.NewMockCheckinRepository()
	} else {
		userRepo = repository.NewUserRepository(deps.DB)
		oauthRepo = repository.NewOAuthRepository(deps.DB)
		placeRepo = repository.NewPlaceRepository(deps.DB)
		routeRepo = repository.NewRouteRepository(deps.DB)
		guideRepo = repository.NewGuideRepository(deps.DB)
		checkinRepo = repository.NewCheckinRepository(deps.DB)
	}

	// 行程仓库始终用真实实现（Mock 模式下用独立 DB 实例或后续扩展）
	var tripRepo repository.TripRepository
	if deps.DB != nil {
		tripRepo = repository.NewTripRepository(deps.DB)
	}

	// === 依赖注入：服务层 ===
	wechatCli := wechat.NewClient(deps.Config.WeChat)
	authSvc := service.NewAuthService(deps.DB, userRepo, oauthRepo, wechatCli, deps.JWTMgr)
	userSvc := service.NewUserService(userRepo)
	contentSvc := service.NewContentService(placeRepo, routeRepo, guideRepo, checkinRepo)
	tripSvc := service.NewTripService(tripRepo, routeRepo, placeRepo)

	// === 依赖注入：控制器层 ===
	authCtrl := controller.NewAuthController(authSvc)
	userCtrl := controller.NewUserController(userSvc)
	contentCtrl := controller.NewContentController(contentSvc)
	tripCtrl := controller.NewTripController(tripSvc)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	// 全局限流：每IP每秒60次
	v1.Use(middleware.RateLimitByIP(deps.RDB, 60, time.Second))

	// --- 公开接口（无需鉴权）---
	pub := v1.Group("")
	{
		pub.GET("/status", healthCheck(deps))

		// 认证相关
		auth := pub.Group("/auth")
		{
			// 微信登录
			auth.POST("/login/wechat", authCtrl.WechatLogin)
			// 刷新 Token
			auth.POST("/refresh-token", authCtrl.RefreshToken)
			// TODO: 后续接入
			// auth.POST("/send-code", smsCtrl.SendCode)         // 发送验证码
			// auth.POST("/login/phone", authCtrl.PhoneLogin)    // 手机号验证码登录
			// auth.POST("/login/apple", authCtrl.AppleLogin)    // Apple ID 登录
			// auth.POST("/login/google", authCtrl.GoogleLogin)  // Google 登录
			// auth.POST("/login/email", authCtrl.EmailLogin)    // Email 登录
		}
	}

	// --- 可选鉴权接口（登录态可选，用于个性化推荐）---
	optional := v1.Group("")
	optional.Use(middleware.OptionalAuth(deps.JWTMgr))
	{
		// 获取其他用户公开信息（游客也可访问，未登录 is_following=false）
		optional.GET("/user/:user_id/profile", userCtrl.GetOtherProfile)

		// 内容浏览（Feed 流 + 详情），游客可访问
		content := optional.Group("/content")
		{
			// Feed 流
			feed := content.Group("/feed")
			{
				feed.GET("/recommend", contentCtrl.RecommendFeed) // 推荐流
				feed.GET("/guide", contentCtrl.GuideFeed)         // 攻略流
				feed.GET("/place", contentCtrl.PlaceFeed)         // 打卡地流
				feed.GET("/route", contentCtrl.RouteFeed)         // 路线流
				feed.GET("/nearby", contentCtrl.NearbyFeed)       // 附近流
				// TODO: 后续接入
				// feed.GET("/following", contentCtrl.FollowingFeed) // 关注流 @auth
				// feed.GET("/city", contentCtrl.CityFeed)           // 城市流
			}
			// 详情
			content.GET("/route/:route_id", contentCtrl.RouteDetail) // 路线详情
			content.GET("/place/:place_id", contentCtrl.PlaceDetail) // 地点详情
			content.GET("/guide/:guide_id", contentCtrl.GuideDetail) // 攻略详情
		}
	}

	// --- 鉴权接口（必须登录）---
	authGroup := v1.Group("")
	authGroup.Use(middleware.Auth(deps.JWTMgr))
	{
		// 用户信息
		user := authGroup.Group("/user")
		{
			user.GET("/profile", userCtrl.GetProfile)       // 获取当前用户信息
			user.PUT("/profile", userCtrl.UpdateProfile)    // 更新用户资料
			// TODO: 后续接入
			// user.POST("/bind-phone", userCtrl.BindPhone)   // 绑定手机号
			// user.POST("/bind-email", userCtrl.BindEmail)   // 绑定邮箱
			// user.POST("/logout", authCtrl.Logout)          // 退出登录
		}
		// TODO: 后续接入
		// authGroup.POST("/checkins", checkinCtrl.Create)         // 发布打卡
		// authGroup.POST("/comments", commentCtrl.Create)         // 发表评论
		// authGroup.POST("/follows/:user_id", followCtrl.Follow)  // 关注用户
		// authGroup.POST("/favorites", favoriteCtrl.Create)       // 收藏
		// authGroup.POST("/likes", likeCtrl.Create)               // 点赞

		// 行程管理（必须登录）
		trips := authGroup.Group("/trips")
		{
			trips.POST("/create", tripCtrl.Create)                          // 4.1 创建行程
			trips.GET("/list", tripCtrl.List)                               // 4.2 行程列表
			trips.GET("/:trip_id", tripCtrl.Detail)                        // 4.3 行程详情
			trips.PUT("/:trip_id", tripCtrl.Update)                        // 4.14 编辑行程
			trips.DELETE("/:trip_id", tripCtrl.Delete)                     // 4.10 删除行程
			trips.POST("/:trip_id/start", tripCtrl.Start)                  // 4.4 开始行程
			trips.POST("/:trip_id/complete", tripCtrl.Complete)            // 4.15 结束行程
			trips.POST("/:trip_id/points", tripCtrl.AddPoint)              // 4.12 添加点位
			trips.DELETE("/:trip_id/point/:point_id", tripCtrl.DeletePoint) // 4.13 删除点位
			trips.POST("/:trip_id/point/:point_id/skip", tripCtrl.SkipPoint)     // 4.8 跳过点位
			trips.POST("/:trip_id/point/:point_id/unskip", tripCtrl.UnskipPoint) // 4.9 取消跳过
			// TODO: 4.5/4.6 打卡 / 4.7 打卡同步 / 4.11 AI规划 / 4.16 清单管理
		}
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
