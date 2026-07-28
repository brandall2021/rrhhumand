package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/rrhhumand/api/internal/payroll/features/repository"
)

type BankService struct {
	bankRepo *repository.BankRepo
}

func NewBankService(bankRepo *repository.BankRepo) *BankService {
	return &BankService{bankRepo: bankRepo}
}

func bankSvcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("bank_svc.%s: %w", op, err)
}

func (s *BankService) CreateBatch(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, bankCode, paymentType string, paymentDate time.Time, userID uuid.UUID) (*domain.BankBatch, error) {
	b := &domain.BankBatch{
		ID:             uuid.New(),
		CompanyID:      companyID,
		RunID:          runID,
		BatchNumber:    fmt.Sprintf("BATCH-%s", uuid.New().String()[:8]),
		BankCode:       bankCode,
		PaymentType:    paymentType,
		TotalAmount:    decimal.Zero,
		TotalEmployees: 0,
		Currency:       "ARS",
		PaymentDate:    paymentDate,
		Status:         "DRAFT",
		FileFormat:     "TXT",
		GeneratedBy:    userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.bankRepo.CreateBatch(ctx, b); err != nil {
		return nil, bankSvcErr("CreateBatch", err)
	}
	return b, nil
}

func (s *BankService) GetBatch(ctx context.Context, companyID, batchID uuid.UUID) (*domain.BankBatch, error) {
	return s.bankRepo.GetBatch(ctx, companyID, batchID)
}

func (s *BankService) ListBatches(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, status string) ([]domain.BankBatch, error) {
	_ = status
	return s.bankRepo.ListBatches(ctx, companyID, runID, 0, 0)
}

func (s *BankService) GenerateBankFile(ctx context.Context, batchID uuid.UUID, format string) error {
	fields := map[string]any{
		"file_format": format,
	}
	return s.bankRepo.UpdateBatchStatus(ctx, batchID, "FILE_GENERATED", fields)
}

func (s *BankService) SendBatch(ctx context.Context, batchID uuid.UUID) error {
	fields := map[string]any{
		"sent_at": time.Now(),
	}
	return s.bankRepo.UpdateBatchStatus(ctx, batchID, "SENT", fields)
}

func (s *BankService) UpdateItemStatus(ctx context.Context, itemID uuid.UUID, status string, errorMsg string) error {
	fields := map[string]any{}
	if errorMsg != "" {
		fields["error_message"] = errorMsg
	}
	return s.bankRepo.UpdateBatchItemStatus(ctx, itemID, status, fields)
}

func (s *BankService) GetBatchItems(ctx context.Context, batchID uuid.UUID) ([]domain.BankBatchItem, error) {
	return s.bankRepo.ListBatchItems(ctx, batchID)
}

func (s *BankService) GetBatchSummary(ctx context.Context, batchID uuid.UUID) (*repository.BatchSummary, error) {
	return s.bankRepo.GetBatchSummary(ctx, batchID)
}
