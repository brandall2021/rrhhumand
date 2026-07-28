package recruitment

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Requisitions
func (r *Repository) CreateRequisition(ctx context.Context, companyID, requestedBy string, req *CreateRequisitionRequest) (*JobRequisition, error) {
	vacancies := 1
	if req.Vacancies != nil { vacancies = *req.Vacancies }
	rec := &JobRequisition{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_requisitions (company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, reason)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 RETURNING id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, reason, status, created_at, updated_at`,
		companyID, req.PositionID, req.DepartmentID, requestedBy, req.HiringManagerID,
		req.Title, req.Description, vacancies, req.EmploymentType, req.WorkMode, req.Location,
		req.SalaryMin, req.SalaryMax, req.Currency, req.Reason,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode, &rec.Location,
		&rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Reason, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) GetRequisition(ctx context.Context, companyID, id string) (*JobRequisition, error) {
	rec := &JobRequisition{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, reason, status, created_at, updated_at
		 FROM job_requisitions WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode, &rec.Location,
		&rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Reason, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) ListRequisitions(ctx context.Context, companyID string, filters RecruitmentFilters) ([]JobRequisition, error) {
	query := `SELECT id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, reason, status, created_at, updated_at
		 FROM job_requisitions WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var recs []JobRequisition
	for rows.Next() {
		var rec JobRequisition
		rows.Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
			&rec.Title, &rec.Description, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode, &rec.Location,
			&rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Reason, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
		recs = append(recs, rec)
	}
	return recs, nil
}

func (r *Repository) UpdateRequisition(ctx context.Context, companyID, id string, req *UpdateRequisitionRequest) (*JobRequisition, error) {
	rec := &JobRequisition{}
	err := r.pool.QueryRow(ctx,
		`UPDATE job_requisitions SET
		 title=COALESCE($3,title), description=COALESCE($4,description), vacancies=COALESCE($5,vacancies),
		 employment_type=COALESCE($6,employment_type), work_mode=COALESCE($7,work_mode), location=COALESCE($8,location),
		 salary_min=COALESCE($9,salary_min), salary_max=COALESCE($10,salary_max), currency=COALESCE($11,currency),
		 reason=COALESCE($12,reason), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, reason, status, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.Vacancies, req.EmploymentType, req.WorkMode,
		req.Location, req.SalaryMin, req.SalaryMax, req.Currency, req.Reason,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode, &rec.Location,
		&rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Reason, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *Repository) UpdateRequisitionStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_requisitions SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

// Postings
func (r *Repository) CreatePosting(ctx context.Context, companyID string, req *CreatePostingRequest) (*JobPosting, error) {
	p := &JobPosting{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_postings (company_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, company_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, published_at, closing_at, status, created_at`,
		companyID, req.RequisitionID, req.Title, req.Description, req.Requirements, req.Responsibilities,
		req.Benefits, req.EmploymentType, req.WorkMode, req.Location,
	).Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.Description, &p.Requirements, &p.Responsibilities,
		&p.Benefits, &p.EmploymentType, &p.WorkMode, &p.Location, &p.PublishedAt, &p.ClosingAt, &p.Status, &p.CreatedAt)
	return p, err
}

func (r *Repository) GetPosting(ctx context.Context, companyID, id string) (*JobPosting, error) {
	p := &JobPosting{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, published_at, closing_at, status, created_at
		 FROM job_postings WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.Description, &p.Requirements, &p.Responsibilities,
		&p.Benefits, &p.EmploymentType, &p.WorkMode, &p.Location, &p.PublishedAt, &p.ClosingAt, &p.Status, &p.CreatedAt)
	return p, err
}

func (r *Repository) ListPostings(ctx context.Context, companyID string, filters RecruitmentFilters) ([]JobPosting, error) {
	query := `SELECT id, company_id, requisition_id, title, description, requirements, responsibilities, benefits, employment_type, work_mode, location, published_at, closing_at, status, created_at
		 FROM job_postings WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var postings []JobPosting
	for rows.Next() {
		var p JobPosting
		rows.Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.Description, &p.Requirements, &p.Responsibilities,
			&p.Benefits, &p.EmploymentType, &p.WorkMode, &p.Location, &p.PublishedAt, &p.ClosingAt, &p.Status, &p.CreatedAt)
		postings = append(postings, p)
	}
	return postings, nil
}

func (r *Repository) UpdatePostingStatus(ctx context.Context, companyID, id, status string) error {
	now := time.Now()
	var err error
	if status == "PUBLISHED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3, published_at=$4, created_at=created_at WHERE company_id=$1 AND id=$2`,
			companyID, id, status, now)
	} else if status == "CLOSED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3, closing_at=$4 WHERE company_id=$1 AND id=$2`,
			companyID, id, status, now)
	} else {
		_, err = r.pool.Exec(ctx,
			`UPDATE job_postings SET status=$3 WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	}
	return err
}

// Candidates
func (r *Repository) CreateCandidate(ctx context.Context, companyID string, req *CreateCandidateRequest) (*Candidate, error) {
	c := &Candidate{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidates (company_id, first_name, last_name, email, phone, document_number, location, linkedin_url, portfolio_url, source)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 ON CONFLICT (company_id, email) DO UPDATE SET first_name=$2, last_name=$3, phone=$5, updated_at=NOW()
		 RETURNING id, company_id, first_name, last_name, email, phone, document_number, location, linkedin_url, portfolio_url, source, status, created_at, updated_at`,
		companyID, req.FirstName, req.LastName, req.Email, req.Phone, req.DocumentNumber,
		req.Location, req.LinkedInURL, req.PortfolioURL, req.Source,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.DocumentNumber,
		&c.Location, &c.LinkedInURL, &c.PortfolioURL, &c.Source, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) GetCandidate(ctx context.Context, companyID, id string) (*Candidate, error) {
	c := &Candidate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, first_name, last_name, email, phone, document_number, location, linkedin_url, portfolio_url, source, status, created_at, updated_at
		 FROM candidates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.DocumentNumber,
		&c.Location, &c.LinkedInURL, &c.PortfolioURL, &c.Source, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) ListCandidates(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Candidate, error) {
	query := `SELECT id, company_id, first_name, last_name, email, phone, document_number, location, linkedin_url, portfolio_url, source, status, created_at, updated_at
		 FROM candidates WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.Source != "" {
		query += fmt.Sprintf(" AND source=$%d", argIdx)
		args = append(args, filters.Source)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var candidates []Candidate
	for rows.Next() {
		var c Candidate
		rows.Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.DocumentNumber,
			&c.Location, &c.LinkedInURL, &c.PortfolioURL, &c.Source, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		candidates = append(candidates, c)
	}
	return candidates, nil
}

func (r *Repository) UpdateCandidate(ctx context.Context, companyID, id string, req *UpdateCandidateRequest) (*Candidate, error) {
	c := &Candidate{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidates SET
		 first_name=COALESCE($3,first_name), last_name=COALESCE($4,last_name), phone=COALESCE($5,phone),
		 location=COALESCE($6,location), linkedin_url=COALESCE($7,linkedin_url), portfolio_url=COALESCE($8,portfolio_url),
		 updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, first_name, last_name, email, phone, document_number, location, linkedin_url, portfolio_url, source, status, created_at, updated_at`,
		companyID, id, req.FirstName, req.LastName, req.Phone, req.Location, req.LinkedInURL, req.PortfolioURL,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.DocumentNumber,
		&c.Location, &c.LinkedInURL, &c.PortfolioURL, &c.Source, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

// Applications
func (r *Repository) CreateApplication(ctx context.Context, companyID string, req *CreateApplicationRequest) (*Application, error) {
	a := &Application{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO applications (company_id, candidate_id, job_posting_id)
		 VALUES ($1,$2,$3)
		 RETURNING id, company_id, candidate_id, job_posting_id, status, applied_at, rejected_at, rejection_reason, hired_at`,
		companyID, req.CandidateID, req.JobPostingID,
	).Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.JobPostingID, &a.Status, &a.AppliedAt, &a.RejectedAt, &a.RejectionReason, &a.HiredAt)
	return a, err
}

func (r *Repository) GetApplication(ctx context.Context, companyID, id string) (*Application, error) {
	a := &Application{}
	err := r.pool.QueryRow(ctx,
		`SELECT a.id, a.company_id, a.candidate_id, COALESCE(c.first_name||' '||c.last_name,''),
		 a.job_posting_id, COALESCE(jp.title,''), a.status, a.applied_at, a.rejected_at, a.rejection_reason, a.hired_at
		 FROM applications a
		 LEFT JOIN candidates c ON a.candidate_id=c.id
		 LEFT JOIN job_postings jp ON a.job_posting_id=jp.id
		 WHERE a.company_id=$1 AND a.id=$2`, companyID, id,
	).Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.CandidateName,
		&a.JobPostingID, &a.PostingTitle, &a.Status, &a.AppliedAt, &a.RejectedAt, &a.RejectionReason, &a.HiredAt)
	return a, err
}

func (r *Repository) ListApplications(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Application, error) {
	query := `SELECT a.id, a.company_id, a.candidate_id, COALESCE(c.first_name||' '||c.last_name,''),
		 a.job_posting_id, COALESCE(jp.title,''), a.status, a.applied_at, a.rejected_at, a.rejection_reason, a.hired_at
		 FROM applications a
		 LEFT JOIN candidates c ON a.candidate_id=c.id
		 LEFT JOIN job_postings jp ON a.job_posting_id=jp.id
		 WHERE a.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.CandidateID != "" {
		query += fmt.Sprintf(" AND a.candidate_id=$%d", argIdx)
		args = append(args, filters.CandidateID)
		argIdx++
	}
	if filters.PostingID != "" {
		query += fmt.Sprintf(" AND a.job_posting_id=$%d", argIdx)
		args = append(args, filters.PostingID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND a.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY a.applied_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var apps []Application
	for rows.Next() {
		var a Application
		rows.Scan(&a.ID, &a.CompanyID, &a.CandidateID, &a.CandidateName,
			&a.JobPostingID, &a.PostingTitle, &a.Status, &a.AppliedAt, &a.RejectedAt, &a.RejectionReason, &a.HiredAt)
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *Repository) UpdateApplicationStatus(ctx context.Context, companyID, id, status string) error {
	var err error
	if status == "REJECTED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, rejected_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	} else if status == "HIRED" {
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3, hired_at=NOW() WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	} else {
		_, err = r.pool.Exec(ctx,
			`UPDATE applications SET status=$3 WHERE company_id=$1 AND id=$2`,
			companyID, id, status)
	}
	return err
}

func (r *Repository) RejectApplication(ctx context.Context, companyID, id, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE applications SET status='REJECTED', rejected_at=NOW(), rejection_reason=$3 WHERE company_id=$1 AND id=$2`,
		companyID, id, reason)
	return err
}

// Stage history
func (r *Repository) AddStageHistory(ctx context.Context, companyID, applicationID string, fromStage *string, toStage string, changedBy *string, notes *string) (*CandidateStageHistory, error) {
	h := &CandidateStageHistory{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_stage_history (company_id, application_id, from_stage, to_stage, changed_by, notes)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, application_id, from_stage, to_stage, changed_by, notes, changed_at`,
		companyID, applicationID, fromStage, toStage, changedBy, notes,
	).Scan(&h.ID, &h.CompanyID, &h.ApplicationID, &h.FromStage, &h.ToStage, &h.ChangedBy, &h.Notes, &h.ChangedAt)
	return h, err
}

func (r *Repository) GetStageHistory(ctx context.Context, applicationID string) ([]CandidateStageHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, application_id, from_stage, to_stage, changed_by, notes, changed_at
		 FROM candidate_stage_history WHERE application_id=$1 ORDER BY changed_at ASC`, applicationID)
	if err != nil { return nil, err }
	defer rows.Close()

	var history []CandidateStageHistory
	for rows.Next() {
		var h CandidateStageHistory
		rows.Scan(&h.ID, &h.CompanyID, &h.ApplicationID, &h.FromStage, &h.ToStage, &h.ChangedBy, &h.Notes, &h.ChangedAt)
		history = append(history, h)
	}
	return history, nil
}

// Candidate documents
func (r *Repository) CreateCandidateDocument(ctx context.Context, companyID, candidateID string, docType, fileName, mimeType string, sizeBytes int64, storageProvider, storageKey string) (*CandidateDocument, error) {
	cd := &CandidateDocument{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_documents (company_id, candidate_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, candidate_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key, parsed_data, created_at`,
		companyID, candidateID, docType, fileName, mimeType, sizeBytes, storageProvider, storageKey,
	).Scan(&cd.ID, &cd.CompanyID, &cd.CandidateID, &cd.DocumentType, &cd.FileName, &cd.MimeType, &cd.SizeBytes,
		&cd.StorageProvider, &cd.StorageKey, &cd.ParsedData, &cd.CreatedAt)
	return cd, err
}

func (r *Repository) ListCandidateDocuments(ctx context.Context, candidateID string) ([]CandidateDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, candidate_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key, parsed_data, created_at
		 FROM candidate_documents WHERE candidate_id=$1 ORDER BY created_at DESC`, candidateID)
	if err != nil { return nil, err }
	defer rows.Close()

	var docs []CandidateDocument
	for rows.Next() {
		var d CandidateDocument
		rows.Scan(&d.ID, &d.CompanyID, &d.CandidateID, &d.DocumentType, &d.FileName, &d.MimeType, &d.SizeBytes,
			&d.StorageProvider, &d.StorageKey, &d.ParsedData, &d.CreatedAt)
		docs = append(docs, d)
	}
	return docs, nil
}

// Screening
func (r *Repository) CreateScreeningQuestion(ctx context.Context, companyID string, req *CreateScreeningQuestionRequest) (*ScreeningQuestion, error) {
	q := &ScreeningQuestion{}
	qType := "BOOLEAN"
	if req.QuestionType != nil { qType = *req.QuestionType }
	required := true
	if req.Required != nil { required = *req.Required }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO screening_questions (company_id, job_posting_id, question, question_type, required)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, company_id, job_posting_id, question, question_type, required, sort_order, active, created_at`,
		companyID, req.JobPostingID, req.Question, qType, required,
	).Scan(&q.ID, &q.CompanyID, &q.JobPostingID, &q.Question, &q.QuestionType, &q.Required, &q.SortOrder, &q.Active, &q.CreatedAt)
	return q, err
}

func (r *Repository) ListScreeningQuestions(ctx context.Context, postingID string) ([]ScreeningQuestion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, job_posting_id, question, question_type, required, sort_order, active, created_at
		 FROM screening_questions WHERE job_posting_id=$1 AND active=TRUE ORDER BY sort_order`, postingID)
	if err != nil { return nil, err }
	defer rows.Close()

	var questions []ScreeningQuestion
	for rows.Next() {
		var q ScreeningQuestion
		rows.Scan(&q.ID, &q.CompanyID, &q.JobPostingID, &q.Question, &q.QuestionType, &q.Required, &q.SortOrder, &q.Active, &q.CreatedAt)
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *Repository) CreateScreeningAnswer(ctx context.Context, companyID, applicationID, questionID, answer string) (*ScreeningAnswer, error) {
	sa := &ScreeningAnswer{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO screening_answers (company_id, application_id, question_id, answer)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, application_id, question_id, answer, created_at`,
		companyID, applicationID, questionID, answer,
	).Scan(&sa.ID, &sa.CompanyID, &sa.ApplicationID, &sa.QuestionID, &sa.Answer, &sa.CreatedAt)
	return sa, err
}

// Interviews
func (r *Repository) CreateInterview(ctx context.Context, companyID string, req *CreateInterviewRequest) (*Interview, error) {
	i := &Interview{}
	var scheduledAt *time.Time
	if req.ScheduledAt != nil {
		t, _ := time.Parse(time.RFC3339, *req.ScheduledAt)
		scheduledAt = &t
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interviews (company_id, application_id, interviewer_id, interview_type, scheduled_at, duration_minutes, meeting_url, location)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, application_id, interviewer_id, interview_type, scheduled_at, duration_minutes, meeting_url, location, status, notes, created_at`,
		companyID, req.ApplicationID, req.InterviewerID, req.InterviewType, scheduledAt,
		req.DurationMinutes, req.MeetingURL, req.Location,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewerID, &i.InterviewType,
		&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.Location, &i.Status, &i.Notes, &i.CreatedAt)
	return i, err
}

func (r *Repository) GetInterview(ctx context.Context, companyID, id string) (*Interview, error) {
	i := &Interview{}
	err := r.pool.QueryRow(ctx,
		`SELECT iv.id, iv.company_id, iv.application_id, COALESCE(c.first_name||' '||c.last_name,''),
		 iv.interviewer_id, COALESCE(e.first_name||' '||e.last_name,''),
		 iv.interview_type, iv.scheduled_at, iv.duration_minutes, iv.meeting_url, iv.location, iv.status, iv.notes, iv.created_at
		 FROM interviews iv
		 LEFT JOIN applications a ON iv.application_id=a.id
		 LEFT JOIN candidates c ON a.candidate_id=c.id
		 LEFT JOIN employees e ON iv.interviewer_id=e.id
		 WHERE iv.company_id=$1 AND iv.id=$2`, companyID, id,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.CandidateName,
		&i.InterviewerID, &i.InterviewerName,
		&i.InterviewType, &i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.Location, &i.Status, &i.Notes, &i.CreatedAt)
	return i, err
}

func (r *Repository) ListInterviews(ctx context.Context, companyID string, filters RecruitmentFilters) ([]Interview, error) {
	query := `SELECT iv.id, iv.company_id, iv.application_id, COALESCE(c.first_name||' '||c.last_name,''),
		 iv.interviewer_id, COALESCE(e.first_name||' '||e.last_name,''),
		 iv.interview_type, iv.scheduled_at, iv.duration_minutes, iv.meeting_url, iv.location, iv.status, iv.notes, iv.created_at
		 FROM interviews iv
		 LEFT JOIN applications a ON iv.application_id=a.id
		 LEFT JOIN candidates c ON a.candidate_id=c.id
		 LEFT JOIN employees e ON iv.interviewer_id=e.id
		 WHERE iv.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.ApplicationID != "" {
		query += fmt.Sprintf(" AND iv.application_id=$%d", argIdx)
		args = append(args, filters.ApplicationID)
		argIdx++
	}
	if filters.InterviewerID != "" {
		query += fmt.Sprintf(" AND iv.interviewer_id=$%d", argIdx)
		args = append(args, filters.InterviewerID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND iv.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY iv.scheduled_at DESC NULLS LAST"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var interviews []Interview
	for rows.Next() {
		var i Interview
		rows.Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.CandidateName,
			&i.InterviewerID, &i.InterviewerName,
			&i.InterviewType, &i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.Location, &i.Status, &i.Notes, &i.CreatedAt)
		interviews = append(interviews, i)
	}
	return interviews, nil
}

func (r *Repository) UpdateInterview(ctx context.Context, companyID, id string, req *UpdateInterviewRequest) (*Interview, error) {
	i := &Interview{}
	var scheduledAt *time.Time
	if req.ScheduledAt != nil {
		t, _ := time.Parse(time.RFC3339, *req.ScheduledAt)
		scheduledAt = &t
	}
	err := r.pool.QueryRow(ctx,
		`UPDATE interviews SET
		 scheduled_at=COALESCE($3,scheduled_at), duration_minutes=COALESCE($4,duration_minutes),
		 meeting_url=COALESCE($5,meeting_url), location=COALESCE($6,location), status=COALESCE($7,status)
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, application_id, interviewer_id, interview_type, scheduled_at, duration_minutes, meeting_url, location, status, notes, created_at`,
		companyID, id, scheduledAt, req.DurationMinutes, req.MeetingURL, req.Location, req.Status,
	).Scan(&i.ID, &i.CompanyID, &i.ApplicationID, &i.InterviewerID, &i.InterviewType,
		&i.ScheduledAt, &i.DurationMinutes, &i.MeetingURL, &i.Location, &i.Status, &i.Notes, &i.CreatedAt)
	return i, err
}

// Interview feedback
func (r *Repository) CreateInterviewFeedback(ctx context.Context, companyID, interviewID, interviewerID string, req *CreateInterviewFeedbackRequest) (*InterviewFeedback, error) {
	fb := &InterviewFeedback{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO interview_feedback (company_id, interview_id, interviewer_id, score, comments, recommendation)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, interview_id, interviewer_id, score, comments, recommendation, created_at`,
		companyID, interviewID, interviewerID, req.Score, req.Comments, req.Recommendation,
	).Scan(&fb.ID, &fb.CompanyID, &fb.InterviewID, &fb.InterviewerID, &fb.Score, &fb.Comments, &fb.Recommendation, &fb.CreatedAt)
	return fb, err
}

func (r *Repository) ListInterviewFeedback(ctx context.Context, interviewID string) ([]InterviewFeedback, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, interview_id, interviewer_id, score, comments, recommendation, created_at
		 FROM interview_feedback WHERE interview_id=$1 ORDER BY created_at`, interviewID)
	if err != nil { return nil, err }
	defer rows.Close()

	var feedbacks []InterviewFeedback
	for rows.Next() {
		var fb InterviewFeedback
		rows.Scan(&fb.ID, &fb.CompanyID, &fb.InterviewID, &fb.InterviewerID, &fb.Score, &fb.Comments, &fb.Recommendation, &fb.CreatedAt)
		feedbacks = append(feedbacks, fb)
	}
	return feedbacks, nil
}

// Assessments
func (r *Repository) CreateAssessment(ctx context.Context, companyID string, req *CreateAssessmentRequest) (*Assessment, error) {
	a := &Assessment{}
	aType := "TECHNICAL"
	if req.AssessmentType != nil { aType = *req.AssessmentType }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO assessments (company_id, application_id, assessment_type, title, description, max_score, duration_minutes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, company_id, application_id, assessment_type, title, description, max_score, score, duration_minutes, status, result, created_at`,
		companyID, req.ApplicationID, aType, req.Title, req.Description, req.MaxScore, req.DurationMinutes,
	).Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description, &a.MaxScore,
		&a.Score, &a.DurationMinutes, &a.Status, &a.Result, &a.CreatedAt)
	return a, err
}

func (r *Repository) ListAssessments(ctx context.Context, applicationID string) ([]Assessment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, application_id, assessment_type, title, description, max_score, score, duration_minutes, status, result, created_at
		 FROM assessments WHERE application_id=$1 ORDER BY created_at`, applicationID)
	if err != nil { return nil, err }
	defer rows.Close()

	var assessments []Assessment
	for rows.Next() {
		var a Assessment
		rows.Scan(&a.ID, &a.CompanyID, &a.ApplicationID, &a.AssessmentType, &a.Title, &a.Description, &a.MaxScore,
			&a.Score, &a.DurationMinutes, &a.Status, &a.Result, &a.CreatedAt)
		assessments = append(assessments, a)
	}
	return assessments, nil
}

// Offers
func (r *Repository) CreateOffer(ctx context.Context, companyID, createdBy string, req *CreateOfferRequest) (*JobOffer, error) {
	o := &JobOffer{}
	var startDate *time.Time
	var deadline *time.Time
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &t
	}
	if req.ResponseDeadline != nil {
		t, _ := time.Parse("2006-01-02", *req.ResponseDeadline)
		deadline = &t
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_offers (company_id, application_id, position_title, department_id, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, benefits, conditions, response_deadline, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		 RETURNING id, company_id, application_id, position_title, department_id, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, benefits, conditions, response_deadline, status, created_by, created_at, updated_at`,
		companyID, req.ApplicationID, req.PositionTitle, req.DepartmentID, startDate,
		req.EmploymentType, req.WorkMode, req.SalaryAmount, req.SalaryCurrency, req.SalaryPeriod,
		req.Benefits, req.Conditions, deadline, createdBy,
	).Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.PositionTitle, &o.DepartmentID, &o.StartDate,
		&o.EmploymentType, &o.WorkMode, &o.SalaryAmount, &o.SalaryCurrency, &o.SalaryPeriod,
		&o.Benefits, &o.Conditions, &o.ResponseDeadline, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repository) GetOffer(ctx context.Context, companyID, id string) (*JobOffer, error) {
	o := &JobOffer{}
	err := r.pool.QueryRow(ctx,
		`SELECT jo.id, jo.company_id, jo.application_id, COALESCE(c.first_name||' '||c.last_name,''),
		 jo.position_title, jo.department_id, jo.start_date, jo.employment_type, jo.work_mode,
		 jo.salary_amount, jo.salary_currency, jo.salary_period, jo.benefits, jo.conditions,
		 jo.response_deadline, jo.status, jo.created_by, jo.created_at, jo.updated_at
		 FROM job_offers jo
		 LEFT JOIN applications a ON jo.application_id=a.id
		 LEFT JOIN candidates c ON a.candidate_id=c.id
		 WHERE jo.company_id=$1 AND jo.id=$2`, companyID, id,
	).Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.CandidateName,
		&o.PositionTitle, &o.DepartmentID, &o.StartDate, &o.EmploymentType, &o.WorkMode,
		&o.SalaryAmount, &o.SalaryCurrency, &o.SalaryPeriod, &o.Benefits, &o.Conditions,
		&o.ResponseDeadline, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *Repository) UpdateOfferStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_offers SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

// Referrals
func (r *Repository) CreateReferral(ctx context.Context, companyID, referrerID string, req *CreateReferralRequest) (*EmployeeReferral, error) {
	ref := &EmployeeReferral{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_referrals (company_id, referrer_employee_id, candidate_id, notes)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, referrer_employee_id, candidate_id, application_id, status, reward_status, notes, created_at`,
		companyID, referrerID, req.CandidateID, req.Notes,
	).Scan(&ref.ID, &ref.CompanyID, &ref.ReferrerEmployeeID, &ref.CandidateID, &ref.ApplicationID,
		&ref.Status, &ref.RewardStatus, &ref.Notes, &ref.CreatedAt)
	return ref, err
}

func (r *Repository) ListReferrals(ctx context.Context, companyID string) ([]EmployeeReferral, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT er.id, er.company_id, er.referrer_employee_id, COALESCE(e.first_name||' '||e.last_name,''),
		 er.candidate_id, COALESCE(c.first_name||' '||c.last_name,''),
		 er.application_id, er.status, er.reward_status, er.notes, er.created_at
		 FROM employee_referrals er
		 LEFT JOIN employees e ON er.referrer_employee_id=e.id
		 LEFT JOIN candidates c ON er.candidate_id=c.id
		 WHERE er.company_id=$1 ORDER BY er.created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var refs []EmployeeReferral
	for rows.Next() {
		var ref EmployeeReferral
		rows.Scan(&ref.ID, &ref.CompanyID, &ref.ReferrerEmployeeID, &ref.ReferrerName,
			&ref.CandidateID, &ref.CandidateName,
			&ref.ApplicationID, &ref.Status, &ref.RewardStatus, &ref.Notes, &ref.CreatedAt)
		refs = append(refs, ref)
	}
	return refs, nil
}

// Audit
func (r *Repository) CreateAuditLog(ctx context.Context, companyID, userID, candidateID, entityType, entityID, action string, oldVal, newVal []byte, ipAddress string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO recruitment_audit_log (company_id, user_id, candidate_id, entity_type, entity_id, action, old_value, new_value, ip_address)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		companyID, userID, candidateID, entityType, entityID, action, oldVal, newVal, ipAddress)
	return err
}

// Dashboard
func (r *Repository) GetDashboard(ctx context.Context, companyID string) (*RecruitmentDashboard, error) {
	dash := &RecruitmentDashboard{}

	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM job_requisitions WHERE company_id=$1 AND status IN ('OPEN','APPROVED')`, companyID).Scan(&dash.OpenRequisitions)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM candidates WHERE company_id=$1`, companyID).Scan(&dash.TotalCandidates)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM applications WHERE company_id=$1 AND applied_at >= NOW() - INTERVAL '7 days'`, companyID).Scan(&dash.ApplicationsThisWeek)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM job_offers WHERE company_id=$1 AND status IN ('DRAFT','PENDING_APPROVAL','SENT')`, companyID).Scan(&dash.PendingOffers)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM applications WHERE company_id=$1 AND status='HIRED' AND hired_at >= date_trunc('month', NOW())`, companyID).Scan(&dash.HiresThisMonth)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM interviews WHERE company_id=$1`, companyID).Scan(&dash.TotalInterviews)

	stageRows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM applications WHERE company_id=$1 GROUP BY status ORDER BY COUNT(*) DESC`, companyID)
	if err == nil {
		defer stageRows.Close()
		for stageRows.Next() {
			var sc StageCount
			stageRows.Scan(&sc.Stage, &sc.Count)
			dash.FunnelByStage = append(dash.FunnelByStage, sc)
		}
	}

	return dash, nil
}
