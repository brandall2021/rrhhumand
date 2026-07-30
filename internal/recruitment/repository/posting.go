package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type PostingRepo struct {
	pool *pgxpool.Pool
}

func NewPostingRepo(pool *pgxpool.Pool) *PostingRepo {
	return &PostingRepo{pool: pool}
}

func (r *PostingRepo) Create(ctx context.Context, companyID string, req *domain.Posting) (*domain.Posting, error) {
	p := &domain.Posting{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_postings (company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, is_public, external_url)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		 RETURNING id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at`,
		companyID, req.PositionID, req.RequisitionID, req.Title, req.Description, req.Requirements,
		req.Responsibilities, req.Benefits, req.EmploymentType, req.WorkMode, req.Location,
		req.SalaryMin, req.SalaryMax, req.Currency, req.IsPublic, req.ExternalURL,
	).Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
		&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
		&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
		&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostingRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Posting, error) {
	p := &domain.Posting{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at
		 FROM job_postings WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
		&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
		&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
		&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostingRepo) List(ctx context.Context, companyID string, status string) ([]domain.Posting, error) {
	query := `SELECT id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at
		 FROM job_postings WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

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

	var postings []domain.Posting
	for rows.Next() {
		var p domain.Posting
		rows.Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
			&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
			&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
			&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		postings = append(postings, p)
	}
	return postings, nil
}

func (r *PostingRepo) Update(ctx context.Context, companyID, id string, req *domain.Posting) (*domain.Posting, error) {
	p := &domain.Posting{}
	err := r.pool.QueryRow(ctx,
		`UPDATE job_postings SET
		 title=COALESCE($3,title), description=COALESCE($4,description), requirements=COALESCE($5,requirements),
		 responsibilities=COALESCE($6,responsibilities), benefits=COALESCE($7,benefits),
		 employment_type=COALESCE($8,employment_type), work_mode=COALESCE($9,work_mode),
		 location=COALESCE($10,location), salary_min=COALESCE($11,salary_min), salary_max=COALESCE($12,salary_max),
		 currency=COALESCE($13,currency), is_public=COALESCE($14,is_public), external_url=COALESCE($15,external_url),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.Requirements, req.Responsibilities,
		req.Benefits, req.EmploymentType, req.WorkMode, req.Location, req.SalaryMin, req.SalaryMax,
		req.Currency, req.IsPublic, req.ExternalURL,
	).Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
		&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
		&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
		&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostingRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	var err error
	now := time.Now()
	if status == "PUBLISHED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3, published_at=$4, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status, now)
	} else if status == "CLOSED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3, closing_at=$4, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status, now)
	} else {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	}
	return err
}

func (r *PostingRepo) AddBoardPost(ctx context.Context, req *domain.PostingBoardPost) (*domain.PostingBoardPost, error) {
	bp := &domain.PostingBoardPost{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO posting_board_posts (posting_id, board_id, external_id, posted_at, status, error_message)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, posting_id, board_id, external_id, posted_at, status, error_message, created_at`,
		req.PostingID, req.BoardID, req.ExternalID, req.PostedAt, req.Status, req.ErrorMsg,
	).Scan(&bp.ID, &bp.PostingID, &bp.BoardID, &bp.ExternalID, &bp.PostedAt, &bp.Status, &bp.ErrorMsg, &bp.CreatedAt)
	return bp, err
}

func (r *PostingRepo) UpdateBoardPost(ctx context.Context, id string, req *domain.PostingBoardPost) (*domain.PostingBoardPost, error) {
	bp := &domain.PostingBoardPost{}
	err := r.pool.QueryRow(ctx,
		`UPDATE posting_board_posts SET external_id=COALESCE($2,external_id), posted_at=COALESCE($3,posted_at),
		 status=COALESCE($4,status), error_message=COALESCE($5,error_message) WHERE id=$1
		 RETURNING id, posting_id, board_id, external_id, posted_at, status, error_message, created_at`,
		id, req.ExternalID, req.PostedAt, req.Status, req.ErrorMsg,
	).Scan(&bp.ID, &bp.PostingID, &bp.BoardID, &bp.ExternalID, &bp.PostedAt, &bp.Status, &bp.ErrorMsg, &bp.CreatedAt)
	return bp, err
}

func (r *PostingRepo) ListBoardPosts(ctx context.Context, postingID string) ([]domain.PostingBoardPost, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, posting_id, board_id, external_id, posted_at, status, error_message, created_at
		 FROM posting_board_posts WHERE posting_id=$1 ORDER BY created_at`, postingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.PostingBoardPost
	for rows.Next() {
		var bp domain.PostingBoardPost
		rows.Scan(&bp.ID, &bp.PostingID, &bp.BoardID, &bp.ExternalID, &bp.PostedAt, &bp.Status, &bp.ErrorMsg, &bp.CreatedAt)
		posts = append(posts, bp)
	}
	return posts, nil
}

func (r *PostingRepo) AddScreeningQuestion(ctx context.Context, req *domain.PostingScreeningQuestion) (*domain.PostingScreeningQuestion, error) {
	q := &domain.PostingScreeningQuestion{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO posting_screening_questions (posting_id, question, question_type, options, required, sort_order, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, posting_id, question, question_type, options, required, sort_order, active, created_at`,
		req.PostingID, req.Question, req.QuestionType, req.Options, req.Required, req.SortOrder, req.Active,
	).Scan(&q.ID, &q.PostingID, &q.Question, &q.QuestionType, &q.Options, &q.Required, &q.SortOrder, &q.Active, &q.CreatedAt)
	return q, err
}

func (r *PostingRepo) UpdateScreeningQuestion(ctx context.Context, id string, req *domain.PostingScreeningQuestion) (*domain.PostingScreeningQuestion, error) {
	q := &domain.PostingScreeningQuestion{}
	err := r.pool.QueryRow(ctx,
		`UPDATE posting_screening_questions SET question=COALESCE($2,question), question_type=COALESCE($3,question_type),
		 options=COALESCE($4,options), required=COALESCE($5,required), sort_order=COALESCE($6,sort_order),
		 active=COALESCE($7,active) WHERE id=$1
		 RETURNING id, posting_id, question, question_type, options, required, sort_order, active, created_at`,
		id, req.Question, req.QuestionType, req.Options, req.Required, req.SortOrder, req.Active,
	).Scan(&q.ID, &q.PostingID, &q.Question, &q.QuestionType, &q.Options, &q.Required, &q.SortOrder, &q.Active, &q.CreatedAt)
	return q, err
}

func (r *PostingRepo) DeleteScreeningQuestion(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM posting_screening_questions WHERE id=$1`, id)
	return err
}

func (r *PostingRepo) ListPublic(ctx context.Context) ([]domain.Posting, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at
		 FROM job_postings WHERE status='PUBLISHED' AND is_public=TRUE ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postings []domain.Posting
	for rows.Next() {
		var p domain.Posting
		rows.Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
			&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
			&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
			&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		postings = append(postings, p)
	}
	return postings, nil
}

func (r *PostingRepo) GetPublicByID(ctx context.Context, id string) (*domain.Posting, error) {
	p := &domain.Posting{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, position_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, salary_min, salary_max, currency, published_at, closing_at, is_public, external_url, status, created_at, updated_at
		 FROM job_postings WHERE id=$1 AND status='PUBLISHED' AND is_public=TRUE`, id,
	).Scan(&p.ID, &p.CompanyID, &p.PositionID, &p.RequisitionID, &p.Title, &p.Description,
		&p.Requirements, &p.Responsibilities, &p.Benefits, &p.EmploymentType, &p.WorkMode,
		&p.Location, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.PublishedAt, &p.ClosingAt,
		&p.IsPublic, &p.ExternalURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostingRepo) ListScreeningQuestions(ctx context.Context, postingID string) ([]domain.PostingScreeningQuestion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, posting_id, question, question_type, options, required, sort_order, active, created_at
		 FROM posting_screening_questions WHERE posting_id=$1 AND active=TRUE ORDER BY sort_order`, postingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.PostingScreeningQuestion
	for rows.Next() {
		var q domain.PostingScreeningQuestion
		rows.Scan(&q.ID, &q.PostingID, &q.Question, &q.QuestionType, &q.Options, &q.Required, &q.SortOrder, &q.Active, &q.CreatedAt)
		questions = append(questions, q)
	}
	return questions, nil
}
