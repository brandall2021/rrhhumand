package leave

import "time"

type CreateLeaveTypeRequest struct {
	Name              string  `json:"name" binding:"required"`
	Code              string  `json:"code" binding:"required"`
	Description       *string `json:"description"`
	Category          string  `json:"category" binding:"required"`
	RequiresApproval  *bool   `json:"requires_approval"`
	RequiresDocument  *bool   `json:"requires_document"`
	AffectsBalance    *bool   `json:"affects_balance"`
	IsPaid            *bool   `json:"is_paid"`
}

type UpdateLeaveTypeRequest struct {
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	Category          *string `json:"category"`
	RequiresApproval  *bool   `json:"requires_approval"`
	RequiresDocument  *bool   `json:"requires_document"`
	AffectsBalance    *bool   `json:"affects_balance"`
	IsPaid            *bool   `json:"is_paid"`
	IsActive          *bool   `json:"is_active"`
}

type CreateLeavePolicyRequest struct {
	LeaveTypeID             string   `json:"leave_type_id" binding:"required"`
	Name                    string   `json:"name" binding:"required"`
	DaysPerYear             *float64 `json:"days_per_year"`
	MinimumDaysBeforeRequest *int    `json:"minimum_days_before_request"`
	MaximumDaysPerRequest   *float64 `json:"maximum_days_per_request"`
	MaximumAccumulatedDays  *float64 `json:"maximum_accumulated_days"`
	AllowNegativeBalance    *bool    `json:"allow_negative_balance"`
	UseBusinessDays         *bool    `json:"use_business_days"`
	RequiresManagerApproval *bool    `json:"requires_manager_approval"`
	RequiresHRApproval      *bool    `json:"requires_hr_approval"`
}

type UpdateLeavePolicyRequest struct {
	Name                    *string  `json:"name"`
	DaysPerYear             *float64 `json:"days_per_year"`
	MinimumDaysBeforeRequest *int    `json:"minimum_days_before_request"`
	MaximumDaysPerRequest   *float64 `json:"maximum_days_per_request"`
	MaximumAccumulatedDays  *float64 `json:"maximum_accumulated_days"`
	AllowNegativeBalance    *bool    `json:"allow_negative_balance"`
	UseBusinessDays         *bool    `json:"use_business_days"`
	RequiresManagerApproval *bool    `json:"requires_manager_approval"`
	RequiresHRApproval      *bool    `json:"requires_hr_approval"`
	IsActive                *bool    `json:"is_active"`
}

type CreateHolidayRequest struct {
	Date        string  `json:"date" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	IsRecurring *bool   `json:"is_recurring"`
}

type CreateLeaveRequestRequest struct {
	LeaveTypeID string  `json:"leave_type_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	Reason      *string `json:"reason"`
	DocumentID  *string `json:"document_id"`
}

type ApproveRequest struct {
	Comments *string `json:"comments"`
}

type RejectRequest struct {
	Comments *string `json:"comments" binding:"required"`
}

type AdjustBalanceRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	LeaveTypeID   string  `json:"leave_type_id" binding:"required"`
	Year          int     `json:"year" binding:"required"`
	AdjustmentDays float64 `json:"adjustment_days" binding:"required"`
	Reason        string  `json:"reason" binding:"required"`
}

type LeaveFilters struct {
	EmployeeID   string
	LeaveTypeID  string
	Status       string
	DateFrom     string
	DateTo       string
	DepartmentID string
}

type CalendarFilters struct {
	DateFrom     string
	DateTo       string
	DepartmentID string
	BranchID     string
	EmployeeID   string
}

type LeaveBalanceResponse struct {
	LeaveTypeID    string  `json:"leave_type_id"`
	LeaveTypeName  string  `json:"leave_type_name"`
	Year           int     `json:"year"`
	AllocatedDays  float64 `json:"allocated_days"`
	CarriedOverDays float64 `json:"carried_over_days"`
	AdjustmentDays float64 `json:"adjustment_days"`
	UsedDays       float64 `json:"used_days"`
	ReservedDays   float64 `json:"reserved_days"`
	AvailableDays  float64 `json:"available_days"`
}

type CalendarDay struct {
	Date      time.Time        `json:"date"`
	IsWeekend bool             `json:"is_weekend"`
	IsHoliday bool             `json:"is_holiday"`
	Absences  []CalendarAbsence `json:"absences,omitempty"`
}

type CalendarAbsence struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	LeaveType    string `json:"leave_type"`
	Status       string `json:"status"`
}
