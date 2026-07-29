package integration

import "context"

type EmployeeInfo struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
	PositionID   *string `json:"position_id,omitempty"`
	ManagerID    *string `json:"manager_id,omitempty"`
	HireDate     string  `json:"hire_date"`
	Status       string  `json:"status"`
}

type CreateEmployeeRequest struct {
	EmployeeNumber string `json:"employee_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DepartmentID   string `json:"department_id"`
	PositionID     string `json:"position_id"`
	ManagerID      string `json:"manager_id"`
	HireDate       string `json:"hire_date"`
}

type EmployeeService interface {
	Create(ctx context.Context, companyID string, req *CreateEmployeeRequest) (*EmployeeInfo, error)
	GetByID(ctx context.Context, id, companyID string) (*EmployeeInfo, error)
	Update(ctx context.Context, id, companyID string, req interface{}) (*EmployeeInfo, error)
	ExistsByEmail(ctx context.Context, companyID, email string) (bool, error)
}

type DocumentService interface {
	Upload(ctx context.Context, companyID, employeeID, docType, fileName string, content []byte) (string, error)
	GetDownloadURL(ctx context.Context, storageKey string) (string, error)
	Delete(ctx context.Context, storageKey string) error
}

type AssetService interface {
	GetEmployeeAssets(ctx context.Context, companyID, employeeID string) ([]AssetInfo, error)
	AssignAsset(ctx context.Context, companyID, employeeID, assetType, description string) error
	ReturnAsset(ctx context.Context, companyID, assetID string, condition string) error
}

type AssetInfo struct {
	ID          string `json:"id"`
	AssetType   string `json:"asset_type"`
	Description string `json:"description"`
	SerialNumber string `json:"serial_number"`
	Status      string `json:"status"`
}

type TrainingService interface {
	AssignCourse(ctx context.Context, companyID, employeeID, courseName string, mandatory bool, dueDate string) error
	GetCourseStatus(ctx context.Context, employeeID, courseName string) (string, error)
}

type NotificationService interface {
	Send(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) error
	SendToRole(ctx context.Context, companyID, role, title, body, notifType, refType, refID string) error
}

type PayrollService interface {
	StartFinalSettlement(ctx context.Context, companyID, employeeID string, terminationType string, lastWorkingDate string) error
	GetSettlementStatus(ctx context.Context, employeeID string) (string, error)
}

type AccessProvisioningService interface {
	CreateAccount(ctx context.Context, employeeID, systemName, accessType string) error
	DisableAccount(ctx context.Context, employeeID, systemName string) error
	RevokeAccess(ctx context.Context, employeeID, systemName string) error
}

type SignatureService interface {
	SendDocument(ctx context.Context, documentID, employeeID, documentName string) error
	GetStatus(ctx context.Context, requestID string) (string, error)
}

type CalendarService interface {
	CreateEvent(ctx context.Context, companyID, employeeID, title, description, startDate string) error
}
