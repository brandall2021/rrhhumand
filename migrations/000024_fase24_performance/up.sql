-- FASE 24: Performance Management (DDD)
-- Reemplaza completamente la migración 000013

-- ============================================================
-- 0. DROP OLD TABLES (000013)
-- ============================================================
DROP TABLE IF EXISTS performance_development_actions CASCADE;
DROP TABLE IF EXISTS performance_development_plans CASCADE;
DROP TABLE IF EXISTS performance_improvement_actions CASCADE;
DROP TABLE IF EXISTS performance_improvement_plans CASCADE;
DROP TABLE IF EXISTS performance_results CASCADE;
DROP TABLE IF EXISTS performance_evidence CASCADE;
DROP TABLE IF EXISTS performance_feedback CASCADE;
DROP TABLE IF EXISTS performance_evaluation_answers CASCADE;
DROP TABLE IF EXISTS performance_evaluations CASCADE;
DROP TABLE IF EXISTS performance_evaluators CASCADE;
DROP TABLE IF EXISTS performance_kpis CASCADE;
DROP TABLE IF EXISTS performance_objectives CASCADE;
DROP TABLE IF EXISTS performance_scoring_rules CASCADE;
DROP TABLE IF EXISTS competencies CASCADE;
DROP TABLE IF EXISTS rating_scale_levels CASCADE;
DROP TABLE IF EXISTS rating_scales CASCADE;
DROP TABLE IF EXISTS template_section_items CASCADE;
DROP TABLE IF EXISTS template_sections CASCADE;
DROP TABLE IF EXISTS evaluation_templates CASCADE;
DROP TABLE IF EXISTS performance_cycles CASCADE;
DROP TABLE IF EXISTS performance_audit_log CASCADE;

-- ============================================================
-- 1. PERFORMANCE CYCLES
-- ============================================================
CREATE TABLE performance_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    cycle_type VARCHAR(50) NOT NULL DEFAULT 'ANNUAL',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    start_date DATE,
    end_date DATE,
    evaluation_start_date DATE,
    evaluation_end_date DATE,
    review_start_date DATE,
    review_end_date DATE,
    calibration_start_date DATE,
    calibration_end_date DATE,
    template_id UUID,
    objective_weight NUMERIC(5,2) DEFAULT 60.00,
    competency_weight NUMERIC(5,2) DEFAULT 40.00,
    min_anonymous_responses INT DEFAULT 3,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 2. PERFORMANCE TEMPLATES
-- ============================================================
CREATE TABLE performance_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    evaluation_type VARCHAR(50) NOT NULL DEFAULT 'MANAGER',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE template_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES performance_templates(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    section_type VARCHAR(50) NOT NULL DEFAULT 'SCALE',
    weight NUMERIC(5,2) DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE template_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES performance_templates(id) ON DELETE CASCADE,
    section_id UUID REFERENCES template_sections(id) ON DELETE SET NULL,
    question TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL DEFAULT 'SCALE',
    required BOOLEAN NOT NULL DEFAULT TRUE,
    weight NUMERIC(5,2) DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================
-- 3. RATING SCALES
-- ============================================================
CREATE TABLE rating_scales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    min_value NUMERIC(5,2) NOT NULL DEFAULT 1,
    max_value NUMERIC(5,2) NOT NULL DEFAULT 5,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rating_scale_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scale_id UUID NOT NULL REFERENCES rating_scales(id) ON DELETE CASCADE,
    value NUMERIC(5,2) NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

-- ============================================================
-- 4. COMPETENCIES
-- ============================================================
CREATE TABLE competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    competency_type VARCHAR(50) NOT NULL DEFAULT 'BEHAVIORAL',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE competency_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    level INT NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE position_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    position_id UUID NOT NULL,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    expected_level INT NOT NULL,
    weight NUMERIC(5,2) DEFAULT 0,
    UNIQUE(company_id, position_id, competency_id)
);

CREATE TABLE cycle_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    expected_level INT NOT NULL DEFAULT 3,
    weight NUMERIC(5,2) DEFAULT 0,
    UNIQUE(cycle_id, employee_id, competency_id)
);

-- ============================================================
-- 5. OBJECTIVES & KEY RESULTS
-- ============================================================
CREATE TABLE performance_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    parent_objective_id UUID REFERENCES performance_objectives(id) ON DELETE SET NULL,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    objective_type VARCHAR(50) NOT NULL DEFAULT 'INDIVIDUAL',
    weight NUMERIC(5,2) DEFAULT 0,
    start_date DATE,
    due_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'NOT_STARTED',
    progress NUMERIC(5,2) DEFAULT 0,
    target_value NUMERIC,
    current_value NUMERIC,
    unit VARCHAR(50),
    progress_type VARCHAR(30) NOT NULL DEFAULT 'PERCENTAGE',
    notes TEXT,
    risk_notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE objective_key_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    objective_id UUID NOT NULL REFERENCES performance_objectives(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    weight NUMERIC(5,2) DEFAULT 0,
    target_value NUMERIC,
    current_value NUMERIC DEFAULT 0,
    unit VARCHAR(50),
    progress NUMERIC(5,2) DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'NOT_STARTED',
    sort_order INT NOT NULL DEFAULT 0
);

-- ============================================================
-- 6. EVALUATION PARTICIPANTS
-- ============================================================
CREATE TABLE performance_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    evaluator_id UUID NOT NULL,
    evaluation_type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    submitted_at TIMESTAMPTZ,
    UNIQUE(cycle_id, employee_id, evaluator_id, evaluation_type)
);

-- ============================================================
-- 7. EVALUATIONS
-- ============================================================
CREATE TABLE performance_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    evaluator_id UUID NOT NULL,
    evaluation_type VARCHAR(50) NOT NULL,
    template_id UUID REFERENCES performance_templates(id),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    overall_score NUMERIC(5,2),
    strengths TEXT,
    improvement_areas TEXT,
    summary TEXT,
    submitted_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ,
    locked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE evaluation_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evaluation_id UUID NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    question_id UUID REFERENCES template_questions(id) ON DELETE SET NULL,
    numeric_value NUMERIC(5,2),
    text_value TEXT,
    selected_value VARCHAR(200),
    boolean_value BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 8. OBJECTIVE & COMPETENCY EVALUATIONS
-- ============================================================
CREATE TABLE objective_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    objective_id UUID NOT NULL REFERENCES performance_objectives(id) ON DELETE CASCADE,
    evaluation_id UUID NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    score NUMERIC(5,2),
    comment TEXT,
    UNIQUE(objective_id, evaluation_id)
);

CREATE TABLE competency_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    evaluation_id UUID NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    score NUMERIC(5,2),
    expected_level INT,
    comment TEXT,
    UNIQUE(competency_id, evaluation_id)
);

-- ============================================================
-- 9. FEEDBACK
-- ============================================================
CREATE TABLE performance_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID REFERENCES performance_cycles(id) ON DELETE SET NULL,
    employee_id UUID NOT NULL,
    author_id UUID NOT NULL,
    feedback_type VARCHAR(50) NOT NULL DEFAULT 'GENERAL',
    visibility VARCHAR(30) NOT NULL DEFAULT 'EMPLOYEE',
    content TEXT NOT NULL,
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE performance_recognitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    author_id UUID NOT NULL,
    recognition_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 10. CHECK-INS
-- ============================================================
CREATE TABLE performance_checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    manager_id UUID NOT NULL,
    cycle_id UUID REFERENCES performance_cycles(id) ON DELETE SET NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    employee_notes TEXT,
    manager_notes TEXT,
    achievements TEXT,
    blockers TEXT,
    next_steps TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 11. REVIEWS
-- ============================================================
CREATE TABLE performance_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    manager_id UUID NOT NULL,
    summary TEXT,
    strengths TEXT,
    improvement_areas TEXT,
    final_score NUMERIC(5,2),
    final_rating VARCHAR(50),
    employee_comments TEXT,
    manager_comments TEXT,
    employee_agreement VARCHAR(30) DEFAULT 'PENDING',
    disagreement_reason TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(cycle_id, employee_id)
);

-- ============================================================
-- 12. CALIBRATION
-- ============================================================
CREATE TABLE calibration_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE calibration_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES calibration_sessions(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    original_score NUMERIC(5,2),
    adjusted_score NUMERIC(5,2),
    original_rating VARCHAR(50),
    adjusted_rating VARCHAR(50),
    reason TEXT,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    UNIQUE(session_id, employee_id)
);

-- ============================================================
-- 13. IMPROVEMENT PLANS
-- ============================================================
CREATE TABLE performance_improvement_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID REFERENCES performance_cycles(id) ON DELETE SET NULL,
    created_by UUID,
    reason TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    success_criteria TEXT,
    final_result TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE improvement_plan_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES performance_improvement_plans(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    responsible_id UUID,
    due_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    progress NUMERIC(5,2) DEFAULT 0,
    evidence TEXT,
    completed_at TIMESTAMPTZ
);

-- ============================================================
-- 14. DEVELOPMENT PLANS
-- ============================================================
CREATE TABLE performance_development_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID REFERENCES performance_cycles(id) ON DELETE SET NULL,
    created_by UUID,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    career_goal TEXT,
    current_level INT,
    target_level INT,
    competency_id UUID REFERENCES competencies(id) ON DELETE SET NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE development_plan_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES performance_development_plans(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    action_type VARCHAR(50) NOT NULL DEFAULT 'COURSE',
    due_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ
);

-- ============================================================
-- 15. EVIDENCE
-- ============================================================
CREATE TABLE performance_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    evaluation_id UUID REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    objective_id UUID REFERENCES performance_objectives(id) ON DELETE CASCADE,
    feedback_id UUID REFERENCES performance_feedback(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT,
    evidence_type VARCHAR(50) NOT NULL DEFAULT 'FILE',
    storage_key TEXT,
    file_name VARCHAR(250),
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    url TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT evidence_at_least_one_target CHECK (
        (evaluation_id IS NOT NULL)::INT +
        (objective_id IS NOT NULL)::INT +
        (feedback_id IS NOT NULL)::INT >= 1
    )
);

-- ============================================================
-- 16. RESULTS (denormalized for fast reads)
-- ============================================================
CREATE TABLE performance_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES performance_cycles(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    objective_score NUMERIC(5,2),
    competency_score NUMERIC(5,2),
    self_score NUMERIC(5,2),
    manager_score NUMERIC(5,2),
    peer_score NUMERIC(5,2),
    hr_score NUMERIC(5,2),
    final_score NUMERIC(5,2),
    final_rating VARCHAR(50),
    final_rating_label VARCHAR(100),
    strengths TEXT,
    improvement_areas TEXT,
    summary TEXT,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(cycle_id, employee_id)
);

-- ============================================================
-- 17. AUDIT LOG
-- ============================================================
CREATE TABLE performance_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 18. OUTBOX EVENTS (module-specific)
-- ============================================================
CREATE TABLE performance_outbox (
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

-- ============================================================
-- 19. PERFORMANCE PERMISSIONS
-- ============================================================
INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'performance.read', 'performance', 'read', 'View performance data', NOW()),
    (gen_random_uuid(), 'performance.cycles.read', 'performance', 'cycles.read', 'View cycles', NOW()),
    (gen_random_uuid(), 'performance.cycles.create', 'performance', 'cycles.create', 'Create cycles', NOW()),
    (gen_random_uuid(), 'performance.cycles.manage', 'performance', 'cycles.manage', 'Manage cycles (start, close, finalize)', NOW()),
    (gen_random_uuid(), 'performance.objectives.read', 'performance', 'objectives.read', 'View objectives', NOW()),
    (gen_random_uuid(), 'performance.objectives.create', 'performance', 'objectives.create', 'Create objectives', NOW()),
    (gen_random_uuid(), 'performance.objectives.update', 'performance', 'objectives.update', 'Update objectives', NOW()),
    (gen_random_uuid(), 'performance.objectives.manage', 'performance', 'objectives.manage', 'Manage objectives', NOW()),
    (gen_random_uuid(), 'performance.evaluations.read', 'performance', 'evaluations.read', 'View evaluations', NOW()),
    (gen_random_uuid(), 'performance.evaluations.create', 'performance', 'evaluations.create', 'Create evaluations', NOW()),
    (gen_random_uuid(), 'performance.evaluations.submit', 'performance', 'evaluations.submit', 'Submit evaluations', NOW()),
    (gen_random_uuid(), 'performance.evaluations.approve', 'performance', 'evaluations.approve', 'Approve evaluations', NOW()),
    (gen_random_uuid(), 'performance.evaluations.lock', 'performance', 'evaluations.lock', 'Lock evaluations', NOW()),
    (gen_random_uuid(), 'performance.feedback.read', 'performance', 'feedback.read', 'View feedback', NOW()),
    (gen_random_uuid(), 'performance.feedback.create', 'performance', 'feedback.create', 'Create feedback', NOW()),
    (gen_random_uuid(), 'performance.calibration.read', 'performance', 'calibration.read', 'View calibration', NOW()),
    (gen_random_uuid(), 'performance.calibration.manage', 'performance', 'calibration.manage', 'Manage calibration', NOW()),
    (gen_random_uuid(), 'performance.improvement.read', 'performance', 'improvement.read', 'View improvement plans', NOW()),
    (gen_random_uuid(), 'performance.improvement.create', 'performance', 'improvement.create', 'Create improvement plans', NOW()),
    (gen_random_uuid(), 'performance.improvement.manage', 'performance', 'improvement.manage', 'Manage improvement plans', NOW()),
    (gen_random_uuid(), 'performance.development.read', 'performance', 'development.read', 'View development plans', NOW()),
    (gen_random_uuid(), 'performance.development.create', 'performance', 'development.create', 'Create development plans', NOW()),
    (gen_random_uuid(), 'performance.development.manage', 'performance', 'development.manage', 'Manage development plans', NOW()),
    (gen_random_uuid(), 'performance.ai.use', 'performance', 'ai.use', 'Use AI assistants', NOW()),
    (gen_random_uuid(), 'performance.admin', 'performance', 'admin', 'Full performance administration', NOW())
ON CONFLICT (name) DO NOTHING;

-- ============================================================
-- 20. INDEXES
-- ============================================================
CREATE INDEX idx_perf_cycles_company_status ON performance_cycles(company_id, status);
CREATE INDEX idx_perf_cycles_type ON performance_cycles(cycle_type);
CREATE INDEX idx_perf_objectives_employee ON performance_objectives(employee_id);
CREATE INDEX idx_perf_objectives_cycle ON performance_objectives(cycle_id);
CREATE INDEX idx_perf_objectives_status ON performance_objectives(status);
CREATE INDEX idx_perf_objectives_parent ON performance_objectives(parent_objective_id);
CREATE INDEX idx_perf_key_results_objective ON objective_key_results(objective_id);
CREATE INDEX idx_perf_participants_cycle ON performance_participants(cycle_id);
CREATE INDEX idx_perf_participants_employee ON performance_participants(employee_id);
CREATE INDEX idx_perf_participants_evaluator ON performance_participants(evaluator_id);
CREATE INDEX idx_perf_evaluations_employee ON performance_evaluations(employee_id);
CREATE INDEX idx_perf_evaluations_cycle ON performance_evaluations(cycle_id);
CREATE INDEX idx_perf_evaluations_evaluator ON performance_evaluations(evaluator_id);
CREATE INDEX idx_perf_evaluations_status ON performance_evaluations(status);
CREATE INDEX idx_perf_feedback_employee ON performance_feedback(employee_id);
CREATE INDEX idx_perf_feedback_author ON performance_feedback(author_id);
CREATE INDEX idx_perf_checkins_employee ON performance_checkins(employee_id);
CREATE INDEX idx_perf_checkins_manager ON performance_checkins(manager_id);
CREATE INDEX idx_perf_reviews_cycle ON performance_reviews(cycle_id);
CREATE INDEX idx_perf_reviews_employee ON performance_reviews(employee_id);
CREATE INDEX idx_perf_calibration_session ON calibration_items(session_id);
CREATE INDEX idx_perf_improvement_employee ON performance_improvement_plans(employee_id);
CREATE INDEX idx_perf_development_employee ON performance_development_plans(employee_id);
CREATE INDEX idx_perf_results_cycle ON performance_results(cycle_id);
CREATE INDEX idx_perf_results_employee ON performance_results(employee_id);
CREATE INDEX idx_perf_evidence_evaluation ON performance_evidence(evaluation_id);
CREATE INDEX idx_perf_evidence_objective ON performance_evidence(objective_id);
CREATE INDEX idx_perf_audit_company ON performance_audit_log(company_id);
CREATE INDEX idx_perf_audit_entity ON performance_audit_log(entity_type, entity_id);
CREATE INDEX idx_perf_outbox_status ON performance_outbox(status);
CREATE INDEX idx_perf_competencies_category ON competencies(category);
CREATE INDEX idx_perf_position_competencies ON position_competencies(position_id);
