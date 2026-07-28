package domain

import (
    "time"
)

type RequisitionStatus string

const (
    ReqStatusDraft           RequisitionStatus = "DRAFT"
    ReqStatusPendingApproval RequisitionStatus = "PENDING_APPROVAL"
    ReqStatusApproved        RequisitionStatus = "APPROVED"
    ReqStatusOpen            RequisitionStatus = "OPEN"
    ReqStatusClosed          RequisitionStatus = "CLOSED"
    ReqStatusCancelled       RequisitionStatus = "CANCELLED"
)

type RequisitionUrgency string

const (
    UrgencyLow    RequisitionUrgency = "LOW"
    UrgencyNormal RequisitionUrgency = "NORMAL"
    UrgencyHigh   RequisitionUrgency = "HIGH"
    UrgencyCritical RequisitionUrgency = "CRITICAL"
)

type Requisition struct {
    ID              string             `json:"id"`
    CompanyID       string             `json:"company_id"`
    PositionID      *string            `json:"position_id,omitempty"`
    DepartmentID    *string            `json:"department_id,omitempty"`
    RequestedBy     string             `json:"requested_by"`
    HiringManagerID *string            `json:"hiring_manager_id,omitempty"`
    Title           string             `json:"title"`
    Description     *string            `json:"description,omitempty"`
    Justification   *string            `json:"justification,omitempty"`
    Vacancies       int                `json:"vacancies"`
    EmploymentType  *string            `json:"employment_type,omitempty"`
    WorkMode        *string            `json:"work_mode,omitempty"`
    Location        *string            `json:"location,omitempty"`
    SalaryMin       *float64           `json:"salary_min,omitempty"`
    SalaryMax       *float64           `json:"salary_max,omitempty"`
    Currency        *string            `json:"currency,omitempty"`
    Urgency         RequisitionUrgency `json:"urgency"`
    Reason          *string            `json:"reason,omitempty"`
    Status          RequisitionStatus  `json:"status"`
    ApprovedAt      *time.Time         `json:"approved_at,omitempty"`
    OpenedAt        *time.Time         `json:"opened_at,omitempty"`
    ClosedAt        *time.Time         `json:"closed_at,omitempty"`
    ClosedReason    *string            `json:"closed_reason,omitempty"`
    CreatedAt       time.Time          `json:"created_at"`
    UpdatedAt       time.Time          `json:"updated_at"`
    Skills          []RequisitionSkill `json:"skills,omitempty"`
}

type RequisitionSkill struct {
    ID            string  `json:"id"`
    RequisitionID string  `json:"requisition_id"`
    Skill         string  `json:"skill"`
    Category      *string `json:"category,omitempty"`
    Required      bool    `json:"required"`
    MinYears      *int    `json:"min_years,omitempty"`
}

type RequisitionApproval struct {
    ID            string     `json:"id"`
    RequisitionID string     `json:"requisition_id"`
    ApproverID    string     `json:"approver_id"`
    StepOrder     int        `json:"step_order"`
    Status        string     `json:"status"`
    Comment       *string    `json:"comment,omitempty"`
    DecidedAt     *time.Time `json:"decided_at,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
}
