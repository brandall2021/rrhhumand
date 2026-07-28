package scheduling

import (
	"context"
	"time"
)

type ConflictDetector struct {
	repo *Repository
}

func NewConflictDetector(repo *Repository) *ConflictDetector {
	return &ConflictDetector{repo: repo}
}

type ScheduleConflict struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	EntityID    string `json:"entity_id"`
	Date        time.Time `json:"date"`
}

func (d *ConflictDetector) CheckShiftSwap(ctx context.Context, companyID, requesterID, targetID string, requesterDate, targetDate time.Time) ([]ScheduleConflict, error) {
	var conflicts []ScheduleConflict

	// Check if requester already has a shift on target date
	_, err := d.repo.GetEmployeeShiftAssignment(ctx, requesterID, targetDate)
	if err == nil {
		conflicts = append(conflicts, ScheduleConflict{
			Type:        "OVERLAPPING_SHIFT",
			Description: "Requester already has a shift on target date",
			Date:        targetDate,
		})
	}

	// Check if target already has a shift on requester date
	_, err = d.repo.GetEmployeeShiftAssignment(ctx, targetID, requesterDate)
	if err == nil {
		conflicts = append(conflicts, ScheduleConflict{
			Type:        "OVERLAPPING_SHIFT",
			Description: "Target already has a shift on requester date",
			Date:        requesterDate,
		})
	}

	return conflicts, nil
}

func (d *ConflictDetector) CheckShiftAssignment(ctx context.Context, employeeID, shiftID string, workDate time.Time) ([]ScheduleConflict, error) {
	var conflicts []ScheduleConflict

	existing, err := d.repo.GetEmployeeShiftAssignment(ctx, employeeID, workDate)
	if err == nil && existing != nil {
		conflicts = append(conflicts, ScheduleConflict{
			Type:        "DUPLICATE_ASSIGNMENT",
			Description: "Employee already has a shift assigned on this date",
			Date:        workDate,
		})
	}

	return conflicts, nil
}
