package repository

import (
	"admin-api/db"
	"admin-api/entity"

	"gorm.io/gorm"
)

type (
	AdminUserRepository interface {
		FindByEmail(email string) (*entity.AdminUser, error)
		Create(adminUser *entity.AdminUser) error
		UpdateRefreshToken(id uint, token string) error
		FindByRefreshToken(refreshToken string) (*entity.AdminUser, error)
		GetAllAdminEmail() ([]string, error)
	}
	adminUserRepository struct {
		db *gorm.DB
	}
)

func NewAdminUserRepository() AdminUserRepository {
	return &adminUserRepository{
		db: db.GetDB(),
	}
}

func (repo *adminUserRepository) FindByEmail(email string) (*entity.AdminUser, error) {
	var user entity.AdminUser
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *adminUserRepository) Create(adminUser *entity.AdminUser) error {
	return repo.db.Create(adminUser).Error
}

func (repo *adminUserRepository) UpdateRefreshToken(id uint, token string) error {
	return repo.db.Model(&entity.AdminUser{}).Where("id = ?", id).Update("RefreshToken", token).Error
}

func (repo *adminUserRepository) FindByRefreshToken(refreshToken string) (*entity.AdminUser, error) {
	var user entity.AdminUser
	if err := repo.db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *adminUserRepository) GetAllAdminEmail() ([]string, error) {
	var emails []string
	if err := repo.db.Model(&entity.AdminUser{}).Pluck("email", &emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}
