package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetPolicy(ctx context.Context, companyID string) (*AttendancePolicy, error) {
	p := &AttendancePolicy{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, tolerance_in_minutes, tolerance_out_minutes, allow_mobile, allow_web, allow_kiosk, require_gps, allow_remote, calculate_overtime, require_correction_approval, max_consecutive_absences, work_start_time::TEXT, work_end_time::TEXT, break_duration_minutes, work_days, is_active, created_at, updated_at
		 FROM attendance_policies WHERE company_id=$1 AND is_active=true LIMIT 1`, companyID,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.ToleranceInMinutes, &p.ToleranceOutMinutes, &p.AllowMobile, &p.AllowWeb, &p.AllowKiosk, &p.RequireGPS, &p.AllowRemote, &p.CalculateOvertime, &p.RequireCorrectionApproval, &p.MaxConsecutiveAbsences, &p.WorkStartTime, &p.WorkEndTime, &p.BreakDurationMinutes, &p.WorkDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) CreatePolicy(ctx context.Context, companyID string, req *CreatePolicyRequest) (*AttendancePolicy, error) {
	p := &AttendancePolicy{}
	tolIn := 0
	tolOut := 0
	allowMobile := true
	allowWeb := true
	allowKiosk := true
	requireGPS := false
	allowRemote := true
	calcOT := true
	requireCorrApproval := true
	workStart := "08:00:00"
	workEnd := "16:00:00"
	breakMin := 60
	workDays := []int{1, 2, 3, 4, 5}
	if req.ToleranceInMinutes != nil {
		tolIn = *req.ToleranceInMinutes
	}
	if req.ToleranceOutMinutes != nil {
		tolOut = *req.ToleranceOutMinutes
	}
	if req.AllowMobile != nil {
		allowMobile = *req.AllowMobile
	}
	if req.AllowWeb != nil {
		allowWeb = *req.AllowWeb
	}
	if req.AllowKiosk != nil {
		allowKiosk = *req.AllowKiosk
	}
	if req.RequireGPS != nil {
		requireGPS = *req.RequireGPS
	}
	if req.AllowRemote != nil {
		allowRemote = *req.AllowRemote
	}
	if req.CalculateOvertime != nil {
		calcOT = *req.CalculateOvertime
	}
	if req.RequireCorrectionApproval != nil {
		requireCorrApproval = *req.RequireCorrectionApproval
	}
	if req.WorkStartTime != nil {
		workStart = *req.WorkStartTime
	}
	if req.WorkEndTime != nil {
		workEnd = *req.WorkEndTime
	}
	if req.BreakDurationMinutes != nil {
		breakMin = *req.BreakDurationMinutes
	}
	if req.WorkDays != nil {
		workDays = *req.WorkDays
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO attendance_policies (company_id, name, tolerance_in_minutes, tolerance_out_minutes, allow_mobile, allow_web, allow_kiosk, require_gps, allow_remote, calculate_overtime, require_correction_approval, work_start_time, work_end_time, break_duration_minutes, work_days)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 RETURNING id, company_id, name, tolerance_in_minutes, tolerance_out_minutes, allow_mobile, allow_web, allow_kiosk, require_gps, allow_remote, calculate_overtime, require_correction_approval, max_consecutive_absences, work_start_time::TEXT, work_end_time::TEXT, break_duration_minutes, work_days, is_active, created_at, updated_at`,
		companyID, req.Name, tolIn, tolOut, allowMobile, allowWeb, allowKiosk, requireGPS, allowRemote, calcOT, requireCorrApproval, workStart, workEnd, breakMin, workDays,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.ToleranceInMinutes, &p.ToleranceOutMinutes, &p.AllowMobile, &p.AllowWeb, &p.AllowKiosk, &p.RequireGPS, &p.AllowRemote, &p.CalculateOvertime, &p.RequireCorrectionApproval, &p.MaxConsecutiveAbsences, &p.WorkStartTime, &p.WorkEndTime, &p.BreakDurationMinutes, &p.WorkDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetOrCreateRecord(ctx context.Context, tx pgx.Tx, companyID, employeeID string, workDate time.Time) (*AttendanceRecord, error) {
	rec := &AttendanceRecord{}
	err := tx.QueryRow(ctx,
		`INSERT INTO attendance_records (id, company_id, employee_id, work_date, status)
		 VALUES ($1,$2,$3,$4,'INCOMPLETE')
		 ON CONFLICT (employee_id, work_date) DO UPDATE SET updated_at=NOW()
		 RETURNING id, company_id, employee_id, work_date, scheduled_start, scheduled_end, actual_start, actual_end, scheduled_minutes, worked_minutes, late_minutes, effective_late_minutes, early_leave_minutes, overtime_minutes, break_minutes, status, notes, created_at, updated_at`,
		uuid.New().String(), companyID, employeeID, workDate,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
		&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) GetRecordByID(ctx context.Context, companyID, id string) (*AttendanceRecord, error) {
	rec := &AttendanceRecord{}
	err := r.pool.QueryRow(ctx,
		`SELECT ar.id, ar.company_id, ar.employee_id, e.first_name || ' ' || e.last_name, ar.work_date, ar.scheduled_start, ar.scheduled_end, ar.actual_start, ar.actual_end, ar.scheduled_minutes, ar.worked_minutes, ar.late_minutes, ar.effective_late_minutes, ar.early_leave_minutes, ar.overtime_minutes, ar.break_minutes, ar.status, ar.notes, ar.created_at, ar.updated_at
		 FROM attendance_records ar JOIN employees e ON ar.employee_id=e.id
		 WHERE ar.company_id=$1 AND ar.id=$2`, companyID, id,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
		&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) GetRecordByEmployeeDate(ctx context.Context, employeeID string, workDate time.Time) (*AttendanceRecord, error) {
	rec := &AttendanceRecord{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, work_date, scheduled_start, scheduled_end, actual_start, actual_end, scheduled_minutes, worked_minutes, late_minutes, effective_late_minutes, early_leave_minutes, overtime_minutes, break_minutes, status, notes, created_at, updated_at
		 FROM attendance_records WHERE employee_id=$1 AND work_date=$2`, employeeID, workDate,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
		&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) UpdateRecord(ctx context.Context, rec *AttendanceRecord) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE attendance_records SET
		 scheduled_start=$1, scheduled_end=$2, actual_start=$3, actual_end=$4,
		 scheduled_minutes=$5, worked_minutes=$6, late_minutes=$7, effective_late_minutes=$8,
		 early_leave_minutes=$9, overtime_minutes=$10, break_minutes=$11, status=$12, notes=$13, updated_at=NOW()
		 WHERE id=$14`,
		rec.ScheduledStart, rec.ScheduledEnd, rec.ActualStart, rec.ActualEnd,
		rec.ScheduledMinutes, rec.WorkedMinutes, rec.LateMinutes, rec.EffectiveLateMinutes,
		rec.EarlyLeaveMinutes, rec.OvertimeMinutes, rec.BreakMinutes, rec.Status, rec.Notes, rec.ID)
	return err
}

func (r *Repository) ListRecords(ctx context.Context, companyID string, filters AttendanceFilters, offset, limit int) ([]AttendanceRecord, int64, error) {
	query := `SELECT ar.id, ar.company_id, ar.employee_id, e.first_name || ' ' || e.last_name, ar.work_date, ar.scheduled_start, ar.scheduled_end, ar.actual_start, ar.actual_end, ar.scheduled_minutes, ar.worked_minutes, ar.late_minutes, ar.effective_late_minutes, ar.early_leave_minutes, ar.overtime_minutes, ar.break_minutes, ar.status, ar.notes, ar.created_at, ar.updated_at
		 FROM attendance_records ar JOIN employees e ON ar.employee_id=e.id
		 WHERE ar.company_id=$1`
	countQuery := `SELECT COUNT(*) FROM attendance_records ar WHERE ar.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND ar.employee_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ar.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.DepartmentID != "" {
		query += fmt.Sprintf(" AND e.department_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.department_id=$%d", argIdx)
		args = append(args, filters.DepartmentID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND ar.status=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ar.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND ar.work_date>=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ar.work_date>=$%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND ar.work_date<=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ar.work_date<=$%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY ar.work_date DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
			&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	return records, total, nil
}

func (r *Repository) GetEmployeeSchedule(ctx context.Context, companyID, employeeID string, workDate time.Time) (time.Time, time.Time, error) {
	startTime := "08:00:00"
	endTime := "16:00:00"
	policy, err := r.GetPolicy(ctx, companyID)
	if err == nil {
		startTime = policy.WorkStartTime
		if startTime == "" {
			startTime = "08:00:00"
		}
		endTime = policy.WorkEndTime
		if endTime == "" {
			endTime = "16:00:00"
		}
	}

	startParsed, _ := time.Parse("15:04:05", startTime)
	endParsed, _ := time.Parse("15:04:05", endTime)

	scheduledStart := time.Date(workDate.Year(), workDate.Month(), workDate.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
	scheduledEnd := time.Date(workDate.Year(), workDate.Month(), workDate.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)

	return scheduledStart, scheduledEnd, nil
}

func (r *Repository) GetEmployeeByID(ctx context.Context, companyID, employeeID string) (string, string, error) {
	var firstName, lastName string
	err := r.pool.QueryRow(ctx,
		`SELECT first_name, last_name FROM employees WHERE id=$1 AND company_id=$2`, employeeID, companyID,
	).Scan(&firstName, &lastName)
	return firstName, lastName, err
}

func (r *Repository) GetUserEmployeeID(ctx context.Context, companyID, userID string) (string, error) {
	var employeeID string
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND email=(SELECT email FROM users WHERE id=$2)`, companyID, userID,
	).Scan(&employeeID)
	return employeeID, err
}

func (r *Repository) GetLocations(ctx context.Context, companyID string) ([]AttendanceLocation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, latitude, longitude, radius_meters, branch_id, is_active, created_at
		 FROM attendance_locations WHERE company_id=$1 AND is_active=true`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []AttendanceLocation
	for rows.Next() {
		var loc AttendanceLocation
		if err := rows.Scan(&loc.ID, &loc.CompanyID, &loc.Name, &loc.Latitude, &loc.Longitude, &loc.RadiusMeters, &loc.BranchID, &loc.IsActive, &loc.CreatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func (r *Repository) CreateLocation(ctx context.Context, companyID string, req *CreateLocationRequest) (*AttendanceLocation, error) {
	loc := &AttendanceLocation{}
	radius := 150
	if req.RadiusMeters != nil {
		radius = *req.RadiusMeters
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO attendance_locations (company_id, name, latitude, longitude, radius_meters, branch_id)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, company_id, name, latitude, longitude, radius_meters, branch_id, is_active, created_at`,
		companyID, req.Name, req.Latitude, req.Longitude, radius, req.BranchID,
	).Scan(&loc.ID, &loc.CompanyID, &loc.Name, &loc.Latitude, &loc.Longitude, &loc.RadiusMeters, &loc.BranchID, &loc.IsActive, &loc.CreatedAt)
	return loc, err
}

func (r *Repository) CreateDevice(ctx context.Context, companyID string, req *CreateDeviceRequest) (*AttendanceDevice, error) {
	dev := &AttendanceDevice{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO attendance_devices (company_id, device_id, name, location, branch_id)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id, company_id, device_id, name, location, branch_id, is_active, created_at`,
		companyID, req.DeviceID, req.Name, req.Location, req.BranchID,
	).Scan(&dev.ID, &dev.CompanyID, &dev.DeviceID, &dev.Name, &dev.Location, &dev.BranchID, &dev.IsActive, &dev.CreatedAt)
	return dev, err
}

func (r *Repository) CreateCorrection(ctx context.Context, corr *AttendanceCorrection) (*AttendanceCorrection, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO attendance_corrections (id, company_id, employee_id, attendance_id, requested_by, correction_type, requested_value, original_value, reason, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, company_id, employee_id, attendance_id, requested_by, correction_type, requested_value, original_value, reason, status, created_at, resolved_at`,
		corr.ID, corr.CompanyID, corr.EmployeeID, corr.AttendanceID, corr.RequestedBy, corr.CorrectionType, corr.RequestedValue, corr.OriginalValue, corr.Reason, corr.Status,
	).Scan(&corr.ID, &corr.CompanyID, &corr.EmployeeID, &corr.AttendanceID, &corr.RequestedBy, &corr.CorrectionType, &corr.RequestedValue, &corr.OriginalValue, &corr.Reason, &corr.Status, &corr.CreatedAt, &corr.ResolvedAt)
	return corr, err
}

func (r *Repository) GetCorrection(ctx context.Context, companyID, id string) (*AttendanceCorrection, error) {
	corr := &AttendanceCorrection{}
	err := r.pool.QueryRow(ctx,
		`SELECT ac.id, ac.company_id, ac.employee_id, e.first_name || ' ' || e.last_name, ac.attendance_id, ac.requested_by, u.first_name || ' ' || u.last_name, ac.approved_by, ac.correction_type, ac.requested_value, ac.original_value, ac.reason, ac.status, ac.created_at, ac.resolved_at
		 FROM attendance_corrections ac
		 JOIN employees e ON ac.employee_id=e.id
		 JOIN users u ON ac.requested_by=u.id
		 WHERE ac.company_id=$1 AND ac.id=$2`, companyID, id,
	).Scan(&corr.ID, &corr.CompanyID, &corr.EmployeeID, &corr.EmployeeName, &corr.AttendanceID, &corr.RequestedBy, &corr.RequestedByName, &corr.ApprovedBy, &corr.CorrectionType, &corr.RequestedValue, &corr.OriginalValue, &corr.Reason, &corr.Status, &corr.CreatedAt, &corr.ResolvedAt)
	return corr, err
}

func (r *Repository) ListCorrections(ctx context.Context, companyID string, status string, offset, limit int) ([]AttendanceCorrection, int64, error) {
	query := `SELECT ac.id, ac.company_id, ac.employee_id, e.first_name || ' ' || e.last_name, ac.attendance_id, ac.requested_by, u.first_name || ' ' || u.last_name, ac.approved_by, ac.correction_type, ac.requested_value, ac.original_value, ac.reason, ac.status, ac.created_at, ac.resolved_at
		 FROM attendance_corrections ac
		 JOIN employees e ON ac.employee_id=e.id
		 JOIN users u ON ac.requested_by=u.id
		 WHERE ac.company_id=$1`
	countQuery := `SELECT COUNT(*) FROM attendance_corrections ac WHERE ac.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND ac.status=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ac.status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY ac.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var corrections []AttendanceCorrection
	for rows.Next() {
		var c AttendanceCorrection
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.EmployeeName, &c.AttendanceID, &c.RequestedBy, &c.RequestedByName, &c.ApprovedBy, &c.CorrectionType, &c.RequestedValue, &c.OriginalValue, &c.Reason, &c.Status, &c.CreatedAt, &c.ResolvedAt); err != nil {
			return nil, 0, err
		}
		corrections = append(corrections, c)
	}
	return corrections, total, nil
}

func (r *Repository) UpdateCorrectionStatus(ctx context.Context, id, status string, approvedBy *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE attendance_corrections SET status=$1, approved_by=$2, resolved_at=NOW() WHERE id=$3`, status, approvedBy, id)
	return err
}

func (r *Repository) GetEmployeesForAbsenceCheck(ctx context.Context, companyID string, workDate time.Time) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND status='active'`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) HasRecordForDate(ctx context.Context, employeeID string, workDate time.Time) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM attendance_records WHERE employee_id=$1 AND work_date=$2)`, employeeID, workDate).Scan(&exists)
	return exists, err
}

func (r *Repository) GetTeamRecords(ctx context.Context, companyID, managerID string, workDate time.Time) ([]AttendanceRecord, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ar.id, ar.company_id, ar.employee_id, e.first_name || ' ' || e.last_name, ar.work_date, ar.scheduled_start, ar.scheduled_end, ar.actual_start, ar.actual_end, ar.scheduled_minutes, ar.worked_minutes, ar.late_minutes, ar.effective_late_minutes, ar.early_leave_minutes, ar.overtime_minutes, ar.break_minutes, ar.status, ar.notes, ar.created_at, ar.updated_at
		 FROM attendance_records ar
		 JOIN employees e ON ar.employee_id=e.id
		 WHERE ar.company_id=$1 AND e.manager_id=$2 AND ar.work_date=$3
		 ORDER BY e.first_name, e.last_name`, companyID, managerID, workDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
			&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *Repository) GetDashboard(ctx context.Context, companyID string, workDate time.Time) (*AttendanceDashboard, error) {
	dash := &AttendanceDashboard{}

	// Total active employees
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE company_id=$1 AND status='active'`, companyID).Scan(&dash.TotalEmployees)

	// Count by status
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='PRESENT'`, companyID, workDate).Scan(&dash.Present)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='LATE'`, companyID, workDate).Scan(&dash.Late)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='EARLY_LEAVE'`, companyID, workDate).Scan(&dash.EarlyLeave)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='VACATION'`, companyID, workDate).Scan(&dash.OnVacation)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='LEAVE'`, companyID, workDate).Scan(&dash.OnLeave)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='HOLIDAY'`, companyID, workDate).Scan(&dash.Holiday)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND status='REMOTE'`, companyID, workDate).Scan(&dash.Remote)

	dash.Absent = dash.TotalEmployees - dash.Present - dash.Late - dash.EarlyLeave - dash.OnVacation - dash.OnLeave - dash.Holiday - dash.Remote
	if dash.Absent < 0 {
		dash.Absent = 0
	}

	// Average clock in/out
	var avgClockIn, avgClockOut string
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(TO_CHAR(AVG(actual_start), 'HH24:MI'), '-') FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND actual_start IS NOT NULL`,
		companyID, workDate).Scan(&avgClockIn)
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(TO_CHAR(AVG(actual_end), 'HH24:MI'), '-') FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND actual_end IS NOT NULL`,
		companyID, workDate).Scan(&avgClockOut)
	dash.AverageClockIn = avgClockIn
	dash.AverageClockOut = avgClockOut

	var avgHours float64
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(worked_minutes)/60.0, 0) FROM attendance_records WHERE company_id=$1 AND work_date=$2 AND worked_minutes > 0`,
		companyID, workDate).Scan(&avgHours)
	dash.AverageHours = avgHours

	return dash, nil
}

func (r *Repository) GetMyAttendance(ctx context.Context, companyID, employeeID string, from, to time.Time) ([]AttendanceRecord, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ar.id, ar.company_id, ar.employee_id, e.first_name || ' ' || e.last_name, ar.work_date, ar.scheduled_start, ar.scheduled_end, ar.actual_start, ar.actual_end, ar.scheduled_minutes, ar.worked_minutes, ar.late_minutes, ar.effective_late_minutes, ar.early_leave_minutes, ar.overtime_minutes, ar.break_minutes, ar.status, ar.notes, ar.created_at, ar.updated_at
		 FROM attendance_records ar JOIN employees e ON ar.employee_id=e.id
		 WHERE ar.company_id=$1 AND ar.employee_id=$2 AND ar.work_date>=$3 AND ar.work_date<=$4
		 ORDER BY ar.work_date`, companyID, employeeID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd, &rec.ActualStart, &rec.ActualEnd,
			&rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes, &rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}
