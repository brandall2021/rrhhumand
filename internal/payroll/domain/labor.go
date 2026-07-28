package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type LaborAgreement struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Activity      *string    `json:"activity,omitempty"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Status        string     `json:"status"`
	CreatedBy     string     `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type LaborCategory struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	AgreementID   *string    `json:"agreement_id,omitempty"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type SalaryScale struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	AgreementID    *string         `json:"agreement_id,omitempty"`
	CategoryID     *string         `json:"category_id,omitempty"`
	MinimumSalary  decimal.Decimal `json:"minimum_salary"`
	MaximumSalary  *decimal.Decimal `json:"maximum_salary,omitempty"`
	ReferenceValue *decimal.Decimal `json:"reference_value,omitempty"`
	EffectiveFrom  time.Time       `json:"effective_from"`
	EffectiveTo    *time.Time      `json:"effective_to,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type StatutoryMinimumWage struct {
	ID            string          `json:"id"`
	Country       string          `json:"country"`
	Jurisdiction  *string         `json:"jurisdiction,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Source        *string         `json:"source,omitempty"`
	EffectiveFrom time.Time       `json:"effective_from"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type PayrollLimit struct {
	ID            string          `json:"id"`
	CompanyID     string          `json:"company_id"`
	ConceptID     *string         `json:"concept_id,omitempty"`
	LimitType     string          `json:"limit_type"`
	MinimumAmount *decimal.Decimal `json:"minimum_amount,omitempty"`
	MaximumAmount *decimal.Decimal `json:"maximum_amount,omitempty"`
	EffectiveFrom time.Time       `json:"effective_from"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
