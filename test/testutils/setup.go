// test/testutils/setup.go
package testutils

import (
	"log"
	"os"
	"payslip-generation-system/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupTestDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		log.Fatalf("TEST_DATABASE_DSN not set: %v", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	DB = db
	seedTestAdmin()
	seedTestUser()
}

func seedTestUser() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	DB.Create(&model.User{
		ID:           uuid.New(),
		Username:     "employee999",
		PasswordHash: string(hash),
		Role:         "employee",
		Salary:       7000000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
}

func seedTestAdmin() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	DB.Create(&model.User{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
		Salary:       7000000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
}
