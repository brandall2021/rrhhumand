package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type PlanRepo struct {
	pool *pgxpool.Pool
}

func NewPlanRepo(pool *pgxpool.Pool) *PlanRepo {
	return &PlanRepo{pool: pool}
}

func (r *PlanRepo) CreatePlan(ctx context.Context, p *domain.BenefitPlan) error {
	q := `INSERT INTO benefit_plans (id,company_id,benefit_id,name,description,plan_type,monthly_cost_employee,
		monthly_cost_employer,annual_cost_employee,annual_cost_employer,currency,coverage_limit,coverage_details,
		max_dependents,dependent_type,enrollment_fee,waiting_period_days,minimum_age,maximum_age,is_default,
		is_active,sort_order,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.BenefitID, p.Name, p.Description, p.PlanType,
		p.MonthlyCostEmployee, p.MonthlyCostEmployer, p.AnnualCostEmployee, p.AnnualCostEmployer,
		p.Currency, p.CoverageLimit, p.CoverageDetails, p.MaxDependents, p.DependentType,
		p.EnrollmentFee, p.WaitingPeriodDays, p.MinimumAge, p.MaximumAge, p.IsDefault,
		p.IsActive, p.SortOrder, p.CreatedBy)
	return repoErr("CreatePlan", err)
}

func (r *PlanRepo) GetPlan(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitPlan, error) {
	q := `SELECT id,company_id,benefit_id,name,description,plan_type,monthly_cost_employee,monthly_cost_employer,
		annual_cost_employee,annual_cost_employer,currency,coverage_limit,coverage_details,max_dependents,
		dependent_type,enrollment_fee,waiting_period_days,minimum_age,maximum_age,is_default,is_active,
		sort_order,created_by,created_at,updated_at
		FROM benefit_plans WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p domain.BenefitPlan
	err := row.Scan(&p.ID, &p.CompanyID, &p.BenefitID, &p.Name, &p.Description, &p.PlanType,
		&p.MonthlyCostEmployee, &p.MonthlyCostEmployer, &p.AnnualCostEmployee, &p.AnnualCostEmployer,
		&p.Currency, &p.CoverageLimit, &p.CoverageDetails, &p.MaxDependents, &p.DependentType,
		&p.EnrollmentFee, &p.WaitingPeriodDays, &p.MinimumAge, &p.MaximumAge, &p.IsDefault, &p.IsActive,
		&p.SortOrder, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetPlan", err)
	}
	return &p, nil
}

func (r *PlanRepo) ListPlans(ctx context.Context, benefitID uuid.UUID) ([]domain.BenefitPlan, error) {
	q := `SELECT id,company_id,benefit_id,name,description,plan_type,monthly_cost_employee,monthly_cost_employer,
		annual_cost_employee,annual_cost_employer,currency,coverage_limit,coverage_details,max_dependents,
		dependent_type,enrollment_fee,waiting_period_days,minimum_age,maximum_age,is_default,is_active,
		sort_order,created_by,created_at,updated_at
		FROM benefit_plans WHERE benefit_id=$1 ORDER BY sort_order,name`
	rows, err := r.pool.Query(ctx, q, benefitID)
	if err != nil {
		return nil, repoErr("ListPlans", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitPlan, error) {
		var p domain.BenefitPlan
		err := row.Scan(&p.ID, &p.CompanyID, &p.BenefitID, &p.Name, &p.Description, &p.PlanType,
			&p.MonthlyCostEmployee, &p.MonthlyCostEmployer, &p.AnnualCostEmployee, &p.AnnualCostEmployer,
			&p.Currency, &p.CoverageLimit, &p.CoverageDetails, &p.MaxDependents, &p.DependentType,
			&p.EnrollmentFee, &p.WaitingPeriodDays, &p.MinimumAge, &p.MaximumAge, &p.IsDefault, &p.IsActive,
			&p.SortOrder, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *PlanRepo) UpdatePlan(ctx context.Context, p *domain.BenefitPlan) error {
	q := `UPDATE benefit_plans SET name=$1,description=$2,plan_type=$3,monthly_cost_employee=$4,
		monthly_cost_employer=$5,annual_cost_employee=$6,annual_cost_employer=$7,currency=$8,
		coverage_limit=$9,coverage_details=$10,max_dependents=$11,dependent_type=$12,enrollment_fee=$13,
		waiting_period_days=$14,minimum_age=$15,maximum_age=$16,is_default=$17,is_active=$18,
		sort_order=$19,updated_at=NOW() WHERE id=$20 AND company_id=$21`
	_, err := r.pool.Exec(ctx, q, p.Name, p.Description, p.PlanType,
		p.MonthlyCostEmployee, p.MonthlyCostEmployer, p.AnnualCostEmployee, p.AnnualCostEmployer,
		p.Currency, p.CoverageLimit, p.CoverageDetails, p.MaxDependents, p.DependentType,
		p.EnrollmentFee, p.WaitingPeriodDays, p.MinimumAge, p.MaximumAge, p.IsDefault, p.IsActive,
		p.SortOrder, p.ID, p.CompanyID)
	return repoErr("UpdatePlan", err)
}

func (r *PlanRepo) DeletePlan(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_plans WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeletePlan", err)
}

func (r *PlanRepo) scanPlan(row pgx.CollectableRow) (domain.BenefitPlan, error) {
	var p domain.BenefitPlan
	err := row.Scan(&p.ID, &p.CompanyID, &p.BenefitID, &p.Name, &p.Description, &p.PlanType,
		&p.MonthlyCostEmployee, &p.MonthlyCostEmployer, &p.AnnualCostEmployee, &p.AnnualCostEmployer,
		&p.Currency, &p.CoverageLimit, &p.CoverageDetails, &p.MaxDependents, &p.DependentType,
		&p.EnrollmentFee, &p.WaitingPeriodDays, &p.MinimumAge, &p.MaximumAge, &p.IsDefault, &p.IsActive,
		&p.SortOrder, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
