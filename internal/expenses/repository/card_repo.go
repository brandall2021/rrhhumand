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
	q := `INSERT INTO corporate_cards (id,company_id,employee_id,card_number,card_holder_name,
		card_type,issuer,expiry_month,expiry_year,credit_limit,currency,status,issued_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.EmployeeID, c.CardNumber, c.CardHolderName,
		c.CardType, c.Issuer, c.ExpiryMonth, c.ExpiryYear, c.CreditLimit, c.Currency, c.Status, c.IssuedBy)
	return repoErr("CreateCard", err)
}

func (r *CardRepo) GetCard(ctx context.Context, companyID, id uuid.UUID) (*domain.CorporateCard, error) {
	q := `SELECT id,company_id,employee_id,card_number,card_holder_name,
		card_type,issuer,expiry_month,expiry_year,credit_limit,currency,status,issued_by,created_at,updated_at
		FROM corporate_cards WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.CorporateCard
	err := row.Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.CardNumber, &c.CardHolderName,
		&c.CardType, &c.Issuer, &c.ExpiryMonth, &c.ExpiryYear, &c.CreditLimit, &c.Currency, &c.Status, &c.IssuedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCard", err)
	}
	return &c, nil
}

func (r *CardRepo) ListCards(ctx context.Context, companyID uuid.UUID) ([]domain.CorporateCard, error) {
	q := `SELECT id,company_id,employee_id,card_number,card_holder_name,
		card_type,issuer,expiry_month,expiry_year,credit_limit,currency,status,issued_by,created_at,updated_at
		FROM corporate_cards WHERE company_id=$1 ORDER BY card_holder_name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCards", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.CorporateCard, error) {
		var c domain.CorporateCard
		err := row.Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.CardNumber, &c.CardHolderName,
			&c.CardType, &c.Issuer, &c.ExpiryMonth, &c.ExpiryYear, &c.CreditLimit, &c.Currency, &c.Status, &c.IssuedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *CardRepo) UpdateCard(ctx context.Context, c *domain.CorporateCard) error {
	q := `UPDATE corporate_cards SET employee_id=$1,card_holder_name=$2,card_type=$3,issuer=$4,
		expiry_month=$5,expiry_year=$6,credit_limit=$7,currency=$8,status=$9,updated_at=NOW()
		WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, c.EmployeeID, c.CardHolderName, c.CardType, c.Issuer,
		c.ExpiryMonth, c.ExpiryYear, c.CreditLimit, c.Currency, c.Status, c.ID, c.CompanyID)
	return repoErr("UpdateCard", err)
}

func (r *CardRepo) CreateTransaction(ctx context.Context, t *domain.CorporateCardTransaction) error {
	q := `INSERT INTO corporate_card_transactions (id,card_id,merchant,amount,currency,transaction_date,
		description,category,reference,status,matched_expense_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CardID, t.Merchant, t.Amount, t.Currency, t.TransactionDate,
		t.Description, t.Category, t.Reference, t.Status, t.MatchedExpenseID)
	return repoErr("CreateTransaction", err)
}

func (r *CardRepo) ListTransactions(ctx context.Context, cardID uuid.UUID) ([]domain.CorporateCardTransaction, error) {
	q := `SELECT id,card_id,merchant,amount,currency,transaction_date,description,category,reference,
		status,matched_expense_id,created_at,updated_at
		FROM corporate_card_transactions WHERE card_id=$1 ORDER BY transaction_date DESC`
	rows, err := r.pool.Query(ctx, q, cardID)
	if err != nil {
		return nil, repoErr("ListTransactions", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.CorporateCardTransaction, error) {
		var t domain.CorporateCardTransaction
		err := row.Scan(&t.ID, &t.CardID, &t.Merchant, &t.Amount, &t.Currency, &t.TransactionDate,
			&t.Description, &t.Category, &t.Reference, &t.Status, &t.MatchedExpenseID, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *CardRepo) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string, matchedExpenseID *uuid.UUID) error {
	q := `UPDATE corporate_card_transactions SET status=$1,matched_expense_id=$2,updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, status, matchedExpenseID, id)
	return repoErr("UpdateTransactionStatus", err)
}
