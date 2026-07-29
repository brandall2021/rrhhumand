package domain

type OutboxEvent struct {
	ID            string
	CompanyID     string
	EventType     string
	AggregateType string
	AggregateID   string
	Payload       *string
	Status        string
	RetryCount    int
	LastError     *string
	CreatedAt     string
	ProcessedAt   *string
}

type OutboxStatus string

const (
	OutboxPending   OutboxStatus = "PENDING"
	OutboxProcessed OutboxStatus = "PROCESSED"
	OutboxFailed    OutboxStatus = "FAILED"
)

type EventType string

const (
	EventOnboardingCreated             EventType = "onboarding.created"
	EventOnboardingStarted             EventType = "onboarding.started"
	EventOnboardingTaskCreated         EventType = "onboarding.task.created"
	EventOnboardingTaskCompleted       EventType = "onboarding.task.completed"
	EventOnboardingTaskOverdue         EventType = "onboarding.task.overdue"
	EventOnboardingDocumentUploaded    EventType = "onboarding.document.uploaded"
	EventOnboardingDocumentApproved    EventType = "onboarding.document.approved"
	EventOnboardingDocumentRejected    EventType = "onboarding.document.rejected"
	EventOnboardingAssetAssigned       EventType = "onboarding.asset.assigned"
	EventOnboardingAssetReturned       EventType = "onboarding.asset.returned"
	EventOnboardingTrainingAssigned    EventType = "onboarding.training.assigned"
	EventOnboardingTrainingCompleted   EventType = "onboarding.training.completed"
	EventOnboardingProbationDue        EventType = "onboarding.probation.due"
	EventOnboardingCompleted           EventType = "onboarding.completed"
	EventOffboardingCreated            EventType = "offboarding.created"
	EventOffboardingApproved           EventType = "offboarding.approved"
	EventOffboardingStarted            EventType = "offboarding.started"
	EventOffboardingTaskCreated        EventType = "offboarding.task.created"
	EventOffboardingTaskCompleted      EventType = "offboarding.task.completed"
	EventOffboardingAssetReturnRequested EventType = "offboarding.asset.return_requested"
	EventOffboardingAssetReturned      EventType = "offboarding.asset.returned"
	EventOffboardingAccessRevocationRequested EventType = "offboarding.access_revocation_requested"
	EventOffboardingAccessRevoked      EventType = "offboarding.access_revoked"
	EventOffboardingExitInterviewCompleted EventType = "offboarding.exit_interview.completed"
	EventOffboardingFinalSettlementRequested EventType = "offboarding.final_settlement_requested"
	EventOffboardingCompleted          EventType = "offboarding.completed"
	EventEmployeeTerminated            EventType = "employee.terminated"
)
