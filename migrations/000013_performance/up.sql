-- FASE 14: Evaluación de Desempeño

-- Ciclos de evaluación
CREATE TABLE performance_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    evaluation_deadline DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_perf_cycles_company ON performance_cycles(company_id);
CREATE INDEX idx_perf_cycles_status ON performance_cycles(status);

-- Plantillas de evaluación
CREATE TABLE evaluation_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_eval_templates_company ON evaluation_templates(company_id);

-- Secciones de plantilla
CREATE TABLE template_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES evaluation_templates(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    section_type VARCHAR(30) NOT NULL,
    weight NUMERIC(5,2) DEFAULT 0,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT TRUE
);
CREATE INDEX idx_template_sections_template ON template_sections(template_id);

-- Ítems de sección (competencias/objetivos/KPI de la plantilla)
CREATE TABLE template_section_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL REFERENCES template_sections(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    item_type VARCHAR(30) NOT NULL DEFAULT 'TEXT',
    weight NUMERIC(5,2) DEFAULT 0,
    sort_order INT DEFAULT 0
);
CREATE INDEX idx_template_items_section ON template_section_items(section_id);

-- Escalas de calificación
CREATE TABLE rating_scales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    min_value NUMERIC(8,2) NOT NULL,
    max_value NUMERIC(8,2) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_rating_scales_company ON rating_scales(company_id);

-- Niveles de escala
CREATE TABLE rating_scale_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scale_id UUID NOT NULL REFERENCES rating_scales(id) ON DELETE CASCADE,
    value NUMERIC(8,2) NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0
);
CREATE INDEX idx_rating_levels_scale ON rating_scale_levels(scale_id);

-- Competencias
CREATE TABLE competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_competencies_company ON competencies(company_id);

-- Reglas de puntuación por empresa
CREATE TABLE performance_scoring_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID,
    objective_weight NUMERIC(5,2) NOT NULL DEFAULT 40,
    competency_weight NUMERIC(5,2) NOT NULL DEFAULT 30,
    kpi_weight NUMERIC(5,2) NOT NULL DEFAULT 20,
    self_eval_weight NUMERIC(5,2) NOT NULL DEFAULT 10,
    manager_weight NUMERIC(5,2) NOT NULL DEFAULT 60,
    peer_weight NUMERIC(5,2) NOT NULL DEFAULT 20,
    hr_weight NUMERIC(5,2) NOT NULL DEFAULT 10,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_scoring_rules_company ON performance_scoring_rules(company_id);

-- Objetivos de desempeño
CREATE TABLE performance_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    metric VARCHAR(100),
    target_value NUMERIC(14,4),
    current_value NUMERIC(14,4),
    unit VARCHAR(50),
    weight NUMERIC(5,2) DEFAULT 0,
    start_date DATE,
    due_date DATE,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_objectives_company ON performance_objectives(company_id);
CREATE INDEX idx_perf_objectives_employee ON performance_objectives(employee_id);
CREATE INDEX idx_perf_objectives_cycle ON performance_objectives(cycle_id);

-- KPI
CREATE TABLE performance_kpis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    target_value NUMERIC(14,4),
    current_value NUMERIC(14,4),
    unit VARCHAR(50),
    weight NUMERIC(5,2) DEFAULT 0,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_kpis_company ON performance_kpis(company_id);
CREATE INDEX idx_perf_kpis_employee ON performance_kpis(employee_id);
CREATE INDEX idx_perf_kpis_cycle ON performance_kpis(cycle_id);

-- Evaluadores asignados
CREATE TABLE performance_evaluators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    evaluator_id UUID NOT NULL,
    evaluator_type VARCHAR(30) NOT NULL,
    status VARCHAR(30) DEFAULT 'PENDING',
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX idx_perf_evaluators_company ON performance_evaluators(company_id);
CREATE INDEX idx_perf_evaluators_cycle ON performance_evaluators(cycle_id);
CREATE INDEX idx_perf_evaluators_employee ON performance_evaluators(employee_id);
CREATE UNIQUE INDEX idx_perf_evaluators_unique ON performance_evaluators(cycle_id, employee_id, evaluator_id, evaluator_type);

-- Evaluaciones
CREATE TABLE performance_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    evaluator_id UUID NOT NULL,
    evaluator_type VARCHAR(30) NOT NULL,
    template_id UUID,
    status VARCHAR(30) DEFAULT 'DRAFT',
    overall_score NUMERIC(8,2),
    comments TEXT,
    submitted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_evaluations_company ON performance_evaluations(company_id);
CREATE INDEX idx_perf_evaluations_cycle ON performance_evaluations(cycle_id);
CREATE INDEX idx_perf_evaluations_employee ON performance_evaluations(employee_id);
CREATE INDEX idx_perf_evaluations_evaluator ON performance_evaluations(evaluator_id);

-- Respuestas de evaluación
CREATE TABLE performance_evaluation_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evaluation_id UUID NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    section_name VARCHAR(100),
    item_name VARCHAR(150) NOT NULL,
    item_type VARCHAR(30) NOT NULL DEFAULT 'TEXT',
    score NUMERIC(8,2),
    value TEXT,
    comments TEXT,
    weight NUMERIC(5,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_answers_evaluation ON performance_evaluation_answers(evaluation_id);

-- Feedback continuo
CREATE TABLE performance_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID,
    from_user_id UUID NOT NULL,
    feedback_type VARCHAR(30) NOT NULL DEFAULT 'GENERAL',
    message TEXT NOT NULL,
    is_private BOOLEAN DEFAULT FALSE,
    visible_to_employee BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_feedback_company ON performance_feedback(company_id);
CREATE INDEX idx_perf_feedback_employee ON performance_feedback(employee_id);

-- Evidencias
CREATE TABLE performance_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    evaluation_id UUID NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    evidence_type VARCHAR(30) NOT NULL DEFAULT 'COMMENT',
    storage_provider VARCHAR(30),
    storage_key TEXT,
    file_name VARCHAR(255),
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    url TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_evidence_evaluation ON performance_evidence(evaluation_id);

-- Resultados finales
CREATE TABLE performance_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    cycle_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    objective_score NUMERIC(8,2),
    competency_score NUMERIC(8,2),
    kpi_score NUMERIC(8,2),
    self_score NUMERIC(8,2),
    manager_score NUMERIC(8,2),
    peer_score NUMERIC(8,2),
    hr_score NUMERIC(8,2),
    final_score NUMERIC(8,2),
    rating VARCHAR(50),
    rating_label VARCHAR(100),
    strengths TEXT,
    areas_to_improve TEXT,
    calculated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(cycle_id, employee_id)
);
CREATE INDEX idx_perf_results_company ON performance_results(company_id);
CREATE INDEX idx_perf_results_cycle ON performance_results(cycle_id);
CREATE INDEX idx_perf_results_employee ON performance_results(employee_id);

-- Planes de mejora
CREATE TABLE performance_improvement_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID,
    result_id UUID,
    title VARCHAR(200) NOT NULL,
    problem_description TEXT,
    objective TEXT,
    responsible_id UUID,
    due_date DATE,
    status VARCHAR(30) DEFAULT 'OPEN',
    outcome TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_imp_plans_company ON performance_improvement_plans(company_id);
CREATE INDEX idx_perf_imp_plans_employee ON performance_improvement_plans(employee_id);

-- Acciones de plan de mejora
CREATE TABLE performance_improvement_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES performance_improvement_plans(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    due_date DATE,
    status VARCHAR(30) DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_imp_actions_plan ON performance_improvement_actions(plan_id);

-- Planes de desarrollo individual (IDP)
CREATE TABLE performance_development_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    cycle_id UUID,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    career_goal TEXT,
    timeline_months INT DEFAULT 12,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_dev_plans_company ON performance_development_plans(company_id);
CREATE INDEX idx_perf_dev_plans_employee ON performance_development_plans(employee_id);

-- Acciones de plan de desarrollo
CREATE TABLE performance_development_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES performance_development_plans(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    action_type VARCHAR(50) NOT NULL DEFAULT 'TRAINING',
    due_date DATE,
    status VARCHAR(30) DEFAULT 'PENDING',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_dev_actions_plan ON performance_development_actions(plan_id);

-- Log de auditoría
CREATE TABLE performance_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    user_id UUID,
    employee_id UUID,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_perf_audit_company ON performance_audit_log(company_id);
CREATE INDEX idx_perf_audit_entity ON performance_audit_log(entity_type, entity_id);

-- Permisos de performance
INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'performance.read', 'performance', 'read', 'View performance data', NOW()),
    (gen_random_uuid(), 'performance.manage_cycles', 'performance', 'manage_cycles', 'Create and manage evaluation cycles', NOW()),
    (gen_random_uuid(), 'performance.manage_templates', 'performance', 'manage_templates', 'Manage evaluation templates', NOW()),
    (gen_random_uuid(), 'performance.manage_competencies', 'performance', 'manage_competencies', 'Manage competencies', NOW()),
    (gen_random_uuid(), 'performance.manage_scales', 'performance', 'manage_scales', 'Manage rating scales', NOW()),
    (gen_random_uuid(), 'performance.create_objectives', 'performance', 'create_objectives', 'Create and manage objectives', NOW()),
    (gen_random_uuid(), 'performance.manage_kpis', 'performance', 'manage_kpis', 'Create and manage KPIs', NOW()),
    (gen_random_uuid(), 'performance.assign_evaluators', 'performance', 'assign_evaluators', 'Assign evaluators to employees', NOW()),
    (gen_random_uuid(), 'performance.evaluate', 'performance', 'evaluate', 'Perform evaluations', NOW()),
    (gen_random_uuid(), 'performance.submit_evaluation', 'performance', 'submit_evaluation', 'Submit evaluations for review', NOW()),
    (gen_random_uuid(), 'performance.approve_evaluation', 'performance', 'approve_evaluation', 'Approve finalized evaluations', NOW()),
    (gen_random_uuid(), 'performance.create_feedback', 'performance', 'create_feedback', 'Create feedback entries', NOW()),
    (gen_random_uuid(), 'performance.manage_plans', 'performance', 'manage_plans', 'Manage improvement and development plans', NOW()),
    (gen_random_uuid(), 'performance.view_reports', 'performance', 'view_reports', 'View performance reports and dashboards', NOW()),
    (gen_random_uuid(), 'performance.manage_development', 'performance', 'manage_development', 'Manage development plans', NOW())
ON CONFLICT (name) DO NOTHING;
