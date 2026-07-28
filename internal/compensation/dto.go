package compensation

import "github.com/shopspring/decimal"

// Structures
type CreateStructureRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	Currency      string  `json:"currency"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	EffectiveTo   *string `json:"effective_to"`
}

type UpdateStructureRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Currency      *string `json:"currency"`
	EffectiveFrom *string `json:"effective_from"`
	EffectiveTo   *string `json:"effective_to"`
	Status        *string `json:"status"`
}

// Grades
type CreateGradeRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	SortOrder *int   `json:"sort_order"`
}

type UpdateGradeRequest struct {
	Code      *string `json:"code"`
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
	Status    *string `json:"status"`
}

// Bands
type CreateBandRequest struct {
	GradeID       *string `json:"grade_id"`
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	MinimumAmount float64 `json:"minimum_amount" binding:"required"`
	MidpointAmount float64 `json:"midpoint_amount" binding:"required"`
	MaximumAmount float64 `json:"maximum_amount" binding:"required"`
	Currency      string  `json:"currency"`
}

type UpdateBandRequest struct {
	Name          *string  `json:"name"`
	GradeID       *string  `json:"grade_id"`
	MinimumAmount *float64 `json:"minimum_amount"`
	MidpointAmount *float64 `json:"midpoint_amount"`
	MaximumAmount *float64 `json:"maximum_amount"`
	Currency      *string  `json:"currency"`
	Status        *string  `json:"status"`
}

// Position-Band
type AssignPositionBandRequest struct {
	SalaryBandID  string  `json:"salary_band_id" binding:"required"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	EffectiveTo   *string `json:"effective_to"`
}

// Employee Compensation
type SetEmployeeCompensationRequest struct {
	SalaryBandID  *string `json:"salary_band_id"`
	BaseAmount    float64 `json:"base_amount" binding:"required"`
	Currency      string  `json:"currency"`
	PayFrequency  string  `json:"pay_frequency"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
}

type UpdateEmployeeCompensationRequest struct {
	BaseAmount    *float64 `json:"base_amount"`
	SalaryBandID  *string  `json:"salary_band_id"`
	Currency      *string  `json:"currency"`
	PayFrequency  *string  `json:"pay_frequency"`
	Status        *string  `json:"status"`
}

// Components
type CreateComponentRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	ComponentType string  `json:"component_type"`
	Taxable       *bool   `json:"taxable"`
	Recurring     *bool   `json:"recurring"`
}

type UpdateComponentRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	ComponentType *string `json:"component_type"`
	Taxable       *bool   `json:"taxable"`
	Recurring     *bool   `json:"recurring"`
	Active        *bool   `json:"active"`
}

// Employee Components
type AssignComponentRequest struct {
	ComponentID   string  `json:"component_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	EffectiveTo   *string `json:"effective_to"`
}

// Adjustments
type CreateAdjustmentRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	AdjustmentType string `json:"adjustment_type" binding:"required"`
	Value         float64 `json:"value" binding:"required"`
	Currency      string  `json:"currency"`
	Reason        string  `json:"reason" binding:"required"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	Notes         *string `json:"notes"`
}

// Proposals
type CreateProposalRequest struct {
	ReviewID        *string `json:"review_id"`
	EmployeeID      string  `json:"employee_id" binding:"required"`
	CurrentAmount   float64 `json:"current_amount" binding:"required"`
	ProposedAmount  float64 `json:"proposed_amount" binding:"required"`
	Reason          string  `json:"reason" binding:"required"`
	PerformanceScore *float64 `json:"performance_score"`
	MarketPosition  *string `json:"market_position"`
	ManagerComment  *string `json:"manager_comment"`
}

// Bonus Plans
type CreateBonusPlanRequest struct {
	Name              string   `json:"name" binding:"required"`
	Description       *string  `json:"description"`
	Period            string   `json:"period"`
	TargetPercentage  *float64 `json:"target_percentage"`
	MaximumPercentage *float64 `json:"maximum_percentage"`
	EligibilityRules  *string  `json:"eligibility_rules"`
}

// Bonuses
type CreateBonusRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	BonusPlanID *string `json:"bonus_plan_id"`
	Name        string  `json:"name" binding:"required"`
	BonusType   string  `json:"bonus_type"`
	Amount      float64 `json:"amount" binding:"required"`
	Currency    string  `json:"currency"`
	Period      *string `json:"period"`
	Reason      *string `json:"reason"`
}

// Benefits
type CreateBenefitRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	BenefitType string  `json:"benefit_type"`
	Provider    *string `json:"provider"`
	CostAmount  *float64 `json:"cost_amount"`
	CostCurrency string `json:"cost_currency"`
	Frequency   string  `json:"frequency"`
	Taxable     *bool   `json:"taxable"`
}

type UpdateBenefitRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	BenefitType *string  `json:"benefit_type"`
	Provider    *string  `json:"provider"`
	CostAmount  *float64 `json:"cost_amount"`
	CostCurrency *string `json:"cost_currency"`
	Frequency   *string  `json:"frequency"`
	Taxable     *bool    `json:"taxable"`
	Active      *bool    `json:"active"`
}

// Employee Benefits
type AssignBenefitRequest struct {
	BenefitID     string  `json:"benefit_id" binding:"required"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	EffectiveTo   *string `json:"effective_to"`
	EmployeeCost  float64 `json:"employee_cost"`
	CompanyCost   float64 `json:"company_cost"`
	Currency      string  `json:"currency"`
}

// Reviews
type CreateReviewRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Period      string  `json:"period"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date"`
	Budget      *float64 `json:"budget"`
	Currency    string  `json:"currency"`
}

// Budgets
type CreateBudgetRequest struct {
	Year         int      `json:"year" binding:"required"`
	DepartmentID *string  `json:"department_id"`
	BudgetAmount float64  `json:"budget_amount" binding:"required"`
	Currency     string   `json:"currency"`
}

type UpdateBudgetRequest struct {
	BudgetAmount *float64 `json:"budget_amount"`
	Status       *string  `json:"status"`
}

// Reports
type EquityAnalysisRequest struct {
	DepartmentID *string `json:"department_id"`
	PositionID   *string `json:"position_id"`
	GradeID      *string `json:"grade_id"`
}

type EquityAnalysisResult struct {
	Groups          []EquityGroup `json:"groups"`
	OverallMedian   *float64      `json:"overall_median"`
	OverallAverage  *float64      `json:"overall_average"`
	Outliers        []string      `json:"outliers"`
}

type EquityGroup struct {
	Label            string   `json:"label"`
	EmployeeCount    int      `json:"employee_count"`
	MinCompensation  float64  `json:"min_compensation"`
	MaxCompensation  float64  `json:"max_compensation"`
	AverageCompensation float64 `json:"average_compensation"`
	MedianCompensation  float64 `json:"median_compensation"`
}

// AI
type AIRecommendationRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Context    string `json:"context"`
}

type AIAdjustmentRecommendation struct {
	EmployeeID         string  `json:"employee_id"`
	CurrentSalary      float64 `json:"current_salary"`
	RecommendedSalary  float64 `json:"recommended_salary"`
	IncreasePercentage float64 `json:"increase_percentage"`
	CompaRatio         float64 `json:"compa_ratio"`
	Reason             string  `json:"reason"`
	Confidence         string  `json:"confidence"`
}

// Filters
type CompensationFilter struct {
	Status     *string `json:"status"`
	EmployeeID *string `json:"employee_id"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type BonusFilter struct {
	EmployeeID *string `json:"employee_id"`
	Status     *string `json:"status"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type BenefitFilter struct {
	Active   *bool   `json:"active"`
	BenefitType *string `json:"benefit_type"`
}

type AdjustmentFilter struct {
	EmployeeID *string `json:"employee_id"`
	Status     *string `json:"status"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type ProposalFilter struct {
	ReviewID   *string `json:"review_id"`
	EmployeeID *string `json:"employee_id"`
	Status     *string `json:"status"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

var _ = decimal.Zero
