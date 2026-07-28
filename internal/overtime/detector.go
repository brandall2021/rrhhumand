package overtime

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Detector struct {
	repo       *Repository
	calculator *Calculator
	policy     *PolicyEngine
}

func NewDetector(repo *Repository) *Detector {
	return &Detector{
		repo:       repo,
		calculator: NewCalculator(),
		policy:     NewPolicyEngine(),
	}
}

func (d *Detector) DetectForEmployee(ctx context.Context, companyID, employeeID string, workDate time.Time, policy *OvertimePolicy) (*OvertimeRecord, error) {
	record, err := d.repo.GetAttendanceRecord(ctx, companyID, employeeID, workDate)
	if err != nil {
		return nil, err
	}

	if record.ActualStart == nil || record.ActualEnd == nil {
		return nil, nil
	}

	planned := record.ScheduledMinutes
	if planned == 0 {
		planned = 480
	}
	actual := record.WorkedMinutes
	if actual == 0 {
		actual = int(record.ActualEnd.Sub(*record.ActualStart).Minutes())
	}
	late := record.LateMinutes
	earlyLeave := record.EarlyLeaveMinutes

	isWeekend := workDate.Weekday() == time.Saturday || workDate.Weekday() == time.Sunday
	isHoliday := d.isHoliday(ctx, companyID, workDate)

	result := d.calculator.Calculate(planned, actual, late, earlyLeave, policy, isWeekend, isHoliday)

	if result.PotentialOvertimeMinutes <= 0 {
		return nil, nil
	}

	existing, _ := d.repo.GetOvertimeRecordByDateAndType(ctx, companyID, employeeID, workDate, result.OvertimeType)
	if existing != nil {
		existing.OvertimeMinutes = result.RoundedOvertimeMinutes
		existing.ActualMinutes = actual
		existing.PlannedMinutes = planned
		existing.LateMinutes = late
		existing.EarlyLeaveMinutes = earlyLeave
		existing.Status = "DETECTED"
		d.repo.UpdateOvertimeRecord(ctx, existing)
		return existing, nil
	}

	overtimeRecord := &OvertimeRecord{
		ID:                 uuid.New().String(),
		CompanyID:          companyID,
		EmployeeID:         employeeID,
		AttendanceID:       &record.ID,
		WorkDate:           workDate,
		PlannedMinutes:     planned,
		ActualMinutes:      actual,
		LateMinutes:        late,
		EarlyLeaveMinutes:  earlyLeave,
		OvertimeMinutes:    result.RoundedOvertimeMinutes,
		OvertimeType:       result.OvertimeType,
		Status:             "DETECTED",
		IsWeekend:          isWeekend,
		IsHoliday:          isHoliday,
		IsNight:            result.IsNight,
	}

	if err := d.repo.CreateOvertimeRecord(ctx, overtimeRecord); err != nil {
		return nil, err
	}

	return overtimeRecord, nil
}

func (d *Detector) DetectForDateRange(ctx context.Context, companyID string, from, to time.Time, policy *OvertimePolicy) ([]OvertimeRecord, int, error) {
	employees, err := d.repo.GetActiveEmployees(ctx, companyID)
	if err != nil {
		return nil, 0, err
	}

	var records []OvertimeRecord
	count := 0
	current := from
	for !current.After(to) {
		for _, empID := range employees {
			rec, err := d.DetectForEmployee(ctx, companyID, empID, current, policy)
			if err == nil && rec != nil {
				records = append(records, *rec)
				count++
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	return records, count, nil
}

func (d *Detector) isHoliday(ctx context.Context, companyID string, date time.Time) bool {
	_, err := d.repo.GetHoliday(ctx, companyID, date)
	return err == nil
}
