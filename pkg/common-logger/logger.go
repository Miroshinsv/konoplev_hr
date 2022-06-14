package logger

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

const traceIdFieldName = "trace_id"

type Fields map[string]interface{}

type Wrapper interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Fatal(msg string, err error, fields map[string]interface{})
}

type Logger struct {
	Log     Wrapper
	Context context.Context
}

func NewLogger(ctx context.Context, log Wrapper) Logger {
	return Logger{
		Log:     log,
		Context: ctx,
	}
}

func (l Logger) Debug(msg string, fields ...Fields) {
	l.Log.Debug(msg, l.prepareFields(fields...))
}

func (l Logger) Info(msg string, fields ...Fields) {
	l.Log.Info(msg, l.prepareFields(fields...))
}

func (l Logger) Warn(msg string, fields ...Fields) {
	l.Log.Warn(msg, l.prepareFields(fields...))
}

func (l Logger) Error(msg string, err error, fields ...Fields) {
	l.Log.Error(msg, err, l.prepareFields(fields...))
}

func (l Logger) Fatal(msg string, err error, fields ...Fields) {
	l.Log.Fatal(msg, err, l.prepareFields(fields...))
}

func (l Logger) LogWithTimeout(infoMessage string, deadlineError error, td time.Duration, fields ...Fields) context.CancelFunc {
	l.Log.Info(infoMessage, l.prepareFields(fields...))

	wd, cancel := context.WithTimeout(l.Context, td)
	go func() {
		// waiting for timeout or cancel
		<-wd.Done()

		if errors.Is(wd.Err(), context.DeadlineExceeded) {
			l.Log.Error("failed function with timeout", deadlineError, l.prepareFields(fields...))

			return
		}

		l.Log.Info("execution function with timeout succeed", l.prepareFields(fields...))
	}()

	return cancel
}

func (l Logger) prepareFields(fields ...Fields) Fields {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = make(Fields)
	}

	if f == nil {
		f = make(Fields)
	}

	traceID := l.Context.Value(TraceIDKey)
	if traceID != nil {
		f[traceIdFieldName] = traceID
	}

	return f
}
