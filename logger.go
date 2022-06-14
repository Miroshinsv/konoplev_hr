package main

import (
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
	"github.com/meBazil/hr-mvp/pkg/process-manager"
)

type PMLogger struct {
	Logger *logger.Logger
}

func (l *PMLogger) Info(msg string, fields ...process.LogFields) {
	l.Logger.Info(msg, convertFields(fields)...)
}

func (l *PMLogger) Error(msg string, err error, fields ...process.LogFields) {
	l.Logger.Error(msg, err, convertFields(fields)...)
}

func convertFields(fields []process.LogFields) []logger.Fields {
	result := make([]logger.Fields, 0, len(fields))

	for _, value := range fields {
		result = append(result, logger.Fields(value))
	}

	return result
}
