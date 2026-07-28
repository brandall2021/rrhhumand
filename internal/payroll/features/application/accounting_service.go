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

type AccountingService struct {
	accountingRepo *repository.AccountingRepo
}

func NewAccountingService(accountingRepo *repository.AccountingRepo) *AccountingService {
	return &AccountingService{accountingRepo: accountingRepo}
}

func acctSvcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("accounting_svc.%s: %w", op, err)
}

func (s *AccountingService) CreateMapping(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, m *domain.AccountingAccountMapping) (*domain.AccountingAccountMapping, error) {
	m.ID = uuid.New()
	m.CompanyID = companyID
	m.CreatedBy = userID
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	if err := s.accountingRepo.CreateMapping(ctx, m); err != nil {
		return nil, acctSvcErr("CreateMapping", err)
	}
	return m, nil
}

func (s *AccountingService) GetMapping(ctx context.Context, companyID, id uuid.UUID) (*domain.AccountingAccountMapping, error) {
	return s.accountingRepo.GetMapping(ctx, companyID, id)
}

func (s *AccountingService) ListMappings(ctx context.Context, companyID uuid.UUID) ([]domain.AccountingAccountMapping, error) {
	return s.accountingRepo.ListMappings(ctx, companyID)
}

func (s *AccountingService) UpdateMapping(ctx context.Context, companyID uuid.UUID, m *domain.AccountingAccountMapping) (*domain.AccountingAccountMapping, error) {
	m.CompanyID = companyID
	m.UpdatedAt = time.Now()
	if err := s.accountingRepo.UpdateMapping(ctx, m); err != nil {
		return nil, acctSvcErr("UpdateMapping", err)
	}
	return m, nil
}

func (s *AccountingService) GetActiveMappings(ctx context.Context, companyID uuid.UUID, conceptID uuid.UUID, date time.Time) ([]domain.AccountingAccountMapping, error) {
	_ = conceptID
	return s.accountingRepo.GetActiveForConcept(ctx, companyID, date)
}

func (s *AccountingService) GenerateAccountingExport(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, exportType, fileFormat string, userID uuid.UUID) (*domain.AccountingExport, error) {
	e := &domain.AccountingExport{
		ID:            uuid.New(),
		CompanyID:     companyID,
		RunID:         runID,
		ExportType:    exportType,
		FileFormat:    fileFormat,
		FileName:      fmt.Sprintf("ACCT_%s_%s.%s", exportType, time.Now().Format("20060102_150405"), fileFormat),
		Status:        "GENERATED",
		EmployeeCount: 0,
		TotalDebit:    decimal.Zero,
		TotalCredit:   decimal.Zero,
		EntryCount:    0,
		GeneratedBy:   userID,
		GeneratedAt:   time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := s.accountingRepo.CreateExport(ctx, e); err != nil {
		return nil, acctSvcErr("GenerateAccountingExport", err)
	}
	return e, nil
}

func (s *AccountingService) GetExport(ctx context.Context, companyID, exportID uuid.UUID) (*domain.AccountingExport, error) {
	return s.accountingRepo.GetExport(ctx, companyID, exportID)
}

func (s *AccountingService) ListExports(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID) ([]domain.AccountingExport, error) {
	return s.accountingRepo.ListExports(ctx, companyID, runID, 0, 0)
}

func (s *AccountingService) GetEntries(ctx context.Context, exportID uuid.UUID) ([]domain.AccountingEntry, error) {
	return s.accountingRepo.ListEntriesByExport(ctx, exportID)
}
