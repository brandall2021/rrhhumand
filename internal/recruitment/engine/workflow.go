package engine

import (
	"context"
	"time"

	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type WorkflowEngine struct {
	workflowRepo       *repository.WorkflowRepo
	applicationRepo    *repository.ApplicationRepo
}

func NewWorkflowEngine(workflowRepo *repository.WorkflowRepo, applicationRepo *repository.ApplicationRepo) *WorkflowEngine {
	return &WorkflowEngine{
		workflowRepo:    workflowRepo,
		applicationRepo: applicationRepo,
	}
}

func (e *WorkflowEngine) EvaluateTransition(ctx context.Context, application domain.Application, toStageID string) (bool, error) {
	workflows, err := e.workflowRepo.FindByEntityType(ctx, application.CompanyID, domain.WfEntityApplication)
	if err != nil {
		return false, err
	}

	if len(workflows) == 0 {
		return true, nil
	}

	for _, wf := range workflows {
		if !wf.Active {
			continue
		}

		stages, err := e.workflowRepo.ListStages(ctx, wf.ID)
		if err != nil {
			continue
		}

		for _, stage := range stages {
			if stage.StageID == toStageID {
				if len(stage.RequiredActions) == 0 {
					return true, nil
				}
				return false, nil
			}
		}
	}

	return true, nil
}

func (e *WorkflowEngine) AutoAdvance(ctx context.Context, applicationID string) error {
	app, err := e.applicationRepo.GetByID(ctx, applicationID)
	if err != nil {
		return err
	}

	workflows, err := e.workflowRepo.FindByEntityType(ctx, app.CompanyID, domain.WfEntityApplication)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		if !wf.Active {
			continue
		}

		stages, err := e.workflowRepo.ListStages(ctx, wf.ID)
		if err != nil {
			continue
		}

		for _, stage := range stages {
			if stage.StageID == *app.CurrentStageID && stage.AutoAdvance {
				nextStage := e.findNextStage(stages, stage.SortOrder)
				if nextStage != nil {
					app.CurrentStageID = &nextStage.StageID
					app.UpdatedAt = time.Now()
					return e.applicationRepo.Update(ctx, app)
				}
			}
		}
	}

	return nil
}

func (e *WorkflowEngine) ExecuteRule(ctx context.Context, rule domain.WorkflowRule, entity interface{}) error {
	switch rule.ActionType {
	case "send_email":
		return e.executeSendEmail(ctx, rule, entity)
	case "update_status":
		return e.executeUpdateStatus(ctx, rule, entity)
	case "notify":
		return e.executeNotify(ctx, rule, entity)
	case "create_task":
		return e.executeCreateTask(ctx, rule, entity)
	default:
		return nil
	}
}

func (e *WorkflowEngine) executeSendEmail(ctx context.Context, rule domain.WorkflowRule, entity interface{}) error {
	return nil
}

func (e *WorkflowEngine) executeUpdateStatus(ctx context.Context, rule domain.WorkflowRule, entity interface{}) error {
	return nil
}

func (e *WorkflowEngine) executeNotify(ctx context.Context, rule domain.WorkflowRule, entity interface{}) error {
	return nil
}

func (e *WorkflowEngine) executeCreateTask(ctx context.Context, rule domain.WorkflowRule, entity interface{}) error {
	return nil
}

func (e *WorkflowEngine) findNextStage(stages []domain.WorkflowStage, currentOrder int) *domain.WorkflowStage {
	var next *domain.WorkflowStage
	for i := range stages {
		if stages[i].SortOrder > currentOrder {
			if next == nil || stages[i].SortOrder < next.SortOrder {
				next = &stages[i]
			}
		}
	}
	return next
}
