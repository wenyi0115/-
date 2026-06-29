package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	std  *zap.Logger
	once sync.Once
)

// Config 日志配置
type Config struct {
	Level    string // debug / info / warn / error
	Format   string // json / text
	Output   string // stdout / file / both
	FilePath string // 文件路径（output=file/both 时生效）
}

// Init 初始化全局日志器
func Init(cfg Config) error {
	var initErr error
	once.Do(func() {
		level, err := parseLevel(cfg.Level)
		if err != nil {
			initErr = err
			return
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		var encoder zapcore.Encoder
		if cfg.Format == "text" {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}

		cores := make([]zapcore.Core, 0, 2)
		// stdout 输出
		if cfg.Output == "stdout" || cfg.Output == "both" {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level))
		}
		// 文件输出（带滚动切割）
		if cfg.Output == "file" || cfg.Output == "both" {
			if cfg.FilePath == "" {
				cfg.FilePath = "./logs/xunlvji.log"
			}
			if dir := filepath.Dir(cfg.FilePath); dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}
			lj := &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    200, // MB
				MaxBackups: 7,
				MaxAge:     30, // 天
				Compress:   true,
			}
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(lj), level))
		}
		if len(cores) == 0 {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level))
		}

		core := zapcore.NewTee(cores...)
		std = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0), zap.AddStacktrace(zapcore.ErrorLevel))
	})
	return initErr
}

// L 返回全局日志器（未初始化时返回 Nop）
func L() *zap.Logger {
	if std == nil {
		return zap.NewNop()
	}
	return std
}

// S 返回全局 SugaredLogger
func S() *zap.SugaredLogger {
	if std == nil {
		return zap.NewNop().Sugar()
	}
	return std.Sugar()
}

// Sync 刷新缓冲
func Sync() {
	if std != nil {
		_ = std.Sync()
	}
}

// WithRequestID 返回带 request_id 字段的 logger
func WithRequestID(requestID string) *zap.Logger {
	return L().With(zap.String("request_id", requestID))
}

func parseLevel(s string) (zapcore.Level, error) {
	switch s {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "":
		return zapcore.InfoLevel, nil
	default:
		return 0, fmt.Errorf("未知的日志级别: %s", s)
	}
}
