-- ============================================================
-- 000031 Fix training tables
-- ============================================================
-- The internal/training Go package references training_course_versions
-- and training_enrollments, but no migration ever created them.
-- 000026 only patched training_course_offerings / training_events.
-- Schemas derived from internal/training/repository.go + models.go.

CREATE TABLE IF NOT EXISTS training_course_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES training_courses(id) ON DELETE CASCADE,
    version VARCHAR(50) NOT NULL DEFAULT 'v1',
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    is_published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tcv_course ON training_course_versions(course_id);

CREATE TABLE IF NOT EXISTS training_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    offering_id UUID NOT NULL REFERENCES training_course_offerings(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    assignment_type VARCHAR(30) NOT NULL DEFAULT 'OPTIONAL',
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    enrollment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    dropped_at TIMESTAMPTZ,
    final_score NUMERIC(5,2),
    passed BOOLEAN,
    progress_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    certificate_url TEXT,
    certificate_issued_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tenr_company ON training_enrollments(company_id);
CREATE INDEX IF NOT EXISTS idx_tenr_offering ON training_enrollments(offering_id);
CREATE INDEX IF NOT EXISTS idx_tenr_employee ON training_enrollments(employee_id);
