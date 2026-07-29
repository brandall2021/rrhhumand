package http

type OnboardingDashboardResponse struct {
	ActiveOnboardings     int     `json:"active_onboardings"`
	PendingOnboardings    int     `json:"pending_onboardings"`
	OverdueOnboardings    int     `json:"overdue_onboardings"`
	CompletedOnboardings  int     `json:"completed_onboardings"`
	AverageProgress       float64 `json:"average_progress"`
	TasksDueToday         int     `json:"tasks_due_today"`
	PendingDocuments      int     `json:"pending_documents"`
	PendingAssets         int     `json:"pending_assets"`
	PendingAccess         int     `json:"pending_access"`
	PendingTraining       int     `json:"pending_training"`
}

type OffboardingDashboardResponse struct {
	ActiveOffboardings    int `json:"active_offboardings"`
	PendingOffboardings   int `json:"pending_offboardings"`
	OverdueOffboardings   int `json:"overdue_offboardings"`
	CompletedOffboardings int `json:"completed_offboardings"`
	AssetsNotReturned     int `json:"assets_not_returned"`
	AccessNotRevoked      int `json:"access_not_revoked"`
}

type EmployeeDashboardResponse struct {
	Status         string  `json:"status"`
	Progress       float64 `json:"progress"`
	TasksTotal     int     `json:"tasks_total"`
	TasksCompleted int     `json:"tasks_completed"`
	DocumentsTotal int     `json:"documents_total"`
	DocsApproved   int     `json:"docs_approved"`
	NextDueDate    *string `json:"next_due_date,omitempty"`
}
