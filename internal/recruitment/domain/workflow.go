package domain

import "time"

type WorkflowEntityType string

const (
    WfEntityApplication  WorkflowEntityType = "APPLICATION"
    WfEntityRequisition  WorkflowEntityType = "REQUISITION"
    WfEntityOffer        WorkflowEntityType = "OFFER"
    WfEntityHiringProcess WorkflowEntityType = "HIRING_PROCESS"
)

type Workflow struct {
    ID          string             `json:"id"`
    CompanyID   string             `json:"company_id"`
    Name        string             `json:"name"`
    Description *string            `json:"description,omitempty"`
    EntityType  WorkflowEntityType `json:"entity_type"`
    IsDefault   bool               `json:"is_default"`
    Active      bool               `json:"active"`
    CreatedAt   time.Time          `json:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at"`
}

type WorkflowStage struct {
    ID                 string    `json:"id"`
    WorkflowID         string    `json:"workflow_id"`
    StageID            string    `json:"stage_id"`
    SortOrder          int       `json:"sort_order"`
    RequiredActions    []string  `json:"required_actions,omitempty"`
    AutoAdvance        bool      `json:"auto_advance"`
    AutoAdvanceDelayH  *int      `json:"auto_advance_delay_hours,omitempty"`
    CreatedAt          time.Time `json:"created_at"`
}

type WorkflowRule struct {
    ID                string    `json:"id"`
    WorkflowID        string    `json:"workflow_id"`
    TriggerEvent      string    `json:"trigger_event"`
    ConditionExpr     *string   `json:"condition_expression,omitempty"`
    ActionType        string    `json:"action_type"`
    ActionConfig      *string   `json:"action_config,omitempty"`
    SortOrder         int       `json:"sort_order"`
    Active            bool      `json:"active"`
    CreatedAt         time.Time `json:"created_at"`
}
