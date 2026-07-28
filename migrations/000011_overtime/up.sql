CREATE TABLE overtime_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    max_daily_minutes INTEGER DEFAULT 120,
    max_weekly_minutes INTEGER DEFAULT 480,
    max_monthly_minutes INTEGER DEFAULT 1920,
    requires_approval BOOLEAN NOT NULL DEFAULT TRUE,
    allows_compensation BOOLEAN NOT NULL DEFAULT TRUE,
    allows_payment BOOLEAN NOT NULL DEFAULT TRUE,
    minimum_overtime_minutes INTEGER DEFAULT 0,
    rounding_minutes INTEGER DEFAULT 1,
    overtime_expiration_days INTEGER DEFAULT 0,
    night_start TIME DEFAULT '22:00:00',
    night_end TIME DEFAULT '06:00:00',
    weekend_multiplier NUMERIC(3,2) DEFAULT 1.50,
    holiday_multiplier NUMERIC(3,2) DEFAULT 2.00,
    night_multiplier NUMERIC(3,2) DEFAULT 1.50,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_overtime_policies_company ON overtime_policies(company_id);

CREATE TABLE overtime_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    attendance_id UUID,
    work_date DATE NOT NULL,
    planned_minutes INTEGER NOT NULL DEFAULT 0,
    actual_minutes INTEGER NOT NULL DEFAULT 0,
    late_minutes INTEGER NOT NULL DEFAULT 0,
    early_leave_minutes INTEGER NOT NULL DEFAULT 0,
    overtime_minutes INTEGER NOT NULL DEFAULT 0,
    approved_minutes INTEGER NOT NULL DEFAULT 0,
    compensated_minutes INTEGER NOT NULL DEFAULT 0,
    paid_minutes INTEGER NOT NULL DEFAULT 0,
    overtime_type VARCHAR(50) NOT NULL DEFAULT 'REGULAR',
    status VARCHAR(30) NOT NULL DEFAULT 'DETECTED',
    is_weekend BOOLEAN DEFAULT FALSE,
    is_holiday BOOLEAN DEFAULT FALSE,
    is_night BOOLEAN DEFAULT FALSE,
    reason TEXT,
    rejection_reason TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_overtime_records_company ON overtime_records(company_id);
CREATE INDEX idx_overtime_records_employee ON overtime_records(employee_id);
CREATE INDEX idx_overtime_records_date ON overtime_records(work_date);
CREATE INDEX idx_overtime_records_status ON overtime_records(status);
CREATE UNIQUE INDEX idx_overtime_records_unique ON overtime_records(employee_id, work_date, overtime_type);

CREATE TABLE overtime_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    overtime_record_id UUID,
    work_date DATE NOT NULL,
    requested_minutes INTEGER NOT NULL,
    approved_minutes INTEGER DEFAULT 0,
    reason TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT
);
CREATE INDEX idx_overtime_requests_company ON overtime_requests(company_id);
CREATE INDEX idx_overtime_requests_employee ON overtime_requests(employee_id);
CREATE INDEX idx_overtime_requests_status ON overtime_requests(status);

CREATE TABLE compensation_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    work_date DATE NOT NULL,
    minutes INTEGER NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT
);
CREATE INDEX idx_compensation_requests_company ON compensation_requests(company_id);
CREATE INDEX idx_compensation_requests_employee ON compensation_requests(employee_id);

CREATE TABLE employee_time_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    balance_minutes INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, employee_id)
);

CREATE TABLE time_balance_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    overtime_record_id UUID,
    transaction_type VARCHAR(30) NOT NULL,
    minutes INTEGER NOT NULL,
    reason TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_time_balance_tx_company ON time_balance_transactions(company_id);
CREATE INDEX idx_time_balance_tx_employee ON time_balance_transactions(employee_id);

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'overtime.read', 'overtime', 'read', 'View overtime records', NOW()),
    (gen_random_uuid(), 'overtime.create', 'overtime', 'create', 'Create overtime records', NOW()),
    (gen_random_uuid(), 'overtime.update', 'overtime', 'update', 'Update overtime records', NOW()),
    (gen_random_uuid(), 'overtime.delete', 'overtime', 'delete', 'Delete overtime records', NOW()),
    (gen_random_uuid(), 'overtime.approve', 'overtime', 'approve', 'Approve/reject overtime', NOW()),
    (gen_random_uuid(), 'overtime.request', 'overtime', 'request', 'Request overtime hours', NOW()),
    (gen_random_uuid(), 'overtime.detect', 'overtime', 'detect', 'Detect overtime from attendance', NOW()),
    (gen_random_uuid(), 'overtime.manage_policy', 'overtime', 'manage_policy', 'Manage overtime policies', NOW()),
    (gen_random_uuid(), 'overtime.compensate', 'overtime', 'compensate', 'Compensate overtime hours', NOW()),
    (gen_random_uuid(), 'overtime.view_balance', 'overtime', 'view_balance', 'View time balances', NOW()),
    (gen_random_uuid(), 'overtime.adjust_balance', 'overtime', 'adjust_balance', 'Adjust time balances manually', NOW()),
    (gen_random_uuid(), 'overtime.view_reports', 'overtime', 'view_reports', 'View overtime reports', NOW())
ON CONFLICT (name) DO NOTHING;
