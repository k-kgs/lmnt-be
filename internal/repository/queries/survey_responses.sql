-- name: UpsertSurveyResponse :one
-- Fired on every screen transition, not just completion — client_id (generated
-- once per browser, persisted in localStorage) is what lets us tell "this
-- respondent progressed from question 3 to question 4" apart from "a new
-- respondent started". completed=false rows that never get updated again are
-- exactly the drop-off data this exists to capture.
INSERT INTO survey_responses (
    client_id, completed, last_screen, track, track_other, pivot_importance, pivot_satisfaction, branch,
    positive_reasons, neutral_reasons, negative_reasons, negative_reason_other, reward_kano,
    monetization, age, gender, email, email_choice, persona_key
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
ON CONFLICT (client_id) DO UPDATE SET
    completed             = EXCLUDED.completed,
    last_screen           = EXCLUDED.last_screen,
    track                 = EXCLUDED.track,
    track_other           = EXCLUDED.track_other,
    pivot_importance      = EXCLUDED.pivot_importance,
    pivot_satisfaction    = EXCLUDED.pivot_satisfaction,
    branch                = EXCLUDED.branch,
    positive_reasons      = EXCLUDED.positive_reasons,
    neutral_reasons       = EXCLUDED.neutral_reasons,
    negative_reasons      = EXCLUDED.negative_reasons,
    negative_reason_other = EXCLUDED.negative_reason_other,
    reward_kano           = EXCLUDED.reward_kano,
    monetization          = EXCLUDED.monetization,
    age                   = EXCLUDED.age,
    gender                = EXCLUDED.gender,
    email                 = EXCLUDED.email,
    email_choice          = EXCLUDED.email_choice,
    persona_key           = EXCLUDED.persona_key,
    updated_at            = now()
RETURNING id, created_at, updated_at;
