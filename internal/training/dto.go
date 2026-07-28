package training

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Active      *bool   `json:"active"`
}

type CreateCourseRequest struct {
	Code                   string   `json:"code" binding:"required"`
	Name                   string   `json:"name" binding:"required"`
	CategoryID             *string  `json:"category_id"`
	ShortDescription       *string  `json:"short_description"`
	Description            *string  `json:"description"`
	Objectives             *string  `json:"objectives"`
	Difficulty             string   `json:"difficulty"`
	DurationMinutes        int      `json:"duration_minutes"`
	Modality               string   `json:"modality"`
	Mandatory              *bool    `json:"mandatory"`
	PassingScore           *float64 `json:"passing_score"`
	CertificateEnabled     *bool    `json:"certificate_enabled"`
	MinAttendancePercentage *float64 `json:"min_attendance_percentage"`
}

type UpdateCourseRequest struct {
	Name                   *string  `json:"name"`
	CategoryID             *string  `json:"category_id"`
	ShortDescription       *string  `json:"short_description"`
	Description            *string  `json:"description"`
	Objectives             *string  `json:"objectives"`
	Difficulty             *string  `json:"difficulty"`
	DurationMinutes        *int     `json:"duration_minutes"`
	Modality               *string  `json:"modality"`
	Status                 *string  `json:"status"`
	Mandatory              *bool    `json:"mandatory"`
	PassingScore           *float64 `json:"passing_score"`
	CertificateEnabled     *bool    `json:"certificate_enabled"`
	MinAttendancePercentage *float64 `json:"min_attendance_percentage"`
}

type CreateVersionRequest struct {
	Version     string  `json:"version" binding:"required"`
	Description *string `json:"description"`
}

type CreateContentRequest struct {
	Title           string  `json:"title" binding:"required"`
	Description     *string `json:"description"`
	ContentType     string  `json:"content_type"`
	ExternalURL     *string `json:"external_url"`
	DurationSeconds *int    `json:"duration_seconds"`
	SortOrder       *int    `json:"sort_order"`
	Required        *bool   `json:"required"`
}

type UpdateContentRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	ContentType     *string `json:"content_type"`
	ExternalURL     *string `json:"external_url"`
	DurationSeconds *int    `json:"duration_seconds"`
	SortOrder       *int    `json:"sort_order"`
	Required        *bool   `json:"required"`
	Published       *bool   `json:"published"`
}

type CreateOfferingRequest struct {
	CourseID        string  `json:"course_id" binding:"required"`
	CourseVersionID *string `json:"course_version_id"`
	Name            string  `json:"name" binding:"required"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	Capacity        *int    `json:"capacity"`
	Modality        *string `json:"modality"`
	Location        *string `json:"location"`
	MeetingURL      *string `json:"meeting_url"`
	InstructorID    *string `json:"instructor_id"`
	ProviderID      *string `json:"provider_id"`
	CostAmount      *float64 `json:"cost_amount"`
	CostCurrency    *string `json:"cost_currency"`
}

type UpdateOfferingRequest struct {
	Name            *string  `json:"name"`
	StartDate       *string  `json:"start_date"`
	EndDate         *string  `json:"end_date"`
	Capacity        *int     `json:"capacity"`
	Modality        *string  `json:"modality"`
	Location        *string  `json:"location"`
	MeetingURL      *string  `json:"meeting_url"`
	InstructorID    *string  `json:"instructor_id"`
	ProviderID      *string  `json:"provider_id"`
	CostAmount      *float64 `json:"cost_amount"`
	CostCurrency    *string  `json:"cost_currency"`
	Status          *string  `json:"status"`
}

type CreateSessionRequest struct {
	Title        *string `json:"title"`
	SessionDate  string  `json:"session_date" binding:"required"`
	StartTime    *string `json:"start_time"`
	EndTime      *string `json:"end_time"`
	Location     *string `json:"location"`
	MeetingURL   *string `json:"meeting_url"`
	InstructorID *string `json:"instructor_id"`
}

type EnrollRequest struct {
	EmployeeID     string `json:"employee_id" binding:"required"`
	AssignmentType string `json:"assignment_type"`
}

type CreateAssignmentRequest struct {
	CourseID       string `json:"course_id" binding:"required"`
	AssigneeType   string `json:"assignee_type" binding:"required"`
	AssigneeID     string `json:"assignee_id"`
	AssignmentType string `json:"assignment_type"`
	DueDate        *string `json:"due_date"`
}

type CreateAssignmentRuleRequest struct {
	Name           string `json:"name" binding:"required"`
	CriteriaField  string `json:"criteria_field" binding:"required"`
	CriteriaValue  string `json:"criteria_value" binding:"required"`
	CourseID       string `json:"course_id" binding:"required"`
	AssignmentType string `json:"assignment_type"`
}

type UpdateProgressRequest struct {
	ProgressPercentage int `json:"progress_percentage"`
	TimeSpentSeconds   int `json:"time_spent_seconds"`
	LastPosition       int `json:"last_position"`
}

type CreateAssessmentRequest struct {
	Title              string   `json:"title" binding:"required"`
	Description        *string  `json:"description"`
	AssessmentType     string   `json:"assessment_type"`
	AttemptsAllowed    *int     `json:"attempts_allowed"`
	PassingScore       *float64 `json:"passing_score"`
	TimeLimitMinutes   *int     `json:"time_limit_minutes"`
	RandomizeQuestions *bool    `json:"randomize_questions"`
	ShowResults        *bool    `json:"show_results"`
}

type UpdateAssessmentRequest struct {
	Title              *string  `json:"title"`
	Description        *string  `json:"description"`
	AttemptsAllowed    *int     `json:"attempts_allowed"`
	PassingScore       *float64 `json:"passing_score"`
	TimeLimitMinutes   *int     `json:"time_limit_minutes"`
	RandomizeQuestions *bool    `json:"randomize_questions"`
	ShowResults        *bool    `json:"show_results"`
	Status             *string  `json:"status"`
}

type CreateQuestionRequest struct {
	Question     string   `json:"question" binding:"required"`
	QuestionType string   `json:"question_type"`
	Points       *float64 `json:"points"`
	Options      []CreateOptionRequest `json:"options"`
}

type CreateOptionRequest struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  *bool  `json:"is_correct"`
}

type SubmitAttemptRequest struct {
	Answers []AnswerRequest `json:"answers" binding:"required"`
}

type AnswerRequest struct {
	QuestionID      string  `json:"question_id" binding:"required"`
	SelectedOptionID *string `json:"selected_option_id"`
	TextAnswer      *string `json:"text_answer"`
	NumericAnswer   *float64 `json:"numeric_answer"`
}

type CreateInstructorRequest struct {
	EmployeeID     *string `json:"employee_id"`
	InstructorType string  `json:"instructor_type"`
	Name           string  `json:"name" binding:"required"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Specialization *string `json:"specialization"`
	Bio            *string `json:"bio"`
}

type CreateProviderRequest struct {
	Name        string  `json:"name" binding:"required"`
	TaxID       *string `json:"tax_id"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Website     *string `json:"website"`
	ContactName *string `json:"contact_name"`
	Notes       *string `json:"notes"`
}

type CreateCompetencyRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    *string `json:"description"`
	CompetencyType string `json:"competency_type"`
	Levels         []CreateLevelRequest `json:"levels"`
}

type CreateLevelRequest struct {
	Level       int    `json:"level" binding:"required"`
	Label       string `json:"label" binding:"required"`
	Description *string `json:"description"`
}

type AssignCompetencyRequest struct {
	Level      int     `json:"level" binding:"required"`
	Source     string  `json:"source"`
	Verified   *bool   `json:"verified"`
	VerifiedBy *string `json:"verified_by"`
}

type CreateTrainingNeedRequest struct {
	EmployeeID   *string `json:"employee_id"`
	CompetencyID *string `json:"competency_id"`
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description"`
	Priority     string  `json:"priority"`
	Source       *string `json:"source"`
	SourceID     *string `json:"source_id"`
}

type CreateTrainingPlanRequest struct {
	EmployeeID    *string  `json:"employee_id"`
	Name          string   `json:"name" binding:"required"`
	Description   *string  `json:"description"`
	Objectives    *string  `json:"objectives"`
	PeriodStart   *string  `json:"period_start"`
	PeriodEnd     *string  `json:"period_end"`
	BudgetAmount  *float64 `json:"budget_amount"`
	BudgetCurrency *string `json:"budget_currency"`
	CourseIDs     []string `json:"course_ids"`
}

type CreateLearningPathRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Objectives  *string  `json:"objectives"`
	DurationDays *int    `json:"duration_days"`
	CourseIDs   []string `json:"course_ids"`
}

type CreateFeedbackRequest struct {
	InstructorRating  *int    `json:"instructor_rating"`
	ContentRating     *int    `json:"content_rating"`
	OrganizationRating *int   `json:"organization_rating"`
	PlatformRating    *int    `json:"platform_rating"`
	OverallRating     *float64 `json:"overall_rating"`
	Comments          *string `json:"comments"`
}

type CreateAttendanceRequest struct {
	Status   string  `json:"status" binding:"required"`
	CheckIn  *string `json:"check_in"`
	CheckOut *string `json:"check_out"`
	Notes    *string `json:"notes"`
}

type AIRecommendationRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Objective  string `json:"objective"`
}

type AIRecommendation struct {
	CourseID       string `json:"course_id"`
	CourseName     string `json:"course_name"`
	Reason         string `json:"reason"`
	CompetencyID   string `json:"competency_id,omitempty"`
	CompetencyName string `json:"competency_name,omitempty"`
	ExpectedLevel  int    `json:"expected_level"`
	Priority       string `json:"priority"`
}

type CourseWithDetails struct {
	Course       Course              `json:"course"`
	Versions     []CourseVersion     `json:"versions,omitempty"`
	Contents     []CourseContent     `json:"contents,omitempty"`
	Offerings    []CourseOffering    `json:"offerings,omitempty"`
	Assessments  []Assessment        `json:"assessments,omitempty"`
	Competencies []CourseCompetency  `json:"competencies,omitempty"`
}
