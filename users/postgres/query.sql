-- Table: User 

-- name: CreateUser :one
INSERT INTO "User" (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CheckUserExists :one
SELECT EXISTS (
  SELECT 1 FROM "User" WHERE id = $1
) AS user_exists;

-- name: CheckUserEmailExists :one
SELECT EXISTS (
  SELECT 1 FROM "User" WHERE email = $1
) AS email_exists;

-- name: GetUserById :one
SELECT * FROM "User" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "User" WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM "User" ORDER BY created_at;

-- name: UpdateUser :exec
UPDATE "User" SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "User" WHERE id = $1;