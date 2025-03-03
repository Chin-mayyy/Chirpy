-- name: CreateUser :one
INSERT INTO users (email)
VALUES (
    $1
)

RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: AddPassword :exec
UPDATE users SET hashed_password = $1 WHERE email = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
