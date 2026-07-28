package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BenefitCost struct {
	ID            uuid.UUID       `json:"id"`
	CompanyID     uuid.UUID       `json:"company_id"`
	BenefitID     uuid.UUID       `json:"benefit_id"`
	PlanID        *uuid.UUID      `json:"plan_id,omitempty"`
	CostType      string          `json:"cost_type"`
	EmployeeCost  decimal.Decimal `json:"employee_cost"`
	EmployerCost  decimal.Decimal `json:"employer_cost"`
	TotalCost     decimal.Decimal `json:"total_cost"`
	Currency      string          `json:"currency"`
	Frequency     string          `json:"frequency"`
	BillingCycleDay *int         `json:"billing_cycle_day,omitempty"`
	EffectiveFrom time.Time      `json:"effective_from"`
	EffectiveTo   *time.Time     `json:"effective_to,omitempty"`
	IsActive      bool            `json:"is_active"`
	CreatedBy     uuid.UUID      `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type BenefitCostSchedule struct {
	ID        uuid.UUID       `json:"id"`
	CompanyID uuid.UUID       `json:"company_id"`
	BenefitID uuid.UUID       `json:"benefit_id"`
	CostID    *uuid.UUID      `json:"cost_id,omitempty"`
	ScheduleDate time.Time    `json:"schedule_date"`
	Amount    decimal.Decimal `json:"amount"`
	Currency  string          `json:"currency"`
	Status    string          `json:"status"`
	PaidAt    *time.Time     `json:"paid_at,omitempty"`
	PaymentReference *string `json:"payment_reference,omitempty"`
	Notes     *string        `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type BenefitFlexiblePlan struct {
	ID                  uuid.UUID       `json:"id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	Name                string          `json:"name"`
	Description         *string         `json:"description,omitempty"`
	PlanType            string          `json:"plan_type"`
	AnnualAmount        decimal.Decimal `json:"annual_amount"`
	MonthlyAmount       decimal.Decimal `json:"monthly_amount"`
	Currency            string          `json:"currency"`
	EmployerContribution decimal.Decimal `json:"employer_contribution"`
	EmployeeContribution decimal.Decimal `json:"employee_contribution"`
	ContributionFrequency *string       `json:"contribution_frequency,omitempty"`
	MaxRolloverAmount   decimal.Decimal `json:"max_rollover_amount"`
	RolloverExpiryMonths int            `json:"rollover_expiry_months"`
	AllowReimbursement  bool            `json:"allow_reimbursement"`
	AllowPrepaidCard    bool            `json:"allow_prepaid_card"`
	RequireReceipts     bool            `json:"require_receipts"`
	ReceiptMinAmount    decimal.Decimal `json:"receipt_min_amount"`
	EligibleCategories  []string        `json:"eligible_categories,omitempty"`
	TaxExempt           bool            `json:"tax_exempt"`
	IsActive            bool            `json:"is_active"`
	CreatedBy           uuid.UUID       `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type BenefitFlexibleBudget struct {
	ID                    uuid.UUID       `json:"id"`
	CompanyID             uuid.UUID       `json:"company_id"`
	EmployeeID            uuid.UUID       `json:"employee_id"`
	FlexiblePlanID        uuid.UUID       `json:"flexible_plan_id"`
	FiscalYear            int             `json:"fiscal_year"`
	TotalAmount           decimal.Decimal `json:"total_amount"`
	UsedAmount            decimal.Decimal `json:"used_amount"`
	PendingAmount         decimal.Decimal `json:"pending_amount"`
	RolledOverFromPrevious decimal.Decimal `json:"rolled_over_from_previous"`
	Currency              string          `json:"currency"`
	Status                string          `json:"status"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type EmployeeBenefitWallet struct {
	ID               uuid.UUID       `json:"id"`
	CompanyID        uuid.UUID       `json:"company_id"`
	EmployeeID       uuid.UUID       `json:"employee_id"`
	BenefitID        *uuid.UUID      `json:"benefit_id,omitempty"`
	WalletType       string          `json:"wallet_type"`
	Balance          decimal.Decimal `json:"balance"`
	Currency         string          `json:"currency"`
	LastTransactionAt *time.Time     `json:"last_transaction_at,omitempty"`
	IsActive         bool            `json:"is_active"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type BenefitWalletTransaction struct {
	ID              uuid.UUID       `json:"id"`
	WalletID        uuid.UUID       `json:"wallet_id"`
	TransactionType string          `json:"transaction_type"`
	Amount          decimal.Decimal `json:"amount"`
	BalanceBefore   decimal.Decimal `json:"balance_before"`
	BalanceAfter    decimal.Decimal `json:"balance_after"`
	Currency        string          `json:"currency"`
	ReferenceType   *string         `json:"reference_type,omitempty"`
	ReferenceID     *uuid.UUID      `json:"reference_id,omitempty"`
	Description     *string         `json:"description,omitempty"`
	ReceiptURL      *string         `json:"receipt_url,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type BenefitReimbursement struct {
	ID              uuid.UUID       `json:"id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	EmployeeID      uuid.UUID       `json:"employee_id"`
	BenefitID       *uuid.UUID      `json:"benefit_id,omitempty"`
	FlexiblePlanID  *uuid.UUID      `json:"flexible_plan_id,omitempty"`
	WalletID        *uuid.UUID      `json:"wallet_id,omitempty"`
	RequestID       *uuid.UUID      `json:"request_id,omitempty"`
	Category        string          `json:"category"`
	Description     string          `json:"description"`
	Amount          decimal.Decimal `json:"amount"`
	ApprovedAmount  *decimal.Decimal `json:"approved_amount,omitempty"`
	Currency        string          `json:"currency"`
	ReceiptDate     time.Time       `json:"receipt_date"`
	ExpenseDate     time.Time       `json:"expense_date"`
	Status          string          `json:"status"`
	RejectionReason *string         `json:"rejection_reason,omitempty"`
	PaymentMethod   *string         `json:"payment_method,omitempty"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	PaymentReference *string        `json:"payment_reference,omitempty"`
	SubmittedBy     *uuid.UUID      `json:"submitted_by,omitempty"`
	ReviewedBy      *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type BenefitReimbursementDocument struct {
	ID              uuid.UUID  `json:"id"`
	ReimbursementID uuid.UUID  `json:"reimbursement_id"`
	FileName        string     `json:"file_name"`
	FileType        string     `json:"file_type"`
	FileSize        int        `json:"file_size"`
	StoragePath     string     `json:"storage_path"`
	OCRText         *string    `json:"ocr_text,omitempty"`
	IsVerified      bool       `json:"is_verified"`
	UploadedBy      uuid.UUID  `json:"uploaded_by"`
	UploadedAt      time.Time  `json:"uploaded_at"`
}
