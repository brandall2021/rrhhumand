package performance

type CreateCycleRequest struct {
	Name               string  `json:"name" binding:"required"`
	Description        *string `json:"description"`
	StartDate          string  `json:"start_date" binding:"required"`
	EndDate            string  `json:"end_date" binding:"required"`
	EvaluationDeadline *string `json:"evaluation_deadline"`
}

type UpdateCycleRequest struct {
	Name               *string `json:"name"`
	Description        *string `json:"description"`
	StartDate          *string `json:"start_date"`
	EndDate            *string `json:"end_date"`
	EvaluationDeadline *string `json:"evaluation_deadline"`
}

type CreateTemplateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsDefault   *bool   `json:"is_default"`
	Sections    []CreateSectionRequest `json:"sections"`
}

type CreateSectionRequest struct {
	Name        string                `json:"name" binding:"required"`
	Description *string               `json:"description"`
	SectionType string                `json:"section_type" binding:"required"`
	Weight      *float64              `json:"weight"`
	Items       []CreateSectionItemRequest `json:"items"`
}

type CreateSectionItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	ItemType    *string `json:"item_type"`
	Weight      *float64 `json:"weight"`
}

type CreateScaleRequest struct {
	Name        string                `json:"name" binding:"required"`
	MinValue    float64               `json:"min_value" binding:"required"`
	MaxValue    float64               `json:"max_value" binding:"required"`
	Description *string               `json:"description"`
	Levels      []CreateScaleLevelRequest `json:"levels"`
}

type CreateScaleLevelRequest struct {
	Value       float64 `json:"value" binding:"required"`
	Label       string  `json:"label" binding:"required"`
	Description *string `json:"description"`
}

type CreateCompetencyRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
}

type UpdateCompetencyRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Active      *bool   `json:"active"`
}

type CreateObjectiveRequest struct {
	EmployeeID   string  `json:"employee_id" binding:"required"`
	CycleID      string  `json:"cycle_id" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description"`
	Metric       *string `json:"metric"`
	TargetValue  *float64 `json:"target_value"`
	Unit         *string `json:"unit"`
	Weight       *float64 `json:"weight"`
	StartDate    *string `json:"start_date"`
	DueDate      *string `json:"due_date"`
}

type UpdateObjectiveRequest struct {
	Title        *string  `json:"title"`
	Description  *string  `json:"description"`
	Metric       *string  `json:"metric"`
	TargetValue  *float64 `json:"target_value"`
	CurrentValue *float64 `json:"current_value"`
	Unit         *string  `json:"unit"`
	Weight       *float64 `json:"weight"`
	Status       *string  `json:"status"`
}

type UpdateProgressRequest struct {
	CurrentValue float64 `json:"current_value" binding:"required"`
}

type CreateKPIRequest struct {
	EmployeeID  string   `json:"employee_id" binding:"required"`
	CycleID     string   `json:"cycle_id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	TargetValue *float64 `json:"target_value"`
	Unit        *string  `json:"unit"`
	Weight      *float64 `json:"weight"`
}

type UpdateKPIRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	TargetValue  *float64 `json:"target_value"`
	CurrentValue *float64 `json:"current_value"`
	Unit         *string  `json:"unit"`
	Weight       *float64 `json:"weight"`
	Status       *string  `json:"status"`
}

type AssignEvaluatorsRequest struct {
	EmployeeID string                 `json:"employee_id" binding:"required"`
	CycleID    string                 `json:"cycle_id" binding:"required"`
	Evaluators []EvaluatorAssignment  `json:"evaluators" binding:"required"`
}

type EvaluatorAssignment struct {
	EvaluatorID   string `json:"evaluator_id" binding:"required"`
	EvaluatorType string `json:"evaluator_type" binding:"required"`
}

type CreateEvaluationRequest struct {
	CycleID       string                `json:"cycle_id" binding:"required"`
	EmployeeID    string                `json:"employee_id" binding:"required"`
	EvaluatorID   string                `json:"evaluator_id" binding:"required"`
	EvaluatorType string                `json:"evaluator_type" binding:"required"`
	TemplateID    *string               `json:"template_id"`
	Comments      *string               `json:"comments"`
	Answers       []CreateAnswerRequest  `json:"answers"`
}

type UpdateEvaluationRequest struct {
	Comments *string `json:"comments"`
}

type CreateAnswerRequest struct {
	SectionName *string  `json:"section_name"`
	ItemName    string   `json:"item_name" binding:"required"`
	ItemType    *string  `json:"item_type"`
	Score       *float64 `json:"score"`
	Value       *string  `json:"value"`
	Comments    *string  `json:"comments"`
	Weight      *float64 `json:"weight"`
}

type CreateFeedbackRequest struct {
	EmployeeID        string  `json:"employee_id" binding:"required"`
	CycleID           *string `json:"cycle_id"`
	FeedbackType      string  `json:"feedback_type" binding:"required"`
	Message           string  `json:"message" binding:"required"`
	IsPrivate         *bool   `json:"is_private"`
	VisibleToEmployee *bool   `json:"visible_to_employee"`
}

type CreateEvidenceRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description"`
	EvidenceType string  `json:"evidence_type"`
	URL          *string `json:"url"`
}

type CreateImprovementPlanRequest struct {
	EmployeeID         string  `json:"employee_id" binding:"required"`
	CycleID            *string `json:"cycle_id"`
	ResultID           *string `json:"result_id"`
	Title              string  `json:"title" binding:"required"`
	ProblemDescription *string `json:"problem_description"`
	Objective          *string `json:"objective"`
	ResponsibleID      *string `json:"responsible_id"`
	DueDate            *string `json:"due_date"`
	Actions            []CreatePlanActionRequest `json:"actions"`
}

type CreatePlanActionRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
}

type UpdateImprovementPlanRequest struct {
	Title              *string `json:"title"`
	ProblemDescription *string `json:"problem_description"`
	Objective          *string `json:"objective"`
	ResponsibleID      *string `json:"responsible_id"`
	DueDate            *string `json:"due_date"`
	Status             *string `json:"status"`
	Outcome            *string `json:"outcome"`
}

type CreateDevelopmentPlanRequest struct {
	EmployeeID     string  `json:"employee_id" binding:"required"`
	CycleID        *string `json:"cycle_id"`
	Title          string  `json:"title" binding:"required"`
	Description    *string `json:"description"`
	CareerGoal     *string `json:"career_goal"`
	TimelineMonths *int    `json:"timeline_months"`
	Actions        []CreateDevActionRequest `json:"actions"`
}

type CreateDevActionRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	ActionType  *string `json:"action_type"`
	DueDate     *string `json:"due_date"`
}

type UpdateDevelopmentPlanRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	CareerGoal     *string `json:"career_goal"`
	TimelineMonths *int    `json:"timeline_months"`
	Status         *string `json:"status"`
}

type UpdateScoringRulesRequest struct {
	ObjectiveWeight  *float64 `json:"objective_weight"`
	CompetencyWeight *float64 `json:"competency_weight"`
	KPIWeight        *float64 `json:"kpi_weight"`
	SelfEvalWeight   *float64 `json:"self_eval_weight"`
	ManagerWeight    *float64 `json:"manager_weight"`
	PeerWeight       *float64 `json:"peer_weight"`
	HRWeight         *float64 `json:"hr_weight"`
}

type CalculateResultRequest struct {
	CycleID    string `json:"cycle_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
}
