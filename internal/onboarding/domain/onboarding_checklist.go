package domain

type ChecklistSection string

const (
	ChecklistDocumentacion  ChecklistSection = "DOCUMENTACION"
	ChecklistTecnologia     ChecklistSection = "TECNOLOGIA"
	ChecklistCapacitacion   ChecklistSection = "CAPACITACION"
	ChecklistFirstDay       ChecklistSection = "PRIMER_DIA"
	ChecklistFirstWeek      ChecklistSection = "PRIMERA_SEMANA"
	ChecklistFirstMonth     ChecklistSection = "PRIMER_MES"
)

type ChecklistStatus string

const (
	ChecklistPending    ChecklistStatus = "PENDING"
	ChecklistCompleted  ChecklistStatus = "COMPLETED"
)

type OnboardingChecklistItem struct {
	ID           string
	CompanyID    string
	OnboardingID string
	EmployeeID   string
	Section      ChecklistSection
	Title        string
	Status       ChecklistStatus
	CompletedAt  *string
	CompletedBy  *string
	SortOrder    int
	CreatedAt    string
	UpdatedAt    string
}
