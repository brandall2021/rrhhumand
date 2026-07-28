package documents

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rrhhumand/api/internal/config"
)

type StorageService struct {
	client *minio.Client
	bucket string
	region string
}

func NewStorageService(cfg config.MinIOConfig) (*StorageService, error) {
	var opts *minio.Options
	if cfg.UseSSL {
		opts = &minio.Options{
			Creds:  nil,
			Secure: true,
			Region: cfg.Region,
		}
	} else {
		opts = &minio.Options{
			Creds:  nil,
			Secure: false,
			Region: cfg.Region,
		}
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  nil,
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	_ = opts
	return &StorageService{client: client, bucket: cfg.Bucket, region: cfg.Region}, nil
}

func (s *StorageService) GenerateStorageKey(companyID, employeeID, documentID, filename string) string {
	parts := []string{"companies", companyID}
	if employeeID != "" {
		parts = append(parts, "employees", employeeID)
	}
	parts = append(parts, "documents", documentID, filename)
	return filepath.Join(parts...)
}

func (s *StorageService) Upload(ctx context.Context, storageKey string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, storageKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *StorageService) Download(ctx context.Context, storageKey string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, storageKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *StorageService) Delete(ctx context.Context, storageKey string) error {
	return s.client.RemoveObject(ctx, s.bucket, storageKey, minio.RemoveObjectOptions{})
}

func (s *StorageService) GetPresignedURL(ctx context.Context, storageKey string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, storageKey, expiry, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (s *StorageService) GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func (s *StorageService) CalculateChecksum(reader io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (s *StorageService) IsAllowedFileType(filename string) bool {
	allowed := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true,
		".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
		".txt": true, ".csv": true,
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return allowed[ext]
}

func (s *StorageService) GenerateDocumentID() string {
	return uuid.New().String()
}

func (s *StorageService) GetClient() *minio.Client {
	return s.client
}

func (s *StorageService) GetBucket() string {
	return s.bucket
}
