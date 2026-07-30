package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreatePostingReq struct {
	PositionID       string     `json:"position_id"`
	RequisitionID    *string    `json:"requisition_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Requirements     *string    `json:"requirements"`
	Responsibilities *string    `json:"responsibilities"`
	Benefits         *string    `json:"benefits"`
	EmploymentType   *string    `json:"employment_type"`
	WorkMode         *string    `json:"work_mode"`
	Location         *string    `json:"location"`
	SalaryMin        *float64   `json:"salary_min"`
	SalaryMax        *float64   `json:"salary_max"`
	Currency         *string    `json:"currency"`
	IsPublic         bool       `json:"is_public"`
	ClosingAt        *time.Time `json:"closing_at"`
}

type PostingService struct {
	postingRepo  *repository.PostingRepo
	positionRepo *repository.PositionRepo
}

func NewPostingService(postingRepo *repository.PostingRepo, positionRepo *repository.PositionRepo) *PostingService {
	return &PostingService{
		postingRepo:  postingRepo,
		positionRepo: positionRepo,
	}
}

func (s *PostingService) Create(ctx context.Context, companyID string, req *CreatePostingReq) (*domain.Posting, error) {
	const op = "CreatePosting"
	now := time.Now()
	p := &domain.Posting{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		PositionID:       req.PositionID,
		RequisitionID:    req.RequisitionID,
		Title:            req.Title,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Benefits:         req.Benefits,
		EmploymentType:   req.EmploymentType,
		WorkMode:         req.WorkMode,
		Location:         req.Location,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Currency:         req.Currency,
		IsPublic:         req.IsPublic,
		ClosingAt:         req.ClosingAt,
		Status:           domain.PostStatusDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	result, err := s.postingRepo.Create(ctx, companyID, p)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PostingService) GetByID(ctx context.Context, companyID, id string) (*domain.Posting, error) {
	const op = "GetPosting"
	return s.postingRepo.GetByID(ctx, companyID, id)
}

func (s *PostingService) List(ctx context.Context, companyID, status string) ([]domain.Posting, error) {
	const op = "ListPostings"
	return s.postingRepo.List(ctx, companyID, status)
}

func (s *PostingService) Update(ctx context.Context, companyID, id string, req *domain.Posting) (*domain.Posting, error) {
	const op = "UpdatePosting"
	req.UpdatedAt = time.Now()
	result, err := s.postingRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PostingService) Publish(ctx context.Context, companyID, id string) error {
	const op = "PublishPosting"
	return s.postingRepo.UpdateStatus(ctx, companyID, id, string(domain.PostStatusPublished))
}

func (s *PostingService) Close(ctx context.Context, companyID, id string) error {
	const op = "ClosePosting"
	return s.postingRepo.UpdateStatus(ctx, companyID, id, string(domain.PostStatusClosed))
}

func (s *PostingService) AddScreeningQuestion(ctx context.Context, companyID, postingID string, q domain.PostingScreeningQuestion) (*domain.PostingScreeningQuestion, error) {
	const op = "AddScreeningQuestion"
	q.ID = uuid.New().String()
	q.PostingID = postingID
	q.Active = true
	q.CreatedAt = time.Now()
	result, err := s.postingRepo.AddScreeningQuestion(ctx, &q)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PostingService) UpdateScreeningQuestion(ctx context.Context, companyID, postingID, questionID string, q domain.PostingScreeningQuestion) error {
	const op = "UpdateScreeningQuestion"
	_, err := s.postingRepo.UpdateScreeningQuestion(ctx, questionID, &q)
	return err
}

func (s *PostingService) DeleteScreeningQuestion(ctx context.Context, companyID, postingID, questionID string) error {
	const op = "DeleteScreeningQuestion"
	return s.postingRepo.DeleteScreeningQuestion(ctx, questionID)
}

func (s *PostingService) ListScreeningQuestions(ctx context.Context, companyID, postingID string) ([]domain.PostingScreeningQuestion, error) {
	const op = "ListScreeningQuestions"
	return s.postingRepo.ListScreeningQuestions(ctx, postingID)
}

func (s *PostingService) ListPublic(ctx context.Context) ([]domain.Posting, error) {
	const op = "ListPublicPostings"
	return s.postingRepo.ListPublic(ctx)
}

func (s *PostingService) GetPublicByID(ctx context.Context, id string) (*domain.Posting, error) {
	const op = "GetPublicPosting"
	return s.postingRepo.GetPublicByID(ctx, id)
}

type PublicApplyReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}
