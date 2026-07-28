package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type WorkflowEngine struct {
	workflowRepo   *repository.WorkflowRepo
	assignmentRepo *repository.AssignmentRepo
}

func NewWorkflowEngine(workflowRepo *repository.WorkflowRepo, assignmentRepo *repository.AssignmentRepo) *WorkflowEngine {
	return &WorkflowEngine{
		workflowRepo:   workflowRepo,
		assignmentRepo: assignmentRepo,
	}
}

func (e *WorkflowEngine) ProcessRequest(ctx context.Context, req *domain.BenefitRequest, reviewerID uuid.UUID, action string, comment string) (*domain.BenefitRequest, error) {
	workflowID := uuid.Nil
	if req.ID != uuid.Nil {
		workflows, err := e.workflowRepo.ListWorkflows(ctx, req.CompanyID, &req.BenefitID)
		if err != nil {
			return nil, fmt.Errorf("benefits_engine.workflow.ProcessRequest: %w", err)
		}
		for _, w := range workflows {
			if w.IsActive {
				workflowID = w.ID
				break
			}
		}
	}

	now := time.Now()
	switch action {
	case "approve":
		req.Status = "APPROVED"
	case "reject":
		req.Status = "REJECTED"
	default:
		req.Status = action
	}
	req.ResolvedBy = &reviewerID
	req.ResolvedAt = &now
	req.ResolutionNotes = &comment
	req.UpdatedAt = now

	rev := &domain.BenefitRequestReview{
		ID:         uuid.New(),
		RequestID:  req.ID,
		ReviewerID: reviewerID,
		ReviewType: action,
		Comment:    &comment,
		ReviewedAt: now,
	}
	if err := e.assignmentRepo.CreateReview(ctx, rev); err != nil {
		return nil, fmt.Errorf("benefits_engine.workflow.ProcessRequest: %w", err)
	}

	if workflowID != uuid.Nil {
		steps, err := e.workflowRepo.ListSteps(ctx, workflowID)
		if err != nil {
			return nil, fmt.Errorf("benefits_engine.workflow.ProcessRequest: %w", err)
		}
		if len(steps) > 0 {
			var reqStepID *uuid.UUID
			for _, step := range steps {
				if reqStepID == nil && action == "approve" {
					reqStepID = &step.ID
				}
			}
			if reqStepID != nil {
				rev.StepID = reqStepID
				if err := e.assignmentRepo.CreateReview(ctx, rev); err != nil {
					return nil, fmt.Errorf("benefits_engine.workflow.ProcessRequest: %w", err)
				}
			}
		}
	}

	return req, nil
}

func (e *WorkflowEngine) GetNextStep(ctx context.Context, workflowID uuid.UUID, currentStep int) (*domain.BenefitWorkflowStep, error) {
	steps, err := e.workflowRepo.ListSteps(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("benefits_engine.workflow.GetNextStep: %w", err)
	}
	for _, step := range steps {
		if step.StepOrder == currentStep+1 {
			return &step, nil
		}
	}
	return nil, nil
}

func (e *WorkflowEngine) AutoApprove(ctx context.Context, req *domain.BenefitRequest) error {
	workflows, err := e.workflowRepo.ListWorkflows(ctx, req.CompanyID, &req.BenefitID)
	if err != nil {
		return fmt.Errorf("benefits_engine.workflow.AutoApprove: %w", err)
	}
	for _, w := range workflows {
		if w.IsActive && w.AutoApprove {
			now := time.Now()
			req.Status = "APPROVED"
			req.ResolvedAt = &now
			req.UpdatedAt = now
			return nil
		}
	}
	return nil
}
