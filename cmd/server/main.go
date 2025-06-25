package main

import (
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/middleware"

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

	// admin router
	adminMux := http.NewServeMux()
	adminMux.Handle("/attendance-period", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateAttendancePeriodHandler(db))))

	http.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	// employee router
	employeeMux := http.NewServeMux()
	employeeMux.Handle("/attendance", middleware.AuthMiddleware(http.HandlerFunc(handler.SubmitAttendanceHanlder(db))))
	employeeMux.Handle("/overtime", middleware.AuthMiddleware(http.HandlerFunc(handler.SubmitOvertimeHandler(db))))

	http.Handle("/employee/", http.StripPrefix("/employee", employeeMux))

	log.Println("Server running on :8081")
	http.ListenAndServe(":8081", nil)
}
