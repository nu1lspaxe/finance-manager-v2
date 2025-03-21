CREATE TABLE IF NOT EXISTS "FM_Record_Bank" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,    
    account_number VARCHAR(20) NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    transaction_id BIGINT NOT NULL, -- Bank System Transaction id
    transaction_date TIMESTAMPTZ NOT NULL,
    record_type TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_record_bank_user_id_date_account 
ON "FM_Record_Bank" (user_id, account_number, transaction_date);

CREATE INDEX IF NOT EXISTS idx_record_bank_user_id_type 
ON "FM_Record_Bank" (user_id, record_type);