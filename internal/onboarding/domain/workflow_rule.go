package domain

type WorkflowType string

const (
	WorkflowOnboarding  WorkflowType = "ONBOARDING"
	WorkflowOffboarding WorkflowType = "OFFBOARDING"
)

type WorkflowRule struct {
	ID           string
	CompanyID    string
	WorkflowType WorkflowType
	Name         string
	Conditions   string
	Actions      string
	Priority     int
	Active       bool
	CreatedAt    string
	UpdatedAt    string
}

type WorkflowCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type WorkflowAction struct {
	Type string      `json:"type"`
	Task string      `json:"task,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
