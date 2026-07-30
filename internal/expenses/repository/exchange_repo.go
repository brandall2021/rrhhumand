package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ExchangeRepo struct {
	pool *pgxpool.Pool
}

func NewExchangeRepo(pool *pgxpool.Pool) *ExchangeRepo {
	return &ExchangeRepo{pool: pool}
}

func (r *ExchangeRepo) Create(ctx context.Context, e *domain.ExchangeRate) error {
	q := `INSERT INTO exchange_rates (id,from_currency,to_currency,rate,rate_date,source)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.FromCurrency, e.ToCurrency, e.Rate, e.EffectiveDate, e.Source)
	return repoErr("ExchangeRepo.Create", err)
}

func (r *ExchangeRepo) GetLatest(ctx context.Context, fromCurrency, toCurrency string) (*domain.ExchangeRate, error) {
	q := `SELECT id,from_currency,to_currency,rate,rate_date,source,created_at
		FROM exchange_rates WHERE from_currency=$1 AND to_currency=$2 ORDER BY rate_date DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, fromCurrency, toCurrency)
	var e domain.ExchangeRate
	err := row.Scan(&e.ID, &e.FromCurrency, &e.ToCurrency, &e.Rate, &e.EffectiveDate, &e.Source, &e.CreatedAt)
	if err != nil {
		return nil, repoErr("ExchangeRepo.GetLatest", err)
	}
	return &e, nil
}

func (r *ExchangeRepo) GetByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*domain.ExchangeRate, error) {
	q := `SELECT id,from_currency,to_currency,rate,rate_date,source,created_at
		FROM exchange_rates WHERE from_currency=$1 AND to_currency=$2 AND rate_date=$3
		ORDER BY rate_date DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, fromCurrency, toCurrency, date)
	var e domain.ExchangeRate
	err := row.Scan(&e.ID, &e.FromCurrency, &e.ToCurrency, &e.Rate, &e.EffectiveDate, &e.Source, &e.CreatedAt)
	if err != nil {
		return nil, repoErr("ExchangeRepo.GetByDate", err)
	}
	return &e, nil
}

func (r *ExchangeRepo) ListByDateRange(ctx context.Context, fromCurrency, toCurrency string, from, to time.Time) ([]domain.ExchangeRate, error) {
	q := `SELECT id,from_currency,to_currency,rate,rate_date,source,created_at
		FROM exchange_rates
		WHERE from_currency=$1 AND to_currency=$2 AND rate_date>=$3 AND rate_date<=$4
		ORDER BY rate_date`
	rows, err := r.pool.Query(ctx, q, fromCurrency, toCurrency, from, to)
	if err != nil {
		return nil, repoErr("ExchangeRepo.ListByDateRange", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExchangeRate, error) {
		var e domain.ExchangeRate
		err := row.Scan(&e.ID, &e.FromCurrency, &e.ToCurrency, &e.Rate, &e.EffectiveDate, &e.Source, &e.CreatedAt)
		return e, err
	})
}
