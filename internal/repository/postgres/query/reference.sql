-- name: ListAllergens :many
SELECT id, name, description, common_names, created_at
FROM allergens
ORDER BY name;

-- name: ListHealthConditions :many
SELECT id, name, description, nutrient_restrictions, nutrient_recommendations, created_at
FROM health_conditions
ORDER BY name;