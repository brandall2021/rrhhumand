package domain

import "time"

type PositionStatus string

const (
    PosStatusDraft     PositionStatus = "DRAFT"
    PosStatusActive    PositionStatus = "ACTIVE"
    PosStatusFilled    PositionStatus = "FILLED"
    PosStatusCancelled PositionStatus = "CANCELLED"
)

type Position struct {
    ID              string         `json:"id"`
    CompanyID       string         `json:"company_id"`
    RequisitionID   *string        `json:"requisition_id,omitempty"`
    Title           string         `json:"title"`
    DepartmentID    *string        `json:"department_id,omitempty"`
    LocationID      *string        `json:"location_id,omitempty"`
    EmploymentType  *string        `json:"employment_type,omitempty"`
    WorkMode        *string        `json:"work_mode,omitempty"`
    Description     *string        `json:"description,omitempty"`
    Requirements    *string        `json:"requirements,omitempty"`
    Responsibilities *string       `json:"responsibilities,omitempty"`
    Benefits        *string        `json:"benefits,omitempty"`
    SalaryMin       *float64       `json:"salary_min,omitempty"`
    SalaryMax       *float64       `json:"salary_max,omitempty"`
    Currency        *string        `json:"currency,omitempty"`
    Vacancies       int            `json:"vacancies"`
    VacanciesFilled int            `json:"vacancies_filled"`
    Status          PositionStatus `json:"status"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    Skills          []PositionSkill `json:"skills,omitempty"`
}

type PositionSkill struct {
    ID         string   `json:"id"`
    PositionID string   `json:"position_id"`
    Skill      string   `json:"skill"`
    Category   *string  `json:"category,omitempty"`
    Required   bool     `json:"required"`
    MinYears   *int     `json:"min_years,omitempty"`
    Weight     *float64 `json:"weight,omitempty"`
}
