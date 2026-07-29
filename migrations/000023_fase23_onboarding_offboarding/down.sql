-- FASE 23 Down
DROP TABLE IF EXISTS offboarding_audit_log CASCADE;
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS workflow_rules CASCADE;
DROP TABLE IF EXISTS exit_interview_answers CASCADE;
DROP TABLE IF EXISTS exit_interview_questions CASCADE;
DROP TABLE IF EXISTS exit_interviews CASCADE;
DROP TABLE IF EXISTS employee_handovers CASCADE;
DROP TABLE IF EXISTS employee_access_revocations CASCADE;
DROP TABLE IF EXISTS offboarding_assets CASCADE;
DROP TABLE IF EXISTS offboarding_tasks CASCADE;
DROP TABLE IF EXISTS offboarding_processes CASCADE;
DROP TABLE IF EXISTS offboarding_template_tasks CASCADE;
DROP TABLE IF EXISTS offboarding_templates CASCADE;
DROP TABLE IF EXISTS employee_exit_reasons CASCADE;
DROP TABLE IF EXISTS onboarding_notes CASCADE;
DROP TABLE IF EXISTS onboarding_document_versions CASCADE;
DROP TABLE IF EXISTS onboarding_task_dependencies CASCADE;
DROP TABLE IF EXISTS onboarding_checklists CASCADE;

ALTER TABLE onboarding_processes
  DROP COLUMN IF EXISTS candidate_id,
  DROP COLUMN IF EXISTS application_id,
  DROP COLUMN IF EXISTS job_offer_id,
  DROP COLUMN IF EXISTS expected_completion_date,
  DROP COLUMN IF EXISTS actual_completion_date,
  DROP COLUMN IF EXISTS employee_type,
  DROP COLUMN IF EXISTS work_mode,
  DROP COLUMN IF EXISTS probation_start_date,
  DROP COLUMN IF EXISTS probation_end_date,
  DROP COLUMN IF EXISTS probation_status;

DELETE FROM permissions WHERE resource IN ('offboarding', 'exit_interview');
DELETE FROM permissions WHERE name IN (
    'onboarding.admin', 'onboarding.checklists.read', 'onboarding.checklists.complete',
    'onboarding.documents.delete', 'onboarding.assets.return',
    'onboarding.notes.read', 'onboarding.notes.create',
    'onboarding.probation.read', 'onboarding.probation.update'
);
