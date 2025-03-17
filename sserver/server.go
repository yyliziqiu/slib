package sserver

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sserver/sresp"
)

func Run(config Config, routes ...func(engine *gin.Engine)) error {
	config = config.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	setGinWriter(config)

	engine := createEngine()
	for _, v := range routes {
		v(engine)
	}

	return engine.Run(config.Listen)
}

func setGinWriter(config Config) {
	if _accessLogger == nil && !config.DisableAccessLog {
		_accessLogger = slog.New3(config.AccessLog)
	}
	if _accessLogger != nil {
		gin.DefaultWriter = _accessLogger.Writer()
	} else {
		gin.DefaultWriter = io.Discard
	}

	if _errorLogger == nil {
		_errorLogger = slog.New3(config.ErrorLog)
	}
	gin.DefaultErrorWriter = _errorLogger.Writer()
}

func createEngine() *gin.Engine {
	engine := gin.New()
	engine.NoRoute(sresp.AbortNotFound)
	engine.NoMethod(sresp.AbortMethodNotAllowed)
	engine.Use(gin.LoggerWithFormatter(logFormatter))
	engine.Use(gin.CustomRecovery(recovery))
	return engine
}

func logFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("%3d | %13v | %15s |%-7s %#v\n%s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

func recovery(ctx *gin.Context, err interface{}) {
	_errorLogger.Warnf("Server panic, path: %s, error: %v", ctx.FullPath(), err)
	sresp.AbortInternalServerError(ctx)
}

func GetErrorLogger() *logrus.Logger {
	return _errorLogger
}

func GetAccessLogger() *logrus.Logger {
	return _accessLogger
}

func SetErrorLogger(logger *logrus.Logger) {
	_errorLogger = logger
}

func SetAccessLogger(logger *logrus.Logger) {
	_accessLogger = logger
}
