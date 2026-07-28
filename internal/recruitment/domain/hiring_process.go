package domain

import "time"

type HiringProcessStatus string

const (
    HireStatusPending     HiringProcessStatus = "PENDING"
    HireStatusInProgress  HiringProcessStatus = "IN_PROGRESS"
    HireStatusCompleted   HiringProcessStatus = "COMPLETED"
    HireStatusCancelled   HiringProcessStatus = "CANCELLED"
)

type HiringProcess struct {
    ID                      string              `json:"id"`
    CompanyID               string              `json:"company_id"`
    OfferID                 *string             `json:"offer_id,omitempty"`
    ApplicationID           string              `json:"application_id"`
    CandidateID             string              `json:"candidate_id"`
    EmployeeID              *string             `json:"employee_id,omitempty"`
    Status                  HiringProcessStatus `json:"status"`
    BackgroundCheckStatus   string              `json:"background_check_status"`
    BackgroundCheckResult   *string             `json:"background_check_result,omitempty"`
    MedicalCheckStatus      string              `json:"medical_check_status"`
    MedicalCheckResult      *string             `json:"medical_check_result,omitempty"`
    DocVerificationStatus   string              `json:"document_verification_status"`
    StartDate               *time.Time          `json:"start_date,omitempty"`
    OnboardingStatus        string              `json:"onboarding_status"`
    OnboardingID            *string             `json:"onboarding_id,omitempty"`
    Notes                   *string             `json:"notes,omitempty"`
    CreatedBy               *string             `json:"created_by,omitempty"`
    CreatedAt               time.Time           `json:"created_at"`
    UpdatedAt               time.Time           `json:"updated_at"`
    Tasks                   []HiringProcessTask `json:"tasks,omitempty"`
}

type HiringProcessTask struct {
    ID          string     `json:"id"`
    ProcessID   string     `json:"process_id"`
    TaskType    string     `json:"task_type"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    AssignedTo  *string    `json:"assigned_to,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    Status      string     `json:"status"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}

type HiringProcessDocument struct {
    ID           string     `json:"id"`
    ProcessID    string     `json:"process_id"`
    DocumentType string     `json:"document_type"`
    FileName     string     `json:"file_name"`
    StorageKey   string     `json:"storage_key"`
    Verified     bool       `json:"verified"`
    VerifiedBy   *string    `json:"verified_by,omitempty"`
    VerifiedAt   *time.Time `json:"verified_at,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
}
