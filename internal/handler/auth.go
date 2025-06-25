package handler

import (
	"encoding/json"
	"net/http"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/utils"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var user model.User
		err := db.Where("username = ?", req.Username).First(&user).Error
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateToken(user.ID.String(), user.Role)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{Token: token})
	}
}
