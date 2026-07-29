package domain

import "time"

type PlanStatus string

const (
	PlanStatusDraft     PlanStatus = "DRAFT"
	PlanStatusActive    PlanStatus = "ACTIVE"
	PlanStatusCompleted PlanStatus = "COMPLETED"
	PlanStatusCancelled PlanStatus = "CANCELLED"
	PlanStatusSuccess   PlanStatus = "SUCCESS"
	PlanStatusFailed    PlanStatus = "FAILED"
)

// Improvement Plans
type ImprovementPlan struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	CycleID         *string    `json:"cycle_id,omitempty"`
	CreatedBy       string     `json:"created_by"`
	Reason          string     `json:"reason"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Status          PlanStatus `json:"status"`
	SuccessCriteria *string    `json:"success_criteria,omitempty"`
	FinalResult     *string    `json:"final_result,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Actions []ImprovementPlanAction `json:"actions,omitempty"`
}

type ImprovementPlanAction struct {
	ID           string     `json:"id"`
	PlanID       string     `json:"plan_id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	ResponsibleID *string   `json:"responsible_id,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Status       PlanStatus `json:"status"`
	Progress     float64    `json:"progress"`
	Evidence     *string    `json:"evidence,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// Development Plans
type DevelopmentPlan struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	CycleID      *string    `json:"cycle_id,omitempty"`
	CreatedBy    string     `json:"created_by"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	CareerGoal   *string    `json:"career_goal,omitempty"`
	CurrentLevel *int       `json:"current_level,omitempty"`
	TargetLevel  *int       `json:"target_level,omitempty"`
	CompetencyID *string    `json:"competency_id,omitempty"`
	Status       PlanStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Actions []DevelopmentPlanAction `json:"actions,omitempty"`
}

type DevelopmentPlanAction struct {
	ID          string     `json:"id"`
	PlanID      string     `json:"plan_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	ActionType  string     `json:"action_type"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      PlanStatus `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type PlanFilter struct {
	CompanyID  string
	EmployeeID string
	CycleID    string
	Status     PlanStatus
	Search     string
}
