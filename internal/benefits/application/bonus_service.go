package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type BonusService struct {
	bonusRepo *repository.BonusRepo
}

func NewBonusService(bonusRepo *repository.BonusRepo) *BonusService {
	return &BonusService{bonusRepo: bonusRepo}
}

func (s *BonusService) CreateBonus(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, b *domain.EmployeeBonus) (*domain.EmployeeBonus, error) {
	b.ID = uuid.New()
	b.CompanyID = companyID
	b.CreatedBy = userID
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	if err := s.bonusRepo.CreateBonus(ctx, b); err != nil {
		return nil, svcErr("CreateBonus", err)
	}
	return b, nil
}

func (s *BonusService) GetBonus(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeBonus, error) {
	return s.bonusRepo.GetBonus(ctx, companyID, id)
}

func (s *BonusService) ListBonuses(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string) ([]domain.EmployeeBonus, error) {
	if employeeID == nil {
		return []domain.EmployeeBonus{}, nil
	}
	return s.bonusRepo.ListBonuses(ctx, *employeeID, status)
}

func (s *BonusService) UpdateBonus(ctx context.Context, companyID uuid.UUID, b *domain.EmployeeBonus) (*domain.EmployeeBonus, error) {
	b.CompanyID = companyID
	b.UpdatedAt = time.Now()
	if err := s.bonusRepo.UpdateBonus(ctx, b); err != nil {
		return nil, svcErr("UpdateBonus", err)
	}
	return b, nil
}

func (s *BonusService) ApproveBonus(ctx context.Context, id, approvedBy uuid.UUID) error {
	b, err := s.bonusRepo.GetBonus(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("ApproveBonus", err)
	}
	now := time.Now()
	b.Status = "APPROVED"
	b.ApprovedBy = &approvedBy
	b.ApprovedAt = &now
	b.UpdatedAt = now
	if err := s.bonusRepo.UpdateBonus(ctx, b); err != nil {
		return svcErr("ApproveBonus", err)
	}
	return nil
}

func (s *BonusService) PayBonus(ctx context.Context, id, payrollRunID uuid.UUID) error {
	b, err := s.bonusRepo.GetBonus(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("PayBonus", err)
	}
	now := time.Now()
	b.Status = "PAID"
	b.PaidInPayroll = true
	b.PayrollRunID = &payrollRunID
	b.PaymentDate = &now
	b.UpdatedAt = now
	if err := s.bonusRepo.UpdateBonus(ctx, b); err != nil {
		return svcErr("PayBonus", err)
	}
	return nil
}

func (s *BonusService) CancelBonus(ctx context.Context, id uuid.UUID, reason string) error {
	b, err := s.bonusRepo.GetBonus(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("CancelBonus", err)
	}
	b.Status = "CANCELLED"
	b.Notes = &reason
	b.UpdatedAt = time.Now()
	if err := s.bonusRepo.UpdateBonus(ctx, b); err != nil {
		return svcErr("CancelBonus", err)
	}
	return nil
}

func (s *BonusService) CreateIncentive(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, i *domain.EmployeeIncentive) (*domain.EmployeeIncentive, error) {
	i.ID = uuid.New()
	i.CompanyID = companyID
	if employeeID != nil {
		i.EmployeeID = *employeeID
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	if err := s.bonusRepo.CreateIncentive(ctx, i); err != nil {
		return nil, svcErr("CreateIncentive", err)
	}
	return i, nil
}

func (s *BonusService) GetIncentive(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeIncentive, error) {
	return s.bonusRepo.GetIncentive(ctx, companyID, id)
}

func (s *BonusService) ListIncentives(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string) ([]domain.EmployeeIncentive, error) {
	if employeeID == nil {
		return []domain.EmployeeIncentive{}, nil
	}
	return s.bonusRepo.ListIncentives(ctx, *employeeID, status)
}

func (s *BonusService) RedeemIncentive(ctx context.Context, id uuid.UUID) error {
	i, err := s.bonusRepo.GetIncentive(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("RedeemIncentive", err)
	}
	now := time.Now()
	i.Status = "REDEEMED"
	i.RedemptionDate = &now
	i.UpdatedAt = now
	if err := s.bonusRepo.UpdateIncentive(ctx, i); err != nil {
		return svcErr("RedeemIncentive", err)
	}
	return nil
}

func (s *BonusService) CreatePayrollMapping(ctx context.Context, m *domain.BenefitPayrollMapping) (*domain.BenefitPayrollMapping, error) {
	m.ID = uuid.New()
	m.SyncStatus = "PENDING"
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	if err := s.bonusRepo.CreatePayrollMapping(ctx, m); err != nil {
		return nil, svcErr("CreatePayrollMapping", err)
	}
	return m, nil
}

func (s *BonusService) ListPayrollMappings(ctx context.Context, companyID uuid.UUID, benefitID *uuid.UUID, mappingType *string) ([]domain.BenefitPayrollMapping, error) {
	if benefitID == nil {
		return []domain.BenefitPayrollMapping{}, nil
	}
	return s.bonusRepo.ListPayrollMappings(ctx, *benefitID, mappingType)
}

func (s *BonusService) SyncToPayroll(ctx context.Context, id uuid.UUID) error {
	pm, err := s.bonusRepo.GetPayrollMapping(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("SyncToPayroll", err)
	}
	pm.SyncStatus = "SYNCED"
	now := time.Now()
	pm.LastSyncedAt = &now
	pm.UpdatedAt = now
	if err := s.bonusRepo.UpdateSyncStatus(ctx, id, "SYNCED", nil); err != nil {
		return svcErr("SyncToPayroll", err)
	}
	return nil
}
