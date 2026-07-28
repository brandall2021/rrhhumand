CREATE TABLE attendance_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    tolerance_in_minutes INTEGER NOT NULL DEFAULT 0,
    tolerance_out_minutes INTEGER NOT NULL DEFAULT 0,
    allow_mobile BOOLEAN NOT NULL DEFAULT true,
    allow_web BOOLEAN NOT NULL DEFAULT true,
    allow_kiosk BOOLEAN NOT NULL DEFAULT true,
    require_gps BOOLEAN NOT NULL DEFAULT false,
    allow_remote BOOLEAN NOT NULL DEFAULT true,
    calculate_overtime BOOLEAN NOT NULL DEFAULT true,
    require_correction_approval BOOLEAN NOT NULL DEFAULT true,
    max_consecutive_absences INTEGER DEFAULT 3,
    work_start_time TIME NOT NULL DEFAULT '08:00',
    work_end_time TIME NOT NULL DEFAULT '16:00',
    break_duration_minutes INTEGER NOT NULL DEFAULT 60,
    work_days INTEGER[] NOT NULL DEFAULT '{1,2,3,4,5}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_attendance_policies_company ON attendance_policies(company_id);

CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    work_date DATE NOT NULL,
    scheduled_start TIMESTAMPTZ,
    scheduled_end TIMESTAMPTZ,
    actual_start TIMESTAMPTZ,
    actual_end TIMESTAMPTZ,
    scheduled_minutes INTEGER DEFAULT 0,
    worked_minutes INTEGER DEFAULT 0,
    late_minutes INTEGER DEFAULT 0,
    effective_late_minutes INTEGER DEFAULT 0,
    early_leave_minutes INTEGER DEFAULT 0,
    overtime_minutes INTEGER DEFAULT 0,
    break_minutes INTEGER DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'INCOMPLETE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(employee_id, work_date)
);
CREATE INDEX idx_attendance_records_company ON attendance_records(company_id);
CREATE INDEX idx_attendance_records_employee ON attendance_records(employee_id);
CREATE INDEX idx_attendance_records_date ON attendance_records(work_date);
CREATE INDEX idx_attendance_records_status ON attendance_records(status);

CREATE TABLE attendance_punches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    attendance_id UUID,
    punch_type VARCHAR(30) NOT NULL,
    punched_at TIMESTAMPTZ NOT NULL,
    source VARCHAR(30) NOT NULL DEFAULT 'WEB',
    latitude NUMERIC(10,7),
    longitude NUMERIC(10,7),
    ip_address INET,
    device_id VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_attendance_punches_company ON attendance_punches(company_id);
CREATE INDEX idx_attendance_punches_employee ON attendance_punches(employee_id);
CREATE INDEX idx_attendance_punches_attendance ON attendance_punches(attendance_id);

CREATE TABLE attendance_corrections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    attendance_id UUID,
    requested_by UUID NOT NULL,
    approved_by UUID,
    correction_type VARCHAR(50) NOT NULL,
    requested_value TIMESTAMPTZ,
    original_value TIMESTAMPTZ,
    reason TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);
CREATE INDEX idx_attendance_corrections_company ON attendance_corrections(company_id);
CREATE INDEX idx_attendance_corrections_employee ON attendance_corrections(employee_id);
CREATE INDEX idx_attendance_corrections_status ON attendance_corrections(status);

CREATE TABLE attendance_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    latitude NUMERIC(10,7) NOT NULL,
    longitude NUMERIC(10,7) NOT NULL,
    radius_meters INTEGER NOT NULL DEFAULT 150,
    branch_id UUID,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_attendance_locations_company ON attendance_locations(company_id);

CREATE TABLE attendance_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    name VARCHAR(150) NOT NULL,
    location VARCHAR(255),
    branch_id UUID,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_attendance_devices_company_id ON attendance_devices(company_id, device_id);

INSERT INTO permissions (id, name, resource, action, description, created_at) VALUES
    (gen_random_uuid(), 'attendance.read', 'attendance', 'read', 'View attendance records', NOW()),
    (gen_random_uuid(), 'attendance.clock_in', 'attendance', 'clock_in', 'Clock in', NOW()),
    (gen_random_uuid(), 'attendance.clock_out', 'attendance', 'clock_out', 'Clock out', NOW()),
    (gen_random_uuid(), 'attendance.manage', 'attendance', 'manage', 'Manage attendance records', NOW()),
    (gen_random_uuid(), 'attendance.view_team', 'attendance', 'view_team', 'View team attendance', NOW()),
    (gen_random_uuid(), 'attendance.view_reports', 'attendance', 'view_reports', 'View attendance reports', NOW()),
    (gen_random_uuid(), 'attendance.export', 'attendance', 'export', 'Export attendance reports', NOW()),
    (gen_random_uuid(), 'attendance.correct', 'attendance', 'correct', 'Request attendance correction', NOW()),
    (gen_random_uuid(), 'attendance.approve_correction', 'attendance', 'approve_correction', 'Approve attendance correction', NOW()),
    (gen_random_uuid(), 'attendance.manage_policies', 'attendance', 'manage_policies', 'Manage attendance policies', NOW()),
    (gen_random_uuid(), 'attendance.manage_devices', 'attendance', 'manage_devices', 'Manage attendance devices', NOW()),
    (gen_random_uuid(), 'attendance.view_calendar', 'attendance', 'view_calendar', 'View attendance calendar', NOW())
ON CONFLICT (name) DO NOTHING;
