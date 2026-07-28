package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Punches struct {
	pool *pgxpool.Pool
}

func NewPunches(pool *pgxpool.Pool) *Punches {
	return &Punches{pool: pool}
}

func (p *Punches) CreatePunch(ctx context.Context, punch *AttendancePunch) (*AttendancePunch, error) {
	err := p.pool.QueryRow(ctx,
		`INSERT INTO attendance_punches (id, company_id, employee_id, attendance_id, punch_type, punched_at, source, latitude, longitude, ip_address, device_id, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 RETURNING id, company_id, employee_id, attendance_id, punch_type, punched_at, source, latitude, longitude, ip_address, device_id, notes, created_at`,
		punch.ID, punch.CompanyID, punch.EmployeeID, punch.AttendanceID, punch.PunchType, punch.PunchedAt, punch.Source,
		punch.Latitude, punch.Longitude, punch.IPAddress, punch.DeviceID, punch.Notes,
	).Scan(&punch.ID, &punch.CompanyID, &punch.EmployeeID, &punch.AttendanceID, &punch.PunchType, &punch.PunchedAt, &punch.Source,
		&punch.Latitude, &punch.Longitude, &punch.IPAddress, &punch.DeviceID, &punch.Notes, &punch.CreatedAt)
	return punch, err
}

func (p *Punches) GetPunchesForAttendance(ctx context.Context, attendanceID string) ([]AttendancePunch, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT id, company_id, employee_id, attendance_id, punch_type, punched_at, source, latitude, longitude, ip_address, device_id, notes, created_at
		 FROM attendance_punches WHERE attendance_id=$1 ORDER BY punched_at`, attendanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var punches []AttendancePunch
	for rows.Next() {
		var punch AttendancePunch
		if err := rows.Scan(&punch.ID, &punch.CompanyID, &punch.EmployeeID, &punch.AttendanceID, &punch.PunchType, &punch.PunchedAt, &punch.Source,
			&punch.Latitude, &punch.Longitude, &punch.IPAddress, &punch.DeviceID, &punch.Notes, &punch.CreatedAt); err != nil {
			return nil, err
		}
		punches = append(punches, punch)
	}
	return punches, nil
}

func (p *Punches) GetLastPunch(ctx context.Context, employeeID string, date time.Time) (*AttendancePunch, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	punch := &AttendancePunch{}
	err := p.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, attendance_id, punch_type, punched_at, source, latitude, longitude, ip_address, device_id, notes, created_at
		 FROM attendance_punches WHERE employee_id=$1 AND punched_at>=$2 AND punched_at<$3
		 ORDER BY punched_at DESC LIMIT 1`, employeeID, startOfDay, endOfDay,
	).Scan(&punch.ID, &punch.CompanyID, &punch.EmployeeID, &punch.AttendanceID, &punch.PunchType, &punch.PunchedAt, &punch.Source,
		&punch.Latitude, &punch.Longitude, &punch.IPAddress, &punch.DeviceID, &punch.Notes, &punch.CreatedAt)
	if err != nil {
		return nil, err
	}
	return punch, nil
}

func (p *Punches) HasClockInToday(ctx context.Context, employeeID string, date time.Time) (bool, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var exists bool
	err := p.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM attendance_punches WHERE employee_id=$1 AND punch_type='CLOCK_IN' AND punched_at>=$2 AND punched_at<$3)`,
		employeeID, startOfDay, endOfDay).Scan(&exists)
	return exists, err
}

func (p *Punches) GetBreakMinutes(ctx context.Context, employeeID string, date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	rows, err := p.pool.Query(ctx,
		`SELECT punched_at, punch_type FROM attendance_punches
		 WHERE employee_id=$1 AND punch_type IN ('BREAK_START','BREAK_END') AND punched_at>=$2 AND punched_at<$3
		 ORDER BY punched_at`, employeeID, startOfDay, endOfDay)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var breakStart time.Time
	totalBreak := 0
	for rows.Next() {
		var punchedAt time.Time
		var punchType string
		if err := rows.Scan(&punchedAt, &punchType); err != nil {
			return 0, err
		}
		if punchType == "BREAK_START" {
			breakStart = punchedAt
		} else if punchType == "BREAK_END" && !breakStart.IsZero() {
			totalBreak += int(punchedAt.Sub(breakStart).Minutes())
			breakStart = time.Time{}
		}
	}
	return totalBreak, nil
}

func (p *Punches) GetBreakMinutesBetween(ctx context.Context, employeeID string, start, end time.Time) (int, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT punched_at, punch_type FROM attendance_punches
		 WHERE employee_id=$1 AND punch_type IN ('BREAK_START','BREAK_END') AND punched_at>=$2 AND punched_at<=$3
		 ORDER BY punched_at`, employeeID, start, end)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var breakStart time.Time
	totalBreak := 0
	for rows.Next() {
		var punchedAt time.Time
		var punchType string
		if err := rows.Scan(&punchedAt, &punchType); err != nil {
			return 0, err
		}
		if punchType == "BREAK_START" {
			breakStart = punchedAt
		} else if punchType == "BREAK_END" && !breakStart.IsZero() {
			totalBreak += int(punchedAt.Sub(breakStart).Minutes())
			breakStart = time.Time{}
		}
	}
	return totalBreak, nil
}

func (p *Punches) GetOpenBreak(ctx context.Context, employeeID string, date time.Time) (bool, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var exists bool
	err := p.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM attendance_punches p1
			WHERE p1.employee_id=$1 AND p1.punch_type='BREAK_START' AND p1.punched_at>=$2 AND p1.punched_at<$3
			AND NOT EXISTS(
				SELECT 1 FROM attendance_punches p2
				WHERE p2.employee_id=$1 AND p2.punch_type='BREAK_END' AND p2.punched_at>p1.punched_at AND p2.punched_at<$3
			)
		)`, employeeID, startOfDay, endOfDay).Scan(&exists)
	return exists, err
}

func (p *Punches) ValidateSource(source string, policy *AttendancePolicy) error {
	switch source {
	case "MOBILE":
		if !policy.AllowMobile {
			return fmt.Errorf("mobile clock-in not allowed")
		}
	case "WEB":
		if !policy.AllowWeb {
			return fmt.Errorf("web clock-in not allowed")
		}
	case "KIOSK":
		if !policy.AllowKiosk {
			return fmt.Errorf("kiosk clock-in not allowed")
		}
	}
	return nil
}
