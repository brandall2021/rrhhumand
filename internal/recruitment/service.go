package recruitment

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Requisitions
func (s *Service) CreateRequisition(ctx context.Context, companyID, userID string, req *CreateRequisitionRequest) (*JobRequisition, error) {
	return s.repo.CreateRequisition(ctx, companyID, userID, req)
}

func (s *Service) GetRequisition(ctx context.Context, companyID, id string) (*JobRequisition, error) {
	return s.repo.GetRequisition(ctx, companyID, id)
}

func (s *Service) ListRequisitions(ctx context.Context, companyID string, filters RecruitmentFilters) ([]JobRequisition, error) {
	return s.repo.ListRequisitions(ctx, companyID, filters)
}

func (s *Service) UpdateRequisition(ctx context.Context, companyID, id string, req *UpdateRequisitionRequest) (*JobRequisition, error) {
	rec, err := s.repo.GetRequisition(ctx, companyID, id)
	if err != nil { return nil, err }
	if rec.Status != "DRAFT" {
		return nil, fmt.Errorf("can only update DRAFT requisitions")
	}
	return s.repo.UpdateRequisition(ctx, companyID, id, req)
}

func (s *Service) SubmitRequisition(ctx context.Context, companyID, id string) error {
	rec, err := s.repo.GetRequisition(ctx, companyID, id)
	if err != nil { return err }
	if rec.Status != "DRAFT" {
		return fmt.Errorf("can only submit DRAFT requisitions")
	}
	return s.repo.UpdateRequisitionStatus(ctx, companyID, id, "PENDING_APPROVAL")
}

func (s *Service) ApproveRequisition(ctx context.Context, companyID, id string) error {
	rec, err := s.repo.GetRequisition(ctx, companyID, id)
	if err != nil { return err }
	if rec.Status != "PENDING_APPROVAL" {
		return fmt.Errorf("can only approve PENDING_APPROVAL requisitions")
	}
	return s.repo.UpdateRequisitionStatus(ctx, companyID, id, "APPROVED")
}

func (s *Service) OpenRequisition(ctx context.Context, companyID, id string) error {
	rec, err := s.repo.GetRequisition(ctx, companyID, id)
	if err != nil { return err }
	if rec.Status != "APPROVED" {
		return fmt.Errorf("can only open APPROVED requisitions")
	}
	return s.repo.UpdateRequisitionStatus(ctx, companyID, id, "OPEN")
}

func (s *Service) CloseRequisition(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateRequisitionStatus(ctx, companyID, id, "CLOSED")
}

func (s *Service) CancelRequisition(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateRequisitionStatus(ctx, companyID, id, "CANCELLED")
}

// Postings
func (s *Service) CreatePosting(ctx context.Context, companyID string, req *CreatePostingRequest) (*JobPosting, error) {
	return s.repo.CreatePosting(ctx, companyID, req)
}

func (s *Service) GetPosting(ctx context.Context, companyID, id string) (*JobPosting, error) {
	return s.repo.GetPosting(ctx, companyID, id)
}

func (s *Service) ListPostings(ctx context.Context, companyID string, filters RecruitmentFilters) ([]JobPosting, error) {
	return s.repo.ListPostings(ctx, companyID, filters)
}

func (s *Service) PublishPosting(ctx context.Context, companyID, id string) error {
	p, err := s.repo.GetPosting(ctx, companyID, id)
	if err != nil { return err }
	if p.Status != "DRAFT" {
		return fmt.Errorf("can only publish DRAFT postings")
	}
	return s.repo.UpdatePostingStatus(ctx, companyID, id, "PUBLISHED")
}

func (s *Service) ClosePosting(ctx context.Context, companyID, id string) error {
	return s.repo.UpdatePostingStatus(ctx, companyID, id, "CLOSED")
}

// Candidates
func (s *Service) CreateCandidate(ctx context.Context, companyID string, req *CreateCandidateRequest) (*Candidate, error) {
	return s.repo.CreateCandidate(ctx, companyID, req)
}

func (s *Service) GetCandidate(ctx context.Context, companyID, id string) (*Candidate, error) {
	return s.repo.GetCandidate(ctx, companyID, id)
}

func (s *Service) ListCandidates(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Candidate, error) {
	return s.repo.ListCandidates(ctx, companyID, filters)
}

func (s *Service) UpdateCandidate(ctx context.Context, companyID, id string, req *UpdateCandidateRequest) (*Candidate, error) {
	return s.repo.UpdateCandidate(ctx, companyID, id, req)
}

// Applications
func (s *Service) CreateApplication(ctx context.Context, companyID string, req *CreateApplicationRequest) (*Application, error) {
	return s.repo.CreateApplication(ctx, companyID, req)
}

func (s *Service) GetApplication(ctx context.Context, companyID, id string) (*Application, error) {
	return s.repo.GetApplication(ctx, companyID, id)
}

func (s *Service) ListApplications(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Application, error) {
	return s.repo.ListApplications(ctx, companyID, filters)
}

func (s *Service) MoveStage(ctx context.Context, companyID, applicationID, userID string, req *MoveStageRequest) (*CandidateStageHistory, error) {
	app, err := s.repo.GetApplication(ctx, companyID, applicationID)
	if err != nil { return nil, err }

	if err := ValidateStageTransition(app.Status, req.ToStage); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateApplicationStatus(ctx, companyID, applicationID, req.ToStage); err != nil {
		return nil, err
	}

	return s.repo.AddStageHistory(ctx, companyID, applicationID, &app.Status, req.ToStage, &userID, req.Notes)
}

func (s *Service) RejectApplication(ctx context.Context, companyID, applicationID string, req *RejectApplicationRequest) error {
	return s.repo.RejectApplication(ctx, companyID, applicationID, req.Reason)
}

func (s *Service) WithdrawApplication(ctx context.Context, companyID, applicationID string) error {
	return s.repo.UpdateApplicationStatus(ctx, companyID, applicationID, "WITHDRAWN")
}

func (s *Service) GetStageHistory(ctx context.Context, applicationID string) ([]CandidateStageHistory, error) {
	return s.repo.GetStageHistory(ctx, applicationID)
}

// Screening
func (s *Service) CreateScreeningQuestion(ctx context.Context, companyID string, req *CreateScreeningQuestionRequest) (*ScreeningQuestion, error) {
	return s.repo.CreateScreeningQuestion(ctx, companyID, req)
}

func (s *Service) ListScreeningQuestions(ctx context.Context, postingID string) ([]ScreeningQuestion, error) {
	return s.repo.ListScreeningQuestions(ctx, postingID)
}

func (s *Service) SubmitScreeningAnswers(ctx context.Context, companyID, applicationID string, answers map[string]string) error {
	for questionID, answer := range answers {
		_, err := s.repo.CreateScreeningAnswer(ctx, companyID, applicationID, questionID, answer)
		if err != nil { return err }
	}
	return nil
}

// Interviews
func (s *Service) CreateInterview(ctx context.Context, companyID string, req *CreateInterviewRequest) (*Interview, error) {
	return s.repo.CreateInterview(ctx, companyID, req)
}

func (s *Service) GetInterview(ctx context.Context, companyID, id string) (*Interview, error) {
	return s.repo.GetInterview(ctx, companyID, id)
}

func (s *Service) ListInterviews(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Interview, error) {
	return s.repo.ListInterviews(ctx, companyID, filters)
}

func (s *Service) UpdateInterview(ctx context.Context, companyID, id string, req *UpdateInterviewRequest) (*Interview, error) {
	return s.repo.UpdateInterview(ctx, companyID, id, req)
}

func (s *Service) CreateInterviewFeedback(ctx context.Context, companyID, interviewID, interviewerID string, req *CreateInterviewFeedbackRequest) (*InterviewFeedback, error) {
	return s.repo.CreateInterviewFeedback(ctx, companyID, interviewID, interviewerID, req)
}

func (s *Service) ListInterviewFeedback(ctx context.Context, interviewID string) ([]InterviewFeedback, error) {
	return s.repo.ListInterviewFeedback(ctx, interviewID)
}

// Assessments
func (s *Service) CreateAssessment(ctx context.Context, companyID string, req *CreateAssessmentRequest) (*Assessment, error) {
	return s.repo.CreateAssessment(ctx, companyID, req)
}

func (s *Service) ListAssessments(ctx context.Context, applicationID string) ([]Assessment, error) {
	return s.repo.ListAssessments(ctx, applicationID)
}

// Offers
func (s *Service) CreateOffer(ctx context.Context, companyID, userID string, req *CreateOfferRequest) (*JobOffer, error) {
	return s.repo.CreateOffer(ctx, companyID, userID, req)
}

func (s *Service) GetOffer(ctx context.Context, companyID, id string) (*JobOffer, error) {
	return s.repo.GetOffer(ctx, companyID, id)
}

func (s *Service) SendOffer(ctx context.Context, companyID, id string) error {
	o, err := s.repo.GetOffer(ctx, companyID, id)
	if err != nil { return err }
	if o.Status != "DRAFT" && o.Status != "PENDING_APPROVAL" {
		return fmt.Errorf("can only send DRAFT/PENDING_APPROVAL offers")
	}
	return s.repo.UpdateOfferStatus(ctx, companyID, id, "SENT")
}

func (s *Service) AcceptOffer(ctx context.Context, companyID, id string) error {
	o, err := s.repo.GetOffer(ctx, companyID, id)
	if err != nil { return err }
	if o.Status != "SENT" {
		return fmt.Errorf("can only accept SENT offers")
	}
	if err := s.repo.UpdateOfferStatus(ctx, companyID, id, "ACCEPTED"); err != nil {
		return err
	}
	return s.repo.UpdateApplicationStatus(ctx, companyID, o.ApplicationID, "HIRED")
}

func (s *Service) RejectOffer(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateOfferStatus(ctx, companyID, id, "REJECTED")
}

// Referrals
func (s *Service) CreateReferral(ctx context.Context, companyID, referrerID string, req *CreateReferralRequest) (*EmployeeReferral, error) {
	return s.repo.CreateReferral(ctx, companyID, referrerID, req)
}

func (s *Service) ListReferrals(ctx context.Context, companyID string) ([]EmployeeReferral, error) {
	return s.repo.ListReferrals(ctx, companyID)
}

// Audit
func (s *Service) AuditLog(ctx context.Context, companyID, userID, candidateID, entityType, entityID, action string, oldVal, newVal interface{}, ip string) error {
	var oldB, newB []byte
	if oldVal != nil { oldB, _ = json.Marshal(oldVal) }
	if newVal != nil { newB, _ = json.Marshal(newVal) }
	return s.repo.CreateAuditLog(ctx, companyID, userID, candidateID, entityType, entityID, action, oldB, newB, ip)
}

// Dashboard
func (s *Service) GetDashboard(ctx context.Context, companyID string) (*RecruitmentDashboard, error) {
	return s.repo.GetDashboard(ctx, companyID)
}

// Hire - convert candidate to employee (integration point with FASE 3)
func (s *Service) HireCandidate(ctx context.Context, companyID, applicationID string) (map[string]interface{}, error) {
	app, err := s.repo.GetApplication(ctx, companyID, applicationID)
	if err != nil { return nil, err }
	if app.Status != "HIRED" {
		return nil, fmt.Errorf("application must be HIRED before conversion")
	}

	candidate, err := s.repo.GetCandidate(ctx, companyID, app.CandidateID)
	if err != nil { return nil, err }

	return map[string]interface{}{
		"candidate_id": candidate.ID,
		"first_name":   candidate.FirstName,
		"last_name":    candidate.LastName,
		"email":        candidate.Email,
		"phone":        candidate.Phone,
		"posting_id":   app.JobPostingID,
		"posting_title": app.PostingTitle,
		"hired_at":     app.HiredAt,
		"message":      "Candidate ready for employee creation - integrate with FASE 3",
	}, nil
}
