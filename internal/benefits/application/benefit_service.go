package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type BenefitService struct {
	benefitRepo  *repository.BenefitRepo
	catalogRepo  *repository.CatalogRepo
	planRepo     *repository.PlanRepo
}

func NewBenefitService(benefitRepo *repository.BenefitRepo, catalogRepo *repository.CatalogRepo, planRepo *repository.PlanRepo) *BenefitService {
	return &BenefitService{
		benefitRepo: benefitRepo,
		catalogRepo: catalogRepo,
		planRepo:    planRepo,
	}
}

func (s *BenefitService) CreateBenefit(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, b *domain.Benefit) (*domain.Benefit, error) {
	b.ID = uuid.New()
	b.CompanyID = companyID
	b.CreatedBy = userID
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	if err := s.benefitRepo.Create(ctx, b); err != nil {
		return nil, svcErr("CreateBenefit", err)
	}
	return b, nil
}

func (s *BenefitService) GetBenefit(ctx context.Context, companyID, id uuid.UUID) (*domain.Benefit, error) {
	return s.benefitRepo.Get(ctx, companyID, id)
}

func (s *BenefitService) ListBenefits(ctx context.Context, companyID uuid.UUID, status *string, typeID *uuid.UUID, visibility *string, limit, offset int) ([]domain.Benefit, error) {
	var typeIDStr *string
	if typeID != nil {
		s := typeID.String()
		typeIDStr = &s
	}
	return s.benefitRepo.List(ctx, companyID, status, typeIDStr, visibility, limit, offset)
}

func (s *BenefitService) UpdateBenefit(ctx context.Context, companyID uuid.UUID, b *domain.Benefit) (*domain.Benefit, error) {
	b.CompanyID = companyID
	b.UpdatedAt = time.Now()
	if err := s.benefitRepo.Update(ctx, b); err != nil {
		return nil, svcErr("UpdateBenefit", err)
	}
	return b, nil
}

func (s *BenefitService) DeleteBenefit(ctx context.Context, companyID, id uuid.UUID) error {
	return s.benefitRepo.Delete(ctx, companyID, id)
}

func (s *BenefitService) SearchBenefits(ctx context.Context, companyID uuid.UUID, query string) ([]domain.Benefit, error) {
	return s.benefitRepo.SearchBenefits(ctx, companyID, query)
}

func (s *BenefitService) CreatePlan(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, p *domain.BenefitPlan) (*domain.BenefitPlan, error) {
	p.ID = uuid.New()
	p.CompanyID = companyID
	p.CreatedBy = userID
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err := s.planRepo.CreatePlan(ctx, p); err != nil {
		return nil, svcErr("CreatePlan", err)
	}
	return p, nil
}

func (s *BenefitService) GetPlan(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitPlan, error) {
	return s.planRepo.GetPlan(ctx, companyID, id)
}

func (s *BenefitService) ListPlans(ctx context.Context, companyID uuid.UUID, benefitID uuid.UUID) ([]domain.BenefitPlan, error) {
	return s.planRepo.ListPlans(ctx, benefitID)
}

func (s *BenefitService) UpdatePlan(ctx context.Context, companyID uuid.UUID, p *domain.BenefitPlan) (*domain.BenefitPlan, error) {
	p.CompanyID = companyID
	p.UpdatedAt = time.Now()
	if err := s.planRepo.UpdatePlan(ctx, p); err != nil {
		return nil, svcErr("UpdatePlan", err)
	}
	return p, nil
}

func (s *BenefitService) DeletePlan(ctx context.Context, companyID, id uuid.UUID) error {
	return s.planRepo.DeletePlan(ctx, companyID, id)
}
