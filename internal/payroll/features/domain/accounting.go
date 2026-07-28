package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountingAccountMapping struct {
	ID                   uuid.UUID  `json:"id"`
	CompanyID            uuid.UUID  `json:"company_id"`
	ConceptID            *uuid.UUID `json:"concept_id,omitempty"`
	MappingType          string     `json:"mapping_type"`
	DebitAccount         string     `json:"debit_account"`
	CreditAccount        string     `json:"credit_account"`
	CostCenterRequired   bool       `json:"cost_center_required"`
	DescriptionTemplate  *string    `json:"description_template,omitempty"`
	Priority             int        `json:"priority"`
	EffectiveFrom        time.Time  `json:"effective_from"`
	EffectiveTo          *time.Time `json:"effective_to,omitempty"`
	IsActive             bool       `json:"is_active"`
	CreatedBy            uuid.UUID  `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type AccountingExport struct {
	ID            uuid.UUID       `json:"id"`
	CompanyID     uuid.UUID       `json:"company_id"`
	RunID         uuid.UUID       `json:"run_id"`
	PeriodID      *uuid.UUID      `json:"period_id,omitempty"`
	ExportType    string          `json:"export_type"`
	FileFormat    string          `json:"file_format"`
	FileName      string          `json:"file_name"`
	FileContent   *string         `json:"file_content,omitempty"`
	StoragePath   *string         `json:"storage_path,omitempty"`
	FileSize      *int            `json:"file_size,omitempty"`
	Status        string          `json:"status"`
	EmployeeCount int             `json:"employee_count"`
	TotalDebit    decimal.Decimal `json:"total_debit"`
	TotalCredit   decimal.Decimal `json:"total_credit"`
	EntryCount    int             `json:"entry_count"`
	ErrorMessage  *string         `json:"error_message,omitempty"`
	GeneratedBy   uuid.UUID       `json:"generated_by"`
	GeneratedAt   time.Time       `json:"generated_at"`
	CreatedAt     time.Time       `json:"created_at"`
}

type AccountingEntry struct {
	ID             uuid.UUID       `json:"id"`
	ExportID       uuid.UUID       `json:"export_id"`
	EntryNumber    int             `json:"entry_number"`
	AccountCode    string          `json:"account_code"`
	AccountName    *string         `json:"account_name,omitempty"`
	CostCenter     *string         `json:"cost_center,omitempty"`
	DebitAmount    decimal.Decimal `json:"debit_amount"`
	CreditAmount   decimal.Decimal `json:"credit_amount"`
	ConceptCode    *string         `json:"concept_code,omitempty"`
	ConceptName    *string         `json:"concept_name,omitempty"`
	EmployeeID     *uuid.UUID      `json:"employee_id,omitempty"`
	EmployeeName   *string         `json:"employee_name,omitempty"`
	DocumentType   *string         `json:"document_type,omitempty"`
	DocumentNumber *string         `json:"document_number,omitempty"`
	Reference      *string         `json:"reference,omitempty"`
	SortOrder      int             `json:"sort_order"`
	CreatedAt      time.Time       `json:"created_at"`
}
