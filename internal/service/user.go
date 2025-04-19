package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/pkg/auth"
)

// Import errors from the user domain
var (
	ErrEmailUpdateNotAllowed = user.ErrEmailUpdateNotAllowed
)

type userService struct {
	repo            user.Repository
	passwordService *auth.PasswordService
	logger          zerolog.Logger
}

func NewUserService(
	repo user.Repository,
	passwordService *auth.PasswordService,
	logger zerolog.Logger,
) UserService {
	return &userService{
		repo:            repo,
		passwordService: passwordService,
		logger:          logger,
	}
}

func (s *userService) Register(email, password, firstName, lastName string) (*user.User, error) {
	hashedPassword, err := s.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id uuid.UUID) (*user.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetByEmail(email string) (*user.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) Update(user *user.User) error {
	return s.repo.Update(user)
}

func (s *userService) UpdatePassword(id uuid.UUID, currentPassword, newPassword string) error {
	foundUser, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	match, err := s.passwordService.VerifyPassword(foundUser.PasswordHash, currentPassword)
	if err != nil {
		return err
	}
	if !match {
		return user.ErrInvalidPassword
	}

	hashedPassword, err := s.passwordService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(id, hashedPassword)
}

func (s *userService) VerifyEmail(id uuid.UUID) error {
	return s.repo.UpdateEmailVerification(id, true)
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// Adapter methods to implement the service.UserService interface
func (s *userService) GetProfile(ctx context.Context, userID string) (*user.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *userService) UpdateProfile(ctx context.Context, userID, name, email string) error {
	s.logger.Info().Str("user_id", userID).Str("name", name).Str("email", email).Msg("Updating user profile")

	id, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to parse user ID")
		return err
	}

	user, err := s.GetByID(id)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user by ID")
		return err
	}

	// Check if email is being changed
	if email != "" && email != user.Email {
		s.logger.Warn().Str("user_id", userID).Str("current_email", user.Email).Str("requested_email", email).Msg("Email update attempted but not allowed")
		// Access the error from the package, not from the user instance
		return ErrEmailUpdateNotAllowed
	}

	// Split name into first and last name
	if name != "" {
		names := strings.Split(name, " ")
		firstName := names[0]
		lastName := ""
		if len(names) > 1 {
			lastName = strings.Join(names[1:], " ")
		}

		user.FirstName = firstName
		user.LastName = lastName
	}

	s.logger.Info().Str("user_id", userID).Str("first_name", user.FirstName).Str("last_name", user.LastName).Msg("Updating user name")
	return s.Update(user)
}
