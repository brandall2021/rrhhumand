package domain

import "time"

type AssessmentStatus string

const (
    AssStatusPending    AssessmentStatus = "PENDING"
    AssStatusSent       AssessmentStatus = "SENT"
    AssStatusInProgress AssessmentStatus = "IN_PROGRESS"
    AssStatusCompleted  AssessmentStatus = "COMPLETED"
    AssStatusExpired    AssessmentStatus = "EXPIRED"
    AssStatusCancelled  AssessmentStatus = "CANCELLED"
)

type Assessment struct {
    ID              string           `json:"id"`
    CompanyID       string           `json:"company_id"`
    ApplicationID   string           `json:"application_id"`
    AssessmentType  string           `json:"assessment_type"`
    Title           string           `json:"title"`
    Description     *string          `json:"description,omitempty"`
    MaxScore        *float64         `json:"max_score,omitempty"`
    PassingScore    *float64         `json:"passing_score,omitempty"`
    DurationMinutes *int             `json:"duration_minutes,omitempty"`
    DueAt           *time.Time       `json:"due_at,omitempty"`
    Status          AssessmentStatus `json:"status"`
    Score           *float64         `json:"score,omitempty"`
    Result          *string          `json:"result,omitempty"`
    ResultSummary   *string          `json:"result_summary,omitempty"`
    CompletedAt     *time.Time       `json:"completed_at,omitempty"`
    CreatedBy       *string          `json:"created_by,omitempty"`
    CreatedAt       time.Time        `json:"created_at"`
    UpdatedAt       time.Time        `json:"updated_at"`
    Sections        []AssessmentSection `json:"sections,omitempty"`
}

type AssessmentSection struct {
    ID           string   `json:"id"`
    AssessmentID string   `json:"assessment_id"`
    Name         string   `json:"name"`
    Description  *string  `json:"description,omitempty"`
    MaxScore     *float64 `json:"max_score,omitempty"`
    Weight       *float64 `json:"weight,omitempty"`
    SortOrder    int      `json:"sort_order"`
}

type AssessmentResult struct {
    ID           string   `json:"id"`
    AssessmentID string   `json:"assessment_id"`
    SectionID    *string  `json:"section_id,omitempty"`
    Score        *float64 `json:"score,omitempty"`
    MaxScore     *float64 `json:"max_score,omitempty"`
    Comment      *string  `json:"comment,omitempty"`
    GradedBy     *string  `json:"graded_by,omitempty"`
    GradedAt     *time.Time `json:"graded_at,omitempty"`
}
