package repository

import (
	"payslip-generation-system/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	FindAttendancePeriod(periodID uuid.UUID) (bool, error)
	IsPayrollRun(periodID uuid.UUID) (bool, error)
	GetAttendances(periodID uuid.UUID) ([]model.Attendance, error)
	GetOvertimes(periodID uuid.UUID) ([]model.Overtime, error)
	GetReimbursements(periodID uuid.UUID) ([]model.Reimbursement, error)
	GetUserSalary(userIDs []uuid.UUID) ([]model.User, error)
	CreateAuditLog(log *model.AuditLog) error
	CreatePayroll(payroll *model.Payroll) error
	CreatePayslip(payslip *model.Payslip) error
}

type PayrollRepositoryImpl struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &PayrollRepositoryImpl{db: db}
}

func (pr *PayrollRepositoryImpl) FindAttendancePeriod(periodID uuid.UUID) (bool, error) {
	var found int64
	err := pr.db.Model(&model.Payroll{}).Where("period_id = ?", periodID).Count(&found).Error
	return found > 0, err
}

func (pr *PayrollRepositoryImpl) IsPayrollRun(periodID uuid.UUID) (bool, error) {
	var count int64
	err := pr.db.Model(&model.Payroll{}).Where("period_id = ?", periodID).Count(&count).Error
	return count > 0, err
}

func (pr *PayrollRepositoryImpl) GetAttendances(periodID uuid.UUID) ([]model.Attendance, error) {
	var result []model.Attendance
	err := pr.db.Raw(`
		SELECT a.* FROM attendances a
		JOIN attendance_periods p ON a.date BETWEEN p.start_date AND p.end_date
		WHERE p.id = ?
	`, periodID).Scan(&result).Error
	return result, err
}

func (pr *PayrollRepositoryImpl) GetOvertimes(periodID uuid.UUID) ([]model.Overtime, error) {
	var result []model.Overtime
	err := pr.db.Raw(`
		SELECT o.* FROM overtimes o
		JOIN attendance_periods p ON o.date BETWEEN p.start_date AND p.end_date
		WHERE p.id = ?
	`, periodID).Scan(&result).Error
	return result, err
}

func (pr *PayrollRepositoryImpl) GetReimbursements(periodID uuid.UUID) ([]model.Reimbursement, error) {
	var result []model.Reimbursement
	err := pr.db.Raw(`
		SELECT r.* FROM reimbursements r
		JOIN attendance_periods p ON r.created_at::date BETWEEN p.start_date AND p.end_date
		WHERE p.id = ?
	`, periodID).Scan(&result).Error
	return result, err
}

func (pr *PayrollRepositoryImpl) GetUserSalary(userIDs []uuid.UUID) ([]model.User, error) {
	var users []model.User
	if err := pr.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (pr *PayrollRepositoryImpl) CreateAuditLog(log *model.AuditLog) error {
	return pr.db.Create(&log).Error
}

func (pr *PayrollRepositoryImpl) CreatePayroll(payroll *model.Payroll) error {
	payroll.CreatedAt = time.Now()
	return pr.db.Create(&payroll).Error
}

func (pr *PayrollRepositoryImpl) CreatePayslip(payslip *model.Payslip) error {
	return pr.db.Create(&payslip).Error
}
