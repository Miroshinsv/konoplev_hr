package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type ChangePasswordRequest struct {
	Password    string `json:"password"`
	PasswordNew string `json:"password_new"`
}

type ChangePasswordForm struct {
	Password    string
	PasswordNew string
}

func NewChangePasswordForm() *ChangePasswordForm {
	return &ChangePasswordForm{}
}

func (r *ChangePasswordForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *ChangePasswordRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	r.validateAndSetPassword(request, errors)
	r.validateAndSetPasswordNew(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return r, nil
}

func (r *ChangePasswordForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"password":     r.Password,
		"password_new": r.PasswordNew,
	}
}

func (r *ChangePasswordForm) validateAndSetPassword(request *ChangePasswordRequest, errors map[string]response.ErrorMessage) {
	if request.Password == "" {
		errors["password"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Password = request.Password
}

func (r *ChangePasswordForm) validateAndSetPasswordNew(request *ChangePasswordRequest, errors map[string]response.ErrorMessage) {
	if request.PasswordNew == "" {
		errors["password_new"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.PasswordNew = request.PasswordNew
}
