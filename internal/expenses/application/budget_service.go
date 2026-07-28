package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type BudgetRepository interface {
	Create(ctx context.Context, budget *domain.ExpenseBudget) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseBudget, error)
	List(ctx context.Context, companyID uuid.UUID, fiscalYear *int) ([]domain.ExpenseBudget, error)
	Update(ctx context.Context, budget *domain.ExpenseBudget) error
}

type BudgetService struct {
	budgetRepo BudgetRepository
}

func NewBudgetService(budgetRepo BudgetRepository) *BudgetService {
	return &BudgetService{budgetRepo: budgetRepo}
}

func (s *BudgetService) CreateBudget(ctx context.Context, companyID, userID uuid.UUID, b *domain.ExpenseBudget) (*domain.ExpenseBudget, error) {
	const op = "CreateBudget"
	now := time.Now()
	b.ID = uuid.New()
	b.CompanyID = companyID
	b.CreatedBy = userID
	b.CreatedAt = now
	b.UpdatedAt = now
	if err := s.budgetRepo.Create(ctx, b); err != nil {
		return nil, svcErr(op, err)
	}
	return b, nil
}

func (s *BudgetService) GetBudget(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseBudget, error) {
	const op = "GetBudget"
	budget, err := s.budgetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if budget.CompanyID != companyID {
		return nil, svcErr(op, domain.ErrNotFound)
	}
	return budget, nil
}

func (s *BudgetService) ListBudgets(ctx context.Context, companyID uuid.UUID, fiscalYear *int) ([]domain.ExpenseBudget, error) {
	const op = "ListBudgets"
	budgets, err := s.budgetRepo.List(ctx, companyID, fiscalYear)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return budgets, nil
}

func (s *BudgetService) CheckAvailability(ctx context.Context, companyID uuid.UUID, amount decimal.Decimal, costCenterID, projectID, categoryID *uuid.UUID) (bool, *decimal.Decimal, error) {
	const op = "CheckAvailability"
	budgets, err := s.budgetRepo.List(ctx, companyID, nil)
	if err != nil {
		return false, nil, svcErr(op, err)
	}

	var remaining decimal.Decimal
	for _, b := range budgets {
		if !b.IsActive {
			continue
		}
		if !matchesBudgetFilter(b, costCenterID, projectID, categoryID) {
			continue
		}
		avail := b.TotalAmount.Sub(b.UsedAmount).Sub(b.ReservedAmount)
		remaining = remaining.Add(avail)
	}

	available := amount.LessThanOrEqual(remaining)
	return available, &remaining, nil
}

func matchesBudgetFilter(b domain.ExpenseBudget, costCenterID, projectID, categoryID *uuid.UUID) bool {
	if costCenterID != nil && (b.CostCenterID == nil || *b.CostCenterID != *costCenterID) {
		return false
	}
	if projectID != nil && (b.ProjectID == nil || *b.ProjectID != *projectID) {
		return false
	}
	if categoryID != nil && (b.CategoryID == nil || *b.CategoryID != *categoryID) {
		return false
	}
	return true
}

func (s *BudgetService) ReserveAmount(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	const op = "ReserveAmount"
	budget, err := s.budgetRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	budget.ReservedAmount = budget.ReservedAmount.Add(amount)
	budget.UpdatedAt = time.Now()
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return svcErr(op, err)
	}
	return nil
}

func (s *BudgetService) ReleaseAmount(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	const op = "ReleaseAmount"
	budget, err := s.budgetRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	newReserved := budget.ReservedAmount.Sub(amount)
	if newReserved.IsNegative() {
		newReserved = decimal.Zero
	}
	budget.ReservedAmount = newReserved
	budget.UpdatedAt = time.Now()
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return svcErr(op, err)
	}
	return nil
}

func (s *BudgetService) ConsumeAmount(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	const op = "ConsumeAmount"
	budget, err := s.budgetRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	budget.UsedAmount = budget.UsedAmount.Add(amount)
	budget.UpdatedAt = time.Now()
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return svcErr(op, err)
	}
	return nil
}
