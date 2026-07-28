package overtime

import (
	"time"
)

type OvertimePolicy struct {
	ID                     string    `json:"id"`
	CompanyID              string    `json:"company_id"`
	Name                   string    `json:"name"`
	Description            *string   `json:"description,omitempty"`
	MaxDailyMinutes        int       `json:"max_daily_minutes"`
	MaxWeeklyMinutes       int       `json:"max_weekly_minutes"`
	MaxMonthlyMinutes      int       `json:"max_monthly_minutes"`
	RequiresApproval       bool      `json:"requires_approval"`
	AllowsCompensation     bool      `json:"allows_compensation"`
	AllowsPayment          bool      `json:"allows_payment"`
	MinimumOvertimeMinutes int       `json:"minimum_overtime_minutes"`
	RoundingMinutes        int       `json:"rounding_minutes"`
	OvertimeExpirationDays int       `json:"overtime_expiration_days"`
	NightStart             string    `json:"night_start"`
	NightEnd               string    `json:"night_end"`
	WeekendMultiplier      float64   `json:"weekend_multiplier"`
	HolidayMultiplier      float64   `json:"holiday_multiplier"`
	NightMultiplier        float64   `json:"night_multiplier"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type OvertimeRecord struct {
	ID                string     `json:"id"`
	CompanyID         string     `json:"company_id"`
	EmployeeID        string     `json:"employee_id"`
	EmployeeName      string     `json:"employee_name,omitempty"`
	AttendanceID      *string    `json:"attendance_id,omitempty"`
	WorkDate          time.Time  `json:"work_date"`
	PlannedMinutes    int        `json:"planned_minutes"`
	ActualMinutes     int        `json:"actual_minutes"`
	LateMinutes       int        `json:"late_minutes"`
	EarlyLeaveMinutes int        `json:"early_leave_minutes"`
	OvertimeMinutes   int        `json:"overtime_minutes"`
	ApprovedMinutes   int        `json:"approved_minutes"`
	CompensatedMinutes int      `json:"compensated_minutes"`
	PaidMinutes       int        `json:"paid_minutes"`
	OvertimeType      string     `json:"overtime_type"`
	Status            string     `json:"status"`
	IsWeekend         bool       `json:"is_weekend"`
	IsHoliday         bool       `json:"is_holiday"`
	IsNight           bool       `json:"is_night"`
	Reason            *string    `json:"reason,omitempty"`
	RejectionReason   *string    `json:"rejection_reason,omitempty"`
	CreatedBy         *string    `json:"created_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type OvertimeRequest struct {
	ID               string     `json:"id"`
	CompanyID        string     `json:"company_id"`
	EmployeeID       string     `json:"employee_id"`
	EmployeeName     string     `json:"employee_name,omitempty"`
	OvertimeRecordID *string    `json:"overtime_record_id,omitempty"`
	WorkDate         time.Time  `json:"work_date"`
	RequestedMinutes int        `json:"requested_minutes"`
	ApprovedMinutes  int        `json:"approved_minutes"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`
	RequestedAt      time.Time  `json:"requested_at"`
	ApprovedBy       *string    `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	RejectionReason  *string    `json:"rejection_reason,omitempty"`
}

type CompensationRequest struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	EmployeeName    string     `json:"employee_name,omitempty"`
	WorkDate        time.Time  `json:"work_date"`
	Minutes         int        `json:"minutes"`
	Reason          string     `json:"reason"`
	Status          string     `json:"status"`
	RequestedAt     time.Time  `json:"requested_at"`
	ApprovedBy      *string    `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
}

type EmployeeTimeBalance struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	EmployeeID      string    `json:"employee_id"`
	EmployeeName    string    `json:"employee_name,omitempty"`
	BalanceMinutes  int       `json:"balance_minutes"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TimeBalanceTransaction struct {
	ID                string    `json:"id"`
	CompanyID         string    `json:"company_id"`
	EmployeeID        string    `json:"employee_id"`
	OvertimeRecordID  *string   `json:"overtime_record_id,omitempty"`
	TransactionType   string    `json:"transaction_type"`
	Minutes           int       `json:"minutes"`
	Reason            *string   `json:"reason,omitempty"`
	CreatedBy         *string   `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type OvertimeDashboard struct {
	TotalDetected     int              `json:"total_detected"`
	TotalPending      int              `json:"total_pending"`
	TotalApproved     int              `json:"total_approved"`
	TotalRejected     int              `json:"total_rejected"`
	TotalCompensated  int              `json:"total_compensated"`
	TotalPaid         int              `json:"total_paid"`
	TotalMinutes      int              `json:"total_minutes"`
	BalanceMinutes    int              `json:"balance_minutes"`
	Records           []OvertimeRecord `json:"records,omitempty"`
}

type OvertimeFilters struct {
	EmployeeID   string
	DepartmentID string
	Status       string
	OvertimeType string
	DateFrom     string
	DateTo       string
}

type CalculationResult struct {
	PlannedMinutes          int
	ActualMinutes           int
	LateMinutes             int
	EarlyLeaveMinutes       int
	PotentialOvertimeMinutes int
	AllowedOvertimeMinutes  int
	RoundedOvertimeMinutes  int
	OvertimeType            string
	IsWeekend               bool
	IsHoliday               bool
	IsNight                 bool
}
