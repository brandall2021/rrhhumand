package domain

import "time"

type EvaluationStatus string

const (
	EvaluationStatusPending   EvaluationStatus = "PENDING"
	EvaluationStatusDraft     EvaluationStatus = "DRAFT"
	EvaluationStatusSubmitted EvaluationStatus = "SUBMITTED"
	EvaluationStatusApproved  EvaluationStatus = "APPROVED"
	EvaluationStatusLocked    EvaluationStatus = "LOCKED"
	EvaluationStatusReopened  EvaluationStatus = "REOPENED"
)

type EvaluationType string

const (
	EvalTypeSelf     EvaluationType = "SELF"
	EvalTypeManager  EvaluationType = "MANAGER"
	EvalTypePeer     EvaluationType = "PEER"
	EvalTypeSubordinate EvaluationType = "SUBORDINATE"
	EvalTypeHR       EvaluationType = "HR"
	EvalTypeExternal EvaluationType = "EXTERNAL"
	EvalType360      EvaluationType = "360"
)

type PerformanceParticipant struct {
	ID             string         `json:"id"`
	CompanyID      string         `json:"company_id"`
	CycleID        string         `json:"cycle_id"`
	EmployeeID     string         `json:"employee_id"`
	EvaluatorID    string         `json:"evaluator_id"`
	EvaluationType EvaluationType `json:"evaluation_type"`
	Status         EvaluationStatus `json:"status"`
	AssignedAt     time.Time      `json:"assigned_at"`
	StartedAt      *time.Time     `json:"started_at,omitempty"`
	SubmittedAt    *time.Time     `json:"submitted_at,omitempty"`
}

type PerformanceEvaluation struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	CycleID        string          `json:"cycle_id"`
	EmployeeID     string          `json:"employee_id"`
	EvaluatorID    string          `json:"evaluator_id"`
	EvaluationType EvaluationType  `json:"evaluation_type"`
	TemplateID     *string         `json:"template_id,omitempty"`
	Status         EvaluationStatus `json:"status"`
	OverallScore   *float64        `json:"overall_score,omitempty"`
	Strengths      *string         `json:"strengths,omitempty"`
	ImprovementAreas *string       `json:"improvement_areas,omitempty"`
	Summary        *string         `json:"summary,omitempty"`
	SubmittedAt    *time.Time      `json:"submitted_at,omitempty"`
	ReviewedAt     *time.Time      `json:"reviewed_at,omitempty"`
	LockedAt       *time.Time      `json:"locked_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	Answers []EvaluationAnswer `json:"answers,omitempty"`
}

type EvaluationAnswer struct {
	ID           string   `json:"id"`
	EvaluationID string   `json:"evaluation_id"`
	QuestionID   *string  `json:"question_id,omitempty"`
	NumericValue *float64 `json:"numeric_value,omitempty"`
	TextValue    *string  `json:"text_value,omitempty"`
	SelectedValue *string `json:"selected_value,omitempty"`
	BooleanValue *bool    `json:"boolean_value,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ObjectiveEvaluation struct {
	ObjectiveID  string   `json:"objective_id"`
	EvaluationID string   `json:"evaluation_id"`
	Score        *float64 `json:"score,omitempty"`
	Comment      *string  `json:"comment,omitempty"`
}

type CompetencyEvaluation struct {
	CompetencyID string   `json:"competency_id"`
	EvaluationID string   `json:"evaluation_id"`
	Score        *float64 `json:"score,omitempty"`
	ExpectedLevel *int    `json:"expected_level,omitempty"`
	Comment      *string  `json:"comment,omitempty"`
}

type PerformanceReview struct {
	ID                string          `json:"id"`
	CompanyID         string          `json:"company_id"`
	CycleID           string          `json:"cycle_id"`
	EmployeeID        string          `json:"employee_id"`
	ManagerID         string          `json:"manager_id"`
	Summary           *string         `json:"summary,omitempty"`
	Strengths         *string         `json:"strengths,omitempty"`
	ImprovementAreas  *string         `json:"improvement_areas,omitempty"`
	FinalScore        *float64        `json:"final_score,omitempty"`
	FinalRating       *string         `json:"final_rating,omitempty"`
	EmployeeComments  *string         `json:"employee_comments,omitempty"`
	ManagerComments   *string         `json:"manager_comments,omitempty"`
	EmployeeAgreement string          `json:"employee_agreement"`
	DisagreementReason *string        `json:"disagreement_reason,omitempty"`
	Status            EvaluationStatus `json:"status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type EvaluationFilter struct {
	CompanyID    string
	CycleID      string
	EmployeeID   string
	EvaluatorID  string
	Status       EvaluationStatus
	EvaluationType EvaluationType
	Search       string
}
