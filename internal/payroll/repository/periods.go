package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreatePeriod(ctx context.Context, p *domain.PayrollPeriod) error {
	q := `INSERT INTO payroll_periods (id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Year, p.Month, p.PeriodType, p.Name, p.StartDate, p.EndDate, p.PaymentDate, p.Status, p.CreatedBy)
	return repoErr("CreatePeriod", err)
}

func (r *Repository) GetPeriod(ctx context.Context, companyID, id string) (*domain.PayrollPeriod, error) {
	q := `SELECT id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,closed_at,created_by,created_at,updated_at
		FROM payroll_periods WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p domain.PayrollPeriod
	err := row.Scan(&p.ID, &p.CompanyID, &p.Year, &p.Month, &p.PeriodType, &p.Name, &p.StartDate, &p.EndDate, &p.PaymentDate,
		&p.Status, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetPeriod", err)
	}
	return &p, nil
}

func (r *Repository) UpdatePeriod(ctx context.Context, p *domain.PayrollPeriod) error {
	q := `UPDATE payroll_periods SET name=$1,payment_date=$2,updated_at=NOW() WHERE id=$3 AND company_id=$4`
	_, err := r.pool.Exec(ctx, q, p.Name, p.PaymentDate, p.ID, p.CompanyID)
	return repoErr("UpdatePeriod", err)
}

func (r *Repository) ListPeriods(ctx context.Context, companyID string, limit, offset int) ([]domain.PayrollPeriod, error) {
	q := `SELECT id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,closed_at,created_by,created_at,updated_at
		FROM payroll_periods WHERE company_id=$1 ORDER BY year DESC, month DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, companyID, limit, offset)
	if err != nil {
		return nil, repoErr("ListPeriods", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollPeriod, error) {
		var p domain.PayrollPeriod
		err := row.Scan(&p.ID, &p.CompanyID, &p.Year, &p.Month, &p.PeriodType, &p.Name, &p.StartDate, &p.EndDate, &p.PaymentDate,
			&p.Status, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *Repository) UpdatePeriodStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_periods SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdatePeriodStatus", err)
}

func (r *Repository) ClosePeriod(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_periods SET status='CLOSED',closed_at=NOW(),updated_at=NOW() WHERE id=$1`, id)
	return repoErr("ClosePeriod", err)
}

func (r *Repository) CreateRun(ctx context.Context, run *domain.PayrollRun) error {
	q := `INSERT INTO payroll_runs (id,company_id,period_id,run_number,run_type,status,engine_version,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, run.ID, run.CompanyID, run.PeriodID, run.RunNumber, run.RunType, run.Status, run.EngineVersion, run.CreatedBy)
	return repoErr("CreateRun", err)
}

func (r *Repository) GetRun(ctx context.Context, companyID, id string) (*domain.PayrollRun, error) {
	q := `SELECT id,company_id,period_id,run_number,run_type,status,engine_version,started_at,finished_at,
		created_by,approved_by,approved_at,closed_by,closed_at,created_at,updated_at
		FROM payroll_runs WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var run domain.PayrollRun
	err := row.Scan(&run.ID, &run.CompanyID, &run.PeriodID, &run.RunNumber, &run.RunType, &run.Status, &run.EngineVersion,
		&run.StartedAt, &run.FinishedAt, &run.CreatedBy, &run.ApprovedBy, &run.ApprovedAt, &run.ClosedBy, &run.ClosedAt, &run.CreatedAt, &run.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetRun", err)
	}
	return &run, nil
}

func (r *Repository) ListRuns(ctx context.Context, companyID string, periodID, runType, status *string, limit, offset int) ([]domain.PayrollRun, error) {
	q := `SELECT id,company_id,period_id,run_number,run_type,status,engine_version,started_at,finished_at,
		created_by,approved_by,approved_at,closed_by,closed_at,created_at,updated_at
		FROM payroll_runs WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if periodID != nil {
		q += fmt.Sprintf(" AND period_id=$%d", n)
		args = append(args, *periodID)
		n++
	}
	if runType != nil {
		q += fmt.Sprintf(" AND run_type=$%d", n)
		args = append(args, *runType)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	q += " ORDER BY created_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListRuns", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollRun, error) {
		var run domain.PayrollRun
		err := row.Scan(&run.ID, &run.CompanyID, &run.PeriodID, &run.RunNumber, &run.RunType, &run.Status, &run.EngineVersion,
			&run.StartedAt, &run.FinishedAt, &run.CreatedBy, &run.ApprovedBy, &run.ApprovedAt, &run.ClosedBy, &run.ClosedAt, &run.CreatedAt, &run.UpdatedAt)
		return run, err
	})
}

func (r *Repository) UpdateRunStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdateRunStatus", err)
}

func (r *Repository) UpdateRunTimestamps(ctx context.Context, id, status string, startedAt, finishedAt *time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status=$1,started_at=$2,finished_at=$3,updated_at=NOW() WHERE id=$4`, status, startedAt, finishedAt, id)
	return repoErr("UpdateRunTimestamps", err)
}

func (r *Repository) ApproveRun(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status='APPROVED',approved_by=$1,approved_at=NOW(),updated_at=NOW() WHERE id=$2`, approvedBy, id)
	return repoErr("ApproveRun", err)
}

func (r *Repository) CloseRun(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status='CLOSED',closed_at=NOW(),updated_at=NOW() WHERE id=$1`, id)
	return repoErr("CloseRun", err)
}

func (r *Repository) GetRunNumber(ctx context.Context, periodID string) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(MAX(run_number),0)+1 FROM payroll_runs WHERE period_id=$1`, periodID).Scan(&n)
	return n, repoErr("GetRunNumber", err)
}
