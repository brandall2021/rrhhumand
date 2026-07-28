package surveys

import "time"

type CreateSurveyRequest struct {
	Title             string     `json:"title" validate:"required"`
	Description       *string    `json:"description,omitempty"`
	Type              string     `json:"type" validate:"required"`
	Anonymous         *bool      `json:"anonymous,omitempty"`
	MultipleResponses *bool      `json:"multiple_responses,omitempty"`
	StartsAt          *time.Time `json:"starts_at,omitempty"`
	EndsAt            *time.Time `json:"ends_at,omitempty"`
}

type UpdateSurveyRequest struct {
	Title             *string    `json:"title,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Type              *string    `json:"type,omitempty"`
	Anonymous         *bool      `json:"anonymous,omitempty"`
	MultipleResponses *bool      `json:"multiple_responses,omitempty"`
	StartsAt          *time.Time `json:"starts_at,omitempty"`
	EndsAt            *time.Time `json:"ends_at,omitempty"`
}

type CreateQuestionRequest struct {
	Question string `json:"question" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Position *int   `json:"position,omitempty"`
	Required *bool  `json:"required,omitempty"`
}

type UpdateQuestionRequest struct {
	Question *string `json:"question,omitempty"`
	Type     *string `json:"type,omitempty"`
	Position *int    `json:"position,omitempty"`
	Required *bool   `json:"required,omitempty"`
}

type CreateOptionRequest struct {
	OptionText string `json:"option_text" validate:"required"`
	Position   *int   `json:"position,omitempty"`
}

type SetTargetsRequest struct {
	Targets []TargetItem `json:"targets" validate:"required"`
}

type TargetItem struct {
	TargetType string  `json:"target_type" validate:"required"`
	TargetID   *string `json:"target_id,omitempty"`
}

type RespondSurveyRequest struct {
	Answers []AnswerItem `json:"answers" validate:"required"`
}

type AnswerItem struct {
	QuestionID string   `json:"question_id" validate:"required"`
	Text       *string  `json:"text,omitempty"`
	Number     *float64 `json:"number,omitempty"`
	OptionID   *string  `json:"option_id,omitempty"`
	OptionIDs  []string `json:"option_ids,omitempty"`
}

type SurveyFilters struct {
	Status string
	Type   string
	Search string
}

type SurveyStats struct {
	TotalTargeted   int                    `json:"total_targeted"`
	TotalResponded  int                    `json:"total_responded"`
	ParticipationRate float64              `json:"participation_rate"`
	Questions       []QuestionStats        `json:"questions"`
}

type QuestionStats struct {
	QuestionID   string                 `json:"question_id"`
	Question     string                 `json:"question"`
	Type         string                 `json:"type"`
	TotalAnswers int                    `json:"total_answers"`
	Average      *float64              `json:"average,omitempty"`
	Min          *float64              `json:"min,omitempty"`
	Max          *float64              `json:"max,omitempty"`
	Distribution []OptionDistribution  `json:"distribution,omitempty"`
	YesCount     *int                  `json:"yes_count,omitempty"`
	NoCount      *int                  `json:"no_count,omitempty"`
	YesPercentage *float64             `json:"yes_percentage,omitempty"`
	SampleTexts  []string              `json:"sample_texts,omitempty"`
}

type OptionDistribution struct {
	OptionID   string  `json:"option_id"`
	OptionText string  `json:"option_text"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}
