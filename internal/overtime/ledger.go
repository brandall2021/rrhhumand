package overtime

import (
	"context"
)

type Ledger struct {
	repo *Repository
}

func NewLedger(repo *Repository) *Ledger {
	return &Ledger{repo: repo}
}

func (l *Ledger) CreditOvertime(ctx context.Context, companyID, employeeID string, minutes int, overtimeRecordID string, createdBy string) error {
	reason := "Overtime approved"
	return l.repo.CreditTimeBalance(ctx, companyID, employeeID, minutes, overtimeRecordID, reason, createdBy)
}

func (l *Ledger) DebitCompensation(ctx context.Context, companyID, employeeID string, minutes int, reason string, createdBy string) error {
	return l.repo.DebitTimeBalance(ctx, companyID, employeeID, minutes, reason, createdBy)
}

func (l *Ledger) GetTransactions(ctx context.Context, companyID, employeeID string) ([]TimeBalanceTransaction, error) {
	return l.repo.ListBalanceTransactions(ctx, companyID, employeeID)
}
