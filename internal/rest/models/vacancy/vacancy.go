package vacancy

import (
	"time"

	"github.com/meBazil/hr-mvp/internal/rest/models/profile"
)

type Vacancy struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Address     string           `json:"address"`
	Lat         float64          `json:"lat"`
	Long        float64          `json:"long"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Profile     *profile.Profile `json:"profile"`
}
