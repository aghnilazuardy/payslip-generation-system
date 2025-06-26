package repository

import (
	"payslip-generation-system/internal/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	SaveAttendancePeriod(attendancePeriod *model.AttendancePeriod) error
}

type AdminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &AdminRepositoryImpl{db: db}
}

func (ar *AdminRepositoryImpl) SaveAttendancePeriod(attendancePeriod *model.AttendancePeriod) error {
	return ar.db.Create(&attendancePeriod).Error
}
