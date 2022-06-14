package vacancy

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type CreateRequest struct {
	Title       string  `json:"title"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}

type CreateForm struct {
	Title       string  `json:"title"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}

func NewCreateForm() *CreateForm {
	return &CreateForm{}
}

func (r *CreateForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *CreateRequest
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

func (r *CreateForm) parseAndValidateTitle(request *CreateRequest, errs map[string]response.ErrorMessage) {
	if request.Title == "" {
		errs["title"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Title = request.Title
}

func (r *CreateForm) parseAndValidateAddress(request *CreateRequest, errs map[string]response.ErrorMessage) {
	if request.Address == "" {
		errs["address"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Address = request.Address
}

func (r *CreateForm) parseAndValidateDescription(request *CreateRequest, errs map[string]response.ErrorMessage) {
	if request.Description == "" {
		errs["description"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Description = request.Description
}

func (r *CreateForm) parseAndValidateLat(request *CreateRequest, errs map[string]response.ErrorMessage) {
	if request.Lat == 0 {
		errs["lat"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Lat = request.Lat
}

func (r *CreateForm) parseAndValidateLong(request *CreateRequest, errs map[string]response.ErrorMessage) {
	if request.Long == 0 {
		errs["long"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	r.Long = request.Long
}

func (r *CreateForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       r.Title,
		"address":     r.Address,
		"description": r.Description,
		"lat":         r.Lat,
		"long":        r.Long,
	}
}
