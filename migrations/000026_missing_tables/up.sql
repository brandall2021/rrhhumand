-- 000026: Create missing worker tables (training_events, training_course_offerings, scoring_matches)

-- ============================================================
-- 1. training_course_offerings
-- ============================================================
CREATE TABLE IF NOT EXISTS training_course_offerings (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    course_id UUID NOT NULL,
    course_version_id UUID,
    name VARCHAR(255) NOT NULL,
    start_date DATE,
    end_date DATE,
    enrollment_start DATE,
    enrollment_end DATE,
    capacity INT NOT NULL DEFAULT 0,
    enrolled_count INT NOT NULL DEFAULT 0,
    modality VARCHAR(50),
    location VARCHAR(255),
    meeting_url TEXT,
    instructor_id UUID,
    provider_id UUID,
    cost_amount NUMERIC(14,2),
    cost_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tco_company ON training_course_offerings(company_id);
CREATE INDEX IF NOT EXISTS idx_tco_course ON training_course_offerings(course_id);
CREATE INDEX IF NOT EXISTS idx_tco_status ON training_course_offerings(status);
CREATE INDEX IF NOT EXISTS idx_tco_instructor ON training_course_offerings(instructor_id);

-- ============================================================
-- 2. training_events
-- ============================================================
CREATE TABLE IF NOT EXISTS training_events (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    employee_id UUID,
    enrollment_id UUID,
    offering_id UUID,
    severity VARCHAR(20) NOT NULL DEFAULT 'info',
    scheduled_for TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    metadata JSONB,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_te_company ON training_events(company_id);
CREATE INDEX IF NOT EXISTS idx_te_processed ON training_events(company_id, processed_at);
CREATE INDEX IF NOT EXISTS idx_te_scheduled ON training_events(company_id, scheduled_for);
CREATE INDEX IF NOT EXISTS idx_te_severity ON training_events(severity);

-- ============================================================
-- 3. scoring_matches
-- ============================================================
CREATE TABLE IF NOT EXISTS scoring_matches (
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    job_posting_id UUID,
    candidate_id UUID,
    score NUMERIC(5,2),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sm_status ON scoring_matches(status);
CREATE INDEX IF NOT EXISTS idx_sm_posting ON scoring_matches(job_posting_id);
CREATE INDEX IF NOT EXISTS idx_sm_candidate ON scoring_matches(candidate_id);
