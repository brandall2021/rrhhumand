package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type ApplicationRepo struct {
	pool *pgxpool.Pool
}

func NewApplicationRepo(pool *pgxpool.Pool) *ApplicationRepo {
	return &ApplicationRepo{pool: pool}
}

func (r *ApplicationRepo) Create(ctx context.Context, companyID string, req *domain.Application) (*domain.Application, error) {
	a := &domain.Application{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO applications (company_id, candidate_id, posting_id, current_stage_id, status, source, source_detail, is_internal_mobility, consent_given, consent_at, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 RETURNING id, company_id, candidate_id, posting_id, current_stage_id, status, score, applied_at, reviewed_at, rejected_at, rejection_reason_id, rejection_reason_text, hired_at, withdrawn_at, withdraw_reason, source, source_detail, is_internal_mobility, consent_given, consent_at, notes, created_at, updated_at`,
		companyID, req.CandidateID, req.PostingID, req.CurrentStageID, req.Status, req.Source,
		req.SourceDetail, req.IsInternalMobility, req.ConsentGiven, req.ConsentAt, req.Notes,
	).Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.PostingID, &a.CurrentStageID, &a.Status, &a.Score,
		&a.AppliedAt, &a.ReviewedAt, &a.RejectedAt, &a.RejectionReasonID, &a.RejectionText,
		&a.HiredAt, &a.WithdrawnAt, &a.WithdrawReason, &a.Source, &a.SourceDetail,
		&a.IsInternalMobility, &a.ConsentGiven, &a.ConsentAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *ApplicationRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Application, error) {
	a := &domain.Application{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, candidate_id, posting_id, current_stage_id, status, score, applied_at, reviewed_at, rejected_at, rejection_reason_id, rejection_reason_text, hired_at, withdrawn_at, withdraw_reason, source, source_detail, is_internal_mobility, consent_given, consent_at, notes, created_at, updated_at
		 FROM applications WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.PostingID, &a.CurrentStageID, &a.Status, &a.Score,
		&a.AppliedAt, &a.ReviewedAt, &a.RejectedAt, &a.RejectionReasonID, &a.RejectionText,
		&a.HiredAt, &a.WithdrawnAt, &a.WithdrawReason, &a.Source, &a.SourceDetail,
		&a.IsInternalMobility, &a.ConsentGiven, &a.ConsentAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *ApplicationRepo) List(ctx context.Context, companyID string, candidateID, postingID, status string) ([]domain.Application, error) {
	query := `SELECT id, company_id, candidate_id, posting_id, current_stage_id, status, score, applied_at, reviewed_at, rejected_at, rejection_reason_id, rejection_reason_text, hired_at, withdrawn_at, withdraw_reason, source, source_detail, is_internal_mobility, consent_given, consent_at, notes, created_at, updated_at
		 FROM applications WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if candidateID != "" {
		query += fmt.Sprintf(" AND candidate_id=$%d", argIdx)
		args = append(args, candidateID)
		argIdx++
	}
	if postingID != "" {
		query += fmt.Sprintf(" AND posting_id=$%d", argIdx)
		args = append(args, postingID)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	query += " ORDER BY applied_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []domain.Application
	for rows.Next() {
		var a domain.Application
		rows.Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.PostingID, &a.CurrentStageID, &a.Status, &a.Score,
			&a.AppliedAt, &a.ReviewedAt, &a.RejectedAt, &a.RejectionReasonID, &a.RejectionText,
			&a.HiredAt, &a.WithdrawnAt, &a.WithdrawReason, &a.Source, &a.SourceDetail,
			&a.IsInternalMobility, &a.ConsentGiven, &a.ConsentAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *ApplicationRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	var err error
	switch status {
	case "REJECTED":
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, rejected_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	case "HIRED":
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, hired_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	case "WITHDRAWN":
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, withdrawn_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	default:
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	}
	return err
}

func (r *ApplicationRepo) UpdateStage(ctx context.Context, companyID, id, stageID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE applications SET current_stage_id=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, stageID)
	return err
}

func (r *ApplicationRepo) GetByCandidateAndPosting(ctx context.Context, companyID, candidateID, postingID string) (*domain.Application, error) {
	a := &domain.Application{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, candidate_id, posting_id, current_stage_id, status, score, applied_at, reviewed_at, rejected_at, rejection_reason_id, rejection_reason_text, hired_at, withdrawn_at, withdraw_reason, source, source_detail, is_internal_mobility, consent_given, consent_at, notes, created_at, updated_at
		 FROM applications WHERE company_id=$1 AND candidate_id=$2 AND posting_id=$3`,
		companyID, candidateID, postingID,
	).Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.PostingID, &a.CurrentStageID, &a.Status, &a.Score,
		&a.AppliedAt, &a.ReviewedAt, &a.RejectedAt, &a.RejectionReasonID, &a.RejectionText,
		&a.HiredAt, &a.WithdrawnAt, &a.WithdrawReason, &a.Source, &a.SourceDetail,
		&a.IsInternalMobility, &a.ConsentGiven, &a.ConsentAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *ApplicationRepo) AddStageHistory(ctx context.Context, req *domain.ApplicationStageHistory) (*domain.ApplicationStageHistory, error) {
	h := &domain.ApplicationStageHistory{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO application_stage_history (application_id, from_stage_id, to_stage_id, changed_by, reason, auto_transition)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, application_id, from_stage_id, to_stage_id, changed_by, reason, auto_transition, created_at`,
		req.ApplicationID, req.FromStageID, req.ToStageID, req.ChangedBy, req.Reason, req.AutoTransition,
	).Scan(&h.ID, &h.ApplicationID, &h.FromStageID, &h.ToStageID, &h.ChangedBy, &h.Reason, &h.AutoTransition, &h.CreatedAt)
	return h, err
}

func (r *ApplicationRepo) ListStageHistory(ctx context.Context, applicationID string) ([]domain.ApplicationStageHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, application_id, from_stage_id, to_stage_id, changed_by, reason, auto_transition, created_at
		 FROM application_stage_history WHERE application_id=$1 ORDER BY created_at ASC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.ApplicationStageHistory
	for rows.Next() {
		var h domain.ApplicationStageHistory
		rows.Scan(&h.ID, &h.ApplicationID, &h.FromStageID, &h.ToStageID, &h.ChangedBy, &h.Reason, &h.AutoTransition, &h.CreatedAt)
		history = append(history, h)
	}
	return history, nil
}

func (r *ApplicationRepo) AddRating(ctx context.Context, req *domain.ApplicationRating) (*domain.ApplicationRating, error) {
	rt := &domain.ApplicationRating{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO application_ratings (application_id, rated_by, rating, comment)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, application_id, rated_by, rating, comment, created_at`,
		req.ApplicationID, req.RatedBy, req.Rating, req.Comment,
	).Scan(&rt.ID, &rt.ApplicationID, &rt.RatedBy, &rt.Rating, &rt.Comment, &rt.CreatedAt)
	return rt, err
}

func (r *ApplicationRepo) ListRatings(ctx context.Context, applicationID string) ([]domain.ApplicationRating, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, application_id, rated_by, rating, comment, created_at
		 FROM application_ratings WHERE application_id=$1 ORDER BY created_at DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []domain.ApplicationRating
	for rows.Next() {
		var rt domain.ApplicationRating
		rows.Scan(&rt.ID, &rt.ApplicationID, &rt.RatedBy, &rt.Rating, &rt.Comment, &rt.CreatedAt)
		ratings = append(ratings, rt)
	}
	return ratings, nil
}

func (r *ApplicationRepo) AddNote(ctx context.Context, req *domain.ApplicationNote) (*domain.ApplicationNote, error) {
	n := &domain.ApplicationNote{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO application_notes (application_id, author_id, content, is_private)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, application_id, author_id, content, is_private, created_at, updated_at`,
		req.ApplicationID, req.AuthorID, req.Content, req.IsPrivate,
	).Scan(&n.ID, &n.ApplicationID, &n.AuthorID, &n.Content, &n.IsPrivate, &n.CreatedAt, &n.UpdatedAt)
	return n, err
}

func (r *ApplicationRepo) UpdateNote(ctx context.Context, id string, req *domain.ApplicationNote) (*domain.ApplicationNote, error) {
	n := &domain.ApplicationNote{}
	err := r.pool.QueryRow(ctx,
		`UPDATE application_notes SET content=COALESCE($2,content), is_private=COALESCE($3,is_private),
		 updated_at=NOW() WHERE id=$1
		 RETURNING id, application_id, author_id, content, is_private, created_at, updated_at`,
		id, req.Content, req.IsPrivate,
	).Scan(&n.ID, &n.ApplicationID, &n.AuthorID, &n.Content, &n.IsPrivate, &n.CreatedAt, &n.UpdatedAt)
	return n, err
}

func (r *ApplicationRepo) ListNotes(ctx context.Context, applicationID string) ([]domain.ApplicationNote, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, application_id, author_id, content, is_private, created_at, updated_at
		 FROM application_notes WHERE application_id=$1 ORDER BY created_at DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.ApplicationNote
	for rows.Next() {
		var n domain.ApplicationNote
		rows.Scan(&n.ID, &n.ApplicationID, &n.AuthorID, &n.Content, &n.IsPrivate, &n.CreatedAt, &n.UpdatedAt)
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *ApplicationRepo) AddScreeningAnswer(ctx context.Context, req *domain.ApplicationScreeningAnswer) (*domain.ApplicationScreeningAnswer, error) {
	sa := &domain.ApplicationScreeningAnswer{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO application_screening_answers (application_id, question_id, answer)
		 VALUES ($1,$2,$3)
		 RETURNING id, application_id, question_id, answer, created_at`,
		req.ApplicationID, req.QuestionID, req.Answer,
	).Scan(&sa.ID, &sa.ApplicationID, &sa.QuestionID, &sa.Answer, &sa.CreatedAt)
	return sa, err
}

func (r *ApplicationRepo) ListScreeningAnswers(ctx context.Context, applicationID string) ([]domain.ApplicationScreeningAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, application_id, question_id, answer, created_at
		 FROM application_screening_answers WHERE application_id=$1`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.ApplicationScreeningAnswer
	for rows.Next() {
		var sa domain.ApplicationScreeningAnswer
		rows.Scan(&sa.ID, &sa.ApplicationID, &sa.QuestionID, &sa.Answer, &sa.CreatedAt)
		answers = append(answers, sa)
	}
	return answers, nil
}
