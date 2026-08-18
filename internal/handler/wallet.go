package handler

import (
	"net/http"

	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/repository"
)

type WalletHandler struct {
	Queries *repository.Queries
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	balance, err := h.Queries.GetBalance(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	transactions, err := h.Queries.ListWalletTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"balance":      balance,
		"transactions": transactions,
	})
}
