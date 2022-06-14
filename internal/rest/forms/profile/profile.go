package profile

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type UpdateRequest struct {
	Name       *string `json:"name"`
	MiddleName *string `json:"middle_name"`
	SureName   *string `json:"sure_name"`
}

type UpdateForm struct {
	Name       *string
	MiddleName *string
	SureName   *string
}

func NewUpdateForm() *UpdateForm {
	return &UpdateForm{}
}

func (f *UpdateForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *UpdateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	f.validateAndSetName(request)
	f.validateAndSetMiddleName(request)
	f.validateAndSetSureName(request)

	return f, nil
}

func (f *UpdateForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        f.Name,
		"middle_name": f.MiddleName,
		"sure_name":   f.SureName,
	}
}

func (f *UpdateForm) validateAndSetName(request *UpdateRequest) {
	if request.Name == nil {
		return
	}

	f.Name = request.Name
}

func (f *UpdateForm) validateAndSetMiddleName(request *UpdateRequest) {
	if request.MiddleName == nil {
		return
	}

	f.MiddleName = request.MiddleName
}

func (f *UpdateForm) validateAndSetSureName(request *UpdateRequest) {
	if request.SureName == nil {
		return
	}

	f.SureName = request.SureName
}
