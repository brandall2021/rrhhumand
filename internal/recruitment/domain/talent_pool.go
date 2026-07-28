package domain

import "time"

type TalentPool struct {
    ID          string    `json:"id"`
    CompanyID   string    `json:"company_id"`
    Name        string    `json:"name"`
    Description *string   `json:"description,omitempty"`
    Criteria    *string   `json:"criteria,omitempty"`
    IsAuto      bool      `json:"is_auto"`
    CreatedAt   time.Time `json:"created_at"`
}

type TalentPoolCandidate struct {
    ID          string    `json:"id"`
    PoolID      string    `json:"pool_id"`
    CandidateID string    `json:"candidate_id"`
    AddedBy     *string   `json:"added_by,omitempty"`
    AddedReason *string   `json:"added_reason,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
