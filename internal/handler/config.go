package handler

import (
	"encoding/json"
	"net/http"

	"kayam-be/internal/repository"
)

type ConfigHandler struct {
	Queries *repository.Queries
}

type configResponse struct {
	Verticals       []repository.ListVerticalsRow             `json:"verticals"`
	Challenges      []repository.ListActiveChallengesRow      `json:"challenges"`
	RedemptionItems []repository.ListActiveRedemptionItemsRow `json:"redemption_items"`
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	verticals, err := h.Queries.ListVerticals(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	challenges, err := h.Queries.ListActiveChallenges(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	items, err := h.Queries.ListActiveRedemptionItems(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, configResponse{
		Verticals:       verticals,
		Challenges:      challenges,
		RedemptionItems: items,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
