package scheduling

type CreateScheduleRequest struct {
	Name         string               `json:"name" binding:"required"`
	Description  *string              `json:"description"`
	ScheduleType string               `json:"schedule_type" binding:"required"`
	Timezone     *string              `json:"timezone"`
	WeeklyHours  *int                 `json:"weekly_hours"`
	Days         []CreateDayRequest    `json:"days"`
}

type CreateDayRequest struct {
	Weekday      int                   `json:"weekday" binding:"required"`
	IsWorkingDay *bool                 `json:"is_working_day"`
	StartTime    *string               `json:"start_time"`
	EndTime      *string               `json:"end_time"`
	BreakMinutes *int                  `json:"break_minutes"`
	Intervals    []CreateIntervalRequest `json:"intervals"`
}

type CreateIntervalRequest struct {
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
	IntervalType string `json:"interval_type"`
	Sequence     *int   `json:"sequence"`
}

type UpdateScheduleRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Timezone     *string `json:"timezone"`
	WeeklyHours  *int    `json:"weekly_hours"`
	Status       *string `json:"status"`
}

type CreateShiftRequest struct {
	Name             string  `json:"name" binding:"required"`
	Code             *string `json:"code"`
	StartTime        string  `json:"start_time" binding:"required"`
	EndTime          string  `json:"end_time" binding:"required"`
	CrossesMidnight  *bool   `json:"crosses_midnight"`
	BreakMinutes     *int    `json:"break_minutes"`
	Color            *string `json:"color"`
}

type UpdateShiftRequest struct {
	Name            *string `json:"name"`
	Code            *string `json:"code"`
	StartTime       *string `json:"start_time"`
	EndTime         *string `json:"end_time"`
	CrossesMidnight *bool   `json:"crosses_midnight"`
	BreakMinutes    *int    `json:"break_minutes"`
	Color           *string `json:"color"`
	Status          *string `json:"status"`
}

type AssignScheduleRequest struct {
	ScheduleID    string  `json:"schedule_id" binding:"required"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	EffectiveTo   *string `json:"effective_to"`
	Priority      *int    `json:"priority"`
}

type AssignShiftRequest struct {
	ShiftID  string `json:"shift_id" binding:"required"`
	WorkDate string `json:"work_date" binding:"required"`
	Notes    *string `json:"notes"`
}

type CreateRotationTemplateRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Description *string                  `json:"description"`
	Days        []CreateRotationDayRequest `json:"days" binding:"required"`
}

type CreateRotationDayRequest struct {
	DayPosition int     `json:"day_position" binding:"required"`
	ShiftID     *string `json:"shift_id"`
	IsRestDay   *bool   `json:"is_rest_day"`
}

type AssignRotationRequest struct {
	TemplateID    string  `json:"template_id" binding:"required"`
	StartDate     string  `json:"start_date" binding:"required"`
	CyclePosition *int    `json:"cycle_position"`
	EffectiveTo   *string `json:"effective_to"`
}

type GenerateCalendarRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	From       string `json:"from" binding:"required"`
	To         string `json:"to" binding:"required"`
}

type CreateExceptionRequest struct {
	EmployeeID    *string `json:"employee_id"`
	ExceptionDate string  `json:"exception_date" binding:"required"`
	ExceptionType string  `json:"exception_type" binding:"required"`
	StartTime     *string `json:"start_time"`
	EndTime       *string `json:"end_time"`
	ShiftID       *string `json:"shift_id"`
	Reason        *string `json:"reason"`
}

type SwapShiftRequest struct {
	TargetEmployeeID string `json:"target_employee_id" binding:"required"`
	RequesterDate    string `json:"requester_date" binding:"required"`
	TargetDate       string `json:"target_date" binding:"required"`
	Reason           *string `json:"reason"`
}
