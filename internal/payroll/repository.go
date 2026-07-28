package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func repoErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("payroll_repo.%s: %w", op, err)
}

type Repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) *Repository {
	return &Repository{pool: pool, log: log}
}

// ========================================================================
// PERIODS
// ========================================================================

func (r *Repository) CreatePeriod(ctx context.Context, p *PayrollPeriod) error {
	q := `INSERT INTO payroll_periods (id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Year, p.Month, p.PeriodType, p.Name, p.StartDate, p.EndDate, p.PaymentDate, p.Status, p.CreatedBy)
	return repoErr("CreatePeriod", err)
}

func (r *Repository) UpdatePeriod(ctx context.Context, p *PayrollPeriod) error {
	q := `UPDATE payroll_periods SET name=$1,payment_date=$2,updated_at=NOW() WHERE id=$3 AND company_id=$4`
	_, err := r.pool.Exec(ctx, q, p.Name, p.PaymentDate, p.ID, p.CompanyID)
	return repoErr("UpdatePeriod", err)
}

func (r *Repository) GetPeriod(ctx context.Context, companyID, id string) (*PayrollPeriod, error) {
	q := `SELECT id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,closed_at,created_by,created_at,updated_at
		FROM payroll_periods WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p PayrollPeriod
	err := row.Scan(&p.ID, &p.CompanyID, &p.Year, &p.Month, &p.PeriodType, &p.Name, &p.StartDate, &p.EndDate, &p.PaymentDate,
		&p.Status, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetPeriod", err)
	}
	return &p, nil
}

func (r *Repository) ListPeriods(ctx context.Context, companyID string, limit, offset int) ([]PayrollPeriod, error) {
	q := `SELECT id,company_id,year,month,period_type,name,start_date,end_date,payment_date,status,closed_at,created_by,created_at,updated_at
		FROM payroll_periods WHERE company_id=$1 ORDER BY year DESC, month DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, companyID, limit, offset)
	if err != nil {
		return nil, repoErr("ListPeriods", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollPeriod, error) {
		var p PayrollPeriod
		err := row.Scan(&p.ID, &p.CompanyID, &p.Year, &p.Month, &p.PeriodType, &p.Name, &p.StartDate, &p.EndDate, &p.PaymentDate,
			&p.Status, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *Repository) UpdatePeriodStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_periods SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdatePeriodStatus", err)
}

func (r *Repository) ClosePeriod(ctx context.Context, id, closedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_periods SET status='CLOSED',closed_at=NOW(),updated_at=NOW() WHERE id=$1 AND closed_by=$2`, id, closedBy)
	return repoErr("ClosePeriod", err)
}

// ========================================================================
// RUNS
// ========================================================================

func (r *Repository) CreateRun(ctx context.Context, run *PayrollRun) error {
	q := `INSERT INTO payroll_runs (id,company_id,period_id,run_number,run_type,status,engine_version,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, run.ID, run.CompanyID, run.PeriodID, run.RunNumber, run.RunType, run.Status, run.EngineVersion, run.CreatedBy)
	return repoErr("CreateRun", err)
}

func (r *Repository) GetRun(ctx context.Context, companyID, id string) (*PayrollRun, error) {
	q := `SELECT id,company_id,period_id,run_number,run_type,status,engine_version,started_at,finished_at,
		created_by,approved_by,approved_at,closed_by,closed_at,created_at,updated_at
		FROM payroll_runs WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var run PayrollRun
	err := row.Scan(&run.ID, &run.CompanyID, &run.PeriodID, &run.RunNumber, &run.RunType, &run.Status, &run.EngineVersion,
		&run.StartedAt, &run.FinishedAt, &run.CreatedBy, &run.ApprovedBy, &run.ApprovedAt, &run.ClosedBy, &run.ClosedAt, &run.CreatedAt, &run.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetRun", err)
	}
	return &run, nil
}

func (r *Repository) ListRuns(ctx context.Context, companyID string, filter RunFilter) ([]PayrollRun, error) {
	q := `SELECT id,company_id,period_id,run_number,run_type,status,engine_version,started_at,finished_at,
		created_by,approved_by,approved_at,closed_by,closed_at,created_at,updated_at
		FROM payroll_runs WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if filter.PeriodID != nil {
		q += fmt.Sprintf(" AND period_id=$%d", n)
		args = append(args, *filter.PeriodID)
		n++
	}
	if filter.RunType != nil {
		q += fmt.Sprintf(" AND run_type=$%d", n)
		args = append(args, *filter.RunType)
		n++
	}
	if filter.Status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *filter.Status)
		n++
	}
	q += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListRuns", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollRun, error) {
		var run PayrollRun
		err := row.Scan(&run.ID, &run.CompanyID, &run.PeriodID, &run.RunNumber, &run.RunType, &run.Status, &run.EngineVersion,
			&run.StartedAt, &run.FinishedAt, &run.CreatedBy, &run.ApprovedBy, &run.ApprovedAt, &run.ClosedBy, &run.ClosedAt, &run.CreatedAt, &run.UpdatedAt)
		return run, err
	})
}

func (r *Repository) UpdateRunStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdateRunStatus", err)
}

func (r *Repository) UpdateRunTimestamps(ctx context.Context, id, status string, startedAt, finishedAt *time.Time) error {
	q := `UPDATE payroll_runs SET status=$1,started_at=$2,finished_at=$3,updated_at=NOW() WHERE id=$4`
	_, err := r.pool.Exec(ctx, q, status, startedAt, finishedAt, id)
	return repoErr("UpdateRunTimestamps", err)
}

func (r *Repository) ApproveRun(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status='APPROVED',approved_by=$1,approved_at=NOW(),updated_at=NOW() WHERE id=$2`, approvedBy, id)
	return repoErr("ApproveRun", err)
}

func (r *Repository) CloseRun(ctx context.Context, id, closedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status='CLOSED',closed_by=$1,closed_at=NOW(),updated_at=NOW() WHERE id=$2`, closedBy, id)
	return repoErr("CloseRun", err)
}

func (r *Repository) GetRunNumber(ctx context.Context, periodID string) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(MAX(run_number),0)+1 FROM payroll_runs WHERE period_id=$1`, periodID).Scan(&n)
	return n, repoErr("GetRunNumber", err)
}

// ========================================================================
// RUN EMPLOYEES
// ========================================================================

func (r *Repository) AddRunEmployee(ctx context.Context, re *PayrollRunEmployee) error {
	q := `INSERT INTO payroll_run_employees (id,run_id,employee_id,employment_id,status,currency)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, re.ID, re.RunID, re.EmployeeID, re.EmploymentID, re.Status, re.Currency)
	return repoErr("AddRunEmployee", err)
}

func (r *Repository) GetRunEmployee(ctx context.Context, runID, employeeID string) (*PayrollRunEmployee, error) {
	q := `SELECT id,run_id,employee_id,employment_id,status,gross_remunerative,gross_non_remunerative,
		deductions_amount,employer_contributions,employer_cost,net_amount,currency,calculation_version,error_message,calculated_at,created_at
		FROM payroll_run_employees WHERE run_id=$1 AND employee_id=$2`
	row := r.pool.QueryRow(ctx, q, runID, employeeID)
	var re PayrollRunEmployee
	err := row.Scan(&re.ID, &re.RunID, &re.EmployeeID, &re.EmploymentID, &re.Status,
		&re.GrossRemunerative, &re.GrossNonRemunerative, &re.DeductionsAmount, &re.EmployerContributions,
		&re.EmployerCost, &re.NetAmount, &re.Currency, &re.CalculationVersion, &re.ErrorMessage, &re.CalculatedAt, &re.CreatedAt)
	if err != nil {
		return nil, repoErr("GetRunEmployee", err)
	}
	return &re, nil
}

func (r *Repository) ListRunEmployees(ctx context.Context, runID string) ([]PayrollRunEmployee, error) {
	q := `SELECT id,run_id,employee_id,employment_id,status,gross_remunerative,gross_non_remunerative,
		deductions_amount,employer_contributions,employer_cost,net_amount,currency,calculation_version,error_message,calculated_at,created_at
		FROM payroll_run_employees WHERE run_id=$1 ORDER BY employee_id`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListRunEmployees", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollRunEmployee, error) {
		var re PayrollRunEmployee
		err := row.Scan(&re.ID, &re.RunID, &re.EmployeeID, &re.EmploymentID, &re.Status,
			&re.GrossRemunerative, &re.GrossNonRemunerative, &re.DeductionsAmount, &re.EmployerContributions,
			&re.EmployerCost, &re.NetAmount, &re.Currency, &re.CalculationVersion, &re.ErrorMessage, &re.CalculatedAt, &re.CreatedAt)
		return re, err
	})
}

func (r *Repository) UpdateRunEmployeeResult(ctx context.Context, id string, re *PayrollRunEmployee) error {
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

// ========================================================================
// SNAPSHOTS
// ========================================================================

func (r *Repository) CreateSnapshot(ctx context.Context, s *EmployeeSnapshot) error {
	q := `INSERT INTO payroll_employee_snapshots (id,run_employee_id,employee_data,employment_data,
		position_data,category_data,agreement_data,salary_data,benefits_data,tax_config_data,social_security_data)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.RunEmployeeID, s.EmployeeData, s.EmploymentData,
		s.PositionData, s.CategoryData, s.AgreementData, s.SalaryData, s.BenefitsData, s.TaxConfigData, s.SocialSecurityData)
	return repoErr("CreateSnapshot", err)
}

// ========================================================================
// CONCEPTS
// ========================================================================

func (r *Repository) CreateConcept(ctx context.Context, c *PayrollConcept) error {
	q := `INSERT INTO payroll_concepts (id,company_id,code,name,description,concept_type,taxability,
		calculation_type,base_concept_id,active,effective_from,effective_to,sort_order,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Code, c.Name, c.Description, c.ConceptType, c.Taxability,
		c.CalculationType, c.BaseConceptID, c.Active, c.EffectiveFrom, c.EffectiveTo, c.SortOrder, c.CreatedBy)
	return repoErr("CreateConcept", err)
}

func (r *Repository) UpdateConcept(ctx context.Context, c *PayrollConcept) error {
	q := `UPDATE payroll_concepts SET name=$1,description=$2,concept_type=$3,taxability=$4,
		calculation_type=$5,base_concept_id=$6,active=$7,sort_order=$8,updated_at=NOW() WHERE id=$9 AND company_id=$10`
	_, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.ConceptType, c.Taxability,
		c.CalculationType, c.BaseConceptID, c.Active, c.SortOrder, c.ID, c.CompanyID)
	return repoErr("UpdateConcept", err)
}

func (r *Repository) GetConcept(ctx context.Context, companyID, id string) (*PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c PayrollConcept
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
		&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetConcept", err)
	}
	return &c, nil
}

func (r *Repository) GetConceptByCode(ctx context.Context, companyID, code string) (*PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE company_id=$1 AND code=$2`
	row := r.pool.QueryRow(ctx, q, companyID, code)
	var c PayrollConcept
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
		&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetConceptByCode", err)
	}
	return &c, nil
}

func (r *Repository) ListConcepts(ctx context.Context, companyID string, filter ConceptFilter) ([]PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if filter.Active != nil {
		q += fmt.Sprintf(" AND active=$%d", n)
		args = append(args, *filter.Active)
		n++
	}
	if filter.ConceptType != nil {
		q += fmt.Sprintf(" AND concept_type=$%d", n)
		args = append(args, *filter.ConceptType)
		n++
	}
	if filter.Taxability != nil {
		q += fmt.Sprintf(" AND taxability=$%d", n)
		args = append(args, *filter.Taxability)
		n++
	}
	q += " ORDER BY sort_order, code"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListConcepts", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollConcept, error) {
		var c PayrollConcept
		err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
			&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

// ========================================================================
// RULES
// ========================================================================

func (r *Repository) CreateRule(ctx context.Context, rule *PayrollRule) error {
	params, _ := json.Marshal(rule.Parameters)
	q := `INSERT INTO payroll_rules (id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,
		rule_type,formula,parameters,priority,effective_from,effective_to,version,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, rule.ID, rule.CompanyID, rule.Country, rule.Jurisdiction, rule.AgreementID, rule.CategoryID,
		rule.ConceptID, rule.RuleType, rule.Formula, params, rule.Priority, rule.EffectiveFrom, rule.EffectiveTo, rule.Version, rule.CreatedBy)
	return repoErr("CreateRule", err)
}

func (r *Repository) UpdateRule(ctx context.Context, rule *PayrollRule) error {
	params, _ := json.Marshal(rule.Parameters)
	q := `UPDATE payroll_rules SET rule_type=$1,formula=$2,parameters=$3,priority=$4,active=$5,effective_to=$6,version=version+1,updated_at=NOW()
		WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, q, rule.RuleType, rule.Formula, params, rule.Priority, rule.Active, rule.EffectiveTo, rule.ID, rule.CompanyID)
	return repoErr("UpdateRule", err)
}

func (r *Repository) GetRule(ctx context.Context, companyID, id string) (*PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at FROM payroll_rules WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var rule PayrollRule
	err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
		&rule.ConceptID, &rule.RuleType, &rule.Formula, &rule.Parameters, &rule.Priority, &rule.EffectiveFrom,
		&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetRule", err)
	}
	return &rule, nil
}

func (r *Repository) ListRules(ctx context.Context, companyID string) ([]PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at
		FROM payroll_rules WHERE company_id=$1 ORDER BY priority, created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollRule, error) {
		var rule PayrollRule
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
			&rule.ConceptID, &rule.RuleType, &rule.Formula, &rule.Parameters, &rule.Priority, &rule.EffectiveFrom,
			&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

func (r *Repository) GetActiveRulesForConcept(ctx context.Context, companyID, conceptID string, date time.Time) ([]PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at
		FROM payroll_rules WHERE company_id=$1 AND concept_id=$2 AND active=true AND effective_from<=$3 AND (effective_to IS NULL OR effective_to>=$3)
		ORDER BY priority DESC`
	rows, err := r.pool.Query(ctx, q, companyID, conceptID, date)
	if err != nil {
		return nil, repoErr("GetActiveRulesForConcept", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollRule, error) {
		var rule PayrollRule
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
			&rule.ConceptID, &rule.RuleType, &rule.Formula, &rule.Parameters, &rule.Priority, &rule.EffectiveFrom,
			&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

func (r *Repository) GetActiveRulesByConceptIDs(ctx context.Context, companyID string, conceptIDs []string, date time.Time) ([]PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at
		FROM payroll_rules WHERE company_id=$1 AND concept_id=ANY($2) AND active=true AND effective_from<=$3 AND (effective_to IS NULL OR effective_to>=$3)
		ORDER BY concept_id, priority DESC`
	rows, err := r.pool.Query(ctx, q, companyID, conceptIDs, date)
	if err != nil {
		return nil, repoErr("GetActiveRulesByConceptIDs", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollRule, error) {
		var rule PayrollRule
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
			&rule.ConceptID, &rule.RuleType, &rule.Formula, &rule.Parameters, &rule.Priority, &rule.EffectiveFrom,
			&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

// ========================================================================
// NOVELTIES
// ========================================================================

func (r *Repository) CreateNovelty(ctx context.Context, n *PayrollNovelty) error {
	q := `INSERT INTO payroll_novelties (id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,
		unit_value,multiplier,start_date,end_date,description,source,source_reference_id,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, n.ID, n.CompanyID, n.EmployeeID, n.PeriodID, n.NoveltyType,
		n.Quantity, n.Unit, n.Amount, n.UnitValue, n.Multiplier, n.StartDate, n.EndDate, n.Description,
		n.Source, n.SourceReferenceID, n.Status, n.CreatedBy)
	return repoErr("CreateNovelty", err)
}

func (r *Repository) UpdateNovelty(ctx context.Context, n *PayrollNovelty) error {
	q := `UPDATE payroll_novelties SET quantity=$1,amount=$2,description=$3,status=$4,updated_at=NOW() WHERE id=$5 AND company_id=$6`
	_, err := r.pool.Exec(ctx, q, n.Quantity, n.Amount, n.Description, n.Status, n.ID, n.CompanyID)
	return repoErr("UpdateNovelty", err)
}

func (r *Repository) GetNovelty(ctx context.Context, companyID, id string) (*PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var n PayrollNovelty
	err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
		&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
		&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetNovelty", err)
	}
	return &n, nil
}

func (r *Repository) ListNovelties(ctx context.Context, companyID string, filter NoveltyFilter) ([]PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if filter.EmployeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *filter.EmployeeID)
		n++
	}
	if filter.PeriodID != nil {
		q += fmt.Sprintf(" AND period_id=$%d", n)
		args = append(args, *filter.PeriodID)
		n++
	}
	if filter.NoveltyType != nil {
		q += fmt.Sprintf(" AND novelty_type=$%d", n)
		args = append(args, *filter.NoveltyType)
		n++
	}
	if filter.Status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *filter.Status)
		n++
	}
	if filter.Source != nil {
		q += fmt.Sprintf(" AND source=$%d", n)
		args = append(args, *filter.Source)
		n++
	}
	q += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, filter.Limit)
		n++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListNovelties", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollNovelty, error) {
		var n PayrollNovelty
		err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
			&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
			&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
		return n, err
	})
}

func (r *Repository) DeleteNovelty(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_novelties WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteNovelty", err)
}

func (r *Repository) ApproveNovelty(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_novelties SET status='APPROVED',approved_by=$1,approved_at=NOW() WHERE id=$2`, approvedBy, id)
	return repoErr("ApproveNovelty", err)
}

func (r *Repository) GetNoveltiesForEmployeePeriod(ctx context.Context, companyID, employeeID, periodID string) ([]PayrollNovelty, error) {
	q := `SELECT id,company_id,employee_id,period_id,novelty_type,quantity,unit,amount,unit_value,multiplier,
		start_date,end_date,description,source,source_reference_id,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM payroll_novelties WHERE company_id=$1 AND employee_id=$2 AND period_id=$3 AND status='APPROVED' ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID, periodID)
	if err != nil {
		return nil, repoErr("GetNoveltiesForEmployeePeriod", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollNovelty, error) {
		var n PayrollNovelty
		err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.PeriodID, &n.NoveltyType, &n.Quantity, &n.Unit, &n.Amount,
			&n.UnitValue, &n.Multiplier, &n.StartDate, &n.EndDate, &n.Description, &n.Source, &n.SourceReferenceID,
			&n.Status, &n.ApprovedBy, &n.ApprovedAt, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt)
		return n, err
	})
}

// ========================================================================
// ITEMS
// ========================================================================

func (r *Repository) BulkCreateItems(ctx context.Context, items []PayrollItem) error {
	if len(items) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_items (id,run_employee_id,concept_id,quantity,unit_value,base_amount,rate,amount,
		is_remunerative,is_deduction,is_employer_contribution,calculation_detail,sort_order) VALUES `
	args := []any{}
	n := 1
	for _, it := range items {
		detail, _ := json.Marshal(it.CalculationDetail)
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12)
		args = append(args, it.ID, it.RunEmployeeID, it.ConceptID, it.Quantity, it.UnitValue, it.BaseAmount, it.Rate, it.Amount,
			it.IsRemunerative, it.IsDeduction, it.IsEmployerContribution, detail, it.SortOrder)
		n += 13
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateItems", err)
}

func (r *Repository) ListItems(ctx context.Context, runEmployeeID string) ([]PayrollItem, error) {
	q := `SELECT id,run_employee_id,concept_id,quantity,unit_value,base_amount,rate,amount,
		is_remunerative,is_deduction,is_employer_contribution,calculation_detail,sort_order,created_at
		FROM payroll_items WHERE run_employee_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q, runEmployeeID)
	if err != nil {
		return nil, repoErr("ListItems", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollItem, error) {
		var it PayrollItem
		err := row.Scan(&it.ID, &it.RunEmployeeID, &it.ConceptID, &it.Quantity, &it.UnitValue, &it.BaseAmount, &it.Rate,
			&it.Amount, &it.IsRemunerative, &it.IsDeduction, &it.IsEmployerContribution, &it.CalculationDetail, &it.SortOrder, &it.CreatedAt)
		return it, err
	})
}

func (r *Repository) DeleteItemsForRunEmployee(ctx context.Context, runEmployeeID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_items WHERE run_employee_id=$1`, runEmployeeID)
	return repoErr("DeleteItemsForRunEmployee", err)
}

// ========================================================================
// BASES
// ========================================================================

func (r *Repository) BulkCreateBases(ctx context.Context, bases []PayrollBase) error {
	if len(bases) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_bases (id,run_employee_id,base_type,base_amount,concept_ids,calculation_detail) VALUES `
	args := []any{}
	n := 1
	for _, b := range bases {
		detail, _ := json.Marshal(b.CalculationDetail)
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5)
		args = append(args, b.ID, b.RunEmployeeID, b.BaseType, b.BaseAmount, b.ConceptIDs, detail)
		n += 6
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateBases", err)
}

func (r *Repository) ListBases(ctx context.Context, runEmployeeID string) ([]PayrollBase, error) {
	q := `SELECT id,run_employee_id,base_type,base_amount,concept_ids,calculation_detail,created_at
		FROM payroll_bases WHERE run_employee_id=$1 ORDER BY base_type`
	rows, err := r.pool.Query(ctx, q, runEmployeeID)
	if err != nil {
		return nil, repoErr("ListBases", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollBase, error) {
		var b PayrollBase
		err := row.Scan(&b.ID, &b.RunEmployeeID, &b.BaseType, &b.BaseAmount, &b.ConceptIDs, &b.CalculationDetail, &b.CreatedAt)
		return b, err
	})
}

// ========================================================================
// DEDUCTIONS
// ========================================================================

func (r *Repository) BulkCreateDeductions(ctx context.Context, deductions []PayrollDeduction) error {
	if len(deductions) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_deductions (id,run_employee_id,concept_id,base_amount,rate,amount,cap_amount,exceeded_amount) VALUES `
	args := []any{}
	n := 1
	for _, d := range deductions {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		args = append(args, d.ID, d.RunEmployeeID, d.ConceptID, d.BaseAmount, d.Rate, d.Amount, d.CapAmount, d.ExceededAmount)
		n += 8
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateDeductions", err)
}

// ========================================================================
// CONTRIBUTIONS
// ========================================================================

func (r *Repository) BulkCreateContributions(ctx context.Context, contributions []PayrollContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_contributions (id,run_employee_id,concept_id,base_amount,rate,amount,cap_amount,exceeded_amount) VALUES `
	args := []any{}
	n := 1
	for _, c := range contributions {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		args = append(args, c.ID, c.RunEmployeeID, c.ConceptID, c.BaseAmount, c.Rate, c.Amount, c.CapAmount, c.ExceededAmount)
		n += 8
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateContributions", err)
}

// ========================================================================
// AGREEMENTS
// ========================================================================

func (r *Repository) CreateAgreement(ctx context.Context, a *LaborAgreement) error {
	q := `INSERT INTO labor_agreements (id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.Code, a.Name, a.Description, a.Activity, a.EffectiveFrom, a.EffectiveTo, a.Status, a.CreatedBy)
	return repoErr("CreateAgreement", err)
}

func (r *Repository) GetAgreement(ctx context.Context, companyID, id string) (*LaborAgreement, error) {
	q := `SELECT id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by,created_at,updated_at
		FROM labor_agreements WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a LaborAgreement
	err := row.Scan(&a.ID, &a.CompanyID, &a.Code, &a.Name, &a.Description, &a.Activity, &a.EffectiveFrom, &a.EffectiveTo, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetAgreement", err)
	}
	return &a, nil
}

func (r *Repository) ListAgreements(ctx context.Context, companyID string) ([]LaborAgreement, error) {
	q := `SELECT id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by,created_at,updated_at
		FROM labor_agreements WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListAgreements", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (LaborAgreement, error) {
		var a LaborAgreement
		err := row.Scan(&a.ID, &a.CompanyID, &a.Code, &a.Name, &a.Description, &a.Activity, &a.EffectiveFrom, &a.EffectiveTo, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

// ========================================================================
// CATEGORIES
// ========================================================================

func (r *Repository) CreateCategory(ctx context.Context, c *LaborCategory) error {
	q := `INSERT INTO labor_categories (id,company_id,agreement_id,code,name,description,effective_from,effective_to)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.AgreementID, c.Code, c.Name, c.Description, c.EffectiveFrom, c.EffectiveTo)
	return repoErr("CreateCategory", err)
}

func (r *Repository) ListCategories(ctx context.Context, companyID string) ([]LaborCategory, error) {
	q := `SELECT id,company_id,agreement_id,code,name,description,effective_from,effective_to,created_at,updated_at
		FROM labor_categories WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCategories", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (LaborCategory, error) {
		var c LaborCategory
		err := row.Scan(&c.ID, &c.CompanyID, &c.AgreementID, &c.Code, &c.Name, &c.Description, &c.EffectiveFrom, &c.EffectiveTo, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

// ========================================================================
// SALARY SCALES
// ========================================================================

func (r *Repository) CreateSalaryScale(ctx context.Context, s *SalaryScale) error {
	q := `INSERT INTO salary_scales (id,company_id,agreement_id,category_id,minimum_salary,maximum_salary,reference_value,effective_from,effective_to)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.AgreementID, s.CategoryID, s.MinimumSalary, s.MaximumSalary, s.ReferenceValue, s.EffectiveFrom, s.EffectiveTo)
	return repoErr("CreateSalaryScale", err)
}

func (r *Repository) ListSalaryScales(ctx context.Context, companyID string) ([]SalaryScale, error) {
	q := `SELECT id,company_id,agreement_id,category_id,minimum_salary,maximum_salary,reference_value,effective_from,effective_to,created_at,updated_at
		FROM salary_scales WHERE company_id=$1 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListSalaryScales", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (SalaryScale, error) {
		var s SalaryScale
		err := row.Scan(&s.ID, &s.CompanyID, &s.AgreementID, &s.CategoryID, &s.MinimumSalary, &s.MaximumSalary, &s.ReferenceValue, &s.EffectiveFrom, &s.EffectiveTo, &s.CreatedAt, &s.UpdatedAt)
		return s, err
	})
}

// ========================================================================
// MINIMUM WAGES
// ========================================================================

func (r *Repository) GetMinimumWage(ctx context.Context, country string, date time.Time) (*StatutoryMinimumWage, error) {
	q := `SELECT id,country,jurisdiction,amount,currency,source,effective_from,effective_to,created_at
		FROM statutory_minimum_wages WHERE country=$1 AND effective_from<=$2 AND (effective_to IS NULL OR effective_to>=$2)
		ORDER BY effective_from DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, country, date)
	var w StatutoryMinimumWage
	err := row.Scan(&w.ID, &w.Country, &w.Jurisdiction, &w.Amount, &w.Currency, &w.Source, &w.EffectiveFrom, &w.EffectiveTo, &w.CreatedAt)
	if err != nil {
		return nil, repoErr("GetMinimumWage", err)
	}
	return &w, nil
}

// ========================================================================
// LIMITS
// ========================================================================

func (r *Repository) GetActiveLimits(ctx context.Context, companyID string, date time.Time) ([]PayrollLimit, error) {
	q := `SELECT id,company_id,concept_id,limit_type,minimum_amount,maximum_amount,effective_from,effective_to,created_at,updated_at
		FROM payroll_limits WHERE company_id=$1 AND effective_from<=$2 AND (effective_to IS NULL OR effective_to>=$2)`
	rows, err := r.pool.Query(ctx, q, companyID, date)
	if err != nil {
		return nil, repoErr("GetActiveLimits", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollLimit, error) {
		var l PayrollLimit
		err := row.Scan(&l.ID, &l.CompanyID, &l.ConceptID, &l.LimitType, &l.MinimumAmount, &l.MaximumAmount, &l.EffectiveFrom, &l.EffectiveTo, &l.CreatedAt, &l.UpdatedAt)
		return l, err
	})
}

// ========================================================================
// ADVANCES
// ========================================================================

func (r *Repository) CreateAdvance(ctx context.Context, a *EmployeeAdvance) error {
	q := `INSERT INTO employee_advances (id,company_id,employee_id,amount,currency,request_date,installments,
		installment_amount,remaining_amount,reason,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EmployeeID, a.Amount, a.Currency, a.RequestDate,
		a.Installments, a.InstallmentAmount, a.RemainingAmount, a.Reason, a.Status, a.CreatedBy)
	return repoErr("CreateAdvance", err)
}

func (r *Repository) ListAdvances(ctx context.Context, companyID, employeeID string) ([]EmployeeAdvance, error) {
	q := `SELECT id,company_id,employee_id,amount,currency,request_date,installments,installment_amount,
		remaining_amount,reason,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM employee_advances WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListAdvances", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (EmployeeAdvance, error) {
		var a EmployeeAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.Amount, &a.Currency, &a.RequestDate, &a.Installments,
			&a.InstallmentAmount, &a.RemainingAmount, &a.Reason, &a.Status, &a.ApprovedBy, &a.ApprovedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

// ========================================================================
// GARNISHMENTS
// ========================================================================

func (r *Repository) CreateGarnishment(ctx context.Context, g *PayrollGarnishment) error {
	q := `INSERT INTO payroll_garnishments (id,company_id,employee_id,court_order_number,court_name,type,percentage,
		fixed_amount,priority,effective_from,effective_to,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, g.ID, g.CompanyID, g.EmployeeID, g.CourtOrderNumber, g.CourtName, g.Type,
		g.Percentage, g.FixedAmount, g.Priority, g.EffectiveFrom, g.EffectiveTo, g.Status, g.CreatedBy)
	return repoErr("CreateGarnishment", err)
}

func (r *Repository) ListGarnishments(ctx context.Context, companyID, employeeID string) ([]PayrollGarnishment, error) {
	q := `SELECT id,company_id,employee_id,court_order_number,court_name,type,percentage,fixed_amount,priority,
		effective_from,effective_to,status,notes,created_by,created_at,updated_at
		FROM payroll_garnishments WHERE company_id=$1 AND employee_id=$2 ORDER BY priority, created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListGarnishments", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollGarnishment, error) {
		var g PayrollGarnishment
		err := row.Scan(&g.ID, &g.CompanyID, &g.EmployeeID, &g.CourtOrderNumber, &g.CourtName, &g.Type, &g.Percentage,
			&g.FixedAmount, &g.Priority, &g.EffectiveFrom, &g.EffectiveTo, &g.Status, &g.Notes, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
		return g, err
	})
}

// ========================================================================
// ERRORS
// ========================================================================

func (r *Repository) CreateError(ctx context.Context, e *PayrollError) error {
	q := `INSERT INTO payroll_errors (id,run_id,employee_id,severity,code,message,field) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.RunID, e.EmployeeID, e.Severity, e.Code, e.Message, e.Field)
	return repoErr("CreateError", err)
}

func (r *Repository) ListErrors(ctx context.Context, runID string) ([]PayrollError, error) {
	q := `SELECT id,run_id,employee_id,severity,code,message,field,resolved,resolved_at,created_at
		FROM payroll_errors WHERE run_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListErrors", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollError, error) {
		var e PayrollError
		err := row.Scan(&e.ID, &e.RunID, &e.EmployeeID, &e.Severity, &e.Code, &e.Message, &e.Field, &e.Resolved, &e.ResolvedAt, &e.CreatedAt)
		return e, err
	})
}

func (r *Repository) ListBlockingErrors(ctx context.Context, runID string) ([]PayrollError, error) {
	q := `SELECT id,run_id,employee_id,severity,code,message,field,resolved,resolved_at,created_at
		FROM payroll_errors WHERE run_id=$1 AND severity='BLOCKING' AND resolved=false ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListBlockingErrors", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (PayrollError, error) {
		var e PayrollError
		err := row.Scan(&e.ID, &e.RunID, &e.EmployeeID, &e.Severity, &e.Code, &e.Message, &e.Field, &e.Resolved, &e.ResolvedAt, &e.CreatedAt)
		return e, err
	})
}

// ========================================================================
// AUDIT
// ========================================================================

func (r *Repository) LogAudit(ctx context.Context, log *PayrollAuditLog) error {
	oldV, _ := json.Marshal(log.OldValues)
	newV, _ := json.Marshal(log.NewValues)
	q := `INSERT INTO payroll_audit_logs (id,company_id,user_id,action,entity_type,entity_id,old_values,new_values,ip_address,user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, log.ID, log.CompanyID, log.UserID, log.Action, log.EntityType, log.EntityID, oldV, newV, log.IPAddress, log.UserAgent)
	return repoErr("LogAudit", err)
}

// ========================================================================
// DASHBOARD & SUMMARY
// ========================================================================

func (r *Repository) GetRunSummary(ctx context.Context, runID string) (*PayrollSummary, error) {
	q := `SELECT COUNT(*),COUNT(*) FILTER (WHERE status IN ('CALCULATED','VALIDATED','APPROVED')),
		COUNT(*) FILTER (WHERE status='ERROR'),
		COALESCE(SUM(gross_remunerative+gross_non_remunerative),0),
		COALESCE(SUM(deductions_amount),0),
		COALESCE(SUM(net_amount),0),
		COALESCE(SUM(employer_contributions),0),
		COALESCE(SUM(employer_cost),0)
		FROM payroll_run_employees WHERE run_id=$1`
	row := r.pool.QueryRow(ctx, q, runID)
	var s PayrollSummary
	err := row.Scan(&s.TotalEmployees, &s.CalculatedEmployees, &s.ErrorEmployees,
		&s.TotalGross, &s.TotalDeductions, &s.TotalNet, &s.TotalContributions, &s.TotalEmployerCost)
	if err != nil {
		return nil, repoErr("GetRunSummary", err)
	}
	return &s, nil
}

func (r *Repository) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	var s DashboardStats
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payroll_periods WHERE company_id=$1 AND status NOT IN ('CLOSED','CANCELLED')`, companyID).Scan(&s.ActivePeriods)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payroll_runs WHERE company_id=$1 AND status NOT IN ('CLOSED','CANCELLED')`, companyID).Scan(&s.PendingRuns)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payroll_errors pr JOIN payroll_runs r ON r.id=pr.run_id WHERE r.company_id=$1`, companyID).Scan(&s.TotalErrors)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payroll_errors pr JOIN payroll_runs r ON r.id=pr.run_id WHERE r.company_id=$1 AND pr.severity='BLOCKING' AND pr.resolved=false`, companyID).Scan(&s.BlockingErrors)
	return &s, nil
}

// ========================================================================
// EMPLOYEE QUERIES
// ========================================================================

func (r *Repository) GetEmployeeCompensation(ctx context.Context, companyID, employeeID string) (decimal.Decimal, string, error) {
	var amount decimal.Decimal
	var currency string
	err := r.pool.QueryRow(ctx,
		`SELECT base_amount, COALESCE(currency,'ARS') FROM employee_compensations WHERE company_id=$1 AND employee_id=$2 AND status='active' ORDER BY created_at DESC LIMIT 1`,
		companyID, employeeID).Scan(&amount, &currency)
	if err != nil {
		return decimal.Zero, "ARS", repoErr("GetEmployeeCompensation", err)
	}
	return amount, currency, nil
}

func (r *Repository) GetEmployeeAgreementCategory(ctx context.Context, companyID, employeeID string) (agreementID, categoryID *string, err error) {
	q := `SELECT e.agreement_id, e.category_id FROM employee_positions ep
		JOIN positions p ON p.id=ep.position_id
		LEFT JOIN employees e ON e.id=ep.employee_id
		WHERE ep.employee_id=$1 AND ep.company_id=$2 AND (ep.end_date IS NULL OR ep.end_date>=CURRENT_DATE)
		LIMIT 1`
	row := r.pool.QueryRow(ctx, q, employeeID, companyID)
	err = row.Scan(&agreementID, &categoryID)
	if err != nil {
		return nil, nil, repoErr("GetEmployeeAgreementCategory", err)
	}
	return
}
