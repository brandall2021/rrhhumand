package integration

import (
	"context"
	"fmt"
	"log"
)

type SignatureAdapter struct{}

func NewSignatureAdapter() *SignatureAdapter {
	return &SignatureAdapter{}
}

func (a *SignatureAdapter) SendDocument(ctx context.Context, documentID, employeeID, documentName string) error {
	log.Printf("[SignatureAdapter] SendDocument document=%s employee=%s name=%s", documentID, employeeID, documentName)
	return nil
}

func (a *SignatureAdapter) GetStatus(ctx context.Context, requestID string) (string, error) {
	log.Printf("[SignatureAdapter] GetStatus request=%s", requestID)
	return fmt.Sprintf("https://sign.example.com/doc/%s", requestID), nil
}
