-- FASE 19B — Payroll Features: Recibos, ARCA, Libro de Sueldos, Bancos, Contabilidad, Reportes
-- ============================================================

-- ============================================================
-- 1. RECIBOS DE SUELDO (digitales)
-- ============================================================
CREATE TABLE payroll_receipt_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_html TEXT NOT NULL,
    template_css TEXT,
    orientation VARCHAR(10) NOT NULL DEFAULT 'PORTRAIT' CHECK (orientation IN ('PORTRAIT','LANDSCAPE')),
    paper_size VARCHAR(20) NOT NULL DEFAULT 'A4',
    show_logo BOOLEAN NOT NULL DEFAULT true,
    show_signature BOOLEAN NOT NULL DEFAULT true,
    show_qr BOOLEAN NOT NULL DEFAULT false,
    show_barcode BOOLEAN NOT NULL DEFAULT false,
    font_family VARCHAR(50) NOT NULL DEFAULT 'Arial',
    font_size INT NOT NULL DEFAULT 10,
    primary_color VARCHAR(7) NOT NULL DEFAULT '#003366',
    secondary_color VARCHAR(7) NOT NULL DEFAULT '#CC0000',
    margin_top DECIMAL(4,1) NOT NULL DEFAULT 15.0,
    margin_bottom DECIMAL(4,1) NOT NULL DEFAULT 15.0,
    margin_left DECIMAL(4,1) NOT NULL DEFAULT 15.0,
    margin_right DECIMAL(4,1) NOT NULL DEFAULT 15.0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_receipt_template_company UNIQUE (company_id, name)
);

CREATE TABLE payroll_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    template_id UUID REFERENCES payroll_receipt_templates(id),
    receipt_number VARCHAR(20) NOT NULL,
    cuit VARCHAR(13) NOT NULL,
    employee_cuil VARCHAR(13) NOT NULL,
    period_name VARCHAR(255) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    payment_date DATE,
    gross_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    gross_non_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    deductions_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    contributions_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    amount_in_words VARCHAR(500),
    digital_token VARCHAR(255),
    qr_code TEXT,
    barcode TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATED' CHECK (status IN (
        'GENERATED','VIEWED','DOWNLOADED','EMAILED','ACKNOWLEDGED','CANCELLED'
    )),
    acknowledged_at TIMESTAMPTZ,
    acknowledged_ip VARCHAR(45),
    viewed_at TIMESTAMPTZ,
    downloaded_at TIMESTAMPTZ,
    emailed_at TIMESTAMPTZ,
    storage_path TEXT,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_receipt_employee_run UNIQUE (employee_id, run_id)
);

CREATE INDEX idx_receipts_run ON payroll_receipts(run_id);
CREATE INDEX idx_receipts_employee ON payroll_receipts(employee_id, company_id);
CREATE INDEX idx_receipts_company ON payroll_receipts(company_id, generated_at DESC);
CREATE INDEX idx_receipts_status ON payroll_receipts(status);

CREATE TABLE payroll_receipt_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_id UUID NOT NULL REFERENCES payroll_receipts(id),
    concept_code VARCHAR(50) NOT NULL,
    concept_name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit_value DECIMAL(14,2) NOT NULL DEFAULT 0,
    base_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    rate DECIMAL(10,6),
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    is_remunerative BOOLEAN NOT NULL DEFAULT false,
    is_deduction BOOLEAN NOT NULL DEFAULT false,
    is_contribution BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_receipt_items_receipt ON payroll_receipt_items(receipt_id);

-- ============================================================
-- 2. ARCA (ex AFIP) — Mapeo de conceptos y exportaciones
-- ============================================================
CREATE TABLE arca_concept_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    concept_id UUID NOT NULL REFERENCES payroll_concepts(id),
    arca_concept_code VARCHAR(20) NOT NULL,
    arca_concept_name VARCHAR(255),
    mapping_type VARCHAR(30) NOT NULL CHECK (mapping_type IN (
        'SICOSS_REMUNERACION','SICOSS_DEDUCCION','SICOSS_CONTRIBUCION',
        'SIAP_REMUNERACION','SIAP_DEDUCCION','SIAP_CONTRIBUCION',
        'F931_REMUNERACION','F931_DEDUCCION','F931_CONTRIBUCION',
        'F1357_REMUNERACION','GENERAL'
    )),
    percentage DECIMAL(5,2),
    is_taxable BOOLEAN NOT NULL DEFAULT true,
    is_contributable BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_arca_mapping UNIQUE (company_id, concept_id, mapping_type, effective_from)
);

CREATE INDEX idx_arca_mappings_concept ON arca_concept_mappings(concept_id);
CREATE INDEX idx_arca_mappings_type ON arca_concept_mappings(mapping_type, effective_from);

CREATE TABLE arca_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    export_type VARCHAR(30) NOT NULL CHECK (export_type IN (
        'SICOSS','SIAP','F931','F1357','CEL','GENERAL'
    )),
    period_id UUID REFERENCES payroll_periods(id),
    file_name VARCHAR(255) NOT NULL,
    file_content TEXT,
    storage_path TEXT,
    file_size INT,
    checksum VARCHAR(64),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','GENERATED','VALIDATED','ERROR','SUBMITTED','ACKNOWLEDGED'
    )),
    error_message TEXT,
    submission_date TIMESTAMPTZ,
    acknowledgement_code VARCHAR(100),
    employee_count INT DEFAULT 0,
    total_amount DECIMAL(14,2) DEFAULT 0,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_arca_exports_run ON arca_exports(run_id);
CREATE INDEX idx_arca_exports_company ON arca_exports(company_id, export_type, generated_at DESC);

-- ============================================================
-- 3. LIBRO DE SUELDOS DIGITAL
-- ============================================================
CREATE TABLE payroll_book_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    run_employee_id UUID NOT NULL REFERENCES payroll_run_employees(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    entry_type VARCHAR(30) NOT NULL CHECK (entry_type IN (
        'ORDINARY','ADJUSTMENT','VACATION','SAC','FINAL','LIQUIDATION'
    )),
    cuil VARCHAR(13) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    nationality VARCHAR(50),
    birth_date DATE,
    sex CHAR(1) CHECK (sex IN ('M','F','X')),
    admission_date DATE NOT NULL,
    discharge_date DATE,
    category_code VARCHAR(50),
    category_name VARCHAR(255),
    agreement_code VARCHAR(20),
    agreement_name VARCHAR(255),
    work_type VARCHAR(50),
    work_place VARCHAR(255),
    gross_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    gross_non_remunerative DECIMAL(14,2) NOT NULL DEFAULT 0,
    deductions_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    contributions_total DECIMAL(14,2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    employer_cost DECIMAL(14,2) NOT NULL DEFAULT 0,
    days_worked INT NOT NULL DEFAULT 0,
    hours_worked DECIMAL(8,2) DEFAULT 0,
    absences INT DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    book_number INT,
    page_number INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_book_entries_run ON payroll_book_entries(run_id);
CREATE INDEX idx_book_entries_employee ON payroll_book_entries(employee_id, company_id);
CREATE UNIQUE INDEX uq_book_entry_run_employee ON payroll_book_entries(run_id, run_employee_id);

CREATE TABLE payroll_book_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    period_id UUID REFERENCES payroll_periods(id),
    year INT NOT NULL,
    month INT NOT NULL,
    export_type VARCHAR(20) NOT NULL DEFAULT 'DIGITAL' CHECK (export_type IN ('DIGITAL','PRINT','CSV','PDF','CEL')),
    file_name VARCHAR(255) NOT NULL,
    file_content TEXT,
    storage_path TEXT,
    file_size INT,
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATED' CHECK (status IN (
        'GENERATED','SUBMITTED','ACKNOWLEDGED','REJECTED'
    )),
    submission_date TIMESTAMPTZ,
    acknowledgement_code VARCHAR(100),
    employee_count INT DEFAULT 0,
    total_gross DECIMAL(14,2) DEFAULT 0,
    total_deductions DECIMAL(14,2) DEFAULT 0,
    total_net DECIMAL(14,2) DEFAULT 0,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 4. BANCOS — Lotes de pago
-- ============================================================
CREATE TABLE bank_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    batch_number VARCHAR(30) NOT NULL,
    bank_code VARCHAR(20) NOT NULL,
    bank_name VARCHAR(255),
    payment_type VARCHAR(30) NOT NULL DEFAULT 'SALARY' CHECK (payment_type IN (
        'SALARY','SAC','VACATION','ADVANCE','BONUS','LIQUIDATION','OTHER'
    )),
    total_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    total_employees INT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','GENERATED','SENT','PROCESSED','PARTIALLY_PAID','PAID','REJECTED','CANCELLED'
    )),
    file_format VARCHAR(20) NOT NULL DEFAULT 'CSV' CHECK (file_format IN (
        'CSV','TXT','XML','XLS','COBIS','NACION','BBVA','MACRO','SANTANDER','GALICIA','HSBC','CUSTOM'
    )),
    file_name VARCHAR(255),
    storage_path TEXT,
    file_content TEXT,
    sent_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    notes TEXT,
    generated_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bank_batches_run ON bank_batches(run_id);
CREATE INDEX idx_bank_batches_company ON bank_batches(company_id, status, payment_date);

CREATE TABLE bank_batch_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES bank_batches(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    run_employee_id UUID REFERENCES payroll_run_employees(id),
    cuil VARCHAR(13) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    bank_code VARCHAR(20),
    bank_name VARCHAR(255),
    branch_code VARCHAR(20),
    account_type VARCHAR(20) CHECK (account_type IN ('CAJA_AHORRO','CUENTA_CORRIENTE','CUENTA_SUELDO','OTHER')),
    account_number VARCHAR(30),
    cbu VARCHAR(22),
    alias VARCHAR(50),
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    concept VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING','PAID','REJECTED','RETURNED'
    )),
    error_message TEXT,
    payment_date DATE,
    transaction_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bank_items_batch ON bank_batch_items(batch_id);
CREATE INDEX idx_bank_items_employee ON bank_batch_items(employee_id);

-- ============================================================
-- 5. CONTABILIDAD — Exportación de asientos contables
-- ============================================================
CREATE TABLE accounting_account_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    concept_id UUID REFERENCES payroll_concepts(id),
    mapping_type VARCHAR(30) NOT NULL CHECK (mapping_type IN (
        'EARNING','DEDUCTION','CONTRIBUTION','PROVISION','EMPLOYER_COST','NET_PAY'
    )),
    debit_account VARCHAR(20) NOT NULL,
    credit_account VARCHAR(20) NOT NULL,
    cost_center_required BOOLEAN NOT NULL DEFAULT false,
    description_template TEXT,
    priority INT NOT NULL DEFAULT 0,
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_accounting_mapping UNIQUE (company_id, concept_id, mapping_type)
);

CREATE TABLE accounting_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    period_id UUID REFERENCES payroll_periods(id),
    export_type VARCHAR(30) NOT NULL DEFAULT 'JOURNAL' CHECK (export_type IN (
        'JOURNAL','LEDGER','TRIAL_BALANCE','COST_CENTER','BUDGET'
    )),
    file_format VARCHAR(20) NOT NULL DEFAULT 'CSV' CHECK (file_format IN (
        'CSV','TXT','XML','XLS','PDF','SIAP_CONTABLE'
    )),
    file_name VARCHAR(255) NOT NULL,
    file_content TEXT,
    storage_path TEXT,
    file_size INT,
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATED' CHECK (status IN (
        'GENERATED','SUBMITTED','IMPORTED','ERROR'
    )),
    employee_count INT DEFAULT 0,
    total_debit DECIMAL(14,2) DEFAULT 0,
    total_credit DECIMAL(14,2) DEFAULT 0,
    entry_count INT DEFAULT 0,
    error_message TEXT,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounting_exports_run ON accounting_exports(run_id);

CREATE TABLE accounting_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    export_id UUID NOT NULL REFERENCES accounting_exports(id),
    entry_number INT NOT NULL,
    account_code VARCHAR(20) NOT NULL,
    account_name VARCHAR(255),
    cost_center VARCHAR(20),
    debit_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    credit_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    concept_code VARCHAR(50),
    concept_name VARCHAR(255),
    employee_id UUID REFERENCES employees(id),
    employee_name VARCHAR(255),
    document_type VARCHAR(20),
    document_number VARCHAR(50),
    reference TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounting_entries_export ON accounting_entries(export_id);
CREATE INDEX idx_accounting_entries_account ON accounting_entries(account_code, export_id);

-- ============================================================
-- 6. REPORTES CONFIGURABLES
-- ============================================================
CREATE TABLE payroll_report_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    report_type VARCHAR(30) NOT NULL CHECK (report_type IN (
        'SUMMARY','DETAILED','COMPARATIVE','COST_CENTER','DEPARTMENT',
        'AGREEMENT','CATEGORY','EVOLUTION','CUSTOM'
    )),
    config JSONB NOT NULL DEFAULT '{}',
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_report_template UNIQUE (company_id, name)
);

CREATE TABLE payroll_report_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    run_id UUID REFERENCES payroll_runs(id),
    template_id UUID REFERENCES payroll_report_templates(id),
    report_type VARCHAR(30) NOT NULL,
    file_format VARCHAR(10) NOT NULL CHECK (file_format IN ('PDF','CSV','XLS','HTML','JSON')),
    file_name VARCHAR(255) NOT NULL,
    file_content TEXT,
    storage_path TEXT,
    file_size INT,
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATED' CHECK (status IN ('GENERATED','ERROR')),
    error_message TEXT,
    config JSONB,
    generated_by UUID NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- RBAC Permissions for 19B
-- ============================================================
INSERT INTO role_permissions (name, description, module) VALUES
    ('payroll.receipt.generate', 'Generate payroll receipts', 'payroll'),
    ('payroll.receipt.view', 'View payroll receipts', 'payroll'),
    ('payroll.receipt.email', 'Email payroll receipts', 'payroll'),
    ('payroll.receipt.cancel', 'Cancel payroll receipts', 'payroll'),
    ('payroll.receipt.template.manage', 'Manage receipt templates', 'payroll'),
    ('payroll.arca.export', 'Export ARCA files', 'payroll'),
    ('payroll.arca.mapping.manage', 'Manage ARCA concept mappings', 'payroll'),
    ('payroll.arca.view', 'View ARCA exports', 'payroll'),
    ('payroll.book.generate', 'Generate libro de sueldos', 'payroll'),
    ('payroll.book.export', 'Export libro de sueldos', 'payroll'),
    ('payroll.bank.batch.create', 'Create bank batches', 'payroll'),
    ('payroll.bank.batch.send', 'Send bank batches', 'payroll'),
    ('payroll.bank.batch.view', 'View bank batches', 'payroll'),
    ('payroll.accounting.export', 'Export accounting entries', 'payroll'),
    ('payroll.accounting.mapping.manage', 'Manage accounting mappings', 'payroll'),
    ('payroll.report.generate', 'Generate payroll reports', 'payroll')
ON CONFLICT (name) DO NOTHING;
