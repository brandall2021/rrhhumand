package events

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Emit(ctx context.Context, eventType, resourceID, description string) error {
	return nil
}
