package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EmployeeBonus struct {
	ID                   uuid.UUID       `json:"id"`
	CompanyID            uuid.UUID       `json:"company_id"`
	EmployeeID           uuid.UUID       `json:"employee_id"`
	BonusType            string          `json:"bonus_type"`
	Name                 string          `json:"name"`
	Description          *string         `json:"description,omitempty"`
	Amount               decimal.Decimal `json:"amount"`
	Currency             string          `json:"currency"`
	PaymentType          string          `json:"payment_type"`
	InstallmentCount     int             `json:"installment_count"`
	InstallmentAmount    *decimal.Decimal `json:"installment_amount,omitempty"`
	Frequency            *string         `json:"frequency,omitempty"`
	GrantDate            time.Time       `json:"grant_date"`
	VestingStart         *time.Time      `json:"vesting_start,omitempty"`
	VestingEnd           *time.Time      `json:"vesting_end,omitempty"`
	PaymentDate          *time.Time      `json:"payment_date,omitempty"`
	Status               string          `json:"status"`
	ClawbackAmount       decimal.Decimal `json:"clawback_amount"`
	ClawbackReason       *string         `json:"clawback_reason,omitempty"`
	PerformancePeriodStart *time.Time    `json:"performance_period_start,omitempty"`
	PerformancePeriodEnd  *time.Time     `json:"performance_period_end,omitempty"`
	PerformanceScore     *decimal.Decimal `json:"performance_score,omitempty"`
	IsTaxable            bool            `json:"is_taxable"`
	TaxWithholding       decimal.Decimal `json:"tax_withholding"`
	NetAmount            *decimal.Decimal `json:"net_amount,omitempty"`
	ApprovedBy           *uuid.UUID      `json:"approved_by,omitempty"`
	ApprovedAt           *time.Time      `json:"approved_at,omitempty"`
	PaidInPayroll        bool            `json:"paid_in_payroll"`
	PayrollRunID         *uuid.UUID      `json:"payroll_run_id,omitempty"`
	Notes                *string         `json:"notes,omitempty"`
	CreatedBy            uuid.UUID       `json:"created_by"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type EmployeeIncentive struct {
	ID             uuid.UUID       `json:"id"`
	CompanyID      uuid.UUID       `json:"company_id"`
	EmployeeID     uuid.UUID       `json:"employee_id"`
	IncentiveType  string          `json:"incentive_type"`
	Name           string          `json:"name"`
	Description    *string         `json:"description,omitempty"`
	Value          decimal.Decimal `json:"value"`
	Currency       string          `json:"currency"`
	AwardDate      time.Time       `json:"award_date"`
	ExpiryDate     *time.Time      `json:"expiry_date,omitempty"`
	RedemptionDate *time.Time      `json:"redemption_date,omitempty"`
	Status         string          `json:"status"`
	PointsCost     int             `json:"points_cost"`
	IsTaxable      bool            `json:"is_taxable"`
	AwardedBy      *uuid.UUID      `json:"awarded_by,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type BenefitPayrollMapping struct {
	ID               uuid.UUID       `json:"id"`
	CompanyID        uuid.UUID       `json:"company_id"`
	BenefitID        *uuid.UUID      `json:"benefit_id,omitempty"`
	FlexiblePlanID   *uuid.UUID      `json:"flexible_plan_id,omitempty"`
	EmployeeBenefitID *uuid.UUID     `json:"employee_benefit_id,omitempty"`
	BonusID          *uuid.UUID      `json:"bonus_id,omitempty"`
	MappingType      string          `json:"mapping_type"`
	PayrollConceptID *uuid.UUID      `json:"payroll_concept_id,omitempty"`
	Amount           decimal.Decimal `json:"amount"`
	Currency         string          `json:"currency"`
	Frequency        string          `json:"frequency"`
	EffectiveFrom    time.Time       `json:"effective_from"`
	EffectiveTo      *time.Time      `json:"effective_to,omitempty"`
	IsActive         bool            `json:"is_active"`
	LastSyncedAt     *time.Time      `json:"last_synced_at,omitempty"`
	SyncStatus       string          `json:"sync_status"`
	SyncError        *string         `json:"sync_error,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
