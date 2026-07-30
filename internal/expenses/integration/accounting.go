package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

type JournalEntry struct {
	AccountCode string          `json:"account_code"`
	Description string          `json:"description"`
	Debit       decimal.Decimal `json:"debit"`
	Credit      decimal.Decimal `json:"credit"`
	CostCenter  *string         `json:"cost_center,omitempty"`
	ProjectID   *uuid.UUID      `json:"project_id,omitempty"`
}

type AccountingAdapter struct {
	pool *pgxpool.Pool
}

func NewAccountingAdapter(pool *pgxpool.Pool) *AccountingAdapter {
	return &AccountingAdapter{pool: pool}
}

func (a *AccountingAdapter) CreateExpenseEntry(ctx context.Context, expense *domain.Expense) error {
	_, err := a.pool.Exec(ctx, `
		INSERT INTO expense_integration_logs (id,company_id,expense_id,integration_type,status,payload,created_at)
		VALUES ($1,$2,$3,'ACCOUNTING','PENDING',$4,NOW())`,
		uuid.New(), expense.CompanyID, expense.ID, map[string]any{
			"expense_id":    expense.ID,
			"amount":        expense.BaseAmount,
			"currency":      expense.BaseCurrency,
			"category_id":   expense.CategoryID,
			"description":   expense.Description,
			"expense_date":  expense.ExpenseDate,
		})
	return integErr("CreateExpenseEntry", err)
}

func (a *AccountingAdapter) CreateJournalEntry(ctx context.Context, entries []JournalEntry) error {
	for _, entry := range entries {
		_, err := a.pool.Exec(ctx, `
			INSERT INTO expense_integration_logs (id,integration_type,status,payload,created_at)
			VALUES ($1,'ACCOUNTING','PENDING',$2,NOW())`,
			uuid.New(), map[string]any{
				"account_code": entry.AccountCode,
				"description":  entry.Description,
				"debit":        entry.Debit,
				"credit":       entry.Credit,
				"cost_center":  entry.CostCenter,
				"project_id":   entry.ProjectID,
			})
		if err != nil {
			return integErr("CreateJournalEntry", err)
		}
	}
	return nil
}
