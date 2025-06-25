package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"type:text;not null"`
	Salary       int       `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AttendancePeriod struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	CreatedBy uuid.UUID
	RequestIP string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Attendance struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID
	Date      time.Time `gorm:"type:date"`
	CreatedBy uuid.UUID
	RequestIP string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Overtime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID
	Date      time.Time `gorm:"type:date"`
	Hours     int
	CreatedBy uuid.UUID
	RequestIP string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reimbursement struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID
	Amount      int
	Description string
	CreatedBy   uuid.UUID
	RequestIP   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Payroll struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PeriodID  uuid.UUID
	CreatedBy uuid.UUID
	RequestIP string
	CreatedAt time.Time
}

type Payslip struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PayrollID          uuid.UUID
	UserID             uuid.UUID
	BaseSalary         int
	AttendaceDays      int
	ProratedSalary     int
	OvertimeHours      int
	OvertimePay        int
	ReimbursementTotal int
	TakeHomePay        int
}

type AuditLog struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TableName   string
	RecordID    uuid.UUID
	Action      string
	PerformedBy uuid.UUID
	RequestIP   string
	RequestID   string
	Timestamp   time.Time
}
