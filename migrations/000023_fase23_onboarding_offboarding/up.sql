-- FASE 23: Onboarding + Offboarding completo
-- Extiende la migración 000015 con offboarding, workflow engine, outbox, etc.

-- ============================================================
-- 1. EXTENDER ONBOARDING_PROCESSES (campos que faltaban)
-- ============================================================
ALTER TABLE onboarding_processes
  ADD COLUMN IF NOT EXISTS candidate_id UUID,
  ADD COLUMN IF NOT EXISTS application_id UUID,
  ADD COLUMN IF NOT EXISTS job_offer_id UUID,
  ADD COLUMN IF NOT EXISTS expected_completion_date DATE,
  ADD COLUMN IF NOT EXISTS actual_completion_date DATE,
  ADD COLUMN IF NOT EXISTS employee_type VARCHAR(50),
  ADD COLUMN IF NOT EXISTS work_mode VARCHAR(30),
  ADD COLUMN IF NOT EXISTS probation_start_date DATE,
  ADD COLUMN IF NOT EXISTS probation_end_date DATE,
  ADD COLUMN IF NOT EXISTS probation_status VARCHAR(30) DEFAULT 'PENDING';

-- ============================================================
-- 2. ONBOARDING CHECKLISTS
-- ============================================================
CREATE TABLE IF NOT EXISTS onboarding_checklists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    section VARCHAR(100) NOT NULL,
    title VARCHAR(250) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    completed_by UUID,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_oc_onboarding ON onboarding_checklists(onboarding_id);
CREATE INDEX IF NOT EXISTS idx_oc_company ON onboarding_checklists(company_id);
CREATE INDEX IF NOT EXISTS idx_oc_status ON onboarding_checklists(company_id, status);

-- ============================================================
-- 3. ONBOARDING TASK DEPENDENCIES
-- ============================================================
CREATE TABLE IF NOT EXISTS onboarding_task_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES onboarding_tasks(id) ON DELETE CASCADE,
    depends_on_task_id UUID NOT NULL REFERENCES onboarding_tasks(id) ON DELETE CASCADE,
    UNIQUE(task_id, depends_on_task_id)
);
CREATE INDEX IF NOT EXISTS idx_otd_task ON onboarding_task_dependencies(task_id);
CREATE INDEX IF NOT EXISTS idx_otd_depends ON onboarding_task_dependencies(depends_on_task_id);

-- ============================================================
-- 4. DOCUMENT VERSIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS onboarding_document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    document_id UUID NOT NULL REFERENCES onboarding_documents(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    file_name VARCHAR(255),
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    storage_key TEXT,
    checksum VARCHAR(64),
    uploaded_by UUID,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);
CREATE INDEX IF NOT EXISTS idx_odv_document ON onboarding_document_versions(document_id);
CREATE INDEX IF NOT EXISTS idx_odv_company ON onboarding_document_versions(company_id);

-- ============================================================
-- 5. ONBOARDING NOTES
-- ============================================================
CREATE TABLE IF NOT EXISTS onboarding_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    content TEXT NOT NULL,
    visibility VARCHAR(30) NOT NULL DEFAULT 'HR_ONLY',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_onotes_onboarding ON onboarding_notes(onboarding_id);
CREATE INDEX IF NOT EXISTS idx_onotes_company ON onboarding_notes(company_id);

-- ============================================================
-- 6. OFFBOARDING — PLANTILLAS
-- ============================================================
CREATE TABLE IF NOT EXISTS offboarding_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_offt_company ON offboarding_templates(company_id);

-- Tareas de plantilla
CREATE TABLE IF NOT EXISTS offboarding_template_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES offboarding_templates(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    task_type VARCHAR(50) NOT NULL DEFAULT 'OTHER',
    assigned_role VARCHAR(50),
    required BOOLEAN NOT NULL DEFAULT TRUE,
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_offtt_template ON offboarding_template_tasks(template_id);

-- ============================================================
-- 7. OFFBOARDING — PROCESOS
-- ============================================================
CREATE TABLE IF NOT EXISTS offboarding_processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    template_id UUID REFERENCES offboarding_templates(id),
    requested_by UUID NOT NULL,
    termination_type VARCHAR(50) NOT NULL,
    reason_id UUID,
    notice_date DATE NOT NULL,
    last_working_date DATE NOT NULL,
    termination_effective_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    progress NUMERIC(5,2) NOT NULL DEFAULT 0,
    employee_status_after VARCHAR(30) DEFAULT 'INACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_offp_company ON offboarding_processes(company_id);
CREATE INDEX IF NOT EXISTS idx_offp_employee ON offboarding_processes(employee_id);
CREATE INDEX IF NOT EXISTS idx_offp_status ON offboarding_processes(company_id, status);
CREATE INDEX IF NOT EXISTS idx_offp_last_working ON offboarding_processes(last_working_date);

-- ============================================================
-- 8. OFFBOARDING — TAREAS
-- ============================================================
CREATE TABLE IF NOT EXISTS offboarding_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    offboarding_id UUID NOT NULL REFERENCES offboarding_processes(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    task_type VARCHAR(50) NOT NULL DEFAULT 'OTHER',
    assigned_to UUID,
    assigned_role VARCHAR(50),
    required BOOLEAN NOT NULL DEFAULT TRUE,
    due_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    completed_by UUID,
    comments TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_offt_offboarding ON offboarding_tasks(offboarding_id);
CREATE INDEX IF NOT EXISTS idx_offt_company ON offboarding_tasks(company_id);
CREATE INDEX IF NOT EXISTS idx_offt_assigned ON offboarding_tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_offt_status ON offboarding_tasks(company_id, status);

-- ============================================================
-- 9. OFFBOARDING — ACTIVOS A RECUPERAR
-- ============================================================
CREATE TABLE IF NOT EXISTS offboarding_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    offboarding_id UUID NOT NULL REFERENCES offboarding_processes(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    description TEXT,
    serial_number VARCHAR(100),
    condition_on_delivery TEXT,
    condition_on_return TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_RETURN',
    returned_at TIMESTAMPTZ,
    returned_to UUID,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_offa_offboarding ON offboarding_assets(offboarding_id);
CREATE INDEX IF NOT EXISTS idx_offa_company ON offboarding_assets(company_id);
CREATE INDEX IF NOT EXISTS idx_offa_employee ON offboarding_assets(employee_id);

-- ============================================================
-- 10. OFFBOARDING — REVOCACIÓN DE ACCESOS
-- ============================================================
CREATE TABLE IF NOT EXISTS employee_access_revocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    offboarding_id UUID REFERENCES offboarding_processes(id) ON DELETE CASCADE,
    system_name VARCHAR(150) NOT NULL,
    access_type VARCHAR(100) NOT NULL DEFAULT 'STANDARD',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    performed_by UUID,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ear_company ON employee_access_revocations(company_id);
CREATE INDEX IF NOT EXISTS idx_ear_employee ON employee_access_revocations(employee_id);
CREATE INDEX IF NOT EXISTS idx_ear_offboarding ON employee_access_revocations(offboarding_id);
CREATE INDEX IF NOT EXISTS idx_ear_status ON employee_access_revocations(status);

-- ============================================================
-- 11. OFFBOARDING — TRANSFERENCIA DE INFORMACIÓN
-- ============================================================
CREATE TABLE IF NOT EXISTS employee_handovers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    offboarding_id UUID REFERENCES offboarding_processes(id) ON DELETE CASCADE,
    handover_to UUID NOT NULL,
    description TEXT,
    projects JSONB,
    pending_tasks JSONB,
    documents JSONB,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_eh_company ON employee_handovers(company_id);
CREATE INDEX IF NOT EXISTS idx_eh_employee ON employee_handovers(employee_id);
CREATE INDEX IF NOT EXISTS idx_eh_handover_to ON employee_handovers(handover_to);

-- ============================================================
-- 12. OFFBOARDING — ENTREVISTA DE SALIDA
-- ============================================================
CREATE TABLE IF NOT EXISTS exit_interviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    offboarding_id UUID NOT NULL REFERENCES offboarding_processes(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    interviewer_id UUID,
    scheduled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    reason TEXT,
    feedback TEXT,
    recommendation VARCHAR(50),
    rating NUMERIC(5,2),
    anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ei_company ON exit_interviews(company_id);
CREATE INDEX IF NOT EXISTS idx_ei_offboarding ON exit_interviews(offboarding_id);

-- Preguntas configurables
CREATE TABLE IF NOT EXISTS exit_interview_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    question TEXT NOT NULL,
    question_type VARCHAR(30) NOT NULL DEFAULT 'TEXT',
    sort_order INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_eiq_company ON exit_interview_questions(company_id);

-- Respuestas
CREATE TABLE IF NOT EXISTS exit_interview_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exit_interview_id UUID NOT NULL REFERENCES exit_interviews(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES exit_interview_questions(id) ON DELETE CASCADE,
    answer TEXT,
    rating INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_eia_interview ON exit_interview_answers(exit_interview_id);

-- ============================================================
-- 13. MOTIVOS DE SALIDA (catálogo configurable)
-- ============================================================
CREATE TABLE IF NOT EXISTS employee_exit_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_eer_company ON employee_exit_reasons(company_id);

-- ============================================================
-- 14. WORKFLOW RULES (motor de reglas)
-- ============================================================
CREATE TABLE IF NOT EXISTS workflow_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    workflow_type VARCHAR(30) NOT NULL,
    name VARCHAR(200) NOT NULL,
    conditions JSONB NOT NULL,
    actions JSONB NOT NULL,
    priority INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_wr_company ON workflow_rules(company_id);
CREATE INDEX IF NOT EXISTS idx_wr_type ON workflow_rules(company_id, workflow_type);

-- ============================================================
-- 15. OUTBOX PATTERN
-- ============================================================
CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_oe_status ON outbox_events(status);
CREATE INDEX IF NOT EXISTS idx_oe_type ON outbox_events(event_type);
CREATE INDEX IF NOT EXISTS idx_oe_company ON outbox_events(company_id);
CREATE INDEX IF NOT EXISTS idx_oe_aggregate ON outbox_events(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_oe_created ON outbox_events(created_at);

-- ============================================================
-- 16. OFFBOARDING AUDIT LOG
-- ============================================================
CREATE TABLE IF NOT EXISTS offboarding_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_offal_company ON offboarding_audit_log(company_id);
CREATE INDEX IF NOT EXISTS idx_offal_action ON offboarding_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_offal_entity ON offboarding_audit_log(entity_type, entity_id);

-- ============================================================
-- 17. NUEVOS PERMISOS (offboarding + onboarding extendido)
-- ============================================================
INSERT INTO permissions (name, resource, action, description, created_at) VALUES
    ('offboarding.read', 'offboarding', 'read', 'View offboarding processes', NOW()),
    ('offboarding.create', 'offboarding', 'create', 'Create offboarding processes', NOW()),
    ('offboarding.update', 'offboarding', 'update', 'Update offboarding processes', NOW()),
    ('offboarding.approve', 'offboarding', 'approve', 'Approve offboarding requests', NOW()),
    ('offboarding.start', 'offboarding', 'start', 'Start offboarding processes', NOW()),
    ('offboarding.complete', 'offboarding', 'complete', 'Complete offboarding processes', NOW()),
    ('offboarding.cancel', 'offboarding', 'cancel', 'Cancel offboarding processes', NOW()),
    ('offboarding.tasks.read', 'offboarding', 'tasks.read', 'View offboarding tasks', NOW()),
    ('offboarding.tasks.create', 'offboarding', 'tasks.create', 'Create offboarding tasks', NOW()),
    ('offboarding.tasks.complete', 'offboarding', 'tasks.complete', 'Complete offboarding tasks', NOW()),
    ('offboarding.assets.read', 'offboarding', 'assets.read', 'View offboarding assets', NOW()),
    ('offboarding.assets.return', 'offboarding', 'assets.return', 'Return offboarding assets', NOW()),
    ('offboarding.access.read', 'offboarding', 'access.read', 'View access revocations', NOW()),
    ('offboarding.access.revoke', 'offboarding', 'access.revoke', 'Revoke access', NOW()),
    ('offboarding.templates.read', 'offboarding', 'templates.read', 'View offboarding templates', NOW()),
    ('offboarding.templates.create', 'offboarding', 'templates.create', 'Create offboarding templates', NOW()),
    ('offboarding.handover.read', 'offboarding', 'handover.read', 'View employee handovers', NOW()),
    ('offboarding.handover.create', 'offboarding', 'handover.create', 'Create employee handovers', NOW()),
    ('exit_interview.read', 'exit_interview', 'read', 'View exit interviews', NOW()),
    ('exit_interview.create', 'exit_interview', 'create', 'Create exit interviews', NOW()),
    ('exit_interview.update', 'exit_interview', 'update', 'Update exit interviews', NOW()),
    ('exit_interview.analyze', 'exit_interview', 'analyze', 'Analyze exit interview data', NOW()),
    ('onboarding.admin', 'onboarding', 'admin', 'Full admin access to onboarding', NOW()),
    ('offboarding.admin', 'offboarding', 'admin', 'Full admin access to offboarding', NOW()),
    ('onboarding.checklists.read', 'onboarding', 'checklists.read', 'View onboarding checklists', NOW()),
    ('onboarding.checklists.complete', 'onboarding', 'checklists.complete', 'Complete onboarding checklist items', NOW()),
    ('onboarding.documents.delete', 'onboarding', 'documents.delete', 'Delete onboarding documents', NOW()),
    ('onboarding.assets.return', 'onboarding', 'assets.return', 'Return onboarding assets', NOW()),
    ('onboarding.notes.read', 'onboarding', 'notes.read', 'View onboarding notes', NOW()),
    ('onboarding.notes.create', 'onboarding', 'notes.create', 'Create onboarding notes', NOW()),
    ('onboarding.probation.read', 'onboarding', 'probation.read', 'View probation status', NOW()),
    ('onboarding.probation.update', 'onboarding', 'probation.update', 'Update probation status', NOW())
ON CONFLICT (name) DO NOTHING;
