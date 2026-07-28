package recruitment

type CreateRequisitionRequest struct {
	PositionID       *string  `json:"position_id"`
	DepartmentID     *string  `json:"department_id"`
	HiringManagerID  *string  `json:"hiring_manager_id"`
	Title            string   `json:"title" binding:"required"`
	Description      *string  `json:"description"`
	Vacancies        *int     `json:"vacancies"`
	EmploymentType   *string  `json:"employment_type"`
	WorkMode         *string  `json:"work_mode"`
	Location         *string  `json:"location"`
	SalaryMin        *float64 `json:"salary_min"`
	SalaryMax        *float64 `json:"salary_max"`
	Currency         *string  `json:"currency"`
	Reason           *string  `json:"reason"`
}

type UpdateRequisitionRequest struct {
	Title            *string  `json:"title"`
	Description      *string  `json:"description"`
	Vacancies        *int     `json:"vacancies"`
	EmploymentType   *string  `json:"employment_type"`
	WorkMode         *string  `json:"work_mode"`
	Location         *string  `json:"location"`
	SalaryMin        *float64 `json:"salary_min"`
	SalaryMax        *float64 `json:"salary_max"`
	Currency         *string  `json:"currency"`
	Reason           *string  `json:"reason"`
}

type CreatePostingRequest struct {
	RequisitionID    string  `json:"requisition_id" binding:"required"`
	Title            string  `json:"title" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	Requirements     *string `json:"requirements"`
	Responsibilities *string `json:"responsibilities"`
	Benefits         *string `json:"benefits"`
	EmploymentType   *string `json:"employment_type"`
	WorkMode         *string `json:"work_mode"`
	Location         *string `json:"location"`
}

type UpdatePostingRequest struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	Requirements     *string `json:"requirements"`
	Responsibilities *string `json:"responsibilities"`
	Benefits         *string `json:"benefits"`
	EmploymentType   *string `json:"employment_type"`
	WorkMode         *string `json:"work_mode"`
	Location         *string `json:"location"`
}

type CreateCandidateRequest struct {
	FirstName      string  `json:"first_name" binding:"required"`
	LastName       string  `json:"last_name" binding:"required"`
	Email          string  `json:"email" binding:"required"`
	Phone          *string `json:"phone"`
	DocumentNumber *string `json:"document_number"`
	Location       *string `json:"location"`
	LinkedInURL    *string `json:"linkedin_url"`
	PortfolioURL   *string `json:"portfolio_url"`
	Source         *string `json:"source"`
}

type UpdateCandidateRequest struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Phone          *string `json:"phone"`
	Location       *string `json:"location"`
	LinkedInURL    *string `json:"linkedin_url"`
	PortfolioURL   *string `json:"portfolio_url"`
}

type CreateApplicationRequest struct {
	CandidateID   string `json:"candidate_id" binding:"required"`
	JobPostingID  string `json:"job_posting_id" binding:"required"`
}

type MoveStageRequest struct {
	ToStage string `json:"to_stage" binding:"required"`
	Notes   *string `json:"notes"`
}

type RejectApplicationRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type CreateInterviewRequest struct {
	ApplicationID   string  `json:"application_id" binding:"required"`
	InterviewerID   string  `json:"interviewer_id" binding:"required"`
	InterviewType   string  `json:"interview_type" binding:"required"`
	ScheduledAt     *string `json:"scheduled_at"`
	DurationMinutes *int    `json:"duration_minutes"`
	MeetingURL      *string `json:"meeting_url"`
	Location        *string `json:"location"`
}

type UpdateInterviewRequest struct {
	ScheduledAt     *string `json:"scheduled_at"`
	DurationMinutes *int    `json:"duration_minutes"`
	MeetingURL      *string `json:"meeting_url"`
	Location        *string `json:"location"`
	Status          *string `json:"status"`
}

type CreateInterviewFeedbackRequest struct {
	Score          *float64 `json:"score"`
	Comments       *string  `json:"comments"`
	Recommendation *string  `json:"recommendation" binding:"required"`
}

type CreateAssessmentRequest struct {
	ApplicationID  string   `json:"application_id" binding:"required"`
	AssessmentType *string  `json:"assessment_type"`
	Title          string   `json:"title" binding:"required"`
	Description    *string  `json:"description"`
	MaxScore       *float64 `json:"max_score"`
	DurationMinutes *int    `json:"duration_minutes"`
}

type CreateOfferRequest struct {
	ApplicationID    string   `json:"application_id" binding:"required"`
	PositionTitle    string   `json:"position_title" binding:"required"`
	DepartmentID     *string  `json:"department_id"`
	StartDate        *string  `json:"start_date"`
	EmploymentType   *string  `json:"employment_type"`
	WorkMode         *string  `json:"work_mode"`
	SalaryAmount     *float64 `json:"salary_amount"`
	SalaryCurrency   *string  `json:"salary_currency"`
	SalaryPeriod     *string  `json:"salary_period"`
	Benefits         *string  `json:"benefits"`
	Conditions       *string  `json:"conditions"`
	ResponseDeadline *string  `json:"response_deadline"`
}

type CreateReferralRequest struct {
	CandidateID string  `json:"candidate_id" binding:"required"`
	Notes       *string `json:"notes"`
}

type CreateScreeningQuestionRequest struct {
	JobPostingID string  `json:"job_posting_id" binding:"required"`
	Question     string  `json:"question" binding:"required"`
	QuestionType *string `json:"question_type"`
	Required     *bool   `json:"required"`
}

type CreateWorkflowRequest struct {
	Name  string              `json:"name" binding:"required"`
	Steps []WorkflowStepRequest `json:"steps" binding:"required"`
}

type WorkflowStepRequest struct {
	StepOrder    int     `json:"step_order" binding:"required"`
	ApproverRole *string `json:"approver_role"`
	ApproverID   *string `json:"approver_id"`
	Required     *bool   `json:"required"`
}
