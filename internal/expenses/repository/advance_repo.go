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
	q := `INSERT INTO expense_advances (id,company_id,employee_id,travel_id,requested_amount,approved_amount,settled_amount,currency,
		request_date,approved_date,paid_date,settled_date,status,rejection_reason,idempotency_key,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EmployeeID, a.TravelID, a.RequestedAmount,
		a.ApprovedAmount, a.SettledAmount, a.Currency, a.RequestDate, a.ApprovedDate, a.PaidDate, a.SettledDate, a.Status, a.RejectionReason, a.IdempotencyKey, a.CreatedBy)
	return repoErr("AdvanceRepo.Create", err)
}

func (r *AdvanceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,requested_amount,approved_amount,settled_amount,currency,
		request_date,approved_date,paid_date,settled_date,status,rejection_reason,idempotency_key,created_by,created_at,updated_at
		FROM expense_advances WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var a domain.ExpenseAdvance
	err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.RequestedAmount,
		&a.ApprovedAmount, &a.SettledAmount, &a.Currency, &a.RequestDate, &a.ApprovedDate, &a.PaidDate, &a.SettledDate, &a.Status, &a.RejectionReason, &a.IdempotencyKey, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("AdvanceRepo.GetByID", err)
	}
	return &a, nil
}

func (r *AdvanceRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,requested_amount,approved_amount,settled_amount,currency,
		request_date,approved_date,paid_date,settled_date,status,rejection_reason,idempotency_key,created_by,created_at,updated_at
		FROM expense_advances WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a domain.ExpenseAdvance
	err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.RequestedAmount,
		&a.ApprovedAmount, &a.SettledAmount, &a.Currency, &a.RequestDate, &a.ApprovedDate, &a.PaidDate, &a.SettledDate, &a.Status, &a.RejectionReason, &a.IdempotencyKey, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("AdvanceRepo.Get", err)
	}
	return &a, nil
}

func (r *AdvanceRepo) List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,requested_amount,approved_amount,settled_amount,currency,
		request_date,approved_date,paid_date,settled_date,status,rejection_reason,idempotency_key,created_by,created_at,updated_at
		FROM expense_advances WHERE company_id=$1`
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
	q += " ORDER BY request_date DESC"
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
		return nil, repoErr("AdvanceRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseAdvance, error) {
		var a domain.ExpenseAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.RequestedAmount,
			&a.ApprovedAmount, &a.SettledAmount, &a.Currency, &a.RequestDate, &a.ApprovedDate, &a.PaidDate, &a.SettledDate, &a.Status, &a.RejectionReason, &a.IdempotencyKey, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *AdvanceRepo) Update(ctx context.Context, a *domain.ExpenseAdvance) error {
	q := `UPDATE expense_advances SET travel_id=$1,requested_amount=$2,approved_amount=$3,settled_amount=$4,currency=$5,
		request_date=$6,approved_date=$7,paid_date=$8,settled_date=$9,status=$10,rejection_reason=$11,idempotency_key=$12,updated_at=NOW()
		WHERE id=$13 AND company_id=$14`
	_, err := r.pool.Exec(ctx, q, a.TravelID, a.RequestedAmount, a.ApprovedAmount, a.SettledAmount, a.Currency,
		a.RequestDate, a.ApprovedDate, a.PaidDate, a.SettledDate, a.Status, a.RejectionReason, a.IdempotencyKey, a.ID, a.CompanyID)
	return repoErr("AdvanceRepo.Update", err)
}

func (r *AdvanceRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_advances SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("AdvanceRepo.UpdateStatus", err)
}

func (r *AdvanceRepo) ListByEmployee(ctx context.Context, companyID, employeeID uuid.UUID) ([]domain.ExpenseAdvance, error) {
	q := `SELECT id,company_id,employee_id,travel_id,requested_amount,approved_amount,settled_amount,currency,
		request_date,approved_date,paid_date,settled_date,status,rejection_reason,idempotency_key,created_by,created_at,updated_at
		FROM expense_advances WHERE company_id=$1 AND employee_id=$2 ORDER BY request_date DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("AdvanceRepo.ListByEmployee", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseAdvance, error) {
		var a domain.ExpenseAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TravelID, &a.RequestedAmount,
			&a.ApprovedAmount, &a.SettledAmount, &a.Currency, &a.RequestDate, &a.ApprovedDate, &a.PaidDate, &a.SettledDate, &a.Status, &a.RejectionReason, &a.IdempotencyKey, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}
