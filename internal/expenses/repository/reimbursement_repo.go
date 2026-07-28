package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ReimbursementRepo struct {
	pool *pgxpool.Pool
}

func NewReimbursementRepo(pool *pgxpool.Pool) *ReimbursementRepo {
	return &ReimbursementRepo{pool: pool}
}

func (r *ReimbursementRepo) Create(ctx context.Context, re *domain.ExpenseReimbursement) error {
	q := `INSERT INTO expense_reimbursements (id,company_id,employee_id,report_id,amount,currency,
		payment_method,payment_date,status,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, re.ID, re.CompanyID, re.EmployeeID, re.ReportID, re.Amount, re.Currency,
		re.PaymentMethod, re.PaymentDate, re.Status, re.Notes, re.CreatedBy)
	return repoErr("ReimbursementRepo.Create", err)
}

func (r *ReimbursementRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseReimbursement, error) {
	q := `SELECT id,company_id,employee_id,report_id,amount,currency,
		payment_method,payment_date,status,notes,created_by,created_at,updated_at
		FROM expense_reimbursements WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var re domain.ExpenseReimbursement
	err := row.Scan(&re.ID, &re.CompanyID, &re.EmployeeID, &re.ReportID, &re.Amount, &re.Currency,
		&re.PaymentMethod, &re.PaymentDate, &re.Status, &re.Notes, &re.CreatedBy, &re.CreatedAt, &re.UpdatedAt)
	if err != nil {
		return nil, repoErr("ReimbursementRepo.Get", err)
	}
	return &re, nil
}

func (r *ReimbursementRepo) List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string) ([]domain.ExpenseReimbursement, error) {
	q := `SELECT id,company_id,employee_id,report_id,amount,currency,
		payment_method,payment_date,status,notes,created_by,created_at,updated_at
		FROM expense_reimbursements WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ReimbursementRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseReimbursement, error) {
		var re domain.ExpenseReimbursement
		err := row.Scan(&re.ID, &re.CompanyID, &re.EmployeeID, &re.ReportID, &re.Amount, &re.Currency,
			&re.PaymentMethod, &re.PaymentDate, &re.Status, &re.Notes, &re.CreatedBy, &re.CreatedAt, &re.UpdatedAt)
		return re, err
	})
}

func (r *ReimbursementRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_reimbursements SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("ReimbursementRepo.UpdateStatus", err)
}
