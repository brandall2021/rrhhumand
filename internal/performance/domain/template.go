package domain

import "time"

type PerformanceTemplate struct {
	ID             string         `json:"id"`
	CompanyID      string         `json:"company_id"`
	Name           string         `json:"name"`
	Description    *string        `json:"description,omitempty"`
	EvaluationType EvaluationType `json:"evaluation_type"`
	Active         bool           `json:"active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	Sections  []TemplateSection  `json:"sections,omitempty"`
	Questions []TemplateQuestion `json:"questions,omitempty"`
}

type SectionType string

const (
	SectionTypeScale      SectionType = "SCALE"
	SectionTypeText       SectionType = "TEXT"
	SectionTypeRating     SectionType = "RATING"
	SectionTypeYesNo      SectionType = "YES_NO"
	SectionTypeMultiple   SectionType = "MULTIPLE"
	SectionTypeFreeText   SectionType = "FREE_TEXT"
)

type TemplateSection struct {
	ID          string      `json:"id"`
	TemplateID  string      `json:"template_id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	SectionType SectionType `json:"section_type"`
	Weight      float64     `json:"weight"`
	SortOrder   int         `json:"sort_order"`
	Active      bool        `json:"active"`
}

type QuestionType string

const (
	QuestionTypeScale    QuestionType = "SCALE"
	QuestionTypeRating   QuestionType = "RATING"
	QuestionTypeText     QuestionType = "TEXT"
	QuestionTypeYesNo    QuestionType = "YES_NO"
	QuestionTypeMultiple QuestionType = "MULTIPLE"
	QuestionTypeBoolean  QuestionType = "BOOLEAN"
)

type TemplateQuestion struct {
	ID           string       `json:"id"`
	TemplateID   string       `json:"template_id"`
	SectionID    *string      `json:"section_id,omitempty"`
	Question     string       `json:"question"`
	QuestionType QuestionType `json:"question_type"`
	Required     bool         `json:"required"`
	Weight       float64      `json:"weight"`
	SortOrder    int          `json:"sort_order"`
	Active       bool         `json:"active"`
}

type RatingScale struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	MinValue    float64   `json:"min_value"`
	MaxValue    float64   `json:"max_value"`
	Description *string   `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Levels []RatingScaleLevel `json:"levels,omitempty"`
}

type RatingScaleLevel struct {
	ID          string  `json:"id"`
	ScaleID     string  `json:"scale_id"`
	Value       float64 `json:"value"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	SortOrder   int     `json:"sort_order"`
}

type OutboxEvent struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	EventType     string    `json:"event_type"`
	AggregateType string    `json:"aggregate_type"`
	AggregateID   string    `json:"aggregate_id"`
	Payload       []byte    `json:"payload,omitempty"`
	Status        string    `json:"status"`
	RetryCount    int       `json:"retry_count"`
	LastError     *string   `json:"last_error,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
}
