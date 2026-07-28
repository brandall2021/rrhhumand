package payroll

import (
	"time"
)

type PayrollPeriod struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Status      string     `json:"status"`
	CalculatedAt *time.Time `json:"calculated_at,omitempty"`
	ApprovedBy  *string    `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	ClosedBy    *string    `json:"closed_by,omitempty"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PayrollConcept struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	CalculationType string    `json:"calculation_type"`
	Taxable         bool      `json:"taxable"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
}

type EmployeeCompensation struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	BaseAmount   float64    `json:"base_amount"`
	Currency     string     `json:"currency"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo  *time.Time `json:"effective_to,omitempty"`
	Reason       *string    `json:"reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type PayrollItem struct {
	ID              string    `json:"id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	EmployeeID      string  `json:"employee_id"`
	ConceptID       string  `json:"concept_id"`
	ConceptCode     string  `json:"concept_code,omitempty"`
	ConceptName     string  `json:"concept_name,omitempty"`
	Quantity        float64 `json:"quantity"`
	UnitAmount      float64 `json:"unit_amount"`
	Amount          float64 `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}

type Benefit struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	BenefitType  *string   `json:"benefit_type,omitempty"`
	DefaultAmount float64 `json:"default_amount"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

type EmployeeBenefit struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	EmployeeID    string     `json:"employee_id"`
	EmployeeName  string     `json:"employee_name,omitempty"`
	BenefitID     string     `json:"benefit_id"`
	BenefitName   string     `json:"benefit_name,omitempty"`
	Amount        *float64   `json:"amount,omitempty"`
	Currency      string     `json:"currency"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}

type PayrollBonus struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	BonusType    string     `json:"bonus_type"`
	Amount       float64    `json:"amount"`
	Currency     string     `json:"currency"`
	Reason       *string    `json:"reason,omitempty"`
	PeriodStart  *time.Time `json:"period_start,omitempty"`
	PeriodEnd    *time.Time `json:"period_end,omitempty"`
	Status       string     `json:"status"`
	ApprovedBy   *string    `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type PayrollAdvance struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	Amount       float64    `json:"amount"`
	Currency     string     `json:"currency"`
	RequestDate  time.Time  `json:"request_date"`
	Reason       *string    `json:"reason,omitempty"`
	Status       string     `json:"status"`
	ApprovedBy   *string    `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type PayrollDeduction struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	Concept      string     `json:"concept"`
	Amount       float64    `json:"amount"`
	Currency     string     `json:"currency"`
	Reason       *string    `json:"reason,omitempty"`
	PeriodStart  *time.Time `json:"period_start,omitempty"`
	PeriodEnd    *time.Time `json:"period_end,omitempty"`
	Status       string     `json:"status"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type PayrollAdjustment struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	EmployeeID      string  `json:"employee_id"`
	Amount          float64 `json:"amount"`
	Reason          string  `json:"reason"`
	Type            string  `json:"type"`
	CreatedBy       *string `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type PayrollLedgerEntry struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	EmployeeID      string  `json:"employee_id"`
	TransactionType string  `json:"transaction_type"`
	ConceptCode     *string `json:"concept_code,omitempty"`
	Amount          float64 `json:"amount"`
	BalanceAfter    float64 `json:"balance_after"`
	Description     *string `json:"description,omitempty"`
	CreatedBy       *string `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type PayrollSnapshot struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	SnapshotData    []byte   `json:"snapshot_data"`
	CreatedAt       time.Time `json:"created_at"`
}

type EmployeePayrollSummary struct {
	EmployeeID      string        `json:"employee_id"`
	EmployeeName    string        `json:"employee_name"`
	BaseSalary      float64       `json:"base_salary"`
	Currency        string        `json:"currency"`
	OvertimeAmount  float64       `json:"overtime_amount"`
	BonusTotal      float64       `json:"bonus_total"`
	BenefitTotal    float64       `json:"benefit_total"`
	DeductionTotal  float64       `json:"deduction_total"`
	AdvanceTotal    float64       `json:"advance_total"`
	GrossTotal      float64       `json:"gross_total"`
	NetTotal        float64       `json:"net_total"`
	Items           []PayrollItem `json:"items,omitempty"`
}

type PayrollReview struct {
	TotalEmployees    int     `json:"total_employees"`
	Calculated        int     `json:"calculated"`
	Pending           int     `json:"pending"`
	Errors            int     `json:"errors"`
	Warnings          int     `json:"warnings"`
	TotalGross        float64 `json:"total_gross"`
	TotalNet          float64 `json:"total_net"`
	ErrorDetails      []PayrollError   `json:"error_details,omitempty"`
	WarningDetails    []PayrollWarning `json:"warning_details,omitempty"`
}

type PayrollError struct {
	EmployeeID string `json:"employee_id"`
	Message    string `json:"message"`
}

type PayrollWarning struct {
	EmployeeID string `json:"employee_id"`
	Message    string `json:"message"`
}

type PayrollDashboard struct {
	TotalEmployees int     `json:"total_employees"`
	TotalGross     float64 `json:"total_gross"`
	TotalNet       float64 `json:"total_net"`
	TotalOvertime  float64 `json:"total_overtime"`
	TotalBonuses   float64 `json:"total_bonuses"`
	TotalBenefits  float64 `json:"total_benefits"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalAdvances  float64 `json:"total_advances"`
}

type PayrollFilters struct {
	EmployeeID string
	Status     string
	DateFrom   string
	DateTo     string
}
