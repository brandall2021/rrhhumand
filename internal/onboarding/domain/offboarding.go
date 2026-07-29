package domain

type OffboardingStatus string

const (
	OffboardingDraft          OffboardingStatus = "DRAFT"
	OffboardingPendingApproval OffboardingStatus = "PENDING_APPROVAL"
	OffboardingApproved       OffboardingStatus = "APPROVED"
	OffboardingInProgress     OffboardingStatus = "IN_PROGRESS"
	OffboardingBlocked        OffboardingStatus = "BLOCKED"
	OffboardingCompleted      OffboardingStatus = "COMPLETED"
	OffboardingCancelled      OffboardingStatus = "CANCELLED"
)

type TerminationType string

const (
	TermResignation     TerminationType = "RESIGNATION"
	TermRetirement      TerminationType = "RETIREMENT"
	TermTermination     TerminationType = "TERMINATION"
	TermEndOfContract   TerminationType = "END_OF_CONTRACT"
	TermMutualAgreement TerminationType = "MUTUAL_AGREEMENT"
	TermLayoff          TerminationType = "LAYOFF"
	TermTransfer        TerminationType = "TRANSFER"
	TermOther           TerminationType = "OTHER"
)

type OffboardingProcess struct {
	ID                      string
	CompanyID               string
	EmployeeID              string
	TemplateID              *string
	RequestedBy             string
	TerminationType         TerminationType
	ReasonID                *string
	NoticeDate              string
	LastWorkingDate         string
	TerminationEffectiveDate *string
	Status                  OffboardingStatus
	Progress                float64
	EmployeeStatusAfter     string
	CreatedAt               string
	UpdatedAt               string
	CompletedAt             *string
}

type OffboardingTemplate struct {
	ID          string
	CompanyID   string
	Name        string
	Description *string
	Active      bool
	CreatedBy   *string
	CreatedAt   string
	UpdatedAt   string
}

type OffboardingTemplateTask struct {
	ID            string
	TemplateID    string
	Title         string
	Description   *string
	TaskType      string
	AssignedRole  *string
	Required      bool
	OrderIndex    int
	CreatedAt     string
	UpdatedAt     string
}

type TerminationReason struct {
	ID          string
	CompanyID   string
	Name        string
	Description *string
	Active      bool
}
