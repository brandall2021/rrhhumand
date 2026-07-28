package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/shopspring/decimal"
)

type BookRepo struct {
	pool *pgxpool.Pool
}

func NewBookRepo(pool *pgxpool.Pool) *BookRepo {
	return &BookRepo{pool: pool}
}

func (r *BookRepo) CreateEntry(ctx context.Context, e *domain.BookEntry) error {
	q := `INSERT INTO payroll_book_entries (id,company_id,run_id,run_employee_id,employee_id,entry_type,
		cuil,surname,name,nationality,birth_date,sex,admission_date,discharge_date,
		category_code,category_name,agreement_code,agreement_name,work_type,work_place,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		days_worked,hours_worked,absences,status,book_number,page_number)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.RunID, e.RunEmployeeID, e.EmployeeID, e.EntryType,
		e.CUIL, e.Surname, e.Name, e.Nationality, e.BirthDate, e.Sex, e.AdmissionDate, e.DischargeDate,
		e.CategoryCode, e.CategoryName, e.AgreementCode, e.AgreementName, e.WorkType, e.WorkPlace,
		e.GrossRemunerative, e.GrossNonRemunerative, e.DeductionsTotal, e.ContributionsTotal, e.NetAmount, e.EmployerCost,
		e.DaysWorked, e.HoursWorked, e.Absences, e.Status, e.BookNumber, e.PageNumber)
	return repoErr("CreateEntry", err)
}

func (r *BookRepo) BulkCreateEntries(ctx context.Context, entries []domain.BookEntry) error {
	if len(entries) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_book_entries (id,company_id,run_id,run_employee_id,employee_id,entry_type,
		cuil,surname,name,nationality,birth_date,sex,admission_date,discharge_date,
		category_code,category_name,agreement_code,agreement_name,work_type,work_place,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		days_worked,hours_worked,absences,status,book_number,page_number) VALUES `
	args := []any{}
	n := 1
	for _, e := range entries {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14, n+15, n+16, n+17, n+18, n+19,
			n+20, n+21, n+22, n+23, n+24, n+25, n+26, n+27, n+28, n+29, n+30, n+31)
		args = append(args, e.ID, e.CompanyID, e.RunID, e.RunEmployeeID, e.EmployeeID, e.EntryType,
			e.CUIL, e.Surname, e.Name, e.Nationality, e.BirthDate, e.Sex, e.AdmissionDate, e.DischargeDate,
			e.CategoryCode, e.CategoryName, e.AgreementCode, e.AgreementName, e.WorkType, e.WorkPlace,
			e.GrossRemunerative, e.GrossNonRemunerative, e.DeductionsTotal, e.ContributionsTotal, e.NetAmount, e.EmployerCost,
			e.DaysWorked, e.HoursWorked, e.Absences, e.Status, e.BookNumber, e.PageNumber)
		n += 32
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateEntries", err)
}

func (r *BookRepo) GetEntry(ctx context.Context, companyID, id uuid.UUID) (*domain.BookEntry, error) {
	q := `SELECT id,company_id,run_id,run_employee_id,employee_id,entry_type,
		cuil,surname,name,nationality,birth_date,sex,admission_date,discharge_date,
		category_code,category_name,agreement_code,agreement_name,work_type,work_place,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		days_worked,hours_worked,absences,status,book_number,page_number,created_at,updated_at
		FROM payroll_book_entries WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.BookEntry
	err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.RunEmployeeID, &e.EmployeeID, &e.EntryType,
		&e.CUIL, &e.Surname, &e.Name, &e.Nationality, &e.BirthDate, &e.Sex, &e.AdmissionDate, &e.DischargeDate,
		&e.CategoryCode, &e.CategoryName, &e.AgreementCode, &e.AgreementName, &e.WorkType, &e.WorkPlace,
		&e.GrossRemunerative, &e.GrossNonRemunerative, &e.DeductionsTotal, &e.ContributionsTotal, &e.NetAmount, &e.EmployerCost,
		&e.DaysWorked, &e.HoursWorked, &e.Absences, &e.Status, &e.BookNumber, &e.PageNumber, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetEntry", err)
	}
	return &e, nil
}

func (r *BookRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]domain.BookEntry, error) {
	q := `SELECT id,company_id,run_id,run_employee_id,employee_id,entry_type,
		cuil,surname,name,nationality,birth_date,sex,admission_date,discharge_date,
		category_code,category_name,agreement_code,agreement_name,work_type,work_place,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		days_worked,hours_worked,absences,status,book_number,page_number,created_at,updated_at
		FROM payroll_book_entries WHERE run_id=$1 ORDER BY surname, name`
	rows, err := r.pool.Query(ctx, q, runID)
	if err != nil {
		return nil, repoErr("ListByRun", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BookEntry, error) {
		var e domain.BookEntry
		err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.RunEmployeeID, &e.EmployeeID, &e.EntryType,
			&e.CUIL, &e.Surname, &e.Name, &e.Nationality, &e.BirthDate, &e.Sex, &e.AdmissionDate, &e.DischargeDate,
			&e.CategoryCode, &e.CategoryName, &e.AgreementCode, &e.AgreementName, &e.WorkType, &e.WorkPlace,
			&e.GrossRemunerative, &e.GrossNonRemunerative, &e.DeductionsTotal, &e.ContributionsTotal, &e.NetAmount, &e.EmployerCost,
			&e.DaysWorked, &e.HoursWorked, &e.Absences, &e.Status, &e.BookNumber, &e.PageNumber, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *BookRepo) ListByEmployee(ctx context.Context, companyID, employeeID uuid.UUID, limit, offset int) ([]domain.BookEntry, error) {
	q := `SELECT id,company_id,run_id,run_employee_id,employee_id,entry_type,
		cuil,surname,name,nationality,birth_date,sex,admission_date,discharge_date,
		category_code,category_name,agreement_code,agreement_name,work_type,work_place,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		days_worked,hours_worked,absences,status,book_number,page_number,created_at,updated_at
		FROM payroll_book_entries WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC`
	args := []any{companyID, employeeID}
	n := 3
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListByEmployee", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BookEntry, error) {
		var e domain.BookEntry
		err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.RunEmployeeID, &e.EmployeeID, &e.EntryType,
			&e.CUIL, &e.Surname, &e.Name, &e.Nationality, &e.BirthDate, &e.Sex, &e.AdmissionDate, &e.DischargeDate,
			&e.CategoryCode, &e.CategoryName, &e.AgreementCode, &e.AgreementName, &e.WorkType, &e.WorkPlace,
			&e.GrossRemunerative, &e.GrossNonRemunerative, &e.DeductionsTotal, &e.ContributionsTotal, &e.NetAmount, &e.EmployerCost,
			&e.DaysWorked, &e.HoursWorked, &e.Absences, &e.Status, &e.BookNumber, &e.PageNumber, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *BookRepo) CreateExport(ctx context.Context, e *domain.BookExport) error {
	q := `INSERT INTO payroll_book_exports (id,company_id,period_id,year,month,export_type,file_name,
		file_content,storage_path,file_size,status,submission_date,acknowledgement_code,
		employee_count,total_gross,total_deductions,total_net,generated_by,generated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.PeriodID, e.Year, e.Month, e.ExportType, e.FileName,
		e.FileContent, e.StoragePath, e.FileSize, e.Status, e.SubmissionDate, e.AcknowledgementCode,
		e.EmployeeCount, e.TotalGross, e.TotalDeductions, e.TotalNet, e.GeneratedBy, e.GeneratedAt)
	return repoErr("CreateExport", err)
}

func (r *BookRepo) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.BookExport, error) {
	q := `SELECT id,company_id,period_id,year,month,export_type,file_name,file_content,storage_path,
		file_size,status,submission_date,acknowledgement_code,employee_count,total_gross,total_deductions,total_net,
		generated_by,generated_at,created_at
		FROM payroll_book_exports WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.BookExport
	err := row.Scan(&e.ID, &e.CompanyID, &e.PeriodID, &e.Year, &e.Month, &e.ExportType, &e.FileName, &e.FileContent, &e.StoragePath,
		&e.FileSize, &e.Status, &e.SubmissionDate, &e.AcknowledgementCode, &e.EmployeeCount,
		&e.TotalGross, &e.TotalDeductions, &e.TotalNet,
		&e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
	if err != nil {
		return nil, repoErr("GetExport", err)
	}
	return &e, nil
}

func (r *BookRepo) ListExports(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]domain.BookExport, error) {
	q := `SELECT id,company_id,period_id,year,month,export_type,file_name,file_content,storage_path,
		file_size,status,submission_date,acknowledgement_code,employee_count,total_gross,total_deductions,total_net,
		generated_by,generated_at,created_at
		FROM payroll_book_exports WHERE company_id=$1 ORDER BY year DESC, month DESC`
	args := []any{companyID}
	n := 2
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListExports", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BookExport, error) {
		var e domain.BookExport
		err := row.Scan(&e.ID, &e.CompanyID, &e.PeriodID, &e.Year, &e.Month, &e.ExportType, &e.FileName, &e.FileContent, &e.StoragePath,
			&e.FileSize, &e.Status, &e.SubmissionDate, &e.AcknowledgementCode, &e.EmployeeCount,
			&e.TotalGross, &e.TotalDeductions, &e.TotalNet,
			&e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
		return e, err
	})
}

func (r *BookRepo) UpdateExportStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_book_exports SET status=$1`
	args := []any{status}
	n := 2
	for k, v := range fields {
		q += fmt.Sprintf(",%s=$%d", k, n)
		args = append(args, v)
		n++
	}
	q += fmt.Sprintf(" WHERE id=$%d", n)
	args = append(args, id)
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("UpdateExportStatus", err)
}

type BookSummary struct {
	TotalEmployees      int             `json:"total_employees"`
	TotalRemunerative   decimal.Decimal `json:"total_remunerative"`
	TotalNonRemunerative decimal.Decimal `json:"total_non_remunerative"`
	TotalDeductions     decimal.Decimal `json:"total_deductions"`
	TotalContributions  decimal.Decimal `json:"total_contributions"`
	TotalNet            decimal.Decimal `json:"total_net"`
	TotalEmployerCost   decimal.Decimal `json:"total_employer_cost"`
}

func (r *BookRepo) GetBookSummary(ctx context.Context, runID uuid.UUID) (*BookSummary, error) {
	q := `SELECT COUNT(*),
		COALESCE(SUM(gross_remunerative),0),
		COALESCE(SUM(gross_non_remunerative),0),
		COALESCE(SUM(deductions_total),0),
		COALESCE(SUM(contributions_total),0),
		COALESCE(SUM(net_amount),0),
		COALESCE(SUM(employer_cost),0)
		FROM payroll_book_entries WHERE run_id=$1`
	row := r.pool.QueryRow(ctx, q, runID)
	var s BookSummary
	err := row.Scan(&s.TotalEmployees, &s.TotalRemunerative, &s.TotalNonRemunerative,
		&s.TotalDeductions, &s.TotalContributions, &s.TotalNet, &s.TotalEmployerCost)
	if err != nil {
		return nil, repoErr("GetBookSummary", err)
	}
	return &s, nil
}
