package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type ChangeEmailRequest struct {
	EmailNew string `json:"email_new"`
	Password string `json:"password"`
}

type ChangeEmailForm struct {
	EmailNew string
	Password string
}

func NewChangeEmailForm() *ChangeEmailForm {
	return &ChangeEmailForm{}
}

func (r *ChangeEmailForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *ChangeEmailRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	r.validateAndSetEmailNew(request, errors)
	r.validateAndSetPassword(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return r, nil
}

func (r *ChangeEmailForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"email_new": r.EmailNew,
		"password":  r.Password,
	}
}

func (r *ChangeEmailForm) validateAndSetEmailNew(request *ChangeEmailRequest, errors map[string]response.ErrorMessage) {
	if request.EmailNew == "" {
		errors["email_new"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.EmailNew = request.EmailNew
}

func (r *ChangeEmailForm) validateAndSetPassword(request *ChangeEmailRequest, errors map[string]response.ErrorMessage) {
	if request.Password == "" {
		errors["password"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Password = request.Password
}
