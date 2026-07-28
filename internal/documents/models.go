package documents

import (
	"time"
)

type CreateDocumentRequest struct {
	Title            string  `json:"title" validate:"required"`
	Description      *string `json:"description,omitempty"`
	CategoryID       *string `json:"category_id,omitempty"`
	EmployeeID       *string `json:"employee_id,omitempty"`
	DepartmentID     *string `json:"department_id,omitempty"`
	IsPublic         *bool   `json:"is_public,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type UpdateDocumentRequest struct {
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	CategoryID       *string    `json:"category_id,omitempty"`
	IsPublic         *bool      `json:"is_public,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
}

type CreateVersionRequest struct {
}

type DocumentFilters struct {
	Status       string
	CategoryID   string
	EmployeeID   string
	DepartmentID string
	MimeType     string
	Search       string
	CreatedFrom  string
	CreatedTo    string
	Tag          string
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type SetDocumentPermissionsRequest struct {
	Permissions []DocumentPermissionItem `json:"permissions" validate:"required"`
}

type DocumentPermissionItem struct {
	GranteeType string `json:"grantee_type" validate:"required"`
	GranteeID   string `json:"grantee_id" validate:"required"`
	CanRead     bool   `json:"can_read"`
	CanDownload bool   `json:"can_download"`
	CanShare    bool   `json:"can_share"`
	CanManage   bool   `json:"can_manage"`
}

type ShareDocumentRequest struct {
	SharedWithType string     `json:"shared_with_type" validate:"required"`
	SharedWithID   string     `json:"shared_with_id" validate:"required"`
	CanRead        bool       `json:"can_read"`
	CanDownload    bool       `json:"can_download"`
	CanShare       bool       `json:"can_share"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

type CreateShareLinkRequest struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	MaxUses   *int       `json:"max_uses,omitempty"`
}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}
