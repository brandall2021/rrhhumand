CREATE TABLE work_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    schedule_type VARCHAR(30) NOT NULL DEFAULT 'FIXED',
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    weekly_hours INTEGER DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_work_schedules_company ON work_schedules(company_id);

CREATE TABLE work_schedule_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL,
    weekday INTEGER NOT NULL,
    is_working_day BOOLEAN NOT NULL DEFAULT TRUE,
    start_time TIME,
    end_time TIME,
    break_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_work_schedule_days_schedule ON work_schedule_days(schedule_id);
CREATE UNIQUE INDEX idx_work_schedule_days_unique ON work_schedule_days(schedule_id, weekday);

CREATE TABLE work_schedule_intervals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_day_id UUID NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    interval_type VARCHAR(30) NOT NULL DEFAULT 'WORK',
    sequence INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_work_schedule_intervals_day ON work_schedule_intervals(schedule_day_id);

CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    crosses_midnight BOOLEAN NOT NULL DEFAULT FALSE,
    break_minutes INTEGER DEFAULT 0,
    color VARCHAR(7),
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_shifts_company ON shifts(company_id);

CREATE TABLE employee_schedule_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    schedule_id UUID NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    priority INTEGER DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_emp_sched_assign_company ON employee_schedule_assignments(company_id);
CREATE INDEX idx_emp_sched_assign_employee ON employee_schedule_assignments(employee_id);

CREATE TABLE employee_shift_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    shift_id UUID NOT NULL,
    work_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(employee_id, work_date)
);
CREATE INDEX idx_emp_shift_assign_company ON employee_shift_assignments(company_id);
CREATE INDEX idx_emp_shift_assign_employee ON employee_shift_assignments(employee_id);
CREATE INDEX idx_emp_shift_assign_date ON employee_shift_assignments(work_date);

CREATE TABLE rotation_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    cycle_length INTEGER NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rotation_templates_company ON rotation_templates(company_id);

CREATE TABLE rotation_template_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    day_position INTEGER NOT NULL,
    shift_id UUID,
    is_rest_day BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rotation_template_days_template ON rotation_template_days(template_id);
CREATE UNIQUE INDEX idx_rotation_template_days_unique ON rotation_template_days(template_id, day_position);

CREATE TABLE employee_rotation_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    template_id UUID NOT NULL,
    start_date DATE NOT NULL,
    cycle_position INTEGER NOT NULL DEFAULT 1,
    effective_to DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_emp_rotation_assign_company ON employee_rotation_assignments(company_id);
CREATE INDEX idx_emp_rotation_assign_employee ON employee_rotation_assignments(employee_id);

CREATE TABLE employee_work_calendar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    work_date DATE NOT NULL,
    schedule_id UUID,
    shift_id UUID,
    planned_start TIMESTAMPTZ,
    planned_end TIMESTAMPTZ,
    planned_break_minutes INTEGER DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED',
    source VARCHAR(30) NOT NULL DEFAULT 'MANUAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(employee_id, work_date)
);
CREATE INDEX idx_emp_work_calendar_company ON employee_work_calendar(company_id);
CREATE INDEX idx_emp_work_calendar_employee ON employee_work_calendar(employee_id);
CREATE INDEX idx_emp_work_calendar_date ON employee_work_calendar(work_date);

CREATE TABLE schedule_exceptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID,
    schedule_id UUID,
    exception_date DATE NOT NULL,
    exception_type VARCHAR(50) NOT NULL,
    start_time TIME,
    end_time TIME,
    shift_id UUID,
    reason TEXT,
    approved_by UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_schedule_exceptions_company ON schedule_exceptions(company_id);
CREATE INDEX idx_schedule_exceptions_employee ON schedule_exceptions(employee_id);

CREATE TABLE shift_swaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    requester_employee_id UUID NOT NULL,
    target_employee_id UUID NOT NULL,
    requester_date DATE NOT NULL,
    target_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    approved_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_shift_swaps_company ON shift_swaps(company_id);

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'scheduling.read', 'scheduling', 'read', 'View schedules and shifts', NOW()),
    (gen_random_uuid(), 'scheduling.create_schedule', 'scheduling', 'create_schedule', 'Create work schedules', NOW()),
    (gen_random_uuid(), 'scheduling.update_schedule', 'scheduling', 'update_schedule', 'Update work schedules', NOW()),
    (gen_random_uuid(), 'scheduling.delete_schedule', 'scheduling', 'delete_schedule', 'Delete work schedules', NOW()),
    (gen_random_uuid(), 'scheduling.create_shift', 'scheduling', 'create_shift', 'Create shifts', NOW()),
    (gen_random_uuid(), 'scheduling.update_shift', 'scheduling', 'update_shift', 'Update shifts', NOW()),
    (gen_random_uuid(), 'scheduling.delete_shift', 'scheduling', 'delete_shift', 'Delete shifts', NOW()),
    (gen_random_uuid(), 'scheduling.assign', 'scheduling', 'assign', 'Assign schedules/shifts to employees', NOW()),
    (gen_random_uuid(), 'scheduling.generate_calendar', 'scheduling', 'generate_calendar', 'Generate work calendar', NOW()),
    (gen_random_uuid(), 'scheduling.manage_exceptions', 'scheduling', 'manage_exceptions', 'Manage schedule exceptions', NOW()),
    (gen_random_uuid(), 'scheduling.approve_swap', 'scheduling', 'approve_swap', 'Approve shift swaps', NOW()),
    (gen_random_uuid(), 'scheduling.view_team', 'scheduling', 'view_team', 'View team schedules', NOW())
ON CONFLICT (name) DO NOTHING;
