package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

const (
	AdvanceStatusRequested = "REQUESTED"
	AdvanceStatusApproved  = "APPROVED"
	AdvanceStatusPaid      = "PAID"
	AdvanceStatusSettled   = "SETTLED"
	AdvanceStatusRejected  = "REJECTED"
	AdvanceStatusCancelled = "CANCELLED"
)

type AdvanceRepository interface {
	Create(ctx context.Context, advance *domain.ExpenseAdvance) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseAdvance, error)
	List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseAdvance, error)
	Update(ctx context.Context, advance *domain.ExpenseAdvance) error
}

type AdvanceService struct {
	advanceRepo AdvanceRepository
	auditRepo   AuditRepository
}

func NewAdvanceService(advanceRepo AdvanceRepository, auditRepo AuditRepository) *AdvanceService {
	return &AdvanceService{advanceRepo: advanceRepo, auditRepo: auditRepo}
}

func (s *AdvanceService) RequestAdvance(ctx context.Context, companyID, employeeID, userID uuid.UUID, a *domain.ExpenseAdvance) (*domain.ExpenseAdvance, error) {
	const op = "RequestAdvance"
	now := time.Now()
	a.ID = uuid.New()
	a.CompanyID = companyID
	a.EmployeeID = employeeID
	a.CreatedBy = userID
	a.Status = AdvanceStatusRequested
	a.RequestDate = now
	a.CreatedAt = now
	a.UpdatedAt = now
	if err := s.advanceRepo.Create(ctx, a); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "advance.requested", EntityType: "expense_advance", EntityID: a.ID, CreatedAt: now,
	})
	return a, nil
}

func (s *AdvanceService) GetAdvance(ctx context.Context, id uuid.UUID) (*domain.ExpenseAdvance, error) {
	const op = "GetAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return advance, nil
}

func (s *AdvanceService) ListAdvances(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseAdvance, error) {
	const op = "ListAdvances"
	advances, err := s.advanceRepo.List(ctx, companyID, employeeID, status, limit, offset)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return advances, nil
}

func (s *AdvanceService) ApproveAdvance(ctx context.Context, id, approverID uuid.UUID, approvedAmount decimal.Decimal) error {
	const op = "ApproveAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if advance.Status != AdvanceStatusRequested {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	advance.Status = AdvanceStatusApproved
	advance.ApprovedAmount = &approvedAmount
	advance.ApprovedDate = &now
	advance.UpdatedAt = now
	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: advance.CompanyID, UserID: approverID,
		Action: "advance.approved", EntityType: "expense_advance", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *AdvanceService) PayAdvance(ctx context.Context, id, payerID uuid.UUID) error {
	const op = "PayAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if advance.Status != AdvanceStatusApproved {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	advance.Status = AdvanceStatusPaid
	advance.PaidDate = &now
	advance.UpdatedAt = now
	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: advance.CompanyID, UserID: payerID,
		Action: "advance.paid", EntityType: "expense_advance", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *AdvanceService) SettleAdvance(ctx context.Context, id uuid.UUID, settledAmount decimal.Decimal) error {
	const op = "SettleAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if advance.Status != AdvanceStatusPaid {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	advance.Status = AdvanceStatusSettled
	advance.SettledAmount = settledAmount
	advance.SettledDate = &now
	advance.UpdatedAt = now
	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: advance.CompanyID,
		Action: "advance.settled", EntityType: "expense_advance", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *AdvanceService) RejectAdvance(ctx context.Context, id, approverID uuid.UUID, reason string) error {
	const op = "RejectAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if advance.Status != AdvanceStatusRequested {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	advance.Status = AdvanceStatusRejected
	advance.RejectionReason = &reason
	advance.UpdatedAt = now
	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: advance.CompanyID, UserID: approverID,
		Action: "advance.rejected", EntityType: "expense_advance", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *AdvanceService) CancelAdvance(ctx context.Context, id, userID uuid.UUID) error {
	const op = "CancelAdvance"
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if advance.Status == AdvanceStatusSettled || advance.Status == AdvanceStatusCancelled {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	advance.Status = AdvanceStatusCancelled
	advance.UpdatedAt = now
	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: advance.CompanyID, UserID: userID,
		Action: "advance.cancelled", EntityType: "expense_advance", EntityID: id, CreatedAt: now,
	})
	return nil
}
