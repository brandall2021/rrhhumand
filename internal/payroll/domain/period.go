package domain

import "time"

type PayrollPeriod struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	Year        int        `json:"year"`
	Month       int        `json:"month"`
	PeriodType  string     `json:"period_type"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
	Status      string     `json:"status"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PayrollRun struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	PeriodID      string     `json:"period_id"`
	RunNumber     int        `json:"run_number"`
	RunType       string     `json:"run_type"`
	Status        string     `json:"status"`
	EngineVersion string     `json:"engine_version"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	CreatedBy     string     `json:"created_by"`
	ApprovedBy    *string    `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ClosedBy      *string    `json:"closed_by,omitempty"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
