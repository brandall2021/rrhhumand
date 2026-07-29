package integration

import (
	"context"
	"log"
)

type PayrollAdapter struct{}

func NewPayrollAdapter() *PayrollAdapter {
	return &PayrollAdapter{}
}

func (a *PayrollAdapter) StartFinalSettlement(ctx context.Context, companyID, employeeID string, terminationType string, lastWorkingDate string) error {
	log.Printf("[PayrollAdapter] StartFinalSettlement company=%s employee=%s type=%s date=%s", companyID, employeeID, terminationType, lastWorkingDate)
	return nil
}

func (a *PayrollAdapter) GetSettlementStatus(ctx context.Context, employeeID string) (string, error) {
	log.Printf("[PayrollAdapter] GetSettlementStatus employee=%s", employeeID)
	return "PENDING", nil
}
