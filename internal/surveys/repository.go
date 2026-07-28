package surveys

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type SurveyRepository struct {
	pool *pgxpool.Pool
}

func NewSurveyRepository(pool *pgxpool.Pool) *SurveyRepository {
	return &SurveyRepository{pool: pool}
}

func (r *SurveyRepository) CreateSurvey(ctx context.Context, s *models.Survey) error {
	query := `
		INSERT INTO surveys (id, company_id, title, description, type, status, anonymous, multiple_responses, starts_at, ends_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		s.ID, s.CompanyID, s.Title, s.Description, s.Type, s.Status,
		s.Anonymous, s.MultipleResponses, s.StartsAt, s.EndsAt, s.CreatedBy,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *SurveyRepository) GetSurveyByID(ctx context.Context, id, companyID string) (*models.Survey, error) {
	query := `
		SELECT
			s.id, s.company_id, s.title, s.description, s.type, s.status,
			s.anonymous, s.multiple_responses, s.starts_at, s.ends_at,
			s.created_by, u.first_name || ' ' || u.last_name,
			s.created_at, s.updated_at
		FROM surveys s
		JOIN users u ON u.id = s.created_by
		WHERE s.id = $1 AND s.company_id = $2`

	survey := &models.Survey{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&survey.ID, &survey.CompanyID, &survey.Title, &survey.Description,
		&survey.Type, &survey.Status, &survey.Anonymous, &survey.MultipleResponses,
		&survey.StartsAt, &survey.EndsAt,
		&survey.CreatedBy, &survey.CreatedByName,
		&survey.CreatedAt, &survey.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return survey, nil
}

func (r *SurveyRepository) ListSurveys(ctx context.Context, companyID string, filters SurveyFilters, offset, limit int) ([]models.Survey, int64, error) {
	where := []string{`s.company_id = $1`}
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		where = append(where, fmt.Sprintf(`s.status = $%d`, argIdx))
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.Type != "" {
		where = append(where, fmt.Sprintf(`s.type = $%d`, argIdx))
		args = append(args, filters.Type)
		argIdx++
	}
	if filters.Search != "" {
		where = append(where, fmt.Sprintf(`s.title ILIKE $%d`, argIdx))
		args = append(args, "%"+filters.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM surveys s WHERE %s`, whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT
			s.id, s.company_id, s.title, s.description, s.type, s.status,
			s.anonymous, s.multiple_responses, s.starts_at, s.ends_at,
			s.created_by, u.first_name || ' ' || u.last_name,
			s.created_at, s.updated_at
		FROM surveys s
		JOIN users u ON u.id = s.created_by
		WHERE %s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var surveys []models.Survey
	for rows.Next() {
		var s models.Survey
		if err := rows.Scan(
			&s.ID, &s.CompanyID, &s.Title, &s.Description,
			&s.Type, &s.Status, &s.Anonymous, &s.MultipleResponses,
			&s.StartsAt, &s.EndsAt,
			&s.CreatedBy, &s.CreatedByName,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		surveys = append(surveys, s)
	}
	return surveys, total, nil
}

func (r *SurveyRepository) UpdateSurvey(ctx context.Context, s *models.Survey) error {
	query := `
		UPDATE surveys SET title=$1, description=$2, type=$3, anonymous=$4,
			multiple_responses=$5, starts_at=$6, ends_at=$7, updated_at=NOW()
		WHERE id=$8 AND company_id=$9`
	_, err := r.pool.Exec(ctx, query,
		s.Title, s.Description, s.Type, s.Anonymous,
		s.MultipleResponses, s.StartsAt, s.EndsAt,
		s.ID, s.CompanyID,
	)
	return err
}

func (r *SurveyRepository) UpdateSurveyStatus(ctx context.Context, id, companyID, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE surveys SET status=$1, updated_at=NOW() WHERE id=$2 AND company_id=$3`,
		status, id, companyID,
	)
	return err
}

func (r *SurveyRepository) DeleteSurvey(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM surveys WHERE id=$1 AND company_id=$2`, id, companyID,
	)
	return err
}

func (r *SurveyRepository) CreateQuestion(ctx context.Context, q *models.SurveyQuestion) error {
	query := `
		INSERT INTO survey_questions (id, survey_id, question, type, position, required)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query,
		q.ID, q.SurveyID, q.Question, q.Type, q.Position, q.Required,
	).Scan(&q.CreatedAt)
}

func (r *SurveyRepository) GetQuestionByID(ctx context.Context, id string) (*models.SurveyQuestion, error) {
	query := `
		SELECT id, survey_id, question, type, position, required, created_at
		FROM survey_questions WHERE id=$1`
	q := &models.SurveyQuestion{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.SurveyID, &q.Question, &q.Type, &q.Position, &q.Required, &q.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (r *SurveyRepository) ListQuestionsBySurveyID(ctx context.Context, surveyID string) ([]models.SurveyQuestion, error) {
	query := `
		SELECT id, survey_id, question, type, position, required, created_at
		FROM survey_questions WHERE survey_id=$1 ORDER BY position ASC`
	rows, err := r.pool.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.SurveyQuestion
	for rows.Next() {
		var q models.SurveyQuestion
		if err := rows.Scan(&q.ID, &q.SurveyID, &q.Question, &q.Type, &q.Position, &q.Required, &q.CreatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *SurveyRepository) UpdateQuestion(ctx context.Context, q *models.SurveyQuestion) error {
	query := `
		UPDATE survey_questions SET question=$1, type=$2, position=$3, required=$4
		WHERE id=$5`
	_, err := r.pool.Exec(ctx, query, q.Question, q.Type, q.Position, q.Required, q.ID)
	return err
}

func (r *SurveyRepository) DeleteQuestion(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM survey_questions WHERE id=$1`, id)
	return err
}

func (r *SurveyRepository) GetQuestionResponseCount(ctx context.Context, questionID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_answers WHERE question_id=$1`, questionID,
	).Scan(&count)
	return count, err
}

func (r *SurveyRepository) CreateOption(ctx context.Context, o *models.SurveyOption) error {
	query := `
		INSERT INTO survey_options (id, question_id, option_text, position)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query, o.ID, o.QuestionID, o.OptionText, o.Position).Scan(&o.CreatedAt)
}

func (r *SurveyRepository) ListOptionsByQuestionID(ctx context.Context, questionID string) ([]models.SurveyOption, error) {
	query := `
		SELECT id, question_id, option_text, position, created_at
		FROM survey_options WHERE question_id=$1 ORDER BY position ASC`
	rows, err := r.pool.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.SurveyOption
	for rows.Next() {
		var o models.SurveyOption
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.OptionText, &o.Position, &o.CreatedAt); err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, nil
}

func (r *SurveyRepository) DeleteOption(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM survey_options WHERE id=$1`, id)
	return err
}

func (r *SurveyRepository) SetTargets(ctx context.Context, surveyID string, targets []models.SurveyTarget) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM survey_targets WHERE survey_id=$1`, surveyID)
	if err != nil {
		return err
	}

	for _, t := range targets {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO survey_targets (id, survey_id, target_type, target_id) VALUES (gen_random_uuid(), $1, $2, $3)`,
			surveyID, t.TargetType, t.TargetID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SurveyRepository) ListTargetsBySurveyID(ctx context.Context, surveyID string) ([]models.SurveyTarget, error) {
	query := `
		SELECT id, survey_id, target_type, target_id, created_at
		FROM survey_targets WHERE survey_id=$1`
	rows, err := r.pool.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []models.SurveyTarget
	for rows.Next() {
		var t models.SurveyTarget
		if err := rows.Scan(&t.ID, &t.SurveyID, &t.TargetType, &t.TargetID, &t.CreatedAt); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (r *SurveyRepository) IsEmployeeTargeted(ctx context.Context, surveyID, employeeID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM survey_targets st
			WHERE st.survey_id = $1
			AND (
				st.target_type = 'ALL'
				OR (st.target_type = 'EMPLOYEE' AND st.target_id = $2)
				OR (st.target_type = 'DEPARTMENT' AND st.target_id = (
					SELECT department_id FROM employees WHERE id = $2
				))
				OR (st.target_type = 'BRANCH' AND st.target_id = (
					SELECT branch_id FROM employees WHERE id = $2
				))
				OR (st.target_type = 'POSITION' AND st.target_id = (
					SELECT position_id FROM employees WHERE id = $2
				))
			)
		)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, surveyID, employeeID).Scan(&exists)
	return exists, err
}

func (r *SurveyRepository) HasEmployeeResponded(ctx context.Context, surveyID, employeeID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM survey_responses WHERE survey_id=$1 AND employee_id=$2)`,
		surveyID, employeeID,
	).Scan(&exists)
	return exists, err
}

func (r *SurveyRepository) CreateResponse(ctx context.Context, resp *models.SurveyResponse) error {
	query := `
		INSERT INTO survey_responses (id, survey_id, employee_id, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING submitted_at`
	return r.pool.QueryRow(ctx, query,
		resp.ID, resp.SurveyID, resp.EmployeeID, nil, nil,
	).Scan(&resp.SubmittedAt)
}

func (r *SurveyRepository) CreateAnswer(ctx context.Context, a *models.SurveyAnswer) error {
	query := `
		INSERT INTO survey_answers (id, response_id, question_id, text_value, number_value, option_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query,
		a.ID, a.ResponseID, a.QuestionID, a.TextValue, a.NumberValue, a.OptionID,
	).Scan(&a.CreatedAt)
}

func (r *SurveyRepository) CreateAnswerOption(ctx context.Context, ao *models.SurveyAnswerOption) error {
	query := `
		INSERT INTO survey_answer_options (id, answer_id, option_id)
		VALUES (gen_random_uuid(), $1, $2)
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, ao.AnswerID, ao.OptionID).Scan(&ao.ID, &ao.CreatedAt)
}

func (r *SurveyRepository) GetResponseCount(ctx context.Context, surveyID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_responses WHERE survey_id=$1`, surveyID,
	).Scan(&count)
	return count, err
}

func (r *SurveyRepository) GetTargetedEmployeeCount(ctx context.Context, surveyID string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT e.id)
		FROM employees e
		JOIN survey_targets st ON st.survey_id = $1
		WHERE e.company_id = (SELECT company_id FROM surveys WHERE id = $1)
		AND e.status = 'active'
		AND (
			st.target_type = 'ALL'
			OR (st.target_type = 'EMPLOYEE' AND st.target_id = e.id)
			OR (st.target_type = 'DEPARTMENT' AND st.target_id = e.department_id)
			OR (st.target_type = 'BRANCH' AND st.target_id = e.branch_id)
			OR (st.target_type = 'POSITION' AND st.target_id = e.position_id)
		)`
	var count int
	err := r.pool.QueryRow(ctx, query, surveyID).Scan(&count)
	return count, err
}

func (r *SurveyRepository) ListAvailableSurveys(ctx context.Context, companyID, employeeID string) ([]models.Survey, error) {
	query := `
		SELECT DISTINCT
			s.id, s.company_id, s.title, s.description, s.type, s.status,
			s.anonymous, s.multiple_responses, s.starts_at, s.ends_at,
			s.created_by, u.first_name || ' ' || u.last_name,
			s.created_at, s.updated_at
		FROM surveys s
		JOIN users u ON u.id = s.created_by
		JOIN survey_targets st ON st.survey_id = s.id
		LEFT JOIN survey_responses sr ON sr.survey_id = s.id AND sr.employee_id = $2
		WHERE s.company_id = $1
		AND s.status = 'PUBLISHED'
		AND sr.id IS NULL
		AND (
			s.starts_at IS NULL OR s.starts_at <= NOW()
		)
		AND (
			s.ends_at IS NULL OR s.ends_at >= NOW()
		)
		AND (
			st.target_type = 'ALL'
			OR (st.target_type = 'EMPLOYEE' AND st.target_id = $2)
			OR (st.target_type = 'DEPARTMENT' AND st.target_id = (
				SELECT department_id FROM employees WHERE id = $2
			))
			OR (st.target_type = 'BRANCH' AND st.target_id = (
				SELECT branch_id FROM employees WHERE id = $2
			))
			OR (st.target_type = 'POSITION' AND st.target_id = (
				SELECT position_id FROM employees WHERE id = $2
			))
		)
		ORDER BY s.created_at DESC`

	rows, err := r.pool.Query(ctx, query, companyID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []models.Survey
	for rows.Next() {
		var s models.Survey
		if err := rows.Scan(
			&s.ID, &s.CompanyID, &s.Title, &s.Description,
			&s.Type, &s.Status, &s.Anonymous, &s.MultipleResponses,
			&s.StartsAt, &s.EndsAt,
			&s.CreatedBy, &s.CreatedByName,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		surveys = append(surveys, s)
	}
	return surveys, nil
}

func (r *SurveyRepository) GetRatingStats(ctx context.Context, questionID string) (avg, min, max float64, count int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(number_value),0), COALESCE(MIN(number_value),0), COALESCE(MAX(number_value),0), COUNT(*)
		FROM survey_answers WHERE question_id=$1 AND number_value IS NOT NULL`, questionID,
	).Scan(&avg, &min, &max, &count)
	return
}

func (r *SurveyRepository) GetOptionDistribution(ctx context.Context, questionID string) ([]OptionDistribution, error) {
	query := `
		SELECT so.id, so.option_text, COUNT(sa.id) as cnt
		FROM survey_options so
		LEFT JOIN survey_answers sa ON sa.option_id = so.id
		WHERE so.question_id = $1
		GROUP BY so.id, so.option_text, so.position
		ORDER BY so.position`
	rows, err := r.pool.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dist []OptionDistribution
	for rows.Next() {
		var d OptionDistribution
		if err := rows.Scan(&d.OptionID, &d.OptionText, &d.Count); err != nil {
			return nil, err
		}
		dist = append(dist, d)
	}
	return dist, nil
}

func (r *SurveyRepository) GetMultipleChoiceDistribution(ctx context.Context, questionID string) ([]OptionDistribution, error) {
	query := `
		SELECT so.id, so.option_text, COUNT(sao.id) as cnt
		FROM survey_options so
		LEFT JOIN survey_answer_options sao ON sao.option_id = so.id
		LEFT JOIN survey_answers sa ON sa.id = sao.answer_id AND sa.question_id = $1
		WHERE so.question_id = $1
		GROUP BY so.id, so.option_text, so.position
		ORDER BY so.position`
	rows, err := r.pool.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dist []OptionDistribution
	for rows.Next() {
		var d OptionDistribution
		if err := rows.Scan(&d.OptionID, &d.OptionText, &d.Count); err != nil {
			return nil, err
		}
		dist = append(dist, d)
	}
	return dist, nil
}

func (r *SurveyRepository) GetYesNoStats(ctx context.Context, questionID string) (yesCount, noCount int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN option_id IN (
				SELECT id FROM survey_options WHERE question_id=$1 AND option_text ILIKE 'si' OR option_text ILIKE 'sí' OR option_text ILIKE 'yes'
			) THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN option_id IN (
				SELECT id FROM survey_options WHERE question_id=$1 AND option_text ILIKE 'no'
			) THEN 1 ELSE 0 END), 0)
		FROM survey_answers WHERE question_id=$1`, questionID,
	).Scan(&yesCount, &noCount)
	return
}

func (r *SurveyRepository) GetNumberStats(ctx context.Context, questionID string) (avg, min, max float64, count int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(number_value),0), COALESCE(MIN(number_value),0), COALESCE(MAX(number_value),0), COUNT(*)
		FROM survey_answers WHERE question_id=$1 AND number_value IS NOT NULL`, questionID,
	).Scan(&avg, &min, &max, &count)
	return
}

func (r *SurveyRepository) GetTextSamples(ctx context.Context, questionID string, limit int) ([]string, error) {
	query := `
		SELECT text_value FROM survey_answers
		WHERE question_id=$1 AND text_value IS NOT NULL
		ORDER BY created_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, query, questionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		texts = append(texts, t)
	}
	return texts, nil
}

func (r *SurveyRepository) GetSurveyResults(ctx context.Context, surveyID string) ([]models.SurveyResponse, error) {
	query := `
		SELECT sr.id, sr.survey_id, sr.employee_id, sr.submitted_at
		FROM survey_responses sr
		WHERE sr.survey_id = $1
		ORDER BY sr.submitted_at`
	rows, err := r.pool.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.SurveyResponse
	for rows.Next() {
		var resp models.SurveyResponse
		if err := rows.Scan(&resp.ID, &resp.SurveyID, &resp.EmployeeID, &resp.SubmittedAt); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (r *SurveyRepository) GetAnswersByResponseID(ctx context.Context, responseID string) ([]models.SurveyAnswer, error) {
	query := `
		SELECT id, response_id, question_id, text_value, number_value, option_id, created_at
		FROM survey_answers WHERE response_id=$1`
	rows, err := r.pool.Query(ctx, query, responseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.SurveyAnswer
	for rows.Next() {
		var a models.SurveyAnswer
		if err := rows.Scan(&a.ID, &a.ResponseID, &a.QuestionID, &a.TextValue, &a.NumberValue, &a.OptionID, &a.CreatedAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *SurveyRepository) GetAnswerOptions(ctx context.Context, answerID string) ([]models.SurveyAnswerOption, error) {
	query := `
		SELECT id, answer_id, option_id, created_at
		FROM survey_answer_options WHERE answer_id=$1`
	rows, err := r.pool.Query(ctx, query, answerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opts []models.SurveyAnswerOption
	for rows.Next() {
		var o models.SurveyAnswerOption
		if err := rows.Scan(&o.ID, &o.AnswerID, &o.OptionID, &o.CreatedAt); err != nil {
			return nil, err
		}
		opts = append(opts, o)
	}
	return opts, nil
}

func (r *SurveyRepository) GetTotalAnswerCount(ctx context.Context, questionID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_answers WHERE question_id=$1`, questionID,
	).Scan(&count)
	return count, err
}

func (r *SurveyRepository) GetResponseCountForRating(ctx context.Context, questionID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_answers WHERE question_id=$1 AND number_value IS NOT NULL`, questionID,
	).Scan(&count)
	return count, err
}

func (r *SurveyRepository) GetResponseCountForOption(ctx context.Context, questionID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_answers WHERE question_id=$1 AND option_id IS NOT NULL`, questionID,
	).Scan(&count)
	return count, err
}

func (r *SurveyRepository) GetResponsesByEmployeeForExport(ctx context.Context, surveyID string) ([]models.SurveyResponse, error) {
	query := `
		SELECT sr.id, sr.survey_id, sr.employee_id, sr.submitted_at
		FROM survey_responses sr
		WHERE sr.survey_id = $1
		ORDER BY sr.submitted_at`
	rows, err := r.pool.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.SurveyResponse
	for rows.Next() {
		var resp models.SurveyResponse
		if err := rows.Scan(&resp.ID, &resp.SurveyID, &resp.EmployeeID, &resp.SubmittedAt); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (r *SurveyRepository) GetAnswersForExport(ctx context.Context, responseIDs []string) (map[string][]models.SurveyAnswer, error) {
	if len(responseIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT id, response_id, question_id, text_value, number_value, option_id, created_at
		FROM survey_answers WHERE response_id = ANY($1)`
	rows, err := r.pool.Query(ctx, query, responseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]models.SurveyAnswer)
	for rows.Next() {
		var a models.SurveyAnswer
		if err := rows.Scan(&a.ID, &a.ResponseID, &a.QuestionID, &a.TextValue, &a.NumberValue, &a.OptionID, &a.CreatedAt); err != nil {
			return nil, err
		}
		result[a.ResponseID] = append(result[a.ResponseID], a)
	}
	return result, nil
}

func (r *SurveyRepository) GetEmployeeName(ctx context.Context, employeeID string) (string, error) {
	var name string
	err := r.pool.QueryRow(ctx,
		`SELECT first_name || ' ' || last_name FROM employees WHERE id=$1`, employeeID,
	).Scan(&name)
	return name, err
}

var _ = time.Now
var _ = pgx.ErrNoRows
