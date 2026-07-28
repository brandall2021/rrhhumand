package domain

import (
    "time"
)

type CandidateStatus string

const (
    CandStatusActive      CandidateStatus = "ACTIVE"
    CandStatusInactive    CandidateStatus = "INACTIVE"
    CandStatusBlacklisted CandidateStatus = "BLACKLISTED"
    CandStatusHired       CandidateStatus = "HIRED"
)

type Candidate struct {
    ID                   string          `json:"id"`
    CompanyID            string          `json:"company_id"`
    FirstName            string          `json:"first_name"`
    LastName             string          `json:"last_name"`
    Email                string          `json:"email"`
    Phone                *string         `json:"phone,omitempty"`
    PhoneCountryCode     *string         `json:"phone_country_code,omitempty"`
    DocumentType         *string         `json:"document_type,omitempty"`
    DocumentNumber       *string         `json:"document_number,omitempty"`
    BirthDate            *time.Time      `json:"birth_date,omitempty"`
    Location             *string         `json:"location,omitempty"`
    Nationality          *string         `json:"nationality,omitempty"`
    Gender               *string         `json:"gender,omitempty"`
    LinkedInURL          *string         `json:"linkedin_url,omitempty"`
    PortfolioURL         *string         `json:"portfolio_url,omitempty"`
    GithubURL            *string         `json:"github_url,omitempty"`
    PersonalWebsite      *string         `json:"personal_website,omitempty"`
    CurrentCompany       *string         `json:"current_company,omitempty"`
    CurrentPosition      *string         `json:"current_position,omitempty"`
    NoticePeriod         *int            `json:"notice_period,omitempty"`
    SalaryExpectMin      *float64        `json:"salary_expectation_min,omitempty"`
    SalaryExpectMax      *float64        `json:"salary_expectation_max,omitempty"`
    SalaryCurrency       *string         `json:"salary_currency,omitempty"`
    Availability         *string         `json:"availability,omitempty"`
    Source               *string         `json:"source,omitempty"`
    SourceDetail         *string         `json:"source_detail,omitempty"`
    IsReferral           bool            `json:"is_employee_referral"`
    ReferrerEmployeeID   *string         `json:"referrer_employee_id,omitempty"`
    Status               CandidateStatus `json:"status"`
    Blacklisted          bool            `json:"blacklisted"`
    BlacklistReason      *string         `json:"blacklist_reason,omitempty"`
    Tags                 []string        `json:"tags,omitempty"`
    Notes                *string         `json:"notes,omitempty"`
    CreatedAt            time.Time       `json:"created_at"`
    UpdatedAt            time.Time       `json:"updated_at"`
    Education            []CandidateEducation   `json:"education,omitempty"`
    Experience           []CandidateExperience  `json:"experience,omitempty"`
    Skills               []CandidateSkill       `json:"skills,omitempty"`
    Certifications       []CandidateCertification `json:"certifications,omitempty"`
    Languages            []CandidateLanguage    `json:"languages,omitempty"`
    Documents            []CandidateDocument    `json:"documents,omitempty"`
}

type CandidateEducation struct {
    ID          string     `json:"id"`
    CandidateID string     `json:"candidate_id"`
    Institution string     `json:"institution"`
    Degree      *string    `json:"degree,omitempty"`
    FieldOfStudy *string   `json:"field_of_study,omitempty"`
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    IsCurrent   bool       `json:"is_current"`
    Grade       *string    `json:"grade,omitempty"`
    Description *string    `json:"description,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}

type CandidateExperience struct {
    ID          string     `json:"id"`
    CandidateID string     `json:"candidate_id"`
    Company     string     `json:"company"`
    Position    string     `json:"position"`
    Location    *string    `json:"location,omitempty"`
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    IsCurrent   bool       `json:"is_current"`
    Description *string    `json:"description,omitempty"`
    Achievements []string  `json:"achievements,omitempty"`
    Industry    *string    `json:"industry,omitempty"`
    Salary      *float64   `json:"salary,omitempty"`
    SalaryCurrency *string `json:"salary_currency,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}

type CandidateSkill struct {
    ID             string  `json:"id"`
    CandidateID    string  `json:"candidate_id"`
    Skill          string  `json:"skill"`
    Category       *string `json:"category,omitempty"`
    Proficiency    string  `json:"proficiency"`
    YearsExp       *float64 `json:"years_experience,omitempty"`
    IsPrimary      bool    `json:"is_primary"`
    CreatedAt      time.Time `json:"created_at"`
}

type CandidateCertification struct {
    ID           string     `json:"id"`
    CandidateID  string     `json:"candidate_id"`
    Name         string     `json:"name"`
    Issuer       *string    `json:"issuer,omitempty"`
    IssueDate    *time.Time `json:"issue_date,omitempty"`
    ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
    CredentialID *string    `json:"credential_id,omitempty"`
    CredentialURL *string   `json:"credential_url,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
}

type CandidateLanguage struct {
    ID          string `json:"id"`
    CandidateID string `json:"candidate_id"`
    Language    string `json:"language"`
    Proficiency string `json:"proficiency"`
    IsNative    bool   `json:"is_native"`
    CreatedAt   time.Time `json:"created_at"`
}

type CandidateDocument struct {
    ID              string     `json:"id"`
    CandidateID     string     `json:"candidate_id"`
    CompanyID       string     `json:"company_id"`
    DocumentType    string     `json:"document_type"`
    FileName        string     `json:"file_name"`
    MimeType        *string    `json:"mime_type,omitempty"`
    SizeBytes       *int64     `json:"size_bytes,omitempty"`
    StorageProvider *string    `json:"storage_provider,omitempty"`
    StorageKey      *string    `json:"storage_key,omitempty"`
    ParsedData      *string    `json:"parsed_data,omitempty"`
    ParsedAt        *time.Time `json:"parsed_at,omitempty"`
    CreatedAt       time.Time  `json:"created_at"`
}

type CandidateParsedData struct {
    ID               string     `json:"id"`
    CandidateID      string     `json:"candidate_id"`
    RawText          *string    `json:"raw_text,omitempty"`
    StructuredData   *string    `json:"structured_data,omitempty"`
    SkillsFound      []string   `json:"skills_found,omitempty"`
    EducationFound   *string    `json:"education_found,omitempty"`
    ExperienceFound  *string    `json:"experience_found,omitempty"`
    LanguagesFound   *string    `json:"languages_found,omitempty"`
    CertsFound       *string    `json:"certifications_found,omitempty"`
    Summary          *string    `json:"summary,omitempty"`
    Score            *float64   `json:"score,omitempty"`
    ParserVersion    *string    `json:"parser_version,omitempty"`
    ParsedAt         time.Time  `json:"parsed_at"`
    CreatedAt        time.Time  `json:"created_at"`
}
