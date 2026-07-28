package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/shopspring/decimal"
)

type BudgetRepo struct {
	pool *pgxpool.Pool
}

func NewBudgetRepo(pool *pgxpool.Pool) *BudgetRepo {
	return &BudgetRepo{pool: pool}
}

func (r *BudgetRepo) Create(ctx context.Context, b *domain.ExpenseBudget) error {
	q := `INSERT INTO expense_budgets (id,company_id,category_id,fiscal_year,period,
		budget_amount,used_amount,reserved_amount,currency,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.CategoryID, b.FiscalYear, b.Period,
		b.BudgetAmount, b.UsedAmount, b.ReservedAmount, b.Currency, b.Notes, b.CreatedBy)
	return repoErr("BudgetRepo.Create", err)
}

func (r *BudgetRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseBudget, error) {
	q := `SELECT id,company_id,category_id,fiscal_year,period,
		budget_amount,used_amount,reserved_amount,currency,notes,created_by,created_at,updated_at
		FROM expense_budgets WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var b domain.ExpenseBudget
	err := row.Scan(&b.ID, &b.CompanyID, &b.CategoryID, &b.FiscalYear, &b.Period,
		&b.BudgetAmount, &b.UsedAmount, &b.ReservedAmount, &b.Currency, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("BudgetRepo.Get", err)
	}
	return &b, nil
}

func (r *BudgetRepo) List(ctx context.Context, companyID uuid.UUID, fiscalYear int) ([]domain.ExpenseBudget, error) {
	q := `SELECT id,company_id,category_id,fiscal_year,period,
		budget_amount,used_amount,reserved_amount,currency,notes,created_by,created_at,updated_at
		FROM expense_budgets WHERE company_id=$1 AND fiscal_year=$2 ORDER BY category_id,period`
	rows, err := r.pool.Query(ctx, q, companyID, fiscalYear)
	if err != nil {
		return nil, repoErr("BudgetRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseBudget, error) {
		var b domain.ExpenseBudget
		err := row.Scan(&b.ID, &b.CompanyID, &b.CategoryID, &b.FiscalYear, &b.Period,
			&b.BudgetAmount, &b.UsedAmount, &b.ReservedAmount, &b.Currency, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
		return b, err
	})
}

func (r *BudgetRepo) UpdateUsage(ctx context.Context, id uuid.UUID, usedDelta, reservedDelta decimal.Decimal) error {
	q := `UPDATE expense_budgets SET used_amount=used_amount+$1,reserved_amount=reserved_amount+$2,updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, usedDelta, reservedDelta, id)
	return repoErr("BudgetRepo.UpdateUsage", err)
}

func (r *BudgetRepo) CheckAvailability(ctx context.Context, budgetID uuid.UUID, amount decimal.Decimal) (bool, error) {
	q := `SELECT (budget_amount - used_amount - reserved_amount) >= $1 FROM expense_budgets WHERE id=$2`
	var available bool
	err := r.pool.QueryRow(ctx, q, amount, budgetID).Scan(&available)
	if err != nil {
		return false, repoErr("BudgetRepo.CheckAvailability", err)
	}
	return available, nil
}
