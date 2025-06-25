package main

import (
	"log"
	"os"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/middleware"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.POST("/login", handler.LoginHandler(db))

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.POST("/attendance-period", handler.CreateAttendancePeriodHandler(db))

	log.Println("Server running on :8081")
	r.Run(":8081")
}
