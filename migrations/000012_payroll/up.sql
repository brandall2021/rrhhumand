CREATE TABLE payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN',
    calculated_at TIMESTAMPTZ,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    closed_by UUID,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, start_date, end_date)
);
CREATE INDEX idx_payroll_periods_company ON payroll_periods(company_id);
CREATE INDEX idx_payroll_periods_status ON payroll_periods(status);

CREATE TABLE payroll_concepts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    type VARCHAR(30) NOT NULL,
    calculation_type VARCHAR(30) NOT NULL DEFAULT 'FIXED',
    taxable BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);
CREATE INDEX idx_payroll_concepts_company ON payroll_concepts(company_id);

CREATE TABLE employee_compensations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    base_amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    effective_from DATE NOT NULL,
    effective_to DATE,
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_emp_comp_company ON employee_compensations(company_id);
CREATE INDEX idx_emp_comp_employee ON employee_compensations(employee_id);

CREATE TABLE payroll_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_period_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    concept_id UUID NOT NULL,
    quantity NUMERIC(14,4) DEFAULT 1,
    unit_amount NUMERIC(14,2) DEFAULT 0,
    amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_items_period ON payroll_items(payroll_period_id);
CREATE INDEX idx_payroll_items_employee ON payroll_items(employee_id);

CREATE TABLE benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    benefit_type VARCHAR(50),
    default_amount NUMERIC(14,2) DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);
CREATE INDEX idx_benefits_company ON benefits(company_id);

CREATE TABLE employee_benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    benefit_id UUID NOT NULL,
    amount NUMERIC(14,2),
    currency VARCHAR(10) DEFAULT 'USD',
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_emp_benefits_company ON employee_benefits(company_id);
CREATE INDEX idx_emp_benefits_employee ON employee_benefits(employee_id);

CREATE TABLE payroll_bonuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    bonus_type VARCHAR(50) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    reason TEXT,
    period_start DATE,
    period_end DATE,
    status VARCHAR(30) DEFAULT 'PENDING',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_bonuses_company ON payroll_bonuses(company_id);
CREATE INDEX idx_payroll_bonuses_employee ON payroll_bonuses(employee_id);

CREATE TABLE payroll_advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    request_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(30) DEFAULT 'PENDING',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_advances_company ON payroll_advances(company_id);
CREATE INDEX idx_payroll_advances_employee ON payroll_advances(employee_id);

CREATE TABLE payroll_deductions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    concept VARCHAR(100) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    reason TEXT,
    period_start DATE,
    period_end DATE,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_deductions_company ON payroll_deductions(company_id);
CREATE INDEX idx_payroll_deductions_employee ON payroll_deductions(employee_id);

CREATE TABLE payroll_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    payroll_period_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    reason TEXT NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'CREDIT',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_adjustments_period ON payroll_adjustments(payroll_period_id);

CREATE TABLE payroll_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    payroll_period_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    concept_code VARCHAR(50),
    amount NUMERIC(14,2) NOT NULL,
    balance_after NUMERIC(14,2) DEFAULT 0,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_ledger_period ON payroll_ledger(payroll_period_id);
CREATE INDEX idx_payroll_ledger_employee ON payroll_ledger(employee_id);

CREATE TABLE payroll_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    payroll_period_id UUID NOT NULL,
    snapshot_data JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_payroll_snapshots_period ON payroll_snapshots(payroll_period_id);

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'payroll.read', 'payroll', 'read', 'View payroll data', NOW()),
    (gen_random_uuid(), 'payroll.create_period', 'payroll', 'create_period', 'Create payroll periods', NOW()),
    (gen_random_uuid(), 'payroll.update_period', 'payroll', 'update_period', 'Update payroll periods', NOW()),
    (gen_random_uuid(), 'payroll.calculate', 'payroll', 'calculate', 'Calculate payroll', NOW()),
    (gen_random_uuid(), 'payroll.approve', 'payroll', 'approve', 'Approve payroll', NOW()),
    (gen_random_uuid(), 'payroll.close', 'payroll', 'close', 'Close payroll period', NOW()),
    (gen_random_uuid(), 'payroll.manage_concepts', 'payroll', 'manage_concepts', 'Manage payroll concepts', NOW()),
    (gen_random_uuid(), 'payroll.manage_compensation', 'payroll', 'manage_compensation', 'Manage employee compensation', NOW()),
    (gen_random_uuid(), 'payroll.manage_benefits', 'payroll', 'manage_benefits', 'Manage benefits', NOW()),
    (gen_random_uuid(), 'payroll.manage_bonuses', 'payroll', 'manage_bonuses', 'Manage bonuses', NOW()),
    (gen_random_uuid(), 'payroll.manage_deductions', 'payroll', 'manage_deductions', 'Manage deductions', NOW()),
    (gen_random_uuid(), 'payroll.manage_advances', 'payroll', 'manage_advances', 'Manage advances', NOW()),
    (gen_random_uuid(), 'payroll.adjust', 'payroll', 'adjust', 'Make payroll adjustments', NOW()),
    (gen_random_uuid(), 'payroll.view_snapshot', 'payroll', 'view_snapshot', 'View payroll snapshots', NOW()),
    (gen_random_uuid(), 'payroll.view_reports', 'payroll', 'view_reports', 'View payroll reports', NOW())
ON CONFLICT (name) DO NOTHING;
