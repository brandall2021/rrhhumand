package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TotalRewardsItem struct {
	ID                uuid.UUID       `json:"id"`
	CompanyID         uuid.UUID       `json:"company_id"`
	Name              string          `json:"name"`
	Category          string          `json:"category"`
	Description       *string         `json:"description,omitempty"`
	AmountType        string          `json:"amount_type"`
	AmountValue       decimal.Decimal `json:"amount_value"`
	AmountPercentage  *decimal.Decimal `json:"amount_percentage,omitempty"`
	Currency          string          `json:"currency"`
	Frequency         string          `json:"frequency"`
	DisplayOrder      int             `json:"display_order"`
	IsMonetary        bool            `json:"is_monetary"`
	IncludeInStatement bool           `json:"include_in_statement"`
	Icon              *string         `json:"icon,omitempty"`
	Color             *string         `json:"color,omitempty"`
	IsActive          bool            `json:"is_active"`
	CreatedBy         uuid.UUID       `json:"created_by"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type TotalRewardsSnapshot struct {
	ID                  uuid.UUID       `json:"id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	EmployeeID          uuid.UUID       `json:"employee_id"`
	SnapshotDate        time.Time       `json:"snapshot_date"`
	FiscalYear          int             `json:"fiscal_year"`
	PeriodName          *string         `json:"period_name,omitempty"`
	BaseSalary          decimal.Decimal `json:"base_salary"`
	VariablePay         decimal.Decimal `json:"variable_pay"`
	BonusesTotal        decimal.Decimal `json:"bonuses_total"`
	IncentivesTotal     decimal.Decimal `json:"incentives_total"`
	BenefitsTotal       decimal.Decimal `json:"benefits_total"`
	EmployerContributions decimal.Decimal `json:"employer_contributions"`
	FlexibleSpending    decimal.Decimal `json:"flexible_spending"`
	InsuranceValue      decimal.Decimal `json:"insurance_value"`
	DevelopmentValue    decimal.Decimal `json:"development_value"`
	WellnessValue       decimal.Decimal `json:"wellness_value"`
	RecognitionValue    decimal.Decimal `json:"recognition_value"`
	PerksValue          decimal.Decimal `json:"perks_value"`
	TotalRewards        decimal.Decimal `json:"total_rewards"`
	Currency            string          `json:"currency"`
	Items               []map[string]any `json:"items,omitempty"`
	Metadata            map[string]any  `json:"metadata,omitempty"`
	GeneratedBy         uuid.UUID       `json:"generated_by"`
	GeneratedAt         time.Time       `json:"generated_at"`
	CreatedAt           time.Time       `json:"created_at"`
}

type BenefitNotificationLog struct {
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

type BenefitReportDefinition struct {
	ID          uuid.UUID      `json:"id"`
	CompanyID   uuid.UUID      `json:"company_id"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	ReportType  string         `json:"report_type"`
	Config      map[string]any `json:"config,omitempty"`
	ScheduleCron *string       `json:"schedule_cron,omitempty"`
	Recipients  []string       `json:"recipients,omitempty"`
	IsActive    bool           `json:"is_active"`
	CreatedBy   uuid.UUID      `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type BenefitReportResult struct {
	ID          uuid.UUID  `json:"id"`
	DefinitionID *uuid.UUID `json:"definition_id,omitempty"`
	CompanyID   uuid.UUID  `json:"company_id"`
	ReportType  string     `json:"report_type"`
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`
	FileName    *string    `json:"file_name,omitempty"`
	FileContent *string    `json:"file_content,omitempty"`
	StoragePath *string    `json:"storage_path,omitempty"`
	FileSize    *int       `json:"file_size,omitempty"`
	Format      string     `json:"format"`
	Status      string     `json:"status"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	GeneratedBy uuid.UUID  `json:"generated_by"`
	GeneratedAt time.Time  `json:"generated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
