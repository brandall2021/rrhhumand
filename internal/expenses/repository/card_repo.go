package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type CardRepo struct {
	pool *pgxpool.Pool
}

func NewCardRepo(pool *pgxpool.Pool) *CardRepo {
	return &CardRepo{pool: pool}
}

func (r *CardRepo) CreateCard(ctx context.Context, c *domain.CorporateCard) error {
	q := `INSERT INTO corporate_cards (id,company_id,employee_id,card_number_masked,cardholder_name,
		provider,credit_limit,currency,expiration_date,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.EmployeeID, c.CardNumberMasked, c.CardholderName,
		c.Provider, c.CreditLimit, c.Currency, c.ExpirationDate, c.IsActive, c.CreatedBy)
	return repoErr("CreateCard", err)
}

func (r *CardRepo) GetCard(ctx context.Context, companyID, id uuid.UUID) (*domain.CorporateCard, error) {
	q := `SELECT id,company_id,employee_id,card_number_masked,cardholder_name,
		provider,credit_limit,currency,expiration_date,is_active,created_by,created_at,updated_at
		FROM corporate_cards WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.CorporateCard
	err := row.Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.CardNumberMasked, &c.CardholderName,
		&c.Provider, &c.CreditLimit, &c.Currency, &c.ExpirationDate, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCard", err)
	}
	return &c, nil
}

func (r *CardRepo) ListCards(ctx context.Context, companyID uuid.UUID) ([]domain.CorporateCard, error) {
	q := `SELECT id,company_id,employee_id,card_number_masked,cardholder_name,
		provider,credit_limit,currency,expiration_date,is_active,created_by,created_at,updated_at
		FROM corporate_cards WHERE company_id=$1 ORDER BY cardholder_name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCards", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.CorporateCard, error) {
		var c domain.CorporateCard
		err := row.Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.CardNumberMasked, &c.CardholderName,
			&c.Provider, &c.CreditLimit, &c.Currency, &c.ExpirationDate, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *CardRepo) UpdateCard(ctx context.Context, c *domain.CorporateCard) error {
	q := `UPDATE corporate_cards SET employee_id=$1,cardholder_name=$2,provider=$3,
		credit_limit=$4,currency=$5,expiration_date=$6,is_active=$7,updated_at=NOW()
		WHERE id=$8 AND company_id=$9`
	_, err := r.pool.Exec(ctx, q, c.EmployeeID, c.CardholderName, c.Provider,
		c.CreditLimit, c.Currency, c.ExpirationDate, c.IsActive, c.ID, c.CompanyID)
	return repoErr("UpdateCard", err)
}

func (r *CardRepo) CreateTransaction(ctx context.Context, t *domain.CorporateCardTransaction) error {
	q := `INSERT INTO corporate_card_transactions (id,card_id,company_id,expense_id,transaction_date,
		merchant_name,amount,currency,reference,status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CardID, t.CompanyID, t.ExpenseID, t.TransactionDate,
		t.MerchantName, t.Amount, t.Currency, t.Reference, t.Status)
	return repoErr("CreateTransaction", err)
}

func (r *CardRepo) ListTransactions(ctx context.Context, cardID uuid.UUID) ([]domain.CorporateCardTransaction, error) {
	q := `SELECT id,card_id,company_id,expense_id,transaction_date,merchant_name,amount,currency,reference,
		status,created_at
		FROM corporate_card_transactions WHERE card_id=$1 ORDER BY transaction_date DESC`
	rows, err := r.pool.Query(ctx, q, cardID)
	if err != nil {
		return nil, repoErr("ListTransactions", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.CorporateCardTransaction, error) {
		var t domain.CorporateCardTransaction
		err := row.Scan(&t.ID, &t.CardID, &t.CompanyID, &t.ExpenseID, &t.TransactionDate,
			&t.MerchantName, &t.Amount, &t.Currency, &t.Reference, &t.Status, &t.CreatedAt)
		return t, err
	})
}

func (r *CardRepo) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string, expenseID *uuid.UUID) error {
	q := `UPDATE corporate_card_transactions SET status=$1,expense_id=$2 WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, status, expenseID, id)
	return repoErr("UpdateTransactionStatus", err)
}
