package engine

import (
	"context"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/shopspring/decimal"
)

type RuleEngine interface {
	Evaluate(ctx context.Context, input RuleInput) (*RuleOutput, error)
}

type RuleInput struct {
	Employee   domain.EmployeeSnapshot
	Period     domain.PayrollPeriod
	Run        domain.PayrollRun
	Concepts   []domain.PayrollConcept
	Rules      []domain.PayrollRule
	Novelties  []domain.PayrollNovelty
	Limits     []domain.PayrollLimit
	BaseSalary decimal.Decimal
	MinWage    *domain.StatutoryMinimumWage
}

type RuleOutput struct {
	Items                []domain.PayrollItem
	Bases                []domain.PayrollBase
	Deductions           []domain.PayrollDeduction
	Contributions        []domain.PayrollContribution
	GrossRemunerative    decimal.Decimal
	GrossNonRemunerative decimal.Decimal
	TotalDeductions      decimal.Decimal
	TotalContributions   decimal.Decimal
	Net                  decimal.Decimal
	EmployerCost         decimal.Decimal
	Warnings             []string
	Errors               []string
}
