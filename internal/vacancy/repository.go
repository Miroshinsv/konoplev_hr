package vacancy

import (
	"gorm.io/gorm"

	"github.com/meBazil/hr-mvp/internal/vacancy/types"
)

type connectionManager interface {
	GetReadConnection() *gorm.DB
	GetWriteConnection() *gorm.DB
}

type Repository struct {
	conn connectionManager
}

func NewRepository(conn connectionManager) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r Repository) Create(vacancy *types.Vacancy) error {
	return r.conn.GetWriteConnection().Create(&vacancy).Error
}

func (r Repository) Update(vacancy *types.Vacancy) error {
	return r.conn.GetWriteConnection().Save(&vacancy).Error
}

func (r Repository) Get(id uint) (*types.Vacancy, error) {
	var result *types.Vacancy
	err := r.conn.GetReadConnection().
		Preload("Profile").
		Where("id = ?", id).
		Take(&result).
		Error

	return result, err
}

func (r Repository) Delete(vacancy *types.Vacancy) {
	r.conn.GetWriteConnection().Delete(&types.Vacancy{}, vacancy.ID)
}
