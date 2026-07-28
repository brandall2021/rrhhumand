package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

func (r *ReportRepo) CreateTemplate(ctx context.Context, t *domain.ReportTemplate) error {
	q := `INSERT INTO payroll_report_templates (id,company_id,name,description,report_type,config,
		is_default,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CompanyID, t.Name, t.Description, t.ReportType, t.Config,
		t.IsDefault, t.IsActive, t.CreatedBy)
	return repoErr("CreateTemplate", err)
}

func (r *ReportRepo) GetTemplate(ctx context.Context, companyID, id uuid.UUID) (*domain.ReportTemplate, error) {
	q := `SELECT id,company_id,name,description,report_type,config,is_default,is_active,created_by,created_at,updated_at
		FROM payroll_report_templates WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var t domain.ReportTemplate
	err := row.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.ReportType, &t.Config,
		&t.IsDefault, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetTemplate", err)
	}
	return &t, nil
}

func (r *ReportRepo) ListTemplates(ctx context.Context, companyID uuid.UUID) ([]domain.ReportTemplate, error) {
	q := `SELECT id,company_id,name,description,report_type,config,is_default,is_active,created_by,created_at,updated_at
		FROM payroll_report_templates WHERE company_id=$1 AND is_active=true ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListTemplates", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ReportTemplate, error) {
		var t domain.ReportTemplate
		err := row.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.ReportType, &t.Config,
			&t.IsDefault, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *ReportRepo) UpdateTemplate(ctx context.Context, t *domain.ReportTemplate) error {
	q := `UPDATE payroll_report_templates SET name=$1,description=$2,report_type=$3,config=$4,
		is_default=$5,is_active=$6,updated_at=NOW() WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, q, t.Name, t.Description, t.ReportType, t.Config,
		t.IsDefault, t.IsActive, t.ID, t.CompanyID)
	return repoErr("UpdateTemplate", err)
}

func (r *ReportRepo) DeleteTemplate(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_report_templates WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteTemplate", err)
}

func (r *ReportRepo) CreateExport(ctx context.Context, e *domain.ReportExport) error {
	q := `INSERT INTO payroll_report_exports (id,company_id,run_id,template_id,report_type,file_format,
		file_name,file_content,storage_path,file_size,status,error_message,config,generated_by,generated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.RunID, e.TemplateID, e.ReportType, e.FileFormat,
		e.FileName, e.FileContent, e.StoragePath, e.FileSize, e.Status, e.ErrorMessage, e.Config, e.GeneratedBy, e.GeneratedAt)
	return repoErr("CreateExport", err)
}

func (r *ReportRepo) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.ReportExport, error) {
	q := `SELECT id,company_id,run_id,template_id,report_type,file_format,file_name,file_content,
		storage_path,file_size,status,error_message,config,generated_by,generated_at,created_at
		FROM payroll_report_exports WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.ReportExport
	err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.TemplateID, &e.ReportType, &e.FileFormat, &e.FileName, &e.FileContent,
		&e.StoragePath, &e.FileSize, &e.Status, &e.ErrorMessage, &e.Config, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
	if err != nil {
		return nil, repoErr("GetExport", err)
	}
	return &e, nil
}

func (r *ReportRepo) ListExports(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]domain.ReportExport, error) {
	q := `SELECT id,company_id,run_id,template_id,report_type,file_format,file_name,file_content,
		storage_path,file_size,status,error_message,config,generated_by,generated_at,created_at
		FROM payroll_report_exports WHERE company_id=$1 ORDER BY created_at DESC`
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
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ReportExport, error) {
		var e domain.ReportExport
		err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.TemplateID, &e.ReportType, &e.FileFormat, &e.FileName, &e.FileContent,
			&e.StoragePath, &e.FileSize, &e.Status, &e.ErrorMessage, &e.Config, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt)
		return e, err
	})
}
