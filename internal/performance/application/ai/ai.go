package ai

import (
	"context"
	"fmt"
)

type SuggestionType string

const (
	SuggestionSummary     SuggestionType = "SUMMARY"
	SuggestionObjective   SuggestionType = "OBJECTIVE"
	SuggestionFeedback    SuggestionType = "FEEDBACK"
	SuggestionEvaluation  SuggestionType = "EVALUATION"
	SuggestionPlan        SuggestionType = "PLAN"
)

type AIService struct {
	enabled bool
}

func NewAIService() *AIService {
	return &AIService{enabled: false}
}

func (s *AIService) IsEnabled() bool {
	return s.enabled
}

func (s *AIService) GenerateSummary(ctx context.Context, data map[string]interface{}) (string, error) {
	if !s.enabled {
		return "", nil
	}
	return "", fmt.Errorf("AI no implementado")
}

func (s *AIService) GenerateSmartSuggestions(ctx context.Context, evalType string, scores map[string]float64) ([]string, error) {
	if !s.enabled {
		return nil, nil
	}
	return nil, fmt.Errorf("AI no implementado")
}

func (s *AIService) GenerateWritingCoach(ctx context.Context, draft string) (string, error) {
	if !s.enabled {
		return draft, nil
	}
	return draft, fmt.Errorf("AI no implementado")
}

func (s *AIService) GenerateObjectiveSuggestions(ctx context.Context, role string) ([]string, error) {
	if !s.enabled {
		return nil, nil
	}
	return nil, fmt.Errorf("AI no implementado")
}

func (s *AIService) GenerateFeedbackSuggestions(ctx context.Context, employeeID string) ([]string, error) {
	if !s.enabled {
		return nil, nil
	}
	return nil, fmt.Errorf("AI no implementado")
}
