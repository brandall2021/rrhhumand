package domain

type ExitInterview struct {
	ID            string
	CompanyID     string
	OffboardingID string
	EmployeeID    string
	InterviewerID *string
	ScheduledAt   *string
	CompletedAt   *string
	Reason        *string
	Feedback      *string
	Recommendation *string
	Rating        *float64
	Anonymous     bool
	CreatedAt     string
	UpdatedAt     string
}

type ExitInterviewQuestion struct {
	ID           string
	CompanyID    string
	Question     string
	QuestionType string
	SortOrder    int
	Active       bool
	CreatedAt    string
}

type ExitInterviewAnswer struct {
	ID             string
	ExitInterviewID string
	QuestionID     string
	Answer         *string
	Rating         *int
	CreatedAt      string
}
