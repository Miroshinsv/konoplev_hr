package profile

import (
	"gorm.io/gorm"

	"github.com/meBazil/hr-mvp/internal/profile/types"
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

func (r Repository) Create(profile *types.Profile) error {
	return r.conn.GetWriteConnection().Create(&profile).Error
}

func (r Repository) Update(profile *types.Profile) error {
	return r.conn.GetWriteConnection().Save(&profile).Error
}

func (r Repository) Get(id uint) (*types.Profile, error) {
	var result *types.Profile
	err := r.conn.GetReadConnection().
		Where("id = ?", id).
		Take(&result).
		Error

	return result, err
}

func (r Repository) GetByEmail(email string) (*types.Profile, error) {
	var result *types.Profile
	err := r.conn.GetReadConnection().Where("email = ?", email).Take(&result).Error

	return result, err
}

func (r Repository) GetByMobile(mobile string) (*types.Profile, error) {
	var result *types.Profile
	err := r.conn.GetReadConnection().Where("mobile = ?", mobile).Take(&result).Error

	return result, err
}

func (r Repository) Delete(profile *types.Profile) {
	r.conn.GetWriteConnection().Delete(&types.Profile{}, profile.ID)
}
