package engine

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/rrhhumand/api/internal/expenses/repository"
	"github.com/shopspring/decimal"
)

type PolicyEngine struct {
	policyRepo *repository.PolicyRepo
}

func NewPolicyEngine(pr *repository.PolicyRepo) *PolicyEngine {
	return &PolicyEngine{policyRepo: pr}
}

func engErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("expenses_engine.%s: %w", op, err)
}

func (e *PolicyEngine) Evaluate(ctx context.Context, expense domain.ExpenseContext) (domain.PolicyResult, error) {
	companyID, err := uuid.Parse(expense.CompanyID)
	if err != nil {
		return domain.PolicyResult{}, engErr("Evaluate.parseCompanyID", err)
	}

	rules, err := e.policyRepo.GetActiveRules(ctx, companyID)
	if err != nil {
		return domain.PolicyResult{}, engErr("Evaluate.getRules", err)
	}

	var violations []domain.PolicyViolation
	var reasons []string

	for _, rule := range rules {
		if rule.CategoryID != nil && rule.CategoryID.String() != expense.Category {
			continue
		}
		if rule.EmployeeCategory != nil && *rule.EmployeeCategory != expense.EmployeeCategory {
			continue
		}

		if v := e.CheckMaxAmount(rule, expense.Amount); v != nil {
			violations = append(violations, *v)
			reasons = append(reasons, v.Message)
		}
		if v := e.CheckReceiptRequirement(rule, expense.HasReceipt); v != nil {
			violations = append(violations, *v)
			reasons = append(reasons, v.Message)
		}
	}

	return domain.PolicyResult{
		Compliant:  len(violations) == 0,
		Reasons:    reasons,
		Violations: violations,
	}, nil
}

func (e *PolicyEngine) CheckMaxAmount(rule domain.ExpensePolicyRule, amount decimal.Decimal) *domain.PolicyViolation {
	if rule.MaxAmount == nil {
		return nil
	}
	if amount.GreaterThan(*rule.MaxAmount) {
		excess := amount.Sub(*rule.MaxAmount)
		return &domain.PolicyViolation{
			RuleID:  rule.ID.String(),
			Limit:   *rule.MaxAmount,
			Actual:  amount,
			Excess:  excess,
			Message: fmt.Sprintf("amount %.2f exceeds max %.2f by %.2f", amount, *rule.MaxAmount, excess),
		}
	}
	return nil
}

func (e *PolicyEngine) CheckReceiptRequirement(rule domain.ExpensePolicyRule, hasReceipt bool) *domain.PolicyViolation {
	if !rule.RequiresReceipt {
		return nil
	}
	if !hasReceipt {
		return &domain.PolicyViolation{
			RuleID:  rule.ID.String(),
			Message: "receipt is required but not provided",
		}
	}
	return nil
}
