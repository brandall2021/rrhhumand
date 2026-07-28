package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Cycles
func (r *Repository) CreateCycle(ctx context.Context, companyID, createdBy string, req *CreateCycleRequest) (*PerformanceCycle, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	c := &PerformanceCycle{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_cycles (company_id, name, description, start_date, end_date, evaluation_deadline, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, company_id, name, description, start_date, end_date, evaluation_deadline, status, created_by, created_at, updated_at`,
		companyID, req.Name, req.Description, startDate, endDate, req.EvaluationDeadline, createdBy,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.EvaluationDeadline,
		&c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) GetCycle(ctx context.Context, companyID, id string) (*PerformanceCycle, error) {
	c := &PerformanceCycle{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, start_date, end_date, evaluation_deadline, status, created_by, created_at, updated_at
		 FROM performance_cycles WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.EvaluationDeadline,
		&c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) ListCycles(ctx context.Context, companyID string) ([]PerformanceCycle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, start_date, end_date, evaluation_deadline, status, created_by, created_at, updated_at
		 FROM performance_cycles WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var cycles []PerformanceCycle
	for rows.Next() {
		var c PerformanceCycle
		rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.EvaluationDeadline,
			&c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		cycles = append(cycles, c)
	}
	return cycles, nil
}

func (r *Repository) UpdateCycle(ctx context.Context, companyID, id string, req *UpdateCycleRequest) (*PerformanceCycle, error) {
	c := &PerformanceCycle{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_cycles SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 start_date=COALESCE($5,start_date), end_date=COALESCE($6,end_date),
		 evaluation_deadline=COALESCE($7,evaluation_deadline), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, start_date, end_date, evaluation_deadline, status, created_by, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.StartDate, req.EndDate, req.EvaluationDeadline,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.EvaluationDeadline,
		&c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) UpdateCycleStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_cycles SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

// Templates
func (r *Repository) CreateTemplate(ctx context.Context, companyID string, req *CreateTemplateRequest) (*EvaluationTemplate, error) {
	t := &EvaluationTemplate{}
	isDefault := false
	if req.IsDefault != nil { isDefault = *req.IsDefault }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO evaluation_templates (company_id, name, description, is_default)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, name, description, is_default, active, created_at`,
		companyID, req.Name, req.Description, isDefault,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.IsDefault, &t.Active, &t.CreatedAt)
	return t, err
}

func (r *Repository) CreateTemplateSection(ctx context.Context, templateID string, req *CreateSectionRequest) (*TemplateSection, error) {
	s := &TemplateSection{}
	weight := 0.0
	if req.Weight != nil { weight = *req.Weight }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO template_sections (template_id, name, description, section_type, weight)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, template_id, name, description, section_type, weight, sort_order, active`,
		templateID, req.Name, req.Description, req.SectionType, weight,
	).Scan(&s.ID, &s.TemplateID, &s.Name, &s.Description, &s.SectionType, &s.Weight, &s.SortOrder, &s.Active)
	return s, err
}

func (r *Repository) CreateTemplateSectionItem(ctx context.Context, sectionID string, req *CreateSectionItemRequest) (*TemplateSectionItem, error) {
	it := &TemplateSectionItem{}
	itemType := "TEXT"
	weight := 0.0
	if req.ItemType != nil { itemType = *req.ItemType }
	if req.Weight != nil { weight = *req.Weight }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO template_section_items (section_id, name, description, item_type, weight)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, section_id, name, description, item_type, weight, sort_order`,
		sectionID, req.Name, req.Description, itemType, weight,
	).Scan(&it.ID, &it.SectionID, &it.Name, &it.Description, &it.ItemType, &it.Weight, &it.SortOrder)
	return it, err
}

func (r *Repository) ListTemplates(ctx context.Context, companyID string) ([]EvaluationTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, is_default, active, created_at
		 FROM evaluation_templates WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var templates []EvaluationTemplate
	for rows.Next() {
		var t EvaluationTemplate
		rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.IsDefault, &t.Active, &t.CreatedAt)
		templates = append(templates, t)
	}
	return templates, nil
}

// Scales
func (r *Repository) CreateScale(ctx context.Context, companyID string, req *CreateScaleRequest) (*RatingScale, error) {
	s := &RatingScale{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO rating_scales (company_id, name, min_value, max_value, description)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, company_id, name, min_value, max_value, description, active, created_at`,
		companyID, req.Name, req.MinValue, req.MaxValue, req.Description,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.MinValue, &s.MaxValue, &s.Description, &s.Active, &s.CreatedAt)
	return s, err
}

func (r *Repository) CreateScaleLevel(ctx context.Context, scaleID string, req *CreateScaleLevelRequest) (*RatingScaleLevel, error) {
	l := &RatingScaleLevel{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO rating_scale_levels (scale_id, value, label, description)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, scale_id, value, label, description, sort_order`,
		scaleID, req.Value, req.Label, req.Description,
	).Scan(&l.ID, &l.ScaleID, &l.Value, &l.Label, &l.Description, &l.SortOrder)
	return l, err
}

func (r *Repository) ListScales(ctx context.Context, companyID string) ([]RatingScale, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, min_value, max_value, description, active, created_at
		 FROM rating_scales WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var scales []RatingScale
	for rows.Next() {
		var s RatingScale
		rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.MinValue, &s.MaxValue, &s.Description, &s.Active, &s.CreatedAt)
		scales = append(scales, s)
	}
	return scales, nil
}

// Competencies
func (r *Repository) CreateCompetency(ctx context.Context, companyID string, req *CreateCompetencyRequest) (*Competency, error) {
	c := &Competency{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO competencies (company_id, name, description, category)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, name, description, category, active, created_at`,
		companyID, req.Name, req.Description, req.Category,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Category, &c.Active, &c.CreatedAt)
	return c, err
}

func (r *Repository) ListCompetencies(ctx context.Context, companyID string) ([]Competency, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, category, active, created_at
		 FROM competencies WHERE company_id=$1 AND active=TRUE ORDER BY name`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var competencies []Competency
	for rows.Next() {
		var c Competency
		rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Category, &c.Active, &c.CreatedAt)
		competencies = append(competencies, c)
	}
	return competencies, nil
}

func (r *Repository) UpdateCompetency(ctx context.Context, companyID, id string, req *UpdateCompetencyRequest) (*Competency, error) {
	c := &Competency{}
	err := r.pool.QueryRow(ctx,
		`UPDATE competencies SET name=COALESCE($3,name), description=COALESCE($4,description),
		 category=COALESCE($5,category), active=COALESCE($6,active)
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, category, active, created_at`,
		companyID, id, req.Name, req.Description, req.Category, req.Active,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Category, &c.Active, &c.CreatedAt)
	return c, err
}

// Objectives
func (r *Repository) CreateObjective(ctx context.Context, companyID, createdBy string, req *CreateObjectiveRequest) (*PerformanceObjective, error) {
	o := &PerformanceObjective{}
	weight := 0.0
	if req.Weight != nil { weight = *req.Weight }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_objectives (company_id, employee_id, cycle_id, title, description, metric, target_value, unit, weight, start_date, due_date, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 RETURNING id, company_id, employee_id, cycle_id, title, description, metric, target_value, current_value, unit, weight, start_date, due_date, status, created_by, created_at, updated_at`,
		companyID, req.EmployeeID, req.CycleID, req.Title, req.Description, req.Metric, req.TargetValue, req.Unit, weight, req.StartDate, req.DueDate, createdBy,
	).Scan(&o.ID, &o.CompanyID, &o.EmployeeID, &o.CycleID, &o.Title, &o.Description, &o.Metric, &o.TargetValue, &o.CurrentValue, &o.Unit, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repository) GetObjective(ctx context.Context, companyID, id string) (*PerformanceObjective, error) {
	o := &PerformanceObjective{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, cycle_id, title, description, metric, target_value, current_value, unit, weight, start_date, due_date, status, created_by, created_at, updated_at
		 FROM performance_objectives WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&o.ID, &o.CompanyID, &o.EmployeeID, &o.CycleID, &o.Title, &o.Description, &o.Metric, &o.TargetValue, &o.CurrentValue, &o.Unit, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repository) ListObjectives(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceObjective, error) {
	query := `SELECT id, company_id, employee_id, cycle_id, title, description, metric, target_value, current_value, unit, weight, start_date, due_date, status, created_by, created_at, updated_at
		 FROM performance_objectives WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.CycleID != "" {
		query += fmt.Sprintf(" AND cycle_id=$%d", argIdx)
		args = append(args, filters.CycleID)
		argIdx++
	}
	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var objectives []PerformanceObjective
	for rows.Next() {
		var o PerformanceObjective
		rows.Scan(&o.ID, &o.CompanyID, &o.EmployeeID, &o.CycleID, &o.Title, &o.Description, &o.Metric, &o.TargetValue, &o.CurrentValue, &o.Unit, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
		objectives = append(objectives, o)
	}
	return objectives, nil
}

func (r *Repository) UpdateObjective(ctx context.Context, companyID, id string, req *UpdateObjectiveRequest) (*PerformanceObjective, error) {
	o := &PerformanceObjective{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_objectives SET
		 title=COALESCE($3,title), description=COALESCE($4,description), metric=COALESCE($5,metric),
		 target_value=COALESCE($6,target_value), current_value=COALESCE($7,current_value),
		 unit=COALESCE($8,unit), weight=COALESCE($9,weight), status=COALESCE($10,status),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, employee_id, cycle_id, title, description, metric, target_value, current_value, unit, weight, start_date, due_date, status, created_by, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.Metric, req.TargetValue, req.CurrentValue, req.Unit, req.Weight, req.Status,
	).Scan(&o.ID, &o.CompanyID, &o.EmployeeID, &o.CycleID, &o.Title, &o.Description, &o.Metric, &o.TargetValue, &o.CurrentValue, &o.Unit, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

// KPIs
func (r *Repository) CreateKPI(ctx context.Context, companyID, createdBy string, req *CreateKPIRequest) (*PerformanceKPI, error) {
	k := &PerformanceKPI{}
	weight := 0.0
	if req.Weight != nil { weight = *req.Weight }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_kpis (company_id, employee_id, cycle_id, name, description, target_value, unit, weight, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, company_id, employee_id, cycle_id, name, description, target_value, current_value, unit, weight, status, created_by, created_at, updated_at`,
		companyID, req.EmployeeID, req.CycleID, req.Name, req.Description, req.TargetValue, req.Unit, weight, createdBy,
	).Scan(&k.ID, &k.CompanyID, &k.EmployeeID, &k.CycleID, &k.Name, &k.Description, &k.TargetValue, &k.CurrentValue, &k.Unit, &k.Weight, &k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt)
	return k, err
}

func (r *Repository) ListKPIs(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceKPI, error) {
	query := `SELECT id, company_id, employee_id, cycle_id, name, description, target_value, current_value, unit, weight, status, created_by, created_at, updated_at
		 FROM performance_kpis WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.CycleID != "" {
		query += fmt.Sprintf(" AND cycle_id=$%d", argIdx)
		args = append(args, filters.CycleID)
		argIdx++
	}
	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var kpis []PerformanceKPI
	for rows.Next() {
		var k PerformanceKPI
		rows.Scan(&k.ID, &k.CompanyID, &k.EmployeeID, &k.CycleID, &k.Name, &k.Description, &k.TargetValue, &k.CurrentValue, &k.Unit, &k.Weight, &k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt)
		kpis = append(kpis, k)
	}
	return kpis, nil
}

func (r *Repository) UpdateKPI(ctx context.Context, companyID, id string, req *UpdateKPIRequest) (*PerformanceKPI, error) {
	k := &PerformanceKPI{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_kpis SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 target_value=COALESCE($5,target_value), current_value=COALESCE($6,current_value),
		 unit=COALESCE($7,unit), weight=COALESCE($8,weight), status=COALESCE($9,status),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, employee_id, cycle_id, name, description, target_value, current_value, unit, weight, status, created_by, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.TargetValue, req.CurrentValue, req.Unit, req.Weight, req.Status,
	).Scan(&k.ID, &k.CompanyID, &k.EmployeeID, &k.CycleID, &k.Name, &k.Description, &k.TargetValue, &k.CurrentValue, &k.Unit, &k.Weight, &k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt)
	return k, err
}

// Evaluators
func (r *Repository) AssignEvaluators(ctx context.Context, companyID string, req *AssignEvaluatorsRequest) ([]PerformanceEvaluator, error) {
	var evaluators []PerformanceEvaluator
	for _, ev := range req.Evaluators {
		e := &PerformanceEvaluator{}
		err := r.pool.QueryRow(ctx,
			`INSERT INTO performance_evaluators (company_id, cycle_id, employee_id, evaluator_id, evaluator_type)
			 VALUES ($1,$2,$3,$4,$5)
			 ON CONFLICT (cycle_id, employee_id, evaluator_id, evaluator_type) DO UPDATE SET status='PENDING'
			 RETURNING id, company_id, cycle_id, employee_id, evaluator_id, evaluator_type, status, assigned_at`,
			companyID, req.CycleID, req.EmployeeID, ev.EvaluatorID, ev.EvaluatorType,
		).Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.EmployeeID, &e.EvaluatorID, &e.EvaluatorType, &e.Status, &e.AssignedAt)
		if err != nil { return nil, err }
		evaluators = append(evaluators, *e)
	}
	return evaluators, nil
}

func (r *Repository) ListEvaluators(ctx context.Context, companyID, cycleID string) ([]PerformanceEvaluator, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pe.id, pe.company_id, pe.cycle_id, pe.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pe.evaluator_id, COALESCE(e2.first_name||' '||e2.last_name,''), pe.evaluator_type, pe.status, pe.assigned_at, pe.completed_at
		 FROM performance_evaluators pe
		 LEFT JOIN employees e ON pe.employee_id=e.id
		 LEFT JOIN employees e2 ON pe.evaluator_id=e2.id
		 WHERE pe.company_id=$1 AND pe.cycle_id=$2
		 ORDER BY pe.employee_id, pe.evaluator_type`, companyID, cycleID)
	if err != nil { return nil, err }
	defer rows.Close()

	var evaluators []PerformanceEvaluator
	for rows.Next() {
		var ev PerformanceEvaluator
		rows.Scan(&ev.ID, &ev.CompanyID, &ev.CycleID, &ev.EmployeeID, &ev.EmployeeName,
			&ev.EvaluatorID, &ev.EvaluatorName, &ev.EvaluatorType, &ev.Status, &ev.AssignedAt, &ev.CompletedAt)
		evaluators = append(evaluators, ev)
	}
	return evaluators, nil
}

func (r *Repository) UpdateEvaluatorStatus(ctx context.Context, evaluatorID, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluators SET status=$2, completed_at=CASE WHEN $2='COMPLETED' THEN NOW() ELSE NULL END WHERE id=$1`,
		evaluatorID, status)
	return err
}

// Evaluations
func (r *Repository) CreateEvaluation(ctx context.Context, companyID string, req *CreateEvaluationRequest) (*PerformanceEvaluation, error) {
	e := &PerformanceEvaluation{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_evaluations (company_id, cycle_id, employee_id, evaluator_id, evaluator_type, template_id, comments)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, company_id, cycle_id, employee_id, evaluator_id, evaluator_type, template_id, status, overall_score, comments, created_at, updated_at`,
		companyID, req.CycleID, req.EmployeeID, req.EvaluatorID, req.EvaluatorType, req.TemplateID, req.Comments,
	).Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.EmployeeID, &e.EvaluatorID, &e.EvaluatorType, &e.TemplateID,
		&e.Status, &e.OverallScore, &e.Comments, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *Repository) GetEvaluation(ctx context.Context, companyID, id string) (*PerformanceEvaluation, error) {
	e := &PerformanceEvaluation{}
	err := r.pool.QueryRow(ctx,
		`SELECT pe.id, pe.company_id, pe.cycle_id, pc.name, pe.employee_id, COALESCE(em.first_name||' '||em.last_name,''),
		 pe.evaluator_id, COALESCE(ev.first_name||' '||ev.last_name,''), pe.evaluator_type, pe.template_id,
		 pe.status, pe.overall_score, pe.comments, pe.submitted_at, pe.created_at, pe.updated_at
		 FROM performance_evaluations pe
		 LEFT JOIN performance_cycles pc ON pe.cycle_id=pc.id
		 LEFT JOIN employees em ON pe.employee_id=em.id
		 LEFT JOIN employees ev ON pe.evaluator_id=ev.id
		 WHERE pe.company_id=$1 AND pe.id=$2`, companyID, id,
	).Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.CycleName, &e.EmployeeID, &e.EmployeeName,
		&e.EvaluatorID, &e.EvaluatorName, &e.EvaluatorType, &e.TemplateID,
		&e.Status, &e.OverallScore, &e.Comments, &e.SubmittedAt, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *Repository) ListEvaluations(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceEvaluation, error) {
	query := `SELECT pe.id, pe.company_id, pe.cycle_id, pc.name, pe.employee_id, COALESCE(em.first_name||' '||em.last_name,''),
		 pe.evaluator_id, COALESCE(ev.first_name||' '||ev.last_name,''), pe.evaluator_type, pe.template_id,
		 pe.status, pe.overall_score, pe.comments, pe.submitted_at, pe.created_at, pe.updated_at
		 FROM performance_evaluations pe
		 LEFT JOIN performance_cycles pc ON pe.cycle_id=pc.id
		 LEFT JOIN employees em ON pe.employee_id=em.id
		 LEFT JOIN employees ev ON pe.evaluator_id=ev.id
		 WHERE pe.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.CycleID != "" {
		query += fmt.Sprintf(" AND pe.cycle_id=$%d", argIdx)
		args = append(args, filters.CycleID)
		argIdx++
	}
	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND pe.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND pe.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY pe.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var evaluations []PerformanceEvaluation
	for rows.Next() {
		var e PerformanceEvaluation
		rows.Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.CycleName, &e.EmployeeID, &e.EmployeeName,
			&e.EvaluatorID, &e.EvaluatorName, &e.EvaluatorType, &e.TemplateID,
			&e.Status, &e.OverallScore, &e.Comments, &e.SubmittedAt, &e.CreatedAt, &e.UpdatedAt)
		evaluations = append(evaluations, e)
	}
	return evaluations, nil
}

func (r *Repository) UpdateEvaluationStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET status=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) SubmitEvaluation(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET status='SUBMITTED', submitted_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2 AND status='DRAFT'`,
		companyID, id)
	return err
}

func (r *Repository) SetEvaluationScore(ctx context.Context, companyID, id string, score float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET overall_score=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, score)
	return err
}

// Answers
func (r *Repository) CreateAnswer(ctx context.Context, evaluationID string, req *CreateAnswerRequest) (*EvaluationAnswer, error) {
	a := &EvaluationAnswer{}
	itemType := "TEXT"
	weight := 0.0
	if req.ItemType != nil { itemType = *req.ItemType }
	if req.Weight != nil { weight = *req.Weight }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_evaluation_answers (evaluation_id, section_name, item_name, item_type, score, value, comments, weight)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, evaluation_id, section_name, item_name, item_type, score, value, comments, weight, created_at`,
		evaluationID, req.SectionName, req.ItemName, itemType, req.Score, req.Value, req.Comments, weight,
	).Scan(&a.ID, &a.EvaluationID, &a.SectionName, &a.ItemName, &a.ItemType, &a.Score, &a.Value, &a.Comments, &a.Weight, &a.CreatedAt)
	return a, err
}

func (r *Repository) DeleteAnswersByEvaluation(ctx context.Context, evaluationID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_evaluation_answers WHERE evaluation_id=$1`, evaluationID)
	return err
}

func (r *Repository) ListAnswersByEvaluation(ctx context.Context, evaluationID string) ([]EvaluationAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, evaluation_id, section_name, item_name, item_type, score, value, comments, weight, created_at
		 FROM performance_evaluation_answers WHERE evaluation_id=$1 ORDER BY created_at`, evaluationID)
	if err != nil { return nil, err }
	defer rows.Close()

	var answers []EvaluationAnswer
	for rows.Next() {
		var a EvaluationAnswer
		rows.Scan(&a.ID, &a.EvaluationID, &a.SectionName, &a.ItemName, &a.ItemType, &a.Score, &a.Value, &a.Comments, &a.Weight, &a.CreatedAt)
		answers = append(answers, a)
	}
	return answers, nil
}

// Feedback
func (r *Repository) CreateFeedback(ctx context.Context, companyID, fromUserID string, req *CreateFeedbackRequest) (*PerformanceFeedback, error) {
	f := &PerformanceFeedback{}
	isPrivate := false
	visibleToEmployee := true
	if req.IsPrivate != nil { isPrivate = *req.IsPrivate }
	if req.VisibleToEmployee != nil { visibleToEmployee = *req.VisibleToEmployee }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_feedback (company_id, employee_id, cycle_id, from_user_id, feedback_type, message, is_private, visible_to_employee)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, employee_id, cycle_id, from_user_id, feedback_type, message, is_private, visible_to_employee, created_at`,
		companyID, req.EmployeeID, req.CycleID, fromUserID, req.FeedbackType, req.Message, isPrivate, visibleToEmployee,
	).Scan(&f.ID, &f.CompanyID, &f.EmployeeID, &f.CycleID, &f.FromUserID, &f.FeedbackType, &f.Message, &f.IsPrivate, &f.VisibleToEmployee, &f.CreatedAt)
	return f, err
}

func (r *Repository) ListFeedback(ctx context.Context, companyID, employeeID string) ([]PerformanceFeedback, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pf.id, pf.company_id, pf.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pf.cycle_id, pf.from_user_id, COALESCE(e2.first_name||' '||e2.last_name,''),
		 pf.feedback_type, pf.message, pf.is_private, pf.visible_to_employee, pf.created_at
		 FROM performance_feedback pf
		 LEFT JOIN employees e ON pf.employee_id=e.id
		 LEFT JOIN employees e2 ON pf.from_user_id=e2.id
		 WHERE pf.company_id=$1 AND pf.employee_id=$2
		 ORDER BY pf.created_at DESC`, companyID, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	var feedbacks []PerformanceFeedback
	for rows.Next() {
		var f PerformanceFeedback
		rows.Scan(&f.ID, &f.CompanyID, &f.EmployeeID, &f.EmployeeName,
			&f.CycleID, &f.FromUserID, &f.FromUserName,
			&f.FeedbackType, &f.Message, &f.IsPrivate, &f.VisibleToEmployee, &f.CreatedAt)
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

// Evidence
func (r *Repository) CreateEvidence(ctx context.Context, companyID, evaluationID, createdBy string, req *CreateEvidenceRequest) (*PerformanceEvidence, error) {
	evidenceType := "COMMENT"
	if req.EvidenceType != "" { evidenceType = req.EvidenceType }
	ev := &PerformanceEvidence{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_evidence (company_id, evaluation_id, title, description, evidence_type, url, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, company_id, evaluation_id, title, description, evidence_type, url, created_by, created_at`,
		companyID, evaluationID, req.Title, req.Description, evidenceType, req.URL, createdBy,
	).Scan(&ev.ID, &ev.CompanyID, &ev.EvaluationID, &ev.Title, &ev.Description, &ev.EvidenceType, &ev.URL, &ev.CreatedBy, &ev.CreatedAt)
	return ev, err
}

func (r *Repository) ListEvidenceByEvaluation(ctx context.Context, evaluationID string) ([]PerformanceEvidence, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, evaluation_id, title, description, evidence_type, storage_provider, storage_key, file_name, mime_type, size_bytes, url, created_by, created_at
		 FROM performance_evidence WHERE evaluation_id=$1 ORDER BY created_at`, evaluationID)
	if err != nil { return nil, err }
	defer rows.Close()

	var evidence []PerformanceEvidence
	for rows.Next() {
		var ev PerformanceEvidence
		rows.Scan(&ev.ID, &ev.CompanyID, &ev.EvaluationID, &ev.Title, &ev.Description, &ev.EvidenceType,
			&ev.StorageProvider, &ev.StorageKey, &ev.FileName, &ev.MimeType, &ev.SizeBytes, &ev.URL, &ev.CreatedBy, &ev.CreatedAt)
		evidence = append(evidence, ev)
	}
	return evidence, nil
}

// Results
func (r *Repository) UpsertResult(ctx context.Context, companyID, cycleID, employeeID string, score *PerformanceScore) (*PerformanceResult, error) {
	res := &PerformanceResult{}
	ratings, _ := json.Marshal(map[string]interface{}{
		"objective": score.ObjectiveScore,
		"competency": score.CompetencyScore,
		"kpi": score.KPIScore,
		"self": score.SelfScore,
		"manager": score.ManagerScore,
		"peer": score.PeerScore,
		"hr": score.HRScore,
	})
	_ = ratings
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_results (company_id, cycle_id, employee_id, objective_score, competency_score, kpi_score,
		 self_score, manager_score, peer_score, hr_score, final_score, rating, rating_label, strengths, areas_to_improve)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 ON CONFLICT (cycle_id, employee_id) DO UPDATE SET
		 objective_score=$4, competency_score=$5, kpi_score=$6, self_score=$7, manager_score=$8, peer_score=$9,
		 hr_score=$10, final_score=$11, rating=$12, rating_label=$13, strengths=$14, areas_to_improve=$15, calculated_at=NOW(), updated_at=NOW()
		 RETURNING id, company_id, cycle_id, employee_id, objective_score, competency_score, kpi_score,
		 self_score, manager_score, peer_score, hr_score, final_score, rating, rating_label, strengths, areas_to_improve, calculated_at, created_at, updated_at`,
		companyID, cycleID, employeeID, score.ObjectiveScore, score.CompetencyScore, score.KPIScore,
		score.SelfScore, score.ManagerScore, score.PeerScore, score.HRScore, score.FinalScore,
		score.Rating, score.RatingLabel, score.Strengths, score.AreasToImprove,
	).Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.EmployeeID, &res.ObjectiveScore, &res.CompetencyScore, &res.KPIScore,
		&res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore, &res.FinalScore, &res.Rating, &res.RatingLabel,
		&res.Strengths, &res.AreasToImprove, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
	return res, err
}

func (r *Repository) GetResult(ctx context.Context, companyID, cycleID, employeeID string) (*PerformanceResult, error) {
	res := &PerformanceResult{}
	err := r.pool.QueryRow(ctx,
		`SELECT pr.id, pr.company_id, pr.cycle_id, pc.name, pr.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pr.objective_score, pr.competency_score, pr.kpi_score, pr.self_score, pr.manager_score, pr.peer_score, pr.hr_score,
		 pr.final_score, pr.rating, pr.rating_label, pr.strengths, pr.areas_to_improve, pr.calculated_at, pr.created_at, pr.updated_at
		 FROM performance_results pr
		 LEFT JOIN performance_cycles pc ON pr.cycle_id=pc.id
		 LEFT JOIN employees e ON pr.employee_id=e.id
		 WHERE pr.company_id=$1 AND pr.cycle_id=$2 AND pr.employee_id=$3`, companyID, cycleID, employeeID,
	).Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.CycleName, &res.EmployeeID, &res.EmployeeName,
		&res.ObjectiveScore, &res.CompetencyScore, &res.KPIScore, &res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore,
		&res.FinalScore, &res.Rating, &res.RatingLabel, &res.Strengths, &res.AreasToImprove, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
	return res, err
}

func (r *Repository) ListResults(ctx context.Context, companyID, cycleID string) ([]PerformanceResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pr.id, pr.company_id, pr.cycle_id, pc.name, pr.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pr.objective_score, pr.competency_score, pr.kpi_score, pr.self_score, pr.manager_score, pr.peer_score, pr.hr_score,
		 pr.final_score, pr.rating, pr.rating_label, pr.strengths, pr.areas_to_improve, pr.calculated_at, pr.created_at, pr.updated_at
		 FROM performance_results pr
		 LEFT JOIN performance_cycles pc ON pr.cycle_id=pc.id
		 LEFT JOIN employees e ON pr.employee_id=e.id
		 WHERE pr.company_id=$1 AND pr.cycle_id=$2
		 ORDER BY pr.final_score DESC NULLS LAST`, companyID, cycleID)
	if err != nil { return nil, err }
	defer rows.Close()

	var results []PerformanceResult
	for rows.Next() {
		var res PerformanceResult
		rows.Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.CycleName, &res.EmployeeID, &res.EmployeeName,
			&res.ObjectiveScore, &res.CompetencyScore, &res.KPIScore, &res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore,
			&res.FinalScore, &res.Rating, &res.RatingLabel, &res.Strengths, &res.AreasToImprove, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
		results = append(results, res)
	}
	return results, nil
}

// Scoring rules
func (r *Repository) GetScoringRules(ctx context.Context, companyID string) (*ScoringRule, error) {
	sr := &ScoringRule{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, objective_weight, competency_weight, kpi_weight,
		 self_eval_weight, manager_weight, peer_weight, hr_weight, active, created_at
		 FROM performance_scoring_rules WHERE company_id=$1 AND active=TRUE ORDER BY created_at DESC LIMIT 1`, companyID,
	).Scan(&sr.ID, &sr.CompanyID, &sr.CycleID, &sr.ObjectiveWeight, &sr.CompetencyWeight, &sr.KPIWeight,
		&sr.SelfEvalWeight, &sr.ManagerWeight, &sr.PeerWeight, &sr.HRWeight, &sr.Active, &sr.CreatedAt)
	if err != nil {
		return &ScoringRule{
			CompanyID:        companyID,
			ObjectiveWeight:  40,
			CompetencyWeight: 30,
			KPIWeight:        20,
			SelfEvalWeight:   10,
			ManagerWeight:    60,
			PeerWeight:       20,
			HRWeight:         10,
			Active:           true,
		}, nil
	}
	return sr, err
}

func (r *Repository) UpdateScoringRules(ctx context.Context, companyID string, req *UpdateScoringRulesRequest) (*ScoringRule, error) {
	sr := &ScoringRule{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_scoring_rules SET
		 objective_weight=COALESCE($3,objective_weight), competency_weight=COALESCE($4,competency_weight),
		 kpi_weight=COALESCE($5,kpi_weight), self_eval_weight=COALESCE($6,self_eval_weight),
		 manager_weight=COALESCE($7,manager_weight), peer_weight=COALESCE($8,peer_weight),
		 hr_weight=COALESCE($9,hr_weight)
		 WHERE company_id=$1 AND active=TRUE
		 RETURNING id, company_id, cycle_id, objective_weight, competency_weight, kpi_weight,
		 self_eval_weight, manager_weight, peer_weight, hr_weight, active, created_at`,
		companyID, req.ObjectiveWeight, req.CompetencyWeight, req.KPIWeight, req.SelfEvalWeight,
		req.ManagerWeight, req.PeerWeight, req.HRWeight,
	).Scan(&sr.ID, &sr.CompanyID, &sr.CycleID, &sr.ObjectiveWeight, &sr.CompetencyWeight, &sr.KPIWeight,
		&sr.SelfEvalWeight, &sr.ManagerWeight, &sr.PeerWeight, &sr.HRWeight, &sr.Active, &sr.CreatedAt)
	return sr, err
}

// Improvement Plans
func (r *Repository) CreateImprovementPlan(ctx context.Context, companyID, createdBy string, req *CreateImprovementPlanRequest) (*ImprovementPlan, error) {
	p := &ImprovementPlan{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_improvement_plans (company_id, employee_id, cycle_id, result_id, title, problem_description, objective, responsible_id, due_date, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, company_id, employee_id, cycle_id, result_id, title, problem_description, objective, responsible_id, due_date, status, outcome, created_by, created_at, updated_at`,
		companyID, req.EmployeeID, req.CycleID, req.ResultID, req.Title, req.ProblemDescription, req.Objective, req.ResponsibleID, req.DueDate, createdBy,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.ResultID, &p.Title, &p.ProblemDescription, &p.Objective, &p.ResponsibleID, &p.DueDate, &p.Status, &p.Outcome, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) CreatePlanAction(ctx context.Context, planID string, req *CreatePlanActionRequest) (*ImprovementAction, error) {
	a := &ImprovementAction{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_improvement_actions (plan_id, title, description, due_date)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, plan_id, title, description, due_date, status, completed_at, created_at`,
		planID, req.Title, req.Description, req.DueDate,
	).Scan(&a.ID, &a.PlanID, &a.Title, &a.Description, &a.DueDate, &a.Status, &a.CompletedAt, &a.CreatedAt)
	return a, err
}

func (r *Repository) ListImprovementPlans(ctx context.Context, companyID string, filters PerformanceFilters) ([]ImprovementPlan, error) {
	query := `SELECT pip.id, pip.company_id, pip.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pip.cycle_id, pip.result_id, pip.title, pip.problem_description, pip.objective, pip.responsible_id, pip.due_date,
		 pip.status, pip.outcome, pip.created_by, pip.created_at, pip.updated_at
		 FROM performance_improvement_plans pip
		 LEFT JOIN employees e ON pip.employee_id=e.id
		 WHERE pip.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND pip.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND pip.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY pip.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var plans []ImprovementPlan
	for rows.Next() {
		var p ImprovementPlan
		rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.EmployeeName,
			&p.CycleID, &p.ResultID, &p.Title, &p.ProblemDescription, &p.Objective, &p.ResponsibleID, &p.DueDate,
			&p.Status, &p.Outcome, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *Repository) GetImprovementPlan(ctx context.Context, companyID, id string) (*ImprovementPlan, error) {
	p := &ImprovementPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT pip.id, pip.company_id, pip.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pip.cycle_id, pip.result_id, pip.title, pip.problem_description, pip.objective, pip.responsible_id, pip.due_date,
		 pip.status, pip.outcome, pip.created_by, pip.created_at, pip.updated_at
		 FROM performance_improvement_plans pip
		 LEFT JOIN employees e ON pip.employee_id=e.id
		 WHERE pip.company_id=$1 AND pip.id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.EmployeeName,
		&p.CycleID, &p.ResultID, &p.Title, &p.ProblemDescription, &p.Objective, &p.ResponsibleID, &p.DueDate,
		&p.Status, &p.Outcome, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) UpdateImprovementPlan(ctx context.Context, companyID, id string, req *UpdateImprovementPlanRequest) (*ImprovementPlan, error) {
	p := &ImprovementPlan{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_improvement_plans SET
		 title=COALESCE($3,title), problem_description=COALESCE($4,problem_description),
		 objective=COALESCE($5,objective), responsible_id=COALESCE($6,responsible_id),
		 due_date=COALESCE($7,due_date), status=COALESCE($8,status), outcome=COALESCE($9,outcome),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, employee_id, cycle_id, result_id, title, problem_description, objective, responsible_id, due_date, status, outcome, created_by, created_at, updated_at`,
		companyID, id, req.Title, req.ProblemDescription, req.Objective, req.ResponsibleID, req.DueDate, req.Status, req.Outcome,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.ResultID, &p.Title, &p.ProblemDescription, &p.Objective, &p.ResponsibleID, &p.DueDate, &p.Status, &p.Outcome, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// Development Plans
func (r *Repository) CreateDevelopmentPlan(ctx context.Context, companyID, createdBy string, req *CreateDevelopmentPlanRequest) (*DevelopmentPlan, error) {
	p := &DevelopmentPlan{}
	timeline := 12
	if req.TimelineMonths != nil { timeline = *req.TimelineMonths }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_development_plans (company_id, employee_id, cycle_id, title, description, career_goal, timeline_months, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, employee_id, cycle_id, title, description, career_goal, timeline_months, status, created_by, created_at, updated_at`,
		companyID, req.EmployeeID, req.CycleID, req.Title, req.Description, req.CareerGoal, timeline, createdBy,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.Title, &p.Description, &p.CareerGoal, &p.TimelineMonths, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) CreateDevAction(ctx context.Context, planID string, req *CreateDevActionRequest) (*DevelopmentAction, error) {
	actionType := "TRAINING"
	if req.ActionType != nil { actionType = *req.ActionType }
	a := &DevelopmentAction{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_development_actions (plan_id, title, description, action_type, due_date)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, plan_id, title, description, action_type, due_date, status, completed_at, created_at`,
		planID, req.Title, req.Description, actionType, req.DueDate,
	).Scan(&a.ID, &a.PlanID, &a.Title, &a.Description, &a.ActionType, &a.DueDate, &a.Status, &a.CompletedAt, &a.CreatedAt)
	return a, err
}

func (r *Repository) ListDevelopmentPlans(ctx context.Context, companyID, employeeID string) ([]DevelopmentPlan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pdp.id, pdp.company_id, pdp.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pdp.cycle_id, pdp.title, pdp.description, pdp.career_goal, pdp.timeline_months, pdp.status, pdp.created_by, pdp.created_at, pdp.updated_at
		 FROM performance_development_plans pdp
		 LEFT JOIN employees e ON pdp.employee_id=e.id
		 WHERE pdp.company_id=$1 AND pdp.employee_id=$2
		 ORDER BY pdp.created_at DESC`, companyID, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	var plans []DevelopmentPlan
	for rows.Next() {
		var p DevelopmentPlan
		rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.EmployeeName,
			&p.CycleID, &p.Title, &p.Description, &p.CareerGoal, &p.TimelineMonths, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *Repository) GetDevelopmentPlan(ctx context.Context, companyID, id string) (*DevelopmentPlan, error) {
	p := &DevelopmentPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT pdp.id, pdp.company_id, pdp.employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 pdp.cycle_id, pdp.title, pdp.description, pdp.career_goal, pdp.timeline_months, pdp.status, pdp.created_by, pdp.created_at, pdp.updated_at
		 FROM performance_development_plans pdp
		 LEFT JOIN employees e ON pdp.employee_id=e.id
		 WHERE pdp.company_id=$1 AND pdp.id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.EmployeeName,
		&p.CycleID, &p.Title, &p.Description, &p.CareerGoal, &p.TimelineMonths, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) UpdateDevelopmentPlan(ctx context.Context, companyID, id string, req *UpdateDevelopmentPlanRequest) (*DevelopmentPlan, error) {
	p := &DevelopmentPlan{}
	err := r.pool.QueryRow(ctx,
		`UPDATE performance_development_plans SET
		 title=COALESCE($3,title), description=COALESCE($4,description),
		 career_goal=COALESCE($5,career_goal), timeline_months=COALESCE($6,timeline_months),
		 status=COALESCE($7,status), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, employee_id, cycle_id, title, description, career_goal, timeline_months, status, created_by, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.CareerGoal, req.TimelineMonths, req.Status,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.Title, &p.Description, &p.CareerGoal, &p.TimelineMonths, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// Audit
func (r *Repository) CreateAuditLog(ctx context.Context, companyID, userID, entityType, entityID, action string, oldVal, newVal []byte, ipAddress string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO performance_audit_log (company_id, user_id, entity_type, entity_id, action, old_value, new_value, ip_address)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		companyID, userID, entityType, entityID, action, oldVal, newVal, ipAddress)
	return err
}

// Dashboard
func (r *Repository) GetDashboard(ctx context.Context, companyID string) (*PerformanceDashboard, error) {
	dash := &PerformanceDashboard{}

	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_cycles WHERE company_id=$1`, companyID).Scan(&dash.TotalCycles)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_cycles WHERE company_id=$1 AND status IN ('OPEN','IN_PROGRESS')`, companyID).Scan(&dash.ActiveCycles)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1`, companyID).Scan(&dash.TotalEvaluations)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1 AND status='SUBMITTED'`, companyID).Scan(&dash.CompletedEvaluations)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1 AND status IN ('DRAFT','PENDING')`, companyID).Scan(&dash.PendingEvaluations)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(final_score),0) FROM performance_results WHERE company_id=$1`, companyID).Scan(&dash.AverageScore)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_objectives WHERE company_id=$1`, companyID).Scan(&dash.TotalObjectives)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_objectives WHERE company_id=$1 AND status='COMPLETED'`, companyID).Scan(&dash.CompletedObjectives)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_kpis WHERE company_id=$1`, companyID).Scan(&dash.TotalKPIs)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_feedback WHERE company_id=$1`, companyID).Scan(&dash.TotalFeedback)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_improvement_plans WHERE company_id=$1`, companyID).Scan(&dash.TotalImprovementPlans)

	rows, err := r.pool.Query(ctx,
		`SELECT COALESCE(rating,'UNRATED'), COUNT(*) FROM performance_results WHERE company_id=$1 GROUP BY rating ORDER BY COUNT(*) DESC`, companyID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rc RatingCount
			rows.Scan(&rc.Rating, &rc.Count)
			dash.RatingDistribution = append(dash.RatingDistribution, rc)
		}
	}

	return dash, nil
}

// Helpers
func (r *Repository) GetObjectiveWeightTotal(ctx context.Context, cycleID, employeeID string) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(weight),0) FROM performance_objectives WHERE cycle_id=$1 AND employee_id=$2 AND status != 'CANCELLED'`,
		cycleID, employeeID).Scan(&total)
	return total, err
}

func (r *Repository) GetAnswersScoreForType(ctx context.Context, evaluationID, evaluatorType string) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(score * weight / 100),0) FROM performance_evaluation_answers WHERE evaluation_id=$1 AND score IS NOT NULL`,
		evaluationID).Scan(&total)
	return total, err
}

func (r *Repository) GetEvaluationsByType(ctx context.Context, companyID, cycleID, employeeID, evaluatorType string) ([]PerformanceEvaluation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluator_type, status, overall_score, comments, created_at, updated_at
		 FROM performance_evaluations
		 WHERE company_id=$1 AND cycle_id=$2 AND employee_id=$3 AND evaluator_type=$4
		 ORDER BY created_at`, companyID, cycleID, employeeID, evaluatorType)
	if err != nil { return nil, err }
	defer rows.Close()

	var evaluations []PerformanceEvaluation
	for rows.Next() {
		var e PerformanceEvaluation
		rows.Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.EmployeeID, &e.EvaluatorID, &e.EvaluatorType, &e.Status, &e.OverallScore, &e.Comments, &e.CreatedAt, &e.UpdatedAt)
		evaluations = append(evaluations, e)
	}
	return evaluations, nil
}
