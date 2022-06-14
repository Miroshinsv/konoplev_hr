package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

func Timeout(dt time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			details, _ := json.Marshal(map[string]string{
				"message": "timeout",
			})

			h := PanicReportTimeoutHandler(next, dt, string(details))

			h.ServeHTTP(w, r)
		})
	}
}

// PanicReportTimeoutHandler replaces http.TimeoutHandler with PanicReportTimeoutHandler
func PanicReportTimeoutHandler(h http.Handler, dt time.Duration, msg string) http.Handler {
	return http.TimeoutHandler(&panicReporterHandler{handler: h}, dt, msg)
}

type panicReporterHandler struct {
	handler http.Handler
}

func (h *panicReporterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		var message string
		switch v := rec.(type) {
		case string:
			message = v
		case error:
			message = v.Error()
		default:
			message = "unknown error"
		}

		body, _ := io.ReadAll(r.Body)

		logger.Error(message, nil, map[string]interface{}{
			"request": fmt.Sprintf("%s %s?%s", r.Method, r.URL.String(), r.URL.Query().Encode()),
			"body":    string(body),
			"stack":   string(debug.Stack()),
		})

		w.WriteHeader(http.StatusInternalServerError)
	}()

	h.handler.ServeHTTP(w, r)
}
