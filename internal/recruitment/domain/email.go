package domain

import "time"

type EmailTemplate struct {
    ID        string    `json:"id"`
    CompanyID string    `json:"company_id"`
    Name      string    `json:"name"`
    Code      string    `json:"code"`
    Subject   string    `json:"subject"`
    BodyHTML  string    `json:"body_html"`
    BodyText  *string   `json:"body_text,omitempty"`
    Variables []string  `json:"variables,omitempty"`
    Category  *string   `json:"category,omitempty"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type EmailLog struct {
    ID             string    `json:"id"`
    CompanyID      string    `json:"company_id"`
    TemplateID     *string   `json:"template_id,omitempty"`
    ApplicationID  *string   `json:"application_id,omitempty"`
    CandidateID    *string   `json:"candidate_id,omitempty"`
    RecipientEmail string    `json:"recipient_email"`
    Subject        string    `json:"subject"`
    Body           *string   `json:"body,omitempty"`
    Status         string    `json:"status"`
    ErrorMessage   *string   `json:"error_message,omitempty"`
    SentAt         time.Time `json:"sent_at"`
}

type ReferralReward struct {
    ID         string     `json:"id"`
    CompanyID  string     `json:"company_id"`
    ReferralID string     `json:"referral_id"`
    RewardType string     `json:"reward_type"`
    Amount     *float64   `json:"amount,omitempty"`
    Currency   *string    `json:"currency,omitempty"`
    Status     string     `json:"status"`
    PaidAt     *time.Time `json:"paid_at,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
}
