package wrapper

import (
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	log *logrus.Entry
}

func (s Logrus) Debug(msg string, fields map[string]interface{}) {
	s.log.WithFields(fields).Debug(msg)
}

func (s Logrus) Info(msg string, fields map[string]interface{}) {
	s.log.WithFields(fields).Info(msg)
}

func (s Logrus) Warn(msg string, fields map[string]interface{}) {
	s.log.WithFields(fields).Warn(msg)
}

func (s Logrus) Error(msg string, err error, fields map[string]interface{}) {
	l := s.log

	if fields != nil {
		l = l.WithFields(fields)
	}

	if err != nil {
		l = l.WithError(err)
	}

	l.Error(msg)
}

func (s Logrus) Fatal(msg string, err error, fields map[string]interface{}) {
	l := s.log

	if fields != nil {
		l = l.WithFields(fields)
	}

	if err != nil {
		l = l.WithError(err)
	}

	l.Fatal(msg)
}

func WrapLogrusEntry(l *logrus.Entry) Logrus {
	return Logrus{log: l}
}

func WrapLogrusLogger(l *logrus.Logger) Logrus {
	return Logrus{log: logrus.NewEntry(l)}
}
