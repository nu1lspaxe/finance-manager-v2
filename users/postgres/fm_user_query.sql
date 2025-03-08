-- Table: User

-- name: CreateUser :one
INSERT INTO "FM_User" (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CheckUserExists :one
SELECT EXISTS (
  SELECT 1 FROM "FM_User" WHERE id = $1
) AS user_exists;

-- name: CheckUserEmailExists :one
SELECT EXISTS (
  SELECT 1 FROM "FM_User" WHERE email = $1
) AS email_exists;

-- name: GetUserById :one
SELECT * FROM "FM_User" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "FM_User" WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM "FM_User" ORDER BY created_at OFFSET $1 LIMIT $2;

-- name: UpdateUser :exec
UPDATE "FM_User" SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "FM_User" WHERE id = $1;

-- name: AddAccount :one
INSERT INTO "FM_Account" (user_id, id_number, balance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserAccounts :many
SELECT * FROM "FM_Account" WHERE user_id = $1;

-- name: CheckAccountExists :one
SELECT EXISTS (
  SELECT 1 FROM "FM_Account" WHERE user_id = $1 AND id_number = $2
) AS account_exists;

-- name: UpdateAccountBalance :exec
UPDATE "FM_Account" SET balance = $2, updated_at = NOW() WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM "FM_Account" WHERE id_number = $1;