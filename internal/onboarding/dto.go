package onboarding

import "time"

type CreateTemplateRequest struct {
	Name                string  `json:"name" binding:"required"`
	Description         *string `json:"description"`
	DefaultDurationDays *int    `json:"default_duration_days"`
}

type UpdateTemplateRequest struct {
	Name                *string `json:"name"`
	Description         *string `json:"description"`
	Status              *string `json:"status"`
	DefaultDurationDays *int    `json:"default_duration_days"`
}

type CreateTemplateTaskRequest struct {
	Title             string  `json:"title" binding:"required"`
	Description       *string `json:"description"`
	Category          string  `json:"category"`
	ResponsibleType   string  `json:"responsible_type"`
	ResponsibleUserID *string `json:"responsible_user_id"`
	Required          *bool   `json:"required"`
	DaysOffset        int     `json:"days_offset"`
	SortOrder         *int    `json:"sort_order"`
	EstimatedMinutes  *int    `json:"estimated_minutes"`
}

type UpdateTemplateTaskRequest struct {
	Title             *string `json:"title"`
	Description       *string `json:"description"`
	Category          *string `json:"category"`
	ResponsibleType   *string `json:"responsible_type"`
	ResponsibleUserID *string `json:"responsible_user_id"`
	Required          *bool   `json:"required"`
	DaysOffset        *int    `json:"days_offset"`
	SortOrder         *int    `json:"sort_order"`
	EstimatedMinutes  *int    `json:"estimated_minutes"`
}

type CreateOnboardingRequest struct {
	EmployeeID       string  `json:"employee_id" binding:"required"`
	TemplateID       *string `json:"template_id"`
	StartDate        string  `json:"start_date" binding:"required"`
	CompletionPolicy *string `json:"completion_policy"`
}

type UpdateOnboardingRequest struct {
	TemplateID       *string `json:"template_id"`
	CompletionPolicy *string `json:"completion_policy"`
}

type CreateTaskRequest struct {
	Title             string  `json:"title" binding:"required"`
	Description       *string `json:"description"`
	Category          string  `json:"category"`
	ResponsibleType   string  `json:"responsible_type"`
	ResponsibleID     *string `json:"responsible_id"`
	DueDate           string  `json:"due_date" binding:"required"`
	Required          *bool   `json:"required"`
	SortOrder         *int    `json:"sort_order"`
	EstimatedMinutes  *int    `json:"estimated_minutes"`
}

type UpdateTaskRequest struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	Category         *string `json:"category"`
	ResponsibleType  *string `json:"responsible_type"`
	ResponsibleID    *string `json:"responsible_id"`
	DueDate          *string `json:"due_date"`
	Required         *bool   `json:"required"`
	SortOrder        *int    `json:"sort_order"`
	EstimatedMinutes *int    `json:"estimated_minutes"`
	Status           *string `json:"status"`
	BlockedReason    *string `json:"blocked_reason"`
}

type UploadDocumentRequest struct {
	DocumentType string `form:"document_type" binding:"required"`
	Required     *bool  `form:"required"`
}

type ReviewDocumentRequest struct {
	Status string `json:"status" binding:"required"`
	Reason *string `json:"reason"`
}

type CreateAssetRequest struct {
	AssetType    string  `json:"asset_type" binding:"required"`
	Description  *string `json:"description"`
	SerialNumber *string `json:"serial_number"`
	Notes        *string `json:"notes"`
}

type UpdateAssetRequest struct {
	Status       *string `json:"status"`
	SerialNumber *string `json:"serial_number"`
	Notes        *string `json:"notes"`
}

type CreateAccessRequest struct {
	SystemName string  `json:"system_name" binding:"required"`
	AccessType string  `json:"access_type"`
	Notes      *string `json:"notes"`
}

type CreateMilestoneRequest struct {
	MilestoneType   string `json:"milestone_type" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Description     *string `json:"description"`
	DaysOffset      int    `json:"days_offset"`
	ResponsibleType string `json:"responsible_type"`
	ResponsibleID   *string `json:"responsible_id"`
}

type UpdateMilestoneRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	DaysOffset      *int    `json:"days_offset"`
	ResponsibleType *string `json:"responsible_type"`
	ResponsibleID   *string `json:"responsible_id"`
	Status          *string `json:"status"`
}

type CreateFeedbackRequest struct {
	FeedbackType       string  `json:"feedback_type" binding:"required"`
	AdaptationScore    *int    `json:"adaptation_score"`
	TeamScore          *int    `json:"team_score"`
	KnowledgeScore     *int    `json:"knowledge_score"`
	CommunicationScore *int    `json:"communication_score"`
	OverallScore       *float64 `json:"overall_score"`
	Comments           *string `json:"comments"`
}

type UpdateFeedbackRequest struct {
	AdaptationScore    *int     `json:"adaptation_score"`
	TeamScore          *int     `json:"team_score"`
	KnowledgeScore     *int     `json:"knowledge_score"`
	CommunicationScore *int     `json:"communication_score"`
	OverallScore       *float64 `json:"overall_score"`
	Comments           *string  `json:"comments"`
}

type AssignBuddyRequest struct {
	BuddyEmployeeID string `json:"buddy_employee_id" binding:"required"`
	Notes           *string `json:"notes"`
}

type CreateExceptionRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
}

type CreateTrainingAssignmentRequest struct {
	CourseName       string  `json:"course_name" binding:"required"`
	Description      *string `json:"description"`
	TrainingType     string  `json:"training_type"`
	DueDate          *string `json:"due_date"`
	ExternalProvider *string `json:"external_provider"`
	ExternalCourseID *string `json:"external_course_id"`
}

type CompleteOnboardingRequest struct {
	Reason *string `json:"reason"`
}

type CancelOnboardingRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type CandidateHiredEvent struct {
	CandidateID  string `json:"candidate_id"`
	CompanyID    string `json:"company_id"`
	PositionID   string `json:"position_id"`
	DepartmentID string `json:"department_id"`
	ManagerID    string `json:"manager_id"`
	StartDate    string `json:"start_date"`
	EmployeeID   string `json:"employee_id"`
	EmploymentType string `json:"employment_type"`
	WorkMode     string `json:"work_mode"`
}

type OnboardingCreatedEvent struct {
	OnboardingID string `json:"onboarding_id"`
	CompanyID    string `json:"company_id"`
	EmployeeID   string `json:"employee_id"`
	StartDate    string `json:"start_date"`
}

type IAOnboardingRequest struct {
	Position   string `json:"position" binding:"required"`
	Department string `json:"department" binding:"required"`
	WorkMode   string `json:"work_mode"`
}

type IATemplateProposal struct {
	Tasks []IATaskProposal `json:"tasks"`
}

type IATaskProposal struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	ResponsibleType string `json:"responsible_type"`
	DaysOffset      int    `json:"days_offset"`
}

type TemplateWithTasks struct {
	Template OnboardingTemplate       `json:"template"`
	Tasks    []OnboardingTemplateTask `json:"tasks,omitempty"`
}

type ProcessWithDetails struct {
	Process      OnboardingProcess       `json:"process"`
	Tasks        []OnboardingTask        `json:"tasks,omitempty"`
	Documents    []OnboardingDocument    `json:"documents,omitempty"`
	Assets       []OnboardingAsset       `json:"assets,omitempty"`
	Access       []AccessRequest         `json:"access,omitempty"`
	Milestones   []OnboardingMilestone   `json:"milestones,omitempty"`
	Feedback     []OnboardingFeedback    `json:"feedback,omitempty"`
	Buddies      []OnboardingBuddy       `json:"buddies,omitempty"`
	Training     []TrainingAssignment    `json:"training,omitempty"`
	Progress     int                     `json:"progress"`
}

func parseDate(s string) (time.Time, error) {
	formats := []string{"2006-01-02", time.RFC3339}
	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}
