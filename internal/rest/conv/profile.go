package conv

import (
	"strconv"

	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/internal/rest/models/profile"
)

func ConvertProfile(prof *ptypes.Profile) *profile.Profile {
	return &profile.Profile{
		ID:         strconv.Itoa(int(prof.ID)),
		Email:      prof.Email,
		Mobile:     prof.Mobile,
		IsActive:   prof.IsActive,
		Name:       prof.Name,
		MiddleName: prof.MiddleName,
		SureName:   prof.SureName,
		CreatedAt:  prof.CreatedAt,
		UpdatedAt:  prof.UpdatedAt,
		Roles:      prof.Roles,
	}
}
