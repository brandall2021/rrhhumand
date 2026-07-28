package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type AssessmentRepo struct {
	pool *pgxpool.Pool
}

func NewAssessmentRepo(pool *pgxpool.Pool) *AssessmentRepo {
	return &AssessmentRepo{pool: pool}
}

func (r *AssessmentRepo) Create(ctx context.Context, companyID string, req *domain.Assessment) (*domain.Assessment, error) {
	a := &domain.Assessment{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO assessments (company_id, application_id, assessment_type, title, description, max_score, passing_score, duration_minutes, due_at, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, company_id, application_id, assessment_type, title, description, max_score, passing_score, duration_minutes, due_at, status, score, result, result_summary, completed_at, created_by, created_at, updated_at`,
		companyID, req.ApplicationID, req.AssessmentType, req.Title, req.Description,
		req.MaxScore, req.PassingScore, req.DurationMinutes, req.DueAt, req.CreatedBy,
	).Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description,
		&a.MaxScore, &a.PassingScore, &a.DurationMinutes, &a.DueAt, &a.Status, &a.Score,
		&a.Result, &a.ResultSummary, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *AssessmentRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Assessment, error) {
	a := &domain.Assessment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, application_id, assessment_type, title, description, max_score, passing_score, duration_minutes, due_at, status, score, result, result_summary, completed_at, created_by, created_at, updated_at
		 FROM assessments WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description,
		&a.MaxScore, &a.PassingScore, &a.DurationMinutes, &a.DueAt, &a.Status, &a.Score,
		&a.Result, &a.ResultSummary, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *AssessmentRepo) List(ctx context.Context, companyID string, applicationID, status string) ([]domain.Assessment, error) {
	query := `SELECT id, company_id, application_id, assessment_type, title, description, max_score, passing_score, duration_minutes, due_at, status, score, result, result_summary, completed_at, created_by, created_at, updated_at
		 FROM assessments WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if applicationID != "" {
		query += fmt.Sprintf(" AND application_id=$%d", argIdx)
		args = append(args, applicationID)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []domain.Assessment
	for rows.Next() {
		var a domain.Assessment
		rows.Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description,
			&a.MaxScore, &a.PassingScore, &a.DurationMinutes, &a.DueAt, &a.Status, &a.Score,
			&a.Result, &a.ResultSummary, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		assessments = append(assessments, a)
	}
	return assessments, nil
}

func (r *AssessmentRepo) Update(ctx context.Context, companyID, id string, req *domain.Assessment) (*domain.Assessment, error) {
	a := &domain.Assessment{}
	err := r.pool.QueryRow(ctx,
		`UPDATE assessments SET
		 title=COALESCE($3,title), description=COALESCE($4,description),
		 max_score=COALESCE($5,max_score), passing_score=COALESCE($6,passing_score),
		 duration_minutes=COALESCE($7,duration_minutes), due_at=COALESCE($8,due_at),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, application_id, assessment_type, title, description, max_score, passing_score, duration_minutes, due_at, status, score, result, result_summary, completed_at, created_by, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.MaxScore, req.PassingScore,
		req.DurationMinutes, req.DueAt,
	).Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description,
		&a.MaxScore, &a.PassingScore, &a.DurationMinutes, &a.DueAt, &a.Status, &a.Score,
		&a.Result, &a.ResultSummary, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *AssessmentRepo) UpdateStatus(ctx context.Context, companyID, id, status string, score *float64, result *string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`UPDATE assessments SET status=$3, score=COALESCE($4,score), result=COALESCE($5,result),
		 completed_at=CASE WHEN $3 IN ('COMPLETED','EXPIRED') THEN $6 ELSE completed_at END,
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, status, score, result, now)
	return err
}

func (r *AssessmentRepo) AddSection(ctx context.Context, req *domain.AssessmentSection) (*domain.AssessmentSection, error) {
	s := &domain.AssessmentSection{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO assessment_sections (assessment_id, name, description, max_score, weight, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, assessment_id, name, description, max_score, weight, sort_order`,
		req.AssessmentID, req.Name, req.Description, req.MaxScore, req.Weight, req.SortOrder,
	).Scan(&s.ID, &s.AssessmentID, &s.Name, &s.Description, &s.MaxScore, &s.Weight, &s.SortOrder)
	return s, err
}

func (r *AssessmentRepo) ListSections(ctx context.Context, assessmentID string) ([]domain.AssessmentSection, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, assessment_id, name, description, max_score, weight, sort_order
		 FROM assessment_sections WHERE assessment_id=$1 ORDER BY sort_order`, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []domain.AssessmentSection
	for rows.Next() {
		var s domain.AssessmentSection
		rows.Scan(&s.ID, &s.AssessmentID, &s.Name, &s.Description, &s.MaxScore, &s.Weight, &s.SortOrder)
		sections = append(sections, s)
	}
	return sections, nil
}

func (r *AssessmentRepo) AddResult(ctx context.Context, req *domain.AssessmentResult) (*domain.AssessmentResult, error) {
	ar := &domain.AssessmentResult{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO assessment_results (assessment_id, section_id, score, max_score, comment, graded_by, graded_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, assessment_id, section_id, score, max_score, comment, graded_by, graded_at`,
		req.AssessmentID, req.SectionID, req.Score, req.MaxScore, req.Comment, req.GradedBy, time.Now(),
	).Scan(&ar.ID, &ar.AssessmentID, &ar.SectionID, &ar.Score, &ar.MaxScore, &ar.Comment, &ar.GradedBy, &ar.GradedAt)
	return ar, err
}

func (r *AssessmentRepo) ListResults(ctx context.Context, assessmentID string) ([]domain.AssessmentResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, assessment_id, section_id, score, max_score, comment, graded_by, graded_at
		 FROM assessment_results WHERE assessment_id=$1`, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AssessmentResult
	for rows.Next() {
		var ar domain.AssessmentResult
		rows.Scan(&ar.ID, &ar.AssessmentID, &ar.SectionID, &ar.Score, &ar.MaxScore, &ar.Comment, &ar.GradedBy, &ar.GradedAt)
		results = append(results, ar)
	}
	return results, nil
}
