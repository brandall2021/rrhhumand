package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type CostRepo struct {
	pool *pgxpool.Pool
}

func NewCostRepo(pool *pgxpool.Pool) *CostRepo {
	return &CostRepo{pool: pool}
}

func (r *CostRepo) CreateCost(ctx context.Context, c *domain.BenefitCost) error {
	q := `INSERT INTO benefit_costs (id,company_id,benefit_id,plan_id,cost_type,employee_cost,employer_cost,
		total_cost,currency,frequency,billing_cycle_day,effective_from,effective_to,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.BenefitID, c.PlanID, c.CostType,
		c.EmployeeCost, c.EmployerCost, c.TotalCost, c.Currency, c.Frequency, c.BillingCycleDay,
		c.EffectiveFrom, c.EffectiveTo, c.IsActive, c.CreatedBy)
	return repoErr("CreateCost", err)
}

func (r *CostRepo) GetCost(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitCost, error) {
	q := `SELECT id,company_id,benefit_id,plan_id,cost_type,employee_cost,employer_cost,total_cost,currency,
		frequency,billing_cycle_day,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM benefit_costs WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.BenefitCost
	err := row.Scan(&c.ID, &c.CompanyID, &c.BenefitID, &c.PlanID, &c.CostType,
		&c.EmployeeCost, &c.EmployerCost, &c.TotalCost, &c.Currency, &c.Frequency, &c.BillingCycleDay,
		&c.EffectiveFrom, &c.EffectiveTo, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCost", err)
	}
	return &c, nil
}

func (r *CostRepo) ListCosts(ctx context.Context, benefitID uuid.UUID) ([]domain.BenefitCost, error) {
	q := `SELECT id,company_id,benefit_id,plan_id,cost_type,employee_cost,employer_cost,total_cost,currency,
		frequency,billing_cycle_day,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM benefit_costs WHERE benefit_id=$1 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, benefitID)
	if err != nil {
		return nil, repoErr("ListCosts", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitCost, error) {
		var c domain.BenefitCost
		err := row.Scan(&c.ID, &c.CompanyID, &c.BenefitID, &c.PlanID, &c.CostType,
			&c.EmployeeCost, &c.EmployerCost, &c.TotalCost, &c.Currency, &c.Frequency, &c.BillingCycleDay,
			&c.EffectiveFrom, &c.EffectiveTo, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *CostRepo) UpdateCost(ctx context.Context, c *domain.BenefitCost) error {
	q := `UPDATE benefit_costs SET plan_id=$1,cost_type=$2,employee_cost=$3,employer_cost=$4,total_cost=$5,
		currency=$6,frequency=$7,billing_cycle_day=$8,effective_from=$9,effective_to=$10,is_active=$11,
		updated_at=NOW() WHERE id=$12 AND company_id=$13`
	_, err := r.pool.Exec(ctx, q, c.PlanID, c.CostType, c.EmployeeCost, c.EmployerCost, c.TotalCost,
		c.Currency, c.Frequency, c.BillingCycleDay, c.EffectiveFrom, c.EffectiveTo, c.IsActive,
		c.ID, c.CompanyID)
	return repoErr("UpdateCost", err)
}

func (r *CostRepo) CreateSchedule(ctx context.Context, s *domain.BenefitCostSchedule) error {
	q := `INSERT INTO benefit_cost_schedules (id,company_id,benefit_id,cost_id,schedule_date,amount,currency,status,paid_at,payment_reference,notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.BenefitID, s.CostID, s.ScheduleDate, s.Amount, s.Currency, s.Status, s.PaidAt, s.PaymentReference, s.Notes)
	return repoErr("CreateSchedule", err)
}

func (r *CostRepo) ListSchedules(ctx context.Context, benefitID uuid.UUID, dateFrom, dateTo *time.Time) ([]domain.BenefitCostSchedule, error) {
	q := `SELECT id,company_id,benefit_id,cost_id,schedule_date,amount,currency,status,paid_at,payment_reference,notes,created_at,updated_at
		FROM benefit_cost_schedules WHERE benefit_id=$1`
	args := []any{benefitID}
	n := 2
	if dateFrom != nil {
		q += fmt.Sprintf(" AND schedule_date>=$%d", n)
		args = append(args, *dateFrom)
		n++
	}
	if dateTo != nil {
		q += fmt.Sprintf(" AND schedule_date<=$%d", n)
		args = append(args, *dateTo)
		n++
	}
	q += " ORDER BY schedule_date"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListSchedules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitCostSchedule, error) {
		var s domain.BenefitCostSchedule
		err := row.Scan(&s.ID, &s.CompanyID, &s.BenefitID, &s.CostID, &s.ScheduleDate, &s.Amount, &s.Currency, &s.Status, &s.PaidAt, &s.PaymentReference, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
		return s, err
	})
}

func (r *CostRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, paidAt *time.Time, paymentReference *string) error {
	q := `UPDATE benefit_cost_schedules SET status=$1,paid_at=$2,payment_reference=$3,updated_at=NOW() WHERE id=$4`
	_, err := r.pool.Exec(ctx, q, status, paidAt, paymentReference, id)
	return repoErr("UpdateStatus", err)
}
