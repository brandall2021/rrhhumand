package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/rrhhumand/api/internal/payroll/repository"
)

type CreatePeriodInput struct {
	Year        int
	Month       int
	PeriodType  string
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	PaymentDate *time.Time
}

type UpdatePeriodInput struct {
	Name        string
	PaymentDate *time.Time
}

func (s *PayrollService) CreatePeriod(ctx context.Context, companyID, userID string, req CreatePeriodInput) (*domain.PayrollPeriod, error) {
	p := &domain.PayrollPeriod{
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
		return nil, fmt.Errorf("create period: %w", err)
	}
	return p, nil
}

func (s *PayrollService) UpdatePeriod(ctx context.Context, companyID, id string, req UpdatePeriodInput) (*domain.PayrollPeriod, error) {
	p, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("update period: get: %w", err)
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.PaymentDate != nil {
		p.PaymentDate = req.PaymentDate
	}
	if err := s.repo.UpdatePeriod(ctx, p); err != nil {
		return nil, fmt.Errorf("update period: save: %w", err)
	}
	return p, nil
}

func (s *PayrollService) GetPeriod(ctx context.Context, companyID, id string) (*domain.PayrollPeriod, error) {
	p, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("get period: %w", err)
	}
	return p, nil
}

func (s *PayrollService) ListPeriods(ctx context.Context, companyID string, limit, offset int) ([]domain.PayrollPeriod, error) {
	if limit <= 0 {
		limit = 20
	}
	periods, err := s.repo.ListPeriods(ctx, companyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list periods: %w", err)
	}
	return periods, nil
}

func (s *PayrollService) ClosePeriod(ctx context.Context, companyID, id, userID string) error {
	p, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil {
		return fmt.Errorf("close period: get: %w", err)
	}
	if p.Status == "CLOSED" {
		return fmt.Errorf("close period: period already closed")
	}
	if err := s.repo.ClosePeriod(ctx, id, userID); err != nil {
		return fmt.Errorf("close period: %w", err)
	}
	return nil
}

func (s *PayrollService) CreateRun(ctx context.Context, companyID, periodID, userID string, runType string) (*domain.PayrollRun, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, fmt.Errorf("create run: get period: %w", err)
	}
	if period.Status == "CLOSED" {
		return nil, fmt.Errorf("create run: cannot create run on closed period")
	}
	runNumber, err := s.repo.GetRunNumber(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("create run: get run number: %w", err)
	}
	run := &domain.PayrollRun{
		ID:        uuid.NewString(),
		CompanyID: companyID,
		PeriodID:  periodID,
		RunNumber: runNumber,
		RunType:   runType,
		Status:    "OPEN",
		CreatedBy: userID,
	}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, fmt.Errorf("create run: %w", err)
	}
	return run, nil
}

func (s *PayrollService) GetRun(ctx context.Context, companyID, id string) (*domain.PayrollRun, error) {
	run, err := s.repo.GetRun(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("get run: %w", err)
	}
	return run, nil
}

func (s *PayrollService) ListRuns(ctx context.Context, companyID string, periodID, runType, status *string, limit, offset int) ([]domain.PayrollRun, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListRuns(ctx, companyID, repository.RunFilter{
		PeriodID: periodID,
		RunType:  runType,
		Status:   status,
		Limit:    limit,
		Offset:   offset,
	})
}
