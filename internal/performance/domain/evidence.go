package domain

import "time"

type EvidenceType string

const (
	EvidenceTypeFile    EvidenceType = "FILE"
	EvidenceTypeLink    EvidenceType = "LINK"
	EvidenceTypeComment EvidenceType = "COMMENT"
	EvidenceTypeImage   EvidenceType = "IMAGE"
	EvidenceTypeDocument EvidenceType = "DOCUMENT"
)

type PerformanceEvidence struct {
	ID           string       `json:"id"`
	CompanyID    string       `json:"company_id"`
	EvaluationID *string      `json:"evaluation_id,omitempty"`
	ObjectiveID  *string      `json:"objective_id,omitempty"`
	FeedbackID   *string      `json:"feedback_id,omitempty"`
	Title        string       `json:"title"`
	Description  *string      `json:"description,omitempty"`
	EvidenceType EvidenceType `json:"evidence_type"`
	StorageKey   *string      `json:"storage_key,omitempty"`
	FileName     *string      `json:"file_name,omitempty"`
	MimeType     *string      `json:"mime_type,omitempty"`
	SizeBytes    *int64       `json:"size_bytes,omitempty"`
	URL          *string      `json:"url,omitempty"`
	CreatedBy    string       `json:"created_by"`
	CreatedAt    time.Time    `json:"created_at"`
}

type PerformanceResult struct {
	ID               string    `json:"id"`
	CompanyID        string    `json:"company_id"`
	CycleID          string    `json:"cycle_id"`
	EmployeeID       string    `json:"employee_id"`
	ObjectiveScore   *float64  `json:"objective_score,omitempty"`
	CompetencyScore  *float64  `json:"competency_score,omitempty"`
	SelfScore        *float64  `json:"self_score,omitempty"`
	ManagerScore     *float64  `json:"manager_score,omitempty"`
	PeerScore        *float64  `json:"peer_score,omitempty"`
	HRScore          *float64  `json:"hr_score,omitempty"`
	FinalScore       *float64  `json:"final_score,omitempty"`
	FinalRating      *string   `json:"final_rating,omitempty"`
	FinalRatingLabel *string   `json:"final_rating_label,omitempty"`
	Strengths        *string   `json:"strengths,omitempty"`
	ImprovementAreas *string   `json:"improvement_areas,omitempty"`
	Summary          *string   `json:"summary,omitempty"`
	CalculatedAt     time.Time `json:"calculated_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PerformanceDashboard struct {
	TotalCycles           int            `json:"total_cycles"`
	ActiveCycles          int            `json:"active_cycles"`
	TotalEvaluations      int            `json:"total_evaluations"`
	CompletedEvaluations  int            `json:"completed_evaluations"`
	PendingEvaluations    int            `json:"pending_evaluations"`
	AverageScore          float64        `json:"average_score"`
	TotalObjectives       int            `json:"total_objectives"`
	CompletedObjectives   int            `json:"completed_objectives"`
	TotalFeedback         int            `json:"total_feedback"`
	TotalImprovementPlans int            `json:"total_improvement_plans"`
	TotalDevelopmentPlans int            `json:"total_development_plans"`
	RatingDistribution    []RatingCount  `json:"rating_distribution,omitempty"`
}

type RatingCount struct {
	Rating string `json:"rating"`
	Count  int    `json:"count"`
}

type PerformanceAuditLog struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	UserID     *string   `json:"user_id,omitempty"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	OldValues  []byte    `json:"old_values,omitempty"`
	NewValues  []byte    `json:"new_values,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
