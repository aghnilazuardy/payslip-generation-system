package main

import (
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/repository"

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

	// authorization route
	userRepo := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepo)

	http.HandleFunc("/login", authHandler.LoginHandler())

	// admin route
	adminRepo := repository.NewAdminRepository(db)
	adminHandler := handler.NewAdminHandler(adminRepo)

	adminMux := http.NewServeMux()
	adminMux.Handle("/attendance-period", middleware.AuthMiddleware(http.HandlerFunc(adminHandler.CreateAttendancePeriodHandler())))
	adminMux.Handle("/payroll-run", middleware.AuthMiddleware(http.HandlerFunc(adminHandler.CreateAttendancePeriodHandler())))
	http.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	// employee route
	employeeRepo := repository.NewEmployeeRepository(db)

	employeeHandler := handler.NewEmployeeHandler(employeeRepo)

	employeeMux := http.NewServeMux()
	employeeMux.Handle("/attendance", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.SubmitAttendanceHanlder())))
	employeeMux.Handle("/overtime", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.SubmitOvertimeHandler())))
	employeeMux.Handle("/reimbursement", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.SubmitReimbursementHandler())))
	http.Handle("/employee/", http.StripPrefix("/employee", employeeMux))

	log.Println("Server running on :8081")
	http.ListenAndServe(":8081", nil)
}
