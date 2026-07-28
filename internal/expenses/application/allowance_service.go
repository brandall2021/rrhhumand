package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type AllowanceRepository interface {
	CreateRule(ctx context.Context, rule *domain.DailyAllowanceRule) error
	GetRule(ctx context.Context, companyID, id uuid.UUID) (*domain.DailyAllowanceRule, error)
	ListRules(ctx context.Context, companyID uuid.UUID) ([]domain.DailyAllowanceRule, error)
	UpdateRule(ctx context.Context, rule *domain.DailyAllowanceRule) error
}

type AllowanceService struct {
	allowanceRepo AllowanceRepository
}

func NewAllowanceService(allowanceRepo AllowanceRepository) *AllowanceService {
	return &AllowanceService{allowanceRepo: allowanceRepo}
}

func (s *AllowanceService) CreateRule(ctx context.Context, companyID, userID uuid.UUID, rule *domain.DailyAllowanceRule) (*domain.DailyAllowanceRule, error) {
	const op = "CreateRule"
	now := time.Now()
	rule.ID = uuid.New()
	rule.CompanyID = companyID
	rule.CreatedBy = userID
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.allowanceRepo.CreateRule(ctx, rule); err != nil {
		return nil, svcErr(op, err)
	}
	return rule, nil
}

func (s *AllowanceService) GetRule(ctx context.Context, companyID, id uuid.UUID) (*domain.DailyAllowanceRule, error) {
	const op = "GetRule"
	rule, err := s.allowanceRepo.GetRule(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return rule, nil
}

func (s *AllowanceService) ListRules(ctx context.Context, companyID uuid.UUID) ([]domain.DailyAllowanceRule, error) {
	const op = "ListRules"
	rules, err := s.allowanceRepo.ListRules(ctx, companyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return rules, nil
}

func (s *AllowanceService) UpdateRule(ctx context.Context, companyID uuid.UUID, rule *domain.DailyAllowanceRule) (*domain.DailyAllowanceRule, error) {
	const op = "UpdateRule"
	existing, err := s.allowanceRepo.GetRule(ctx, companyID, rule.ID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	rule.CompanyID = companyID
	rule.CreatedAt = existing.CreatedAt
	rule.CreatedBy = existing.CreatedBy
	rule.UpdatedAt = time.Now()
	if err := s.allowanceRepo.UpdateRule(ctx, rule); err != nil {
		return nil, svcErr(op, err)
	}
	return rule, nil
}

func (s *AllowanceService) CalculateAllowance(ctx context.Context, companyID uuid.UUID, destination string, days int, employeeCategory string) (dailyAmount decimal.Decimal, totalAmount decimal.Decimal, err error) {
	const op = "CalculateAllowance"
	rules, err := s.allowanceRepo.ListRules(ctx, companyID)
	if err != nil {
		return decimal.Zero, decimal.Zero, svcErr(op, err)
	}

	var matchedRule *domain.DailyAllowanceRule
	for _, r := range rules {
		if !r.IsActive {
			continue
		}
		if time.Now().Before(r.EffectiveFrom) {
			continue
		}
		if r.EffectiveTo != nil && time.Now().After(*r.EffectiveTo) {
			continue
		}
		if r.EmployeeCategory != nil && *r.EmployeeCategory != employeeCategory {
			continue
		}
		if r.Country != nil && *r.Country != destination {
			continue
		}
		matchedRule = &r
		break
	}

	if matchedRule == nil {
		return decimal.Zero, decimal.Zero, svcErr(op, domain.ErrNotFound)
	}

	daily := matchedRule.DailyAmount
	total := daily.Mul(decimal.NewFromInt(int64(days)))
	return daily, total, nil
}
