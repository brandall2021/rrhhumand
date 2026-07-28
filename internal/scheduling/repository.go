package scheduling

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateSchedule(ctx context.Context, companyID string, req *CreateScheduleRequest) (*WorkSchedule, error) {
	s := &WorkSchedule{}
	tz := "UTC"
	weeklyHours := 0
	if req.Timezone != nil {
		tz = *req.Timezone
	}
	if req.WeeklyHours != nil {
		weeklyHours = *req.WeeklyHours
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO work_schedules (company_id, name, description, schedule_type, timezone, weekly_hours)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, name, description, schedule_type, timezone, weekly_hours, status, created_at, updated_at`,
		companyID, req.Name, req.Description, req.ScheduleType, tz, weeklyHours,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.ScheduleType, &s.Timezone, &s.WeeklyHours, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, day := range req.Days {
		isWorking := true
		if day.IsWorkingDay != nil {
			isWorking = *day.IsWorkingDay
		}
		breakMin := 0
		if day.BreakMinutes != nil {
			breakMin = *day.BreakMinutes
		}
		var dayID string
		err := r.pool.QueryRow(ctx,
			`INSERT INTO work_schedule_days (schedule_id, weekday, is_working_day, start_time, end_time, break_minutes)
			 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
			s.ID, day.Weekday, isWorking, day.StartTime, day.EndTime, breakMin,
		).Scan(&dayID)
		if err != nil {
			return nil, err
		}

		for _, interval := range day.Intervals {
			intervalType := "WORK"
			if interval.IntervalType != "" {
				intervalType = interval.IntervalType
			}
			seq := 1
			if interval.Sequence != nil {
				seq = *interval.Sequence
			}
			r.pool.Exec(ctx,
				`INSERT INTO work_schedule_intervals (schedule_day_id, start_time, end_time, interval_type, sequence)
				 VALUES ($1,$2,$3,$4,$5)`, dayID, interval.StartTime, interval.EndTime, intervalType, seq)
		}
	}

	return s, nil
}

func (r *Repository) GetSchedule(ctx context.Context, companyID, id string) (*WorkSchedule, error) {
	s := &WorkSchedule{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, schedule_type, timezone, weekly_hours, status, created_at, updated_at
		 FROM work_schedules WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.ScheduleType, &s.Timezone, &s.WeeklyHours, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	s.Days, _ = r.GetScheduleDays(ctx, s.ID)
	return s, nil
}

func (r *Repository) ListSchedules(ctx context.Context, companyID string) ([]WorkSchedule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, schedule_type, timezone, weekly_hours, status, created_at, updated_at
		 FROM work_schedules WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []WorkSchedule
	for rows.Next() {
		var s WorkSchedule
		if err := rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.ScheduleType, &s.Timezone, &s.WeeklyHours, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *Repository) UpdateSchedule(ctx context.Context, companyID, id string, req *UpdateScheduleRequest) (*WorkSchedule, error) {
	s := &WorkSchedule{}
	err := r.pool.QueryRow(ctx,
		`UPDATE work_schedules SET
		 name=COALESCE($3,name), description=COALESCE($4,description), timezone=COALESCE($5,timezone),
		 weekly_hours=COALESCE($6,weekly_hours), status=COALESCE($7,status), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, schedule_type, timezone, weekly_hours, status, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.Timezone, req.WeeklyHours, req.Status,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.ScheduleType, &s.Timezone, &s.WeeklyHours, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *Repository) DeleteSchedule(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM work_schedules WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *Repository) GetScheduleDays(ctx context.Context, scheduleID string) ([]WorkScheduleDay, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, schedule_id, weekday, is_working_day, start_time::TEXT, end_time::TEXT, break_minutes, created_at
		 FROM work_schedule_days WHERE schedule_id=$1 ORDER BY weekday`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []WorkScheduleDay
	for rows.Next() {
		var d WorkScheduleDay
		if err := rows.Scan(&d.ID, &d.ScheduleID, &d.Weekday, &d.IsWorkingDay, &d.StartTime, &d.EndTime, &d.BreakMinutes, &d.CreatedAt); err != nil {
			return nil, err
		}
		d.Intervals, _ = r.GetIntervals(ctx, d.ID)
		days = append(days, d)
	}
	return days, nil
}

func (r *Repository) GetIntervals(ctx context.Context, dayID string) ([]WorkScheduleInterval, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, schedule_day_id, start_time::TEXT, end_time::TEXT, interval_type, sequence
		 FROM work_schedule_intervals WHERE schedule_day_id=$1 ORDER BY sequence`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var intervals []WorkScheduleInterval
	for rows.Next() {
		var i WorkScheduleInterval
		if err := rows.Scan(&i.ID, &i.ScheduleDayID, &i.StartTime, &i.EndTime, &i.IntervalType, &i.Sequence); err != nil {
			return nil, err
		}
		intervals = append(intervals, i)
	}
	return intervals, nil
}

func (r *Repository) CreateShift(ctx context.Context, companyID string, req *CreateShiftRequest) (*Shift, error) {
	s := &Shift{}
	crossMidnight := false
	breakMin := 0
	if req.CrossesMidnight != nil {
		crossMidnight = *req.CrossesMidnight
	}
	if req.BreakMinutes != nil {
		breakMin = *req.BreakMinutes
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO shifts (company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color, status, created_at, updated_at`,
		companyID, req.Name, req.Code, req.StartTime, req.EndTime, crossMidnight, breakMin, req.Color,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Code, &s.StartTime, &s.EndTime, &s.CrossesMidnight, &s.BreakMinutes, &s.Color, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *Repository) GetShift(ctx context.Context, companyID, id string) (*Shift, error) {
	s := &Shift{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color, status, created_at, updated_at
		 FROM shifts WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Code, &s.StartTime, &s.EndTime, &s.CrossesMidnight, &s.BreakMinutes, &s.Color, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *Repository) ListShifts(ctx context.Context, companyID string) ([]Shift, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color, status, created_at, updated_at
		 FROM shifts WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []Shift
	for rows.Next() {
		var s Shift
		if err := rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.Code, &s.StartTime, &s.EndTime, &s.CrossesMidnight, &s.BreakMinutes, &s.Color, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		shifts = append(shifts, s)
	}
	return shifts, nil
}

func (r *Repository) UpdateShift(ctx context.Context, companyID, id string, req *UpdateShiftRequest) (*Shift, error) {
	s := &Shift{}
	err := r.pool.QueryRow(ctx,
		`UPDATE shifts SET
		 name=COALESCE($3,name), code=COALESCE($4,code), start_time=COALESCE($5,start_time),
		 end_time=COALESCE($6,end_time), crosses_midnight=COALESCE($7,crosses_midnight),
		 break_minutes=COALESCE($8,break_minutes), color=COALESCE($9,color), status=COALESCE($10,status), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color, status, created_at, updated_at`,
		companyID, id, req.Name, req.Code, req.StartTime, req.EndTime, req.CrossesMidnight, req.BreakMinutes, req.Color, req.Status,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Code, &s.StartTime, &s.EndTime, &s.CrossesMidnight, &s.BreakMinutes, &s.Color, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *Repository) DeleteShift(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM shifts WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *Repository) AssignSchedule(ctx context.Context, companyID, employeeID string, req *AssignScheduleRequest) (*EmployeeScheduleAssignment, error) {
	a := &EmployeeScheduleAssignment{}
	effectiveFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}

	// Deactivate previous active assignment
	r.pool.Exec(ctx,
		`UPDATE employee_schedule_assignments SET status='INACTIVE' WHERE employee_id=$1 AND status='ACTIVE'`, employeeID)

	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_schedule_assignments (company_id, employee_id, schedule_id, effective_from, effective_to, priority)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, employee_id, schedule_id, effective_from, effective_to, priority, status, created_at`,
		companyID, employeeID, req.ScheduleID, effectiveFrom, req.EffectiveTo, priority,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.ScheduleID, &a.EffectiveFrom, &a.EffectiveTo, &a.Priority, &a.Status, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetEmployeeScheduleAssignment(ctx context.Context, employeeID string, date time.Time) (*EmployeeScheduleAssignment, error) {
	a := &EmployeeScheduleAssignment{}
	err := r.pool.QueryRow(ctx,
		`SELECT esa.id, esa.company_id, esa.employee_id, ws.name, esa.schedule_id, esa.effective_from, esa.effective_to, esa.priority, esa.status, esa.created_at
		 FROM employee_schedule_assignments esa
		 JOIN work_schedules ws ON esa.schedule_id=ws.id
		 WHERE esa.employee_id=$1 AND esa.status='ACTIVE'
		 AND esa.effective_from<=$2 AND (esa.effective_to IS NULL OR esa.effective_to>=$2)
		 ORDER BY esa.priority DESC LIMIT 1`, employeeID, date,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.ScheduleName, &a.ScheduleID, &a.EffectiveFrom, &a.EffectiveTo, &a.Priority, &a.Status, &a.CreatedAt)
	return a, err
}

func (r *Repository) AssignShift(ctx context.Context, companyID, employeeID string, req *AssignShiftRequest) (*EmployeeShiftAssignment, error) {
	workDate, _ := time.Parse("2006-01-02", req.WorkDate)
	a := &EmployeeShiftAssignment{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_shift_assignments (company_id, employee_id, shift_id, work_date, notes)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (employee_id, work_date) DO UPDATE SET shift_id=$3, notes=$5
		 RETURNING id, company_id, employee_id, shift_id, work_date, status, notes, created_at`,
		companyID, employeeID, req.ShiftID, workDate, req.Notes,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.ShiftID, &a.WorkDate, &a.Status, &a.Notes, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetEmployeeShiftAssignment(ctx context.Context, employeeID string, date time.Time) (*EmployeeShiftAssignment, error) {
	a := &EmployeeShiftAssignment{}
	err := r.pool.QueryRow(ctx,
		`SELECT esa.id, esa.company_id, esa.employee_id, s.name, esa.shift_id, esa.work_date, esa.status, esa.notes, esa.created_at
		 FROM employee_shift_assignments esa
		 JOIN shifts s ON esa.shift_id=s.id
		 WHERE esa.employee_id=$1 AND esa.work_date=$2`, employeeID, date,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.ShiftName, &a.ShiftID, &a.WorkDate, &a.Status, &a.Notes, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetShiftByID(ctx context.Context, id string) (*Shift, error) {
	s := &Shift{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, code, start_time, end_time, crosses_midnight, break_minutes, color, status, created_at, updated_at
		 FROM shifts WHERE id=$1`, id,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Code, &s.StartTime, &s.EndTime, &s.CrossesMidnight, &s.BreakMinutes, &s.Color, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *Repository) CreateRotationTemplate(ctx context.Context, companyID string, req *CreateRotationTemplateRequest) (*RotationTemplate, error) {
	t := &RotationTemplate{}
	cycleLength := len(req.Days)
	err := r.pool.QueryRow(ctx,
		`INSERT INTO rotation_templates (company_id, name, description, cycle_length)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, name, description, cycle_length, status, created_at`,
		companyID, req.Name, req.Description, cycleLength,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.CycleLength, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, day := range req.Days {
		isRest := false
		if day.IsRestDay != nil {
			isRest = *day.IsRestDay
		}
		r.pool.Exec(ctx,
			`INSERT INTO rotation_template_days (template_id, day_position, shift_id, is_rest_day)
			 VALUES ($1,$2,$3,$4)`, t.ID, day.DayPosition, day.ShiftID, isRest)
	}
	return t, nil
}

func (r *Repository) GetRotationTemplate(ctx context.Context, companyID, id string) (*RotationTemplate, error) {
	t := &RotationTemplate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, cycle_length, status, created_at
		 FROM rotation_templates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.CycleLength, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, _ := r.pool.Query(ctx,
		`SELECT rtd.id, rtd.template_id, rtd.day_position, rtd.shift_id, s.name, rtd.is_rest_day
		 FROM rotation_template_days rtd
		 LEFT JOIN shifts s ON rtd.shift_id=s.id
		 WHERE rtd.template_id=$1 ORDER BY rtd.day_position`, t.ID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var d RotationTemplateDay
			rows.Scan(&d.ID, &d.TemplateID, &d.DayPosition, &d.ShiftID, &d.ShiftName, &d.IsRestDay)
			t.Days = append(t.Days, d)
		}
	}
	return t, nil
}

func (r *Repository) ListRotationTemplates(ctx context.Context, companyID string) ([]RotationTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, cycle_length, status, created_at
		 FROM rotation_templates WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []RotationTemplate
	for rows.Next() {
		var t RotationTemplate
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.CycleLength, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *Repository) AssignRotation(ctx context.Context, companyID, employeeID string, req *AssignRotationRequest) (*EmployeeRotationAssignment, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	cyclePos := 1
	if req.CyclePosition != nil {
		cyclePos = *req.CyclePosition
	}
	a := &EmployeeRotationAssignment{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_rotation_assignments (company_id, employee_id, template_id, start_date, cycle_position, effective_to)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, employee_id, template_id, start_date, cycle_position, effective_to, status, created_at`,
		companyID, employeeID, req.TemplateID, startDate, cyclePos, req.EffectiveTo,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TemplateID, &a.StartDate, &a.CyclePosition, &a.EffectiveTo, &a.Status, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetEmployeeRotation(ctx context.Context, employeeID string, date time.Time) (*EmployeeRotationAssignment, error) {
	a := &EmployeeRotationAssignment{}
	err := r.pool.QueryRow(ctx,
		`SELECT era.id, era.company_id, era.employee_id, era.template_id, rt.name, era.start_date, era.cycle_position, era.effective_to, era.status, era.created_at
		 FROM employee_rotation_assignments era
		 JOIN rotation_templates rt ON era.template_id=rt.id
		 WHERE era.employee_id=$1 AND era.status='ACTIVE'
		 AND era.start_date<=$2 AND (era.effective_to IS NULL OR era.effective_to>=$2)
		 LIMIT 1`, employeeID, date,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.TemplateID, &a.TemplateName, &a.StartDate, &a.CyclePosition, &a.EffectiveTo, &a.Status, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetCalendarEntry(ctx context.Context, employeeID string, date time.Time) (*EmployeeWorkCalendar, error) {
	c := &EmployeeWorkCalendar{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, work_date, schedule_id, shift_id, planned_start, planned_end, planned_break_minutes, status, source, created_at
		 FROM employee_work_calendar WHERE employee_id=$1 AND work_date=$2`, employeeID, date,
	).Scan(&c.ID, &c.CompanyID, &c.EmployeeID, &c.WorkDate, &c.ScheduleID, &c.ShiftID, &c.PlannedStart, &c.PlannedEnd, &c.PlannedBreakMinutes, &c.Status, &c.Source, &c.CreatedAt)
	return c, err
}

func (r *Repository) UpsertCalendarEntry(ctx context.Context, entry *EmployeeWorkCalendar) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO employee_work_calendar (id, company_id, employee_id, work_date, schedule_id, shift_id, planned_start, planned_end, planned_break_minutes, status, source)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT (employee_id, work_date) DO UPDATE SET
		 schedule_id=$5, shift_id=$6, planned_start=$7, planned_end=$8, planned_break_minutes=$9, status=$10, source=$11`,
		entry.ID, entry.CompanyID, entry.EmployeeID, entry.WorkDate, entry.ScheduleID, entry.ShiftID,
		entry.PlannedStart, entry.PlannedEnd, entry.PlannedBreakMinutes, entry.Status, entry.Source)
	return err
}

func (r *Repository) ListCalendarEntries(ctx context.Context, companyID string, filters CalendarFilters, offset, limit int) ([]EmployeeWorkCalendar, int64, error) {
	query := `SELECT ewc.id, ewc.company_id, ewc.employee_id, ewc.work_date, ewc.schedule_id, ewc.shift_id, COALESCE(s.name,''), ewc.planned_start, ewc.planned_end, ewc.planned_break_minutes, ewc.status, ewc.source, ewc.created_at
		 FROM employee_work_calendar ewc
		 LEFT JOIN shifts s ON ewc.shift_id=s.id
		 WHERE ewc.company_id=$1`
	countQuery := `SELECT COUNT(*) FROM employee_work_calendar ewc WHERE ewc.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND ewc.employee_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ewc.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND ewc.work_date>=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ewc.work_date>=$%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND ewc.work_date<=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ewc.work_date<=$%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND ewc.status=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND ewc.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY ewc.work_date LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []EmployeeWorkCalendar
	for rows.Next() {
		var e EmployeeWorkCalendar
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.WorkDate, &e.ScheduleID, &e.ShiftID, &e.ShiftName, &e.PlannedStart, &e.PlannedEnd, &e.PlannedBreakMinutes, &e.Status, &e.Source, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}

func (r *Repository) CreateException(ctx context.Context, companyID string, req *CreateExceptionRequest) (*ScheduleException, error) {
	exceptionDate, _ := time.Parse("2006-01-02", req.ExceptionDate)
	e := &ScheduleException{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO schedule_exceptions (company_id, employee_id, exception_date, exception_type, start_time, end_time, shift_id, reason)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, employee_id, exception_date, exception_type, start_time, end_time, shift_id, reason, status, created_at`,
		companyID, req.EmployeeID, exceptionDate, req.ExceptionType, req.StartTime, req.EndTime, req.ShiftID, req.Reason,
	).Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.ExceptionDate, &e.ExceptionType, &e.StartTime, &e.EndTime, &e.ShiftID, &e.Reason, &e.Status, &e.CreatedAt)
	return e, err
}

func (r *Repository) GetException(ctx context.Context, companyID, employeeID string, date time.Time) (*ScheduleException, error) {
	e := &ScheduleException{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, exception_date, exception_type, start_time, end_time, shift_id, reason, status, created_at
		 FROM schedule_exceptions WHERE company_id=$1 AND employee_id=$2 AND exception_date=$3 AND status='APPROVED'
		 LIMIT 1`, companyID, employeeID, date,
	).Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.ExceptionDate, &e.ExceptionType, &e.StartTime, &e.EndTime, &e.ShiftID, &e.Reason, &e.Status, &e.CreatedAt)
	return e, err
}

func (r *Repository) ListExceptions(ctx context.Context, companyID string) ([]ScheduleException, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, exception_date, exception_type, start_time, end_time, shift_id, reason, status, created_at
		 FROM schedule_exceptions WHERE company_id=$1 ORDER BY exception_date`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exceptions []ScheduleException
	for rows.Next() {
		var e ScheduleException
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EmployeeID, &e.ExceptionDate, &e.ExceptionType, &e.StartTime, &e.EndTime, &e.ShiftID, &e.Reason, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		exceptions = append(exceptions, e)
	}
	return exceptions, nil
}

func (r *Repository) CreateSwap(ctx context.Context, companyID, requesterID string, req *SwapShiftRequest) (*ShiftSwap, error) {
	requesterDate, _ := time.Parse("2006-01-02", req.RequesterDate)
	targetDate, _ := time.Parse("2006-01-02", req.TargetDate)
	s := &ShiftSwap{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO shift_swaps (company_id, requester_employee_id, target_employee_id, requester_date, target_date, reason)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, requester_employee_id, target_employee_id, requester_date, target_date, reason, status, created_at`,
		companyID, requesterID, req.TargetEmployeeID, requesterDate, targetDate, req.Reason,
	).Scan(&s.ID, &s.CompanyID, &s.RequesterEmployeeID, &s.TargetEmployeeID, &s.RequesterDate, &s.TargetDate, &s.Reason, &s.Status, &s.CreatedAt)
	return s, err
}

func (r *Repository) UpdateSwapStatus(ctx context.Context, id, status, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE shift_swaps SET status=$1, approved_by=$2 WHERE id=$3`, status, approvedBy, id)
	return err
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}
