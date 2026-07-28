package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type AuditRepo struct {
	pool *pgxpool.Pool
}

func NewAuditRepo(pool *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{pool: pool}
}

func (r *AuditRepo) Create(ctx context.Context, a *domain.ExpenseAuditLog) error {
	q := `INSERT INTO expense_audit_logs (id,company_id,entity_type,entity_id,action,actor_id,
		changes,ip_address,user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EntityType, a.EntityID, a.Action, a.ActorID,
		a.Changes, a.IPAddress, a.UserAgent)
	return repoErr("AuditRepo.Create", err)
}

func (r *AuditRepo) List(ctx context.Context, companyID uuid.UUID, entityType, entityID, action *string, from, to *time.Time, limit, offset int) ([]domain.ExpenseAuditLog, error) {
	q := `SELECT id,company_id,entity_type,entity_id,action,actor_id,changes,ip_address,user_agent,created_at
		FROM expense_audit_logs WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if entityType != nil {
		q += fmt.Sprintf(" AND entity_type=$%d", n)
		args = append(args, *entityType)
		n++
	}
	if entityID != nil {
		q += fmt.Sprintf(" AND entity_id=$%d", n)
		args = append(args, *entityID)
		n++
	}
	if action != nil {
		q += fmt.Sprintf(" AND action=$%d", n)
		args = append(args, *action)
		n++
	}
	if from != nil {
		q += fmt.Sprintf(" AND created_at>=$%d", n)
		args = append(args, *from)
		n++
	}
	if to != nil {
		q += fmt.Sprintf(" AND created_at<=$%d", n)
		args = append(args, *to)
		n++
	}
	q += " ORDER BY created_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("AuditRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseAuditLog, error) {
		var a domain.ExpenseAuditLog
		err := row.Scan(&a.ID, &a.CompanyID, &a.EntityType, &a.EntityID, &a.Action, &a.ActorID,
			&a.Changes, &a.IPAddress, &a.UserAgent, &a.CreatedAt)
		return a, err
	})
}
