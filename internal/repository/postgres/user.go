package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) user.Repository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) Create(user *user.User) error {
	var dateOfBirth sql.NullTime
	if user.DateOfBirth != nil {
		dateOfBirth = sql.NullTime{Time: *user.DateOfBirth, Valid: true}
	}

	_, err := r.queries.CreateUser(context.Background(), db.CreateUserParams{
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		FirstName:     sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
		LastName:      sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		DateOfBirth:   dateOfBirth,
		Gender:        sql.NullString{String: user.Gender, Valid: user.Gender != ""},
		HeightCm:      sql.NullString{String: fmt.Sprintf("%.2f", user.HeightCm), Valid: user.HeightCm != 0},
		WeightKg:      sql.NullString{String: fmt.Sprintf("%.2f", user.WeightKg), Valid: user.WeightKg != 0},
		ActivityLevel: sql.NullString{String: user.ActivityLevel, Valid: user.ActivityLevel != ""},
	})
	return err
}

func (r *userRepository) GetByID(id uuid.UUID) (*user.User, error) {
	u, err := r.queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return mapDbUserToDomain(&u), nil
}

func (r *userRepository) GetByEmail(email string) (*user.User, error) {
	u, err := r.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return mapDbUserToDomain(&u), nil
}

func (r *userRepository) Update(user *user.User) error {
	var dateOfBirth sql.NullTime
	if user.DateOfBirth != nil {
		dateOfBirth = sql.NullTime{Time: *user.DateOfBirth, Valid: true}
	}

	_, err := r.queries.UpdateUser(context.Background(), db.UpdateUserParams{
		ID:            user.ID,
		FirstName:     sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
		LastName:      sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		DateOfBirth:   dateOfBirth,
		Gender:        sql.NullString{String: user.Gender, Valid: user.Gender != ""},
		HeightCm:      sql.NullString{String: fmt.Sprintf("%.2f", user.HeightCm), Valid: user.HeightCm != 0},
		WeightKg:      sql.NullString{String: fmt.Sprintf("%.2f", user.WeightKg), Valid: user.WeightKg != 0},
		ActivityLevel: sql.NullString{String: user.ActivityLevel, Valid: user.ActivityLevel != ""},
	})
	return err
}

func (r *userRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	return r.queries.UpdateUserPassword(context.Background(), db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
}

func (r *userRepository) UpdateLastLogin(id uuid.UUID) error {
	return r.queries.UpdateUserLastLogin(context.Background(), id)
}

func (r *userRepository) UpdateEmailVerification(id uuid.UUID, verified bool) error {
	return r.queries.UpdateUserEmailVerification(context.Background(), db.UpdateUserEmailVerificationParams{
		ID:            id,
		EmailVerified: sql.NullBool{Bool: verified, Valid: true},
	})
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.queries.DeleteUser(context.Background(), id)
}

func mapDbUserToDomain(u *db.User) *user.User {
	var dateOfBirth *time.Time
	if u.DateOfBirth.Valid {
		dateOfBirth = &u.DateOfBirth.Time
	}

	var lastLogin *time.Time
	if u.LastLogin.Valid {
		lastLogin = &u.LastLogin.Time
	}

	heightCm, _ := strconv.ParseFloat(u.HeightCm.String, 64)
	weightKg, _ := strconv.ParseFloat(u.WeightKg.String, 64)

	return &user.User{
		ID:            u.ID,
		Email:         u.Email,
		PasswordHash:  u.PasswordHash,
		FirstName:     u.FirstName.String,
		LastName:      u.LastName.String,
		DateOfBirth:   dateOfBirth,
		Gender:        u.Gender.String,
		HeightCm:      heightCm,
		WeightKg:      weightKg,
		ActivityLevel: u.ActivityLevel.String,
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		LastLogin:     lastLogin,
		AccountStatus: u.AccountStatus.String,
		EmailVerified: u.EmailVerified.Bool,
		MFAEnabled:    u.MfaEnabled.Bool,
		MFASecret:     u.MfaSecret.String,
	}
}
