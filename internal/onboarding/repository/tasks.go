package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/onboarding/domain"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, t *domain.OnboardingTask) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_tasks
		 (company_id, onboarding_id, employee_id, title, description, category,
		  responsible_type, responsible_id, due_date, status, required, sort_order, estimated_minutes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING id, created_at, updated_at`,
		t.CompanyID, t.TemplateID, "", t.Title, t.Description, t.TaskType,
		"", "", 0, "PENDING", t.Required, t.OrderIndex, nil,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepo) CreateAssignment(ctx context.Context, a *domain.OnboardingTaskAssignment) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_task_assignments
		 (onboarding_id, task_id, assigned_to, assigned_role, status, due_date, comments)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`,
		a.OnboardingID, a.TaskID, a.AssignedTo, a.AssignedRole, a.Status, a.DueDate, a.Comments,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *TaskRepo) GetAssignment(ctx context.Context, companyID, id string) (*domain.OnboardingTaskAssignment, error) {
	a := &domain.OnboardingTaskAssignment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, task_id, assigned_to, assigned_role, status, due_date,
		        completed_at, completed_by, comments, created_at, updated_at
		 FROM onboarding_task_assignments WHERE id=$1`, id,
	).Scan(&a.ID, &a.OnboardingID, &a.TaskID, &a.AssignedTo, &a.AssignedRole,
		&a.Status, &a.DueDate, &a.CompletedAt, &a.CompletedBy, &a.Comments,
		&a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *TaskRepo) ListAssignments(ctx context.Context, onboardingID string) ([]domain.OnboardingTaskAssignment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, task_id, assigned_to, assigned_role, status, due_date,
		        completed_at, completed_by, comments, created_at, updated_at
		 FROM onboarding_task_assignments WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []domain.OnboardingTaskAssignment
	for rows.Next() {
		var a domain.OnboardingTaskAssignment
		if err := rows.Scan(&a.ID, &a.OnboardingID, &a.TaskID, &a.AssignedTo, &a.AssignedRole,
			&a.Status, &a.DueDate, &a.CompletedAt, &a.CompletedBy, &a.Comments,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func (r *TaskRepo) UpdateAssignmentStatus(ctx context.Context, id string, status domain.TaskStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_task_assignments SET status=$2, updated_at=NOW() WHERE id=$1`,
		id, status)
	return err
}

func (r *TaskRepo) CompleteAssignment(ctx context.Context, id, completedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_task_assignments SET status='COMPLETED', completed_at=NOW(), completed_by=$2, updated_at=NOW()
		 WHERE id=$1`, id, completedBy)
	return err
}

func (r *TaskRepo) BlockAssignment(ctx context.Context, id, comments string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_task_assignments SET status='BLOCKED', comments=$2, updated_at=NOW()
		 WHERE id=$1`, id, comments)
	return err
}

func (r *TaskRepo) CreateDependency(ctx context.Context, d *domain.OnboardingTaskDependency) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_task_dependencies (task_id, depends_on_task_id)
		 VALUES ($1,$2) RETURNING id`, d.TaskID, d.DependsOnTaskID,
	).Scan(&d.ID)
}

func (r *TaskRepo) GetDependencies(ctx context.Context, taskID string) ([]domain.OnboardingTaskDependency, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, task_id, depends_on_task_id FROM onboarding_task_dependencies WHERE task_id=$1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ds []domain.OnboardingTaskDependency
	for rows.Next() {
		var d domain.OnboardingTaskDependency
		if err := rows.Scan(&d.ID, &d.TaskID, &d.DependsOnTaskID); err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *TaskRepo) AreDependenciesMet(ctx context.Context, taskID string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_task_dependencies d
		 JOIN onboarding_task_assignments a ON d.depends_on_task_id = a.task_id
		 WHERE d.task_id=$1 AND a.status != 'COMPLETED'`, taskID).Scan(&count)
	return count == 0, err
}

func (r *TaskRepo) GetCounts(ctx context.Context, onboardingID string) (total, completed int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(*) FILTER (WHERE status='COMPLETED')
		 FROM onboarding_task_assignments WHERE onboarding_id=$1`, onboardingID,
	).Scan(&total, &completed)
	return
}

func (r *TaskRepo) GetByResponsible(ctx context.Context, companyID, responsibleID string, status string) ([]domain.OnboardingTaskAssignment, error) {
	query := `SELECT a.id, a.onboarding_id, a.task_id, a.assigned_to, a.assigned_role,
	                 a.status, a.due_date, a.completed_at, a.completed_by, a.comments, a.created_at, a.updated_at
	          FROM onboarding_task_assignments a
	          JOIN onboarding_processes p ON a.onboarding_id = p.id
	          WHERE p.company_id=$1 AND a.assigned_to=$2`
	args := []interface{}{companyID, responsibleID}
	if status != "" {
		query += " AND a.status=$3"
		args = append(args, status)
	}
	query += " ORDER BY a.due_date ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []domain.OnboardingTaskAssignment
	for rows.Next() {
		var a domain.OnboardingTaskAssignment
		if err := rows.Scan(&a.ID, &a.OnboardingID, &a.TaskID, &a.AssignedTo, &a.AssignedRole,
			&a.Status, &a.DueDate, &a.CompletedAt, &a.CompletedBy, &a.Comments,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func (r *TaskRepo) CountOverdue(ctx context.Context, companyID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_task_assignments a
		 JOIN onboarding_processes p ON a.onboarding_id = p.id
		 WHERE p.company_id=$1 AND a.due_date < CURRENT_DATE AND a.status NOT IN ('COMPLETED','CANCELLED')`,
		companyID).Scan(&count)
	return count, err
}

func (r *TaskRepo) CountByStatus(ctx context.Context, companyID string) (dueToday int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_tasks
		 WHERE company_id=$1 AND due_date=CURRENT_DATE AND status NOT IN ('COMPLETED','CANCELLED')`,
		companyID).Scan(&dueToday)
	return
}

func fmtConditions(conditions map[string]string) string {
	s := ""
	for k, v := range conditions {
		if s != "" {
			s += " AND "
		}
		s += fmt.Sprintf("%s='%s'", k, v)
	}
	return s
}
