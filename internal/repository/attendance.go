package repository

import (
	"payslip-generation-system/internal/model"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	SaveAttendance(attendancePeriod *model.Attendance) error
}

type AttendanceRepositoryImpl struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{db: db}
}

func (ar *AttendanceRepositoryImpl) SaveAttendance(attendance *model.Attendance) error {
	return ar.db.Create(&attendance).Error
}
