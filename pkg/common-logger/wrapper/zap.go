package wrapper

import (
	"go.uber.org/zap"
)

type Zap struct {
	log *zap.Logger
}

func (s Zap) Debug(msg string, fields map[string]interface{}) {
	s.log.Debug(msg, s.prepareFields(fields)...)
}

func (s Zap) Info(msg string, fields map[string]interface{}) {
	s.log.Info(msg, s.prepareFields(fields)...)
}

func (s Zap) Warn(msg string, fields map[string]interface{}) {
	s.log.Warn(msg, s.prepareFields(fields)...)
}

func (s Zap) Error(msg string, err error, fields map[string]interface{}) {
	zf := s.prepareFields(fields)
	zf = append(zf, zap.Error(err))

	s.log.Error(msg, zf...)
}

func (s Zap) Fatal(msg string, err error, fields map[string]interface{}) {
	zf := s.prepareFields(fields)
	zf = append(zf, zap.Error(err))

	s.log.Fatal(msg, zf...)
}

func (s Zap) prepareFields(fields map[string]interface{}) []zap.Field {
	zf := make([]zap.Field, len(fields))
	for key, val := range fields {
		zf = append(zf, zap.Any(key, val))
	}

	return zf
}

func WrapZap(l *zap.Logger) Zap {
	return Zap{log: l}
}
