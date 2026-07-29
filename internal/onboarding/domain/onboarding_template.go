package domain

type OnboardingTemplate struct {
	ID           string
	CompanyID    string
	Name         string
	Description  *string
	EmploymentType *EmploymentType
	DepartmentID *string
	PositionID   *string
	LocationID   *string
	Active       bool
	CreatedAt    string
	UpdatedAt    string
}
