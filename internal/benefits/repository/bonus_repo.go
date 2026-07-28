package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type BonusRepo struct {
	pool *pgxpool.Pool
}

func NewBonusRepo(pool *pgxpool.Pool) *BonusRepo {
	return &BonusRepo{pool: pool}
}

func scanBonus(row pgx.CollectableRow) (domain.EmployeeBonus, error) {
	var b domain.EmployeeBonus
	err := row.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BonusType, &b.Name, &b.Description,
		&b.Amount, &b.Currency, &b.PaymentType, &b.InstallmentCount, &b.InstallmentAmount,
		&b.Frequency, &b.GrantDate, &b.VestingStart, &b.VestingEnd, &b.PaymentDate,
		&b.Status, &b.ClawbackAmount, &b.ClawbackReason, &b.PerformancePeriodStart, &b.PerformancePeriodEnd,
		&b.PerformanceScore, &b.IsTaxable, &b.TaxWithholding, &b.NetAmount, &b.ApprovedBy, &b.ApprovedAt,
		&b.PaidInPayroll, &b.PayrollRunID, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *BonusRepo) CreateBonus(ctx context.Context, b *domain.EmployeeBonus) error {
	q := `INSERT INTO employee_bonuses (id,company_id,employee_id,bonus_type,name,description,amount,currency,
		payment_type,installment_count,installment_amount,frequency,grant_date,vesting_start,vesting_end,
		payment_date,status,clawback_amount,clawback_reason,performance_period_start,performance_period_end,
		performance_score,is_taxable,tax_withholding,net_amount,approved_by,approved_at,paid_in_payroll,
		payroll_run_id,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.EmployeeID, b.BonusType, b.Name, b.Description,
		b.Amount, b.Currency, b.PaymentType, b.InstallmentCount, b.InstallmentAmount,
		b.Frequency, b.GrantDate, b.VestingStart, b.VestingEnd, b.PaymentDate,
		b.Status, b.ClawbackAmount, b.ClawbackReason, b.PerformancePeriodStart, b.PerformancePeriodEnd,
		b.PerformanceScore, b.IsTaxable, b.TaxWithholding, b.NetAmount, b.ApprovedBy, b.ApprovedAt,
		b.PaidInPayroll, b.PayrollRunID, b.Notes, b.CreatedBy)
	return repoErr("CreateBonus", err)
}

func (r *BonusRepo) GetBonus(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeBonus, error) {
	q := `SELECT id,company_id,employee_id,bonus_type,name,description,amount,currency,payment_type,
		installment_count,installment_amount,frequency,grant_date,vesting_start,vesting_end,payment_date,
		status,clawback_amount,clawback_reason,performance_period_start,performance_period_end,
		performance_score,is_taxable,tax_withholding,net_amount,approved_by,approved_at,paid_in_payroll,
		payroll_run_id,notes,created_by,created_at,updated_at
		FROM employee_bonuses WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	b, err := scanBonus(row)
	if err != nil {
		return nil, repoErr("GetBonus", err)
	}
	return &b, nil
}

func (r *BonusRepo) ListBonuses(ctx context.Context, employeeID uuid.UUID, status *string) ([]domain.EmployeeBonus, error) {
	q := `SELECT id,company_id,employee_id,bonus_type,name,description,amount,currency,payment_type,
		installment_count,installment_amount,frequency,grant_date,vesting_start,vesting_end,payment_date,
		status,clawback_amount,clawback_reason,performance_period_start,performance_period_end,
		performance_score,is_taxable,tax_withholding,net_amount,approved_by,approved_at,paid_in_payroll,
		payroll_run_id,notes,created_by,created_at,updated_at
		FROM employee_bonuses WHERE employee_id=$1`
	args := []any{employeeID}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", 2)
		args = append(args, *status)
	}
	q += " ORDER BY grant_date DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBonuses", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanBonus)
}

func (r *BonusRepo) UpdateBonus(ctx context.Context, b *domain.EmployeeBonus) error {
	q := `UPDATE employee_bonuses SET bonus_type=$1,name=$2,description=$3,amount=$4,currency=$5,
		payment_type=$6,installment_count=$7,installment_amount=$8,frequency=$9,grant_date=$10,
		vesting_start=$11,vesting_end=$12,payment_date=$13,status=$14,clawback_amount=$15,
		clawback_reason=$16,performance_period_start=$17,performance_period_end=$18,performance_score=$19,
		is_taxable=$20,tax_withholding=$21,net_amount=$22,approved_by=$23,approved_at=$24,
		paid_in_payroll=$25,payroll_run_id=$26,notes=$27,updated_at=NOW() WHERE id=$28 AND company_id=$29`
	_, err := r.pool.Exec(ctx, q, b.BonusType, b.Name, b.Description, b.Amount, b.Currency,
		b.PaymentType, b.InstallmentCount, b.InstallmentAmount, b.Frequency, b.GrantDate,
		b.VestingStart, b.VestingEnd, b.PaymentDate, b.Status, b.ClawbackAmount,
		b.ClawbackReason, b.PerformancePeriodStart, b.PerformancePeriodEnd, b.PerformanceScore,
		b.IsTaxable, b.TaxWithholding, b.NetAmount, b.ApprovedBy, b.ApprovedAt,
		b.PaidInPayroll, b.PayrollRunID, b.Notes, b.ID, b.CompanyID)
	return repoErr("UpdateBonus", err)
}

func (r *BonusRepo) UpdateBonusStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE employee_bonuses SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdateBonusStatus", err)
}

func scanIncentive(row pgx.CollectableRow) (domain.EmployeeIncentive, error) {
	var i domain.EmployeeIncentive
	err := row.Scan(&i.ID, &i.CompanyID, &i.EmployeeID, &i.IncentiveType, &i.Name, &i.Description,
		&i.Value, &i.Currency, &i.AwardDate, &i.ExpiryDate, &i.RedemptionDate, &i.Status,
		&i.PointsCost, &i.IsTaxable, &i.AwardedBy, &i.Notes, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

func (r *BonusRepo) CreateIncentive(ctx context.Context, i *domain.EmployeeIncentive) error {
	q := `INSERT INTO employee_incentives (id,company_id,employee_id,incentive_type,name,description,value,
		currency,award_date,expiry_date,redemption_date,status,points_cost,is_taxable,awarded_by,notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.pool.Exec(ctx, q, i.ID, i.CompanyID, i.EmployeeID, i.IncentiveType, i.Name, i.Description,
		i.Value, i.Currency, i.AwardDate, i.ExpiryDate, i.RedemptionDate, i.Status,
		i.PointsCost, i.IsTaxable, i.AwardedBy, i.Notes)
	return repoErr("CreateIncentive", err)
}

func (r *BonusRepo) GetIncentive(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeIncentive, error) {
	q := `SELECT id,company_id,employee_id,incentive_type,name,description,value,currency,award_date,
		expiry_date,redemption_date,status,points_cost,is_taxable,awarded_by,notes,created_at,updated_at
		FROM employee_incentives WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	i, err := scanIncentive(row)
	if err != nil {
		return nil, repoErr("GetIncentive", err)
	}
	return &i, nil
}

func (r *BonusRepo) ListIncentives(ctx context.Context, employeeID uuid.UUID, status *string) ([]domain.EmployeeIncentive, error) {
	q := `SELECT id,company_id,employee_id,incentive_type,name,description,value,currency,award_date,
		expiry_date,redemption_date,status,points_cost,is_taxable,awarded_by,notes,created_at,updated_at
		FROM employee_incentives WHERE employee_id=$1`
	args := []any{employeeID}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", 2)
		args = append(args, *status)
	}
	q += " ORDER BY award_date DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListIncentives", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanIncentive)
}

func (r *BonusRepo) UpdateIncentive(ctx context.Context, i *domain.EmployeeIncentive) error {
	q := `UPDATE employee_incentives SET incentive_type=$1,name=$2,description=$3,value=$4,currency=$5,
		award_date=$6,expiry_date=$7,redemption_date=$8,status=$9,points_cost=$10,is_taxable=$11,
		awarded_by=$12,notes=$13,updated_at=NOW() WHERE id=$14 AND company_id=$15`
	_, err := r.pool.Exec(ctx, q, i.IncentiveType, i.Name, i.Description, i.Value, i.Currency,
		i.AwardDate, i.ExpiryDate, i.RedemptionDate, i.Status, i.PointsCost, i.IsTaxable,
		i.AwardedBy, i.Notes, i.ID, i.CompanyID)
	return repoErr("UpdateIncentive", err)
}

func (r *BonusRepo) UpdateIncentiveStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE employee_incentives SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdateIncentiveStatus", err)
}

func scanPayrollMapping(row pgx.CollectableRow) (domain.BenefitPayrollMapping, error) {
	var pm domain.BenefitPayrollMapping
	err := row.Scan(&pm.ID, &pm.CompanyID, &pm.BenefitID, &pm.FlexiblePlanID, &pm.EmployeeBenefitID,
		&pm.BonusID, &pm.MappingType, &pm.PayrollConceptID, &pm.Amount, &pm.Currency, &pm.Frequency,
		&pm.EffectiveFrom, &pm.EffectiveTo, &pm.IsActive, &pm.LastSyncedAt, &pm.SyncStatus, &pm.SyncError,
		&pm.CreatedAt, &pm.UpdatedAt)
	return pm, err
}

func (r *BonusRepo) CreatePayrollMapping(ctx context.Context, pm *domain.BenefitPayrollMapping) error {
	q := `INSERT INTO benefit_payroll_mappings (id,company_id,benefit_id,flexible_plan_id,employee_benefit_id,
		bonus_id,mapping_type,payroll_concept_id,amount,currency,frequency,effective_from,effective_to,
		is_active,sync_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, pm.ID, pm.CompanyID, pm.BenefitID, pm.FlexiblePlanID, pm.EmployeeBenefitID,
		pm.BonusID, pm.MappingType, pm.PayrollConceptID, pm.Amount, pm.Currency, pm.Frequency,
		pm.EffectiveFrom, pm.EffectiveTo, pm.IsActive, pm.SyncStatus)
	return repoErr("CreatePayrollMapping", err)
}

func (r *BonusRepo) GetPayrollMapping(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitPayrollMapping, error) {
	q := `SELECT id,company_id,benefit_id,flexible_plan_id,employee_benefit_id,bonus_id,mapping_type,
		payroll_concept_id,amount,currency,frequency,effective_from,effective_to,is_active,last_synced_at,
		sync_status,sync_error,created_at,updated_at
		FROM benefit_payroll_mappings WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	pm, err := scanPayrollMapping(row)
	if err != nil {
		return nil, repoErr("GetPayrollMapping", err)
	}
	return &pm, nil
}

func (r *BonusRepo) ListPayrollMappings(ctx context.Context, benefitID uuid.UUID, mappingType *string) ([]domain.BenefitPayrollMapping, error) {
	q := `SELECT id,company_id,benefit_id,flexible_plan_id,employee_benefit_id,bonus_id,mapping_type,
		payroll_concept_id,amount,currency,frequency,effective_from,effective_to,is_active,last_synced_at,
		sync_status,sync_error,created_at,updated_at
		FROM benefit_payroll_mappings WHERE benefit_id=$1`
	args := []any{benefitID}
	if mappingType != nil {
		q += fmt.Sprintf(" AND mapping_type=$%d", 2)
		args = append(args, *mappingType)
	}
	q += " ORDER BY effective_from DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListPayrollMappings", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanPayrollMapping)
}

func (r *BonusRepo) UpdatePayrollMapping(ctx context.Context, pm *domain.BenefitPayrollMapping) error {
	q := `UPDATE benefit_payroll_mappings SET mapping_type=$1,payroll_concept_id=$2,amount=$3,currency=$4,
		frequency=$5,effective_from=$6,effective_to=$7,is_active=$8,updated_at=NOW()
		WHERE id=$9 AND company_id=$10`
	_, err := r.pool.Exec(ctx, q, pm.MappingType, pm.PayrollConceptID, pm.Amount, pm.Currency, pm.Frequency,
		pm.EffectiveFrom, pm.EffectiveTo, pm.IsActive, pm.ID, pm.CompanyID)
	return repoErr("UpdatePayrollMapping", err)
}

func (r *BonusRepo) UpdateSyncStatus(ctx context.Context, id uuid.UUID, syncStatus string, syncError *string) error {
	q := `UPDATE benefit_payroll_mappings SET sync_status=$1,sync_error=$2,last_synced_at=NOW(),updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, syncStatus, syncError, id)
	return repoErr("UpdateSyncStatus", err)
}
