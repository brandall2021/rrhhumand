package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type evaluationRepo struct {
	pool *pgxpool.Pool
}

func (r *evaluationRepo) Create(ctx context.Context, e *domain.PerformanceEvaluation) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_evaluations (company_id, cycle_id, employee_id, evaluator_id, evaluation_type, template_id, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		e.CompanyID, e.CycleID, e.EmployeeID, e.EvaluatorID, e.EvaluationType, e.TemplateID, e.Status, now, now,
	).Scan(&e.ID)
}

func (r *evaluationRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceEvaluation, error) {
	e := &domain.PerformanceEvaluation{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, template_id, status,
		 overall_score, strengths, improvement_areas, summary, submitted_at, reviewed_at, locked_at, created_at, updated_at
		 FROM performance_evaluations WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.EmployeeID, &e.EvaluatorID, &e.EvaluationType, &e.TemplateID,
		&e.Status, &e.OverallScore, &e.Strengths, &e.ImprovementAreas, &e.Summary,
		&e.SubmittedAt, &e.ReviewedAt, &e.LockedAt, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *evaluationRepo) List(ctx context.Context, filter domain.EvaluationFilter) ([]domain.PerformanceEvaluation, error) {
	query := `SELECT id, company_id, cycle_id, employee_id, evaluator_id, evaluation_type, template_id, status,
		 overall_score, strengths, improvement_areas, summary, submitted_at, reviewed_at, locked_at, created_at, updated_at
		 FROM performance_evaluations WHERE company_id=$1`
	args := []interface{}{filter.CompanyID}
	argIdx := 2

	if filter.CycleID != "" {
		query += fmt.Sprintf(" AND cycle_id=$%d", argIdx)
		args = append(args, filter.CycleID)
		argIdx++
	}
	if filter.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filter.EmployeeID)
		argIdx++
	}
	if filter.EvaluatorID != "" {
		query += fmt.Sprintf(" AND evaluator_id=$%d", argIdx)
		args = append(args, filter.EvaluatorID)
		argIdx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.EvaluationType != "" {
		query += fmt.Sprintf(" AND evaluation_type=$%d", argIdx)
		args = append(args, filter.EvaluationType)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []domain.PerformanceEvaluation
	for rows.Next() {
		var e domain.PerformanceEvaluation
		rows.Scan(&e.ID, &e.CompanyID, &e.CycleID, &e.EmployeeID, &e.EvaluatorID, &e.EvaluationType, &e.TemplateID,
			&e.Status, &e.OverallScore, &e.Strengths, &e.ImprovementAreas, &e.Summary,
			&e.SubmittedAt, &e.ReviewedAt, &e.LockedAt, &e.CreatedAt, &e.UpdatedAt)
		evaluations = append(evaluations, e)
	}
	return evaluations, nil
}

func (r *evaluationRepo) Update(ctx context.Context, e *domain.PerformanceEvaluation) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET overall_score=$3, strengths=$4, improvement_areas=$5, summary=$6, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		e.CompanyID, e.ID, e.OverallScore, e.Strengths, e.ImprovementAreas, e.Summary)
	return err
}

func (r *evaluationRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.EvaluationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET status=$3, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`, companyID, id, status)
	return err
}

func (r *evaluationRepo) Submit(ctx context.Context, companyID, id string, score float64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_evaluations SET status='SUBMITTED', overall_score=$3, submitted_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2 AND status='DRAFT'`,
		companyID, id, score)
	return err
}

func (r *evaluationRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_evaluations WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *evaluationRepo) CreateAnswer(ctx context.Context, a *domain.EvaluationAnswer) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO evaluation_answers (evaluation_id, question_id, numeric_value, text_value, selected_value, boolean_value, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		a.EvaluationID, a.QuestionID, a.NumericValue, a.TextValue, a.SelectedValue, a.BooleanValue, now, now,
	).Scan(&a.ID)
}

func (r *evaluationRepo) BulkCreateAnswers(ctx context.Context, answers []domain.EvaluationAnswer) error {
	for i := range answers {
		if err := r.CreateAnswer(ctx, &answers[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *evaluationRepo) ListAnswersByEvaluation(ctx context.Context, evaluationID string) ([]domain.EvaluationAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, evaluation_id, question_id, numeric_value, text_value, selected_value, boolean_value, created_at, updated_at
		 FROM evaluation_answers WHERE evaluation_id=$1 ORDER BY created_at`, evaluationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.EvaluationAnswer
	for rows.Next() {
		var a domain.EvaluationAnswer
		rows.Scan(&a.ID, &a.EvaluationID, &a.QuestionID, &a.NumericValue, &a.TextValue, &a.SelectedValue, &a.BooleanValue, &a.CreatedAt, &a.UpdatedAt)
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *evaluationRepo) UpdateAnswer(ctx context.Context, a *domain.EvaluationAnswer) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE evaluation_answers SET numeric_value=$2, text_value=$3, selected_value=$4, boolean_value=$5, updated_at=NOW()
		 WHERE id=$1`,
		a.ID, a.NumericValue, a.TextValue, a.SelectedValue, a.BooleanValue)
	return err
}

func (r *evaluationRepo) DeleteAnswersByEvaluation(ctx context.Context, evaluationID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM evaluation_answers WHERE evaluation_id=$1`, evaluationID)
	return err
}

func (r *evaluationRepo) CreateObjectiveEvaluation(ctx context.Context, oe *domain.ObjectiveEvaluation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO objective_evaluations (objective_id, evaluation_id, score, comment)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (objective_id, evaluation_id) DO UPDATE SET score=$3, comment=$4`,
		oe.ObjectiveID, oe.EvaluationID, oe.Score, oe.Comment)
	return err
}

func (r *evaluationRepo) ListObjectiveEvaluations(ctx context.Context, evaluationID string) ([]domain.ObjectiveEvaluation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT objective_id, evaluation_id, score, comment
		 FROM objective_evaluations WHERE evaluation_id=$1`, evaluationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oes []domain.ObjectiveEvaluation
	for rows.Next() {
		var oe domain.ObjectiveEvaluation
		rows.Scan(&oe.ObjectiveID, &oe.EvaluationID, &oe.Score, &oe.Comment)
		oes = append(oes, oe)
	}
	return oes, nil
}

func (r *evaluationRepo) CreateCompetencyEvaluation(ctx context.Context, ce *domain.CompetencyEvaluation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO competency_evaluations (competency_id, evaluation_id, score, expected_level, comment)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (competency_id, evaluation_id) DO UPDATE SET score=$3, expected_level=$4, comment=$5`,
		ce.CompetencyID, ce.EvaluationID, ce.Score, ce.ExpectedLevel, ce.Comment)
	return err
}

func (r *evaluationRepo) ListCompetencyEvaluations(ctx context.Context, evaluationID string) ([]domain.CompetencyEvaluation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT competency_id, evaluation_id, score, expected_level, comment
		 FROM competency_evaluations WHERE evaluation_id=$1`, evaluationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ces []domain.CompetencyEvaluation
	for rows.Next() {
		var ce domain.CompetencyEvaluation
		rows.Scan(&ce.CompetencyID, &ce.EvaluationID, &ce.Score, &ce.ExpectedLevel, &ce.Comment)
		ces = append(ces, ce)
	}
	return ces, nil
}

// ReviewRepository

type reviewRepo struct {
	pool *pgxpool.Pool
}

func (r *reviewRepo) Create(ctx context.Context, rev *domain.PerformanceReview) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_reviews (company_id, cycle_id, employee_id, manager_id, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		rev.CompanyID, rev.CycleID, rev.EmployeeID, rev.ManagerID, rev.Status, now, now,
	).Scan(&rev.ID)
}

func (r *reviewRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceReview, error) {
	rev := &domain.PerformanceReview{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, manager_id, summary, strengths, improvement_areas,
		 final_score, final_rating, employee_comments, manager_comments, employee_agreement, disagreement_reason, status, created_at, updated_at
		 FROM performance_reviews WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&rev.ID, &rev.CompanyID, &rev.CycleID, &rev.EmployeeID, &rev.ManagerID, &rev.Summary, &rev.Strengths,
		&rev.ImprovementAreas, &rev.FinalScore, &rev.FinalRating, &rev.EmployeeComments, &rev.ManagerComments,
		&rev.EmployeeAgreement, &rev.DisagreementReason, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt)
	return rev, err
}

func (r *reviewRepo) GetByCycleEmployee(ctx context.Context, companyID, cycleID, employeeID string) (*domain.PerformanceReview, error) {
	rev := &domain.PerformanceReview{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, manager_id, summary, strengths, improvement_areas,
		 final_score, final_rating, employee_comments, manager_comments, employee_agreement, disagreement_reason, status, created_at, updated_at
		 FROM performance_reviews WHERE company_id=$1 AND cycle_id=$2 AND employee_id=$3`,
		companyID, cycleID, employeeID,
	).Scan(&rev.ID, &rev.CompanyID, &rev.CycleID, &rev.EmployeeID, &rev.ManagerID, &rev.Summary, &rev.Strengths,
		&rev.ImprovementAreas, &rev.FinalScore, &rev.FinalRating, &rev.EmployeeComments, &rev.ManagerComments,
		&rev.EmployeeAgreement, &rev.DisagreementReason, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt)
	return rev, err
}

func (r *reviewRepo) ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceReview, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, manager_id, summary, strengths, improvement_areas,
		 final_score, final_rating, employee_comments, manager_comments, employee_agreement, disagreement_reason, status, created_at, updated_at
		 FROM performance_reviews WHERE company_id=$1 AND cycle_id=$2 ORDER BY created_at`, companyID, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.PerformanceReview
	for rows.Next() {
		var rev domain.PerformanceReview
		rows.Scan(&rev.ID, &rev.CompanyID, &rev.CycleID, &rev.EmployeeID, &rev.ManagerID, &rev.Summary, &rev.Strengths,
			&rev.ImprovementAreas, &rev.FinalScore, &rev.FinalRating, &rev.EmployeeComments, &rev.ManagerComments,
			&rev.EmployeeAgreement, &rev.DisagreementReason, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt)
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

func (r *reviewRepo) Update(ctx context.Context, rev *domain.PerformanceReview) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_reviews SET summary=$3, strengths=$4, improvement_areas=$5, final_score=$6, final_rating=$7,
		 employee_comments=$8, manager_comments=$9, employee_agreement=$10, disagreement_reason=$11, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		rev.CompanyID, rev.ID, rev.Summary, rev.Strengths, rev.ImprovementAreas, rev.FinalScore, rev.FinalRating,
		rev.EmployeeComments, rev.ManagerComments, rev.EmployeeAgreement, rev.DisagreementReason)
	return err
}

func (r *reviewRepo) UpdateStatus(ctx context.Context, companyID, id string, status domain.EvaluationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_reviews SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func NewEvaluationRepository(pool *pgxpool.Pool) EvaluationRepository {
	return &evaluationRepo{pool: pool}
}

func NewReviewRepository(pool *pgxpool.Pool) ReviewRepository {
	return &reviewRepo{pool: pool}
}
