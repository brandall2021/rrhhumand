package domain

import "time"

type FeedbackType string

const (
	FeedbackTypeGeneral      FeedbackType = "GENERAL"
	FeedbackTypePositive     FeedbackType = "POSITIVE"
	FeedbackTypeConstructive FeedbackType = "CONSTRUCTIVE"
	FeedbackTypePraise       FeedbackType = "PRAISE"
	FeedbackTypeSuggestion   FeedbackType = "SUGGESTION"
	FeedbackTypeRecognition  FeedbackType = "RECOGNITION"
)

type FeedbackVisibility string

const (
	VisibilityEmployee  FeedbackVisibility = "EMPLOYEE"
	VisibilityManager   FeedbackVisibility = "MANAGER"
	VisibilityHR        FeedbackVisibility = "HR"
	VisibilityPublic    FeedbackVisibility = "PUBLIC"
)

type PerformanceFeedback struct {
	ID           string             `json:"id"`
	CompanyID    string             `json:"company_id"`
	CycleID      *string            `json:"cycle_id,omitempty"`
	EmployeeID   string             `json:"employee_id"`
	AuthorID     string             `json:"author_id"`
	FeedbackType FeedbackType       `json:"feedback_type"`
	Visibility   FeedbackVisibility `json:"visibility"`
	Content      string             `json:"content"`
	IsAnonymous  bool               `json:"is_anonymous"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type RecognitionType string

const (
	RecognitionTypeValues     RecognitionType = "VALUES"
	RecognitionTypeAchievement RecognitionType = "ACHIEVEMENT"
	RecognitionTypeInnovation RecognitionType = "INNOVATION"
	RecognitionTypeTeamwork   RecognitionType = "TEAMWORK"
	RecognitionTypeLeadership RecognitionType = "LEADERSHIP"
)

type PerformanceRecognition struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	EmployeeID      string          `json:"employee_id"`
	AuthorID        string          `json:"author_id"`
	RecognitionType RecognitionType `json:"recognition_type"`
	Message         string          `json:"message"`
	CreatedAt       time.Time       `json:"created_at"`
}

type CheckInStatus string

const (
	CheckInPending    CheckInStatus = "PENDING"
	CheckInCompleted  CheckInStatus = "COMPLETED"
	CheckInCancelled  CheckInStatus = "CANCELLED"
)

type PerformanceCheckIn struct {
	ID           string        `json:"id"`
	CompanyID    string        `json:"company_id"`
	EmployeeID   string        `json:"employee_id"`
	ManagerID    string        `json:"manager_id"`
	CycleID      *string       `json:"cycle_id,omitempty"`
	ScheduledAt  time.Time     `json:"scheduled_at"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	EmployeeNotes *string      `json:"employee_notes,omitempty"`
	ManagerNotes *string       `json:"manager_notes,omitempty"`
	Achievements *string       `json:"achievements,omitempty"`
	Blockers     *string       `json:"blockers,omitempty"`
	NextSteps    *string       `json:"next_steps,omitempty"`
	Status       CheckInStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type FeedbackFilter struct {
	CompanyID    string
	EmployeeID   string
	AuthorID     string
	FeedbackType FeedbackType
	Visibility   FeedbackVisibility
	FromDate     *time.Time
	ToDate       *time.Time
}
