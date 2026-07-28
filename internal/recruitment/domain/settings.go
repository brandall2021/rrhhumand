package domain

import "time"

type RecruitmentSource struct {
    ID        string    `json:"id"`
    CompanyID string    `json:"company_id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    Config    *string   `json:"config,omitempty"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}

type RecruitmentStage struct {
    ID        string `json:"id"`
    CompanyID string `json:"company_id"`
    Name      string `json:"name"`
    Category  string `json:"category"`
    SortOrder int    `json:"sort_order"`
    Color     *string `json:"color,omitempty"`
    Active    bool   `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}

type StageTransition struct {
    ID              string   `json:"id"`
    CompanyID       string   `json:"company_id"`
    FromStageID     string   `json:"from_stage_id"`
    ToStageID       string   `json:"to_stage_id"`
    RequiredActions []string `json:"required_actions,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
}

type RejectionReason struct {
    ID        string    `json:"id"`
    CompanyID string    `json:"company_id"`
    Name      string    `json:"name"`
    Category  *string   `json:"category,omitempty"`
    Active    bool      `json:"active"`
    SortOrder int       `json:"sort_order"`
    CreatedAt time.Time `json:"created_at"`
}
