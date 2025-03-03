-- Table: Record

-- name: CreateRecord :one
INSERT INTO "Fianace_Manager.Record" (user_id, amount, transaction_date, record_type, detail)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRecord :one
SELECT * FROM "Fianace_Manager.Record" WHERE id = $1;

-- name: GetUserRecords :many
SELECT * FROM "Fianace_Manager.Record" WHERE user_id = $1;

-- name: GetUserRecordsWithPeriod :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND transaction_date BETWEEN $2 AND $3;

-- name: GetUserRecordsFromDate :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND transaction_date >= $2;

-- name: GetUserRecordsToDate :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND transaction_date <= $2;

-- name: GetUserRecordsByType :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND record_type = $2;

-- name: GetUserRecordsByTypeWithPeriod :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND record_type = $2 AND transaction_date BETWEEN $3 AND $4;

-- name: GetUserRecordsByTypeFromDate :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND record_type = $2 AND transaction_date >= $3;

-- name: GetUserRecordsByTypeToDate :many
SELECT * FROM "Fianace_Manager.Record"
WHERE user_id = $1 AND record_type = $2 AND transaction_date <= $3;

-- name: UpdateRecord :exec
UPDATE "Fianace_Manager.Record" SET amount = $2, transaction_date = $3, record_type = $4, detail = $5, updated_time = NOW() WHERE id = $1;

-- name: DeleteRecord :exec
DELETE FROM "Fianace_Manager.Record" WHERE id = $1;
