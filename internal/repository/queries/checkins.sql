-- name: InsertCheckin :one
INSERT INTO checkins (user_challenge_id, metric_data, verification_status)
VALUES ($1, $2, 'auto_approved')
RETURNING id, user_challenge_id, date, verification_status, metric_data, created_at;

-- name: ListCheckinsForUserChallenge :many
SELECT id, date, metric_data
FROM checkins
WHERE user_challenge_id = $1
ORDER BY date;
