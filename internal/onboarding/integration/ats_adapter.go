package integration

import (
	"context"
	"log"
)

type ATSAdapter struct{}

func NewATSAdapter() *ATSAdapter {
	return &ATSAdapter{}
}

func (a *ATSAdapter) GetCandidate(ctx context.Context, companyID, candidateID string) (*CandidateInfo, error) {
	log.Printf("[ATSAdapter] GetCandidate company=%s candidate=%s", companyID, candidateID)
	return &CandidateInfo{
		ID:        candidateID,
		CompanyID: companyID,
		FirstName: "Candidate",
		LastName:  "Name",
	}, nil
}

func (a *ATSAdapter) GetHiringProcess(ctx context.Context, companyID, applicationID string) (*HiringProcessInfo, error) {
	log.Printf("[ATSAdapter] GetHiringProcess company=%s application=%s", companyID, applicationID)
	return &HiringProcessInfo{
		ID:            applicationID,
		CandidateID:   "cand-" + applicationID,
		ApplicationID: applicationID,
		Status:        "HIRED",
	}, nil
}

func (a *ATSAdapter) MarkOnboardingStarted(ctx context.Context, companyID, applicationID, onboardingID, employeeID string) error {
	log.Printf("[ATSAdapter] MarkOnboardingStarted company=%s application=%s onboarding=%s employee=%s", companyID, applicationID, onboardingID, employeeID)
	return nil
}

func (a *ATSAdapter) MarkOnboardingCompleted(ctx context.Context, companyID, applicationID string) error {
	log.Printf("[ATSAdapter] MarkOnboardingCompleted company=%s application=%s", companyID, applicationID)
	return nil
}

func (a *ATSAdapter) GetReadyForOnboarding(ctx context.Context, companyID string) ([]CandidateInfo, error) {
	log.Printf("[ATSAdapter] GetReadyForOnboarding company=%s", companyID)
	return nil, nil
}
