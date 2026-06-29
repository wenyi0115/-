package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xunlvji/xunlvji/internal/config"
	"github.com/xunlvji/xunlvji/internal/router"
	"github.com/xunlvji/xunlvji/pkg/database"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(logger.Config{
		Level:    cfg.Log.Level,
		Format:   cfg.Log.Format,
		Output:   cfg.Log.Output,
		FilePath: cfg.Log.FilePath,
	}); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.L().Info("寻旅记后端启动中",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("mode", cfg.App.Mode),
	)

	// 3. 设置 Gin 模式
	gin.SetMode(cfg.App.Mode)

	// 4. 初始化 MySQL（失败不阻断启动，允许无DB调试）
	db, err := database.InitMySQL(cfg.MySQL, cfg.App.Mode)
	if err != nil {
		logger.L().Warn("MySQL 连接失败，服务以无DB模式启动（仅供开发调试）", zap.Error(err))
	} else {
		logger.L().Info("MySQL 连接成功",
			zap.String("host", cfg.MySQL.Host),
			zap.Int("port", cfg.MySQL.Port),
			zap.String("database", cfg.MySQL.Database),
		)
	}

	// 5. 初始化 Redis
	rdb, err := database.InitRedis(cfg.Redis)
	if err != nil {
		logger.L().Warn("Redis 连接失败，限流/缓存功能将降级", zap.Error(err))
	} else {
		logger.L().Info("Redis 连接成功",
			zap.String("addr", cfg.Redis.Addr()),
			zap.Int("db", cfg.Redis.DB),
		)
	}

	// 6. 初始化 JWT 管理器
	jwtMgr := jwt.NewManager(cfg.JWT)

	// 7. 注册路由
	r := router.New(router.Dependencies{
		DB:     db,
		RDB:    rdb,
		JWTMgr: jwtMgr,
		Config: cfg,
	})

	// 8. 启动 HTTP 服务
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.L().Info("HTTP 服务监听", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Fatal("HTTP 服务异常退出", zap.Error(err))
		}
	}()

	// 9. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.L().Info("收到退出信号，开始优雅关闭", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error("HTTP 服务关闭超时", zap.Error(err))
	}
	if db != nil {
		if err := database.CloseMySQL(db); err != nil {
			logger.L().Error("关闭 MySQL 连接失败", zap.Error(err))
		}
	}
	if rdb != nil {
		if err := database.CloseRedis(rdb); err != nil {
			logger.L().Error("关闭 Redis 连接失败", zap.Error(err))
		}
	}

	logger.L().Info("寻旅记后端已关闭")
}
