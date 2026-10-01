CREATE TABLE IF NOT EXISTS banks (
    bank_uuid  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(150) NOT NULL,
    code       VARCHAR(20)  NOT NULL,
    type       VARCHAR(50)  NOT NULL,
    status     VARCHAR(10)  NOT NULL DEFAULT 'Active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_banks_status CHECK (status IN ('Active', 'Inactive')),
    CONSTRAINT uq_banks_code UNIQUE (code)
);