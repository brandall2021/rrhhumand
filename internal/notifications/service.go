package notifications

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, companyID, userID, title, message, notificationType string) error {
	return nil
}
