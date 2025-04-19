-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    first_name,
    last_name,
    date_of_birth,
    gender,
    height_cm,
    weight_kg,
    activity_level
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    date_of_birth = COALESCE($4, date_of_birth),
    gender = COALESCE($5, gender),
    height_cm = COALESCE($6, height_cm),
    weight_kg = COALESCE($7, weight_kg),
    activity_level = COALESCE($8, activity_level),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
    last_login = NOW()
WHERE id = $1;

-- name: UpdateUserEmailVerification :exec
UPDATE users
SET
    email_verified = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
