package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

const (
	ReimbursementStatusPending  = "PENDING"
	ReimbursementStatusApproved = "APPROVED"
	ReimbursementStatusPaid     = "PAID"
	ReimbursementStatusRejected = "REJECTED"
	ReimbursementStatusCancelled = "CANCELLED"
)

type ReimbursementRepository interface {
	Create(ctx context.Context, reimbursement *domain.ExpenseReimbursement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseReimbursement, error)
	List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseReimbursement, error)
	Update(ctx context.Context, reimbursement *domain.ExpenseReimbursement) error
}

type ReimbursementService struct {
	reimbursementRepo ReimbursementRepository
	auditRepo         AuditRepository
}

func NewReimbursementService(reimbursementRepo ReimbursementRepository, auditRepo AuditRepository) *ReimbursementService {
	return &ReimbursementService{reimbursementRepo: reimbursementRepo, auditRepo: auditRepo}
}

func (s *ReimbursementService) CreateReimbursement(ctx context.Context, companyID, employeeID, userID uuid.UUID, r *domain.ExpenseReimbursement) (*domain.ExpenseReimbursement, error) {
	const op = "CreateReimbursement"
	now := time.Now()
	r.ID = uuid.New()
	r.CompanyID = companyID
	r.EmployeeID = employeeID
	r.CreatedBy = userID
	r.Status = ReimbursementStatusPending
	r.CreatedAt = now
	r.UpdatedAt = now
	if err := s.reimbursementRepo.Create(ctx, r); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "reimbursement.created", EntityType: "expense_reimbursement", EntityID: r.ID, CreatedAt: now,
	})
	return r, nil
}

func (s *ReimbursementService) GetReimbursement(ctx context.Context, id uuid.UUID) (*domain.ExpenseReimbursement, error) {
	const op = "GetReimbursement"
	r, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return r, nil
}

func (s *ReimbursementService) ListReimbursements(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseReimbursement, error) {
	const op = "ListReimbursements"
	reimbursements, err := s.reimbursementRepo.List(ctx, companyID, employeeID, status, limit, offset)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return reimbursements, nil
}

func (s *ReimbursementService) ApproveReimbursement(ctx context.Context, id, approverID uuid.UUID) error {
	const op = "ApproveReimbursement"
	r, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != ReimbursementStatusPending {
		return svcErr(op, domain.ErrInvalidInput)
	}
	r.Status = ReimbursementStatusApproved
	r.UpdatedAt = time.Now()
	if err := s.reimbursementRepo.Update(ctx, r); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: r.CompanyID, UserID: approverID,
		Action: "reimbursement.approved", EntityType: "expense_reimbursement", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ReimbursementService) PayReimbursement(ctx context.Context, id uuid.UUID, method string, payrollRunID *uuid.UUID) error {
	const op = "PayReimbursement"
	r, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != ReimbursementStatusApproved {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	r.Status = ReimbursementStatusPaid
	r.PaymentMethod = &method
	r.PayrollRunID = payrollRunID
	r.PaidAt = &now
	r.UpdatedAt = now
	if err := s.reimbursementRepo.Update(ctx, r); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: r.CompanyID,
		Action: "reimbursement.paid", EntityType: "expense_reimbursement", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *ReimbursementService) RejectReimbursement(ctx context.Context, id, approverID uuid.UUID, reason string) error {
	const op = "RejectReimbursement"
	r, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != ReimbursementStatusPending {
		return svcErr(op, domain.ErrInvalidInput)
	}
	r.Status = ReimbursementStatusRejected
	r.RejectionReason = &reason
	r.UpdatedAt = time.Now()
	if err := s.reimbursementRepo.Update(ctx, r); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: r.CompanyID, UserID: approverID,
		Action: "reimbursement.rejected", EntityType: "expense_reimbursement", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ReimbursementService) CancelReimbursement(ctx context.Context, id uuid.UUID) error {
	const op = "CancelReimbursement"
	r, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status == ReimbursementStatusPaid || r.Status == ReimbursementStatusCancelled {
		return svcErr(op, domain.ErrInvalidInput)
	}
	r.Status = ReimbursementStatusCancelled
	r.UpdatedAt = time.Now()
	if err := s.reimbursementRepo.Update(ctx, r); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: r.CompanyID,
		Action: "reimbursement.cancelled", EntityType: "expense_reimbursement", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}
