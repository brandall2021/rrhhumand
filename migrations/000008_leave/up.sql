CREATE TABLE leave_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    requires_approval BOOLEAN NOT NULL DEFAULT true,
    requires_document BOOLEAN NOT NULL DEFAULT false,
    affects_balance BOOLEAN NOT NULL DEFAULT true,
    is_paid BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_leave_types_company ON leave_types(company_id);
CREATE UNIQUE INDEX idx_leave_types_company_code ON leave_types(company_id, code);

CREATE TABLE leave_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    leave_type_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    days_per_year NUMERIC(8,2),
    minimum_days_before_request INTEGER DEFAULT 0,
    maximum_days_per_request NUMERIC(8,2),
    maximum_accumulated_days NUMERIC(8,2),
    allow_negative_balance BOOLEAN NOT NULL DEFAULT false,
    use_business_days BOOLEAN NOT NULL DEFAULT true,
    requires_manager_approval BOOLEAN NOT NULL DEFAULT true,
    requires_hr_approval BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_leave_policies_company ON leave_policies(company_id);
CREATE INDEX idx_leave_policies_type ON leave_policies(leave_type_id);

CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID,
    date DATE NOT NULL,
    name VARCHAR(150) NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_holidays_company ON holidays(company_id);
CREATE INDEX idx_holidays_date ON holidays(date);

CREATE TABLE leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    leave_type_id UUID NOT NULL,
    year INTEGER NOT NULL,
    allocated_days NUMERIC(8,2) NOT NULL DEFAULT 0,
    carried_over_days NUMERIC(8,2) NOT NULL DEFAULT 0,
    adjustment_days NUMERIC(8,2) NOT NULL DEFAULT 0,
    used_days NUMERIC(8,2) NOT NULL DEFAULT 0,
    reserved_days NUMERIC(8,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (employee_id, leave_type_id, year)
);
CREATE INDEX idx_leave_balances_company ON leave_balances(company_id);
CREATE INDEX idx_leave_balances_employee ON leave_balances(employee_id);

CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    leave_type_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    requested_days NUMERIC(8,2) NOT NULL,
    reason TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    document_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_leave_requests_company ON leave_requests(company_id);
CREATE INDEX idx_leave_requests_employee ON leave_requests(employee_id);
CREATE INDEX idx_leave_requests_status ON leave_requests(status);
CREATE INDEX idx_leave_requests_dates ON leave_requests(start_date, end_date);

CREATE TABLE leave_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    leave_request_id UUID NOT NULL,
    approver_id UUID NOT NULL,
    level INTEGER NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    comments TEXT,
    decided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_leave_approvals_company ON leave_approvals(company_id);
CREATE INDEX idx_leave_approvals_request ON leave_approvals(leave_request_id);

CREATE TABLE leave_request_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    leave_request_id UUID NOT NULL,
    action VARCHAR(30) NOT NULL,
    old_status VARCHAR(30),
    new_status VARCHAR(30),
    performed_by UUID NOT NULL,
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_leave_request_history_request ON leave_request_history(leave_request_id);

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'leave.read', 'leave', 'read', 'View leave requests and balances', NOW()),
    (gen_random_uuid(), 'leave.create', 'leave', 'create', 'Create leave requests', NOW()),
    (gen_random_uuid(), 'leave.update', 'leave', 'update', 'Update leave requests', NOW()),
    (gen_random_uuid(), 'leave.cancel', 'leave', 'cancel', 'Cancel leave requests', NOW()),
    (gen_random_uuid(), 'leave.approve', 'leave', 'approve', 'Approve leave requests', NOW()),
    (gen_random_uuid(), 'leave.reject', 'leave', 'reject', 'Reject leave requests', NOW()),
    (gen_random_uuid(), 'leave.manage_types', 'leave', 'manage_types', 'Manage leave types', NOW()),
    (gen_random_uuid(), 'leave.manage_policies', 'leave', 'manage_policies', 'Manage leave policies', NOW()),
    (gen_random_uuid(), 'leave.manage_balances', 'leave', 'manage_balances', 'Manage leave balances', NOW()),
    (gen_random_uuid(), 'leave.view_team', 'leave', 'view_team', 'View team leave requests', NOW()),
    (gen_random_uuid(), 'leave.view_reports', 'leave', 'view_reports', 'View leave reports', NOW()),
    (gen_random_uuid(), 'leave.export', 'leave', 'export', 'Export leave reports', NOW())
ON CONFLICT (name) DO NOTHING;
