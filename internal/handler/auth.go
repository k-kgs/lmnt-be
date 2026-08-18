package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kayam-be/internal/service"
)

type AuthHandler struct {
	Service *service.AuthService
}

type loginRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, errors.New("name and email are required"))
		return
	}

	user, err := h.Service.Login(r.Context(), req.Name, req.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": user.ID.String(), // see middleware.RequireAuth for what this means
	})
}
