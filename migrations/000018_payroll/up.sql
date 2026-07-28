-- FASE 19 — Payroll / Liquidación de Haberes (19A: Core Engine)
-- Principio: tasas, topes, mínimos y reglas NUNCA hardcodeadas.
-- Motor parametrizable con reglas versionadas por vigencia.

-- ============================================================
-- 1. PERIODOS DE LIQUIDACIÓN
-- ============================================================
CREATE TABLE payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    year INT NOT NULL CHECK (year >= 2020),
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN (
        'MONTHLY','BIWEEKLY','WEEKLY','SPECIAL','BONUS','SAC','VACATION','FINAL','ADJUSTMENT'
    )),
    name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN (
        'OPEN','PROCESSING','CALCULATED','VALIDATED','APPROVED','CLOSED','CANCELLED'
    )),
    closed_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_payroll_period_company UNIQUE (company_id, year, month, period_type)
);

CREATE INDEX idx_payroll_periods_company ON payroll_periods(company_id, year, month);

-- ============================================================
-- 2. CORRIDAS DE LIQUIDACIÓN
-- ============================================================
CREATE TABLE payroll_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    period_id UUID NOT NULL REFERENCES payroll_periods(id),
    run_number INT NOT NULL,
    run_type VARCHAR(20) NOT NULL DEFAULT 'ORDINARY' CHECK (run_type IN (
        'ORDINARY','ADJUSTMENT','COMPLEMENTARY','FINAL','SAC','VACATION'
    )),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN (
        'OPEN','LOADING','CALCULATING','CALCULATED','VALIDATING','VALIDATED','APPROVED','CLOSED','CANCELLED'
    )),
    engine_version VARCHAR(20) NOT NULL DEFAULT '1.0.0',
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    closed_by UUID,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payroll_runs_period ON payroll_runs(company_id, period_id);

-- ============================================================
-- 3. EMPLEADOS DE LA CORRIDA
-- ============================================================
CREATE TABLE payroll_run_employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    employment_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','CALCULATING','CALCULATED','ERROR','VALIDATED','APPROVED'
    )),
    gross_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    gross_non_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    deductions_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_contributions DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    calculation_version INT NOT NULL DEFAULT 1,
    error_message TEXT,
    calculated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pe_run_employee ON payroll_run_employees(run_id, employee_id);
CREATE UNIQUE INDEX uq_pe_run_employee ON payroll_run_employees(run_id, employee_id);

-- ============================================================
-- 4. SNAPSHOT DEL EMPLEADO (congela datos al momento del cálculo)
-- ============================================================
CREATE TABLE payroll_employee_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    employee_data JSONB NOT NULL,
    employment_data JSONB,
    position_data JSONB,
    category_data JSONB,
    agreement_data JSONB,
    salary_data JSONB,
    benefits_data JSONB,
    tax_config_data JSONB,
    social_security_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_snapshot_run_employee ON payroll_employee_snapshots(run_employee_id);

-- ============================================================
-- 5. CONCEPTOS DE LIQUIDACIÓN
-- ============================================================
CREATE TABLE payroll_concepts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    concept_type VARCHAR(30) NOT NULL CHECK (concept_type IN (
        'EARNING','DEDUCTION','EMPLOYER_CONTRIBUTION','INFORMATIONAL'
    )),
    taxability VARCHAR(30) NOT NULL DEFAULT 'REMUNERATIVO' CHECK (taxability IN (
        'REMUNERATIVO','NO_REMUNERATIVO'
    )),
    calculation_type VARCHAR(30) NOT NULL DEFAULT 'AMOUNT' CHECK (calculation_type IN (
        'AMOUNT','PERCENTAGE','FORMULA','TABLE','DAILY','HOURLY','UNIT'
    )),
    base_concept_id UUID REFERENCES payroll_concepts(id),
    active BOOLEAN NOT NULL DEFAULT true,
    effective_from DATE NOT NULL,
    effective_to DATE,
    sort_order INT NOT NULL DEFAULT 0,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_payroll_concept_code UNIQUE (company_id, code)
);

CREATE INDEX idx_payroll_concepts_company ON payroll_concepts(company_id);

-- ============================================================
-- 6. REGLAS DE CÁLCULO (parametrizables, versionadas)
-- ============================================================
CREATE TABLE payroll_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    country VARCHAR(5) NOT NULL DEFAULT 'AR',
    jurisdiction VARCHAR(10),
    agreement_id UUID,
    category_id UUID,
    concept_id UUID NOT NULL REFERENCES payroll_concepts(id),
    rule_type VARCHAR(30) NOT NULL CHECK (rule_type IN (
        'PERCENTAGE','AMOUNT','FORMULA','TABLE','MULTIPLIER','THRESHOLD','EXEMPTION'
    )),
    formula TEXT,
    parameters JSONB NOT NULL DEFAULT '{}',
    priority INT NOT NULL DEFAULT 0,
    effective_from DATE NOT NULL,
    effective_to DATE,
    version INT NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payroll_rules_concept ON payroll_rules(company_id, concept_id, effective_from, effective_to);
CREATE INDEX idx_payroll_rules_agreement ON payroll_rules(agreement_id);

-- ============================================================
-- 7. NOVEDADES
-- ============================================================
CREATE TABLE payroll_novelties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    period_id UUID REFERENCES payroll_periods(id),
    novelty_type VARCHAR(50) NOT NULL,
    quantity DECIMAL(10,2),
    unit VARCHAR(20),
    amount DECIMAL(14,2),
    unit_value DECIMAL(14,2),
    multiplier DECIMAL(10,4) DEFAULT 1,
    start_date DATE,
    end_date DATE,
    description TEXT,
    source VARCHAR(30) NOT NULL DEFAULT 'MANUAL' CHECK (source IN (
        'MANUAL','ATTENDANCE','LEAVE','VACATION','PERFORMANCE','COMPENSATION','IMPORT','API','INTEGRATION'
    )),
    source_reference_id VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','APPROVED','REJECTED','APPLIED'
    )),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_novelties_employee_period ON payroll_novelties(company_id, employee_id, period_id);
CREATE INDEX idx_novelties_period ON payroll_novelties(period_id);

-- ============================================================
-- 8. ÍTEMS DE LIQUIDACIÓN (resultado del cálculo)
-- ============================================================
CREATE TABLE payroll_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    concept_id UUID NOT NULL REFERENCES payroll_concepts(id),
    quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit_value DECIMAL(14,4) NOT NULL DEFAULT 0,
    base_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    rate DECIMAL(10,6),
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    is_remunerative BOOLEAN NOT NULL DEFAULT false,
    is_deduction BOOLEAN NOT NULL DEFAULT false,
    is_employer_contribution BOOLEAN NOT NULL DEFAULT false,
    calculation_detail JSONB,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_items_run_employee ON payroll_items(run_employee_id);

-- ============================================================
-- 9. BASES IMPONIBLES
-- ============================================================
CREATE TABLE payroll_bases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    base_type VARCHAR(50) NOT NULL,
    base_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    concept_ids UUID[] NOT NULL DEFAULT '{}',
    calculation_detail JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bases_run_employee ON payroll_bases(run_employee_id);

-- ============================================================
-- 10. APORTES DEL EMPLEADO (seguimiento)
-- ============================================================
CREATE TABLE payroll_deductions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    concept_id UUID NOT NULL REFERENCES payroll_concepts(id),
    base_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    rate DECIMAL(10,6),
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    cap_amount DECIMAL(14,2),
    exceeded_amount DECIMAL(14,2) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deductions_run_employee ON payroll_deductions(run_employee_id);

-- ============================================================
-- 11. CONTRIBUCIONES PATRONALES (seguimiento)
-- ============================================================
CREATE TABLE payroll_contributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    concept_id UUID NOT NULL REFERENCES payroll_concepts(id),
    base_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    rate DECIMAL(10,6),
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    cap_amount DECIMAL(14,2),
    exceeded_amount DECIMAL(14,2) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contributions_run_employee ON payroll_contributions(run_employee_id);

-- ============================================================
-- 12. CONVENIOS COLECTIVOS
-- ============================================================
CREATE TABLE labor_agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    activity VARCHAR(255),
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE')),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_labor_agreement_code UNIQUE (company_id, code)
);

-- ============================================================
-- 13. CATEGORÍAS LABORALES
-- ============================================================
CREATE TABLE labor_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    agreement_id UUID REFERENCES labor_agreements(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_labor_category UNIQUE (company_id, agreement_id, code)
);

-- ============================================================
-- 14. ESCALAS SALARIALES
-- ============================================================
CREATE TABLE salary_scales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    agreement_id UUID REFERENCES labor_agreements(id),
    category_id UUID REFERENCES labor_categories(id),
    minimum_salary DECIMAL(14,2) NOT NULL,
    maximum_salary DECIMAL(14,2),
    reference_value DECIMAL(14,2),
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_salary_scales_category ON salary_scales(company_id, category_id, effective_from, effective_to);

-- ============================================================
-- 15. SALARIO MÍNIMO VITAL Y MÓVIL (parametrizable)
-- ============================================================
CREATE TABLE statutory_minimum_wages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country VARCHAR(5) NOT NULL DEFAULT 'AR',
    jurisdiction VARCHAR(10),
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    source VARCHAR(255),
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_min_wages_effective ON statutory_minimum_wages(country, effective_from, effective_to);

-- ============================================================
-- 16. TOPES Y MÍNIMOS
-- ============================================================
CREATE TABLE payroll_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    concept_id UUID REFERENCES payroll_concepts(id),
    limit_type VARCHAR(20) NOT NULL CHECK (limit_type IN (
        'MINIMUM','MAXIMUM','BASE_MINIMUM','BASE_MAXIMUM','EXEMPTION'
    )),
    minimum_amount DECIMAL(14,2),
    maximum_amount DECIMAL(14,2),
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_limits_concept ON payroll_limits(company_id, concept_id, effective_from, effective_to);

-- ============================================================
-- 17. ADELANTOS
-- ============================================================
CREATE TABLE employee_advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    request_date DATE NOT NULL,
    installments INT NOT NULL DEFAULT 1,
    installment_amount DECIMAL(14,2),
    remaining_amount DECIMAL(14,2) NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','APPROVED','PAYING','PAID','CANCELLED'
    )),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_advances_employee ON employee_advances(company_id, employee_id);

-- ============================================================
-- 18. EMBARGOS
-- ============================================================
CREATE TABLE payroll_garnishments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    court_order_number VARCHAR(100) NOT NULL,
    court_name VARCHAR(255),
    type VARCHAR(20) NOT NULL DEFAULT 'PERCENTAGE' CHECK (type IN (
        'PERCENTAGE','FIXED_AMOUNT','BOTH'
    )),
    percentage DECIMAL(5,2),
    fixed_amount DECIMAL(14,2),
    priority INT NOT NULL DEFAULT 1,
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','SUSPENDED','FINISHED','CANCELLED')),
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_garnishments_employee ON payroll_garnishments(company_id, employee_id);

-- ============================================================
-- 19. RETENCIONES CONFIGURABLES
-- ============================================================
CREATE TABLE payroll_withholdings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    withholding_type VARCHAR(30) NOT NULL CHECK (withholding_type IN (
        'TAX','UNION','ADVANCE','GARNISHMENT','OTHER'
    )),
    concept_id UUID REFERENCES payroll_concepts(id),
    base_amount DECIMAL(14,2),
    rate DECIMAL(10,6),
    amount DECIMAL(14,2),
    priority INT NOT NULL DEFAULT 0,
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE')),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withholdings_employee ON payroll_withholdings(company_id, employee_id);

-- ============================================================
-- 20. ERRORES DE LIQUIDACIÓN
-- ============================================================
CREATE TABLE payroll_errors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    employee_id UUID REFERENCES employees(id),
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('INFO','WARNING','ERROR','BLOCKING')),
    code VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    field VARCHAR(100),
    resolved BOOLEAN NOT NULL DEFAULT false,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_errors_run ON payroll_errors(run_id);
CREATE INDEX idx_errors_employee ON payroll_errors(employee_id);

-- ============================================================
-- 21. AUDITORÍA
-- ============================================================
CREATE TABLE payroll_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_company ON payroll_audit_logs(company_id, created_at);

-- ============================================================
-- RBAC Permissions
-- ============================================================
INSERT INTO role_permissions (name, description, module) VALUES
    ('payroll.period.list', 'List payroll periods', 'payroll'),
    ('payroll.period.create', 'Create payroll period', 'payroll'),
    ('payroll.period.view', 'View payroll period', 'payroll'),
    ('payroll.period.close', 'Close payroll period', 'payroll'),
    ('payroll.run.list', 'List payroll runs', 'payroll'),
    ('payroll.run.create', 'Create payroll run', 'payroll'),
    ('payroll.run.view', 'View payroll run', 'payroll'),
    ('payroll.run.calculate', 'Calculate payroll run', 'payroll'),
    ('payroll.run.validate', 'Validate payroll run', 'payroll'),
    ('payroll.run.approve', 'Approve payroll run', 'payroll'),
    ('payroll.run.close', 'Close payroll run', 'payroll'),
    ('payroll.concept.list', 'List concepts', 'payroll'),
    ('payroll.concept.create', 'Create concept', 'payroll'),
    ('payroll.concept.update', 'Update concept', 'payroll'),
    ('payroll.rule.list', 'List rules', 'payroll'),
    ('payroll.rule.create', 'Create rule', 'payroll'),
    ('payroll.rule.update', 'Update rule', 'payroll'),
    ('payroll.novelty.list', 'List novelties', 'payroll'),
    ('payroll.novelty.create', 'Create novelty', 'payroll'),
    ('payroll.novelty.update', 'Update novelty', 'payroll'),
    ('payroll.novelty.delete', 'Delete novelty', 'payroll'),
    ('payroll.novelty.import', 'Import novelties', 'payroll'),
    ('payroll.novelty.approve', 'Approve novelty', 'payroll'),
    ('payroll.report.view', 'View payroll reports', 'payroll')
ON CONFLICT (name) DO NOTHING;
