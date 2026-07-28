package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type WorkflowService struct {
	workflowRepo *repository.WorkflowRepo
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepo) *WorkflowService {
	return &WorkflowService{workflowRepo: workflowRepo}
}

func (s *WorkflowService) Create(ctx context.Context, companyID string, wf *domain.Workflow) (*domain.Workflow, error) {
	const op = "CreateWorkflow"
	wf.ID = uuid.New().String()
	wf.CompanyID = companyID
	now := time.Now()
	wf.CreatedAt = now
	wf.UpdatedAt = now
	result, err := s.workflowRepo.Create(ctx, companyID, wf)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *WorkflowService) GetByID(ctx context.Context, companyID, id string) (*domain.Workflow, error) {
	const op = "GetWorkflow"
	return s.workflowRepo.GetByID(ctx, companyID, id)
}

func (s *WorkflowService) List(ctx context.Context, companyID, entityType string) ([]domain.Workflow, error) {
	const op = "ListWorkflows"
	return s.workflowRepo.List(ctx, companyID, entityType)
}

func (s *WorkflowService) Update(ctx context.Context, companyID, id string, wf *domain.Workflow) (*domain.Workflow, error) {
	const op = "UpdateWorkflow"
	wf.UpdatedAt = time.Now()
	result, err := s.workflowRepo.Update(ctx, companyID, id, wf)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *WorkflowService) Activate(ctx context.Context, companyID, id string) error {
	const op = "ActivateWorkflow"
	return s.workflowRepo.Activate(ctx, companyID, id)
}

func (s *WorkflowService) Deactivate(ctx context.Context, companyID, id string) error {
	const op = "DeactivateWorkflow"
	return s.workflowRepo.Deactivate(ctx, companyID, id)
}

func (s *WorkflowService) AddStage(ctx context.Context, companyID, workflowID string, stage domain.WorkflowStage) (*domain.WorkflowStage, error) {
	const op = "AddWorkflowStage"
	stage.ID = uuid.New().String()
	stage.WorkflowID = workflowID
	stage.CreatedAt = time.Now()
	result, err := s.workflowRepo.AddStage(ctx, &stage)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *WorkflowService) RemoveStage(ctx context.Context, companyID, workflowID, stageID string) error {
	const op = "RemoveWorkflowStage"
	return s.workflowRepo.RemoveStage(ctx, stageID)
}

func (s *WorkflowService) ListStages(ctx context.Context, companyID, workflowID string) ([]domain.WorkflowStage, error) {
	const op = "ListWorkflowStages"
	return s.workflowRepo.ListStages(ctx, workflowID)
}

func (s *WorkflowService) ReorderStages(ctx context.Context, companyID, workflowID string, stageIDs []string) error {
	const op = "ReorderWorkflowStages"
	return s.workflowRepo.ReorderStages(ctx, stageIDs)
}

func (s *WorkflowService) AddRule(ctx context.Context, companyID, workflowID string, rule domain.WorkflowRule) (*domain.WorkflowRule, error) {
	const op = "AddWorkflowRule"
	rule.ID = uuid.New().String()
	rule.WorkflowID = workflowID
	rule.CreatedAt = time.Now()
	result, err := s.workflowRepo.AddRule(ctx, &rule)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *WorkflowService) UpdateRule(ctx context.Context, companyID, workflowID string, rule domain.WorkflowRule) error {
	const op = "UpdateWorkflowRule"
	_, err := s.workflowRepo.UpdateRule(ctx, rule.ID, &rule)
	return err
}

func (s *WorkflowService) DeleteRule(ctx context.Context, companyID, workflowID, ruleID string) error {
	const op = "DeleteWorkflowRule"
	return s.workflowRepo.DeleteRule(ctx, ruleID)
}

func (s *WorkflowService) ListRules(ctx context.Context, companyID, workflowID string) ([]domain.WorkflowRule, error) {
	const op = "ListWorkflowRules"
	return s.workflowRepo.ListRules(ctx, workflowID)
}
