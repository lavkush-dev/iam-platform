package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"iam-platform/internal/dto"
	"iam-platform/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(s *services.AuthService, l *zap.Logger) *AuthHandler {
	return &AuthHandler{s, l}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}
