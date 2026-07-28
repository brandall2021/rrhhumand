package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/shopspring/decimal"
)

type CostService struct {
	costRepo *repository.CostRepo
}

func NewCostService(costRepo *repository.CostRepo) *CostService {
	return &CostService{costRepo: costRepo}
}

func (s *CostService) CreateCost(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, c *domain.BenefitCost) (*domain.BenefitCost, error) {
	c.ID = uuid.New()
	c.CompanyID = companyID
	c.CreatedBy = userID
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if err := s.costRepo.CreateCost(ctx, c); err != nil {
		return nil, svcErr("CreateCost", err)
	}
	return c, nil
}

func (s *CostService) ListCosts(ctx context.Context, companyID uuid.UUID, benefitID uuid.UUID) ([]domain.BenefitCost, error) {
	return s.costRepo.ListCosts(ctx, benefitID)
}

func (s *CostService) UpdateCost(ctx context.Context, companyID uuid.UUID, c *domain.BenefitCost) (*domain.BenefitCost, error) {
	c.CompanyID = companyID
	c.UpdatedAt = time.Now()
	if err := s.costRepo.UpdateCost(ctx, c); err != nil {
		return nil, svcErr("UpdateCost", err)
	}
	return c, nil
}

func (s *CostService) CreateSchedule(ctx context.Context, companyID, benefitID uuid.UUID, date time.Time, amount decimal.Decimal, currency string, notes *string) (*domain.BenefitCostSchedule, error) {
	sch := &domain.BenefitCostSchedule{
		ID:           uuid.New(),
		CompanyID:    companyID,
		BenefitID:    benefitID,
		ScheduleDate: date,
		Amount:       amount,
		Currency:     currency,
		Status:       "PENDING",
		Notes:        notes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.costRepo.CreateSchedule(ctx, sch); err != nil {
		return nil, svcErr("CreateSchedule", err)
	}
	return sch, nil
}

func (s *CostService) ListSchedules(ctx context.Context, companyID, benefitID uuid.UUID, from, to *time.Time) ([]domain.BenefitCostSchedule, error) {
	return s.costRepo.ListSchedules(ctx, benefitID, from, to)
}

func (s *CostService) MarkSchedulePaid(ctx context.Context, id uuid.UUID, reference string) error {
	now := time.Now()
	return s.costRepo.UpdateStatus(ctx, id, "PAID", &now, &reference)
}
