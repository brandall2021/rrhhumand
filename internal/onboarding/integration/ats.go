package integration

import "context"

type CandidateInfo struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Phone        *string `json:"phone,omitempty"`
	PositionID   *string `json:"position_id,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
	ManagerID    *string `json:"manager_id,omitempty"`
	JobOfferID   *string `json:"job_offer_id,omitempty"`
	ApplicationID *string `json:"application_id,omitempty"`
	StartDate    *string `json:"start_date,omitempty"`
	WorkMode     *string `json:"work_mode,omitempty"`
	EmployeeType *string `json:"employee_type,omitempty"`
	Salary       *string `json:"salary,omitempty"`
	LocationID   *string `json:"location_id,omitempty"`
}

type HiringProcessInfo struct {
	ID           string `json:"id"`
	CandidateID  string `json:"candidate_id"`
	ApplicationID string `json:"application_id"`
	OfferID      string `json:"offer_id"`
	Status       string `json:"status"`
	EmployeeID   *string `json:"employee_id,omitempty"`
	OnboardingID *string `json:"onboarding_id,omitempty"`
}

type ATSIntegration interface {
	GetCandidate(ctx context.Context, companyID, candidateID string) (*CandidateInfo, error)
	GetHiringProcess(ctx context.Context, companyID, applicationID string) (*HiringProcessInfo, error)
	MarkOnboardingStarted(ctx context.Context, companyID, applicationID, onboardingID, employeeID string) error
	MarkOnboardingCompleted(ctx context.Context, companyID, applicationID string) error
	GetReadyForOnboarding(ctx context.Context, companyID string) ([]CandidateInfo, error)
}
