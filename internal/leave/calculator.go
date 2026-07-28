package leave

import (
	"context"
	"math"
	"time"
)

type DayCalculator struct {
	holidayRepo *Repository
}

func NewDayCalculator(repo *Repository) *DayCalculator {
	return &DayCalculator{holidayRepo: repo}
}

func (c *DayCalculator) CalculateBusinessDays(startDate, endDate time.Time, useBusinessDays bool, companyID string) float64 {
	if !useBusinessDays {
		return float64(endDate.Sub(startDate).Hours()/24) + 1
	}

	ctx := context.Background()
	holidays, _ := c.holidayRepo.GetHolidays(ctx, companyID, startDate, endDate)
	holidayMap := make(map[time.Time]bool)
	for _, h := range holidays {
		holidayMap[time.Date(h.Date.Year(), h.Date.Month(), h.Date.Day(), 0, 0, 0, 0, time.UTC)] = true
	}

	days := 0.0
	current := startDate
	for !current.After(endDate) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			if !holidayMap[time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)] {
				days++
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	return days
}

func (c *DayCalculator) CalculateCalendarDays(startDate, endDate time.Time) float64 {
	return float64(endDate.Sub(startDate).Hours()/24) + 1
}

func (c *DayCalculator) IsWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}

func (c *DayCalculator) IsHoliday(date time.Time, companyID string) bool {
	ctx := context.Background()
	holidays, _ := c.holidayRepo.GetHolidays(ctx, companyID, date, date)
	return len(holidays) > 0
}

func (c *DayCalculator) GetBusinessDaysInPeriod(startDate, endDate time.Time, companyID string) []time.Time {
	ctx := context.Background()
	holidays, _ := c.holidayRepo.GetHolidays(ctx, companyID, startDate, endDate)
	holidayMap := make(map[time.Time]bool)
	for _, h := range holidays {
		holidayMap[time.Date(h.Date.Year(), h.Date.Month(), h.Date.Day(), 0, 0, 0, 0, time.UTC)] = true
	}

	var days []time.Time
	current := startDate
	for !current.After(endDate) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			if !holidayMap[time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)] {
				days = append(days, current)
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	return days
}

func (c *DayCalculator) CalculateBusinessDaysExcludingWeekends(startDate, endDate time.Time) float64 {
	days := 0.0
	current := startDate
	for !current.After(endDate) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}
	return days
}

func (c *DayCalculator) RoundDays(days float64) float64 {
	return math.Round(days*100) / 100
}
