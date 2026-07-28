package models

import "time"

type LeaveType struct {
	ID                string    `json:"id"`
	CompanyID         string    `json:"company_id"`
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	Description       *string   `json:"description,omitempty"`
	Category          string    `json:"category"`
	RequiresApproval  bool      `json:"requires_approval"`
	RequiresDocument  bool      `json:"requires_document"`
	AffectsBalance    bool      `json:"affects_balance"`
	IsPaid            bool      `json:"is_paid"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type LeavePolicy struct {
	ID                      string    `json:"id"`
	CompanyID               string    `json:"company_id"`
	LeaveTypeID             string    `json:"leave_type_id"`
	LeaveTypeName           string    `json:"leave_type_name,omitempty"`
	Name                    string    `json:"name"`
	DaysPerYear             *float64  `json:"days_per_year,omitempty"`
	MinimumDaysBeforeRequest int      `json:"minimum_days_before_request"`
	MaximumDaysPerRequest   *float64  `json:"maximum_days_per_request,omitempty"`
	MaximumAccumulatedDays  *float64  `json:"maximum_accumulated_days,omitempty"`
	AllowNegativeBalance    bool      `json:"allow_negative_balance"`
	UseBusinessDays         bool      `json:"use_business_days"`
	RequiresManagerApproval bool      `json:"requires_manager_approval"`
	RequiresHRApproval      bool      `json:"requires_hr_approval"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type Holiday struct {
	ID          string    `json:"id"`
	CompanyID   *string   `json:"company_id,omitempty"`
	Date        time.Time `json:"date"`
	Name        string    `json:"name"`
	IsRecurring bool      `json:"is_recurring"`
	CreatedAt   time.Time `json:"created_at"`
}

type LeaveBalance struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	EmployeeID      string    `json:"employee_id"`
	LeaveTypeID     string    `json:"leave_type_id"`
	LeaveTypeName   string    `json:"leave_type_name,omitempty"`
	Year            int       `json:"year"`
	AllocatedDays   float64   `json:"allocated_days"`
	CarriedOverDays float64   `json:"carried_over_days"`
	AdjustmentDays  float64   `json:"adjustment_days"`
	UsedDays        float64   `json:"used_days"`
	ReservedDays    float64   `json:"reserved_days"`
	AvailableDays   float64   `json:"available_days"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LeaveRequest struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	EmployeeID     string          `json:"employee_id"`
	EmployeeName   string          `json:"employee_name,omitempty"`
	LeaveTypeID    string          `json:"leave_type_id"`
	LeaveTypeName  string          `json:"leave_type_name,omitempty"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	RequestedDays  float64         `json:"requested_days"`
	Reason         *string         `json:"reason,omitempty"`
	Status         string          `json:"status"`
	DocumentID     *string         `json:"document_id,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Approvals      []LeaveApproval `json:"approvals,omitempty"`
}

type LeaveApproval struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	LeaveRequestID string     `json:"leave_request_id"`
	ApproverID     string     `json:"approver_id"`
	ApproverName   string     `json:"approver_name,omitempty"`
	Level          int        `json:"level"`
	Status         string     `json:"status"`
	Comments       *string    `json:"comments,omitempty"`
	DecidedAt      *time.Time `json:"decided_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type LeaveRequestHistory struct {
	ID             string     `json:"id"`
	LeaveRequestID string     `json:"leave_request_id"`
	Action         string     `json:"action"`
	OldStatus      *string    `json:"old_status,omitempty"`
	NewStatus      *string    `json:"new_status,omitempty"`
	PerformedBy    string     `json:"performed_by"`
	PerformedByName string   `json:"performed_by_name,omitempty"`
	Comments       *string    `json:"comments,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
