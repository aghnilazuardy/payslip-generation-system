package repository

import (
	"payslip-generation-system/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
