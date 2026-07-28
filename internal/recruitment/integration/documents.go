package integration

import (
	"context"
)

type DocumentResult struct {
	ID             string  `json:"id"`
	URL            string  `json:"url"`
	FileName       string  `json:"file_name"`
	MimeType       string  `json:"mime_type"`
	SizeBytes      int64   `json:"size_bytes"`
	StorageKey     string  `json:"storage_key"`
}

type DocumentAdapter struct{}

func NewDocumentAdapter() *DocumentAdapter {
	return &DocumentAdapter{}
}

func (a *DocumentAdapter) UploadDocument(ctx context.Context, candidateID string, file []byte, docType string) (*DocumentResult, error) {
	return &DocumentResult{
		ID:         "doc-" + candidateID,
		URL:        "https://storage.example.com/documents/" + candidateID,
		FileName:   "document.pdf",
		MimeType:   "application/pdf",
		SizeBytes:  int64(len(file)),
		StorageKey: "recruitment/" + candidateID + "/document.pdf",
	}, nil
}

func (a *DocumentAdapter) GetDocumentURL(ctx context.Context, docID string) string {
	return "https://storage.example.com/documents/" + docID
}

func (a *DocumentAdapter) DeleteDocument(ctx context.Context, docID string) error {
	return nil
}

func (a *DocumentAdapter) CopyDocument(ctx context.Context, sourceID, targetEntity string) (*DocumentResult, error) {
	return &DocumentResult{
		ID:    "doc-copy-" + sourceID,
		URL:   "https://storage.example.com/documents/copy-" + sourceID,
		FileName: "copy.pdf",
	}, nil
}
