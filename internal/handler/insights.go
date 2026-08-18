package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"kayam-be/internal/repository"
)

type InsightHandler struct {
	Queries *repository.Queries
}

// Trend projects a single numeric field (e.g. "weight_kg", "pages") out of
// metric_data across a user_challenge's check-in history — the field key is
// caller-supplied via query param since it differs per vertical (root design
// doc §6).
func (h *InsightHandler) Trend(w http.ResponseWriter, r *http.Request) {
	var userChallengeID pgtype.UUID
	if err := userChallengeID.Scan(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	fieldKey := r.URL.Query().Get("field")
	if fieldKey == "" {
		writeError(w, http.StatusBadRequest, errMissingFieldParam)
		return
	}

	rows, err := h.Queries.TrendForUserChallenge(r.Context(), repository.TrendForUserChallengeParams{
		UserChallengeID: userChallengeID,
		FieldKey:        fieldKey,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, rows)
}

func (h *InsightHandler) Adherence(w http.ResponseWriter, r *http.Request) {
	var userChallengeID pgtype.UUID
	if err := userChallengeID.Scan(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := h.Queries.AdherenceByDayOfWeek(r.Context(), userChallengeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, rows)
}
