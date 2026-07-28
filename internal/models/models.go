package models

import (
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"strconv"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	Active       bool       `json:"active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CompanyID *string   `json:"company_id,omitempty"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

type Company struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	LogoURL   *string   `json:"logo_url,omitempty"`
	Plan      string    `json:"plan"`
	Settings  []byte    `json:"settings,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Branch struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Name      string    `json:"name"`
	Code      *string   `json:"code,omitempty"`
	Address   *string   `json:"address,omitempty"`
	City      *string   `json:"city,omitempty"`
	State     *string   `json:"state,omitempty"`
	Country   string    `json:"country"`
	Phone     *string   `json:"phone,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Timezone  string    `json:"timezone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Department struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Code        *string   `json:"code,omitempty"`
	Description *string   `json:"description,omitempty"`
	BranchID    *string   `json:"branch_id,omitempty"`
	ParentID    *string   `json:"parent_id,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Position struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Name         string    `json:"name"`
	Code         *string   `json:"code,omitempty"`
	Description  *string   `json:"description,omitempty"`
	DepartmentID *string   `json:"department_id,omitempty"`
	Level        int       `json:"level"`
	MinSalary    *float64  `json:"min_salary,omitempty"`
	MaxSalary    *float64  `json:"max_salary,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Employee struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeNumber  string     `json:"employee_number"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DNI             *string    `json:"dni,omitempty"`
	Email           *string    `json:"email,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	BirthDate       *string    `json:"birth_date,omitempty"`
	PhotoURL        *string    `json:"photo_url,omitempty"`
	BranchID        *string    `json:"branch_id,omitempty"`
	DepartmentID    *string    `json:"department_id,omitempty"`
	PositionID      *string    `json:"position_id,omitempty"`
	ManagerID       *string    `json:"manager_id,omitempty"`
	HireDate        string     `json:"hire_date"`
	TerminationDate *string    `json:"termination_date,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	BranchName     *string `json:"branch_name,omitempty"`
	DepartmentName *string `json:"department_name,omitempty"`
	PositionName   *string `json:"position_name,omitempty"`
	ManagerName    *string `json:"manager_name,omitempty"`
}

type EmployeeContact struct {
	ID           string `json:"id"`
	EmployeeID   string `json:"employee_id"`
	ContactType  string `json:"contact_type"`
	ContactValue string `json:"contact_value"`
	IsPrimary    bool   `json:"is_primary"`
}

type EmployeeAddress struct {
	ID            string  `json:"id"`
	EmployeeID    string  `json:"employee_id"`
	AddressType   string  `json:"address_type"`
	Street        *string `json:"street,omitempty"`
	StreetNumber  *string `json:"street_number,omitempty"`
	Apartment     *string `json:"apartment,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	Country       string  `json:"country"`
	PostalCode    *string `json:"postal_code,omitempty"`
	IsPrimary     bool    `json:"is_primary"`
}

type EmployeeEmergencyContact struct {
	ID           string  `json:"id"`
	EmployeeID   string  `json:"employee_id"`
	Name         string  `json:"name"`
	Relationship *string `json:"relationship,omitempty"`
	Phone        string  `json:"phone"`
	AltPhone     *string `json:"alt_phone,omitempty"`
	IsPrimary    bool    `json:"is_primary"`
}

type EmployeeHistory struct {
	ID          string  `json:"id"`
	EmployeeID  string  `json:"employee_id"`
	EventType   string  `json:"event_type"`
	OldValue    *string `json:"old_value,omitempty"`
	NewValue    *string `json:"new_value,omitempty"`
	Description *string `json:"description,omitempty"`
	PerformedBy *string `json:"performed_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationParams struct {
	Page    int
	PerPage int
	Offset  int
}

func NewPaginationParams(c *gin.Context) *PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return &PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

func (p *PaginationParams) ToMeta(total int64) map[string]interface{} {
	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))
	return map[string]interface{}{
		"page":        p.Page,
		"per_page":    p.PerPage,
		"total":       total,
		"total_pages": totalPages,
	}
}

type Post struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	AuthorID     string    `json:"author_id"`
	AuthorName   string    `json:"author_name,omitempty"`
	AuthorPhoto  *string   `json:"author_photo,omitempty"`
	Content      string    `json:"content"`
	Visibility   string    `json:"visibility"`
	Pinned       bool      `json:"pinned"`
	Media        []PostMedia `json:"media,omitempty"`
	Comments     []Comment `json:"comments,omitempty"`
	Reactions    []Reaction `json:"reactions,omitempty"`
	ReactionCounts map[string]int `json:"reaction_counts,omitempty"`
	UserReaction *string   `json:"user_reaction,omitempty"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostMedia struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	MediaType string    `json:"media_type"`
	URL       string    `json:"url"`
	Filename  *string   `json:"filename,omitempty"`
	FileSize  *int      `json:"file_size,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID         string    `json:"id"`
	PostID     string    `json:"post_id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name,omitempty"`
	AuthorPhoto *string  `json:"author_photo,omitempty"`
	ParentID   *string   `json:"parent_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Reaction struct {
	ID         string    `json:"id"`
	PostID     string    `json:"post_id"`
	EmployeeID string    `json:"employee_id"`
	ReactionType string  `json:"reaction_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type Mention struct {
	ID                 string    `json:"id"`
	PostID             *string   `json:"post_id,omitempty"`
	CommentID          *string   `json:"comment_id,omitempty"`
	MentionedEmployeeID string   `json:"mentioned_employee_id"`
	CreatedAt          time.Time `json:"created_at"`
}

type Survey struct {
	ID                 string     `json:"id"`
	CompanyID          string     `json:"company_id"`
	Title              string     `json:"title"`
	Description        *string    `json:"description,omitempty"`
	Type               string     `json:"type"`
	Status             string     `json:"status"`
	Anonymous          bool       `json:"anonymous"`
	MultipleResponses  bool       `json:"multiple_responses"`
	StartsAt           *time.Time `json:"starts_at,omitempty"`
	EndsAt             *time.Time `json:"ends_at,omitempty"`
	CreatedBy          string     `json:"created_by"`
	CreatedByName      string     `json:"created_by_name,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Questions          []SurveyQuestion `json:"questions,omitempty"`
	Targets            []SurveyTarget   `json:"targets,omitempty"`
	ResponseCount      int              `json:"response_count,omitempty"`
	TargetCount        int              `json:"target_count,omitempty"`
	ParticipationRate  float64          `json:"participation_rate,omitempty"`
}

type SurveyQuestion struct {
	ID         string           `json:"id"`
	SurveyID   string           `json:"survey_id"`
	Question   string           `json:"question"`
	Type       string           `json:"type"`
	Position   int              `json:"position"`
	Required   bool             `json:"required"`
	CreatedAt  time.Time        `json:"created_at"`
	Options    []SurveyOption   `json:"options,omitempty"`
}

type SurveyOption struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	OptionText string    `json:"option_text"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

type SurveyTarget struct {
	ID         string    `json:"id"`
	SurveyID   string    `json:"survey_id"`
	TargetType string    `json:"target_type"`
	TargetID   *string   `json:"target_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type SurveyResponse struct {
	ID         string    `json:"id"`
	SurveyID   string    `json:"survey_id"`
	EmployeeID string    `json:"employee_id,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
	Answers    []SurveyAnswer `json:"answers,omitempty"`
}

type SurveyAnswer struct {
	ID         string    `json:"id"`
	ResponseID string    `json:"response_id"`
	QuestionID string    `json:"question_id"`
	TextValue  *string   `json:"text_value,omitempty"`
	NumberValue *float64 `json:"number_value,omitempty"`
	OptionID   *string   `json:"option_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Options    []SurveyAnswerOption `json:"selected_options,omitempty"`
}

type SurveyAnswerOption struct {
	ID         string    `json:"id"`
	AnswerID   string    `json:"answer_id"`
	OptionID   string    `json:"option_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type DocumentCategory struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	ParentID    *string   `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Children    []DocumentCategory `json:"children,omitempty"`
}

type Document struct {
	ID               string     `json:"id"`
	CompanyID        string     `json:"company_id"`
	CategoryID       *string    `json:"category_id,omitempty"`
	CategoryName     *string    `json:"category_name,omitempty"`
	EmployeeID       *string    `json:"employee_id,omitempty"`
	DepartmentID     *string    `json:"department_id,omitempty"`
	UploadedBy       string     `json:"uploaded_by"`
	UploadedByName   string     `json:"uploaded_by_name,omitempty"`
	Title            string     `json:"title"`
	Description      *string    `json:"description,omitempty"`
	OriginalFilename string     `json:"original_filename"`
	StorageKey       string     `json:"-"`
	MimeType         string     `json:"mime_type"`
	FileSize         int64      `json:"file_size"`
	Checksum         *string    `json:"checksum,omitempty"`
	Status           string     `json:"status"`
	IsPublic         bool       `json:"is_public"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Versions         []DocumentVersion  `json:"versions,omitempty"`
	Tags             []DocumentTag      `json:"tags,omitempty"`
	CurrentVersion   int               `json:"current_version,omitempty"`
}

type DocumentVersion struct {
	ID               string    `json:"id"`
	DocumentID       string    `json:"document_id"`
	Version          int       `json:"version"`
	StorageKey       string    `json:"-"`
	OriginalFilename string    `json:"original_filename"`
	MimeType         string    `json:"mime_type"`
	FileSize         int64     `json:"file_size"`
	Checksum         *string   `json:"checksum,omitempty"`
	UploadedBy       string    `json:"uploaded_by"`
	CreatedAt        time.Time `json:"created_at"`
}

type DocumentPermission struct {
	ID           string    `json:"id"`
	DocumentID   string    `json:"document_id"`
	GranteeType  string    `json:"grantee_type"`
	GranteeID    string    `json:"grantee_id"`
	CanRead      bool      `json:"can_read"`
	CanDownload  bool      `json:"can_download"`
	CanShare     bool      `json:"can_share"`
	CanManage    bool      `json:"can_manage"`
	CreatedAt    time.Time `json:"created_at"`
}

type DocumentTag struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type DocumentAccessLog struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	UserAgent  *string   `json:"user_agent,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type DocumentShare struct {
	ID              string     `json:"id"`
	DocumentID      string     `json:"document_id"`
	SharedBy        string     `json:"shared_by"`
	SharedWithType  string     `json:"shared_with_type"`
	SharedWithID    string     `json:"shared_with_id"`
	CanRead         bool       `json:"can_read"`
	CanDownload     bool       `json:"can_download"`
	CanShare        bool       `json:"can_share"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	Token           *string    `json:"token,omitempty"`
	TokenExpiresAt  *time.Time `json:"token_expires_at,omitempty"`
	MaxUses         *int       `json:"max_uses,omitempty"`
	UseCount        int        `json:"use_count"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
}
