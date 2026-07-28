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

type FlexibleRepo struct {
	pool *pgxpool.Pool
}

func NewFlexibleRepo(pool *pgxpool.Pool) *FlexibleRepo {
	return &FlexibleRepo{pool: pool}
}

func (r *FlexibleRepo) CreatePlan(ctx context.Context, p *domain.BenefitFlexiblePlan) error {
	q := `INSERT INTO benefit_flexible_plans (id,company_id,name,description,plan_type,annual_amount,monthly_amount,
		currency,employer_contribution,employee_contribution,contribution_frequency,max_rollover_amount,
		rollover_expiry_months,allow_reimbursement,allow_prepaid_card,require_receipts,receipt_min_amount,
		eligible_categories,tax_exempt,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Name, p.Description, p.PlanType,
		p.AnnualAmount, p.MonthlyAmount, p.Currency, p.EmployerContribution, p.EmployeeContribution,
		p.ContributionFrequency, p.MaxRolloverAmount, p.RolloverExpiryMonths,
		p.AllowReimbursement, p.AllowPrepaidCard, p.RequireReceipts, p.ReceiptMinAmount,
		p.EligibleCategories, p.TaxExempt, p.IsActive, p.CreatedBy)
	return repoErr("CreatePlan", err)
}

func (r *FlexibleRepo) GetPlan(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitFlexiblePlan, error) {
	q := `SELECT id,company_id,name,description,plan_type,annual_amount,monthly_amount,currency,
		employer_contribution,employee_contribution,contribution_frequency,max_rollover_amount,
		rollover_expiry_months,allow_reimbursement,allow_prepaid_card,require_receipts,receipt_min_amount,
		eligible_categories,tax_exempt,is_active,created_by,created_at,updated_at
		FROM benefit_flexible_plans WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p domain.BenefitFlexiblePlan
	err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.PlanType,
		&p.AnnualAmount, &p.MonthlyAmount, &p.Currency, &p.EmployerContribution, &p.EmployeeContribution,
		&p.ContributionFrequency, &p.MaxRolloverAmount, &p.RolloverExpiryMonths,
		&p.AllowReimbursement, &p.AllowPrepaidCard, &p.RequireReceipts, &p.ReceiptMinAmount,
		&p.EligibleCategories, &p.TaxExempt, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetPlan", err)
	}
	return &p, nil
}

func (r *FlexibleRepo) ListPlans(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitFlexiblePlan, error) {
	q := `SELECT id,company_id,name,description,plan_type,annual_amount,monthly_amount,currency,
		employer_contribution,employee_contribution,contribution_frequency,max_rollover_amount,
		rollover_expiry_months,allow_reimbursement,allow_prepaid_card,require_receipts,receipt_min_amount,
		eligible_categories,tax_exempt,is_active,created_by,created_at,updated_at
		FROM benefit_flexible_plans WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListPlans", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitFlexiblePlan, error) {
		var p domain.BenefitFlexiblePlan
		err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.PlanType,
			&p.AnnualAmount, &p.MonthlyAmount, &p.Currency, &p.EmployerContribution, &p.EmployeeContribution,
			&p.ContributionFrequency, &p.MaxRolloverAmount, &p.RolloverExpiryMonths,
			&p.AllowReimbursement, &p.AllowPrepaidCard, &p.RequireReceipts, &p.ReceiptMinAmount,
			&p.EligibleCategories, &p.TaxExempt, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *FlexibleRepo) UpdatePlan(ctx context.Context, p *domain.BenefitFlexiblePlan) error {
	q := `UPDATE benefit_flexible_plans SET name=$1,description=$2,plan_type=$3,annual_amount=$4,monthly_amount=$5,
		currency=$6,employer_contribution=$7,employee_contribution=$8,contribution_frequency=$9,
		max_rollover_amount=$10,rollover_expiry_months=$11,allow_reimbursement=$12,allow_prepaid_card=$13,
		require_receipts=$14,receipt_min_amount=$15,eligible_categories=$16,tax_exempt=$17,is_active=$18,
		updated_at=NOW() WHERE id=$19 AND company_id=$20`
	_, err := r.pool.Exec(ctx, q, p.Name, p.Description, p.PlanType,
		p.AnnualAmount, p.MonthlyAmount, p.Currency, p.EmployerContribution, p.EmployeeContribution,
		p.ContributionFrequency, p.MaxRolloverAmount, p.RolloverExpiryMonths,
		p.AllowReimbursement, p.AllowPrepaidCard, p.RequireReceipts, p.ReceiptMinAmount,
		p.EligibleCategories, p.TaxExempt, p.IsActive, p.ID, p.CompanyID)
	return repoErr("UpdatePlan", err)
}

func (r *FlexibleRepo) CreateBudget(ctx context.Context, b *domain.BenefitFlexibleBudget) error {
	q := `INSERT INTO benefit_flexible_budgets (id,company_id,employee_id,flexible_plan_id,fiscal_year,total_amount,
		used_amount,pending_amount,rolled_over_from_previous,currency,status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.EmployeeID, b.FlexiblePlanID, b.FiscalYear,
		b.TotalAmount, b.UsedAmount, b.PendingAmount, b.RolledOverFromPrevious, b.Currency, b.Status)
	return repoErr("CreateBudget", err)
}

func (r *FlexibleRepo) GetBudget(ctx context.Context, employeeID, flexiblePlanID uuid.UUID, fiscalYear int) (*domain.BenefitFlexibleBudget, error) {
	q := `SELECT id,company_id,employee_id,flexible_plan_id,fiscal_year,total_amount,used_amount,pending_amount,
		rolled_over_from_previous,currency,status,created_at,updated_at
		FROM benefit_flexible_budgets WHERE employee_id=$1 AND flexible_plan_id=$2 AND fiscal_year=$3`
	row := r.pool.QueryRow(ctx, q, employeeID, flexiblePlanID, fiscalYear)
	var b domain.BenefitFlexibleBudget
	err := row.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.FlexiblePlanID, &b.FiscalYear,
		&b.TotalAmount, &b.UsedAmount, &b.PendingAmount, &b.RolledOverFromPrevious, &b.Currency, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBudget", err)
	}
	return &b, nil
}

func (r *FlexibleRepo) ListBudgets(ctx context.Context, employeeID, planID *uuid.UUID) ([]domain.BenefitFlexibleBudget, error) {
	q := `SELECT id,company_id,employee_id,flexible_plan_id,fiscal_year,total_amount,used_amount,pending_amount,
		rolled_over_from_previous,currency,status,created_at,updated_at
		FROM benefit_flexible_budgets WHERE 1=1`
	args := []any{}
	n := 1
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if planID != nil {
		q += fmt.Sprintf(" AND flexible_plan_id=$%d", n)
		args = append(args, *planID)
		n++
	}
	q += " ORDER BY fiscal_year DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBudgets", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitFlexibleBudget, error) {
		var b domain.BenefitFlexibleBudget
		err := row.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.FlexiblePlanID, &b.FiscalYear,
			&b.TotalAmount, &b.UsedAmount, &b.PendingAmount, &b.RolledOverFromPrevious, &b.Currency, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		return b, err
	})
}

func (r *FlexibleRepo) UpdateUsage(ctx context.Context, id uuid.UUID, usedAmount, pendingAmount decimal.Decimal) error {
	q := `UPDATE benefit_flexible_budgets SET used_amount=used_amount+$1,pending_amount=pending_amount-$1,updated_at=NOW() WHERE id=$2`
	_, err := r.pool.Exec(ctx, q, usedAmount, id)
	return repoErr("UpdateUsage", err)
}
