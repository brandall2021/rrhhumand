-- FASE 21 — Gestión de Gastos y Viáticos
-- ============================================================

-- ============================================================
-- 1. CATÁLOGO DE GASTOS
-- ============================================================

CREATE TABLE expense_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES expense_categories(id),
    requires_receipt BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_expense_category_code UNIQUE (company_id, code)
);

-- 2. MÉTODOS DE PAGO

CREATE TABLE expense_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    is_corporate BOOLEAN NOT NULL DEFAULT false,
    requires_receipt BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_expense_payment_method UNIQUE (company_id, code)
);

-- 3. TIPOS DE CAMBIO

CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(18,8) NOT NULL,
    effective_date DATE NOT NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'MANUAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exchange_rates_pair ON exchange_rates(company_id, from_currency, to_currency, effective_date DESC);

-- 4. GASTOS

CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    travel_id UUID,
    expense_report_id UUID,
    category_id UUID NOT NULL REFERENCES expense_categories(id),
    expense_date DATE NOT NULL,
    description TEXT NOT NULL,
    original_amount DECIMAL(18,2) NOT NULL,
    original_currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    exchange_rate DECIMAL(18,8) DEFAULT 1,
    base_amount DECIMAL(18,2) NOT NULL,
    base_currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_method_id UUID REFERENCES expense_payment_methods(id),
    payment_method_other VARCHAR(50),
    merchant_name VARCHAR(200),
    merchant_tax_id VARCHAR(50),
    receipt_number VARCHAR(100),
    is_reimbursable BOOLEAN NOT NULL DEFAULT true,
    is_policy_compliant BOOLEAN NOT NULL DEFAULT true,
    policy_status VARCHAR(30) NOT NULL DEFAULT 'PENDING' CHECK (policy_status IN ('PENDING','COMPLIANT','NON_COMPLIANT','EXCEPTION','OVERRIDDEN')),
    policy_override_reason TEXT,
    policy_override_by UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT','SUBMITTED','UNDER_REVIEW','APPROVED','REJECTED','OBSERVED','PAID','CANCELLED'
    )),
    rejection_reason TEXT,
    observation TEXT,
    cost_center_id UUID,
    project_id UUID,
    receipt_missing_reason TEXT,
    idempotency_key VARCHAR(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expenses_employee ON expenses(employee_id, company_id);
CREATE INDEX idx_expenses_travel ON expenses(travel_id);
CREATE INDEX idx_expenses_report ON expenses(expense_report_id);
CREATE INDEX idx_expenses_category ON expenses(category_id);
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_expenses_date ON expenses(company_id, expense_date);
CREATE UNIQUE INDEX uq_expense_idempotency ON expenses(company_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

-- 5. VIAJES

CREATE TABLE travels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    title VARCHAR(200) NOT NULL,
    purpose TEXT,
    origin VARCHAR(150) NOT NULL,
    destination VARCHAR(150) NOT NULL,
    departure_date DATE NOT NULL,
    return_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT','REQUESTED','APPROVED','REJECTED','IN_PROGRESS','COMPLETED','CANCELLED'
    )),
    estimated_budget DECIMAL(18,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    cost_center_id UUID,
    project_id UUID,
    rejection_reason TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_travels_employee ON travels(employee_id, company_id);
CREATE INDEX idx_travels_status ON travels(status);

-- 6. PARTICIPANTES DE VIAJE

CREATE TABLE travel_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    travel_id UUID NOT NULL REFERENCES travels(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    role VARCHAR(50) NOT NULL DEFAULT 'PARTICIPANT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_travel_participant UNIQUE (travel_id, employee_id)
);

-- 7. RENDICIONES

CREATE TABLE expense_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    travel_id UUID REFERENCES travels(id),
    advance_id UUID,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    total_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    advance_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    reimbursable_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    employee_refund_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    company_owes_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT','SUBMITTED','UNDER_REVIEW','APPROVED','REJECTED','OBSERVED','PAID','CANCELLED'
    )),
    submitted_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    rejection_reason TEXT,
    observation TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_reports_employee ON expense_reports(employee_id, company_id);
CREATE INDEX idx_expense_reports_travel ON expense_reports(travel_id);
CREATE INDEX idx_expense_reports_status ON expense_reports(status);

-- 8. ANTICIPOS

CREATE TABLE expense_advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    travel_id UUID REFERENCES travels(id),
    requested_amount DECIMAL(18,2) NOT NULL,
    approved_amount DECIMAL(18,2),
    settled_amount DECIMAL(18,2) DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    request_date DATE NOT NULL DEFAULT CURRENT_DATE,
    approved_date TIMESTAMPTZ,
    paid_date TIMESTAMPTZ,
    settled_date TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'REQUESTED' CHECK (status IN (
        'REQUESTED','APPROVED','REJECTED','PAID','PARTIALLY_SETTLED','SETTLED','CANCELLED'
    )),
    rejection_reason TEXT,
    idempotency_key VARCHAR(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_advances_employee ON expense_advances(employee_id, company_id);
CREATE INDEX idx_expense_advances_travel ON expense_advances(travel_id);
CREATE INDEX idx_expense_advances_status ON expense_advances(status);

-- 9. REEMBOLSOS

CREATE TABLE expense_reimbursements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    expense_report_id UUID REFERENCES expense_reports(id),
    advance_id UUID REFERENCES expense_advances(id),
    amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_method VARCHAR(30) CHECK (payment_method IN ('BANK_TRANSFER','PAYROLL','CHECK','CASH','WALLET','OTHER')),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','APPROVED','PROCESSING','PAID','REJECTED','CANCELLED'
    )),
    payroll_id UUID,
    payroll_run_id UUID REFERENCES payroll_runs(id),
    paid_at TIMESTAMPTZ,
    rejection_reason TEXT,
    idempotency_key VARCHAR(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_reimbursements_employee ON expense_reimbursements(employee_id, company_id);
CREATE INDEX idx_expense_reimbursements_report ON expense_reimbursements(expense_report_id);

-- 10. COMPROBANTES

CREATE TABLE expense_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    expense_id UUID NOT NULL REFERENCES expenses(id),
    storage_key VARCHAR(500) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    hash VARCHAR(128),
    ocr_text TEXT,
    ocr_processed BOOLEAN NOT NULL DEFAULT false,
    ocr_data JSONB DEFAULT '{}',
    uploaded_by UUID NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_receipts_expense ON expense_receipts(expense_id);

-- 11. DETECCIÓN DE DUPLICADOS

CREATE TABLE expense_duplicate_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    expense_id UUID NOT NULL REFERENCES expenses(id),
    duplicate_expense_id UUID REFERENCES expenses(id),
    match_reason VARCHAR(100) NOT NULL,
    match_score DECIMAL(5,2),
    status VARCHAR(30) NOT NULL DEFAULT 'POTENTIAL_DUPLICATE' CHECK (status IN ('POTENTIAL_DUPLICATE','CONFIRMED_DUPLICATE','FALSE_POSITIVE')),
    resolved_by UUID,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 12. POLÍTICAS

CREATE TABLE expense_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE expense_policy_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES expense_policies(id),
    category_id UUID REFERENCES expense_categories(id),
    employee_category VARCHAR(50),
    max_amount DECIMAL(18,2),
    currency VARCHAR(3),
    requires_receipt BOOLEAN NOT NULL DEFAULT true,
    requires_approval BOOLEAN NOT NULL DEFAULT false,
    allowed_payment_methods JSONB DEFAULT '[]',
    daily_allowance_category VARCHAR(50),
    conditions JSONB DEFAULT '{}',
    priority INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 13. WORKFLOWS

CREATE TABLE expense_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    workflow_type VARCHAR(30) NOT NULL CHECK (workflow_type IN (
        'EXPENSE','TRAVEL','ADVANCE','REIMBURSEMENT','GENERAL'
    )),
    min_amount DECIMAL(18,2) DEFAULT 0,
    max_amount DECIMAL(18,2),
    requires_chain_approval BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE expense_workflow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES expense_workflows(id),
    step_order INT NOT NULL,
    approval_type VARCHAR(30) NOT NULL CHECK (approval_type IN (
        'MANAGER','HR','FINANCE','DIRECTOR','ADMIN','AUTO'
    )),
    approver_role_id UUID REFERENCES roles(id),
    max_rejection_count INT NOT NULL DEFAULT 3,
    is_required BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_expense_workflow_step UNIQUE (workflow_id, step_order)
);

-- 14. APROBACIONES

CREATE TABLE expense_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    entity_type VARCHAR(30) NOT NULL CHECK (entity_type IN ('EXPENSE','TRAVEL','EXPENSE_REPORT','ADVANCE','REIMBURSEMENT')),
    entity_id UUID NOT NULL,
    step_id UUID REFERENCES expense_workflow_steps(id),
    approver_id UUID NOT NULL,
    step_order INT NOT NULL DEFAULT 1,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','APPROVED','REJECTED','SKIPPED','OBSERVED')),
    comment TEXT,
    approved_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    observed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_approvals_entity ON expense_approvals(entity_type, entity_id);
CREATE INDEX idx_expense_approvals_approver ON expense_approvals(approver_id, status);

-- 15. VIÁTICOS DIARIOS

CREATE TABLE daily_allowance_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(150) NOT NULL,
    country VARCHAR(3),
    region VARCHAR(100),
    city VARCHAR(150),
    employee_category VARCHAR(50),
    daily_amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    meal_percentage DECIMAL(5,2),
    lodging_percentage DECIMAL(5,2),
    transport_percentage DECIMAL(5,2),
    other_percentage DECIMAL(5,2),
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 16. PRESUPUESTOS

CREATE TABLE expense_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(150) NOT NULL,
    fiscal_year INT NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_amount DECIMAL(18,2) NOT NULL,
    used_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    reserved_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    cost_center_id UUID,
    project_id UUID,
    category_id UUID REFERENCES expense_categories(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 17. TARJETAS CORPORATIVAS

CREATE TABLE corporate_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID REFERENCES employees(id),
    card_number_masked VARCHAR(30) NOT NULL,
    cardholder_name VARCHAR(150) NOT NULL,
    provider VARCHAR(150),
    credit_limit DECIMAL(18,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    expiration_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE corporate_card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES corporate_cards(id),
    company_id UUID NOT NULL REFERENCES companies(id),
    expense_id UUID REFERENCES expenses(id),
    transaction_date TIMESTAMPTZ NOT NULL,
    merchant_name VARCHAR(200),
    amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    reference VARCHAR(100),
    status VARCHAR(30) NOT NULL DEFAULT 'IMPORTED' CHECK (status IN ('IMPORTED','MATCHED','UNMATCHED','CLASSIFIED','DISPUTED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 18. AUDITORÍA

CREATE TABLE expense_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    user_id UUID NOT NULL,
    employee_id UUID,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(30) NOT NULL,
    entity_id UUID NOT NULL,
    old_values JSONB DEFAULT '{}',
    new_values JSONB DEFAULT '{}',
    ip VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expense_audit_entity ON expense_audit_logs(entity_type, entity_id);
CREATE INDEX idx_expense_audit_company ON expense_audit_logs(company_id, created_at DESC);

-- 19. LOG DE INTEGRACIÓN

CREATE TABLE expense_integration_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    integration_type VARCHAR(30) NOT NULL CHECK (integration_type IN (
        'PAYROLL','ACCOUNTING','BENEFITS','OCR','STORAGE','NOTIFICATION','OTHER'
    )),
    entity_type VARCHAR(30),
    entity_id UUID,
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('OUTBOUND','INBOUND')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('SUCCESS','ERROR','PENDING','RETRY')),
    request_body TEXT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 20. NOTIFICACIONES

CREATE TABLE expense_notification_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID,
    notification_type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('IN_APP','EMAIL','PUSH','WHATSAPP')),
    title VARCHAR(255) NOT NULL,
    body TEXT,
    metadata JSONB DEFAULT '{}',
    read_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- RBAC Permissions
-- ============================================================
INSERT INTO role_permissions (name, description, module) VALUES
    -- Read
    ('expenses.read', 'View expenses', 'expenses'),
    ('expenses.read.own', 'View own expenses', 'expenses'),
    ('expenses.create', 'Create expenses', 'expenses'),
    ('expenses.update', 'Update expenses', 'expenses'),
    ('expenses.delete', 'Delete expenses', 'expenses'),
    ('expenses.submit', 'Submit expenses', 'expenses'),
    ('expenses.approve', 'Approve expenses', 'expenses'),
    ('expenses.reject', 'Reject expenses', 'expenses'),
    ('expenses.observe', 'Observe expenses', 'expenses'),
    ('expenses.reimburse', 'Process reimbursements', 'expenses'),
    ('expenses.report', 'Generate expense reports', 'expenses'),
    -- Travel
    ('travels.read', 'View travels', 'expenses'),
    ('travels.create', 'Create travels', 'expenses'),
    ('travels.approve', 'Approve travels', 'expenses'),
    ('travels.reject', 'Reject travels', 'expenses'),
    -- Advances
    ('advances.read', 'View advances', 'expenses'),
    ('advances.create', 'Request advances', 'expenses'),
    ('advances.approve', 'Approve advances', 'expenses'),
    ('advances.pay', 'Pay advances', 'expenses'),
    ('advances.settle', 'Settle advances', 'expenses'),
    -- Policies
    ('policies.read', 'View expense policies', 'expenses'),
    ('policies.create', 'Create expense policies', 'expenses'),
    ('policies.update', 'Update expense policies', 'expenses'),
    ('policies.delete', 'Delete expense policies', 'expenses'),
    -- Budget
    ('budgets.read', 'View budgets', 'expenses'),
    ('budgets.manage', 'Manage budgets', 'expenses'),
    ('expenses.admin', 'Full expense administration', 'expenses')
ON CONFLICT (name) DO NOTHING;
