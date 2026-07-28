package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/rrhhumand/api/internal/payroll/repository"
)

type CreateNoveltyInput struct {
	EmployeeID  string
	NoveltyType string
	Quantity    *float64
	Unit        *string
	Amount      *float64
	UnitValue   *float64
	Multiplier  *float64
	StartDate   *string
	EndDate     *string
	Description *string
	Source      string
}

type UpdateNoveltyInput struct {
	Quantity    *float64
	Amount      *float64
	Description *string
	Status      *string
}

func (s *PayrollService) CreateNovelty(ctx context.Context, companyID, userID string, req CreateNoveltyInput) (*domain.PayrollNovelty, error) {
	qty := decimalFromFloatPtr(req.Quantity)
	amt := decimalFromFloatPtr(req.Amount)
	uv := decimalFromFloatPtr(req.UnitValue)
	mult := decimalFromFloatPtr(req.Multiplier)

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("create novelty: parse start_date: %w", err)
		}
		startDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("create novelty: parse end_date: %w", err)
		}
		endDate = &t
	}

	n := &domain.PayrollNovelty{
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
		return nil, fmt.Errorf("create novelty: %w", err)
	}
	return n, nil
}

func (s *PayrollService) UpdateNovelty(ctx context.Context, companyID, id string, req UpdateNoveltyInput) (*domain.PayrollNovelty, error) {
	n, err := s.repo.GetNovelty(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("update novelty: get: %w", err)
	}
	if req.Quantity != nil {
		n.Quantity = decimalFromFloatPtr(req.Quantity)
	}
	if req.Amount != nil {
		n.Amount = decimalFromFloatPtr(req.Amount)
	}
	if req.Description != nil {
		n.Description = req.Description
	}
	if req.Status != nil {
		n.Status = *req.Status
	}
	if err := s.repo.UpdateNovelty(ctx, n); err != nil {
		return nil, fmt.Errorf("update novelty: save: %w", err)
	}
	return n, nil
}

func (s *PayrollService) GetNovelty(ctx context.Context, companyID, id string) (*domain.PayrollNovelty, error) {
	n, err := s.repo.GetNovelty(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("get novelty: %w", err)
	}
	return n, nil
}

func (s *PayrollService) ListNovelties(ctx context.Context, companyID string, employeeID, periodID, noveltyType, status, source *string, limit, offset int) ([]domain.PayrollNovelty, error) {
	if limit <= 0 {
		limit = 20
	}
	novelties, err := s.repo.ListNovelties(ctx, companyID, repository.NoveltyFilter{
		EmployeeID:  employeeID,
		PeriodID:    periodID,
		NoveltyType: noveltyType,
		Status:      status,
		Source:      source,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list novelties: %w", err)
	}
	return novelties, nil
}

func (s *PayrollService) DeleteNovelty(ctx context.Context, companyID, id string) error {
	if err := s.repo.DeleteNovelty(ctx, companyID, id); err != nil {
		return fmt.Errorf("delete novelty: %w", err)
	}
	return nil
}

func (s *PayrollService) ApproveNovelty(ctx context.Context, companyID, id, userID string) error {
	n, err := s.repo.GetNovelty(ctx, companyID, id)
	if err != nil {
		return fmt.Errorf("approve novelty: get: %w", err)
	}
	if n.Status != "PENDING" {
		return fmt.Errorf("approve novelty: novelty is not pending, current: %s", n.Status)
	}
	if err := s.repo.ApproveNovelty(ctx, id, userID); err != nil {
		return fmt.Errorf("approve novelty: %w", err)
	}
	return nil
}

func (s *PayrollService) ImportNovelties(ctx context.Context, companyID, userID string, novelties []CreateNoveltyInput) ([]domain.PayrollNovelty, error) {
	result := make([]domain.PayrollNovelty, 0, len(novelties))
	for _, nr := range novelties {
		n, err := s.CreateNovelty(ctx, companyID, userID, nr)
		if err != nil {
			return result, fmt.Errorf("import novelties: %w", err)
		}
		result = append(result, *n)
	}
	return result, nil
}

func decimalFromFloatPtr(f *float64) *decimal.Decimal {
	if f == nil {
		return nil
	}
	d := decimal.NewFromFloat(*f)
	return &d
}
