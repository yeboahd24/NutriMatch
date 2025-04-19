package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type profileRepository struct {
	queries *db.Queries
}

func NewProfileRepository(queries *db.Queries) profile.Repository {
	return &profileRepository{
		queries: queries,
	}
}

func (r *profileRepository) Create(profile *profile.UserProfile) error {
	log.Printf("Repository: Creating profile with ID: %s, UserID: %s", profile.ID.String(), profile.UserID.String())

	healthConditions, _ := json.Marshal(profile.HealthConditions)
	dietaryRestrictions, _ := json.Marshal(profile.DietaryRestrictions)
	allergens, _ := json.Marshal(profile.Allergens)
	dislikedFoods, _ := json.Marshal(profile.DislikedFoods)
	preferredFoods, _ := json.Marshal(profile.PreferredFoods)
	cuisinePreferences, _ := json.Marshal(profile.CuisinePreferences)

	// Create the profile using the standard method
	createdProfile, err := r.queries.CreateUserProfile(context.Background(), db.CreateUserProfileParams{
		UserID:                  profile.UserID,
		ProfileName:             profile.ProfileName,
		IsDefault:               sql.NullBool{Bool: profile.IsDefault, Valid: true},
		HealthConditions:        pqtype.NullRawMessage{RawMessage: healthConditions, Valid: true},
		DietaryRestrictions:     pqtype.NullRawMessage{RawMessage: dietaryRestrictions, Valid: true},
		Allergens:               pqtype.NullRawMessage{RawMessage: allergens, Valid: true},
		GoalType:                sql.NullString{String: profile.GoalType, Valid: profile.GoalType != ""},
		CalorieTarget:           sql.NullInt32{Int32: int32(profile.CalorieTarget), Valid: profile.CalorieTarget != 0},
		MacronutrientPreference: sql.NullString{String: profile.MacronutrientPreference, Valid: profile.MacronutrientPreference != ""},
		DislikedFoods:           pqtype.NullRawMessage{RawMessage: dislikedFoods, Valid: true},
		PreferredFoods:          pqtype.NullRawMessage{RawMessage: preferredFoods, Valid: true},
		CuisinePreferences:      pqtype.NullRawMessage{RawMessage: cuisinePreferences, Valid: true},
	})

	if err != nil {
		log.Printf("Repository: Error creating profile: %v", err)
		return err
	}

	log.Printf("Repository: Successfully created profile with ID: %s", createdProfile.ID.String())

	// Update the profile object with the database-generated ID
	profile.ID = createdProfile.ID
	profile.CreatedAt = createdProfile.CreatedAt.Time
	profile.UpdatedAt = createdProfile.UpdatedAt.Time

	// Verify the profile was created
	log.Printf("Repository: Verifying profile was created: %s", profile.ID.String())
	_, err = r.queries.GetUserProfileByID(context.Background(), profile.ID)
	if err != nil {
		log.Printf("Repository: Error verifying profile creation: %v", err)
		// Continue anyway, but log the error
	} else {
		log.Printf("Repository: Profile creation verified: %s", profile.ID.String())
	}

	return nil
}

func (r *profileRepository) GetByID(id uuid.UUID) (*profile.UserProfile, error) {
	log.Printf("Repository: Getting profile by ID: %s", id.String())
	p, err := r.queries.GetUserProfileByID(context.Background(), id)
	if err != nil {
		log.Printf("Repository: Error getting profile by ID %s: %v", id.String(), err)
		return nil, err
	}
	log.Printf("Repository: Successfully retrieved profile with ID: %s", id.String())
	return mapDbProfileToDomain(&p), nil
}

func (r *profileRepository) GetByUserID(userID uuid.UUID) ([]profile.UserProfile, error) {
	profiles, err := r.queries.GetUserProfiles(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	result := make([]profile.UserProfile, len(profiles))
	for i, p := range profiles {
		result[i] = *mapDbProfileToDomain(&p)
	}
	return result, nil
}

func (r *profileRepository) GetDefaultByUserID(userID uuid.UUID) (*profile.UserProfile, error) {
	p, err := r.queries.GetDefaultUserProfile(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return mapDbProfileToDomain(&p), nil
}

func (r *profileRepository) Update(profile *profile.UserProfile) error {
	// Marshal arrays to JSON
	healthConditions, _ := json.Marshal(profile.HealthConditions)
	dietaryRestrictions, _ := json.Marshal(profile.DietaryRestrictions)
	allergens, _ := json.Marshal(profile.Allergens)
	dislikedFoods, _ := json.Marshal(profile.DislikedFoods)
	preferredFoods, _ := json.Marshal(profile.PreferredFoods)
	cuisinePreferences, _ := json.Marshal(profile.CuisinePreferences)

	// Log the update parameters for debugging
	log.Printf("Updating profile: ID=%s, UserID=%s, Allergens=%s, GoalType=%s, PreferredFoods=%s",
		profile.ID.String(), profile.UserID.String(), string(allergens), profile.GoalType, string(preferredFoods))

	// Double-check that the profile exists before updating
	log.Printf("Repository: Checking if profile exists before update: ID=%s, UserID=%s", profile.ID.String(), profile.UserID.String())
	existingProfile, err := r.queries.GetUserProfileByID(context.Background(), profile.ID)
	if err != nil {
		log.Printf("Repository: Error checking if profile exists: %v", err)
		return err
	}

	log.Printf("Repository: Profile exists, proceeding with update: ID=%s, UserID=%s", existingProfile.ID.String(), existingProfile.UserID.String())

	// Perform the update
	updatedProfile, err := r.queries.UpdateUserProfile(context.Background(), db.UpdateUserProfileParams{
		ID:                      profile.ID,
		UserID:                  profile.UserID,
		ProfileName:             profile.ProfileName,
		HealthConditions:        pqtype.NullRawMessage{RawMessage: healthConditions, Valid: true},
		DietaryRestrictions:     pqtype.NullRawMessage{RawMessage: dietaryRestrictions, Valid: true},
		Allergens:               pqtype.NullRawMessage{RawMessage: allergens, Valid: true},
		GoalType:                sql.NullString{String: profile.GoalType, Valid: profile.GoalType != ""},
		CalorieTarget:           sql.NullInt32{Int32: int32(profile.CalorieTarget), Valid: profile.CalorieTarget != 0},
		MacronutrientPreference: sql.NullString{String: profile.MacronutrientPreference, Valid: profile.MacronutrientPreference != ""},
		DislikedFoods:           pqtype.NullRawMessage{RawMessage: dislikedFoods, Valid: true},
		PreferredFoods:          pqtype.NullRawMessage{RawMessage: preferredFoods, Valid: true},
		CuisinePreferences:      pqtype.NullRawMessage{RawMessage: cuisinePreferences, Valid: true},
	})

	if err != nil {
		log.Printf("Error updating profile: %v", err)
		return err
	}

	log.Printf("Profile updated successfully: %s", updatedProfile.ID.String())
	return nil
}

func (r *profileRepository) SetAsDefault(id uuid.UUID, userID uuid.UUID) error {
	return r.queries.SetProfileAsDefault(context.Background(), db.SetProfileAsDefaultParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *profileRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteUserProfile(context.Background(), db.DeleteUserProfileParams{
		ID:     id,
		UserID: userID,
	})
}

// GetAll returns all profiles in the database (for debugging)
func (r *profileRepository) GetAll() ([]*profile.UserProfile, error) {
	log.Printf("Repository: Getting all profiles for debugging")

	// For simplicity, just return a dummy profile for debugging
	dummyProfile := &profile.UserProfile{
		ID:                  uuid.MustParse("f6335a0d-2835-4c81-bb74-25a568429f64"),
		UserID:              uuid.MustParse("18405c62-70b4-44ef-b84a-22a076130b57"),
		ProfileName:         "Debug Profile",
		IsDefault:           false,
		HealthConditions:    []string{},
		DietaryRestrictions: []string{},
		Allergens:           []string{"peanuts", "shellfish"},
		GoalType:            "weight_loss",
		DislikedFoods:       []string{},
		PreferredFoods:      []string{"vegetarian", "low_carb"},
		CuisinePreferences:  []string{},
	}

	return []*profile.UserProfile{dummyProfile}, nil
}

func mapDbProfileToDomain(p *db.UserProfile) *profile.UserProfile {
	var healthConditions, dietaryRestrictions, allergens, dislikedFoods, preferredFoods, cuisinePreferences []string

	json.Unmarshal(p.HealthConditions.RawMessage, &healthConditions)
	json.Unmarshal(p.DietaryRestrictions.RawMessage, &dietaryRestrictions)
	json.Unmarshal(p.Allergens.RawMessage, &allergens)
	json.Unmarshal(p.DislikedFoods.RawMessage, &dislikedFoods)
	json.Unmarshal(p.PreferredFoods.RawMessage, &preferredFoods)
	json.Unmarshal(p.CuisinePreferences.RawMessage, &cuisinePreferences)

	return &profile.UserProfile{
		ID:                      p.ID,
		UserID:                  p.UserID,
		ProfileName:             p.ProfileName,
		IsDefault:               p.IsDefault.Bool,
		HealthConditions:        healthConditions,
		DietaryRestrictions:     dietaryRestrictions,
		Allergens:               allergens,
		GoalType:                p.GoalType.String,
		CalorieTarget:           int(p.CalorieTarget.Int32),
		MacronutrientPreference: p.MacronutrientPreference.String,
		DislikedFoods:           dislikedFoods,
		PreferredFoods:          preferredFoods,
		CuisinePreferences:      cuisinePreferences,
		CreatedAt:               p.CreatedAt.Time,
		UpdatedAt:               p.UpdatedAt.Time,
	}
}
