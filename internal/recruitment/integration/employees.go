package integration

import (
	"context"
)

type EmployeeInfo struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Phone        *string `json:"phone,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
	PositionID   *string `json:"position_id,omitempty"`
	JobTitle     *string `json:"job_title,omitempty"`
	IsActive     bool    `json:"is_active"`
}

type EmployeeAdapter struct{}

func NewEmployeeAdapter() *EmployeeAdapter {
	return &EmployeeAdapter{}
}

func (a *EmployeeAdapter) GetEmployeeByID(ctx context.Context, id string) (*EmployeeInfo, error) {
	return &EmployeeInfo{
		ID:        id,
		FirstName: "Employee",
		LastName:  "Name",
		Email:     "employee@company.com",
		IsActive:  true,
	}, nil
}

func (a *EmployeeAdapter) SearchEmployees(ctx context.Context, query string) ([]*EmployeeInfo, error) {
	return []*EmployeeInfo{
		{
			ID:        "emp-1",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@company.com",
			IsActive:  true,
		},
	}, nil
}

func (a *EmployeeAdapter) IsHiringManager(ctx context.Context, employeeID string) (bool, error) {
	return true, nil
}

func (a *EmployeeAdapter) GetEmployeeByEmail(ctx context.Context, email string) (*EmployeeInfo, error) {
	return &EmployeeInfo{
		ID:        "emp-1",
		FirstName: "Employee",
		LastName:  "Name",
		Email:     email,
		IsActive:  true,
	}, nil
}
