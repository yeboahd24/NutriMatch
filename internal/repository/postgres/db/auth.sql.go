// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: auth.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token,
    expires_at
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, token, expires_at, created_at, revoked, revoked_at
`

type CreateRefreshTokenParams struct {
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.queryRow(ctx, q.createRefreshTokenStmt, createRefreshToken, arg.UserID, arg.Token, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.Revoked,
		&i.RevokedAt,
	)
	return i, err
}

const deleteExpiredRefreshTokens = `-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW()
`

func (q *Queries) DeleteExpiredRefreshTokens(ctx context.Context) error {
	_, err := q.exec(ctx, q.deleteExpiredRefreshTokensStmt, deleteExpiredRefreshTokens)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, created_at, revoked, revoked_at FROM refresh_tokens
WHERE token = $1 AND revoked = false AND expires_at > NOW()
LIMIT 1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.queryRow(ctx, q.getRefreshTokenStmt, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.Revoked,
		&i.RevokedAt,
	)
	return i, err
}

const revokeAllUserRefreshTokens = `-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET
    revoked = true,
    revoked_at = NOW()
WHERE user_id = $1 AND revoked = false
`

func (q *Queries) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.revokeAllUserRefreshTokensStmt, revokeAllUserRefreshTokens, userID)
	return err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
    revoked = true,
    revoked_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.exec(ctx, q.revokeRefreshTokenStmt, revokeRefreshToken, token)
	return err
}
