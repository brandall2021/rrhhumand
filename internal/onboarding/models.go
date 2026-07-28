package onboarding

import (
	"encoding/json"
	"time"
)

type OnboardingTemplate struct {
	ID                  string     `json:"id"`
	CompanyID           string     `json:"company_id"`
	Name                string     `json:"name"`
	Description         *string    `json:"description,omitempty"`
	Status              string     `json:"status"`
	DefaultDurationDays int        `json:"default_duration_days"`
	CreatedBy           string     `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type OnboardingTemplateTask struct {
	ID               string  `json:"id"`
	TemplateID        string  `json:"template_id"`
	Title             string  `json:"title"`
	Description       *string `json:"description,omitempty"`
	Category          string  `json:"category"`
	ResponsibleType   string  `json:"responsible_type"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty"`
	Required          bool    `json:"required"`
	DaysOffset        int     `json:"days_offset"`
	SortOrder         int     `json:"sort_order"`
	EstimatedMinutes  *int    `json:"estimated_minutes,omitempty"`
}

type OnboardingProcess struct {
	ID                   string     `json:"id"`
	CompanyID            string     `json:"company_id"`
	EmployeeID           string     `json:"employee_id"`
	TemplateID           *string    `json:"template_id,omitempty"`
	StartDate            time.Time  `json:"start_date"`
	TargetCompletionDate time.Time  `json:"target_completion_date"`
	Status               string     `json:"status"`
	ProgressPercentage   int        `json:"progress_percentage"`
	CompletionPolicy     string     `json:"completion_policy"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	CancelledAt          *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason   *string    `json:"cancellation_reason,omitempty"`
	CreatedBy            string     `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type OnboardingTask struct {
	ID              string     `json:"id"`
	OnboardingID    string     `json:"onboarding_id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Category        string     `json:"category"`
	ResponsibleType string     `json:"responsible_type"`
	ResponsibleID   *string    `json:"responsible_id,omitempty"`
	DueDate         time.Time  `json:"due_date"`
	Status          string     `json:"status"`
	Required        bool       `json:"required"`
	SortOrder       int        `json:"sort_order"`
	EstimatedMinutes *int      `json:"estimated_minutes,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	BlockedReason   *string    `json:"blocked_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OnboardingDocument struct {
	ID              string     `json:"id"`
	OnboardingID    string     `json:"onboarding_id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	DocumentType    string     `json:"document_type"`
	FileName        *string    `json:"file_name,omitempty"`
	MimeType        *string    `json:"mime_type,omitempty"`
	SizeBytes       *int64     `json:"size_bytes,omitempty"`
	Checksum        *string    `json:"checksum,omitempty"`
	StorageProvider *string    `json:"storage_provider,omitempty"`
	StorageKey      *string    `json:"storage_key,omitempty"`
	Status          string     `json:"status"`
	Required        bool       `json:"required"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	UploadedAt      *time.Time `json:"uploaded_at,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy      *string    `json:"reviewed_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OnboardingAsset struct {
	ID           string     `json:"id"`
	OnboardingID string     `json:"onboarding_id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	AssetType    string     `json:"asset_type"`
	Description  *string    `json:"description,omitempty"`
	SerialNumber *string    `json:"serial_number,omitempty"`
	Status       string     `json:"status"`
	AssignedBy   *string    `json:"assigned_by,omitempty"`
	AssignedAt   *time.Time `json:"assigned_at,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type AccessRequest struct {
	ID           string     `json:"id"`
	OnboardingID string     `json:"onboarding_id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	SystemName   string     `json:"system_name"`
	AccessType   string     `json:"access_type"`
	Status       string     `json:"status"`
	RequestedAt  time.Time  `json:"requested_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ApprovedBy   *string    `json:"approved_by,omitempty"`
	ActivatedAt  *time.Time `json:"activated_at,omitempty"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type OnboardingMilestone struct {
	ID              string     `json:"id"`
	OnboardingID    string     `json:"onboarding_id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	MilestoneType   string     `json:"milestone_type"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	DaysOffset      int        `json:"days_offset"`
	DueDate         time.Time  `json:"due_date"`
	ResponsibleType string     `json:"responsible_type"`
	ResponsibleID   *string    `json:"responsible_id,omitempty"`
	Status          string     `json:"status"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OnboardingFeedback struct {
	ID               string     `json:"id"`
	OnboardingID     string     `json:"onboarding_id"`
	CompanyID        string     `json:"company_id"`
	EmployeeID       string     `json:"employee_id"`
	FeedbackType     string     `json:"feedback_type"`
	SubmittedBy      string     `json:"submitted_by"`
	AdaptationScore  *int       `json:"adaptation_score,omitempty"`
	TeamScore        *int       `json:"team_score,omitempty"`
	KnowledgeScore   *int       `json:"knowledge_score,omitempty"`
	CommunicationScore *int     `json:"communication_score,omitempty"`
	OverallScore     *float64   `json:"overall_score,omitempty"`
	Comments         *string    `json:"comments,omitempty"`
	SubmittedAt      time.Time  `json:"submitted_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type OnboardingBuddy struct {
	ID              string     `json:"id"`
	OnboardingID    string     `json:"onboarding_id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	BuddyEmployeeID string     `json:"buddy_employee_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Status          string     `json:"status"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OnboardingException struct {
	ID          string    `json:"id"`
	OnboardingID string   `json:"onboarding_id"`
	CompanyID   string    `json:"company_id"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	Reason      string    `json:"reason"`
	CreatedBy   string    `json:"created_by"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type TrainingAssignment struct {
	ID               string     `json:"id"`
	OnboardingID     string     `json:"onboarding_id"`
	CompanyID        string     `json:"company_id"`
	EmployeeID       string     `json:"employee_id"`
	CourseName       string     `json:"course_name"`
	Description      *string    `json:"description,omitempty"`
	TrainingType     string     `json:"training_type"`
	Status           string     `json:"status"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	ExternalProvider *string    `json:"external_provider,omitempty"`
	ExternalCourseID *string    `json:"external_course_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type Notification struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	UserID          string     `json:"user_id"`
	Title           string     `json:"title"`
	Body            *string    `json:"body,omitempty"`
	NotificationType string    `json:"notification_type"`
	Channel         string     `json:"channel"`
	ReferenceType   *string    `json:"reference_type,omitempty"`
	ReferenceID     *string    `json:"reference_id,omitempty"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type DomainEvent struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	CompanyID     string          `json:"company_id"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	Payload       json.RawMessage `json:"payload,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
}

type OnboardingAuditLog struct {
	ID         string          `json:"id"`
	CompanyID  string          `json:"company_id"`
	UserID     *string         `json:"user_id,omitempty"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	OldValue   json.RawMessage `json:"old_value,omitempty"`
	NewValue   json.RawMessage `json:"new_value,omitempty"`
	IPAddress  *string         `json:"ip_address,omitempty"`
	UserAgent  *string         `json:"user_agent,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type OnboardingDashboard struct {
	Active               int     `json:"active"`
	Pending              int     `json:"pending"`
	Completed            int     `json:"completed"`
	Overdue              int     `json:"overdue"`
	AverageProgress      float64 `json:"average_progress"`
	TasksDueToday        int     `json:"tasks_due_today"`
	DocumentsPendingReview int   `json:"documents_pending_review"`
}

type EmployeeDashboard struct {
	Status         string        `json:"status"`
	Progress       int           `json:"progress"`
	TasksTotal     int           `json:"tasks_total"`
	TasksCompleted int           `json:"tasks_completed"`
	PendingTasks   int           `json:"pending_tasks"`
	NextMilestone  *MilestoneRef `json:"next_milestone,omitempty"`
}

type MilestoneRef struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type OnboardingFilters struct {
	Status     string
	EmployeeID string
	TemplateID string
	DateFrom   string
	DateTo     string
	Search     string
}
