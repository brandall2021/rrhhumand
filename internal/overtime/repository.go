package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Policies
func (r *Repository) CreatePolicy(ctx context.Context, companyID string, req *CreateOvertimePolicyRequest) (*OvertimePolicy, error) {
	p := &OvertimePolicy{}
	maxDaily := 120
	maxWeekly := 480
	maxMonthly := 1920
	requiresApproval := true
	allowsCompensation := true
	allowsPayment := true
	minOT := 0
	rounding := 1
	expiration := 0
	nightStart := "22:00:00"
	nightEnd := "06:00:00"
	weekendMult := 1.5
	holidayMult := 2.0
	nightMult := 1.5

	if req.MaxDailyMinutes != nil { maxDaily = *req.MaxDailyMinutes }
	if req.MaxWeeklyMinutes != nil { maxWeekly = *req.MaxWeeklyMinutes }
	if req.MaxMonthlyMinutes != nil { maxMonthly = *req.MaxMonthlyMinutes }
	if req.RequiresApproval != nil { requiresApproval = *req.RequiresApproval }
	if req.AllowsCompensation != nil { allowsCompensation = *req.AllowsCompensation }
	if req.AllowsPayment != nil { allowsPayment = *req.AllowsPayment }
	if req.MinimumOvertimeMinutes != nil { minOT = *req.MinimumOvertimeMinutes }
	if req.RoundingMinutes != nil { rounding = *req.RoundingMinutes }
	if req.OvertimeExpirationDays != nil { expiration = *req.OvertimeExpirationDays }
	if req.NightStart != nil { nightStart = *req.NightStart }
	if req.NightEnd != nil { nightEnd = *req.NightEnd }
	if req.WeekendMultiplier != nil { weekendMult = *req.WeekendMultiplier }
	if req.HolidayMultiplier != nil { holidayMult = *req.HolidayMultiplier }
	if req.NightMultiplier != nil { nightMult = *req.NightMultiplier }

	err := r.pool.QueryRow(ctx,
		`INSERT INTO overtime_policies (company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING id, company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier, status, created_at, updated_at`,
		companyID, req.Name, req.Description, maxDaily, maxWeekly, maxMonthly,
		requiresApproval, allowsCompensation, allowsPayment, minOT, rounding, expiration,
		nightStart, nightEnd, weekendMult, holidayMult, nightMult,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.MaxDailyMinutes, &p.MaxWeeklyMinutes, &p.MaxMonthlyMinutes,
		&p.RequiresApproval, &p.AllowsCompensation, &p.AllowsPayment, &p.MinimumOvertimeMinutes, &p.RoundingMinutes, &p.OvertimeExpirationDays,
		&p.NightStart, &p.NightEnd, &p.WeekendMultiplier, &p.HolidayMultiplier, &p.NightMultiplier,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetPolicy(ctx context.Context, companyID, id string) (*OvertimePolicy, error) {
	p := &OvertimePolicy{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier, status, created_at, updated_at
		 FROM overtime_policies WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.MaxDailyMinutes, &p.MaxWeeklyMinutes, &p.MaxMonthlyMinutes,
		&p.RequiresApproval, &p.AllowsCompensation, &p.AllowsPayment, &p.MinimumOvertimeMinutes, &p.RoundingMinutes, &p.OvertimeExpirationDays,
		&p.NightStart, &p.NightEnd, &p.WeekendMultiplier, &p.HolidayMultiplier, &p.NightMultiplier,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetActivePolicy(ctx context.Context, companyID string) (*OvertimePolicy, error) {
	p := &OvertimePolicy{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier, status, created_at, updated_at
		 FROM overtime_policies WHERE company_id=$1 AND status='ACTIVE' ORDER BY created_at DESC LIMIT 1`, companyID,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.MaxDailyMinutes, &p.MaxWeeklyMinutes, &p.MaxMonthlyMinutes,
		&p.RequiresApproval, &p.AllowsCompensation, &p.AllowsPayment, &p.MinimumOvertimeMinutes, &p.RoundingMinutes, &p.OvertimeExpirationDays,
		&p.NightStart, &p.NightEnd, &p.WeekendMultiplier, &p.HolidayMultiplier, &p.NightMultiplier,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ListPolicies(ctx context.Context, companyID string) ([]OvertimePolicy, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier, status, created_at, updated_at
		 FROM overtime_policies WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var policies []OvertimePolicy
	for rows.Next() {
		var p OvertimePolicy
		rows.Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.MaxDailyMinutes, &p.MaxWeeklyMinutes, &p.MaxMonthlyMinutes,
			&p.RequiresApproval, &p.AllowsCompensation, &p.AllowsPayment, &p.MinimumOvertimeMinutes, &p.RoundingMinutes, &p.OvertimeExpirationDays,
			&p.NightStart, &p.NightEnd, &p.WeekendMultiplier, &p.HolidayMultiplier, &p.NightMultiplier,
			&p.Status, &p.CreatedAt, &p.UpdatedAt)
		policies = append(policies, p)
	}
	return policies, nil
}

func (r *Repository) UpdatePolicy(ctx context.Context, companyID, id string, req *UpdateOvertimePolicyRequest) (*OvertimePolicy, error) {
	p := &OvertimePolicy{}
	err := r.pool.QueryRow(ctx,
		`UPDATE overtime_policies SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 max_daily_minutes=COALESCE($5,max_daily_minutes), max_weekly_minutes=COALESCE($6,max_weekly_minutes),
		 max_monthly_minutes=COALESCE($7,max_monthly_minutes), requires_approval=COALESCE($8,requires_approval),
		 allows_compensation=COALESCE($9,allows_compensation), allows_payment=COALESCE($10,allows_payment),
		 minimum_overtime_minutes=COALESCE($11,minimum_overtime_minutes), rounding_minutes=COALESCE($12,rounding_minutes),
		 overtime_expiration_days=COALESCE($13,overtime_expiration_days), night_start=COALESCE($14,night_start),
		 night_end=COALESCE($15,night_end), weekend_multiplier=COALESCE($16,weekend_multiplier),
		 holiday_multiplier=COALESCE($17,holiday_multiplier), night_multiplier=COALESCE($18,night_multiplier),
		 status=COALESCE($19,status), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, max_daily_minutes, max_weekly_minutes, max_monthly_minutes,
		 requires_approval, allows_compensation, allows_payment, minimum_overtime_minutes, rounding_minutes, overtime_expiration_days,
		 night_start, night_end, weekend_multiplier, holiday_multiplier, night_multiplier, status, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.MaxDailyMinutes, req.MaxWeeklyMinutes, req.MaxMonthlyMinutes,
		req.RequiresApproval, req.AllowsCompensation, req.AllowsPayment, req.MinimumOvertimeMinutes, req.RoundingMinutes,
		req.OvertimeExpirationDays, req.NightStart, req.NightEnd, req.WeekendMultiplier, req.HolidayMultiplier, req.NightMultiplier, req.Status,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.MaxDailyMinutes, &p.MaxWeeklyMinutes, &p.MaxMonthlyMinutes,
		&p.RequiresApproval, &p.AllowsCompensation, &p.AllowsPayment, &p.MinimumOvertimeMinutes, &p.RoundingMinutes, &p.OvertimeExpirationDays,
		&p.NightStart, &p.NightEnd, &p.WeekendMultiplier, &p.HolidayMultiplier, &p.NightMultiplier,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) DeletePolicy(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM overtime_policies WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

// Overtime Records
func (r *Repository) CreateOvertimeRecord(ctx context.Context, rec *OvertimeRecord) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO overtime_records (id, company_id, employee_id, attendance_id, work_date, planned_minutes, actual_minutes,
		 late_minutes, early_leave_minutes, overtime_minutes, approved_minutes, compensated_minutes, paid_minutes,
		 overtime_type, status, is_weekend, is_holiday, is_night, reason, rejection_reason, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`,
		rec.ID, rec.CompanyID, rec.EmployeeID, rec.AttendanceID, rec.WorkDate, rec.PlannedMinutes, rec.ActualMinutes,
		rec.LateMinutes, rec.EarlyLeaveMinutes, rec.OvertimeMinutes, rec.ApprovedMinutes, rec.CompensatedMinutes, rec.PaidMinutes,
		rec.OvertimeType, rec.Status, rec.IsWeekend, rec.IsHoliday, rec.IsNight, rec.Reason, rec.RejectionReason, rec.CreatedBy)
	return err
}

func (r *Repository) GetOvertimeRecord(ctx context.Context, companyID, id string) (*OvertimeRecord, error) {
	rec := &OvertimeRecord{}
	err := r.pool.QueryRow(ctx,
		`SELECT o.id, o.company_id, o.employee_id, COALESCE(e.first_name||' '||e.last_name,''), o.attendance_id, o.work_date,
		 o.planned_minutes, o.actual_minutes, o.late_minutes, o.early_leave_minutes, o.overtime_minutes,
		 o.approved_minutes, o.compensated_minutes, o.paid_minutes, o.overtime_type, o.status,
		 o.is_weekend, o.is_holiday, o.is_night, o.reason, o.rejection_reason, o.created_by, o.created_at, o.updated_at
		 FROM overtime_records o
		 LEFT JOIN employees e ON o.employee_id=e.id
		 WHERE o.company_id=$1 AND o.id=$2`, companyID, id,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.AttendanceID, &rec.WorkDate,
		&rec.PlannedMinutes, &rec.ActualMinutes, &rec.LateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes,
		&rec.ApprovedMinutes, &rec.CompensatedMinutes, &rec.PaidMinutes, &rec.OvertimeType, &rec.Status,
		&rec.IsWeekend, &rec.IsHoliday, &rec.IsNight, &rec.Reason, &rec.RejectionReason, &rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) GetOvertimeRecordByDateAndType(ctx context.Context, companyID, employeeID string, workDate time.Time, overtimeType string) (*OvertimeRecord, error) {
	rec := &OvertimeRecord{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, attendance_id, work_date, planned_minutes, actual_minutes,
		 late_minutes, early_leave_minutes, overtime_minutes, approved_minutes, compensated_minutes, paid_minutes,
		 overtime_type, status, is_weekend, is_holiday, is_night, reason, rejection_reason, created_by, created_at, updated_at
		 FROM overtime_records WHERE company_id=$1 AND employee_id=$2 AND work_date=$3 AND overtime_type=$4`, companyID, employeeID, workDate, overtimeType,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.AttendanceID, &rec.WorkDate,
		&rec.PlannedMinutes, &rec.ActualMinutes, &rec.LateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes,
		&rec.ApprovedMinutes, &rec.CompensatedMinutes, &rec.PaidMinutes, &rec.OvertimeType, &rec.Status,
		&rec.IsWeekend, &rec.IsHoliday, &rec.IsNight, &rec.Reason, &rec.RejectionReason, &rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) UpdateOvertimeRecord(ctx context.Context, rec *OvertimeRecord) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE overtime_records SET
		 attendance_id=$3, planned_minutes=$4, actual_minutes=$5, late_minutes=$6, early_leave_minutes=$7,
		 overtime_minutes=$8, approved_minutes=$9, compensated_minutes=$10, paid_minutes=$11,
		 overtime_type=$12, status=$13, is_weekend=$14, is_holiday=$15, is_night=$16,
		 reason=$17, rejection_reason=$18, updated_at=NOW()
		 WHERE id=$1 AND company_id=$2`,
		rec.ID, rec.CompanyID, rec.AttendanceID, rec.PlannedMinutes, rec.ActualMinutes, rec.LateMinutes, rec.EarlyLeaveMinutes,
		rec.OvertimeMinutes, rec.ApprovedMinutes, rec.CompensatedMinutes, rec.PaidMinutes,
		rec.OvertimeType, rec.Status, rec.IsWeekend, rec.IsHoliday, rec.IsNight,
		rec.Reason, rec.RejectionReason)
	return err
}

func (r *Repository) ListOvertimeRecords(ctx context.Context, companyID string, filters OvertimeFilters) ([]OvertimeRecord, error) {
	query := `SELECT o.id, o.company_id, o.employee_id, COALESCE(e.first_name||' '||e.last_name,''), o.attendance_id, o.work_date,
		 o.planned_minutes, o.actual_minutes, o.late_minutes, o.early_leave_minutes, o.overtime_minutes,
		 o.approved_minutes, o.compensated_minutes, o.paid_minutes, o.overtime_type, o.status,
		 o.is_weekend, o.is_holiday, o.is_night, o.reason, o.rejection_reason, o.created_by, o.created_at, o.updated_at
		 FROM overtime_records o
		 LEFT JOIN employees e ON o.employee_id=e.id
		 WHERE o.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND o.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND o.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.OvertimeType != "" {
		query += fmt.Sprintf(" AND o.overtime_type=$%d", argIdx)
		args = append(args, filters.OvertimeType)
		argIdx++
	}
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND o.work_date>=$%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND o.work_date<=$%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}

	query += " ORDER BY o.work_date DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var records []OvertimeRecord
	for rows.Next() {
		var rec OvertimeRecord
		rows.Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.EmployeeName, &rec.AttendanceID, &rec.WorkDate,
			&rec.PlannedMinutes, &rec.ActualMinutes, &rec.LateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes,
			&rec.ApprovedMinutes, &rec.CompensatedMinutes, &rec.PaidMinutes, &rec.OvertimeType, &rec.Status,
			&rec.IsWeekend, &rec.IsHoliday, &rec.IsNight, &rec.Reason, &rec.RejectionReason, &rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt)
		records = append(records, rec)
	}
	return records, nil
}

func (r *Repository) GetWeeklyOvertimeMinutes(ctx context.Context, companyID, employeeID string, weekStart time.Time) (int, error) {
	var total int
	weekEnd := weekStart.AddDate(0, 0, 7)
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(overtime_minutes),0) FROM overtime_records
		 WHERE company_id=$1 AND employee_id=$2 AND work_date>=$3 AND work_date<$4 AND status IN ('DETECTED','PENDING','REQUESTED','SUBMITTED','APPROVED')`,
		companyID, employeeID, weekStart, weekEnd).Scan(&total)
	return total, err
}

func (r *Repository) GetMonthlyOvertimeMinutes(ctx context.Context, companyID, employeeID string, monthStart time.Time) (int, error) {
	var total int
	monthEnd := time.Date(monthStart.Year(), monthStart.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(overtime_minutes),0) FROM overtime_records
		 WHERE company_id=$1 AND employee_id=$2 AND work_date>=$3 AND work_date<$4 AND status IN ('DETECTED','PENDING','REQUESTED','SUBMITTED','APPROVED')`,
		companyID, employeeID, monthStart, monthEnd).Scan(&total)
	return total, err
}

func (r *Repository) GetOvertimeDashboard(ctx context.Context, companyID string) (*OvertimeDashboard, error) {
	dash := &OvertimeDashboard{}
	err := r.pool.QueryRow(ctx,
		`SELECT
		 COALESCE(SUM(CASE WHEN status='DETECTED' THEN overtime_minutes ELSE 0 END),0),
		 COALESCE(SUM(CASE WHEN status IN ('PENDING','REQUESTED','SUBMITTED') THEN overtime_minutes ELSE 0 END),0),
		 COALESCE(SUM(CASE WHEN status='APPROVED' THEN approved_minutes ELSE 0 END),0),
		 COALESCE(SUM(CASE WHEN status='REJECTED' THEN overtime_minutes ELSE 0 END),0),
		 COALESCE(SUM(compensated_minutes),0),
		 COALESCE(SUM(paid_minutes),0),
		 COALESCE(SUM(overtime_minutes),0)
		 FROM overtime_records WHERE company_id=$1`, companyID,
	).Scan(&dash.TotalDetected, &dash.TotalPending, &dash.TotalApproved, &dash.TotalRejected,
		&dash.TotalCompensated, &dash.TotalPaid, &dash.TotalMinutes)
	return dash, err
}

// Overtime Requests
func (r *Repository) CreateOvertimeRequest(ctx context.Context, req *OvertimeRequest) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO overtime_requests (id, company_id, employee_id, overtime_record_id, work_date, requested_minutes, approved_minutes, reason, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		req.ID, req.CompanyID, req.EmployeeID, req.OvertimeRecordID, req.WorkDate, req.RequestedMinutes, req.ApprovedMinutes, req.Reason, req.Status)
	return err
}

func (r *Repository) GetOvertimeRequest(ctx context.Context, companyID, id string) (*OvertimeRequest, error) {
	req := &OvertimeRequest{}
	err := r.pool.QueryRow(ctx,
		`SELECT o.id, o.company_id, o.employee_id, COALESCE(e.first_name||' '||e.last_name,''), o.overtime_record_id, o.work_date,
		 o.requested_minutes, o.approved_minutes, o.reason, o.status, o.requested_at, o.approved_by, o.approved_at, o.rejection_reason
		 FROM overtime_requests o
		 LEFT JOIN employees e ON o.employee_id=e.id
		 WHERE o.company_id=$1 AND o.id=$2`, companyID, id,
	).Scan(&req.ID, &req.CompanyID, &req.EmployeeID, &req.EmployeeName, &req.OvertimeRecordID, &req.WorkDate,
		&req.RequestedMinutes, &req.ApprovedMinutes, &req.Reason, &req.Status, &req.RequestedAt,
		&req.ApprovedBy, &req.ApprovedAt, &req.RejectionReason)
	return req, err
}

func (r *Repository) UpdateOvertimeRequest(ctx context.Context, req *OvertimeRequest) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE overtime_requests SET approved_minutes=$3, status=$4, approved_by=$5, approved_at=$6, rejection_reason=$7 WHERE id=$1 AND company_id=$2`,
		req.ID, req.CompanyID, req.ApprovedMinutes, req.Status, req.ApprovedBy, req.ApprovedAt, req.RejectionReason)
	return err
}

func (r *Repository) ListOvertimeRequests(ctx context.Context, companyID string, filters OvertimeFilters) ([]OvertimeRequest, error) {
	query := `SELECT o.id, o.company_id, o.employee_id, COALESCE(e.first_name||' '||e.last_name,''), o.overtime_record_id, o.work_date,
		 o.requested_minutes, o.approved_minutes, o.reason, o.status, o.requested_at, o.approved_by, o.approved_at, o.rejection_reason
		 FROM overtime_requests o
		 LEFT JOIN employees e ON o.employee_id=e.id
		 WHERE o.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND o.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND o.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND o.work_date>=$%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND o.work_date<=$%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}

	query += " ORDER BY o.requested_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var requests []OvertimeRequest
	for rows.Next() {
		var req OvertimeRequest
		rows.Scan(&req.ID, &req.CompanyID, &req.EmployeeID, &req.EmployeeName, &req.OvertimeRecordID, &req.WorkDate,
			&req.RequestedMinutes, &req.ApprovedMinutes, &req.Reason, &req.Status, &req.RequestedAt,
			&req.ApprovedBy, &req.ApprovedAt, &req.RejectionReason)
		requests = append(requests, req)
	}
	return requests, nil
}

// Compensation Requests
func (r *Repository) CreateCompensationRequest(ctx context.Context, comp *CompensationRequest) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO compensation_requests (id, company_id, employee_id, work_date, minutes, reason, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		comp.ID, comp.CompanyID, comp.EmployeeID, comp.WorkDate, comp.Minutes, comp.Reason, comp.Status)
	return err
}

func (r *Repository) GetCompensationRequest(ctx context.Context, companyID, id string) (*CompensationRequest, error) {
	comp := &CompensationRequest{}
	err := r.pool.QueryRow(ctx,
		`SELECT c.id, c.company_id, c.employee_id, COALESCE(e.first_name||' '||e.last_name,''), c.work_date,
		 c.minutes, c.reason, c.status, c.requested_at, c.approved_by, c.approved_at, c.rejection_reason
		 FROM compensation_requests c
		 LEFT JOIN employees e ON c.employee_id=e.id
		 WHERE c.company_id=$1 AND c.id=$2`, companyID, id,
	).Scan(&comp.ID, &comp.CompanyID, &comp.EmployeeID, &comp.EmployeeName, &comp.WorkDate,
		&comp.Minutes, &comp.Reason, &comp.Status, &comp.RequestedAt,
		&comp.ApprovedBy, &comp.ApprovedAt, &comp.RejectionReason)
	return comp, err
}

func (r *Repository) UpdateCompensationRequest(ctx context.Context, comp *CompensationRequest) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE compensation_requests SET status=$3, approved_by=$4, approved_at=$5, rejection_reason=$6 WHERE id=$1 AND company_id=$2`,
		comp.ID, comp.CompanyID, comp.Status, comp.ApprovedBy, comp.ApprovedAt, comp.RejectionReason)
	return err
}

func (r *Repository) ListCompensationRequests(ctx context.Context, companyID string, filters OvertimeFilters) ([]CompensationRequest, error) {
	query := `SELECT c.id, c.company_id, c.employee_id, COALESCE(e.first_name||' '||e.last_name,''), c.work_date,
		 c.minutes, c.reason, c.status, c.requested_at, c.approved_by, c.approved_at, c.rejection_reason
		 FROM compensation_requests c
		 LEFT JOIN employees e ON c.employee_id=e.id
		 WHERE c.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND c.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND c.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}

	query += " ORDER BY c.requested_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var comps []CompensationRequest
	for rows.Next() {
		var comp CompensationRequest
		rows.Scan(&comp.ID, &comp.CompanyID, &comp.EmployeeID, &comp.EmployeeName, &comp.WorkDate,
			&comp.Minutes, &comp.Reason, &comp.Status, &comp.RequestedAt,
			&comp.ApprovedBy, &comp.ApprovedAt, &comp.RejectionReason)
		comps = append(comps, comp)
	}
	return comps, nil
}

// Time Balance
func (r *Repository) GetTimeBalance(ctx context.Context, companyID, employeeID string) (*EmployeeTimeBalance, error) {
	b := &EmployeeTimeBalance{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, balance_minutes, updated_at
		 FROM employee_time_balances WHERE company_id=$1 AND employee_id=$2`, companyID, employeeID,
	).Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BalanceMinutes, &b.UpdatedAt)
	return b, err
}

func (r *Repository) CreateTimeBalance(ctx context.Context, b *EmployeeTimeBalance) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO employee_time_balances (id, company_id, employee_id, balance_minutes) VALUES ($1,$2,$3,$4)
		 ON CONFLICT (company_id, employee_id) DO UPDATE SET balance_minutes=$4, updated_at=NOW()`,
		b.ID, b.CompanyID, b.EmployeeID, b.BalanceMinutes)
	return err
}

func (r *Repository) UpdateTimeBalance(ctx context.Context, b *EmployeeTimeBalance) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE employee_time_balances SET balance_minutes=$3, updated_at=NOW() WHERE id=$1 AND company_id=$2`,
		b.ID, b.CompanyID, b.BalanceMinutes)
	return err
}

func (r *Repository) CreditTimeBalance(ctx context.Context, companyID, employeeID string, minutes int, overtimeRecordID, reason, createdBy string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO employee_time_balances (id, company_id, employee_id, balance_minutes)
		 VALUES (gen_random_uuid(), $1, $2, $3)
		 ON CONFLICT (company_id, employee_id) DO UPDATE SET balance_minutes=employee_time_balances.balance_minutes+$3, updated_at=NOW()`,
		companyID, employeeID, minutes)
	if err != nil { return err }

	_, err = tx.Exec(ctx,
		`INSERT INTO time_balance_transactions (id, company_id, employee_id, overtime_record_id, transaction_type, minutes, reason, created_by)
		 VALUES (gen_random_uuid(), $1, $2, $3, 'OVERTIME_CREDIT', $4, $5, $6)`,
		companyID, employeeID, overtimeRecordID, minutes, reason, createdBy)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *Repository) DebitTimeBalance(ctx context.Context, companyID, employeeID string, minutes int, reason, createdBy string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	var current int
	err = tx.QueryRow(ctx,
		`SELECT balance_minutes FROM employee_time_balances WHERE company_id=$1 AND employee_id=$2 FOR UPDATE`,
		companyID, employeeID).Scan(&current)
	if err != nil { return err }
	if current < minutes { return fmt.Errorf("insufficient balance") }

	_, err = tx.Exec(ctx,
		`UPDATE employee_time_balances SET balance_minutes=balance_minutes-$3, updated_at=NOW() WHERE company_id=$1 AND employee_id=$2`,
		companyID, employeeID, minutes)
	if err != nil { return err }

	_, err = tx.Exec(ctx,
		`INSERT INTO time_balance_transactions (id, company_id, employee_id, transaction_type, minutes, reason, created_by)
		 VALUES (gen_random_uuid(), $1, $2, 'COMPENSATION_DEBIT', -$3, $4, $5)`,
		companyID, employeeID, minutes, reason, createdBy)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *Repository) CreateBalanceTransaction(ctx context.Context, tx *TimeBalanceTransaction) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO time_balance_transactions (id, company_id, employee_id, overtime_record_id, transaction_type, minutes, reason, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		tx.ID, tx.CompanyID, tx.EmployeeID, tx.OvertimeRecordID, tx.TransactionType, tx.Minutes, tx.Reason, tx.CreatedBy)
	return err
}

func (r *Repository) ListBalanceTransactions(ctx context.Context, companyID, employeeID string) ([]TimeBalanceTransaction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, overtime_record_id, transaction_type, minutes, reason, created_by, created_at
		 FROM time_balance_transactions WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC`, companyID, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	var txs []TimeBalanceTransaction
	for rows.Next() {
		var tx TimeBalanceTransaction
		rows.Scan(&tx.ID, &tx.CompanyID, &tx.EmployeeID, &tx.OvertimeRecordID, &tx.TransactionType, &tx.Minutes, &tx.Reason, &tx.CreatedBy, &tx.CreatedAt)
		txs = append(txs, tx)
	}
	return txs, nil
}

// Helpers
func (r *Repository) GetAttendanceRecord(ctx context.Context, companyID, employeeID string, workDate time.Time) (*AttendanceRecord, error) {
	rec := &AttendanceRecord{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, work_date, scheduled_start, scheduled_end, actual_start, actual_end,
		 scheduled_minutes, worked_minutes, late_minutes, effective_late_minutes, early_leave_minutes, overtime_minutes, break_minutes, status
		 FROM attendance_records WHERE company_id=$1 AND employee_id=$2 AND work_date=$3`, companyID, employeeID, workDate,
	).Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.WorkDate, &rec.ScheduledStart, &rec.ScheduledEnd,
		&rec.ActualStart, &rec.ActualEnd, &rec.ScheduledMinutes, &rec.WorkedMinutes, &rec.LateMinutes,
		&rec.EffectiveLateMinutes, &rec.EarlyLeaveMinutes, &rec.OvertimeMinutes, &rec.BreakMinutes, &rec.Status)
	return rec, err
}

type AttendanceRecord struct {
	ID                   string     `json:"id"`
	CompanyID            string     `json:"company_id"`
	EmployeeID           string     `json:"employee_id"`
	WorkDate             time.Time  `json:"work_date"`
	ScheduledStart       *time.Time `json:"scheduled_start,omitempty"`
	ScheduledEnd         *time.Time `json:"scheduled_end,omitempty"`
	ActualStart          *time.Time `json:"actual_start,omitempty"`
	ActualEnd            *time.Time `json:"actual_end,omitempty"`
	ScheduledMinutes     int        `json:"scheduled_minutes"`
	WorkedMinutes        int        `json:"worked_minutes"`
	LateMinutes          int        `json:"late_minutes"`
	EffectiveLateMinutes int        `json:"effective_late_minutes"`
	EarlyLeaveMinutes    int        `json:"early_leave_minutes"`
	OvertimeMinutes      int        `json:"overtime_minutes"`
	BreakMinutes         int        `json:"break_minutes"`
	Status               string     `json:"status"`
}

func (r *Repository) GetHoliday(ctx context.Context, companyID string, date time.Time) (string, error) {
	var name string
	err := r.pool.QueryRow(ctx,
		`SELECT name FROM holidays WHERE company_id=$1 AND date=$2 LIMIT 1`, companyID, date,
	).Scan(&name)
	return name, err
}

func (r *Repository) GetActiveEmployees(ctx context.Context, companyID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND status='active'`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) GetEmployeeIDFromUser(ctx context.Context, companyID, userID string) (string, error) {
	var empID string
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND user_id=$2`, companyID, userID,
	).Scan(&empID)
	return empID, err
}
