-- FASE 20 — Benefits & Total Rewards
-- ============================================================

-- ============================================================
-- 20A — Catálogo de Beneficios y Elegibilidad
-- ============================================================

CREATE TABLE benefit_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(7),
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_benefit_category UNIQUE (company_id, name)
);

CREATE TABLE benefit_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    category_id UUID REFERENCES benefit_categories(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    code VARCHAR(50) NOT NULL,
    nature VARCHAR(30) NOT NULL CHECK (nature IN (
        'MONETARY','NON_MONETARY','SERVICE','DISCOUNT','INSURANCE','FLEXIBLE','WELLNESS','RECOGNITION','OTHER'
    )),
    tax_treatment VARCHAR(30) NOT NULL CHECK (tax_treatment IN (
        'TAXABLE_EXEMPT','TAXABLE_UPTO','FULLY_TAXABLE','NON_TAXABLE','FRINGE_BENEFIT'
    )),
    requires_approval BOOLEAN NOT NULL DEFAULT false,
    is_reimbursable BOOLEAN NOT NULL DEFAULT false,
    is_flexible BOOLEAN NOT NULL DEFAULT false,
    has_wallet BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    config_schema JSONB DEFAULT '{}',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_benefit_type_code UNIQUE (company_id, code)
);

CREATE TABLE benefit_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    tax_id VARCHAR(50),
    provider_type VARCHAR(30) NOT NULL CHECK (provider_type IN (
        'INSURANCE','MEDICAL','DENTAL','VISION','GYM','FOOD','TRANSPORT','EDUCATION','RECREATION','TECHNOLOGY','WELLNESS','FINANCIAL','OTHER'
    )),
    contact_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    website VARCHAR(255),
    address TEXT,
    service_region VARCHAR(100),
    contract_start DATE,
    contract_end DATE,
    contract_file_path TEXT,
    billing_cycle VARCHAR(20) CHECK (billing_cycle IN ('MONTHLY','QUARTERLY','BIANNUAL','ANNUAL','PER_EMPLOYEE')),
    billing_currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    rating DECIMAL(2,1) CHECK (rating >= 0 AND rating <= 5),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_benefit_provider UNIQUE (company_id, name)
);

CREATE TABLE benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    type_id UUID NOT NULL REFERENCES benefit_types(id),
    provider_id UUID REFERENCES benefit_providers(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description VARCHAR(255),
    coverage_details TEXT,
    eligibility_summary TEXT,
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    website_url VARCHAR(500),
    terms_url VARCHAR(500),
    documentation_url VARCHAR(500),
    provider_reference VARCHAR(100),
    start_date DATE,
    end_date DATE,
    max_beneficiaries INT,
    current_beneficiaries INT NOT NULL DEFAULT 0,
    waiting_period_days INT NOT NULL DEFAULT 0,
    minimum_service_months INT DEFAULT 0,
    deductible_amount DECIMAL(14,2),
    deductible_period VARCHAR(10) CHECK (deductible_period IN ('MONTHLY','QUARTERLY','ANNUAL','ONE_TIME')),
    copay_percentage DECIMAL(5,2),
    max_coverage_amount DECIMAL(14,2),
    max_coverage_period VARCHAR(10) CHECK (max_coverage_period IN ('MONTHLY','QUARTERLY','ANNUAL','PER_EVENT','LIFETIME')),
    auto_enroll BOOLEAN NOT NULL DEFAULT false,
    enrollment_deadline_days INT,
    requires_evidence BOOLEAN NOT NULL DEFAULT false,
    evidence_description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('DRAFT','ACTIVE','INACTIVE','ARCHIVED')),
    visibility VARCHAR(20) NOT NULL DEFAULT 'ALL' CHECK (visibility IN ('ALL','MANAGERS_ONLY','HR_ONLY','ADMIN_ONLY')),
    sort_order INT NOT NULL DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_benefit_code UNIQUE (company_id, code)
);

CREATE INDEX idx_benefits_type ON benefits(type_id);
CREATE INDEX idx_benefits_provider ON benefits(provider_id);
CREATE INDEX idx_benefits_company ON benefits(company_id, status);
CREATE INDEX idx_benefits_visibility ON benefits(visibility, status);

CREATE TABLE benefit_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    plan_type VARCHAR(30) NOT NULL CHECK (plan_type IN (
        'INDIVIDUAL','FAMILY','INDIVIDUAL_PLUS','TITULAR_ONLY','TITULAR_PLUS_FAMILY','CUSTOM'
    )),
    monthly_cost_employee DECIMAL(14,2) NOT NULL DEFAULT 0,
    monthly_cost_employer DECIMAL(14,2) NOT NULL DEFAULT 0,
    annual_cost_employee DECIMAL(14,2) NOT NULL DEFAULT 0,
    annual_cost_employer DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    coverage_limit DECIMAL(14,2),
    coverage_details JSONB DEFAULT '{}',
    max_dependents INT DEFAULT 0,
    dependent_type VARCHAR(30) CHECK (dependent_type IN ('SPOUSE','CHILDREN','PARENTS','SIBLINGS','OTHER')),
    enrollment_fee DECIMAL(14,2),
    waiting_period_days INT DEFAULT 0,
    minimum_age INT,
    maximum_age INT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_benefit_plan UNIQUE (company_id, benefit_id, name)
);

CREATE INDEX idx_benefit_plans_benefit ON benefit_plans(benefit_id);

CREATE TABLE benefit_eligibility_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    rule_type VARCHAR(30) NOT NULL CHECK (rule_type IN (
        'JOB_LEVEL','DEPARTMENT','BRANCH','POSITION','EMPLOYMENT_TYPE',
        'SENIORITY','AGE','CONTRACT_TYPE','WORK_SCHEDULE','LOCATION',
        'PERFORMANCE_RATING','TENURE','SALARY_BAND','CUSTOM','GENDER'
    )),
    operator VARCHAR(10) NOT NULL CHECK (operator IN ('EQ','NEQ','GT','GTE','LT','LTE','IN','NOT_IN','BETWEEN','CONTAINS')),
    value TEXT NOT NULL,
    value_to TEXT,
    logic_group INT NOT NULL DEFAULT 0,
    logic_operator VARCHAR(5) NOT NULL DEFAULT 'AND' CHECK (logic_operator IN ('AND','OR','NOT')),
    priority INT NOT NULL DEFAULT 0,
    error_message TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_eligibility_rules_benefit ON benefit_eligibility_rules(benefit_id);

-- ============================================================
-- 20B — Workflows de Solicitudes
-- ============================================================

CREATE TABLE benefit_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID REFERENCES benefits(id),
    workflow_type VARCHAR(30) NOT NULL CHECK (workflow_type IN (
        'ENROLLMENT','CANCELLATION','CHANGE','REIMBURSEMENT','EXCEPTION','FLEXIBLE_SPENDING','GENERAL'
    )),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    requires_chain_approval BOOLEAN NOT NULL DEFAULT false,
    auto_approve BOOLEAN NOT NULL DEFAULT false,
    auto_approve_if_no_manager BOOLEAN NOT NULL DEFAULT true,
    escalation_hours INT DEFAULT 48,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefit_workflows_benefit ON benefit_workflows(benefit_id);

CREATE TABLE benefit_workflow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES benefit_workflows(id),
    step_order INT NOT NULL,
    approval_type VARCHAR(20) NOT NULL CHECK (approval_type IN (
        'MANAGER','HR','COMPENSATION','FINANCE','DIRECTOR','ADMIN','AUTO'
    )),
    approver_role_id UUID REFERENCES roles(id),
    max_rejection_count INT NOT NULL DEFAULT 3,
    is_required BOOLEAN NOT NULL DEFAULT true,
    notification_template TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_workflow_step UNIQUE (workflow_id, step_order)
);

-- ============================================================
-- 20B — Asignaciones y Solicitudes
-- ============================================================

CREATE TABLE employee_benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    plan_id UUID REFERENCES benefit_plans(id),
    provider_id UUID REFERENCES benefit_providers(id),
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN (
        'PENDING','ACTIVE','SUSPENDED','CANCELLED','EXPIRED','WAITING_PERIOD','DECLINED'
    )),
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    start_date DATE,
    end_date DATE,
    cancellation_date DATE,
    cancellation_reason TEXT,
    auto_renew BOOLEAN NOT NULL DEFAULT true,
    renewal_date DATE,
    coverage_level VARCHAR(30) DEFAULT 'INDIVIDUAL' CHECK (coverage_level IN (
        'INDIVIDUAL','FAMILY','INDIVIDUAL_PLUS','TITULAR','TITULAR_PLUS_FAMILY','CUSTOM'
    )),
    dependents JSONB DEFAULT '[]',
    emergency_contact JSONB DEFAULT '{}',
    beneficiary_info JSONB DEFAULT '{}',
    employee_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payroll_deduction_enabled BOOLEAN NOT NULL DEFAULT true,
    payroll_deduction_amount DECIMAL(14,2) DEFAULT 0,
    external_member_id VARCHAR(100),
    external_policy_number VARCHAR(100),
    external_group_number VARCHAR(100),
    documents JSONB DEFAULT '[]',
    notes TEXT,
    source VARCHAR(20) NOT NULL DEFAULT 'ADMIN' CHECK (source IN ('ADMIN','SELF_SERVICE','ONBOARDING','MIGRATION','API')),
    enrolled_by UUID NOT NULL,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emp_benefits_employee ON employee_benefits(employee_id, company_id);
CREATE INDEX idx_emp_benefits_benefit ON employee_benefits(benefit_id);
CREATE INDEX idx_emp_benefits_status ON employee_benefits(status);
CREATE UNIQUE INDEX uq_emp_benefit_active ON employee_benefits(employee_id, benefit_id) WHERE status IN ('ACTIVE','PENDING','WAITING_PERIOD');

CREATE TABLE employee_benefit_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_benefit_id UUID NOT NULL REFERENCES employee_benefits(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    action VARCHAR(30) NOT NULL CHECK (action IN (
        'ENROLLED','ACTIVATED','CHANGED_PLAN','CHANGED_COVERAGE','SUSPENDED','REINSTATED',
        'CANCELLED','EXPIRED','UPDATED_COST','UPDATED_DEPENDENTS','UPDATED_BENEFICIARY',
        'DECLINED','AUTO_RENEWED','PAYROLL_UPDATED','EXCEPTION_GRANTED'
    )),
    previous_value JSONB,
    new_value JSONB,
    change_reason TEXT,
    changed_by UUID NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emp_benefit_history_benefit ON employee_benefit_history(employee_benefit_id);
CREATE INDEX idx_emp_benefit_history_employee ON employee_benefit_history(employee_id);

CREATE TABLE benefit_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    plan_id UUID REFERENCES benefit_plans(id),
    employee_benefit_id UUID REFERENCES employee_benefits(id),
    request_type VARCHAR(30) NOT NULL CHECK (request_type IN (
        'ENROLLMENT','CANCELLATION','CHANGE_PLAN','CHANGE_COVERAGE','EXCEPTION','REIMBURSEMENT','FLEXIBLE_SPENDING','OTHER'
    )),
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT','SUBMITTED','PENDING_APPROVAL','APPROVED','REJECTED','CANCELLED','WAITING_INFO','ESCALATED'
    )),
    request_data JSONB DEFAULT '{}',
    justification TEXT,
    coverage_level VARCHAR(30),
    dependents JSONB DEFAULT '[]',
    effective_date DATE,
    notes TEXT,
    submitted_by UUID,
    submitted_at TIMESTAMPTZ,
    resolved_by UUID,
    resolved_at TIMESTAMPTZ,
    resolution_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefit_requests_employee ON benefit_requests(employee_id, company_id);
CREATE INDEX idx_benefit_requests_benefit ON benefit_requests(benefit_id);
CREATE INDEX idx_benefit_requests_status ON benefit_requests(status);

CREATE TABLE benefit_request_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES benefit_requests(id),
    step_id UUID REFERENCES benefit_workflow_steps(id),
    reviewer_id UUID NOT NULL,
    review_type VARCHAR(20) NOT NULL CHECK (review_type IN ('APPROVAL','REJECTION','ESCALATION','INFO_REQUEST','AUTO_APPROVED')),
    comment TEXT,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefit_reviews_request ON benefit_request_reviews(request_id);

-- ============================================================
-- 20C — Costos, Wallet Flexible y Reembolsos
-- ============================================================

CREATE TABLE benefit_costs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    plan_id UUID REFERENCES benefit_plans(id),
    cost_type VARCHAR(30) NOT NULL CHECK (cost_type IN (
        'MONTHLY_PREMIUM','ANNUAL_PREMIUM','ENROLLMENT_FEE','ADMIN_FEE','PER_EMPLOYEE','PER_DEPENDENT','DISCOUNT_PERCENTAGE','FIXED_AMOUNT','VARIABLE'
    )),
    employee_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    total_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('MONTHLY','QUARTERLY','BIANNUAL','ANNUAL','ONE_TIME','PER_USAGE')),
    billing_cycle_day INT,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefit_costs_benefit ON benefit_costs(benefit_id);

CREATE TABLE benefit_cost_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID NOT NULL REFERENCES benefits(id),
    cost_id UUID REFERENCES benefit_costs(id),
    schedule_date DATE NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','PAID','CANCELLED','WAIVED')),
    paid_at TIMESTAMPTZ,
    payment_reference VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cost_schedules_benefit ON benefit_cost_schedules(benefit_id);
CREATE INDEX idx_cost_schedules_date ON benefit_cost_schedules(schedule_date);

CREATE TABLE benefit_flexible_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    plan_type VARCHAR(30) NOT NULL CHECK (plan_type IN (
        'MEAL','TRANSPORT','EDUCATION','HEALTH','WELLNESS','CHILD_CARE','OLDER_CARE','COMMUNICATION','TECHNOLOGY','CUSTOM'
    )),
    annual_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    monthly_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    employer_contribution DECIMAL(14,2) NOT NULL DEFAULT 0,
    employee_contribution DECIMAL(14,2) NOT NULL DEFAULT 0,
    contribution_frequency VARCHAR(10) CHECK (contribution_frequency IN ('MONTHLY','QUARTERLY','ANNUAL')),
    max_rollover_amount DECIMAL(14,2) DEFAULT 0,
    rollover_expiry_months INT DEFAULT 12,
    allow_reimbursement BOOLEAN NOT NULL DEFAULT true,
    allow_prepaid_card BOOLEAN NOT NULL DEFAULT false,
    require_receipts BOOLEAN NOT NULL DEFAULT true,
    receipt_min_amount DECIMAL(14,2) DEFAULT 0,
    eligible_categories JSONB DEFAULT '[]',
    tax_exempt BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_flexible_plan UNIQUE (company_id, name)
);

CREATE TABLE benefit_flexible_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    flexible_plan_id UUID NOT NULL REFERENCES benefit_flexible_plans(id),
    fiscal_year INT NOT NULL,
    total_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    used_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    pending_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    rolled_over_from_previous DECIMAL(14,2) DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','EXHAUSTED','EXPIRED','CLOSED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_flexible_budget UNIQUE (employee_id, flexible_plan_id, fiscal_year)
);

CREATE TABLE employee_benefit_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    benefit_id UUID REFERENCES benefits(id),
    wallet_type VARCHAR(30) NOT NULL CHECK (wallet_type IN (
        'FLEXIBLE_SPENDING','MEAL','TRANSPORT','EDUCATION','HEALTH','WELLNESS','GENERAL'
    )),
    balance DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    last_transaction_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefit_wallets_employee ON employee_benefit_wallets(employee_id, company_id);

CREATE TABLE benefit_wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES employee_benefit_wallets(id),
    transaction_type VARCHAR(30) NOT NULL CHECK (transaction_type IN (
        'CREDIT_EMPLOYER','CREDIT_EMPLOYEE','CREDIT_ROLLOVER','CREDIT_ADJUSTMENT',
        'DEBIT_REIMBURSEMENT','DEBIT_PURCHASE','DEBIT_CARD','DEBIT_ADJUSTMENT','DEBIT_EXPIRY','DEBIT_TRANSFER'
    )),
    amount DECIMAL(14,2) NOT NULL,
    balance_before DECIMAL(14,2) NOT NULL,
    balance_after DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    reference_type VARCHAR(50),
    reference_id UUID,
    description TEXT,
    receipt_url VARCHAR(500),
    transaction_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wallet_tx_wallet ON benefit_wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_tx_date ON benefit_wallet_transactions(transaction_date);

CREATE TABLE benefit_reimbursements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    benefit_id UUID REFERENCES benefits(id),
    flexible_plan_id UUID REFERENCES benefit_flexible_plans(id),
    wallet_id UUID REFERENCES employee_benefit_wallets(id),
    request_id UUID REFERENCES benefit_requests(id),
    category VARCHAR(30) NOT NULL CHECK (category IN (
        'MEDICAL','DENTAL','VISION','EDUCATION','GYM','WELLNESS','TRANSPORT','MEAL','CHILD_CARE','OLDER_CARE','COMMUNICATION','TECHNOLOGY','OTHER'
    )),
    description TEXT NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    approved_amount DECIMAL(14,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    receipt_date DATE NOT NULL,
    expense_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT','SUBMITTED','PENDING_REVIEW','APPROVED','REJECTED','PAID','CANCELLED'
    )),
    rejection_reason TEXT,
    payment_method VARCHAR(20) CHECK (payment_method IN ('BANK_TRANSFER','PAYROLL','WALLET','CHECK','CASH')),
    paid_at TIMESTAMPTZ,
    payment_reference VARCHAR(100),
    submitted_by UUID,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reimbursements_employee ON benefit_reimbursements(employee_id, company_id);
CREATE INDEX idx_reimbursements_status ON benefit_reimbursements(status);

CREATE TABLE benefit_reimbursement_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reimbursement_id UUID NOT NULL REFERENCES benefit_reimbursements(id),
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INT NOT NULL,
    storage_path TEXT NOT NULL,
    ocr_text TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    uploaded_by UUID NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reimb_docs_reimbursement ON benefit_reimbursement_documents(reimbursement_id);

-- ============================================================
-- 20D — Bonos, Incentivos e Integración Payroll
-- ============================================================

CREATE TABLE employee_bonuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    bonus_type VARCHAR(30) NOT NULL CHECK (bonus_type IN (
        'PERFORMANCE','SIGNING','RETENTION','REFERRAL','PROJECT','PROFIT_SHARING','COMMISSION','SPOT','HOLIDAY','PRODUCTIVITY','OTHER'
    )),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('LUMP_SUM','INSTALLMENTS','PAYROLL_INCLUSION')),
    installment_count INT DEFAULT 1,
    installment_amount DECIMAL(14,2),
    frequency VARCHAR(10) CHECK (frequency IN ('MONTHLY','QUARTERLY','ANNUAL','ONE_TIME')),
    grant_date DATE NOT NULL DEFAULT CURRENT_DATE,
    vesting_start DATE,
    vesting_end DATE,
    payment_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','APPROVED','PAID','CANCELLED','FORFEITED','CLAWBACK'
    )),
    clawback_amount DECIMAL(14,2) DEFAULT 0,
    clawback_reason TEXT,
    performance_period_start DATE,
    performance_period_end DATE,
    performance_score DECIMAL(5,2),
    is_taxable BOOLEAN NOT NULL DEFAULT true,
    tax_withholding DECIMAL(14,2) DEFAULT 0,
    net_amount DECIMAL(14,2),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    paid_in_payroll BOOLEAN NOT NULL DEFAULT false,
    payroll_run_id UUID REFERENCES payroll_runs(id),
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emp_bonuses_employee ON employee_bonuses(employee_id, company_id);
CREATE INDEX idx_emp_bonuses_status ON employee_bonuses(status);

CREATE TABLE employee_incentives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    incentive_type VARCHAR(30) NOT NULL CHECK (incentive_type IN (
        'SPOT_AWARD','RECOGNITION','GIFT_CARD','MERCHANDISE','EXPERIENCE','EXTRA_VACATION','FLEXIBLE_HOURS','TRAINING','CERTIFICATION','OTHER'
    )),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    award_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date DATE,
    redemption_date DATE,
    status VARCHAR(20) NOT NULL CHECK (status IN (
        'AWARDED','ACCEPTED','REDEEMED','EXPIRED','CANCELLED','DECLINED'
    )),
    points_cost INT DEFAULT 0,
    is_taxable BOOLEAN NOT NULL DEFAULT false,
    awarded_by UUID,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emp_incentives_employee ON employee_incentives(employee_id, company_id);

CREATE TABLE benefit_payroll_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    benefit_id UUID REFERENCES benefits(id),
    flexible_plan_id UUID REFERENCES benefit_flexible_plans(id),
    employee_benefit_id UUID REFERENCES employee_benefits(id),
    bonus_id UUID REFERENCES employee_bonuses(id),
    mapping_type VARCHAR(30) NOT NULL CHECK (mapping_type IN (
        'BENEFIT_DEDUCTION','BENEFIT_CONTRIBUTION','FLEXIBLE_CREDIT','BONUS_PAYMENT','INCENTIVE','REIMBURSEMENT'
    )),
    payroll_concept_id UUID REFERENCES payroll_concepts(id),
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('MONTHLY','ONE_TIME','QUARTERLY','ANNUAL')),
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_synced_at TIMESTAMPTZ,
    sync_status VARCHAR(20) DEFAULT 'PENDING' CHECK (sync_status IN ('PENDING','SYNCED','ERROR','SKIPPED')),
    sync_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payroll_mappings_benefit ON benefit_payroll_mappings(benefit_id);
CREATE INDEX idx_payroll_mappings_concept ON benefit_payroll_mappings(payroll_concept_id);

-- ============================================================
-- 20E — Total Rewards, Reportes y Notificaciones
-- ============================================================

CREATE TABLE total_rewards_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (category IN (
        'BASE_SALARY','VARIABLE_PAY','BENEFITS','BONUSES','INCENTIVES','EQUITY','RETIREMENT','FLEXIBLE_SPENDING',
        'INSURANCE','PERKS','RECOGNITION','DEVELOPMENT','WELLNESS','OTHER'
    )),
    description TEXT,
    amount_type VARCHAR(20) NOT NULL CHECK (amount_type IN ('FIXED','PERCENTAGE','RANGE','FORMULA')),
    amount_value DECIMAL(14,2) DEFAULT 0,
    amount_percentage DECIMAL(5,2),
    currency VARCHAR(3) DEFAULT 'ARS',
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('MONTHLY','ANNUAL','ONE_TIME','QUARTERLY')),
    display_order INT NOT NULL DEFAULT 0,
    is_monetary BOOLEAN NOT NULL DEFAULT true,
    include_in_statement BOOLEAN NOT NULL DEFAULT true,
    icon VARCHAR(50),
    color VARCHAR(7),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE total_rewards_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    snapshot_date DATE NOT NULL DEFAULT CURRENT_DATE,
    fiscal_year INT NOT NULL,
    period_name VARCHAR(100),
    base_salary DECIMAL(14,2) NOT NULL DEFAULT 0,
    variable_pay DECIMAL(14,2) NOT NULL DEFAULT 0,
    bonuses_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    incentives_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    benefits_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_contributions DECIMAL(14,2) NOT NULL DEFAULT 0,
    flexible_spending DECIMAL(14,2) NOT NULL DEFAULT 0,
    insurance_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    development_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    wellness_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    recognition_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    perks_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    total_rewards DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    items JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rewards_snapshot_employee ON total_rewards_snapshots(employee_id, company_id);
CREATE INDEX idx_rewards_snapshot_period ON total_rewards_snapshots(company_id, fiscal_year);

CREATE TABLE benefit_notification_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID REFERENCES employees(id),
    notification_type VARCHAR(30) NOT NULL CHECK (notification_type IN (
        'ENROLLMENT_OPEN','ENROLLMENT_REMINDER','ENROLLMENT_CONFIRMATION','ENROLLMENT_DEADLINE',
        'APPROVAL_REQUEST','APPROVAL_REMINDER','REQUEST_APPROVED','REQUEST_REJECTED',
        'BENEFIT_ACTIVATED','BENEFIT_SUSPENDED','BENEFIT_CANCELLED','BENEFIT_EXPIRING',
        'REIMBURSEMENT_APPROVED','REIMBURSEMENT_PAID','REIMBURSEMENT_REJECTED',
        'COST_CHANGE','PLAN_CHANGE','PROVIDER_CHANGE','FLEXIBLE_BUDGET_LOW',
        'WALLET_TRANSACTION','BONUS_AWARDED','INCENTIVE_AWARDED','STATEMENT_AVAILABLE',
        'DEADLINE_REMINDER','GENERAL'
    )),
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('EMAIL','PUSH','IN_APP','SMS','DASHBOARD')),
    title VARCHAR(255) NOT NULL,
    body TEXT,
    metadata JSONB DEFAULT '{}',
    read_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_log_employee ON benefit_notification_log(employee_id);
CREATE INDEX idx_notification_log_type ON benefit_notification_log(notification_type, sent_at);

CREATE TABLE benefit_report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    report_type VARCHAR(30) NOT NULL CHECK (report_type IN (
        'COST_ANALYSIS','UTILIZATION','ENROLLMENT','ELIGIBILITY','PROVIDER_PERFORMANCE',
        'TOTAL_REWARDS','BUDGET_VS_ACTUAL','FLEXIBLE_SPENDING','BONUS_TRACKING','INCENTIVE_ANALYSIS',
        'EMPLOYEE_SATISFACTION','ROI_ANALYSIS','COMPLIANCE','CUSTOM'
    )),
    config JSONB DEFAULT '{}',
    schedule_cron VARCHAR(100),
    recipients JSONB DEFAULT '[]',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE benefit_report_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    definition_id UUID REFERENCES benefit_report_definitions(id),
    company_id UUID NOT NULL REFERENCES companies(id),
    report_type VARCHAR(30) NOT NULL,
    period_start DATE,
    period_end DATE,
    file_name VARCHAR(255),
    file_content TEXT,
    storage_path TEXT,
    file_size INT,
    format VARCHAR(10) NOT NULL DEFAULT 'JSON' CHECK (format IN ('JSON','CSV','PDF','XLS','HTML')),
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATED' CHECK (status IN ('GENERATED','ERROR')),
    error_message TEXT,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- RBAC Permissions
-- ============================================================
INSERT INTO role_permissions (name, description, module) VALUES
    -- Benefits catalog
    ('benefits.category.manage', 'Manage benefit categories', 'benefits'),
    ('benefits.type.manage', 'Manage benefit types', 'benefits'),
    ('benefits.benefit.create', 'Create benefits', 'benefits'),
    ('benefits.benefit.view', 'View benefits', 'benefits'),
    ('benefits.benefit.edit', 'Edit benefits', 'benefits'),
    ('benefits.benefit.delete', 'Delete benefits', 'benefits'),
    ('benefits.provider.manage', 'Manage benefit providers', 'benefits'),
    ('benefits.plan.manage', 'Manage benefit plans', 'benefits'),
    ('benefits.eligibility.manage', 'Manage eligibility rules', 'benefits'),
    ('benefits.catalog.view', 'View benefit catalog', 'benefits'),
    -- Employee assignment
    ('benefits.employee.assign', 'Assign benefits to employees', 'benefits'),
    ('benefits.employee.view', 'View employee benefits', 'benefits'),
    ('benefits.employee.modify', 'Modify employee benefits', 'benefits'),
    ('benefits.employee.remove', 'Remove employee benefits', 'benefits'),
    ('benefits.employee.history', 'View employee benefit history', 'benefits'),
    -- Requests & workflows
    ('benefits.request.create', 'Create benefit requests', 'benefits'),
    ('benefits.request.view', 'View benefit requests', 'benefits'),
    ('benefits.request.approve', 'Approve benefit requests', 'benefits'),
    ('benefits.request.reject', 'Reject benefit requests', 'benefits'),
    ('benefits.workflow.manage', 'Manage approval workflows', 'benefits'),
    -- Wallet & reimbursements
    ('benefits.wallet.view', 'View benefit wallets', 'benefits'),
    ('benefits.reimbursement.create', 'Create reimbursements', 'benefits'),
    ('benefits.reimbursement.approve', 'Approve reimbursements', 'benefits'),
    ('benefits.reimbursement.view', 'View reimbursements', 'benefits'),
    ('benefits.flexible.manage', 'Manage flexible plans', 'benefits'),
    -- Costs & billing
    ('benefits.cost.manage', 'Manage benefit costs', 'benefits'),
    ('benefits.cost.view', 'View benefit costs', 'benefits'),
    -- Bonuses & incentives
    ('benefits.bonus.create', 'Create bonuses', 'benefits'),
    ('benefits.bonus.view', 'View bonuses', 'benefits'),
    ('benefits.bonus.approve', 'Approve bonuses', 'benefits'),
    ('benefits.incentive.create', 'Create incentives', 'benefits'),
    ('benefits.incentive.view', 'View incentives', 'benefits'),
    -- Payroll integration
    ('benefits.payroll.sync', 'Sync benefits with payroll', 'benefits'),
    ('benefits.payroll.view', 'View payroll mappings', 'benefits'),
    -- Total rewards
    ('benefits.rewards.view', 'View total rewards statements', 'benefits'),
    ('benefits.rewards.generate', 'Generate total rewards statements', 'benefits'),
    ('benefits.rewards.manage', 'Manage rewards items', 'benefits'),
    -- Reports
    ('benefits.report.generate', 'Generate benefits reports', 'benefits'),
    ('benefits.report.view', 'View benefits reports', 'benefits'),
    -- Employee self-service
    ('benefits.me.view', 'View own benefits', 'benefits'),
    ('benefits.me.enroll', 'Self-enroll in benefits', 'benefits'),
    ('benefits.me.request', 'Create benefit requests', 'benefits'),
    ('benefits.me.reimbursement', 'Submit reimbursements', 'benefits')
ON CONFLICT (name) DO NOTHING;
