-- name: LeaderboardForChallenge :many
SELECT uc.id AS user_challenge_id, u.id AS user_id, u.name,
       COALESCE(s.current_streak, 0) AS current_streak
FROM user_challenges uc
JOIN users u ON u.id = uc.user_id
LEFT JOIN streaks s ON s.user_challenge_id = uc.id
WHERE uc.challenge_id = $1
ORDER BY current_streak DESC, uc.joined_at ASC
LIMIT 50;
