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

type ArcaRepo struct {
	pool *pgxpool.Pool
}

func NewArcaRepo(pool *pgxpool.Pool) *ArcaRepo {
	return &ArcaRepo{pool: pool}
}

func (r *ArcaRepo) CreateMapping(ctx context.Context, m *domain.ArcaConceptMapping) error {
	q := `INSERT INTO payroll_arca_concept_mappings (id,company_id,concept_id,arca_concept_code,arca_concept_name,
		mapping_type,percentage,is_taxable,is_contributable,notes,effective_from,effective_to,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, m.ID, m.CompanyID, m.ConceptID, m.ArcaConceptCode, m.ArcaConceptName,
		m.MappingType, m.Percentage, m.IsTaxable, m.IsContributable, m.Notes, m.EffectiveFrom, m.EffectiveTo, m.IsActive, m.CreatedBy)
	return repoErr("CreateMapping", err)
}

func (r *ArcaRepo) GetMapping(ctx context.Context, companyID, id uuid.UUID) (*domain.ArcaConceptMapping, error) {
	q := `SELECT id,company_id,concept_id,arca_concept_code,arca_concept_name,mapping_type,percentage,
		is_taxable,is_contributable,notes,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM payroll_arca_concept_mappings WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var m domain.ArcaConceptMapping
	err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.ArcaConceptCode, &m.ArcaConceptName,
		&m.MappingType, &m.Percentage, &m.IsTaxable, &m.IsContributable, &m.Notes,
		&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetMapping", err)
	}
	return &m, nil
}

func (r *ArcaRepo) ListMappings(ctx context.Context, companyID uuid.UUID) ([]domain.ArcaConceptMapping, error) {
	q := `SELECT id,company_id,concept_id,arca_concept_code,arca_concept_name,mapping_type,percentage,
		is_taxable,is_contributable,notes,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM payroll_arca_concept_mappings WHERE company_id=$1 ORDER BY arca_concept_code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListMappings", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ArcaConceptMapping, error) {
		var m domain.ArcaConceptMapping
		err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.ArcaConceptCode, &m.ArcaConceptName,
			&m.MappingType, &m.Percentage, &m.IsTaxable, &m.IsContributable, &m.Notes,
			&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *ArcaRepo) UpdateMapping(ctx context.Context, m *domain.ArcaConceptMapping) error {
	q := `UPDATE payroll_arca_concept_mappings SET arca_concept_code=$1,arca_concept_name=$2,mapping_type=$3,
		percentage=$4,is_taxable=$5,is_contributable=$6,notes=$7,effective_from=$8,effective_to=$9,
		is_active=$10,updated_at=NOW() WHERE id=$11 AND company_id=$12`
	_, err := r.pool.Exec(ctx, q, m.ArcaConceptCode, m.ArcaConceptName, m.MappingType,
		m.Percentage, m.IsTaxable, m.IsContributable, m.Notes, m.EffectiveFrom, m.EffectiveTo,
		m.IsActive, m.ID, m.CompanyID)
	return repoErr("UpdateMapping", err)
}

func (r *ArcaRepo) DeleteMapping(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_arca_concept_mappings WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteMapping", err)
}

func (r *ArcaRepo) GetActiveMappingsForConcept(ctx context.Context, companyID, conceptID uuid.UUID, date time.Time) ([]domain.ArcaConceptMapping, error) {
	q := `SELECT id,company_id,concept_id,arca_concept_code,arca_concept_name,mapping_type,percentage,
		is_taxable,is_contributable,notes,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM payroll_arca_concept_mappings WHERE company_id=$1 AND concept_id=$2 AND is_active=true
		AND effective_from<=$3 AND (effective_to IS NULL OR effective_to>=$3) ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID, conceptID, date)
	if err != nil {
		return nil, repoErr("GetActiveMappingsForConcept", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ArcaConceptMapping, error) {
		var m domain.ArcaConceptMapping
		err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.ArcaConceptCode, &m.ArcaConceptName,
			&m.MappingType, &m.Percentage, &m.IsTaxable, &m.IsContributable, &m.Notes,
			&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *ArcaRepo) GetMappingsByType(ctx context.Context, companyID uuid.UUID, mappingType string) ([]domain.ArcaConceptMapping, error) {
	q := `SELECT id,company_id,concept_id,arca_concept_code,arca_concept_name,mapping_type,percentage,
		is_taxable,is_contributable,notes,effective_from,effective_to,is_active,created_by,created_at,updated_at
		FROM payroll_arca_concept_mappings WHERE company_id=$1 AND mapping_type=$2 AND is_active=true ORDER BY arca_concept_code`
	rows, err := r.pool.Query(ctx, q, companyID, mappingType)
	if err != nil {
		return nil, repoErr("GetMappingsByType", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ArcaConceptMapping, error) {
		var m domain.ArcaConceptMapping
		err := row.Scan(&m.ID, &m.CompanyID, &m.ConceptID, &m.ArcaConceptCode, &m.ArcaConceptName,
			&m.MappingType, &m.Percentage, &m.IsTaxable, &m.IsContributable, &m.Notes,
			&m.EffectiveFrom, &m.EffectiveTo, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *ArcaRepo) CreateExport(ctx context.Context, e *domain.ArcaExport) error {
	q := `INSERT INTO payroll_arca_exports (id,company_id,run_id,export_type,period_id,file_name,
		file_content,storage_path,file_size,checksum,status,error_message,submission_date,
		acknowledgement_code,employee_count,total_amount,generated_by,generated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.RunID, e.ExportType, e.PeriodID, e.FileName,
		e.FileContent, e.StoragePath, e.FileSize, e.Checksum, e.Status, e.ErrorMessage, e.SubmissionDate,
		e.AcknowledgementCode, e.EmployeeCount, e.TotalAmount, e.GeneratedBy, e.GeneratedAt)
	return repoErr("CreateExport", err)
}

func (r *ArcaRepo) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.ArcaExport, error) {
	q := `SELECT id,company_id,run_id,export_type,period_id,file_name,file_content,storage_path,
		file_size,checksum,status,error_message,submission_date,acknowledgement_code,
		employee_count,total_amount,generated_by,generated_at,created_at,updated_at
		FROM payroll_arca_exports WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var e domain.ArcaExport
	err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.ExportType, &e.PeriodID, &e.FileName, &e.FileContent, &e.StoragePath,
		&e.FileSize, &e.Checksum, &e.Status, &e.ErrorMessage, &e.SubmissionDate, &e.AcknowledgementCode,
		&e.EmployeeCount, &e.TotalAmount, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetExport", err)
	}
	return &e, nil
}

func (r *ArcaRepo) ListExports(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, limit, offset int) ([]domain.ArcaExport, error) {
	q := `SELECT id,company_id,run_id,export_type,period_id,file_name,file_content,storage_path,
		file_size,checksum,status,error_message,submission_date,acknowledgement_code,
		employee_count,total_amount,generated_by,generated_at,created_at,updated_at
		FROM payroll_arca_exports WHERE company_id=$1`
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
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ArcaExport, error) {
		var e domain.ArcaExport
		err := row.Scan(&e.ID, &e.CompanyID, &e.RunID, &e.ExportType, &e.PeriodID, &e.FileName, &e.FileContent, &e.StoragePath,
			&e.FileSize, &e.Checksum, &e.Status, &e.ErrorMessage, &e.SubmissionDate, &e.AcknowledgementCode,
			&e.EmployeeCount, &e.TotalAmount, &e.GeneratedBy, &e.GeneratedAt, &e.CreatedAt, &e.UpdatedAt)
		return e, err
	})
}

func (r *ArcaRepo) UpdateExportStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_arca_exports SET status=$1`
	args := []any{status}
	n := 2
	for k, v := range fields {
		q += fmt.Sprintf(",%s=$%d", k, n)
		args = append(args, v)
		n++
	}
	q += fmt.Sprintf(",updated_at=NOW() WHERE id=$%d", n)
	args = append(args, id)
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("UpdateExportStatus", err)
}
