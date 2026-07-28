package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type PolicyRepository interface {
	CreatePolicy(ctx context.Context, policy *domain.ExpensePolicy) error
	GetPolicy(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpensePolicy, error)
	ListPolicies(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePolicy, error)
	UpdatePolicy(ctx context.Context, policy *domain.ExpensePolicy) error

	CreateRule(ctx context.Context, rule *domain.ExpensePolicyRule) error
	ListRules(ctx context.Context, policyID uuid.UUID) ([]domain.ExpensePolicyRule, error)
	UpdateRule(ctx context.Context, rule *domain.ExpensePolicyRule) error
	DeleteRule(ctx context.Context, policyID, ruleID uuid.UUID) error
}

type PolicyService struct {
	policyRepo  PolicyRepository
	expenseRepo ExpenseRepository
}

func NewPolicyService(policyRepo PolicyRepository, expenseRepo ExpenseRepository) *PolicyService {
	return &PolicyService{policyRepo: policyRepo, expenseRepo: expenseRepo}
}

func (s *PolicyService) CreatePolicy(ctx context.Context, companyID, userID uuid.UUID, p *domain.ExpensePolicy) (*domain.ExpensePolicy, error) {
	const op = "CreatePolicy"
	now := time.Now()
	p.ID = uuid.New()
	p.CompanyID = companyID
	p.CreatedBy = userID
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := s.policyRepo.CreatePolicy(ctx, p); err != nil {
		return nil, svcErr(op, err)
	}
	return p, nil
}

func (s *PolicyService) GetPolicy(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpensePolicy, error) {
	const op = "GetPolicy"
	policy, err := s.policyRepo.GetPolicy(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return policy, nil
}

func (s *PolicyService) ListPolicies(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePolicy, error) {
	const op = "ListPolicies"
	policies, err := s.policyRepo.ListPolicies(ctx, companyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return policies, nil
}

func (s *PolicyService) UpdatePolicy(ctx context.Context, companyID uuid.UUID, p *domain.ExpensePolicy) (*domain.ExpensePolicy, error) {
	const op = "UpdatePolicy"
	existing, err := s.policyRepo.GetPolicy(ctx, companyID, p.ID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	p.CompanyID = companyID
	p.CreatedAt = existing.CreatedAt
	p.CreatedBy = existing.CreatedBy
	p.UpdatedAt = time.Now()
	if err := s.policyRepo.UpdatePolicy(ctx, p); err != nil {
		return nil, svcErr(op, err)
	}
	return p, nil
}

func (s *PolicyService) CreateRule(ctx context.Context, policyID uuid.UUID, rule *domain.ExpensePolicyRule) (*domain.ExpensePolicyRule, error) {
	const op = "CreateRule"
	now := time.Now()
	rule.ID = uuid.New()
	rule.PolicyID = policyID
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.policyRepo.CreateRule(ctx, rule); err != nil {
		return nil, svcErr(op, err)
	}
	return rule, nil
}

func (s *PolicyService) ListRules(ctx context.Context, policyID uuid.UUID) ([]domain.ExpensePolicyRule, error) {
	const op = "ListRules"
	rules, err := s.policyRepo.ListRules(ctx, policyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return rules, nil
}

func (s *PolicyService) UpdateRule(ctx context.Context, rule *domain.ExpensePolicyRule) (*domain.ExpensePolicyRule, error) {
	const op = "UpdateRule"
	rule.UpdatedAt = time.Now()
	if err := s.policyRepo.UpdateRule(ctx, rule); err != nil {
		return nil, svcErr(op, err)
	}
	return rule, nil
}

func (s *PolicyService) DeleteRule(ctx context.Context, policyID, ruleID uuid.UUID) error {
	const op = "DeleteRule"
	if err := s.policyRepo.DeleteRule(ctx, policyID, ruleID); err != nil {
		return svcErr(op, err)
	}
	return nil
}

func (s *PolicyService) EvaluateExpense(ctx context.Context, companyID uuid.UUID, expense *domain.Expense, employeeCategory string) (domain.PolicyResult, error) {
	const op = "EvaluateExpense"
	policies, err := s.policyRepo.ListPolicies(ctx, companyID)
	if err != nil {
		return domain.PolicyResult{}, svcErr(op, err)
	}

	var allViolations []domain.PolicyViolation
	compliant := true

	for _, policy := range policies {
		if !policy.IsActive {
			continue
		}
		if time.Now().Before(policy.EffectiveFrom) {
			continue
		}
		if policy.EffectiveTo != nil && time.Now().After(*policy.EffectiveTo) {
			continue
		}

		rules, err := s.policyRepo.ListRules(ctx, policy.ID)
		if err != nil {
			return domain.PolicyResult{}, svcErr(op, err)
		}

		for _, rule := range rules {
			if !rule.IsActive {
				continue
			}
			if rule.CategoryID != nil && *rule.CategoryID != expense.CategoryID {
				continue
			}
			if rule.EmployeeCategory != nil && *rule.EmployeeCategory != employeeCategory {
				continue
			}
			if rule.Currency != nil && *rule.Currency != expense.BaseCurrency {
				continue
			}
			if rule.MaxAmount != nil && expense.BaseAmount.GreaterThan(*rule.MaxAmount) {
				compliant = false
				excess := expense.BaseAmount.Sub(*rule.MaxAmount)
				allViolations = append(allViolations, domain.PolicyViolation{
					RuleID:  rule.ID.String(),
					Limit:   *rule.MaxAmount,
					Actual:  expense.BaseAmount,
					Excess:  excess,
					Message: "expense exceeds maximum allowed amount",
				})
			}
			if rule.RequiresReceipt && !expense.IsReimbursable {
				compliant = false
				allViolations = append(allViolations, domain.PolicyViolation{
					RuleID:  rule.ID.String(),
					Message: "receipt is required but not provided",
				})
			}
		}
	}

	reasons := make([]string, 0, len(allViolations))
	for _, v := range allViolations {
		reasons = append(reasons, v.Message)
	}

	return domain.PolicyResult{
		Compliant:  compliant,
		Reasons:    reasons,
		Violations: allViolations,
	}, nil
}

func (s *PolicyService) OverridePolicy(ctx context.Context, expenseID uuid.UUID, reason string, overriddenBy uuid.UUID) error {
	const op = "OverridePolicy"
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return svcErr(op, err)
	}
	expense.IsPolicyCompliant = true
	expense.PolicyStatus = "OVERRIDDEN"
	expense.PolicyOverrideReason = &reason
	expense.PolicyOverrideBy = &overriddenBy
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	return nil
}
