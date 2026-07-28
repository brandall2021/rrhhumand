package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/shopspring/decimal"
)

type CostEngine struct {
	costRepo *repository.CostRepo
}

func NewCostEngine(costRepo *repository.CostRepo) *CostEngine {
	return &CostEngine{costRepo: costRepo}
}

func (e *CostEngine) CalculateCost(ctx context.Context, benefitID uuid.UUID, employeeCount int) (*decimal.Decimal, error) {
	costs, err := e.costRepo.ListCosts(ctx, benefitID)
	if err != nil {
		return nil, fmt.Errorf("benefits_engine.cost.CalculateCost: %w", err)
	}

	total := decimal.Zero
	for _, c := range costs {
		if c.IsActive {
			total = total.Add(c.TotalCost.Mul(decimal.NewFromInt(int64(employeeCount))))
		}
	}
	return &total, nil
}

func (e *CostEngine) GenerateSchedule(ctx context.Context, benefitID uuid.UUID, from, to time.Time) ([]domain.BenefitCostSchedule, error) {
	schedules, err := e.costRepo.ListSchedules(ctx, benefitID, &from, &to)
	if err != nil {
		return nil, fmt.Errorf("benefits_engine.cost.GenerateSchedule: %w", err)
	}
	return schedules, nil
}
