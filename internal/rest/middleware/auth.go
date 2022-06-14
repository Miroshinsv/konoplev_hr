package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/meBazil/hr-mvp/internal/auth"
	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	helper "github.com/meBazil/hr-mvp/internal/rest/helpers"
)

const (
	ContextUserField ContextField = "__user"
)

type ContextField string

type AuthService interface {
	AuthByToken(token string) (*ptypes.Profile, error)
}

func Auth(service AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if helper.IsPublicRoute(mux.CurrentRoute(r)) {
				next.ServeHTTP(w, r)
				return
			}

			token := helper.ExtractTokenFromHeaders(r.Header)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			prof, err := service.AuthByToken(token)
			if err != nil {
				switch {
				case errors.Is(err, helper.ErrInvalidToken),
					errors.Is(err, helper.ErrInvalidTokenVersion),
					errors.Is(err, helper.ErrInvalidTokenType):
					w.WriteHeader(http.StatusUnauthorized)
					return
				case errors.Is(err, auth.ErrInActiveProfile):
					w.WriteHeader(http.StatusForbidden)
					return
				default:
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ContextUserField, prof)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
