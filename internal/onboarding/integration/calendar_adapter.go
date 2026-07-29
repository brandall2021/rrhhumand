package integration

import (
	"context"
	"log"
)

type CalendarAdapter struct{}

func NewCalendarAdapter() *CalendarAdapter {
	return &CalendarAdapter{}
}

func (a *CalendarAdapter) CreateEvent(ctx context.Context, companyID, employeeID, title, description, startDate string) error {
	log.Printf("[CalendarAdapter] CreateEvent company=%s employee=%s title=%s date=%s", companyID, employeeID, title, startDate)
	return nil
}
