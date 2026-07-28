package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type HiringProcessRepo struct {
	pool *pgxpool.Pool
}

func NewHiringProcessRepo(pool *pgxpool.Pool) *HiringProcessRepo {
	return &HiringProcessRepo{pool: pool}
}

func (r *HiringProcessRepo) Create(ctx context.Context, companyID string, req *domain.HiringProcess) (*domain.HiringProcess, error) {
	h := &domain.HiringProcess{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO hiring_processes (company_id, offer_id, application_id, candidate_id, employee_id, start_date, notes, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, offer_id, application_id, candidate_id, employee_id, status, background_check_status, background_check_result, medical_check_status, medical_check_result, document_verification_status, start_date, onboarding_status, onboarding_id, notes, created_by, created_at, updated_at`,
		companyID, req.OfferID, req.ApplicationID, req.CandidateID, req.EmployeeID,
		req.StartDate, req.Notes, req.CreatedBy,
	).Scan(&h.ID, &h.CompanyID, &h.OfferID, &h.ApplicationID, &h.CandidateID, &h.EmployeeID,
		&h.Status, &h.BackgroundCheckStatus, &h.BackgroundCheckResult, &h.MedicalCheckStatus,
		&h.MedicalCheckResult, &h.DocVerificationStatus, &h.StartDate, &h.OnboardingStatus,
		&h.OnboardingID, &h.Notes, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

func (r *HiringProcessRepo) GetByID(ctx context.Context, companyID, id string) (*domain.HiringProcess, error) {
	h := &domain.HiringProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, offer_id, application_id, candidate_id, employee_id, status, background_check_status, background_check_result, medical_check_status, medical_check_result, document_verification_status, start_date, onboarding_status, onboarding_id, notes, created_by, created_at, updated_at
		 FROM hiring_processes WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&h.ID, &h.CompanyID, &h.OfferID, &h.ApplicationID, &h.CandidateID, &h.EmployeeID,
		&h.Status, &h.BackgroundCheckStatus, &h.BackgroundCheckResult, &h.MedicalCheckStatus,
		&h.MedicalCheckResult, &h.DocVerificationStatus, &h.StartDate, &h.OnboardingStatus,
		&h.OnboardingID, &h.Notes, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

func (r *HiringProcessRepo) ListByCompany(ctx context.Context, companyID, status string) ([]domain.HiringProcess, error) {
	query := `SELECT id, company_id, offer_id, application_id, candidate_id, employee_id, status, background_check_status, background_check_result, medical_check_status, medical_check_result, document_verification_status, start_date, onboarding_status, onboarding_id, notes, created_by, created_at, updated_at
		 FROM hiring_processes WHERE company_id=$1`
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

	var processes []domain.HiringProcess
	for rows.Next() {
		var h domain.HiringProcess
		rows.Scan(&h.ID, &h.CompanyID, &h.OfferID, &h.ApplicationID, &h.CandidateID, &h.EmployeeID,
			&h.Status, &h.BackgroundCheckStatus, &h.BackgroundCheckResult, &h.MedicalCheckStatus,
			&h.MedicalCheckResult, &h.DocVerificationStatus, &h.StartDate, &h.OnboardingStatus,
			&h.OnboardingID, &h.Notes, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
		processes = append(processes, h)
	}
	return processes, nil
}

func (r *HiringProcessRepo) Update(ctx context.Context, companyID, id string, req *domain.HiringProcess) (*domain.HiringProcess, error) {
	h := &domain.HiringProcess{}
	err := r.pool.QueryRow(ctx,
		`UPDATE hiring_processes SET
		 notes=COALESCE($3,notes), start_date=COALESCE($4,start_date),
		 employee_id=COALESCE($5,employee_id), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, offer_id, application_id, candidate_id, employee_id, status, background_check_status, background_check_result, medical_check_status, medical_check_result, document_verification_status, start_date, onboarding_status, onboarding_id, notes, created_by, created_at, updated_at`,
		companyID, id, req.Notes, req.StartDate, req.EmployeeID,
	).Scan(&h.ID, &h.CompanyID, &h.OfferID, &h.ApplicationID, &h.CandidateID, &h.EmployeeID,
		&h.Status, &h.BackgroundCheckStatus, &h.BackgroundCheckResult, &h.MedicalCheckStatus,
		&h.MedicalCheckResult, &h.DocVerificationStatus, &h.StartDate, &h.OnboardingStatus,
		&h.OnboardingID, &h.Notes, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

func (r *HiringProcessRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE hiring_processes SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *HiringProcessRepo) UpdateBackgroundCheck(ctx context.Context, companyID, id, status string, result *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE hiring_processes SET background_check_status=$3, background_check_result=COALESCE($4,background_check_result), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status, result)
	return err
}

func (r *HiringProcessRepo) UpdateMedicalCheck(ctx context.Context, companyID, id, status string, result *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE hiring_processes SET medical_check_status=$3, medical_check_result=COALESCE($4,medical_check_result), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status, result)
	return err
}

func (r *HiringProcessRepo) UpdateDocVerification(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE hiring_processes SET document_verification_status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *HiringProcessRepo) UpdateOnboarding(ctx context.Context, companyID, id, status string, onboardingID *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE hiring_processes SET onboarding_status=$3, onboarding_id=COALESCE($4,onboarding_id), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status, onboardingID)
	return err
}

func (r *HiringProcessRepo) AddTask(ctx context.Context, req *domain.HiringProcessTask) (*domain.HiringProcessTask, error) {
	t := &domain.HiringProcessTask{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO hiring_process_tasks (process_id, task_type, title, description, assigned_to, due_date, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, process_id, task_type, title, description, assigned_to, due_date, status, completed_at, created_at`,
		req.ProcessID, req.TaskType, req.Title, req.Description, req.AssignedTo, req.DueDate, req.Status,
	).Scan(&t.ID, &t.ProcessID, &t.TaskType, &t.Title, &t.Description, &t.AssignedTo, &t.DueDate, &t.Status, &t.CompletedAt, &t.CreatedAt)
	return t, err
}

func (r *HiringProcessRepo) UpdateTask(ctx context.Context, id string, req *domain.HiringProcessTask) (*domain.HiringProcessTask, error) {
	t := &domain.HiringProcessTask{}
	now := time.Now()
	err := r.pool.QueryRow(ctx,
		`UPDATE hiring_process_tasks SET title=COALESCE($2,title), description=COALESCE($3,description),
		 assigned_to=COALESCE($4,assigned_to), due_date=COALESCE($5,due_date),
		 status=COALESCE($6,status), completed_at=CASE WHEN $6='COMPLETED' THEN $7 ELSE completed_at END
		 WHERE id=$1
		 RETURNING id, process_id, task_type, title, description, assigned_to, due_date, status, completed_at, created_at`,
		id, req.Title, req.Description, req.AssignedTo, req.DueDate, req.Status, now,
	).Scan(&t.ID, &t.ProcessID, &t.TaskType, &t.Title, &t.Description, &t.AssignedTo, &t.DueDate, &t.Status, &t.CompletedAt, &t.CreatedAt)
	return t, err
}

func (r *HiringProcessRepo) ListTasks(ctx context.Context, processID string) ([]domain.HiringProcessTask, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, process_id, task_type, title, description, assigned_to, due_date, status, completed_at, created_at
		 FROM hiring_process_tasks WHERE process_id=$1 ORDER BY created_at`, processID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.HiringProcessTask
	for rows.Next() {
		var t domain.HiringProcessTask
		rows.Scan(&t.ID, &t.ProcessID, &t.TaskType, &t.Title, &t.Description, &t.AssignedTo, &t.DueDate, &t.Status, &t.CompletedAt, &t.CreatedAt)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *HiringProcessRepo) AddDocument(ctx context.Context, req *domain.HiringProcessDocument) (*domain.HiringProcessDocument, error) {
	d := &domain.HiringProcessDocument{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO hiring_process_documents (process_id, document_type, file_name, storage_key)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, process_id, document_type, file_name, storage_key, verified, verified_by, verified_at, created_at`,
		req.ProcessID, req.DocumentType, req.FileName, req.StorageKey,
	).Scan(&d.ID, &d.ProcessID, &d.DocumentType, &d.FileName, &d.StorageKey, &d.Verified, &d.VerifiedBy, &d.VerifiedAt, &d.CreatedAt)
	return d, err
}

func (r *HiringProcessRepo) ListDocuments(ctx context.Context, processID string) ([]domain.HiringProcessDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, process_id, document_type, file_name, storage_key, verified, verified_by, verified_at, created_at
		 FROM hiring_process_documents WHERE process_id=$1 ORDER BY created_at`, processID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.HiringProcessDocument
	for rows.Next() {
		var d domain.HiringProcessDocument
		rows.Scan(&d.ID, &d.ProcessID, &d.DocumentType, &d.FileName, &d.StorageKey, &d.Verified, &d.VerifiedBy, &d.VerifiedAt, &d.CreatedAt)
		docs = append(docs, d)
	}
	return docs, nil
}
