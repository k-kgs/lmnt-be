-- name: ListVerticals :many
SELECT id, key, label, icon, input_schema
FROM verticals
ORDER BY label;

-- name: ListActiveChallenges :many
SELECT c.id, c.title, c.vertical_id, v.key AS vertical_key, c.influencer_handle,
       c.member_count, c.difficulty_stat, c.is_template
FROM challenges c
JOIN verticals v ON v.id = c.vertical_id
ORDER BY c.member_count DESC;

-- name: ListActiveRedemptionItems :many
SELECT id, type, title, coin_cost, metadata
FROM redemption_items
WHERE active = true
ORDER BY coin_cost;
