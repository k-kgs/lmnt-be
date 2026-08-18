package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/service"
)

type CheckinHandler struct {
	Service *service.CheckinService
}

type createCheckinRequest struct {
	UserChallengeID string         `json:"user_challenge_id"`
	MetricData      map[string]any `json:"metric_data"`
}

func (h *CheckinHandler) Create(w http.ResponseWriter, r *http.Request) {
	if _, ok := kmw.UserIDFromContext(r.Context()); !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	var req createCheckinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var userChallengeID pgtype.UUID
	if err := userChallengeID.Scan(req.UserChallengeID); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.Service.CreateCheckin(r.Context(), userChallengeID, req.MetricData)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, service.ErrAlreadyCheckedIn):
			writeError(w, http.StatusConflict, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}
