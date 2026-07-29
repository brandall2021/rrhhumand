package integration

import (
	"context"
	"fmt"
	"log"
)

type DocumentAdapter struct{}

func NewDocumentAdapter() *DocumentAdapter {
	return &DocumentAdapter{}
}

func (a *DocumentAdapter) Upload(ctx context.Context, companyID, employeeID, docType, fileName string, content []byte) (string, error) {
	log.Printf("[DocumentAdapter] Upload company=%s employee=%s type=%s file=%s", companyID, employeeID, docType, fileName)
	return fmt.Sprintf("storage/%s/%s/%s", companyID, employeeID, fileName), nil
}

func (a *DocumentAdapter) GetDownloadURL(ctx context.Context, storageKey string) (string, error) {
	log.Printf("[DocumentAdapter] GetDownloadURL key=%s", storageKey)
	return fmt.Sprintf("https://storage.example.com/%s", storageKey), nil
}

func (a *DocumentAdapter) Delete(ctx context.Context, storageKey string) error {
	log.Printf("[DocumentAdapter] Delete key=%s", storageKey)
	return nil
}
