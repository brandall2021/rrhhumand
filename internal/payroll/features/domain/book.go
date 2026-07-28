package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BookEntry struct {
	ID                   uuid.UUID       `json:"id"`
	CompanyID            uuid.UUID       `json:"company_id"`
	RunID                uuid.UUID       `json:"run_id"`
	RunEmployeeID        uuid.UUID       `json:"run_employee_id"`
	EmployeeID           uuid.UUID       `json:"employee_id"`
	EntryType            string          `json:"entry_type"`
	CUIL                 string          `json:"cuil"`
	Surname              string          `json:"surname"`
	Name                 string          `json:"name"`
	Nationality          *string         `json:"nationality,omitempty"`
	BirthDate            *time.Time      `json:"birth_date,omitempty"`
	Sex                  *string         `json:"sex,omitempty"`
	AdmissionDate        time.Time       `json:"admission_date"`
	DischargeDate        *time.Time      `json:"discharge_date,omitempty"`
	CategoryCode         *string         `json:"category_code,omitempty"`
	CategoryName         *string         `json:"category_name,omitempty"`
	AgreementCode        *string         `json:"agreement_code,omitempty"`
	AgreementName        *string         `json:"agreement_name,omitempty"`
	WorkType             *string         `json:"work_type,omitempty"`
	WorkPlace            *string         `json:"work_place,omitempty"`
	GrossRemunerative    decimal.Decimal `json:"gross_remunerative"`
	GrossNonRemunerative decimal.Decimal `json:"gross_non_remunerative"`
	DeductionsTotal      decimal.Decimal `json:"deductions_total"`
	ContributionsTotal   decimal.Decimal `json:"contributions_total"`
	NetAmount            decimal.Decimal `json:"net_amount"`
	EmployerCost         decimal.Decimal `json:"employer_cost"`
	DaysWorked           int             `json:"days_worked"`
	HoursWorked          decimal.Decimal `json:"hours_worked"`
	Absences             int             `json:"absences"`
	Status               string          `json:"status"`
	BookNumber           *int            `json:"book_number,omitempty"`
	PageNumber           *int            `json:"page_number,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type BookExport struct {
	ID                  uuid.UUID       `json:"id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	PeriodID            *uuid.UUID      `json:"period_id,omitempty"`
	Year                int             `json:"year"`
	Month               int             `json:"month"`
	ExportType          string          `json:"export_type"`
	FileName            string          `json:"file_name"`
	FileContent         *string         `json:"file_content,omitempty"`
	StoragePath         *string         `json:"storage_path,omitempty"`
	FileSize            *int            `json:"file_size,omitempty"`
	Status              string          `json:"status"`
	SubmissionDate      *time.Time      `json:"submission_date,omitempty"`
	AcknowledgementCode *string         `json:"acknowledgement_code,omitempty"`
	EmployeeCount       int             `json:"employee_count"`
	TotalGross          decimal.Decimal `json:"total_gross"`
	TotalDeductions     decimal.Decimal `json:"total_deductions"`
	TotalNet            decimal.Decimal `json:"total_net"`
	GeneratedBy         uuid.UUID       `json:"generated_by"`
	GeneratedAt         time.Time       `json:"generated_at"`
	CreatedAt           time.Time       `json:"created_at"`
}
