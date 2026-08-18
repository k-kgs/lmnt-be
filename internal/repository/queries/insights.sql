-- name: TrendForUserChallenge :many
-- Generic trend series: one numeric value per check-in date, picked from
-- metric_data by field key (e.g. "weight_kg", "pages") — the caller supplies
-- which key to project since that differs per vertical.
SELECT date, (metric_data ->> sqlc.arg(field_key)::text)::float8 AS value
FROM checkins
WHERE user_challenge_id = $1
  AND metric_data ? sqlc.arg(field_key)::text
ORDER BY date;

-- name: AdherenceByDayOfWeek :many
SELECT EXTRACT(dow FROM date)::int AS day_of_week, COUNT(*)::int AS checkin_count
FROM checkins
WHERE user_challenge_id = $1
GROUP BY EXTRACT(dow FROM date)
ORDER BY day_of_week;
