-- Table: FM_Record_Bank

-- name: CreateBankRecordsBulk :many
INSERT INTO "FM_Record_Bank" (
    user_id, 
    account_number, 
    amount, 
    record_type, 
    transaction_id, 
    transaction_date, 
    detail
) VALUES (
    $1::bigint,                -- single user_id
    $2::text,                  -- single account_number
    UNNEST($3::numeric[]),     -- amount array
    UNNEST($4::text[]),        -- record_type array
    UNNEST($5::bigint[]),      -- transaction_id array
    UNNEST($6::TIMESTAMPTZ[]), -- transaction_date array
    UNNEST($7::text[])         -- detail array
)
RETURNING *;

-- name: GetBankRecord :one
SELECT * FROM "FM_Record_Bank" WHERE id = $1;

-- name: GetUserBankRecords :many
SELECT * FROM "FM_Record_Bank" WHERE user_id = $1 AND account_number = $2;

-- name: GetUserBankRecordsWithPeriod :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND created_at BETWEEN $3 AND $4;

-- name: GetUserBankRecordsFromDate :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND created_at >= $3;

-- name: GetUserBankRecordsToDate :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND created_at <= $3;

-- name: GetUserBankRecordsByType :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND record_type = $3;

-- name: GetUserBankRecordsByTypeWithPeriod :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND record_type = $3 AND created_at BETWEEN $4 AND $5;

-- name: GetUserBankRecordsByTypeFromDate :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND record_type = $3 AND created_at >= $4;

-- name: GetUserBankRecordsByTypeToDate :many
SELECT * FROM "FM_Record_Bank"
WHERE user_id = $1 AND account_number = $2 AND record_type = $3 AND created_at <= $4;

-- name: GetExistingBankRecords :many
SELECT transaction_id FROM "FM_Record_Bank" WHERE transaction_id = ANY($1::bigint[]);
