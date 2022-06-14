package sqlutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type Logger struct {
	Logger                    *logger.Logger
	SlowThreshold             time.Duration
	TraceQueries              bool
	IgnoreRecordNotFoundError bool
}

func (l *Logger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(_ context.Context, msg string, params ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, params...), logger.Fields{
		"filename": utils.FileWithLineNum(),
	})
}

func (l *Logger) Warn(_ context.Context, msg string, params ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, params...), logger.Fields{
		"filename": utils.FileWithLineNum(),
	})
}

func (l *Logger) Error(_ context.Context, msg string, params ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, params...), nil, logger.Fields{
		"filename": utils.FileWithLineNum(),
	})
}

func (l *Logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin).Truncate(time.Millisecond)

	switch {
	case err != nil && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()

		l.Logger.Error(err.Error(), nil, logger.Fields{
			"file":    utils.FileWithLineNum(),
			"elapsed": elapsed,
			"sql":     sql,
			"rows":    rows,
		})
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()

		logger.Warn(fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold), logger.Fields{
			"file":    utils.FileWithLineNum(),
			"elapsed": elapsed,
			"sql":     sql,
			"rows":    rows,
		})
	case l.TraceQueries:
		sql, rows := fc()

		logger.Info("gorm trace", logger.Fields{
			"file":    utils.FileWithLineNum(),
			"elapsed": elapsed,
			"sql":     sql,
			"rows":    rows,
		})
	}
}
