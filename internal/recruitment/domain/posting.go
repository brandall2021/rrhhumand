package domain

import "time"

type PostingStatus string

const (
    PostStatusDraft     PostingStatus = "DRAFT"
    PostStatusPublished PostingStatus = "PUBLISHED"
    PostStatusClosed    PostingStatus = "CLOSED"
    PostStatusCancelled PostingStatus = "CANCELLED"
)

type Posting struct {
    ID              string        `json:"id"`
    CompanyID       string        `json:"company_id"`
    PositionID      string        `json:"position_id"`
    RequisitionID   *string       `json:"requisition_id,omitempty"`
    Title           string        `json:"title"`
    Description     string        `json:"description"`
    Requirements    *string       `json:"requirements,omitempty"`
    Responsibilities *string      `json:"responsibilities,omitempty"`
    Benefits        *string       `json:"benefits,omitempty"`
    EmploymentType  *string       `json:"employment_type,omitempty"`
    WorkMode        *string       `json:"work_mode,omitempty"`
    Location        *string       `json:"location,omitempty"`
    SalaryMin       *float64      `json:"salary_min,omitempty"`
    SalaryMax       *float64      `json:"salary_max,omitempty"`
    Currency        *string       `json:"currency,omitempty"`
    PublishedAt     *time.Time    `json:"published_at,omitempty"`
    ClosingAt       *time.Time    `json:"closing_at,omitempty"`
    IsPublic        bool          `json:"is_public"`
    ExternalURL     *string       `json:"external_url,omitempty"`
    Status          PostingStatus `json:"status"`
    CreatedAt       time.Time     `json:"created_at"`
    UpdatedAt       time.Time     `json:"updated_at"`
}

type PostingBoard struct {
    ID        string    `json:"id"`
    CompanyID string    `json:"company_id"`
    Name      string    `json:"name"`
    Platform  string    `json:"platform"`
    Config    *string   `json:"config,omitempty"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}

type PostingBoardPost struct {
    ID          string     `json:"id"`
    PostingID   string     `json:"posting_id"`
    BoardID     string     `json:"board_id"`
    ExternalID  *string    `json:"external_id,omitempty"`
    PostedAt    *time.Time `json:"posted_at,omitempty"`
    Status      string     `json:"status"`
    ErrorMsg    *string    `json:"error_message,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}

type PostingScreeningQuestion struct {
    ID          string   `json:"id"`
    PostingID   string   `json:"posting_id"`
    Question    string   `json:"question"`
    QuestionType string  `json:"question_type"`
    Options     *string  `json:"options,omitempty"`
    Required    bool     `json:"required"`
    SortOrder   int      `json:"sort_order"`
    Active      bool     `json:"active"`
    CreatedAt   time.Time `json:"created_at"`
}
