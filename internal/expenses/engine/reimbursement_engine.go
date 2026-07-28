package engine

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/rrhhumand/api/internal/expenses/repository"
	"github.com/shopspring/decimal"
)

type ReimbursementEngine struct {
	expenseRepo *repository.ExpenseRepo
	pool        *pgxpool.Pool
}

func NewReimbursementEngine(er *repository.ExpenseRepo, pool *pgxpool.Pool) *ReimbursementEngine {
	return &ReimbursementEngine{expenseRepo: er, pool: pool}
}

func (e *ReimbursementEngine) CalculateFromReport(ctx context.Context, report *domain.ExpenseReport) (*domain.SettlementResult, error) {
	expenses, err := e.expenseRepo.ListByReport(ctx, report.ID)
	if err != nil {
		return nil, engErr("CalculateFromReport.listExpenses", err)
	}

	total := decimal.Zero
	for _, exp := range expenses {
		if exp.IsReimbursable {
			total = total.Add(exp.BaseAmount)
		}
	}

	companyOwes := decimal.Zero
	employeeOwes := decimal.Zero

	if total.GreaterThan(report.AdvanceAmount) {
		companyOwes = total.Sub(report.AdvanceAmount)
	} else {
		employeeOwes = report.AdvanceAmount.Sub(total)
	}

	return &domain.SettlementResult{
		TotalExpenses:  total,
		AdvanceAmount:  report.AdvanceAmount,
		CompanyOwes:    companyOwes,
		EmployeeOwes:   employeeOwes,
		Currency:       report.Currency,
	}, nil
}

func (e *ReimbursementEngine) GeneratePayrollItem(ctx context.Context, reimbursement *domain.ExpenseReimbursement) error {
	itemID := uuid.New()
	_, err := e.pool.Exec(ctx, `
		INSERT INTO benefit_payroll_mappings (id, company_id, employee_id, entity_type, entity_id, amount, currency, sync_status, created_at, updated_at)
		VALUES ($1,$2,$3,'expense_reimbursement',$4,$5,$6,'PENDING',NOW(),NOW())`,
		itemID, reimbursement.CompanyID, reimbursement.EmployeeID, reimbursement.ID, reimbursement.Amount, reimbursement.Currency)
	return engErr("GeneratePayrollItem", err)
}
