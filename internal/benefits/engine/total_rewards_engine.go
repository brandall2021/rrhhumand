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

type TotalRewardsEngine struct {
	bonusRepo      *repository.BonusRepo
	assignmentRepo *repository.AssignmentRepo
	costRepo       *repository.CostRepo
}

func NewTotalRewardsEngine(bonusRepo *repository.BonusRepo, assignmentRepo *repository.AssignmentRepo, costRepo *repository.CostRepo) *TotalRewardsEngine {
	return &TotalRewardsEngine{
		bonusRepo:      bonusRepo,
		assignmentRepo: assignmentRepo,
		costRepo:       costRepo,
	}
}

func (e *TotalRewardsEngine) CalculateSnapshot(ctx context.Context, companyID, employeeID uuid.UUID, fiscalYear int) (*domain.TotalRewardsSnapshot, error) {
	bonuses, err := e.bonusRepo.ListBonuses(ctx, employeeID, nil)
	if err != nil {
		return nil, fmt.Errorf("benefits_engine.rewards.CalculateSnapshot: %w", err)
	}

	bonusTotal := decimal.Zero
	for _, b := range bonuses {
		if b.Status == "APPROVED" || b.Status == "PAID" {
			bonusTotal = bonusTotal.Add(b.Amount)
		}
	}

	benefits, err := e.assignmentRepo.List(ctx, &companyID, &employeeID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("benefits_engine.rewards.CalculateSnapshot: %w", err)
	}

	benefitsTotal := decimal.Zero
	for _, eb := range benefits {
		if eb.Status == "ACTIVE" {
			benefitsTotal = benefitsTotal.Add(eb.EmployerCost)
		}
	}

	now := time.Now()
	snapshot := &domain.TotalRewardsSnapshot{
		ID:                  uuid.New(),
		CompanyID:           companyID,
		EmployeeID:          employeeID,
		SnapshotDate:        now,
		FiscalYear:          fiscalYear,
		BonusesTotal:        bonusTotal,
		BenefitsTotal:       benefitsTotal,
		TotalRewards:        bonusTotal.Add(benefitsTotal),
		Currency:            "ARS",
		GeneratedBy:         employeeID,
		GeneratedAt:         now,
		CreatedAt:           now,
	}

	return snapshot, nil
}

func (e *TotalRewardsEngine) GetBreakdown(ctx context.Context, snapshotID uuid.UUID) (map[string]decimal.Decimal, error) {
	return map[string]decimal.Decimal{
		"base_salary": decimal.Zero,
		"bonuses":     decimal.Zero,
		"benefits":    decimal.Zero,
		"incentives":  decimal.Zero,
		"total":       decimal.Zero,
	}, nil
}
