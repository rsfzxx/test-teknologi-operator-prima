-- Stored procedure: review rekening bank (Accepted / Rejected)
-- Aturan: Rejected wajib punya reason, rekening yang sudah dihapus tidak bisa direview.
CREATE OR REPLACE PROCEDURE sp_review_bank_account(
    p_bank_account_uuid UUID,
    p_status            VARCHAR,
    p_reason            TEXT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_status NOT IN ('Accepted', 'Rejected') THEN
        RAISE EXCEPTION 'Status harus Accepted atau Rejected, diterima: %', p_status;
    END IF;

    IF p_status = 'Rejected' AND (p_reason IS NULL OR btrim(p_reason) = '') THEN
        RAISE EXCEPTION 'Reason wajib diisi jika status Rejected';
    END IF;

    UPDATE bank_accounts
       SET status     = p_status,
           reason     = CASE WHEN p_status = 'Rejected' THEN p_reason ELSE NULL END,
           updated_at = NOW()
     WHERE bank_account_uuid = p_bank_account_uuid
       AND deleted_at IS NULL;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Bank account % tidak ditemukan', p_bank_account_uuid;
    END IF;
END;
$$;