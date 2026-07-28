package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreateRunEmployee(ctx context.Context, re *domain.PayrollRunEmployee) error {
	q := `INSERT INTO payroll_run_employees (id,run_id,employee_id,employment_id,status,currency)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, re.ID, re.RunID, re.EmployeeID, re.EmploymentID, re.Status, re.Currency)
	return repoErr("CreateRunEmployee", err)
}

func (r *Repository) GetRunEmployee(ctx context.Context, runID, employeeID string) (*domain.PayrollRunEmployee, error) {
	q := `SELECT id,run_id,employee_id,employment_id,status,gross_remunerative,gross_non_remunerative,
		deductions_amount,employer_contributions,employer_cost,net_amount,currency,calculation_version,error_message,calculated_at,created_at
		FROM payroll_run_employees WHERE run_id=$1 AND employee_id=$2`
	row := r.pool.QueryRow(ctx, q, runID, employeeID)
	var re domain.PayrollRunEmployee
	err := row.Scan(&re.ID, &re.RunID, &re.EmployeeID, &re.EmploymentID, &re.Status,
		&re.GrossRemunerative, &re.GrossNonRemunerative, &re.DeductionsAmount, &re.EmployerContributions,
		&re.EmployerCost, &re.NetAmount, &re.Currency, &re.CalculationVersion, &re.ErrorMessage, &re.CalculatedAt, &re.CreatedAt)
	if err != nil {
		return nil, repoErr("GetRunEmployee", err)
	}
	return &re, nil
}

func (r *Repository) ListRunEmployees(ctx context.Context, runID string) ([]domain.PayrollRunEmployee, error) {
	q := `SELECT id,run_id,employee_id,employment_id,status,gross_remunerative,gross_non_remunerative,
		deductions_amount,employer_contributions,employer_cost,net_amount,currency,calculation_version,error_message,calculated_at,created_at
		FROM payroll_run_employees WHERE run_id=$1 ORDER BY employee_id`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListRunEmployees", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollRunEmployee, error) {
		var re domain.PayrollRunEmployee
		err := row.Scan(&re.ID, &re.RunID, &re.EmployeeID, &re.EmploymentID, &re.Status,
			&re.GrossRemunerative, &re.GrossNonRemunerative, &re.DeductionsAmount, &re.EmployerContributions,
			&re.EmployerCost, &re.NetAmount, &re.Currency, &re.CalculationVersion, &re.ErrorMessage, &re.CalculatedAt, &re.CreatedAt)
		return re, err
	})
}

func (r *Repository) UpdateRunEmployeeResult(ctx context.Context, id string, re *domain.PayrollRunEmployee) error {
	q := `UPDATE payroll_run_employees SET status=$1,gross_remunerative=$2,gross_non_remunerative=$3,
		deductions_amount=$4,employer_contributions=$5,employer_cost=$6,net_amount=$7,
		calculation_version=$8,error_message=$9,calculated_at=$10 WHERE id=$11`
	_, err := r.pool.Exec(ctx, q, re.Status, re.GrossRemunerative, re.GrossNonRemunerative, re.DeductionsAmount,
		re.EmployerContributions, re.EmployerCost, re.NetAmount, re.CalculationVersion, re.ErrorMessage, re.CalculatedAt, id)
	return repoErr("UpdateRunEmployeeResult", err)
}

func (r *Repository) BulkUpdateRunEmployeeStatus(ctx context.Context, runID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_run_employees SET status=$1 WHERE run_id=$2`, status, runID)
	return repoErr("BulkUpdateRunEmployeeStatus", err)
}

func (r *Repository) CreateSnapshot(ctx context.Context, s *domain.EmployeeSnapshot) error {
	empData, _ := json.Marshal(s.EmployeeData)
	emplData, _ := json.Marshal(s.EmploymentData)
	posData, _ := json.Marshal(s.PositionData)
	catData, _ := json.Marshal(s.CategoryData)
	agrData, _ := json.Marshal(s.AgreementData)
	salData, _ := json.Marshal(s.SalaryData)
	benData, _ := json.Marshal(s.BenefitsData)
	taxData, _ := json.Marshal(s.TaxConfigData)
	ssData, _ := json.Marshal(s.SocialSecurityData)

	q := `INSERT INTO payroll_employee_snapshots (id,run_employee_id,employee_data,employment_data,
		position_data,category_data,agreement_data,salary_data,benefits_data,tax_config_data,social_security_data)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.RunEmployeeID,
		empData, emplData, posData, catData, agrData,
		salData, benData, taxData, ssData)
	return repoErr("CreateSnapshot", err)
}
