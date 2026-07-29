package integration

import (
	"context"
	"log"
)

type NotificationAdapter struct{}

func NewNotificationAdapter() *NotificationAdapter {
	return &NotificationAdapter{}
}

func (a *NotificationAdapter) Send(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) error {
	log.Printf("[NotificationAdapter] Send company=%s user=%s title=%s type=%s", companyID, userID, title, notifType)
	return nil
}

func (a *NotificationAdapter) SendToRole(ctx context.Context, companyID, role, title, body, notifType, refType, refID string) error {
	log.Printf("[NotificationAdapter] SendToRole company=%s role=%s title=%s", companyID, role, title)
	return nil
}
