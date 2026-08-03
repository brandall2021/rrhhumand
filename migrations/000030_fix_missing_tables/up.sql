-- ============================================================
-- 000030 Fix missing tables
-- ============================================================
-- compensation_domain_events is defined in 000017_compensation but that
-- migration wraps everything in a transaction, so a single conflicting
-- CREATE TABLE ("employee_compensations" already exists) aborts the whole
-- file and rolls it back. Recreate the table here idempotently.

CREATE TABLE IF NOT EXISTS compensation_domain_events (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    payload JSONB,
    created_by UUID,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cde_company ON compensation_domain_events(company_id, processed_at);
