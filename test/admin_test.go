package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/repository"
	"payslip-generation-system/internal/service"
	"payslip-generation-system/test/testutils"
	"testing"
	"time"
)

func TestSubmitAttendancePeriod_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewAdminRepository(db)
	payrollRepo := repository.NewPayrollRepository(db)
	service := service.NewPayrollService(payrollRepo)
	adminHandler := handler.NewAdminHandler(repo, service)
	h := adminHandler.CreateAttendancePeriodHandler()
	protected := middleware.AuthMiddleware(h)

	token := testutils.GetTokenFor(t, "admin", "password")

	body := map[string]interface{}{
		"startDate": time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
		"endDate":   time.Now().Format("2006-01-02"),
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/admin/attendance-period", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	protected.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestRunPayroll_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewAdminRepository(db)
	payrollRepo := repository.NewPayrollRepository(db)
	service := service.NewPayrollService(payrollRepo)
	adminHandler := handler.NewAdminHandler(repo, service)
	runHandler := adminHandler.RunPayroll()

	token := testutils.GetTokenFor(t, "admin", "password")

	body := map[string]interface{}{
		"attendancePeriodId": "ae2c633c-ffa3-4038-b828-4dbb0403b7b6",
	}
	jsonBody, _ := json.Marshal(body)

	// Run payroll
	runReq := httptest.NewRequest(http.MethodPost, "/admin/payroll-run", bytes.NewReader(jsonBody))
	runReq.Header.Set("Authorization", "Bearer "+token)
	runReq.Header.Set("X-Real-IP", "127.0.0.1")
	runW := httptest.NewRecorder()
	middleware.AuthMiddleware(runHandler).ServeHTTP(runW, runReq)

	if runW.Code != http.StatusCreated {
		t.Errorf("expected payroll run status 201, got %d", runW.Code)
	}
}

func TestSummaryPayslips_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewAdminRepository(db)
	payrollRepo := repository.NewPayrollRepository(db)
	service := service.NewPayrollService(payrollRepo)
	adminHandler := handler.NewAdminHandler(repo, service)
	h := adminHandler.GetPayslipSummaryHandler()
	protected := middleware.AuthMiddleware(h)

	token := testutils.GetTokenFor(t, "admin", "password")

	body := map[string]interface{}{
		"payrollID": "335b1d33-eddb-4a8d-9acf-eef397f6e7e2",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodGet, "/admin/payslips", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	protected.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp["data"].([]interface{}); !ok {
		t.Error("expected summary payslip data as list")
	}
}
