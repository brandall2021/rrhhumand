package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type AllowanceRepo struct {
	pool *pgxpool.Pool
}

func NewAllowanceRepo(pool *pgxpool.Pool) *AllowanceRepo {
	return &AllowanceRepo{pool: pool}
}

func (r *AllowanceRepo) Create(ctx context.Context, a *domain.DailyAllowanceRule) error {
	q := `INSERT INTO daily_allowance_rules (id,company_id,country,city,amount,currency,
		meal_included,accommodation_included,description,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.Country, a.City, a.Amount, a.Currency,
		a.MealIncluded, a.AccommodationIncluded, a.Description, a.IsActive, a.CreatedBy)
	return repoErr("AllowanceRepo.Create", err)
}

func (r *AllowanceRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.DailyAllowanceRule, error) {
	q := `SELECT id,company_id,country,city,amount,currency,
		meal_included,accommodation_included,description,is_active,created_by,created_at,updated_at
		FROM daily_allowance_rules WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a domain.DailyAllowanceRule
	err := row.Scan(&a.ID, &a.CompanyID, &a.Country, &a.City, &a.Amount, &a.Currency,
		&a.MealIncluded, &a.AccommodationIncluded, &a.Description, &a.IsActive, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("AllowanceRepo.Get", err)
	}
	return &a, nil
}

func (r *AllowanceRepo) List(ctx context.Context, companyID uuid.UUID, country, city *string) ([]domain.DailyAllowanceRule, error) {
	q := `SELECT id,company_id,country,city,amount,currency,
		meal_included,accommodation_included,description,is_active,created_by,created_at,updated_at
		FROM daily_allowance_rules WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if country != nil {
		q += fmt.Sprintf(" AND country=$%d", n)
		args = append(args, *country)
		n++
	}
	if city != nil {
		q += fmt.Sprintf(" AND city=$%d", n)
		args = append(args, *city)
		n++
	}
	q += " ORDER BY country,city"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("AllowanceRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.DailyAllowanceRule, error) {
		var a domain.DailyAllowanceRule
		err := row.Scan(&a.ID, &a.CompanyID, &a.Country, &a.City, &a.Amount, &a.Currency,
			&a.MealIncluded, &a.AccommodationIncluded, &a.Description, &a.IsActive, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *AllowanceRepo) Update(ctx context.Context, a *domain.DailyAllowanceRule) error {
	q := `UPDATE daily_allowance_rules SET country=$1,city=$2,amount=$3,currency=$4,
		meal_included=$5,accommodation_included=$6,description=$7,is_active=$8,updated_at=NOW()
		WHERE id=$9 AND company_id=$10`
	_, err := r.pool.Exec(ctx, q, a.Country, a.City, a.Amount, a.Currency,
		a.MealIncluded, a.AccommodationIncluded, a.Description, a.IsActive, a.ID, a.CompanyID)
	return repoErr("AllowanceRepo.Update", err)
}
