package scheduling

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CalendarGenerator struct {
	repo     *Repository
	resolver *ScheduleResolver
}

func NewCalendarGenerator(repo *Repository, resolver *ScheduleResolver) *CalendarGenerator {
	return &CalendarGenerator{repo: repo, resolver: resolver}
}

func (g *CalendarGenerator) Generate(ctx context.Context, companyID, employeeID string, from, to time.Time) ([]EmployeeWorkCalendar, error) {
	var entries []EmployeeWorkCalendar
	current := from

	for !current.After(to) {
		resolved, err := g.resolver.Resolve(ctx, employeeID, current)
		if err != nil {
			current = current.AddDate(0, 0, 1)
			continue
		}

		entry := &EmployeeWorkCalendar{
			ID:                 uuid.New().String(),
			CompanyID:          companyID,
			EmployeeID:         employeeID,
			WorkDate:           current,
			PlannedBreakMinutes: resolved.BreakMinutes,
			Status:             resolved.Status,
			Source:             "GENERATED",
		}

		if resolved.ShiftID != nil {
			entry.ShiftID = resolved.ShiftID
		}
		if resolved.PlannedStart != nil {
			entry.PlannedStart = resolved.PlannedStart
		}
		if resolved.PlannedEnd != nil {
			entry.PlannedEnd = resolved.PlannedEnd
		}

		if err := g.repo.UpsertCalendarEntry(ctx, entry); err == nil {
			entries = append(entries, *entry)
		}

		current = current.AddDate(0, 0, 1)
	}

	return entries, nil
}

func (g *CalendarGenerator) GenerateForCompany(ctx context.Context, companyID string, from, to time.Time) error {
	// Get all active employees for the company
	rows, err := g.repo.pool.Query(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND status='active'`, companyID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var employeeID string
		if err := rows.Scan(&employeeID); err != nil {
			continue
		}
		g.Generate(ctx, companyID, employeeID, from, to)
	}
	return nil
}
