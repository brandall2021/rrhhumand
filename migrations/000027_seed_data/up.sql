-- 000027: Seed data for development / demo
-- Uses ON CONFLICT to be idempotent (safe for re-runs)

-- ============================================================
-- 0. Create missing training tables if not exist
-- ============================================================
CREATE TABLE IF NOT EXISTS training_categories (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tcat_company ON training_categories(company_id);
CREATE INDEX IF NOT EXISTS idx_tcat_parent ON training_categories(parent_id);

CREATE TABLE IF NOT EXISTS training_courses (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id UUID REFERENCES training_categories(id) ON DELETE SET NULL,
    short_description TEXT,
    description TEXT,
    objectives TEXT,
    difficulty VARCHAR(20) DEFAULT 'beginner',
    duration_minutes INT NOT NULL DEFAULT 0,
    modality VARCHAR(30) DEFAULT 'online',
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    mandatory BOOLEAN NOT NULL DEFAULT false,
    passing_score NUMERIC(5,2) DEFAULT 70.00,
    certificate_enabled BOOLEAN NOT NULL DEFAULT false,
    min_attendance_percentage NUMERIC(5,2) DEFAULT 80.00,
    created_by UUID,
    published_by UUID,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tcrs_company ON training_courses(company_id);
CREATE INDEX IF NOT EXISTS idx_tcrs_category ON training_courses(category_id);
CREATE INDEX IF NOT EXISTS idx_tcrs_code ON training_courses(company_id, code);
CREATE INDEX IF NOT EXISTS idx_tcrs_status ON training_courses(status);

-- ============================================================
-- 1. Demo Company
-- ============================================================
INSERT INTO companies (id, name, slug, plan, settings, active)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'TechCorp Argentina',
    'techcorp-ar',
    'enterprise',
    '{"timezone":"America/Argentina/Buenos_Aires","currency":"ARS","locale":"es-AR"}',
    true
) ON CONFLICT (slug) DO NOTHING;

-- ============================================================
-- 2. Branches
-- ============================================================
INSERT INTO branches (id, company_id, name, address, city, country, active)
VALUES
    ('a0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000001', 'Oficina Buenos Aires', 'Av. Corrientes 1234', 'CABA', 'Argentina', true),
    ('a0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000001', 'Oficina Córdoba', 'Av. Colón 567', 'Córdoba', 'Argentina', true)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 3. Departments
-- ============================================================
INSERT INTO departments (id, company_id, name, code, active)
VALUES
    ('a0000000-0000-0000-0000-000000000020', 'a0000000-0000-0000-0000-000000000001', 'Ingeniería', 'ENG', true),
    ('a0000000-0000-0000-0000-000000000021', 'a0000000-0000-0000-0000-000000000001', 'Recursos Humanos', 'HR', true),
    ('a0000000-0000-0000-0000-000000000022', 'a0000000-0000-0000-0000-000000000001', 'Ventas', 'SALES', true),
    ('a0000000-0000-0000-0000-000000000023', 'a0000000-0000-0000-0000-000000000001', 'Marketing', 'MKT', true),
    ('a0000000-0000-0000-0000-000000000024', 'a0000000-0000-0000-0000-000000000001', 'Finanzas', 'FIN', true)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 4. Positions
-- ============================================================
INSERT INTO positions (id, company_id, department_id, name, code, level, active)
VALUES
    ('a0000000-0000-0000-0000-000000000030', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000020', 'Desarrollador Senior', 'SR-DEV', 4, true),
    ('a0000000-0000-0000-0000-000000000031', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000020', 'Desarrollador Junior', 'JR-DEV', 1, true),
    ('a0000000-0000-0000-0000-000000000032', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000021', 'Analista de RRHH', 'HR-ANL', 3, true),
    ('a0000000-0000-0000-0000-000000000033', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000022', 'Ejecutivo de Ventas', 'SALES-EXEC', 3, true),
    ('a0000000-0000-0000-0000-000000000034', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000024', 'Analista Financiero', 'FIN-ANL', 3, true)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 5. Admin User + password (admin123)
-- ============================================================
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000040',
    'admin@techcorp.com',
    crypt('admin123', gen_salt('bf', 10)),
    'Admin',
    'TechCorp',
    true
) ON CONFLICT (email) DO NOTHING;

-- Assign user to company
INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000040', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- Assign COMPANY_ADMIN role
INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000040', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'COMPANY_ADMIN'
ON CONFLICT DO NOTHING;

-- ============================================================
-- 6. HR User + password (hr123)
-- ============================================================
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000041',
    'hr@techcorp.com',
    crypt('hr123', gen_salt('bf', 10)),
    'Laura',
    'Rodríguez',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000041', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000041', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'HR_ADMIN'
ON CONFLICT DO NOTHING;

-- ============================================================
-- 7. Employees (with user accounts)
-- ============================================================

-- Employee 1: Manager
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000050',
    'carlos.garcia@techcorp.com',
    crypt('emp123', gen_salt('bf', 10)),
    'Carlos',
    'García',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000050', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000050', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'MANAGER'
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000060',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-001',
    'Carlos',
    'García',
    'carlos.garcia@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000020',
    'a0000000-0000-0000-0000-000000000030',
    '2020-03-15',
    'active'
) ON CONFLICT (company_id, employee_number) DO NOTHING;

-- Employee 2: Developer
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000051',
    'maria.lopez@techcorp.com',
    crypt('emp123', gen_salt('bf', 10)),
    'María',
    'López',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000051', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000051', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'EMPLOYEE'
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000061',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-002',
    'María',
    'López',
    'maria.lopez@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000020',
    'a0000000-0000-0000-0000-000000000031',
    'a0000000-0000-0000-0000-000000000060',
    '2022-06-01',
    'active'
) ON CONFLICT (company_id, employee_number) DO NOTHING;

-- Employee 3: HR
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000052',
    'ana.martinez@techcorp.com',
    crypt('emp123', gen_salt('bf', 10)),
    'Ana',
    'Martínez',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000052', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000052', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'EMPLOYEE'
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000062',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-003',
    'Ana',
    'Martínez',
    'ana.martinez@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000021',
    'a0000000-0000-0000-0000-000000000032',
    'a0000000-0000-0000-0000-000000000060',
    '2023-01-10',
    'active'
) ON CONFLICT (company_id, employee_number) DO NOTHING;

-- Employee 4: Sales
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000053',
    'pedro.fernandez@techcorp.com',
    crypt('emp123', gen_salt('bf', 10)),
    'Pedro',
    'Fernández',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000053', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000053', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'EMPLOYEE'
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000063',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-004',
    'Pedro',
    'Fernández',
    'pedro.fernandez@techcorp.com',
    'a0000000-0000-0000-0000-000000000011',
    'a0000000-0000-0000-0000-000000000022',
    'a0000000-0000-0000-0000-000000000033',
    'a0000000-0000-0000-0000-000000000060',
    '2023-03-20',
    'active'
) ON CONFLICT (company_id, employee_number) DO NOTHING;

-- Employee 5: Finance
INSERT INTO users (id, email, password_hash, first_name, last_name, active)
VALUES (
    'a0000000-0000-0000-0000-000000000054',
    'lucia.torres@techcorp.com',
    crypt('emp123', gen_salt('bf', 10)),
    'Lucía',
    'Torres',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_companies (user_id, company_id)
VALUES ('a0000000-0000-0000-0000-000000000054', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, company_id)
SELECT 'a0000000-0000-0000-0000-000000000054', id, 'a0000000-0000-0000-0000-000000000001'
FROM roles WHERE name = 'EMPLOYEE'
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, company_id, employee_number, first_name, last_name, email, branch_id, department_id, position_id, manager_id, hire_date, status)
VALUES (
    'a0000000-0000-0000-0000-000000000064',
    'a0000000-0000-0000-0000-000000000001',
    'EMP-005',
    'Lucía',
    'Torres',
    'lucia.torres@techcorp.com',
    'a0000000-0000-0000-0000-000000000010',
    'a0000000-0000-0000-0000-000000000024',
    'a0000000-0000-0000-0000-000000000034',
    'a0000000-0000-0000-0000-000000000060',
    '2023-08-01',
    'active'
) ON CONFLICT (company_id, employee_number) DO NOTHING;

-- ============================================================
-- 8. Training courses
-- ============================================================
INSERT INTO training_categories (id, company_id, name, description, active)
VALUES
    ('a0000000-0000-0000-0000-000000000070', 'a0000000-0000-0000-0000-000000000001', 'Desarrollo Profesional', 'Cursos para desarrollo de habilidades profesionales', true),
    ('a0000000-0000-0000-0000-000000000071', 'a0000000-0000-0000-0000-000000000001', 'Liderazgo', 'Programas de formación en liderazgo', true),
    ('a0000000-0000-0000-0000-000000000072', 'a0000000-0000-0000-0000-000000000001', 'Tecnología', 'Cursos técnicos y de programación', true)
ON CONFLICT DO NOTHING;

INSERT INTO training_courses (id, company_id, category_id, code, name, short_description, modality, duration_minutes, status, created_by)
VALUES
    ('a0000000-0000-0000-0000-000000000080', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000072', 'GIT-101', 'Git y Control de Versiones', 'Curso introductorio a Git y GitHub', 'online', 480, 'published', 'a0000000-0000-0000-0000-000000000040'),
    ('a0000000-0000-0000-0000-000000000081', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000070', 'COM-101', 'Comunicación Efectiva', 'Técnicas de comunicación en el ámbito laboral', 'in_person', 240, 'published', 'a0000000-0000-0000-0000-000000000040'),
    ('a0000000-0000-0000-0000-000000000082', 'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000071', 'LDR-101', 'Liderazgo de Equipos', 'Fundamentos de liderazgo y gestión de equipos', 'hybrid', 960, 'draft', 'a0000000-0000-0000-0000-000000000040')
ON CONFLICT DO NOTHING;
