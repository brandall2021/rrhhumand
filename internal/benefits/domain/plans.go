package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BenefitPlan struct {
	ID                 uuid.UUID       `json:"id"`
	CompanyID          uuid.UUID       `json:"company_id"`
	BenefitID          uuid.UUID       `json:"benefit_id"`
	Name               string          `json:"name"`
	Description        *string         `json:"description,omitempty"`
	PlanType           string          `json:"plan_type"`
	MonthlyCostEmployee decimal.Decimal `json:"monthly_cost_employee"`
	MonthlyCostEmployer decimal.Decimal `json:"monthly_cost_employer"`
	AnnualCostEmployee  decimal.Decimal `json:"annual_cost_employee"`
	AnnualCostEmployer  decimal.Decimal `json:"annual_cost_employer"`
	Currency            string          `json:"currency"`
	CoverageLimit       *decimal.Decimal `json:"coverage_limit,omitempty"`
	CoverageDetails     map[string]any  `json:"coverage_details,omitempty"`
	MaxDependents       int             `json:"max_dependents"`
	DependentType       *string         `json:"dependent_type,omitempty"`
	EnrollmentFee       *decimal.Decimal `json:"enrollment_fee,omitempty"`
	WaitingPeriodDays   int             `json:"waiting_period_days"`
	MinimumAge          *int            `json:"minimum_age,omitempty"`
	MaximumAge          *int            `json:"maximum_age,omitempty"`
	IsDefault           bool            `json:"is_default"`
	IsActive            bool            `json:"is_active"`
	SortOrder           int             `json:"sort_order"`
	CreatedBy           uuid.UUID       `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type BenefitEligibilityRule struct {
	ID            uuid.UUID  `json:"id"`
	CompanyID     uuid.UUID  `json:"company_id"`
	BenefitID     uuid.UUID  `json:"benefit_id"`
	RuleType      string     `json:"rule_type"`
	Operator      string     `json:"operator"`
	Value         string     `json:"value"`
	ValueTo       *string    `json:"value_to,omitempty"`
	LogicGroup    int        `json:"logic_group"`
	LogicOperator string     `json:"logic_operator"`
	Priority      int        `json:"priority"`
	ErrorMessage  *string    `json:"error_message,omitempty"`
	IsActive      bool       `json:"is_active"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type BenefitWorkflow struct {
	ID                      uuid.UUID  `json:"id"`
	CompanyID               uuid.UUID  `json:"company_id"`
	BenefitID               *uuid.UUID `json:"benefit_id,omitempty"`
	WorkflowType            string     `json:"workflow_type"`
	Name                    string     `json:"name"`
	Description             *string    `json:"description,omitempty"`
	RequiresChainApproval   bool       `json:"requires_chain_approval"`
	AutoApprove             bool       `json:"auto_approve"`
	AutoApproveIfNoManager  bool       `json:"auto_approve_if_no_manager"`
	EscalationHours         int        `json:"escalation_hours"`
	IsActive                bool       `json:"is_active"`
	CreatedBy               uuid.UUID  `json:"created_by"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type BenefitWorkflowStep struct {
	ID                    uuid.UUID  `json:"id"`
	WorkflowID            uuid.UUID  `json:"workflow_id"`
	StepOrder             int        `json:"step_order"`
	ApprovalType          string     `json:"approval_type"`
	ApproverRoleID        *uuid.UUID `json:"approver_role_id,omitempty"`
	MaxRejectionCount     int        `json:"max_rejection_count"`
	IsRequired            bool       `json:"is_required"`
	NotificationTemplate  *string    `json:"notification_template,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}
