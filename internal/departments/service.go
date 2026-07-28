package departments

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type DepartmentService struct {
	repo *DepartmentRepository
}

func NewDepartmentService(repo *DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

type CreateDepartmentRequest struct {
	Name        string  `json:"name" validate:"required"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	BranchID    *string `json:"branch_id,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	BranchID    *string `json:"branch_id,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

func (s *DepartmentService) Create(ctx context.Context, companyID string, req *CreateDepartmentRequest) (*models.Department, error) {
	dept := &models.Department{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		BranchID:    req.BranchID,
		ParentID:    req.ParentID,
		Active:      true,
	}

	if err := s.repo.Create(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DepartmentService) GetByID(ctx context.Context, id, companyID string) (*models.Department, error) {
	dept, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return dept, nil
}

func (s *DepartmentService) List(ctx context.Context, companyID string, params *models.PaginationParams, search string) ([]models.Department, int64, error) {
	return s.repo.List(ctx, companyID, params.Offset, params.PerPage, search)
}

func (s *DepartmentService) Update(ctx context.Context, id, companyID string, req *UpdateDepartmentRequest) (*models.Department, error) {
	dept, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	if req.Name != nil {
		dept.Name = *req.Name
	}
	if req.Code != nil {
		dept.Code = req.Code
	}
	if req.Description != nil {
		dept.Description = req.Description
	}
	if req.BranchID != nil {
		dept.BranchID = req.BranchID
	}
	if req.ParentID != nil {
		dept.ParentID = req.ParentID
	}
	if req.Active != nil {
		dept.Active = *req.Active
	}

	if err := s.repo.Update(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DepartmentService) Delete(ctx context.Context, id, companyID string) error {
	_, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return errors.New("department not found")
	}
	return s.repo.Delete(ctx, id, companyID)
}
