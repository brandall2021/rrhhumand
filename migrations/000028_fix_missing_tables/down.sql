-- 000028: Revert - Fix missing tables, columns, and seed data

-- 1. Remove seed employees
DELETE FROM employees WHERE id IN (
    'a0000000-0000-0000-0000-000000000060',
    'a0000000-0000-0000-0000-000000000061',
    'a0000000-0000-0000-0000-000000000062',
    'a0000000-0000-0000-0000-000000000063',
    'a0000000-0000-0000-0000-000000000064'
);

-- 2. Drop compensation tables
DROP TABLE IF EXISTS compensation_audit_logs CASCADE;
DROP TABLE IF EXISTS compensation_equity_snapshots CASCADE;
DROP TABLE IF EXISTS compensation_reviews CASCADE;
DROP TABLE IF EXISTS bonuses CASCADE;
DROP TABLE IF EXISTS bonus_plans CASCADE;
DROP TABLE IF EXISTS salary_adjustment_proposals CASCADE;
DROP TABLE IF EXISTS compensation_adjustments CASCADE;
DROP TABLE IF EXISTS compensation_history CASCADE;
DROP TABLE IF EXISTS employee_compensation_components CASCADE;
DROP TABLE IF EXISTS compensation_components CASCADE;
DROP TABLE IF EXISTS position_salary_bands CASCADE;
DROP TABLE IF EXISTS salary_bands CASCADE;
DROP TABLE IF EXISTS salary_grades CASCADE;
DROP TABLE IF EXISTS compensation_structures CASCADE;

-- 3. Remove added columns from benefits
ALTER TABLE benefits DROP COLUMN IF EXISTS updated_at;
ALTER TABLE benefits DROP COLUMN IF EXISTS created_at;
ALTER TABLE benefits DROP COLUMN IF EXISTS created_by;
ALTER TABLE benefits DROP COLUMN IF EXISTS taxable;
ALTER TABLE benefits DROP COLUMN IF EXISTS frequency;
ALTER TABLE benefits DROP COLUMN IF EXISTS cost_currency;
ALTER TABLE benefits DROP COLUMN IF EXISTS cost_amount;
ALTER TABLE benefits DROP COLUMN IF EXISTS provider;

-- 4. Remove added columns from employee_compensations
ALTER TABLE employee_compensations DROP COLUMN IF EXISTS updated_at;
ALTER TABLE employee_compensations DROP COLUMN IF EXISTS created_by;
ALTER TABLE employee_compensations DROP COLUMN IF EXISTS status;
ALTER TABLE employee_compensations DROP COLUMN IF EXISTS pay_frequency;
ALTER TABLE employee_compensations DROP COLUMN IF EXISTS salary_band_id;

-- 5. Remove added columns from job_postings
ALTER TABLE job_postings DROP COLUMN IF EXISTS updated_at;
ALTER TABLE job_postings DROP COLUMN IF EXISTS external_url;
ALTER TABLE job_postings DROP COLUMN IF EXISTS is_public;
ALTER TABLE job_postings DROP COLUMN IF EXISTS currency;
ALTER TABLE job_postings DROP COLUMN IF EXISTS salary_max;
ALTER TABLE job_postings DROP COLUMN IF EXISTS salary_min;
ALTER TABLE job_postings DROP COLUMN IF EXISTS position_id;

-- 6. Rename report_id back to expense_report_id in expense_reimbursements
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='expense_reimbursements' AND column_name='report_id'
    ) THEN
        ALTER TABLE expense_reimbursements RENAME COLUMN report_id TO expense_report_id;
    END IF;
END $$;

-- 7. Rename report_id back to expense_report_id in expenses
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='expenses' AND column_name='report_id'
    ) THEN
        ALTER TABLE expenses RENAME COLUMN report_id TO expense_report_id;
    END IF;
END $$;
