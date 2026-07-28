package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/shopspring/decimal"
)

type ReimbursementService struct {
	reimbursementRepo *repository.ReimbursementRepo
	walletRepo        *repository.WalletRepo
}

func NewReimbursementService(reimbursementRepo *repository.ReimbursementRepo, walletRepo *repository.WalletRepo) *ReimbursementService {
	return &ReimbursementService{
		reimbursementRepo: reimbursementRepo,
		walletRepo:        walletRepo,
	}
}

func (s *ReimbursementService) CreateReimbursement(ctx context.Context, companyID, employeeID uuid.UUID, r *domain.BenefitReimbursement) (*domain.BenefitReimbursement, error) {
	r.ID = uuid.New()
	r.CompanyID = companyID
	r.EmployeeID = employeeID
	r.Status = "PENDING"
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	if err := s.reimbursementRepo.Create(ctx, r); err != nil {
		return nil, svcErr("CreateReimbursement", err)
	}
	return r, nil
}

func (s *ReimbursementService) GetReimbursement(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitReimbursement, error) {
	return s.reimbursementRepo.Get(ctx, companyID, id)
}

func (s *ReimbursementService) ListReimbursements(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string) ([]domain.BenefitReimbursement, error) {
	return s.reimbursementRepo.List(ctx, &companyID, employeeID, nil, status)
}

func (s *ReimbursementService) ApproveReimbursement(ctx context.Context, id, reviewerID uuid.UUID, approvedAmount *decimal.Decimal) error {
	r, err := s.reimbursementRepo.Get(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("ApproveReimbursement", err)
	}
	r.Status = "APPROVED"
	r.ApprovedAmount = approvedAmount
	r.ReviewedBy = &reviewerID
	now := time.Now()
	r.ReviewedAt = &now
	r.UpdatedAt = now
	if err := s.reimbursementRepo.UpdateStatus(ctx, id, "APPROVED", &reviewerID); err != nil {
		return svcErr("ApproveReimbursement", err)
	}
	return nil
}

func (s *ReimbursementService) RejectReimbursement(ctx context.Context, id, reviewerID uuid.UUID, reason string) error {
	r, err := s.reimbursementRepo.Get(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("RejectReimbursement", err)
	}
	now := time.Now()
	r.Status = "REJECTED"
	r.RejectionReason = &reason
	r.ReviewedBy = &reviewerID
	r.ReviewedAt = &now
	r.UpdatedAt = now
	if err := s.reimbursementRepo.UpdateStatus(ctx, id, "REJECTED", &reviewerID); err != nil {
		return svcErr("RejectReimbursement", err)
	}
	return nil
}

func (s *ReimbursementService) PayReimbursement(ctx context.Context, id uuid.UUID, method, reference string) error {
	r, err := s.reimbursementRepo.Get(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("PayReimbursement", err)
	}
	now := time.Now()
	r.Status = "PAID"
	r.PaymentMethod = &method
	r.PaymentReference = &reference
	r.PaidAt = &now
	r.UpdatedAt = now
	if err := s.reimbursementRepo.UpdateStatus(ctx, id, "PAID", nil); err != nil {
		return svcErr("PayReimbursement", err)
	}
	if r.WalletID != nil {
		if err := s.walletRepo.UpdateBalance(ctx, *r.WalletID, decimal.Zero); err != nil {
			return svcErr("PayReimbursement", err)
		}
	}
	return nil
}

func (s *ReimbursementService) CancelReimbursement(ctx context.Context, id uuid.UUID) error {
	r, err := s.reimbursementRepo.Get(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("CancelReimbursement", err)
	}
	r.Status = "CANCELLED"
	r.UpdatedAt = time.Now()
	if err := s.reimbursementRepo.UpdateStatus(ctx, id, "CANCELLED", nil); err != nil {
		return svcErr("CancelReimbursement", err)
	}
	return nil
}

func (s *ReimbursementService) UploadDocument(ctx context.Context, reimbursementID, uploadedBy uuid.UUID, fileName, fileType, storagePath string, fileSize int) (*domain.BenefitReimbursementDocument, error) {
	d := &domain.BenefitReimbursementDocument{
		ID:              uuid.New(),
		ReimbursementID: reimbursementID,
		FileName:        fileName,
		FileType:        fileType,
		FileSize:        fileSize,
		StoragePath:     storagePath,
		UploadedBy:      uploadedBy,
		UploadedAt:      time.Now(),
	}
	if err := s.reimbursementRepo.CreateDocument(ctx, d); err != nil {
		return nil, svcErr("UploadDocument", err)
	}
	return d, nil
}

func (s *ReimbursementService) ListDocuments(ctx context.Context, reimbursementID uuid.UUID) ([]domain.BenefitReimbursementDocument, error) {
	return s.reimbursementRepo.ListDocuments(ctx, reimbursementID)
}
