-- Remove derived allergens (excluding seed data)
DELETE FROM allergens WHERE name NOT IN (
    'peanuts', 'tree_nuts', 'milk', 'eggs', 'soy', 
    'wheat', 'fish', 'shellfish', 'sesame'
);

-- Reset updated health condition data
UPDATE health_conditions 
SET nutrient_restrictions = nutrient_restrictions - 'carbohydrates',
    nutrient_recommendations = nutrient_recommendations
WHERE name = 'type_2_diabetes';

-- Remove derived health condition
DELETE FROM health_conditions WHERE name = 'weight_management';