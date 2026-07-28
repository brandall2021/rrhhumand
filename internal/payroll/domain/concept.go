package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayrollConcept struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	Code            string     `json:"code"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	ConceptType     string     `json:"concept_type"`
	Taxability      string     `json:"taxability"`
	CalculationType string     `json:"calculation_type"`
	BaseConceptID   *string    `json:"base_concept_id,omitempty"`
	Active          bool       `json:"active"`
	EffectiveFrom   time.Time  `json:"effective_from"`
	EffectiveTo     *time.Time `json:"effective_to,omitempty"`
	SortOrder       int        `json:"sort_order"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type PayrollRule struct {
	ID            string         `json:"id"`
	CompanyID     string         `json:"company_id"`
	Country       string         `json:"country"`
	Jurisdiction  *string        `json:"jurisdiction,omitempty"`
	AgreementID   *string        `json:"agreement_id,omitempty"`
	CategoryID    *string        `json:"category_id,omitempty"`
	ConceptID     string         `json:"concept_id"`
	RuleType      string         `json:"rule_type"`
	Formula       *string        `json:"formula,omitempty"`
	Parameters    map[string]any `json:"parameters"`
	Priority      int            `json:"priority"`
	EffectiveFrom time.Time      `json:"effective_from"`
	EffectiveTo   *time.Time     `json:"effective_to,omitempty"`
	Version       int            `json:"version"`
	Active        bool           `json:"active"`
	CreatedBy     string         `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type PayrollNovelty struct {
	ID                string           `json:"id"`
	CompanyID         string           `json:"company_id"`
	EmployeeID        string           `json:"employee_id"`
	PeriodID          *string          `json:"period_id,omitempty"`
	NoveltyType       string           `json:"novelty_type"`
	Quantity          *decimal.Decimal `json:"quantity,omitempty"`
	Unit              *string          `json:"unit,omitempty"`
	Amount            *decimal.Decimal `json:"amount,omitempty"`
	UnitValue         *decimal.Decimal `json:"unit_value,omitempty"`
	Multiplier        *decimal.Decimal `json:"multiplier,omitempty"`
	StartDate         *time.Time       `json:"start_date,omitempty"`
	EndDate           *time.Time       `json:"end_date,omitempty"`
	Description       *string          `json:"description,omitempty"`
	Source            string           `json:"source"`
	SourceReferenceID *string          `json:"source_reference_id,omitempty"`
	Status            string           `json:"status"`
	ApprovedBy        *string          `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time       `json:"approved_at,omitempty"`
	CreatedBy         string           `json:"created_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
