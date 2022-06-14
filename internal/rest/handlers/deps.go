package handlers

import (
	"github.com/meBazil/hr-mvp/internal/auth"
	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/internal/vacancy/types"
)

type authService interface {
	AuthByEmail(email string, pwd string) (*ptypes.Profile, error)
	AuthByMobile(mobile string, pwd string) (*ptypes.Profile, error)
	Register(email string, mobile string, pwd string) (*ptypes.Profile, error)
	GenerateAccessToken(profile *ptypes.Profile) (string, error)
	GenerateRefreshToken(profile *ptypes.Profile) (string, error)
	GetValidator() *auth.Validator
	ChangePassword(profile *ptypes.Profile, pwd string, pwdNew string) error
	SetPassword(prof *ptypes.Profile, pwdNew string) error
	IsPasswordCorrect(prof *ptypes.Profile, pwd string) bool
	HasRole(prof *ptypes.Profile, role ptypes.Role) bool
}

type profileService interface {
	RegisterUser(email string, mobile string, pwd string) (*ptypes.Profile, error)
	ChangePassword(profile *ptypes.Profile, pwd string) error
	Get(id uint) (*ptypes.Profile, error)
	GetByEmail(email string) (*ptypes.Profile, error)
	GetByMobile(mobile string) (*ptypes.Profile, error)
	Block(profile *ptypes.Profile) error
	Update(profile *ptypes.Profile) error
	Delete(profile *ptypes.Profile)
	GeneratePasswordHash(pwd string) (string, error)
	IsPasswordCorrect(prof *ptypes.Profile, pwd string) bool
}

type vacancyService interface {
	Get(id uint) (*types.Vacancy, error)
	Update(vacancy *types.Vacancy) error
	Delete(vacancy *types.Vacancy)
	Create(vacancy *types.Vacancy) (*types.Vacancy, error)
}
