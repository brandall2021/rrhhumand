package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type SettingsService struct {
	settingsRepo *repository.SettingsRepo
}

func NewSettingsService(settingsRepo *repository.SettingsRepo) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *SettingsService) CreateSource(ctx context.Context, companyID string, src *domain.RecruitmentSource) (*domain.RecruitmentSource, error) {
	const op = "CreateSource"
	src.ID = uuid.New().String()
	src.CompanyID = companyID
	src.Active = true
	src.CreatedAt = time.Now()
	result, err := s.settingsRepo.CreateSource(ctx, companyID, src)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) ListSources(ctx context.Context, companyID string) ([]domain.RecruitmentSource, error) {
	const op = "ListSources"
	return s.settingsRepo.ListSources(ctx, companyID)
}

func (s *SettingsService) UpdateSource(ctx context.Context, companyID string, src *domain.RecruitmentSource) (*domain.RecruitmentSource, error) {
	const op = "UpdateSource"
	result, err := s.settingsRepo.UpdateSource(ctx, companyID, src.ID, src)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) CreateStage(ctx context.Context, companyID string, stage *domain.RecruitmentStage) (*domain.RecruitmentStage, error) {
	const op = "CreateStage"
	stage.ID = uuid.New().String()
	stage.CompanyID = companyID
	stage.Active = true
	stage.CreatedAt = time.Now()
	result, err := s.settingsRepo.CreateStage(ctx, companyID, stage)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) ListStages(ctx context.Context, companyID string) ([]domain.RecruitmentStage, error) {
	const op = "ListStages"
	return s.settingsRepo.ListStages(ctx, companyID)
}

func (s *SettingsService) UpdateStage(ctx context.Context, companyID string, stage *domain.RecruitmentStage) (*domain.RecruitmentStage, error) {
	const op = "UpdateStage"
	result, err := s.settingsRepo.UpdateStage(ctx, companyID, stage.ID, stage)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) ReorderStages(ctx context.Context, companyID string, stageIDs []string) error {
	const op = "ReorderStages"
	return s.settingsRepo.ReorderStages(ctx, stageIDs)
}

func (s *SettingsService) CreateTransition(ctx context.Context, companyID string, t *domain.StageTransition) (*domain.StageTransition, error) {
	const op = "CreateTransition"
	t.ID = uuid.New().String()
	t.CompanyID = companyID
	t.CreatedAt = time.Now()
	result, err := s.settingsRepo.CreateTransition(ctx, t)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) ListTransitions(ctx context.Context, companyID string) ([]domain.StageTransition, error) {
	const op = "ListTransitions"
	return s.settingsRepo.ListTransitions(ctx, companyID)
}

func (s *SettingsService) DeleteTransition(ctx context.Context, companyID, transitionID string) error {
	const op = "DeleteTransition"
	return s.settingsRepo.DeleteTransition(ctx, companyID, transitionID)
}

func (s *SettingsService) CreateRejectionReason(ctx context.Context, companyID string, r *domain.RejectionReason) (*domain.RejectionReason, error) {
	const op = "CreateRejectionReason"
	r.ID = uuid.New().String()
	r.CompanyID = companyID
	r.Active = true
	r.CreatedAt = time.Now()
	result, err := s.settingsRepo.CreateRejectionReason(ctx, companyID, r)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *SettingsService) ListRejectionReasons(ctx context.Context, companyID string) ([]domain.RejectionReason, error) {
	const op = "ListRejectionReasons"
	return s.settingsRepo.ListRejectionReasons(ctx, companyID)
}

func (s *SettingsService) UpdateRejectionReason(ctx context.Context, companyID string, r *domain.RejectionReason) (*domain.RejectionReason, error) {
	const op = "UpdateRejectionReason"
	result, err := s.settingsRepo.UpdateRejectionReason(ctx, companyID, r.ID, r)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}
