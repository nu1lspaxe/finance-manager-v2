CREATE TABLE IF NOT EXISTS "Record" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount NUMERIC(9,2) NOT NULL,
    transaction_date DATE NOT NULL,
    record_type TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_record_user_id_transaction_date ON "Record" (user_id, transaction_date);

CREATE INDEX idx_record_user_id_type ON "Record" (user_id, record_type);