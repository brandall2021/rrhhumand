package domain

import "time"

type ScoringModel struct {
    ID          string    `json:"id"`
    CompanyID   string    `json:"company_id"`
    Name        string    `json:"name"`
    Description *string   `json:"description,omitempty"`
    Config      *string   `json:"config,omitempty"`
    IsDefault   bool      `json:"is_default"`
    Active      bool      `json:"active"`
    CreatedAt   time.Time `json:"created_at"`
}

type ScoringCriterion struct {
    ID          string   `json:"id"`
    ModelID     string   `json:"model_id"`
    Name        string   `json:"name"`
    Field       string   `json:"field"`
    Weight      float64  `json:"weight"`
    ScoringType string   `json:"scoring_type"`
    Config      *string  `json:"config,omitempty"`
    SortOrder   int      `json:"sort_order"`
}

type MatchingResult struct {
    ID              string   `json:"id"`
    CandidateID     string   `json:"candidate_id"`
    PositionID      string   `json:"position_id"`
    OverallScore    *float64 `json:"overall_score,omitempty"`
    SkillScore      *float64 `json:"skill_score,omitempty"`
    ExperienceScore *float64 `json:"experience_score,omitempty"`
    EducationScore  *float64 `json:"education_score,omitempty"`
    CultureScore    *float64 `json:"culture_score,omitempty"`
    Details         *string  `json:"details,omitempty"`
    MatchedAt       time.Time `json:"matched_at"`
}
