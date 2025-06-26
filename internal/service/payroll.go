package service

import (
	"errors"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PayrollService interface {
	ProcessPayroll(periodID, createdBy uuid.UUID, ip, requestID string) error
}

type PayrollServiceImpl struct {
	PayrollRepo repository.PayrollRepository
}

func NewPayrollService(repo repository.PayrollRepository) PayrollService {
	return &PayrollServiceImpl{PayrollRepo: repo}
}

func (s *PayrollServiceImpl) ProcessPayroll(periodID, createdBy uuid.UUID, ip, requestID string) error {
	// check if attendance period is exist
	found, err := s.PayrollRepo.FindAttendancePeriod(periodID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("attendance period not found")
	}

	// check if payroll already processed
	exists, err := s.PayrollRepo.IsPayrollRun(periodID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("payroll already processed for this period")
	}

	// get all employees attendance in given period
	attendances, err := s.PayrollRepo.GetAttendances(periodID)
	if err != nil {
		return err
	}

	// get all employees overtime hours in given period
	overtimes, err := s.PayrollRepo.GetOvertimes(periodID)
	if err != nil {
		return err
	}

	// get all reimbursement of employee in given period
	reimbursements, err := s.PayrollRepo.GetReimbursements(periodID)
	if err != nil {
		return err
	}

	// aggregate data
	attendanceMap := map[uuid.UUID]int{}
	overtimeMap := map[uuid.UUID]int{}
	reimbursementMap := map[uuid.UUID]int{}
	baseSalaryMap := map[uuid.UUID]int{}
	uniqueUserIDs := map[uuid.UUID]bool{}

	// mapping the attendance of employee
	for _, a := range attendances {
		attendanceMap[a.UserID]++
		uniqueUserIDs[a.UserID] = true
	}

	// mapping the overtime of employee
	for _, o := range overtimes {
		overtimeMap[o.UserID] += o.Hours
		uniqueUserIDs[o.UserID] = true
	}

	// mapping the reimbursement of employee
	for _, r := range reimbursements {
		reimbursementMap[r.UserID] += r.Amount
		uniqueUserIDs[r.UserID] = true
	}

	// mapping uniqe employee who did attendancy, overtime, and reimbursement
	userIDs := []uuid.UUID{}
	for id := range uniqueUserIDs {
		userIDs = append(userIDs, id)
	}

	// get base salary of each employee
	users, err := s.PayrollRepo.GetUserSalary(userIDs)
	if err != nil {
		return err
	}

	// mapping the base salary
	for _, u := range users {
		baseSalaryMap[u.ID] = u.Salary
	}

	// input the payroll
	payroll := &model.Payroll{
		ID:        uuid.New(),
		PeriodID:  periodID,
		CreatedBy: createdBy,
		RequestIP: ip,
		CreatedAt: time.Now(),
	}
	if err := s.PayrollRepo.CreatePayroll(payroll); err != nil {
		return err
	}

	// calculate overtime, reimburse, and prorated salary to input payslip of each employee
	for userID, days := range attendanceMap {
		salary := baseSalaryMap[userID]
		overtimeHrs := overtimeMap[userID]
		reimburse := reimbursementMap[userID]

		// assuming 20 working days per month
		prorated := (salary * days) / 20
		hourlyRate := salary / (20 * 8)
		overtimePay := 2 * hourlyRate * overtimeHrs
		total := prorated + overtimePay + reimburse

		p := &model.Payslip{
			ID:                 uuid.New(),
			PayrollID:          payroll.ID,
			UserID:             userID,
			BaseSalary:         salary,
			AttendanceDays:     days,
			ProratedSalary:     prorated,
			OvertimeHours:      overtimeHrs,
			OvertimePay:        overtimePay,
			ReimbursementTotal: reimburse,
			TakeHomePay:        total,
		}
		s.PayrollRepo.CreatePayslip(p)
	}

	// logging the process for audit purpose
	audit := model.AuditLog{
		ID:          uuid.New(),
		TableName:   "payrolls",
		RecordID:    payroll.ID,
		Action:      "CREATE",
		PerformedBy: createdBy,
		RequestIP:   ip,
		RequestID:   requestID,
		Timestamp:   time.Now(),
	}
	_ = s.PayrollRepo.CreateAuditLog(&audit).Error

	return nil
}
