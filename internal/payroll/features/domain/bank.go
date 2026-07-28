package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BankBatch struct {
	ID             uuid.UUID       `json:"id"`
	CompanyID      uuid.UUID       `json:"company_id"`
	RunID          uuid.UUID       `json:"run_id"`
	BatchNumber    string          `json:"batch_number"`
	BankCode       string          `json:"bank_code"`
	BankName       *string         `json:"bank_name,omitempty"`
	PaymentType    string          `json:"payment_type"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	TotalEmployees int             `json:"total_employees"`
	Currency       string          `json:"currency"`
	PaymentDate    time.Time       `json:"payment_date"`
	Status         string          `json:"status"`
	FileFormat     string          `json:"file_format"`
	FileName       *string         `json:"file_name,omitempty"`
	StoragePath    *string         `json:"storage_path,omitempty"`
	FileContent    *string         `json:"file_content,omitempty"`
	SentAt         *time.Time      `json:"sent_at,omitempty"`
	ProcessedAt    *time.Time      `json:"processed_at,omitempty"`
	ErrorMessage   *string         `json:"error_message,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	GeneratedBy    uuid.UUID       `json:"generated_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type BankBatchItem struct {
	ID            uuid.UUID       `json:"id"`
	BatchID       uuid.UUID       `json:"batch_id"`
	EmployeeID    uuid.UUID       `json:"employee_id"`
	RunEmployeeID *uuid.UUID      `json:"run_employee_id,omitempty"`
	CUIL          string          `json:"cuil"`
	Surname       string          `json:"surname"`
	Name          string          `json:"name"`
	BankCode      *string         `json:"bank_code,omitempty"`
	BankName      *string         `json:"bank_name,omitempty"`
	BranchCode    *string         `json:"branch_code,omitempty"`
	AccountType   *string         `json:"account_type,omitempty"`
	AccountNumber *string         `json:"account_number,omitempty"`
	CBU           *string         `json:"cbu,omitempty"`
	Alias         *string         `json:"alias,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Concept       *string         `json:"concept,omitempty"`
	Status        string          `json:"status"`
	ErrorMessage  *string         `json:"error_message,omitempty"`
	PaymentDate   *time.Time      `json:"payment_date,omitempty"`
	TransactionID *string         `json:"transaction_id,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
