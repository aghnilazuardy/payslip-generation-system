package handler

import (
	"encoding/json"
	"net/http"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendancePeriodRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type PayrollRequest struct {
	PeriodID string `json:"attendancePeriodId"`
}

func CreateAttendancePeriodHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		if middleware.GetUserRole(r) != "admin" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
			return
		}

		var req AttendancePeriodRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		startDate, err1 := time.Parse("2006-01-02", req.StartDate)
		endDate, err2 := time.Parse("2006-01-02", req.EndDate)
		if err1 != nil || err2 != nil || !startDate.Before(endDate) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid date range"})
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid user ID"})
			return
		}

		period := model.AttendancePeriod{
			StartDate: startDate,
			EndDate:   endDate,
			CreatedBy: userID,
			RequestIP: r.RemoteAddr,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&period).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to create period"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "attendance period created successfully"})
	}
}
