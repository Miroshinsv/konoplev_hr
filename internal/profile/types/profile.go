package types

import (
	"gorm.io/gorm"

	"github.com/meBazil/hr-mvp/pkg/sqlutil"
)

type Profile struct {
	gorm.Model

	Name       string
	MiddleName string
	SureName   string
	Email      string `gorm:"unique;index"`
	Mobile     string
	Password   string

	IsActive bool `gorm:"index"`

	Roles sqlutil.StringArray `sql:"type=json"`
}
