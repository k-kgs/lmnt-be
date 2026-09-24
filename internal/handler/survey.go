package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kayam-be/internal/repository"
)

var errMissingClientID = errors.New("clientId is required")

type SurveyHandler struct {
	Queries *repository.Queries
}

// surveyResponseInput mirrors lmnt-fe's checkpoint payload (src/api/survey.ts,
// built from src/hooks/useSurveyFlow.ts's SurveyAnswers) field for field. The
// frontend posts one of these on every screen transition, not just at the end
// — clientID is what makes this idempotent per respondent (see UpsertSurveyResponse).
type surveyResponseInput struct {
	ClientID             string   `json:"clientId"`
	Completed            bool     `json:"completed"`
	LastScreen           *string  `json:"lastScreen"`
	Track                *string  `json:"track"`
	TrackOther           *string  `json:"trackOther"`
	TrackingTool         *string  `json:"trackingTool"`
	TrackingToolOther    *string  `json:"trackingToolOther"`
	TrackingSatisfaction *string  `json:"trackingSatisfaction"`
	TrackingAppFeedback  *string  `json:"trackingAppFeedback"`
	PivotImportance      *string  `json:"pivotImportance"`
	PivotSatisfaction    *string  `json:"pivotSatisfaction"`
	Branch               *string  `json:"branch"`
	PositiveReasons      []string `json:"positiveReasons"`
	NeutralReasons       []string `json:"neutralReasons"`
	NegativeReasons      []string `json:"negativeReasons"`
	NegativeReasonOther  *string  `json:"negativeReasonOther"`
	RewardKano           *string  `json:"rewardKano"`
	Monetization         *string  `json:"monetization"`
	Age                  *string  `json:"age"`
	Gender               *string  `json:"gender"`
	Email                *string  `json:"email"`
	EmailChoice          *string  `json:"emailChoice"`
	PersonaKey           *string  `json:"personaKey"`
}

func marshalStringSlice(v []string) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

// Upsert stores (or updates) one survey response, keyed on the client-generated
// clientId. Fired on every screen transition — a row with completed=false that
// never gets a later update with completed=true is exactly the drop-off signal
// this exists to capture. Public, unauthenticated: a respondent may or may not
// be logged in, and no user_id is recorded.
func (h *SurveyHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var in surveyResponseInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if in.ClientID == "" {
		writeError(w, http.StatusBadRequest, errMissingClientID)
		return
	}

	positiveReasons, err := marshalStringSlice(in.PositiveReasons)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	neutralReasons, err := marshalStringSlice(in.NeutralReasons)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	negativeReasons, err := marshalStringSlice(in.NegativeReasons)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	row, err := h.Queries.UpsertSurveyResponse(r.Context(), repository.UpsertSurveyResponseParams{
		ClientID:             in.ClientID,
		Completed:            in.Completed,
		LastScreen:           in.LastScreen,
		Track:                in.Track,
		TrackOther:           in.TrackOther,
		TrackingTool:         in.TrackingTool,
		TrackingToolOther:    in.TrackingToolOther,
		TrackingSatisfaction: in.TrackingSatisfaction,
		TrackingAppFeedback:  in.TrackingAppFeedback,
		PivotImportance:      in.PivotImportance,
		PivotSatisfaction:    in.PivotSatisfaction,
		Branch:               in.Branch,
		PositiveReasons:      positiveReasons,
		NeutralReasons:       neutralReasons,
		NegativeReasons:      negativeReasons,
		NegativeReasonOther:  in.NegativeReasonOther,
		RewardKano:           in.RewardKano,
		Monetization:         in.Monetization,
		Age:                  in.Age,
		Gender:               in.Gender,
		Email:                in.Email,
		EmailChoice:          in.EmailChoice,
		PersonaKey:           in.PersonaKey,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"id": row.ID})
}
