package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Travel struct {
	ID              uuid.UUID       `json:"id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	EmployeeID      uuid.UUID       `json:"employee_id"`
	Title           string          `json:"title"`
	Purpose         *string         `json:"purpose,omitempty"`
	Origin          string          `json:"origin"`
	Destination     string          `json:"destination"`
	DepartureDate   time.Time       `json:"departure_date"`
	ReturnDate      time.Time       `json:"return_date"`
	Status          string          `json:"status"`
	EstimatedBudget *decimal.Decimal `json:"estimated_budget,omitempty"`
	Currency        string          `json:"currency"`
	CostCenterID    *uuid.UUID      `json:"cost_center_id,omitempty"`
	ProjectID       *uuid.UUID      `json:"project_id,omitempty"`
	RejectionReason *string         `json:"rejection_reason,omitempty"`
	Notes           *string         `json:"notes,omitempty"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type TravelParticipant struct {
	ID         uuid.UUID `json:"id"`
	TravelID   uuid.UUID `json:"travel_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	Role       string    `json:"role"`
	Notes      *string   `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ExpenseReport struct {
	ID                 uuid.UUID       `json:"id"`
	CompanyID          uuid.UUID       `json:"company_id"`
	EmployeeID         uuid.UUID       `json:"employee_id"`
	TravelID           *uuid.UUID      `json:"travel_id,omitempty"`
	AdvanceID          *uuid.UUID      `json:"advance_id,omitempty"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	TotalAmount        decimal.Decimal `json:"total_amount"`
	AdvanceAmount      decimal.Decimal `json:"advance_amount"`
	ReimbursableAmount decimal.Decimal `json:"reimbursable_amount"`
	EmployeeRefundAmount decimal.Decimal `json:"employee_refund_amount"`
	CompanyOwesAmount  decimal.Decimal `json:"company_owes_amount"`
	Currency           string          `json:"currency"`
	Status             string          `json:"status"`
	SubmittedAt        *time.Time      `json:"submitted_at,omitempty"`
	ApprovedAt         *time.Time      `json:"approved_at,omitempty"`
	PaidAt             *time.Time      `json:"paid_at,omitempty"`
	RejectedAt         *time.Time      `json:"rejected_at,omitempty"`
	RejectionReason    *string         `json:"rejection_reason,omitempty"`
	Observation        *string         `json:"observation,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type ExpenseAdvance struct {
	ID              uuid.UUID       `json:"id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	EmployeeID      uuid.UUID       `json:"employee_id"`
	TravelID        *uuid.UUID      `json:"travel_id,omitempty"`
	RequestedAmount decimal.Decimal `json:"requested_amount"`
	ApprovedAmount  *decimal.Decimal `json:"approved_amount,omitempty"`
	SettledAmount   decimal.Decimal `json:"settled_amount"`
	Currency        string          `json:"currency"`
	RequestDate     time.Time       `json:"request_date"`
	ApprovedDate    *time.Time      `json:"approved_date,omitempty"`
	PaidDate        *time.Time      `json:"paid_date,omitempty"`
	SettledDate     *time.Time      `json:"settled_date,omitempty"`
	Status          string          `json:"status"`
	RejectionReason *string         `json:"rejection_reason,omitempty"`
	IdempotencyKey  *string         `json:"idempotency_key,omitempty"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ExpenseReimbursement struct {
	ID              uuid.UUID       `json:"id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	EmployeeID      uuid.UUID       `json:"employee_id"`
	ExpenseReportID *uuid.UUID      `json:"expense_report_id,omitempty"`
	AdvanceID       *uuid.UUID      `json:"advance_id,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	PaymentMethod   *string         `json:"payment_method,omitempty"`
	Status          string          `json:"status"`
	PayrollID       *uuid.UUID      `json:"payroll_id,omitempty"`
	PayrollRunID    *uuid.UUID      `json:"payroll_run_id,omitempty"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	RejectionReason *string         `json:"rejection_reason,omitempty"`
	Notes           *string         `json:"notes,omitempty"`
	IdempotencyKey  *string         `json:"idempotency_key,omitempty"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ExpenseReceipt struct {
	ID           uuid.UUID      `json:"id"`
	CompanyID    uuid.UUID      `json:"company_id"`
	ExpenseID    uuid.UUID      `json:"expense_id"`
	StorageKey   string         `json:"storage_key"`
	Filename     string         `json:"filename"`
	MimeType     string         `json:"mime_type"`
	Size         int64          `json:"size"`
	Hash         *string        `json:"hash,omitempty"`
	OCRText      *string        `json:"ocr_text,omitempty"`
	OCRProcessed bool           `json:"ocr_processed"`
	OCRData      map[string]any `json:"ocr_data,omitempty"`
	UploadedBy   uuid.UUID      `json:"uploaded_by"`
	UploadedAt   time.Time      `json:"uploaded_at"`
}

type ExpenseDuplicateCheck struct {
	ID                uuid.UUID  `json:"id"`
	CompanyID         uuid.UUID  `json:"company_id"`
	ExpenseID         uuid.UUID  `json:"expense_id"`
	DuplicateExpenseID *uuid.UUID `json:"duplicate_expense_id,omitempty"`
	MatchReason       string     `json:"match_reason"`
	MatchScore        *float64   `json:"match_score,omitempty"`
	Status            string     `json:"status"`
	ResolvedBy        *uuid.UUID `json:"resolved_by,omitempty"`
	ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}
