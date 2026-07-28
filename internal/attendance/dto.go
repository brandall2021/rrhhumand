package attendance

import (
	"time"
)

type ClockInRequest struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Source    string   `json:"source" binding:"required"`
	DeviceID  *string  `json:"device_id"`
	Notes     *string  `json:"notes"`
}

type ClockOutRequest struct {
	Source   string  `json:"source"`
	DeviceID *string `json:"device_id"`
	Notes    *string `json:"notes"`
}

type BreakStartRequest struct {
	Source   string  `json:"source"`
	DeviceID *string `json:"device_id"`
}

type BreakEndRequest struct {
	Source   string  `json:"source"`
	DeviceID *string `json:"device_id"`
}

type CreateCorrectionRequest struct {
	AttendanceID   *string `json:"attendance_id"`
	CorrectionType string  `json:"correction_type" binding:"required"`
	RequestedValue string  `json:"requested_value" binding:"required"`
	Reason         string  `json:"reason" binding:"required"`
}

type ApproveCorrectionRequest struct {
	Comments *string `json:"comments"`
}

type CreatePolicyRequest struct {
	Name                     string   `json:"name" binding:"required"`
	ToleranceInMinutes       *int     `json:"tolerance_in_minutes"`
	ToleranceOutMinutes      *int     `json:"tolerance_out_minutes"`
	AllowMobile              *bool    `json:"allow_mobile"`
	AllowWeb                 *bool    `json:"allow_web"`
	AllowKiosk               *bool    `json:"allow_kiosk"`
	RequireGPS               *bool    `json:"require_gps"`
	AllowRemote              *bool    `json:"allow_remote"`
	CalculateOvertime        *bool    `json:"calculate_overtime"`
	RequireCorrectionApproval *bool   `json:"require_correction_approval"`
	WorkStartTime            *string  `json:"work_start_time"`
	WorkEndTime              *string  `json:"work_end_time"`
	BreakDurationMinutes     *int     `json:"break_duration_minutes"`
	WorkDays                 *[]int   `json:"work_days"`
}

type UpdatePolicyRequest struct {
	Name                     *string  `json:"name"`
	ToleranceInMinutes       *int     `json:"tolerance_in_minutes"`
	ToleranceOutMinutes      *int     `json:"tolerance_out_minutes"`
	AllowMobile              *bool    `json:"allow_mobile"`
	AllowWeb                 *bool    `json:"allow_web"`
	AllowKiosk               *bool    `json:"allow_kiosk"`
	RequireGPS               *bool    `json:"require_gps"`
	AllowRemote              *bool    `json:"allow_remote"`
	CalculateOvertime        *bool    `json:"calculate_overtime"`
	RequireCorrectionApproval *bool   `json:"require_correction_approval"`
	WorkStartTime            *string  `json:"work_start_time"`
	WorkEndTime              *string  `json:"work_end_time"`
	BreakDurationMinutes     *int     `json:"break_duration_minutes"`
	WorkDays                 *[]int   `json:"work_days"`
	IsActive                 *bool    `json:"is_active"`
}

type CreateLocationRequest struct {
	Name         string  `json:"name" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	RadiusMeters *int    `json:"radius_meters"`
	BranchID     *string `json:"branch_id"`
}

type CreateDeviceRequest struct {
	DeviceID string  `json:"device_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Location *string `json:"location"`
	BranchID *string `json:"branch_id"`
}

type KioskPunchRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Pin        string `json:"pin"`
	Source     string `json:"source"`
}

type CalendarDay struct {
	Date      time.Time        `json:"date"`
	IsWeekend bool             `json:"is_weekend"`
	IsHoliday bool             `json:"is_holiday"`
	Record    *AttendanceRecord `json:"record,omitempty"`
}
