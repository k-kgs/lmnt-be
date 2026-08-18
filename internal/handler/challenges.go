package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/repository"
)

type ChallengeHandler struct {
	Queries *repository.Queries
}

type joinChallengeRequest struct {
	CustomGoal json.RawMessage `json:"custom_goal,omitempty"`
}

// Join implements the "Choose Your Path" convergence point from the
// wireframes — both "Join a Challenge" and "Set My Own Goal" land here,
// the latter simply passing a custom_goal payload alongside a challenge
// (a system-provided template challenge stands in for a pure custom goal
// in this prototype; see kayam-fe design.md for the FE side of this).
func (h *ChallengeHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	challengeIDStr := chi.URLParam(r, "id")
	var challengeID pgtype.UUID
	if err := challengeID.Scan(challengeIDStr); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req joinChallengeRequest
	_ = json.NewDecoder(r.Body).Decode(&req) // body is optional

	uc, err := h.Queries.CreateUserChallenge(r.Context(), repository.CreateUserChallengeParams{
		UserID:      userID,
		ChallengeID: challengeID,
		CustomGoal:  req.CustomGoal,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, uc)
}

// MyChallenges lists the current user's joined challenges with their live
// streak — this is what drives the Home/dashboard and Calendar screens.
func (h *ChallengeHandler) MyChallenges(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	rows, err := h.Queries.ListUserChallengesForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, rows)
}
