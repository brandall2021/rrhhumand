-- FASE 18: Compensation & Benefits
-- Migration: 000017

BEGIN;

-- 1. compensation_structures
CREATE TABLE compensation_structures (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. salary_grades
CREATE TABLE salary_grades (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    structure_id UUID NOT NULL REFERENCES compensation_structures(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

-- 3. salary_bands
CREATE TABLE salary_bands (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    structure_id UUID NOT NULL REFERENCES compensation_structures(id) ON DELETE CASCADE,
    grade_id UUID REFERENCES salary_grades(id) ON DELETE SET NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    minimum_amount NUMERIC(15,2) NOT NULL,
    midpoint_amount NUMERIC(15,2) NOT NULL,
    maximum_amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code),
    CONSTRAINT chk_band_range CHECK (minimum_amount <= midpoint_amount AND midpoint_amount <= maximum_amount AND minimum_amount < maximum_amount)
);

-- 4. position_salary_bands
CREATE TABLE position_salary_bands (
    id UUID PRIMARY KEY,
    position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    salary_band_id UUID NOT NULL REFERENCES salary_bands(id) ON DELETE CASCADE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(position_id, salary_band_id, effective_from)
);

-- 5. employee_compensations
CREATE TABLE employee_compensations (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    salary_band_id UUID REFERENCES salary_bands(id) ON DELETE SET NULL,
    base_amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    pay_frequency VARCHAR(20) NOT NULL DEFAULT 'monthly',
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_base_amount CHECK (base_amount > 0),
    CONSTRAINT chk_pay_frequency CHECK (pay_frequency IN ('hourly','daily','weekly','biweekly','monthly','annual'))
);

-- 6. compensation_components (catalog)
CREATE TABLE compensation_components (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    component_type VARCHAR(30) NOT NULL DEFAULT 'salary',
    taxable BOOLEAN NOT NULL DEFAULT true,
    recurring BOOLEAN NOT NULL DEFAULT false,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code),
    CONSTRAINT chk_component_type CHECK (component_type IN ('salary','bonus','commission','allowance','overtime','premium','incentive','benefit','other'))
);

-- 7. employee_compensation_components
CREATE TABLE employee_compensation_components (
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_comp_comp_amount CHECK (amount > 0)
);

-- 8. compensation_history
CREATE TABLE compensation_history (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    previous_amount NUMERIC(15,2),
    new_amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    reason VARCHAR(50) NOT NULL,
    effective_from DATE NOT NULL,
    approved_by UUID,
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_reason CHECK (reason IN ('annual_review','promotion','role_change','market_adjustment','merit_increase','collective_adjustment','correction','other'))
);

-- 9. compensation_adjustments
CREATE TABLE compensation_adjustments (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    adjustment_type VARCHAR(20) NOT NULL,
    value NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    reason VARCHAR(50) NOT NULL,
    effective_from DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    applied_by UUID,
    applied_at TIMESTAMPTZ,
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_adjustment_type CHECK (adjustment_type IN ('percentage','fixed_amount','new_salary')),
    CONSTRAINT chk_adjustment_status CHECK (status IN ('draft','pending_approval','approved','rejected','applied','cancelled'))
);

-- 10. salary_adjustment_proposals
CREATE TABLE salary_adjustment_proposals (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    review_id UUID REFERENCES compensation_reviews(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    current_amount NUMERIC(15,2) NOT NULL,
    proposed_amount NUMERIC(15,2) NOT NULL,
    increase_percentage NUMERIC(5,2),
    reason VARCHAR(50) NOT NULL,
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_proposal_status CHECK (status IN ('draft','submitted','approved','rejected','applied','cancelled'))
);

-- 11. bonus_plans
CREATE TABLE bonus_plans (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    period VARCHAR(20) NOT NULL DEFAULT 'annual',
    target_percentage NUMERIC(5,2),
    maximum_percentage NUMERIC(5,2),
    eligibility_rules JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_bonus_plan_status CHECK (status IN ('draft','active','inactive','archived'))
);

-- 12. bonuses
CREATE TABLE bonuses (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    bonus_plan_id UUID REFERENCES bonus_plans(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    bonus_type VARCHAR(30) NOT NULL DEFAULT 'discretionary',
    amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    period VARCHAR(20),
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_bonus_amount CHECK (amount > 0),
    CONSTRAINT chk_bonus_type CHECK (bonus_type IN ('discretionary','performance','commission','signing','retention','referral','project','other')),
    CONSTRAINT chk_bonus_status CHECK (status IN ('draft','pending_approval','approved','rejected','paid','cancelled'))
);

-- 13. benefits (catalog)
CREATE TABLE benefits (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    benefit_type VARCHAR(30) NOT NULL DEFAULT 'other',
    provider VARCHAR(255),
    cost_amount NUMERIC(15,2),
    cost_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    frequency VARCHAR(20) NOT NULL DEFAULT 'monthly',
    taxable BOOLEAN NOT NULL DEFAULT false,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code),
    CONSTRAINT chk_benefit_type CHECK (benefit_type IN ('health_insurance','life_insurance','meal','transportation','gym','education','internet','phone','vehicle','vouchers','extra_days','other'))
);

-- 14. employee_benefits
CREATE TABLE employee_benefits (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    benefit_id UUID NOT NULL REFERENCES benefits(id) ON DELETE CASCADE,
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    employee_cost NUMERIC(15,2) NOT NULL DEFAULT 0,
    company_cost NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_employee_benefit_status CHECK (status IN ('active','suspended','expired','cancelled'))
);

-- 15. compensation_reviews
CREATE TABLE compensation_reviews (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    period VARCHAR(20) NOT NULL DEFAULT 'annual',
    start_date DATE NOT NULL,
    end_date DATE,
    budget NUMERIC(15,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_review_status CHECK (status IN ('draft','open','in_review','approved','closed','cancelled'))
);

-- 16. compensation_budgets
CREATE TABLE compensation_budgets (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    year INT NOT NULL,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    budget_amount NUMERIC(15,2) NOT NULL,
    committed_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    spent_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, year, department_id),
    CONSTRAINT chk_budget_status CHECK (status IN ('draft','active','closed','archived'))
);

-- 17. compensation_equity_snapshots
CREATE TABLE compensation_equity_snapshots (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    position_id UUID REFERENCES positions(id) ON DELETE SET NULL,
    grade_id UUID REFERENCES salary_grades(id) ON DELETE SET NULL,
    employee_count INT NOT NULL,
    median_compensation NUMERIC(15,2),
    average_compensation NUMERIC(15,2),
    min_compensation NUMERIC(15,2),
    max_compensation NUMERIC(15,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    metadata JSONB,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 18. audit_logs (compensation-specific)
CREATE TABLE compensation_audit_logs (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 19. domain_events
CREATE TABLE compensation_domain_events (
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

-- Indexes
CREATE INDEX idx_ec_employee ON employee_compensations(company_id, employee_id);
CREATE INDEX idx_ch_employee ON compensation_history(company_id, employee_id);
CREATE INDEX idx_sb_company ON salary_bands(company_id);
CREATE INDEX idx_sb_grade ON salary_bands(company_id, grade_id);
CREATE INDEX idx_eb_employee ON employee_benefits(company_id, employee_id);
CREATE INDEX idx_bonus_employee ON bonuses(company_id, employee_id);
CREATE INDEX idx_cr_company ON compensation_reviews(company_id);
CREATE INDEX idx_cb_year ON compensation_budgets(company_id, year);
CREATE INDEX idx_cb_dept_year ON compensation_budgets(company_id, department_id, year);
CREATE INDEX idx_psb_position ON position_salary_bands(position_id);
CREATE INDEX idx_ecc_employee ON employee_compensation_components(company_id, employee_id);
CREATE INDEX idx_sap_review ON salary_adjustment_proposals(review_id);
CREATE INDEX idx_sap_employee ON salary_adjustment_proposals(company_id, employee_id);
CREATE INDEX idx_ces_snapshot ON compensation_equity_snapshots(company_id, snapshot_date);
CREATE INDEX idx_cal_company ON compensation_audit_logs(company_id, created_at);
CREATE INDEX idx_cde_company ON compensation_domain_events(company_id, processed_at);

-- RBAC Permissions
INSERT INTO rbac_permissions (id, code, name, module, created_at) VALUES
    (gen_random_uuid(), 'compensation.read', 'View compensations', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation.create', 'Create compensations', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation.update', 'Update compensations', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation.delete', 'Delete compensations', 'compensation', NOW()),
    (gen_random_uuid(), 'salary_band.read', 'View salary bands', 'compensation', NOW()),
    (gen_random_uuid(), 'salary_band.manage', 'Manage salary bands', 'compensation', NOW()),
    (gen_random_uuid(), 'salary_adjustment.create', 'Create salary adjustments', 'compensation', NOW()),
    (gen_random_uuid(), 'salary_adjustment.approve', 'Approve salary adjustments', 'compensation', NOW()),
    (gen_random_uuid(), 'salary_adjustment.apply', 'Apply salary adjustments', 'compensation', NOW()),
    (gen_random_uuid(), 'bonus.read', 'View bonuses', 'compensation', NOW()),
    (gen_random_uuid(), 'bonus.create', 'Create bonuses', 'compensation', NOW()),
    (gen_random_uuid(), 'bonus.approve', 'Approve bonuses', 'compensation', NOW()),
    (gen_random_uuid(), 'benefit.read', 'View benefits', 'compensation', NOW()),
    (gen_random_uuid(), 'benefit.manage', 'Manage benefits', 'compensation', NOW()),
    (gen_random_uuid(), 'benefit.assign', 'Assign benefits', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation_review.create', 'Create reviews', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation_review.manage', 'Manage reviews', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation_budget.read', 'View budgets', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation_budget.manage', 'Manage budgets', 'compensation', NOW()),
    (gen_random_uuid(), 'compensation_reports.read', 'View reports', 'compensation', NOW())
ON CONFLICT (code) DO NOTHING;

COMMIT;
