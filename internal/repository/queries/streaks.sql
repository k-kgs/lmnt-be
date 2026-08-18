-- name: GetStreak :one
SELECT user_challenge_id, current_streak, longest_streak, last_checkin_date
FROM streaks
WHERE user_challenge_id = $1;

-- name: UpsertStreak :one
INSERT INTO streaks (user_challenge_id, current_streak, longest_streak, last_checkin_date, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (user_challenge_id)
DO UPDATE SET current_streak = $2, longest_streak = $3, last_checkin_date = $4, updated_at = now()
RETURNING user_challenge_id, current_streak, longest_streak, last_checkin_date;
