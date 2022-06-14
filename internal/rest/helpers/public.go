package helper

import (
	"strings"

	"github.com/gorilla/mux"
)

const apiPublicRoutePrefix = "public_"

func IsPublicRoute(route *mux.Route) bool {
	return strings.Contains(route.GetName(), apiPublicRoutePrefix)
}
