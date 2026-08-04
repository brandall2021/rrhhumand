# Frontend Full Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the RRH-HumanD React frontend module by module so every page is a real, working CRUD screen backed by the existing Go API.

**Architecture:** Backend (Go/Gin at `/api/v1`) is ~95% complete. Work is 100% frontend: React 19 + Vite + Tailwind 4 + Radix UI (dialog/select/tabs) + react-router 7 + axios. Each page lives in `frontend/src/pages/<module>.tsx`, uses `@/lib/api` (axios instance, `data.data` envelope), `@/components/ui/*` primitives, and Spanish UI copy. Existing pattern in `employees.tsx` / `departments.tsx` is the reference: `useState` + `fetchData` + table + `Dialog` modal form. Keep that pattern; add Radix `Tabs` for multi-section pages.

**Tech Stack:** React 19, TypeScript, Vite 8, Tailwind CSS 4, @radix-ui/react-dialog/select/tabs, lucide-react icons, axios, react-router-dom 7.

## Global Constraints

- Frontend only — do not modify Go backend unless an endpoint is genuinely missing (log the gap in the task notes; fix only when the page cannot be built without it).
- API envelope is `{success, data, error}`; lists: `data` is an array (some endpoints return `{items, total}` — handle both with a helper).
- UI copy in Spanish (rioplatense, voseo as in existing pages: "Completá", "Seleccionar...", "N°").
- Reuse existing UI primitives; do not add new component libraries.
- `npm run build` (`tsc -b && vite build`) and `npm run lint` must pass before a module is considered done.
- Follow existing file conventions (function components, `export default`, `interface` types at top, lucide icons).

---
## Priority order (from product owner)

1. Real CRUD on current screens (documents, training, expenses, compensation, recruitment)
2. Complete employee detail page (ficha)
3. Leaves & vacations
4. Attendance
5. Payroll / nómina
6. Onboarding / offboarding
7. Performance
8. Benefits / total rewards
9. Expenses & travel complete
10. Recruitment complete
11. Training / LMS complete
12. "Mi cuenta" employee portal
13. Feed & surveys
14. Admin multi-tenant
15. Reports & audit

---

## Task 1: Real CRUD on current screens

Reference API map: `docs` note — full route inventory from explore agent (stored in session). Backend routes per module:

### Task 1a: Documents (`frontend/src/pages/documents.tsx`)
Endpoints available: `GET/POST /documents`, `GET/PUT/DELETE /documents/:id`, `GET /documents/:id/download`, `POST /documents/:id/archive`, `POST /documents/:id/restore`, `GET/POST /documents/:id/versions`, `GET /documents/stats`, `GET /documents/expiring`, `GET/POST /document-categories`, `GET/PUT/DELETE /document-categories/:id`.
- [ ] Keep list (name, category, employee, size, status) + Nuevo button
- [ ] Upload modal: name, description, category select, employee select, file input (multipart)
- [ ] Row actions: Descargar (download), Editar (PUT), Archivar/Restaurar, Versiones (list + create), Eliminar
- [ ] Categories section: list + create/edit/delete modal
- [ ] Stats header (total, expiring soon)

### Task 1b: Training (`frontend/src/pages/training.tsx`)
Endpoints: `GET/POST /training/categories`, `PUT /training/categories/:id`, `GET/POST /training/courses`, `GET/PUT /training/courses/:id`, `POST /training/courses/:id/publish`, `GET /training/courses/:id/details`, `GET/POST /training/courses/:id/versions`, `GET/POST /training/courses/:id/assessments`, `GET /training/dashboard/stats`.
- [ ] Categories management (create/edit)
- [ ] Course list with search + Nuevo
- [ ] Course create/edit modal (full field set from `training_courses` columns)
- [ ] Course detail drawer/dialog: versions (create/list), publish button, contents list if available
- [ ] Dashboard stats header

### Task 1c: Expenses (`frontend/src/pages/expenses.tsx`)
Endpoints: `GET/POST /expenses`, `GET/PUT /expenses/:id`, `POST /expenses/:id/submit`, `POST /expenses/:id/approve|reject|observe|cancel`, `DELETE /expenses/:id`, `GET/POST /expense-categories` (+CRUD), `GET/POST /expenses/payment-methods`, `GET /expenses/travels` (list), `GET /expenses/reports` (list), `GET /expenses/advances` (list), `GET/POST /expense-policies`.
- [ ] Expense list with filters (status) + Nuevo
- [ ] Expense create/edit modal (category, amount, currency, date, description)
- [ ] Workflow actions: submit / approve / reject / observe / cancel
- [ ] Receipt upload per expense
- [ ] Categories + payment methods CRUD
- [ ] Tabs: Gastos | Viajes | Reportes | Anticipos | Políticas (list + create for each, approve/reject where applicable)

### Task 1d: Compensation (`frontend/src/pages/compensation.tsx`)
Endpoints: `GET/POST /compensation/structures`, `PUT /compensation/structures/:id`, `GET/POST /compensation/grades`, `PUT /compensation/grades/:id`, `GET/POST /compensation/bands`, `PUT /compensation/bands/:id`, `GET/POST /compensation/components`, `GET/POST /compensation/adjustments` (+ approve/reject/apply), `GET/POST /compensation/bonus-plans`, `GET/POST /compensation/bonuses` (+approve/reject), `GET/POST /compensation/benefits` (+assign), `GET/POST /compensation/reviews`, `GET/POST /compensation/budgets`, dashboard stats.
- [ ] Tabs: Dashboard | Estructuras | Grados | Bandas | Componentes | Ajustes | Bonos | Beneficios | Presupuestos
- [ ] Each tab: table + Nuevo modal + edit where endpoint exists

### Task 1e: Recruitment (`frontend/src/pages/recruitment.tsx`)
Endpoints: `GET/POST /recruitment/positions`, `PUT /recruitment/positions/:id`, `GET/POST /recruitment/postings`, `PUT /recruitment/postings/:id`, `POST /recruitment/postings/:id/publish`, `POST /recruitment/postings/:id/close`, `GET/POST /recruitment/candidates`, `PUT /recruitment/candidates/:id`, `GET /recruitment/candidates/search`, `GET/POST /recruitment/applications`, `POST /recruitment/applications/:id/move-stage|reject`, `GET/POST /recruitment/interviews`, `GET/POST /recruitment/offers`, `GET /recruitment/dashboard`, `GET /recruitment/funnel`.
- [ ] Tabs: Dashboard | Posiciones | Publicaciones | Candidatos | Postulaciones | Entrevistas | Ofertas
- [ ] Each tab: table + create/edit modal
- [ ] Postings: publish/close actions
- [ ] Applications: move stage, reject
- [ ] Dashboard: funnel chart (simple bars)

## Task 2: Complete employee detail (`frontend/src/pages/employee-detail.tsx`)
Endpoints: existing `/employees/:id`, plus `/employees/:id/contacts` (GET/PUT), `/employees/:id/addresses` (GET/PUT), `/employees/:id/emergency-contacts` (GET/PUT), `/employees/:id/history`, `/employees/:id/documents`, `/compensation/employee/:id` (GET/PUT via `GetEmployeeCompensation`/`SetEmployeeCompensation`), `/leave/requests?employee_id=`, `/attendance?employee_id=`, `/payroll/runs/:id/employees/:eid/...`, `/training/employees/:id/certificates`.
- [ ] Header card: photo placeholder, name, employee_number, status badge, position/department/branch/manager
- [ ] Tabs: Datos personales | Contactos | Direcciones | Emergencia | Historial | Documentos | Compensación | Licencias | Asistencia | Payroll | Capacitación
- [ ] Editable forms per tab using the PUT endpoints above
- [ ] Leave balance widget, attendance summary widget

## Task 3+: Remaining priorities

Each subsequent priority (3-15) is a separate plan once 1-2 land. Sketch per priority in `docs/superpowers/plans/`. Endpoints already inventoried in session for: leaves, attendance, payroll (+features), onboarding/offboarding, performance, benefits, expenses full, recruitment full, training full, /me/* portal, feed, surveys, settings/admin, reports.
