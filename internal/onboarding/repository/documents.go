package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/onboarding/domain"
)

type DocumentRepo struct {
	pool *pgxpool.Pool
}

func NewDocumentRepo(pool *pgxpool.Pool) *DocumentRepo {
	return &DocumentRepo{pool: pool}
}

func (r *DocumentRepo) Create(ctx context.Context, d *domain.OnboardingDocument) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_documents (company_id, onboarding_id, employee_id, document_type, name, required, status, expiration_date)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`,
		d.CompanyID, d.OnboardingID, d.EmployeeID, d.DocumentType, d.Name, d.Required, d.Status, d.ExpirationDate,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *DocumentRepo) GetByID(ctx context.Context, companyID, id string) (*domain.OnboardingDocument, error) {
	d := &domain.OnboardingDocument{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, onboarding_id, employee_id, document_type, name, required, status,
		        storage_key, mime_type, uploaded_at, verified_at, verified_by, expiration_date, created_at, updated_at
		 FROM onboarding_documents WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&d.ID, &d.CompanyID, &d.OnboardingID, &d.EmployeeID, &d.DocumentType, &d.Name,
		&d.Required, &d.Status, &d.StorageKey, &d.MimeType, &d.UploadedAt,
		&d.VerifiedAt, &d.VerifiedBy, &d.ExpirationDate, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

func (r *DocumentRepo) ListByOnboarding(ctx context.Context, onboardingID string) ([]domain.OnboardingDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, onboarding_id, employee_id, document_type, name, required, status,
		        storage_key, mime_type, uploaded_at, verified_at, verified_by, expiration_date, created_at, updated_at
		 FROM onboarding_documents WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ds []domain.OnboardingDocument
	for rows.Next() {
		var d domain.OnboardingDocument
		if err := rows.Scan(&d.ID, &d.CompanyID, &d.OnboardingID, &d.EmployeeID, &d.DocumentType, &d.Name,
			&d.Required, &d.Status, &d.StorageKey, &d.MimeType, &d.UploadedAt,
			&d.VerifiedAt, &d.VerifiedBy, &d.ExpirationDate, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *DocumentRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.DocStatus, reviewedBy *string) error {
	if reviewedBy != nil && status == domain.DocApproved {
		_, err := r.pool.Exec(ctx,
			`UPDATE onboarding_documents SET status=$3, verified_by=$4, verified_at=NOW(), updated_at=NOW()
			 WHERE company_id=$1 AND id=$2`, companyID, id, status, *reviewedBy)
		return err
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET status=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, status)
	return err
}

func (r *DocumentRepo) UpdateStorage(ctx context.Context, companyID, id, storageKey, mimeType string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET storage_key=$3, mime_type=$4, status='UPLOADED', uploaded_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, storageKey, mimeType)
	return err
}

func (r *DocumentRepo) CountPendingReview(ctx context.Context, companyID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_documents WHERE company_id=$1 AND status='UPLOADED'`, companyID).Scan(&count)
	return count, err
}

func (r *DocumentRepo) CreateVersion(ctx context.Context, v *domain.OnboardingDocumentVersion) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_document_versions (company_id, document_id, version, file_name, mime_type, size_bytes, storage_key, checksum, uploaded_by, notes)
		 VALUES ($1,$2,(SELECT COALESCE(MAX(version),0)+1 FROM onboarding_document_versions WHERE document_id=$2),$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, uploaded_at`,
		v.CompanyID, v.DocumentID, v.FileName, v.MimeType, v.SizeBytes, v.StorageKey, v.Checksum, v.UploadedBy, v.Notes,
	).Scan(&v.ID, &v.UploadedAt)
}

func (r *DocumentRepo) ListVersions(ctx context.Context, documentID string) ([]domain.OnboardingDocumentVersion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, document_id, version, file_name, mime_type, size_bytes, storage_key, checksum, uploaded_by, uploaded_at, notes
		 FROM onboarding_document_versions WHERE document_id=$1 ORDER BY version DESC`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vs []domain.OnboardingDocumentVersion
	for rows.Next() {
		var v domain.OnboardingDocumentVersion
		if err := rows.Scan(&v.ID, &v.CompanyID, &v.DocumentID, &v.Version, &v.FileName, &v.MimeType,
			&v.SizeBytes, &v.StorageKey, &v.Checksum, &v.UploadedBy, &v.UploadedAt, &v.Notes); err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
