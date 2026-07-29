package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/onboarding/domain"
)

type OffboardingRepo struct {
	pool *pgxpool.Pool
}

func NewOffboardingRepo(pool *pgxpool.Pool) *OffboardingRepo {
	return &OffboardingRepo{pool: pool}
}

func (r *OffboardingRepo) CreateProcess(ctx context.Context, p *domain.OffboardingProcess) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO offboarding_processes
		 (company_id, employee_id, template_id, requested_by, termination_type, reason_id,
		  notice_date, last_working_date, termination_effective_date, status, employee_status_after)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 RETURNING id, created_at, updated_at`,
		p.CompanyID, p.EmployeeID, p.TemplateID, p.RequestedBy, p.TerminationType, p.ReasonID,
		p.NoticeDate, p.LastWorkingDate, p.TerminationEffectiveDate, p.Status, p.EmployeeStatusAfter,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *OffboardingRepo) GetProcessByID(ctx context.Context, companyID, id string) (*domain.OffboardingProcess, error) {
	p := &domain.OffboardingProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, template_id, requested_by, termination_type, reason_id,
		        notice_date, last_working_date, termination_effective_date, status, progress,
		        employee_status_after, created_at, updated_at, completed_at
		 FROM offboarding_processes WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.TemplateID, &p.RequestedBy, &p.TerminationType,
		&p.ReasonID, &p.NoticeDate, &p.LastWorkingDate, &p.TerminationEffectiveDate,
		&p.Status, &p.Progress, &p.EmployeeStatusAfter, &p.CreatedAt, &p.UpdatedAt, &p.CompletedAt)
	return p, err
}

func (r *OffboardingRepo) ListProcesses(ctx context.Context, companyID string, status, employeeID string) ([]domain.OffboardingProcess, error) {
	query := `SELECT id, company_id, employee_id, template_id, requested_by, termination_type, reason_id,
	                 notice_date, last_working_date, termination_effective_date, status, progress,
	                 employee_status_after, created_at, updated_at, completed_at
	          FROM offboarding_processes WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if employeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, employeeID)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ps []domain.OffboardingProcess
	for rows.Next() {
		var p domain.OffboardingProcess
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.TemplateID, &p.RequestedBy,
			&p.TerminationType, &p.ReasonID, &p.NoticeDate, &p.LastWorkingDate,
			&p.TerminationEffectiveDate, &p.Status, &p.Progress,
			&p.EmployeeStatusAfter, &p.CreatedAt, &p.UpdatedAt, &p.CompletedAt); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (r *OffboardingRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.OffboardingStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offboarding_processes SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *OffboardingRepo) UpdateProgress(ctx context.Context, companyID, id string, progress float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offboarding_processes SET progress=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, progress)
	return err
}

func (r *OffboardingRepo) Complete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offboarding_processes SET status='COMPLETED', progress=100, completed_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *OffboardingRepo) HasActiveProcess(ctx context.Context, companyID, employeeID string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM offboarding_processes
		 WHERE company_id=$1 AND employee_id=$2 AND status NOT IN ('COMPLETED','CANCELLED')`,
		companyID, employeeID).Scan(&count)
	return count > 0, err
}

func (r *OffboardingRepo) CreateTask(ctx context.Context, t *domain.OffboardingTask) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO offboarding_tasks (company_id, offboarding_id, title, description, task_type,
		 assigned_to, assigned_role, required, due_date, status, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`,
		t.CompanyID, t.OffboardingID, t.Title, t.Description, t.TaskType,
		t.AssignedTo, t.AssignedRole, t.Required, t.DueDate, t.Status, t.SortOrder,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *OffboardingRepo) GetTask(ctx context.Context, companyID, id string) (*domain.OffboardingTask, error) {
	t := &domain.OffboardingTask{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, offboarding_id, title, description, task_type, assigned_to,
		        assigned_role, required, due_date, status, completed_at, completed_by, comments,
		        sort_order, created_at, updated_at
		 FROM offboarding_tasks WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&t.ID, &t.CompanyID, &t.OffboardingID, &t.Title, &t.Description, &t.TaskType,
		&t.AssignedTo, &t.AssignedRole, &t.Required, &t.DueDate, &t.Status,
		&t.CompletedAt, &t.CompletedBy, &t.Comments, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *OffboardingRepo) ListTasks(ctx context.Context, offboardingID string) ([]domain.OffboardingTask, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, offboarding_id, title, description, task_type, assigned_to,
		        assigned_role, required, due_date, status, completed_at, completed_by, comments,
		        sort_order, created_at, updated_at
		 FROM offboarding_tasks WHERE offboarding_id=$1 ORDER BY sort_order`, offboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []domain.OffboardingTask
	for rows.Next() {
		var t domain.OffboardingTask
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.OffboardingID, &t.Title, &t.Description, &t.TaskType,
			&t.AssignedTo, &t.AssignedRole, &t.Required, &t.DueDate, &t.Status,
			&t.CompletedAt, &t.CompletedBy, &t.Comments, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (r *OffboardingRepo) CompleteTask(ctx context.Context, companyID, id, completedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offboarding_tasks SET status='COMPLETED', completed_at=NOW(), completed_by=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, completedBy)
	return err
}

func (r *OffboardingRepo) GetTaskCounts(ctx context.Context, offboardingID string) (total, completed int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(*) FILTER (WHERE status='COMPLETED')
		 FROM offboarding_tasks WHERE offboarding_id=$1`, offboardingID,
	).Scan(&total, &completed)
	return
}

func (r *OffboardingRepo) CreateAsset(ctx context.Context, a *domain.OffboardingAsset) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO offboarding_assets (company_id, offboarding_id, employee_id, asset_type, description,
		 serial_number, condition_on_delivery, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at, updated_at`,
		a.CompanyID, a.OffboardingID, a.EmployeeID, a.AssetType, a.Description,
		a.SerialNumber, a.ConditionOnDelivery, a.Status, a.Notes,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *OffboardingRepo) ListAssets(ctx context.Context, offboardingID string) ([]domain.OffboardingAsset, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, offboarding_id, employee_id, asset_type, description,
		        serial_number, condition_on_delivery, condition_on_return, status,
		        returned_at, returned_to, notes, created_at, updated_at
		 FROM offboarding_assets WHERE offboarding_id=$1 ORDER BY created_at`, offboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []domain.OffboardingAsset
	for rows.Next() {
		var a domain.OffboardingAsset
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.OffboardingID, &a.EmployeeID, &a.AssetType,
			&a.Description, &a.SerialNumber, &a.ConditionOnDelivery, &a.ConditionOnReturn,
			&a.Status, &a.ReturnedAt, &a.ReturnedTo, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func (r *OffboardingRepo) UpdateAssetStatus(ctx context.Context, companyID, id string, status domain.OffboardingAssetStatus, conditionOnReturn *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offboarding_assets SET status=$3, condition_on_return=$4, returned_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, status, conditionOnReturn)
	return err
}

func (r *OffboardingRepo) CreateAccessRevocation(ctx context.Context, a *domain.AccessRevocation) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO employee_access_revocations (company_id, employee_id, offboarding_id, system_name, access_type, status)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, requested_at, created_at, updated_at`,
		a.CompanyID, a.EmployeeID, a.OffboardingID, a.SystemName, a.AccessType, a.Status,
	).Scan(&a.ID, &a.RequestedAt, &a.CreatedAt, &a.UpdatedAt)
}

func (r *OffboardingRepo) ListAccessRevocations(ctx context.Context, offboardingID string) ([]domain.AccessRevocation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, offboarding_id, system_name, access_type,
		        requested_at, revoked_at, status, performed_by, error_message, created_at, updated_at
		 FROM employee_access_revocations WHERE offboarding_id=$1 ORDER BY created_at`, offboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []domain.AccessRevocation
	for rows.Next() {
		var a domain.AccessRevocation
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.OffboardingID, &a.SystemName,
			&a.AccessType, &a.RequestedAt, &a.RevokedAt, &a.Status, &a.PerformedBy,
			&a.ErrorMessage, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func (r *OffboardingRepo) UpdateAccessRevocation(ctx context.Context, companyID, id string, status string, performedBy *string, errMsg *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE employee_access_revocations SET status=$3, performed_by=$4, error_message=$5, revoked_at=CASE WHEN $3='REVOKED' THEN NOW() ELSE revoked_at END, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, status, performedBy, errMsg)
	return err
}

func (r *OffboardingRepo) CreateHandover(ctx context.Context, h *domain.EmployeeHandover) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO employee_handovers (company_id, employee_id, offboarding_id, handover_to, description, projects, pending_tasks, documents, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at, updated_at`,
		h.CompanyID, h.EmployeeID, h.OffboardingID, h.HandoverTo, h.Description,
		h.Projects, h.PendingTasks, h.Documents, h.Status,
	).Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)
}

func (r *OffboardingRepo) GetHandover(ctx context.Context, companyID, offboardingID string) (*domain.EmployeeHandover, error) {
	h := &domain.EmployeeHandover{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, offboarding_id, handover_to, description,
		        projects, pending_tasks, documents, status, completed_at, created_at, updated_at
		 FROM employee_handovers WHERE company_id=$1 AND offboarding_id=$2`, companyID, offboardingID,
	).Scan(&h.ID, &h.CompanyID, &h.EmployeeID, &h.OffboardingID, &h.HandoverTo, &h.Description,
		&h.Projects, &h.PendingTasks, &h.Documents, &h.Status, &h.CompletedAt, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

func (r *OffboardingRepo) CompleteHandover(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE employee_handovers SET status='COMPLETED', completed_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *OffboardingRepo) CreateExitInterview(ctx context.Context, e *domain.ExitInterview) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO exit_interviews (company_id, offboarding_id, employee_id, interviewer_id, scheduled_at, anonymous)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		e.CompanyID, e.OffboardingID, e.EmployeeID, e.InterviewerID, e.ScheduledAt, e.Anonymous,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *OffboardingRepo) GetExitInterview(ctx context.Context, companyID, offboardingID string) (*domain.ExitInterview, error) {
	e := &domain.ExitInterview{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, offboarding_id, employee_id, interviewer_id, scheduled_at,
		        completed_at, reason, feedback, recommendation, rating, anonymous, created_at, updated_at
		 FROM exit_interviews WHERE company_id=$1 AND offboarding_id=$2`, companyID, offboardingID,
	).Scan(&e.ID, &e.CompanyID, &e.OffboardingID, &e.EmployeeID, &e.InterviewerID,
		&e.ScheduledAt, &e.CompletedAt, &e.Reason, &e.Feedback, &e.Recommendation,
		&e.Rating, &e.Anonymous, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *OffboardingRepo) CompleteExitInterview(ctx context.Context, companyID, offboardingID string, reason, feedback string, recommendation *string, rating *float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE exit_interviews SET completed_at=NOW(), reason=$3, feedback=$4, recommendation=$5, rating=$6, updated_at=NOW()
		 WHERE company_id=$1 AND offboarding_id=$2`, companyID, offboardingID, reason, feedback, recommendation, rating)
	return err
}

func (r *OffboardingRepo) ListExitReasons(ctx context.Context, companyID string) ([]domain.TerminationReason, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, active FROM employee_exit_reasons WHERE company_id=$1 AND active=true ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rs []domain.TerminationReason
	for rows.Next() {
		var r domain.TerminationReason
		if err := rows.Scan(&r.ID, &r.CompanyID, &r.Name, &r.Description, &r.Active); err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (r *OffboardingRepo) GetDashboardStats(ctx context.Context, companyID string) (active, pending, completed, overdue int, err error) {
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM offboarding_processes WHERE company_id=$1 AND status='IN_PROGRESS'`, companyID).Scan(&active)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM offboarding_processes WHERE company_id=$1 AND status='DRAFT'`, companyID).Scan(&pending)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM offboarding_processes WHERE company_id=$1 AND status='COMPLETED'`, companyID).Scan(&completed)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM offboarding_processes WHERE company_id=$1 AND status IN ('IN_PROGRESS','DRAFT') AND last_working_date < CURRENT_DATE`, companyID).Scan(&overdue)
	return
}
