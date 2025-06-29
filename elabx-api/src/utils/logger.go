// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/9 9:28
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : logger.go
// @Software: GoLand
package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func SetupLogger(outputPath string, loglevel string) *zap.Logger {
	// 日志分割
	// MaxSize:定义了日志文件的最大大小，单位是MB。
	// MaxBackups:定义了最多保留的备份文件数量。当备份文件数量超过MaxBackups后，lumberjack会自动删除最旧的备份文件。
	// MaxAge:定义了备份文件的最大保存天数。当备份文件的保存天数超过MaxAge后，lumberjack会自动删除备份文件。
	// Compress:定义了备份文件是否需要压缩。如果设置为true，备份的日志文件会被压缩为.gz格式。
	hook := lumberjack.Logger{
		Filename:   outputPath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    10,         // 每个日志文件保存10M，默认 100M
		MaxBackups: 30,         // 保留30个备份，默认不限
		MaxAge:     7,          // 保留7天，默认不限
		Compress:   false,      // 是否压缩，默认不压缩
	}

	write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
		}, //输出的时间格式
		EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写日志级别
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 短路径的调用者
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), write, level)

	// 创建 Logger
	//logger, err := config.Build()

	// 设置初始化字段,如：添加一个服务器名称
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	// 创建 Logger
	logger := zap.New(
		core,
		zap.AddCaller(),                   // 添加调用者信息
		zap.AddStacktrace(zap.ErrorLevel), // 从错误级别开始记录堆栈跟踪
	)

	return logger
}

var Apis []gin.RouteInfo

func GetAllRoutes(engine *gin.Engine) {
	for _, r := range engine.Routes() {
		Apis = append(Apis, r)
	}
	return
}

var Logger = SetupLogger("log/elabx_svr.log", "info")
