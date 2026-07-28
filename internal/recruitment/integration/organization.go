package integration

import (
	"context"
)

type DepartmentInfo struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Code        *string `json:"code,omitempty"`
	ManagerID   *string `json:"manager_id,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	CostCenter  *string `json:"cost_center,omitempty"`
	Active      bool    `json:"active"`
}

type PositionInfo struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	DepartmentID *string `json:"department_id,omitempty"`
	Title        string  `json:"title"`
	Grade        *string `json:"grade,omitempty"`
	Active       bool    `json:"active"`
}

type OrgAdapter struct{}

func NewOrgAdapter() *OrgAdapter {
	return &OrgAdapter{}
}

func (a *OrgAdapter) GetDepartmentByID(ctx context.Context, id string) (*DepartmentInfo, error) {
	return &DepartmentInfo{
		ID:     id,
		Name:   "Engineering",
		Active: true,
	}, nil
}

func (a *OrgAdapter) GetPositionByID(ctx context.Context, id string) (*PositionInfo, error) {
	return &PositionInfo{
		ID:     id,
		Title:  "Software Engineer",
		Active: true,
	}, nil
}

func (a *OrgAdapter) ListDepartments(ctx context.Context, companyID string) ([]DepartmentInfo, error) {
	return []DepartmentInfo{
		{ID: "dept-1", CompanyID: companyID, Name: "Engineering", Active: true},
		{ID: "dept-2", CompanyID: companyID, Name: "Marketing", Active: true},
		{ID: "dept-3", CompanyID: companyID, Name: "Sales", Active: true},
	}, nil
}

func (a *OrgAdapter) GetDepartmentTree(ctx context.Context, companyID string) ([]DepartmentInfo, error) {
	return a.ListDepartments(ctx, companyID)
}
