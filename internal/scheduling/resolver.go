package scheduling

import (
	"context"
	"time"
)

type ScheduleResolver struct {
	repo *Repository
}

func NewScheduleResolver(repo *Repository) *ScheduleResolver {
	return &ScheduleResolver{repo: repo}
}

func (r *ScheduleResolver) Resolve(ctx context.Context, employeeID string, date time.Time) (*ResolvedSchedule, error) {
	result := &ResolvedSchedule{
		EmployeeID:   employeeID,
		Date:         date,
		Status:       "DAY_OFF",
		IsWorkingDay: false,
		Timezone:     "UTC",
	}

	// 1. Check exception first (highest priority)
	// We need companyID - get it from the assignment
	scheduleAssign, err := r.repo.GetEmployeeScheduleAssignment(ctx, employeeID, date)
	if err == nil && scheduleAssign != nil {
		result.Timezone = "UTC"

		// Check exception
		exception, err := r.repo.GetException(ctx, scheduleAssign.CompanyID, employeeID, date)
		if err == nil && exception != nil {
			return r.resolveException(result, exception), nil
		}

		// Check shift assignment
		shiftAssign, err := r.repo.GetEmployeeShiftAssignment(ctx, employeeID, date)
		if err == nil && shiftAssign != nil {
			return r.resolveFromShiftAssignment(ctx, result, shiftAssign), nil
		}

		// Resolve from schedule
		return r.resolveFromSchedule(ctx, result, scheduleAssign), nil
	}

	// 2. Check rotation assignment
	rotationAssign, err := r.repo.GetEmployeeRotation(ctx, employeeID, date)
	if err == nil && rotationAssign != nil {
		return r.resolveFromRotation(ctx, result, rotationAssign), nil
	}

	return result, nil
}

func (r *ScheduleResolver) resolveException(result *ResolvedSchedule, exc *ScheduleException) *ResolvedSchedule {
	switch exc.ExceptionType {
	case "DAY_OFF":
		result.Status = "DAY_OFF"
		result.IsWorkingDay = false
	case "SPECIAL_SCHEDULE":
		result.Status = "WORKING"
		result.IsWorkingDay = true
		if exc.StartTime != nil && exc.EndTime != nil {
			startParsed, _ := time.Parse("15:04:05", *exc.StartTime)
			endParsed, _ := time.Parse("15:04:05", *exc.EndTime)
			start := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
			end := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)
			result.PlannedStart = &start
			result.PlannedEnd = &end
		}
	case "HOLIDAY_WORK":
		result.Status = "WORKING"
		result.IsWorkingDay = true
	}
	return result
}

func (r *ScheduleResolver) resolveFromShiftAssignment(ctx context.Context, result *ResolvedSchedule, assign *EmployeeShiftAssignment) *ResolvedSchedule {
	shift, err := r.repo.GetShiftByID(ctx, assign.ShiftID)
	if err != nil {
		return result
	}

	result.ShiftID = &assign.ShiftID
	result.ShiftName = shift.Name
	result.BreakMinutes = shift.BreakMinutes
	result.Status = "WORKING"
	result.IsWorkingDay = true

	startParsed, _ := time.Parse("15:04:05", shift.StartTime)
	endParsed, _ := time.Parse("15:04:05", shift.EndTime)

	start := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
	end := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)

	if shift.CrossesMidnight {
		end = end.AddDate(0, 0, 1)
	}

	result.PlannedStart = &start
	result.PlannedEnd = &end
	return result
}

func (r *ScheduleResolver) resolveFromSchedule(ctx context.Context, result *ResolvedSchedule, assign *EmployeeScheduleAssignment) *ResolvedSchedule {
	schedule, err := r.repo.GetSchedule(ctx, assign.CompanyID, assign.ScheduleID)
	if err != nil {
		return result
	}

	result.Timezone = schedule.Timezone
	weekday := int(result.Date.Weekday())

	for _, day := range schedule.Days {
		if day.Weekday == weekday {
			if !day.IsWorkingDay {
				result.Status = "DAY_OFF"
				result.IsWorkingDay = false
				return result
			}
			result.Status = "WORKING"
			result.IsWorkingDay = true
			result.BreakMinutes = day.BreakMinutes

			if day.StartTime != nil && day.EndTime != nil {
				startParsed, _ := time.Parse("15:04:05", *day.StartTime)
				endParsed, _ := time.Parse("15:04:05", *day.EndTime)
				start := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
				end := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)
				result.PlannedStart = &start
				result.PlannedEnd = &end
			}
			break
		}
	}
	return result
}

func (r *ScheduleResolver) resolveFromRotation(ctx context.Context, result *ResolvedSchedule, assign *EmployeeRotationAssignment) *ResolvedSchedule {
	template, err := r.repo.GetRotationTemplate(ctx, assign.CompanyID, assign.TemplateID)
	if err != nil {
		return result
	}

	daysSinceStart := int(result.Date.Sub(assign.StartDate).Hours() / 24)
	position := ((daysSinceStart + assign.CyclePosition - 1) % template.CycleLength) + 1

	for _, day := range template.Days {
		if day.DayPosition == position {
			if day.IsRestDay {
				result.Status = "DAY_OFF"
				result.IsWorkingDay = false
				return result
			}
			if day.ShiftID != nil {
				shift, err := r.repo.GetShiftByID(ctx, *day.ShiftID)
				if err == nil {
					result.ShiftID = day.ShiftID
					result.ShiftName = shift.Name
					result.BreakMinutes = shift.BreakMinutes
					result.Status = "WORKING"
					result.IsWorkingDay = true

					startParsed, _ := time.Parse("15:04:05", shift.StartTime)
					endParsed, _ := time.Parse("15:04:05", shift.EndTime)
					start := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
					end := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)
					if shift.CrossesMidnight {
						end = end.AddDate(0, 0, 1)
					}
					result.PlannedStart = &start
					result.PlannedEnd = &end
				}
			}
			break
		}
	}
	return result
}
