package compensation

import (
	"time"

	"github.com/shopspring/decimal"
)

type CompensationStructure struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Currency    string     `json:"currency"`
	EffectiveFrom string   `json:"effective_from"`
	EffectiveTo *string    `json:"effective_to"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SalaryGrade struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	StructureID string    `json:"structure_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	SortOrder   int       `json:"sort_order"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SalaryBand struct {
	ID            string          `json:"id"`
	CompanyID     string          `json:"company_id"`
	StructureID   string          `json:"structure_id"`
	GradeID       *string         `json:"grade_id"`
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	MinimumAmount decimal.Decimal `json:"minimum_amount"`
	MidpointAmount decimal.Decimal `json:"midpoint_amount"`
	MaximumAmount decimal.Decimal `json:"maximum_amount"`
	Currency      string          `json:"currency"`
	Status        string          `json:"status"`
	CreatedBy     string          `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type PositionSalaryBand struct {
	ID           string    `json:"id"`
	PositionID   string    `json:"position_id"`
	SalaryBandID string    `json:"salary_band_id"`
	EffectiveFrom string   `json:"effective_from"`
	EffectiveTo  *string   `json:"effective_to"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type EmployeeCompensation struct {
	ID            string          `json:"id"`
	CompanyID     string          `json:"company_id"`
	EmployeeID    string          `json:"employee_id"`
	SalaryBandID  *string         `json:"salary_band_id"`
	BaseAmount    decimal.Decimal `json:"base_amount"`
	Currency      string          `json:"currency"`
	PayFrequency  string          `json:"pay_frequency"`
	EffectiveFrom string          `json:"effective_from"`
	EffectiveTo   *string         `json:"effective_to"`
	Status        string          `json:"status"`
	CreatedBy     string          `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CompensationComponent struct {
	ID            string  `json:"id"`
	CompanyID     string  `json:"company_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	ComponentType string  `json:"component_type"`
	Taxable       bool    `json:"taxable"`
	Recurring     bool    `json:"recurring"`
	Active        bool    `json:"active"`
	CreatedBy     string  `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmployeeCompensationComponent struct {
	ID           string          `json:"id"`
	CompanyID    string          `json:"company_id"`
	EmployeeID   string          `json:"employee_id"`
	ComponentID  string          `json:"component_id"`
	Amount       decimal.Decimal `json:"amount"`
	Currency     string          `json:"currency"`
	EffectiveFrom string         `json:"effective_from"`
	EffectiveTo  *string         `json:"effective_to"`
	CreatedBy    string          `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CompensationHistory struct {
	ID            string          `json:"id"`
	CompanyID     string          `json:"company_id"`
	EmployeeID    string          `json:"employee_id"`
	PreviousAmount *decimal.Decimal `json:"previous_amount"`
	NewAmount     decimal.Decimal `json:"new_amount"`
	Currency      string          `json:"currency"`
	Reason        string          `json:"reason"`
	EffectiveFrom string          `json:"effective_from"`
	ApprovedBy    *string         `json:"approved_by"`
	Notes         *string         `json:"notes"`
	CreatedBy     string          `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CompensationAdjustment struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	EmployeeID      string          `json:"employee_id"`
	AdjustmentType  string          `json:"adjustment_type"`
	Value           decimal.Decimal `json:"value"`
	Currency        string          `json:"currency"`
	Reason          string          `json:"reason"`
	EffectiveFrom   string          `json:"effective_from"`
	Status          string          `json:"status"`
	ApprovedBy      *string         `json:"approved_by"`
	ApprovedAt      *time.Time      `json:"approved_at"`
	AppliedBy       *string         `json:"applied_by"`
	AppliedAt       *time.Time      `json:"applied_at"`
	Notes           *string         `json:"notes"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type SalaryAdjustmentProposal struct {
	ID               string          `json:"id"`
	CompanyID        string          `json:"company_id"`
	ReviewID         *string         `json:"review_id"`
	EmployeeID       string          `json:"employee_id"`
	CurrentAmount    decimal.Decimal `json:"current_amount"`
	ProposedAmount   decimal.Decimal `json:"proposed_amount"`
	IncreasePercentage *decimal.Decimal `json:"increase_percentage"`
	Reason           string          `json:"reason"`
	PerformanceScore *decimal.Decimal `json:"performance_score"`
	MarketPosition   *string         `json:"market_position"`
	ManagerComment   *string         `json:"manager_comment"`
	HRComment        *string         `json:"hr_comment"`
	Status           string          `json:"status"`
	SubmittedBy      *string         `json:"submitted_by"`
	ApprovedBy       *string         `json:"approved_by"`
	ApprovedAt       *time.Time      `json:"approved_at"`
	RejectedBy       *string         `json:"rejected_by"`
	RejectedAt       *time.Time      `json:"rejected_at"`
	RejectionReason  *string         `json:"rejection_reason"`
	CreatedBy        string          `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type BonusPlan struct {
	ID                string    `json:"id"`
	CompanyID         string    `json:"company_id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description"`
	Period            string    `json:"period"`
	TargetPercentage  *decimal.Decimal `json:"target_percentage"`
	MaximumPercentage *decimal.Decimal `json:"maximum_percentage"`
	EligibilityRules  *string   `json:"eligibility_rules"`
	Status            string    `json:"status"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Bonus struct {
	ID          string          `json:"id"`
	CompanyID   string          `json:"company_id"`
	EmployeeID  string          `json:"employee_id"`
	BonusPlanID *string         `json:"bonus_plan_id"`
	Name        string          `json:"name"`
	BonusType   string          `json:"bonus_type"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string          `json:"currency"`
	Period      *string         `json:"period"`
	Reason      *string         `json:"reason"`
	Status      string          `json:"status"`
	ApprovedBy  *string         `json:"approved_by"`
	ApprovedAt  *time.Time      `json:"approved_at"`
	PaidAt      *time.Time      `json:"paid_at"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Benefit struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	BenefitType string  `json:"benefit_type"`
	Provider    *string `json:"provider"`
	CostAmount  *decimal.Decimal `json:"cost_amount"`
	CostCurrency string `json:"cost_currency"`
	Frequency   string  `json:"frequency"`
	Taxable     bool    `json:"taxable"`
	Active      bool    `json:"active"`
	CreatedBy   string  `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeeBenefit struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	EmployeeID     string          `json:"employee_id"`
	BenefitID      string          `json:"benefit_id"`
	EnrollmentDate string          `json:"enrollment_date"`
	EffectiveFrom  string          `json:"effective_from"`
	EffectiveTo    *string         `json:"effective_to"`
	EmployeeCost   decimal.Decimal `json:"employee_cost"`
	CompanyCost    decimal.Decimal `json:"company_cost"`
	Currency       string          `json:"currency"`
	Status         string          `json:"status"`
	CreatedBy      string          `json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CompensationReview struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Period      string    `json:"period"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	Budget      *decimal.Decimal `json:"budget"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CompensationBudget struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	Year           int             `json:"year"`
	DepartmentID   *string         `json:"department_id"`
	BudgetAmount   decimal.Decimal `json:"budget_amount"`
	CommittedAmount decimal.Decimal `json:"committed_amount"`
	SpentAmount    decimal.Decimal `json:"spent_amount"`
	Currency       string          `json:"currency"`
	Status         string          `json:"status"`
	CreatedBy      string          `json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CompensationEquitySnapshot struct {
	ID                  string          `json:"id"`
	CompanyID           string          `json:"company_id"`
	SnapshotDate        string          `json:"snapshot_date"`
	DepartmentID        *string         `json:"department_id"`
	PositionID          *string         `json:"position_id"`
	GradeID             *string         `json:"grade_id"`
	EmployeeCount       int             `json:"employee_count"`
	MedianCompensation  *decimal.Decimal `json:"median_compensation"`
	AverageCompensation *decimal.Decimal `json:"average_compensation"`
	MinCompensation     *decimal.Decimal `json:"min_compensation"`
	MaxCompensation     *decimal.Decimal `json:"max_compensation"`
	Currency            string          `json:"currency"`
	Metadata            *string         `json:"metadata"`
	CreatedBy           string          `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
}

type CompensationAuditLog struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *string   `json:"entity_id"`
	OldValue   *string   `json:"old_value"`
	NewValue   *string   `json:"new_value"`
	IPAddress  *string   `json:"ip_address"`
	UserAgent  *string   `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

type CompensationDomainEvent struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	EventType   string    `json:"event_type"`
	EntityType  string    `json:"entity_type"`
	EntityID    *string   `json:"entity_id"`
	Payload     *string   `json:"payload"`
	CreatedBy   *string   `json:"created_by"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Derived types for calculations

type TotalCompensation struct {
	BaseSalary           decimal.Decimal `json:"base_salary"`
	FixedComponents      decimal.Decimal `json:"fixed_components"`
	VariableCompensation decimal.Decimal `json:"variable_compensation"`
	Benefits             decimal.Decimal `json:"benefits"`
	Total                decimal.Decimal `json:"total"`
	Currency             string          `json:"currency"`
}

type CompaRatioResult struct {
	Ratio       decimal.Decimal `json:"ratio"`
	Category    string          `json:"category"`
	Salary      decimal.Decimal `json:"salary"`
	Midpoint    decimal.Decimal `json:"midpoint"`
}

type RangePenetrationResult struct {
	Penetration decimal.Decimal `json:"penetration"`
	Salary      decimal.Decimal `json:"salary"`
	Minimum     decimal.Decimal `json:"minimum"`
	Maximum     decimal.Decimal `json:"maximum"`
}

type BandAnalysis struct {
	Band              SalaryBand               `json:"band"`
	EmployeeCount     int                      `json:"employee_count"`
	AverageSalary     *decimal.Decimal          `json:"average_salary"`
	MedianSalary      *decimal.Decimal          `json:"median_salary"`
	MinSalary         *decimal.Decimal          `json:"min_salary"`
	MaxSalary         *decimal.Decimal          `json:"max_salary"`
	BelowRange        int                      `json:"below_range"`
	InRange           int                      `json:"in_range"`
	AboveRange        int                      `json:"above_range"`
}

type DashboardStats struct {
	TotalSalaryCost      decimal.Decimal `json:"total_salary_cost"`
	AverageCompensation  decimal.Decimal `json:"average_compensation"`
	BenefitCost          decimal.Decimal `json:"benefit_cost"`
	TotalBonuses         decimal.Decimal `json:"total_bonuses"`
	BudgetTotal          decimal.Decimal `json:"budget_total"`
	BudgetUsed           decimal.Decimal `json:"budget_used"`
	PendingProposals     int             `json:"pending_proposals"`
	EmployeesOutOfBand   int             `json:"employees_out_of_band"`
	ActiveReviews        int             `json:"active_reviews"`
	Currency             string          `json:"currency"`
}
