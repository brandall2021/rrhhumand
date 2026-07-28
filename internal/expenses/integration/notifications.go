package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationAdapter struct {
	pool *pgxpool.Pool
}

func NewNotificationAdapter(pool *pgxpool.Pool) *NotificationAdapter {
	return &NotificationAdapter{pool: pool}
}

func (a *NotificationAdapter) SendNotification(ctx context.Context, companyID, employeeID *uuid.UUID, notifType, channel, title, body string, metadata map[string]any) error {
	_, err := a.pool.Exec(ctx, `
		INSERT INTO expense_notification_log (id,company_id,employee_id,notification_type,channel,title,body,metadata,sent_at,created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())`,
		uuid.New(), companyID, employeeID, notifType, channel, title, body, metadata)
	return integErr("SendNotification", err)
}

func (a *NotificationAdapter) NotifyApprover(ctx context.Context, approverID uuid.UUID, entityType string, entityID uuid.UUID) error {
	title := "Approval Request"
	body := "You have a pending " + entityType + " requiring your approval."
	metadata := map[string]any{"entity_type": entityType, "entity_id": entityID.String()}
	return a.SendNotification(ctx, nil, &approverID, "APPROVAL_REQUEST", "in_app", title, body, metadata)
}

func (a *NotificationAdapter) NotifyEmployee(ctx context.Context, employeeID uuid.UUID, notifType, title, body string) error {
	return a.SendNotification(ctx, nil, &employeeID, notifType, "in_app", title, body, nil)
}
