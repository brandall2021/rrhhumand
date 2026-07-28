package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo       *Repository
	punches    *Punches
	calculator *Calculator
	geofence   *GeoFence
}

func NewService(repo *Repository, punches *Punches, calculator *Calculator, geofence *GeoFence) *Service {
	return &Service{
		repo:       repo,
		punches:    punches,
		calculator: calculator,
		geofence:   geofence,
	}
}

func (s *Service) ClockIn(ctx context.Context, companyID, employeeID string, req *ClockInRequest) (*AttendanceRecord, error) {
	// Get employee
	_, _, err := s.repo.GetEmployeeByID(ctx, companyID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found")
	}

	// Get policy
	policy, err := s.repo.GetPolicy(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("attendance policy not configured")
	}

	// Validate source
	if err := s.punches.ValidateSource(req.Source, policy); err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Check if already clocked in
	hasClockIn, _ := s.punches.HasClockInToday(ctx, employeeID, today)
	if hasClockIn {
		return nil, fmt.Errorf("already clocked in today")
	}

	// Get schedule
	scheduledStart, _, _ := s.repo.GetEmployeeSchedule(ctx, companyID, employeeID, today)

	// Validate geofence
	if policy.RequireGPS && req.Latitude != nil && req.Longitude != nil {
		locations, _ := s.repo.GetLocations(ctx, companyID)
		valid, msg := s.geofence.ValidatePunch(req.Latitude, req.Longitude, locations, policy.RequireGPS)
		if !valid {
			return nil, fmt.Errorf("geofence validation failed: %s", msg)
		}
	}

	// Create or get record
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	record, err := s.repo.GetOrCreateRecord(ctx, tx, companyID, employeeID, today)
	if err != nil {
		return nil, err
	}

	// Update actual start
	_, err = tx.Exec(ctx, `UPDATE attendance_records SET actual_start=$1, updated_at=NOW() WHERE id=$2`, now, record.ID)
	if err != nil {
		return nil, err
	}

	// Create punch
	punchID := uuid.New().String()
	punch := &AttendancePunch{
		ID:           punchID,
		CompanyID:    companyID,
		EmployeeID:   employeeID,
		AttendanceID: &record.ID,
		PunchType:    "CLOCK_IN",
		PunchedAt:    now,
		Source:       req.Source,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		DeviceID:     req.DeviceID,
		Notes:        req.Notes,
	}
	if _, err := s.punches.CreatePunch(ctx, punch); err != nil {
		return nil, err
	}

	// Calculate late
	lateMinutes, effectiveLate := s.calculator.CalculateLateMinutes(scheduledStart, now, policy.ToleranceInMinutes)
	_, err = tx.Exec(ctx, `UPDATE attendance_records SET late_minutes=$1, effective_late_minutes=$2, updated_at=NOW() WHERE id=$3`, lateMinutes, effectiveLate, record.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Return updated record
	return s.repo.GetRecordByID(ctx, companyID, record.ID)
}

func (s *Service) ClockOut(ctx context.Context, companyID, employeeID string, req *ClockOutRequest) (*AttendanceRecord, error) {
	policy, err := s.repo.GetPolicy(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("attendance policy not configured")
	}

	if req.Source != "" {
		if err := s.punches.ValidateSource(req.Source, policy); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Get today's record
	record, err := s.repo.GetRecordByEmployeeDate(ctx, employeeID, today)
	if err != nil {
		return nil, fmt.Errorf("no attendance record found for today")
	}

	if record.ActualStart == nil {
		return nil, fmt.Errorf("no clock-in found for today")
	}

	if record.ActualEnd != nil {
		return nil, fmt.Errorf("already clocked out today")
	}

	// Get schedule
	scheduledStart, scheduledEnd, _ := s.repo.GetEmployeeSchedule(ctx, companyID, employeeID, today)

	// Get break minutes
	breakMinutes, _ := s.punches.GetBreakMinutes(ctx, employeeID, today)

	// Calculate worked minutes
	workedMinutes := s.calculator.CalculateWorkedMinutes(*record.ActualStart, now, breakMinutes)

	// Calculate early leave
	earlyLeave := s.calculator.CalculateEarlyLeaveMinutes(scheduledEnd, now, policy.ToleranceOutMinutes)

	// Calculate overtime
	overtime := s.calculator.CalculateOvertimeMinutes(scheduledEnd, now)

	// Determine status
	status := s.calculator.DetermineStatus(record.LateMinutes, earlyLeave, overtime, false, false, false, false)

	// Create punch
	punchID := uuid.New().String()
	source := req.Source
	if source == "" {
		source = "WEB"
	}
	punch := &AttendancePunch{
		ID:           punchID,
		CompanyID:    companyID,
		EmployeeID:   employeeID,
		AttendanceID: &record.ID,
		PunchType:    "CLOCK_OUT",
		PunchedAt:    now,
		Source:       source,
		DeviceID:     req.DeviceID,
		Notes:        req.Notes,
	}
	if _, err := s.punches.CreatePunch(ctx, punch); err != nil {
		return nil, err
	}

	// Update record
	record.ActualEnd = &now
	record.WorkedMinutes = workedMinutes
	record.EarlyLeaveMinutes = earlyLeave
	record.OvertimeMinutes = overtime
	record.BreakMinutes = breakMinutes
	record.ScheduledMinutes = s.calculator.CalculateScheduledMinutes(scheduledStart, scheduledEnd)
	record.Status = status

	if err := s.repo.UpdateRecord(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) StartBreak(ctx context.Context, companyID, employeeID string, req *BreakStartRequest) error {
	today := time.Now()
	hasOpen, _ := s.punches.GetOpenBreak(ctx, employeeID, today)
	if hasOpen {
		return fmt.Errorf("break already in progress")
	}

	record, err := s.repo.GetRecordByEmployeeDate(ctx, employeeID, today)
	if err != nil {
		return fmt.Errorf("no attendance record for today")
	}

	punch := &AttendancePunch{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		EmployeeID:   employeeID,
		AttendanceID: &record.ID,
		PunchType:    "BREAK_START",
		PunchedAt:    time.Now(),
		Source:       req.Source,
		DeviceID:     req.DeviceID,
	}
	_, err = s.punches.CreatePunch(ctx, punch)
	return err
}

func (s *Service) EndBreak(ctx context.Context, companyID, employeeID string, req *BreakEndRequest) error {
	hasOpen, _ := s.punches.GetOpenBreak(ctx, employeeID, time.Now())
	if !hasOpen {
		return fmt.Errorf("no open break")
	}

	record, err := s.repo.GetRecordByEmployeeDate(ctx, employeeID, time.Now())
	if err != nil {
		return fmt.Errorf("no attendance record for today")
	}

	punch := &AttendancePunch{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		EmployeeID:   employeeID,
		AttendanceID: &record.ID,
		PunchType:    "BREAK_END",
		PunchedAt:    time.Now(),
		Source:       req.Source,
		DeviceID:     req.DeviceID,
	}
	_, err = s.punches.CreatePunch(ctx, punch)
	return err
}

func (s *Service) GetMyAttendance(ctx context.Context, companyID, employeeID string, from, to time.Time) ([]AttendanceRecord, error) {
	return s.repo.GetMyAttendance(ctx, companyID, employeeID, from, to)
}

func (s *Service) ListRecords(ctx context.Context, companyID string, filters AttendanceFilters, offset, limit int) ([]AttendanceRecord, int64, error) {
	return s.repo.ListRecords(ctx, companyID, filters, offset, limit)
}

func (s *Service) GetRecordByID(ctx context.Context, companyID, id string) (*AttendanceRecord, error) {
	rec, err := s.repo.GetRecordByID(ctx, companyID, id)
	if err != nil {
		return nil, err
	}
	rec.Punches, _ = s.punches.GetPunchesForAttendance(ctx, rec.ID)
	return rec, nil
}

func (s *Service) GetDashboard(ctx context.Context, companyID string) (*AttendanceDashboard, error) {
	return s.repo.GetDashboard(ctx, companyID, time.Now())
}

func (s *Service) GetTeamRecords(ctx context.Context, companyID, managerID string, workDate time.Time) ([]AttendanceRecord, error) {
	return s.repo.GetTeamRecords(ctx, companyID, managerID, workDate)
}

func (s *Service) CreatePolicy(ctx context.Context, companyID string, req *CreatePolicyRequest) (*AttendancePolicy, error) {
	return s.repo.CreatePolicy(ctx, companyID, req)
}

func (s *Service) CreateLocation(ctx context.Context, companyID string, req *CreateLocationRequest) (*AttendanceLocation, error) {
	return s.repo.CreateLocation(ctx, companyID, req)
}

func (s *Service) CreateDevice(ctx context.Context, companyID string, req *CreateDeviceRequest) (*AttendanceDevice, error) {
	return s.repo.CreateDevice(ctx, companyID, req)
}

func (s *Service) CreateCorrection(ctx context.Context, companyID, employeeID string, req *CreateCorrectionRequest) (*AttendanceCorrection, error) {
	requestedValue, err := time.Parse(time.RFC3339, req.RequestedValue)
	if err != nil {
		return nil, fmt.Errorf("invalid requested_value format")
	}

	corr := &AttendanceCorrection{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		EmployeeID:      employeeID,
		AttendanceID:    req.AttendanceID,
		RequestedBy:     employeeID,
		CorrectionType:  req.CorrectionType,
		RequestedValue:  &requestedValue,
		Reason:          req.Reason,
		Status:          "PENDING",
	}
	return s.repo.CreateCorrection(ctx, corr)
}

func (s *Service) ApproveCorrection(ctx context.Context, companyID, correctionID, approverID string) error {
	corr, err := s.repo.GetCorrection(ctx, companyID, correctionID)
	if err != nil {
		return fmt.Errorf("correction not found")
	}
	if corr.Status != "PENDING" {
		return fmt.Errorf("correction is not pending")
	}

	if err := s.repo.UpdateCorrectionStatus(ctx, correctionID, "APPROVED", &approverID); err != nil {
		return err
	}

	// Apply correction to attendance record
	if corr.AttendanceID != nil && corr.RequestedValue != nil {
		switch corr.CorrectionType {
		case "CLOCK_IN":
			s.repo.pool.Exec(ctx, `UPDATE attendance_records SET actual_start=$1, updated_at=NOW() WHERE id=$2`, corr.RequestedValue, corr.AttendanceID)
		case "CLOCK_OUT":
			s.repo.pool.Exec(ctx, `UPDATE attendance_records SET actual_end=$1, updated_at=NOW() WHERE id=$2`, corr.RequestedValue, corr.AttendanceID)
		}
	}

	return nil
}

func (s *Service) RejectCorrection(ctx context.Context, companyID, correctionID, approverID string) error {
	corr, err := s.repo.GetCorrection(ctx, companyID, correctionID)
	if err != nil {
		return fmt.Errorf("correction not found")
	}
	if corr.Status != "PENDING" {
		return fmt.Errorf("correction is not pending")
	}
	return s.repo.UpdateCorrectionStatus(ctx, correctionID, "REJECTED", &approverID)
}

func (s *Service) ListCorrections(ctx context.Context, companyID, status string, offset, limit int) ([]AttendanceCorrection, int64, error) {
	return s.repo.ListCorrections(ctx, companyID, status, offset, limit)
}

func (s *Service) GetCalendar(ctx context.Context, companyID, employeeID string, from, to time.Time) ([]AttendanceRecord, error) {
	return s.repo.GetMyAttendance(ctx, companyID, employeeID, from, to)
}

func (s *Service) ProcessAbsences(ctx context.Context, companyID string, workDate time.Time) error {
	employeeIDs, err := s.repo.GetEmployeesForAbsenceCheck(ctx, companyID, workDate)
	if err != nil {
		return err
	}

	policy, _ := s.repo.GetPolicy(ctx, companyID)
	if policy != nil && !s.calculator.IsWorkDay(workDate, policy.WorkDays) {
		return nil
	}

	for _, empID := range employeeIDs {
		hasRecord, _ := s.repo.HasRecordForDate(ctx, empID, workDate)
		if hasRecord {
			continue
		}

		tx, err := s.repo.BeginTx(ctx)
		if err != nil {
			continue
		}

		scheduledStart, scheduledEnd, _ := s.repo.GetEmployeeSchedule(ctx, companyID, empID, workDate)

		rec := &AttendanceRecord{
			ID:             uuid.New().String(),
			CompanyID:      companyID,
			EmployeeID:     empID,
			WorkDate:       workDate,
			ScheduledStart: &scheduledStart,
			ScheduledEnd:   &scheduledEnd,
			Status:         "ABSENT",
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO attendance_records (id, company_id, employee_id, work_date, scheduled_start, scheduled_end, status)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			rec.ID, rec.CompanyID, rec.EmployeeID, rec.WorkDate, rec.ScheduledStart, rec.ScheduledEnd, rec.Status)
		if err != nil {
			tx.Rollback(ctx)
			continue
		}
		tx.Commit(ctx)
	}

	return nil
}

func (s *Service) ExportCSV(ctx context.Context, companyID string, filters AttendanceFilters) ([]AttendanceRecord, error) {
	records, _, err := s.repo.ListRecords(ctx, companyID, filters, 0, 10000)
	return records, err
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
