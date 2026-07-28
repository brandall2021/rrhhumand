package performance

import (
	"time"
)

type PerformanceCycle struct {
	ID                string     `json:"id"`
	CompanyID         string     `json:"company_id"`
	Name              string     `json:"name"`
	Description       *string    `json:"description,omitempty"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	EvaluationDeadline *time.Time `json:"evaluation_deadline,omitempty"`
	Status            string     `json:"status"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type EvaluationTemplate struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsDefault   bool      `json:"is_default"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type TemplateSection struct {
	ID          string  `json:"id"`
	TemplateID  string  `json:"template_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	SectionType string  `json:"section_type"`
	Weight      float64 `json:"weight"`
	SortOrder   int     `json:"sort_order"`
	Active      bool    `json:"active"`
}

type TemplateSectionItem struct {
	ID          string  `json:"id"`
	SectionID   string  `json:"section_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ItemType    string  `json:"item_type"`
	Weight      float64 `json:"weight"`
	SortOrder   int     `json:"sort_order"`
}

type RatingScale struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	Name       string    `json:"name"`
	MinValue   float64   `json:"min_value"`
	MaxValue   float64   `json:"max_value"`
	Description *string  `json:"description,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type RatingScaleLevel struct {
	ID        string  `json:"id"`
	ScaleID   string  `json:"scale_id"`
	Value     float64 `json:"value"`
	Label     string  `json:"label"`
	Description *string `json:"description,omitempty"`
	SortOrder int     `json:"sort_order"`
}

type Competency struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type ScoringRule struct {
	ID               string  `json:"id"`
	CompanyID        string  `json:"company_id"`
	CycleID          *string `json:"cycle_id,omitempty"`
	ObjectiveWeight  float64 `json:"objective_weight"`
	CompetencyWeight float64 `json:"competency_weight"`
	KPIWeight        float64 `json:"kpi_weight"`
	SelfEvalWeight   float64 `json:"self_eval_weight"`
	ManagerWeight    float64 `json:"manager_weight"`
	PeerWeight       float64 `json:"peer_weight"`
	HRWeight         float64 `json:"hr_weight"`
	Active           bool    `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
}

type PerformanceObjective struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	CycleID      string     `json:"cycle_id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	Metric       *string    `json:"metric,omitempty"`
	TargetValue  *float64   `json:"target_value,omitempty"`
	CurrentValue *float64   `json:"current_value,omitempty"`
	Unit         *string    `json:"unit,omitempty"`
	Weight       float64    `json:"weight"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Status       string     `json:"status"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PerformanceKPI struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"`
	CycleID      string     `json:"cycle_id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	TargetValue  *float64   `json:"target_value,omitempty"`
	CurrentValue *float64   `json:"current_value,omitempty"`
	Unit         *string    `json:"unit,omitempty"`
	Weight       float64    `json:"weight"`
	Status       string     `json:"status"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PerformanceEvaluator struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	CycleID        string     `json:"cycle_id"`
	EmployeeID     string     `json:"employee_id"`
	EmployeeName   string     `json:"employee_name,omitempty"`
	EvaluatorID    string     `json:"evaluator_id"`
	EvaluatorName  string     `json:"evaluator_name,omitempty"`
	EvaluatorType  string     `json:"evaluator_type"`
	Status         string     `json:"status"`
	AssignedAt     time.Time  `json:"assigned_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type PerformanceEvaluation struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	CycleID        string     `json:"cycle_id"`
	CycleName      string     `json:"cycle_name,omitempty"`
	EmployeeID     string     `json:"employee_id"`
	EmployeeName   string     `json:"employee_name,omitempty"`
	EvaluatorID    string     `json:"evaluator_id"`
	EvaluatorName  string     `json:"evaluator_name,omitempty"`
	EvaluatorType  string     `json:"evaluator_type"`
	TemplateID     *string    `json:"template_id,omitempty"`
	Status         string     `json:"status"`
	OverallScore   *float64   `json:"overall_score,omitempty"`
	Comments       *string    `json:"comments,omitempty"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type EvaluationAnswer struct {
	ID           string    `json:"id"`
	EvaluationID string    `json:"evaluation_id"`
	SectionName  *string   `json:"section_name,omitempty"`
	ItemName     string    `json:"item_name"`
	ItemType     string    `json:"item_type"`
	Score        *float64  `json:"score,omitempty"`
	Value        *string   `json:"value,omitempty"`
	Comments     *string   `json:"comments,omitempty"`
	Weight       float64   `json:"weight"`
	CreatedAt    time.Time `json:"created_at"`
}

type PerformanceFeedback struct {
	ID                string     `json:"id"`
	CompanyID         string     `json:"company_id"`
	EmployeeID        string     `json:"employee_id"`
	EmployeeName      string     `json:"employee_name,omitempty"`
	CycleID           *string    `json:"cycle_id,omitempty"`
	FromUserID        string     `json:"from_user_id"`
	FromUserName      string     `json:"from_user_name,omitempty"`
	FeedbackType      string     `json:"feedback_type"`
	Message           string     `json:"message"`
	IsPrivate         bool       `json:"is_private"`
	VisibleToEmployee bool       `json:"visible_to_employee"`
	CreatedAt         time.Time  `json:"created_at"`
}

type PerformanceEvidence struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	EvaluationID   string     `json:"evaluation_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	EvidenceType   string     `json:"evidence_type"`
	StorageProvider *string   `json:"storage_provider,omitempty"`
	StorageKey     *string    `json:"storage_key,omitempty"`
	FileName       *string    `json:"file_name,omitempty"`
	MimeType       *string    `json:"mime_type,omitempty"`
	SizeBytes      *int64     `json:"size_bytes,omitempty"`
	URL            *string    `json:"url,omitempty"`
	CreatedBy      *string    `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type PerformanceResult struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	CycleID         string     `json:"cycle_id"`
	CycleName       string     `json:"cycle_name,omitempty"`
	EmployeeID      string     `json:"employee_id"`
	EmployeeName    string     `json:"employee_name,omitempty"`
	ObjectiveScore  *float64   `json:"objective_score,omitempty"`
	CompetencyScore *float64   `json:"competency_score,omitempty"`
	KPIScore        *float64   `json:"kpi_score,omitempty"`
	SelfScore       *float64   `json:"self_score,omitempty"`
	ManagerScore    *float64   `json:"manager_score,omitempty"`
	PeerScore       *float64   `json:"peer_score,omitempty"`
	HRScore         *float64   `json:"hr_score,omitempty"`
	FinalScore      *float64   `json:"final_score,omitempty"`
	Rating          *string    `json:"rating,omitempty"`
	RatingLabel     *string    `json:"rating_label,omitempty"`
	Strengths       *string    `json:"strengths,omitempty"`
	AreasToImprove  *string    `json:"areas_to_improve,omitempty"`
	CalculatedAt    time.Time  `json:"calculated_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ImprovementPlan struct {
	ID                string     `json:"id"`
	CompanyID         string     `json:"company_id"`
	EmployeeID        string     `json:"employee_id"`
	EmployeeName      string     `json:"employee_name,omitempty"`
	CycleID           *string    `json:"cycle_id,omitempty"`
	ResultID          *string    `json:"result_id,omitempty"`
	Title             string     `json:"title"`
	ProblemDescription *string   `json:"problem_description,omitempty"`
	Objective         *string    `json:"objective,omitempty"`
	ResponsibleID     *string    `json:"responsible_id,omitempty"`
	DueDate           *time.Time `json:"due_date,omitempty"`
	Status            string     `json:"status"`
	Outcome           *string    `json:"outcome,omitempty"`
	CreatedBy         *string    `json:"created_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ImprovementAction struct {
	ID          string     `json:"id"`
	PlanID      string     `json:"plan_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type DevelopmentPlan struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	EmployeeID     string     `json:"employee_id"`
	EmployeeName   string     `json:"employee_name,omitempty"`
	CycleID        *string    `json:"cycle_id,omitempty"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	CareerGoal     *string    `json:"career_goal,omitempty"`
	TimelineMonths int        `json:"timeline_months"`
	Status         string     `json:"status"`
	CreatedBy      *string    `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type DevelopmentAction struct {
	ID          string     `json:"id"`
	PlanID      string     `json:"plan_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	ActionType  string     `json:"action_type"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PerformanceAuditLog struct {
	ID         string     `json:"id"`
	CompanyID  string     `json:"company_id"`
	UserID     *string    `json:"user_id,omitempty"`
	EmployeeID *string    `json:"employee_id,omitempty"`
	EntityType string     `json:"entity_type"`
	EntityID   string     `json:"entity_id"`
	Action     string     `json:"action"`
	OldValue   []byte     `json:"old_value,omitempty"`
	NewValue   []byte     `json:"new_value,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type PerformanceScore struct {
	ObjectiveScore   float64 `json:"objective_score"`
	CompetencyScore  float64 `json:"competency_score"`
	KPIScore         float64 `json:"kpi_score"`
	SelfScore        float64 `json:"self_score"`
	ManagerScore     float64 `json:"manager_score"`
	PeerScore        float64 `json:"peer_score"`
	HRScore          float64 `json:"hr_score"`
	FinalScore       float64 `json:"final_score"`
	Rating           string  `json:"rating"`
	RatingLabel      string  `json:"rating_label"`
	Strengths        string  `json:"strengths"`
	AreasToImprove   string  `json:"areas_to_improve"`
}

type PerformanceDashboard struct {
	TotalCycles          int              `json:"total_cycles"`
	ActiveCycles         int              `json:"active_cycles"`
	TotalEvaluations     int              `json:"total_evaluations"`
	CompletedEvaluations int              `json:"completed_evaluations"`
	PendingEvaluations   int              `json:"pending_evaluations"`
	AverageScore         float64          `json:"average_score"`
	TotalObjectives      int              `json:"total_objectives"`
	CompletedObjectives  int              `json:"completed_objectives"`
	TotalKPIs            int              `json:"total_kpis"`
	TotalFeedback        int              `json:"total_feedback"`
	TotalImprovementPlans int             `json:"total_improvement_plans"`
	RatingDistribution   []RatingCount    `json:"rating_distribution,omitempty"`
}

type RatingCount struct {
	Rating string `json:"rating"`
	Count  int    `json:"count"`
}

type PerformanceFilters struct {
	CycleID     string
	EmployeeID  string
	EvaluatorID string
	Status      string
	Rating      string
	DateFrom    string
	DateTo      string
}
