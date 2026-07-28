package branches

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type BranchService struct {
	repo *BranchRepository
}

func NewBranchService(repo *BranchRepository) *BranchService {
	return &BranchService{repo: repo}
}

type CreateBranchRequest struct {
	Name     string  `json:"name" validate:"required"`
	Code     *string `json:"code,omitempty"`
	Address  *string `json:"address,omitempty"`
	City     *string `json:"city,omitempty"`
	State    *string `json:"state,omitempty"`
	Country  *string `json:"country,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type UpdateBranchRequest struct {
	Name     *string `json:"name,omitempty"`
	Code     *string `json:"code,omitempty"`
	Address  *string `json:"address,omitempty"`
	City     *string `json:"city,omitempty"`
	State    *string `json:"state,omitempty"`
	Country  *string `json:"country,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}

func (s *BranchService) Create(ctx context.Context, companyID string, req *CreateBranchRequest) (*models.Branch, error) {
	branch := &models.Branch{
		ID:        uuid.New().String(),
		CompanyID: companyID,
		Name:      req.Name,
		Code:      req.Code,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		Country:   "Argentina",
		Phone:     req.Phone,
		Email:     req.Email,
		Timezone:  "America/Argentina/Buenos_Aires",
		Active:    true,
	}
	if req.Country != nil {
		branch.Country = *req.Country
	}
	if req.Timezone != nil {
		branch.Timezone = *req.Timezone
	}

	if err := s.repo.Create(ctx, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) GetByID(ctx context.Context, id, companyID string) (*models.Branch, error) {
	branch, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) List(ctx context.Context, companyID string, params *models.PaginationParams, search string) ([]models.Branch, int64, error) {
	return s.repo.List(ctx, companyID, params.Offset, params.PerPage, search)
}

func (s *BranchService) Update(ctx context.Context, id, companyID string, req *UpdateBranchRequest) (*models.Branch, error) {
	branch, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Code != nil {
		branch.Code = req.Code
	}
	if req.Address != nil {
		branch.Address = req.Address
	}
	if req.City != nil {
		branch.City = req.City
	}
	if req.State != nil {
		branch.State = req.State
	}
	if req.Country != nil {
		branch.Country = *req.Country
	}
	if req.Phone != nil {
		branch.Phone = req.Phone
	}
	if req.Email != nil {
		branch.Email = req.Email
	}
	if req.Timezone != nil {
		branch.Timezone = *req.Timezone
	}
	if req.Active != nil {
		branch.Active = *req.Active
	}

	if err := s.repo.Update(ctx, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) Delete(ctx context.Context, id, companyID string) error {
	_, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return errors.New("branch not found")
	}
	return s.repo.Delete(ctx, id, companyID)
}
