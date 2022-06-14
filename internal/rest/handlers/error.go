package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/response"
)

func HandleError(err response.Error, w http.ResponseWriter) {
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(err.GetHTTPStatus())
	_ = json.NewEncoder(w).Encode(ParseError(err))
}

// ParseError determines the error type and creates a map with the error description.
func ParseError(err response.Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	switch e := err.(type) { // nolint:gocritic
	case *response.ValidationError:
		return map[string]interface{}{
			"errors":  e.Errors(),
			"message": e.PublicMessage(),
		}

	case *response.UnprocessableEntityError:
		return map[string]interface{}{
			"errors":  e.Errors(),
			"message": e.PublicMessage(),
		}

	case *response.RateLimitedError:
		return map[string]interface{}{
			"message":     e.PublicMessage(),
			"retry_after": e.RetryAfter,
		}

	default:
		return map[string]interface{}{
			"message": err.PublicMessage(),
		}
	}
}
