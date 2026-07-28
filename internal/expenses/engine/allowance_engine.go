package engine

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/rrhhumand/api/internal/expenses/repository"
	"github.com/shopspring/decimal"
)

type AllowanceEngine struct {
	allowanceRepo *repository.AllowanceRepo
}

func NewAllowanceEngine(ar *repository.AllowanceRepo) *AllowanceEngine {
	return &AllowanceEngine{allowanceRepo: ar}
}

func (e *AllowanceEngine) CalculateDaily(ctx context.Context, companyID uuid.UUID, destination string, employeeCategory string) (*decimal.Decimal, error) {
	country := &destination
	rules, err := e.allowanceRepo.List(ctx, companyID, country, nil)
	if err != nil {
		return nil, engErr("CalculateDaily.list", err)
	}

	rule := findMatchingRule(rules, destination, employeeCategory)
	if rule == nil {
		return nil, nil
	}
	return &rule.DailyAmount, nil
}

func (e *AllowanceEngine) CalculateTrip(ctx context.Context, companyID uuid.UUID, destination string, employeeCategory string, days int, mealsProvided []string) (*decimal.Decimal, error) {
	country := &destination
	rules, err := e.allowanceRepo.List(ctx, companyID, country, nil)
	if err != nil {
		return nil, engErr("CalculateTrip.list", err)
	}

	rule := findMatchingRule(rules, destination, employeeCategory)
	if rule == nil {
		return nil, nil
	}

	total := rule.DailyAmount.Mul(decimal.NewFromInt(int64(days)))

	if rule.MealPercentage != nil && len(mealsProvided) > 0 {
		mealDeduction := rule.DailyAmount.Mul(*rule.MealPercentage).Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(int64(len(mealsProvided))))
		total = total.Sub(mealDeduction)
	}

	return &total, nil
}

func findMatchingRule(rules []domain.DailyAllowanceRule, destination, employeeCategory string) *domain.DailyAllowanceRule {
	for i := range rules {
		rule := &rules[i]
		if rule.City != nil && *rule.City == destination {
			return rule
		}
		if rule.Country != nil && *rule.Country == destination {
			return rule
		}
	}
	if len(rules) > 0 {
		return &rules[0]
	}
	return nil
}
