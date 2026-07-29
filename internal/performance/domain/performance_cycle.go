package domain

import "time"

type CycleStatus string

const (
	CycleStatusDraft      CycleStatus = "DRAFT"
	CycleStatusOpen       CycleStatus = "OPEN"
	CycleStatusInProgress CycleStatus = "IN_PROGRESS"
	CycleStatusReview     CycleStatus = "REVIEW"
	CycleStatusCalibration CycleStatus = "CALIBRATION"
	CycleStatusClosed     CycleStatus = "CLOSED"
	CycleStatusCancelled  CycleStatus = "CANCELLED"
)

type CycleType string

const (
	CycleTypeAnnual    CycleType = "ANNUAL"
	CycleTypeSemiannual CycleType = "SEMIANNUAL"
	CycleTypeQuarterly CycleType = "QUARTERLY"
	CycleTypeMonthly   CycleType = "MONTHLY"
	CycleTypeCustom    CycleType = "CUSTOM"
)

type PerformanceCycle struct {
	ID                    string     `json:"id"`
	CompanyID             string     `json:"company_id"`
	Name                  string     `json:"name"`
	Description           *string    `json:"description,omitempty"`
	CycleType             CycleType  `json:"cycle_type"`
	Status                CycleStatus `json:"status"`
	StartDate             *time.Time `json:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	EvaluationStartDate   *time.Time `json:"evaluation_start_date,omitempty"`
	EvaluationEndDate     *time.Time `json:"evaluation_end_date,omitempty"`
	ReviewStartDate       *time.Time `json:"review_start_date,omitempty"`
	ReviewEndDate         *time.Time `json:"review_end_date,omitempty"`
	CalibrationStartDate  *time.Time `json:"calibration_start_date,omitempty"`
	CalibrationEndDate    *time.Time `json:"calibration_end_date,omitempty"`
	TemplateID            *string    `json:"template_id,omitempty"`
	ObjectiveWeight       float64    `json:"objective_weight"`
	CompetencyWeight      float64    `json:"competency_weight"`
	MinAnonymousResponses int        `json:"min_anonymous_responses"`
	CreatedBy             string     `json:"created_by"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type PerformanceCycleFilter struct {
	CompanyID string
	Status    CycleStatus
	Type      CycleType
	Search    string
}
