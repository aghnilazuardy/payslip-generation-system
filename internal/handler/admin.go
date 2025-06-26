package handler

import (
	"encoding/json"
	"net/http"
	"payslip-generation-system/internal/helper"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/repository"
	"payslip-generation-system/internal/service"
	"time"

	"github.com/google/uuid"
)

type AttendancePeriodRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type PayrollRequest struct {
	PeriodID string `json:"attendancePeriodId"`
}

type AdminHandler struct {
	AdminRepo      repository.AdminRepository
	PayrollService service.PayrollService
}

func NewAdminHandler(adminRepo repository.AdminRepository, payrollService service.PayrollService) *AdminHandler {
	return &AdminHandler{AdminRepo: adminRepo, PayrollService: payrollService}
}

func (adh *AdminHandler) CreateAttendancePeriodHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil, nil))
			return
		}

		if middleware.GetUserRole(r) != "admin" {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusForbidden, "forbidden", nil, nil))
			return
		}

		var req AttendancePeriodRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid request", nil, nil))
			return
		}

		startDate, err1 := time.Parse("2006-01-02", req.StartDate)
		endDate, err2 := time.Parse("2006-01-02", req.EndDate)
		if err1 != nil || err2 != nil || !startDate.Before(endDate) {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid date range", nil, nil))
			return
		}

		userID, err := uuid.Parse(middleware.GetUserID(r))
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "invalid user ID", nil, nil))
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

		saveErr := adh.AdminRepo.SaveAttendancePeriod(&period)
		if saveErr != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "failed to create period", nil, nil))
			return
		}

		json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusCreated, "attendance period created successfully", nil, nil))
	}
}

func (adh *AdminHandler) RunPayroll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil, nil))
			return
		}

		if middleware.GetUserRole(r) != "admin" {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusForbidden, "forbidden", nil, nil))
			return
		}

		var req PayrollRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid JSON", nil, nil))
			return
		}

		periodID, err := uuid.Parse(req.PeriodID)
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid period ID", nil, nil))
			return
		}

		err = adh.PayrollService.ProcessPayroll(
			periodID,
			uuid.MustParse(middleware.GetUserID(r)),
			r.RemoteAddr,
			middleware.GetRequestID(r),
		)
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil))
			return
		}

		json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusCreated, "payroll processed", nil, nil))
	}
}
