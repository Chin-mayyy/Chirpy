// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
)

const addPassword = `-- name: AddPassword :exec
UPDATE users SET hashed_password = $1 WHERE email = $2
`

type AddPasswordParams struct {
	HashedPassword string
	Email          string
}

func (q *Queries) AddPassword(ctx context.Context, arg AddPasswordParams) error {
	_, err := q.db.ExecContext(ctx, addPassword, arg.HashedPassword, arg.Email)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email)
VALUES (
    $1
)

RETURNING id, created_at, updated_at, email, hashed_password
`

func (q *Queries) CreateUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}
