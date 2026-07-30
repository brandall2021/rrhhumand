package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

type CreateAgreementInput struct {
	Code          string
	Name          string
	Description   *string
	Activity      *string
	EffectiveFrom string
}

type CreateCategoryInput struct {
	AgreementID *string
	Code        string
	Name        string
	Description *string
}

type CreateSalaryScaleInput struct {
	AgreementID   *string
	CategoryID    *string
	MinimumSalary decimal.Decimal
	MaximumSalary *decimal.Decimal
}

type CreateAdvanceInput struct {
	EmployeeID   string
	Amount       decimal.Decimal
	RequestDate  string
	Installments int
	Reason       *string
}

type CreateGarnishmentInput struct {
	EmployeeID       string
	CourtOrderNumber string
	CourtName        *string
	Type             string
	Percentage       *decimal.Decimal
	FixedAmount      *decimal.Decimal
	Priority         int
	EffectiveFrom    string
}

func (s *PayrollService) CreateAgreement(ctx context.Context, companyID, userID string, req CreateAgreementInput) (*domain.LaborAgreement, error) {
	effFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return nil, fmt.Errorf("create agreement: parse effective_from: %w", err)
	}
	a := &domain.LaborAgreement{
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
		return nil, fmt.Errorf("create agreement: %w", err)
	}
	return a, nil
}

func (s *PayrollService) ListAgreements(ctx context.Context, companyID string) ([]domain.LaborAgreement, error) {
	agreements, err := s.repo.ListAgreements(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list agreements: %w", err)
	}
	return agreements, nil
}

func (s *PayrollService) CreateCategory(ctx context.Context, companyID string, req CreateCategoryInput) (*domain.LaborCategory, error) {
	c := &domain.LaborCategory{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		AgreementID:   req.AgreementID,
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		EffectiveFrom: time.Now().AddDate(0, -1, 0),
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return c, nil
}

func (s *PayrollService) ListCategories(ctx context.Context, companyID string) ([]domain.LaborCategory, error) {
	categories, err := s.repo.ListCategories(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (s *PayrollService) CreateSalaryScale(ctx context.Context, companyID string, req CreateSalaryScaleInput) (*domain.SalaryScale, error) {
	sc := &domain.SalaryScale{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		AgreementID:   req.AgreementID,
		CategoryID:    req.CategoryID,
		MinimumSalary: req.MinimumSalary,
		MaximumSalary: req.MaximumSalary,
		EffectiveFrom: time.Now().AddDate(0, -1, 0),
	}
	if err := s.repo.CreateSalaryScale(ctx, sc); err != nil {
		return nil, fmt.Errorf("create salary scale: %w", err)
	}
	return sc, nil
}

func (s *PayrollService) ListSalaryScales(ctx context.Context, companyID string) ([]domain.SalaryScale, error) {
	scales, err := s.repo.ListSalaryScales(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list salary scales: %w", err)
	}
	return scales, nil
}

func (s *PayrollService) CreateAdvance(ctx context.Context, companyID, userID string, req CreateAdvanceInput) (*domain.EmployeeAdvance, error) {
	reqDate, err := time.Parse("2006-01-02", req.RequestDate)
	if err != nil {
		return nil, fmt.Errorf("create advance: parse request_date: %w", err)
	}
	installmentAmt := req.Amount.Div(decimal.NewFromInt(int64(req.Installments)))
	a := &domain.EmployeeAdvance{
		ID:                uuid.NewString(),
		CompanyID:         companyID,
		EmployeeID:        req.EmployeeID,
		Amount:            req.Amount,
		Currency:          "ARS",
		RequestDate:       reqDate,
		Installments:      req.Installments,
		InstallmentAmount: &installmentAmt,
		RemainingAmount:   req.Amount,
		Reason:            req.Reason,
		Status:            "PENDING",
		CreatedBy:         userID,
	}
	if err := s.repo.CreateAdvance(ctx, a); err != nil {
		return nil, fmt.Errorf("create advance: %w", err)
	}
	return a, nil
}

func (s *PayrollService) ListAdvances(ctx context.Context, companyID, employeeID string) ([]domain.EmployeeAdvance, error) {
	advances, err := s.repo.ListAdvances(ctx, companyID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list advances: %w", err)
	}
	return advances, nil
}

func (s *PayrollService) CreateGarnishment(ctx context.Context, companyID, userID string, req CreateGarnishmentInput) (*domain.PayrollGarnishment, error) {
	effFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return nil, fmt.Errorf("create garnishment: parse effective_from: %w", err)
	}
	g := &domain.PayrollGarnishment{
		ID:               uuid.NewString(),
		CompanyID:        companyID,
		EmployeeID:       req.EmployeeID,
		CourtOrderNumber: req.CourtOrderNumber,
		CourtName:        req.CourtName,
		Type:             req.Type,
		Percentage:       req.Percentage,
		FixedAmount:      req.FixedAmount,
		Priority:         req.Priority,
		EffectiveFrom:    effFrom,
		Status:           "ACTIVE",
		CreatedBy:        userID,
	}
	if err := s.repo.CreateGarnishment(ctx, g); err != nil {
		return nil, fmt.Errorf("create garnishment: %w", err)
	}
	return g, nil
}

func (s *PayrollService) ListGarnishments(ctx context.Context, companyID, employeeID string) ([]domain.PayrollGarnishment, error) {
	garnishments, err := s.repo.ListGarnishments(ctx, companyID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list garnishments: %w", err)
	}
	return garnishments, nil
}
