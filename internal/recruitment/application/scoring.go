package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/engine"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type ScoringService struct {
	scoringRepo    *repository.ScoringRepo
	scoringEngine  *engine.ScoringEngine
	matchingEngine *engine.MatchingEngine
	candidateRepo  *repository.CandidateRepo
	positionRepo   *repository.PositionRepo
}

func NewScoringService(
	scoringRepo *repository.ScoringRepo,
	scoringEngine *engine.ScoringEngine,
	matchingEngine *engine.MatchingEngine,
	candidateRepo *repository.CandidateRepo,
	positionRepo *repository.PositionRepo,
) *ScoringService {
	return &ScoringService{
		scoringRepo:    scoringRepo,
		scoringEngine:  scoringEngine,
		matchingEngine: matchingEngine,
		candidateRepo:  candidateRepo,
		positionRepo:   positionRepo,
	}
}

func (s *ScoringService) CreateModel(ctx context.Context, companyID string, model *domain.ScoringModel) (*domain.ScoringModel, error) {
	const op = "CreateScoringModel"
	model.ID = uuid.New().String()
	model.CompanyID = companyID
	model.CreatedAt = time.Now()
	if !model.Active {
		model.Active = true
	}
	result, err := s.scoringRepo.CreateScoringModel(ctx, companyID, model)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ScoringService) GetModel(ctx context.Context, companyID, id string) (*domain.ScoringModel, error) {
	const op = "GetScoringModel"
	model, err := s.scoringRepo.GetScoringModel(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return model, nil
}

func (s *ScoringService) ListModels(ctx context.Context, companyID string) ([]domain.ScoringModel, error) {
	const op = "ListScoringModels"
	models, err := s.scoringRepo.ListScoringModels(ctx, companyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return models, nil
}

func (s *ScoringService) UpdateModel(ctx context.Context, companyID, id string, model *domain.ScoringModel) (*domain.ScoringModel, error) {
	const op = "UpdateScoringModel"
	result, err := s.scoringRepo.UpdateScoringModel(ctx, companyID, id, model)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ScoringService) DeleteModel(ctx context.Context, companyID, id string) error {
	const op = "DeleteScoringModel"
	return s.scoringRepo.DeleteScoringModel(ctx, companyID, id)
}

func (s *ScoringService) AddCriterion(ctx context.Context, companyID, modelID string, criterion domain.ScoringCriterion) (*domain.ScoringCriterion, error) {
	const op = "AddCriterion"
	model, err := s.scoringRepo.GetScoringModel(ctx, companyID, modelID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if model.CompanyID != companyID {
		return nil, svcErr(op, domain.ErrNotFound)
	}
	criterion.ID = uuid.New().String()
	criterion.ModelID = modelID
	result, err := s.scoringRepo.AddCriterion(ctx, &criterion)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *ScoringService) UpdateCriterion(ctx context.Context, companyID, modelID string, criterion domain.ScoringCriterion) error {
	const op = "UpdateCriterion"
	model, err := s.scoringRepo.GetScoringModel(ctx, companyID, modelID)
	if err != nil {
		return svcErr(op, err)
	}
	if model.CompanyID != companyID {
		return svcErr(op, domain.ErrNotFound)
	}
	_, err = s.scoringRepo.UpdateCriterion(ctx, criterion.ID, &criterion)
	return err
}

func (s *ScoringService) DeleteCriterion(ctx context.Context, companyID, modelID, criterionID string) error {
	const op = "DeleteCriterion"
	model, err := s.scoringRepo.GetScoringModel(ctx, companyID, modelID)
	if err != nil {
		return svcErr(op, err)
	}
	if model.CompanyID != companyID {
		return svcErr(op, domain.ErrNotFound)
	}
	return s.scoringRepo.DeleteCriterion(ctx, criterionID)
}

func (s *ScoringService) ListCriteria(ctx context.Context, companyID, modelID string) ([]domain.ScoringCriterion, error) {
	const op = "ListCriteria"
	model, err := s.scoringRepo.GetScoringModel(ctx, companyID, modelID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if model.CompanyID != companyID {
		return nil, svcErr(op, domain.ErrNotFound)
	}
	return s.scoringRepo.ListCriteria(ctx, modelID)
}

func (s *ScoringService) ScoreCandidate(ctx context.Context, companyID, candidateID, positionID string) (*domain.MatchingResult, error) {
	const op = "ScoreCandidate"
	result, err := s.matchingEngine.Match(ctx, companyID, candidateID, positionID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}
