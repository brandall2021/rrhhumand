package domain

type DocStatus string

const (
	DocPending     DocStatus = "PENDING"
	DocUploaded    DocStatus = "UPLOADED"
	DocUnderReview DocStatus = "UNDER_REVIEW"
	DocApproved    DocStatus = "APPROVED"
	DocRejected    DocStatus = "REJECTED"
	DocExpired     DocStatus = "EXPIRED"
)

type SignatureStatus string

const (
	SignatureNotRequired SignatureStatus = "NOT_REQUIRED"
	SignaturePending     SignatureStatus = "PENDING"
	SignatureSent        SignatureStatus = "SENT"
	SignatureViewed      SignatureStatus = "VIEWED"
	SignatureSigned      SignatureStatus = "SIGNED"
	SignatureRejected    SignatureStatus = "REJECTED"
	SignatureExpired     SignatureStatus = "EXPIRED"
)

type OnboardingDocument struct {
	ID             string
	CompanyID      string
	OnboardingID   string
	EmployeeID     string
	DocumentType   string
	Name           string
	Required       bool
	Status         DocStatus
	StorageKey     *string
	MimeType       *string
	UploadedAt     *string
	VerifiedAt     *string
	VerifiedBy     *string
	ExpirationDate *string
	CreatedAt      string
	UpdatedAt      string
}

type OnboardingDocumentVersion struct {
	ID         string
	CompanyID  string
	DocumentID string
	Version    int
	FileName   string
	MimeType   string
	SizeBytes  int64
	StorageKey string
	Checksum   string
	UploadedBy string
	UploadedAt string
	Notes      *string
}

type SignatureProvider interface {
	CreateRequest(docID, employeeID, companyID string) error
	GetStatus(requestID string) (SignatureStatus, error)
	Download(requestID string) ([]byte, error)
}
