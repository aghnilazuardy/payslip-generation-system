package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9/helper"
	"gorm.io/gorm"
)

type OvertimeRequest struct {
	Hours int `json:"hours"`
}

type ReimbursementRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type EmployeeHandler struct {
	EmployeeRepo repository.EmployeeRepository
}

func NewEmployeeHandler(employeeRepo repository.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{EmployeeRepo: employeeRepo}
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

		saveErr := emh.EmployeeRepo.SaveAttendance(&attendance)
		if saveErr != nil {
			if strings.Contains(saveErr.Error(), "duplicate key") {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusConflict, "already submitted today", nil, nil))
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

		saveErr := emh.EmployeeRepo.SaveOvertime(&overtime)
		if saveErr != nil {
			if strings.Contains(saveErr.Error(), "duplicate key") {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusConflict, "already submitted today", nil, nil))
			} else {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "failed to submit overtime", nil, nil))
			}
			return
		}

		json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusCreated, "overtime submitted successfully", nil, nil))
	}
}

func (emh *EmployeeHandler) SubmitReimbursementHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil, nil))
			return
		}

		// to check user role
		if middleware.GetUserRole(r) != "employee" {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusForbidden, "forbidden", nil, nil))
			return
		}

		var req ReimbursementRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid request", nil, nil))
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "invalid user ID", nil, nil))
			return
		}

		hourNow := time.Now().Hour()
		fmt.Printf("%v\n", hourNow)
		endOfWorkingHour, _ := strconv.Atoi(os.Getenv("WORK_HOUR_END"))
		fmt.Printf("%v-%v\n", hourNow, endOfWorkingHour)
		if hourNow < endOfWorkingHour {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "overtime can only be submitted after 5PM", nil, nil))
			return
		}

		reimburse := model.Reimbursement{
			UserID:      userID,
			Amount:      req.Amount,
			Description: req.Description,
			CreatedBy:   userID,
			RequestIP:   r.RemoteAddr,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		saveErr := emh.EmployeeRepo.SaveReimbursement(&reimburse)
		if saveErr != nil {
			if strings.Contains(saveErr.Error(), "duplicate key") {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusConflict, "already submitted today", nil, nil))
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Overtime successfully submitted"})
	}
}

func SubmitReimbursementHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// to check user role
		if middleware.GetUserRole(c) != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		var req ReimbursementRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement"})
		}

		userID, err := uuid.Parse(middleware.GetUserID(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		reimbursement := model.Reimbursement{
			UserID:      userID,
			Amount:      req.Amount,
			Description: req.Description,
			CreatedBy:   userID,
			RequestIP:   c.ClientIP(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&reimbursement); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reimbursement"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Reimbursement successfully submitted"})
	}
}
