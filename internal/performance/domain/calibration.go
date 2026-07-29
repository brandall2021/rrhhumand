package domain

import "time"

type CalibrationStatus string

const (
	CalibrationDraft     CalibrationStatus = "DRAFT"
	CalibrationInProgress CalibrationStatus = "IN_PROGRESS"
	CalibrationCompleted  CalibrationStatus = "COMPLETED"
	CalibrationCancelled  CalibrationStatus = "CANCELLED"
)

type CalibrationSession struct {
	ID          string            `json:"id"`
	CompanyID   string            `json:"company_id"`
	CycleID     string            `json:"cycle_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Status      CalibrationStatus `json:"status"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	Items []CalibrationItem `json:"items,omitempty"`
}

type CalibrationItem struct {
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id"`
	EmployeeID     string    `json:"employee_id"`
	OriginalScore  *float64  `json:"original_score,omitempty"`
	AdjustedScore  *float64  `json:"adjusted_score,omitempty"`
	OriginalRating *string   `json:"original_rating,omitempty"`
	AdjustedRating *string   `json:"adjusted_rating,omitempty"`
	Reason         *string   `json:"reason,omitempty"`
	ApprovedBy     *string   `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
}
