package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/meBazil/hr-mvp/internal/rest/middleware"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
	"github.com/meBazil/hr-mvp/pkg/process-manager"
)

func NewHealthCheckServer(listen, path string, handler http.Handler) *http.Server {
	router := mux.NewRouter()
	router.Use(middleware.Panic)
	router.Handle(path, handler)

	server := &http.Server{
		Addr:    listen,
		Handler: router,
	}

	return server
}

func DefaultHandler(manager *process.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"process_manager": manager.IsRunning(),
		}

		body, err := json.Marshal(resp)
		if err != nil {
			logger.Warn("unable to marshal health check")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}
