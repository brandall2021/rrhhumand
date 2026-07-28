DROP TABLE IF EXISTS onboarding_audit_log CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS domain_events CASCADE;
DROP TABLE IF EXISTS training_assignments CASCADE;
DROP TABLE IF EXISTS onboarding_exceptions CASCADE;
DROP TABLE IF EXISTS onboarding_buddies CASCADE;
DROP TABLE IF EXISTS onboarding_feedback CASCADE;
DROP TABLE IF EXISTS onboarding_milestones CASCADE;
DROP TABLE IF EXISTS access_requests CASCADE;
DROP TABLE IF EXISTS onboarding_assets CASCADE;
DROP TABLE IF EXISTS onboarding_documents CASCADE;
DROP TABLE IF EXISTS onboarding_tasks CASCADE;
DROP TABLE IF EXISTS onboarding_processes CASCADE;
DROP TABLE IF EXISTS onboarding_template_tasks CASCADE;
DROP TABLE IF EXISTS onboarding_templates CASCADE;

DELETE FROM permissions WHERE resource = 'onboarding';
