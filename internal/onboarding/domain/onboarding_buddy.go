package domain

type BuddyStatus string

const (
	BuddyActive   BuddyStatus = "ACTIVE"
	BuddyCompleted BuddyStatus = "COMPLETED"
)

type OnboardingBuddy struct {
	ID              string
	CompanyID       string
	OnboardingID    string
	EmployeeID      string
	BuddyEmployeeID string
	StartDate       string
	EndDate         *string
	Status          BuddyStatus
	Notes           *string
	CreatedAt       string
	UpdatedAt       string
}
