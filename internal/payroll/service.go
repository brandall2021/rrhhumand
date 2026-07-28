package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("payroll_svc.%s: %w", op, err)
}

type Service struct {
	repo *Repository
	log  *zap.Logger
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// ========================================================================
// PERIODS
// ========================================================================

func (s *Service) CreatePeriod(ctx context.Context, companyID, userID string, req CreatePeriodRequest) (*PayrollPeriod, error) {
	p := &PayrollPeriod{
		ID:          uuid.NewString(),
		CompanyID:   companyID,
		Year:        req.Year,
		Month:       req.Month,
		PeriodType:  req.PeriodType,
		Name:        req.Name,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		PaymentDate: req.PaymentDate,
		Status:      "OPEN",
		CreatedBy:   userID,
	}
	if err := s.repo.CreatePeriod(ctx, p); err != nil {
		return nil, svcErr("CreatePeriod", err)
	}
	return p, nil
}

func (s *Service) UpdatePeriod(ctx context.Context, companyID, id string, req UpdatePeriodRequest) (*PayrollPeriod, error) {
	p, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdatePeriod", err)
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.PaymentDate != nil {
		p.PaymentDate = req.PaymentDate
	}
	if err := s.repo.UpdatePeriod(ctx, p); err != nil {
		return nil, svcErr("UpdatePeriod", err)
	}
	return p, nil
}

func (s *Service) GetPeriod(ctx context.Context, companyID, id string) (*PayrollPeriod, error) {
	return s.repo.GetPeriod(ctx, companyID, id)
}

func (s *Service) ListPeriods(ctx context.Context, companyID string, limit, offset int) ([]PayrollPeriod, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListPeriods(ctx, companyID, limit, offset)
}

func (s *Service) ClosePeriod(ctx context.Context, companyID, id, userID string) error {
	p, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil {
		return svcErr("ClosePeriod", err)
	}
	if p.Status == "CLOSED" {
		return svcErr("ClosePeriod", fmt.Errorf("period already closed"))
	}
	return s.repo.ClosePeriod(ctx, id, userID)
}

// ========================================================================
// RUNS
// ========================================================================

func (s *Service) CreateRun(ctx context.Context, companyID, periodID, userID string, req CreateRunRequest) (*PayrollRun, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, svcErr("CreateRun", err)
	}
	if period.Status == "CLOSED" {
		return nil, svcErr("CreateRun", fmt.Errorf("cannot create run on closed period"))
	}
	runNumber, err := s.repo.GetRunNumber(ctx, periodID)
	if err != nil {
		return nil, svcErr("CreateRun", err)
	}
	run := &PayrollRun{
		ID:        uuid.NewString(),
		CompanyID: companyID,
		PeriodID:  periodID,
		RunNumber: runNumber,
		RunType:   req.RunType,
		Status:    "OPEN",
		CreatedBy: userID,
	}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, svcErr("CreateRun", err)
	}
	return run, nil
}

func (s *Service) GetRun(ctx context.Context, companyID, id string) (*PayrollRun, error) {
	return s.repo.GetRun(ctx, companyID, id)
}

func (s *Service) ListRuns(ctx context.Context, companyID string, filter RunFilter) ([]PayrollRun, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.repo.ListRuns(ctx, companyID, filter)
}

func (s *Service) CalculateRun(ctx context.Context, companyID, runID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return svcErr("CalculateRun", err)
	}
	if run.Status != "OPEN" && run.Status != "CALCULATED" {
		return svcErr("CalculateRun", fmt.Errorf("run status %s does not allow calculation", run.Status))
	}

	s.repo.UpdateRunStatus(ctx, runID, "CALCULATING")
	now := time.Now()
	s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATING", &now, nil)

	employees, err := s.repo.ListRunEmployees(ctx, runID)
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "OPEN")
		return svcErr("CalculateRun", err)
	}

	period, err := s.repo.GetPeriod(ctx, companyID, run.PeriodID)
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "OPEN")
		return svcErr("CalculateRun", err)
	}

	concepts, err := s.repo.ListConcepts(ctx, companyID, ConceptFilter{Active: boolPtr(true)})
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "OPEN")
		return svcErr("CalculateRun", err)
	}

	if len(concepts) == 0 {
		s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATED", nil, timePtr(time.Now()))
		return nil
	}

	conceptIDs := make([]string, len(concepts))
	for i, c := range concepts {
		conceptIDs[i] = c.ID
	}

	rules, err := s.repo.GetActiveRulesByConceptIDs(ctx, companyID, conceptIDs, period.StartDate)
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "OPEN")
		return svcErr("CalculateRun", err)
	}

	limits, err := s.repo.GetActiveLimits(ctx, companyID, period.StartDate)
	if err != nil {
		s.log.Warn("no limits found", zap.Error(err))
	}
	minWage, _ := s.repo.GetMinimumWage(ctx, "AR", period.StartDate)

	ruleMap := mapRulesByConcept(rules)
	limitMap := mapLimitsByConcept(limits)

	for i := range employees {
		emp := &employees[i]
		emp.Status = "CALCULATING"
		s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)

		novelties, nErr := s.repo.GetNoveltiesForEmployeePeriod(ctx, companyID, emp.EmployeeID, run.PeriodID)
		if nErr != nil {
			emp.Status = "ERROR"
			errMsg := nErr.Error()
			emp.ErrorMessage = &errMsg
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			s.recordRunError(ctx, runID, &emp.EmployeeID, "ERROR", "NOVELTY_FETCH", nErr.Error(), nil)
			continue
		}

		result, calcErr := s.calculateEmployee(ctx, companyID, emp, period, concepts, ruleMap, novelties, limits, limitMap, minWage)
		if calcErr != nil {
			emp.Status = "ERROR"
			errMsg := calcErr.Error()
			emp.ErrorMessage = &errMsg
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			s.recordRunError(ctx, runID, &emp.EmployeeID, "ERROR", "CALCULATION", calcErr.Error(), nil)
			continue
		}

		emp.Status = "CALCULATED"
		emp.GrossRemunerative = result.GrossRemunerative
		emp.GrossNonRemunerative = result.GrossNonRemunerative
		emp.DeductionsAmount = result.EmployeeDeductions
		emp.EmployerContributions = result.EmployerContributions
		emp.EmployerCost = result.EmployerCost
		emp.NetAmount = result.Net
		emp.CalculationVersion++
		calcAt := time.Now()
		emp.CalculatedAt = &calcAt
		s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)

		if len(result.Items) > 0 {
			s.repo.DeleteItemsForRunEmployee(ctx, emp.ID)
			s.repo.BulkCreateItems(ctx, result.Items)
		}
		if len(result.Bases) > 0 {
			s.repo.BulkCreateBases(ctx, result.Bases)
		}
		for _, w := range result.Warnings {
			s.recordRunError(ctx, runID, &emp.EmployeeID, "WARNING", "CALC_WARN", w, nil)
		}
		for _, e := range result.Errors {
			s.recordRunError(ctx, runID, &emp.EmployeeID, "ERROR", "CALC_ERR", e, nil)
		}
	}

	finished := time.Now()
	s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATED", nil, &finished)
	return nil
}

func (s *Service) calculateEmployee(ctx context.Context, companyID string, emp *PayrollRunEmployee,
	period *PayrollPeriod, concepts []PayrollConcept, ruleMap map[string][]PayrollRule,
	novelties []PayrollNovelty, limits []PayrollLimit, limitMap map[string][]PayrollLimit,
	minWage *StatutoryMinimumWage) (*PayrollResult, error) {

	baseSalary, currency, _ := s.repo.GetEmployeeCompensation(ctx, companyID, emp.EmployeeID)
	emp.Currency = currency

	result := &PayrollResult{
		EmployeeID: emp.EmployeeID,
		Items:      []PayrollItem{},
		Bases:      []PayrollBase{},
		Warnings:   []string{},
		Errors:     []string{},
	}

	var remunerativeTotal decimal.Decimal
	var nonRemunerativeTotal decimal.Decimal
	var deductionTotal decimal.Decimal
	var contributionTotal decimal.Decimal

	itemOrder := 0

	noveltyMap := mapNoveltiesByType(novelties)

	for _, concept := range concepts {
		if !concept.Active {
			continue
		}
		if concept.EffectiveTo != nil && period.StartDate.After(*concept.EffectiveTo) {
			continue
		}
		if concept.EffectiveFrom.After(period.StartDate) {
			continue
		}

		itemOrder++
		item := PayrollItem{
			ID:            uuid.NewString(),
			RunEmployeeID: emp.ID,
			ConceptID:     concept.ID,
			Quantity:      decimal.NewFromInt(1),
			SortOrder:     itemOrder,
		}

		novelties := noveltyMap[concept.Code]

		switch concept.CalculationType {
		case "AMOUNT":
			amount, err := s.applyRules(ctx, concept.ID, ruleMap[concept.ID], nil, period, limitMap[concept.ID])
			if err == nil && !amount.IsZero() {
				item.Amount = amount
			} else if len(novelties) > 0 {
				for _, n := range novelties {
					if n.Amount != nil {
						item.Amount = item.Amount.Add(*n.Amount)
					}
				}
			}
			if !item.Amount.IsZero() && item.Amount.GreaterThan(decimal.Zero) {
				item.BaseAmount = item.Amount
			}

		case "PERCENTAGE":
			baseConcept := s.findConcept(concept.BaseConceptID, concepts)
			baseItems := result.Items
			var baseAmount decimal.Decimal
			if baseConcept != nil {
				for _, it := range baseItems {
					if it.ConceptID == baseConcept.ID {
						baseAmount = baseAmount.Add(it.Amount)
						break
					}
				}
			}
			if baseAmount.IsZero() {
				baseAmount = remunerativeTotal
			}

			rules := ruleMap[concept.ID]
			rate, ruleErr := s.getPercentageRate(rules, period)
			if ruleErr != nil {
				baseAmount = decimal.Zero
			}
			if !baseAmount.IsZero() {
				item.BaseAmount = baseAmount
				item.Rate = &rate
				item.Amount = s.applyCap(baseAmount.Mul(rate.Div(decimal.NewFromInt(100))), limitMap[concept.ID], &item)
				if !item.Amount.IsZero() {
					deductionTotal = deductionTotal.Add(item.Amount)
				}
			}

		case "HOURLY":
			hourlyValue := s.calcHourlyValue(baseSalary, period)
			if len(novelties) > 0 {
				var totalHours decimal.Decimal
				for _, n := range novelties {
					if n.Quantity != nil {
						totalHours = totalHours.Add(*n.Quantity)
					}
				}
				if !totalHours.IsZero() {
					item.Quantity = totalHours
					item.UnitValue = hourlyValue
					item.BaseAmount = hourlyValue
					rules := ruleMap[concept.ID]
					multiplier := decimal.NewFromInt(1)
					for _, rule := range rules {
						if rule.RuleType == "MULTIPLIER" {
							if pct, ok := getDecimalParam(rule.Parameters, "multiplier"); ok {
								multiplier = pct
							}
						}
					}
					item.Amount = totalHours.Mul(hourlyValue).Mul(multiplier)
				}
			}

		case "DAILY":
			dailyValue := s.calcDailyValue(baseSalary, period)
			if len(novelties) > 0 {
				var totalDays decimal.Decimal
				for _, n := range novelties {
					if n.Quantity != nil {
						totalDays = totalDays.Add(*n.Quantity)
					}
				}
				if totalDays.IsZero() && len(novelties) > 0 {
					for _, n := range novelties {
						if n.StartDate != nil && n.EndDate != nil {
							days := decimal.NewFromInt(int64(n.EndDate.Sub(*n.StartDate).Hours()/24) + 1)
							totalDays = totalDays.Add(days)
						}
					}
				}
				if !totalDays.IsZero() {
					item.Quantity = totalDays
					item.UnitValue = dailyValue
					item.BaseAmount = dailyValue
					item.Amount = totalDays.Mul(dailyValue)
				}
			}

		case "UNIT":
			if len(novelties) > 0 {
				var totalUnits decimal.Decimal
				var unitVal decimal.Decimal
				for _, n := range novelties {
					if n.Quantity != nil {
						totalUnits = totalUnits.Add(*n.Quantity)
					}
					if n.UnitValue != nil {
						unitVal = *n.UnitValue
					}
				}
				if !totalUnits.IsZero() {
					item.Quantity = totalUnits
					item.UnitValue = unitVal
					item.BaseAmount = unitVal
					item.Amount = totalUnits.Mul(unitVal)
				}
			}
		}

		if item.Amount.IsZero() {
			continue
		}

		item.IsRemunerative = concept.Taxability == "REMUNERATIVO" && concept.ConceptType == "EARNING"
		item.IsDeduction = concept.ConceptType == "DEDUCTION"
		item.IsEmployerContribution = concept.ConceptType == "EMPLOYER_CONTRIBUTION"

		switch concept.ConceptType {
		case "EARNING":
			if concept.Taxability == "REMUNERATIVO" {
				remunerativeTotal = remunerativeTotal.Add(item.Amount)
			} else {
				nonRemunerativeTotal = nonRemunerativeTotal.Add(item.Amount)
			}
		case "DEDUCTION":
			deductionTotal = deductionTotal.Add(item.Amount)
		case "EMPLOYER_CONTRIBUTION":
			contributionTotal = contributionTotal.Add(item.Amount)
		}

		result.Items = append(result.Items, item)
	}

	result.GrossRemunerative = remunerativeTotal
	result.GrossNonRemunerative = nonRemunerativeTotal
	result.EmployeeDeductions = deductionTotal
	result.EmployerContributions = contributionTotal
	result.Net = remunerativeTotal.Add(nonRemunerativeTotal).Sub(deductionTotal)
	result.EmployerCost = remunerativeTotal.Add(contributionTotal)

	result.Bases = s.buildBases(emp.ID, result.Items, remunerativeTotal)

	return result, nil
}

func (s *Service) applyRules(ctx context.Context, conceptID string, rules []PayrollRule, baseAmount *decimal.Decimal, period *PayrollPeriod, limits []PayrollLimit) (decimal.Decimal, error) {
	if len(rules) == 0 {
		return decimal.Zero, fmt.Errorf("no rules for concept %s", conceptID)
	}
	rule := rules[0]
	switch rule.RuleType {
	case "AMOUNT":
		if amt, ok := getDecimalParam(rule.Parameters, "amount"); ok {
			capAmount := s.applyCap(amt, limits, nil)
			return capAmount, nil
		}
	case "PERCENTAGE":
		if pct, ok := getDecimalParam(rule.Parameters, "rate"); ok {
			if baseAmount != nil && !baseAmount.IsZero() {
				amt := baseAmount.Mul(pct).Div(decimal.NewFromInt(100))
				return s.applyCap(amt, limits, nil), nil
			}
		}
	case "FORMULA":
		if rule.Formula != nil && *rule.Formula == "base_salary" {
			return decimal.Zero, fmt.Errorf("use base salary from compensation")
		}
	}
	return decimal.Zero, fmt.Errorf("no applicable rule for concept %s", conceptID)
}

func (s *Service) getPercentageRate(rules []PayrollRule, period *PayrollPeriod) (decimal.Decimal, error) {
	for _, rule := range rules {
		if rule.RuleType == "PERCENTAGE" || rule.RuleType == "AMOUNT" {
			if pct, ok := getDecimalParam(rule.Parameters, "rate"); ok {
				return pct, nil
			}
			if pct, ok := getDecimalParam(rule.Parameters, "percentage"); ok {
				return pct, nil
			}
		}
	}
	if len(rules) > 0 {
		return decimal.Zero, fmt.Errorf("no rate parameter on rules for concept")
	}
	return decimal.Zero, fmt.Errorf("no rules found")
}

func (s *Service) applyCap(amount decimal.Decimal, limits []PayrollLimit, item *PayrollItem) decimal.Decimal {
	for _, l := range limits {
		switch l.LimitType {
		case "MAXIMUM":
			if l.MaximumAmount != nil && amount.GreaterThan(*l.MaximumAmount) {
				if item != nil {
					detail := map[string]any{"capped": true, "original": amount, "max": *l.MaximumAmount}
					item.CalculationDetail, _ = json.Marshal(detail)
				}
				return *l.MaximumAmount
			}
		case "MINIMUM":
			if l.MinimumAmount != nil && amount.LessThan(*l.MinimumAmount) {
				if item != nil {
					detail := map[string]any{"adjusted": true, "original": amount, "min": *l.MinimumAmount}
					item.CalculationDetail, _ = json.Marshal(detail)
				}
				return *l.MinimumAmount
			}
		}
	}
	return amount
}

func (s *Service) calcHourlyValue(baseSalary decimal.Decimal, period *PayrollPeriod) decimal.Decimal {
	workDays := s.workDaysInPeriod(period)
	if workDays == 0 {
		return decimal.Zero
	}
	hoursPerDay := decimal.NewFromInt(8)
	return baseSalary.Div(decimal.NewFromInt(int64(workDays))).Div(hoursPerDay).Round(2)
}

func (s *Service) calcDailyValue(baseSalary decimal.Decimal, period *PayrollPeriod) decimal.Decimal {
	workDays := s.workDaysInPeriod(period)
	if workDays == 0 {
		return decimal.Zero
	}
	return baseSalary.Div(decimal.NewFromInt(int64(workDays))).Round(2)
}

func (s *Service) workDaysInPeriod(period *PayrollPeriod) int {
	days := 0
	current := period.StartDate
	for !current.After(period.EndDate) {
		if current.Weekday() != time.Sunday {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}
	if days == 0 {
		days = 30
	}
	return days
}

func (s *Service) findConcept(id *string, concepts []PayrollConcept) *PayrollConcept {
	if id == nil {
		return nil
	}
	for _, c := range concepts {
		if c.ID == *id {
			return &c
		}
	}
	return nil
}

func (s *Service) buildBases(runEmployeeID string, items []PayrollItem, remunerativeTotal decimal.Decimal) []PayrollBase {
	var bases []PayrollBase
	conceptIDs := make([]string, 0, len(items))
	for _, it := range items {
		conceptIDs = append(conceptIDs, it.ConceptID)
	}

	baseTypes := []string{"BASE_REMUNERATIVA", "BASE_JUBILACION", "BASE_OBRA_SOCIAL", "BASE_ART", "BASE_SINDICAL", "BASE_ASIGNACIONES"}
	for _, bt := range baseTypes {
		base := PayrollBase{
			ID:            uuid.NewString(),
			RunEmployeeID: runEmployeeID,
			BaseType:      bt,
			BaseAmount:    remunerativeTotal,
			ConceptIDs:    conceptIDs,
		}
		bases = append(bases, base)
	}
	return bases
}

// ========================================================================
// VALIDATE
// ========================================================================

func (s *Service) ValidateRun(ctx context.Context, companyID, runID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return svcErr("ValidateRun", err)
	}
	if run.Status != "CALCULATED" {
		return svcErr("ValidateRun", fmt.Errorf("run must be CALCULATED before validation, current: %s", run.Status))
	}

	s.repo.UpdateRunStatus(ctx, runID, "VALIDATING")

	employees, err := s.repo.ListRunEmployees(ctx, runID)
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "CALCULATED")
		return svcErr("ValidateRun", err)
	}

	hasBlocking := false
	for _, emp := range employees {
		if emp.Status == "ERROR" {
			s.recordRunError(ctx, runID, &emp.EmployeeID, "BLOCKING", "EMPLOYEE_ERROR",
				fmt.Sprintf("employee %s has calculation error", emp.EmployeeID), nil)
			hasBlocking = true
		}
		if emp.NetAmount.IsNegative() {
			s.recordRunError(ctx, runID, &emp.EmployeeID, "WARNING", "NEGATIVE_NET",
				fmt.Sprintf("employee %s has negative net amount", emp.EmployeeID), nil)
		}
	}

	summary, err := s.repo.GetRunSummary(ctx, runID)
	if err == nil && summary.TotalGross.IsNegative() {
		s.recordRunError(ctx, runID, nil, "BLOCKING", "NEGATIVE_GROSS", "total gross is negative", nil)
		hasBlocking = true
	}

	if hasBlocking {
		s.repo.UpdateRunStatus(ctx, runID, "CALCULATED")
		return svcErr("ValidateRun", fmt.Errorf("validation found blocking errors"))
	}

	s.repo.UpdateRunStatus(ctx, runID, "VALIDATED")
	return nil
}

// ========================================================================
// APPROVE
// ========================================================================

func (s *Service) ApproveRun(ctx context.Context, companyID, runID, userID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return svcErr("ApproveRun", err)
	}
	if run.Status != "VALIDATED" {
		return svcErr("ApproveRun", fmt.Errorf("run must be VALIDATED before approval, current: %s", run.Status))
	}
	blocking, err := s.repo.ListBlockingErrors(ctx, runID)
	if err != nil {
		return svcErr("ApproveRun", err)
	}
	if len(blocking) > 0 {
		return svcErr("ApproveRun", fmt.Errorf("cannot approve run with blocking errors"))
	}
	return s.repo.ApproveRun(ctx, runID, userID)
}

// ========================================================================
// CLOSE
// ========================================================================

func (s *Service) CloseRun(ctx context.Context, companyID, runID, userID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return svcErr("CloseRun", err)
	}
	if run.Status != "APPROVED" {
		return svcErr("CloseRun", fmt.Errorf("run must be APPROVED before closing, current: %s", run.Status))
	}
	if err := s.repo.CloseRun(ctx, runID, userID); err != nil {
		return svcErr("CloseRun", err)
	}
	// Auto-close period if no other open runs
	runs, _ := s.repo.ListRuns(ctx, companyID, RunFilter{PeriodID: &run.PeriodID, Status: strPtr("APPROVED")})
	if len(runs) == 0 {
		s.repo.UpdatePeriodStatus(ctx, run.PeriodID, "CLOSED")
	}
	return nil
}

// ========================================================================
// CONCEPTS
// ========================================================================

func (s *Service) CreateConcept(ctx context.Context, companyID, userID string, req CreateConceptRequest) (*PayrollConcept, error) {
	c := &PayrollConcept{
		ID:              uuid.NewString(),
		CompanyID:       companyID,
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		ConceptType:     req.ConceptType,
		Taxability:      req.Taxability,
		CalculationType: req.CalculationType,
		BaseConceptID:   req.BaseConceptID,
		Active:          true,
		EffectiveFrom:   time.Now().AddDate(0, -1, 0),
		SortOrder:       req.SortOrder,
		CreatedBy:       userID,
	}
	if err := s.repo.CreateConcept(ctx, c); err != nil {
		return nil, svcErr("CreateConcept", err)
	}
	return c, nil
}

func (s *Service) UpdateConcept(ctx context.Context, companyID, id string, req UpdateConceptRequest) (*PayrollConcept, error) {
	c, err := s.repo.GetConcept(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateConcept", err)
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.ConceptType != nil {
		c.ConceptType = *req.ConceptType
	}
	if req.Taxability != nil {
		c.Taxability = *req.Taxability
	}
	if req.CalculationType != nil {
		c.CalculationType = *req.CalculationType
	}
	if req.BaseConceptID != nil {
		c.BaseConceptID = req.BaseConceptID
	}
	if req.Active != nil {
		c.Active = *req.Active
	}
	if req.SortOrder != nil {
		c.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateConcept(ctx, c); err != nil {
		return nil, svcErr("UpdateConcept", err)
	}
	return c, nil
}

// ========================================================================
// RULES
// ========================================================================

func (s *Service) CreateRule(ctx context.Context, companyID, userID string, req CreateRuleRequest) (*PayrollRule, error) {
	effFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	var effTo *time.Time
	if req.EffectiveTo != nil {
		t, _ := time.Parse("2006-01-02", *req.EffectiveTo)
		effTo = &t
	}
	params, _ := json.Marshal(req.Parameters)
	rule := &PayrollRule{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		Country:       "AR",
		ConceptID:     req.ConceptID,
		RuleType:      req.RuleType,
		Formula:       req.Formula,
		Parameters:    params,
		Priority:      req.Priority,
		EffectiveFrom: effFrom,
		EffectiveTo:   effTo,
		Version:       1,
		Active:        true,
		CreatedBy:     userID,
	}
	if err := s.repo.CreateRule(ctx, rule); err != nil {
		return nil, svcErr("CreateRule", err)
	}
	return rule, nil
}

func (s *Service) UpdateRule(ctx context.Context, companyID, id string, req UpdateRuleRequest) (*PayrollRule, error) {
	rule, err := s.repo.GetRule(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateRule", err)
	}
	if req.RuleType != nil {
		rule.RuleType = *req.RuleType
	}
	if req.Formula != nil {
		rule.Formula = req.Formula
	}
	if req.Parameters != nil {
		params, _ := json.Marshal(req.Parameters)
		rule.Parameters = params
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Active != nil {
		rule.Active = *req.Active
	}
	if req.EffectiveTo != nil {
		t, _ := time.Parse("2006-01-02", *req.EffectiveTo)
		rule.EffectiveTo = &t
	}
	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		return nil, svcErr("UpdateRule", err)
	}
	return rule, nil
}

// ========================================================================
// NOVELTIES
// ========================================================================

func (s *Service) CreateNovelty(ctx context.Context, companyID, userID string, req CreateNoveltyRequest) (*PayrollNovelty, error) {
	qty := decimalPtrFromFloat(req.Quantity)
	amt := decimalPtrFromFloat(req.Amount)
	uv := decimalPtrFromFloat(req.UnitValue)
	mult := decimalPtrFromFloat(req.Multiplier)

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &t
	}

	n := &PayrollNovelty{
		ID:          uuid.NewString(),
		CompanyID:   companyID,
		EmployeeID:  req.EmployeeID,
		NoveltyType: req.NoveltyType,
		Quantity:    qty,
		Unit:        req.Unit,
		Amount:      amt,
		UnitValue:   uv,
		Multiplier:  mult,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
		Source:      req.Source,
		Status:      "PENDING",
		CreatedBy:   userID,
	}
	if err := s.repo.CreateNovelty(ctx, n); err != nil {
		return nil, svcErr("CreateNovelty", err)
	}
	return n, nil
}

func (s *Service) UpdateNovelty(ctx context.Context, companyID, id string, req UpdateNoveltyRequest) (*PayrollNovelty, error) {
	n, err := s.repo.GetNovelty(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateNovelty", err)
	}
	if req.Quantity != nil {
		n.Quantity = decimalPtrFromFloat(req.Quantity)
	}
	if req.Amount != nil {
		n.Amount = decimalPtrFromFloat(req.Amount)
	}
	if req.Description != nil {
		n.Description = req.Description
	}
	if req.Status != nil {
		n.Status = *req.Status
	}
	if err := s.repo.UpdateNovelty(ctx, n); err != nil {
		return nil, svcErr("UpdateNovelty", err)
	}
	return n, nil
}

func (s *Service) ImportNovelties(ctx context.Context, companyID, userID string, req ImportNoveltiesRequest) ([]PayrollNovelty, error) {
	result := make([]PayrollNovelty, 0, len(req.Novelties))
	for _, nr := range req.Novelties {
		n, err := s.CreateNovelty(ctx, companyID, userID, nr)
		if err != nil {
			return result, svcErr("ImportNovelties", err)
		}
		result = append(result, *n)
	}
	return result, nil
}

func (s *Service) ApproveNovelty(ctx context.Context, companyID, id, userID string) error {
	return s.repo.ApproveNovelty(ctx, id, userID)
}

// ========================================================================
// EMPLOYEES ON RUN
// ========================================================================

func (s *Service) AddEmployeeToRun(ctx context.Context, companyID, runID, employeeID string) (*PayrollRunEmployee, error) {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return nil, svcErr("AddEmployeeToRun", err)
	}
	if run.Status != "OPEN" {
		return nil, svcErr("AddEmployeeToRun", fmt.Errorf("run is not open"))
	}
	re := &PayrollRunEmployee{
		ID:         uuid.NewString(),
		RunID:      runID,
		EmployeeID: employeeID,
		Status:     "PENDING",
		Currency:   "ARS",
	}
	if err := s.repo.AddRunEmployee(ctx, re); err != nil {
		return nil, svcErr("AddEmployeeToRun", err)
	}
	return re, nil
}

func (s *Service) ListRunEmployees(ctx context.Context, companyID, runID string) ([]PayrollRunEmployee, error) {
	return s.repo.ListRunEmployees(ctx, runID)
}

func (s *Service) GetRunEmployee(ctx context.Context, companyID, runID, employeeID string) (*PayrollRunEmployee, error) {
	return s.repo.GetRunEmployee(ctx, runID, employeeID)
}

func (s *Service) GetRunSummary(ctx context.Context, companyID, runID string) (*PayrollSummary, error) {
	return s.repo.GetRunSummary(ctx, runID)
}

// ========================================================================
// MISC SERVICES
// ========================================================================

func (s *Service) CreateAdvance(ctx context.Context, companyID, userID string, req CreateAdvanceRequest) (*EmployeeAdvance, error) {
	reqDate, _ := time.Parse("2006-01-02", req.RequestDate)
	installmentAmt := decimal.NewFromFloat(req.Amount).Div(decimal.NewFromInt(int64(req.Installments)))
	a := &EmployeeAdvance{
		ID:                uuid.NewString(),
		CompanyID:         companyID,
		EmployeeID:        req.EmployeeID,
		Amount:            decimal.NewFromFloat(req.Amount),
		RequestDate:       reqDate,
		Installments:      req.Installments,
		InstallmentAmount: &installmentAmt,
		RemainingAmount:   decimal.NewFromFloat(req.Amount),
		Reason:            req.Reason,
		Status:            "PENDING",
		CreatedBy:         userID,
	}
	if err := s.repo.CreateAdvance(ctx, a); err != nil {
		return nil, svcErr("CreateAdvance", err)
	}
	return a, nil
}

func (s *Service) CreateGarnishment(ctx context.Context, companyID, userID string, req CreateGarnishmentRequest) (*PayrollGarnishment, error) {
	effFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	g := &PayrollGarnishment{
		ID:               uuid.NewString(),
		CompanyID:        companyID,
		EmployeeID:       req.EmployeeID,
		CourtOrderNumber: req.CourtOrderNumber,
		CourtName:        req.CourtName,
		Type:             req.Type,
		Percentage:       decimalPtrFromFloat(req.Percentage),
		FixedAmount:      decimalPtrFromFloat(req.FixedAmount),
		Priority:         req.Priority,
		EffectiveFrom:    effFrom,
		Status:           "ACTIVE",
		CreatedBy:        userID,
	}
	if err := s.repo.CreateGarnishment(ctx, g); err != nil {
		return nil, svcErr("CreateGarnishment", err)
	}
	return g, nil
}

func (s *Service) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, companyID)
}

func (s *Service) GetItems(ctx context.Context, runEmployeeID string) ([]PayrollItem, error) {
	return s.repo.ListItems(ctx, runEmployeeID)
}

func (s *Service) GetBases(ctx context.Context, runEmployeeID string) ([]PayrollBase, error) {
	return s.repo.ListBases(ctx, runEmployeeID)
}

func (s *Service) GetErrors(ctx context.Context, runID string) ([]PayrollError, error) {
	return s.repo.ListErrors(ctx, runID)
}

// ========================================================================
// HELPERS
// ========================================================================

func (s *Service) recordRunError(ctx context.Context, runID string, employeeID *string, severity, code, message string, field *string) {
	e := &PayrollError{
		ID:         uuid.NewString(),
		RunID:      runID,
		EmployeeID: employeeID,
		Severity:   severity,
		Code:       code,
		Message:    message,
		Field:      field,
	}
	if err := s.repo.CreateError(ctx, e); err != nil {
		s.log.Warn("failed to record payroll error", zap.Error(err))
	}
}

func mapRulesByConcept(rules []PayrollRule) map[string][]PayrollRule {
	m := make(map[string][]PayrollRule)
	for _, r := range rules {
		m[r.ConceptID] = append(m[r.ConceptID], r)
	}
	return m
}

func mapLimitsByConcept(limits []PayrollLimit) map[string][]PayrollLimit {
	m := make(map[string][]PayrollLimit)
	for _, l := range limits {
		if l.ConceptID != nil {
			m[*l.ConceptID] = append(m[*l.ConceptID], l)
		}
	}
	return m
}

func mapNoveltiesByType(novelties []PayrollNovelty) map[string][]PayrollNovelty {
	m := make(map[string][]PayrollNovelty)
	for _, n := range novelties {
		m[n.NoveltyType] = append(m[n.NoveltyType], n)
	}
	// Also map by concept code patterns
	for _, n := range novelties {
		m[n.NoveltyType] = append(m[n.NoveltyType], n)
	}
	return m
}

func getDecimalParam(params json.RawMessage, key string) (decimal.Decimal, bool) {
	var m map[string]any
	if err := json.Unmarshal(params, &m); err != nil {
		return decimal.Zero, false
	}
	v, ok := m[key]
	if !ok {
		return decimal.Zero, false
	}
	switch val := v.(type) {
	case float64:
		return decimal.NewFromFloat(val), true
	case string:
		d, err := decimal.NewFromString(val)
		return d, err == nil
	}
	return decimal.Zero, false
}

func decimalPtrFromFloat(f *float64) *decimal.Decimal {
	if f == nil {
		return nil
	}
	d := decimal.NewFromFloat(*f)
	return &d
}

func boolPtr(b bool) *bool { return &b }

func timePtr(t time.Time) *time.Time { return &t }

func strPtr(s string) *string { return &s }

// ========================================================================
// AGREEMENTS & CATEGORIES
// ========================================================================

func (s *Service) CreateAgreement(ctx context.Context, companyID, userID string, req CreateAgreementRequest) (*LaborAgreement, error) {
	effFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	a := &LaborAgreement{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		Activity:      req.Activity,
		EffectiveFrom: effFrom,
		Status:        "ACTIVE",
		CreatedBy:     userID,
	}
	if err := s.repo.CreateAgreement(ctx, a); err != nil {
		return nil, svcErr("CreateAgreement", err)
	}
	return a, nil
}

func (s *Service) ListAgreements(ctx context.Context, companyID string) ([]LaborAgreement, error) {
	return s.repo.ListAgreements(ctx, companyID)
}

func (s *Service) CreateCategory(ctx context.Context, companyID string, req CreateCategoryRequest) (*LaborCategory, error) {
	c := &LaborCategory{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		AgreementID:   req.AgreementID,
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		EffectiveFrom: time.Now().AddDate(0, -1, 0),
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, svcErr("CreateCategory", err)
	}
	return c, nil
}

func (s *Service) ListCategories(ctx context.Context, companyID string) ([]LaborCategory, error) {
	return s.repo.ListCategories(ctx, companyID)
}

func (s *Service) CreateSalaryScale(ctx context.Context, companyID string, req CreateSalaryScaleRequest) (*SalaryScale, error) {
	sc := &SalaryScale{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		AgreementID:   req.AgreementID,
		CategoryID:    req.CategoryID,
		MinimumSalary: decimal.NewFromFloat(req.MinimumSalary),
		EffectiveFrom: time.Now().AddDate(0, -1, 0),
	}
	if req.MaximumSalary != nil {
		m := decimal.NewFromFloat(*req.MaximumSalary)
		sc.MaximumSalary = &m
	}
	if err := s.repo.CreateSalaryScale(ctx, sc); err != nil {
		return nil, svcErr("CreateSalaryScale", err)
	}
	return sc, nil
}

func (s *Service) ListSalaryScales(ctx context.Context, companyID string) ([]SalaryScale, error) {
	return s.repo.ListSalaryScales(ctx, companyID)
}

// ========================================================================
// IDEMPOTENT CALCULATION
// ========================================================================

func (s *Service) GetOrCreateRunEmployee(ctx context.Context, companyID, runID, employeeID string) (*PayrollRunEmployee, error) {
	existing, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err == nil {
		return existing, nil
	}
	return s.AddEmployeeToRun(ctx, companyID, runID, employeeID)
}

// Employee result helpers

func (s *Service) GetEmployeeResult(ctx context.Context, companyID, runID, employeeID string) (*EmployeeResult, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, svcErr("GetEmployeeResult", err)
	}
	items, _ := s.repo.ListItems(ctx, re.ID)
	return &EmployeeResult{
		EmployeeID:          re.EmployeeID,
		Status:              re.Status,
		GrossRemunerative:   re.GrossRemunerative.InexactFloat64(),
		GrossNonRemunerative: re.GrossNonRemunerative.InexactFloat64(),
		DeductionsAmount:    re.DeductionsAmount.InexactFloat64(),
		NetAmount:           re.NetAmount.InexactFloat64(),
		EmployerContributions: re.EmployerContributions.InexactFloat64(),
		EmployerCost:        re.EmployerCost.InexactFloat64(),
		Error:               re.ErrorMessage,
	}
}

func (s *Service) GetEmployeeItems(ctx context.Context, companyID, runID, employeeID string) ([]PayrollItem, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, svcErr("GetEmployeeItems", err)
	}
	return s.repo.ListItems(ctx, re.ID)
}

// Employee self-service
func (s *Service) GetMyPeriods(ctx context.Context, companyID, employeeID string) ([]PayrollPeriod, error) {
	return s.repo.ListPeriods(ctx, companyID, 12, 0)
}

func (s *Service) GetMyReceiptData(ctx context.Context, companyID, employeeID, runID string) (*PayrollRunEmployee, error) {
	return s.repo.GetRunEmployee(ctx, runID, employeeID)
}

func (s *Service) GetMyItems(ctx context.Context, companyID, employeeID, runID string) ([]PayrollItem, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, svcErr("GetMyItems", err)
	}
	return s.repo.ListItems(ctx, re.ID)
}
