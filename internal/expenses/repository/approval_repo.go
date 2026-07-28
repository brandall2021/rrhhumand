package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ApprovalRepo struct {
	pool *pgxpool.Pool
}

func NewApprovalRepo(pool *pgxpool.Pool) *ApprovalRepo {
	return &ApprovalRepo{pool: pool}
}

func (r *ApprovalRepo) Create(ctx context.Context, a *domain.ExpenseApproval) error {
	q := `INSERT INTO expense_approvals (id,entity_type,entity_id,step,approver_id,status,comments,assigned_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.EntityType, a.EntityID, a.Step, a.ApproverID, a.Status, a.Comments, a.AssignedAt)
	return repoErr("ApprovalRepo.Create", err)
}

func (r *ApprovalRepo) Get(ctx context.Context, id uuid.UUID) (*domain.ExpenseApproval, error) {
	q := `SELECT id,entity_type,entity_id,step,approver_id,status,comments,assigned_at,decided_at,created_at,updated_at
		FROM expense_approvals WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var a domain.ExpenseApproval
	err := row.Scan(&a.ID, &a.EntityType, &a.EntityID, &a.Step, &a.ApproverID, &a.Status, &a.Comments,
		&a.AssignedAt, &a.DecidedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("ApprovalRepo.Get", err)
	}
	return &a, nil
}

func (r *ApprovalRepo) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.ExpenseApproval, error) {
	q := `SELECT id,entity_type,entity_id,step,approver_id,status,comments,assigned_at,decided_at,created_at,updated_at
		FROM expense_approvals WHERE entity_type=$1 AND entity_id=$2 ORDER BY step`
	rows, err := r.pool.Query(ctx, q, entityType, entityID)
	if err != nil {
		return nil, repoErr("ApprovalRepo.ListByEntity", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseApproval, error) {
		var a domain.ExpenseApproval
		err := row.Scan(&a.ID, &a.EntityType, &a.EntityID, &a.Step, &a.ApproverID, &a.Status, &a.Comments,
			&a.AssignedAt, &a.DecidedAt, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *ApprovalRepo) ListPendingByApprover(ctx context.Context, approverID uuid.UUID) ([]domain.ExpenseApproval, error) {
	q := `SELECT id,entity_type,entity_id,step,approver_id,status,comments,assigned_at,decided_at,created_at,updated_at
		FROM expense_approvals WHERE approver_id=$1 AND status='pending' ORDER BY assigned_at`
	rows, err := r.pool.Query(ctx, q, approverID)
	if err != nil {
		return nil, repoErr("ApprovalRepo.ListPendingByApprover", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseApproval, error) {
		var a domain.ExpenseApproval
		err := row.Scan(&a.ID, &a.EntityType, &a.EntityID, &a.Step, &a.ApproverID, &a.Status, &a.Comments,
			&a.AssignedAt, &a.DecidedAt, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *ApprovalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, comments *string) error {
	q := `UPDATE expense_approvals SET status=$1,comments=$2,decided_at=NOW(),updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, status, comments, id)
	return repoErr("ApprovalRepo.UpdateStatus", err)
}

func (r *ApprovalRepo) GetCurrentStep(ctx context.Context, entityType string, entityID uuid.UUID) (*domain.ExpenseApproval, error) {
	q := `SELECT id,entity_type,entity_id,step,approver_id,status,comments,assigned_at,decided_at,created_at,updated_at
		FROM expense_approvals WHERE entity_type=$1 AND entity_id=$2 AND status='pending' ORDER BY step LIMIT 1`
	row := r.pool.QueryRow(ctx, q, entityType, entityID)
	var a domain.ExpenseApproval
	err := row.Scan(&a.ID, &a.EntityType, &a.EntityID, &a.Step, &a.ApproverID, &a.Status, &a.Comments,
		&a.AssignedAt, &a.DecidedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("ApprovalRepo.GetCurrentStep", err)
	}
	return &a, nil
}
