package handler

import (
	"encoding/json"
	"net/http"
	"payslip-generation-system/internal/helper"
	"payslip-generation-system/internal/repository"
	"payslip-generation-system/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	UserRepo repository.UserRepository
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

func (lh *AuthHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil, nil))
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusBadRequest, "invalid request", nil, nil))
			return
		}

		user, err := lh.UserRepo.FindByUsername(req.Username)
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusUnauthorized, "invalid credentials", nil, nil))
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusUnauthorized, "invalid credentials", nil, nil))
			return
		}

		token, err := utils.GenerateToken(user.ID.String(), user.Role)
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusInternalServerError, "failed to generate token", nil, nil))
			return
		}

		json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusOK, "login success", LoginResponse{Token: token}, nil))
	}
}
