package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseCategory struct {
	ID             uuid.UUID  `json:"id"`
	CompanyID      uuid.UUID  `json:"company_id"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	ParentID       *uuid.UUID `json:"parent_id,omitempty"`
	RequiresReceipt bool      `json:"requires_receipt"`
	IsActive       bool       `json:"is_active"`
	SortOrder      int        `json:"sort_order"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ExpensePaymentMethod struct {
	ID             uuid.UUID `json:"id"`
	CompanyID      uuid.UUID `json:"company_id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	IsCorporate    bool      `json:"is_corporate"`
	RequiresReceipt bool     `json:"requires_receipt"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type ExchangeRate struct {
	ID            uuid.UUID       `json:"id"`
	CompanyID     uuid.UUID       `json:"company_id"`
	FromCurrency  string          `json:"from_currency"`
	ToCurrency    string          `json:"to_currency"`
	Rate          decimal.Decimal `json:"rate"`
	EffectiveDate time.Time       `json:"effective_date"`
	Source        string          `json:"source"`
	CreatedAt     time.Time       `json:"created_at"`
}

type Expense struct {
	ID                 uuid.UUID       `json:"id"`
	CompanyID          uuid.UUID       `json:"company_id"`
	EmployeeID         uuid.UUID       `json:"employee_id"`
	TravelID           *uuid.UUID      `json:"travel_id,omitempty"`
	ExpenseReportID    *uuid.UUID      `json:"expense_report_id,omitempty"`
	CategoryID         uuid.UUID       `json:"category_id"`
	ExpenseDate        time.Time       `json:"expense_date"`
	Description        string          `json:"description"`
	OriginalAmount     decimal.Decimal `json:"original_amount"`
	OriginalCurrency   string          `json:"original_currency"`
	ExchangeRate       decimal.Decimal `json:"exchange_rate"`
	BaseAmount         decimal.Decimal `json:"base_amount"`
	BaseCurrency       string          `json:"base_currency"`
	TaxAmount          decimal.Decimal `json:"tax_amount"`
	TotalAmount        decimal.Decimal `json:"total_amount"`
	PaymentMethodID    *uuid.UUID      `json:"payment_method_id,omitempty"`
	PaymentMethodOther *string         `json:"payment_method_other,omitempty"`
	MerchantName       *string         `json:"merchant_name,omitempty"`
	MerchantTaxID      *string         `json:"merchant_tax_id,omitempty"`
	ReceiptNumber      *string         `json:"receipt_number,omitempty"`
	IsReimbursable     bool            `json:"is_reimbursable"`
	IsPolicyCompliant  bool            `json:"is_policy_compliant"`
	PolicyStatus       string          `json:"policy_status"`
	PolicyOverrideReason *string       `json:"policy_override_reason,omitempty"`
	PolicyOverrideBy   *uuid.UUID      `json:"policy_override_by,omitempty"`
	Status             string          `json:"status"`
	RejectionReason    *string         `json:"rejection_reason,omitempty"`
	Observation        *string         `json:"observation,omitempty"`
	CostCenterID       *uuid.UUID      `json:"cost_center_id,omitempty"`
	ProjectID          *uuid.UUID      `json:"project_id,omitempty"`
	ReceiptRequired      bool            `json:"receipt_required"`
	IsBillable           bool            `json:"is_billable"`
	BillableClient       *string         `json:"billable_client,omitempty"`
	ReceiptMissingReason *string         `json:"receipt_missing_reason,omitempty"`
	IdempotencyKey     *string         `json:"idempotency_key,omitempty"`
	CreatedBy          uuid.UUID       `json:"created_by"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
