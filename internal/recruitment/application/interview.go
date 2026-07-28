package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreateInterviewReq struct {
	ApplicationID   string     `json:"application_id"`
	InterviewType   string     `json:"interview_type"`
	Title           *string    `json:"title"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	DurationMinutes *int       `json:"duration_minutes"`
	MeetingURL      *string    `json:"meeting_url"`
	MeetingPassword *string    `json:"meeting_password"`
	Location        *string    `json:"location"`
	Instructions    *string    `json:"instructions"`
	CreatedBy       *string    `json:"created_by"`
}

type SubmitFeedbackReq struct {
	Score          *float64                  `json:"score"`
	Comments       *string                   `json:"comments"`
	Strengths      []string                  `json:"strengths"`
	Weaknesses     []string                  `json:"weaknesses"`
	Recommendation *string                   `json:"recommendation"`
	Questions      []SubmitFeedbackQuestionReq `json:"questions,omitempty"`
}

type SubmitFeedbackQuestionReq struct {
	Question string   `json:"question"`
	Score    *float64 `json:"score"`
	Comment  *string  `json:"comment"`
}

type InterviewService struct {
	interviewRepo   *repository.InterviewRepo
	applicationRepo *repository.ApplicationRepo
}

func NewInterviewService(interviewRepo *repository.InterviewRepo, applicationRepo *repository.ApplicationRepo) *InterviewService {
	return &InterviewService{
		interviewRepo:   interviewRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *InterviewService) Create(ctx context.Context, companyID string, req *CreateInterviewReq) (*domain.Interview, error) {
	const op = "CreateInterview"
	now := time.Now()
	i := &domain.Interview{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		ApplicationID:   req.ApplicationID,
		InterviewType:   req.InterviewType,
		Title:           req.Title,
		ScheduledAt:     req.ScheduledAt,
		DurationMinutes: req.DurationMinutes,
		MeetingURL:      req.MeetingURL,
		MeetingPassword: req.MeetingPassword,
		Location:        req.Location,
		Instructions:    req.Instructions,
		Status:          domain.IntStatusScheduled,
		CreatedBy:       req.CreatedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	result, err := s.interviewRepo.Create(ctx, companyID, i)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *InterviewService) GetByID(ctx context.Context, companyID, id string) (*domain.Interview, error) {
	const op = "GetInterview"
	return s.interviewRepo.GetByID(ctx, companyID, id)
}

func (s *InterviewService) List(ctx context.Context, companyID, applicationID, status string) ([]domain.Interview, error) {
	const op = "ListInterviews"
	return s.interviewRepo.List(ctx, companyID, applicationID, status)
}

func (s *InterviewService) Update(ctx context.Context, companyID, id string, req *domain.Interview) (*domain.Interview, error) {
	const op = "UpdateInterview"
	req.UpdatedAt = time.Now()
	result, err := s.interviewRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *InterviewService) Cancel(ctx context.Context, companyID, id string) error {
	const op = "CancelInterview"
	return s.interviewRepo.UpdateStatus(ctx, companyID, id, domain.IntStatusCancelled)
}

func (s *InterviewService) Complete(ctx context.Context, companyID, id string) error {
	const op = "CompleteInterview"
	return s.interviewRepo.UpdateStatus(ctx, companyID, id, domain.IntStatusCompleted)
}

func (s *InterviewService) AddPanelMember(ctx context.Context, companyID, interviewID string, member domain.InterviewPanelMember) (*domain.InterviewPanelMember, error) {
	const op = "AddPanelMember"
	member.ID = uuid.New().String()
	member.InterviewID = interviewID
	member.CreatedAt = time.Now()
	result, err := s.interviewRepo.AddPanelMember(ctx, &member)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *InterviewService) RemovePanelMember(ctx context.Context, companyID, interviewID, memberID string) error {
	const op = "RemovePanelMember"
	return s.interviewRepo.RemovePanelMember(ctx, memberID)
}

func (s *InterviewService) ListPanelMembers(ctx context.Context, companyID, interviewID string) ([]domain.InterviewPanelMember, error) {
	const op = "ListPanelMembers"
	return s.interviewRepo.ListPanelMembers(ctx, interviewID)
}

func (s *InterviewService) SubmitFeedback(ctx context.Context, companyID, interviewID string, panelistID string, req *SubmitFeedbackReq) (*domain.InterviewFeedback, error) {
	const op = "SubmitFeedback"
	now := time.Now()
	fb := &domain.InterviewFeedback{
		ID:             uuid.New().String(),
		InterviewID:    interviewID,
		PanelistID:     panelistID,
		Score:          req.Score,
		Comments:       req.Comments,
		Strengths:      req.Strengths,
		Weaknesses:     req.Weaknesses,
		Recommendation: req.Recommendation,
		SubmittedAt:    &now,
		CreatedAt:      now,
	}
	for _, q := range req.Questions {
		fb.Questions = append(fb.Questions, domain.InterviewFeedbackQuestion{
			ID:       uuid.New().String(),
			Question: q.Question,
			Score:    q.Score,
			Comment:  q.Comment,
		})
	}
	result, err := s.interviewRepo.AddFeedback(ctx, fb)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *InterviewService) ListFeedback(ctx context.Context, companyID, interviewID string) ([]domain.InterviewFeedback, error) {
	const op = "ListFeedback"
	return s.interviewRepo.ListFeedback(ctx, interviewID)
}
