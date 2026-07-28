package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ArcaConceptMapping struct {
	ID              uuid.UUID       `json:"id"`
	CompanyID       uuid.UUID       `json:"company_id"`
	ConceptID       uuid.UUID       `json:"concept_id"`
	ArcaConceptCode string          `json:"arca_concept_code"`
	ArcaConceptName *string         `json:"arca_concept_name,omitempty"`
	MappingType     string          `json:"mapping_type"`
	Percentage      *decimal.Decimal `json:"percentage,omitempty"`
	IsTaxable       bool            `json:"is_taxable"`
	IsContributable bool            `json:"is_contributable"`
	Notes           *string         `json:"notes,omitempty"`
	EffectiveFrom   time.Time       `json:"effective_from"`
	EffectiveTo     *time.Time      `json:"effective_to,omitempty"`
	IsActive        bool            `json:"is_active"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ArcaExport struct {
	ID                  uuid.UUID       `json:"id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	RunID               uuid.UUID       `json:"run_id"`
	ExportType          string          `json:"export_type"`
	PeriodID            *uuid.UUID      `json:"period_id,omitempty"`
	FileName            string          `json:"file_name"`
	FileContent         *string         `json:"file_content,omitempty"`
	StoragePath         *string         `json:"storage_path,omitempty"`
	FileSize            *int            `json:"file_size,omitempty"`
	Checksum            *string         `json:"checksum,omitempty"`
	Status              string          `json:"status"`
	ErrorMessage        *string         `json:"error_message,omitempty"`
	SubmissionDate      *time.Time      `json:"submission_date,omitempty"`
	AcknowledgementCode *string         `json:"acknowledgement_code,omitempty"`
	EmployeeCount       int             `json:"employee_count"`
	TotalAmount         decimal.Decimal `json:"total_amount"`
	GeneratedBy         uuid.UUID       `json:"generated_by"`
	GeneratedAt         time.Time       `json:"generated_at"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
