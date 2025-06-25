package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OvertimeRequest struct {
	Hours int `json:"hours"`
}

func SubmitAttendanceHanlder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// to check user role
		if middleware.GetUserRole(c) != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// to check is today weekend or weekday
		today := time.Now().Truncate(24 * time.Hour)
		weekday := today.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot submit on weekend"})
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		log.Printf("Attendance submission: user_id=%s date=%v", userID.String(), today.Format("2006-01-02"))

		attendance := model.Attendance{
			UserID:    userID,
			Date:      today,
			CreatedBy: userID,
			RequestIP: c.ClientIP(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&attendance).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				c.JSON(http.StatusConflict, gin.H{"error": "already submitted today"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Attendance successfully submitted"})
	}
}

func SubmitOvertimeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// to check user role
		if middleware.GetUserRole(c) != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		var req OvertimeRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Hours < 1 || req.Hours > 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hours"})
		}

		userID, err := uuid.Parse(middleware.GetUserID(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		hourNow := time.Now().Hour()
		fmt.Printf("%v\n", hourNow)
		endOfWorkingHour, _ := strconv.Atoi(os.Getenv("WORK_HOUR_END"))
		fmt.Printf("%v-%v\n", hourNow, endOfWorkingHour)
		if hourNow < endOfWorkingHour {
			c.JSON(http.StatusBadRequest, gin.H{"error": "overtime can only be submitted after 5PM"})
			return
		}

		today := time.Now().Truncate(24 * time.Hour)

		overtime := model.Overtime{
			UserID:    userID,
			Date:      today,
			Hours:     req.Hours,
			CreatedBy: userID,
			RequestIP: c.ClientIP(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&overtime).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				c.JSON(http.StatusConflict, gin.H{"error": "already submitted today"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Overtime successfully submitted"})
	}
}
