-- name: CreateUser :one
INSERT INTO "BK_User" (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT username, email, created_at, updated_at
FROM "BK_User"
WHERE id = $1;

-- name: GetUserAccounts :many
SELECT id, id_number, balance
FROM "BK_Account"
WHERE user_id = $1;

-- name: GetAllUsers :many
SELECT id, username
FROM "BK_User";

-- name: CheckUserEmailExists :one
SELECT EXISTS (
  SELECT 1 FROM "BK_User" WHERE email = $1
) AS email_exists;

-- name: UpdateUser :one
UPDATE "BK_User"
SET username = $2, email = $3, password = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateAccount :one
INSERT INTO "BK_Account" (user_id, id_number)
VALUES ($1, generate_account_number())
RETURNING *;

-- name: CheckAccountIDNumberExists :one
SELECT EXISTS (
  SELECT 1 FROM "BK_Account" WHERE id_number = $1
) AS id_number_exists;

-- name: GetAccountByIDNumber :one
SELECT id, user_id, balance, created_at, updated_at
FROM "BK_Account"
WHERE id_number = $1;

-- name: GetAccountBalance :one
SELECT balance
FROM "BK_Account"
WHERE id = $1;

-- name: GetAccountTransactionsByIDNumber :many
SELECT t.id, t.amount, t.tx_type, t.detail, t.created_at
FROM "BK_Transaction" t
JOIN "BK_Account" a ON t.account_id = a.id
WHERE a.id_number = $1
ORDER BY t.created_at DESC;

-- name: GetAllAccounts :many
SELECT id, id_number, user_id
FROM "BK_Account";

-- name: WithdrawFromAccount :one
-- Withdraws the specified amount from the account balance.
-- $1: account ID
-- $2: amount to withdraw
-- $3: Transaction detail
SELECT new_balance::NUMERIC, transaction_id::BIGINT FROM withdraw_from_account($1, $2, $3);

-- name: DepositToAccount :one
-- $1: account ID
-- $2: amount to deposit
-- $3: Transaction detail
SELECT new_balance::NUMERIC, transaction_id::BIGINT FROM deposit_to_account($1, $2, $3);

-- name: CreateTransaction :one
INSERT INTO "BK_Transaction" (
    account_id, 
    amount, 
    balance_after, 
    tx_type, 
    detail
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTransactionByID :one
SELECT id, account_id, amount, balance_after, tx_type, detail, created_at
FROM "BK_Transaction"
WHERE id = $1;
