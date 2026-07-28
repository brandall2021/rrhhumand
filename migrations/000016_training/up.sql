-- FASE 17: Capacitación, Formación y LMS

-- Categorías de cursos (jerárquicas)
CREATE TABLE course_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    parent_id UUID REFERENCES course_categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cc_company ON course_categories(company_id);
CREATE INDEX idx_cc_parent ON course_categories(parent_id);

-- Cursos
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    category_id UUID REFERENCES course_categories(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    short_description TEXT,
    description TEXT,
    objectives TEXT,
    difficulty VARCHAR(30) DEFAULT 'BEGINNER',
    duration_minutes INT NOT NULL DEFAULT 0,
    modality VARCHAR(30) NOT NULL DEFAULT 'ONLINE',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    mandatory BOOLEAN DEFAULT FALSE,
    passing_score NUMERIC(5,2),
    certificate_enabled BOOLEAN DEFAULT FALSE,
    min_attendance_percentage NUMERIC(5,2) DEFAULT 80.00,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_courses_company ON courses(company_id);
CREATE INDEX idx_courses_status ON courses(company_id, status);
CREATE INDEX idx_courses_category ON courses(category_id);
CREATE UNIQUE INDEX idx_courses_code ON courses(company_id, code);

-- Versiones de curso
CREATE TABLE course_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    version VARCHAR(20) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cv_course ON course_versions(course_id);

-- Contenidos del curso
CREATE TABLE course_contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_version_id UUID NOT NULL REFERENCES course_versions(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    content_type VARCHAR(30) NOT NULL DEFAULT 'VIDEO',
    storage_provider VARCHAR(30),
    storage_key TEXT,
    external_url TEXT,
    duration_seconds INT DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    required BOOLEAN DEFAULT TRUE,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cc_version ON course_contents(course_version_id);

-- Prerrequisitos de cursos
CREATE TABLE course_prerequisites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    prerequisite_course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    required BOOLEAN DEFAULT TRUE,
    UNIQUE(course_id, prerequisite_course_id)
);
CREATE INDEX idx_cp_course ON course_prerequisites(course_id);
CREATE INDEX idx_cp_prereq ON course_prerequisites(prerequisite_course_id);

-- Instructores
CREATE TABLE instructors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID,
    instructor_type VARCHAR(30) NOT NULL DEFAULT 'EMPLOYEE',
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200),
    phone VARCHAR(50),
    specialization TEXT,
    bio TEXT,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_instructors_company ON instructors(company_id);

-- Proveedores de capacitación
CREATE TABLE training_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    tax_id VARCHAR(50),
    email VARCHAR(200),
    phone VARCHAR(50),
    website TEXT,
    contact_name VARCHAR(200),
    status VARCHAR(30) DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tp_company ON training_providers(company_id);

-- Ediciones de curso (offerings)
CREATE TABLE course_offerings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    course_version_id UUID REFERENCES course_versions(id),
    name VARCHAR(200) NOT NULL,
    start_date DATE,
    end_date DATE,
    enrollment_start TIMESTAMPTZ,
    enrollment_end TIMESTAMPTZ,
    capacity INT DEFAULT 0,
    enrolled_count INT DEFAULT 0,
    modality VARCHAR(30),
    location TEXT,
    meeting_url TEXT,
    instructor_id UUID REFERENCES instructors(id),
    provider_id UUID REFERENCES training_providers(id),
    cost_amount NUMERIC(14,2),
    cost_currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_co_company ON course_offerings(company_id);
CREATE INDEX idx_co_course ON course_offerings(course_id);
CREATE INDEX idx_co_dates ON course_offerings(company_id, start_date, end_date);
CREATE INDEX idx_co_status ON course_offerings(status);

-- Sesiones de una edición
CREATE TABLE training_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    title VARCHAR(200),
    session_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    location TEXT,
    meeting_url TEXT,
    instructor_id UUID REFERENCES instructors(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ts_offering ON training_sessions(offering_id);

-- Inscripciones
CREATE TABLE course_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    assignment_type VARCHAR(30) DEFAULT 'OPTIONAL',
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    enrollment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    dropped_at TIMESTAMPTZ,
    final_score NUMERIC(5,2),
    passed BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ce_company ON course_enrollments(company_id);
CREATE INDEX idx_ce_offering ON course_enrollments(offering_id);
CREATE INDEX idx_ce_employee ON course_enrollments(employee_id);
CREATE INDEX idx_ce_status ON course_enrollments(status);

-- Asignaciones masivas de capacitación
CREATE TABLE training_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    assignee_type VARCHAR(30) NOT NULL,
    assignee_id UUID,
    assignment_type VARCHAR(30) DEFAULT 'MANDATORY',
    due_date DATE,
    active BOOLEAN DEFAULT TRUE,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ta_company ON training_assignments(company_id);
CREATE INDEX idx_ta_course ON training_assignments(course_id);
CREATE INDEX idx_ta_assignee ON training_assignments(assignee_type, assignee_id);

-- Reglas de asignación automática
CREATE TABLE training_assignment_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    criteria_field VARCHAR(50) NOT NULL,
    criteria_value VARCHAR(200) NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    assignment_type VARCHAR(30) DEFAULT 'MANDATORY',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tar_company ON training_assignment_rules(company_id);

-- Progreso del contenido
CREATE TABLE course_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES course_enrollments(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES course_contents(id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL DEFAULT 'NOT_STARTED',
    progress_percentage INT DEFAULT 0,
    time_spent_seconds INT DEFAULT 0,
    last_position INT DEFAULT 0,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(enrollment_id, content_id)
);
CREATE INDEX idx_cp_enrollment ON course_progress(enrollment_id);

-- Evaluaciones
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    assessment_type VARCHAR(30) NOT NULL DEFAULT 'QUIZ',
    attempts_allowed INT DEFAULT 1,
    passing_score NUMERIC(5,2),
    time_limit_minutes INT,
    randomize_questions BOOLEAN DEFAULT FALSE,
    show_results BOOLEAN DEFAULT TRUE,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_assessments_course ON assessments(course_id);

-- Preguntas de evaluación
CREATE TABLE assessment_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    question_type VARCHAR(30) NOT NULL DEFAULT 'SINGLE_CHOICE',
    points NUMERIC(5,2) DEFAULT 1,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_aq_assessment ON assessment_questions(assessment_id);

-- Opciones de pregunta
CREATE TABLE assessment_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES assessment_questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    sort_order INT DEFAULT 0
);
CREATE INDEX idx_ao_question ON assessment_options(question_id);

-- Intentos de evaluación
CREATE TABLE assessment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES course_enrollments(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    attempt_number INT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'IN_PROGRESS',
    score NUMERIC(5,2),
    passed BOOLEAN,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_aa_assessment ON assessment_attempts(assessment_id);
CREATE INDEX idx_aa_enrollment ON assessment_attempts(enrollment_id);
CREATE INDEX idx_aa_employee ON assessment_attempts(employee_id);

-- Respuestas del intento
CREATE TABLE assessment_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES assessment_questions(id) ON DELETE CASCADE,
    selected_option_id UUID REFERENCES assessment_options(id),
    text_answer TEXT,
    numeric_answer NUMERIC(14,4),
    is_correct BOOLEAN,
    score NUMERIC(5,2) DEFAULT 0
);
CREATE INDEX idx_aa_attempt ON assessment_answers(attempt_id);

-- Certificados
CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES course_enrollments(id) ON DELETE CASCADE,
    certificate_no VARCHAR(50) NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    storage_provider VARCHAR(30),
    storage_key TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(certificate_no)
);
CREATE INDEX idx_cert_company ON certificates(company_id);
CREATE INDEX idx_cert_employee ON certificates(employee_id);
CREATE INDEX idx_cert_expires ON certificates(expires_at);
CREATE INDEX idx_cert_status ON certificates(status);

-- Competencias
CREATE TABLE competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    competency_type VARCHAR(30) NOT NULL DEFAULT 'TECHNICAL',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comp_company ON competencies(company_id);
CREATE INDEX idx_comp_type ON competencies(competency_type);

-- Niveles de competencia
CREATE TABLE competency_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    level INT NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    UNIQUE(competency_id, level)
);
CREATE INDEX idx_cl_competency ON competency_levels(competency_id);

-- Relación curso → competencia
CREATE TABLE course_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    acquired_level INT DEFAULT 1,
    UNIQUE(course_id, competency_id)
);
CREATE INDEX idx_cc_comp ON course_competencies(competency_id);

-- Competencias de empleados
CREATE TABLE employee_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    level INT NOT NULL DEFAULT 1,
    source VARCHAR(30) NOT NULL DEFAULT 'SELF',
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(employee_id, competency_id)
);
CREATE INDEX idx_ec_company ON employee_competencies(company_id);
CREATE INDEX idx_ec_employee ON employee_competencies(employee_id);
CREATE INDEX idx_ec_competency ON employee_competencies(competency_id);

-- Brechas de competencia
CREATE TABLE competency_gaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    required_level INT NOT NULL,
    current_level INT NOT NULL DEFAULT 0,
    gap INT NOT NULL DEFAULT 0,
    source VARCHAR(50),
    source_id UUID,
    status VARCHAR(30) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cg_company ON competency_gaps(company_id);
CREATE INDEX idx_cg_employee ON competency_gaps(employee_id);

-- Necesidades de capacitación
CREATE TABLE training_needs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID,
    competency_id UUID REFERENCES competencies(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    priority VARCHAR(30) DEFAULT 'MEDIUM',
    source VARCHAR(50),
    source_id UUID,
    status VARCHAR(30) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tn_company ON training_needs(company_id);
CREATE INDEX idx_tn_employee ON training_needs(employee_id);
CREATE INDEX idx_tn_status ON training_needs(status);

-- Planes de capacitación
CREATE TABLE training_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    objectives TEXT,
    period_start DATE,
    period_end DATE,
    budget_amount NUMERIC(14,2),
    budget_currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(30) DEFAULT 'DRAFT',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tp_company ON training_plans(company_id);
CREATE INDEX idx_tp_employee ON training_plans(employee_id);
CREATE INDEX idx_tp_status ON training_plans(status);

-- Cursos del plan de capacitación
CREATE TABLE training_plan_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES training_plans(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    priority VARCHAR(30) DEFAULT 'MEDIUM',
    sort_order INT DEFAULT 0,
    UNIQUE(plan_id, course_id)
);

-- Rutas de aprendizaje
CREATE TABLE learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    objectives TEXT,
    duration_days INT,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_lp_company ON learning_paths(company_id);

-- Cursos de la ruta
CREATE TABLE learning_path_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    required BOOLEAN DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    UNIQUE(path_id, course_id)
);
CREATE INDEX idx_lpc_path ON learning_path_courses(path_id);

-- Inscripciones a rutas
CREATE TABLE learning_path_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    status VARCHAR(30) DEFAULT 'ACTIVE',
    progress_percentage INT DEFAULT 0,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(path_id, employee_id)
);
CREATE INDEX idx_lpe_company ON learning_path_enrollments(company_id);
CREATE INDEX idx_lpe_employee ON learning_path_enrollments(employee_id);

-- Asistencia a sesiones
CREATE TABLE training_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    session_id UUID NOT NULL REFERENCES training_sessions(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES course_enrollments(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PRESENT',
    check_in TIMESTAMPTZ,
    check_out TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, enrollment_id)
);
CREATE INDEX idx_ta_session ON training_attendance(session_id);
CREATE INDEX idx_ta_enrollment ON training_attendance(enrollment_id);

-- Feedback de cursos
CREATE TABLE course_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES course_enrollments(id) ON DELETE CASCADE,
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    instructor_rating INT,
    content_rating INT,
    organization_rating INT,
    platform_rating INT,
    overall_rating NUMERIC(3,1),
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(enrollment_id)
);
CREATE INDEX idx_cf_enrollment ON course_feedback(enrollment_id);

-- Costos de capacitación
CREATE TABLE training_costs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    enrollment_id UUID REFERENCES course_enrollments(id),
    offering_id UUID REFERENCES course_offerings(id),
    cost_type VARCHAR(30) NOT NULL,
    description TEXT,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    incurred_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tc_company ON training_costs(company_id);
CREATE INDEX idx_tc_enrollment ON training_costs(enrollment_id);

-- Presupuestos de capacitación
CREATE TABLE training_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    year INT NOT NULL,
    department_id UUID,
    total_amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tb_company ON training_budgets(company_id);

-- Permisos de capacitación
INSERT INTO permissions (name, resource, action, description, created_at) VALUES
    ('training.read', 'training', 'read', 'View training data', NOW()),
    ('training.create', 'training', 'create', 'Create courses and content', NOW()),
    ('training.update', 'training', 'update', 'Update courses and content', NOW()),
    ('training.delete', 'training', 'delete', 'Delete courses', NOW()),
    ('training.assign', 'training', 'assign', 'Assign training to employees', NOW()),
    ('training.enroll', 'training', 'enroll', 'Enroll in courses', NOW()),
    ('training.progress', 'training', 'progress', 'Track progress', NOW()),
    ('training.assessment', 'training', 'assessment', 'Manage assessments', NOW()),
    ('training.grade', 'training', 'grade', 'Grade assessments', NOW()),
    ('training.certificate.read', 'training', 'certificate.read', 'View certificates', NOW()),
    ('training.certificate.revoke', 'training', 'certificate.revoke', 'Revoke certificates', NOW()),
    ('training.competency.read', 'training', 'competency.read', 'View competencies', NOW()),
    ('training.competency.manage', 'training', 'competency.manage', 'Manage competencies', NOW()),
    ('training.plan.read', 'training', 'plan.read', 'View training plans', NOW()),
    ('training.plan.create', 'training', 'plan.create', 'Create training plans', NOW()),
    ('training.plan.update', 'training', 'plan.update', 'Update training plans', NOW()),
    ('training.provider.manage', 'training', 'provider.manage', 'Manage providers', NOW()),
    ('training.instructor.manage', 'training', 'instructor.manage', 'Manage instructors', NOW()),
    ('training.analytics.read', 'training', 'analytics.read', 'View training analytics', NOW())
ON CONFLICT (name) DO NOTHING;
