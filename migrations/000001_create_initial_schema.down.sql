-- Drop tables in reverse order of creation (to handle foreign key constraints)
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS health_conditions;
DROP TABLE IF EXISTS allergens;
DROP TABLE IF EXISTS user_saved_foods;
DROP TABLE IF EXISTS food_ratings;
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;
