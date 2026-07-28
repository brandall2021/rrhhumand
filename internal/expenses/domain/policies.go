package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpensePolicy struct {
	ID           uuid.UUID  `json:"id"`
	CompanyID    uuid.UUID  `json:"company_id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	Version      int        `json:"version"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo  *time.Time `json:"effective_to,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ExpensePolicyRule struct {
	ID                    uuid.UUID       `json:"id"`
	PolicyID              uuid.UUID       `json:"policy_id"`
	CategoryID            *uuid.UUID      `json:"category_id,omitempty"`
	EmployeeCategory      *string         `json:"employee_category,omitempty"`
	MaxAmount             *decimal.Decimal `json:"max_amount,omitempty"`
	Currency              *string         `json:"currency,omitempty"`
	RequiresReceipt       bool            `json:"requires_receipt"`
	RequiresApproval      bool            `json:"requires_approval"`
	AllowedPaymentMethods []string        `json:"allowed_payment_methods,omitempty"`
	DailyAllowanceCategory *string        `json:"daily_allowance_category,omitempty"`
	Conditions            map[string]any  `json:"conditions,omitempty"`
	Priority              int             `json:"priority"`
	IsActive              bool            `json:"is_active"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type ExpenseWorkflow struct {
	ID                   uuid.UUID  `json:"id"`
	CompanyID            uuid.UUID  `json:"company_id"`
	Name                 string     `json:"name"`
	Description          *string    `json:"description,omitempty"`
	WorkflowType         string     `json:"workflow_type"`
	MinAmount            decimal.Decimal `json:"min_amount"`
	MaxAmount            *decimal.Decimal `json:"max_amount,omitempty"`
	RequiresChainApproval bool      `json:"requires_chain_approval"`
	IsActive             bool       `json:"is_active"`
	CreatedBy            uuid.UUID  `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type ExpenseWorkflowStep struct {
	ID               uuid.UUID  `json:"id"`
	WorkflowID       uuid.UUID  `json:"workflow_id"`
	StepOrder        int        `json:"step_order"`
	ApprovalType     string     `json:"approval_type"`
	ApproverRoleID   *uuid.UUID `json:"approver_role_id,omitempty"`
	MaxRejectionCount int       `json:"max_rejection_count"`
	IsRequired       bool       `json:"is_required"`
	CreatedAt        time.Time  `json:"created_at"`
}

type ExpenseApproval struct {
	ID         uuid.UUID  `json:"id"`
	CompanyID  uuid.UUID  `json:"company_id"`
	EntityType string     `json:"entity_type"`
	EntityID   uuid.UUID  `json:"entity_id"`
	StepID     *uuid.UUID `json:"step_id,omitempty"`
	ApproverID uuid.UUID  `json:"approver_id"`
	StepOrder  int        `json:"step_order"`
	Status     string     `json:"status"`
	Comment    *string    `json:"comment,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	RejectedAt *time.Time `json:"rejected_at,omitempty"`
	ObservedAt *time.Time `json:"observed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ExpenseBudget struct {
	ID            uuid.UUID       `json:"id"`
	CompanyID     uuid.UUID       `json:"company_id"`
	Name          string          `json:"name"`
	FiscalYear    int             `json:"fiscal_year"`
	PeriodStart   time.Time       `json:"period_start"`
	PeriodEnd     time.Time       `json:"period_end"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	UsedAmount    decimal.Decimal `json:"used_amount"`
	ReservedAmount decimal.Decimal `json:"reserved_amount"`
	Currency      string          `json:"currency"`
	CostCenterID  *uuid.UUID      `json:"cost_center_id,omitempty"`
	ProjectID     *uuid.UUID      `json:"project_id,omitempty"`
	CategoryID    *uuid.UUID      `json:"category_id,omitempty"`
	IsActive      bool            `json:"is_active"`
	CreatedBy     uuid.UUID      `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type DailyAllowanceRule struct {
	ID                  uuid.UUID       `json:"id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	Name                string          `json:"name"`
	Country             *string         `json:"country,omitempty"`
	Region              *string         `json:"region,omitempty"`
	City                *string         `json:"city,omitempty"`
	EmployeeCategory    *string         `json:"employee_category,omitempty"`
	DailyAmount         decimal.Decimal `json:"daily_amount"`
	Currency            string          `json:"currency"`
	MealPercentage      *decimal.Decimal `json:"meal_percentage,omitempty"`
	LodgingPercentage   *decimal.Decimal `json:"lodging_percentage,omitempty"`
	TransportPercentage *decimal.Decimal `json:"transport_percentage,omitempty"`
	OtherPercentage     *decimal.Decimal `json:"other_percentage,omitempty"`
	EffectiveFrom       time.Time       `json:"effective_from"`
	EffectiveTo         *time.Time      `json:"effective_to,omitempty"`
	IsActive            bool            `json:"is_active"`
	CreatedBy           uuid.UUID       `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type CorporateCard struct {
	ID               uuid.UUID  `json:"id"`
	CompanyID        uuid.UUID  `json:"company_id"`
	EmployeeID       *uuid.UUID `json:"employee_id,omitempty"`
	CardNumberMasked string     `json:"card_number_masked"`
	CardholderName   string     `json:"cardholder_name"`
	Provider         *string    `json:"provider,omitempty"`
	CreditLimit      *decimal.Decimal `json:"credit_limit,omitempty"`
	Currency         string     `json:"currency"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CorporateCardTransaction struct {
	ID              uuid.UUID       `json:"id"`
	CardID          uuid.UUID       `json:"card_id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	ExpenseID       *uuid.UUID      `json:"expense_id,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
	MerchantName    *string         `json:"merchant_name,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Reference       *string         `json:"reference,omitempty"`
	Status          string          `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
}

type ExpenseAuditLog struct {
	ID         uuid.UUID      `json:"id"`
	CompanyID  uuid.UUID      `json:"company_id"`
	UserID     uuid.UUID      `json:"user_id"`
	EmployeeID *uuid.UUID     `json:"employee_id,omitempty"`
	Action     string         `json:"action"`
	EntityType string         `json:"entity_type"`
	EntityID   uuid.UUID      `json:"entity_id"`
	OldValues  map[string]any `json:"old_values,omitempty"`
	NewValues  map[string]any `json:"new_values,omitempty"`
	IP         *string        `json:"ip,omitempty"`
	UserAgent  *string        `json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

type ExpenseNotificationLog struct {
	ID               uuid.UUID      `json:"id"`
	CompanyID        uuid.UUID      `json:"company_id"`
	EmployeeID       *uuid.UUID     `json:"employee_id,omitempty"`
	NotificationType string         `json:"notification_type"`
	Channel          string         `json:"channel"`
	Title            string         `json:"title"`
	Body             *string        `json:"body,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	ReadAt           *time.Time     `json:"read_at,omitempty"`
	SentAt           time.Time      `json:"sent_at"`
	CreatedAt        time.Time      `json:"created_at"`
}
