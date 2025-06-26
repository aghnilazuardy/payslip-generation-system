package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"payslip-generation-system/internal/helper"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OvertimeRequest struct {
	Hours int `json:"hours"`
}

type ReimbursementRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type EmployeeHandler struct {
	AttendanceRepo repository.AttendanceRepository
	OvertimeRepo   repository.OvertimeRepository
}

func NewEmployeeHandler(attendanceRepo repository.AttendanceRepository, overtimeRepo repository.OvertimeRepository) *EmployeeHandler {
	return &EmployeeHandler{AttendanceRepo: attendanceRepo, OvertimeRepo: overtimeRepo}
}

func (emh *EmployeeHandler) SubmitAttendanceHanlder() http.HandlerFunc {
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

		// to check is today weekend or weekday
		today := time.Now().Truncate(24 * time.Hour)
		weekday := today.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "cannot submit on weekend", nil, nil))
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "invalid user ID", nil, nil))
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

		saveErr := emh.AttendanceRepo.SaveAttendance(&attendance)
		if saveErr != nil {
			if strings.Contains(saveErr.Error(), "duplicate key") {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusConflict, "already submitted today", nil, nil))
			} else {
				json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "failed to create attendance", nil, nil))
			}
			return
		}

		json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusCreated, "attendance submitted successfully", nil, nil))
	}
}

func (emh *EmployeeHandler) SubmitOvertimeHandler() http.HandlerFunc {
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

		var req OvertimeRequest
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

		saveErr := emh.OvertimeRepo.SaveOvertime(&overtime)
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
