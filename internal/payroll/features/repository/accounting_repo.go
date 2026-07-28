package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

type AccountingRepo struct {
	pool *pgxpool.Pool
}

func NewAccountingRepo(pool *pgxpool.Pool) *AccountingRepo {
	return &AccountingRepo{pool: pool}
}

func (r *AccountingRepo) CreateMapping(ctx context.Context, m *domain.AccountingAccountMapping) error {
	q := `INSERT INTO payroll_accounting_account_mappings (id,company_id,concept_id,mapping_type,
		debit_account,credit_account,cost_center_required,description_template,priority,
		effective_from,effective_to,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, m.ID, m.CompanyID, m.ConceptID, m.MappingType,
		m.DebitAccount, m.CreditAccount, m.CostCenterRequired, m.DescriptionTemplate, m.Priority,
		m.EffectiveFrom, m.EffectiveTo, m.IsActive, m.CreatedBy)
	return repoErr("CreateMapping", err)
}

func (r *AccountingRepo) GetMapping(ctx context.Context, companyID, id uuid.UUID) (*domain.AccountingAccountMapping, error) {
	q := `SELECT id,company_id,concept_id,mapping_type,debit_account,credit_account,
		cost_center_required,description_template,priority,effective_from,effective_to,
		is_active,created_by,created_at,updated_at
		FROM payroll_accounting_account_mappings WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var m domain.AccountingAccountMapping
	err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.MappingType,
		&m.DebitAccount, &m.CreditAccount, &m.CostCenterRequired, &m.DescriptionTemplate, &m.Priority,
		&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetMapping", err)
	}
	return &m, nil
}

func (r *AccountingRepo) ListMappings(ctx context.Context, companyID uuid.UUID) ([]domain.AccountingAccountMapping, error) {
	q := `SELECT id,company_id,concept_id,mapping_type,debit_account,credit_account,
		cost_center_required,description_template,priority,effective_from,effective_to,
		is_active,created_by,created_at,updated_at
		FROM payroll_accounting_account_mappings WHERE company_id=$1 ORDER BY priority, mapping_type`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListMappings", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.AccountingAccountMapping, error) {
		var m domain.AccountingAccountMapping
		err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.MappingType,
			&m.DebitAccount, &m.CreditAccount, &m.CostCenterRequired, &m.DescriptionTemplate, &m.Priority,
			&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *AccountingRepo) UpdateMapping(ctx context.Context, m *domain.AccountingAccountMapping) error {
	q := `UPDATE payroll_accounting_account_mappings SET mapping_type=$1,debit_account=$2,credit_account=$3,
		cost_center_required=$4,description_template=$5,priority=$6,effective_from=$7,effective_to=$8,
		is_active=$9,updated_at=NOW() WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, m.MappingType, m.DebitAccount, m.CreditAccount,
		m.CostCenterRequired, m.DescriptionTemplate, m.Priority, m.EffectiveFrom, m.EffectiveTo,
		m.IsActive, m.ID, m.CompanyID)
	return repoErr("UpdateMapping", err)
}

func (r *AccountingRepo) GetActiveForConcept(ctx context.Context, companyID uuid.UUID, date time.Time) ([]domain.AccountingAccountMapping, error) {
	q := `SELECT id,company_id,concept_id,mapping_type,debit_account,credit_account,
		cost_center_required,description_template,priority,effective_from,effective_to,
		is_active,created_by,created_at,updated_at
		FROM payroll_accounting_account_mappings WHERE company_id=$1 AND is_active=true
		AND effective_from<=$2 AND (effective_to IS NULL OR effective_to>=$2) ORDER BY priority, mapping_type`
	rows, err := r.pool.Query(ctx, q, companyID, date)
	if err != nil {
		return nil, repoErr("GetActiveForConcept", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.AccountingAccountMapping, error) {
		var m domain.AccountingAccountMapping
		err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.MappingType,
			&m.DebitAccount, &m.CreditAccount, &m.CostCenterRequired, &m.DescriptionTemplate, &m.Priority,
			&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *AccountingRepo) CreateExport(ctx context.Context, e *domain.AccountingExport) error {
	q := `INSERT INTO payroll_accounting_exports (id,company_id,run_id,period_id,export_type,file_format,
		file_name,file_content,storage_path,file_size,status,employee_count,total_debit,total_credit,
		entry_count,error_message,generated_by,generated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.RunID, e.PeriodID, e.ExportType, e.FileFormat,
		e.FileName, e.FileContent, e.StoragePath, e.FileSize, e.Status, e.EmployeeCount,
		e.TotalDebit, e.TotalCredit, e.EntryCount, e.ErrorMessage, e.GeneratedBy, e.GeneratedAt)
	return repoErr("CreateExport", err)
}

func (r *AccountingRepo) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.AccountingExport, error) {
	q := `SELECT id,company_id,run_id,period_id,export_type,file_format,file_name,file_content,
		storage_path,file_size,status,employee_count,total_debit,total_credit,entry_count,
		error_message,generated_by,generated_at,created_at
		FROM payroll_accounting_exports WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.AccountingExport
	err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.PeriodID, &e.ExportType, &e.FileFormat, &e.FileName, &e.FileContent,
		&e.StoragePath, &e.FileSize, &e.Status, &e.EmployeeCount, &e.TotalDebit, &e.TotalCredit, &e.EntryCount,
		&e.ErrorMessage, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
	if err != nil {
		return nil, repoErr("GetExport", err)
	}
	return &e, nil
}

func (r *AccountingRepo) ListExports(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, limit, offset int) ([]domain.AccountingExport, error) {
	q := `SELECT id,company_id,run_id,period_id,export_type,file_format,file_name,file_content,
		storage_path,file_size,status,employee_count,total_debit,total_credit,entry_count,
		error_message,generated_by,generated_at,created_at
		FROM payroll_accounting_exports WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if runID != nil {
		q += fmt.Sprintf(" AND run_id=$%d", n)
		args = append(args, *runID)
		n++
	}
	q += " ORDER BY created_at DESC"
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
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.AccountingExport, error) {
		var e domain.AccountingExport
		err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.PeriodID, &e.ExportType, &e.FileFormat, &e.FileName, &e.FileContent,
			&e.StoragePath, &e.FileSize, &e.Status, &e.EmployeeCount, &e.TotalDebit, &e.TotalCredit, &e.EntryCount,
			&e.ErrorMessage, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
		return e, err
	})
}

func (r *AccountingRepo) UpdateExportStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_accounting_exports SET status=$1`
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

func (r *AccountingRepo) CreateEntry(ctx context.Context, e *domain.AccountingEntry) error {
	q := `INSERT INTO payroll_accounting_entries (id,export_id,entry_number,account_code,account_name,
		cost_center,debit_amount,credit_amount,concept_code,concept_name,employee_id,employee_name,
		document_type,document_number,reference,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.ExportID, e.EntryNumber, e.AccountCode, e.AccountName,
		e.CostCenter, e.DebitAmount, e.CreditAmount, e.ConceptCode, e.ConceptName, e.EmployeeID, e.EmployeeName,
		e.DocumentType, e.DocumentNumber, e.Reference, e.SortOrder)
	return repoErr("CreateEntry", err)
}

func (r *AccountingRepo) BulkCreateEntries(ctx context.Context, entries []domain.AccountingEntry) error {
	if len(entries) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_accounting_entries (id,export_id,entry_number,account_code,account_name,
		cost_center,debit_amount,credit_amount,concept_code,concept_name,employee_id,employee_name,
		document_type,document_number,reference,sort_order) VALUES `
	args := []any{}
	n := 1
	for _, e := range entries {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14, n+15)
		args = append(args, e.ID, e.ExportID, e.EntryNumber, e.AccountCode, e.AccountName,
			e.CostCenter, e.DebitAmount, e.CreditAmount, e.ConceptCode, e.ConceptName, e.EmployeeID, e.EmployeeName,
			e.DocumentType, e.DocumentNumber, e.Reference, e.SortOrder)
		n += 16
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateEntries", err)
}

func (r *AccountingRepo) ListEntriesByExport(ctx context.Context, exportID uuid.UUID) ([]domain.AccountingEntry, error) {
	q := `SELECT id,export_id,entry_number,account_code,account_name,cost_center,debit_amount,credit_amount,
		concept_code,concept_name,employee_id,employee_name,document_type,document_number,reference,sort_order,created_at
		FROM payroll_accounting_entries WHERE export_id=$1 ORDER BY entry_number, sort_order`
	rows, err := r.pool.Query(ctx, q, exportID)
	if err != nil {
		return nil, repoErr("ListEntriesByExport", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.AccountingEntry, error) {
		var e domain.AccountingEntry
		err := row.Scan(&e.ID, &e.ExportID, &e.EntryNumber, &e.AccountCode, &e.AccountName,
			&e.CostCenter, &e.DebitAmount, &e.CreditAmount, &e.ConceptCode, &e.ConceptName,
			&e.EmployeeID, &e.EmployeeName, &e.DocumentType, &e.DocumentNumber, &e.Reference, &e.SortOrder, &e.CreatedAt)
		return e, err
	})
}
