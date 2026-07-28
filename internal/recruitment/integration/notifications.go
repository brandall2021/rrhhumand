package integration

import (
	"context"
	"log"
)

type NotificationAdapter struct{}

func NewNotificationAdapter() *NotificationAdapter {
	return &NotificationAdapter{}
}

func (a *NotificationAdapter) SendNotification(ctx context.Context, userID, title, body string) error {
	log.Printf("[NOTIFICATION] user=%s title=%s body=%s", userID, title, body)
	return nil
}

func (a *NotificationAdapter) SendEmailNotification(ctx context.Context, email, subject, body string) error {
	log.Printf("[EMAIL] to=%s subject=%s body=%s", email, subject, body)
	return nil
}

func (a *NotificationAdapter) SendBatchNotification(ctx context.Context, userIDs []string, title, body string) error {
	log.Printf("[BATCH NOTIFICATION] users=%v title=%s body=%s", userIDs, title, body)
	return nil
}

func (a *NotificationAdapter) SendTemplateNotification(ctx context.Context, userID, templateName string, vars map[string]string) error {
	log.Printf("[TEMPLATE NOTIFICATION] user=%s template=%s vars=%v", userID, templateName, vars)
	return nil
}
