package vacancy

import (
	"github.com/meBazil/hr-mvp/internal/vacancy/types"
)

type vacancyRepo interface {
	Create(vacancy *types.Vacancy) error
	Update(vacancy *types.Vacancy) error
	Get(ID uint) (*types.Vacancy, error)
	Delete(vacancy *types.Vacancy)
}

type Service struct {
	vacancyRepo vacancyRepo
}

func NewService(r vacancyRepo) (*Service, error) {
	service := Service{
		vacancyRepo: r,
	}

	return &service, nil
}

func (s *Service) Get(id uint) (*types.Vacancy, error) {
	profile, err := s.vacancyRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) Update(vacancy *types.Vacancy) error {
	return s.vacancyRepo.Update(vacancy)
}

func (s *Service) Delete(vacancy *types.Vacancy) {
	s.vacancyRepo.Delete(vacancy)
}

func (s *Service) Create(vacancy *types.Vacancy) (*types.Vacancy, error) {
	err := s.vacancyRepo.Create(vacancy)

	return vacancy, err
}
