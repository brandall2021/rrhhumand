package attendance

import (
	"time"
)

type AttendancePolicy struct {
	ID                       string    `json:"id"`
	CompanyID                string    `json:"company_id"`
	Name                     string    `json:"name"`
	ToleranceInMinutes       int       `json:"tolerance_in_minutes"`
	ToleranceOutMinutes      int       `json:"tolerance_out_minutes"`
	AllowMobile              bool      `json:"allow_mobile"`
	AllowWeb                 bool      `json:"allow_web"`
	AllowKiosk               bool      `json:"allow_kiosk"`
	RequireGPS               bool      `json:"require_gps"`
	AllowRemote              bool      `json:"allow_remote"`
	CalculateOvertime        bool      `json:"calculate_overtime"`
	RequireCorrectionApproval bool     `json:"require_correction_approval"`
	MaxConsecutiveAbsences   *int      `json:"max_consecutive_absences,omitempty"`
	WorkStartTime            string    `json:"work_start_time"`
	WorkEndTime              string    `json:"work_end_time"`
	BreakDurationMinutes     int       `json:"break_duration_minutes"`
	WorkDays                 []int     `json:"work_days"`
	IsActive                 bool      `json:"is_active"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type AttendanceRecord struct {
	ID                   string     `json:"id"`
	CompanyID            string     `json:"company_id"`
	EmployeeID           string     `json:"employee_id"`
	EmployeeName         string     `json:"employee_name,omitempty"`
	WorkDate             time.Time  `json:"work_date"`
	ScheduledStart       *time.Time `json:"scheduled_start,omitempty"`
	ScheduledEnd         *time.Time `json:"scheduled_end,omitempty"`
	ActualStart          *time.Time `json:"actual_start,omitempty"`
	ActualEnd            *time.Time `json:"actual_end,omitempty"`
	ScheduledMinutes     int        `json:"scheduled_minutes"`
	WorkedMinutes        int        `json:"worked_minutes"`
	LateMinutes          int        `json:"late_minutes"`
	EffectiveLateMinutes int        `json:"effective_late_minutes"`
	EarlyLeaveMinutes    int        `json:"early_leave_minutes"`
	OvertimeMinutes      int        `json:"overtime_minutes"`
	BreakMinutes         int        `json:"break_minutes"`
	Status               string     `json:"status"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Punches              []AttendancePunch `json:"punches,omitempty"`
}

type AttendancePunch struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EmployeeID   string    `json:"employee_id"`
	AttendanceID *string   `json:"attendance_id,omitempty"`
	PunchType    string    `json:"punch_type"`
	PunchedAt    time.Time `json:"punched_at"`
	Source       string    `json:"source"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	IPAddress    *string   `json:"ip_address,omitempty"`
	DeviceID     *string   `json:"device_id,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type AttendanceCorrection struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	EmployeeName    string     `json:"employee_name,omitempty"`
	AttendanceID    *string    `json:"attendance_id,omitempty"`
	RequestedBy     string     `json:"requested_by"`
	RequestedByName string     `json:"requested_by_name,omitempty"`
	ApprovedBy      *string    `json:"approved_by,omitempty"`
	CorrectionType  string     `json:"correction_type"`
	RequestedValue  *time.Time `json:"requested_value,omitempty"`
	OriginalValue   *time.Time `json:"original_value,omitempty"`
	Reason          string     `json:"reason"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

type AttendanceLocation struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	Name          string    `json:"name"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	RadiusMeters  int       `json:"radius_meters"`
	BranchID      *string   `json:"branch_id,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type AttendanceDevice struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	DeviceID   string    `json:"device_id"`
	Name       string    `json:"name"`
	Location   *string   `json:"location,omitempty"`
	BranchID   *string   `json:"branch_id,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type AttendanceFilters struct {
	EmployeeID   string
	DepartmentID string
	Status       string
	DateFrom     string
	DateTo       string
	BranchID     string
}

type AttendanceDashboard struct {
	TotalEmployees    int              `json:"total_employees"`
	Present           int              `json:"present"`
	Absent            int              `json:"absent"`
	Late              int              `json:"late"`
	EarlyLeave        int              `json:"early_leave"`
	OnVacation        int              `json:"on_vacation"`
	OnLeave           int              `json:"on_leave"`
	Holiday           int              `json:"holiday"`
	Remote            int              `json:"remote"`
	AverageClockIn    string           `json:"average_clock_in"`
	AverageClockOut   string           `json:"average_clock_out"`
	AverageHours      float64          `json:"average_hours"`
	Records           []AttendanceRecord `json:"records,omitempty"`
}

type AttendanceSummary struct {
	DaysWorked     int     `json:"days_worked"`
	DaysLate       int     `json:"days_late"`
	DaysAbsent     int     `json:"days_absent"`
	TotalHours     float64 `json:"total_hours"`
	TotalOvertime  float64 `json:"total_overtime"`
	TotalLateMin   int     `json:"total_late_minutes"`
	TotalEarlyMin  int     `json:"total_early_leave_minutes"`
}
