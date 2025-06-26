package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/repository"
	"payslip-generation-system/test/testutils"
	"testing"
)

func TestMain(m *testing.M) {
	testutils.SetupTestDB(nil)
	os.Exit(m.Run())
}

func TestLoginHandler_Success(t *testing.T) {
	repo := repository.NewUserRepository(testutils.DB)
	authHandler := handler.NewAuthHandler(repo)

	handler := authHandler.LoginHandler()

	body := map[string]string{
		"username": "employee001",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'data' field in response")
	}

	if _, ok := data["token"]; !ok {
		t.Error("expected token in response")
	}
}
