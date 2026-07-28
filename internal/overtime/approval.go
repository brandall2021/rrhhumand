package overtime

import (
	"context"
	"fmt"
	"time"
)

type ApprovalManager struct {
	repo    *Repository
	balance *BalanceManager
	ledger  *Ledger
}

func NewApprovalManager(repo *Repository) *ApprovalManager {
	return &ApprovalManager{
		repo:    repo,
		balance: NewBalanceManager(repo),
		ledger:  NewLedger(repo),
	}
}

func (am *ApprovalManager) Approve(ctx context.Context, companyID, recordID string, approvedMinutes int, approvedBy string) error {
	record, err := am.repo.GetOvertimeRecord(ctx, companyID, recordID)
	if err != nil {
		return fmt.Errorf("overtime record not found")
	}

	if record.Status != "DETECTED" && record.Status != "PENDING" && record.Status != "REQUESTED" && record.Status != "SUBMITTED" {
		return fmt.Errorf("cannot approve record in status %s", record.Status)
	}

	if approvedMinutes > record.OvertimeMinutes {
		return fmt.Errorf("approved_minutes cannot exceed overtime_minutes")
	}

	if approvedMinutes < 0 {
		return fmt.Errorf("approved_minutes cannot be negative")
	}

	if approvedMinutes == 0 {
		approvedMinutes = record.OvertimeMinutes
	}

	record.ApprovedMinutes = approvedMinutes
	record.Status = "APPROVED"
	record.UpdatedAt = time.Now()

	if err := am.repo.UpdateOvertimeRecord(ctx, record); err != nil {
		return err
	}

	am.ledger.CreditOvertime(ctx, companyID, record.EmployeeID, approvedMinutes, recordID, approvedBy)

	return nil
}

func (am *ApprovalManager) Reject(ctx context.Context, companyID, recordID, reason, rejectedBy string) error {
	record, err := am.repo.GetOvertimeRecord(ctx, companyID, recordID)
	if err != nil {
		return fmt.Errorf("overtime record not found")
	}

	if record.Status != "DETECTED" && record.Status != "PENDING" && record.Status != "REQUESTED" && record.Status != "SUBMITTED" {
		return fmt.Errorf("cannot reject record in status %s", record.Status)
	}

	record.Status = "REJECTED"
	record.RejectionReason = &reason
	record.UpdatedAt = time.Now()

	return am.repo.UpdateOvertimeRecord(ctx, record)
}

func (am *ApprovalManager) ApproveRequest(ctx context.Context, companyID, requestID string, approvedMinutes int, approvedBy string) error {
	req, err := am.repo.GetOvertimeRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("overtime request not found")
	}

	if req.Status != "PENDING" && req.Status != "SUBMITTED" {
		return fmt.Errorf("cannot approve request in status %s", req.Status)
	}

	if approvedMinutes <= 0 {
		approvedMinutes = req.RequestedMinutes
	}

	if approvedMinutes > req.RequestedMinutes {
		return fmt.Errorf("approved_minutes cannot exceed requested_minutes")
	}

	now := time.Now()
	req.ApprovedMinutes = approvedMinutes
	req.Status = "APPROVED"
	req.ApprovedBy = &approvedBy
	req.ApprovedAt = &now

	if err := am.repo.UpdateOvertimeRequest(ctx, req); err != nil {
		return err
	}

	if req.OvertimeRecordID != nil {
		am.Approve(ctx, companyID, *req.OvertimeRecordID, approvedMinutes, approvedBy)
	}

	return nil
}

func (am *ApprovalManager) RejectRequest(ctx context.Context, companyID, requestID, reason, rejectedBy string) error {
	req, err := am.repo.GetOvertimeRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("overtime request not found")
	}

	if req.Status != "PENDING" && req.Status != "SUBMITTED" {
		return fmt.Errorf("cannot reject request in status %s", req.Status)
	}

	now := time.Now()
	req.Status = "REJECTED"
	req.RejectionReason = &reason
	req.ApprovedBy = &rejectedBy
	req.ApprovedAt = &now

	return am.repo.UpdateOvertimeRequest(ctx, req)
}
