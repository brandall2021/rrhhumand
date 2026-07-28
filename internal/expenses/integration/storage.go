package integration

import "fmt"

type StorageAdapter struct {
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

func NewStorageAdapter(endpoint, accessKey, secretKey string, useSSL bool) *StorageAdapter {
	return &StorageAdapter{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		useSSL:    useSSL,
	}
}

func (a *StorageAdapter) UploadFile(key string, content []byte, contentType string) error {
	_ = a.endpoint
	_ = a.accessKey
	_ = a.secretKey
	_ = a.useSSL
	_ = contentType
	return fmt.Errorf("UploadFile: not implemented (placeholder for key=%s, size=%d)", key, len(content))
}

func (a *StorageAdapter) GetFile(key string) ([]byte, error) {
	_ = a.endpoint
	_ = a.accessKey
	_ = a.secretKey
	_ = a.useSSL
	return nil, fmt.Errorf("GetFile: not implemented (placeholder for key=%s)", key)
}
