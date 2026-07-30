package domain

import "time"

type OfferStatus string

const (
    OfferStatusDraft     OfferStatus = "DRAFT"
    OfferStatusPendingApproval OfferStatus = "PENDING_APPROVAL"
    OfferStatusApproved  OfferStatus = "APPROVED"
    OfferStatusSent      OfferStatus = "SENT"
    OfferStatusAccepted  OfferStatus = "ACCEPTED"
    OfferStatusRejected  OfferStatus = "REJECTED"
    OfferStatusExpired   OfferStatus = "EXPIRED"
    OfferStatusWithdrawn OfferStatus = "WITHDRAWN"
    OfferStatusNegotiating OfferStatus = "NEGOTIATING"
)

type Offer struct {
    ID                 string      `json:"id"`
    CompanyID          string      `json:"company_id"`
    ApplicationID      string      `json:"application_id"`
    PositionTitle      string      `json:"position_title"`
    DepartmentID       *string     `json:"department_id,omitempty"`
    OfferType          string      `json:"offer_type"`
    StartDate          *time.Time  `json:"start_date,omitempty"`
    EmploymentType     *string     `json:"employment_type,omitempty"`
    WorkMode           *string     `json:"work_mode,omitempty"`
    SalaryAmount       *float64    `json:"salary_amount,omitempty"`
    SalaryCurrency     *string     `json:"salary_currency,omitempty"`
    SalaryPeriod       *string     `json:"salary_period,omitempty"`
    VariableComp       *string     `json:"variable_compensation,omitempty"`
    BenefitsSummary    *string     `json:"benefits_summary,omitempty"`
    EquityTerms        *string     `json:"equity_terms,omitempty"`
    Conditions         *string     `json:"conditions,omitempty"`
    Notes              *string     `json:"notes,omitempty"`
    ResponseDeadline   *time.Time  `json:"response_deadline,omitempty"`
    Status             OfferStatus `json:"status"`
    SentAt             *time.Time  `json:"sent_at,omitempty"`
    AcceptedAt         *time.Time  `json:"accepted_at,omitempty"`
    RejectedAt         *time.Time  `json:"rejected_at,omitempty"`
    RejectionReason    *string     `json:"rejection_reason,omitempty"`
    ExpiredAt          *time.Time  `json:"expired_at,omitempty"`
    CreatedBy          *string     `json:"created_by,omitempty"`
    CreatedAt          time.Time   `json:"created_at"`
    UpdatedAt          time.Time   `json:"updated_at"`
}

type OfferApproval struct {
    ID        string     `json:"id"`
    OfferID   string     `json:"offer_id"`
    ApproverID string    `json:"approver_id"`
    StepOrder int        `json:"step_order"`
    Status    string     `json:"status"`
    Comment   *string    `json:"comment,omitempty"`
    DecidedAt *time.Time `json:"decided_at,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
}

type OfferNegotiation struct {
    ID             string     `json:"id"`
    OfferID        string     `json:"offer_id"`
    RequestedBy    string     `json:"requested_by"`
    Field          string     `json:"field"`
    OriginalValue  *string    `json:"original_value,omitempty"`
    RequestedValue string     `json:"requested_value"`
    CounterValue   *string    `json:"counter_value,omitempty"`
    Status         string     `json:"status"`
    Notes          *string    `json:"notes,omitempty"`
    CreatedAt      time.Time  `json:"created_at"`
    ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}

type OfferDocument struct {
    ID           string     `json:"id"`
    OfferID      string     `json:"offer_id"`
    DocumentType string     `json:"document_type"`
    FileName     string     `json:"file_name"`
    StorageKey   string     `json:"storage_key"`
    SignedAt     *time.Time `json:"signed_at,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
}
