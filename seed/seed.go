package seed

import (
	"fmt"
	"math/rand"
	"payslip-generation-system/internal/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	password := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	for i := 1; i <= 100; i++ {
		employee := model.User{
			ID:           uuid.New(),
			Username:     fmt.Sprintf("employee%03d", i),
			PasswordHash: string(hash),
			Role:         "employee",
			Salary:       rand.Intn(5_000_000) + 5_000_000,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		db.Create(&employee)
	}

	admin := model.User{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
		Salary:       0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return db.Create(&admin).Error
}
