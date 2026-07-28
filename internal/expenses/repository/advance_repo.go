package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type AdvanceRepo struct {
	pool *pgxpool.Pool
}

func NewAdvanceRepo(pool *pgxpool.Pool) *AdvanceRepo {
	return &AdvanceRepo{pool: pool}
}

func (r *AdvanceRepo) Create(ctx context.Context, a *domain.ExpenseAdvance) error {
	q := `INSERT INTO expense_advances (id,company_id,employee_id,travel_id,report_id,amount,currency,
		requested_at,status,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EmployeeID, a.TravelID, a.ReportID,
		a.Amount, a.Currency, a.RequestedAt, a.Status, a.Notes, a.CreatedBy)
	return repoErr("AdvanceRepo.Create", err)
}

func (r *AdvanceRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,report_id,amount,currency,
		requested_at,approved_at,disbursed_at,status,notes,created_by,created_at,updated_at
		FROM expense_advances WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a domain.ExpenseAdvance
	err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.ReportID, &a.Amount, &a.Currency,
		&a.RequestedAt, &a.ApprovedAt, &a.DisbursedAt, &a.Status, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("AdvanceRepo.Get", err)
	}
	return &a, nil
}

func (r *AdvanceRepo) List(ctx context.Context, companyID uuid.UUID, employeeID, travelID *uuid.UUID, status *string) ([]domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,report_id,amount,currency,
		requested_at,approved_at,disbursed_at,status,notes,created_by,created_at,updated_at
		FROM expense_advances WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if travelID != nil {
		q += fmt.Sprintf(" AND travel_id=$%d", n)
		args = append(args, *travelID)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	q += " ORDER BY requested_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("AdvanceRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseAdvance, error) {
		var a domain.ExpenseAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.ReportID, &a.Amount, &a.Currency,
			&a.RequestedAt, &a.ApprovedAt, &a.DisbursedAt, &a.Status, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *AdvanceRepo) Update(ctx context.Context, a *domain.ExpenseAdvance) error {
	q := `UPDATE expense_advances SET travel_id=$1,report_id=$2,amount=$3,currency=$4,
		requested_at=$5,approved_at=$6,disbursed_at=$7,status=$8,notes=$9,updated_at=NOW()
		WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, a.TravelID, a.ReportID, a.Amount, a.Currency,
		a.RequestedAt, a.ApprovedAt, a.DisbursedAt, a.Status, a.Notes, a.ID, a.CompanyID)
	return repoErr("AdvanceRepo.Update", err)
}

func (r *AdvanceRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_advances SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("AdvanceRepo.UpdateStatus", err)
}

func (r *AdvanceRepo) ListByEmployee(ctx context.Context, companyID, employeeID uuid.UUID) ([]domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,report_id,amount,currency,
		requested_at,approved_at,disbursed_at,status,notes,created_by,created_at,updated_at
		FROM expense_advances WHERE company_id=$1 AND employee_id=$2 ORDER BY requested_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("AdvanceRepo.ListByEmployee", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseAdvance, error) {
		var a domain.ExpenseAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.ReportID, &a.Amount, &a.Currency,
			&a.RequestedAt, &a.ApprovedAt, &a.DisbursedAt, &a.Status, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}
