-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: ResetUsers :one
DELETE FROM users
RETURNING *; 

-- name: RetrieveUserByEmail :one
SELECT * from users
WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2
WHERE id = $3
RETURNING *;