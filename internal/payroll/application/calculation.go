package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/rrhhumand/api/internal/payroll/engine"
)

func (s *PayrollService) AddEmployeeToRun(ctx context.Context, companyID, runID, employeeID string) (*domain.PayrollRunEmployee, error) {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return nil, fmt.Errorf("add employee: get run: %w", err)
	}
	if run.Status != "OPEN" {
		return nil, fmt.Errorf("add employee: run is not open")
	}
	re := &domain.PayrollRunEmployee{
		ID:         uuid.NewString(),
		RunID:      runID,
		EmployeeID: employeeID,
		Status:     "PENDING",
		Currency:   "ARS",
	}
	if err := s.repo.CreateRunEmployee(ctx, re); err != nil {
		return nil, fmt.Errorf("add employee: %w", err)
	}
	return re, nil
}

func (s *PayrollService) ListRunEmployees(ctx context.Context, companyID, runID string) ([]domain.PayrollRunEmployee, error) {
	employees, err := s.repo.ListRunEmployees(ctx, runID)
	if err != nil {
		return nil, fmt.Errorf("list run employees: %w", err)
	}
	return employees, nil
}

func (s *PayrollService) GetRunEmployee(ctx context.Context, companyID, runID, employeeID string) (*domain.PayrollRunEmployee, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("get run employee: %w", err)
	}
	return re, nil
}

func (s *PayrollService) CalculateRun(ctx context.Context, companyID, runID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return fmt.Errorf("calculate: load run: %w", err)
	}
	if run.Status != "OPEN" && run.Status != "CALCULATED" {
		return fmt.Errorf("calculate: run status %s does not allow calculation", run.Status)
	}

	now := time.Now()
	if err := s.repo.UpdateRunStatus(ctx, runID, "CALCULATING"); err != nil {
		return fmt.Errorf("calculate: update status: %w", err)
	}
	s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATING", &now, nil)

	rollback := func() {
		s.repo.UpdateRunStatus(ctx, runID, "OPEN")
		s.repo.UpdateRunTimestamps(ctx, runID, "OPEN", nil, nil)
	}

	period, err := s.repo.GetPeriod(ctx, companyID, run.PeriodID)
	if err != nil {
		rollback()
		return fmt.Errorf("calculate: load period: %w", err)
	}

	concepts, err := s.repo.ListConcepts(ctx, companyID, nil, nil, boolPtr(true))
	if err != nil {
		rollback()
		return fmt.Errorf("calculate: load concepts: %w", err)
	}

	employees, err := s.repo.ListRunEmployees(ctx, runID)
	if err != nil {
		rollback()
		return fmt.Errorf("calculate: load employees: %w", err)
	}

	if len(concepts) == 0 || len(employees) == 0 {
		finished := time.Now()
		s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATED", nil, &finished)
		return nil
	}

	conceptIDs := make([]string, len(concepts))
	for i, c := range concepts {
		conceptIDs[i] = c.ID
	}

	rules, err := s.repo.GetActiveRulesByConceptIDs(ctx, companyID, conceptIDs, period.StartDate)
	if err != nil {
		rollback()
		return fmt.Errorf("calculate: load rules: %w", err)
	}

	limits, _ := s.repo.GetActiveLimits(ctx, companyID, period.StartDate)
	minWage, _ := s.repo.GetMinimumWage(ctx, "AR", period.StartDate)

	for i := range employees {
		emp := &employees[i]

		oldStatus := emp.Status
		emp.Status = "CALCULATING"
		emp.CalculationVersion++
		s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)

		snapshot := s.buildSnapshot(ctx, companyID, emp, period)
		if snapshot != nil {
			if err := s.repo.CreateSnapshot(ctx, snapshot); err != nil {
				s.log.Warn("calculate: create snapshot", zap.Error(err))
			}
		}

		baseSalary, currency, err := s.repo.GetEmployeeCompensation(ctx, companyID, emp.EmployeeID)
		if err != nil {
			emp.Status = "ERROR"
			errMsg := err.Error()
			emp.ErrorMessage = &errMsg
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			s.recordError(ctx, runID, &emp.EmployeeID, "ERROR", "COMPENSATION_FETCH", err.Error(), nil)
			continue
		}
		emp.Currency = currency

		novelties, err := s.repo.GetNoveltiesForEmployeePeriod(ctx, companyID, emp.EmployeeID, run.PeriodID)
		if err != nil {
			emp.Status = "ERROR"
			errMsg := err.Error()
			emp.ErrorMessage = &errMsg
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			s.recordError(ctx, runID, &emp.EmployeeID, "ERROR", "NOVELTY_FETCH", err.Error(), nil)
			continue
		}

		input := engine.RuleInput{
			Employee: domain.EmployeeSnapshot{
				ID:            emp.ID,
				RunEmployeeID: emp.ID,
				EmployeeData: map[string]any{
					"employee_id": emp.EmployeeID,
				},
				SalaryData: map[string]any{
					"base_salary": baseSalary,
				},
			},
			Period:     *period,
			Run:        *run,
			Concepts:   concepts,
			Rules:      rules,
			Novelties:  novelties,
			Limits:     limits,
			BaseSalary: baseSalary,
			MinWage:    minWage,
		}

		result, err := s.engine.Evaluate(ctx, input)
		if err != nil {
			emp.Status = "ERROR"
			errMsg := err.Error()
			emp.ErrorMessage = &errMsg
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			s.recordError(ctx, runID, &emp.EmployeeID, "ERROR", "ENGINE_EVALUATE", err.Error(), nil)
			continue
		}

		if result == nil {
			emp.Status = oldStatus
			s.repo.UpdateRunEmployeeResult(ctx, emp.ID, emp)
			continue
		}

		emp.Status = "CALCULATED"
		emp.GrossRemunerative = result.GrossRemunerative
		emp.GrossNonRemunerative = result.GrossNonRemunerative
		emp.DeductionsAmount = result.TotalDeductions
		emp.EmployerContributions = result.TotalContributions
		emp.EmployerCost = result.EmployerCost
		emp.NetAmount = result.Net
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
			s.recordError(ctx, runID, &emp.EmployeeID, "WARNING", "CALC_WARN", w, nil)
		}
		for _, e := range result.Errors {
			s.recordError(ctx, runID, &emp.EmployeeID, "ERROR", "CALC_ERR", e, nil)
		}
	}

	finished := time.Now()
	s.repo.UpdateRunTimestamps(ctx, runID, "CALCULATED", nil, &finished)
	return nil
}

func (s *PayrollService) ValidateRun(ctx context.Context, companyID, runID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return fmt.Errorf("validate: load run: %w", err)
	}
	if run.Status != "CALCULATED" {
		return fmt.Errorf("validate: run must be CALCULATED, current: %s", run.Status)
	}

	s.repo.UpdateRunStatus(ctx, runID, "VALIDATING")

	employees, err := s.repo.ListRunEmployees(ctx, runID)
	if err != nil {
		s.repo.UpdateRunStatus(ctx, runID, "CALCULATED")
		return fmt.Errorf("validate: load employees: %w", err)
	}

	hasBlocking := false
	for _, emp := range employees {
		if emp.Status == "ERROR" {
			s.recordError(ctx, runID, &emp.EmployeeID, "BLOCKING", "EMPLOYEE_ERROR",
				fmt.Sprintf("employee %s has calculation error", emp.EmployeeID), nil)
			hasBlocking = true
		}
		if emp.NetAmount.IsNegative() {
			s.recordError(ctx, runID, &emp.EmployeeID, "WARNING", "NEGATIVE_NET",
				fmt.Sprintf("employee %s has negative net amount %s", emp.EmployeeID, emp.NetAmount.String()), nil)
		}
	}

	summary, err := s.repo.GetRunSummary(ctx, runID)
	if err == nil && summary.TotalGross.IsNegative() {
		s.recordError(ctx, runID, nil, "BLOCKING", "NEGATIVE_GROSS", "total gross is negative", nil)
		hasBlocking = true
	}

	if hasBlocking {
		s.repo.UpdateRunStatus(ctx, runID, "CALCULATED")
		return fmt.Errorf("validate: blocking errors found")
	}

	s.repo.UpdateRunStatus(ctx, runID, "VALIDATED")
	return nil
}

func (s *PayrollService) ApproveRun(ctx context.Context, companyID, runID, userID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return fmt.Errorf("approve: load run: %w", err)
	}
	if run.Status != "VALIDATED" {
		return fmt.Errorf("approve: run must be VALIDATED, current: %s", run.Status)
	}
	blocking, err := s.repo.ListBlockingErrors(ctx, runID)
	if err != nil {
		return fmt.Errorf("approve: list blocking errors: %w", err)
	}
	if len(blocking) > 0 {
		return fmt.Errorf("approve: cannot approve run with %d blocking errors", len(blocking))
	}
	if err := s.repo.ApproveRun(ctx, runID, userID); err != nil {
		return fmt.Errorf("approve: %w", err)
	}
	return nil
}

func (s *PayrollService) CloseRun(ctx context.Context, companyID, runID, userID string) error {
	run, err := s.repo.GetRun(ctx, companyID, runID)
	if err != nil {
		return fmt.Errorf("close run: load run: %w", err)
	}
	if run.Status != "APPROVED" {
		return fmt.Errorf("close run: run must be APPROVED, current: %s", run.Status)
	}
	if err := s.repo.CloseRun(ctx, runID); err != nil {
		return fmt.Errorf("close run: %w", err)
	}
	runs, err := s.repo.ListRuns(ctx, companyID, &run.PeriodID, nil, strPtr("APPROVED"), 0, 0)
	if err != nil || len(runs) == 0 {
		s.repo.UpdatePeriodStatus(ctx, run.PeriodID, "CLOSED")
	}
	return nil
}

func (s *PayrollService) buildSnapshot(ctx context.Context, companyID string, emp *domain.PayrollRunEmployee, period *domain.PayrollPeriod) *domain.EmployeeSnapshot {
	agreementID, categoryID, _ := s.repo.GetEmployeeAgreementCategory(ctx, companyID, emp.EmployeeID)
	baseSalary, _, _ := s.repo.GetEmployeeCompensation(ctx, companyID, emp.EmployeeID)

	return &domain.EmployeeSnapshot{
		ID:            uuid.NewString(),
		RunEmployeeID: emp.ID,
		EmployeeData: map[string]any{
			"id": emp.EmployeeID,
		},
		SalaryData: map[string]any{
			"base_salary": baseSalary,
		},
		AgreementData: map[string]any{
			"agreement_id": agreementID,
		},
		CategoryData: map[string]any{
			"category_id": categoryID,
		},
	}
}

func (s *PayrollService) recordError(ctx context.Context, runID string, employeeID *string, severity, code, message string, field *string) {
	e := &domain.PayrollError{
		ID:         uuid.NewString(),
		RunID:      runID,
		EmployeeID: employeeID,
		Severity:   severity,
		Code:       code,
		Message:    message,
		Field:      field,
	}
	if err := s.repo.CreateError(ctx, e); err != nil {
		s.log.Warn("failed to record payroll error", zap.String("code", code), zap.Error(err))
	}
}


