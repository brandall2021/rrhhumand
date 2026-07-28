package payroll

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type PayrollPeriod struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	Year        int        `json:"year"`
	Month       int        `json:"month"`
	PeriodType  string     `json:"period_type"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
	Status      string     `json:"status"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PayrollRun struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	PeriodID      string     `json:"period_id"`
	RunNumber     int        `json:"run_number"`
	RunType       string     `json:"run_type"`
	Status        string     `json:"status"`
	EngineVersion string     `json:"engine_version"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	CreatedBy     string     `json:"created_by"`
	ApprovedBy    *string    `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ClosedBy      *string    `json:"closed_by,omitempty"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

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
	ID                string          `json:"id"`
	RunEmployeeID     string          `json:"run_employee_id"`
	EmployeeData      json.RawMessage `json:"employee_data"`
	EmploymentData    json.RawMessage `json:"employment_data,omitempty"`
	PositionData      json.RawMessage `json:"position_data,omitempty"`
	CategoryData      json.RawMessage `json:"category_data,omitempty"`
	AgreementData     json.RawMessage `json:"agreement_data,omitempty"`
	SalaryData        json.RawMessage `json:"salary_data,omitempty"`
	BenefitsData      json.RawMessage `json:"benefits_data,omitempty"`
	TaxConfigData     json.RawMessage `json:"tax_config_data,omitempty"`
	SocialSecurityData json.RawMessage `json:"social_security_data,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

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
	ID           string          `json:"id"`
	CompanyID    string          `json:"company_id"`
	Country      string          `json:"country"`
	Jurisdiction *string         `json:"jurisdiction,omitempty"`
	AgreementID  *string         `json:"agreement_id,omitempty"`
	CategoryID   *string         `json:"category_id,omitempty"`
	ConceptID    string          `json:"concept_id"`
	RuleType     string          `json:"rule_type"`
	Formula      *string         `json:"formula,omitempty"`
	Parameters   json.RawMessage `json:"parameters"`
	Priority     int             `json:"priority"`
	EffectiveFrom time.Time      `json:"effective_from"`
	EffectiveTo  *time.Time      `json:"effective_to,omitempty"`
	Version      int             `json:"version"`
	Active       bool            `json:"active"`
	CreatedBy    string          `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type PayrollNovelty struct {
	ID                string     `json:"id"`
	CompanyID         string     `json:"company_id"`
	EmployeeID        string     `json:"employee_id"`
	PeriodID          *string    `json:"period_id,omitempty"`
	NoveltyType       string     `json:"novelty_type"`
	Quantity          *decimal.Decimal `json:"quantity,omitempty"`
	Unit              *string    `json:"unit,omitempty"`
	Amount            *decimal.Decimal `json:"amount,omitempty"`
	UnitValue         *decimal.Decimal `json:"unit_value,omitempty"`
	Multiplier        *decimal.Decimal `json:"multiplier,omitempty"`
	StartDate         *time.Time `json:"start_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Source            string     `json:"source"`
	SourceReferenceID *string    `json:"source_reference_id,omitempty"`
	Status            string     `json:"status"`
	ApprovedBy        *string    `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type PayrollItem struct {
	ID                  string          `json:"id"`
	RunEmployeeID       string          `json:"run_employee_id"`
	ConceptID           string          `json:"concept_id"`
	Quantity            decimal.Decimal `json:"quantity"`
	UnitValue           decimal.Decimal `json:"unit_value"`
	BaseAmount          decimal.Decimal `json:"base_amount"`
	Rate                *decimal.Decimal `json:"rate,omitempty"`
	Amount              decimal.Decimal `json:"amount"`
	IsRemunerative      bool            `json:"is_remunerative"`
	IsDeduction         bool            `json:"is_deduction"`
	IsEmployerContribution bool         `json:"is_employer_contribution"`
	CalculationDetail   json.RawMessage `json:"calculation_detail,omitempty"`
	SortOrder           int             `json:"sort_order"`
	CreatedAt           time.Time       `json:"created_at"`
}

type PayrollBase struct {
	ID              string          `json:"id"`
	RunEmployeeID   string          `json:"run_employee_id"`
	BaseType        string          `json:"base_type"`
	BaseAmount      decimal.Decimal `json:"base_amount"`
	ConceptIDs      []string        `json:"concept_ids"`
	CalculationDetail json.RawMessage `json:"calculation_detail,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type PayrollDeduction struct {
	ID              string          `json:"id"`
	RunEmployeeID   string          `json:"run_employee_id"`
	ConceptID       string          `json:"concept_id"`
	BaseAmount      decimal.Decimal `json:"base_amount"`
	Rate            *decimal.Decimal `json:"rate,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	CapAmount       *decimal.Decimal `json:"cap_amount,omitempty"`
	ExceededAmount  decimal.Decimal `json:"exceeded_amount"`
	CreatedAt       time.Time       `json:"created_at"`
}

type PayrollContribution struct {
	ID              string          `json:"id"`
	RunEmployeeID   string          `json:"run_employee_id"`
	ConceptID       string          `json:"concept_id"`
	BaseAmount      decimal.Decimal `json:"base_amount"`
	Rate            *decimal.Decimal `json:"rate,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	CapAmount       *decimal.Decimal `json:"cap_amount,omitempty"`
	ExceededAmount  decimal.Decimal `json:"exceeded_amount"`
	CreatedAt       time.Time       `json:"created_at"`
}

type LaborAgreement struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Activity    *string    `json:"activity,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type LaborCategory struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	AgreementID *string    `json:"agreement_id,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
	ID            string    `json:"id"`
	Country       string    `json:"country"`
	Jurisdiction  *string   `json:"jurisdiction,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string    `json:"currency"`
	Source        *string   `json:"source,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
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
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	EmployeeID      string          `json:"employee_id"`
	CourtOrderNumber string         `json:"court_order_number"`
	CourtName       *string         `json:"court_name,omitempty"`
	Type            string          `json:"type"`
	Percentage      *decimal.Decimal `json:"percentage,omitempty"`
	FixedAmount     *decimal.Decimal `json:"fixed_amount,omitempty"`
	Priority        int             `json:"priority"`
	EffectiveFrom   time.Time       `json:"effective_from"`
	EffectiveTo     *time.Time      `json:"effective_to,omitempty"`
	Status          string          `json:"status"`
	Notes           *string         `json:"notes,omitempty"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type PayrollWithholding struct {
	ID               string          `json:"id"`
	CompanyID        string          `json:"company_id"`
	EmployeeID       string          `json:"employee_id"`
	WithholdingType  string          `json:"withholding_type"`
	ConceptID        *string         `json:"concept_id,omitempty"`
	BaseAmount       *decimal.Decimal `json:"base_amount,omitempty"`
	Rate             *decimal.Decimal `json:"rate,omitempty"`
	Amount           *decimal.Decimal `json:"amount,omitempty"`
	Priority         int             `json:"priority"`
	EffectiveFrom    time.Time       `json:"effective_from"`
	EffectiveTo      *time.Time      `json:"effective_to,omitempty"`
	Status           string          `json:"status"`
	CreatedBy        string          `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type PayrollError struct {
	ID         string     `json:"id"`
	RunID      string     `json:"run_id"`
	EmployeeID *string    `json:"employee_id,omitempty"`
	Severity   string     `json:"severity"`
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Field      *string    `json:"field,omitempty"`
	Resolved   bool       `json:"resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type PayrollAuditLog struct {
	ID         string          `json:"id"`
	CompanyID  string          `json:"company_id"`
	UserID     string          `json:"user_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   *string         `json:"entity_id,omitempty"`
	OldValues  json.RawMessage `json:"old_values,omitempty"`
	NewValues  json.RawMessage `json:"new_values,omitempty"`
	IPAddress  *string         `json:"ip_address,omitempty"`
	UserAgent  *string         `json:"user_agent,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type RuleContext struct {
	Employee   EmployeeSnapshot
	Period     PayrollPeriod
	Run        PayrollRun
	Concepts   []PayrollConcept
	Novelties  []PayrollNovelty
	Bases      map[string]decimal.Decimal
	Parameters map[string]any
}

type RuleResult struct {
	Amount    decimal.Decimal
	Detail    string
	BaseUsed  decimal.Decimal
	RateUsed  *decimal.Decimal
}

type PayrollResult struct {
	EmployeeID          string
	GrossRemunerative    decimal.Decimal
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
	TotalEmployees      int             `json:"total_employees"`
	CalculatedEmployees  int             `json:"calculated_employees"`
	ErrorEmployees      int             `json:"error_employees"`
	TotalGross          decimal.Decimal `json:"total_gross"`
	TotalDeductions     decimal.Decimal `json:"total_deductions"`
	TotalNet            decimal.Decimal `json:"total_net"`
	TotalContributions  decimal.Decimal `json:"total_contributions"`
	TotalEmployerCost   decimal.Decimal `json:"total_employer_cost"`
}

type DashboardStats struct {
	ActivePeriods      int `json:"active_periods"`
	PendingRuns        int `json:"pending_runs"`
	TotalErrors        int `json:"total_errors"`
	BlockingErrors     int `json:"blocking_errors"`
}
