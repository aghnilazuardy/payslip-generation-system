package repository

import (
	"payslip-generation-system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	SaveAttendancePeriod(attendancePeriod *model.AttendancePeriod) error
	GetPayslipSummary(payrollID uuid.UUID) ([]model.EmployeePayslipSummary, error)
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

func (ar *AdminRepositoryImpl) GetPayslipSummary(payrollID uuid.UUID) ([]model.EmployeePayslipSummary, error) {
	var results []model.EmployeePayslipSummary
	err := ar.db.Raw(`
		SELECT p.user_id, u.username, p.take_home_pay
		FROM payslips p
		JOIN users u ON p.user_id = u.id
		WHERE p.payroll_id = ?`, payrollID).Scan(&results).Error

	return results, err
}
