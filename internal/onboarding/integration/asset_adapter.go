package integration

import (
	"context"
	"fmt"
	"log"
)

type AssetAdapter struct{}

func NewAssetAdapter() *AssetAdapter {
	return &AssetAdapter{}
}

func (a *AssetAdapter) GetEmployeeAssets(ctx context.Context, companyID, employeeID string) ([]AssetInfo, error) {
	log.Printf("[AssetAdapter] GetEmployeeAssets company=%s employee=%s", companyID, employeeID)
	return nil, nil
}

func (a *AssetAdapter) AssignAsset(ctx context.Context, companyID, employeeID, assetType, description string) error {
	log.Printf("[AssetAdapter] AssignAsset company=%s employee=%s type=%s", companyID, employeeID, assetType)
	return nil
}

func (a *AssetAdapter) ReturnAsset(ctx context.Context, companyID, assetID string, condition string) error {
	log.Printf("[AssetAdapter] ReturnAsset company=%s asset=%s condition=%s", companyID, assetID, condition)
	return nil
}

func (a *AssetAdapter) GenerateAssetRequest(ctx context.Context, companyID, employeeID, assetType, description string) error {
	log.Printf("[AssetAdapter] GenerateAssetRequest company=%s employee=%s type=%s", companyID, employeeID, assetType)
	return nil
}

func (a *AssetAdapter) CancelAssetRequest(ctx context.Context, companyID, requestID string) error {
	log.Printf("[AssetAdapter] CancelAssetRequest company=%s request=%s", companyID, requestID)
	return nil
}

func (a *AssetAdapter) GetAssetRequestStatus(ctx context.Context, requestID string) (string, error) {
	return "PENDING", nil
}

func (a *AssetAdapter) GetAssetBySerial(ctx context.Context, companyID, serial string) (*AssetInfo, error) {
	return &AssetInfo{
		ID:          fmt.Sprintf("asset-%s", serial),
		AssetType:   "NOTEBOOK",
		Description: "Standard equipment",
		SerialNumber: serial,
		Status:      "AVAILABLE",
	}, nil
}
