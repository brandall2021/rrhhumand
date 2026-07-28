package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CompensationManager struct {
	repo    *Repository
	balance *BalanceManager
	ledger  *Ledger
}

func NewCompensationManager(repo *Repository) *CompensationManager {
	return &CompensationManager{
		repo:    repo,
		balance: NewBalanceManager(repo),
		ledger:  NewLedger(repo),
	}
}

func (cm *CompensationManager) RequestCompensation(ctx context.Context, companyID, employeeID string, req *RequestCompensationRequest) (*CompensationRequest, error) {
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	balance, err := cm.balance.GetBalance(ctx, companyID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance")
	}

	if balance.BalanceMinutes < req.Minutes {
		return nil, fmt.Errorf("insufficient balance: available %d, requested %d", balance.BalanceMinutes, req.Minutes)
	}

	compensation := &CompensationRequest{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		EmployeeID:   employeeID,
		WorkDate:     workDate,
		Minutes:      req.Minutes,
		Reason:       req.Reason,
		Status:       "PENDING",
		RequestedAt:  time.Now(),
	}

	if err := cm.repo.CreateCompensationRequest(ctx, compensation); err != nil {
		return nil, err
	}

	return compensation, nil
}

func (cm *CompensationManager) ApproveCompensation(ctx context.Context, companyID, requestID, approvedBy string) error {
	comp, err := cm.repo.GetCompensationRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("compensation request not found")
	}

	if comp.Status != "PENDING" {
		return fmt.Errorf("cannot approve compensation in status %s", comp.Status)
	}

	balance, err := cm.balance.GetBalance(ctx, companyID, comp.EmployeeID)
	if err != nil {
		return err
	}

	if balance.BalanceMinutes < comp.Minutes {
		return fmt.Errorf("insufficient balance at approval time")
	}

	now := time.Now()
	comp.Status = "APPROVED"
	comp.ApprovedBy = &approvedBy
	comp.ApprovedAt = &now

	if err := cm.repo.UpdateCompensationRequest(ctx, comp); err != nil {
		return err
	}

	cm.balance.Debit(ctx, companyID, comp.EmployeeID, comp.Minutes, "COMPENSATION_DEBIT", &comp.Reason, &approvedBy)

	return nil
}

func (cm *CompensationManager) RejectCompensation(ctx context.Context, companyID, requestID, reason, rejectedBy string) error {
	comp, err := cm.repo.GetCompensationRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("compensation request not found")
	}

	if comp.Status != "PENDING" {
		return fmt.Errorf("cannot reject compensation in status %s", comp.Status)
	}

	now := time.Now()
	comp.Status = "REJECTED"
	comp.RejectionReason = &reason
	comp.ApprovedBy = &rejectedBy
	comp.ApprovedAt = &now

	return cm.repo.UpdateCompensationRequest(ctx, comp)
}

func (cm *CompensationManager) CancelCompensation(ctx context.Context, companyID, requestID string) error {
	comp, err := cm.repo.GetCompensationRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("compensation request not found")
	}

	if comp.Status != "PENDING" {
		return fmt.Errorf("can only cancel pending compensations")
	}

	comp.Status = "CANCELLED"
	return cm.repo.UpdateCompensationRequest(ctx, comp)
}
