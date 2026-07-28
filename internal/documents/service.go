package documents

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type DocumentService struct {
	repo    *DocumentRepository
	storage *StorageService
}

func NewDocumentService(repo *DocumentRepository, storage *StorageService) *DocumentService {
	return &DocumentService{repo: repo, storage: storage}
}

func (s *DocumentService) UploadDocument(ctx context.Context, companyID, userID string, req *CreateDocumentRequest, file io.Reader, filename string, fileSize int64, mimeType string) (*models.Document, error) {
	if !s.storage.IsAllowedFileType(filename) {
		return nil, errors.New("file type not allowed")
	}

	maxSize := int64(25 * 1024 * 1024)
	if fileSize > maxSize {
		return nil, errors.New("file too large")
	}

	docID := uuid.New().String()
	storageKey := s.storage.GenerateStorageKey(companyID, *req.EmployeeID, docID, filename)

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(buf.Bytes())
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	if err := s.storage.Upload(ctx, storageKey, bytes.NewReader(buf.Bytes()), fileSize, mimeType); err != nil {
		return nil, err
	}

	doc := &models.Document{
		ID:               docID,
		CompanyID:        companyID,
		CategoryID:       req.CategoryID,
		EmployeeID:       req.EmployeeID,
		DepartmentID:     req.DepartmentID,
		UploadedBy:       userID,
		Title:            req.Title,
		Description:      req.Description,
		OriginalFilename: filename,
		StorageKey:       storageKey,
		MimeType:         mimeType,
		FileSize:         fileSize,
		Checksum:         &checksum,
		Status:           "ACTIVE",
		IsPublic:         req.IsPublic != nil && *req.IsPublic,
		ExpiresAt:        req.ExpiresAt,
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		_ = s.storage.Delete(ctx, storageKey)
		return nil, err
	}

	version := &models.DocumentVersion{
		ID:               uuid.New().String(),
		DocumentID:       docID,
		Version:          1,
		StorageKey:       storageKey,
		OriginalFilename: filename,
		MimeType:         mimeType,
		FileSize:         fileSize,
		Checksum:         &checksum,
		UploadedBy:       userID,
	}
	_ = s.repo.CreateVersion(ctx, version)

	if len(req.Tags) > 0 {
		s.SetDocumentTags(ctx, docID, companyID, req.Tags)
	}

	return doc, nil
}

func (s *DocumentService) GetDocumentByID(ctx context.Context, id, companyID string) (*models.Document, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}

	versions, _ := s.repo.ListVersions(ctx, doc.ID)
	doc.Versions = versions
	doc.CurrentVersion = len(versions)

	tags, _ := s.repo.ListTagsByDocumentID(ctx, doc.ID)
	doc.Tags = tags

	return doc, nil
}

func (s *DocumentService) ListDocuments(ctx context.Context, companyID string, filters DocumentFilters, params *models.PaginationParams) ([]models.Document, int64, error) {
	return s.repo.ListDocuments(ctx, companyID, filters, params.Offset, params.PerPage)
}

func (s *DocumentService) UpdateDocument(ctx context.Context, id, companyID string, req *UpdateDocumentRequest) (*models.Document, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if req.Title != nil {
		doc.Title = *req.Title
	}
	if req.Description != nil {
		doc.Description = req.Description
	}
	if req.CategoryID != nil {
		doc.CategoryID = req.CategoryID
	}
	if req.IsPublic != nil {
		doc.IsPublic = *req.IsPublic
	}
	if req.ExpiresAt != nil {
		doc.ExpiresAt = req.ExpiresAt
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return s.GetDocumentByID(ctx, id, companyID)
}

func (s *DocumentService) DeleteDocument(ctx context.Context, id, companyID string) error {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	if doc.Status == "DELETED" {
		return errors.New("document already deleted")
	}

	return s.repo.UpdateDocumentStatus(ctx, id, companyID, "DELETED")
}

func (s *DocumentService) RestoreDocument(ctx context.Context, id, companyID string) error {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	if doc.Status != "DELETED" {
		return errors.New("document not in trash")
	}

	return s.repo.UpdateDocumentStatus(ctx, id, companyID, "ACTIVE")
}

func (s *DocumentService) PermanentDelete(ctx context.Context, id, companyID string) error {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	_ = s.storage.Delete(ctx, doc.StorageKey)

	versions, _ := s.repo.ListVersions(ctx, doc.ID)
	for _, v := range versions {
		_ = s.storage.Delete(ctx, v.StorageKey)
	}

	return s.repo.DeleteDocument(ctx, id, companyID)
}

func (s *DocumentService) ArchiveDocument(ctx context.Context, id, companyID string) error {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	if doc.Status != "ACTIVE" {
		return errors.New("document not active")
	}

	return s.repo.UpdateDocumentStatus(ctx, id, companyID, "ARCHIVED")
}

func (s *DocumentService) DownloadDocument(ctx context.Context, id, companyID, userID string) (io.ReadCloser, *models.Document, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return nil, nil, errors.New("document not found")
	}

	hasAccess, err := s.repo.HasAccess(ctx, id, userID, companyID)
	if err != nil || !hasAccess {
		return nil, nil, errors.New("access denied")
	}

	reader, err := s.storage.Download(ctx, doc.StorageKey)
	if err != nil {
		return nil, nil, err
	}

	_ = s.repo.LogAccess(ctx, &models.DocumentAccessLog{
		DocumentID: id,
		UserID:     userID,
		Action:     "DOCUMENT_DOWNLOADED",
	})

	return reader, doc, nil
}

func (s *DocumentService) GetDownloadURL(ctx context.Context, id, companyID, userID string) (string, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id, companyID)
	if err != nil {
		return "", errors.New("document not found")
	}

	hasAccess, err := s.repo.HasAccess(ctx, id, userID, companyID)
	if err != nil || !hasAccess {
		return "", errors.New("access denied")
	}

	url, err := s.storage.GetPresignedURL(ctx, doc.StorageKey, 15*time.Minute)
	if err != nil {
		return "", err
	}

	_ = s.repo.LogAccess(ctx, &models.DocumentAccessLog{
		DocumentID: id,
		UserID:     userID,
		Action:     "DOCUMENT_VIEWED",
	})

	return url, nil
}

func (s *DocumentService) CreateVersion(ctx context.Context, documentID, companyID, userID string, file io.Reader, filename string, fileSize int64, mimeType string) (*models.DocumentVersion, error) {
	doc, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	latestVersion, err := s.repo.GetLatestVersionNumber(ctx, documentID)
	if err != nil {
		return nil, err
	}

	newVersion := latestVersion + 1
	newStorageKey := s.storage.GenerateStorageKey(companyID, *doc.EmployeeID, documentID, fmt.Sprintf("v%d_%s", newVersion, filename))

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(buf.Bytes())
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	if err := s.storage.Upload(ctx, newStorageKey, bytes.NewReader(buf.Bytes()), fileSize, mimeType); err != nil {
		return nil, err
	}

	version := &models.DocumentVersion{
		ID:               uuid.New().String(),
		DocumentID:       documentID,
		Version:          newVersion,
		StorageKey:       newStorageKey,
		OriginalFilename: filename,
		MimeType:         mimeType,
		FileSize:         fileSize,
		Checksum:         &checksum,
		UploadedBy:       userID,
	}

	if err := s.repo.CreateVersion(ctx, version); err != nil {
		_ = s.storage.Delete(ctx, newStorageKey)
		return nil, err
	}

	_ = s.repo.LogAccess(ctx, &models.DocumentAccessLog{
		DocumentID: documentID,
		UserID:     userID,
		Action:     "DOCUMENT_VERSION_CREATED",
	})

	return version, nil
}

func (s *DocumentService) ListVersions(ctx context.Context, documentID, companyID string) ([]models.DocumentVersion, error) {
	_, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}
	return s.repo.ListVersions(ctx, documentID)
}

func (s *DocumentService) SetPermissions(ctx context.Context, documentID, companyID string, req *SetDocumentPermissionsRequest) error {
	_, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	var perms []models.DocumentPermission
	for _, p := range req.Permissions {
		perms = append(perms, models.DocumentPermission{
			DocumentID:  documentID,
			GranteeType: p.GranteeType,
			GranteeID:   p.GranteeID,
			CanRead:     p.CanRead,
			CanDownload: p.CanDownload,
			CanShare:    p.CanShare,
			CanManage:   p.CanManage,
		})
	}

	return s.repo.SetPermissions(ctx, documentID, perms)
}

func (s *DocumentService) ListPermissions(ctx context.Context, documentID, companyID string) ([]models.DocumentPermission, error) {
	_, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}
	return s.repo.ListPermissions(ctx, documentID)
}

func (s *DocumentService) CreateTag(ctx context.Context, companyID, name string) (*models.DocumentTag, error) {
	tag := &models.DocumentTag{
		ID:        uuid.New().String(),
		CompanyID: companyID,
		Name:      name,
	}
	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *DocumentService) ListTags(ctx context.Context, companyID string) ([]models.DocumentTag, error) {
	return s.repo.ListTagsByCompany(ctx, companyID)
}

func (s *DocumentService) SetDocumentTags(ctx context.Context, documentID, companyID string, tagNames []string) error {
	_, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return errors.New("document not found")
	}

	var tagIDs []string
	for _, name := range tagNames {
		tag := &models.DocumentTag{
			ID:        uuid.New().String(),
			CompanyID: companyID,
			Name:      name,
		}
		_ = s.repo.CreateTag(ctx, tag)
		tagIDs = append(tagIDs, tag.ID)
	}

	return s.repo.SetDocumentTags(ctx, documentID, tagIDs)
}

func (s *DocumentService) CreateShare(ctx context.Context, documentID, companyID, userID string, req *ShareDocumentRequest) (*models.DocumentShare, error) {
	doc, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	share := &models.DocumentShare{
		ID:             uuid.New().String(),
		DocumentID:     documentID,
		SharedBy:       userID,
		SharedWithType: req.SharedWithType,
		SharedWithID:   req.SharedWithID,
		CanRead:        req.CanRead,
		CanDownload:    req.CanDownload,
		CanShare:       req.CanShare,
		ExpiresAt:      req.ExpiresAt,
		IsActive:       true,
	}

	if err := s.repo.CreateShare(ctx, share); err != nil {
		return nil, err
	}

	_ = s.repo.LogAccess(ctx, &models.DocumentAccessLog{
		DocumentID: documentID,
		UserID:     userID,
		Action:     "DOCUMENT_SHARED",
	})

	_ = doc
	return share, nil
}

func (s *DocumentService) CreateShareLink(ctx context.Context, documentID, companyID, userID string, req *CreateShareLinkRequest) (*models.DocumentShare, error) {
	doc, err := s.repo.GetDocumentByID(ctx, documentID, companyID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	token := uuid.New().String()
	tokenExpires := time.Now().Add(24 * time.Hour)
	if req.ExpiresAt != nil {
		tokenExpires = *req.ExpiresAt
	}

	share := &models.DocumentShare{
		ID:              uuid.New().String(),
		DocumentID:      documentID,
		SharedBy:        userID,
		SharedWithType:  "LINK",
		SharedWithID:    "PUBLIC",
		CanRead:         true,
		CanDownload:     true,
		Token:           &token,
		TokenExpiresAt:  &tokenExpires,
		MaxUses:         req.MaxUses,
		IsActive:        true,
	}

	if err := s.repo.CreateShare(ctx, share); err != nil {
		return nil, err
	}

	_ = doc
	return share, nil
}

func (s *DocumentService) AccessShareLink(ctx context.Context, token, userID string) (*models.Document, error) {
	share, err := s.repo.GetShareByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid share link")
	}

	if share.TokenExpiresAt != nil && time.Now().After(*share.TokenExpiresAt) {
		return nil, errors.New("share link expired")
	}

	if share.MaxUses != nil && share.UseCount >= *share.MaxUses {
		return nil, errors.New("share link max uses reached")
	}

	_ = s.repo.IncrementShareUseCount(ctx, share.ID)

	doc, err := s.repo.GetDocumentByID(ctx, share.DocumentID, "")
	if err != nil {
		return nil, errors.New("document not found")
	}

	_ = s.repo.LogAccess(ctx, &models.DocumentAccessLog{
		DocumentID: share.DocumentID,
		UserID:     userID,
		Action:     "DOCUMENT_VIEWED",
	})

	return doc, nil
}

func (s *DocumentService) RevokeShare(ctx context.Context, shareID string) error {
	return s.repo.RevokeShare(ctx, shareID)
}

func (s *DocumentService) ListExpiringDocuments(ctx context.Context, companyID string, withinDays int) ([]models.Document, error) {
	return s.repo.ListExpiringDocuments(ctx, companyID, withinDays)
}

func (s *DocumentService) GetDocumentStats(ctx context.Context, companyID string) (map[string]interface{}, error) {
	return s.repo.GetDocumentStats(ctx, companyID)
}
