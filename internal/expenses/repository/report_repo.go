package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

func (r *ReportRepo) Create(ctx context.Context, rep *domain.ExpenseReport) error {
	q := `INSERT INTO expense_reports (id,company_id,employee_id,travel_id,title,description,
		total_amount,currency,status,submitted_at,approved_at,rejected_at,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, rep.ID, rep.CompanyID, rep.EmployeeID, rep.TravelID, rep.Title, rep.Description,
		rep.TotalAmount, rep.Currency, rep.Status, rep.SubmittedAt, rep.ApprovedAt, rep.RejectedAt, rep.Notes, rep.CreatedBy)
	return repoErr("ReportRepo.Create", err)
}

func (r *ReportRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseReport, error) {
	q := `SELECT id,company_id,employee_id,travel_id,title,description,
		total_amount,currency,status,submitted_at,approved_at,rejected_at,notes,created_by,created_at,updated_at
		FROM expense_reports WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var rep domain.ExpenseReport
	err := row.Scan(&rep.ID, &rep.CompanyID, &rep.EmployeeID, &rep.TravelID, &rep.Title, &rep.Description,
		&rep.TotalAmount, &rep.Currency, &rep.Status, &rep.SubmittedAt, &rep.ApprovedAt, &rep.RejectedAt,
		&rep.Notes, &rep.CreatedBy, &rep.CreatedAt, &rep.UpdatedAt)
	if err != nil {
		return nil, repoErr("ReportRepo.GetByID", err)
	}
	return &rep, nil
}

func (r *ReportRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseReport, error) {
	q := `SELECT id,company_id,employee_id,travel_id,title,description,
		total_amount,currency,status,submitted_at,approved_at,rejected_at,notes,created_by,created_at,updated_at
		FROM expense_reports WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var rep domain.ExpenseReport
	err := row.Scan(&rep.ID, &rep.CompanyID, &rep.EmployeeID, &rep.TravelID, &rep.Title, &rep.Description,
		&rep.TotalAmount, &rep.Currency, &rep.Status, &rep.SubmittedAt, &rep.ApprovedAt, &rep.RejectedAt,
		&rep.Notes, &rep.CreatedBy, &rep.CreatedAt, &rep.UpdatedAt)
	if err != nil {
		return nil, repoErr("ReportRepo.Get", err)
	}
	return &rep, nil
}

func (r *ReportRepo) List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseReport, error) {
	q := `SELECT id,company_id,employee_id,travel_id,title,description,
		total_amount,currency,status,submitted_at,approved_at,rejected_at,notes,created_by,created_at,updated_at
		FROM expense_reports WHERE company_id=$1`
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
		return nil, repoErr("ReportRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseReport, error) {
		var rep domain.ExpenseReport
		err := row.Scan(&rep.ID, &rep.CompanyID, &rep.EmployeeID, &rep.TravelID, &rep.Title, &rep.Description,
			&rep.TotalAmount, &rep.Currency, &rep.Status, &rep.SubmittedAt, &rep.ApprovedAt, &rep.RejectedAt,
			&rep.Notes, &rep.CreatedBy, &rep.CreatedAt, &rep.UpdatedAt)
		return rep, err
	})
}

func (r *ReportRepo) Update(ctx context.Context, rep *domain.ExpenseReport) error {
	q := `UPDATE expense_reports SET title=$1,description=$2,total_amount=$3,currency=$4,
		status=$5,submitted_at=$6,approved_at=$7,rejected_at=$8,notes=$9,updated_at=NOW()
		WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, rep.Title, rep.Description, rep.TotalAmount, rep.Currency,
		rep.Status, rep.SubmittedAt, rep.ApprovedAt, rep.RejectedAt, rep.Notes, rep.ID, rep.CompanyID)
	return repoErr("ReportRepo.Update", err)
}

func (r *ReportRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_reports SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("ReportRepo.UpdateStatus", err)
}

func (r *ReportRepo) GetByAdvance(ctx context.Context, advanceID uuid.UUID) (*domain.ExpenseReport, error) {
	q := `SELECT er.id,er.company_id,er.employee_id,er.travel_id,er.title,er.description,
		er.total_amount,er.currency,er.status,er.submitted_at,er.approved_at,er.rejected_at,er.notes,
		er.created_by,er.created_at,er.updated_at
		FROM expense_reports er
		INNER JOIN expense_advances ea ON ea.report_id=er.id
		WHERE ea.id=$1`
	row := r.pool.QueryRow(ctx, q, advanceID)
	var rep domain.ExpenseReport
	err := row.Scan(&rep.ID, &rep.CompanyID, &rep.EmployeeID, &rep.TravelID, &rep.Title, &rep.Description,
		&rep.TotalAmount, &rep.Currency, &rep.Status, &rep.SubmittedAt, &rep.ApprovedAt, &rep.RejectedAt,
		&rep.Notes, &rep.CreatedBy, &rep.CreatedAt, &rep.UpdatedAt)
	if err != nil {
		return nil, repoErr("ReportRepo.GetByAdvance", err)
	}
	return &rep, nil
}
