-- Extract additional allergens from ingredient_analysis that aren't in seed data
INSERT INTO allergens (name, description, common_names)
SELECT DISTINCT ON (allergen_name) 
    allergen_name,
    'Common allergen found in food products',
    jsonb_build_object(
        'alternatives',
        jsonb_agg(DISTINCT 
            CASE 
                WHEN allergen_info->>'alternate_names' IS NOT NULL THEN allergen_info->>'alternate_names'
                WHEN allergen_info->>'variants' IS NOT NULL THEN allergen_info->>'variants'
                ELSE NULL 
            END
        ) FILTER (WHERE 
            allergen_info->>'alternate_names' IS NOT NULL OR 
            allergen_info->>'variants' IS NOT NULL
        )
    )
FROM foods,
     jsonb_array_elements(ingredient_analysis->'allergens') as allergen_info,
     jsonb_array_elements_text(
         CASE 
             WHEN allergen_info->>'name' IS NOT NULL THEN jsonb_build_array(allergen_info->>'name')
             WHEN allergen_info->>'allergen' IS NOT NULL THEN jsonb_build_array(allergen_info->>'allergen')
             ELSE '[]'::jsonb
         END
     ) as allergen_name
WHERE ingredient_analysis->'allergens' IS NOT NULL
  AND allergen_name IS NOT NULL
  AND allergen_name NOT IN (
    'peanuts', 'tree_nuts', 'milk', 'eggs', 'soy', 
    'wheat', 'fish', 'shellfish', 'sesame'
  )
GROUP BY allergen_name
ON CONFLICT (name) DO NOTHING;

-- Update existing health conditions with additional data
UPDATE health_conditions 
SET nutrient_restrictions = nutrient_restrictions || 
    CASE name
        WHEN 'type_2_diabetes' THEN 
            '{"sugar": {"max": 25, "unit": "g"}, "carbohydrates": {"max": 130, "unit": "g"}}'::jsonb
        ELSE nutrient_restrictions
    END,
    nutrient_recommendations = nutrient_recommendations || 
    CASE name
        WHEN 'type_2_diabetes' THEN 
            '{"fiber": {"min": 25, "unit": "g"}, "protein": {"min": 50, "unit": "g"}}'::jsonb
        ELSE nutrient_recommendations
    END
WHERE name IN ('type_2_diabetes');

-- Insert new health conditions
INSERT INTO health_conditions (name, description, nutrient_restrictions, nutrient_recommendations)
VALUES
('weight_management', 'Balanced nutrition for healthy weight',
    '{"calories": {"max": 600, "unit": "kcal"}, "added_sugars": {"max": 10, "unit": "g"}}',
    '{"protein": {"min": 20, "unit": "g"}, "fiber": {"min": 5, "unit": "g"}}'
)
ON CONFLICT (name) DO NOTHING;