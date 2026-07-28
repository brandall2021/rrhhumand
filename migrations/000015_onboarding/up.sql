-- FASE 16: Onboarding / Incorporación de Empleados

-- Plantillas de onboarding
CREATE TABLE onboarding_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    default_duration_days INT NOT NULL DEFAULT 90,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ot_company ON onboarding_templates(company_id);
CREATE INDEX idx_ot_status ON onboarding_templates(status);

-- Tareas de plantilla
CREATE TABLE onboarding_template_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES onboarding_templates(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL DEFAULT 'PERSONAL',
    responsible_type VARCHAR(30) NOT NULL DEFAULT 'EMPLOYEE',
    responsible_user_id UUID,
    required BOOLEAN NOT NULL DEFAULT TRUE,
    days_offset INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    estimated_minutes INT
);
CREATE INDEX idx_ott_template ON onboarding_template_tasks(template_id);

-- Procesos de onboarding
CREATE TABLE onboarding_processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    template_id UUID,
    start_date DATE NOT NULL,
    target_completion_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'NOT_STARTED',
    progress_percentage INT NOT NULL DEFAULT 0,
    completion_policy VARCHAR(30) NOT NULL DEFAULT 'STRICT',
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_op_company ON onboarding_processes(company_id);
CREATE INDEX idx_op_employee ON onboarding_processes(employee_id);
CREATE INDEX idx_op_status ON onboarding_processes(company_id, status);
CREATE INDEX idx_op_target_date ON onboarding_processes(target_completion_date);

-- Tareas del proceso de onboarding
CREATE TABLE onboarding_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL DEFAULT 'PERSONAL',
    responsible_type VARCHAR(30) NOT NULL DEFAULT 'EMPLOYEE',
    responsible_id UUID,
    due_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    required BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    estimated_minutes INT,
    completed_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    blocked_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ot_onboarding ON onboarding_tasks(onboarding_id);
CREATE INDEX idx_ot_company ON onboarding_tasks(company_id);
CREATE INDEX idx_ot_employee ON onboarding_tasks(employee_id);
CREATE INDEX idx_ot_responsible ON onboarding_tasks(responsible_id);
CREATE INDEX idx_ot_due_date ON onboarding_tasks(company_id, due_date);
CREATE INDEX idx_ot_status ON onboarding_tasks(status);

-- Documentos del onboarding
CREATE TABLE onboarding_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255),
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    checksum VARCHAR(64),
    storage_provider VARCHAR(30),
    storage_key TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'REQUIRED',
    required BOOLEAN NOT NULL DEFAULT TRUE,
    rejection_reason TEXT,
    uploaded_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_od_onboarding ON onboarding_documents(onboarding_id);
CREATE INDEX idx_od_company ON onboarding_documents(company_id);
CREATE INDEX idx_od_status ON onboarding_documents(status);

-- Equipamiento asignado
CREATE TABLE onboarding_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    description TEXT,
    serial_number VARCHAR(100),
    status VARCHAR(30) NOT NULL DEFAULT 'REQUESTED',
    assigned_by UUID,
    assigned_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    returned_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_oa_onboarding ON onboarding_assets(onboarding_id);
CREATE INDEX idx_oa_company ON onboarding_assets(company_id);

-- Solicitudes de acceso a sistemas
CREATE TABLE access_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    system_name VARCHAR(100) NOT NULL,
    access_type VARCHAR(50) NOT NULL DEFAULT 'STANDARD',
    status VARCHAR(30) NOT NULL DEFAULT 'REQUESTED',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    approved_by UUID,
    activated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ar_onboarding ON access_requests(onboarding_id);
CREATE INDEX idx_ar_company ON access_requests(company_id);
CREATE INDEX idx_ar_employee ON access_requests(employee_id);
CREATE INDEX idx_ar_status ON access_requests(status);

-- Milestones del onboarding
CREATE TABLE onboarding_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    milestone_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    days_offset INT NOT NULL DEFAULT 30,
    due_date DATE NOT NULL,
    responsible_type VARCHAR(30) NOT NULL DEFAULT 'MANAGER',
    responsible_id UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_om_onboarding ON onboarding_milestones(onboarding_id);
CREATE INDEX idx_om_company ON onboarding_milestones(company_id);

-- Feedback del onboarding
CREATE TABLE onboarding_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    feedback_type VARCHAR(30) NOT NULL,
    submitted_by UUID NOT NULL,
    adaptation_score INT,
    team_score INT,
    knowledge_score INT,
    communication_score INT,
    overall_score NUMERIC(4,1),
    comments TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_of_onboarding ON onboarding_feedback(onboarding_id);
CREATE INDEX idx_of_company ON onboarding_feedback(company_id);

-- Buddy/Mentor asignado
CREATE TABLE onboarding_buddies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    buddy_employee_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ob_onboarding ON onboarding_buddies(onboarding_id);
CREATE INDEX idx_ob_company ON onboarding_buddies(company_id);
CREATE INDEX idx_ob_buddy ON onboarding_buddies(buddy_employee_id);

-- Excepciones del onboarding
CREATE TABLE onboarding_exceptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    reason TEXT NOT NULL,
    created_by UUID NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_oe_onboarding ON onboarding_exceptions(onboarding_id);
CREATE INDEX idx_oe_company ON onboarding_exceptions(company_id);

-- Capacitaciones asignadas
CREATE TABLE training_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES onboarding_processes(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    course_name VARCHAR(200) NOT NULL,
    description TEXT,
    training_type VARCHAR(50) NOT NULL DEFAULT 'MANDATORY',
    status VARCHAR(30) NOT NULL DEFAULT 'ASSIGNED',
    due_date DATE,
    completed_at TIMESTAMPTZ,
    external_provider VARCHAR(100),
    external_course_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ta_onboarding ON training_assignments(onboarding_id);
CREATE INDEX idx_ta_company ON training_assignments(company_id);
CREATE INDEX idx_ta_employee ON training_assignments(employee_id);

-- Eventos de dominio
CREATE TABLE domain_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    company_id UUID NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);
CREATE INDEX idx_de_type ON domain_events(event_type);
CREATE INDEX idx_de_company ON domain_events(company_id);
CREATE INDEX idx_de_aggregate ON domain_events(aggregate_id);
CREATE INDEX idx_de_processed ON domain_events(processed_at);

-- Notificaciones internas
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    notification_type VARCHAR(50) NOT NULL,
    channel VARCHAR(30) NOT NULL DEFAULT 'IN_APP',
    reference_type VARCHAR(50),
    reference_id UUID,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notif_company ON notifications(company_id);
CREATE INDEX idx_notif_user ON notifications(user_id);
CREATE INDEX idx_notif_read ON notifications(company_id, user_id, read_at);

-- Log de auditoría de onboarding
CREATE TABLE onboarding_audit_log (
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
CREATE INDEX idx_oal_company ON onboarding_audit_log(company_id);
CREATE INDEX idx_oal_action ON onboarding_audit_log(action);
CREATE INDEX idx_oal_entity ON onboarding_audit_log(entity_type, entity_id);

-- Permisos de onboarding
INSERT INTO permissions (name, resource, action, description, created_at) VALUES
    ('onboarding.read', 'onboarding', 'read', 'View onboarding processes', NOW()),
    ('onboarding.create', 'onboarding', 'create', 'Create onboarding processes', NOW()),
    ('onboarding.update', 'onboarding', 'update', 'Update onboarding processes', NOW()),
    ('onboarding.start', 'onboarding', 'start', 'Start onboarding processes', NOW()),
    ('onboarding.complete', 'onboarding', 'complete', 'Complete onboarding processes', NOW()),
    ('onboarding.cancel', 'onboarding', 'cancel', 'Cancel onboarding processes', NOW()),
    ('onboarding.tasks.read', 'onboarding', 'tasks.read', 'View onboarding tasks', NOW()),
    ('onboarding.tasks.create', 'onboarding', 'tasks.create', 'Create onboarding tasks', NOW()),
    ('onboarding.tasks.update', 'onboarding', 'tasks.update', 'Update onboarding tasks', NOW()),
    ('onboarding.tasks.complete', 'onboarding', 'tasks.complete', 'Complete onboarding tasks', NOW()),
    ('onboarding.documents.read', 'onboarding', 'documents.read', 'View onboarding documents', NOW()),
    ('onboarding.documents.upload', 'onboarding', 'documents.upload', 'Upload onboarding documents', NOW()),
    ('onboarding.documents.review', 'onboarding', 'documents.review', 'Review onboarding documents', NOW()),
    ('onboarding.documents.approve', 'onboarding', 'documents.approve', 'Approve onboarding documents', NOW()),
    ('onboarding.assets.read', 'onboarding', 'assets.read', 'View onboarding assets', NOW()),
    ('onboarding.assets.assign', 'onboarding', 'assets.assign', 'Assign onboarding assets', NOW()),
    ('onboarding.assets.deliver', 'onboarding', 'assets.deliver', 'Deliver onboarding assets', NOW()),
    ('onboarding.access.read', 'onboarding', 'access.read', 'View access requests', NOW()),
    ('onboarding.access.request', 'onboarding', 'access.request', 'Request access', NOW()),
    ('onboarding.access.approve', 'onboarding', 'access.approve', 'Approve access requests', NOW()),
    ('onboarding.templates.read', 'onboarding', 'templates.read', 'View onboarding templates', NOW()),
    ('onboarding.templates.create', 'onboarding', 'templates.create', 'Create onboarding templates', NOW()),
    ('onboarding.templates.update', 'onboarding', 'templates.update', 'Update onboarding templates', NOW()),
    ('onboarding.templates.delete', 'onboarding', 'templates.delete', 'Delete onboarding templates', NOW()),
    ('onboarding.analytics.read', 'onboarding', 'analytics.read', 'View onboarding analytics', NOW())
ON CONFLICT (name) DO NOTHING;
