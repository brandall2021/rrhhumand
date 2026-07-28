package payroll

import (
	"context"

	"github.com/google/uuid"
)

type Calculator struct {
	repo *Repository
}

func NewCalculator(repo *Repository) *Calculator {
	return &Calculator{repo: repo}
}

type CalculationResult struct {
	Items    []PayrollItem
	Total    float64
	Errors   []PayrollError
	Warnings []PayrollWarning
}

func (c *Calculator) CalculateForEmployee(ctx context.Context, companyID, employeeID string, period *PayrollPeriod) (*CalculationResult, error) {
	result := &CalculationResult{}

	compensation, err := c.repo.GetCompensation(ctx, companyID, employeeID)
	if err != nil {
		result.Errors = append(result.Errors, PayrollError{
			EmployeeID: employeeID,
			Message:    "No compensation configured",
		})
		return result, nil
	}

	baseConcept, err := c.repo.GetConceptByCode(ctx, companyID, "BASE")
	if err != nil {
		baseConcept = &PayrollConcept{ID: uuid.New().String(), Code: "BASE", Name: "Base Salary", Type: "EARNING"}
	}

	baseItem := PayrollItem{
		ID:              uuid.New().String(),
		PayrollPeriodID: period.ID,
		EmployeeID:      employeeID,
		ConceptID:       baseConcept.ID,
		Quantity:        1,
		UnitAmount:      compensation.BaseAmount,
		Amount:          compensation.BaseAmount,
	}
	result.Items = append(result.Items, baseItem)

	overtimeMinutes, _ := c.repo.GetApprovedOvertimeForPeriod(ctx, companyID, employeeID, period.StartDate, period.EndDate)
	if overtimeMinutes > 0 {
		overtimeHours := overtimeMinutes / 60.0
		rate, _ := c.repo.GetOvertimeRate(ctx, companyID, employeeID)
		hourlyRate := compensation.BaseAmount / 176.0
		overtimeAmount := overtimeHours * hourlyRate * rate

		overtimeConcept, _ := c.repo.GetConceptByCode(ctx, companyID, "OVERTIME")
		if overtimeConcept == nil {
			overtimeConcept = &PayrollConcept{ID: uuid.New().String(), Code: "OVERTIME", Name: "Overtime", Type: "EARNING"}
		}

		overtimeItem := PayrollItem{
			ID:              uuid.New().String(),
			PayrollPeriodID: period.ID,
			EmployeeID:      employeeID,
			ConceptID:       overtimeConcept.ID,
			Quantity:        overtimeHours,
			UnitAmount:      hourlyRate * rate,
			Amount:          overtimeAmount,
		}
		result.Items = append(result.Items, overtimeItem)
	}

	bonuses, _ := c.repo.GetBonusesForPeriod(ctx, companyID, employeeID, period.StartDate, period.EndDate)
	bonusConcept, _ := c.repo.GetConceptByCode(ctx, companyID, "BONUS")
	if bonusConcept == nil {
		bonusConcept = &PayrollConcept{ID: uuid.New().String(), Code: "BONUS", Name: "Bonus", Type: "EARNING"}
	}
	for _, bonus := range bonuses {
		item := PayrollItem{
			ID:              uuid.New().String(),
			PayrollPeriodID: period.ID,
			EmployeeID:      employeeID,
			ConceptID:       bonusConcept.ID,
			Quantity:        1,
			UnitAmount:      bonus.Amount,
			Amount:          bonus.Amount,
		}
		result.Items = append(result.Items, item)
	}

	benefits, _ := c.repo.GetEmployeeBenefits(ctx, companyID, employeeID, period.StartDate)
	benefitConcept, _ := c.repo.GetConceptByCode(ctx, companyID, "BENEFIT")
	if benefitConcept == nil {
		benefitConcept = &PayrollConcept{ID: uuid.New().String(), Code: "BENEFIT", Name: "Benefit", Type: "BENEFIT"}
	}
	for _, ben := range benefits {
		amount := 0.0
		if ben.Amount != nil {
			amount = *ben.Amount
		}
		if amount > 0 {
			item := PayrollItem{
				ID:              uuid.New().String(),
				PayrollPeriodID: period.ID,
				EmployeeID:      employeeID,
				ConceptID:       benefitConcept.ID,
				Quantity:        1,
				UnitAmount:      amount,
				Amount:          amount,
			}
			result.Items = append(result.Items, item)
		}
	}

	deductions, _ := c.repo.GetDeductionsForPeriod(ctx, companyID, employeeID, period.StartDate, period.EndDate)
	deductionConcept, _ := c.repo.GetConceptByCode(ctx, companyID, "DEDUCTION")
	if deductionConcept == nil {
		deductionConcept = &PayrollConcept{ID: uuid.New().String(), Code: "DEDUCTION", Name: "Deduction", Type: "DEDUCTION"}
	}
	for _, ded := range deductions {
		item := PayrollItem{
			ID:              uuid.New().String(),
			PayrollPeriodID: period.ID,
			EmployeeID:      employeeID,
			ConceptID:       deductionConcept.ID,
			Quantity:        1,
			UnitAmount:      -ded.Amount,
			Amount:          -ded.Amount,
		}
		result.Items = append(result.Items, item)
	}

	advances, _ := c.repo.GetAdvancesForPeriod(ctx, companyID, employeeID, period.StartDate, period.EndDate)
	advanceConcept, _ := c.repo.GetConceptByCode(ctx, companyID, "ADVANCE")
	if advanceConcept == nil {
		advanceConcept = &PayrollConcept{ID: uuid.New().String(), Code: "ADVANCE", Name: "Advance", Type: "DEDUCTION"}
	}
	for _, adv := range advances {
		item := PayrollItem{
			ID:              uuid.New().String(),
			PayrollPeriodID: period.ID,
			EmployeeID:      employeeID,
			ConceptID:       advanceConcept.ID,
			Quantity:        1,
			UnitAmount:      -adv.Amount,
			Amount:          -adv.Amount,
		}
		result.Items = append(result.Items, item)
	}

	for _, item := range result.Items {
		result.Total += item.Amount
	}

	return result, nil
}

func (c *Calculator) CalculateForAllEmployees(ctx context.Context, companyID string, period *PayrollPeriod) (*PayrollReview, error) {
	employees, err := c.repo.GetActiveEmployees(ctx, companyID)
	if err != nil {
		return nil, err
	}

	review := &PayrollReview{
		TotalEmployees: len(employees),
	}

	for _, empID := range employees {
		result, err := c.CalculateForEmployee(ctx, companyID, empID, period)
		if err != nil {
			review.ErrorDetails = append(review.ErrorDetails, PayrollError{EmployeeID: empID, Message: err.Error()})
			continue
		}

		for _, item := range result.Items {
			c.repo.CreateItem(ctx, &item)
		}

		for _, e := range result.Errors {
			review.ErrorDetails = append(review.ErrorDetails, e)
		}
		for _, w := range result.Warnings {
			review.WarningDetails = append(review.WarningDetails, w)
		}
	}

	review.Errors = len(review.ErrorDetails)
	review.Warnings = len(review.WarningDetails)
	review.Calculated = review.TotalEmployees - review.Errors

	return review, nil
}

func (c *Calculator) GetPeriodTotalGross(ctx context.Context, periodID string) (float64, error) {
	items, err := c.repo.ListItemsByPeriod(ctx, periodID)
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, item := range items {
		if item.Amount > 0 {
			total += item.Amount
		}
	}
	return total, nil
}

func (c *Calculator) GetPeriodTotalNet(ctx context.Context, periodID string) (float64, error) {
	items, err := c.repo.ListItemsByPeriod(ctx, periodID)
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, item := range items {
		total += item.Amount
	}
	return total, nil
}
