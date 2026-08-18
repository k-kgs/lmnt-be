package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"kayam-be/internal/repository"
)

type CommunityHandler struct {
	Queries *repository.Queries
}

// Leaderboard ranks by current_streak within a challenge. A friends-only
// filter (via a `follows` table) is deferred per root design doc §4 —
// "only build this table if/when community screens are reached" — this
// prototype pass returns the full challenge leaderboard.
func (h *CommunityHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	var challengeID pgtype.UUID
	if err := challengeID.Scan(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := h.Queries.LeaderboardForChallenge(r.Context(), challengeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, rows)
}

// Analytics is the group-level "consistency and growth" view for one
// challenge: live participant counts, average streak, today's check-in
// rate, and retention. See internal/repository/queries/analytics.sql for
// the important caveat that retention is a snapshot approximation, not
// true day-by-day cohort tracking.
func (h *CommunityHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	var challengeID pgtype.UUID
	if err := challengeID.Scan(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	row, err := h.Queries.ChallengeAnalytics(r.Context(), challengeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, row)
}
