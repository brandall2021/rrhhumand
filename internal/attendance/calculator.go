package attendance

import (
	"math"
	"time"
)

type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) CalculateWorkedMinutes(start, end time.Time, breakMinutes int) int {
	total := int(end.Sub(start).Minutes()) - breakMinutes
	if total < 0 {
		return 0
	}
	return total
}

func (c *Calculator) CalculateLateMinutes(scheduledStart, actualStart time.Time, toleranceMinutes int) (int, int) {
	diff := int(actualStart.Sub(scheduledStart).Minutes())
	if diff <= 0 {
		return 0, 0
	}
	effective := diff - toleranceMinutes
	if effective < 0 {
		effective = 0
	}
	return diff, effective
}

func (c *Calculator) CalculateEarlyLeaveMinutes(scheduledEnd, actualEnd time.Time, toleranceMinutes int) int {
	diff := int(scheduledEnd.Sub(actualEnd).Minutes())
	if diff <= 0 {
		return 0
	}
	effective := diff - toleranceMinutes
	if effective < 0 {
		return 0
	}
	return effective
}

func (c *Calculator) CalculateOvertimeMinutes(scheduledEnd, actualEnd time.Time) int {
	diff := int(actualEnd.Sub(scheduledEnd).Minutes())
	if diff <= 0 {
		return 0
	}
	return diff
}

func (c *Calculator) DetermineStatus(lateMinutes, earlyLeaveMinutes, overtimeMinutes int, isHoliday, isDayOff, hasVacation, hasLeave bool) string {
	if isHoliday {
		return "HOLIDAY"
	}
	if isDayOff {
		return "DAY_OFF"
	}
	if hasVacation {
		return "VACATION"
	}
	if hasLeave {
		return "LEAVE"
	}
	if lateMinutes > 0 && earlyLeaveMinutes > 0 {
		return "PARTIAL"
	}
	if lateMinutes > 0 {
		return "LATE"
	}
	if earlyLeaveMinutes > 0 {
		return "EARLY_LEAVE"
	}
	return "PRESENT"
}

func (c *Calculator) CalculateScheduledMinutes(start, end time.Time) int {
	return int(end.Sub(start).Minutes())
}

func (c *Calculator) ParseTimeString(timeStr string, date time.Time) time.Time {
	t, _ := time.Parse("15:04:05", timeStr)
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

func (c *Calculator) IsWorkDay(date time.Time, workDays []int) bool {
	dayOfWeek := int(date.Weekday())
	for _, d := range workDays {
		if d == dayOfWeek {
			return true
		}
	}
	return false
}

func (c *Calculator) FormatMinutes(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return formatHoursFloat(float64(hours) + float64(mins)/60.0)
}

func formatHoursFloat(hours float64) string {
	h := int(hours)
	m := int(math.Round((hours - float64(h)) * 60))
	return formatHours(h, m)
}

func formatHours(h, m int) string {
	if h == 0 && m == 0 {
		return "0h 0m"
	}
	if h == 0 {
		return formatMinutesStr(m)
	}
	if m == 0 {
		return formatHoursStr(h)
	}
	return formatHoursMinutes(h, m)
}

func formatMinutesStr(m int) string {
	return "0h " + intToStr(m) + "m"
}

func formatHoursStr(h int) string {
	return intToStr(h) + "h 0m"
}

func formatHoursMinutes(h, m int) string {
	return intToStr(h) + "h " + intToStr(m) + "m"
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
