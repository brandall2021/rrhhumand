package http

type CreateOnboardingRequest struct {
	EmployeeID       string  `json:"employee_id" binding:"required"`
	CandidateID      *string `json:"candidate_id"`
	ApplicationID    *string `json:"application_id"`
	JobOfferID       *string `json:"job_offer_id"`
	TemplateID       *string `json:"template_id"`
	StartDate        string  `json:"start_date" binding:"required"`
	CompletionPolicy *string `json:"completion_policy"`
	EmployeeType     *string `json:"employee_type"`
	WorkMode         *string `json:"work_mode"`
}

type UpdateOnboardingRequest struct {
	TemplateID       *string `json:"template_id"`
	CompletionPolicy *string `json:"completion_policy"`
	EmployeeType     *string `json:"employee_type"`
	WorkMode         *string `json:"work_mode"`
}

type CreateTaskRequest struct {
	Title           string  `json:"title" binding:"required"`
	Description     *string `json:"description"`
	TaskType        string  `json:"task_type"`
	AssignedTo      *string `json:"assigned_to"`
	AssignedRole    string  `json:"assigned_role"`
	DueDate         string  `json:"due_date"`
	Required        *bool   `json:"required"`
}

type CreateOffboardingRequest struct {
	EmployeeID              string  `json:"employee_id" binding:"required"`
	TerminationType         string  `json:"termination_type" binding:"required"`
	ReasonID                *string `json:"reason_id"`
	NoticeDate              string  `json:"notice_date" binding:"required"`
	LastWorkingDate         string  `json:"last_working_date" binding:"required"`
	TerminationEffectiveDate *string `json:"termination_effective_date"`
	TemplateID              *string `json:"template_id"`
}

type CompleteExitInterviewRequest struct {
	Reason        string   `json:"reason" binding:"required"`
	Feedback      string   `json:"feedback"`
	Recommendation *string  `json:"recommendation"`
	Rating        *float64 `json:"rating"`
	Anonymous     bool     `json:"anonymous"`
}

type WorkflowRuleRequest struct {
	WorkflowType string `json:"workflow_type" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Conditions   string `json:"conditions" binding:"required"`
	Actions      string `json:"actions" binding:"required"`
	Priority     *int   `json:"priority"`
}

type ReturnAssetRequest struct {
	ConditionOnReturn *string `json:"condition_on_return"`
	Status            string  `json:"status" binding:"required"`
}

type RevokeAccessRequest struct {
	SystemName string `json:"system_name" binding:"required"`
	AccessType string `json:"access_type"`
}

type HandoverRequest struct {
	HandoverTo  string  `json:"handover_to" binding:"required"`
	Description *string `json:"description"`
	Projects    *string `json:"projects"`
	PendingTasks *string `json:"pending_tasks"`
	Documents   *string `json:"documents"`
}

type ProbationRequest struct {
	Status string `json:"status" binding:"required"`
}
