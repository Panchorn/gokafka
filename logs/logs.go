package logs

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestID       = "requestID"
	RequestIDLogger = "RequestIDLogger"
)

var log *zap.Logger

func init() {
	fmt.Println("logs init")
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(requestID string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(RequestID, requestID))
	log.Info(msg, fields...)
}

func Debug(requestID string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(RequestID, requestID))
	log.Debug(msg, fields...)
}

func Error(requestID string, msg interface{}, fields ...zap.Field) {
	fields = append(fields, zap.String(RequestID, requestID))
	switch v := msg.(type) {
	case error:
		log.Error(v.Error(), fields...)
	case string:
		log.Error(v, fields...)
	}
}

func CreateLogContext(c echo.Context) echo.Context {
	return NewContext(c, zap.String(RequestID, c.Response().Header().Get("X-Request-ID")))
}

func NewContext(ctx echo.Context, fields ...zap.Field) echo.Context {
	ctx.Set(RequestIDLogger, log.With(fields...))
	return ctx
}
