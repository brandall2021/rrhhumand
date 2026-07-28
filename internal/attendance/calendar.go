package attendance

import (
	"time"
)

type Calendar struct{}

func NewCalendar() *Calendar {
	return &Calendar{}
}

type CalendarMonth struct {
	Year  int            `json:"year"`
	Month int            `json:"month"`
	Days  []CalendarDayFull  `json:"days"`
}

type CalendarDayFull struct {
	Date      time.Time        `json:"date"`
	DayOfWeek string           `json:"day_of_week"`
	IsWeekend bool             `json:"is_weekend"`
	IsHoliday bool             `json:"is_holiday"`
	Status    string           `json:"status"`
	Record    *AttendanceRecord `json:"record,omitempty"`
}

func (cal *Calendar) GetMonth(year, month int) CalendarMonth {
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	days := make([]CalendarDayFull, daysInMonth)
	for d := 1; d <= daysInMonth; d++ {
		date := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
		days[d-1] = CalendarDayFull{
			Date:      date,
			DayOfWeek: date.Weekday().String(),
			IsWeekend: date.Weekday() == time.Saturday || date.Weekday() == time.Sunday,
		}
	}
	return CalendarMonth{Year: year, Month: month, Days: days}
}
