package engine

import (
	"context"
	"fmt"
	"sort"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/shopspring/decimal"
)

type CalculationStrategy interface {
	Supports(concept domain.PayrollConcept) bool
	Calculate(ctx context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, novelties []domain.PayrollNovelty, baseSalary decimal.Decimal, period domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error)
}

type BasicSalaryStrategy struct{}

func (s *BasicSalaryStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.ConceptType == "EARNING" && concept.CalculationType == "AMOUNT"
}

func (s *BasicSalaryStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, _ []domain.PayrollNovelty, baseSalary decimal.Decimal, _ domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	amount := baseSalary
	if val, ok := GetParameter(rule.Parameters, "amount"); ok {
		amount = val
	}

	amount = ApplyCap(amount, limits)

	item := &domain.PayrollItem{
		ConceptID:     concept.ID,
		BaseAmount:    baseSalary,
		Amount:        amount,
		Quantity:      decimal.NewFromInt(1),
		UnitValue:     amount,
		IsRemunerative: concept.Taxability == "REMUNERATIVE",
		SortOrder:     concept.SortOrder,
	}
	return item, nil
}

type PercentageStrategy struct{}

func (s *PercentageStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.CalculationType == "PERCENTAGE"
}

func (s *PercentageStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, _ []domain.PayrollNovelty, baseSalary decimal.Decimal, _ domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	base := baseSalary
	if val, ok := GetParameter(rule.Parameters, "base_amount"); ok {
		base = val
	}

	rate := decimal.NewFromInt(100)
	if val, ok := GetParameter(rule.Parameters, "rate"); ok {
		rate = val
	}
	if val, ok := GetParameter(rule.Parameters, "percentage"); ok {
		rate = val
	}

	result := ApplyPercentage(base, rate, limits)
	r := rate

	item := &domain.PayrollItem{
		ConceptID:  concept.ID,
		BaseAmount: base,
		Rate:       &r,
		Amount:     result,
		Quantity:   decimal.NewFromInt(1),
		UnitValue:  base,
	}
	return item, nil
}

type HourlyStrategy struct{}

func (s *HourlyStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.ConceptType == "EARNING" && concept.CalculationType == "HOURLY"
}

func (s *HourlyStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, novelties []domain.PayrollNovelty, baseSalary decimal.Decimal, period domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	hours := decimal.Zero
	if val, ok := GetParameter(rule.Parameters, "hours"); ok {
		hours = val
	}

	for _, nov := range novelties {
		if nov.NoveltyType == "HOURS" && nov.Quantity != nil {
			hours = hours.Add(*nov.Quantity)
		}
	}

	workDays := 30
	if period.EndDate.Month() == period.StartDate.Month() {
		workDays = CalcWorkDaysInPeriod(period.StartDate, period.EndDate)
	}

	hourlyValue := CalcHourlyValue(baseSalary, workDays)
	if val, ok := GetParameter(rule.Parameters, "hourly_value"); ok {
		hourlyValue = val
	}

	amount := hours.Mul(hourlyValue)
	amount = ApplyCap(amount, limits)

	item := &domain.PayrollItem{
		ConceptID:  concept.ID,
		BaseAmount: baseSalary,
		Amount:     amount,
		Quantity:   hours,
		UnitValue:  hourlyValue,
	}
	return item, nil
}

type DailyStrategy struct{}

func (s *DailyStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.ConceptType == "EARNING" && concept.CalculationType == "DAILY"
}

func (s *DailyStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, novelties []domain.PayrollNovelty, baseSalary decimal.Decimal, period domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	days := decimal.Zero
	if val, ok := GetParameter(rule.Parameters, "days"); ok {
		days = val
	}

	for _, nov := range novelties {
		if nov.NoveltyType == "DAYS" && nov.Quantity != nil {
			days = days.Add(*nov.Quantity)
		}
	}

	workDays := 30
	if period.EndDate.Month() == period.StartDate.Month() {
		workDays = CalcWorkDaysInPeriod(period.StartDate, period.EndDate)
	}

	dailyValue := CalcDailyValue(baseSalary, workDays)
	if val, ok := GetParameter(rule.Parameters, "daily_value"); ok {
		dailyValue = val
	}

	amount := days.Mul(dailyValue)
	amount = ApplyCap(amount, limits)

	item := &domain.PayrollItem{
		ConceptID:  concept.ID,
		BaseAmount: baseSalary,
		Amount:     amount,
		Quantity:   days,
		UnitValue:  dailyValue,
	}
	return item, nil
}

type UnitStrategy struct{}

func (s *UnitStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.ConceptType == "EARNING" && concept.CalculationType == "UNIT"
}

func (s *UnitStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, novelties []domain.PayrollNovelty, baseSalary decimal.Decimal, _ domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	quantity := decimal.NewFromInt(1)
	if val, ok := GetParameter(rule.Parameters, "quantity"); ok {
		quantity = val
	}

	for _, nov := range novelties {
		if nov.NoveltyType == "UNITS" && nov.Quantity != nil {
			quantity = quantity.Add(*nov.Quantity)
		}
		if nov.NoveltyType == "UNIT" && nov.Quantity != nil {
			quantity = quantity.Add(*nov.Quantity)
		}
	}

	unitValue := decimal.Zero
	if val, ok := GetParameter(rule.Parameters, "unit_value"); ok {
		unitValue = val
	}
	for _, nov := range novelties {
		if nov.NoveltyType == "UNIT_VALUE" && nov.UnitValue != nil {
			unitValue = *nov.UnitValue
		}
	}

	amount := quantity.Mul(unitValue)
	amount = ApplyCap(amount, limits)

	item := &domain.PayrollItem{
		ConceptID:  concept.ID,
		BaseAmount: baseSalary,
		Amount:     amount,
		Quantity:   quantity,
		UnitValue:  unitValue,
	}
	return item, nil
}

type FixedDeductionStrategy struct{}

func (s *FixedDeductionStrategy) Supports(concept domain.PayrollConcept) bool {
	return concept.ConceptType == "DEDUCTION" && concept.CalculationType == "AMOUNT"
}

func (s *FixedDeductionStrategy) Calculate(_ context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, _ []domain.PayrollNovelty, baseSalary decimal.Decimal, _ domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule := rules[0]

	amount := decimal.Zero
	if val, ok := GetParameter(rule.Parameters, "amount"); ok {
		amount = val
	}
	if val, ok := GetParameter(rule.Parameters, "value"); ok {
		amount = val
	}

	base := baseSalary
	if concept.BaseConceptID != nil {
		for _, r := range rules {
			if val, ok := GetParameter(r.Parameters, "base_amount"); ok {
				base = val
				break
			}
		}
	}

	amount = ApplyCap(amount, limits)

	item := &domain.PayrollItem{
		ConceptID:   concept.ID,
		BaseAmount:  base,
		Amount:      amount,
		Quantity:    decimal.NewFromInt(1),
		UnitValue:   amount,
		IsDeduction: true,
	}
	return item, nil
}

func GetStrategies() []CalculationStrategy {
	return []CalculationStrategy{
		&BasicSalaryStrategy{},
		&PercentageStrategy{},
		&HourlyStrategy{},
		&DailyStrategy{},
		&UnitStrategy{},
		&FixedDeductionStrategy{},
	}
}

func EvaluateUsingStrategy(ctx context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, novelties []domain.PayrollNovelty, baseSalary decimal.Decimal, period domain.PayrollPeriod, limits []domain.PayrollLimit) (*domain.PayrollItem, error) {
	strategies := GetStrategies()

	for _, strategy := range strategies {
		if strategy.Supports(concept) {
			return strategy.Calculate(ctx, concept, rules, novelties, baseSalary, period, limits)
		}
	}

	return nil, fmt.Errorf("no strategy supports concept %s with type %s calc %s",
		concept.Code, concept.ConceptType, concept.CalculationType)
}

func filterRulesForPeriod(rules []domain.PayrollRule, period domain.PayrollPeriod) []domain.PayrollRule {
	var active []domain.PayrollRule
	for _, r := range rules {
		if !r.Active {
			continue
		}
		if period.StartDate.Before(r.EffectiveFrom) {
			continue
		}
		if r.EffectiveTo != nil && !r.EffectiveTo.IsZero() && period.EndDate.After(*r.EffectiveTo) {
			continue
		}
		active = append(active, r)
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].Priority < active[j].Priority
	})
	return active
}

var _ = sort.Strings
