package integration

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

type BenefitWallet struct {
	ID         uuid.UUID       `json:"id"`
	EmployeeID uuid.UUID       `json:"employee_id"`
	PlanID     uuid.UUID       `json:"plan_id"`
	Balance    decimal.Decimal `json:"balance"`
	Currency   string          `json:"currency"`
	ExpiresAt  *time.Time      `json:"expires_at,omitempty"`
}

type BenefitsAdapter struct {
	pool *pgxpool.Pool
}

func NewBenefitsAdapter(pool *pgxpool.Pool) *BenefitsAdapter {
	return &BenefitsAdapter{pool: pool}
}

func (a *BenefitsAdapter) GetBenefitWallets(ctx context.Context, employeeID uuid.UUID) ([]BenefitWallet, error) {
	rows, err := a.pool.Query(ctx, `
		SELECT id,employee_id,plan_id,balance,currency,expires_at
		FROM employee_benefit_wallets WHERE employee_id=$1 AND balance > 0`, employeeID)
	if err != nil {
		return nil, integErr("GetBenefitWallets", err)
	}
	defer rows.Close()

	var wallets []BenefitWallet
	for rows.Next() {
		var w BenefitWallet
		if err := rows.Scan(&w.ID, &w.EmployeeID, &w.PlanID, &w.Balance, &w.Currency, &w.ExpiresAt); err != nil {
			return nil, integErr("GetBenefitWallets.scan", err)
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (a *BenefitsAdapter) CreateReimbursementFromBenefit(ctx context.Context, expense *domain.Expense) error {
	_, err := a.pool.Exec(ctx, `
		INSERT INTO benefit_reimbursements (id,company_id,employee_id,amount,currency,expense_id,status,created_at)
		VALUES ($1,$2,$3,$4,$5,$6,'PENDING',NOW())`,
		uuid.New(), expense.CompanyID, expense.EmployeeID, expense.BaseAmount, expense.BaseCurrency, expense.ID)
	return integErr("CreateReimbursementFromBenefit", err)
}
