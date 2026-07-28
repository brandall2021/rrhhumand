package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/shopspring/decimal"
)

type RewardsService struct {
	rewardsRepo *repository.RewardsRepo
	bonusRepo   *repository.BonusRepo
}

func NewRewardsService(rewardsRepo *repository.RewardsRepo, bonusRepo *repository.BonusRepo) *RewardsService {
	return &RewardsService{
		rewardsRepo: rewardsRepo,
		bonusRepo:   bonusRepo,
	}
}

func (s *RewardsService) CreateRewardsItem(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, item *domain.TotalRewardsItem) (*domain.TotalRewardsItem, error) {
	item.ID = uuid.New()
	item.CompanyID = companyID
	item.CreatedBy = userID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	if err := s.rewardsRepo.CreateItem(ctx, item); err != nil {
		return nil, svcErr("CreateRewardsItem", err)
	}
	return item, nil
}

func (s *RewardsService) ListRewardsItems(ctx context.Context, companyID uuid.UUID) ([]domain.TotalRewardsItem, error) {
	return s.rewardsRepo.ListItems(ctx, companyID)
}

func (s *RewardsService) UpdateRewardsItem(ctx context.Context, companyID uuid.UUID, item *domain.TotalRewardsItem) (*domain.TotalRewardsItem, error) {
	item.CompanyID = companyID
	item.UpdatedAt = time.Now()
	if err := s.rewardsRepo.UpdateItem(ctx, item); err != nil {
		return nil, svcErr("UpdateRewardsItem", err)
	}
	return item, nil
}

func (s *RewardsService) GenerateSnapshot(ctx context.Context, companyID, employeeID, generatedBy uuid.UUID, fiscalYear int, periodName string) (*domain.TotalRewardsSnapshot, error) {
	snapshot := &domain.TotalRewardsSnapshot{
		ID:                  uuid.New(),
		CompanyID:           companyID,
		EmployeeID:          employeeID,
		SnapshotDate:        time.Now(),
		FiscalYear:          fiscalYear,
		PeriodName:          &periodName,
		BaseSalary:          decimal.Zero,
		VariablePay:         decimal.Zero,
		BonusesTotal:        decimal.Zero,
		IncentivesTotal:     decimal.Zero,
		BenefitsTotal:       decimal.Zero,
		EmployerContributions: decimal.Zero,
		FlexibleSpending:    decimal.Zero,
		InsuranceValue:      decimal.Zero,
		DevelopmentValue:    decimal.Zero,
		WellnessValue:       decimal.Zero,
		RecognitionValue:    decimal.Zero,
		PerksValue:          decimal.Zero,
		TotalRewards:        decimal.Zero,
		Currency:            "ARS",
		GeneratedBy:         generatedBy,
		GeneratedAt:         time.Now(),
		CreatedAt:           time.Now(),
	}
	if err := s.rewardsRepo.CreateSnapshot(ctx, snapshot); err != nil {
		return nil, svcErr("GenerateSnapshot", err)
	}
	return snapshot, nil
}

func (s *RewardsService) GetLatestSnapshot(ctx context.Context, companyID, employeeID uuid.UUID) (*domain.TotalRewardsSnapshot, error) {
	return s.rewardsRepo.GetLatestSnapshot(ctx, employeeID)
}

func (s *RewardsService) ListSnapshots(ctx context.Context, companyID uuid.UUID, fiscalYear *int) ([]domain.TotalRewardsSnapshot, error) {
	fy := 0
	if fiscalYear != nil {
		fy = *fiscalYear
	}
	return s.rewardsRepo.ListSnapshots(ctx, companyID, fy)
}

func (s *RewardsService) CreateReportDefinition(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, d *domain.BenefitReportDefinition) (*domain.BenefitReportDefinition, error) {
	d.ID = uuid.New()
	d.CompanyID = companyID
	d.CreatedBy = userID
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	if err := s.rewardsRepo.CreateReportDefinition(ctx, d); err != nil {
		return nil, svcErr("CreateReportDefinition", err)
	}
	return d, nil
}

func (s *RewardsService) ListReportDefinitions(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitReportDefinition, error) {
	return s.rewardsRepo.ListReportDefinitions(ctx, companyID)
}

func (s *RewardsService) LogNotification(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, notifType, channel, title, body string, metadata map[string]any) (*domain.BenefitNotificationLog, error) {
	n := &domain.BenefitNotificationLog{
		ID:               uuid.New(),
		CompanyID:        companyID,
		EmployeeID:       employeeID,
		NotificationType: notifType,
		Channel:          channel,
		Title:            title,
		Body:             &body,
		Metadata:         metadata,
		SentAt:           time.Now(),
		CreatedAt:        time.Now(),
	}
	if err := s.rewardsRepo.CreateNotification(ctx, n); err != nil {
		return nil, svcErr("LogNotification", err)
	}
	return n, nil
}

func (s *RewardsService) ListNotifications(ctx context.Context, employeeID uuid.UUID, notifType *string, limit, offset int) ([]domain.BenefitNotificationLog, error) {
	return s.rewardsRepo.ListNotifications(ctx, &employeeID, nil, limit, offset)
}

func (s *RewardsService) MarkNotificationRead(ctx context.Context, id uuid.UUID) error {
	return s.rewardsRepo.MarkRead(ctx, id)
}
