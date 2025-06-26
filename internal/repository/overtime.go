package repository

import (
	"payslip-generation-system/internal/model"

	"gorm.io/gorm"
)

type OvertimeRepository interface {
	SaveOvertime(overtime *model.Overtime) error
}

type OvertimeRepositoryImpl struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) OvertimeRepository {
	return &OvertimeRepositoryImpl{db: db}
}

func (or *OvertimeRepositoryImpl) SaveOvertime(overtime *model.Overtime) error {
	return or.db.Create(&overtime).Error
}
