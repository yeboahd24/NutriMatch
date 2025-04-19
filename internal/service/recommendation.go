package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/domain/recommendation"
)

type recommendationService struct {
	foodRepo    food.Repository
	profileRepo profile.Repository
	logger      zerolog.Logger
}

func NewRecommendationService(
	foodRepo food.Repository,
	profileRepo profile.Repository,
	logger zerolog.Logger,
) RecommendationService {
	return &recommendationService{
		foodRepo:    foodRepo,
		profileRepo: profileRepo,
		logger:      logger,
	}
}

func (s *recommendationService) GetRecommendations(userID uuid.UUID, req recommendation.RecommendationRequest) (*recommendation.RecommendationResponse, error) {
	// Get user profile
	var userProfile *profile.UserProfile
	var err error

	if req.ProfileID != nil {
		userProfile, err = s.profileRepo.GetByID(*req.ProfileID)
		if err != nil {
			return nil, err
		}
		// Verify ownership
		if userProfile.UserID != userID {
			return nil, profile.ErrUnauthorized
		}
	} else {
		userProfile, err = s.profileRepo.GetDefaultByUserID(userID)
		if err != nil {
			return nil, err
		}
	}

	// Generate rules from profile
	rules, err := s.GenerateRulesFromProfile(userProfile)
	if err != nil {
		return nil, err
	}

	// Add custom rules if provided
	if len(req.CustomRules) > 0 {
		rules = append(rules, req.CustomRules...)
	}

	// Get all foods (paginated)
	limit := 100
	if req.Limit > 0 {
		limit = req.Limit
	}
	foods, err := s.foodRepo.List(limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Apply rules to filter foods
	var filteredFoods []food.Food
	for _, f := range foods {
		if s.applyRules(f, rules) {
			filteredFoods = append(filteredFoods, f)
		}
	}

	// Get total count for pagination
	totalCount := len(filteredFoods)

	return &recommendation.RecommendationResponse{
		Foods:        filteredFoods,
		TotalCount:   totalCount,
		AppliedRules: rules,
	}, nil
}

func (s *recommendationService) GetAlternatives(userID uuid.UUID, foodID string, limit int) ([]food.Food, error) {
	// Get the original food
	originalFood, err := s.foodRepo.GetByID(foodID)
	if err != nil {
		return nil, err
	}

	// Get foods of the same type
	foods, err := s.foodRepo.ListByType(originalFood.FoodType, limit+1, 0)
	if err != nil {
		return nil, err
	}

	// Filter out the original food and limit results
	var alternatives []food.Food
	for _, f := range foods {
		if f.ID != foodID {
			alternatives = append(alternatives, f)
			if len(alternatives) >= limit {
				break
			}
		}
	}

	return alternatives, nil
}

func (s *recommendationService) GenerateRulesFromProfile(profile *profile.UserProfile) ([]recommendation.Rule, error) {
	var rules []recommendation.Rule

	// Add allergen rules (highest priority)
	for _, allergen := range profile.Allergens {
		rules = append(rules, recommendation.Rule{
			Type:      "allergen",
			Operation: "exclude",
			Target:    allergen,
			Priority:  100,
		})
	}

	// Add dietary restriction rules (high priority)
	for _, restriction := range profile.DietaryRestrictions {
		rules = append(rules, recommendation.Rule{
			Type:      "dietary",
			Operation: "exclude",
			Target:    restriction,
			Priority:  90,
		})
	}

	// Add disliked foods rules (medium priority)
	for _, disliked := range profile.DislikedFoods {
		rules = append(rules, recommendation.Rule{
			Type:      "preference",
			Operation: "exclude",
			Target:    disliked,
			Priority:  50,
		})
	}

	// Add preferred foods rules (medium priority)
	for _, preferred := range profile.PreferredFoods {
		rules = append(rules, recommendation.Rule{
			Type:      "preference",
			Operation: "include",
			Target:    preferred,
			Priority:  60,
		})
	}

	// Add calorie target rule if set (medium priority)
	if profile.CalorieTarget > 0 {
		rules = append(rules, recommendation.Rule{
			Type:      "nutrient",
			Operation: "max",
			Target:    "calories",
			Value:     profile.CalorieTarget,
			Priority:  70,
		})
	}

	// Add cuisine preference rules (low priority)
	for _, cuisine := range profile.CuisinePreferences {
		rules = append(rules, recommendation.Rule{
			Type:      "cuisine",
			Operation: "prefer",
			Target:    cuisine,
			Priority:  40,
		})
	}

	return rules, nil
}

func (s *recommendationService) applyRules(food food.Food, rules []recommendation.Rule) bool {
	// Sort rules by priority (highest first)
	sortRulesByPriority(rules)

	for _, rule := range rules {
		switch rule.Operation {
		case "exclude":
			if s.matchesExclusionRule(food, rule) {
				return false
			}
		case "include":
			if s.matchesInclusionRule(food, rule) {
				return true
			}
		case "max":
			if !s.matchesMaxRule(food, rule) {
				return false
			}
		case "min":
			if !s.matchesMinRule(food, rule) {
				return false
			}
		case "prefer":
			// Preference rules don't exclude foods, they just affect sorting
			continue
		}
	}

	return true
}

func (s *recommendationService) matchesExclusionRule(food food.Food, rule recommendation.Rule) bool {
	switch rule.Type {
	case "allergen":
		// Check ingredients for allergen
		return containsString(food.Labels, rule.Target) ||
			containsString(food.AlternateNames, rule.Target) ||
			contains(food.Ingredients, rule.Target)
	case "dietary":
		// Check if food violates dietary restriction
		return violatesDietaryRestriction(food, rule.Target)
	case "preference":
		// Check if food matches disliked item
		return matchesPreference(food, rule.Target)
	default:
		return false
	}
}

func (s *recommendationService) matchesInclusionRule(food food.Food, rule recommendation.Rule) bool {
	switch rule.Type {
	case "preference":
		return matchesPreference(food, rule.Target)
	default:
		return false
	}
}

func (s *recommendationService) matchesMaxRule(food food.Food, rule recommendation.Rule) bool {
	switch rule.Type {
	case "nutrient":
		value, ok := getNutrientValue(food, rule.Target)
		if !ok {
			return true // If nutrient info not available, don't exclude
		}
		maxValue, ok := rule.Value.(int)
		if !ok {
			return true // If rule value invalid, don't exclude
		}
		return value <= float64(maxValue)
	default:
		return true
	}
}

func (s *recommendationService) matchesMinRule(food food.Food, rule recommendation.Rule) bool {
	switch rule.Type {
	case "nutrient":
		value, ok := getNutrientValue(food, rule.Target)
		if !ok {
			return true // If nutrient info not available, don't exclude
		}
		minValue, ok := rule.Value.(int)
		if !ok {
			return true // If rule value invalid, don't exclude
		}
		return value >= float64(minValue)
	default:
		return true
	}
}

// Helper functions

func sortRulesByPriority(rules []recommendation.Rule) {
	// Sort rules by priority in descending order
	for i := 0; i < len(rules)-1; i++ {
		for j := i + 1; j < len(rules); j++ {
			if rules[i].Priority < rules[j].Priority {
				rules[i], rules[j] = rules[j], rules[i]
			}
		}
	}
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func contains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func violatesDietaryRestriction(food food.Food, restriction string) bool {
	// Implementation depends on how dietary restrictions are stored and checked
	// This is a simplified example
	return containsString(food.Labels, restriction)
}

func matchesPreference(food food.Food, preference string) bool {
	return containsString(food.Labels, preference) ||
		containsString(food.AlternateNames, preference) ||
		contains(food.Name, preference)
}

func getNutrientValue(food food.Food, nutrient string) (float64, bool) {
	// Nutrition100g is already a map[string]interface{}
	if value, ok := food.Nutrition100g[nutrient]; ok {
		// Try to convert the value to float64
		switch v := value.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

// Adapter methods to implement the service.RecommendationService interface
func (s *recommendationService) GetDailyRecommendations(ctx context.Context, profileID string, limit int) ([]food.Food, error) {
	pid, err := uuid.Parse(profileID)
	if err != nil {
		return nil, err
	}

	// Get the profile
	profile, err := s.profileRepo.GetByID(pid)
	if err != nil {
		return nil, err
	}

	// Create a recommendation request
	req := recommendation.RecommendationRequest{
		ProfileID: &pid,
		Limit:     limit,
		Offset:    0,
	}

	// Get recommendations
	resp, err := s.GetRecommendations(profile.UserID, req)
	if err != nil {
		return nil, err
	}

	return resp.Foods, nil
}

func (s *recommendationService) GetMealPlanRecommendations(ctx context.Context, profileID string, days int) (*recommendation.MealPlan, error) {
	pid, err := uuid.Parse(profileID)
	if err != nil {
		return nil, err
	}

	// Get the profile
	profile, err := s.profileRepo.GetByID(pid)
	if err != nil {
		return nil, err
	}

	// Create a meal plan
	mealPlan := &recommendation.MealPlan{
		ProfileID: profileID,
		Days:      make([]recommendation.DailyPlan, days),
		TotalDays: days,
	}

	// For each day, create a daily plan with meals
	for i := 0; i < days; i++ {
		dailyPlan := recommendation.DailyPlan{
			Date:  "Day " + strconv.Itoa(i+1),     // Simplified date representation
			Meals: make([]recommendation.Meal, 3), // Breakfast, lunch, dinner
		}

		// For each meal type, get recommendations
		mealTypes := []string{"breakfast", "lunch", "dinner"}
		for j, mealType := range mealTypes {
			// Create a recommendation request with meal type as a custom rule
			req := recommendation.RecommendationRequest{
				ProfileID: &pid,
				CustomRules: []recommendation.Rule{
					{
						Type:      "meal_type",
						Operation: "include",
						Target:    mealType,
						Priority:  100,
					},
				},
				Limit:  3, // 3 foods per meal
				Offset: 0,
			}

			// Get recommendations
			resp, err := s.GetRecommendations(profile.UserID, req)
			if err != nil {
				return nil, err
			}

			dailyPlan.Meals[j] = recommendation.Meal{
				Type:  mealType,
				Foods: resp.Foods,
			}
		}

		mealPlan.Days[i] = dailyPlan
	}

	return mealPlan, nil
}

func (s *recommendationService) GetFoodAlternatives(ctx context.Context, foodID string, limit int) ([]food.Food, error) {
	// We can use the existing GetAlternatives method, but we need to provide a dummy userID
	// In a real implementation, we would get the user ID from the context
	dummyUserID := uuid.New()
	return s.GetAlternatives(dummyUserID, foodID, limit)
}
