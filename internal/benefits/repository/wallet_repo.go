package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/shopspring/decimal"
)

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepo(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

func (r *WalletRepo) Create(ctx context.Context, w *domain.EmployeeBenefitWallet) error {
	q := `INSERT INTO employee_benefit_wallets (id,company_id,employee_id,benefit_id,wallet_type,balance,currency,is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, w.ID, w.CompanyID, w.EmployeeID, w.BenefitID, w.WalletType, w.Balance, w.Currency, w.IsActive)
	return repoErr("Create", err)
}

func (r *WalletRepo) Get(ctx context.Context, id uuid.UUID) (*domain.EmployeeBenefitWallet, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,wallet_type,balance,currency,last_transaction_at,is_active,created_at,updated_at
		FROM employee_benefit_wallets WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var w domain.EmployeeBenefitWallet
	err := row.Scan(&w.ID, &w.CompanyID, &w.EmployeeID, &w.BenefitID, &w.WalletType, &w.Balance, &w.Currency, &w.LastTransactionAt, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, repoErr("Get", err)
	}
	return &w, nil
}

func (r *WalletRepo) GetByEmployeeAndType(ctx context.Context, employeeID uuid.UUID, walletType string) (*domain.EmployeeBenefitWallet, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,wallet_type,balance,currency,last_transaction_at,is_active,created_at,updated_at
		FROM employee_benefit_wallets WHERE employee_id=$1 AND wallet_type=$2 AND is_active=true`
	row := r.pool.QueryRow(ctx, q, employeeID, walletType)
	var w domain.EmployeeBenefitWallet
	err := row.Scan(&w.ID, &w.CompanyID, &w.EmployeeID, &w.BenefitID, &w.WalletType, &w.Balance, &w.Currency, &w.LastTransactionAt, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetByEmployeeAndType", err)
	}
	return &w, nil
}

func (r *WalletRepo) List(ctx context.Context, employeeID uuid.UUID) ([]domain.EmployeeBenefitWallet, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,wallet_type,balance,currency,last_transaction_at,is_active,created_at,updated_at
		FROM employee_benefit_wallets WHERE employee_id=$1 ORDER BY wallet_type`
	rows, err := r.pool.Query(ctx, q, employeeID)
	if err != nil {
		return nil, repoErr("List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.EmployeeBenefitWallet, error) {
		var w domain.EmployeeBenefitWallet
		err := row.Scan(&w.ID, &w.CompanyID, &w.EmployeeID, &w.BenefitID, &w.WalletType, &w.Balance, &w.Currency, &w.LastTransactionAt, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
		return w, err
	})
}

func (r *WalletRepo) UpdateBalance(ctx context.Context, walletID uuid.UUID, balance decimal.Decimal) error {
	q := `UPDATE employee_benefit_wallets SET balance=$1,last_transaction_at=NOW(),updated_at=NOW() WHERE id=$2`
	_, err := r.pool.Exec(ctx, q, balance, walletID)
	return repoErr("UpdateBalance", err)
}

func (r *WalletRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE employee_benefit_wallets SET is_active=false,updated_at=NOW() WHERE id=$1`, id)
	return repoErr("Deactivate", err)
}

func (r *WalletRepo) CreateTransaction(ctx context.Context, t *domain.BenefitWalletTransaction) error {
	q := `INSERT INTO benefit_wallet_transactions (id,wallet_id,transaction_type,amount,balance_before,balance_after,
		currency,reference_type,reference_id,description,receipt_url,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.WalletID, t.TransactionType, t.Amount, t.BalanceBefore, t.BalanceAfter,
		t.Currency, t.ReferenceType, t.ReferenceID, t.Description, t.ReceiptURL, t.CreatedBy)
	return repoErr("CreateTransaction", err)
}

func (r *WalletRepo) ListTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]domain.BenefitWalletTransaction, error) {
	q := `SELECT id,wallet_id,transaction_type,amount,balance_before,balance_after,currency,reference_type,
		reference_id,description,receipt_url,transaction_date,created_by,created_at
		FROM benefit_wallet_transactions WHERE wallet_id=$1 ORDER BY transaction_date DESC`
	args := []any{walletID}
	n := 2
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListTransactions", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitWalletTransaction, error) {
		var t domain.BenefitWalletTransaction
		err := row.Scan(&t.ID, &t.WalletID, &t.TransactionType, &t.Amount, &t.BalanceBefore, &t.BalanceAfter,
			&t.Currency, &t.ReferenceType, &t.ReferenceID, &t.Description, &t.ReceiptURL, &t.TransactionDate, &t.CreatedBy, &t.CreatedAt)
		return t, err
	})
}
