package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EmployeeBenefit struct {
	ID                     uuid.UUID       `json:"id"`
	CompanyID              uuid.UUID       `json:"company_id"`
	EmployeeID             uuid.UUID       `json:"employee_id"`
	BenefitID              uuid.UUID       `json:"benefit_id"`
	PlanID                 *uuid.UUID      `json:"plan_id,omitempty"`
	ProviderID             *uuid.UUID      `json:"provider_id,omitempty"`
	Status                 string          `json:"status"`
	EnrollmentDate         time.Time       `json:"enrollment_date"`
	StartDate              *time.Time      `json:"start_date,omitempty"`
	EndDate                *time.Time      `json:"end_date,omitempty"`
	CancellationDate       *time.Time      `json:"cancellation_date,omitempty"`
	CancellationReason     *string         `json:"cancellation_reason,omitempty"`
	AutoRenew              bool            `json:"auto_renew"`
	RenewalDate            *time.Time      `json:"renewal_date,omitempty"`
	CoverageLevel          *string         `json:"coverage_level,omitempty"`
	Dependents             []map[string]any `json:"dependents,omitempty"`
	EmergencyContact       map[string]any  `json:"emergency_contact,omitempty"`
	BeneficiaryInfo        map[string]any  `json:"beneficiary_info,omitempty"`
	EmployeeCost           decimal.Decimal `json:"employee_cost"`
	EmployerCost           decimal.Decimal `json:"employer_cost"`
	Currency               string          `json:"currency"`
	PayrollDeductionEnabled bool           `json:"payroll_deduction_enabled"`
	PayrollDeductionAmount decimal.Decimal `json:"payroll_deduction_amount"`
	ExternalMemberID       *string         `json:"external_member_id,omitempty"`
	ExternalPolicyNumber   *string         `json:"external_policy_number,omitempty"`
	ExternalGroupNumber    *string         `json:"external_group_number,omitempty"`
	Documents              []map[string]any `json:"documents,omitempty"`
	Notes                  *string         `json:"notes,omitempty"`
	Source                 string          `json:"source"`
	EnrolledBy             uuid.UUID       `json:"enrolled_by"`
	EnrolledAt             time.Time       `json:"enrolled_at"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

type EmployeeBenefitHistory struct {
	ID               uuid.UUID      `json:"id"`
	EmployeeBenefitID uuid.UUID     `json:"employee_benefit_id"`
	EmployeeID       uuid.UUID      `json:"employee_id"`
	BenefitID        uuid.UUID      `json:"benefit_id"`
	Action           string         `json:"action"`
	PreviousValue    map[string]any `json:"previous_value,omitempty"`
	NewValue         map[string]any `json:"new_value,omitempty"`
	ChangeReason     *string        `json:"change_reason,omitempty"`
	ChangedBy        uuid.UUID      `json:"changed_by"`
	ChangedAt        time.Time      `json:"changed_at"`
}

type BenefitRequest struct {
	ID              uuid.UUID      `json:"id"`
	CompanyID       uuid.UUID      `json:"company_id"`
	EmployeeID      uuid.UUID      `json:"employee_id"`
	BenefitID       uuid.UUID      `json:"benefit_id"`
	PlanID          *uuid.UUID     `json:"plan_id,omitempty"`
	EmployeeBenefitID *uuid.UUID   `json:"employee_benefit_id,omitempty"`
	RequestType     string         `json:"request_type"`
	Status          string         `json:"status"`
	RequestData     map[string]any `json:"request_data,omitempty"`
	Justification   *string        `json:"justification,omitempty"`
	CoverageLevel   *string        `json:"coverage_level,omitempty"`
	Dependents      []map[string]any `json:"dependents,omitempty"`
	EffectiveDate   *time.Time     `json:"effective_date,omitempty"`
	Notes           *string        `json:"notes,omitempty"`
	SubmittedBy     *uuid.UUID     `json:"submitted_by,omitempty"`
	SubmittedAt     *time.Time     `json:"submitted_at,omitempty"`
	ResolvedBy      *uuid.UUID     `json:"resolved_by,omitempty"`
	ResolvedAt      *time.Time     `json:"resolved_at,omitempty"`
	ResolutionNotes *string        `json:"resolution_notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type BenefitRequestReview struct {
	ID         uuid.UUID  `json:"id"`
	RequestID  uuid.UUID  `json:"request_id"`
	StepID     *uuid.UUID `json:"step_id,omitempty"`
	ReviewerID uuid.UUID  `json:"reviewer_id"`
	ReviewType string     `json:"review_type"`
	Comment    *string    `json:"comment,omitempty"`
	ReviewedAt time.Time  `json:"reviewed_at"`
}
