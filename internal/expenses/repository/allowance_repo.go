package repository

import (
	"context"

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

func (r *AllowanceRepo) CreateRule(ctx context.Context, a *domain.DailyAllowanceRule) error {
	q := `INSERT INTO daily_allowance_rules (id,company_id,name,country,region,city,employee_category,daily_amount,currency,
		meal_percentage,lodging_percentage,transport_percentage,other_percentage,effective_from,effective_to,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.Name, a.Country, a.Region, a.City, a.EmployeeCategory, a.DailyAmount, a.Currency,
		a.MealPercentage, a.LodgingPercentage, a.TransportPercentage, a.OtherPercentage, a.EffectiveFrom, a.EffectiveTo, a.IsActive, a.CreatedBy)
	return repoErr("AllowanceRepo.Create", err)
}

func (r *AllowanceRepo) GetRule(ctx context.Context, companyID, id uuid.UUID) (*domain.DailyAllowanceRule, error) {
	q := `SELECT id,company_id,name,country,region,city,employee_category,daily_amount,currency,
		meal_percentage,lodging_percentage,transport_percentage,other_percentage,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM daily_allowance_rules WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a domain.DailyAllowanceRule
	err := row.Scan(&a.ID, &a.CompanyID, &a.Name, &a.Country, &a.Region, &a.City, &a.EmployeeCategory, &a.DailyAmount, &a.Currency,
		&a.MealPercentage, &a.LodgingPercentage, &a.TransportPercentage, &a.OtherPercentage, &a.EffectiveFrom, &a.EffectiveTo, &a.IsActive, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("AllowanceRepo.Get", err)
	}
	return &a, nil
}

func (r *AllowanceRepo) ListRules(ctx context.Context, companyID uuid.UUID) ([]domain.DailyAllowanceRule, error) {
	q := `SELECT id,company_id,name,country,region,city,employee_category,daily_amount,currency,
		meal_percentage,lodging_percentage,transport_percentage,other_percentage,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM daily_allowance_rules WHERE company_id=$1 ORDER BY country,city`
	args := []any{companyID}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("AllowanceRepo.List", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.DailyAllowanceRule, error) {
		var a domain.DailyAllowanceRule
		err := row.Scan(&a.ID, &a.CompanyID, &a.Name, &a.Country, &a.Region, &a.City, &a.EmployeeCategory, &a.DailyAmount, &a.Currency,
			&a.MealPercentage, &a.LodgingPercentage, &a.TransportPercentage, &a.OtherPercentage, &a.EffectiveFrom, &a.EffectiveTo, &a.IsActive, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *AllowanceRepo) UpdateRule(ctx context.Context, a *domain.DailyAllowanceRule) error {
	q := `UPDATE daily_allowance_rules SET name=$1,country=$2,region=$3,city=$4,employee_category=$5,
		daily_amount=$6,currency=$7,meal_percentage=$8,lodging_percentage=$9,transport_percentage=$10,
		other_percentage=$11,effective_from=$12,effective_to=$13,is_active=$14,updated_at=NOW()
		WHERE id=$15 AND company_id=$16`
	_, err := r.pool.Exec(ctx, q, a.Name, a.Country, a.Region, a.City, a.EmployeeCategory,
		a.DailyAmount, a.Currency, a.MealPercentage, a.LodgingPercentage, a.TransportPercentage,
		a.OtherPercentage, a.EffectiveFrom, a.EffectiveTo, a.IsActive, a.ID, a.CompanyID)
	return repoErr("AllowanceRepo.Update", err)
}
