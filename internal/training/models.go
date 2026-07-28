package training

import (
	"encoding/json"
	"time"
)

type CourseCategory struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	ParentID    *string   `json:"parent_id,omitempty"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type Course struct {
	ID                     string     `json:"id"`
	CompanyID              string     `json:"company_id"`
	CategoryID             *string    `json:"category_id,omitempty"`
	Code                   string     `json:"code"`
	Name                   string     `json:"name"`
	ShortDescription       *string    `json:"short_description,omitempty"`
	Description            *string    `json:"description,omitempty"`
	Objectives             *string    `json:"objectives,omitempty"`
	Difficulty             string     `json:"difficulty"`
	DurationMinutes        int        `json:"duration_minutes"`
	Modality               string     `json:"modality"`
	Status                 string     `json:"status"`
	Mandatory              bool       `json:"mandatory"`
	PassingScore           *float64   `json:"passing_score,omitempty"`
	CertificateEnabled     bool       `json:"certificate_enabled"`
	MinAttendancePercentage float64   `json:"min_attendance_percentage"`
	CreatedBy              string     `json:"created_by"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

type CourseVersion struct {
	ID          string     `json:"id"`
	CourseID    string     `json:"course_id"`
	Version     string     `json:"version"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CourseContent struct {
	ID              string    `json:"id"`
	CourseVersionID string    `json:"course_version_id"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty"`
	ContentType     string    `json:"content_type"`
	StorageProvider *string   `json:"storage_provider,omitempty"`
	StorageKey      *string   `json:"storage_key,omitempty"`
	ExternalURL     *string   `json:"external_url,omitempty"`
	DurationSeconds int       `json:"duration_seconds"`
	SortOrder       int       `json:"sort_order"`
	Required        bool      `json:"required"`
	Published       bool      `json:"published"`
	CreatedAt       time.Time `json:"created_at"`
}

type CoursePrerequisite struct {
	ID                   string `json:"id"`
	CourseID             string `json:"course_id"`
	PrerequisiteCourseID string `json:"prerequisite_course_id"`
	Required             bool   `json:"required"`
}

type Instructor struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	EmployeeID      *string   `json:"employee_id,omitempty"`
	InstructorType  string    `json:"instructor_type"`
	Name            string    `json:"name"`
	Email           *string   `json:"email,omitempty"`
	Phone           *string   `json:"phone,omitempty"`
	Specialization  *string   `json:"specialization,omitempty"`
	Bio             *string   `json:"bio,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type TrainingProvider struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	TaxID       *string   `json:"tax_id,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Website     *string   `json:"website,omitempty"`
	ContactName *string   `json:"contact_name,omitempty"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CourseOffering struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	CourseID        string     `json:"course_id"`
	CourseVersionID *string    `json:"course_version_id,omitempty"`
	Name            string     `json:"name"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	EnrollmentStart *time.Time `json:"enrollment_start,omitempty"`
	EnrollmentEnd   *time.Time `json:"enrollment_end,omitempty"`
	Capacity        int        `json:"capacity"`
	EnrolledCount   int        `json:"enrolled_count"`
	Modality        *string    `json:"modality,omitempty"`
	Location        *string    `json:"location,omitempty"`
	MeetingURL      *string    `json:"meeting_url,omitempty"`
	InstructorID    *string    `json:"instructor_id,omitempty"`
	ProviderID      *string    `json:"provider_id,omitempty"`
	CostAmount      *float64   `json:"cost_amount,omitempty"`
	CostCurrency    string     `json:"cost_currency"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TrainingSession struct {
	ID           string    `json:"id"`
	OfferingID   string    `json:"offering_id"`
	Title        *string   `json:"title,omitempty"`
	SessionDate  time.Time `json:"session_date"`
	StartTime    *string   `json:"start_time,omitempty"`
	EndTime      *string   `json:"end_time,omitempty"`
	Location     *string   `json:"location,omitempty"`
	MeetingURL   *string   `json:"meeting_url,omitempty"`
	InstructorID *string   `json:"instructor_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Enrollment struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	OfferingID     string     `json:"offering_id"`
	EmployeeID     string     `json:"employee_id"`
	AssignmentType string     `json:"assignment_type"`
	Status         string     `json:"status"`
	EnrollmentDate time.Time  `json:"enrollment_date"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	DroppedAt      *time.Time `json:"dropped_at,omitempty"`
	FinalScore     *float64   `json:"final_score,omitempty"`
	Passed         *bool      `json:"passed,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type TrainingAssignment struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	CourseID       string    `json:"course_id"`
	AssigneeType   string    `json:"assignee_type"`
	AssigneeID     *string   `json:"assignee_id,omitempty"`
	AssignmentType string    `json:"assignment_type"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	Active         bool      `json:"active"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type TrainingAssignmentRule struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	Name           string    `json:"name"`
	CriteriaField  string    `json:"criteria_field"`
	CriteriaValue  string    `json:"criteria_value"`
	CourseID       string    `json:"course_id"`
	AssignmentType string    `json:"assignment_type"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
}

type CourseProgress struct {
	ID                string     `json:"id"`
	EnrollmentID      string     `json:"enrollment_id"`
	ContentID         string     `json:"content_id"`
	Status            string     `json:"status"`
	ProgressPercentage int        `json:"progress_percentage"`
	TimeSpentSeconds  int        `json:"time_spent_seconds"`
	LastPosition      int        `json:"last_position"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type Assessment struct {
	ID                string    `json:"id"`
	CourseID          string    `json:"course_id"`
	Title             string    `json:"title"`
	Description       *string   `json:"description,omitempty"`
	AssessmentType    string    `json:"assessment_type"`
	AttemptsAllowed   int       `json:"attempts_allowed"`
	PassingScore      *float64  `json:"passing_score,omitempty"`
	TimeLimitMinutes  *int      `json:"time_limit_minutes,omitempty"`
	RandomizeQuestions bool     `json:"randomize_questions"`
	ShowResults       bool      `json:"show_results"`
	Status            string    `json:"status"`
	SortOrder         int       `json:"sort_order"`
	CreatedAt         time.Time `json:"created_at"`
}

type AssessmentQuestion struct {
	ID           string    `json:"id"`
	AssessmentID string    `json:"assessment_id"`
	Question     string    `json:"question"`
	QuestionType string    `json:"question_type"`
	Points       float64   `json:"points"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
}

type AssessmentOption struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	SortOrder  int    `json:"sort_order"`
}

type AssessmentAttempt struct {
	ID            string     `json:"id"`
	AssessmentID  string     `json:"assessment_id"`
	EnrollmentID  string     `json:"enrollment_id"`
	EmployeeID    string     `json:"employee_id"`
	AttemptNumber int        `json:"attempt_number"`
	Status        string     `json:"status"`
	Score         *float64   `json:"score,omitempty"`
	Passed        *bool      `json:"passed,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	SubmittedAt   *time.Time `json:"submitted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type AssessmentAnswer struct {
	ID              string   `json:"id"`
	AttemptID       string   `json:"attempt_id"`
	QuestionID      string   `json:"question_id"`
	SelectedOptionID *string  `json:"selected_option_id,omitempty"`
	TextAnswer      *string  `json:"text_answer,omitempty"`
	NumericAnswer   *float64 `json:"numeric_answer,omitempty"`
	IsCorrect       *bool    `json:"is_correct,omitempty"`
	Score           float64  `json:"score"`
}

type Certificate struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	CourseID        string     `json:"course_id"`
	EnrollmentID    string     `json:"enrollment_id"`
	CertificateNo   string     `json:"certificate_no"`
	IssuedAt        time.Time  `json:"issued_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	Status          string     `json:"status"`
	StorageProvider *string    `json:"storage_provider,omitempty"`
	StorageKey      *string    `json:"storage_key,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Competency struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	CompetencyType string    `json:"competency_type"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
}

type CompetencyLevel struct {
	ID           string  `json:"id"`
	CompetencyID string  `json:"competency_id"`
	Level        int     `json:"level"`
	Label        string  `json:"label"`
	Description  *string `json:"description,omitempty"`
}

type CourseCompetency struct {
	ID            string `json:"id"`
	CourseID      string `json:"course_id"`
	CompetencyID  string `json:"competency_id"`
	AcquiredLevel int    `json:"acquired_level"`
}

type EmployeeCompetency struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	CompetencyID string     `json:"competency_id"`
	Level        int        `json:"level"`
	Source       string     `json:"source"`
	Verified     bool       `json:"verified"`
	VerifiedBy   *string    `json:"verified_by,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CompetencyGap struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EmployeeID   string    `json:"employee_id"`
	CompetencyID string    `json:"competency_id"`
	RequiredLevel int      `json:"required_level"`
	CurrentLevel int       `json:"current_level"`
	Gap          int       `json:"gap"`
	Source       *string   `json:"source,omitempty"`
	SourceID     *string   `json:"source_id,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TrainingNeed struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EmployeeID   *string   `json:"employee_id,omitempty"`
	CompetencyID *string   `json:"competency_id,omitempty"`
	Title        string    `json:"title"`
	Description  *string   `json:"description,omitempty"`
	Priority     string    `json:"priority"`
	Source       *string   `json:"source,omitempty"`
	SourceID     *string   `json:"source_id,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TrainingPlan struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"company_id"`
	EmployeeID    *string    `json:"employee_id,omitempty"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Objectives    *string    `json:"objectives,omitempty"`
	PeriodStart   *time.Time `json:"period_start,omitempty"`
	PeriodEnd     *time.Time `json:"period_end,omitempty"`
	BudgetAmount  *float64   `json:"budget_amount,omitempty"`
	BudgetCurrency string    `json:"budget_currency"`
	Status        string     `json:"status"`
	CreatedBy     string     `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type TrainingPlanCourse struct {
	ID        string `json:"id"`
	PlanID    string `json:"plan_id"`
	CourseID  string `json:"course_id"`
	Priority  string `json:"priority"`
	SortOrder int    `json:"sort_order"`
}

type LearningPath struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Objectives  *string   `json:"objectives,omitempty"`
	DurationDays *int     `json:"duration_days,omitempty"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LearningPathCourse struct {
	ID        string `json:"id"`
	PathID    string `json:"path_id"`
	CourseID  string `json:"course_id"`
	Required  bool   `json:"required"`
	SortOrder int    `json:"sort_order"`
}

type LearningPathEnrollment struct {
	ID                 string     `json:"id"`
	CompanyID          string     `json:"company_id"`
	PathID             string     `json:"path_id"`
	EmployeeID         string     `json:"employee_id"`
	Status             string     `json:"status"`
	ProgressPercentage int        `json:"progress_percentage"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type TrainingAttendance struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	SessionID    string     `json:"session_id"`
	EnrollmentID string     `json:"enrollment_id"`
	EmployeeID   string     `json:"employee_id"`
	Status       string     `json:"status"`
	CheckIn      *time.Time `json:"check_in,omitempty"`
	CheckOut     *time.Time `json:"check_out,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CourseFeedback struct {
	ID                string   `json:"id"`
	EnrollmentID      string   `json:"enrollment_id"`
	CompanyID         string   `json:"company_id"`
	EmployeeID        string   `json:"employee_id"`
	InstructorRating  *int     `json:"instructor_rating,omitempty"`
	ContentRating     *int     `json:"content_rating,omitempty"`
	OrganizationRating *int    `json:"organization_rating,omitempty"`
	PlatformRating    *int     `json:"platform_rating,omitempty"`
	OverallRating     *float64 `json:"overall_rating,omitempty"`
	Comments          *string  `json:"comments,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type TrainingCost struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EnrollmentID *string   `json:"enrollment_id,omitempty"`
	OfferingID   *string   `json:"offering_id,omitempty"`
	CostType     string    `json:"cost_type"`
	Description  *string   `json:"description,omitempty"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	IncurredDate time.Time `json:"incurred_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type TrainingBudget struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Year         int       `json:"year"`
	DepartmentID *string   `json:"department_id,omitempty"`
	TotalAmount  float64   `json:"total_amount"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
}

type TrainingDashboard struct {
	ActiveCourses        int     `json:"active_courses"`
	TrainedEmployees     int     `json:"trained_employees"`
	CompletedCourses     int     `json:"completed_courses"`
	PendingCourses       int     `json:"pending_courses"`
	ActiveCertificates   int     `json:"active_certificates"`
	ExpiringCertificates int     `json:"expiring_certificates"`
	TotalTrainingHours   float64 `json:"total_training_hours"`
	TotalCost            float64 `json:"total_cost"`
	AvgPassRate          float64 `json:"avg_pass_rate"`
	AvgSatisfaction      float64 `json:"avg_satisfaction"`
}

type EmployeeTrainingDashboard struct {
	Enrolled       int     `json:"enrolled"`
	InProgress     int     `json:"in_progress"`
	Completed      int     `json:"completed"`
	Certificates   int     `json:"certificates"`
	TotalHours     float64 `json:"total_hours"`
	Competencies   int     `json:"competencies"`
}

type CourseFilter struct {
	Status     string
	CategoryID string
	Modality   string
	Difficulty string
	Search     string
}
