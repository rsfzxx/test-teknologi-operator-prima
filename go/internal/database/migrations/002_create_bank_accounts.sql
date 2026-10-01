CREATE TABLE IF NOT EXISTS bank_accounts (
    bank_account_uuid UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    account_number    VARCHAR(30)  NOT NULL,
    bank_name         VARCHAR(150) NOT NULL,
    account_name      VARCHAR(150) NOT NULL,
    bank_code         VARCHAR(20)  NOT NULL,
    status            VARCHAR(10)  NOT NULL DEFAULT 'Review',
    reason            TEXT,
    created_at        TIMESTAMPTZ  DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT chk_ba_status         CHECK (status IN ('Accepted', 'Review', 'Rejected')),
    CONSTRAINT chk_ba_account_number CHECK (account_number ~ '^[0-9]+$'),
    CONSTRAINT chk_ba_bank_code      CHECK (bank_code ~ '^[0-9]+$'),
    CONSTRAINT chk_ba_bank_name      CHECK (bank_name ~ '^[A-Za-z .,&''-]+$'),
    CONSTRAINT chk_ba_account_name   CHECK (account_name ~ '^[A-Za-z .,&''-]+$')
);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_status  ON bank_accounts (status)  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bank_accounts_number  ON bank_accounts (account_number) WHERE deleted_at IS NULL;