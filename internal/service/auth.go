package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/pkg/auth"
)

type authService struct {
	userRepo        user.Repository
	authRepo        user.AuthRepository
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
	logger          zerolog.Logger
}

func NewAuthService(
	userRepo user.Repository,
	authRepo user.AuthRepository,
	jwtService *auth.JWTService,
	passwordService *auth.PasswordService,
	logger zerolog.Logger,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		authRepo:        authRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		logger:          logger,
	}
}

func (s *authService) DomainRegister(req user.RegisterRequest) (*user.User, error) {
	// Hash password
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &user.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) DomainLogin(creds user.Credentials) (*user.TokenPair, error) {
	// Get user by email
	foundUser, err := s.userRepo.GetByEmail(creds.Email)
	if err != nil {
		return nil, err
	}

	// Verify password
	match, err := s.passwordService.VerifyPassword(foundUser.PasswordHash, creds.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, user.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, expiry, err := s.jwtService.GenerateAccessToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := s.authRepo.StoreRefreshToken(foundUser.ID, refreshToken, expiry); err != nil {
		return nil, err
	}

	return &user.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiry,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	s.logger.Info().Msg("Refreshing token")

	// Validate refresh token
	token, err := s.authRepo.GetRefreshToken(refreshToken)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get refresh token")
		return "", "", err
	}

	if token.Revoked {
		s.logger.Error().Msg("Token is revoked")
		return "", "", user.ErrInvalidToken
	}

	// Get user
	foundUser, err := s.userRepo.GetByID(token.UserID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get user by ID")
		return "", "", err
	}

	// Generate new tokens
	accessToken, expiry, err := s.jwtService.GenerateAccessToken(foundUser.ID, foundUser.Email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate access token")
		return "", "", err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate refresh token")
		return "", "", err
	}

	// Revoke old token and store new one
	if err := s.authRepo.RevokeRefreshToken(refreshToken); err != nil {
		s.logger.Error().Err(err).Msg("Failed to revoke old refresh token")
		return "", "", err
	}

	if err := s.authRepo.StoreRefreshToken(foundUser.ID, newRefreshToken, expiry); err != nil {
		s.logger.Error().Err(err).Msg("Failed to store new refresh token")
		return "", "", err
	}

	s.logger.Info().Msg("Token refreshed successfully")
	return accessToken, newRefreshToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	s.logger.Info().Msg("Logging out user")

	// Validate that the refresh token exists and is not already revoked
	token, err := s.authRepo.GetRefreshToken(refreshToken)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get refresh token during logout")
		return err
	}

	if token.Revoked {
		s.logger.Error().Msg("Token is already revoked")
		return user.ErrInvalidToken
	}

	// Revoke the token
	return s.authRepo.RevokeRefreshToken(refreshToken)
}

func (s *authService) LogoutAll(userID uuid.UUID) error {
	return s.authRepo.RevokeAllUserTokens(userID)
}

func (s *authService) VerifyPassword(hashedPassword, password string) bool {
	match, err := s.passwordService.VerifyPassword(hashedPassword, password)
	if err != nil {
		return false
	}
	return match
}

func (s *authService) HashPassword(password string) (string, error) {
	return s.passwordService.HashPassword(password)
}

// Adapter methods to implement the service.AuthService interface
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	creds := user.Credentials{
		Email:    email,
		Password: password,
	}

	tokenPair, err := s.DomainLogin(creds)
	if err != nil {
		return "", "", err
	}

	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

func (s *authService) Register(ctx context.Context, email, password, name string) error {
	names := strings.Split(name, " ")
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	req := user.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	_, err := s.DomainRegister(req)
	return err
}
