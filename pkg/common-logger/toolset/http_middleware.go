package toolset

import (
	"context"
	"net/http"

	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

const HeaderTraceKey string = "X-Trace-ID"

func HTTPTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(HeaderTraceKey)
		if val == "" {
			val = logger.GetNewTraceID()
		}

		r = r.WithContext(
			context.WithValue(r.Context(), logger.TraceIDKey, val),
		)

		next.ServeHTTP(w, r)
	})
}
