package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreateNovelty(ctx context.Context, n *domain.PayrollNovelty) error {
	q := `INSERT INTO payroll_novelties (id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,
		unit_value,multiplier,start_date,end_date,description,source,source_reference_id,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, n.ID, n.CompanyID, n.EmployeeID, n.PeriodID, n.NoveltyType,
		n.Quantity, n.Unit, n.Amount, n.UnitValue, n.Multiplier, n.StartDate, n.EndDate, n.Description,
		n.Source, n.SourceReferenceID, n.Status, n.CreatedBy)
	return repoErr("CreateNovelty", err)
}

func (r *Repository) UpdateNovelty(ctx context.Context, n *domain.PayrollNovelty) error {
	q := `UPDATE payroll_novelties SET quantity=$1,amount=$2,description=$3,status=$4,updated_at=NOW() WHERE id=$5 AND company_id=$6`
	_, err := r.pool.Exec(ctx, q, n.Quantity, n.Amount, n.Description, n.Status, n.ID, n.CompanyID)
	return repoErr("UpdateNovelty", err)
}

func (r *Repository) GetNovelty(ctx context.Context, companyID, id string) (*domain.PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var n domain.PayrollNovelty
	err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
		&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
		&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetNovelty", err)
	}
	return &n, nil
}

func (r *Repository) ListNovelties(ctx context.Context, companyID string, employeeID, periodID, noveltyType, status, source *string, limit, offset int) ([]domain.PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if periodID != nil {
		q += fmt.Sprintf(" AND period_id=$%d", n)
		args = append(args, *periodID)
		n++
	}
	if noveltyType != nil {
		q += fmt.Sprintf(" AND novelty_type=$%d", n)
		args = append(args, *noveltyType)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	if source != nil {
		q += fmt.Sprintf(" AND source=$%d", n)
		args = append(args, *source)
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
		return nil, repoErr("ListNovelties", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollNovelty, error) {
		var n domain.PayrollNovelty
		err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
			&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
			&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
		return n, err
	})
}

func (r *Repository) DeleteNovelty(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_novelties WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteNovelty", err)
}

func (r *Repository) ApproveNovelty(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_novelties SET status='APPROVED',approved_by=$1,approved_at=NOW() WHERE id=$2`, approvedBy, id)
	return repoErr("ApproveNovelty", err)
}

func (r *Repository) GetNoveltiesForEmployeePeriod(ctx context.Context, companyID, employeeID, periodID string) ([]domain.PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE company_id=$1 AND employee_id=$2 AND period_id=$3 AND status='APPROVED' ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID, periodID)
	if err != nil {
		return nil, repoErr("GetNoveltiesForEmployeePeriod", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollNovelty, error) {
		var n domain.PayrollNovelty
		err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
			&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
			&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
		return n, err
	})
}
