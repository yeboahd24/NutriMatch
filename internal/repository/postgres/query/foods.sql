-- name: CreateFood :one
INSERT INTO foods (
    id,
    name,
    alternate_names,
    description,
    food_type,
    source,
    serving,
    nutrition_100g,
    ean_13,
    labels,
    package_size,
    ingredients,
    ingredient_analysis
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetFoodByID :one
SELECT * FROM foods
WHERE id = $1 LIMIT 1;

-- name: GetFoodByEAN13 :one
SELECT * FROM foods
WHERE ean_13 = $1 LIMIT 1;

-- name: ListFoods :many
SELECT * FROM foods
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListFoodsByType :many
SELECT * FROM foods
WHERE food_type = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: SearchFoodsByName :many
SELECT * FROM foods
WHERE name ILIKE '%' || $1 || '%'
   OR EXISTS (
       SELECT 1
       FROM jsonb_array_elements_text(alternate_names) AS alt_name
       WHERE alt_name ILIKE '%' || $1 || '%'
   )
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CountFoods :one
SELECT COUNT(*) FROM foods;

-- name: DeleteFood :exec
DELETE FROM foods
WHERE id = $1;
