package domain

import "time"

type InterviewStatus string

const (
    IntStatusScheduled  InterviewStatus = "SCHEDULED"
    IntStatusConfirmed  InterviewStatus = "CONFIRMED"
    IntStatusInProgress InterviewStatus = "IN_PROGRESS"
    IntStatusCompleted  InterviewStatus = "COMPLETED"
    IntStatusCancelled  InterviewStatus = "CANCELLED"
    IntStatusNoShow     InterviewStatus = "NO_SHOW"
    IntStatusRescheduled InterviewStatus = "RESCHEDULED"
)

type Interview struct {
    ID              string          `json:"id"`
    CompanyID       string          `json:"company_id"`
    ApplicationID   string          `json:"application_id"`
    InterviewType   string          `json:"interview_type"`
    Title           *string         `json:"title,omitempty"`
    ScheduledAt     *time.Time      `json:"scheduled_at,omitempty"`
    DurationMinutes *int            `json:"duration_minutes,omitempty"`
    MeetingURL      *string         `json:"meeting_url,omitempty"`
    MeetingPassword *string         `json:"meeting_password,omitempty"`
    Location        *string         `json:"location,omitempty"`
    Instructions    *string         `json:"instructions,omitempty"`
    Status          InterviewStatus `json:"status"`
    Score           *float64        `json:"score,omitempty"`
    Notes           *string         `json:"notes,omitempty"`
    CancelledAt     *time.Time      `json:"cancelled_at,omitempty"`
    CancelReason    *string         `json:"cancel_reason,omitempty"`
    CreatedBy       *string         `json:"created_by,omitempty"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    Panel           []InterviewPanelMember `json:"panel,omitempty"`
}

type InterviewPanelMember struct {
    ID          string    `json:"id"`
    InterviewID string    `json:"interview_id"`
    EmployeeID  string    `json:"employee_id"`
    Role        string    `json:"role"`
    Status      string    `json:"status"`
    ResponseAt  *time.Time `json:"response_at,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type InterviewFeedback struct {
    ID                 string                   `json:"id"`
    InterviewID        string                   `json:"interview_id"`
    PanelistID         string                   `json:"panelist_id"`
    Score              *float64                 `json:"score,omitempty"`
    Comments           *string                  `json:"comments,omitempty"`
    Strengths          []string                 `json:"strengths,omitempty"`
    Weaknesses         []string                 `json:"weaknesses,omitempty"`
    Recommendation     *string                  `json:"recommendation,omitempty"`
    SubmittedAt        *time.Time               `json:"submitted_at,omitempty"`
    CreatedAt          time.Time                `json:"created_at"`
    Questions          []InterviewFeedbackQuestion `json:"questions,omitempty"`
}

type InterviewFeedbackQuestion struct {
    ID                string  `json:"id"`
    InterviewFeedbackID string `json:"interview_feedback_id"`
    Question          string  `json:"question"`
    Score             *float64 `json:"score,omitempty"`
    Comment           *string `json:"comment,omitempty"`
}
