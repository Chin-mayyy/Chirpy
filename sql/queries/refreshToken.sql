-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, user_id, expires_at)
VALUES(
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
LEFT JOIN refresh_tokens on users.id = refresh_tokens.user_id
WHERE refresh_tokens.user_id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE token = $1;
