package scheduling

import (
	"time"
)

type WorkSchedule struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	ScheduleType string    `json:"schedule_type"`
	Timezone     string    `json:"timezone"`
	WeeklyHours  int       `json:"weekly_hours"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Days         []WorkScheduleDay `json:"days,omitempty"`
}

type WorkScheduleDay struct {
	ID           string    `json:"id"`
	ScheduleID   string    `json:"schedule_id"`
	Weekday      int       `json:"weekday"`
	IsWorkingDay bool      `json:"is_working_day"`
	StartTime    *string   `json:"start_time,omitempty"`
	EndTime      *string   `json:"end_time,omitempty"`
	BreakMinutes int       `json:"break_minutes"`
	CreatedAt    time.Time `json:"created_at"`
	Intervals    []WorkScheduleInterval `json:"intervals,omitempty"`
}

type WorkScheduleInterval struct {
	ID            string `json:"id"`
	ScheduleDayID string `json:"schedule_day_id"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	IntervalType  string `json:"interval_type"`
	Sequence      int    `json:"sequence"`
}

type Shift struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	Name           string    `json:"name"`
	Code           *string   `json:"code,omitempty"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	CrossesMidnight bool    `json:"crosses_midnight"`
	BreakMinutes   int       `json:"break_minutes"`
	Color          *string   `json:"color,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EmployeeScheduleAssignment struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EmployeeID   string    `json:"employee_id"`
	EmployeeName string    `json:"employee_name,omitempty"`
	ScheduleID   string    `json:"schedule_id"`
	ScheduleName string    `json:"schedule_name,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo  *time.Time `json:"effective_to,omitempty"`
	Priority     int        `json:"priority"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
}

type EmployeeShiftAssignment struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	EmployeeID   string    `json:"employee_id"`
	EmployeeName string    `json:"employee_name,omitempty"`
	ShiftID      string    `json:"shift_id"`
	ShiftName    string    `json:"shift_name,omitempty"`
	WorkDate     time.Time `json:"work_date"`
	Status       string    `json:"status"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type RotationTemplate struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CycleLength int       `json:"cycle_length"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Days        []RotationTemplateDay `json:"days,omitempty"`
}

type RotationTemplateDay struct {
	ID          string `json:"id"`
	TemplateID  string `json:"template_id"`
	DayPosition int    `json:"day_position"`
	ShiftID     *string `json:"shift_id,omitempty"`
	ShiftName   string  `json:"shift_name,omitempty"`
	IsRestDay   bool   `json:"is_rest_day"`
}

type EmployeeRotationAssignment struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	EmployeeID      string     `json:"employee_id"`
	TemplateID      string     `json:"template_id"`
	TemplateName    string     `json:"template_name,omitempty"`
	StartDate       time.Time  `json:"start_date"`
	CyclePosition   int        `json:"cycle_position"`
	EffectiveTo     *time.Time `json:"effective_to,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

type EmployeeWorkCalendar struct {
	ID                 string     `json:"id"`
	CompanyID          string     `json:"company_id"`
	EmployeeID         string     `json:"employee_id"`
	WorkDate           time.Time  `json:"work_date"`
	ScheduleID         *string    `json:"schedule_id,omitempty"`
	ShiftID            *string    `json:"shift_id,omitempty"`
	ShiftName          string     `json:"shift_name,omitempty"`
	PlannedStart       *time.Time `json:"planned_start,omitempty"`
	PlannedEnd         *time.Time `json:"planned_end,omitempty"`
	PlannedBreakMinutes int      `json:"planned_break_minutes"`
	Status             string     `json:"status"`
	Source             string     `json:"source"`
	CreatedAt          time.Time  `json:"created_at"`
}

type ScheduleException struct {
	ID             string     `json:"id"`
	CompanyID      string     `json:"company_id"`
	EmployeeID     *string    `json:"employee_id,omitempty"`
	EmployeeName   string     `json:"employee_name,omitempty"`
	ScheduleID     *string    `json:"schedule_id,omitempty"`
	ExceptionDate  time.Time  `json:"exception_date"`
	ExceptionType  string     `json:"exception_type"`
	StartTime      *string    `json:"start_time,omitempty"`
	EndTime        *string    `json:"end_time,omitempty"`
	ShiftID        *string    `json:"shift_id,omitempty"`
	Reason         *string    `json:"reason,omitempty"`
	ApprovedBy     *string    `json:"approved_by,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ShiftSwap struct {
	ID                   string     `json:"id"`
	CompanyID            string     `json:"company_id"`
	RequesterEmployeeID  string     `json:"requester_employee_id"`
	RequesterName        string     `json:"requester_name,omitempty"`
	TargetEmployeeID     string     `json:"target_employee_id"`
	TargetName           string     `json:"target_name,omitempty"`
	RequesterDate        time.Time  `json:"requester_date"`
	TargetDate           time.Time  `json:"target_date"`
	Reason               *string    `json:"reason,omitempty"`
	Status               string     `json:"status"`
	ApprovedBy           *string    `json:"approved_by,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
}

type ResolvedSchedule struct {
	EmployeeID     string     `json:"employee_id"`
	Date           time.Time  `json:"date"`
	ShiftID        *string    `json:"shift_id,omitempty"`
	ShiftName      string     `json:"shift_name,omitempty"`
	PlannedStart   *time.Time `json:"planned_start,omitempty"`
	PlannedEnd     *time.Time `json:"planned_end,omitempty"`
	BreakMinutes   int        `json:"break_minutes"`
	Status         string     `json:"status"`
	Timezone       string     `json:"timezone"`
	IsWorkingDay   bool       `json:"is_working_day"`
}

type CalendarFilters struct {
	DateFrom     string
	DateTo       string
	EmployeeID   string
	DepartmentID string
	BranchID     string
	Status       string
}
