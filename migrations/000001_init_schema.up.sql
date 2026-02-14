BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
    account_id BIGINT PRIMARY KEY,
    balance    NUMERIC(20, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT accounts_balance_non_negative CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    id                     BIGSERIAL   PRIMARY KEY,
    source_account_id      BIGINT      NOT NULL REFERENCES accounts(account_id),
    destination_account_id BIGINT      NOT NULL REFERENCES accounts(account_id),
    amount                 NUMERIC(20, 2) NOT NULL,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT transactions_amount_positive CHECK (amount > 0),
    CONSTRAINT transactions_different_accounts CHECK (source_account_id <> destination_account_id)
);

CREATE INDEX idx_transactions_source ON transactions(source_account_id);
CREATE INDEX idx_transactions_destination ON transactions(destination_account_id);

COMMIT;
