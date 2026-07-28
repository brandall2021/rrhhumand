package payroll

import "time"

type CreatePeriodRequest struct {
	Year        int       `json:"year" binding:"required"`
	Month       int       `json:"month" binding:"required"`
	PeriodType  string    `json:"period_type" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
}

type UpdatePeriodRequest struct {
	Name        string     `json:"name,omitempty"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
}

type CreateRunRequest struct {
	RunType string `json:"run_type" binding:"required"`
}

type CreateConceptRequest struct {
	Code            string  `json:"code" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     *string `json:"description,omitempty"`
	ConceptType     string  `json:"concept_type" binding:"required"`
	Taxability      string  `json:"taxability" binding:"required"`
	CalculationType string  `json:"calculation_type" binding:"required"`
	BaseConceptID   *string `json:"base_concept_id,omitempty"`
	SortOrder       int     `json:"sort_order"`
}

type UpdateConceptRequest struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ConceptType     *string `json:"concept_type,omitempty"`
	Taxability      *string `json:"taxability,omitempty"`
	CalculationType *string `json:"calculation_type,omitempty"`
	BaseConceptID   *string `json:"base_concept_id,omitempty"`
	Active          *bool   `json:"active,omitempty"`
	SortOrder       *int    `json:"sort_order,omitempty"`
}

type CreateRuleRequest struct {
	ConceptID    string `json:"concept_id" binding:"required"`
	RuleType     string `json:"rule_type" binding:"required"`
	Formula      *string `json:"formula,omitempty"`
	Parameters   map[string]any `json:"parameters"`
	Priority     int    `json:"priority"`
	EffectiveFrom string `json:"effective_from" binding:"required"`
	EffectiveTo  *string `json:"effective_to,omitempty"`
}

type UpdateRuleRequest struct {
	RuleType     *string        `json:"rule_type,omitempty"`
	Formula      *string        `json:"formula,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	Priority     *int           `json:"priority,omitempty"`
	Active       *bool          `json:"active,omitempty"`
	EffectiveTo  *string        `json:"effective_to,omitempty"`
}

type CreateNoveltyRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	NoveltyType string  `json:"novelty_type" binding:"required"`
	Quantity    *float64 `json:"quantity,omitempty"`
	Unit        *string `json:"unit,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	UnitValue   *float64 `json:"unit_value,omitempty"`
	Multiplier  *float64 `json:"multiplier,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
	Description *string `json:"description,omitempty"`
	Source      string  `json:"source"`
}

type UpdateNoveltyRequest struct {
	Quantity    *float64 `json:"quantity,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type ImportNoveltiesRequest struct {
	Novelties []CreateNoveltyRequest `json:"novelties" binding:"required"`
}

type CreateAdvanceRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	RequestDate string  `json:"request_date"`
	Installments int    `json:"installments"`
	Reason      *string `json:"reason,omitempty"`
}

type CreateGarnishmentRequest struct {
	EmployeeID       string   `json:"employee_id" binding:"required"`
	CourtOrderNumber string   `json:"court_order_number" binding:"required"`
	CourtName        *string  `json:"court_name,omitempty"`
	Type             string   `json:"type"`
	Percentage       *float64 `json:"percentage,omitempty"`
	FixedAmount      *float64 `json:"fixed_amount,omitempty"`
	Priority         int      `json:"priority"`
	EffectiveFrom    string   `json:"effective_from" binding:"required"`
}

type CreateAgreementRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description,omitempty"`
	Activity      *string `json:"activity,omitempty"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
}

type CreateCategoryRequest struct {
	AgreementID *string `json:"agreement_id,omitempty"`
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

type CreateSalaryScaleRequest struct {
	AgreementID   *string  `json:"agreement_id,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	MinimumSalary float64  `json:"minimum_salary" binding:"required"`
	MaximumSalary *float64 `json:"maximum_salary,omitempty"`
}

type RunFilter struct {
	PeriodID   *string `form:"period_id"`
	RunType    *string `form:"run_type"`
	Status     *string `form:"status"`
	Limit      int     `form:"limit"`
	Offset     int     `form:"offset"`
}

type NoveltyFilter struct {
	EmployeeID  *string `form:"employee_id"`
	PeriodID    *string `form:"period_id"`
	NoveltyType *string `form:"novelty_type"`
	Status      *string `form:"status"`
	Source      *string `form:"source"`
	Limit       int     `form:"limit"`
	Offset      int     `form:"offset"`
}

type ConceptFilter struct {
	Active      *bool   `form:"active"`
	ConceptType *string `form:"concept_type"`
	Taxability  *string `form:"taxability"`
}

type EmployeeResult struct {
	EmployeeID          string  `json:"employee_id"`
	EmployeeName        string  `json:"employee_name"`
	Status              string  `json:"status"`
	GrossRemunerative   float64 `json:"gross_remunerative"`
	GrossNonRemunerative float64 `json:"gross_non_remunerative"`
	DeductionsAmount    float64 `json:"deductions_amount"`
	NetAmount           float64 `json:"net_amount"`
	EmployerContributions float64 `json:"employer_contributions"`
	EmployerCost        float64 `json:"employer_cost"`
	Error               *string `json:"error,omitempty"`
}
