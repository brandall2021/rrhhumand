package overtime

import (
	"time"
)

type Calculator struct {
	rounding *Rounding
	limits   *Limits
}

func NewCalculator() *Calculator {
	return &Calculator{
		rounding: NewRounding(),
		limits:   NewLimits(),
	}
}

func (c *Calculator) Calculate(plannedMinutes, actualMinutes, lateMinutes, earlyLeaveMinutes int, policy *OvertimePolicy, isWeekend, isHoliday bool) *CalculationResult {
	result := &CalculationResult{
		PlannedMinutes:    plannedMinutes,
		ActualMinutes:     actualMinutes,
		LateMinutes:       lateMinutes,
		EarlyLeaveMinutes: earlyLeaveMinutes,
		IsWeekend:         isWeekend,
		IsHoliday:         isHoliday,
	}

	grossOvertime := actualMinutes - plannedMinutes
	if grossOvertime < 0 {
		grossOvertime = 0
	}

	result.PotentialOvertimeMinutes = grossOvertime

	result.OvertimeType = c.determineType(isWeekend, isHoliday, policy)

	if policy != nil && policy.MinimumOvertimeMinutes > 0 && grossOvertime < policy.MinimumOvertimeMinutes {
		result.AllowedOvertimeMinutes = 0
		result.RoundedOvertimeMinutes = 0
		return result
	}

	result.AllowedOvertimeMinutes = grossOvertime
	result.RoundedOvertimeMinutes = grossOvertime

	if policy != nil && policy.RoundingMinutes > 1 {
		result.RoundedOvertimeMinutes = c.rounding.Apply(grossOvertime, policy.RoundingMinutes, "DOWN")
	}

	result.IsNight = false

	return result
}

func (c *Calculator) CalculateWithPolicy(plannedMinutes, actualMinutes, lateMinutes, earlyLeaveMinutes, weeklyTotal, monthlyTotal int, policy *OvertimePolicy, isWeekend, isHoliday bool) *CalculationResult {
	result := c.Calculate(plannedMinutes, actualMinutes, lateMinutes, earlyLeaveMinutes, policy, isWeekend, isHoliday)

	if policy != nil && result.PotentialOvertimeMinutes > 0 {
		allowed, excess := c.limits.ApplyLimits(result.PotentialOvertimeMinutes, weeklyTotal, monthlyTotal, policy)
		result.AllowedOvertimeMinutes = allowed
		result.RoundedOvertimeMinutes = c.rounding.Apply(allowed, policy.RoundingMinutes, "DOWN")
		_ = excess
	}

	return result
}

func (c *Calculator) determineType(isWeekend, isHoliday bool, policy *OvertimePolicy) string {
	if isHoliday {
		return "HOLIDAY"
	}
	if isWeekend {
		return "WEEKEND"
	}
	return "REGULAR"
}

func (c *Calculator) GetDailyPlannedMinutes(resolved interface{}, workDate time.Time) int {
	_ = resolved
	return 480
}
