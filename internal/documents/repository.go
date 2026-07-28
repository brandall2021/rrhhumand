package documents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *models.Document) error {
	query := `
		INSERT INTO documents (id, company_id, category_id, employee_id, department_id, uploaded_by,
			title, description, original_filename, storage_key, mime_type, file_size, checksum,
			status, is_public, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		doc.ID, doc.CompanyID, doc.CategoryID, doc.EmployeeID, doc.DepartmentID, doc.UploadedBy,
		doc.Title, doc.Description, doc.OriginalFilename, doc.StorageKey, doc.MimeType,
		doc.FileSize, doc.Checksum, doc.Status, doc.IsPublic, doc.ExpiresAt,
	).Scan(&doc.CreatedAt, &doc.UpdatedAt)
}

func (r *DocumentRepository) GetDocumentByID(ctx context.Context, id, companyID string) (*models.Document, error) {
	query := `
		SELECT
			d.id, d.company_id, d.category_id,
			(SELECT name FROM document_categories WHERE id = d.category_id),
			d.employee_id, d.department_id, d.uploaded_by,
			u.first_name || ' ' || u.last_name,
			d.title, d.description, d.original_filename, d.storage_key,
			d.mime_type, d.file_size, d.checksum, d.status, d.is_public,
			d.expires_at, d.created_at, d.updated_at
		FROM documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE d.id = $1 AND d.company_id = $2`

	doc := &models.Document{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&doc.ID, &doc.CompanyID, &doc.CategoryID, &doc.CategoryName,
		&doc.EmployeeID, &doc.DepartmentID, &doc.UploadedBy, &doc.UploadedByName,
		&doc.Title, &doc.Description, &doc.OriginalFilename, &doc.StorageKey,
		&doc.MimeType, &doc.FileSize, &doc.Checksum, &doc.Status, &doc.IsPublic,
		&doc.ExpiresAt, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *DocumentRepository) ListDocuments(ctx context.Context, companyID string, filters DocumentFilters, offset, limit int) ([]models.Document, int64, error) {
	where := []string{`d.company_id = $1`, `d.status != 'DELETED'`}
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		where = append(where, fmt.Sprintf(`d.status = $%d`, argIdx))
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.CategoryID != "" {
		where = append(where, fmt.Sprintf(`d.category_id = $%d`, argIdx))
		args = append(args, filters.CategoryID)
		argIdx++
	}
	if filters.EmployeeID != "" {
		where = append(where, fmt.Sprintf(`d.employee_id = $%d`, argIdx))
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.DepartmentID != "" {
		where = append(where, fmt.Sprintf(`d.department_id = $%d`, argIdx))
		args = append(args, filters.DepartmentID)
		argIdx++
	}
	if filters.MimeType != "" {
		where = append(where, fmt.Sprintf(`d.mime_type LIKE $%d`, argIdx))
		args = append(args, "%"+filters.MimeType+"%")
		argIdx++
	}
	if filters.Search != "" {
		where = append(where, fmt.Sprintf(`(d.title ILIKE $%d OR d.description ILIKE $%d OR d.original_filename ILIKE $%d)`, argIdx, argIdx, argIdx))
		args = append(args, "%"+filters.Search+"%")
		argIdx++
	}
	if filters.CreatedFrom != "" {
		where = append(where, fmt.Sprintf(`d.created_at >= $%d`, argIdx))
		args = append(args, filters.CreatedFrom)
		argIdx++
	}
	if filters.CreatedTo != "" {
		where = append(where, fmt.Sprintf(`d.created_at <= $%d`, argIdx))
		args = append(args, filters.CreatedTo)
		argIdx++
	}
	if filters.Tag != "" {
		where = append(where, fmt.Sprintf(`EXISTS (SELECT 1 FROM document_tag_relations dtr JOIN document_tags dt ON dt.id = dtr.tag_id WHERE dtr.document_id = d.id AND dt.name = $%d)`, argIdx))
		args = append(args, filters.Tag)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM documents d WHERE %s`, whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT
			d.id, d.company_id, d.category_id,
			(SELECT name FROM document_categories WHERE id = d.category_id),
			d.employee_id, d.department_id, d.uploaded_by,
			u.first_name || ' ' || u.last_name,
			d.title, d.description, d.original_filename, d.storage_key,
			d.mime_type, d.file_size, d.checksum, d.status, d.is_public,
			d.expires_at, d.created_at, d.updated_at
		FROM documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE %s
		ORDER BY d.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.CategoryID, &d.CategoryName,
			&d.EmployeeID, &d.DepartmentID, &d.UploadedBy, &d.UploadedByName,
			&d.Title, &d.Description, &d.OriginalFilename, &d.StorageKey,
			&d.MimeType, &d.FileSize, &d.Checksum, &d.Status, &d.IsPublic,
			&d.ExpiresAt, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		docs = append(docs, d)
	}
	return docs, total, nil
}

func (r *DocumentRepository) UpdateDocument(ctx context.Context, doc *models.Document) error {
	query := `
		UPDATE documents SET title=$1, description=$2, category_id=$3, is_public=$4,
			expires_at=$5, updated_at=NOW()
		WHERE id=$6 AND company_id=$7`
	_, err := r.pool.Exec(ctx, query,
		doc.Title, doc.Description, doc.CategoryID, doc.IsPublic,
		doc.ExpiresAt, doc.ID, doc.CompanyID,
	)
	return err
}

func (r *DocumentRepository) UpdateDocumentStatus(ctx context.Context, id, companyID, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE documents SET status=$1, updated_at=NOW() WHERE id=$2 AND company_id=$3`,
		status, id, companyID,
	)
	return err
}

func (r *DocumentRepository) DeleteDocument(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM documents WHERE id=$1 AND company_id=$2`, id, companyID,
	)
	return err
}

func (r *DocumentRepository) CreateVersion(ctx context.Context, v *models.DocumentVersion) error {
	query := `
		INSERT INTO document_versions (id, document_id, version, storage_key, original_filename, mime_type, file_size, checksum, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query,
		v.ID, v.DocumentID, v.Version, v.StorageKey, v.OriginalFilename,
		v.MimeType, v.FileSize, v.Checksum, v.UploadedBy,
	).Scan(&v.CreatedAt)
}

func (r *DocumentRepository) GetLatestVersionNumber(ctx context.Context, documentID string) (int, error) {
	var version int
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(version), 0) FROM document_versions WHERE document_id=$1`, documentID,
	).Scan(&version)
	return version, err
}

func (r *DocumentRepository) ListVersions(ctx context.Context, documentID string) ([]models.DocumentVersion, error) {
	query := `
		SELECT id, document_id, version, storage_key, original_filename, mime_type, file_size, checksum, uploaded_by, created_at
		FROM document_versions WHERE document_id=$1 ORDER BY version DESC`
	rows, err := r.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.DocumentVersion
	for rows.Next() {
		var v models.DocumentVersion
		if err := rows.Scan(&v.ID, &v.DocumentID, &v.Version, &v.StorageKey, &v.OriginalFilename,
			&v.MimeType, &v.FileSize, &v.Checksum, &v.UploadedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (r *DocumentRepository) GetVersionByID(ctx context.Context, id string) (*models.DocumentVersion, error) {
	query := `
		SELECT id, document_id, version, storage_key, original_filename, mime_type, file_size, checksum, uploaded_by, created_at
		FROM document_versions WHERE id=$1`
	v := &models.DocumentVersion{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.DocumentID, &v.Version, &v.StorageKey, &v.OriginalFilename,
		&v.MimeType, &v.FileSize, &v.Checksum, &v.UploadedBy, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *DocumentRepository) SetPermissions(ctx context.Context, documentID string, perms []models.DocumentPermission) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM document_permissions WHERE document_id=$1`, documentID)
	if err != nil {
		return err
	}

	for _, p := range perms {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO document_permissions (id, document_id, grantee_type, grantee_id, can_read, can_download, can_share, can_manage)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)`,
			documentID, p.GranteeType, p.GranteeID, p.CanRead, p.CanDownload, p.CanShare, p.CanManage,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DocumentRepository) ListPermissions(ctx context.Context, documentID string) ([]models.DocumentPermission, error) {
	query := `
		SELECT id, document_id, grantee_type, grantee_id, can_read, can_download, can_share, can_manage, created_at
		FROM document_permissions WHERE document_id=$1`
	rows, err := r.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.DocumentPermission
	for rows.Next() {
		var p models.DocumentPermission
		if err := rows.Scan(&p.ID, &p.DocumentID, &p.GranteeType, &p.GranteeID,
			&p.CanRead, &p.CanDownload, &p.CanShare, &p.CanManage, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *DocumentRepository) HasAccess(ctx context.Context, documentID, userID, companyID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM documents d
			WHERE d.id = $1 AND d.company_id = $2
			AND (
				d.is_public = true
				OR d.uploaded_by = $3
				OR EXISTS (
					SELECT 1 FROM document_permissions dp
					WHERE dp.document_id = d.id
					AND dp.grantee_type = 'USER'
					AND dp.grantee_id = $3
					AND dp.can_read = true
				)
				OR EXISTS (
					SELECT 1 FROM user_roles ur
					JOIN document_permissions dp ON dp.document_id = d.id
					WHERE ur.user_id = $3
					AND ur.company_id = $2
					AND dp.grantee_type = 'ROLE'
					AND dp.grantee_id = ur.role_id
					AND dp.can_read = true
				)
			)
		)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, documentID, userID, companyID).Scan(&exists)
	return exists, err
}

func (r *DocumentRepository) CanDownload(ctx context.Context, documentID, userID, companyID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM documents d
			WHERE d.id = $1 AND d.company_id = $2
			AND (
				d.uploaded_by = $3
				OR EXISTS (
					SELECT 1 FROM document_permissions dp
					WHERE dp.document_id = d.id
					AND dp.grantee_type = 'USER'
					AND dp.grantee_id = $3
					AND dp.can_download = true
				)
			)
		)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, documentID, userID, companyID).Scan(&exists)
	return exists, err
}

func (r *DocumentRepository) LogAccess(ctx context.Context, log *models.DocumentAccessLog) error {
	query := `
		INSERT INTO document_access_logs (id, document_id, user_id, action, ip_address, user_agent)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query,
		log.DocumentID, log.UserID, log.Action, log.IPAddress, log.UserAgent,
	)
	return err
}

func (r *DocumentRepository) CreateTag(ctx context.Context, tag *models.DocumentTag) error {
	query := `
		INSERT INTO document_tags (id, company_id, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (company_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, tag.ID, tag.CompanyID, tag.Name).Scan(&tag.ID, &tag.CreatedAt)
}

func (r *DocumentRepository) ListTagsByCompany(ctx context.Context, companyID string) ([]models.DocumentTag, error) {
	query := `SELECT id, company_id, name, created_at FROM document_tags WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.DocumentTag
	for rows.Next() {
		var t models.DocumentTag
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *DocumentRepository) SetDocumentTags(ctx context.Context, documentID string, tagIDs []string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM document_tag_relations WHERE document_id=$1`, documentID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO document_tag_relations (id, document_id, tag_id) VALUES (gen_random_uuid(), $1, $2)`,
			documentID, tagID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DocumentRepository) ListTagsByDocumentID(ctx context.Context, documentID string) ([]models.DocumentTag, error) {
	query := `
		SELECT dt.id, dt.company_id, dt.name, dt.created_at
		FROM document_tags dt
		JOIN document_tag_relations dtr ON dtr.tag_id = dt.id
		WHERE dtr.document_id = $1
		ORDER BY dt.name`
	rows, err := r.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.DocumentTag
	for rows.Next() {
		var t models.DocumentTag
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *DocumentRepository) CreateShare(ctx context.Context, share *models.DocumentShare) error {
	query := `
		INSERT INTO document_shares (id, document_id, shared_by, shared_with_type, shared_with_id,
			can_read, can_download, can_share, expires_at, token, token_expires_at, max_uses)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query,
		share.ID, share.DocumentID, share.SharedBy, share.SharedWithType, share.SharedWithID,
		share.CanRead, share.CanDownload, share.CanShare, share.ExpiresAt,
		share.Token, share.TokenExpiresAt, share.MaxUses,
	).Scan(&share.CreatedAt)
}

func (r *DocumentRepository) GetShareByToken(ctx context.Context, token string) (*models.DocumentShare, error) {
	query := `
		SELECT id, document_id, shared_by, shared_with_type, shared_with_id,
			can_read, can_download, can_share, expires_at, token, token_expires_at,
			max_uses, use_count, is_active, created_at
		FROM document_shares WHERE token=$1 AND is_active=true`
	s := &models.DocumentShare{}
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&s.ID, &s.DocumentID, &s.SharedBy, &s.SharedWithType, &s.SharedWithID,
		&s.CanRead, &s.CanDownload, &s.CanShare, &s.ExpiresAt,
		&s.Token, &s.TokenExpiresAt, &s.MaxUses, &s.UseCount, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *DocumentRepository) IncrementShareUseCount(ctx context.Context, shareID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE document_shares SET use_count = use_count + 1 WHERE id=$1`, shareID,
	)
	return err
}

func (r *DocumentRepository) RevokeShare(ctx context.Context, shareID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE document_shares SET is_active = false WHERE id=$1`, shareID,
	)
	return err
}

func (r *DocumentRepository) ListExpiringDocuments(ctx context.Context, companyID string, withinDays int) ([]models.Document, error) {
	query := `
		SELECT
			d.id, d.company_id, d.category_id,
			(SELECT name FROM document_categories WHERE id = d.category_id),
			d.employee_id, d.department_id, d.uploaded_by,
			u.first_name || ' ' || u.last_name,
			d.title, d.description, d.original_filename, d.storage_key,
			d.mime_type, d.file_size, d.checksum, d.status, d.is_public,
			d.expires_at, d.created_at, d.updated_at
		FROM documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE d.company_id = $1
		AND d.status = 'ACTIVE'
		AND d.expires_at IS NOT NULL
		AND d.expires_at <= NOW() + INTERVAL '1 day' * $2
		ORDER BY d.expires_at ASC`

	rows, err := r.pool.Query(ctx, query, companyID, withinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.CategoryID, &d.CategoryName,
			&d.EmployeeID, &d.DepartmentID, &d.UploadedBy, &d.UploadedByName,
			&d.Title, &d.Description, &d.OriginalFilename, &d.StorageKey,
			&d.MimeType, &d.FileSize, &d.Checksum, &d.Status, &d.IsPublic,
			&d.ExpiresAt, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *DocumentRepository) GetDocumentStats(ctx context.Context, companyID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM documents WHERE company_id=$1 AND status != 'DELETED'`, companyID,
	).Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_documents"] = total

	var active int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM documents WHERE company_id=$1 AND status='ACTIVE'`, companyID,
	).Scan(&active)
	if err != nil {
		return nil, err
	}
	stats["active_documents"] = active

	var archived int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM documents WHERE company_id=$1 AND status='ARCHIVED'`, companyID,
	).Scan(&archived)
	if err != nil {
		return nil, err
	}
	stats["archived_documents"] = archived

	var totalSize int64
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(file_size), 0) FROM documents WHERE company_id=$1 AND status != 'DELETED'`, companyID,
	).Scan(&totalSize)
	if err != nil {
		return nil, err
	}
	stats["total_size_bytes"] = totalSize

	var expiringCount int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM documents WHERE company_id=$1 AND status='ACTIVE' AND expires_at IS NOT NULL AND expires_at <= NOW() + INTERVAL '30 days'`, companyID,
	).Scan(&expiringCount)
	if err != nil {
		return nil, err
	}
	stats["expiring_soon"] = expiringCount

	return stats, nil
}

var _ = pgx.ErrNoRows
var _ = time.Now
