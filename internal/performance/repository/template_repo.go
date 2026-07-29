package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type templateRepo struct {
	pool *pgxpool.Pool
}

func (r *templateRepo) Create(ctx context.Context, t *domain.PerformanceTemplate) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_templates (company_id, name, description, evaluation_type, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		t.CompanyID, t.Name, t.Description, t.EvaluationType, t.Active, now, now,
	).Scan(&t.ID)
}

func (r *templateRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceTemplate, error) {
	t := &domain.PerformanceTemplate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, evaluation_type, active, created_at, updated_at
		 FROM performance_templates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.EvaluationType, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *templateRepo) List(ctx context.Context, companyID string) ([]domain.PerformanceTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, evaluation_type, active, created_at, updated_at
		 FROM performance_templates WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []domain.PerformanceTemplate
	for rows.Next() {
		var t domain.PerformanceTemplate
		rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.EvaluationType, &t.Active, &t.CreatedAt, &t.UpdatedAt)
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *templateRepo) Update(ctx context.Context, t *domain.PerformanceTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_templates SET name=$3, description=$4, evaluation_type=$5, active=$6, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		t.CompanyID, t.ID, t.Name, t.Description, t.EvaluationType, t.Active)
	return err
}

func (r *templateRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_templates WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *templateRepo) CreateSection(ctx context.Context, s *domain.TemplateSection) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO template_sections (template_id, name, description, section_type, weight, sort_order, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		s.TemplateID, s.Name, s.Description, s.SectionType, s.Weight, s.SortOrder, s.Active,
	).Scan(&s.ID)
}

func (r *templateRepo) ListSectionsByTemplate(ctx context.Context, templateID string) ([]domain.TemplateSection, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, name, description, section_type, weight, sort_order, active
		 FROM template_sections WHERE template_id=$1 ORDER BY sort_order`, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []domain.TemplateSection
	for rows.Next() {
		var s domain.TemplateSection
		rows.Scan(&s.ID, &s.TemplateID, &s.Name, &s.Description, &s.SectionType, &s.Weight, &s.SortOrder, &s.Active)
		sections = append(sections, s)
	}
	return sections, nil
}

func (r *templateRepo) DeleteSectionsByTemplate(ctx context.Context, templateID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM template_sections WHERE template_id=$1`, templateID)
	return err
}

func (r *templateRepo) CreateQuestion(ctx context.Context, q *domain.TemplateQuestion) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO template_questions (template_id, section_id, question, question_type, required, weight, sort_order, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		q.TemplateID, q.SectionID, q.Question, q.QuestionType, q.Required, q.Weight, q.SortOrder, q.Active,
	).Scan(&q.ID)
}

func (r *templateRepo) ListQuestionsByTemplate(ctx context.Context, templateID string) ([]domain.TemplateQuestion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, section_id, question, question_type, required, weight, sort_order, active
		 FROM template_questions WHERE template_id=$1 ORDER BY sort_order`, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.TemplateQuestion
	for rows.Next() {
		var q domain.TemplateQuestion
		rows.Scan(&q.ID, &q.TemplateID, &q.SectionID, &q.Question, &q.QuestionType, &q.Required, &q.Weight, &q.SortOrder, &q.Active)
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *templateRepo) ListQuestionsBySection(ctx context.Context, sectionID string) ([]domain.TemplateQuestion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, section_id, question, question_type, required, weight, sort_order, active
		 FROM template_questions WHERE section_id=$1 ORDER BY sort_order`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.TemplateQuestion
	for rows.Next() {
		var q domain.TemplateQuestion
		rows.Scan(&q.ID, &q.TemplateID, &q.SectionID, &q.Question, &q.QuestionType, &q.Required, &q.Weight, &q.SortOrder, &q.Active)
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *templateRepo) DeleteQuestionsByTemplate(ctx context.Context, templateID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM template_questions WHERE template_id=$1`, templateID)
	return err
}

// ScaleRepository

type scaleRepo struct {
	pool *pgxpool.Pool
}

func (r *scaleRepo) Create(ctx context.Context, s *domain.RatingScale) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO rating_scales (company_id, name, min_value, max_value, description, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		s.CompanyID, s.Name, s.MinValue, s.MaxValue, s.Description, s.Active, now, now,
	).Scan(&s.ID)
}

func (r *scaleRepo) GetByID(ctx context.Context, companyID, id string) (*domain.RatingScale, error) {
	s := &domain.RatingScale{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, min_value, max_value, description, active, created_at, updated_at
		 FROM rating_scales WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.MinValue, &s.MaxValue, &s.Description, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *scaleRepo) List(ctx context.Context, companyID string) ([]domain.RatingScale, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, min_value, max_value, description, active, created_at, updated_at
		 FROM rating_scales WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scales []domain.RatingScale
	for rows.Next() {
		var s domain.RatingScale
		rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.MinValue, &s.MaxValue, &s.Description, &s.Active, &s.CreatedAt, &s.UpdatedAt)
		scales = append(scales, s)
	}
	return scales, nil
}

func (r *scaleRepo) Update(ctx context.Context, s *domain.RatingScale) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE rating_scales SET name=$3, min_value=$4, max_value=$5, description=$6, active=$7, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		s.CompanyID, s.ID, s.Name, s.MinValue, s.MaxValue, s.Description, s.Active)
	return err
}

func (r *scaleRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM rating_scales WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *scaleRepo) CreateLevel(ctx context.Context, l *domain.RatingScaleLevel) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO rating_scale_levels (scale_id, value, label, description, sort_order)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		l.ScaleID, l.Value, l.Label, l.Description, l.SortOrder,
	).Scan(&l.ID)
}

func (r *scaleRepo) ListLevelsByScale(ctx context.Context, scaleID string) ([]domain.RatingScaleLevel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, scale_id, value, label, description, sort_order
		 FROM rating_scale_levels WHERE scale_id=$1 ORDER BY sort_order`, scaleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []domain.RatingScaleLevel
	for rows.Next() {
		var l domain.RatingScaleLevel
		rows.Scan(&l.ID, &l.ScaleID, &l.Value, &l.Label, &l.Description, &l.SortOrder)
		levels = append(levels, l)
	}
	return levels, nil
}

func (r *scaleRepo) DeleteLevelsByScale(ctx context.Context, scaleID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM rating_scale_levels WHERE scale_id=$1`, scaleID)
	return err
}

// CompetencyRepository

type competencyRepo struct {
	pool *pgxpool.Pool
}

func (r *competencyRepo) Create(ctx context.Context, c *domain.Competency) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO competencies (company_id, name, description, category, competency_type, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		c.CompanyID, c.Name, c.Description, c.Category, c.CompetencyType, c.Active, now, now,
	).Scan(&c.ID)
}

func (r *competencyRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Competency, error) {
	c := &domain.Competency{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, category, competency_type, active, created_at, updated_at
		 FROM competencies WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Category, &c.CompetencyType, &c.Active, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *competencyRepo) List(ctx context.Context, filter domain.CompetencyFilter) ([]domain.Competency, error) {
	query := `SELECT id, company_id, name, description, category, competency_type, active, created_at, updated_at
		 FROM competencies WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category=$%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND competency_type=$%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Active != nil {
		query += fmt.Sprintf(" AND active=$%d", argIdx)
		args = append(args, *filter.Active)
		argIdx++
	}
	query += " ORDER BY name"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competencies []domain.Competency
	for rows.Next() {
		var c domain.Competency
		rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Category, &c.CompetencyType, &c.Active, &c.CreatedAt, &c.UpdatedAt)
		competencies = append(competencies, c)
	}
	return competencies, nil
}

func (r *competencyRepo) Update(ctx context.Context, c *domain.Competency) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE competencies SET name=$3, description=$4, category=$5, competency_type=$6, active=$7, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		c.CompanyID, c.ID, c.Name, c.Description, c.Category, c.CompetencyType, c.Active)
	return err
}

func (r *competencyRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM competencies WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *competencyRepo) CreateLevel(ctx context.Context, l *domain.CompetencyLevel) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO competency_levels (competency_id, level, label, description, sort_order)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		l.CompetencyID, l.Level, l.Label, l.Description, l.SortOrder,
	).Scan(&l.ID)
}

func (r *competencyRepo) ListLevelsByCompetency(ctx context.Context, competencyID string) ([]domain.CompetencyLevel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, competency_id, level, label, description, sort_order
		 FROM competency_levels WHERE competency_id=$1 ORDER BY sort_order`, competencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []domain.CompetencyLevel
	for rows.Next() {
		var l domain.CompetencyLevel
		rows.Scan(&l.ID, &l.CompetencyID, &l.Level, &l.Label, &l.Description, &l.SortOrder)
		levels = append(levels, l)
	}
	return levels, nil
}

func (r *competencyRepo) DeleteLevelsByCompetency(ctx context.Context, competencyID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM competency_levels WHERE competency_id=$1`, competencyID)
	return err
}

func (r *competencyRepo) UpsertPositionCompetency(ctx context.Context, pc *domain.PositionCompetency) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO position_competencies (company_id, position_id, competency_id, expected_level, weight)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (company_id, position_id, competency_id) DO UPDATE SET expected_level=$4, weight=$5
		 RETURNING id`,
		pc.CompanyID, pc.PositionID, pc.CompetencyID, pc.ExpectedLevel, pc.Weight,
	).Scan(&pc.ID)
}

func (r *competencyRepo) ListByPosition(ctx context.Context, companyID, positionID string) ([]domain.PositionCompetency, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, position_id, competency_id, expected_level, weight
		 FROM position_competencies WHERE company_id=$1 AND position_id=$2`, companyID, positionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pcs []domain.PositionCompetency
	for rows.Next() {
		var pc domain.PositionCompetency
		rows.Scan(&pc.ID, &pc.CompanyID, &pc.PositionID, &pc.CompetencyID, &pc.ExpectedLevel, &pc.Weight)
		pcs = append(pcs, pc)
	}
	return pcs, nil
}

func (r *competencyRepo) DeletePositionCompetency(ctx context.Context, companyID, positionID, competencyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM position_competencies WHERE company_id=$1 AND position_id=$2 AND competency_id=$3`,
		companyID, positionID, competencyID)
	return err
}

func (r *competencyRepo) UpsertCycleCompetency(ctx context.Context, cc *domain.CycleCompetency) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO cycle_competencies (cycle_id, employee_id, competency_id, expected_level, weight)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (cycle_id, employee_id, competency_id) DO UPDATE SET expected_level=$4, weight=$5
		 RETURNING id`,
		cc.CycleID, cc.EmployeeID, cc.CompetencyID, cc.ExpectedLevel, cc.Weight,
	).Scan(&cc.ID)
}

func (r *competencyRepo) ListByCycleEmployee(ctx context.Context, cycleID, employeeID string) ([]domain.CycleCompetency, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, cycle_id, employee_id, competency_id, expected_level, weight
		 FROM cycle_competencies WHERE cycle_id=$1 AND employee_id=$2`, cycleID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ccs []domain.CycleCompetency
	for rows.Next() {
		var cc domain.CycleCompetency
		rows.Scan(&cc.ID, &cc.CycleID, &cc.EmployeeID, &cc.CompetencyID, &cc.ExpectedLevel, &cc.Weight)
		ccs = append(ccs, cc)
	}
	return ccs, nil
}

func (r *competencyRepo) DeleteCycleCompetency(ctx context.Context, cycleID, employeeID, competencyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM cycle_competencies WHERE cycle_id=$1 AND employee_id=$2 AND competency_id=$3`,
		cycleID, employeeID, competencyID)
	return err
}

func NewTemplateRepository(pool *pgxpool.Pool) TemplateRepository {
	return &templateRepo{pool: pool}
}

func NewScaleRepository(pool *pgxpool.Pool) ScaleRepository {
	return &scaleRepo{pool: pool}
}

func NewCompetencyRepository(pool *pgxpool.Pool) CompetencyRepository {
	return &competencyRepo{pool: pool}
}
