package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type PayrollRun struct {
	ID        uuid.UUID  `json:"id"`
	CompanyID uuid.UUID  `json:"company_id"`
	Name      string     `json:"name"`
	PeriodID  uuid.UUID  `json:"period_id"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type PayrollAdapter struct {
	pool *pgxpool.Pool
}

func NewPayrollAdapter(pool *pgxpool.Pool) *PayrollAdapter {
	return &PayrollAdapter{pool: pool}
}

func integErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("expenses_integration.%s: %w", op, err)
}

func (a *PayrollAdapter) SyncReimbursementToPayroll(ctx context.Context, reimbursement *domain.ExpenseReimbursement) error {
	var exists int
	err := a.pool.QueryRow(ctx, `SELECT COUNT(1) FROM benefit_payroll_mappings WHERE entity_type='expense_reimbursement' AND entity_id=$1`, reimbursement.ID).Scan(&exists)
	if err != nil {
		return integErr("SyncReimbursementToPayroll.check", err)
	}

	if exists > 0 {
		_, err = a.pool.Exec(ctx, `
			UPDATE benefit_payroll_mappings SET amount=$1,currency=$2,sync_status='PENDING',updated_at=NOW()
			WHERE entity_type='expense_reimbursement' AND entity_id=$3`,
			reimbursement.Amount, reimbursement.Currency, reimbursement.ID)
	} else {
		_, err = a.pool.Exec(ctx, `
			INSERT INTO benefit_payroll_mappings (id,company_id,employee_id,entity_type,entity_id,amount,currency,sync_status,created_at,updated_at)
			VALUES ($1,$2,$3,'expense_reimbursement',$4,$5,$6,'PENDING',NOW(),NOW())`,
			uuid.New(), reimbursement.CompanyID, reimbursement.EmployeeID, reimbursement.ID, reimbursement.Amount, reimbursement.Currency)
	}
	return integErr("SyncReimbursementToPayroll", err)
}

func (a *PayrollAdapter) GetPayrollRun(ctx context.Context, runID uuid.UUID) (*PayrollRun, error) {
	var run PayrollRun
	err := a.pool.QueryRow(ctx, `SELECT id,company_id,name,period_id,status,created_at FROM payroll_runs WHERE id=$1`, runID).
		Scan(&run.ID, &run.CompanyID, &run.Name, &run.PeriodID, &run.Status, &run.CreatedAt)
	if err != nil {
		return nil, integErr("GetPayrollRun", err)
	}
	return &run, nil
}
