package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type CatalogService struct {
	catalogRepo *repository.CatalogRepo
}

func NewCatalogService(catalogRepo *repository.CatalogRepo) *CatalogService {
	return &CatalogService{catalogRepo: catalogRepo}
}

func (s *CatalogService) CreateCategory(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, name string, description *string, icon *string, color *string, sortOrder int) (*domain.BenefitCategory, error) {
	c := &domain.BenefitCategory{
		ID:          uuid.New(),
		CompanyID:   companyID,
		Name:        name,
		Description: description,
		Icon:        icon,
		Color:       color,
		SortOrder:   sortOrder,
		IsActive:    true,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.catalogRepo.CreateCategory(ctx, c); err != nil {
		return nil, svcErr("CreateCategory", err)
	}
	return c, nil
}

func (s *CatalogService) GetCategory(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitCategory, error) {
	return s.catalogRepo.GetCategory(ctx, companyID, id)
}

func (s *CatalogService) ListCategories(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitCategory, error) {
	return s.catalogRepo.ListCategories(ctx, companyID)
}

func (s *CatalogService) UpdateCategory(ctx context.Context, companyID uuid.UUID, category *domain.BenefitCategory) (*domain.BenefitCategory, error) {
	category.CompanyID = companyID
	category.UpdatedAt = time.Now()
	if err := s.catalogRepo.UpdateCategory(ctx, category); err != nil {
		return nil, svcErr("UpdateCategory", err)
	}
	return category, nil
}

func (s *CatalogService) DeleteCategory(ctx context.Context, companyID, id uuid.UUID) error {
	return s.catalogRepo.DeleteCategory(ctx, companyID, id)
}

func (s *CatalogService) CreateType(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, t *domain.BenefitType) (*domain.BenefitType, error) {
	t.ID = uuid.New()
	t.CompanyID = companyID
	t.CreatedBy = userID
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if err := s.catalogRepo.CreateType(ctx, t); err != nil {
		return nil, svcErr("CreateType", err)
	}
	return t, nil
}

func (s *CatalogService) GetType(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitType, error) {
	return s.catalogRepo.GetType(ctx, companyID, id)
}

func (s *CatalogService) ListTypes(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitType, error) {
	return s.catalogRepo.ListTypes(ctx, companyID)
}

func (s *CatalogService) UpdateType(ctx context.Context, companyID uuid.UUID, t *domain.BenefitType) (*domain.BenefitType, error) {
	t.CompanyID = companyID
	t.UpdatedAt = time.Now()
	if err := s.catalogRepo.UpdateType(ctx, t); err != nil {
		return nil, svcErr("UpdateType", err)
	}
	return t, nil
}

func (s *CatalogService) DeleteType(ctx context.Context, companyID, id uuid.UUID) error {
	return s.catalogRepo.DeleteType(ctx, companyID, id)
}

func (s *CatalogService) CreateProvider(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, p *domain.BenefitProvider) (*domain.BenefitProvider, error) {
	p.ID = uuid.New()
	p.CompanyID = companyID
	p.CreatedBy = userID
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err := s.catalogRepo.CreateProvider(ctx, p); err != nil {
		return nil, svcErr("CreateProvider", err)
	}
	return p, nil
}

func (s *CatalogService) GetProvider(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitProvider, error) {
	return s.catalogRepo.GetProvider(ctx, companyID, id)
}

func (s *CatalogService) ListProviders(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitProvider, error) {
	return s.catalogRepo.ListProviders(ctx, companyID)
}

func (s *CatalogService) UpdateProvider(ctx context.Context, companyID uuid.UUID, p *domain.BenefitProvider) (*domain.BenefitProvider, error) {
	p.CompanyID = companyID
	p.UpdatedAt = time.Now()
	if err := s.catalogRepo.UpdateProvider(ctx, p); err != nil {
		return nil, svcErr("UpdateProvider", err)
	}
	return p, nil
}

func (s *CatalogService) DeleteProvider(ctx context.Context, companyID, id uuid.UUID) error {
	return s.catalogRepo.DeleteProvider(ctx, companyID, id)
}
