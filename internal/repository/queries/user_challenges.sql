-- name: CreateUserChallenge :one
INSERT INTO user_challenges (user_id, challenge_id, custom_goal)
VALUES ($1, $2, $3)
RETURNING id, user_id, challenge_id, custom_goal, status, joined_at;

-- name: GetUserChallengeWithVertical :one
SELECT uc.id, uc.user_id, uc.challenge_id, uc.status,
       c.vertical_id, v.key AS vertical_key, v.input_schema
FROM user_challenges uc
JOIN challenges c ON c.id = uc.challenge_id
JOIN verticals v ON v.id = c.vertical_id
WHERE uc.id = $1;

-- name: ListUserChallengesForUser :many
SELECT uc.id, uc.status, uc.joined_at,
       c.title AS challenge_title, c.influencer_handle,
       v.key AS vertical_key, v.label AS vertical_label,
       COALESCE(s.current_streak, 0) AS current_streak,
       COALESCE(s.longest_streak, 0) AS longest_streak,
       s.last_checkin_date
FROM user_challenges uc
JOIN challenges c ON c.id = uc.challenge_id
JOIN verticals v ON v.id = c.vertical_id
LEFT JOIN streaks s ON s.user_challenge_id = uc.id
WHERE uc.user_id = $1
ORDER BY uc.joined_at DESC;
