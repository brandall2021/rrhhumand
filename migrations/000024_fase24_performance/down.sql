-- Revert FASE 24: restore tables from 000013

-- Drop new tables
DROP TABLE IF EXISTS performance_outbox CASCADE;
DROP TABLE IF EXISTS performance_audit_log CASCADE;
DROP TABLE IF EXISTS performance_results CASCADE;
DROP TABLE IF EXISTS performance_evidence CASCADE;
DROP TABLE IF EXISTS development_plan_actions CASCADE;
DROP TABLE IF EXISTS performance_development_plans CASCADE;
DROP TABLE IF EXISTS improvement_plan_actions CASCADE;
DROP TABLE IF EXISTS performance_improvement_plans CASCADE;
DROP TABLE IF EXISTS calibration_items CASCADE;
DROP TABLE IF EXISTS calibration_sessions CASCADE;
DROP TABLE IF EXISTS performance_reviews CASCADE;
DROP TABLE IF EXISTS performance_recognitions CASCADE;
DROP TABLE IF EXISTS performance_checkins CASCADE;
DROP TABLE IF EXISTS performance_feedback CASCADE;
DROP TABLE IF EXISTS competency_evaluations CASCADE;
DROP TABLE IF EXISTS objective_evaluations CASCADE;
DROP TABLE IF EXISTS evaluation_answers CASCADE;
DROP TABLE IF EXISTS performance_evaluations CASCADE;
DROP TABLE IF EXISTS performance_participants CASCADE;
DROP TABLE IF EXISTS objective_key_results CASCADE;
DROP TABLE IF EXISTS performance_objectives CASCADE;
DROP TABLE IF EXISTS cycle_competencies CASCADE;
DROP TABLE IF EXISTS position_competencies CASCADE;
DROP TABLE IF EXISTS competency_levels CASCADE;
DROP TABLE IF EXISTS competencies CASCADE;
DROP TABLE IF EXISTS rating_scale_levels CASCADE;
DROP TABLE IF EXISTS rating_scales CASCADE;
DROP TABLE IF EXISTS template_questions CASCADE;
DROP TABLE IF EXISTS template_sections CASCADE;
DROP TABLE IF EXISTS performance_templates CASCADE;
DROP TABLE IF EXISTS performance_cycles CASCADE;

-- Remove permissions
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE name LIKE 'performance.%'
);
DELETE FROM permissions WHERE name LIKE 'performance.%';

-- Note: to restore 000013 tables, re-run migrations/000013_performance/up.sql
