package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/shopspring/decimal"
)

type ExpenseRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseRepo(pool *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{pool: pool}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *domain.Expense) error {
	q := `INSERT INTO expenses (id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.EmployeeID, e.TravelID, e.CategoryID, e.PaymentMethodID,
		e.ExpenseReportID, e.Description, e.MerchantName, e.OriginalAmount, e.TaxAmount, e.TotalAmount, e.OriginalCurrency,
		e.ExchangeRate, e.ExpenseDate, e.Status, e.Observation, e.ReceiptRequired, e.IsBillable, e.BillableClient, e.CreatedBy)
	return repoErr("ExpenseRepo.Create", err)
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	q := `SELECT id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by,created_at,updated_at
		FROM expenses WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var e domain.Expense
	err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.TravelID, &e.CategoryID, &e.PaymentMethodID, &e.ExpenseReportID,
		&e.Description, &e.MerchantName, &e.OriginalAmount, &e.TaxAmount, &e.TotalAmount, &e.OriginalCurrency,
		&e.ExchangeRate, &e.ExpenseDate, &e.Status, &e.Observation, &e.ReceiptRequired, &e.IsBillable,
		&e.BillableClient, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, repoErr("ExpenseRepo.GetByID", err)
	}
	return &e, nil
}

func (r *ExpenseRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.Expense, error) {
	q := `SELECT id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by,created_at,updated_at
		FROM expenses WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.Expense
	err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.TravelID, &e.CategoryID, &e.PaymentMethodID, &e.ExpenseReportID,
		&e.Description, &e.MerchantName, &e.OriginalAmount, &e.TaxAmount, &e.TotalAmount, &e.OriginalCurrency,
		&e.ExchangeRate, &e.ExpenseDate, &e.Status, &e.Observation, &e.ReceiptRequired, &e.IsBillable,
		&e.BillableClient, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, repoErr("ExpenseRepo.Get", err)
	}
	return &e, nil
}

func (r *ExpenseRepo) List(ctx context.Context, companyID uuid.UUID, employeeID, travelID *uuid.UUID, status *string, dateFrom, dateTo *time.Time, limit, offset int) ([]domain.Expense, error) {
	q := `SELECT id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by,created_at,updated_at
		FROM expenses WHERE company_id=$1`
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
	if dateFrom != nil {
		q += fmt.Sprintf(" AND expense_date>=$%d", n)
		args = append(args, *dateFrom)
		n++
	}
	if dateTo != nil {
		q += fmt.Sprintf(" AND expense_date<=$%d", n)
		args = append(args, *dateTo)
		n++
	}
	q += " ORDER BY expense_date DESC"
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
		return nil, repoErr("ExpenseRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Expense, error) {
		var e domain.Expense
		err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.TravelID, &e.CategoryID, &e.PaymentMethodID, &e.ExpenseReportID,
			&e.Description, &e.MerchantName, &e.OriginalAmount, &e.TaxAmount, &e.TotalAmount, &e.OriginalCurrency,
			&e.ExchangeRate, &e.ExpenseDate, &e.Status, &e.Observation, &e.ReceiptRequired, &e.IsBillable,
			&e.BillableClient, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *ExpenseRepo) Update(ctx context.Context, e *domain.Expense) error {
	q := `UPDATE expenses SET category_id=$1,payment_method_id=$2,description=$3,merchant=$4,
		amount=$5,tax_amount=$6,total_amount=$7,currency=$8,exchange_rate=$9,expense_date=$10,
		status=$11,notes=$12,receipt_required=$13,is_billable=$14,billable_client=$15,updated_at=NOW()
		WHERE id=$16 AND company_id=$17`
	_, err := r.pool.Exec(ctx, q, e.CategoryID, e.PaymentMethodID, e.Description, e.MerchantName,
		e.OriginalAmount, e.TaxAmount, e.TotalAmount, e.OriginalCurrency, e.ExchangeRate, e.ExpenseDate,
		e.Status, e.Observation, e.ReceiptRequired, e.IsBillable, e.BillableClient, e.ID, e.CompanyID)
	return repoErr("ExpenseRepo.Update", err)
}

func (r *ExpenseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expenses SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("ExpenseRepo.UpdateStatus", err)
}

func (r *ExpenseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expenses WHERE id=$1`, id)
	return repoErr("ExpenseRepo.Delete", err)
}

func (r *ExpenseRepo) ListByReport(ctx context.Context, expenseReportID uuid.UUID) ([]domain.Expense, error) {
	q := `SELECT id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by,created_at,updated_at
		FROM expenses WHERE report_id=$1 ORDER BY expense_date`
	rows, err := r.pool.Query(ctx, q, expenseReportID)
	if err != nil {
		return nil, repoErr("ExpenseRepo.ListByReport", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Expense, error) {
		var e domain.Expense
		err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.TravelID, &e.CategoryID, &e.PaymentMethodID, &e.ExpenseReportID,
			&e.Description, &e.MerchantName, &e.OriginalAmount, &e.TaxAmount, &e.TotalAmount, &e.OriginalCurrency,
			&e.ExchangeRate, &e.ExpenseDate, &e.Status, &e.Observation, &e.ReceiptRequired, &e.IsBillable,
			&e.BillableClient, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *ExpenseRepo) Search(ctx context.Context, companyID uuid.UUID, query string) ([]domain.Expense, error) {
	q := `SELECT id,company_id,employee_id,travel_id,category_id,payment_method_id,report_id,
		description,merchant,amount,tax_amount,total_amount,currency,exchange_rate,expense_date,
		status,notes,receipt_required,is_billable,billable_client,created_by,created_at,updated_at
		FROM expenses WHERE company_id=$1 AND (description ILIKE $2 OR merchant ILIKE $2) ORDER BY expense_date DESC`
	searchPattern := "%" + query + "%"
	rows, err := r.pool.Query(ctx, q, companyID, searchPattern)
	if err != nil {
		return nil, repoErr("ExpenseRepo.Search", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Expense, error) {
		var e domain.Expense
		err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.TravelID, &e.CategoryID, &e.PaymentMethodID, &e.ExpenseReportID,
			&e.Description, &e.MerchantName, &e.OriginalAmount, &e.TaxAmount, &e.TotalAmount, &e.OriginalCurrency,
			&e.ExchangeRate, &e.ExpenseDate, &e.Status, &e.Observation, &e.ReceiptRequired, &e.IsBillable,
			&e.BillableClient, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *ExpenseRepo) GetTotalByEmployee(ctx context.Context, companyID, employeeID uuid.UUID, from, to time.Time) (decimal.Decimal, error) {
	q := `SELECT COALESCE(SUM(total_amount),0) FROM expenses
		WHERE company_id=$1 AND employee_id=$2 AND status='approved' AND expense_date>=$3 AND expense_date<=$4`
	var total decimal.Decimal
	err := r.pool.QueryRow(ctx, q, companyID, employeeID, from, to).Scan(&total)
	if err != nil {
		return decimal.Zero, repoErr("ExpenseRepo.GetTotalByEmployee", err)
	}
	return total, nil
}
