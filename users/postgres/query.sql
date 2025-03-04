-- Table: User

-- name: CreateUser :one
INSERT INTO "fianace_manager"."FM_User" (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CheckUserExists :one
SELECT EXISTS (
  SELECT 1 FROM "fianace_manager"."FM_User" WHERE id = $1
) AS user_exists;

-- name: CheckUserEmailExists :one
SELECT EXISTS (
  SELECT 1 FROM "fianace_manager"."FM_User" WHERE email = $1
) AS email_exists;

-- name: GetUserById :one
SELECT * FROM "fianace_manager"."FM_User" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "fianace_manager"."FM_User" WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM "fianace_manager"."FM_User" ORDER BY created_at OFFSET $1 LIMIT $2;

-- name: UpdateUser :exec
UPDATE "fianace_manager"."FM_User" SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "fianace_manager"."FM_User" WHERE id = $1;