package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/auth"
	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type IValidator interface {
	ValidateJWT(tokenString string) (*auth.Claims, error)
}

type RefreshRequest struct {
	Refresh string `json:"refresh"`
}

type RefreshForm struct {
	jwtValidator IValidator

	Claims *auth.Claims
}

func NewRefreshForm(validator IValidator) *RefreshForm {
	return &RefreshForm{
		jwtValidator: validator,
	}
}

func (r *RefreshForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *RefreshRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	r.validateAndSetRefresh(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return r, nil
}

func (r *RefreshForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"claims": r.Claims,
	}
}

func (r *RefreshForm) validateAndSetRefresh(request *RefreshRequest, errors map[string]response.ErrorMessage) {
	if request.Refresh == "" {
		errors["refresh"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	claims, err := r.jwtValidator.ValidateJWT(request.Refresh)
	if err != nil || claims.Type != auth.RefreshToken {
		errors["refresh"] = response.ErrorMessage{
			Code:    response.WrongValue,
			Message: "wrong refresh token",
		}

		return
	}

	r.Claims = claims
}
