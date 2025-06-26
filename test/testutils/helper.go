package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/repository"
	"testing"
)

type key string

const (
	UserIDKey key = "user_id"
	RoleKey   key = "role"
)

func GetTokenFor(t *testing.T, username, password string) string {
	repo := repository.NewUserRepository(DB)
	authHandler := handler.NewAuthHandler(repo)

	h := authHandler.LoginHandler()

	body := map[string]string{
		"username": username,
		"password": password,
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	r.Header.Set("Content-Type", "application/json")

	h.ServeHTTP(w, r)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'data' field in response")
	}

	if _, ok := data["token"]; !ok {
		t.Error("expected token in response")
	}

	return data["token"].(string)
}
