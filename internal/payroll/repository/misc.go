package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreateAdvance(ctx context.Context, a *domain.EmployeeAdvance) error {
	q := `INSERT INTO employee_advances (id,company_id,employee_id,amount,currency,request_date,installments,
		installment_amount,remaining_amount,reason,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EmployeeID, a.Amount, a.Currency, a.RequestDate,
		a.Installments, a.InstallmentAmount, a.RemainingAmount, a.Reason, a.Status, a.CreatedBy)
	return repoErr("CreateAdvance", err)
}

func (r *Repository) ListAdvances(ctx context.Context, companyID, employeeID string) ([]domain.EmployeeAdvance, error) {
	q := `SELECT id,company_id,employee_id,amount,currency,request_date,installments,installment_amount,
		remaining_amount,reason,status,approved_by,approved_at,created_by,created_at,updated_at
		FROM employee_advances WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListAdvances", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.EmployeeAdvance, error) {
		var a domain.EmployeeAdvance
		err := row.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.Amount, &a.Currency, &a.RequestDate, &a.Installments,
			&a.InstallmentAmount, &a.RemainingAmount, &a.Reason, &a.Status, &a.ApprovedBy, &a.ApprovedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *Repository) CreateGarnishment(ctx context.Context, g *domain.PayrollGarnishment) error {
	q := `INSERT INTO payroll_garnishments (id,company_id,employee_id,court_order_number,court_name,type,percentage,
		fixed_amount,priority,effective_from,effective_to,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, g.ID, g.CompanyID, g.EmployeeID, g.CourtOrderNumber, g.CourtName, g.Type,
		g.Percentage, g.FixedAmount, g.Priority, g.EffectiveFrom, g.EffectiveTo, g.Status, g.CreatedBy)
	return repoErr("CreateGarnishment", err)
}

func (r *Repository) ListGarnishments(ctx context.Context, companyID, employeeID string) ([]domain.PayrollGarnishment, error) {
	q := `SELECT id,company_id,employee_id,court_order_number,court_name,type,percentage,fixed_amount,priority,
		effective_from,effective_to,status,notes,created_by,created_at,updated_at
		FROM payroll_garnishments WHERE company_id=$1 AND employee_id=$2 ORDER BY priority, created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListGarnishments", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollGarnishment, error) {
		var g domain.PayrollGarnishment
		err := row.Scan(&g.ID, &g.CompanyID, &g.EmployeeID, &g.CourtOrderNumber, &g.CourtName, &g.Type, &g.Percentage,
			&g.FixedAmount, &g.Priority, &g.EffectiveFrom, &g.EffectiveTo, &g.Status, &g.Notes, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
		return g, err
	})
}

func (r *Repository) CreateError(ctx context.Context, e *domain.PayrollError) error {
	q := `INSERT INTO payroll_errors (id,run_id,employee_id,severity,code,message,field) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.RunID, e.EmployeeID, e.Severity, e.Code, e.Message, e.Field)
	return repoErr("CreateError", err)
}

func (r *Repository) ListErrors(ctx context.Context, runID string) ([]domain.PayrollError, error) {
	q := `SELECT id,run_id,employee_id,severity,code,message,field,resolved,resolved_at,created_at
		FROM payroll_errors WHERE run_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListErrors", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollError, error) {
		var e domain.PayrollError
		err := row.Scan(&e.ID, &e.RunID, &e.EmployeeID, &e.Severity, &e.Code, &e.Message, &e.Field, &e.Resolved, &e.ResolvedAt, &e.CreatedAt)
		return e, err
	})
}

func (r *Repository) ListBlockingErrors(ctx context.Context, runID string) ([]domain.PayrollError, error) {
	q := `SELECT id,run_id,employee_id,severity,code,message,field,resolved,resolved_at,created_at
		FROM payroll_errors WHERE run_id=$1 AND severity='BLOCKING' AND resolved=false ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListBlockingErrors", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollError, error) {
		var e domain.PayrollError
		err := row.Scan(&e.ID, &e.RunID, &e.EmployeeID, &e.Severity, &e.Code, &e.Message, &e.Field, &e.Resolved, &e.ResolvedAt, &e.CreatedAt)
		return e, err
	})
}

func (r *Repository) LogAudit(ctx context.Context, log *domain.PayrollAuditLog) error {
	oldV, _ := json.Marshal(log.OldValues)
	newV, _ := json.Marshal(log.NewValues)
	q := `INSERT INTO payroll_audit_logs (id,company_id,user_id,action,entity_type,entity_id,old_values,new_values,ip_address,user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, log.ID, log.CompanyID, log.UserID, log.Action, log.EntityType, log.EntityID, oldV, newV, log.IPAddress, log.UserAgent)
	return repoErr("LogAudit", err)
}

func (r *Repository) GetRunSummary(ctx context.Context, runID string) (*domain.PayrollSummary, error) {
	q := `SELECT COUNT(*),COUNT(*) FILTER (WHERE status IN ('CALCULATED','VALIDATED','APPROVED')),
		COUNT(*) FILTER (WHERE status='ERROR'),
		COALESCE(SUM(gross_remunerative+gross_non_remunerative),0),
		COALESCE(SUM(deductions_amount),0),
		COALESCE(SUM(net_amount),0),
		COALESCE(SUM(employer_contributions),0),
		COALESCE(SUM(employer_cost),0)
		FROM payroll_run_employees WHERE run_id=$1`
	row := r.pool.QueryRow(ctx, q, runID)
	var s domain.PayrollSummary
	err := row.Scan(&s.TotalEmployees, &s.CalculatedEmployees, &s.ErrorEmployees,
		&s.TotalGross, &s.TotalDeductions, &s.TotalNet, &s.TotalContributions, &s.TotalEmployerCost)
	if err != nil {
		return nil, repoErr("GetRunSummary", err)
	}
	return &s, nil
}

func (r *Repository) GetDashboardStats(ctx context.Context, companyID string) (*domain.DashboardStats, error) {
	var s domain.DashboardStats
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll_periods WHERE company_id=$1 AND status NOT IN ('CLOSED','CANCELLED')`,
		companyID).Scan(&s.ActivePeriods)
	if err != nil {
		return nil, repoErr("GetDashboardStats", err)
	}
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll_runs WHERE company_id=$1 AND status NOT IN ('CLOSED','CANCELLED')`,
		companyID).Scan(&s.PendingRuns)
	if err != nil {
		return nil, repoErr("GetDashboardStats", err)
	}
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll_errors pr JOIN payroll_runs r ON r.id=pr.run_id WHERE r.company_id=$1`,
		companyID).Scan(&s.TotalErrors)
	if err != nil {
		return nil, repoErr("GetDashboardStats", err)
	}
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll_errors pr JOIN payroll_runs r ON r.id=pr.run_id WHERE r.company_id=$1 AND pr.severity='BLOCKING' AND pr.resolved=false`,
		companyID).Scan(&s.BlockingErrors)
	if err != nil {
		return nil, repoErr("GetDashboardStats", err)
	}
	return &s, nil
}
