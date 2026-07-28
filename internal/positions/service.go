package positions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type PositionService struct {
	repo *PositionRepository
}

func NewPositionService(repo *PositionRepository) *PositionService {
	return &PositionService{repo: repo}
}

type CreatePositionRequest struct {
	Name         string   `json:"name" validate:"required"`
	Code         *string  `json:"code,omitempty"`
	Description  *string  `json:"description,omitempty"`
	DepartmentID *string  `json:"department_id,omitempty"`
	Level        *int     `json:"level,omitempty"`
	MinSalary    *float64 `json:"min_salary,omitempty"`
	MaxSalary    *float64 `json:"max_salary,omitempty"`
}

type UpdatePositionRequest struct {
	Name         *string  `json:"name,omitempty"`
	Code         *string  `json:"code,omitempty"`
	Description  *string  `json:"description,omitempty"`
	DepartmentID *string  `json:"department_id,omitempty"`
	Level        *int     `json:"level,omitempty"`
	MinSalary    *float64 `json:"min_salary,omitempty"`
	MaxSalary    *float64 `json:"max_salary,omitempty"`
	Active       *bool    `json:"active,omitempty"`
}

func (s *PositionService) Create(ctx context.Context, companyID string, req *CreatePositionRequest) (*models.Position, error) {
	level := 1
	if req.Level != nil {
		level = *req.Level
	}

	pos := &models.Position{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
		Level:        level,
		MinSalary:    req.MinSalary,
		MaxSalary:    req.MaxSalary,
		Active:       true,
	}

	if err := s.repo.Create(ctx, pos); err != nil {
		return nil, err
	}
	return pos, nil
}

func (s *PositionService) GetByID(ctx context.Context, id, companyID string) (*models.Position, error) {
	pos, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("position not found")
		}
		return nil, err
	}
	return pos, nil
}

func (s *PositionService) List(ctx context.Context, companyID string, params *models.PaginationParams, search, departmentID string) ([]models.Position, int64, error) {
	return s.repo.List(ctx, companyID, params.Offset, params.PerPage, search, departmentID)
}

func (s *PositionService) Update(ctx context.Context, id, companyID string, req *UpdatePositionRequest) (*models.Position, error) {
	pos, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("position not found")
	}

	if req.Name != nil {
		pos.Name = *req.Name
	}
	if req.Code != nil {
		pos.Code = req.Code
	}
	if req.Description != nil {
		pos.Description = req.Description
	}
	if req.DepartmentID != nil {
		pos.DepartmentID = req.DepartmentID
	}
	if req.Level != nil {
		pos.Level = *req.Level
	}
	if req.MinSalary != nil {
		pos.MinSalary = req.MinSalary
	}
	if req.MaxSalary != nil {
		pos.MaxSalary = req.MaxSalary
	}
	if req.Active != nil {
		pos.Active = *req.Active
	}

	if err := s.repo.Update(ctx, pos); err != nil {
		return nil, err
	}
	return pos, nil
}

func (s *PositionService) Delete(ctx context.Context, id, companyID string) error {
	_, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return errors.New("position not found")
	}
	return s.repo.Delete(ctx, id, companyID)
}
