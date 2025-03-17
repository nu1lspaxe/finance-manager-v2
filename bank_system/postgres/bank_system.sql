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

-- Generate a unique account number
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

-- Get account transactions by ID number
CREATE OR REPLACE FUNCTION get_account_transactions_by_id_number(p_id_number VARCHAR(20))
RETURNS TABLE (
    id BIGINT,
    amount NUMERIC(20,2),
    tx_type TEXT,
    detail TEXT,
    created_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.tx_type, t.detail, t.created_at
    FROM "BK_Transaction" t
    JOIN "BK_Account" a ON t.account_id = a.id
    WHERE a.id_number = p_id_number
    ORDER BY t.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Withdraw from account
CREATE OR REPLACE FUNCTION withdraw_from_account(account_id BIGINT, amount NUMERIC(20,2))
RETURNS NUMERIC(100,2) AS $$
DECLARE
    current_balance NUMERIC(100,2);
BEGIN
    -- Lock the row to prevent race conditions
    SELECT balance INTO current_balance
    FROM "BK_Account"
    WHERE id = account_id
    FOR UPDATE;

    -- Check if balance is sufficient
    IF current_balance < amount THEN
        RAISE EXCEPTION 'Insufficient balance. Available: %, Requested: %', current_balance, amount
            USING ERRCODE = 'P0001';
    END IF;

    -- Update balance
    UPDATE "BK_Account"
    SET balance = balance - amount, updated_at = NOW()
    WHERE id = account_id
    RETURNING balance INTO current_balance;

    RETURN current_balance;
END;
$$ LANGUAGE plpgsql;

-- Deposit to account
CREATE OR REPLACE FUNCTION deposit_to_account(account_id BIGINT, amount NUMERIC(20, 2))
RETURNS NUMERIC(100, 2) AS $$
DECLARE
    current_balance NUMERIC(100, 2);
BEGIN
    -- Lock the row to prevent race conditions
    SELECT balance INTO current_balance
    FROM "BK_Account"
    WHERE id = account_id
    FOR UPDATE;

    -- Update the account balance
    UPDATE "BK_Account"
    SET balance = balance + amount, updated_at = NOW()
    WHERE id = account_id
    RETURNING balance INTO current_balance;

    -- Return the new balance
    RETURN current_balance;
END;
$$ LANGUAGE plpgsql;
