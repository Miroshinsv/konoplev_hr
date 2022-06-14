package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/meBazil/hr-mvp/internal/acl"
	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	helper "github.com/meBazil/hr-mvp/internal/rest/helpers"
)

func ACL(aclService *acl.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if helper.IsPublicRoute(mux.CurrentRoute(r)) {
				next.ServeHTTP(w, r)
				return
			}

			profile := r.Context().Value(ContextUserField)
			if profile == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !aclService.Enforce(profile.(*ptypes.Profile).Roles, mux.CurrentRoute(r).GetName(), r.Method) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
