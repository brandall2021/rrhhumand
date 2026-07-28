-- FASE 15: Reclutamiento y Selección (ATS)

-- Solicitud de vacantes
CREATE TABLE job_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    position_id UUID,
    department_id UUID,
    requested_by UUID NOT NULL,
    hiring_manager_id UUID,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    vacancies INTEGER NOT NULL DEFAULT 1,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    location VARCHAR(200),
    salary_min NUMERIC(14,2),
    salary_max NUMERIC(14,2),
    currency VARCHAR(10),
    reason TEXT,
    status VARCHAR(30) DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_requisitions_company ON job_requisitions(company_id);
CREATE INDEX idx_requisitions_status ON job_requisitions(status);
CREATE INDEX idx_requisitions_position ON job_requisitions(position_id);
CREATE INDEX idx_requisitions_department ON job_requisitions(department_id);

-- Workflow de aprobación (configurable por empresa)
CREATE TABLE approval_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL DEFAULT 'REQUISITION',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_approval_wf_company ON approval_workflows(company_id);

-- Pasos del workflow
CREATE TABLE approval_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES approval_workflows(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    approver_role VARCHAR(100),
    approver_id UUID,
    required BOOLEAN DEFAULT TRUE
);
CREATE INDEX idx_approval_steps_workflow ON approval_steps(workflow_id);

-- Instancias de aprobación (registros de aprobación reales)
CREATE TABLE approval_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    workflow_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    current_step INT DEFAULT 1,
    status VARCHAR(30) DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_approval_inst_company ON approval_instances(company_id);
CREATE INDEX idx_approval_inst_entity ON approval_instances(entity_type, entity_id);

-- Publicaciones de trabajo
CREATE TABLE job_postings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    requisition_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT,
    responsibilities TEXT,
    benefits TEXT,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    location VARCHAR(200),
    published_at TIMESTAMPTZ,
    closing_at TIMESTAMPTZ,
    status VARCHAR(30) DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_postings_company ON job_postings(company_id);
CREATE INDEX idx_postings_requisition ON job_postings(requisition_id);
CREATE INDEX idx_postings_status ON job_postings(status);

-- Candidatos
CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(200) NOT NULL,
    phone VARCHAR(50),
    document_number VARCHAR(100),
    location VARCHAR(200),
    linkedin_url TEXT,
    portfolio_url TEXT,
    source VARCHAR(100),
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, email)
);
CREATE INDEX idx_candidates_company ON candidates(company_id);
CREATE INDEX idx_candidates_email ON candidates(company_id, email);
CREATE INDEX idx_candidates_status ON candidates(status);

-- Postulaciones
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    candidate_id UUID NOT NULL,
    job_posting_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'NEW',
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    hired_at TIMESTAMPTZ
);
CREATE INDEX idx_applications_company ON applications(company_id);
CREATE INDEX idx_applications_candidate ON applications(candidate_id);
CREATE INDEX idx_applications_posting ON applications(job_posting_id);
CREATE INDEX idx_applications_status ON applications(status);

-- Historial de etapas del candidato
CREATE TABLE candidate_stage_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL,
    from_stage VARCHAR(50),
    to_stage VARCHAR(50) NOT NULL,
    changed_by UUID,
    notes TEXT,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stage_history_application ON candidate_stage_history(application_id);

-- Documentos/CV de candidatos
CREATE TABLE candidate_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    candidate_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL DEFAULT 'CV',
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    storage_provider VARCHAR(30),
    storage_key TEXT,
    parsed_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_candidate_docs_company ON candidate_documents(company_id);
CREATE INDEX idx_candidate_docs_candidate ON candidate_documents(candidate_id);

-- Preguntas de screening
CREATE TABLE screening_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    job_posting_id UUID NOT NULL,
    question TEXT NOT NULL,
    question_type VARCHAR(30) NOT NULL DEFAULT 'BOOLEAN',
    required BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_screening_questions_posting ON screening_questions(job_posting_id);

-- Respuestas de screening
CREATE TABLE screening_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL,
    question_id UUID NOT NULL REFERENCES screening_questions(id) ON DELETE CASCADE,
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_screening_answers_application ON screening_answers(application_id);
CREATE INDEX idx_screening_answers_question ON screening_answers(question_id);

-- Entrevistas
CREATE TABLE interviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL,
    interviewer_id UUID NOT NULL,
    interview_type VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMPTZ,
    duration_minutes INT,
    meeting_url TEXT,
    location TEXT,
    status VARCHAR(30) DEFAULT 'SCHEDULED',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_interviews_company ON interviews(company_id);
CREATE INDEX idx_interviews_application ON interviews(application_id);
CREATE INDEX idx_interviews_interviewer ON interviews(interviewer_id);
CREATE INDEX idx_interviews_status ON interviews(status);

-- Feedback de entrevistas
CREATE TABLE interview_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    interview_id UUID NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    interviewer_id UUID NOT NULL,
    score NUMERIC(5,2),
    comments TEXT,
    recommendation VARCHAR(30),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_interview_feedback_interview ON interview_feedback(interview_id);

-- Evaluaciones técnicas / assessments
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL,
    assessment_type VARCHAR(50) NOT NULL DEFAULT 'TECHNICAL',
    title VARCHAR(200) NOT NULL,
    description TEXT,
    max_score NUMERIC(8,2),
    score NUMERIC(8,2),
    duration_minutes INT,
    status VARCHAR(30) DEFAULT 'PENDING',
    result TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_assessments_company ON assessments(company_id);
CREATE INDEX idx_assessments_application ON assessments(application_id);

-- Ofertas laborales
CREATE TABLE job_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL,
    position_title VARCHAR(200) NOT NULL,
    department_id UUID,
    start_date DATE,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    salary_amount NUMERIC(14,2),
    salary_currency VARCHAR(10),
    salary_period VARCHAR(30),
    benefits TEXT,
    conditions TEXT,
    response_deadline DATE,
    status VARCHAR(30) DEFAULT 'DRAFT',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_offers_company ON job_offers(company_id);
CREATE INDEX idx_offers_application ON job_offers(application_id);
CREATE INDEX idx_offers_status ON job_offers(status);

-- Referidos de empleados
CREATE TABLE employee_referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    referrer_employee_id UUID NOT NULL,
    candidate_id UUID NOT NULL,
    application_id UUID,
    status VARCHAR(30) DEFAULT 'PENDING',
    reward_status VARCHAR(30) DEFAULT 'NONE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_referrals_company ON employee_referrals(company_id);
CREATE INDEX idx_referrals_referrer ON employee_referrals(referrer_employee_id);
CREATE INDEX idx_referrals_candidate ON employee_referrals(candidate_id);

-- Log de auditoría de reclutamiento
CREATE TABLE recruitment_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID,
    candidate_id UUID,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_rec_audit_company ON recruitment_audit_log(company_id);
CREATE INDEX idx_rec_audit_entity ON recruitment_audit_log(entity_type, entity_id);

-- Permisos de reclutamiento
INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'recruitment.read', 'recruitment', 'read', 'View recruitment data', NOW()),
    (gen_random_uuid(), 'recruitment.create_requisition', 'recruitment', 'create_requisition', 'Create job requisitions', NOW()),
    (gen_random_uuid(), 'recruitment.approve_requisition', 'recruitment', 'approve_requisition', 'Approve job requisitions', NOW()),
    (gen_random_uuid(), 'recruitment.manage_postings', 'recruitment', 'manage_postings', 'Create and manage job postings', NOW()),
    (gen_random_uuid(), 'recruitment.manage_candidates', 'recruitment', 'manage_candidates', 'Create and manage candidates', NOW()),
    (gen_random_uuid(), 'recruitment.conduct_interviews', 'recruitment', 'conduct_interviews', 'Schedule and conduct interviews', NOW()),
    (gen_random_uuid(), 'recruitment.create_offers', 'recruitment', 'create_offers', 'Create and send job offers', NOW()),
    (gen_random_uuid(), 'recruitment.hire', 'recruitment', 'hire', 'Hire candidates', NOW()),
    (gen_random_uuid(), 'recruitment.analytics', 'recruitment', 'analytics', 'View recruitment analytics and dashboard', NOW())
ON CONFLICT (name) DO NOTHING;
