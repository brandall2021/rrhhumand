package domain

type TaskType string

const (
	TaskDocument   TaskType = "DOCUMENT"
	TaskApproval   TaskType = "APPROVAL"
	TaskTraining   TaskType = "TRAINING"
	TaskAccount    TaskType = "ACCOUNT"
	TaskAsset      TaskType = "ASSET"
	TaskMeeting    TaskType = "MEETING"
	TaskSignature  TaskType = "SIGNATURE"
	TaskChecklist  TaskType = "CHECKLIST"
	TaskInfo       TaskType = "INFORMATION"
	TaskSystem     TaskType = "SYSTEM"
	TaskOther      TaskType = "OTHER"
)

type TaskStatus string

const (
	TaskPending    TaskStatus = "PENDING"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskBlocked    TaskStatus = "BLOCKED"
	TaskCompleted  TaskStatus = "COMPLETED"
	TaskCancelled  TaskStatus = "CANCELLED"
	TaskOverdue    TaskStatus = "OVERDUE"
)

type ResponsibleType string

const (
	ResponsibleEmployee  ResponsibleType = "EMPLOYEE"
	ResponsibleManager   ResponsibleType = "MANAGER"
	ResponsibleHR        ResponsibleType = "HR"
	ResponsibleIT        ResponsibleType = "IT"
	ResponsibleFinance   ResponsibleType = "FINANCE"
	ResponsibleSecurity  ResponsibleType = "SECURITY"
	ResponsibleTraining  ResponsibleType = "TRAINING"
	ResponsibleLegal     ResponsibleType = "LEGAL"
	ResponsibleExternal  ResponsibleType = "EXTERNAL"
)

type OnboardingTask struct {
	ID              string
	CompanyID       string
	TemplateID      string
	Title           string
	Description     *string
	TaskType        TaskType
	OrderIndex      int
	DueDays         int
	Required        bool
	RequiresApproval bool
	DepartmentID    *string
	Role            *string
	Active          bool
	CreatedAt       string
	UpdatedAt       string
}

type OnboardingTaskAssignment struct {
	ID            string
	OnboardingID  string
	TaskID        string
	AssignedTo    *string
	AssignedRole  *string
	Status        TaskStatus
	DueDate       *string
	CompletedAt   *string
	CompletedBy   *string
	Comments      *string
	CreatedAt     string
	UpdatedAt     string
}

type OnboardingTaskDependency struct {
	ID               string
	TaskID           string
	DependsOnTaskID  string
}
