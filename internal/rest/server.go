package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/meBazil/hr-mvp/internal/acl"
	"github.com/meBazil/hr-mvp/internal/auth"
	"github.com/meBazil/hr-mvp/internal/config"
	"github.com/meBazil/hr-mvp/internal/rest/handlers"
	"github.com/meBazil/hr-mvp/internal/rest/middleware"
	"github.com/meBazil/hr-mvp/internal/rest/response"
)

type mustHaveLimitedByTimeout interface {
	LimitedByTimeout() bool
}

func NewRestServer(cfg config.REST, authService *auth.Service, aclService *acl.Service, apiHandlers []handlers.APIHandler) *http.Server {
	handler := mux.NewRouter()
	handler.Use(
		middleware.Panic,
		middleware.RateLimitByUser(middleware.DefaultRateLimit),
		middleware.Auth(authService),
		middleware.ACL(aclService),
	)

	baseRouter := handler.PathPrefix(fmt.Sprintf("/%s", cfg.APIVersion)).Subrouter()
	baseRouter.Use(middleware.Timeout(cfg.HandleTimeout))
	baseRouter.Use(middleware.JSON)

	routerWithoutTimeout := handler.PathPrefix(fmt.Sprintf("/%s", cfg.APIVersion)).Subrouter()
	routerWithoutTimeout.Use(middleware.JSON)

	for _, h := range apiHandlers {
		// By default, we should use the timeout middleware.
		// But it can be disabled for some handlers.
		// To disable the timeout middleware a handler must implement LimitedByTimeout that returns false.
		if handler, ok := h.(mustHaveLimitedByTimeout); ok && !handler.LimitedByTimeout() {
			h.EnrichRoutes(routerWithoutTimeout)
		} else {
			h.EnrichRoutes(baseRouter)
		}
	}

	return &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      response.ConfigureCorsHandler(handler),
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}
}
