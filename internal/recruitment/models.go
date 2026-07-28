package recruitment

import (
	"encoding/json"
	"time"
)

type JobRequisition struct {
	ID               string     `json:"id"`
	CompanyID        string     `json:"company_id"`
	PositionID       *string    `json:"position_id,omitempty"`
	DepartmentID     *string    `json:"department_id,omitempty"`
	RequestedBy      string     `json:"requested_by"`
	HiringManagerID  *string    `json:"hiring_manager_id,omitempty"`
	Title            string     `json:"title"`
	Description      *string    `json:"description,omitempty"`
	Vacancies        int        `json:"vacancies"`
	EmploymentType   *string    `json:"employment_type,omitempty"`
	WorkMode         *string    `json:"work_mode,omitempty"`
	Location         *string    `json:"location,omitempty"`
	SalaryMin        *float64   `json:"salary_min,omitempty"`
	SalaryMax        *float64   `json:"salary_max,omitempty"`
	Currency         *string    `json:"currency,omitempty"`
	Reason           *string    `json:"reason,omitempty"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type ApprovalWorkflow struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	Name       string    `json:"name"`
	EntityType string    `json:"entity_type"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type ApprovalStep struct {
	ID            string  `json:"id"`
	WorkflowID    string  `json:"workflow_id"`
	StepOrder     int     `json:"step_order"`
	ApproverRole  *string `json:"approver_role,omitempty"`
	ApproverID    *string `json:"approver_id,omitempty"`
	Required      bool    `json:"required"`
}

type ApprovalInstance struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	WorkflowID  string    `json:"workflow_id"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	CurrentStep int       `json:"current_step"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JobPosting struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	RequisitionID   string     `json:"requisition_id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Requirements    *string    `json:"requirements,omitempty"`
	Responsibilities *string   `json:"responsibilities,omitempty"`
	Benefits        *string    `json:"benefits,omitempty"`
	EmploymentType  *string    `json:"employment_type,omitempty"`
	WorkMode        *string    `json:"work_mode,omitempty"`
	Location        *string    `json:"location,omitempty"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	ClosingAt       *time.Time `json:"closing_at,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Candidate struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          *string   `json:"phone,omitempty"`
	DocumentNumber *string   `json:"document_number,omitempty"`
	Location       *string   `json:"location,omitempty"`
	LinkedInURL    *string   `json:"linkedin_url,omitempty"`
	PortfolioURL   *string   `json:"portfolio_url,omitempty"`
	Source         *string   `json:"source,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Application struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	CandidateID   string     `json:"candidate_id"`
	CandidateName string     `json:"candidate_name,omitempty"`
	JobPostingID  string     `json:"job_posting_id"`
	PostingTitle  string     `json:"posting_title,omitempty"`
	Status        string     `json:"status"`
	AppliedAt     time.Time  `json:"applied_at"`
	RejectedAt    *time.Time `json:"rejected_at,omitempty"`
	RejectionReason *string  `json:"rejection_reason,omitempty"`
	HiredAt       *time.Time `json:"hired_at,omitempty"`
}

type CandidateStageHistory struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	ApplicationID string     `json:"application_id"`
	FromStage     *string    `json:"from_stage,omitempty"`
	ToStage       string     `json:"to_stage"`
	ChangedBy     *string    `json:"changed_by,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	ChangedAt     time.Time  `json:"changed_at"`
}

type CandidateDocument struct {
	ID           string          `json:"id"`
	CompanyID    string          `json:"company_id"`
	CandidateID  string          `json:"candidate_id"`
	DocumentType string          `json:"document_type"`
	FileName     string          `json:"file_name"`
	MimeType     *string         `json:"mime_type,omitempty"`
	SizeBytes    *int64          `json:"size_bytes,omitempty"`
	StorageProvider *string      `json:"storage_provider,omitempty"`
	StorageKey   *string         `json:"storage_key,omitempty"`
	ParsedData   json.RawMessage `json:"parsed_data,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ScreeningQuestion struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	JobPostingID string    `json:"job_posting_id"`
	Question     string    `json:"question"`
	QuestionType string    `json:"question_type"`
	Required     bool      `json:"required"`
	SortOrder    int       `json:"sort_order"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

type ScreeningAnswer struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	ApplicationID string    `json:"application_id"`
	QuestionID    string    `json:"question_id"`
	Answer        string    `json:"answer"`
	CreatedAt     time.Time `json:"created_at"`
}

type Interview struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	ApplicationID   string     `json:"application_id"`
	CandidateName   string     `json:"candidate_name,omitempty"`
	InterviewerID   string     `json:"interviewer_id"`
	InterviewerName string     `json:"interviewer_name,omitempty"`
	InterviewType   string     `json:"interview_type"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	MeetingURL      *string    `json:"meeting_url,omitempty"`
	Location        *string    `json:"location,omitempty"`
	Status          string     `json:"status"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type InterviewFeedback struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	InterviewID     string     `json:"interview_id"`
	InterviewerID   string     `json:"interviewer_id"`
	Score           *float64   `json:"score,omitempty"`
	Comments        *string    `json:"comments,omitempty"`
	Recommendation  *string    `json:"recommendation,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Assessment struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	ApplicationID  string     `json:"application_id"`
	AssessmentType string     `json:"assessment_type"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	MaxScore       *float64   `json:"max_score,omitempty"`
	Score          *float64   `json:"score,omitempty"`
	DurationMinutes *int      `json:"duration_minutes,omitempty"`
	Status         string     `json:"status"`
	Result         *string    `json:"result,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type JobOffer struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	ApplicationID   string     `json:"application_id"`
	CandidateName   string     `json:"candidate_name,omitempty"`
	PositionTitle   string     `json:"position_title"`
	DepartmentID    *string    `json:"department_id,omitempty"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EmploymentType  *string    `json:"employment_type,omitempty"`
	WorkMode        *string    `json:"work_mode,omitempty"`
	SalaryAmount    *float64   `json:"salary_amount,omitempty"`
	SalaryCurrency  *string    `json:"salary_currency,omitempty"`
	SalaryPeriod    *string    `json:"salary_period,omitempty"`
	Benefits        *string    `json:"benefits,omitempty"`
	Conditions      *string    `json:"conditions,omitempty"`
	ResponseDeadline *time.Time `json:"response_deadline,omitempty"`
	Status          string     `json:"status"`
	CreatedBy       *string    `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type EmployeeReferral struct {
	ID                 string    `json:"id"`
	CompanyID          string    `json:"company_id"`
	ReferrerEmployeeID string    `json:"referrer_employee_id"`
	ReferrerName       string    `json:"referrer_name,omitempty"`
	CandidateID        string    `json:"candidate_id"`
	CandidateName      string    `json:"candidate_name,omitempty"`
	ApplicationID      *string   `json:"application_id,omitempty"`
	Status             string    `json:"status"`
	RewardStatus       string    `json:"reward_status"`
	Notes              *string   `json:"notes,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type RecruitmentAuditLog struct {
	ID          string          `json:"id"`
	CompanyID   string          `json:"company_id"`
	UserID      *string         `json:"user_id,omitempty"`
	CandidateID *string         `json:"candidate_id,omitempty"`
	EntityType  string          `json:"entity_type"`
	EntityID    string          `json:"entity_id"`
	Action      string          `json:"action"`
	OldValue    json.RawMessage `json:"old_value,omitempty"`
	NewValue    json.RawMessage `json:"new_value,omitempty"`
	IPAddress   *string         `json:"ip_address,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type RecruitmentDashboard struct {
	OpenRequisitions   int     `json:"open_requisitions"`
	TotalCandidates    int     `json:"total_candidates"`
	ApplicationsThisWeek int   `json:"applications_this_week"`
	PendingOffers      int     `json:"pending_offers"`
	HiresThisMonth     int     `json:"hires_this_month"`
	TotalInterviews    int     `json:"total_interviews"`
	AvgTimeToHire      float64 `json:"avg_time_to_hire"`
	FunnelByStage      []StageCount `json:"funnel_by_stage,omitempty"`
}

type StageCount struct {
	Stage string `json:"stage"`
	Count int    `json:"count"`
}

type RecruitmentFilters struct {
	RequisitionID string
	PostingID     string
	CandidateID   string
	ApplicationID string
	Status        string
	Stage         string
	InterviewerID string
	DateFrom      string
	DateTo        string
	Source        string
}
