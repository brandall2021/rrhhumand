package domain

type OnboardingStatus string

const (
	OnboardingDraft       OnboardingStatus = "DRAFT"
	OnboardingPending     OnboardingStatus = "PENDING"
	OnboardingInProgress  OnboardingStatus = "IN_PROGRESS"
	OnboardingBlocked     OnboardingStatus = "BLOCKED"
	OnboardingCompleted   OnboardingStatus = "COMPLETED"
	OnboardingCancelled   OnboardingStatus = "CANCELLED"
)

type WorkMode string

const (
	WorkModeRemote    WorkMode = "REMOTE"
	WorkModePresencial WorkMode = "PRESENCIAL"
	WorkModeHybrid    WorkMode = "HIBRIDO"
)

type EmploymentType string

const (
	EmploymentTypeAdmin       EmploymentType = "ADMINISTRATIVO"
	EmploymentTypeTeacher     EmploymentType = "DOCENTE"
	EmploymentTypeTechnical   EmploymentType = "TECNICO"
	EmploymentTypeDeveloper   EmploymentType = "DESARROLLADOR"
	EmploymentTypeManager     EmploymentType = "MANAGER"
	EmploymentTypeDirector    EmploymentType = "DIRECTOR"
	EmploymentTypeIntern      EmploymentType = "PASANTE"
	EmploymentTypeContractor  EmploymentType = "CONTRATISTA"
)

type ProbationStatus string

const (
	ProbationPending   ProbationStatus = "PENDING"
	ProbationCompleted ProbationStatus = "COMPLETED"
	ProbationExtended  ProbationStatus = "EXTENDED"
	ProbationTerminated ProbationStatus = "TERMINATED"
)

type OnboardingProcess struct {
	ID                     string
	CompanyID              string
	EmployeeID             string
	CandidateID            *string
	ApplicationID          *string
	JobOfferID             *string
	TemplateID             *string
	Status                 OnboardingStatus
	StartDate              string
	ExpectedCompletionDate *string
	ActualCompletionDate   *string
	Progress               float64
	CompletionPolicy       string
	EmployeeType           *EmploymentType
	WorkMode               *WorkMode
	ProbationStartDate     *string
	ProbationEndDate       *string
	ProbationStatus        ProbationStatus
	CreatedBy              string
	CreatedAt              string
	UpdatedAt              string
}

type OnboardingEvent struct {
	EventType     string
	OnboardingID  string
	EmployeeID    string
	CompanyID     string
	CompletedAt   *string
	LastWorkingDate *string
}
