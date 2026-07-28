package application

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/rrhhumand/api/internal/payroll/features/repository"
)

type ReceiptService struct {
	receiptRepo *repository.ReceiptRepo
}

func NewReceiptService(receiptRepo *repository.ReceiptRepo) *ReceiptService {
	return &ReceiptService{receiptRepo: receiptRepo}
}

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("receipt_svc.%s: %w", op, err)
}

func (s *ReceiptService) CreateTemplate(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, t *domain.ReceiptTemplate) (*domain.ReceiptTemplate, error) {
	t.ID = uuid.New()
	t.CompanyID = companyID
	t.CreatedBy = userID
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if err := s.receiptRepo.CreateTemplate(ctx, t); err != nil {
		return nil, svcErr("CreateTemplate", err)
	}
	return t, nil
}

func (s *ReceiptService) GetTemplate(ctx context.Context, companyID, id uuid.UUID) (*domain.ReceiptTemplate, error) {
	return s.receiptRepo.GetTemplate(ctx, companyID, id)
}

func (s *ReceiptService) ListTemplates(ctx context.Context, companyID uuid.UUID) ([]domain.ReceiptTemplate, error) {
	return s.receiptRepo.ListTemplates(ctx, companyID)
}

func (s *ReceiptService) UpdateTemplate(ctx context.Context, companyID uuid.UUID, t *domain.ReceiptTemplate) (*domain.ReceiptTemplate, error) {
	t.CompanyID = companyID
	t.UpdatedAt = time.Now()
	if err := s.receiptRepo.UpdateTemplate(ctx, t); err != nil {
		return nil, svcErr("UpdateTemplate", err)
	}
	return t, nil
}

func (s *ReceiptService) DeleteTemplate(ctx context.Context, companyID, id uuid.UUID) error {
	return s.receiptRepo.DeleteTemplate(ctx, companyID, id)
}

func (s *ReceiptService) GenerateReceipts(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, employeeIDs []uuid.UUID, userID uuid.UUID) ([]domain.Receipt, error) {
	var receipts []domain.Receipt
	now := time.Now()
	for _, empID := range employeeIDs {
		rec := &domain.Receipt{
			ID:             uuid.New(),
			CompanyID:      companyID,
			RunID:          runID,
			EmployeeID:     empID,
			Status:         "GENERATED",
			GeneratedBy:    userID,
			GeneratedAt:    now,
			CreatedAt:      now,
			UpdatedAt:      now,
			Currency:       "ARS",
			GrossRemunerative: decimal.Zero,
			GrossNonRemunerative: decimal.Zero,
			DeductionsTotal: decimal.Zero,
			ContributionsTotal: decimal.Zero,
			NetAmount:      decimal.Zero,
			EmployerCost:   decimal.Zero,
		}
		if err := s.receiptRepo.CreateReceipt(ctx, rec); err != nil {
			return nil, svcErr("GenerateReceipts", err)
		}
		receipts = append(receipts, *rec)
	}
	return receipts, nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, companyID, id uuid.UUID) (*domain.Receipt, error) {
	return s.receiptRepo.GetReceipt(ctx, companyID, id)
}

func (s *ReceiptService) ListReceipts(ctx context.Context, companyID uuid.UUID, runID, employeeID *uuid.UUID, limit, offset int) ([]domain.Receipt, error) {
	return s.receiptRepo.ListReceipts(ctx, companyID, runID, employeeID, limit, offset)
}

func (s *ReceiptService) GetReceiptItems(ctx context.Context, receiptID uuid.UUID) ([]domain.ReceiptItem, error) {
	return s.receiptRepo.ListReceiptItems(ctx, receiptID)
}

func (s *ReceiptService) AcknowledgeReceipt(ctx context.Context, receiptID uuid.UUID, ip string) error {
	fields := map[string]any{
		"acknowledged_at": time.Now(),
		"acknowledged_ip": ip,
	}
	return s.receiptRepo.UpdateReceiptStatus(ctx, receiptID, "ACKNOWLEDGED", fields)
}

func (s *ReceiptService) AmountInWords(amount decimal.Decimal, currency string) string {
	if amount.IsZero() {
		return "CERO PESOS"
	}
	intPart := amount.IntPart()
	_ = math.Abs(float64(intPart))
	_ = strings.Builder{}
	return fmt.Sprintf("SON %d %s", intPart, currency)
}
