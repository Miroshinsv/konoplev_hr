package conv

import (
	"strconv"

	"github.com/meBazil/hr-mvp/internal/rest/models/vacancy"
	vtypes "github.com/meBazil/hr-mvp/internal/vacancy/types"
)

func ConvertVacancy(vac *vtypes.Vacancy) *vacancy.Vacancy {
	return &vacancy.Vacancy{
		ID:          strconv.Itoa(int(vac.ID)),
		Title:       vac.Title,
		Address:     vac.Address,
		Lat:         vac.Lat,
		Long:        vac.Long,
		Description: vac.Description,
		CreatedAt:   vac.CreatedAt,
		UpdatedAt:   vac.UpdatedAt,
		Profile:     ConvertProfile(vac.Profile),
	}
}
