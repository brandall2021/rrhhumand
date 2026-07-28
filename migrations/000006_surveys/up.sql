-- Surveys
CREATE TABLE surveys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    anonymous BOOLEAN NOT NULL DEFAULT false,
    multiple_responses BOOLEAN NOT NULL DEFAULT false,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_surveys_company ON surveys(company_id, status);
CREATE INDEX idx_surveys_created_by ON surveys(created_by);

-- Survey questions
CREATE TABLE survey_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    position INTEGER NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_survey_questions_survey ON survey_questions(survey_id, position);

-- Survey options
CREATE TABLE survey_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES survey_questions(id) ON DELETE CASCADE,
    option_text VARCHAR(255) NOT NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_survey_options_question ON survey_options(question_id, position);

-- Survey targets
CREATE TABLE survey_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    target_type VARCHAR(30) NOT NULL,
    target_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(survey_id, target_type, target_id)
);

CREATE INDEX idx_survey_targets_survey ON survey_targets(survey_id);

-- Survey responses (one per employee per survey)
CREATE TABLE survey_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    UNIQUE(survey_id, employee_id)
);

CREATE INDEX idx_survey_responses_survey ON survey_responses(survey_id);
CREATE INDEX idx_survey_responses_employee ON survey_responses(employee_id);

-- Survey answers
CREATE TABLE survey_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    response_id UUID NOT NULL REFERENCES survey_responses(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES survey_questions(id) ON DELETE CASCADE,
    text_value TEXT,
    number_value NUMERIC,
    option_id UUID REFERENCES survey_options(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_survey_answers_response ON survey_answers(response_id);
CREATE INDEX idx_survey_answers_question ON survey_answers(question_id);

-- Survey answer options (for MULTIPLE_CHOICE)
CREATE TABLE survey_answer_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    answer_id UUID NOT NULL REFERENCES survey_answers(id) ON DELETE CASCADE,
    option_id UUID NOT NULL REFERENCES survey_options(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(answer_id, option_id)
);

CREATE INDEX idx_survey_answer_options_answer ON survey_answer_options(answer_id);

-- Survey permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    ('surveys.read', 'surveys', 'read', 'View surveys'),
    ('surveys.create', 'surveys', 'create', 'Create surveys'),
    ('surveys.update', 'surveys', 'update', 'Update surveys'),
    ('surveys.delete', 'surveys', 'delete', 'Delete surveys'),
    ('surveys.publish', 'surveys', 'publish', 'Publish surveys'),
    ('surveys.close', 'surveys', 'close', 'Close surveys'),
    ('surveys.results', 'surveys', 'results', 'View survey results'),
    ('surveys.export', 'surveys', 'export', 'Export survey results'),
    ('surveys.respond', 'surveys', 'respond', 'Respond to surveys'),
    ('surveys.targets', 'surveys', 'targets', 'Manage survey targets');
