-- name: InsertWalletTransaction :one
INSERT INTO wallet_transactions (user_id, delta, reason)
VALUES ($1, $2, $3)
RETURNING id, user_id, delta, reason, created_at;

-- name: GetBalance :one
SELECT COALESCE(SUM(delta), 0)::bigint AS balance
FROM wallet_transactions
WHERE user_id = $1;

-- name: ListWalletTransactions :many
SELECT id, delta, reason, created_at
FROM wallet_transactions
WHERE user_id = $1
ORDER BY created_at DESC;
