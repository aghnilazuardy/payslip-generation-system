package repository

import (
	"payslip-generation-system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	SaveAttendance(attendance *model.Attendance) error
	SaveOvertime(overtime *model.Overtime) error
	SaveReimbursement(reimbursement *model.Reimbursement) error
	GetPayslip(userID, payrollID uuid.UUID) (*model.Payslip, error)
}

type EmployeeRepositoryImpl struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{db: db}
}

func (er *EmployeeRepositoryImpl) SaveAttendance(attendance *model.Attendance) error {
	return er.db.Create(&attendance).Error
}

func (er *EmployeeRepositoryImpl) SaveOvertime(overtime *model.Overtime) error {
	return er.db.Create(&overtime).Error
}

func (er *EmployeeRepositoryImpl) SaveReimbursement(reimbursement *model.Reimbursement) error {
	return er.db.Create(&reimbursement).Error
}

func (er *EmployeeRepositoryImpl) GetPayslip(userID, payrollID uuid.UUID) (*model.Payslip, error) {
	var result model.Payslip
	if err := er.db.Where("payroll_id = ? AND user_id = ?", payrollID, userID).Find(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
