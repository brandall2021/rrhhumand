package companies

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type CompanyService struct {
	repo *CompanyRepository
}

func NewCompanyService(repo *CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

type CreateCompanyRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type UpdateCompanyRequest struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
	Plan *string `json:"plan,omitempty"`
}

func (s *CompanyService) Create(ctx context.Context, req *CreateCompanyRequest) (*models.Company, error) {
	existing, _ := s.repo.FindByID(ctx, "")
	_ = existing

	company := &models.Company{
		ID:     uuid.New().String(),
		Name:   req.Name,
		Slug:   req.Slug,
		Plan:   "free",
		Active: true,
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) GetByID(ctx context.Context, id string) (*models.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) Update(ctx context.Context, id string, req *UpdateCompanyRequest) (*models.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("company not found")
	}

	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Slug != nil {
		company.Slug = *req.Slug
	}
	if req.Plan != nil {
		company.Plan = *req.Plan
	}

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) List(ctx context.Context, params *models.PaginationParams) ([]models.Company, int64, error) {
	return s.repo.List(ctx, params.Offset, params.PerPage)
}
