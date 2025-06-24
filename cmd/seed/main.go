package main

import (
	"log"
	"os"
	"payslip-generation-system/internal/seed"

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

	if err := seed.SeedUsers(db); err != nil {
		log.Fatal("failed to seed users: ", err)
	}

	log.Println("Successfully seeded 100 employees and 1 admin user")
}
