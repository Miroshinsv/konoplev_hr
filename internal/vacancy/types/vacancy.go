package types

import (
	"gorm.io/gorm"

	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
)

type Vacancy struct {
	gorm.Model

	Title       string
	Address     string
	Lat         float64
	Long        float64
	Description string

	ProfileID uint
	Profile   *ptypes.Profile
}
