# RRHHumand — Multi-Tenant Employee Management Platform

Sistema SaaS de Gestión de Recursos Humanos con arquitectura de monolito modular, diseñado para soporte multi-empresa y futura extracción a microservicios.

## Stack Tecnológico
Demo company: TechCorp Argentina (slug: techcorp-ar)

     ┌─────────────────────┬──────────────────────────────────┬──────────────┬──────────────────┐
     │Usuario              │Email                             │Password      │Rol               │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │Admin                │admin@techcorp.com                │admin123      │COMPANY_ADMIN     │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │HR                   │hr@techcorp.com                   │hr123         │HR_ADMIN          │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │Carlos García        │carlos.garcia@techcorp.com        │emp123        │MANAGER           │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │María López          │maria.lopez@techcorp.com          │emp123        │EMPLOYEE          │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │Ana Martínez         │ana.martinez@techcorp.com         │emp123        │EMPLOYEE          │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │Pedro Fernández      │pedro.fernandez@techcorp.com      │emp123        │EMPLOYEE          │
     ├─────────────────────┼──────────────────────────────────┼──────────────┼──────────────────┤
     │Lucía Torres         │lucia.torres@techcorp.com         │emp123        │EMPLOYEE          │
     └─────────────────────┴──────────────────────────────────┴──────────────┴──────────────────┘  
| Componente | Tecnología |
|---|---|
| **Lenguaje** | Go 1.25.0 |
| **Framework HTTP** | Gin |
| **Base de datos** | PostgreSQL 16+ |
| **Pool de conexiones** | pgx/v5 |
| **Autenticación** | JWT (HS256) + Refresh Tokens |
| **Object Storage** | MinIO (S3-compatible) |
| **Containerización** | Docker |
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
│   ├── performance/          # Performance Management (DDD)
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

## Quick Start (Desarrollo Local)

### Prerequisitos
- Go 1.25+
- Docker

### 1. Iniciar base de datos
```bash
docker compose up -d db
```

### 2. Aplicar migraciones
```bash
for dir in migrations/*/; do
  PGPASSWORD=Mfcd62!!Mfcd62!! psql -h localhost -p 5433 -U postgres -d rrhhumand -f "${dir}up.sql"
done
```

### 3. Correr la API
```bash
cp .env.example .env
go run cmd/api/main.go
```

### 4. Login inicial
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@rrhh.com","password":"Admin123!"}'
```

---

## Deploy con Dokploy

### Requisitos

- [Dokploy](https://dokploy.com) v0.x+ instalado en un VPS
- Repositorio Git accesible desde Dokploy
- PostgreSQL 16+ como servicio externo o add-on de Dokploy
- (Opcional) MinIO para almacenamiento de documentos

### Variables de Entorno

Crear un `.env` en Dokploy con todas las variables definidas en `.env.example`.

**Obligatorias:**

| Variable | Descripción | Ejemplo |
|---|---|---|
| `DATABASE_HOST` | Host de PostgreSQL | `postgres.internal` |
| `DATABASE_PORT` | Puerto de PostgreSQL | `5432` |
| `DATABASE_USER` | Usuario BD | `postgres` |
| `DATABASE_PASSWORD` | Contraseña BD | `securepassword` |
| `DATABASE_NAME` | Nombre BD | `rrhhumand` |
| `JWT_SECRET` | Clave secreta JWT (mín. 32 caracteres) | `your-256-bit-secret` |
| `APP_ENV` | `production` | `production` |
| `APP_DEBUG` | `false` en producción | `false` |

**Opcionales:**

| Variable | Default | Descripción |
|---|---|---|
| `APP_PORT` | `8080` | Puerto del contenedor |
| `DATABASE_SSLMODE` | `disable` | Modo SSL PostgreSQL |
| `DATABASE_MAX_OPEN_CONNS` | `25` | Conexiones abiertas máx. |
| `DATABASE_MAX_IDLE_CONNS` | `5` | Conexiones inactivas máx. |
| `JWT_EXPIRATION` | `15m` | Expiración token acceso |
| `JWT_REFRESH_EXPIRATION` | `7d` | Expiración refresh token |
| `MINIO_ENDPOINT` | `localhost:9000` | Endpoint MinIO |
| `MINIO_ACCESS_KEY` | `minioadmin` | Access Key MinIO |
| `MINIO_SECRET_KEY` | `minioadmin` | Secret Key MinIO |
| `MINIO_BUCKET` | `rrhhumand` | Bucket MinIO |
| `MINIO_REGION` | `us-east-1` | Región MinIO |
| `MINIO_USE_SSL` | `false` | SSL para MinIO |
| `MAX_UPLOAD_SIZE_MB` | `25` | Tamaño máx. upload (MB) |
| `LOG_LEVEL` | `info` | Nivel de log |
| `LOG_FORMAT` | `json` | Formato de log (`json` o `console`) |

### Configuración en Dokploy

1. **Crear proyecto** en el dashboard de Dokploy.
2. **Conectar repositorio** Git que contiene el código.
3. **Build**: Seleccionar **Dockerfile** (Dokploy lo detecta automáticamente).
   - Puerto: `8080`
4. **Variables de entorno**: Copiar todas las variables de la tabla superior.
5. **Base de datos**: Agregar un servicio PostgreSQL o conectar uno externo. Configurar `DATABASE_*` apuntando al host correcto.
6. **Volumen** (opcional para persistencia de migrations): No es necesario, las migraciones se ejecutan en cada deploy.
7. **Deploy**: Dokploy construye la imagen con `docker build` usando el `Dockerfile` y despliega el contenedor.

### Comportamiento del Entrypoint

El entrypoint (`docker/entrypoint.sh`) hace lo siguiente al iniciar el contenedor:

1. Espera a que PostgreSQL esté disponible (`pg_isready`).
2. Ejecuta todas las migraciones SQL en orden numérico (`migrations/*/up.sql`).
3. Inicia el servidor Go.

### Notas para Producción

- **JWT_SECRET**: Usar una clave de al menos 256 bits. Generar con:
  ```bash
  openssl rand -base64 32
  ```
- **MinIO**: Si no se usa MinIO, el servidor inicia igual pero las funciones de documentos no estarán disponibles.
- **Migraciones**: Son **idempotentes** — se ejecutan en cada deploy sin riesgo de duplicar datos.
- **Workers**: El servidor inicia automáticamente workers en goroutines para payroll, benefits, expenses y recruitment.
- **Base de datos separada**: Para producción, usar una instancia PostgreSQL administrada (RDS, Cloud SQL, Supabase, etc.) en lugar del contenedor local.

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

- **Nota**: Reemplazada por la **FASE 24 (DDD)** — mismo alcance funcional, arquitectura rediseñada con Domain-Driven Design.

---

### FASE 24 — Performance Management (DDD)

**Objetivo**: Refactorización completa del módulo de desempeño con arquitectura DDD en capas, aislando dominio, repositorios, servicios de aplicación, motores de scoring/workflow e integraciones.

#### Arquitectura del Módulo

```
internal/performance/
├── domain/                 # 10 entidades de negocio + errores + filtros
│   ├── performance_cycle.go
│   ├── template.go
│   ├── evaluation.go
│   ├── objective.go
│   ├── competency.go
│   ├── feedback.go
│   ├── evidence.go
│   ├── calibration.go
│   ├── plans.go
│   └── errors.go
├── repository/             # 8 repositorios (interfaces + PostgreSQL)
│   ├── interfaces.go
│   ├── cycle_repo.go
│   ├── template_repo.go
│   ├── evaluation_repo.go
│   ├── objective_repo.go
│   ├── feedback_repo.go
│   ├── plan_repo.go
│   └── misc_repo.go
├── application/            # Servicios + motores
│   ├── service/            # CycleService, ObjectiveService, EvaluationService, etc.
│   ├── scoring/            # Scorer con adapters + DefaultRatingScale
│   ├── workflow/           # StateMachine para ciclo, evaluación y plan
│   └── ai/                 # Stub de integración con IA
├── http/                   # Handler REST + DTOs
│   ├── handler.go          # 87 endpoints
│   └── dto.go              # Request/Response types
└── integration/            # Adapters para integraciones externas
```

**Scoring Engine**: Calcula puntuación final ponderada combinando objetivos (60%), competencias (40%), y evaluaciones multi-fuente (SELF, MANAGER, PEER, HR) con pesos configurables. Incluye `DefaultRatingScale` con 5 niveles.

**Workflow Engine**: Máquinas de estado para ciclo (DRAFT → OPEN → IN_PROGRESS → REVIEW → CALIBRATION → CLOSED), evaluación (DRAFT → SUBMITTED → APPROVED/LOCKED) y plan de mejora (DRAFT → ACTIVE → COMPLETED).

**Tablas**: `performance_cycles`, `evaluation_templates`, `template_sections`, `template_questions`, `rating_scales`, `rating_scale_levels`, `competencies`, `competency_levels`, `position_competencies`, `cycle_competencies`, `performance_objectives`, `objective_key_results`, `performance_participants`, `performance_evaluations`, `evaluation_answers`, `objective_evaluations`, `competency_evaluations`, `performance_reviews`, `performance_feedback`, `performance_recognitions`, `performance_checkins`, `calibration_sessions`, `calibration_items`, `performance_evidence`, `improvement_plans`, `improvement_plan_actions`, `development_plans`, `development_plan_actions`, `performance_results`, `performance_dashboard`, `performance_audit_log`, `outbox_events`

---

### FASE 22 — Reclutamiento y Selección (ATS) — Completo

**Objetivo**: Applicant Tracking System completo con arquitectura en capas, IA integrada y pipeline full lifecycle desde la necesidad hasta la contratación.

#### Arquitectura del Módulo

```
internal/recruitment/
├── domain/          # 18 entidades de negocio
├── repository/      # 16 repositorios PostgreSQL (pgx)
├── application/     # 15 servicios de aplicación
├── engine/          # 4 motores de negocio (scoring, matching, ofertas, workflows)
├── ai/              # 4 archivos (cliente, prompts, config, provider)
├── integration/     # 3 adapters + cliente base + calendario + email
├── http/            # 17 handlers REST
├── worker.go        # Worker de expiración
└── router.go        # Registro de rutas
```

#### Pipeline de Reclutamiento
```
REQUISITION → APPROVAL → POSITION → POSTING → APPLICATION
→ SCREENING → INTERVIEW → ASSESSMENT → OFFER → HIRING → ONBOARDING
```

#### Tablas — 32+

| Tabla | Descripción |
|-------|-------------|
| `job_requisitions` | Solicitudes de vacante con skills, approval workflow |
| `requisition_skills` | Skills requeridas por requisición |
| `requisition_approvals` | Aprobaciones de requisición |
| `positions` | Puestos abiertos (independiente de requisición) |
| `position_skills` | Skills requeridas por puesto |
| `job_postings` | Publicaciones con multiplataforma |
| `posting_board_posts` | Posts en boards externos (LinkedIn, Indeed, etc.) |
| `posting_screening_questions` | Preguntas de screening por posting |
| `candidates` | Candidatos con datos personales, CV, skills |
| `candidate_education` | Historial educativo |
| `candidate_experience` | Experiencia laboral |
| `candidate_skills` | Habilidades del candidato |
| `candidate_certifications` | Certificaciones |
| `candidate_languages` | Idiomas |
| `candidate_documents` | Documentos adjuntos |
| `candidate_parsed_data` | Datos parseados del CV (JSONB) |
| `applications` | Postulaciones candidato → publicación |
| `application_stage_history` | Historial de movimientos de etapa |
| `application_ratings` | Calificaciones de la postulación |
| `application_notes` | Notas internas |
| `application_screening_answers` | Respuestas a screening |
| `interviews` | Entrevistas programadas |
| `interview_panel_members` | Miembros del panel |
| `interview_feedback` | Feedback con score y recomendación |
| `interview_feedback_questions` | Respuestas a preguntas de feedback |
| `assessments` | Evaluaciones técnicas |
| `assessment_sections` | Secciones de evaluación |
| `assessment_results` | Resultados detallados |
| `offers` | Ofertas laborales |
| `offer_approvals` | Aprobaciones de oferta |
| `offer_negotiations` | Negociaciones salariales |
| `offer_documents` | Documentos de oferta |
| `hiring_processes` | Procesos de contratación post-aceptación |
| `hiring_process_tasks` | Tareas de contratación (background check, médico, docs) |
| `hiring_process_documents` | Documentos del proceso |
| `talent_pools` | Pool de talento |
| `talent_pool_candidates` | Candidatos en pool |
| `referral_rewards` | Recompensas por referidos |
| `recruitment_sources` | Fuentes de reclutamiento |
| `recruitment_stages` | Etapas del pipeline |
| `recruitment_stage_transitions` | Transiciones permitidas |
| `rejection_reasons` | Motivos de rechazo |
| `scoring_models` | Modelos de scoring |
| `scoring_criteria` | Criterios de scoring |
| `matching_results` | Resultados de matching |
| `email_templates` | Plantillas de email |
| `email_log` | Log de envíos |
| `workflows` | Workflows configurables |
| `workflow_stages` | Etapas de workflow |
| `workflow_rules` | Reglas de workflow |
| `dashboard_cache` | Cache de dashboard |

#### IA Integrada

| Componente | Descripción |
|------------|-------------|
| `ai/client.go` | Cliente LLM genérico (OpenAI-compatible) |
| `ai/prompts.go` | Prompts para CV parsing, matching y ranking |
| `ai/config.go` | Configuración de proveedor IA |
| `ai/provider.go` | Provider interface extensible |

| Módulo AI | Descripción |
|--------|-------------|
| CV Parsing | Extracción estructurada de skills, experiencia, educación desde CV |
| Matching semántico | Match candidate-position usando embeddings semánticos |
| Ranking IA | Ranking con justificación en lenguaje natural |
| Scoring | Scoring engine con criterios ponderados configurables |

#### Workers

| Worker | Descripción |
|--------|-------------|
| Auto-close postings | Cierra publicaciones con `closing_at` vencido |
| Expire offers | Expira ofertas con `response_deadline` vencido |

#### Integraciones

| Integración | Tipo | Descripción |
|-------------|------|-------------|
| Job Boards | Adapter | LinkedIn, Indeed, Computrabajo (extensible) |
| AI Providers | Provider | OpenAI, Claude, Gemini (extensible) |
| Calendar | Sync | Google Calendar, Outlook |
| Email | SMTP | Notificaciones con templates configurables |
| FASE 3 (Empleados) | Event | Conversión candidato → empleado |
| FASE 16 (Onboarding) | Event | Trigger de checklist al contratar |

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
| 14 | Evaluación de Desempeño | ✅ Completada (reemplazada por FASE 24) |
| 15 | Reclutamiento y Selección (ATS) | ⏳ Reemplazada por FASE 22 |
| 16 | Onboarding | ✅ Completada |
| 17 | Notificaciones | ⏳ Pendiente |
| 18 | Reportes | ⏳ Pendiente |
| 19 | IA General | ⏳ Pendiente |
| 19b | Payroll Features (Recibos, ARC, Libros, Bancos, Contab., Reportes) | ✅ Completada |
| 20 | Benefits & Total Rewards | ✅ Completada |
| 21 | Expenses & Travel | ✅ Completada |
| 22 | Reclutamiento y Selección (ATS Completo) | ✅ Completada |
| 23 | Onboarding & Offboarding (DDD) | ✅ Completada |
| 24 | Performance Management (DDD) | ✅ Completada |

## Estadísticas del Proyecto

| Métrica | Valor |
|---------|-------|
| Fases completadas | 23 / 24 |
| Archivos Go | 120+ |
| Tablas PostgreSQL | 150+ |
| Endpoints API | 350+ |
| Permisos RBAC | 150+ |
| Líneas de código | 40,000+ |

## Licencia

Proyecto privado — RRHHumand Team
