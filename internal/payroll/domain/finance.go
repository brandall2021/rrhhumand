package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type EmployeeAdvance struct {
	ID                string          `json:"id"`
	CompanyID         string          `json:"company_id"`
	EmployeeID        string          `json:"employee_id"`
	Amount            decimal.Decimal `json:"amount"`
	Currency          string          `json:"currency"`
	RequestDate       time.Time       `json:"request_date"`
	Installments      int             `json:"installments"`
	InstallmentAmount *decimal.Decimal `json:"installment_amount,omitempty"`
	RemainingAmount   decimal.Decimal `json:"remaining_amount"`
	Reason            *string         `json:"reason,omitempty"`
	Status            string          `json:"status"`
	ApprovedBy        *string         `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty"`
	CreatedBy         string          `json:"created_by"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type PayrollGarnishment struct {
	ID               string          `json:"id"`
	CompanyID        string          `json:"company_id"`
	EmployeeID       string          `json:"employee_id"`
	CourtOrderNumber string          `json:"court_order_number"`
	CourtName        *string         `json:"court_name,omitempty"`
	Type             string          `json:"type"`
	Percentage       *decimal.Decimal `json:"percentage,omitempty"`
	FixedAmount      *decimal.Decimal `json:"fixed_amount,omitempty"`
	Priority         int             `json:"priority"`
	EffectiveFrom    time.Time       `json:"effective_from"`
	EffectiveTo      *time.Time      `json:"effective_to,omitempty"`
	Status           string          `json:"status"`
	Notes            *string         `json:"notes,omitempty"`
	CreatedBy        string          `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type PayrollWithholding struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	EmployeeID      string          `json:"employee_id"`
	WithholdingType string          `json:"withholding_type"`
	ConceptID       *string         `json:"concept_id,omitempty"`
	BaseAmount      *decimal.Decimal `json:"base_amount,omitempty"`
	Rate            *decimal.Decimal `json:"rate,omitempty"`
	Amount          *decimal.Decimal `json:"amount,omitempty"`
	Priority        int             `json:"priority"`
	EffectiveFrom   time.Time       `json:"effective_from"`
	EffectiveTo     *time.Time      `json:"effective_to,omitempty"`
	Status          string          `json:"status"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
