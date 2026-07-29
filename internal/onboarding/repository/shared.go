package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/onboarding/domain"
)

type SharedRepo struct {
	pool *pgxpool.Pool
}

func NewSharedRepo(pool *pgxpool.Pool) *SharedRepo {
	return &SharedRepo{pool: pool}
}

func (r *SharedRepo) CreateWorkflowRule(ctx context.Context, rule *domain.WorkflowRule) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO workflow_rules (company_id, workflow_type, name, conditions, actions, priority, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`,
		rule.CompanyID, rule.WorkflowType, rule.Name, rule.Conditions, rule.Actions, rule.Priority, rule.Active,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *SharedRepo) ListWorkflowRules(ctx context.Context, companyID string, workflowType domain.WorkflowType) ([]domain.WorkflowRule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, workflow_type, name, conditions, actions, priority, active, created_at, updated_at
		 FROM workflow_rules WHERE company_id=$1 AND workflow_type=$2 AND active=true ORDER BY priority`, companyID, workflowType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rs []domain.WorkflowRule
	for rows.Next() {
		var rule domain.WorkflowRule
		if err := rows.Scan(&rule.ID, &rule.CompanyID, &rule.WorkflowType, &rule.Name,
			&rule.Conditions, &rule.Actions, &rule.Priority, &rule.Active,
			&rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, err
		}
		rs = append(rs, rule)
	}
	return rs, nil
}

func (r *SharedRepo) CreateOutboxEvent(ctx context.Context, e *domain.OutboxEvent) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO outbox_events (company_id, event_type, aggregate_type, aggregate_id, payload, status)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`,
		e.CompanyID, e.EventType, e.AggregateType, e.AggregateID, e.Payload, domain.OutboxPending,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *SharedRepo) GetPendingOutboxEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, event_type, aggregate_type, aggregate_id, payload, status, retry_count, last_error, created_at, processed_at
		 FROM outbox_events WHERE status='PENDING' ORDER BY created_at ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var es []domain.OutboxEvent
	for rows.Next() {
		var e domain.OutboxEvent
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EventType, &e.AggregateType, &e.AggregateID,
			&e.Payload, &e.Status, &e.RetryCount, &e.LastError, &e.CreatedAt, &e.ProcessedAt); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

func (r *SharedRepo) MarkOutboxProcessed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE outbox_events SET status='PROCESSED', processed_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *SharedRepo) MarkOutboxFailed(ctx context.Context, id string, errMsg string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE outbox_events SET status='FAILED', retry_count=retry_count+1, last_error=$2 WHERE id=$1`,
		id, errMsg)
	return err
}

func (r *SharedRepo) CreateAuditLog(ctx context.Context, companyID, userID, action, entityType, entityID, ipAddress, userAgent string, oldVal, newVal interface{}) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO onboarding_audit_log (company_id, user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		companyID, userID, action, entityType, entityID, oldVal, newVal, ipAddress, userAgent)
	return err
}

func (r *SharedRepo) CreateOffboardingAuditLog(ctx context.Context, companyID, userID, action, entityType, entityID, ipAddress, userAgent string, oldVal, newVal interface{}) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO offboarding_audit_log (company_id, user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		companyID, userID, action, entityType, entityID, oldVal, newVal, ipAddress, userAgent)
	return err
}

func (r *SharedRepo) CreateNotification(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications (company_id, user_id, title, body, notification_type, channel, reference_type, reference_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		companyID, userID, title, body, notifType, "IN_APP", refType, refID)
	return err
}
