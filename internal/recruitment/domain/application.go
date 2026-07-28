package domain

import "time"

type ApplicationStatus string

const (
    AppStatusNew         ApplicationStatus = "NEW"
    AppStatusScreening   ApplicationStatus = "SCREENING"
    AppStatusInterview   ApplicationStatus = "INTERVIEW"
    AppStatusAssessment  ApplicationStatus = "ASSESSMENT"
    AppStatusOffer       ApplicationStatus = "OFFER"
    AppStatusHired       ApplicationStatus = "HIRED"
    AppStatusRejected    ApplicationStatus = "REJECTED"
    AppStatusWithdrawn   ApplicationStatus = "WITHDRAWN"
    AppStatusOnHold      ApplicationStatus = "ON_HOLD"
)

type Application struct {
    ID                 string            `json:"id"`
    CompanyID          string            `json:"company_id"`
    CandidateID        string            `json:"candidate_id"`
    PostingID          string            `json:"posting_id"`
    CurrentStageID     *string           `json:"current_stage_id,omitempty"`
    Status             ApplicationStatus `json:"status"`
    Score              *float64          `json:"score,omitempty"`
    AppliedAt          time.Time         `json:"applied_at"`
    ReviewedAt         *time.Time        `json:"reviewed_at,omitempty"`
    RejectedAt         *time.Time        `json:"rejected_at,omitempty"`
    RejectionReasonID  *string           `json:"rejection_reason_id,omitempty"`
    RejectionText      *string           `json:"rejection_reason_text,omitempty"`
    HiredAt            *time.Time        `json:"hired_at,omitempty"`
    WithdrawnAt        *time.Time        `json:"withdrawn_at,omitempty"`
    WithdrawReason     *string           `json:"withdraw_reason,omitempty"`
    Source             *string           `json:"source,omitempty"`
    SourceDetail       *string           `json:"source_detail,omitempty"`
    IsInternalMobility bool              `json:"is_internal_mobility"`
    ConsentGiven       bool              `json:"consent_given"`
    ConsentAt          *time.Time        `json:"consent_at,omitempty"`
    Notes              *string           `json:"notes,omitempty"`
    CreatedAt          time.Time         `json:"created_at"`
    UpdatedAt          time.Time         `json:"updated_at"`
}

type ApplicationScreeningAnswer struct {
    ID            string    `json:"id"`
    ApplicationID string    `json:"application_id"`
    QuestionID    string    `json:"question_id"`
    Answer        string    `json:"answer"`
    CreatedAt     time.Time `json:"created_at"`
}

type ApplicationStageHistory struct {
    ID            string    `json:"id"`
    ApplicationID string    `json:"application_id"`
    FromStageID   *string   `json:"from_stage_id,omitempty"`
    ToStageID     string    `json:"to_stage_id"`
    ChangedBy     *string   `json:"changed_by,omitempty"`
    Reason        *string   `json:"reason,omitempty"`
    AutoTransition bool     `json:"auto_transition"`
    CreatedAt     time.Time `json:"created_at"`
}

type ApplicationRating struct {
    ID            string    `json:"id"`
    ApplicationID string    `json:"application_id"`
    RatedBy       string    `json:"rated_by"`
    Rating        int       `json:"rating"`
    Comment       *string   `json:"comment,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
}

type ApplicationNote struct {
    ID            string    `json:"id"`
    ApplicationID string    `json:"application_id"`
    AuthorID      string    `json:"author_id"`
    Content       string    `json:"content"`
    IsPrivate     bool      `json:"is_private"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
