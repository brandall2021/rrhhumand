package http

import (
	"time"

	"github.com/rrhhumand/api/internal/payroll/application"
)

// ========================================================================
// PERIOD
// ========================================================================

type CreatePeriodReq struct {
	Year        int        `json:"year" binding:"required"`
	Month       int        `json:"month" binding:"required"`
	PeriodType  string     `json:"period_type" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     time.Time  `json:"end_date" binding:"required"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
}

func (r CreatePeriodReq) ToInput() application.CreatePeriodInput {
	return application.CreatePeriodInput{
		Year: r.Year, Month: r.Month, PeriodType: r.PeriodType, Name: r.Name,
		StartDate: r.StartDate, EndDate: r.EndDate, PaymentDate: r.PaymentDate,
	}
}

type UpdatePeriodReq struct {
	Name        string     `json:"name,omitempty"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
}

func (r UpdatePeriodReq) ToInput() application.UpdatePeriodInput {
	return application.UpdatePeriodInput{Name: r.Name, PaymentDate: r.PaymentDate}
}

type CreateRunReq struct {
	RunType string `json:"run_type" binding:"required"`
}

// ========================================================================
// CONCEPT
// ========================================================================

type CreateConceptReq struct {
	Code            string  `json:"code" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     *string `json:"description,omitempty"`
	ConceptType     string  `json:"concept_type" binding:"required"`
	Taxability      string  `json:"taxability" binding:"required"`
	CalculationType string  `json:"calculation_type" binding:"required"`
	BaseConceptID   *string `json:"base_concept_id,omitempty"`
	SortOrder       int     `json:"sort_order"`
}

func (r CreateConceptReq) ToInput() application.CreateConceptInput {
	return application.CreateConceptInput{
		Code: r.Code, Name: r.Name, Description: r.Description,
		ConceptType: r.ConceptType, Taxability: r.Taxability, CalculationType: r.CalculationType,
		BaseConceptID: r.BaseConceptID, SortOrder: r.SortOrder,
	}
}

type UpdateConceptReq struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ConceptType     *string `json:"concept_type,omitempty"`
	Taxability      *string `json:"taxability,omitempty"`
	CalculationType *string `json:"calculation_type,omitempty"`
	BaseConceptID   *string `json:"base_concept_id,omitempty"`
	Active          *bool   `json:"active,omitempty"`
	SortOrder       *int    `json:"sort_order,omitempty"`
}

func (r UpdateConceptReq) ToInput() application.UpdateConceptInput {
	return application.UpdateConceptInput{
		Name: r.Name, Description: r.Description, ConceptType: r.ConceptType,
		Taxability: r.Taxability, CalculationType: r.CalculationType,
		BaseConceptID: r.BaseConceptID, Active: r.Active, SortOrder: r.SortOrder,
	}
}

// ========================================================================
// RULE
// ========================================================================

type CreateRuleReq struct {
	ConceptID     string         `json:"concept_id" binding:"required"`
	RuleType      string         `json:"rule_type" binding:"required"`
	Formula       *string        `json:"formula,omitempty"`
	Parameters    map[string]any `json:"parameters"`
	Priority      int            `json:"priority"`
	EffectiveFrom string         `json:"effective_from" binding:"required"`
	EffectiveTo   *string        `json:"effective_to,omitempty"`
}

func (r CreateRuleReq) ToInput() application.CreateRuleInput {
	return application.CreateRuleInput{
		ConceptID: r.ConceptID, RuleType: r.RuleType, Formula: r.Formula,
		Parameters: r.Parameters, Priority: r.Priority,
		EffectiveFrom: &r.EffectiveFrom, EffectiveTo: r.EffectiveTo,
	}
}

type UpdateRuleReq struct {
	RuleType    *string        `json:"rule_type,omitempty"`
	Formula     *string        `json:"formula,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Priority    *int           `json:"priority,omitempty"`
	Active      *bool          `json:"active,omitempty"`
	EffectiveTo *string        `json:"effective_to,omitempty"`
}

func (r UpdateRuleReq) ToInput() application.UpdateRuleInput {
	return application.UpdateRuleInput{
		RuleType: r.RuleType, Formula: r.Formula, Parameters: r.Parameters,
		Priority: r.Priority, Active: r.Active, EffectiveTo: r.EffectiveTo,
	}
}

// ========================================================================
// NOVELTY
// ========================================================================

type CreateNoveltyReq struct {
	EmployeeID  string   `json:"employee_id" binding:"required"`
	NoveltyType string   `json:"novelty_type" binding:"required"`
	Quantity    *float64 `json:"quantity,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	UnitValue   *float64 `json:"unit_value,omitempty"`
	Multiplier  *float64 `json:"multiplier,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	Description *string  `json:"description,omitempty"`
	Source      string   `json:"source"`
}

func (r CreateNoveltyReq) ToInput() application.CreateNoveltyInput {
	return application.CreateNoveltyInput{
		EmployeeID: r.EmployeeID, NoveltyType: r.NoveltyType,
		Quantity: r.Quantity, Unit: r.Unit, Amount: r.Amount, UnitValue: r.UnitValue,
		Multiplier: r.Multiplier, StartDate: r.StartDate, EndDate: r.EndDate,
		Description: r.Description, Source: r.Source,
	}
}

type UpdateNoveltyReq struct {
	Quantity    *float64 `json:"quantity,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

func (r UpdateNoveltyReq) ToInput() application.UpdateNoveltyInput {
	return application.UpdateNoveltyInput{Quantity: r.Quantity, Amount: r.Amount, Description: r.Description, Status: r.Status}
}

type ImportNoveltiesReq struct {
	Novelties []CreateNoveltyReq `json:"novelties" binding:"required"`
}

// ========================================================================
// ADVANCE
// ========================================================================

type CreateAdvanceReq struct {
	EmployeeID  string   `json:"employee_id" binding:"required"`
	Amount      float64  `json:"amount" binding:"required"`
	RequestDate string   `json:"request_date"`
	Installments int     `json:"installments"`
	Reason      *string  `json:"reason,omitempty"`
}

func (r CreateAdvanceReq) ToInput() application.CreateAdvanceInput {
	return application.CreateAdvanceInput{
		EmployeeID: r.EmployeeID, Amount: r.Amount, RequestDate: r.RequestDate,
		Installments: r.Installments, Reason: r.Reason,
	}
}

// ========================================================================
// GARNISHMENT
// ========================================================================

type CreateGarnishmentReq struct {
	EmployeeID       string   `json:"employee_id" binding:"required"`
	CourtOrderNumber string   `json:"court_order_number" binding:"required"`
	CourtName        *string  `json:"court_name,omitempty"`
	Type             string   `json:"type"`
	Percentage       *float64 `json:"percentage,omitempty"`
	FixedAmount      *float64 `json:"fixed_amount,omitempty"`
	Priority         int      `json:"priority"`
	EffectiveFrom    string   `json:"effective_from" binding:"required"`
}

func (r CreateGarnishmentReq) ToInput() application.CreateGarnishmentInput {
	return application.CreateGarnishmentInput{
		EmployeeID: r.EmployeeID, CourtOrderNumber: r.CourtOrderNumber,
		CourtName: r.CourtName, Type: r.Type, Percentage: r.Percentage,
		FixedAmount: r.FixedAmount, Priority: r.Priority, EffectiveFrom: r.EffectiveFrom,
	}
}

// ========================================================================
// AGREEMENTS & CATEGORIES
// ========================================================================

type CreateAgreementReq struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description,omitempty"`
	Activity      *string `json:"activity,omitempty"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
}

func (r CreateAgreementReq) ToInput() application.CreateAgreementInput {
	return application.CreateAgreementInput{
		Code: r.Code, Name: r.Name, Description: r.Description,
		Activity: r.Activity, EffectiveFrom: r.EffectiveFrom,
	}
}

type CreateCategoryReq struct {
	AgreementID *string `json:"agreement_id,omitempty"`
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

func (r CreateCategoryReq) ToInput() application.CreateCategoryInput {
	return application.CreateCategoryInput{
		AgreementID: r.AgreementID, Code: r.Code, Name: r.Name, Description: r.Description,
	}
}

type CreateSalaryScaleReq struct {
	AgreementID   *string  `json:"agreement_id,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	MinimumSalary float64  `json:"minimum_salary" binding:"required"`
	MaximumSalary *float64 `json:"maximum_salary,omitempty"`
}

func (r CreateSalaryScaleReq) ToInput() application.CreateSalaryScaleInput {
	return application.CreateSalaryScaleInput{
		AgreementID: r.AgreementID, CategoryID: r.CategoryID,
		MinimumSalary: r.MinimumSalary, MaximumSalary: r.MaximumSalary,
	}
}
