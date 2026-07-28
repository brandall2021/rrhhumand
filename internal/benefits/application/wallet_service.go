package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/shopspring/decimal"
)

type WalletService struct {
	walletRepo   *repository.WalletRepo
	flexibleRepo *repository.FlexibleRepo
}

func NewWalletService(walletRepo *repository.WalletRepo, flexibleRepo *repository.FlexibleRepo) *WalletService {
	return &WalletService{
		walletRepo:   walletRepo,
		flexibleRepo: flexibleRepo,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, companyID, employeeID uuid.UUID, walletType string, balance decimal.Decimal, currency string) (*domain.EmployeeBenefitWallet, error) {
	w := &domain.EmployeeBenefitWallet{
		ID:         uuid.New(),
		CompanyID:  companyID,
		EmployeeID: employeeID,
		WalletType: walletType,
		Balance:    balance,
		Currency:   currency,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := s.walletRepo.Create(ctx, w); err != nil {
		return nil, svcErr("CreateWallet", err)
	}
	return w, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (*domain.EmployeeBenefitWallet, error) {
	return s.walletRepo.Get(ctx, id)
}

func (s *WalletService) GetEmployeeWallet(ctx context.Context, companyID, employeeID uuid.UUID, walletType string) (*domain.EmployeeBenefitWallet, error) {
	return s.walletRepo.GetByEmployeeAndType(ctx, employeeID, walletType)
}

func (s *WalletService) ListEmployeeWallets(ctx context.Context, companyID, employeeID uuid.UUID) ([]domain.EmployeeBenefitWallet, error) {
	return s.walletRepo.List(ctx, employeeID)
}

func (s *WalletService) CreditWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, txType, description string, createdBy *uuid.UUID) error {
	w, err := s.walletRepo.Get(ctx, id)
	if err != nil {
		return svcErr("CreditWallet", err)
	}
	balanceBefore := w.Balance
	balanceAfter := balanceBefore.Add(amount)
	now := time.Now()
	tx := &domain.BenefitWalletTransaction{
		ID:              uuid.New(),
		WalletID:        id,
		TransactionType: txType,
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Currency:        w.Currency,
		Description:     &description,
		TransactionDate: now,
		CreatedBy:       createdBy,
		CreatedAt:       now,
	}
	if err := s.walletRepo.UpdateBalance(ctx, id, balanceAfter); err != nil {
		return svcErr("CreditWallet", err)
	}
	if err := s.walletRepo.CreateTransaction(ctx, tx); err != nil {
		return svcErr("CreditWallet", err)
	}
	return nil
}

func (s *WalletService) DebitWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, txType, description string, createdBy *uuid.UUID) error {
	w, err := s.walletRepo.Get(ctx, id)
	if err != nil {
		return svcErr("DebitWallet", err)
	}
	if w.Balance.LessThan(amount) {
		return svcErr("DebitWallet", domain.ErrBalanceExceeded)
	}
	balanceBefore := w.Balance
	balanceAfter := balanceBefore.Sub(amount)
	now := time.Now()
	tx := &domain.BenefitWalletTransaction{
		ID:              uuid.New(),
		WalletID:        id,
		TransactionType: txType,
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Currency:        w.Currency,
		Description:     &description,
		TransactionDate: now,
		CreatedBy:       createdBy,
		CreatedAt:       now,
	}
	if err := s.walletRepo.UpdateBalance(ctx, id, balanceAfter); err != nil {
		return svcErr("DebitWallet", err)
	}
	if err := s.walletRepo.CreateTransaction(ctx, tx); err != nil {
		return svcErr("DebitWallet", err)
	}
	return nil
}

func (s *WalletService) ListTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]domain.BenefitWalletTransaction, error) {
	return s.walletRepo.ListTransactions(ctx, walletID, limit, offset)
}

func (s *WalletService) CreateFlexiblePlan(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, p *domain.BenefitFlexiblePlan) (*domain.BenefitFlexiblePlan, error) {
	p.ID = uuid.New()
	p.CompanyID = companyID
	p.CreatedBy = userID
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err := s.flexibleRepo.CreatePlan(ctx, p); err != nil {
		return nil, svcErr("CreateFlexiblePlan", err)
	}
	return p, nil
}

func (s *WalletService) ListFlexiblePlans(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitFlexiblePlan, error) {
	return s.flexibleRepo.ListPlans(ctx, companyID)
}

func (s *WalletService) CreateBudget(ctx context.Context, companyID, employeeID, planID uuid.UUID, fiscalYear int, totalAmount decimal.Decimal) (*domain.BenefitFlexibleBudget, error) {
	b := &domain.BenefitFlexibleBudget{
		ID:             uuid.New(),
		CompanyID:      companyID,
		EmployeeID:     employeeID,
		FlexiblePlanID: planID,
		FiscalYear:     fiscalYear,
		TotalAmount:    totalAmount,
		UsedAmount:     decimal.Zero,
		PendingAmount:  totalAmount,
		Currency:       "ARS",
		Status:         "ACTIVE",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.flexibleRepo.CreateBudget(ctx, b); err != nil {
		return nil, svcErr("CreateBudget", err)
	}
	return b, nil
}

func (s *WalletService) GetBudget(ctx context.Context, companyID, employeeID, planID uuid.UUID, fiscalYear int) (*domain.BenefitFlexibleBudget, error) {
	return s.flexibleRepo.GetBudget(ctx, employeeID, planID, fiscalYear)
}

func (s *WalletService) ListEmployeeBudgets(ctx context.Context, companyID, employeeID uuid.UUID) ([]domain.BenefitFlexibleBudget, error) {
	return s.flexibleRepo.ListBudgets(ctx, &employeeID, nil)
}
