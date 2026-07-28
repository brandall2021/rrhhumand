# RRHHumand — Multi-Tenant Employee Management Platform

Sistema SaaS de Gestión de Recursos Humanos con arquitectura de monolito modular, diseñado para soporte multi-empresa y futura extracción a microservicios.

## Stack Tecnológico

| Componente | Tecnología |
|---|---|
| **Lenguaje** | Go 1.25.0 |
| **Framework HTTP** | Gin |
| **Base de datos** | PostgreSQL 16+ |
| **Pool de conexiones** | pgx/v5 |
| **Autenticación** | JWT (HS256) + Refresh Tokens |
| **Object Storage** | MinIO (S3-compatible) |
| **Containerización** | Docker + Docker Swarm |
| **Logging** | Zap (uber) |
| **Migraciones** | SQL raw (versionadas) |

## Arquitectura

```
rrhhumand/
├── cmd/api/                  # Entry point (main.go)
├── internal/                 # Módulos de negocio
│   ├── auth/                 # Autenticación + JWT
│   ├── attendance/           # Asistencia y control de jornada
│   ├── branches/             # Sucursales
│   ├── companies/            # Empresas (multi-tenant)
│   ├── departments/          # Departamentos
│   ├── document_categories/  # Categorías de documentos
│   ├── documents/            # Gestión documental + MinIO
│   ├── employees/            # Empleados
│   ├── feed/                 # Feed corporativo
│   ├── handlers/             # Health check
│   ├── leave/                # Vacaciones y licencias
│   ├── middleware/            # Auth, Tenant, RBAC, CORS
│   ├── models/               # Modelos compartidos
│   ├── organization/         # Árbol organizacional
│   ├── overtime/             # Horas extras y compensaciones
│   ├── payroll/              # Nómina y compensaciones económicas
│   ├── performance/          # Evaluación de desempeño
│   ├── permissions/          # Permisos RBAC
│   ├── positions/            # Puestos
│   ├── profile/              # Perfil de empleado
│   ├── recruitment/          # Reclutamiento y selección (ATS)
│   ├── roles/                # Roles del sistema
│   ├── scheduling/           # Turnos, horarios y planificación
│   ├── server/               # Router central
│   ├── surveys/              # Encuestas
│   └── users/                # Usuarios
├── migrations/               # Migraciones SQL versionadas
├── pkg/                      # Paquetes compartidos
│   ├── database/             # Conexión PostgreSQL
│   ├── logger/               # Logger (Zap)
│   ├── pagination/           # Paginación
│   ├── response/             # Respuestas HTTP estandarizadas
│   ├── security/             # Hashing bcrypt
│   └── validator/            # Validaciones
├── docs/                     # Documentación OpenAPI
└── docker/                   # Dockerfiles
```

## Principios de Diseño

- **Multi-tenant por empresa**: Todas las tablas de negocio tienen `company_id`. El tenant se extrae del JWT, nunca se confía del frontend.
- **IDs UUID**: Todas las entidades usan `gen_random_uuid()`.
- **Formato de respuesta estandarizado**:
  ```json
  { "success": true, "data": {} }
  { "success": false, "error": { "code": "...", "message": "..." } }
  ```
- **RBAC granular**: Middleware que verifica permisos en base de datos. SUPER_ADMIN y COMPANY_ADMIN bypass.
- **API versionada**: Todas las rutas bajo `/api/v1/`.

## Quick Start

### Prerequisitos
- Docker + Docker Swarm
- Go 1.25+

### 1. Levantar servicios
```bash
docker swarm init
docker stack deploy -c docker-compose.yml rrhh
```

### 2. Ejecutar migraciones
```bash
# Aplicar migraciones en orden
for dir in migrations/0000*/; do
  docker exec -i postgres_postgres.1.<id> psql -U postgres -d rrhhumand < "$dir/up.sql"
done
```

### 3. Correr la API
```bash
go run cmd/api/main.go
```

### 4. Login inicial
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@rrhh.com","password":"Admin123!"}'
```

---

## Fases Implementadas

### FASE 0 — Bootstrap

**Objetivo**: Estructura base del proyecto, configuración y herramientas fundamentales.

- Proyecto Go con estructura de monolito modular
- Conexión a PostgreSQL con pgxpool
- Dockerfile y docker-compose.yml
- Migraciones SQL versionadas (up/down)
- Health check (`/health`, `/ready`)
- Logger estructurado con Zap
- Manejo centralizado de errores
- Makefile con comandos útiles
- Configuración por variables de entorno (`config.go`)

**Tablas**: `users`, `roles`, `permissions`, `companies`

---

### FASE 1 — Autenticación + Multi-Tenant

**Objetivo**: Sistema de login, JWT, roles, permisos y aislamiento por empresa.

- Login/logout/refresh con JWT (HS256)
- Refresh tokens con expiración
- Hashing bcrypt para contraseñas
- Roles semilla: SUPER_ADMIN, COMPANY_ADMIN, HR_ADMIN, MANAGER, EMPLOYEE
- Permisos granulares por módulo
- Middleware de autenticación (`AuthMiddleware`)
- Middleware de tenant (`TenantMiddleware`)
- Middleware RBAC con `PermissionChecker` interface
- CORS configurado

**Tablas**: `refresh_tokens`, `user_companies`, `user_roles`, `role_permissions`

---

### FASE 2 — Empresas, Sucursales, Departamentos y Puestos

**Objetivo**: Estructura organizacional base.

- CRUD de empresas (multi-tenant)
- CRUD de sucursales concompany_id
- CRUD de departamentos
- CRUD de puestos
- Paginación, búsqueda, filtros y ordenamiento
- Aislamiento de tenant en todas las queries

**Tablas**: `companies`, `branches`, `departments`, `positions`

---

### FASE 3 — Empleados

**Objetivo**: Gestión completa del empleado con información personal, contactos, direcciones y historial.

- CRUD de empleados
- Upsert de contactos, direcciones y contactos de emergencia
- Historial automático de cambios (`employee_history`)
- GET/PUT de información del empleado
- Integración con organización (reporta a)

**Tablas**: `employees`, `employee_contacts`, `employee_addresses`, `employee_emergency_contacts`, `employee_history`

---

### FASE 4 — Árbol Organizacional

**Objetivo**: Estructura jerárquica recursiva de la empresa.

- Árbol recursivo via `manager_id`
- Detección de ciclos al asignar manager
- GET `/organization/tree` retorna el árbol completo

---

### FASE 5 — Perfil del Empleado

**Objetivo**: El empleado puede ver y editar su propio perfil.

- `GET /me/profile` — perfil completo con tenure calculado
- `PUT /me/profile` — actualizar perfil propio
- Cálculo automático de antigüedad (tenure)

---

### FASE 6 — Feed Corporativo

**Objetivo**: Red social interna para comunicación de la empresa.

- Crear, listar, actualizar y eliminar posts
- Comentarios en posts
- Reacciones (like, love, etc.)
- Menciones a usuarios
- Multimedia en posts (`post_media`)

**Tablas**: `posts`, `post_media`, `comments`, `reactions`, `mentions`

---

### FASE 7 — Encuestas

**Objetivo**: Sistema de encuestas con lifecycle, estadísticas y exportación CSV.

- CRUD de encuestas con lifecycle (draft → published → closed → archived)
- Preguntas con opciones de respuesta (5 tipos)
- Destinatarios por empresa, departamento o empleado
- Respuestas de empleados
- Estadísticas por pregunta y总体
- Exportación CSV
- 10 permisos granulares
- RBAC corregido con `PermissionChecker` interface

**Tablas**: `surveys`, `survey_questions`, `survey_options`, `survey_targets`, `survey_responses`, `survey_response_answers`

---

### FASE 8 — Gestión Documental

**Objetivo**: Sistema de documentos con versionado, permisos, shares y almacenamiento MinIO.

- Upload/download con MinIO (S3-compatible)
- Versionado de documentos
- Permisos por documento (view, edit, delete, share)
- Compartir documentos con usuarios o links públicos
- Tags y categorías
- Documentos que expiran (alertas)
- Archivar/restaurar/eliminar permanentemente
- Estadísticas de uso
- 10 permisos granulares

**Tablas**: `documents`, `document_versions`, `document_permissions`, `document_shares`, `document_share_links`, `document_tags`, `document_tag_relations`, `document_access_logs`, `document_categories`

---

### FASE 9 — Vacaciones y Licencias

**Objetivo**: Gestión completa de ausentismo con balance, aprobaciones y calendario.

- Tipos de licencia (vacaciones, enfermedad, personal, etc.)
- Políticas por empresa/departamento
- Feriados
- Balance de días con `SELECT ... FOR UPDATE`
- Solicitud de licencia conoverlap detection
- Workflow de aprobación (approve/reject/cancel)
- Calendario personal y de equipo
- Reportes
- 12 permisos granulares

**Tablas**: `leave_types`, `leave_policies`, `leave_holidays`, `leave_balances`, `leave_requests`, `leave_approvals`, `leave_request_history`

---

### FASE 10 — Asistencia y Control de Jornada

**Objetivo**: Registro de entrada/salida con geolocalización, breaks y dashboards.

- Clock-in / Clock-out
- Inicio/fin de breaks
- Geolocalización con fórmula Haversine
- Correcciones de asistencia con aprobación
- Dashboard general y de equipo
- Calendario de asistencia
- Políticas, ubicaciones y dispositivos
- Exportación CSV
- 12 permisos granulares

**Tablas**: `attendance_records`, `attendance_punches`, `attendance_breaks`, `attendance_policies`, `attendance_locations`, `attendance_devices`, `attendance_corrections`

---

### FASE 11 — Turnos, Horarios y Planificación

**Objetivo**: Gestión completa de horarios con prioridad de resolución, rotación y swap de turnos.

- Horarios base (semanales)
- Turnos individuales
- Asignación de horario/turno por empleado
- Rotación de turnos con templates
- Resolución de conflicto con cadena de prioridad: asignación directa > rotación > horario del empleado > horario default
- Generación de calendario
- Excepciones de horario
- Swap de turnos entre empleados con aprobación
- 12 permisos granulares

**Tablas**: `work_schedules`, `shifts`, `employee_schedules`, `employee_shifts`, `employee_rotations`, `schedule_calendar`, `schedule_exceptions`, `shift_swaps`, `shift_swap_approvals`, `rotation_templates`, `rotation_template_shifts`

---

### FASE 12 — Horas Extras y Compensaciones

**Objetivo**: Control de horas extras con detección automática, redondeo, límites y ledger de balance.

- Detección automática de overtime desde attendance
- Registro manual de horas extras
- Cálculo con configuración de redondeo (5, 10, 15, 30, 60 min)
- Límites diarios, semanales y mensuales
- Solicitud de horas extras con aprobación
- Compensaciones (tiempo libre o pago)
- Ledger de balance con `SELECT ... FOR UPDATE`
- Ajustes manuales de balance
- Dashboard de overtime
- 12 permisos granulares

**Tablas**: `overtime_records`, `overtime_requests`, `overtime_policies`, `overtime_balances`, `overtime_balance_transactions`, `compensation_requests`

---

### FASE 13 — Nómina, Compensaciones Económicas y Beneficios

**Objetivo**: Sistema de nómina con períodos, conceptos, cálculo automático, beneficios, bonos, anticipos y deducciones.

- Períodos de nómina con lifecycle (OPEN → CALCULATING → REVIEW → APPROVED → CLOSED)
- Conceptos de nómina (earning, deduction, benefit)
- Compensación por empleado (sueldo base)
- Cálculo automático de items de nómina
- Beneficios (health, dental, life, etc.)
- Bonos con workflow de aprobación
- Anticipos con aprobación
- Deducciones
- Ajustes por período
- Ledger contable
- Snapshot de nómina
- Dashboard por período
- 15 permisos granulares

**Tablas**: `payroll_periods`, `payroll_concepts`, `employee_compensations`, `payroll_items`, `benefits`, `employee_benefits`, `bonuses`, `advances`, `deductions`, `payroll_adjustments`, `payroll_ledger`, `payroll_snapshots`

---

### FASE 14 — Evaluación de Desempeño

**Objetivo**: Sistema 360° de evaluación con ciclos, plantillas, escalas, competencias, objetivos, KPIs, evaluadores, scoring engine y planes de mejora/desarrollo.

- **Ciclos de evaluación**: Períodos de evaluación con lifecycle
- **Plantillas**: Secciones con items configurables
- **Escalas de calificación**: Niveles personalizados
- **Competencias**: Definición de competencias por empresa
- **Objetivos**: Metas medibles por empleado/ciclo con pesos
- **KPIs**: Indicadores clave con targets
- **Evaluadores 360°**: SELF, MANAGER, PEER, HR
- **Evaluaciones**: Respuestas por sección/item con scoring
- **Feedback continuo**: Comentarios entre empleados
- **Evidencias**: Documentos adjuntos a evaluaciones
- **Resultados**: Score final ponderado con rating
- **Scoring Engine**: Cálculo configurable de pesos (objetivo, competencia, KPI, self, manager, peer, HR)
- **Planes de mejora**: Con acciones y seguimiento
- **Planes de desarrollo (IDP)**: Desarrollo individual con timeline
- **Dashboard**: Métricas y distribución de ratings
- 15 permisos granulares

**Tablas**: `performance_cycles`, `evaluation_templates`, `template_sections`, `template_section_items`, `rating_scales`, `rating_scale_levels`, `competencies`, `performance_scoring_rules`, `performance_objectives`, `performance_kpis`, `performance_evaluators`, `performance_evaluations`, `performance_evaluation_answers`, `performance_feedback`, `performance_evidence`, `performance_results`, `performance_improvement_plans`, `performance_improvement_actions`, `performance_development_plans`, `performance_development_actions`, `performance_audit_log`

---

### FASE 15 — Reclutamiento y Selección (ATS)

**Objetivo**: Applicant Tracking System integrado que gestiona el proceso completo desde la detección de necesidad hasta la contratación.

#### Pipeline de Reclutamiento
```
NECESIDAD DE PERSONAL
        ↓
SOLICITUD DE VACANTE
        ↓
APROBACIÓN
        ↓
PUBLICACIÓN
        ↓
POSTULACIÓN
        ↓
CANDIDATO
        ↓
SCREENING
        ↓
ENTREVISTAS
        ↓
EVALUACIONES
        ↓
SELECCIÓN
        ↓
OFERTA
        ↓
ACEPTACIÓN
        ↓
CONTRATACIÓN
        ↓
FASE 16 — ONBOARDING
        ↓
FASE 3 — EMPLEADO
```

#### Pipeline de Candidatos (Kanban)
```
NEW → SCREENING → PHONE_INTERVIEW → TECHNICAL_INTERVIEW → HR_INTERVIEW → ASSESSMENT → FINALIST → OFFER → HIRED
```
Con estados de salida: `REJECTED`, `WITHDRAWN`, `ON_HOLD`

#### Módulos Implementados

| # | Módulo | Descripción |
|---|--------|-------------|
| 1 | Solicitud de vacante | Crear requisiciones con puesto, departamento, salario, etc. |
| 2 | Aprobación de vacante | Workflow configurable (DRAFT → PENDING_APPROVAL → APPROVED → OPEN) |
| 3 | Puestos abiertos | Publicaciones de trabajo con requisitos y responsabilidades |
| 4 | Publicación de ofertas | Publicar en múltiples canales (abstracción JobBoardProvider) |
| 5 | Portal de candidatos | Endpoints públicos para candidatos |
| 6 | Candidatos | CRUD completo con deduplicación por email |
| 7 | CV/documentos | Almacenamiento de CVs con parsed_data JSONB |
| 8 | Screening | Preguntas configurables por vacante (BOOLEAN, TEXT, NUMBER, SELECT) |
| 9 | Pipeline | Movimiento de etapas con validación de transiciones |
| 10 | Entrevistas | Programación con entrevistador, tipo, fecha, modalidad |
| 11 | Evaluaciones | Evaluaciones técnicas con score y duración |
| 12 | Feedback | Score, comentarios y recomendación (STRONG_YES a STRONG_NO) |
| 13 | Matching | Engine de matching candidato-vacante (stub para IA) |
| 14 | Ranking | Ranking ponderado de candidatos (stub para IA) |
| 15 | Oferta laboral | Crear, enviar, aceptar/rechazar con salary (amount/currency/period) |
| 16 | Aprobación de contratación | Workflow de aprobación de ofertas |
| 17 | Conversión candidato → empleado | Integración con FASE 3 para crear empleado |
| 18 | Estadísticas | Dashboard: vacantes abiertas, candidatos, entrevistas, ofertas, hires |
| 19 | Notificaciones | Eventos del módulo (stub para FASE 16) |
| 20 | Auditoría | Log de acciones por usuario, candidato y empresa |

#### API REST — 36 Endpoints

**Vacantes**
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/recruitment/requisitions` | Listar requisiciones |
| POST | `/api/v1/recruitment/requisitions` | Crear requisición |
| GET | `/api/v1/recruitment/requisitions/:id` | Obtener requisición |
| PUT | `/api/v1/recruitment/requisitions/:id` | Actualizar requisición |
| POST | `/api/v1/recruitment/requisitions/:id/submit` | Enviar a aprobación |
| POST | `/api/v1/recruitment/requisitions/:id/approve` | Aprobar requisición |
| POST | `/api/v1/recruitment/requisitions/:id/open` | Abrir para contratación |
| POST | `/api/v1/recruitment/requisitions/:id/close` | Cerrar requisición |

**Publicaciones**
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/recruitment/jobs` | Listar publicaciones |
| POST | `/api/v1/recruitment/jobs` | Crear publicación |
| GET | `/api/v1/recruitment/jobs/:id` | Obtener publicación |
| POST | `/api/v1/recruitment/jobs/:id/publish` | Publicar |
| POST | `/api/v1/recruitment/jobs/:id/close` | Cerrar publicación |

**Candidatos**
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/recruitment/candidates` | Listar candidatos |
| POST | `/api/v1/recruitment/candidates` | Crear candidato |
| GET | `/api/v1/recruitment/candidates/:id` | Obtener candidato |
| PUT | `/api/v1/recruitment/candidates/:id` | Actualizar candidato |

**Postulaciones**
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/recruitment/applications` | Crear postulación |
| GET | `/api/v1/recruitment/applications/:id` | Obtener postulación |
| GET | `/api/v1/recruitment/applications` | Listar postulaciones |
| POST | `/api/v1/recruitment/applications/:id/stage` | Mover de etapa |
| POST | `/api/v1/recruitment/applications/:id/reject` | Rechazar |
| POST | `/api/v1/recruitment/applications/:id/withdraw` | Retirar postulación |
| GET | `/api/v1/recruitment/applications/:id/history` | Historial de etapas |
| POST | `/api/v1/recruitment/applications/:id/hire` | Convertir a empleado |

**Entrevistas**
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/recruitment/interviews` | Crear entrevista |
| GET | `/api/v1/recruitment/interviews` | Listar entrevistas |
| GET | `/api/v1/recruitment/interviews/:id` | Obtener entrevista |
| PUT | `/api/v1/recruitment/interviews/:id` | Actualizar entrevista |
| POST | `/api/v1/recruitment/interviews/:id/feedback` | Enviar feedback |
| GET | `/api/v1/recruitment/interviews/:id/feedback` | Listar feedback |

**Evaluaciones y Screening**
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/recruitment/assessments` | Crear evaluación técnica |
| GET | `/api/v1/recruitment/assessments/:id` | Listar evaluaciones |
| POST | `/api/v1/recruitment/screening` | Crear pregunta de screening |
| GET | `/api/v1/recruitment/screening/:id` | Listar preguntas |

**Ofertas**
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/recruitment/offers` | Crear oferta |
| GET | `/api/v1/recruitment/offers/:id` | Obtener oferta |
| POST | `/api/v1/recruitment/offers/:id/send` | Enviar oferta |
| POST | `/api/v1/recruitment/offers/:id/accept` | Aceptar oferta |
| POST | `/api/v1/recruitment/offers/:id/reject` | Rechazar oferta |

**Referidos y Dashboard**
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/recruitment/referrals` | Crear referido |
| GET | `/api/v1/recruitment/referrals` | Listar referidos |
| GET | `/api/v1/recruitment/dashboard` | Dashboard de reclutamiento |

#### Base de Datos — 17 Tablas

| Tabla | Descripción |
|-------|-------------|
| `job_requisitions` | Solicitudes de vacante |
| `approval_workflows` | Workflows de aprobación configurables |
| `approval_steps` | Pasos del workflow |
| `approval_instances` | Instancias de aprobación reales |
| `job_postings` | Publicaciones de trabajo |
| `candidates` | Candidatos (unique por email+empresa) |
| `applications` | Postulaciones candidato → publicación |
| `candidate_stage_history` | Historial de movimientos de etapa |
| `candidate_documents` | CVs y documentos con parsed_data JSONB |
| `screening_questions` | Preguntas de screening por publicación |
| `screening_answers` | Respuestas de screening |
| `interviews` | Entrevistas programadas |
| `interview_feedback` | Feedback de entrevistadores |
| `assessments` | Evaluaciones técnicas |
| `job_offers` | Ofertas laborales |
| `employee_referrals` | Programa de referidos |
| `recruitment_audit_log` | Auditoría de acciones |

#### RBAC — 9 Permisos

| Permiso | Descripción |
|---------|-------------|
| `recruitment.read` | Ver datos de reclutamiento |
| `recruitment.create_requisition` | Crear requisiciones |
| `recruitment.approve_requisition` | Aprobar requisiciones |
| `recruitment.manage_postings` | Gestionar publicaciones |
| `recruitment.manage_candidates` | Gestionar candidatos |
| `recruitment.conduct_interviews` | Realizar entrevistas |
| `recruitment.create_offers` | Crear ofertas |
| `recruitment.hire` | Contratar candidatos |
| `recruitment.analytics` | Ver analytics y dashboard |

#### Worker
- Cierra automáticamente publicaciones con `closing_at` vencido
- Expira ofertas con `response_deadline` vencido

#### Integraciones
- **FASE 3 (Empleados)**: Conversión automática candidato → empleado al aceptar oferta
- **FASE 16 (Onboarding)**: Evento `candidate.hired` para trigger de checklist
- **IA (FASE 19)**: Puntos de extensión para CV Parser, Matching y Ranking

---

## Resumen de Fases

| FASE | Nombre | Estado |
|------|--------|--------|
| 0 | Bootstrap | ✅ Completada |
| 1 | Auth + Multi-Tenant | ✅ Completada |
| 2 | Companies (Empresas, Sucursales, Departamentos, Puestos) | ✅ Completada |
| 3 | Employees (Empleados) | ✅ Completada |
| 4 | Organization Chart (Árbol Organizacional) | ✅ Completada |
| 5 | Employee Profile (Perfil) | ✅ Completada |
| 6 | Corporate Feed (Feed Corporativo) | ✅ Completada |
| 7 | Surveys (Encuestas) | ✅ Completada |
| 8 | Document Management (Gestión Documental) | ✅ Completada |
| 9 | Vacaciones y Licencias | ✅ Completada |
| 10 | Asistencia y Control de Jornada | ✅ Completada |
| 11 | Turnos, Horarios y Planificación | ✅ Completada |
| 12 | Horas Extras y Compensaciones | ✅ Completada |
| 13 | Nómina, Compensaciones Económicas y Beneficios | ✅ Completada |
| 14 | Evaluación de Desempeño | ✅ Completada |
| 15 | Reclutamiento y Selección (ATS) | ✅ Completada |
| 16 | Onboarding | ⏳ Pendiente |
| 17 | Notificaciones | ⏳ Pendiente |
| 18 | Reportes | ⏳ Pendiente |
| 19 | IA | ⏳ Pendiente |
| 20 | Integraciones | ⏳ Pendiente |

## Estadísticas del Proyecto

| Métrica | Valor |
|---------|-------|
| Fases completadas | 16 / 21 |
| Archivos Go | 80+ |
| Tablas PostgreSQL | 90+ |
| Endpoints API | 250+ |
| Permisos RBAC | 120+ |
| Líneas de código | 25,000+ |

## Licencia

Proyecto privado — RRHHumand Team
