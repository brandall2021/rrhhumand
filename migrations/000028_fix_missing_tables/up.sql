-- 000028: Fix missing tables, columns, and seed data
-- This migration is idempotent

-- ============================================================
-- 1. Rename expense_report_id to report_id in expenses
-- ============================================================
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='expenses' AND column_name='expense_report_id'
    ) THEN
        ALTER TABLE expenses RENAME COLUMN expense_report_id TO report_id;
    END IF;
END $$;

-- ============================================================
-- 2. Rename expense_report_id to report_id in expense_reimbursements
-- ============================================================
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='expense_reimbursements' AND column_name='expense_report_id'
    ) THEN
        ALTER TABLE expense_reimbursements RENAME COLUMN expense_report_id TO report_id;
    END IF;
END $$;

-- ============================================================
-- 3. Add missing columns to job_postings
-- ============================================================
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS position_id UUID;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS salary_min NUMERIC(15,2);
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS salary_max NUMERIC(15,2);
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS currency VARCHAR(3) NOT NULL DEFAULT 'ARS';
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS is_public BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS external_url TEXT;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_postings_position ON job_postings(position_id);

-- ============================================================
-- 4. Add missing columns to employee_compensations
-- ============================================================
ALTER TABLE employee_compensations ADD COLUMN IF NOT EXISTS salary_band_id UUID;
ALTER TABLE employee_compensations ADD COLUMN IF NOT EXISTS pay_frequency VARCHAR(20) NOT NULL DEFAULT 'monthly';
ALTER TABLE employee_compensations ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';
ALTER TABLE employee_compensations ADD COLUMN IF NOT EXISTS created_by UUID;
ALTER TABLE employee_compensations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_emp_comp_band ON employee_compensations(salary_band_id);
CREATE INDEX IF NOT EXISTS idx_emp_comp_status ON employee_compensations(status);

-- ============================================================
-- 5. Add missing columns to benefits
-- ============================================================
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS provider VARCHAR(200);
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS cost_amount NUMERIC(15,2);
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS cost_currency VARCHAR(3) NOT NULL DEFAULT 'USD';
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS frequency VARCHAR(20) NOT NULL DEFAULT 'monthly';
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS taxable BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS created_by UUID;
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- ============================================================
-- 6. Create missing compensation tables
-- ============================================================

-- 6a. compensation_structures
CREATE TABLE IF NOT EXISTS compensation_structures (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cs_company ON compensation_structures(company_id);
CREATE INDEX IF NOT EXISTS idx_cs_status ON compensation_structures(status);

-- 6b. salary_grades
CREATE TABLE IF NOT EXISTS salary_grades (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    structure_id UUID NOT NULL REFERENCES compensation_structures(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sg_company ON salary_grades(company_id);
CREATE INDEX IF NOT EXISTS idx_sg_structure ON salary_grades(structure_id);

-- 6c. salary_bands
CREATE TABLE IF NOT EXISTS salary_bands (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    structure_id UUID NOT NULL REFERENCES compensation_structures(id) ON DELETE CASCADE,
    grade_id UUID REFERENCES salary_grades(id) ON DELETE SET NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    minimum_amount NUMERIC(15,2) NOT NULL,
    midpoint_amount NUMERIC(15,2) NOT NULL,
    maximum_amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sb_company ON salary_bands(company_id);
CREATE INDEX IF NOT EXISTS idx_sb_structure ON salary_bands(structure_id);
CREATE INDEX IF NOT EXISTS idx_sb_grade ON salary_bands(grade_id);

-- 6d. position_salary_bands
CREATE TABLE IF NOT EXISTS position_salary_bands (
    id UUID PRIMARY KEY,
    position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    salary_band_id UUID NOT NULL REFERENCES salary_bands(id) ON DELETE CASCADE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_psb_position ON position_salary_bands(position_id);
CREATE INDEX IF NOT EXISTS idx_psb_band ON position_salary_bands(salary_band_id);

-- 6e. compensation_components
CREATE TABLE IF NOT EXISTS compensation_components (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    component_type VARCHAR(30) NOT NULL DEFAULT 'fixed',
    taxable BOOLEAN NOT NULL DEFAULT true,
    recurring BOOLEAN NOT NULL DEFAULT false,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cc_company ON compensation_components(company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_cc_code ON compensation_components(company_id, code);

-- 6f. employee_compensation_components
CREATE TABLE IF NOT EXISTS employee_compensation_components (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES compensation_components(id) ON DELETE CASCADE,
    amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ecc_company ON employee_compensation_components(company_id);
CREATE INDEX IF NOT EXISTS idx_ecc_employee ON employee_compensation_components(employee_id);
CREATE INDEX IF NOT EXISTS idx_ecc_component ON employee_compensation_components(component_id);

-- 6g. compensation_history
CREATE TABLE IF NOT EXISTS compensation_history (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    previous_amount NUMERIC(15,2),
    new_amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    reason VARCHAR(200) NOT NULL,
    effective_from DATE NOT NULL,
    approved_by UUID,
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ch_company ON compensation_history(company_id);
CREATE INDEX IF NOT EXISTS idx_ch_employee ON compensation_history(employee_id);

-- 6h. compensation_adjustments
CREATE TABLE IF NOT EXISTS compensation_adjustments (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    adjustment_type VARCHAR(30) NOT NULL,
    value NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    reason VARCHAR(200) NOT NULL,
    effective_from DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    applied_by UUID,
    applied_at TIMESTAMPTZ,
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ca_company ON compensation_adjustments(company_id);
CREATE INDEX IF NOT EXISTS idx_ca_employee ON compensation_adjustments(employee_id);
CREATE INDEX IF NOT EXISTS idx_ca_status ON compensation_adjustments(status);

-- 6i. salary_adjustment_proposals
CREATE TABLE IF NOT EXISTS salary_adjustment_proposals (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    review_id UUID,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    current_amount NUMERIC(15,2) NOT NULL,
    proposed_amount NUMERIC(15,2) NOT NULL,
    increase_percentage NUMERIC(5,2),
    reason TEXT NOT NULL,
    performance_score NUMERIC(5,2),
    market_position VARCHAR(30),
    manager_comment TEXT,
    hr_comment TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    submitted_by UUID,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    rejected_by UUID,
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sap_company ON salary_adjustment_proposals(company_id);
CREATE INDEX IF NOT EXISTS idx_sap_employee ON salary_adjustment_proposals(employee_id);
CREATE INDEX IF NOT EXISTS idx_sap_review ON salary_adjustment_proposals(review_id);
CREATE INDEX IF NOT EXISTS idx_sap_status ON salary_adjustment_proposals(status);

-- 6j. bonus_plans
CREATE TABLE IF NOT EXISTS bonus_plans (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    period VARCHAR(30) NOT NULL DEFAULT 'annual',
    target_percentage NUMERIC(5,2),
    maximum_percentage NUMERIC(5,2),
    eligibility_rules TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bp_company ON bonus_plans(company_id);
CREATE INDEX IF NOT EXISTS idx_bp_status ON bonus_plans(status);

-- 6k. bonuses
CREATE TABLE IF NOT EXISTS bonuses (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    bonus_plan_id UUID REFERENCES bonus_plans(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    bonus_type VARCHAR(30) NOT NULL DEFAULT 'performance',
    amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    period VARCHAR(30),
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bon_company ON bonuses(company_id);
CREATE INDEX IF NOT EXISTS idx_bon_employee ON bonuses(employee_id);
CREATE INDEX IF NOT EXISTS idx_bon_plan ON bonuses(bonus_plan_id);
CREATE INDEX IF NOT EXISTS idx_bon_status ON bonuses(status);

-- 6l. compensation_reviews
CREATE TABLE IF NOT EXISTS compensation_reviews (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    period VARCHAR(30) NOT NULL DEFAULT 'annual',
    start_date DATE NOT NULL,
    end_date DATE,
    budget NUMERIC(15,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cr_company ON compensation_reviews(company_id);
CREATE INDEX IF NOT EXISTS idx_cr_status ON compensation_reviews(status);

-- 6m. compensation_equity_snapshots
CREATE TABLE IF NOT EXISTS compensation_equity_snapshots (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    position_id UUID REFERENCES positions(id) ON DELETE SET NULL,
    grade_id UUID REFERENCES salary_grades(id) ON DELETE SET NULL,
    employee_count INT NOT NULL DEFAULT 0,
    median_compensation NUMERIC(15,2),
    average_compensation NUMERIC(15,2),
    min_compensation NUMERIC(15,2),
    max_compensation NUMERIC(15,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    metadata TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ces_company ON compensation_equity_snapshots(company_id);
CREATE INDEX IF NOT EXISTS idx_ces_date ON compensation_equity_snapshots(company_id, snapshot_date);

-- 6n. compensation_audit_logs
CREATE TABLE IF NOT EXISTS compensation_audit_logs (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_value TEXT,
    new_value TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cal_company ON compensation_audit_logs(company_id);
CREATE INDEX IF NOT EXISTS idx_cal_entity ON compensation_audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_cal_user ON compensation_audit_logs(user_id);

-- ============================================================
-- 7. Insert seed positions (with integer level values)
--    (The original 000027 seed used strings like 'senior' for 
--     an integer column, causing silent 0-row inserts)
-- ============================================================
INSERT INTO positions (id, company_id, department_id, name, code, level, active)
VALUES
    ('a0000000-0000-0000-0000-000000000030', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000020', 'Desarrollador Senior', 'SR-DEV', 4, true),
    ('a0000000-0000-0000-0000-000000000031', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000020', 'Desarrollador Junior', 'JR-DEV', 1, true),
    ('a0000000-0000-0000-0000-000000000032', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000021', 'Analista de RRHH', 'HR-ANL', 3, true),
    ('a0000000-0000-0000-0000-000000000033', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000022', 'Ejecutivo de Ventas', 'SALES-EXEC', 3, true),
    ('a0000000-0000-0000-0000-000000000034', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000024', 'Analista Financiero', 'FIN-ANL', 3, true)
ON CONFLICT (company_id, code) DO NOTHING;

-- ============================================================
-- 8. Re-insert seed employee data
-- ============================================================
DELETE FROM employees WHERE company_id = 'a0000000-0000-0000-0000-000000000001';

-- Employee 1: Manager (Carlos García)
INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000060',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-001',
    'Carlos',
    'García',
    'carlos.garcia@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000020',
    'a0000000-0000-0000-0000-000000000030',
    '2020-03-15',
    'active'
);

-- Employee 2: Developer (María López)
INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000061',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-002',
    'María',
    'López',
    'maria.lopez@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000020',
    'a0000000-0000-0000-0000-000000000031',
    'a0000000-0000-0000-0000-000000000060',
    '2022-06-01',
    'active'
);

-- Employee 3: HR (Ana Martínez)
INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000062',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-003',
    'Ana',
    'Martínez',
    'ana.martinez@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000021',
    'a0000000-0000-0000-0000-000000000032',
    'a0000000-0000-0000-0000-000000000060',
    '2023-01-10',
    'active'
);

-- Employee 4: Sales (Pedro Fernández)
INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000063',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-004',
    'Pedro',
    'Fernández',
    'pedro.fernandez@techcorp.com',
    'a0000000-0000-0000-0000-000000000011',
    'a0000000-0000-0000-0000-000000000022',
    'a0000000-0000-0000-0000-000000000033',
    'a0000000-0000-0000-0000-000000000060',
    '2023-03-20',
    'active'
);

-- Employee 5: Finance (Lucía Torres)
INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000064',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-005',
    'Lucía',
    'Torres',
    'lucia.torres@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000024',
    'a0000000-0000-0000-0000-000000000034',
    'a0000000-0000-0000-0000-000000000060',
    '2023-08-01',
    'active'
);
