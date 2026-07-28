package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayrollItem struct {
	ID                    string          `json:"id"`
	RunEmployeeID         string          `json:"run_employee_id"`
	ConceptID             string          `json:"concept_id"`
	Quantity              decimal.Decimal `json:"quantity"`
	UnitValue             decimal.Decimal `json:"unit_value"`
	BaseAmount            decimal.Decimal `json:"base_amount"`
	Rate                  *decimal.Decimal `json:"rate,omitempty"`
	Amount                decimal.Decimal `json:"amount"`
	IsRemunerative        bool            `json:"is_remunerative"`
	IsDeduction           bool            `json:"is_deduction"`
	IsEmployerContribution bool           `json:"is_employer_contribution"`
	CalculationDetail     map[string]any  `json:"calculation_detail,omitempty"`
	SortOrder             int             `json:"sort_order"`
	CreatedAt             time.Time       `json:"created_at"`
}

type PayrollBase struct {
	ID                string          `json:"id"`
	RunEmployeeID     string          `json:"run_employee_id"`
	BaseType          string          `json:"base_type"`
	BaseAmount        decimal.Decimal `json:"base_amount"`
	ConceptIDs        []string        `json:"concept_ids"`
	CalculationDetail map[string]any  `json:"calculation_detail,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

type PayrollDeduction struct {
	ID             string          `json:"id"`
	RunEmployeeID  string          `json:"run_employee_id"`
	ConceptID      string          `json:"concept_id"`
	BaseAmount     decimal.Decimal `json:"base_amount"`
	Rate           *decimal.Decimal `json:"rate,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	CapAmount      *decimal.Decimal `json:"cap_amount,omitempty"`
	ExceededAmount decimal.Decimal `json:"exceeded_amount"`
	CreatedAt      time.Time       `json:"created_at"`
}

type PayrollContribution struct {
	ID             string          `json:"id"`
	RunEmployeeID  string          `json:"run_employee_id"`
	ConceptID      string          `json:"concept_id"`
	BaseAmount     decimal.Decimal `json:"base_amount"`
	Rate           *decimal.Decimal `json:"rate,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	CapAmount      *decimal.Decimal `json:"cap_amount,omitempty"`
	ExceededAmount decimal.Decimal `json:"exceeded_amount"`
	CreatedAt      time.Time       `json:"created_at"`
}

type PayrollResult struct {
	EmployeeID          string
	GrossRemunerative   decimal.Decimal
	GrossNonRemunerative decimal.Decimal
	EmployeeDeductions  decimal.Decimal
	Net                 decimal.Decimal
	EmployerContributions decimal.Decimal
	EmployerCost        decimal.Decimal
	Items               []PayrollItem
	Bases               []PayrollBase
	Warnings            []string
	Errors              []string
}

type PayrollSummary struct {
	TotalEmployees     int             `json:"total_employees"`
	CalculatedEmployees int            `json:"calculated_employees"`
	ErrorEmployees     int             `json:"error_employees"`
	TotalGross         decimal.Decimal `json:"total_gross"`
	TotalDeductions    decimal.Decimal `json:"total_deductions"`
	TotalNet           decimal.Decimal `json:"total_net"`
	TotalContributions decimal.Decimal `json:"total_contributions"`
	TotalEmployerCost  decimal.Decimal `json:"total_employer_cost"`
}

type DashboardStats struct {
	ActivePeriods  int `json:"active_periods"`
	PendingRuns    int `json:"pending_runs"`
	TotalErrors    int `json:"total_errors"`
	BlockingErrors int `json:"blocking_errors"`
}
