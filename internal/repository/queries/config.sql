-- name: ListVerticals :many
SELECT id, key, label, icon, input_schema
FROM verticals
ORDER BY label;

-- name: ListActiveChallenges :many
-- member_count is computed live from active participants, never stored, so
-- it can't drift from reality the way a manually incremented counter can.
SELECT c.id, c.title, c.vertical_id, v.key AS vertical_key, c.influencer_handle,
       c.difficulty_stat, c.is_template,
       COUNT(uc.id) FILTER (WHERE uc.status = 'active')::int AS member_count
FROM challenges c
JOIN verticals v ON v.id = c.vertical_id
LEFT JOIN user_challenges uc ON uc.challenge_id = c.id
GROUP BY c.id, c.title, c.vertical_id, v.key, c.influencer_handle, c.difficulty_stat, c.is_template
ORDER BY member_count DESC;

-- name: ListActiveRedemptionItems :many
SELECT id, type, title, coin_cost, metadata
FROM redemption_items
WHERE active = true
ORDER BY coin_cost;
