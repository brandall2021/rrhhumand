package domain

type OffboardingTaskStatus string

const (
	OffTaskPending    OffboardingTaskStatus = "PENDING"
	OffTaskInProgress OffboardingTaskStatus = "IN_PROGRESS"
	OffTaskBlocked    OffboardingTaskStatus = "BLOCKED"
	OffTaskCompleted  OffboardingTaskStatus = "COMPLETED"
	OffTaskCancelled  OffboardingTaskStatus = "CANCELLED"
)

type OffboardingTask struct {
	ID            string
	CompanyID     string
	OffboardingID string
	Title         string
	Description   *string
	TaskType      string
	AssignedTo    *string
	AssignedRole  *string
	Required      bool
	DueDate       *string
	Status        OffboardingTaskStatus
	CompletedAt   *string
	CompletedBy   *string
	Comments      *string
	SortOrder     int
	CreatedAt     string
	UpdatedAt     string
}

type OffboardingTaskType string

const (
	OffTaskDocumentacion   OffboardingTaskType = "DOCUMENTACION"
	OffTaskActivos         OffboardingTaskType = "ACTIVOS"
	OffTaskAccesos         OffboardingTaskType = "ACCESOS"
	OffTaskEntrevista      OffboardingTaskType = "ENTREVISTA"
	OffTaskLiquidacion     OffboardingTaskType = "LIQUIDACION"
	OffTaskTransferencia   OffboardingTaskType = "TRANSFERENCIA"
	OffTaskNotificacion    OffboardingTaskType = "NOTIFICACION"
	OffTaskCertificacion   OffboardingTaskType = "CERTIFICACION"
	OffTaskOther           OffboardingTaskType = "OTHER"
)
