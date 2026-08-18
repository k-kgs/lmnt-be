-- name: GetUserByEmail :one
SELECT id, name, email, phone, auth_method, last_login_at, login_count, created_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, auth_method)
VALUES ($1, $2, $3)
RETURNING id, name, email, phone, auth_method, last_login_at, login_count, created_at;

-- name: RecordLogin :exec
UPDATE users
SET last_login_at = now(), login_count = login_count + 1
WHERE id = $1;

-- name: GetUserByID :one
SELECT id, name, email, phone, auth_method, last_login_at, login_count, created_at
FROM users
WHERE id = $1;

-- name: LockUserForRedemption :exec
-- Locks the user's row for the duration of the transaction so two concurrent
-- redemption requests for the same user serialize instead of racing on balance.
SELECT id FROM users WHERE id = $1 FOR UPDATE;
