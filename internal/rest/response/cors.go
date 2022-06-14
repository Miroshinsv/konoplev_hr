package response

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
)

func ConfigureCorsHandler(router *mux.Router) http.Handler {
	handlerMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	})
	handlerCredentials := handlers.AllowCredentials()
	handlerAllowedHeaders := handlers.AllowedHeaders([]string{
		HeaderContentType,
		HeaderAuthorization,
	})
	handlerExposedHeaders := handlers.ExposedHeaders([]string{
		HeaderTotalCount,
		HeaderCurrentOffset,
		HeaderLimit,
	})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	return handlers.CORS(handlerMethods, handlerCredentials, handlerAllowedHeaders, handlerExposedHeaders, allowedOrigins)(router)
}
