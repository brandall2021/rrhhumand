package integration

import (
	"context"
)

type OnboardingAdapter struct{}

func NewOnboardingAdapter() *OnboardingAdapter {
	return &OnboardingAdapter{}
}

type CandidateOnboardingData struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	Phone          *string `json:"phone,omitempty"`
	PositionTitle  string  `json:"position_title"`
	DepartmentID   *string `json:"department_id,omitempty"`
	StartDate      string  `json:"start_date"`
	EmploymentType *string `json:"employment_type,omitempty"`	
	WorkMode       *string `json:"work_mode,omitempty"`
}

func (a *OnboardingAdapter) CreateOnboardingFromHire(ctx context.Context, hiringProcessID string, candidateData *CandidateOnboardingData) (string, error) {
	return "onboarding-" + hiringProcessID, nil
}

func (a *OnboardingAdapter) GetOnboardingStatus(ctx context.Context, onboardingID string) (string, error) {
	return "IN_PROGRESS", nil
}

func (a *OnboardingAdapter) CompleteOnboarding(ctx context.Context, onboardingID string) error {
	return nil
}

func (a *OnboardingAdapter) GetOnboardingTasks(ctx context.Context, onboardingID string) ([]string, error) {
	return []string{
		"Complete personal information",
		"Sign employment contract",
		"IT equipment request",
		"Benefits enrollment",
		"First day orientation",
	}, nil
}
