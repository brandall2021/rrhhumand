package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BenefitCategory struct {
	ID          uuid.UUID  `json:"id"`
	CompanyID   uuid.UUID  `json:"company_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Icon        *string    `json:"icon,omitempty"`
	Color       *string    `json:"color,omitempty"`
	SortOrder   int        `json:"sort_order"`
	IsActive    bool       `json:"is_active"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type BenefitType struct {
	ID               uuid.UUID       `json:"id"`
	CompanyID        uuid.UUID       `json:"company_id"`
	CategoryID       *uuid.UUID      `json:"category_id,omitempty"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	Code             string          `json:"code"`
	Nature           string          `json:"nature"`
	TaxTreatment     string          `json:"tax_treatment"`
	RequiresApproval bool            `json:"requires_approval"`
	IsReimbursable   bool            `json:"is_reimbursable"`
	IsFlexible       bool            `json:"is_flexible"`
	HasWallet        bool            `json:"has_wallet"`
	SortOrder        int             `json:"sort_order"`
	IsActive         bool            `json:"is_active"`
	ConfigSchema     map[string]any  `json:"config_schema,omitempty"`
	CreatedBy        uuid.UUID       `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type BenefitProvider struct {
	ID              uuid.UUID  `json:"id"`
	CompanyID       uuid.UUID  `json:"company_id"`
	Name            string     `json:"name"`
	LegalName       *string    `json:"legal_name,omitempty"`
	TaxID           *string    `json:"tax_id,omitempty"`
	ProviderType    string     `json:"provider_type"`
	ContactName     *string    `json:"contact_name,omitempty"`
	ContactEmail    *string    `json:"contact_email,omitempty"`
	ContactPhone    *string    `json:"contact_phone,omitempty"`
	Website         *string    `json:"website,omitempty"`
	Address         *string    `json:"address,omitempty"`
	ServiceRegion   *string    `json:"service_region,omitempty"`
	ContractStart   *time.Time `json:"contract_start,omitempty"`
	ContractEnd     *time.Time `json:"contract_end,omitempty"`
	ContractFilePath *string   `json:"contract_file_path,omitempty"`
	BillingCycle    *string    `json:"billing_cycle,omitempty"`
	BillingCurrency string     `json:"billing_currency"`
	Rating          *float64   `json:"rating,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Benefit struct {
	ID                    uuid.UUID       `json:"id"`
	CompanyID             uuid.UUID       `json:"company_id"`
	TypeID                uuid.UUID       `json:"type_id"`
	ProviderID            *uuid.UUID      `json:"provider_id,omitempty"`
	Code                  string          `json:"code"`
	Name                  string          `json:"name"`
	Description           *string         `json:"description,omitempty"`
	ShortDescription      *string         `json:"short_description,omitempty"`
	CoverageDetails       *string         `json:"coverage_details,omitempty"`
	EligibilitySummary    *string         `json:"eligibility_summary,omitempty"`
	LogoURL               *string         `json:"logo_url,omitempty"`
	BannerURL             *string         `json:"banner_url,omitempty"`
	WebsiteURL            *string         `json:"website_url,omitempty"`
	TermsURL              *string         `json:"terms_url,omitempty"`
	DocumentationURL      *string         `json:"documentation_url,omitempty"`
	ProviderReference     *string         `json:"provider_reference,omitempty"`
	StartDate             *time.Time      `json:"start_date,omitempty"`
	EndDate               *time.Time      `json:"end_date,omitempty"`
	MaxBeneficiaries      *int            `json:"max_beneficiaries,omitempty"`
	CurrentBeneficiaries  int             `json:"current_beneficiaries"`
	WaitingPeriodDays     int             `json:"waiting_period_days"`
	MinimumServiceMonths  *int            `json:"minimum_service_months,omitempty"`
	DeductibleAmount      *decimal.Decimal `json:"deductible_amount,omitempty"`
	DeductiblePeriod      *string         `json:"deductible_period,omitempty"`
	CopayPercentage       *decimal.Decimal `json:"copay_percentage,omitempty"`
	MaxCoverageAmount     *decimal.Decimal `json:"max_coverage_amount,omitempty"`
	MaxCoveragePeriod     *string         `json:"max_coverage_period,omitempty"`
	AutoEnroll            bool            `json:"auto_enroll"`
	EnrollmentDeadlineDays *int           `json:"enrollment_deadline_days,omitempty"`
	RequiresEvidence      bool            `json:"requires_evidence"`
	EvidenceDescription   *string         `json:"evidence_description,omitempty"`
	Status                string          `json:"status"`
	Visibility            string          `json:"visibility"`
	SortOrder             int             `json:"sort_order"`
	Metadata              map[string]any  `json:"metadata,omitempty"`
	CreatedBy             uuid.UUID       `json:"created_by"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}
