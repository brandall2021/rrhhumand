package overtime

type CreateOvertimePolicyRequest struct {
	Name                   string   `json:"name" binding:"required"`
	Description            *string  `json:"description"`
	MaxDailyMinutes        *int     `json:"max_daily_minutes"`
	MaxWeeklyMinutes       *int     `json:"max_weekly_minutes"`
	MaxMonthlyMinutes      *int     `json:"max_monthly_minutes"`
	RequiresApproval       *bool    `json:"requires_approval"`
	AllowsCompensation     *bool    `json:"allows_compensation"`
	AllowsPayment          *bool    `json:"allows_payment"`
	MinimumOvertimeMinutes *int     `json:"minimum_overtime_minutes"`
	RoundingMinutes        *int     `json:"rounding_minutes"`
	OvertimeExpirationDays *int     `json:"overtime_expiration_days"`
	NightStart             *string  `json:"night_start"`
	NightEnd               *string  `json:"night_end"`
	WeekendMultiplier      *float64 `json:"weekend_multiplier"`
	HolidayMultiplier      *float64 `json:"holiday_multiplier"`
	NightMultiplier        *float64 `json:"night_multiplier"`
}

type UpdateOvertimePolicyRequest struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	MaxDailyMinutes        *int     `json:"max_daily_minutes"`
	MaxWeeklyMinutes       *int     `json:"max_weekly_minutes"`
	MaxMonthlyMinutes      *int     `json:"max_monthly_minutes"`
	RequiresApproval       *bool    `json:"requires_approval"`
	AllowsCompensation     *bool    `json:"allows_compensation"`
	AllowsPayment          *bool    `json:"allows_payment"`
	MinimumOvertimeMinutes *int     `json:"minimum_overtime_minutes"`
	RoundingMinutes        *int     `json:"rounding_minutes"`
	OvertimeExpirationDays *int     `json:"overtime_expiration_days"`
	NightStart             *string  `json:"night_start"`
	NightEnd               *string  `json:"night_end"`
	WeekendMultiplier      *float64 `json:"weekend_multiplier"`
	HolidayMultiplier      *float64 `json:"holiday_multiplier"`
	NightMultiplier        *float64 `json:"night_multiplier"`
	Status                 *string  `json:"status"`
}

type RequestOvertimeRequest struct {
	WorkDate         string `json:"work_date" binding:"required"`
	RequestedMinutes int    `json:"requested_minutes" binding:"required"`
	Reason           string `json:"reason" binding:"required"`
	OvertimeRecordID *string `json:"overtime_record_id"`
}

type ApproveOvertimeRequest struct {
	ApprovedMinutes *int    `json:"approved_minutes"`
	Comments        *string `json:"comments"`
}

type RejectOvertimeRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type RequestCompensationRequest struct {
	WorkDate string `json:"work_date" binding:"required"`
	Minutes  int    `json:"minutes" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

type ApproveCompensationRequest struct {
	Comments *string `json:"comments"`
}

type DetectOvertimeRequest struct {
	DateFrom string `json:"date_from" binding:"required"`
	DateTo   string `json:"date_to" binding:"required"`
}

type RecalculateRequest struct {
	Date string `json:"date" binding:"required"`
}

type AdjustBalanceRequest struct {
	Minutes int    `json:"minutes" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}
