-- Table: Fianace_Manager.User 

-- name: CreateUser :one
INSERT INTO "Fianace_Manager.User" (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CheckUserExists :one
SELECT EXISTS (
  SELECT 1 FROM "Fianace_Manager.User" WHERE id = $1
) AS user_exists;

-- name: CheckUserEmailExists :one
SELECT EXISTS (
  SELECT 1 FROM "Fianace_Manager.User" WHERE email = $1
) AS email_exists;

-- name: GetUserById :one
SELECT * FROM "Fianace_Manager.User" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "Fianace_Manager.User" WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM "Fianace_Manager.User" ORDER BY created_at OFFSET $1 LIMIT $2;

-- name: UpdateUser :exec
UPDATE "Fianace_Manager.User" SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "Fianace_Manager.User" WHERE id = $1;