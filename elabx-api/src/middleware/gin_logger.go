// Package middleware coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/29 11:27
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : gin_logger.go
// @Software: GoLand
package middleware

import (
	"bytes"
	"eLabX/src/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
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

//var Logger = SetupLogger("log/elabx_svr.log", "info")

// GinLogger 接收gin框架默认的日志
//func GinLogger(logger *zap.Logger) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		path := c.Request.URL.Path
//
//		// 使用自定义 ResponseWriter
//		crw := &CustomResponseWriter{
//			body:           bytes.NewBufferString(""),
//			ResponseWriter: c.Writer,
//		}
//		c.Writer = crw
//		// 打印请求信息
//		reqBody, _ := c.GetRawData()
//		// 请求包体写回。
//		if len(reqBody) > 0 {
//			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
//		}
//		c.Next()
//
//		cost := time.Since(start)
//		logger.Info(path,
//			zap.String("method", c.Request.Method), // 请求方法类型 eg: GET
//			zap.String("path", path),               // 请求路径 eg: /test
//			zap.Int("status", c.Writer.Status()),   // 状态码 eg: 200
//			zap.Duration("duration", cost),         // 返回花费时间
//			zap.String("query", string(reqBody)),   // 请求参数 eg: name=1&password=2
//			zap.String("ip", c.ClientIP()),         // 返回真实的客户端IP eg: ::1（这个就是本机IP，ipv6地址）
//			//zap.String("user-agent", c.Request.UserAgent()),                      // 返回客户端的用户代理。 eg: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36
//			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 返回Errors 切片中ErrorTypePrivate类型的错误
//
//		)
//	}
//}

var lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger() (err error) {
	writeSyncer := getLogWriter(utils.GlobalConfig.Output.Logfile, utils.GlobalConfig.Output.MaxSize, utils.GlobalConfig.Output.MaxBackups, utils.GlobalConfig.Output.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(utils.GlobalConfig.Output.Loglevel))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
	} //输出的时间格式
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			//zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
