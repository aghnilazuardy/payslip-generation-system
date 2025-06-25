package handler

import (
	"net/http"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendancePeriodRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func CreateAttendancePeriodHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if middleware.GetUserRole(c) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		var req AttendancePeriodRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		startDate, err1 := time.Parse("2006-01-02", req.StartDate)
		endDate, err2 := time.Parse("2006-01-02", req.EndDate)
		if err1 != nil || err2 != nil || !startDate.Before(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		period := model.AttendancePeriod{
			StartDate: startDate,
			EndDate:   endDate,
			CreatedBy: userID,
			RequestIP: c.ClientIP(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&period).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Attendance successfully created"})
	}
}
