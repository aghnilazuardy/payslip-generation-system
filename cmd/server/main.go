package main

import (
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/handler"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using OS environment")
	}

	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	http.HandleFunc("/login", handler.LoginHandler(db))

	// admin := r.Group("/admin")
	// admin.Use(middleware.AuthMiddleware())
	// admin.POST("/attendance-period", handler.CreateAttendancePeriodHandler(db))

	// employee := r.Group("/employee")
	// employee.Use(middleware.AuthMiddleware())
	// employee.POST("/attendance", handler.SubmitAttendanceHanlder(db))
	// employee.POST("/overtime", handler.SubmitOvertimeHandler(db))
	// employee.POST("/reimbursement", handler.SubmitReimbursementHandler(db))

	log.Println("Server running on :8081")
	http.ListenAndServe(":8081", nil)
}
