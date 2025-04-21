-- Remove seeded data
DELETE FROM allergens WHERE name IN (
    'peanuts', 'tree_nuts', 'milk', 'eggs', 'soy',
    'wheat', 'fish', 'shellfish', 'sesame'
);

DELETE FROM health_conditions WHERE name IN (
    'hypertension', 'type_2_diabetes', 'celiac_disease',
    'high_cholesterol', 'gout'
);