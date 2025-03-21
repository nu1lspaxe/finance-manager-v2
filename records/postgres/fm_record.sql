CREATE TABLE IF NOT EXISTS "FM_Record" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    transaction_date DATE NOT NULL,
    record_type TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_record_user_id_date ON "FM_Record" (user_id, transaction_date);

CREATE INDEX IF NOT EXISTS idx_record_user_id_type ON "FM_Record" (user_id, record_type);