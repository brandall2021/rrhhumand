package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// ----- Cycles -----

func (r *PostgresRepository) CreateCycle(ctx context.Context, c *domain.PerformanceCycle) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_cycles (company_id, name, description, cycle_type, status, start_date, end_date,
		 evaluation_start_date, evaluation_end_date, review_start_date, review_end_date,
		 calibration_start_date, calibration_end_date, template_id, objective_weight, competency_weight,
		 min_anonymous_responses, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		 RETURNING id`,
		c.CompanyID, c.Name, c.Description, c.CycleType, c.Status, c.StartDate, c.EndDate,
		c.EvaluationStartDate, c.EvaluationEndDate, c.ReviewStartDate, c.ReviewEndDate,
		c.CalibrationStartDate, c.CalibrationEndDate, c.TemplateID, c.ObjectiveWeight, c.CompetencyWeight,
		c.MinAnonymousResponses, c.CreatedBy, now, now,
	).Scan(&c.ID)
}

func (r *PostgresRepository) GetCycleByID(ctx context.Context, companyID, id string) (*domain.PerformanceCycle, error) {
	c := &domain.PerformanceCycle{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, cycle_type, status, start_date, end_date,
		 evaluation_start_date, evaluation_end_date, review_start_date, review_end_date,
		 calibration_start_date, calibration_end_date, template_id, objective_weight, competency_weight,
		 min_anonymous_responses, created_by, created_at, updated_at
		 FROM performance_cycles WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CycleType, &c.Status, &c.StartDate, &c.EndDate,
		&c.EvaluationStartDate, &c.EvaluationEndDate, &c.ReviewStartDate, &c.ReviewEndDate,
		&c.CalibrationStartDate, &c.CalibrationEndDate, &c.TemplateID, &c.ObjectiveWeight, &c.CompetencyWeight,
		&c.MinAnonymousResponses, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresRepository) ListCycles(ctx context.Context, filter domain.PerformanceCycleFilter) ([]domain.PerformanceCycle, error) {
	query := `SELECT id, company_id, name, description, cycle_type, status, start_date, end_date,
		 evaluation_start_date, evaluation_end_date, review_start_date, review_end_date,
		 calibration_start_date, calibration_end_date, template_id, objective_weight, competency_weight,
		 min_anonymous_responses, created_by, created_at, updated_at
		 FROM performance_cycles WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND cycle_type=$%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []domain.PerformanceCycle
	for rows.Next() {
		var c domain.PerformanceCycle
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CycleType, &c.Status,
			&c.StartDate, &c.EndDate, &c.EvaluationStartDate, &c.EvaluationEndDate,
			&c.ReviewStartDate, &c.ReviewEndDate, &c.CalibrationStartDate, &c.CalibrationEndDate,
			&c.TemplateID, &c.ObjectiveWeight, &c.CompetencyWeight,
			&c.MinAnonymousResponses, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cycles = append(cycles, c)
	}
	return cycles, nil
}

func (r *PostgresRepository) UpdateCycle(ctx context.Context, c *domain.PerformanceCycle) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_cycles SET name=$3, description=$4, cycle_type=$5, start_date=$6, end_date=$7,
		 evaluation_start_date=$8, evaluation_end_date=$9, review_start_date=$10, review_end_date=$11,
		 calibration_start_date=$12, calibration_end_date=$13, template_id=$14, objective_weight=$15,
		 competency_weight=$16, min_anonymous_responses=$17, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		c.CompanyID, c.ID, c.Name, c.Description, c.CycleType, c.StartDate, c.EndDate,
		c.EvaluationStartDate, c.EvaluationEndDate, c.ReviewStartDate, c.ReviewEndDate,
		c.CalibrationStartDate, c.CalibrationEndDate, c.TemplateID, c.ObjectiveWeight, c.CompetencyWeight,
		c.MinAnonymousResponses)
	return err
}

func (r *PostgresRepository) UpdateCycleStatus(ctx context.Context, companyID, id string, status domain.CycleStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_cycles SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *PostgresRepository) DeleteCycle(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_cycles WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}
