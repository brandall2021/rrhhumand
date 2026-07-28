package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type WorkflowService struct {
	workflowRepo *repository.WorkflowRepo
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepo) *WorkflowService {
	return &WorkflowService{workflowRepo: workflowRepo}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, w *domain.BenefitWorkflow) (*domain.BenefitWorkflow, error) {
	w.ID = uuid.New()
	w.CompanyID = companyID
	w.CreatedBy = userID
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	if err := s.workflowRepo.CreateWorkflow(ctx, w); err != nil {
		return nil, svcErr("CreateWorkflow", err)
	}
	return w, nil
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitWorkflow, error) {
	return s.workflowRepo.GetWorkflow(ctx, companyID, id)
}

func (s *WorkflowService) ListWorkflows(ctx context.Context, companyID uuid.UUID, benefitID *uuid.UUID) ([]domain.BenefitWorkflow, error) {
	return s.workflowRepo.ListWorkflows(ctx, companyID, benefitID)
}

func (s *WorkflowService) UpdateWorkflow(ctx context.Context, companyID uuid.UUID, w *domain.BenefitWorkflow) (*domain.BenefitWorkflow, error) {
	w.CompanyID = companyID
	w.UpdatedAt = time.Now()
	if err := s.workflowRepo.UpdateWorkflow(ctx, w); err != nil {
		return nil, svcErr("UpdateWorkflow", err)
	}
	return w, nil
}

func (s *WorkflowService) AddStep(ctx context.Context, companyID uuid.UUID, step *domain.BenefitWorkflowStep) (*domain.BenefitWorkflowStep, error) {
	step.ID = uuid.New()
	step.CreatedAt = time.Now()
	if err := s.workflowRepo.CreateStep(ctx, step); err != nil {
		return nil, svcErr("AddStep", err)
	}
	return step, nil
}

func (s *WorkflowService) ListSteps(ctx context.Context, workflowID uuid.UUID) ([]domain.BenefitWorkflowStep, error) {
	return s.workflowRepo.ListSteps(ctx, workflowID)
}

func (s *WorkflowService) UpdateStep(ctx context.Context, step *domain.BenefitWorkflowStep) (*domain.BenefitWorkflowStep, error) {
	if err := s.workflowRepo.UpdateStep(ctx, step); err != nil {
		return nil, svcErr("UpdateStep", err)
	}
	return step, nil
}

func (s *WorkflowService) DeleteStep(ctx context.Context, id uuid.UUID) error {
	return s.workflowRepo.DeleteStep(ctx, uuid.UUID{}, id)
}
