package profile

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/pkg/sqlutil"
)

var (
	ErrInactiveProfile = errors.New("user is not active")
)

type profileRepo interface {
	Create(profile *types.Profile) error
	Update(profile *types.Profile) error
	Get(ID uint) (*types.Profile, error)
	GetByEmail(email string) (*types.Profile, error)
	GetByMobile(mobile string) (*types.Profile, error)
	Delete(profile *types.Profile)
}

type Service struct {
	profileRepo profileRepo
}

func NewService(r profileRepo) (*Service, error) {
	service := Service{
		profileRepo: r,
	}

	return &service, nil
}

func (s *Service) RegisterUser(email string, mobile string, pwd string) (*types.Profile, error) {
	hashed, err := s.GeneratePasswordHash(pwd)
	if err != nil {
		return nil, err
	}

	prof := &types.Profile{
		Email:    email,
		Mobile:   mobile,
		Password: hashed,
		IsActive: true,
		Roles:    sqlutil.StringArray{string(types.RoleUser)},
	}

	if err := s.registerUser(prof); err != nil {
		return nil, err
	}

	return prof, nil
}

func (s *Service) registerUser(prof *types.Profile) error {
	if err := s.profileRepo.Create(prof); err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangePassword(profile *types.Profile, pwd string) error {
	hashed, err := s.GeneratePasswordHash(pwd)
	if err != nil {
		return err
	}

	profile.Password = hashed

	err = s.profileRepo.Update(profile)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Get(id uint) (*types.Profile, error) {
	profile, err := s.profileRepo.Get(id)
	if err != nil {
		return nil, err
	}

	if !profile.IsActive {
		return nil, ErrInactiveProfile
	}

	return profile, nil
}

func (s *Service) GetByEmail(email string) (*types.Profile, error) {
	profile, err := s.profileRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if !profile.IsActive {
		return nil, ErrInactiveProfile
	}

	return profile, nil
}

func (s *Service) GetByMobile(mobile string) (*types.Profile, error) {
	profile, err := s.profileRepo.GetByMobile(mobile)
	if err != nil {
		return nil, err
	}

	if !profile.IsActive {
		return nil, ErrInactiveProfile
	}

	return profile, nil
}

func (s *Service) Block(profile *types.Profile) error {
	profile.IsActive = false

	return s.profileRepo.Update(profile)
}

func (s *Service) Update(profile *types.Profile) error {
	return s.profileRepo.Update(profile)
}

func (s *Service) Delete(profile *types.Profile) {
	s.profileRepo.Delete(profile)
}

func (s *Service) GeneratePasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *Service) IsPasswordCorrect(prof *types.Profile, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(prof.Password), []byte(pwd))

	return err == nil
}
