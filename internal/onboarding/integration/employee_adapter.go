package integration

import (
	"context"
	"fmt"
	"log"
)

type EmployeeAdapter struct{}

func NewEmployeeAdapter() *EmployeeAdapter {
	return &EmployeeAdapter{}
}

func (a *EmployeeAdapter) Create(ctx context.Context, companyID string, req *CreateEmployeeRequest) (*EmployeeInfo, error) {
	log.Printf("[EmployeeAdapter] Create company=%s employee=%s %s", companyID, req.FirstName, req.LastName)
	return &EmployeeInfo{
		ID:        fmt.Sprintf("emp-%s", req.EmployeeNumber),
		CompanyID: companyID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     &req.Email,
		Status:    "ACTIVE",
	}, nil
}

func (a *EmployeeAdapter) GetByID(ctx context.Context, id, companyID string) (*EmployeeInfo, error) {
	log.Printf("[EmployeeAdapter] GetByID id=%s company=%s", id, companyID)
	return &EmployeeInfo{
		ID:        id,
		CompanyID: companyID,
		FirstName: "John",
		LastName:  "Doe",
		Status:    "ACTIVE",
	}, nil
}

func (a *EmployeeAdapter) Update(ctx context.Context, id, companyID string, req interface{}) (*EmployeeInfo, error) {
	log.Printf("[EmployeeAdapter] Update id=%s company=%s", id, companyID)
	return &EmployeeInfo{
		ID:        id,
		CompanyID: companyID,
		Status:    "UPDATED",
	}, nil
}

func (a *EmployeeAdapter) ExistsByEmail(ctx context.Context, companyID, email string) (bool, error) {
	log.Printf("[EmployeeAdapter] ExistsByEmail company=%s email=%s", companyID, email)
	return false, nil
}

