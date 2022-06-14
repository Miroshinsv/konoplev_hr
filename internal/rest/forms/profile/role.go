package profile

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/internal/rest/forms"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

type UpdateRoleRequest struct {
	Roles []string `json:"roles"`
}

type UpdateRoleForm struct {
	Roles []string
}

func NewUpdateRoleForm() *UpdateRoleForm {
	return &UpdateRoleForm{}
}

func (f *UpdateRoleForm) ParseAndValidate(req *http.Request) (forms.Former, response.Error) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error("unable to read body", err)
		return nil, response.NewInternalError()
	}

	var request *UpdateRoleRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errs := make(map[string]response.ErrorMessage)
	f.validateAndSetRoles(request, errs)

	if len(errs) > 0 {
		return nil, response.NewValidationError(errs)
	}

	return f, nil
}

func (f *UpdateRoleForm) validateAndSetRoles(request *UpdateRoleRequest, errs map[string]response.ErrorMessage) {
	if len(request.Roles) == 0 {
		errs["roles"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}

		return
	}

	for _, r := range request.Roles {
		if !f.inRolesArray(r) {
			errs["roles"] = response.ErrorMessage{
				Code:    response.WrongValue,
				Message: "wrong value",
			}

			return
		}

		f.Roles = append(f.Roles, r)
	}
}

func (f *UpdateRoleForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"roles": f.Roles,
	}
}

func (f *UpdateRoleForm) inRolesArray(role string) bool {
	for _, ex := range []types.Role{types.RoleAdmin, types.RoleHR, types.RoleUser} {
		if strings.EqualFold(role, string(ex)) {
			return true
		}
	}

	return false
}
