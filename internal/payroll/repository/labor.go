package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreateAgreement(ctx context.Context, a *domain.LaborAgreement) error {
	q := `INSERT INTO labor_agreements (id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.Code, a.Name, a.Description, a.Activity, a.EffectiveFrom, a.EffectiveTo, a.Status, a.CreatedBy)
	return repoErr("CreateAgreement", err)
}

func (r *Repository) GetAgreement(ctx context.Context, companyID, id string) (*domain.LaborAgreement, error) {
	q := `SELECT id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by,created_at,updated_at
		FROM labor_agreements WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var a domain.LaborAgreement
	err := row.Scan(&a.ID, &a.CompanyID, &a.Code, &a.Name, &a.Description, &a.Activity, &a.EffectiveFrom, &a.EffectiveTo, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetAgreement", err)
	}
	return &a, nil
}

func (r *Repository) ListAgreements(ctx context.Context, companyID string) ([]domain.LaborAgreement, error) {
	q := `SELECT id,company_id,code,name,description,activity,effective_from,effective_to,status,created_by,created_at,updated_at
		FROM labor_agreements WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListAgreements", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.LaborAgreement, error) {
		var a domain.LaborAgreement
		err := row.Scan(&a.ID, &a.CompanyID, &a.Code, &a.Name, &a.Description, &a.Activity, &a.EffectiveFrom, &a.EffectiveTo, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		return a, err
	})
}

func (r *Repository) CreateCategory(ctx context.Context, c *domain.LaborCategory) error {
	q := `INSERT INTO labor_categories (id,company_id,agreement_id,code,name,description,effective_from,effective_to)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.AgreementID, c.Code, c.Name, c.Description, c.EffectiveFrom, c.EffectiveTo)
	return repoErr("CreateCategory", err)
}

func (r *Repository) ListCategories(ctx context.Context, companyID string) ([]domain.LaborCategory, error) {
	q := `SELECT id,company_id,agreement_id,code,name,description,effective_from,effective_to,created_at,updated_at
		FROM labor_categories WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCategories", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.LaborCategory, error) {
		var c domain.LaborCategory
		err := row.Scan(&c.ID, &c.CompanyID, &c.AgreementID, &c.Code, &c.Name, &c.Description, &c.EffectiveFrom, &c.EffectiveTo, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *Repository) CreateSalaryScale(ctx context.Context, s *domain.SalaryScale) error {
	q := `INSERT INTO salary_scales (id,company_id,agreement_id,category_id,minimum_salary,maximum_salary,reference_value,effective_from,effective_to)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.AgreementID, s.CategoryID, s.MinimumSalary, s.MaximumSalary, s.ReferenceValue, s.EffectiveFrom, s.EffectiveTo)
	return repoErr("CreateSalaryScale", err)
}

func (r *Repository) ListSalaryScales(ctx context.Context, companyID string) ([]domain.SalaryScale, error) {
	q := `SELECT id,company_id,agreement_id,category_id,minimum_salary,maximum_salary,reference_value,effective_from,effective_to,created_at,updated_at
		FROM salary_scales WHERE company_id=$1 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListSalaryScales", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.SalaryScale, error) {
		var s domain.SalaryScale
		err := row.Scan(&s.ID, &s.CompanyID, &s.AgreementID, &s.CategoryID, &s.MinimumSalary, &s.MaximumSalary, &s.ReferenceValue, &s.EffectiveFrom, &s.EffectiveTo, &s.CreatedAt, &s.UpdatedAt)
		return s, err
	})
}

func (r *Repository) GetMinimumWage(ctx context.Context, country string, date time.Time) (*domain.StatutoryMinimumWage, error) {
	q := `SELECT id,country,jurisdiction,amount,currency,source,effective_from,effective_to,created_at
		FROM statutory_minimum_wages WHERE country=$1 AND effective_from<=$2 AND (effective_to IS NULL OR effective_to>=$2)
		ORDER BY effective_from DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, country, date)
	var w domain.StatutoryMinimumWage
	err := row.Scan(&w.ID, &w.Country, &w.Jurisdiction, &w.Amount, &w.Currency, &w.Source, &w.EffectiveFrom, &w.EffectiveTo, &w.CreatedAt)
	if err != nil {
		return nil, repoErr("GetMinimumWage", err)
	}
	return &w, nil
}

func (r *Repository) GetActiveLimits(ctx context.Context, companyID string, date time.Time) ([]domain.PayrollLimit, error) {
	q := `SELECT id,company_id,concept_id,limit_type,minimum_amount,maximum_amount,effective_from,effective_to,created_at,updated_at
		FROM payroll_limits WHERE company_id=$1 AND effective_from<=$2 AND (effective_to IS NULL OR effective_to>=$2)`
	rows, err := r.pool.Query(ctx, q, companyID, date)
	if err != nil {
		return nil, repoErr("GetActiveLimits", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollLimit, error) {
		var l domain.PayrollLimit
		err := row.Scan(&l.ID, &l.CompanyID, &l.ConceptID, &l.LimitType, &l.MinimumAmount, &l.MaximumAmount, &l.EffectiveFrom, &l.EffectiveTo, &l.CreatedAt, &l.UpdatedAt)
		return l, err
	})
}

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
