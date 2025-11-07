package log

import (
	"mission04/internal/config"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {
	// 配置日志输出位置和滚动策略
	hook := lumberjack.Logger{
		Filename:   config.GetConfig().Log.Path, // 日志文件路径
		MaxSize:    128,                         // 单个日志文件最大尺寸，单位MB
		MaxBackups: 30,                          // 最多保留备份数
		MaxAge:     7,                           // 最多保留天数
		Compress:   true,                        // 是否压缩
	}

	// 设置日志级别
	var levelEnab zapcore.Level
	switch config.GetConfig().Log.Level {
	case "debug":
		levelEnab = zapcore.DebugLevel
	case "info":
		levelEnab = zapcore.InfoLevel
	case "warn":
		levelEnab = zapcore.WarnLevel
	case "error":
		levelEnab = zapcore.ErrorLevel
	default:
		levelEnab = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别大写

	// 配置核心
	// 创建一个复合的zap核心处理器，用于同时处理日志到文件和控制台输出
	// 该核心处理器使用zapcore.NewTee将多个核心处理器组合在一起
	// 第一个核心处理器：将日志以JSON格式写入到指定的hook中（通常是文件）
	// 第二个核心处理器：将日志以控制台友好的格式输出到标准输出
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&hook),
			levelEnab,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			levelEnab,
		),
	)

	// 创建日志实例
	// 初始化全局Logger实例
	// 使用zap库创建一个新的日志记录器，配置了核心输出、调用者信息追踪和调用栈跳过功能
	// core: 日志输出的核心配置，定义了日志的编码格式、输出目的地等基础设置
	// zap.AddCaller(): 添加调用者信息，会在日志中显示调用日志函数的文件名和行号
	// zap.AddCallerSkip(1): 跳过一层调用栈，使调用者信息指向实际调用Logger的代码位置
	// 返回值: 配置完成的zap.Logger实例，赋值给全局Logger变量供整个程序使用
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}
