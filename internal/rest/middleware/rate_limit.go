package middleware

import (
	"encoding/json"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"

	helper "github.com/meBazil/hr-mvp/internal/rest/helpers"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	"github.com/meBazil/hr-mvp/internal/rest/throttle"
)

// DefaultRateLimit is used by default for any non-public endpoint.
var DefaultRateLimit = throttle.LimitConfig{
	Every: throttle.RPS(10), //nolint:gomnd
	Burst: 15,               //nolint:gomnd
}

// RateLimitByUser intercepts authenticated requests, and does rate limit by `user.ID`.
func RateLimitByUser(conf throttle.LimitConfig) func(next http.Handler) http.Handler {
	limiter := throttle.NewLimiter(conf)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// no rate limit for public routes
			if helper.IsPublicRoute(mux.CurrentRoute(r)) {
				next.ServeHTTP(w, r)
				return
			}

			token := helper.ExtractTokenFromHeaders(r.Header)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			type Claims struct {
				jwt.StandardClaims

				UserID uint `json:"uid"`
			}

			parsed, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
			if err != nil {
				log.WithError(err).Error("parse jwt token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			claims, ok := parsed.Claims.(*Claims)
			if !ok {
				log.Warning("unable to parse correct JWT token claims")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ok, retryAfter := limiter.AllowInfo(claims.UserID)
			if !ok {
				seconds := int(math.Ceil(retryAfter.Seconds()))

				err := response.NewRateLimitedError(seconds)
				if err == nil {
					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				w.WriteHeader(err.GetHTTPStatus())
				_ = json.NewEncoder(w).Encode(err)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
