package integration

import (
	"context"
	"time"
)

type CalendarEvent struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	AllDay          bool      `json:"all_day"`
	Location        string    `json:"location,omitempty"`
	MeetingURL      string    `json:"meeting_url,omitempty"`
	OrganizerID     string    `json:"organizer_id,omitempty"`
	Attendees       []string  `json:"attendees,omitempty"`
	Status          string    `json:"status"`
}

type CalendarAdapter struct{}

func NewCalendarAdapter() *CalendarAdapter {
	return &CalendarAdapter{}
}

func (a *CalendarAdapter) CreateCalendarEvent(ctx context.Context, event *CalendarEvent) (string, error) {
	return "cal-event-" + event.Title, nil
}

func (a *CalendarAdapter) CheckAvailability(ctx context.Context, employeeID string, start, end time.Time) (bool, error) {
	return true, nil
}

func (a *CalendarAdapter) UpdateCalendarEvent(ctx context.Context, eventID string, event *CalendarEvent) error {
	return nil
}

func (a *CalendarAdapter) CancelCalendarEvent(ctx context.Context, eventID string) error {
	return nil
}

func (a *CalendarAdapter) GetEmployeeSchedule(ctx context.Context, employeeID string, from, to time.Time) ([]CalendarEvent, error) {
	return []CalendarEvent{}, nil
}
