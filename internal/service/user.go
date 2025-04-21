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
	ErrInvalidCredentials    = user.ErrInvalidCredentials
)

type userService struct {
	repo            user.Repository
	authRepo        user.AuthRepository
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
	logger          zerolog.Logger
}

func NewUserService(
	repo user.Repository,
	authRepo user.AuthRepository,
	jwtService *auth.JWTService,
	passwordService *auth.PasswordService,
	logger zerolog.Logger,
) UserService {
	return &userService{
		repo:            repo,
		authRepo:        authRepo,
		jwtService:      jwtService,
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

func (s *userService) Authenticate(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	foundUser, err := s.GetByEmail(email)
	if err != nil {
		return "", user.ErrInvalidCredentials
	}

	// Verify password
	match, err := s.passwordService.VerifyPassword(foundUser.PasswordHash, password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", user.ErrInvalidCredentials
	}

	// Update last login time
	if err := s.repo.UpdateLastLogin(foundUser.ID); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to update last login time")
		// Don't fail authentication if this fails
	}

	// In this implementation we're returning an empty token since actual token generation
	// should be handled by the auth service. The auth service should be used for actual
	// authentication in production.
	return "", nil
}

func (s *userService) CreateUser(ctx context.Context, input *user.RegisterInput) (*user.User, error) {
	// Create user using the Register method which already has the core logic
	return s.Register(input.Email, input.Password, input.FirstName, input.LastName)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	// Parse the string ID into UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	// Use existing GetByID method which handles the core logic
	return s.GetByID(userID)
}

func (s *userService) UpdateUser(ctx context.Context, userID string, input *user.UpdateUserInput) (*user.User, error) {
	// Parse string ID to UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Get existing user
	existingUser, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided in input
	if input.FirstName != nil {
		existingUser.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		existingUser.LastName = *input.LastName
	}
	if input.DateOfBirth != nil {
		existingUser.DateOfBirth = input.DateOfBirth
	}
	if input.Gender != nil {
		existingUser.Gender = *input.Gender
	}
	if input.HeightCm != nil {
		existingUser.HeightCm = *input.HeightCm
	}
	if input.WeightKg != nil {
		existingUser.WeightKg = *input.WeightKg
	}
	if input.ActivityLevel != nil {
		existingUser.ActivityLevel = *input.ActivityLevel
	}

	// Save updates
	if err := s.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	// Get user by email
	foundUser, err := s.GetByEmail(email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	// Verify password
	match, err := s.passwordService.VerifyPassword(foundUser.PasswordHash, password)
	if err != nil {
		return "", "", err
	}
	if !match {
		return "", "", ErrInvalidCredentials
	}

	// Generate tokens using JWT service
	accessToken, expiry, err := s.jwtService.GenerateAccessToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Store refresh token
	if err := s.authRepo.StoreRefreshToken(foundUser.ID, refreshToken, expiry); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
