-- Seed allergens table
INSERT INTO allergens (name, description, common_names) VALUES
('peanuts', 'A legume commonly causing severe allergic reactions', '{"alternatives": ["ground nuts", "goober peas", "arachis hypogaea"], "scientific": "Arachis hypogaea"}'),
('tree_nuts', 'Various nuts grown on trees', '{"alternatives": ["almonds", "walnuts", "cashews", "pecans", "pistachios"], "types": ["almond", "walnut", "cashew", "pecan", "pistachio"]}'),
('milk', 'Dairy milk and milk products', '{"alternatives": ["dairy", "lactose", "whey", "casein"], "products": ["cheese", "yogurt", "butter", "cream"]}'),
('eggs', 'Chicken eggs and egg products', '{"alternatives": ["egg white", "egg yolk", "albumin"], "products": ["mayonnaise", "meringue"]}'),
('soy', 'Soybeans and soy products', '{"alternatives": ["soya", "edamame"], "products": ["tofu", "tempeh", "miso"]}'),
('wheat', 'Wheat and wheat-containing grains', '{"alternatives": ["gluten", "spelt", "semolina"], "products": ["flour", "bread", "pasta"]}'),
('fish', 'Various species of finned fish', '{"types": ["salmon", "tuna", "cod", "bass"], "products": ["fish sauce", "worcestershire sauce"]}'),
('shellfish', 'Crustaceans and mollusks', '{"types": ["shrimp", "crab", "lobster", "clams", "mussels"], "products": ["seafood flavoring"]}'),
('sesame', 'Sesame seeds and sesame products', '{"alternatives": ["tahini", "sesamol", "gingelly"], "products": ["hummus", "halvah"]}');

-- Seed health conditions table
INSERT INTO health_conditions (name, description, nutrient_restrictions, nutrient_recommendations) VALUES
('hypertension', 'High blood pressure condition requiring dietary management', 
    '{"sodium": {"max": 2000, "unit": "mg"}, "saturated_fat": {"max": 20, "unit": "g"}}',
    '{"potassium": {"min": 3500, "unit": "mg"}, "magnesium": {"min": 400, "unit": "mg"}, "calcium": {"min": 1000, "unit": "mg"}}'
),
('type_2_diabetes', 'Type 2 diabetes mellitus requiring blood sugar control',
    '{"sugar": {"max": 25, "unit": "g"}, "refined_carbs": {"max": 50, "unit": "g"}}',
    '{"fiber": {"min": 25, "unit": "g"}, "protein": {"min": 50, "unit": "g"}}'
),
('celiac_disease', 'Autoimmune disorder requiring strict gluten avoidance',
    '{"gluten": {"max": 0, "unit": "mg"}}',
    '{"iron": {"min": 18, "unit": "mg"}, "vitamin_b12": {"min": 2.4, "unit": "mcg"}, "fiber": {"min": 25, "unit": "g"}}'
),
('high_cholesterol', 'Elevated blood cholesterol levels requiring dietary management',
    '{"saturated_fat": {"max": 15, "unit": "g"}, "trans_fat": {"max": 0, "unit": "g"}, "cholesterol": {"max": 200, "unit": "mg"}}',
    '{"fiber": {"min": 25, "unit": "g"}, "omega_3": {"min": 1000, "unit": "mg"}}'
),
('gout', 'Inflammatory arthritis requiring purine restriction',
    '{"purines": {"max": 200, "unit": "mg"}}',
    '{"vitamin_c": {"min": 500, "unit": "mg"}, "water": {"min": 2000, "unit": "ml"}}'
);