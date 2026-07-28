package overtime

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type BalanceManager struct {
	repo *Repository
}

func NewBalanceManager(repo *Repository) *BalanceManager {
	return &BalanceManager{repo: repo}
}

func (bm *BalanceManager) GetBalance(ctx context.Context, companyID, employeeID string) (*EmployeeTimeBalance, error) {
	balance, err := bm.repo.GetTimeBalance(ctx, companyID, employeeID)
	if err != nil {
		balance = &EmployeeTimeBalance{
			ID:             uuid.New().String(),
			CompanyID:      companyID,
			EmployeeID:     employeeID,
			BalanceMinutes: 0,
		}
		bm.repo.CreateTimeBalance(ctx, balance)
	}
	return balance, nil
}

func (bm *BalanceManager) Credit(ctx context.Context, companyID, employeeID string, minutes int, txType string, overtimeRecordID *string, reason *string, createdBy *string) error {
	balance, err := bm.GetBalance(ctx, companyID, employeeID)
	if err != nil {
		return err
	}

	tx := &TimeBalanceTransaction{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		EmployeeID:       employeeID,
		OvertimeRecordID: overtimeRecordID,
		TransactionType:  txType,
		Minutes:          minutes,
		Reason:           reason,
		CreatedBy:        createdBy,
	}

	if err := bm.repo.CreateBalanceTransaction(ctx, tx); err != nil {
		return err
	}

	balance.BalanceMinutes += minutes
	return bm.repo.UpdateTimeBalance(ctx, balance)
}

func (bm *BalanceManager) Debit(ctx context.Context, companyID, employeeID string, minutes int, txType string, reason *string, createdBy *string) error {
	balance, err := bm.GetBalance(ctx, companyID, employeeID)
	if err != nil {
		return err
	}

	if balance.BalanceMinutes < minutes {
		return fmt.Errorf("insufficient balance: available %d, requested %d", balance.BalanceMinutes, minutes)
	}

	tx := &TimeBalanceTransaction{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		EmployeeID:      employeeID,
		TransactionType: txType,
		Minutes:         -minutes,
		Reason:          reason,
		CreatedBy:       createdBy,
	}

	if err := bm.repo.CreateBalanceTransaction(ctx, tx); err != nil {
		return err
	}

	balance.BalanceMinutes -= minutes
	return bm.repo.UpdateTimeBalance(ctx, balance)
}

func (bm *BalanceManager) Adjust(ctx context.Context, companyID, employeeID string, minutes int, reason string, createdBy string) error {
	balance, err := bm.GetBalance(ctx, companyID, employeeID)
	if err != nil {
		return err
	}

	txType := "MANUAL_CREDIT"
	min := minutes
	if minutes < 0 {
		txType = "MANUAL_DEBIT"
		min = -minutes
	}

	tx := &TimeBalanceTransaction{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		EmployeeID:      employeeID,
		TransactionType: txType,
		Minutes:         min,
		Reason:          &reason,
		CreatedBy:       &createdBy,
	}

	if err := bm.repo.CreateBalanceTransaction(ctx, tx); err != nil {
		return err
	}

	balance.BalanceMinutes += minutes
	return bm.repo.UpdateTimeBalance(ctx, balance)
}

func (bm *BalanceManager) GetTransactions(ctx context.Context, companyID, employeeID string) ([]TimeBalanceTransaction, error) {
	return bm.repo.ListBalanceTransactions(ctx, companyID, employeeID)
}
