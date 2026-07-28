-- FASE 22: Reclutamiento y Seleccion (ATS completo)
-- Reemplaza la FASE 14 con un ATS de capa completa

-- ============================================================
-- CATALOGOS BASE
-- ============================================================

CREATE TABLE recruitment_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'MANUAL',
    config JSONB,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rec_sources_company ON recruitment_sources(company_id);

CREATE TABLE recruitment_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'SCREENING',
    sort_order INT NOT NULL DEFAULT 0,
    color VARCHAR(20),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rec_stages_company ON recruitment_stages(company_id);

CREATE TABLE recruitment_stage_transitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    from_stage_id UUID NOT NULL REFERENCES recruitment_stages(id),
    to_stage_id UUID NOT NULL REFERENCES recruitment_stages(id),
    required_actions TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, from_stage_id, to_stage_id)
);
CREATE INDEX idx_rec_transitions_from ON recruitment_stage_transitions(from_stage_id);
CREATE INDEX idx_rec_transitions_to ON recruitment_stage_transitions(to_stage_id);

CREATE TABLE rejection_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rejection_reasons_company ON rejection_reasons(company_id);

-- ============================================================
-- REQUISICIONES
-- ============================================================

CREATE TABLE job_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    position_id UUID,
    department_id UUID,
    requested_by UUID NOT NULL,
    hiring_manager_id UUID,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    justification TEXT,
    vacancies INT NOT NULL DEFAULT 1,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    location VARCHAR(200),
    salary_min NUMERIC(14,2),
    salary_max NUMERIC(14,2),
    currency VARCHAR(10),
    urgency VARCHAR(30) DEFAULT 'NORMAL',
    reason TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    approved_at TIMESTAMPTZ,
    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    closed_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_jr_company ON job_requisitions(company_id);
CREATE INDEX idx_jr_status ON job_requisitions(status);
CREATE INDEX idx_jr_position ON job_requisitions(position_id);
CREATE INDEX idx_jr_department ON job_requisitions(department_id);
CREATE INDEX idx_jr_requested_by ON job_requisitions(requested_by);

CREATE TABLE job_requisition_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id UUID NOT NULL REFERENCES job_requisitions(id) ON DELETE CASCADE,
    skill VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    required BOOLEAN DEFAULT TRUE,
    min_years INT
);
CREATE INDEX idx_jrs_requisition ON job_requisition_skills(requisition_id);

CREATE TABLE job_requisition_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id UUID NOT NULL REFERENCES job_requisitions(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    step_order INT NOT NULL DEFAULT 1,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    comment TEXT,
    decided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_jra_requisition ON job_requisition_approvals(requisition_id);
CREATE INDEX idx_jra_approver ON job_requisition_approvals(approver_id);

-- ============================================================
-- PUESTOS / VACANTES
-- ============================================================

CREATE TABLE job_positions_ats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    requisition_id UUID REFERENCES job_requisitions(id),
    title VARCHAR(200) NOT NULL,
    department_id UUID,
    location_id UUID,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    description TEXT,
    requirements TEXT,
    responsibilities TEXT,
    benefits TEXT,
    salary_min NUMERIC(14,2),
    salary_max NUMERIC(14,2),
    currency VARCHAR(10),
    vacancies INT NOT NULL DEFAULT 1,
    vacancies_filled INT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_jp_company ON job_positions_ats(company_id);
CREATE INDEX idx_jp_status ON job_positions_ats(status);
CREATE INDEX idx_jp_requisition ON job_positions_ats(requisition_id);

CREATE TABLE job_position_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    position_id UUID NOT NULL REFERENCES job_positions_ats(id) ON DELETE CASCADE,
    skill VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    required BOOLEAN DEFAULT TRUE,
    min_years INT,
    weight NUMERIC(5,2) DEFAULT 1.0
);
CREATE INDEX idx_jps_position ON job_position_skills(position_id);

-- ============================================================
-- PUBLICACIONES
-- ============================================================

CREATE TABLE job_postings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    position_id UUID NOT NULL REFERENCES job_positions_ats(id),
    requisition_id UUID REFERENCES job_requisitions(id),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT,
    responsibilities TEXT,
    benefits TEXT,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    location VARCHAR(200),
    salary_min NUMERIC(14,2),
    salary_max NUMERIC(14,2),
    currency VARCHAR(10),
    published_at TIMESTAMPTZ,
    closing_at TIMESTAMPTZ,
    is_public BOOLEAN DEFAULT FALSE,
    external_url TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_jpo_company ON job_postings(company_id);
CREATE INDEX idx_jpo_position ON job_postings(position_id);
CREATE INDEX idx_jpo_status ON job_postings(status);
CREATE INDEX idx_jpo_public ON job_postings(is_public, status) WHERE is_public = TRUE;

CREATE TABLE posting_boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    config JSONB,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE posting_board_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id UUID NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    board_id UUID NOT NULL REFERENCES posting_boards(id) ON DELETE CASCADE,
    external_id TEXT,
    posted_at TIMESTAMPTZ,
    status VARCHAR(30) DEFAULT 'PENDING',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(posting_id, board_id)
);

CREATE TABLE posting_screening_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id UUID NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    question_type VARCHAR(30) NOT NULL DEFAULT 'BOOLEAN',
    options JSONB,
    required BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_psq_posting ON posting_screening_questions(posting_id);

-- ============================================================
-- CANDIDATOS
-- ============================================================

CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(200) NOT NULL,
    phone VARCHAR(50),
    phone_country_code VARCHAR(5),
    document_type VARCHAR(30),
    document_number VARCHAR(100),
    birth_date DATE,
    location VARCHAR(200),
    nationality VARCHAR(100),
    gender VARCHAR(30),
    linkedin_url TEXT,
    portfolio_url TEXT,
    github_url TEXT,
    personal_website TEXT,
    current_company VARCHAR(200),
    current_position VARCHAR(200),
    notice_period INT,
    salary_expectation_min NUMERIC(14,2),
    salary_expectation_max NUMERIC(14,2),
    salary_currency VARCHAR(10),
    availability VARCHAR(50),
    source VARCHAR(100),
    source_detail TEXT,
    is_employee_referral BOOLEAN DEFAULT FALSE,
    referrer_employee_id UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    blacklisted BOOLEAN DEFAULT FALSE,
    blacklist_reason TEXT,
    tags TEXT[],
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, email)
);
CREATE INDEX idx_cand_company ON candidates(company_id);
CREATE INDEX idx_cand_email ON candidates(company_id, email);
CREATE INDEX idx_cand_status ON candidates(status);
CREATE INDEX idx_cand_source ON candidates(source);
CREATE INDEX idx_cand_tags ON candidates USING gin(tags);
CREATE INDEX idx_cand_referrer ON candidates(referrer_employee_id);

CREATE TABLE candidate_education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    institution VARCHAR(200) NOT NULL,
    degree VARCHAR(200),
    field_of_study VARCHAR(200),
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN DEFAULT FALSE,
    grade VARCHAR(50),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ce_candidate ON candidate_education(candidate_id);

CREATE TABLE candidate_experience (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    company VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    location VARCHAR(200),
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN DEFAULT FALSE,
    description TEXT,
    achievements TEXT[],
    industry VARCHAR(100),
    salary NUMERIC(14,2),
    salary_currency VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cx_candidate ON candidate_experience(candidate_id);

CREATE TABLE candidate_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    skill VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    proficiency VARCHAR(30) DEFAULT 'INTERMEDIATE',
    years_experience NUMERIC(4,1),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(candidate_id, skill)
);
CREATE INDEX idx_cs_candidate ON candidate_skills(candidate_id);
CREATE INDEX idx_cs_skill ON candidate_skills(skill);

CREATE TABLE candidate_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    issuer VARCHAR(200),
    issue_date DATE,
    expiry_date DATE,
    credential_id VARCHAR(200),
    credential_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cc_candidate ON candidate_certifications(candidate_id);

CREATE TABLE candidate_languages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    language VARCHAR(100) NOT NULL,
    proficiency VARCHAR(30) NOT NULL DEFAULT 'BASIC',
    is_native BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(candidate_id, language)
);

CREATE TABLE candidate_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL DEFAULT 'CV',
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    storage_provider VARCHAR(30),
    storage_key TEXT,
    parsed_data JSONB,
    parsed_at TIMESTAMPTZ,
    embedding vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cd_candidate ON candidate_documents(candidate_id);
CREATE INDEX idx_cd_type ON candidate_documents(document_type);

CREATE TABLE candidate_parsed_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE UNIQUE,
    raw_text TEXT,
    structured_data JSONB,
    skills_found TEXT[],
    education_found JSONB,
    experience_found JSONB,
    languages_found JSONB,
    certifications_found JSONB,
    summary TEXT,
    score NUMERIC(5,2),
    parser_version VARCHAR(30),
    parsed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- POSTULACIONES (APPLICATIONS)
-- ============================================================

CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    posting_id UUID NOT NULL REFERENCES job_postings(id),
    current_stage_id UUID REFERENCES recruitment_stages(id),
    status VARCHAR(50) NOT NULL DEFAULT 'NEW',
    score NUMERIC(5,2),
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    rejection_reason_id UUID REFERENCES rejection_reasons(id),
    rejection_reason_text TEXT,
    hired_at TIMESTAMPTZ,
    withdrawn_at TIMESTAMPTZ,
    withdraw_reason TEXT,
    source VARCHAR(100),
    source_detail TEXT,
    is_internal_mobility BOOLEAN DEFAULT FALSE,
    consent_given BOOLEAN DEFAULT FALSE,
    consent_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_app_company ON applications(company_id);
CREATE INDEX idx_app_candidate ON applications(candidate_id);
CREATE INDEX idx_app_posting ON applications(posting_id);
CREATE INDEX idx_app_status ON applications(status);
CREATE INDEX idx_app_stage ON applications(current_stage_id);
CREATE INDEX idx_app_score ON applications(score DESC);
CREATE UNIQUE INDEX idx_app_unique ON applications(candidate_id, posting_id);

CREATE TABLE application_screening_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES posting_screening_questions(id),
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_asa_application ON application_screening_answers(application_id);

CREATE TABLE application_stage_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    from_stage_id UUID REFERENCES recruitment_stages(id),
    to_stage_id UUID NOT NULL REFERENCES recruitment_stages(id),
    changed_by UUID,
    reason TEXT,
    auto_transition BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ash_application ON application_stage_history(application_id);

CREATE TABLE application_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    rated_by UUID NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ar_application ON application_ratings(application_id);
CREATE INDEX idx_ar_rater ON application_ratings(rated_by);

CREATE TABLE application_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_an_application ON application_notes(application_id);

-- ============================================================
-- ENTREVISTAS
-- ============================================================

CREATE TABLE interviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL REFERENCES applications(id),
    interview_type VARCHAR(50) NOT NULL,
    title VARCHAR(200),
    scheduled_at TIMESTAMPTZ,
    duration_minutes INT,
    meeting_url TEXT,
    meeting_password TEXT,
    location TEXT,
    instructions TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED',
    score NUMERIC(5,2),
    notes TEXT,
    cancelled_at TIMESTAMPTZ,
    cancel_reason TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_int_company ON interviews(company_id);
CREATE INDEX idx_int_application ON interviews(application_id);
CREATE INDEX idx_int_status ON interviews(status);
CREATE INDEX idx_int_scheduled ON interviews(scheduled_at) WHERE status = 'SCHEDULED';

CREATE TABLE interview_panel (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_id UUID NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'INTERVIEWER',
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    response_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(interview_id, employee_id)
);
CREATE INDEX idx_ip_interview ON interview_panel(interview_id);

CREATE TABLE interview_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_id UUID NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    panelist_id UUID NOT NULL,
    score NUMERIC(5,2),
    comments TEXT,
    strengths TEXT[],
    weaknesses TEXT[],
    recommendation VARCHAR(30),
    submitted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_if_interview ON interview_feedback(interview_id);

CREATE TABLE interview_feedback_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_feedback_id UUID NOT NULL REFERENCES interview_feedback(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    score NUMERIC(5,2),
    comment TEXT
);
CREATE INDEX idx_ifq_feedback ON interview_feedback_questions(interview_feedback_id);

-- ============================================================
-- EVALUACIONES / ASSESSMENTS
-- ============================================================

CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL REFERENCES applications(id),
    assessment_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    max_score NUMERIC(8,2),
    passing_score NUMERIC(8,2),
    duration_minutes INT,
    due_at TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    score NUMERIC(8,2),
    result TEXT,
    result_summary TEXT,
    completed_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ass_company ON assessments(company_id);
CREATE INDEX idx_ass_application ON assessments(application_id);
CREATE INDEX idx_ass_type ON assessments(assessment_type);

CREATE TABLE assessment_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    max_score NUMERIC(8,2),
    weight NUMERIC(5,2) DEFAULT 1.0,
    sort_order INT DEFAULT 0
);
CREATE INDEX idx_asect_assessment ON assessment_sections(assessment_id);

CREATE TABLE assessment_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    section_id UUID REFERENCES assessment_sections(id),
    score NUMERIC(8,2),
    max_score NUMERIC(8,2),
    comment TEXT,
    graded_by UUID,
    graded_at TIMESTAMPTZ
);

-- ============================================================
-- OFERTAS
-- ============================================================

CREATE TABLE offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    application_id UUID NOT NULL REFERENCES applications(id),
    position_title VARCHAR(200) NOT NULL,
    department_id UUID,
    offer_type VARCHAR(50) NOT NULL DEFAULT 'STANDARD',
    start_date DATE,
    employment_type VARCHAR(50),
    work_mode VARCHAR(50),
    salary_amount NUMERIC(14,2),
    salary_currency VARCHAR(10),
    salary_period VARCHAR(30),
    variable_compensation TEXT,
    benefits_summary TEXT,
    equity_terms TEXT,
    conditions TEXT,
    notes TEXT,
    response_deadline TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    sent_at TIMESTAMPTZ,
    accepted_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    expired_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_off_company ON offers(company_id);
CREATE INDEX idx_off_application ON offers(application_id);
CREATE INDEX idx_off_status ON offers(status);

CREATE TABLE offer_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    offer_id UUID NOT NULL REFERENCES offers(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    step_order INT NOT NULL DEFAULT 1,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    comment TEXT,
    decided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_offapp_offer ON offer_approvals(offer_id);

CREATE TABLE offer_negotiations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    offer_id UUID NOT NULL REFERENCES offers(id) ON DELETE CASCADE,
    requested_by VARCHAR(30) NOT NULL,
    field VARCHAR(100) NOT NULL,
    original_value TEXT,
    requested_value TEXT NOT NULL,
    counter_value TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);
CREATE INDEX idx_offneg_offer ON offer_negotiations(offer_id);

CREATE TABLE offer_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    offer_id UUID NOT NULL REFERENCES offers(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL,
    signed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- HIRING PROCESS (transicion a onboarding)
-- ============================================================

CREATE TABLE hiring_processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    offer_id UUID REFERENCES offers(id),
    application_id UUID NOT NULL REFERENCES applications(id),
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    employee_id UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    background_check_status VARCHAR(30) DEFAULT 'PENDING',
    background_check_result TEXT,
    medical_check_status VARCHAR(30) DEFAULT 'PENDING',
    medical_check_result TEXT,
    document_verification_status VARCHAR(30) DEFAULT 'PENDING',
    start_date DATE,
    onboarding_status VARCHAR(30) DEFAULT 'NOT_STARTED',
    onboarding_id UUID,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_hp_company ON hiring_processes(company_id);
CREATE INDEX idx_hp_offer ON hiring_processes(offer_id);
CREATE INDEX idx_hp_application ON hiring_processes(application_id);
CREATE INDEX idx_hp_candidate ON hiring_processes(candidate_id);
CREATE INDEX idx_hp_status ON hiring_processes(status);

CREATE TABLE hiring_process_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    process_id UUID NOT NULL REFERENCES hiring_processes(id) ON DELETE CASCADE,
    task_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    assigned_to UUID,
    due_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_hpt_process ON hiring_process_tasks(process_id);

CREATE TABLE hiring_process_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    process_id UUID NOT NULL REFERENCES hiring_processes(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- WORKFLOWS / AUTOMATION
-- ============================================================

CREATE TABLE recruitment_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    entity_type VARCHAR(50) NOT NULL DEFAULT 'APPLICATION',
    is_default BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rw_company ON recruitment_workflows(company_id);

CREATE TABLE recruitment_workflow_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES recruitment_workflows(id) ON DELETE CASCADE,
    stage_id UUID NOT NULL REFERENCES recruitment_stages(id),
    sort_order INT NOT NULL DEFAULT 0,
    required_actions TEXT[],
    auto_advance BOOLEAN DEFAULT FALSE,
    auto_advance_delay_hours INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rws_workflow ON recruitment_workflow_stages(workflow_id);

CREATE TABLE recruitment_workflow_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES recruitment_workflows(id) ON DELETE CASCADE,
    trigger_event VARCHAR(50) NOT NULL,
    condition_expression TEXT,
    action_type VARCHAR(50) NOT NULL,
    action_config JSONB,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rwr_workflow ON recruitment_workflow_rules(workflow_id);

-- ============================================================
-- SCORING Y MATCHING
-- ============================================================

CREATE TABLE scoring_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    config JSONB,
    is_default BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_sm_company ON scoring_models(company_id);

CREATE TABLE scoring_criteria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id UUID NOT NULL REFERENCES scoring_models(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    field VARCHAR(100) NOT NULL,
    weight NUMERIC(5,2) NOT NULL DEFAULT 1.0,
    scoring_type VARCHAR(30) NOT NULL DEFAULT 'EXACT_MATCH',
    config JSONB,
    sort_order INT DEFAULT 0
);
CREATE INDEX idx_sc_model ON scoring_criteria(model_id);

CREATE TABLE matching_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    position_id UUID NOT NULL REFERENCES job_positions_ats(id),
    overall_score NUMERIC(5,2),
    skill_score NUMERIC(5,2),
    experience_score NUMERIC(5,2),
    education_score NUMERIC(5,2),
    culture_score NUMERIC(5,2),
    details JSONB,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(candidate_id, position_id)
);
CREATE INDEX idx_mr_candidate ON matching_results(candidate_id);
CREATE INDEX idx_mr_position ON matching_results(position_id);
CREATE INDEX idx_mr_score ON matching_results(overall_score DESC);

-- ============================================================
-- EMAILS / COMUNICACIONES
-- ============================================================

CREATE TABLE email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(100) NOT NULL,
    subject TEXT NOT NULL,
    body_html TEXT NOT NULL,
    body_text TEXT,
    variables TEXT[],
    category VARCHAR(50),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);
CREATE INDEX idx_et_company ON email_templates(company_id);

CREATE TABLE email_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    template_id UUID REFERENCES email_templates(id),
    application_id UUID REFERENCES applications(id),
    candidate_id UUID REFERENCES candidates(id),
    recipient_email VARCHAR(200) NOT NULL,
    subject TEXT NOT NULL,
    body TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'SENT',
    error_message TEXT,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_el_company ON email_log(company_id);
CREATE INDEX idx_el_application ON email_log(application_id);

-- ============================================================
-- REFERRALS
-- ============================================================

CREATE TABLE referral_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    referral_id UUID NOT NULL,
    reward_type VARCHAR(50) NOT NULL,
    amount NUMERIC(14,2),
    currency VARCHAR(10),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rr_company ON referral_rewards(company_id);

-- ============================================================
-- TALENT POOL / CANDIDATE POOL
-- ============================================================

CREATE TABLE talent_pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    criteria JSONB,
    is_auto BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tp_company ON talent_pools(company_id);

CREATE TABLE talent_pool_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pool_id UUID NOT NULL REFERENCES talent_pools(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    added_by UUID,
    added_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(pool_id, candidate_id)
);
CREATE INDEX idx_tpc_pool ON talent_pool_candidates(pool_id);
CREATE INDEX idx_tpc_candidate ON talent_pool_candidates(candidate_id);

-- ============================================================
-- AUDITORIA
-- ============================================================

CREATE TABLE recruitment_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ral_company ON recruitment_audit_log(company_id);
CREATE INDEX idx_ral_entity ON recruitment_audit_log(entity_type, entity_id);
CREATE INDEX idx_ral_created ON recruitment_audit_log(created_at DESC);

-- ============================================================
-- DASHBOARD CACHE / REPORTES
-- ============================================================

CREATE TABLE recruitment_dashboard_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cache_key VARCHAR(100) NOT NULL,
    cache_data JSONB NOT NULL,
    cached_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    UNIQUE(company_id, cache_key)
);

CREATE TABLE recruitment_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,
    filters JSONB,
    schedule_cron VARCHAR(100),
    last_run_at TIMESTAMPTZ,
    last_result JSONB,
    active BOOLEAN DEFAULT TRUE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- RBAC PERMISOS
-- ============================================================

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'ats.requisition.create', 'ats', 'requisition.create', 'Crear requisiciones', NOW()),
    (gen_random_uuid(), 'ats.requisition.read', 'ats', 'requisition.read', 'Ver requisiciones', NOW()),
    (gen_random_uuid(), 'ats.requisition.update', 'ats', 'requisition.update', 'Actualizar requisiciones', NOW()),
    (gen_random_uuid(), 'ats.requisition.approve', 'ats', 'requisition.approve', 'Aprobar requisiciones', NOW()),
    (gen_random_uuid(), 'ats.requisition.delete', 'ats', 'requisition.delete', 'Eliminar requisiciones', NOW()),
    (gen_random_uuid(), 'ats.position.create', 'ats', 'position.create', 'Crear posiciones', NOW()),
    (gen_random_uuid(), 'ats.position.read', 'ats', 'position.read', 'Ver posiciones', NOW()),
    (gen_random_uuid(), 'ats.position.update', 'ats', 'position.update', 'Actualizar posiciones', NOW()),
    (gen_random_uuid(), 'ats.posting.create', 'ats', 'posting.create', 'Crear publicaciones', NOW()),
    (gen_random_uuid(), 'ats.posting.read', 'ats', 'posting.read', 'Ver publicaciones', NOW()),
    (gen_random_uuid(), 'ats.posting.update', 'ats', 'posting.update', 'Actualizar publicaciones', NOW()),
    (gen_random_uuid(), 'ats.posting.publish', 'ats', 'posting.publish', 'Publicar vacantes', NOW()),
    (gen_random_uuid(), 'ats.candidate.create', 'ats', 'candidate.create', 'Crear candidatos', NOW()),
    (gen_random_uuid(), 'ats.candidate.read', 'ats', 'candidate.read', 'Ver candidatos', NOW()),
    (gen_random_uuid(), 'ats.candidate.update', 'ats', 'candidate.update', 'Actualizar candidatos', NOW()),
    (gen_random_uuid(), 'ats.candidate.delete', 'ats', 'candidate.delete', 'Eliminar candidatos', NOW()),
    (gen_random_uuid(), 'ats.candidate.blacklist', 'ats', 'candidate.blacklist', 'Bloquear candidatos', NOW()),
    (gen_random_uuid(), 'ats.application.read', 'ats', 'application.read', 'Ver postulaciones', NOW()),
    (gen_random_uuid(), 'ats.application.update', 'ats', 'application.update', 'Actualizar postulaciones', NOW()),
    (gen_random_uuid(), 'ats.application.stage', 'ats', 'application.stage', 'Cambiar etapa', NOW()),
    (gen_random_uuid(), 'ats.application.reject', 'ats', 'application.reject', 'Rechazar postulaciones', NOW()),
    (gen_random_uuid(), 'ats.interview.create', 'ats', 'interview.create', 'Crear entrevistas', NOW()),
    (gen_random_uuid(), 'ats.interview.read', 'ats', 'interview.read', 'Ver entrevistas', NOW()),
    (gen_random_uuid(), 'ats.interview.update', 'ats', 'interview.update', 'Actualizar entrevistas', NOW()),
    (gen_random_uuid(), 'ats.interview.feedback', 'ats', 'interview.feedback', 'Dar feedback', NOW()),
    (gen_random_uuid(), 'ats.assessment.create', 'ats', 'assessment.create', 'Crear evaluaciones', NOW()),
    (gen_random_uuid(), 'ats.assessment.read', 'ats', 'assessment.read', 'Ver evaluaciones', NOW()),
    (gen_random_uuid(), 'ats.assessment.score', 'ats', 'assessment.score', 'Calificar evaluaciones', NOW()),
    (gen_random_uuid(), 'ats.offer.create', 'ats', 'offer.create', 'Crear ofertas', NOW()),
    (gen_random_uuid(), 'ats.offer.read', 'ats', 'offer.read', 'Ver ofertas', NOW()),
    (gen_random_uuid(), 'ats.offer.update', 'ats', 'offer.update', 'Actualizar ofertas', NOW()),
    (gen_random_uuid(), 'ats.offer.send', 'ats', 'offer.send', 'Enviar ofertas', NOW()),
    (gen_random_uuid(), 'ats.offer.approve', 'ats', 'offer.approve', 'Aprobar ofertas', NOW()),
    (gen_random_uuid(), 'ats.hiring.read', 'ats', 'hiring.read', 'Ver contrataciones', NOW()),
    (gen_random_uuid(), 'ats.hiring.update', 'ats', 'hiring.update', 'Actualizar contrataciones', NOW()),
    (gen_random_uuid(), 'ats.workflow.manage', 'ats', 'workflow.manage', 'Gestionar workflows', NOW()),
    (gen_random_uuid(), 'ats.template.manage', 'ats', 'template.manage', 'Gestionar plantillas email', NOW()),
    (gen_random_uuid(), 'ats.report.read', 'ats', 'report.read', 'Ver reportes', NOW()),
    (gen_random_uuid(), 'ats.analytics', 'ats', 'analytics', 'Ver analytics', NOW()),
    (gen_random_uuid(), 'ats.pool.manage', 'ats', 'pool.manage', 'Gestionar talent pools', NOW()),
    (gen_random_uuid(), 'ats.settings.manage', 'ats', 'settings.manage', 'Gestionar configuracion ATS', NOW()),
    (gen_random_uuid(), 'ats.public.apply', 'ats', 'public.apply', 'Postular via portal publico', NOW()),
    (gen_random_uuid(), 'ats.public.read', 'ats', 'public.read', 'Ver vacantes publicas', NOW())
ON CONFLICT (name) DO NOTHING;
