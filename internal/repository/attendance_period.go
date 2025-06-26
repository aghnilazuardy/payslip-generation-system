package repository

import (
	"payslip-generation-system/internal/model"

	"gorm.io/gorm"
)

type AttendancePeriodRepository interface {
	SaveAttendancePeriod(attendancePeriod *model.AttendancePeriod) error
}

type AttendancePeriodRepositoryImpl struct {
	db *gorm.DB
}

func NewAttendancePeriodRepository(db *gorm.DB) AttendancePeriodRepository {
	return &AttendancePeriodRepositoryImpl{db: db}
}

func (apr *AttendancePeriodRepositoryImpl) SaveAttendancePeriod(attendancePeriod *model.AttendancePeriod) error {
	return apr.db.Create(&attendancePeriod).Error
}
