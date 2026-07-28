package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayrollRunEmployee struct {
	ID                   string          `json:"id"`
	RunID                string          `json:"run_id"`
	EmployeeID           string          `json:"employee_id"`
	EmploymentID         *string         `json:"employment_id,omitempty"`
	Status               string          `json:"status"`
	GrossRemunerative    decimal.Decimal `json:"gross_remunerative"`
	GrossNonRemunerative decimal.Decimal `json:"gross_non_remunerative"`
	DeductionsAmount     decimal.Decimal `json:"deductions_amount"`
	EmployerContributions decimal.Decimal `json:"employer_contributions"`
	EmployerCost         decimal.Decimal `json:"employer_cost"`
	NetAmount            decimal.Decimal `json:"net_amount"`
	Currency             string          `json:"currency"`
	CalculationVersion   int             `json:"calculation_version"`
	ErrorMessage         *string         `json:"error_message,omitempty"`
	CalculatedAt         *time.Time      `json:"calculated_at,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
}

type EmployeeSnapshot struct {
	ID                 string          `json:"id"`
	RunEmployeeID      string          `json:"run_employee_id"`
	EmployeeData       map[string]any  `json:"employee_data"`
	EmploymentData     map[string]any  `json:"employment_data,omitempty"`
	PositionData       map[string]any  `json:"position_data,omitempty"`
	CategoryData       map[string]any  `json:"category_data,omitempty"`
	AgreementData      map[string]any  `json:"agreement_data,omitempty"`
	SalaryData         map[string]any  `json:"salary_data,omitempty"`
	BenefitsData       map[string]any  `json:"benefits_data,omitempty"`
	TaxConfigData      map[string]any  `json:"tax_config_data,omitempty"`
	SocialSecurityData map[string]any  `json:"social_security_data,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}
