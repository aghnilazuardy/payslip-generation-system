package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/middleware"
	"payslip-generation-system/internal/repository"
	"payslip-generation-system/test/testutils"
	"testing"
)

// func TestSubmitAttendance_Success(t *testing.T) {
// 	db := testutils.DB
// 	repo := repository.NewEmployeeRepository(db)
// 	employeeHandler := handler.NewEmployeeHandler(repo)
// 	handler := employeeHandler.SubmitAttendanceHanlder()

// 	// Simulate login
// 	token := testutils.GetTokenFor(t, "employee001", "password")

// 	req := httptest.NewRequest(http.MethodPost, "/employee/attendance", nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("X-Real-IP", "127.0.0.1")

// 	w := httptest.NewRecorder()
// 	protected := middleware.AuthMiddleware(handler)
// 	protected.ServeHTTP(w, req)

// 	if w.Code != http.StatusCreated {
// 		t.Errorf("expected status 201, got %d", w.Code)
// 	}
// }

func TestSubmitOvertime_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewEmployeeRepository(db)
	employeeHandler := handler.NewEmployeeHandler(repo)
	h := employeeHandler.SubmitOvertimeHandler()
	protected := middleware.AuthMiddleware(h)

	token := testutils.GetTokenFor(t, "employee001", "password")

	body := map[string]interface{}{
		"hours": 2,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	protected.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestSubmitReimbursement_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewEmployeeRepository(db)
	employeeHandler := handler.NewEmployeeHandler(repo)
	h := employeeHandler.SubmitReimbursementHandler()
	protected := middleware.AuthMiddleware(h)

	token := testutils.GetTokenFor(t, "employee001", "password")

	body := map[string]interface{}{
		"amount":      250000,
		"description": "Internet allowance",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/employee/reimbursement", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	protected.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestGeneratePayslip_Success(t *testing.T) {
	db := testutils.DB
	repo := repository.NewEmployeeRepository(db)
	employeeHandler := handler.NewEmployeeHandler(repo)
	h := employeeHandler.GetPayslipHandler()

	// Seed attendance + reimbursement + period
	token := testutils.GetTokenFor(t, "employee001", "password")

	body := map[string]interface{}{
		"payrollID": "335b1d33-eddb-4a8d-9acf-eef397f6e7e2",
	}
	jsonBody, _ := json.Marshal(body)

	// Generate payslip
	req := httptest.NewRequest(http.MethodGet, "/employee/payslip", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	middleware.AuthMiddleware(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	if _, ok := response["data"].(map[string]interface{}); !ok {
		t.Error("expected payslip data in response")
	}
}
