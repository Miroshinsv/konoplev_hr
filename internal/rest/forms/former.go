package forms

import (
	"net/http"

	"github.com/meBazil/hr-mvp/internal/rest/response"
)

type Former interface {
	ParseAndValidate(request *http.Request) (Former, response.Error)
	ConvertToMap() map[string]interface{}
}
