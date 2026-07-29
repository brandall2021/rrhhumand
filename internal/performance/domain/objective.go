package domain

import "time"

type ObjectiveStatus string

const (
	ObjectiveStatusNotStarted ObjectiveStatus = "NOT_STARTED"
	ObjectiveStatusInProgress ObjectiveStatus = "IN_PROGRESS"
	ObjectiveStatusCompleted  ObjectiveStatus = "COMPLETED"
	ObjectiveStatusCancelled  ObjectiveStatus = "CANCELLED"
	ObjectiveStatusAtRisk     ObjectiveStatus = "AT_RISK"
	ObjectiveStatusOnTrack    ObjectiveStatus = "ON_TRACK"
)

type ObjectiveType string

const (
	ObjectiveTypeIndividual ObjectiveType = "INDIVIDUAL"
	ObjectiveTypeTeam       ObjectiveType = "TEAM"
	ObjectiveTypeDepartment ObjectiveType = "DEPARTMENT"
	ObjectiveTypeCompany    ObjectiveType = "COMPANY"
)

type ProgressType string

const (
	ProgressTypePercentage ProgressType = "PERCENTAGE"
	ProgressTypeNumeric    ProgressType = "NUMERIC"
	ProgressTypeBoolean    ProgressType = "BOOLEAN"
)

type PerformanceObjective struct {
	ID               string          `json:"id"`
	CompanyID        string          `json:"company_id"`
	CycleID          string          `json:"cycle_id"`
	EmployeeID       string          `json:"employee_id"`
	ParentObjectiveID *string        `json:"parent_objective_id,omitempty"`
	Title            string          `json:"title"`
	Description      *string         `json:"description,omitempty"`
	ObjectiveType    ObjectiveType   `json:"objective_type"`
	Weight           float64         `json:"weight"`
	StartDate        *time.Time      `json:"start_date,omitempty"`
	DueDate          *time.Time      `json:"due_date,omitempty"`
	Status           ObjectiveStatus `json:"status"`
	Progress         float64         `json:"progress"`
	TargetValue      *float64        `json:"target_value,omitempty"`
	CurrentValue     *float64        `json:"current_value,omitempty"`
	Unit             *string         `json:"unit,omitempty"`
	ProgressType     ProgressType    `json:"progress_type"`
	Notes            *string         `json:"notes,omitempty"`
	RiskNotes        *string         `json:"risk_notes,omitempty"`
	CreatedBy        string          `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`

	KeyResults []ObjectiveKeyResult `json:"key_results,omitempty"`
}

type ObjectiveKeyResult struct {
	ID          string          `json:"id"`
	ObjectiveID string          `json:"objective_id"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Weight      float64         `json:"weight"`
	TargetValue *float64        `json:"target_value,omitempty"`
	CurrentValue float64        `json:"current_value"`
	Unit        *string         `json:"unit,omitempty"`
	Progress    float64         `json:"progress"`
	Status      ObjectiveStatus `json:"status"`
	SortOrder   int             `json:"sort_order"`
}

type ObjectiveFilter struct {
	CompanyID  string
	CycleID    string
	EmployeeID string
	Status     ObjectiveStatus
	Search     string
}
