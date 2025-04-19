package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
)

type profileService struct {
	profileRepo profile.Repository
	userRepo    user.Repository
	logger      zerolog.Logger
}

func NewProfileService(
	profileRepo profile.Repository,
	userRepo user.Repository,
	logger zerolog.Logger,
) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *profileService) Create(profile *profile.UserProfile) error {
	// Verify user exists
	if _, err := s.userRepo.GetByID(profile.UserID); err != nil {
		return err
	}

	return s.profileRepo.Create(profile)
}

func (s *profileService) GetByID(id uuid.UUID, userID uuid.UUID) (*profile.UserProfile, error) {
	p, err := s.profileRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if p.UserID != userID {
		return nil, profile.ErrUnauthorized
	}

	return p, nil
}

func (s *profileService) GetByUserID(userID uuid.UUID) ([]profile.UserProfile, error) {
	// Verify user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}

	return s.profileRepo.GetByUserID(userID)
}

func (s *profileService) GetDefaultByUserID(userID uuid.UUID) (*profile.UserProfile, error) {
	// Verify user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}

	return s.profileRepo.GetDefaultByUserID(userID)
}

func (s *profileService) Update(updatedProfile *profile.UserProfile) error {
	// Verify profile exists
	existing, err := s.profileRepo.GetByID(updatedProfile.ID)
	if err != nil {
		return err
	}

	// Ensure we preserve the original user ID
	updatedProfile.UserID = existing.UserID

	return s.profileRepo.Update(updatedProfile)
}

func (s *profileService) SetAsDefault(id uuid.UUID, userID uuid.UUID) error {
	// Verify profile exists and belongs to user
	existing, err := s.profileRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return profile.ErrUnauthorized
	}

	return s.profileRepo.SetAsDefault(id, userID)
}

func (s *profileService) Delete(id uuid.UUID, userID uuid.UUID) error {
	// Verify profile exists and belongs to user
	existing, err := s.profileRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return profile.ErrUnauthorized
	}

	return s.profileRepo.Delete(id, userID)
}

// Adapter methods to implement the service.ProfileService interface
func (s *profileService) CreateProfile(ctx context.Context, userID string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) (*profile.UserProfile, error) {
	s.logger.Info().Str("user_id", userID).Msg("Creating new profile")

	uid, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to parse user ID")
		return nil, err
	}

	// Generate a new profile ID
	profileID := uuid.New()
	s.logger.Info().Str("profile_id", profileID.String()).Str("user_id", userID).Msg("Generated new profile ID")

	newProfile := &profile.UserProfile{
		ID:                  profileID,
		UserID:              uid,
		HealthConditions:    []string{},
		DietaryRestrictions: []string{},
		Allergens:           allergies,
		GoalType:            goals[0],
		PreferredFoods:      preferences,
		DislikedFoods:       []string{},
		CuisinePreferences:  []string{},
		IsDefault:           isDefault,
	}

	// Create the profile
	s.logger.Info().Str("profile_id", profileID.String()).Msg("Creating profile in database")
	if err := s.Create(newProfile); err != nil {
		s.logger.Error().Err(err).Str("profile_id", profileID.String()).Msg("Failed to create profile")
		return nil, err
	}

	// If isDefault is true, set this profile as the default
	if isDefault {
		s.logger.Info().Str("profile_id", newProfile.ID.String()).Msg("Setting profile as default")
		if err := s.profileRepo.SetAsDefault(newProfile.ID, uid); err != nil {
			s.logger.Error().Err(err).Str("profile_id", newProfile.ID.String()).Msg("Failed to set profile as default")
			// Continue anyway, but log the error
		} else {
			s.logger.Info().Str("profile_id", newProfile.ID.String()).Msg("Profile set as default")
			// Update the local object to reflect the change
			newProfile.IsDefault = true
		}
	}

	// The profile ID has been updated with the database-generated ID in the repository layer
	// Log the updated profile ID
	s.logger.Info().Str("profile_id", newProfile.ID.String()).Msg("Profile created with database-generated ID")

	// Verify the profile was created using the updated ID
	_, err = s.profileRepo.GetByID(newProfile.ID)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", newProfile.ID.String()).Msg("Failed to verify profile creation")
		// Continue anyway, but log the error
	} else {
		s.logger.Info().Str("profile_id", newProfile.ID.String()).Msg("Profile creation verified")
	}

	return newProfile, nil
}

func (s *profileService) GetProfile(ctx context.Context, id string) (*profile.UserProfile, error) {
	s.logger.Info().Str("profile_id", id).Msg("Getting profile by ID")

	profileID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to parse profile ID")
		return nil, err
	}

	// Since we don't have the user ID in this context, we'll need to fetch the profile first
	// and then verify ownership in a separate step
	p, err := s.profileRepo.GetByID(profileID)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to get profile by ID")
		return nil, err
	}

	s.logger.Info().Str("profile_id", id).Str("user_id", p.UserID.String()).Msg("Successfully retrieved profile")
	return p, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, id string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) error {
	s.logger.Info().Str("profile_id", id).Msg("Starting profile update")

	profileID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to parse profile ID")
		return err
	}

	// Log the update attempt
	s.logger.Debug().Str("profile_id", id).Interface("goals", goals).Interface("allergies", allergies).Interface("preferences", preferences).Msg("Attempting to update profile")

	// Get existing profile
	p, err := s.profileRepo.GetByID(profileID)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to get profile by ID")
		return err
	}

	s.logger.Info().Str("profile_id", id).Str("user_id", p.UserID.String()).Msg("Found profile to update")

	// Update fields - only update fields that are provided
	if len(allergies) > 0 {
		p.Allergens = allergies
	}
	if len(goals) > 0 {
		p.GoalType = goals[0]
	}
	if len(preferences) > 0 {
		p.PreferredFoods = preferences
	}

	// Update isDefault if provided
	if isDefault && !p.IsDefault {
		s.logger.Info().Str("profile_id", id).Msg("Setting profile as default")
		if err := s.profileRepo.SetAsDefault(profileID, p.UserID); err != nil {
			s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to set profile as default")
			return err
		}
		s.logger.Info().Str("profile_id", id).Msg("Profile set as default")
		p.IsDefault = true
	}

	// Add additional fields from the request
	// These fields might be stored in a different table or not used in the current implementation
	// but we'll log them for debugging purposes
	s.logger.Debug().Int("age", age).Str("gender", gender).Float64("weight", weight).Float64("height", height).Msg("Additional profile fields received")

	// Log the profile before update
	s.logger.Debug().Interface("profile", p).Msg("Profile before update")

	// Update the profile - the Update method will preserve the original user ID
	err = s.Update(p)
	if err != nil {
		s.logger.Error().Err(err).Str("profile_id", id).Msg("Failed to update profile")
	}
	return err
}

func (s *profileService) DeleteProfile(ctx context.Context, id string) error {
	profileID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Get existing profile
	p, err := s.profileRepo.GetByID(profileID)
	if err != nil {
		return err
	}

	return s.Delete(profileID, p.UserID)
}

// GetProfilesByUserID returns all profiles for a specific user
func (s *profileService) GetProfilesByUserID(ctx context.Context, userID string) ([]profile.UserProfile, error) {
	s.logger.Info().Str("user_id", userID).Msg("Getting all profiles for user")

	// Parse the user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to parse user ID")
		return nil, err
	}

	// Get all profiles for the user
	profiles, err := s.profileRepo.GetByUserID(uid)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get profiles for user")
		return nil, err
	}

	s.logger.Info().Str("user_id", userID).Int("profile_count", len(profiles)).Msg("Successfully retrieved profiles for user")
	return profiles, nil
}

// GetAllProfiles returns all profiles in the database (for debugging)
func (s *profileService) GetAllProfiles(ctx context.Context) ([]*profile.UserProfile, error) {
	s.logger.Info().Msg("Getting all profiles for debugging")

	// Get all profiles from the database
	return s.profileRepo.GetAll()
}
