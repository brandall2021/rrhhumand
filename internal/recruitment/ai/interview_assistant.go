package ai

import (
	"context"
	"math/rand"
)

type InterviewAssistant struct{}

func NewInterviewAssistant() *InterviewAssistant {
	return &InterviewAssistant{}
}

// TODO: implement AI integration
func (a *InterviewAssistant) GenerateQuestions(ctx context.Context, applicationID, interviewType string) ([]string, error) {
	questions := map[string][]string{
		"TECHNICAL": {
			"Describe your experience with the primary technologies required for this role.",
			"How do you approach debugging complex issues?",
			"Explain a time you had to make a technical trade-off decision.",
			"What is your experience with system design and architecture?",
			"How do you stay current with industry trends?",
		},
		"BEHAVIORAL": {
			"Tell me about a time you handled a conflict in a team.",
			"Describe a situation where you went above and beyond for a project.",
			"How do you handle constructive criticism?",
			"Tell me about a time you failed and what you learned.",
			"Describe your ideal work environment.",
		},
		"HR": {
			"Why are you interested in this position?",
			"Where do you see yourself in 5 years?",
			"What are your salary expectations?",
			"Why do you want to leave your current position?",
			"Do you have any questions about the company?",
		},
		"CULTURE": {
			"What does a great company culture look like to you?",
			"How do you contribute to team culture?",
			"Describe your preferred working style.",
			"How do you maintain work-life balance?",
			"What values are most important to you in a workplace?",
		},
	}

	if qs, ok := questions[interviewType]; ok {
		return qs, nil
	}
	return questions["TECHNICAL"], nil
}

// TODO: implement AI integration
func (a *InterviewAssistant) SuggestScore(ctx context.Context, feedbackText string) (float64, error) {
	return rand.Float64()*5 + 1, nil
}

// TODO: implement AI integration
func (a *InterviewAssistant) AnalyzeResponse(ctx context.Context, question, response string) (string, error) {
	return "Response analysis not available in stub mode.", nil
}

// TODO: implement AI integration
func (a *InterviewAssistant) GenerateFeedbackSummary(ctx context.Context, feedbacks []string) (string, error) {
	return "Feedback summary not available in stub mode.", nil
}
