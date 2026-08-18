package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/repository"
	"kayam-be/internal/service"
)

type RedemptionHandler struct {
	Service *service.RedemptionService
	Queries *repository.Queries
}

type createRedemptionRequest struct {
	RedemptionItemID string `json:"redemption_item_id"`
}

func (h *RedemptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	var req createRedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var itemID pgtype.UUID
	if err := itemID.Scan(req.RedemptionItemID); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.Service.Redeem(r.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			writeError(w, http.StatusPaymentRequired, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *RedemptionHandler) MyRedemptions(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	rows, err := h.Queries.ListRedemptionsForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, rows)
}
