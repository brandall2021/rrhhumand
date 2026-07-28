package application

import (
	"context"
	"fmt"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (s *PayrollService) GetRunSummary(ctx context.Context, companyID, runID string) (*domain.PayrollSummary, error) {
	summary, err := s.repo.GetRunSummary(ctx, runID)
	if err != nil {
		return nil, fmt.Errorf("get run summary: %w", err)
	}
	return summary, nil
}

func (s *PayrollService) GetDashboardStats(ctx context.Context, companyID string) (*domain.DashboardStats, error) {
	stats, err := s.repo.GetDashboardStats(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("get dashboard stats: %w", err)
	}
	return stats, nil
}

func (s *PayrollService) GetEmployeeResult(ctx context.Context, companyID, runID, employeeID string) (*domain.PayrollRunEmployee, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("get employee result: %w", err)
	}
	return re, nil
}

func (s *PayrollService) GetEmployeeItems(ctx context.Context, companyID, runID, employeeID string) ([]domain.PayrollItem, error) {
	re, err := s.repo.GetRunEmployee(ctx, runID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("get employee items: get run employee: %w", err)
	}
	items, err := s.repo.ListItems(ctx, re.ID)
	if err != nil {
		return nil, fmt.Errorf("get employee items: list items: %w", err)
	}
	return items, nil
}

func (s *PayrollService) GetErrors(ctx context.Context, runID string) ([]domain.PayrollError, error) {
	errors, err := s.repo.ListErrors(ctx, runID)
	if err != nil {
		return nil, fmt.Errorf("get errors: %w", err)
	}
	return errors, nil
}
