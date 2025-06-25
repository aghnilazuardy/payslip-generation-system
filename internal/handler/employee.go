package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OvertimeRequest struct {
	Hours int `json:"hours"`
}

func SubmitAttendanceHanlder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		// to check user role
		if middleware.GetUserRole(r) != "employee" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
			return
		}

		// to check is today weekend or weekday
		today := time.Now().Truncate(24 * time.Hour)
		weekday := today.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "cannot submit on weekend"})
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid user ID"})
			return
		}

		log.Printf("Attendance submission: user_id=%s date=%v", userID.String(), today.Format("2006-01-02"))

		attendance := model.Attendance{
			UserID:    userID,
			Date:      today,
			CreatedBy: userID,
			RequestIP: r.RemoteAddr,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&attendance).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				w.WriteHeader(http.StatusConflict)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": "already submitted today"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": "failed to create attendance"})
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "attendance submitted successfully"})
	}
}

func SubmitOvertimeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		// to check user role
		if middleware.GetUserRole(r) != "employee" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
			return
		}

		var req OvertimeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid user ID"})
			return
		}

		hourNow := time.Now().Hour()
		fmt.Printf("%v\n", hourNow)
		endOfWorkingHour, _ := strconv.Atoi(os.Getenv("WORK_HOUR_END"))
		fmt.Printf("%v-%v\n", hourNow, endOfWorkingHour)
		if hourNow < endOfWorkingHour {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "overtime can only be submitted after 5PM"})
			return
		}

		today := time.Now().Truncate(24 * time.Hour)

		overtime := model.Overtime{
			UserID:    userID,
			Date:      today,
			Hours:     req.Hours,
			CreatedBy: userID,
			RequestIP: r.RemoteAddr,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&overtime).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				w.WriteHeader(http.StatusConflict)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": "already submitted today"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": "failed to submit overtime"})
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "overtime submitted successfully"})
	}
}
