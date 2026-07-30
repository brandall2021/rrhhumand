-- 000025: Down - revert schema fixes

DROP TABLE IF EXISTS compensation_budgets CASCADE;

ALTER TABLE employee_benefits DROP CONSTRAINT IF EXISTS chk_employee_benefit_status;
ALTER TABLE employee_benefits ADD CONSTRAINT chk_employee_benefit_status CHECK (status IN ('active','suspended','expired','cancelled'));

ALTER TABLE employee_benefits DROP COLUMN IF EXISTS enrollment_date;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS employee_cost;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS company_cost;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS created_by;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS updated_at;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS plan_id;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS provider_id;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS start_date;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS end_date;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS cancellation_date;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS cancellation_reason;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS auto_renew;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS renewal_date;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS coverage_level;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS dependents;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS emergency_contact;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS beneficiary_info;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS employer_cost;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS payroll_deduction_enabled;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS payroll_deduction_amount;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS external_member_id;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS external_policy_number;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS external_group_number;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS documents;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS notes;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS source;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS enrolled_by;
ALTER TABLE employee_benefits DROP COLUMN IF EXISTS enrolled_at;
