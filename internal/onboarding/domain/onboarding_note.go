package domain

type NoteVisibility string

const (
	NoteEmployee  NoteVisibility = "EMPLOYEE"
	NoteManager   NoteVisibility = "MANAGER"
	NoteHROnly    NoteVisibility = "HR_ONLY"
	NoteAdminOnly NoteVisibility = "ADMIN_ONLY"
)

type OnboardingNote struct {
	ID           string
	CompanyID    string
	OnboardingID string
	EntityType   string
	EntityID     *string
	Content      string
	Visibility   NoteVisibility
	CreatedBy    string
	CreatedAt    string
	UpdatedAt    string
}
