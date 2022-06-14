package vacancy

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type UpdateRequest struct {
	Title       *string  `json:"tile"`
	Address     *string  `json:"address"`
	Description *string  `json:"description"`
	Lat         *float64 `json:"lat"`
	Long        *float64 `json:"long"`
}

type UpdateForm struct {
	Title       *string  `json:"tile"`
	Address     *string  `json:"address"`
	Description *string  `json:"description"`
	Lat         *float64 `json:"lat"`
	Long        *float64 `json:"long"`
}

func NewUpdateForm() *UpdateForm {
	return &UpdateForm{}
}

func (r *UpdateForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
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

	errors := make(map[string]response.ErrorMessage)
	r.parseAndValidateTitle(request, errors)
	r.parseAndValidateAddress(request, errors)
	r.parseAndValidateDescription(request, errors)
	r.parseAndValidateLat(request, errors)
	r.parseAndValidateLong(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return r, nil
}

func (r *UpdateForm) parseAndValidateTitle(request *UpdateRequest, _ map[string]response.ErrorMessage) {
	r.Title = request.Title
}

func (r *UpdateForm) parseAndValidateAddress(request *UpdateRequest, _ map[string]response.ErrorMessage) {
	r.Address = request.Address
}

func (r *UpdateForm) parseAndValidateDescription(request *UpdateRequest, _ map[string]response.ErrorMessage) {
	r.Description = request.Description
}

func (r *UpdateForm) parseAndValidateLat(request *UpdateRequest, _ map[string]response.ErrorMessage) {
	r.Lat = request.Lat
}

func (r *UpdateForm) parseAndValidateLong(request *UpdateRequest, _ map[string]response.ErrorMessage) {
	r.Long = request.Long
}

func (r *UpdateForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       r.Title,
		"address":     r.Address,
		"description": r.Description,
		"lat":         r.Lat,
		"long":        r.Long,
	}
}
