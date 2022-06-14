package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/meBazil/hr-mvp/pkg/common-logger/wrapper"
)

var defaultLogger = NewLoggerFrom(context.Background(), zerolog.New(os.Stdout))

func NewLoggerFrom(ctx context.Context, log interface{}) Logger {
	var entity Wrapper

	switch t := log.(type) {
	case *logrus.Entry:
		entity = wrapper.WrapLogrusEntry(t)
	case logrus.Entry:
		entity = wrapper.WrapLogrusEntry(&t)
	case *logrus.Logger:
		entity = wrapper.WrapLogrusLogger(t)
	case logrus.Logger:
		entity = wrapper.WrapLogrusLogger(&t)
	case *zap.Logger:
		entity = wrapper.WrapZap(t)
	case zap.Logger:
		entity = wrapper.WrapZap(&t)
	case *zerolog.Logger:
		entity = wrapper.WrapZerolog(t)
	case zerolog.Logger:
		entity = wrapper.WrapZerolog(&t)
	default:
		panic("unsupported type of logger")
	}

	return NewLogger(ctx, entity)
}

func SetDefaultLogger(log Logger) {
	defaultLogger = log
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

func WithContext(ctx context.Context) Logger {
	return NewLogger(ctx, defaultLogger.Log)
}

func WithContextAndLogger(ctx context.Context, logger Logger) Logger {
	return NewLogger(ctx, logger.Log)
}

func Debug(msg string, fields ...Fields) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Fields) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Fields) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, err error, fields ...Fields) {
	defaultLogger.Error(msg, err, fields...)
}

func Fatal(msg string, err error, fields ...Fields) {
	defaultLogger.Fatal(msg, err, fields...)
}

func LogWithTimeout(infoMessage string, deadlineError error, td time.Duration, fields ...Fields) context.CancelFunc {
	return defaultLogger.LogWithTimeout(infoMessage, deadlineError, td, fields...)
}

func GetNewTraceID() string {
	return uuid.NewV4().String()
}
