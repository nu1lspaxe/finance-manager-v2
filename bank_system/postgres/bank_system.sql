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

-- Withdraw from account
CREATE OR REPLACE FUNCTION withdraw_from_account(account_id BIGINT, amount NUMERIC(20,2), tx_detail TEXT)
RETURNS TABLE (
    new_balance NUMERIC(100, 2),
    transaction_id BIGINT
) AS $$
DECLARE
    tx_id BIGINT;
BEGIN
    IF amount <= 0 THEN
        RAISE EXCEPTION 'Amount must be positive, got %', amount USING ERRCODE = 'P0001';
    END IF;

    UPDATE "BK_Account"
    SET balance = balance - amount,
        updated_at = NOW()
    WHERE id = account_id AND balance >= amount
    RETURNING balance INTO new_balance;

    IF NOT FOUND THEN
        PERFORM 1 FROM "BK_Account" WHERE id = account_id;
        IF FOUND THEN
            RAISE EXCEPTION 'Insufficient balance for account %', account_id USING ERRCODE = 'P0001';
        ELSE
            RAISE EXCEPTION 'Account % not found', account_id USING ERRCODE = 'P0001';
        END IF;
    END IF;

    INSERT INTO "BK_Transaction" (account_id, amount, balance_after, tx_type, detail)
    VALUES (account_id, amount, new_balance, 'WITHDRAW', tx_detail)
    RETURNING id INTO tx_id;

    RETURN QUERY SELECT new_balance, tx_id;
END;
$$ LANGUAGE plpgsql;

-- Deposit to account
CREATE OR REPLACE FUNCTION deposit_to_account(account_id BIGINT, amount NUMERIC(20, 2), tx_detail TEXT)
RETURNS TABLE (
    new_balance NUMERIC(100, 2),
    transaction_id BIGINT
) AS $$
DECLARE
    new_balance NUMERIC(100, 2);
    tx_id BIGINT;
BEGIN
    IF amount <= 0 THEN
        RAISE EXCEPTION 'Amount must be positive, got %', amount USING ERRCODE = 'P0001';
    END IF;

    UPDATE "BK_Account"
    SET balance = balance + amount,
        updated_at = NOW()
    WHERE id = account_id
    RETURNING balance INTO new_balance;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Account % not found', account_id USING ERRCODE = 'P0001';
    END IF;

    INSERT INTO "BK_Transaction" (account_id, amount, balance_after, tx_type, detail)
    VALUES (account_id, amount, new_balance, 'DEPOSIT', tx_detail)
    RETURNING id INTO tx_id;

    RETURN QUERY SELECT new_balance, tx_id;
END;
$$ LANGUAGE plpgsql;