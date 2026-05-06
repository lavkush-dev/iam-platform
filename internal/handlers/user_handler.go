package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"iam-platform/internal/dto"
	"iam-platform/internal/services"
)

type UserHandler struct {
	service *services.UserService
	logger  *zap.Logger
}

func NewUserHandler(s *services.UserService, l *zap.Logger) *UserHandler {
	return &UserHandler{s, l}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
