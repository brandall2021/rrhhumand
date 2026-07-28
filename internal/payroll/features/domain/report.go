package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReportTemplate struct {
	ID          uuid.UUID       `json:"id"`
	CompanyID   uuid.UUID       `json:"company_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	ReportType  string          `json:"report_type"`
	Config      json.RawMessage `json:"config"`
	IsDefault   bool            `json:"is_default"`
	IsActive    bool            `json:"is_active"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ReportExport struct {
	ID           uuid.UUID       `json:"id"`
	CompanyID    uuid.UUID       `json:"company_id"`
	RunID        *uuid.UUID      `json:"run_id,omitempty"`
	TemplateID   *uuid.UUID      `json:"template_id,omitempty"`
	ReportType   string          `json:"report_type"`
	FileFormat   string          `json:"file_format"`
	FileName     string          `json:"file_name"`
	FileContent  *string         `json:"file_content,omitempty"`
	StoragePath  *string         `json:"storage_path,omitempty"`
	FileSize     *int            `json:"file_size,omitempty"`
	Status       string          `json:"status"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	Config       json.RawMessage `json:"config,omitempty"`
	GeneratedBy  uuid.UUID       `json:"generated_by"`
	GeneratedAt  time.Time       `json:"generated_at"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ReportFilter struct {
	CompanyID    string
	RunID        *string
	PeriodID     *string
	EmployeeID   *string
	DepartmentID *string
	AgreementID  *string
	CategoryID   *string
	ReportType   string
	Format       string
	FromDate     *time.Time
	ToDate       *time.Time
}
