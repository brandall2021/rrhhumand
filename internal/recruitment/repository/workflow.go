package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type WorkflowRepo struct {
	pool *pgxpool.Pool
}

func NewWorkflowRepo(pool *pgxpool.Pool) *WorkflowRepo {
	return &WorkflowRepo{pool: pool}
}

func (r *WorkflowRepo) Create(ctx context.Context, companyID string, req *domain.Workflow) (*domain.Workflow, error) {
	w := &domain.Workflow{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_workflows (company_id, name, description, entity_type, is_default, active)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, name, description, entity_type, is_default, active, created_at, updated_at`,
		companyID, req.Name, req.Description, req.EntityType, req.IsDefault, req.Active,
	).Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.EntityType, &w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (r *WorkflowRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Workflow, error) {
	w := &domain.Workflow{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, entity_type, is_default, active, created_at, updated_at
		 FROM recruitment_workflows WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.EntityType, &w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (r *WorkflowRepo) List(ctx context.Context, companyID string, entityType string) ([]domain.Workflow, error) {
	query := `SELECT id, company_id, name, description, entity_type, is_default, active, created_at, updated_at
		 FROM recruitment_workflows WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if entityType != "" {
		query += fmt.Sprintf(" AND entity_type=$%d", argIdx)
		args = append(args, entityType)
		argIdx++
	}
	query += " ORDER BY name"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []domain.Workflow
	for rows.Next() {
		var w domain.Workflow
		rows.Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.EntityType, &w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt)
		workflows = append(workflows, w)
	}
	return workflows, nil
}

func (r *WorkflowRepo) Update(ctx context.Context, companyID, id string, req *domain.Workflow) (*domain.Workflow, error) {
	w := &domain.Workflow{}
	err := r.pool.QueryRow(ctx,
		`UPDATE recruitment_workflows SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 entity_type=COALESCE($5,entity_type), is_default=COALESCE($6,is_default),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, entity_type, is_default, active, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.EntityType, req.IsDefault,
	).Scan(&w.ID, &w.CompanyID, &w.Name, &w.Description, &w.EntityType, &w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (r *WorkflowRepo) Activate(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE recruitment_workflows SET active=TRUE, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *WorkflowRepo) Deactivate(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE recruitment_workflows SET active=FALSE, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *WorkflowRepo) AddStage(ctx context.Context, req *domain.WorkflowStage) (*domain.WorkflowStage, error) {
	s := &domain.WorkflowStage{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_workflow_stages (workflow_id, stage_id, sort_order, required_actions, auto_advance, auto_advance_delay_hours)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, workflow_id, stage_id, sort_order, required_actions, auto_advance, auto_advance_delay_hours, created_at`,
		req.WorkflowID, req.StageID, req.SortOrder, req.RequiredActions, req.AutoAdvance, req.AutoAdvanceDelayH,
	).Scan(&s.ID, &s.WorkflowID, &s.StageID, &s.SortOrder, &s.RequiredActions, &s.AutoAdvance, &s.AutoAdvanceDelayH, &s.CreatedAt)
	return s, err
}

func (r *WorkflowRepo) RemoveStage(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM recruitment_workflow_stages WHERE id=$1`, id)
	return err
}

func (r *WorkflowRepo) ListStages(ctx context.Context, workflowID string) ([]domain.WorkflowStage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workflow_id, stage_id, sort_order, required_actions, auto_advance, auto_advance_delay_hours, created_at
		 FROM recruitment_workflow_stages WHERE workflow_id=$1 ORDER BY sort_order`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []domain.WorkflowStage
	for rows.Next() {
		var s domain.WorkflowStage
		rows.Scan(&s.ID, &s.WorkflowID, &s.StageID, &s.SortOrder, &s.RequiredActions, &s.AutoAdvance, &s.AutoAdvanceDelayH, &s.CreatedAt)
		stages = append(stages, s)
	}
	return stages, nil
}

func (r *WorkflowRepo) ReorderStages(ctx context.Context, stageIDs []string) error {
	for i, id := range stageIDs {
		_, err := r.pool.Exec(ctx,
			`UPDATE recruitment_workflow_stages SET sort_order=$2 WHERE id=$1`,
			id, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *WorkflowRepo) AddRule(ctx context.Context, req *domain.WorkflowRule) (*domain.WorkflowRule, error) {
	rl := &domain.WorkflowRule{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_workflow_rules (workflow_id, trigger_event, condition_expression, action_type, action_config, sort_order, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, workflow_id, trigger_event, condition_expression, action_type, action_config, sort_order, active, created_at`,
		req.WorkflowID, req.TriggerEvent, req.ConditionExpr, req.ActionType, req.ActionConfig, req.SortOrder, req.Active,
	).Scan(&rl.ID, &rl.WorkflowID, &rl.TriggerEvent, &rl.ConditionExpr, &rl.ActionType, &rl.ActionConfig, &rl.SortOrder, &rl.Active, &rl.CreatedAt)
	return rl, err
}

func (r *WorkflowRepo) UpdateRule(ctx context.Context, id string, req *domain.WorkflowRule) (*domain.WorkflowRule, error) {
	rl := &domain.WorkflowRule{}
	err := r.pool.QueryRow(ctx,
		`UPDATE recruitment_workflow_rules SET
		 trigger_event=COALESCE($2,trigger_event), condition_expression=COALESCE($3,condition_expression),
		 action_type=COALESCE($4,action_type), action_config=COALESCE($5,action_config),
		 sort_order=COALESCE($6,sort_order), active=COALESCE($7,active) WHERE id=$1
		 RETURNING id, workflow_id, trigger_event, condition_expression, action_type, action_config, sort_order, active, created_at`,
		id, req.TriggerEvent, req.ConditionExpr, req.ActionType, req.ActionConfig, req.SortOrder, req.Active,
	).Scan(&rl.ID, &rl.WorkflowID, &rl.TriggerEvent, &rl.ConditionExpr, &rl.ActionType, &rl.ActionConfig, &rl.SortOrder, &rl.Active, &rl.CreatedAt)
	return rl, err
}

func (r *WorkflowRepo) DeleteRule(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM recruitment_workflow_rules WHERE id=$1`, id)
	return err
}

func (r *WorkflowRepo) ListRules(ctx context.Context, workflowID string) ([]domain.WorkflowRule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workflow_id, trigger_event, condition_expression, action_type, action_config, sort_order, active, created_at
		 FROM recruitment_workflow_rules WHERE workflow_id=$1 ORDER BY sort_order`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []domain.WorkflowRule
	for rows.Next() {
		var rl domain.WorkflowRule
		rows.Scan(&rl.ID, &rl.WorkflowID, &rl.TriggerEvent, &rl.ConditionExpr, &rl.ActionType, &rl.ActionConfig, &rl.SortOrder, &rl.Active, &rl.CreatedAt)
		rules = append(rules, rl)
	}
	return rules, nil
}
