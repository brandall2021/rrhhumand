package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type objectiveRepo struct {
	pool *pgxpool.Pool
}

func (r *objectiveRepo) Create(ctx context.Context, o *domain.PerformanceObjective) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_objectives (company_id, cycle_id, employee_id, parent_objective_id, title, description,
		 objective_type, weight, start_date, due_date, status, progress, target_value, current_value, unit, progress_type,
		 notes, risk_notes, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
		 RETURNING id`,
		o.CompanyID, o.CycleID, o.EmployeeID, o.ParentObjectiveID, o.Title, o.Description,
		o.ObjectiveType, o.Weight, o.StartDate, o.DueDate, o.Status, o.Progress, o.TargetValue, o.CurrentValue,
		o.Unit, o.ProgressType, o.Notes, o.RiskNotes, o.CreatedBy, now, now,
	).Scan(&o.ID)
}

func (r *objectiveRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceObjective, error) {
	o := &domain.PerformanceObjective{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, parent_objective_id, title, description,
		 objective_type, weight, start_date, due_date, status, progress, target_value, current_value, unit, progress_type,
		 notes, risk_notes, created_by, created_at, updated_at
		 FROM performance_objectives WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&o.ID, &o.CompanyID, &o.CycleID, &o.EmployeeID, &o.ParentObjectiveID, &o.Title, &o.Description,
		&o.ObjectiveType, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.Progress, &o.TargetValue, &o.CurrentValue,
		&o.Unit, &o.ProgressType, &o.Notes, &o.RiskNotes, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *objectiveRepo) List(ctx context.Context, filter domain.ObjectiveFilter) ([]domain.PerformanceObjective, error) {
	query := `SELECT id, company_id, cycle_id, employee_id, parent_objective_id, title, description,
		 objective_type, weight, start_date, due_date, status, progress, target_value, current_value, unit, progress_type,
		 notes, risk_notes, created_by, created_at, updated_at
		 FROM performance_objectives WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.CycleID != "" {
		query += fmt.Sprintf(" AND cycle_id=$%d", argIdx)
		args = append(args, filter.CycleID)
		argIdx++
	}
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

	var objectives []domain.PerformanceObjective
	for rows.Next() {
		var o domain.PerformanceObjective
		rows.Scan(&o.ID, &o.CompanyID, &o.CycleID, &o.EmployeeID, &o.ParentObjectiveID, &o.Title, &o.Description,
			&o.ObjectiveType, &o.Weight, &o.StartDate, &o.DueDate, &o.Status, &o.Progress, &o.TargetValue, &o.CurrentValue,
			&o.Unit, &o.ProgressType, &o.Notes, &o.RiskNotes, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
		objectives = append(objectives, o)
	}
	return objectives, nil
}

func (r *objectiveRepo) Update(ctx context.Context, o *domain.PerformanceObjective) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_objectives SET title=$3, description=$4, objective_type=$5, weight=$6, start_date=$7,
		 due_date=$8, status=$9, progress=$10, target_value=$11, current_value=$12, unit=$13, progress_type=$14,
		 notes=$15, risk_notes=$16, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		o.CompanyID, o.ID, o.Title, o.Description, o.ObjectiveType, o.Weight, o.StartDate, o.DueDate,
		o.Status, o.Progress, o.TargetValue, o.CurrentValue, o.Unit, o.ProgressType, o.Notes, o.RiskNotes)
	return err
}

func (r *objectiveRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_objectives WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *objectiveRepo) CreateKeyResult(ctx context.Context, kr *domain.ObjectiveKeyResult) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO objective_key_results (objective_id, title, description, weight, target_value, current_value, unit, progress, status, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		kr.ObjectiveID, kr.Title, kr.Description, kr.Weight, kr.TargetValue, kr.CurrentValue, kr.Unit, kr.Progress, kr.Status, kr.SortOrder,
	).Scan(&kr.ID)
}

func (r *objectiveRepo) ListKeyResultsByObjective(ctx context.Context, objectiveID string) ([]domain.ObjectiveKeyResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, objective_id, title, description, weight, target_value, current_value, unit, progress, status, sort_order
		 FROM objective_key_results WHERE objective_id=$1 ORDER BY sort_order`, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var krs []domain.ObjectiveKeyResult
	for rows.Next() {
		var kr domain.ObjectiveKeyResult
		rows.Scan(&kr.ID, &kr.ObjectiveID, &kr.Title, &kr.Description, &kr.Weight, &kr.TargetValue, &kr.CurrentValue, &kr.Unit, &kr.Progress, &kr.Status, &kr.SortOrder)
		krs = append(krs, kr)
	}
	return krs, nil
}

func (r *objectiveRepo) UpdateKeyResult(ctx context.Context, kr *domain.ObjectiveKeyResult) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE objective_key_results SET title=$2, description=$3, weight=$4, target_value=$5, current_value=$6,
		 unit=$7, progress=$8, status=$9 WHERE id=$1`,
		kr.ID, kr.Title, kr.Description, kr.Weight, kr.TargetValue, kr.CurrentValue, kr.Unit, kr.Progress, kr.Status)
	return err
}

func (r *objectiveRepo) DeleteKeyResult(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM objective_key_results WHERE id=$1`, id)
	return err
}

func (r *objectiveRepo) GetWeightTotal(ctx context.Context, cycleID, employeeID string) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(weight),0) FROM performance_objectives WHERE cycle_id=$1 AND employee_id=$2 AND status!='CANCELLED'`,
		cycleID, employeeID).Scan(&total)
	return total, err
}

// ParticipantRepository

type participantRepo struct {
	pool *pgxpool.Pool
}

func (r *participantRepo) Create(ctx context.Context, p *domain.PerformanceParticipant) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_participants (company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.CompanyID, p.CycleID, p.EmployeeID, p.EvaluatorID, p.EvaluationType, p.Status, p.AssignedAt,
	).Scan(&p.ID)
}

func (r *participantRepo) BulkCreate(ctx context.Context, participants []domain.PerformanceParticipant) error {
	for i := range participants {
		p := &participants[i]
		now := time.Now()
		err := r.pool.QueryRow(ctx,
			`INSERT INTO performance_participants (company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 ON CONFLICT (cycle_id, employee_id, evaluator_id, evaluation_type) DO UPDATE SET status='PENDING'
			 RETURNING id`,
			p.CompanyID, p.CycleID, p.EmployeeID, p.EvaluatorID, p.EvaluationType, p.Status, now,
		).Scan(&p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *participantRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceParticipant, error) {
	p := &domain.PerformanceParticipant{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at, started_at, submitted_at
		 FROM performance_participants WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.CycleID, &p.EmployeeID, &p.EvaluatorID, &p.EvaluationType, &p.Status, &p.AssignedAt, &p.StartedAt, &p.SubmittedAt)
	return p, err
}

func (r *participantRepo) ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at, started_at, submitted_at
		 FROM performance_participants WHERE company_id=$1 AND cycle_id=$2 ORDER BY assigned_at`, companyID, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.PerformanceParticipant
	for rows.Next() {
		var p domain.PerformanceParticipant
		rows.Scan(&p.ID, &p.CompanyID, &p.CycleID, &p.EmployeeID, &p.EvaluatorID, &p.EvaluationType, &p.Status, &p.AssignedAt, &p.StartedAt, &p.SubmittedAt)
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *participantRepo) ListByEmployee(ctx context.Context, companyID, cycleID, employeeID string) ([]domain.PerformanceParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at, started_at, submitted_at
		 FROM performance_participants WHERE company_id=$1 AND cycle_id=$2 AND employee_id=$3 ORDER BY evaluation_type`,
		companyID, cycleID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.PerformanceParticipant
	for rows.Next() {
		var p domain.PerformanceParticipant
		rows.Scan(&p.ID, &p.CompanyID, &p.CycleID, &p.EmployeeID, &p.EvaluatorID, &p.EvaluationType, &p.Status, &p.AssignedAt, &p.StartedAt, &p.SubmittedAt)
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *participantRepo) ListByEvaluator(ctx context.Context, companyID, cycleID, evaluatorID string) ([]domain.PerformanceParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, status, assigned_at, started_at, submitted_at
		 FROM performance_participants WHERE company_id=$1 AND cycle_id=$2 AND evaluator_id=$3 ORDER BY assigned_at`,
		companyID, cycleID, evaluatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.PerformanceParticipant
	for rows.Next() {
		var p domain.PerformanceParticipant
		rows.Scan(&p.ID, &p.CompanyID, &p.CycleID, &p.EmployeeID, &p.EvaluatorID, &p.EvaluationType, &p.Status, &p.AssignedAt, &p.StartedAt, &p.SubmittedAt)
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *participantRepo) UpdateStatus(ctx context.Context, id string, status domain.EvaluationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_participants SET status=$2,
		 started_at=CASE WHEN $2='IN_PROGRESS' AND started_at IS NULL THEN NOW() ELSE started_at END,
		 submitted_at=CASE WHEN $2='SUBMITTED' THEN NOW() ELSE submitted_at END
		 WHERE id=$1`, id, status)
	return err
}

func (r *participantRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_participants WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func NewObjectiveRepository(pool *pgxpool.Pool) ObjectiveRepository {
	return &objectiveRepo{pool: pool}
}

func NewParticipantRepository(pool *pgxpool.Pool) ParticipantRepository {
	return &participantRepo{pool: pool}
}
