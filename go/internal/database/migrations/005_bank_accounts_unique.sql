ALTER TABLE bank_accounts
    ADD CONSTRAINT uq_bank_accounts_number_bank UNIQUE (account_number, bank_name);