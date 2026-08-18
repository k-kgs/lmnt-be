-- name: GetRedemptionItem :one
SELECT id, type, title, coin_cost, metadata, active
FROM redemption_items
WHERE id = $1;

-- name: InsertRedemption :one
INSERT INTO redemptions (user_id, redemption_item_id, code_or_slot)
VALUES ($1, $2, $3)
RETURNING id, user_id, redemption_item_id, code_or_slot, created_at;

-- name: ListRedemptionsForUser :many
SELECT r.id, r.code_or_slot, r.created_at,
       ri.title, ri.type, ri.coin_cost
FROM redemptions r
JOIN redemption_items ri ON ri.id = r.redemption_item_id
WHERE r.user_id = $1
ORDER BY r.created_at DESC;
