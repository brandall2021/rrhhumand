package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type calibrationRepo struct {
	pool *pgxpool.Pool
}

func (r *calibrationRepo) CreateSession(ctx context.Context, s *domain.CalibrationSession) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO calibration_sessions (company_id, cycle_id, name, description, status, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		s.CompanyID, s.CycleID, s.Name, s.Description, s.Status, s.CreatedBy, now, now,
	).Scan(&s.ID)
}

func (r *calibrationRepo) GetSessionByID(ctx context.Context, companyID, id string) (*domain.CalibrationSession, error) {
	s := &domain.CalibrationSession{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, name, description, status, started_at, completed_at, created_by, created_at, updated_at
		 FROM calibration_sessions WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&s.ID, &s.CompanyID, &s.CycleID, &s.Name, &s.Description, &s.Status, &s.StartedAt, &s.CompletedAt, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *calibrationRepo) ListSessionsByCycle(ctx context.Context, companyID, cycleID string) ([]domain.CalibrationSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, name, description, status, started_at, completed_at, created_by, created_at, updated_at
		 FROM calibration_sessions WHERE company_id=$1 AND cycle_id=$2 ORDER BY created_at`, companyID, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.CalibrationSession
	for rows.Next() {
		var s domain.CalibrationSession
		rows.Scan(&s.ID, &s.CompanyID, &s.CycleID, &s.Name, &s.Description, &s.Status, &s.StartedAt, &s.CompletedAt, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *calibrationRepo) UpdateSession(ctx context.Context, s *domain.CalibrationSession) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calibration_sessions SET name=$3, description=$4, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		s.CompanyID, s.ID, s.Name, s.Description)
	return err
}

func (r *calibrationRepo) UpdateSessionStatus(ctx context.Context, companyID, id string, status domain.CalibrationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calibration_sessions SET status=$3,
		 started_at=CASE WHEN $3='IN_PROGRESS' AND started_at IS NULL THEN NOW() ELSE started_at END,
		 completed_at=CASE WHEN $3='COMPLETED' THEN NOW() ELSE completed_at END,
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *calibrationRepo) CreateItem(ctx context.Context, item *domain.CalibrationItem) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO calibration_items (session_id, employee_id, original_score, adjusted_score, original_rating, adjusted_rating, reason)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 ON CONFLICT (session_id, employee_id) DO UPDATE SET adjusted_score=$4, adjusted_rating=$6, reason=$7
		 RETURNING id`,
		item.SessionID, item.EmployeeID, item.OriginalScore, item.AdjustedScore, item.OriginalRating, item.AdjustedRating, item.Reason,
	).Scan(&item.ID)
}

func (r *calibrationRepo) BulkCreateItems(ctx context.Context, items []domain.CalibrationItem) error {
	for i := range items {
		if err := r.CreateItem(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *calibrationRepo) ListItemsBySession(ctx context.Context, sessionID string) ([]domain.CalibrationItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, session_id, employee_id, original_score, adjusted_score, original_rating, adjusted_rating, reason, approved_by, approved_at
		 FROM calibration_items WHERE session_id=$1 ORDER BY employee_id`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CalibrationItem
	for rows.Next() {
		var item domain.CalibrationItem
		rows.Scan(&item.ID, &item.SessionID, &item.EmployeeID, &item.OriginalScore, &item.AdjustedScore, &item.OriginalRating, &item.AdjustedRating, &item.Reason, &item.ApprovedBy, &item.ApprovedAt)
		items = append(items, item)
	}
	return items, nil
}

func (r *calibrationRepo) UpdateItem(ctx context.Context, item *domain.CalibrationItem) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calibration_items SET adjusted_score=$2, adjusted_rating=$3, reason=$4
		 WHERE id=$1`,
		item.ID, item.AdjustedScore, item.AdjustedRating, item.Reason)
	return err
}

func (r *calibrationRepo) ApproveItem(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calibration_items SET approved_by=$2, approved_at=NOW() WHERE id=$1`, id, approvedBy)
	return err
}

// Plan repositories

type improvementPlanRepo struct {
	pool *pgxpool.Pool
}

func (r *improvementPlanRepo) Create(ctx context.Context, p *domain.ImprovementPlan) error {
	now := time.Now()
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_improvement_plans (company_id, employee_id, cycle_id, created_by, reason, start_date, end_date, status, success_criteria)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		p.CompanyID, p.EmployeeID, p.CycleID, p.CreatedBy, p.Reason, p.StartDate, p.EndDate, p.Status, p.SuccessCriteria,
	).Scan(&p.ID)
	p.CreatedAt = now
	p.UpdatedAt = now
	return err
}

func (r *improvementPlanRepo) GetByID(ctx context.Context, companyID, id string) (*domain.ImprovementPlan, error) {
	p := &domain.ImprovementPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, cycle_id, created_by, reason, start_date, end_date, status, success_criteria, final_result, created_at, updated_at
		 FROM performance_improvement_plans WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.CreatedBy, &p.Reason, &p.StartDate, &p.EndDate, &p.Status, &p.SuccessCriteria, &p.FinalResult, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *improvementPlanRepo) List(ctx context.Context, filter domain.PlanFilter) ([]domain.ImprovementPlan, error) {
	query := `SELECT id, company_id, employee_id, cycle_id, created_by, reason, start_date, end_date, status, success_criteria, final_result, created_at, updated_at
		 FROM performance_improvement_plans WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filter.EmployeeID)
		argIdx++
	}
	if filter.CycleID != "" {
		query += fmt.Sprintf(" AND cycle_id=$%d", argIdx)
		args = append(args, filter.CycleID)
		argIdx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []domain.ImprovementPlan
	for rows.Next() {
		var p domain.ImprovementPlan
		rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.CreatedBy, &p.Reason, &p.StartDate, &p.EndDate, &p.Status, &p.SuccessCriteria, &p.FinalResult, &p.CreatedAt, &p.UpdatedAt)
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *improvementPlanRepo) Update(ctx context.Context, p *domain.ImprovementPlan) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_improvement_plans SET reason=$3, start_date=$4, end_date=$5, success_criteria=$6, final_result=$7, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		p.CompanyID, p.ID, p.Reason, p.StartDate, p.EndDate, p.SuccessCriteria, p.FinalResult)
	return err
}

func (r *improvementPlanRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.PlanStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_improvement_plans SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *improvementPlanRepo) CreateAction(ctx context.Context, a *domain.ImprovementPlanAction) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO improvement_plan_actions (plan_id, title, description, responsible_id, due_date, status, progress, evidence)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		a.PlanID, a.Title, a.Description, a.ResponsibleID, a.DueDate, a.Status, a.Progress, a.Evidence,
	).Scan(&a.ID)
}

func (r *improvementPlanRepo) ListActionsByPlan(ctx context.Context, planID string) ([]domain.ImprovementPlanAction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, plan_id, title, description, responsible_id, due_date, status, progress, evidence, completed_at
		 FROM improvement_plan_actions WHERE plan_id=$1 ORDER BY due_date`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.ImprovementPlanAction
	for rows.Next() {
		var a domain.ImprovementPlanAction
		rows.Scan(&a.ID, &a.PlanID, &a.Title, &a.Description, &a.ResponsibleID, &a.DueDate, &a.Status, &a.Progress, &a.Evidence, &a.CompletedAt)
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *improvementPlanRepo) UpdateAction(ctx context.Context, a *domain.ImprovementPlanAction) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE improvement_plan_actions SET title=$2, description=$3, responsible_id=$4, due_date=$5, status=$6, progress=$7, evidence=$8,
		 completed_at=CASE WHEN $6='COMPLETED' THEN NOW() ELSE completed_at END
		 WHERE id=$1`,
		a.ID, a.Title, a.Description, a.ResponsibleID, a.DueDate, a.Status, a.Progress, a.Evidence)
	return err
}

// Development Plan

type devPlanRepo struct {
	pool *pgxpool.Pool
}

func (r *devPlanRepo) Create(ctx context.Context, p *domain.DevelopmentPlan) error {
	now := time.Now()
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_development_plans (company_id, employee_id, cycle_id, created_by, title, description, career_goal, current_level, target_level, competency_id, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
		p.CompanyID, p.EmployeeID, p.CycleID, p.CreatedBy, p.Title, p.Description, p.CareerGoal, p.CurrentLevel, p.TargetLevel, p.CompetencyID, p.Status, now, now,
	).Scan(&p.ID)
	return err
}

func (r *devPlanRepo) GetByID(ctx context.Context, companyID, id string) (*domain.DevelopmentPlan, error) {
	p := &domain.DevelopmentPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, cycle_id, created_by, title, description, career_goal, current_level, target_level, competency_id, status, created_at, updated_at
		 FROM performance_development_plans WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.CreatedBy, &p.Title, &p.Description, &p.CareerGoal, &p.CurrentLevel, &p.TargetLevel, &p.CompetencyID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *devPlanRepo) List(ctx context.Context, filter domain.PlanFilter) ([]domain.DevelopmentPlan, error) {
	query := `SELECT id, company_id, employee_id, cycle_id, created_by, title, description, career_goal, current_level, target_level, competency_id, status, created_at, updated_at
		 FROM performance_development_plans WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filter.EmployeeID)
		argIdx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []domain.DevelopmentPlan
	for rows.Next() {
		var p domain.DevelopmentPlan
		rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CycleID, &p.CreatedBy, &p.Title, &p.Description, &p.CareerGoal, &p.CurrentLevel, &p.TargetLevel, &p.CompetencyID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *devPlanRepo) Update(ctx context.Context, p *domain.DevelopmentPlan) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_development_plans SET title=$3, description=$4, career_goal=$5, current_level=$6, target_level=$7, competency_id=$8, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		p.CompanyID, p.ID, p.Title, p.Description, p.CareerGoal, p.CurrentLevel, p.TargetLevel, p.CompetencyID)
	return err
}

func (r *devPlanRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.PlanStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_development_plans SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *devPlanRepo) CreateAction(ctx context.Context, a *domain.DevelopmentPlanAction) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO development_plan_actions (plan_id, title, description, action_type, due_date, status)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		a.PlanID, a.Title, a.Description, a.ActionType, a.DueDate, a.Status,
	).Scan(&a.ID)
}

func (r *devPlanRepo) ListActionsByPlan(ctx context.Context, planID string) ([]domain.DevelopmentPlanAction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, plan_id, title, description, action_type, due_date, status, completed_at
		 FROM development_plan_actions WHERE plan_id=$1 ORDER BY due_date`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.DevelopmentPlanAction
	for rows.Next() {
		var a domain.DevelopmentPlanAction
		rows.Scan(&a.ID, &a.PlanID, &a.Title, &a.Description, &a.ActionType, &a.DueDate, &a.Status, &a.CompletedAt)
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *devPlanRepo) UpdateAction(ctx context.Context, a *domain.DevelopmentPlanAction) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE development_plan_actions SET title=$2, description=$3, action_type=$4, due_date=$5, status=$6,
		 completed_at=CASE WHEN $6='COMPLETED' THEN NOW() ELSE completed_at END
		 WHERE id=$1`,
		a.ID, a.Title, a.Description, a.ActionType, a.DueDate, a.Status)
	return err
}

func NewCalibrationRepository(pool *pgxpool.Pool) CalibrationRepository {
	return &calibrationRepo{pool: pool}
}

func NewImprovementPlanRepository(pool *pgxpool.Pool) ImprovementPlanRepository {
	return &improvementPlanRepo{pool: pool}
}

func NewDevelopmentPlanRepository(pool *pgxpool.Pool) DevelopmentPlanRepository {
	return &devPlanRepo{pool: pool}
}
