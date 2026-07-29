package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type feedbackRepo struct {
	pool *pgxpool.Pool
}

func (r *feedbackRepo) Create(ctx context.Context, f *domain.PerformanceFeedback) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_feedback (company_id, cycle_id, employee_id, author_id, feedback_type, visibility, content, is_anonymous, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		f.CompanyID, f.CycleID, f.EmployeeID, f.AuthorID, f.FeedbackType, f.Visibility, f.Content, f.IsAnonymous, now, now,
	).Scan(&f.ID)
}

func (r *feedbackRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceFeedback, error) {
	f := &domain.PerformanceFeedback{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, author_id, feedback_type, visibility, content, is_anonymous, created_at, updated_at
		 FROM performance_feedback WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&f.ID, &f.CompanyID, &f.CycleID, &f.EmployeeID, &f.AuthorID, &f.FeedbackType, &f.Visibility, &f.Content, &f.IsAnonymous, &f.CreatedAt, &f.UpdatedAt)
	return f, err
}

func (r *feedbackRepo) List(ctx context.Context, filter domain.FeedbackFilter) ([]domain.PerformanceFeedback, error) {
	query := `SELECT id, company_id, cycle_id, employee_id, author_id, feedback_type, visibility, content, is_anonymous, created_at, updated_at
		 FROM performance_feedback WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filter.EmployeeID)
		argIdx++
	}
	if filter.AuthorID != "" {
		query += fmt.Sprintf(" AND author_id=$%d", argIdx)
		args = append(args, filter.AuthorID)
		argIdx++
	}
	if filter.FeedbackType != "" {
		query += fmt.Sprintf(" AND feedback_type=$%d", argIdx)
		args = append(args, filter.FeedbackType)
		argIdx++
	}
	if filter.Visibility != "" {
		query += fmt.Sprintf(" AND visibility=$%d", argIdx)
		args = append(args, filter.Visibility)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []domain.PerformanceFeedback
	for rows.Next() {
		var f domain.PerformanceFeedback
		rows.Scan(&f.ID, &f.CompanyID, &f.CycleID, &f.EmployeeID, &f.AuthorID, &f.FeedbackType, &f.Visibility, &f.Content, &f.IsAnonymous, &f.CreatedAt, &f.UpdatedAt)
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *feedbackRepo) Update(ctx context.Context, f *domain.PerformanceFeedback) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_feedback SET feedback_type=$3, visibility=$4, content=$5, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		f.CompanyID, f.ID, f.FeedbackType, f.Visibility, f.Content)
	return err
}

func (r *feedbackRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_feedback WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *feedbackRepo) CreateRecognition(ctx context.Context, rec *domain.PerformanceRecognition) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_recognitions (company_id, employee_id, author_id, recognition_type, message, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		rec.CompanyID, rec.EmployeeID, rec.AuthorID, rec.RecognitionType, rec.Message, now,
	).Scan(&rec.ID)
}

func (r *feedbackRepo) ListRecognitionsByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceRecognition, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, author_id, recognition_type, message, created_at
		 FROM performance_recognitions WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC`,
		companyID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recognitions []domain.PerformanceRecognition
	for rows.Next() {
		var rec domain.PerformanceRecognition
		rows.Scan(&rec.ID, &rec.CompanyID, &rec.EmployeeID, &rec.AuthorID, &rec.RecognitionType, &rec.Message, &rec.CreatedAt)
		recognitions = append(recognitions, rec)
	}
	return recognitions, nil
}

// CheckInRepository

type checkInRepo struct {
	pool *pgxpool.Pool
}

func (r *checkInRepo) Create(ctx context.Context, ci *domain.PerformanceCheckIn) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_checkins (company_id, employee_id, manager_id, cycle_id, scheduled_at, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		ci.CompanyID, ci.EmployeeID, ci.ManagerID, ci.CycleID, ci.ScheduledAt, ci.Status, now, now,
	).Scan(&ci.ID)
}

func (r *checkInRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceCheckIn, error) {
	ci := &domain.PerformanceCheckIn{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, manager_id, cycle_id, scheduled_at, completed_at,
		 employee_notes, manager_notes, achievements, blockers, next_steps, status, created_at, updated_at
		 FROM performance_checkins WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&ci.ID, &ci.CompanyID, &ci.EmployeeID, &ci.ManagerID, &ci.CycleID, &ci.ScheduledAt, &ci.CompletedAt,
		&ci.EmployeeNotes, &ci.ManagerNotes, &ci.Achievements, &ci.Blockers, &ci.NextSteps, &ci.Status, &ci.CreatedAt, &ci.UpdatedAt)
	return ci, err
}

func (r *checkInRepo) ListByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceCheckIn, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, manager_id, cycle_id, scheduled_at, completed_at,
		 employee_notes, manager_notes, achievements, blockers, next_steps, status, created_at, updated_at
		 FROM performance_checkins WHERE company_id=$1 AND employee_id=$2 ORDER BY scheduled_at DESC`,
		companyID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkins []domain.PerformanceCheckIn
	for rows.Next() {
		var ci domain.PerformanceCheckIn
		rows.Scan(&ci.ID, &ci.CompanyID, &ci.EmployeeID, &ci.ManagerID, &ci.CycleID, &ci.ScheduledAt, &ci.CompletedAt,
			&ci.EmployeeNotes, &ci.ManagerNotes, &ci.Achievements, &ci.Blockers, &ci.NextSteps, &ci.Status, &ci.CreatedAt, &ci.UpdatedAt)
		checkins = append(checkins, ci)
	}
	return checkins, nil
}

func (r *checkInRepo) ListByManager(ctx context.Context, companyID, managerID string) ([]domain.PerformanceCheckIn, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, manager_id, cycle_id, scheduled_at, completed_at,
		 employee_notes, manager_notes, achievements, blockers, next_steps, status, created_at, updated_at
		 FROM performance_checkins WHERE company_id=$1 AND manager_id=$2 ORDER BY scheduled_at DESC`,
		companyID, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkins []domain.PerformanceCheckIn
	for rows.Next() {
		var ci domain.PerformanceCheckIn
		rows.Scan(&ci.ID, &ci.CompanyID, &ci.EmployeeID, &ci.ManagerID, &ci.CycleID, &ci.ScheduledAt, &ci.CompletedAt,
			&ci.EmployeeNotes, &ci.ManagerNotes, &ci.Achievements, &ci.Blockers, &ci.NextSteps, &ci.Status, &ci.CreatedAt, &ci.UpdatedAt)
		checkins = append(checkins, ci)
	}
	return checkins, nil
}

func (r *checkInRepo) Update(ctx context.Context, ci *domain.PerformanceCheckIn) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_checkins SET employee_notes=$3, manager_notes=$4, achievements=$5, blockers=$6,
		 next_steps=$7, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		ci.CompanyID, ci.ID, ci.EmployeeNotes, ci.ManagerNotes, ci.Achievements, ci.Blockers, ci.NextSteps)
	return err
}

func (r *checkInRepo) Complete(ctx context.Context, companyID, id string, notes map[string]*string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_checkins SET status='COMPLETED', completed_at=NOW(),
		 employee_notes=COALESCE($3,employee_notes), manager_notes=COALESCE($4,manager_notes),
		 achievements=COALESCE($5,achievements), blockers=COALESCE($6,blockers), next_steps=COALESCE($7,next_steps),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, notes["employee_notes"], notes["manager_notes"], notes["achievements"],
		notes["blockers"], notes["next_steps"])
	return err
}

func NewFeedbackRepository(pool *pgxpool.Pool) FeedbackRepository {
	return &feedbackRepo{pool: pool}
}

func NewCheckInRepository(pool *pgxpool.Pool) CheckInRepository {
	return &checkInRepo{pool: pool}
}
