-- Table: FM_User

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
SELECT 
  u.*,
  ARRAY_AGG(COALESCE(a.id_number, ''))::TEXT[] AS account_numbers
FROM "FM_User" u
LEFT JOIN "FM_Account" a ON u.id = a.user_id
WHERE u.id = $1
GROUP BY u.id;

-- name: GetUserByEmail :one
SELECT 
  u.*,
  ARRAY_AGG(COALESCE(a.id_number, ''))::TEXT[] AS account_numbers
FROM "FM_User" u
LEFT JOIN "FM_Account" a ON u.id = a.user_id
WHERE u.email = $1
GROUP BY u.id;

-- name: GetAllUsers :many
SELECT id FROM "FM_User" ORDER BY created_at;

-- name: UpdateUser :exec
UPDATE "FM_User" SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "FM_User" WHERE id = $1;

-- table: FM_Account

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