CREATE TABLE IF NOT EXISTS "BK_User" (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    email VARCHAR(256) NOT NULL UNIQUE,
    password VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "BK_Account" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    id_number VARCHAR(20) NOT NULL UNIQUE,
    balance NUMERIC(100, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "BK_User"(id)
);

CREATE TABLE IF NOT EXISTS "BK_Transaction" (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    balance_after NUMERIC(100, 2) NOT NULL,
    tx_type TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES "BK_Account"(id)
);

CREATE OR REPLACE FUNCTION generate_account_number()
RETURNS VARCHAR(20) AS $$
DECLARE
    result VARCHAR(20);
BEGIN
    result := (
        SELECT STRING_AGG(FLOOR(RANDOM() * 10)::TEXT, '')
        FROM generate_series(1, 20)
    );

    WHILE EXISTS (
        SELECT 1 FROM "BK_Account" WHERE id_number = result
    ) LOOP 
        result := (
            SELECT STRING_AGG(FLOOR(RANDOM() * 10)::TEXT, '')
            FROM generate_series(1, 20)
        );
    END LOOP;

    RETURN result;
END;
$$ LANGUAGE plpgsql;
