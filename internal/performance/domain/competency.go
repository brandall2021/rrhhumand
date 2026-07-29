package domain

import "time"

type CompetencyType string

const (
	CompetencyTypeBehavioral CompetencyType = "BEHAVIORAL"
	CompetencyTypeTechnical  CompetencyType = "TECHNICAL"
	CompetencyTypeLeadership CompetencyType = "LEADERSHIP"
	CompetencyTypeManagerial CompetencyType = "MANAGERIAL"
)

type Competency struct {
	ID              string         `json:"id"`
	CompanyID       string         `json:"company_id"`
	Name            string         `json:"name"`
	Description     *string        `json:"description,omitempty"`
	Category        *string        `json:"category,omitempty"`
	CompetencyType  CompetencyType `json:"competency_type"`
	Active          bool           `json:"active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	Levels []CompetencyLevel `json:"levels,omitempty"`
}

type CompetencyLevel struct {
	ID           string  `json:"id"`
	CompetencyID string  `json:"competency_id"`
	Level        int     `json:"level"`
	Label        string  `json:"label"`
	Description  *string `json:"description,omitempty"`
	SortOrder    int     `json:"sort_order"`
}

type PositionCompetency struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	PositionID   string  `json:"position_id"`
	CompetencyID string  `json:"competency_id"`
	ExpectedLevel int    `json:"expected_level"`
	Weight       float64 `json:"weight"`
}

type CycleCompetency struct {
	ID           string  `json:"id"`
	CycleID      string  `json:"cycle_id"`
	EmployeeID   string  `json:"employee_id"`
	CompetencyID string  `json:"competency_id"`
	ExpectedLevel int    `json:"expected_level"`
	Weight       float64 `json:"weight"`
}

type CompetencyFilter struct {
	CompanyID string
	Category  string
	Type      CompetencyType
	Active    *bool
}
