package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type authRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) *authRepository {
	return &authRepository{
		queries: queries,
	}
}

func (r *authRepository) StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	_, err := r.queries.CreateRefreshToken(context.Background(), db.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

func (r *authRepository) GetRefreshToken(token string) (*user.RefreshToken, error) {
	rt, err := r.queries.GetRefreshToken(context.Background(), token)
	if err != nil {
		return nil, err
	}

	var revokedAt *time.Time
	if rt.RevokedAt.Valid {
		revokedAt = &rt.RevokedAt.Time
	}

	return &user.RefreshToken{
		Token:     rt.Token,
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt.Time,
		Revoked:   rt.Revoked.Bool,
		RevokedAt: revokedAt,
	}, nil
}

func (r *authRepository) RevokeRefreshToken(token string) error {
	return r.queries.RevokeRefreshToken(context.Background(), token)
}

func (r *authRepository) RevokeAllUserRefreshTokens(userID uuid.UUID) error {
	return r.queries.RevokeAllUserRefreshTokens(context.Background(), userID)
}

func (r *authRepository) DeleteExpiredRefreshTokens() error {
	return r.queries.DeleteExpiredRefreshTokens(context.Background())
}

func (r *authRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.queries.RevokeAllUserRefreshTokens(context.Background(), userID)
}
