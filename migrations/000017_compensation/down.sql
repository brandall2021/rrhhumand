BEGIN;

DROP TABLE IF EXISTS compensation_domain_events CASCADE;
DROP TABLE IF EXISTS compensation_audit_logs CASCADE;
DROP TABLE IF EXISTS compensation_equity_snapshots CASCADE;
DROP TABLE IF EXISTS compensation_budgets CASCADE;
DROP TABLE IF EXISTS compensation_reviews CASCADE;
DROP TABLE IF EXISTS employee_benefits CASCADE;
DROP TABLE IF EXISTS benefits CASCADE;
DROP TABLE IF EXISTS bonuses CASCADE;
DROP TABLE IF EXISTS bonus_plans CASCADE;
DROP TABLE IF EXISTS salary_adjustment_proposals CASCADE;
DROP TABLE IF EXISTS compensation_adjustments CASCADE;
DROP TABLE IF EXISTS compensation_history CASCADE;
DROP TABLE IF EXISTS employee_compensation_components CASCADE;
DROP TABLE IF EXISTS compensation_components CASCADE;
DROP TABLE IF EXISTS employee_compensations CASCADE;
DROP TABLE IF EXISTS position_salary_bands CASCADE;
DROP TABLE IF EXISTS salary_bands CASCADE;
DROP TABLE IF EXISTS salary_grades CASCADE;
DROP TABLE IF EXISTS compensation_structures CASCADE;

COMMIT;
