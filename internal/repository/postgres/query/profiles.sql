-- name: CreateUserProfile :one
INSERT INTO user_profiles (
    id,
    user_id,
    profile_name,
    is_default,
    health_conditions,
    dietary_restrictions,
    allergens,
    goal_type,
    calorie_target,
    macronutrient_preference,
    disliked_foods,
    preferred_foods,
    cuisine_preferences
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetUserProfileByID :one
SELECT * FROM user_profiles
WHERE id = $1 LIMIT 1;

-- name: CheckProfileExists :one
SELECT EXISTS(
    SELECT 1 FROM user_profiles WHERE id = $1
) AS exists;

-- name: GetProfileByIDDirect :one
SELECT * FROM user_profiles WHERE id = $1;

-- name: GetUserProfiles :many
SELECT * FROM user_profiles
WHERE user_id = $1
ORDER BY is_default DESC, profile_name;

-- name: GetDefaultUserProfile :one
SELECT * FROM user_profiles
WHERE user_id = $1 AND is_default = true
LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET
    profile_name = COALESCE($3, profile_name),
    health_conditions = COALESCE($4, health_conditions),
    dietary_restrictions = COALESCE($5, dietary_restrictions),
    allergens = COALESCE($6, allergens),
    goal_type = COALESCE($7, goal_type),
    calorie_target = COALESCE($8, calorie_target),
    macronutrient_preference = COALESCE($9, macronutrient_preference),
    disliked_foods = COALESCE($10, disliked_foods),
    preferred_foods = COALESCE($11, preferred_foods),
    cuisine_preferences = COALESCE($12, cuisine_preferences),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SetProfileAsDefault :exec
WITH updated_profile AS (
    UPDATE user_profiles
    SET is_default = true
    WHERE user_profiles.id = $1 AND user_profiles.user_id = $2
)
UPDATE user_profiles
SET is_default = false
WHERE user_profiles.user_id = $2 AND user_profiles.id != $1;

-- name: DeleteUserProfile :exec
DELETE FROM user_profiles
WHERE id = $1 AND user_id = $2;
