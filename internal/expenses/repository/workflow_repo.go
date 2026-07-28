package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/shopspring/decimal"
)

type WorkflowRepo struct {
	pool *pgxpool.Pool
}

func NewWorkflowRepo(pool *pgxpool.Pool) *WorkflowRepo {
	return &WorkflowRepo{pool: pool}
}

func (r *WorkflowRepo) CreateWorkflow(ctx context.Context, w *domain.ExpenseWorkflow) error {
	q := `INSERT INTO expense_workflows (id,company_id,name,description,workflow_type,
		min_amount,max_amount,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, w.ID, w.CompanyID, w.Name, w.Description, w.WorkflowType,
		w.MinAmount, w.MaxAmount, w.IsActive, w.CreatedBy)
	return repoErr("CreateWorkflow", err)
}

func (r *WorkflowRepo) GetWorkflow(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseWorkflow, error) {
	q := `SELECT id,company_id,name,description,workflow_type,
		min_amount,max_amount,is_active,created_by,created_at,updated_at
		FROM expense_workflows WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var w domain.ExpenseWorkflow
	err := row.Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.WorkflowType,
		&w.MinAmount, &w.MaxAmount, &w.IsActive, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetWorkflow", err)
	}
	return &w, nil
}

func (r *WorkflowRepo) ListWorkflows(ctx context.Context, companyID uuid.UUID) ([]domain.ExpenseWorkflow, error) {
	q := `SELECT id,company_id,name,description,workflow_type,
		min_amount,max_amount,is_active,created_by,created_at,updated_at
		FROM expense_workflows WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListWorkflows", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseWorkflow, error) {
		var w domain.ExpenseWorkflow
		err := row.Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.WorkflowType,
			&w.MinAmount, &w.MaxAmount, &w.IsActive, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
		return w, err
	})
}

func (r *WorkflowRepo) UpdateWorkflow(ctx context.Context, w *domain.ExpenseWorkflow) error {
	q := `UPDATE expense_workflows SET name=$1,description=$2,workflow_type=$3,
		min_amount=$4,max_amount=$5,is_active=$6,updated_at=NOW() WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, q, w.Name, w.Description, w.WorkflowType,
		w.MinAmount, w.MaxAmount, w.IsActive, w.ID, w.CompanyID)
	return repoErr("UpdateWorkflow", err)
}

func (r *WorkflowRepo) CreateStep(ctx context.Context, s *domain.ExpenseWorkflowStep) error {
	q := `INSERT INTO expense_workflow_steps (id,workflow_id,step_order,approver_type,approver_id,
		role_name,min_amount,max_amount,required_approvals)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.WorkflowID, s.StepOrder, s.ApproverType, s.ApproverID,
		s.RoleName, s.MinAmount, s.MaxAmount, s.RequiredApprovals)
	return repoErr("CreateStep", err)
}

func (r *WorkflowRepo) ListSteps(ctx context.Context, workflowID uuid.UUID) ([]domain.ExpenseWorkflowStep, error) {
	q := `SELECT id,workflow_id,step_order,approver_type,approver_id,role_name,
		min_amount,max_amount,required_approvals,created_at,updated_at
		FROM expense_workflow_steps WHERE workflow_id=$1 ORDER BY step_order`
	rows, err := r.pool.Query(ctx, q, workflowID)
	if err != nil {
		return nil, repoErr("ListSteps", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseWorkflowStep, error) {
		var s domain.ExpenseWorkflowStep
		err := row.Scan(&s.ID, &s.WorkflowID, &s.StepOrder, &s.ApproverType, &s.ApproverID, &s.RoleName,
			&s.MinAmount, &s.MaxAmount, &s.RequiredApprovals, &s.CreatedAt, &s.UpdatedAt)
		return s, err
	})
}

func (r *WorkflowRepo) UpdateStep(ctx context.Context, s *domain.ExpenseWorkflowStep) error {
	q := `UPDATE expense_workflow_steps SET step_order=$1,approver_type=$2,approver_id=$3,
		role_name=$4,min_amount=$5,max_amount=$6,required_approvals=$7,updated_at=NOW()
		WHERE id=$8 AND workflow_id=$9`
	_, err := r.pool.Exec(ctx, q, s.StepOrder, s.ApproverType, s.ApproverID, s.RoleName,
		s.MinAmount, s.MaxAmount, s.RequiredApprovals, s.ID, s.WorkflowID)
	return repoErr("UpdateStep", err)
}

func (r *WorkflowRepo) DeleteStep(ctx context.Context, workflowID, stepID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_workflow_steps WHERE id=$1 AND workflow_id=$2`, stepID, workflowID)
	return repoErr("DeleteStep", err)
}

func (r *WorkflowRepo) FindWorkflow(ctx context.Context, companyID uuid.UUID, workflowType string, amount decimal.Decimal) (*domain.ExpenseWorkflow, error) {
	q := `SELECT id,company_id,name,description,workflow_type,
		min_amount,max_amount,is_active,created_by,created_at,updated_at
		FROM expense_workflows
		WHERE company_id=$1 AND workflow_type=$2 AND is_active=true
		AND (min_amount IS NULL OR min_amount<=$3)
		AND (max_amount IS NULL OR max_amount>=$3)
		ORDER BY min_amount NULLS FIRST LIMIT 1`
	row := r.pool.QueryRow(ctx, q, companyID, workflowType, amount)
	var w domain.ExpenseWorkflow
	err := row.Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.WorkflowType,
		&w.MinAmount, &w.MaxAmount, &w.IsActive, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, repoErr("FindWorkflow", err)
	}
	return &w, nil
}
