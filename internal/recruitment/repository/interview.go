package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type InterviewRepo struct {
	pool *pgxpool.Pool
}

func NewInterviewRepo(pool *pgxpool.Pool) *InterviewRepo {
	return &InterviewRepo{pool: pool}
}

func (r *InterviewRepo) Create(ctx context.Context, companyID string, req *domain.Interview) (*domain.Interview, error) {
	i := &domain.Interview{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interviews (company_id, application_id, interview_type, title, scheduled_at, duration_minutes, meeting_url, meeting_password, location, instructions, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 RETURNING id, company_id, application_id, interview_type, title, scheduled_at, duration_minutes, meeting_url, meeting_password, location, instructions, status, score, notes, cancelled_at, cancel_reason, created_by, created_at, updated_at`,
		companyID, req.ApplicationID, req.InterviewType, req.Title, req.ScheduledAt,
		req.DurationMinutes, req.MeetingURL, req.MeetingPassword, req.Location, req.Instructions, req.CreatedBy,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewType, &i.Title,
		&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.MeetingPassword, &i.Location,
		&i.Instructions, &i.Status, &i.Score, &i.Notes, &i.CancelledAt, &i.CancelReason,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

func (r *InterviewRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Interview, error) {
	i := &domain.Interview{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, application_id, interview_type, title, scheduled_at, duration_minutes, meeting_url, meeting_password, location, instructions, status, score, notes, cancelled_at, cancel_reason, created_by, created_at, updated_at
		 FROM interviews WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewType, &i.Title,
		&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.MeetingPassword, &i.Location,
		&i.Instructions, &i.Status, &i.Score, &i.Notes, &i.CancelledAt, &i.CancelReason,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

func (r *InterviewRepo) List(ctx context.Context, companyID string, applicationID, status string) ([]domain.Interview, error) {
	query := `SELECT id, company_id, application_id, interview_type, title, scheduled_at, duration_minutes, meeting_url, meeting_password, location, instructions, status, score, notes, cancelled_at, cancel_reason, created_by, created_at, updated_at
		 FROM interviews WHERE company_id=$1`
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
	query += " ORDER BY scheduled_at DESC NULLS LAST"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interviews []domain.Interview
	for rows.Next() {
		var i domain.Interview
		rows.Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewType, &i.Title,
			&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.MeetingPassword, &i.Location,
			&i.Instructions, &i.Status, &i.Score, &i.Notes, &i.CancelledAt, &i.CancelReason,
			&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
		interviews = append(interviews, i)
	}
	return interviews, nil
}

func (r *InterviewRepo) Update(ctx context.Context, companyID, id string, req *domain.Interview) (*domain.Interview, error) {
	i := &domain.Interview{}
	err := r.pool.QueryRow(ctx,
		`UPDATE interviews SET
		 interview_type=COALESCE($3,interview_type), title=COALESCE($4,title),
		 scheduled_at=COALESCE($5,scheduled_at), duration_minutes=COALESCE($6,duration_minutes),
		 meeting_url=COALESCE($7,meeting_url), meeting_password=COALESCE($8,meeting_password),
		 location=COALESCE($9,location), instructions=COALESCE($10,instructions),
		 notes=COALESCE($11,notes), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, application_id, interview_type, title, scheduled_at, duration_minutes, meeting_url, meeting_password, location, instructions, status, score, notes, cancelled_at, cancel_reason, created_by, created_at, updated_at`,
		companyID, id, req.InterviewType, req.Title, req.ScheduledAt, req.DurationMinutes,
		req.MeetingURL, req.MeetingPassword, req.Location, req.Instructions, req.Notes,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewType, &i.Title,
		&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.MeetingPassword, &i.Location,
		&i.Instructions, &i.Status, &i.Score, &i.Notes, &i.CancelledAt, &i.CancelReason,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

func (r *InterviewRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`UPDATE interviews SET status=$3, updated_at=NOW(), cancelled_at=CASE WHEN $3='CANCELLED' THEN $4 ELSE cancelled_at END WHERE company_id=$1 AND id=$2`,
		companyID, id, status, now)
	return err
}

func (r *InterviewRepo) AddPanelMember(ctx context.Context, req *domain.InterviewPanelMember) (*domain.InterviewPanelMember, error) {
	pm := &domain.InterviewPanelMember{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interview_panel (interview_id, employee_id, role, status)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, interview_id, employee_id, role, status, response_at, created_at`,
		req.InterviewID, req.EmployeeID, req.Role, req.Status,
	).Scan(&pm.ID, &pm.InterviewID, &pm.EmployeeID, &pm.Role, &pm.Status, &pm.ResponseAt, &pm.CreatedAt)
	return pm, err
}

func (r *InterviewRepo) UpdatePanelMember(ctx context.Context, id string, req *domain.InterviewPanelMember) (*domain.InterviewPanelMember, error) {
	pm := &domain.InterviewPanelMember{}
	err := r.pool.QueryRow(ctx,
		`UPDATE interview_panel SET role=COALESCE($2,role), status=COALESCE($3,status),
		 response_at=COALESCE($4,response_at) WHERE id=$1
		 RETURNING id, interview_id, employee_id, role, status, response_at, created_at`,
		id, req.Role, req.Status, req.ResponseAt,
	).Scan(&pm.ID, &pm.InterviewID, &pm.EmployeeID, &pm.Role, &pm.Status, &pm.ResponseAt, &pm.CreatedAt)
	return pm, err
}

func (r *InterviewRepo) RemovePanelMember(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM interview_panel WHERE id=$1`, id)
	return err
}

func (r *InterviewRepo) ListPanelMembers(ctx context.Context, interviewID string) ([]domain.InterviewPanelMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, interview_id, employee_id, role, status, response_at, created_at
		 FROM interview_panel WHERE interview_id=$1`, interviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.InterviewPanelMember
	for rows.Next() {
		var pm domain.InterviewPanelMember
		rows.Scan(&pm.ID, &pm.InterviewID, &pm.EmployeeID, &pm.Role, &pm.Status, &pm.ResponseAt, &pm.CreatedAt)
		members = append(members, pm)
	}
	return members, nil
}

func (r *InterviewRepo) AddFeedback(ctx context.Context, req *domain.InterviewFeedback) (*domain.InterviewFeedback, error) {
	fb := &domain.InterviewFeedback{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interview_feedback (interview_id, panelist_id, score, comments, strengths, weaknesses, recommendation, submitted_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, interview_id, panelist_id, score, comments, strengths, weaknesses, recommendation, submitted_at, created_at`,
		req.InterviewID, req.PanelistID, req.Score, req.Comments, req.Strengths, req.Weaknesses,
		req.Recommendation, req.SubmittedAt,
	).Scan(&fb.ID, &fb.InterviewID, &fb.PanelistID, &fb.Score, &fb.Comments, &fb.Strengths,
		&fb.Weaknesses, &fb.Recommendation, &fb.SubmittedAt, &fb.CreatedAt)
	return fb, err
}

func (r *InterviewRepo) ListFeedback(ctx context.Context, interviewID string) ([]domain.InterviewFeedback, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, interview_id, panelist_id, score, comments, strengths, weaknesses, recommendation, submitted_at, created_at
		 FROM interview_feedback WHERE interview_id=$1 ORDER BY created_at`, interviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []domain.InterviewFeedback
	for rows.Next() {
		var fb domain.InterviewFeedback
		rows.Scan(&fb.ID, &fb.InterviewID, &fb.PanelistID, &fb.Score, &fb.Comments, &fb.Strengths,
			&fb.Weaknesses, &fb.Recommendation, &fb.SubmittedAt, &fb.CreatedAt)
		feedbacks = append(feedbacks, fb)
	}
	return feedbacks, nil
}

func (r *InterviewRepo) AddFeedbackQuestion(ctx context.Context, req *domain.InterviewFeedbackQuestion) (*domain.InterviewFeedbackQuestion, error) {
	fq := &domain.InterviewFeedbackQuestion{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interview_feedback_questions (interview_feedback_id, question, score, comment)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, interview_feedback_id, question, score, comment`,
		req.InterviewFeedbackID, req.Question, req.Score, req.Comment,
	).Scan(&fq.ID, &fq.InterviewFeedbackID, &fq.Question, &fq.Score, &fq.Comment)
	return fq, err
}

func (r *InterviewRepo) CheckConflicts(ctx context.Context, companyID, employeeID string, scheduledAt time.Time, durationMinutes int) ([]domain.Interview, error) {
	query := `SELECT i.id, i.company_id, i.application_id, i.interview_type, i.title, i.scheduled_at,
		 i.duration_minutes, i.meeting_url, i.meeting_password, i.location, i.instructions,
		 i.status, i.score, i.notes, i.cancelled_at, i.cancel_reason, i.created_by, i.created_at, i.updated_at
		 FROM interviews i
		 INNER JOIN interview_panel ip ON i.id = ip.interview_id
		 WHERE i.company_id=$1 AND ip.employee_id=$2
		 AND i.status IN ('SCHEDULED', 'CONFIRMED')
		 AND i.scheduled_at < $4
		 AND (i.scheduled_at + (COALESCE(i.duration_minutes, 60) * interval '1 minute')) > $3`
	rows, err := r.pool.Query(ctx, query, companyID, employeeID, scheduledAt, scheduledAt.Add(time.Duration(durationMinutes)*time.Minute))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conflicts []domain.Interview
	for rows.Next() {
		var i domain.Interview
		rows.Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewType, &i.Title,
			&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.MeetingPassword, &i.Location,
			&i.Instructions, &i.Status, &i.Score, &i.Notes, &i.CancelledAt, &i.CancelReason,
			&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
		conflicts = append(conflicts, i)
	}
	return conflicts, nil
}
