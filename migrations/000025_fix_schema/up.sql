-- 000025: Fix schema - add missing columns and tables
-- This migration is idempotent

-- ============================================================
-- 1. Add missing columns to employee_benefits
-- ============================================================
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS employee_cost NUMERIC(15,2) NOT NULL DEFAULT 0;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS company_cost NUMERIC(15,2) NOT NULL DEFAULT 0;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS created_by UUID;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS plan_id UUID;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS provider_id UUID;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS end_date DATE;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS cancellation_date DATE;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS cancellation_reason TEXT;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS renewal_date DATE;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS coverage_level VARCHAR(30) DEFAULT 'INDIVIDUAL';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS dependents JSONB DEFAULT '[]';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS emergency_contact JSONB DEFAULT '{}';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS beneficiary_info JSONB DEFAULT '{}';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS employer_cost NUMERIC(15,2) NOT NULL DEFAULT 0;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS payroll_deduction_enabled BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS payroll_deduction_amount NUMERIC(15,2) DEFAULT 0;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS external_member_id VARCHAR(100);
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS external_policy_number VARCHAR(100);
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS external_group_number VARCHAR(100);
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS documents JSONB DEFAULT '[]';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS source VARCHAR(20) NOT NULL DEFAULT 'ADMIN';
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS enrolled_by UUID;
ALTER TABLE employee_benefits ADD COLUMN IF NOT EXISTS enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Drop old constraint from 000012 and add proper constraint
ALTER TABLE employee_benefits DROP CONSTRAINT IF EXISTS chk_employee_benefit_status;
ALTER TABLE employee_benefits ADD CONSTRAINT chk_employee_benefit_status CHECK (status IN ('active','suspended','expired','cancelled','ACTIVE','PENDING','SUSPENDED','CANCELLED','EXPIRED','WAITING_PERIOD','DECLINED'));

-- ============================================================
-- 2. Create compensation_budgets if not exists
-- ============================================================
CREATE TABLE IF NOT EXISTS compensation_budgets (
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_comp_budgets_unique ON compensation_budgets(company_id, year, department_id);
CREATE INDEX IF NOT EXISTS idx_comp_budgets_company ON compensation_budgets(company_id);
CREATE INDEX IF NOT EXISTS idx_comp_budgets_year ON compensation_budgets(year);
CREATE INDEX IF NOT EXISTS idx_comp_budgets_status ON compensation_budgets(status);
