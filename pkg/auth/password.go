package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/yeboahd24/nutrimatch/internal/config"
)

// PasswordService handles password hashing and verification
type PasswordService struct {
	config config.SecurityConfig
}

// NewPasswordService creates a new password service
func NewPasswordService(config config.SecurityConfig) *PasswordService {
	return &PasswordService{
		config: config,
	}
}

// HashPassword hashes a password using Argon2id
func (s *PasswordService) HashPassword(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      s.config.ArgonMemory,
		Iterations:  s.config.ArgonIterations,
		Parallelism: s.config.ArgonParallelism,
		SaltLength:  s.config.ArgonSaltLength,
		KeyLength:   s.config.ArgonKeyLength,
	}

	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(hash, password string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
