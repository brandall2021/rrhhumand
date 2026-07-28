package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreateAssessmentReq struct {
	ApplicationID   string     `json:"application_id"`
	AssessmentType  string     `json:"assessment_type"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	MaxScore        *float64   `json:"max_score"`
	PassingScore    *float64   `json:"passing_score"`
	DurationMinutes *int       `json:"duration_minutes"`
	DueAt           *time.Time `json:"due_at"`
	CreatedBy       *string    `json:"created_by"`
}

type ScoreAssessmentReq struct {
	Score         *float64 `json:"score"`
	Result        *string  `json:"result"`
	ResultSummary *string  `json:"result_summary"`
}

type AssessmentService struct {
	assessmentRepo  *repository.AssessmentRepo
	applicationRepo *repository.ApplicationRepo
}

func NewAssessmentService(assessmentRepo *repository.AssessmentRepo, applicationRepo *repository.ApplicationRepo) *AssessmentService {
	return &AssessmentService{
		assessmentRepo:  assessmentRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *AssessmentService) Create(ctx context.Context, companyID string, req *CreateAssessmentReq) (*domain.Assessment, error) {
	const op = "CreateAssessment"
	now := time.Now()
	a := &domain.Assessment{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		ApplicationID:   req.ApplicationID,
		AssessmentType:  req.AssessmentType,
		Title:           req.Title,
		Description:     req.Description,
		MaxScore:        req.MaxScore,
		PassingScore:    req.PassingScore,
		DurationMinutes: req.DurationMinutes,
		DueAt:           req.DueAt,
		Status:          domain.AssStatusPending,
		CreatedBy:       req.CreatedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	result, err := s.assessmentRepo.Create(ctx, companyID, a)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *AssessmentService) GetByID(ctx context.Context, companyID, id string) (*domain.Assessment, error) {
	const op = "GetAssessment"
	return s.assessmentRepo.GetByID(ctx, companyID, id)
}

func (s *AssessmentService) List(ctx context.Context, companyID, applicationID, status string) ([]domain.Assessment, error) {
	const op = "ListAssessments"
	return s.assessmentRepo.List(ctx, companyID, applicationID, status)
}

func (s *AssessmentService) Update(ctx context.Context, companyID, id string, req *domain.Assessment) (*domain.Assessment, error) {
	const op = "UpdateAssessment"
	req.UpdatedAt = time.Now()
	result, err := s.assessmentRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *AssessmentService) Send(ctx context.Context, companyID, id string) error {
	const op = "SendAssessment"
	return s.assessmentRepo.UpdateStatus(ctx, companyID, id, domain.AssStatusSent, nil, nil)
}

func (s *AssessmentService) Score(ctx context.Context, companyID, id string, req *ScoreAssessmentReq) error {
	const op = "ScoreAssessment"
	return s.assessmentRepo.UpdateStatus(ctx, companyID, id, domain.AssStatusCompleted, req.Score, req.Result)
}

func (s *AssessmentService) Cancel(ctx context.Context, companyID, id string) error {
	const op = "CancelAssessment"
	return s.assessmentRepo.UpdateStatus(ctx, companyID, id, domain.AssStatusCancelled, nil, nil)
}

func (s *AssessmentService) AddSection(ctx context.Context, companyID, assessmentID string, section domain.AssessmentSection) (*domain.AssessmentSection, error) {
	const op = "AddSection"
	section.ID = uuid.New().String()
	section.AssessmentID = assessmentID
	result, err := s.assessmentRepo.AddSection(ctx, &section)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *AssessmentService) ListSections(ctx context.Context, companyID, assessmentID string) ([]domain.AssessmentSection, error) {
	const op = "ListSections"
	return s.assessmentRepo.ListSections(ctx, assessmentID)
}

func (s *AssessmentService) AddResult(ctx context.Context, companyID, assessmentID string, result domain.AssessmentResult) (*domain.AssessmentResult, error) {
	const op = "AddResult"
	result.ID = uuid.New().String()
	result.AssessmentID = assessmentID
	r, err := s.assessmentRepo.AddResult(ctx, &result)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return r, nil
}

func (s *AssessmentService) ListResults(ctx context.Context, companyID, assessmentID string) ([]domain.AssessmentResult, error) {
	const op = "ListResults"
	return s.assessmentRepo.ListResults(ctx, assessmentID)
}
