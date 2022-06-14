package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Email    string
	Mobile   string
	Password string
}

func NewRegisterForm() *RegisterForm {
	return &RegisterForm{}
}

func (r *RegisterForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *RegisterRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	r.validateAndSetEmail(request, errors)
	r.validateAndSetPassword(request, errors)
	r.validateAndSetMobile(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return r, nil
}

func (r *RegisterForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":  r.Email,
		"mobile": r.Mobile,
	}
}

func (r *RegisterForm) validateAndSetEmail(request *RegisterRequest, errors map[string]response.ErrorMessage) {
	if request.Email == "" {
		errors["email"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Email = request.Email
}

func (r *RegisterForm) validateAndSetPassword(request *RegisterRequest, errors map[string]response.ErrorMessage) {
	if request.Password == "" {
		errors["password"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Password = request.Password
}

func (r *RegisterForm) validateAndSetMobile(request *RegisterRequest, errors map[string]response.ErrorMessage) {
	if request.Mobile == "" {
		errors["mobile"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	var mobileRx = regexp.MustCompile(`^(\+\d{1,3}[- ]?)?\d{10}$`)
	if len(mobileRx.FindStringIndex(request.Mobile)) == 0 {
		errors["mobile"] = response.ErrorMessage{
			Code:    response.WrongValue,
			Message: "wrong value",
		}

		return
	}

	r.Mobile = mobileRx.FindString(request.Mobile)
}
