package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type WorkflowRepo struct {
	pool *pgxpool.Pool
}

func NewWorkflowRepo(pool *pgxpool.Pool) *WorkflowRepo {
	return &WorkflowRepo{pool: pool}
}

func (r *WorkflowRepo) CreateWorkflow(ctx context.Context, w *domain.BenefitWorkflow) error {
	q := `INSERT INTO benefit_workflows (id,company_id,benefit_id,workflow_type,name,description,
		requires_chain_approval,auto_approve,auto_approve_if_no_manager,escalation_hours,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, w.ID, w.CompanyID, w.BenefitID, w.WorkflowType, w.Name, w.Description,
		w.RequiresChainApproval, w.AutoApprove, w.AutoApproveIfNoManager, w.EscalationHours, w.IsActive, w.CreatedBy)
	return repoErr("CreateWorkflow", err)
}

func (r *WorkflowRepo) GetWorkflow(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitWorkflow, error) {
	q := `SELECT id,company_id,benefit_id,workflow_type,name,description,requires_chain_approval,auto_approve,
		auto_approve_if_no_manager,escalation_hours,is_active,created_by,created_at,updated_at
		FROM benefit_workflows WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var w domain.BenefitWorkflow
	err := row.Scan(&w.ID, &w.CompanyID, &w.BenefitID, &w.WorkflowType, &w.Name, &w.Description,
		&w.RequiresChainApproval, &w.AutoApprove, &w.AutoApproveIfNoManager, &w.EscalationHours, &w.IsActive, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetWorkflow", err)
	}
	return &w, nil
}

func (r *WorkflowRepo) ListWorkflows(ctx context.Context, companyID uuid.UUID, benefitID *uuid.UUID) ([]domain.BenefitWorkflow, error) {
	q := `SELECT id,company_id,benefit_id,workflow_type,name,description,requires_chain_approval,auto_approve,
		auto_approve_if_no_manager,escalation_hours,is_active,created_by,created_at,updated_at
		FROM benefit_workflows WHERE company_id=$1`
	args := []any{companyID}
	if benefitID != nil {
		q += fmt.Sprintf(" AND benefit_id=$%d", 2)
		args = append(args, *benefitID)
	}
	q += " ORDER BY name"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListWorkflows", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitWorkflow, error) {
		var w domain.BenefitWorkflow
		err := row.Scan(&w.ID, &w.CompanyID, &w.BenefitID, &w.WorkflowType, &w.Name, &w.Description,
			&w.RequiresChainApproval, &w.AutoApprove, &w.AutoApproveIfNoManager, &w.EscalationHours, &w.IsActive, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
		return w, err
	})
}

func (r *WorkflowRepo) UpdateWorkflow(ctx context.Context, w *domain.BenefitWorkflow) error {
	q := `UPDATE benefit_workflows SET benefit_id=$1,workflow_type=$2,name=$3,description=$4,
		requires_chain_approval=$5,auto_approve=$6,auto_approve_if_no_manager=$7,escalation_hours=$8,
		is_active=$9,updated_at=NOW() WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, w.BenefitID, w.WorkflowType, w.Name, w.Description,
		w.RequiresChainApproval, w.AutoApprove, w.AutoApproveIfNoManager, w.EscalationHours,
		w.IsActive, w.ID, w.CompanyID)
	return repoErr("UpdateWorkflow", err)
}

func (r *WorkflowRepo) CreateStep(ctx context.Context, s *domain.BenefitWorkflowStep) error {
	q := `INSERT INTO benefit_workflow_steps (id,workflow_id,step_order,approval_type,approver_role_id,
		max_rejection_count,is_required,notification_template)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.WorkflowID, s.StepOrder, s.ApprovalType, s.ApproverRoleID,
		s.MaxRejectionCount, s.IsRequired, s.NotificationTemplate)
	return repoErr("CreateStep", err)
}

func (r *WorkflowRepo) ListSteps(ctx context.Context, workflowID uuid.UUID) ([]domain.BenefitWorkflowStep, error) {
	q := `SELECT id,workflow_id,step_order,approval_type,approver_role_id,max_rejection_count,is_required,notification_template,created_at
		FROM benefit_workflow_steps WHERE workflow_id=$1 ORDER BY step_order`
	rows, err := r.pool.Query(ctx, q, workflowID)
	if err != nil {
		return nil, repoErr("ListSteps", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitWorkflowStep, error) {
		var s domain.BenefitWorkflowStep
		err := row.Scan(&s.ID, &s.WorkflowID, &s.StepOrder, &s.ApprovalType, &s.ApproverRoleID,
			&s.MaxRejectionCount, &s.IsRequired, &s.NotificationTemplate, &s.CreatedAt)
		return s, err
	})
}

func (r *WorkflowRepo) UpdateStep(ctx context.Context, s *domain.BenefitWorkflowStep) error {
	q := `UPDATE benefit_workflow_steps SET step_order=$1,approval_type=$2,approver_role_id=$3,
		max_rejection_count=$4,is_required=$5,notification_template=$6 WHERE id=$7 AND workflow_id=$8`
	_, err := r.pool.Exec(ctx, q, s.StepOrder, s.ApprovalType, s.ApproverRoleID,
		s.MaxRejectionCount, s.IsRequired, s.NotificationTemplate, s.ID, s.WorkflowID)
	return repoErr("UpdateStep", err)
}

func (r *WorkflowRepo) DeleteStep(ctx context.Context, workflowID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_workflow_steps WHERE id=$1 AND workflow_id=$2`, id, workflowID)
	return repoErr("DeleteStep", err)
}
