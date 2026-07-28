package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type ApplicationService struct {
	applicationRepo *repository.ApplicationRepo
	candidateRepo   *repository.CandidateRepo
	postingRepo     *repository.PostingRepo
}

func NewApplicationService(applicationRepo *repository.ApplicationRepo, candidateRepo *repository.CandidateRepo, postingRepo *repository.PostingRepo) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		candidateRepo:   candidateRepo,
		postingRepo:     postingRepo,
	}
}

func (s *ApplicationService) Create(ctx context.Context, companyID, candidateID, postingID string) (*domain.Application, error) {
	const op = "CreateApplication"
	candidate, err := s.candidateRepo.GetByID(ctx, companyID, candidateID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if candidate.Status == domain.CandStatusBlacklisted {
		return nil, svcErr(op, domain.ErrCandidateBlacklisted)
	}

	posting, err := s.postingRepo.GetByID(ctx, companyID, postingID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if posting.Status != domain.PostStatusPublished {
		return nil, svcErr(op, domain.ErrPostingClosed)
	}

	existing, err := s.applicationRepo.GetByCandidateAndPosting(ctx, companyID, candidateID, postingID)
	if err == nil && existing != nil {
		return nil, svcErr(op, domain.ErrDuplicateApplication)
	}

	now := time.Now()
	a := &domain.Application{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		CandidateID: candidateID,
		PostingID:   postingID,
		Status:      domain.AppStatusNew,
		AppliedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	result, err := s.applicationRepo.Create(ctx, companyID, a)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ApplicationService) GetByID(ctx context.Context, companyID, id string) (*domain.Application, error) {
	const op = "GetApplication"
	return s.applicationRepo.GetByID(ctx, companyID, id)
}

func (s *ApplicationService) List(ctx context.Context, companyID, candidateID, postingID, status string) ([]domain.Application, error) {
	const op = "ListApplications"
	return s.applicationRepo.List(ctx, companyID, candidateID, postingID, status)
}

func (s *ApplicationService) MoveStage(ctx context.Context, companyID, id, toStageID, changedBy, reason string) error {
	const op = "MoveStage"
	a, err := s.applicationRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if !isValidStageTransition(a.Status, domain.ApplicationStatus(toStageID)) {
		return svcErr(op, domain.ErrStageNotAllowed)
	}
	now := time.Now()
	if err := s.applicationRepo.UpdateStage(ctx, companyID, id, toStageID); err != nil {
		return svcErr(op, err)
	}
	history := &domain.ApplicationStageHistory{
		ID:             uuid.New().String(),
		ApplicationID:  id,
		FromStageID:    a.CurrentStageID,
		ToStageID:      toStageID,
		ChangedBy:      &changedBy,
		Reason:         &reason,
		AutoTransition: false,
		CreatedAt:      now,
	}
	return s.applicationRepo.AddStageHistory(ctx, history)
}

func (s *ApplicationService) Reject(ctx context.Context, companyID, id, reasonID, reasonText string) error {
	const op = "RejectApplication"
	a, err := s.applicationRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if !isValidStageTransition(a.Status, domain.AppStatusRejected) {
		return svcErr(op, domain.ErrStageNotAllowed)
	}
	return s.applicationRepo.UpdateStatus(ctx, companyID, id, string(domain.AppStatusRejected))
}

func (s *ApplicationService) Withdraw(ctx context.Context, companyID, id, reason string) error {
	const op = "WithdrawApplication"
	a, err := s.applicationRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if !isValidStageTransition(a.Status, domain.AppStatusWithdrawn) {
		return svcErr(op, domain.ErrStageNotAllowed)
	}
	return s.applicationRepo.UpdateStatus(ctx, companyID, id, string(domain.AppStatusWithdrawn))
}

func (s *ApplicationService) AddNote(ctx context.Context, companyID, applicationID string, note domain.ApplicationNote) (*domain.ApplicationNote, error) {
	const op = "AddNote"
	note.ID = uuid.New().String()
	note.ApplicationID = applicationID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	result, err := s.applicationRepo.AddNote(ctx, &note)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ApplicationService) UpdateNote(ctx context.Context, companyID, applicationID string, note domain.ApplicationNote) error {
	const op = "UpdateNote"
	note.UpdatedAt = time.Now()
	_, err := s.applicationRepo.UpdateNote(ctx, note.ID, &note)
	return err
}

func (s *ApplicationService) ListNotes(ctx context.Context, companyID, applicationID string) ([]domain.ApplicationNote, error) {
	const op = "ListNotes"
	return s.applicationRepo.ListNotes(ctx, applicationID)
}

func (s *ApplicationService) AddRating(ctx context.Context, companyID, applicationID string, rating domain.ApplicationRating) (*domain.ApplicationRating, error) {
	const op = "AddRating"
	rating.ID = uuid.New().String()
	rating.ApplicationID = applicationID
	rating.CreatedAt = time.Now()
	result, err := s.applicationRepo.AddRating(ctx, &rating)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ApplicationService) ListRatings(ctx context.Context, companyID, applicationID string) ([]domain.ApplicationRating, error) {
	const op = "ListRatings"
	return s.applicationRepo.ListRatings(ctx, applicationID)
}

func (s *ApplicationService) GetStageHistory(ctx context.Context, companyID, applicationID string) ([]domain.ApplicationStageHistory, error) {
	const op = "GetStageHistory"
	return s.applicationRepo.ListStageHistory(ctx, applicationID)
}

func isValidStageTransition(from, to domain.ApplicationStatus) bool {
	allowedTransitions := map[domain.ApplicationStatus][]domain.ApplicationStatus{
		domain.AppStatusNew:        {domain.AppStatusScreening, domain.AppStatusRejected, domain.AppStatusWithdrawn, domain.AppStatusOnHold},
		domain.AppStatusScreening:  {domain.AppStatusInterview, domain.AppStatusAssessment, domain.AppStatusRejected, domain.AppStatusWithdrawn, domain.AppStatusOnHold},
		domain.AppStatusInterview:  {domain.AppStatusAssessment, domain.AppStatusOffer, domain.AppStatusRejected, domain.AppStatusWithdrawn, domain.AppStatusOnHold},
		domain.AppStatusAssessment: {domain.AppStatusInterview, domain.AppStatusOffer, domain.AppStatusRejected, domain.AppStatusWithdrawn, domain.AppStatusOnHold},
		domain.AppStatusOffer:      {domain.AppStatusHired, domain.AppStatusRejected, domain.AppStatusWithdrawn},
		domain.AppStatusOnHold:     {domain.AppStatusScreening, domain.AppStatusInterview, domain.AppStatusAssessment, domain.AppStatusOffer, domain.AppStatusRejected},
	}
	allowed, ok := allowedTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
