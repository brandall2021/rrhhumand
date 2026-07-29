package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/onboarding/domain"
)

type OnboardingRepo struct {
	pool *pgxpool.Pool
}

func NewOnboardingRepo(pool *pgxpool.Pool) *OnboardingRepo {
	return &OnboardingRepo{pool: pool}
}

func (r *OnboardingRepo) Create(ctx context.Context, p *domain.OnboardingProcess) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_processes
		 (company_id, employee_id, candidate_id, application_id, job_offer_id, template_id,
		  status, start_date, expected_completion_date, completion_policy,
		  employee_type, work_mode, probation_start_date, probation_end_date, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 RETURNING id, created_at, updated_at`,
		p.CompanyID, p.EmployeeID, p.CandidateID, p.ApplicationID, p.JobOfferID, p.TemplateID,
		p.Status, p.StartDate, p.ExpectedCompletionDate, p.CompletionPolicy,
		p.EmployeeType, p.WorkMode, p.ProbationStartDate, p.ProbationEndDate, p.CreatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *OnboardingRepo) GetByID(ctx context.Context, companyID, id string) (*domain.OnboardingProcess, error) {
	p := &domain.OnboardingProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, candidate_id, application_id, job_offer_id,
		        template_id, status, start_date, expected_completion_date, actual_completion_date,
		        progress_percentage, completion_policy, employee_type, work_mode,
		        probation_start_date, probation_end_date, probation_status,
		        created_by, created_at, updated_at
		 FROM onboarding_processes WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CandidateID, &p.ApplicationID,
		&p.JobOfferID, &p.TemplateID, &p.Status, &p.StartDate,
		&p.ExpectedCompletionDate, &p.ActualCompletionDate,
		&p.Progress, &p.CompletionPolicy, &p.EmployeeType, &p.WorkMode,
		&p.ProbationStartDate, &p.ProbationEndDate, &p.ProbationStatus,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *OnboardingRepo) GetByEmployee(ctx context.Context, companyID, employeeID string) (*domain.OnboardingProcess, error) {
	p := &domain.OnboardingProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, candidate_id, application_id, job_offer_id,
		        template_id, status, start_date, expected_completion_date, actual_completion_date,
		        progress_percentage, completion_policy, employee_type, work_mode,
		        probation_start_date, probation_end_date, probation_status,
		        created_by, created_at, updated_at
		 FROM onboarding_processes
		 WHERE company_id=$1 AND employee_id=$2 AND status NOT IN ('COMPLETED','CANCELLED')
		 ORDER BY created_at DESC LIMIT 1`, companyID, employeeID,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CandidateID, &p.ApplicationID,
		&p.JobOfferID, &p.TemplateID, &p.Status, &p.StartDate,
		&p.ExpectedCompletionDate, &p.ActualCompletionDate,
		&p.Progress, &p.CompletionPolicy, &p.EmployeeType, &p.WorkMode,
		&p.ProbationStartDate, &p.ProbationEndDate, &p.ProbationStatus,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *OnboardingRepo) List(ctx context.Context, companyID string, status, employeeID, search string) ([]domain.OnboardingProcess, error) {
	query := `SELECT id, company_id, employee_id, candidate_id, application_id, job_offer_id,
	                 template_id, status, start_date, expected_completion_date, actual_completion_date,
	                 progress_percentage, completion_policy, employee_type, work_mode,
	                 probation_start_date, probation_end_date, probation_status,
	                 created_by, created_at, updated_at
	          FROM onboarding_processes WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if employeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, employeeID)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ps []domain.OnboardingProcess
	for rows.Next() {
		var p domain.OnboardingProcess
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.CandidateID, &p.ApplicationID,
			&p.JobOfferID, &p.TemplateID, &p.Status, &p.StartDate,
			&p.ExpectedCompletionDate, &p.ActualCompletionDate,
			&p.Progress, &p.CompletionPolicy, &p.EmployeeType, &p.WorkMode,
			&p.ProbationStartDate, &p.ProbationEndDate, &p.ProbationStatus,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (r *OnboardingRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.OnboardingStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *OnboardingRepo) UpdateProgress(ctx context.Context, companyID, id string, progress float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET progress_percentage=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, progress)
	return err
}

func (r *OnboardingRepo) Complete(ctx context.Context, companyID, id string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status='COMPLETED', actual_completion_date=$3, progress_percentage=100, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, now)
	return err
}

func (r *OnboardingRepo) Cancel(ctx context.Context, companyID, id, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status='CANCELLED', cancellation_reason=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, reason)
	return err
}

func (r *OnboardingRepo) HasActiveProcess(ctx context.Context, companyID, employeeID string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_processes
		 WHERE company_id=$1 AND employee_id=$2 AND status NOT IN ('COMPLETED','CANCELLED')`,
		companyID, employeeID).Scan(&count)
	return count > 0, err
}

func (r *OnboardingRepo) UpdateProbation(ctx context.Context, companyID, id string, status domain.ProbationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET probation_status=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, status)
	return err
}

func (r *OnboardingRepo) GetDashboardStats(ctx context.Context, companyID string) (active, pending, completed, overdue int, avgProgress float64, err error) {
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='IN_PROGRESS'`, companyID).Scan(&active)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='NOT_STARTED'`, companyID).Scan(&pending)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='COMPLETED'`, companyID).Scan(&completed)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status IN ('IN_PROGRESS','NOT_STARTED') AND expected_completion_date < CURRENT_DATE`, companyID).Scan(&overdue)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(progress_percentage),0) FROM onboarding_processes WHERE company_id=$1 AND status IN ('IN_PROGRESS','NOT_STARTED')`, companyID).Scan(&avgProgress)
	return
}
