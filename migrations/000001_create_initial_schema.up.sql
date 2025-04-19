-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    date_of_birth DATE,
    gender VARCHAR(20),
    height_cm NUMERIC(5,2),
    weight_kg NUMERIC(5,2),
    activity_level VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE,
    account_status VARCHAR(20) DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret VARCHAR(255)
);

CREATE INDEX idx_users_email ON users(email);

-- Create user_profiles table
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_name VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    health_conditions JSONB DEFAULT '[]',
    dietary_restrictions JSONB DEFAULT '[]',
    allergens JSONB DEFAULT '[]',
    goal_type VARCHAR(50),
    calorie_target INTEGER,
    macronutrient_preference VARCHAR(50),
    disliked_foods JSONB DEFAULT '[]',
    preferred_foods JSONB DEFAULT '[]',
    cuisine_preferences JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- Create foods table
CREATE TABLE foods (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alternate_names JSONB,
    description TEXT,
    food_type VARCHAR(50),
    source JSONB,
    serving JSONB,
    nutrition_100g JSONB,
    ean_13 VARCHAR(13),
    labels JSONB,
    package_size JSONB,
    ingredients TEXT,
    ingredient_analysis JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_foods_name ON foods(name);
CREATE INDEX idx_foods_food_type ON foods(food_type);
CREATE INDEX idx_foods_ean_13 ON foods(ean_13);
CREATE INDEX idx_foods_nutrition ON foods USING GIN (nutrition_100g);
CREATE INDEX idx_foods_labels ON foods USING GIN (labels);
CREATE INDEX idx_foods_ingredient_analysis ON foods USING GIN (ingredient_analysis);

-- Create food_ratings table
CREATE TABLE food_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_id VARCHAR(50) NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, food_id)
);

CREATE INDEX idx_food_ratings_user_id ON food_ratings(user_id);
CREATE INDEX idx_food_ratings_food_id ON food_ratings(food_id);

-- Create user_saved_foods table
CREATE TABLE user_saved_foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_id VARCHAR(50) NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    list_type VARCHAR(50) NOT NULL, -- favorites, shopping_list, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, food_id, list_type)
);

CREATE INDEX idx_user_saved_foods_user_id ON user_saved_foods(user_id);
CREATE INDEX idx_user_saved_foods_food_id ON user_saved_foods(food_id);

-- Create allergens reference table
CREATE TABLE allergens (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    common_names JSONB, -- alternative names for the allergen
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create health_conditions reference table
CREATE TABLE health_conditions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    nutrient_restrictions JSONB, -- nutrients to limit for this condition
    nutrient_recommendations JSONB, -- nutrients to increase for this condition
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create audit_logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create refresh_tokens table for JWT token management
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
