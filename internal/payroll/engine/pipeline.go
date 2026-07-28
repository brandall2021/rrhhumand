package engine

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/shopspring/decimal"
)

type Pipeline struct {
	stages []Stage
}

type Stage struct {
	Name    string
	Execute func(ctx context.Context, state *PipelineState) error
}

type PipelineState struct {
	Input  RuleInput
	Output *RuleOutput
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stages: []Stage{
			{Name: "LOAD", Execute: loadStage},
			{Name: "SNAPSHOT", Execute: snapshotStage},
			{Name: "EARNINGS", Execute: earningsStage},
			{Name: "BASES", Execute: basesStage},
			{Name: "DEDUCTIONS", Execute: deductionsStage},
			{Name: "CONTRIBUTIONS", Execute: contributionsStage},
			{Name: "NET", Execute: netStage},
			{Name: "COST", Execute: costStage},
			{Name: "VALIDATE", Execute: validateStage},
		},
	}
}

func (p *Pipeline) Execute(ctx context.Context, input RuleInput) (*RuleOutput, error) {
	state := &PipelineState{
		Input: input,
		Output: &RuleOutput{
			Warnings: []string{},
			Errors:   []string{},
		},
	}

	for _, stage := range p.stages {
		select {
		case <-ctx.Done():
			return state.Output, ctx.Err()
		default:
		}

		if err := stage.Execute(ctx, state); err != nil {
			state.Output.Errors = append(state.Output.Errors,
				fmt.Sprintf("stage %s: %v", stage.Name, err))
			return state.Output, err
		}
	}

	return state.Output, nil
}

func loadStage(_ context.Context, state *PipelineState) error {
	input := state.Input

	if input.Period.ID == "" {
		return fmt.Errorf("period is required")
	}
	if input.Run.ID == "" {
		return fmt.Errorf("run is required")
	}
	if input.Employee.ID == "" {
		return fmt.Errorf("employee snapshot is required")
	}
	if input.BaseSalary.LessThan(decimal.Zero) {
		return fmt.Errorf("base salary cannot be negative")
	}

	if state.Output.Warnings == nil {
		state.Output.Warnings = []string{}
	}
	if state.Output.Errors == nil {
		state.Output.Errors = []string{}
	}

	return nil
}

func snapshotStage(_ context.Context, state *PipelineState) error {
	emp := state.Input.Employee
	if emp.SalaryData == nil {
		emp.SalaryData = make(map[string]any)
	}
	emp.SalaryData["base_salary"] = state.Input.BaseSalary.String()
	emp.SalaryData["run_employee_id"] = emp.RunEmployeeID
	emp.EmployeeData["employee_id"] = emp.ID
	return nil
}

func earningsStage(ctx context.Context, state *PipelineState) error {
	conceptsByType := groupConceptsByType(state.Input.Concepts)

	earnings := conceptsByType["EARNING"]
	sort.Slice(earnings, func(i, j int) bool {
		return earnings[i].SortOrder < earnings[j].SortOrder
	})

	ruleByConcept := buildRuleMap(state.Input.Rules)

	var grossRem, grossNonRem decimal.Decimal

	for _, concept := range earnings {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rules, ok := ruleByConcept[concept.ID]
		if !ok {
			continue
		}

		rules = filterActiveRules(rules, state.Input.Period.StartDate)
		if len(rules) == 0 {
			continue
		}

		limits := filterLimitsForConcept(state.Input.Limits, concept.ID)

		item, err := EvaluateConcept(ctx, concept, rules, evaluationContext(state))
		if err != nil {
			state.Output.Warnings = append(state.Output.Warnings,
				fmt.Sprintf("concept %s (%s): %v", concept.Code, concept.Name, err))
			continue
		}
		if item == nil {
			continue
		}

		item.IsDeduction = false
		item.IsEmployerContribution = false
		item.IsRemunerative = concept.Taxability == "REMUNERATIVE"
		item.SortOrder = concept.SortOrder

		item.Amount = applyLimits(item.Amount, limits)

		state.Output.Items = append(state.Output.Items, *item)

		if concept.Taxability == "REMUNERATIVE" {
			grossRem = grossRem.Add(item.Amount)
		} else {
			grossNonRem = grossNonRem.Add(item.Amount)
		}
	}

	state.Output.GrossRemunerative = grossRem
	state.Output.GrossNonRemunerative = grossNonRem

	return nil
}

func basesStage(_ context.Context, state *PipelineState) error {
	if state.Output.GrossRemunerative.IsZero() && state.Output.GrossNonRemunerative.IsZero() {
		return nil
	}

	remunConceptIDs := []string{}
	nonRemunConceptIDs := []string{}

	for _, item := range state.Output.Items {
		if item.IsRemunerative {
			remunConceptIDs = append(remunConceptIDs, item.ConceptID)
		} else {
			nonRemunConceptIDs = append(nonRemunConceptIDs, item.ConceptID)
		}
	}

	if len(remunConceptIDs) > 0 {
		state.Output.Bases = append(state.Output.Bases, domain.PayrollBase{
			RunEmployeeID: state.Input.Employee.RunEmployeeID,
			BaseType:      "REMUNERATIVE",
			BaseAmount:    state.Output.GrossRemunerative,
			ConceptIDs:    remunConceptIDs,
		})
	}

	if len(nonRemunConceptIDs) > 0 {
		state.Output.Bases = append(state.Output.Bases, domain.PayrollBase{
			RunEmployeeID: state.Input.Employee.RunEmployeeID,
			BaseType:      "NON_REMUNERATIVE",
			BaseAmount:    state.Output.GrossNonRemunerative,
			ConceptIDs:    nonRemunConceptIDs,
		})
	}

	return nil
}

func deductionsStage(ctx context.Context, state *PipelineState) error {
	conceptsByType := groupConceptsByType(state.Input.Concepts)

	deductions := conceptsByType["DEDUCTION"]
	sort.Slice(deductions, func(i, j int) bool {
		return deductions[i].SortOrder < deductions[j].SortOrder
	})

	ruleByConcept := buildRuleMap(state.Input.Rules)

	var totalDeductions decimal.Decimal

	for _, concept := range deductions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rules, ok := ruleByConcept[concept.ID]
		if !ok {
			continue
		}

		rules = filterActiveRules(rules, state.Input.Period.StartDate)
		if len(rules) == 0 {
			continue
		}

		limits := filterLimitsForConcept(state.Input.Limits, concept.ID)

		item, err := EvaluateConcept(ctx, concept, rules, evaluationContext(state))
		if err != nil {
			state.Output.Warnings = append(state.Output.Warnings,
				fmt.Sprintf("deduction concept %s (%s): %v", concept.Code, concept.Name, err))
			continue
		}
		if item == nil {
			continue
		}

		item.IsDeduction = true
		item.IsEmployerContribution = false
		item.IsRemunerative = false
		item.SortOrder = concept.SortOrder + 1000

		item.Amount = applyLimits(item.Amount, limits)

		state.Output.Items = append(state.Output.Items, *item)

		deduction := domain.PayrollDeduction{
			RunEmployeeID: state.Input.Employee.RunEmployeeID,
			ConceptID:     concept.ID,
			BaseAmount:    item.BaseAmount,
			Rate:          item.Rate,
			Amount:        item.Amount,
		}
		state.Output.Deductions = append(state.Output.Deductions, deduction)
		totalDeductions = totalDeductions.Add(item.Amount)
	}

	state.Output.TotalDeductions = totalDeductions
	return nil
}

func contributionsStage(ctx context.Context, state *PipelineState) error {
	conceptsByType := groupConceptsByType(state.Input.Concepts)

	contributions := conceptsByType["CONTRIBUTION"]
	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].SortOrder < contributions[j].SortOrder
	})

	ruleByConcept := buildRuleMap(state.Input.Rules)

	var totalContributions decimal.Decimal

	for _, concept := range contributions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rules, ok := ruleByConcept[concept.ID]
		if !ok {
			continue
		}

		rules = filterActiveRules(rules, state.Input.Period.StartDate)
		if len(rules) == 0 {
			continue
		}

		item, err := EvaluateConcept(ctx, concept, rules, evaluationContext(state))
		if err != nil {
			state.Output.Warnings = append(state.Output.Warnings,
				fmt.Sprintf("contribution concept %s (%s): %v", concept.Code, concept.Name, err))
			continue
		}
		if item == nil {
			continue
		}

		item.IsDeduction = false
		item.IsEmployerContribution = true
		item.IsRemunerative = false
		item.SortOrder = concept.SortOrder + 2000

		state.Output.Items = append(state.Output.Items, *item)

		contribution := domain.PayrollContribution{
			RunEmployeeID: state.Input.Employee.RunEmployeeID,
			ConceptID:     concept.ID,
			BaseAmount:    item.BaseAmount,
			Rate:          item.Rate,
			Amount:        item.Amount,
		}
		state.Output.Contributions = append(state.Output.Contributions, contribution)
		totalContributions = totalContributions.Add(item.Amount)
	}

	state.Output.TotalContributions = totalContributions
	return nil
}

func netStage(_ context.Context, state *PipelineState) error {
	totalGross := state.Output.GrossRemunerative.Add(state.Output.GrossNonRemunerative)
	state.Output.Net = totalGross.Sub(state.Output.TotalDeductions)
	if state.Output.Net.LessThan(decimal.Zero) {
		state.Output.Net = decimal.Zero
		state.Output.Warnings = append(state.Output.Warnings,
			"net salary would be negative, floored to zero")
	}
	return nil
}

func costStage(_ context.Context, state *PipelineState) error {
	totalGross := state.Output.GrossRemunerative.Add(state.Output.GrossNonRemunerative)
	state.Output.EmployerCost = totalGross.Add(state.Output.TotalContributions)
	return nil
}

func validateStage(_ context.Context, state *PipelineState) error {
	if state.Output.GrossRemunerative.LessThan(decimal.Zero) {
		state.Output.Warnings = append(state.Output.Warnings,
			"gross remunerative is negative")
	}
	if state.Output.TotalDeductions.GreaterThan(state.Output.GrossRemunerative.Add(state.Output.GrossNonRemunerative)) {
		state.Output.Warnings = append(state.Output.Warnings,
			"total deductions exceed gross salary")
	}
	return nil
}

func groupConceptsByType(concepts []domain.PayrollConcept) map[string][]domain.PayrollConcept {
	result := make(map[string][]domain.PayrollConcept)
	for _, c := range concepts {
		if !c.Active {
			continue
		}
		if c.EffectiveTo != nil && !c.EffectiveTo.IsZero() && time.Now().After(*c.EffectiveTo) {
			continue
		}
		result[c.ConceptType] = append(result[c.ConceptType], c)
	}
	return result
}

func buildRuleMap(rules []domain.PayrollRule) map[string][]domain.PayrollRule {
	result := make(map[string][]domain.PayrollRule)
	for _, r := range rules {
		if !r.Active {
			continue
		}
		result[r.ConceptID] = append(result[r.ConceptID], r)
	}
	return result
}

func filterActiveRules(rules []domain.PayrollRule, date time.Time) []domain.PayrollRule {
	var active []domain.PayrollRule
	for _, r := range rules {
		if !r.Active {
			continue
		}
		if date.Before(r.EffectiveFrom) {
			continue
		}
		if r.EffectiveTo != nil && !r.EffectiveTo.IsZero() && date.After(*r.EffectiveTo) {
			continue
		}
		active = append(active, r)
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].Priority < active[j].Priority
	})
	return active
}

func filterLimitsForConcept(limits []domain.PayrollLimit, conceptID string) []domain.PayrollLimit {
	var result []domain.PayrollLimit
	for _, l := range limits {
		if l.ConceptID == nil || *l.ConceptID == conceptID {
			result = append(result, l)
		}
	}
	return result
}

func applyLimits(amount decimal.Decimal, limits []domain.PayrollLimit) decimal.Decimal {
	result := amount
	for _, l := range limits {
		if l.MinimumAmount != nil && result.LessThan(*l.MinimumAmount) {
			result = *l.MinimumAmount
		}
		if l.MaximumAmount != nil && result.GreaterThan(*l.MaximumAmount) {
			result = *l.MaximumAmount
		}
	}
	return result
}

func evaluationContext(state *PipelineState) map[string]any {
	return map[string]any{
		"employee":            state.Input.Employee,
		"period":              state.Input.Period,
		"run":                 state.Input.Run,
		"base_salary":         state.Input.BaseSalary,
		"min_wage":            state.Input.MinWage,
		"gross_remunerative":  state.Output.GrossRemunerative,
		"gross_non_remun":     state.Output.GrossNonRemunerative,
		"total_deductions":    state.Output.TotalDeductions,
		"total_contributions": state.Output.TotalContributions,
		"net":                 state.Output.Net,
	}
}
