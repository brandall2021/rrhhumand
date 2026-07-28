package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type EligibilityEngine interface {
	Evaluate(ctx context.Context, employeeID uuid.UUID, rules []domain.BenefitEligibilityRule) (bool, []string, error)
}

type EligibilityService struct {
	eligibilityRepo *repository.EligibilityRepo
	benefitRepo     *repository.BenefitRepo
	engine          EligibilityEngine
}

func NewEligibilityService(eligibilityRepo *repository.EligibilityRepo, benefitRepo *repository.BenefitRepo, engine EligibilityEngine) *EligibilityService {
	return &EligibilityService{
		eligibilityRepo: eligibilityRepo,
		benefitRepo:     benefitRepo,
		engine:          engine,
	}
}

func (s *EligibilityService) CreateRule(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, r *domain.BenefitEligibilityRule) (*domain.BenefitEligibilityRule, error) {
	r.ID = uuid.New()
	r.CompanyID = companyID
	r.CreatedBy = userID
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	if err := s.eligibilityRepo.CreateRule(ctx, r); err != nil {
		return nil, svcErr("CreateRule", err)
	}
	return r, nil
}

func (s *EligibilityService) GetRule(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitEligibilityRule, error) {
	return s.eligibilityRepo.GetRule(ctx, companyID, id)
}

func (s *EligibilityService) ListRules(ctx context.Context, companyID uuid.UUID, benefitID uuid.UUID) ([]domain.BenefitEligibilityRule, error) {
	return s.eligibilityRepo.ListRules(ctx, benefitID)
}

func (s *EligibilityService) UpdateRule(ctx context.Context, companyID uuid.UUID, r *domain.BenefitEligibilityRule) (*domain.BenefitEligibilityRule, error) {
	r.CompanyID = companyID
	r.UpdatedAt = time.Now()
	if err := s.eligibilityRepo.UpdateRule(ctx, r); err != nil {
		return nil, svcErr("UpdateRule", err)
	}
	return r, nil
}

func (s *EligibilityService) DeleteRule(ctx context.Context, companyID, id uuid.UUID) error {
	return s.eligibilityRepo.DeleteRule(ctx, companyID, id)
}

func (s *EligibilityService) EvaluateEmployee(ctx context.Context, companyID, employeeID, benefitID uuid.UUID) (bool, []string, error) {
	rules, err := s.eligibilityRepo.ListRules(ctx, benefitID)
	if err != nil {
		return false, nil, svcErr("EvaluateEmployee", err)
	}
	if len(rules) == 0 {
		return true, nil, nil
	}
	eligible, reasons, err := s.engine.Evaluate(ctx, employeeID, rules)
	if err != nil {
		return false, nil, svcErr("EvaluateEmployee", err)
	}
	return eligible, reasons, nil
}

func (s *EligibilityService) ListEligibleBenefits(ctx context.Context, companyID, employeeID uuid.UUID) ([]domain.Benefit, error) {
	allBenefits, err := s.benefitRepo.List(ctx, companyID, nil, nil, nil, 0, 0)
	if err != nil {
		return nil, svcErr("ListEligibleBenefits", err)
	}
	eligible := make([]domain.Benefit, 0, len(allBenefits))
	for _, b := range allBenefits {
		rules, err := s.eligibilityRepo.ListRules(ctx, b.ID)
		if err != nil {
			return nil, svcErr("ListEligibleBenefits", err)
		}
		if len(rules) == 0 {
			eligible = append(eligible, b)
			continue
		}
		ok, _, err := s.engine.Evaluate(ctx, employeeID, rules)
		if err != nil {
			return nil, svcErr("ListEligibleBenefits", err)
		}
		if ok {
			eligible = append(eligible, b)
		}
	}
	return eligible, nil
}
