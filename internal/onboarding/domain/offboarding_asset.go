package domain

type OffboardingAssetStatus string

const (
	AssetPendingReturn OffboardingAssetStatus = "PENDING_RETURN"
	AssetReturned      OffboardingAssetStatus = "RETURNED"
	AssetLost          OffboardingAssetStatus = "LOST"
	AssetDamaged       OffboardingAssetStatus = "DAMAGED"
)

type OffboardingAsset struct {
	ID                 string
	CompanyID          string
	OffboardingID      string
	EmployeeID         string
	AssetType          string
	Description        *string
	SerialNumber       *string
	ConditionOnDelivery *string
	ConditionOnReturn  *string
	Status             OffboardingAssetStatus
	ReturnedAt         *string
	ReturnedTo         *string
	Notes              *string
	CreatedAt          string
	UpdatedAt          string
}

type AccessRevocation struct {
	ID            string
	CompanyID     string
	EmployeeID    string
	OffboardingID *string
	SystemName    string
	AccessType    string
	RequestedAt   string
	RevokedAt     *string
	Status        string
	PerformedBy   *string
	ErrorMessage  *string
	CreatedAt     string
	UpdatedAt     string
}

type AccessRevocationStatus string

const (
	RevokePending    AccessRevocationStatus = "PENDING"
	RevokeInProgress AccessRevocationStatus = "IN_PROGRESS"
	RevokeRevoked    AccessRevocationStatus = "REVOKED"
	RevokeFailed     AccessRevocationStatus = "FAILED"
)

type AccessProvisioningService interface {
	CreateAccount(employeeID, systemName, accessType string) error
	DisableAccount(employeeID, systemName string) error
	EnableAccount(employeeID, systemName string) error
	ResetAccess(employeeID, systemName string) error
}

type EmployeeHandover struct {
	ID           string
	CompanyID    string
	EmployeeID   string
	OffboardingID *string
	HandoverTo   string
	Description  *string
	Projects     *string
	PendingTasks *string
	Documents    *string
	Status       string
	CompletedAt  *string
	CreatedAt    string
	UpdatedAt    string
}

type HandoverStatus string

const (
	HandoverPending   HandoverStatus = "PENDING"
	HandoverCompleted HandoverStatus = "COMPLETED"
)
