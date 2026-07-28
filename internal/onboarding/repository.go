package onboarding

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

// Templates

func (r *Repository) CreateTemplate(ctx context.Context, t *OnboardingTemplate) error {
	days := 90
	if t.DefaultDurationDays > 0 {
		days = t.DefaultDurationDays
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_templates (company_id, name, description, status, default_duration_days, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, created_at, updated_at`,
		t.CompanyID, t.Name, t.Description, "DRAFT", days, t.CreatedBy,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) GetTemplate(ctx context.Context, companyID, id string) (*OnboardingTemplate, error) {
	t := &OnboardingTemplate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, status, default_duration_days, created_by, created_at, updated_at
		 FROM onboarding_templates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Status, &t.DefaultDurationDays, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *Repository) ListTemplates(ctx context.Context, companyID string) ([]OnboardingTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, status, default_duration_days, created_by, created_at, updated_at
		 FROM onboarding_templates WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []OnboardingTemplate
	for rows.Next() {
		var t OnboardingTemplate
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Status, &t.DefaultDurationDays, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (r *Repository) UpdateTemplate(ctx context.Context, companyID, id string, req *UpdateTemplateRequest) (*OnboardingTemplate, error) {
	t := &OnboardingTemplate{}
	err := r.pool.QueryRow(ctx,
		`UPDATE onboarding_templates SET
		 name=COALESCE($3,name), description=COALESCE($4,description), status=COALESCE($5,status),
		 default_duration_days=COALESCE($6,default_duration_days), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, status, default_duration_days, created_by, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.Status, req.DefaultDurationDays,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Status, &t.DefaultDurationDays, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *Repository) DeleteTemplate(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM onboarding_templates WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

// Template Tasks

func (r *Repository) CreateTemplateTask(ctx context.Context, t *OnboardingTemplateTask) error {
	req := false
	if t.Required {
		req = true
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_template_tasks (template_id, title, description, category, responsible_type, responsible_user_id, required, days_offset, sort_order, estimated_minutes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		t.TemplateID, t.Title, t.Description, t.Category, t.ResponsibleType, t.ResponsibleUserID, req, t.DaysOffset, t.SortOrder, t.EstimatedMinutes,
	).Scan(&t.ID)
}

func (r *Repository) ListTemplateTasks(ctx context.Context, templateID string) ([]OnboardingTemplateTask, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, title, description, category, responsible_type, responsible_user_id, required, days_offset, sort_order, estimated_minutes
		 FROM onboarding_template_tasks WHERE template_id=$1 ORDER BY sort_order`, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []OnboardingTemplateTask
	for rows.Next() {
		var t OnboardingTemplateTask
		if err := rows.Scan(&t.ID, &t.TemplateID, &t.Title, &t.Description, &t.Category, &t.ResponsibleType, &t.ResponsibleUserID, &t.Required, &t.DaysOffset, &t.SortOrder, &t.EstimatedMinutes); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (r *Repository) UpdateTemplateTask(ctx context.Context, id string, req *UpdateTemplateTaskRequest) (*OnboardingTemplateTask, error) {
	t := &OnboardingTemplateTask{}
	var required bool
	if req.Required != nil {
		required = *req.Required
	}
	err := r.pool.QueryRow(ctx,
		`UPDATE onboarding_template_tasks SET
		 title=COALESCE($2,title), description=COALESCE($3,description), category=COALESCE($4,category),
		 responsible_type=COALESCE($5,responsible_type), responsible_user_id=COALESCE($6,responsible_user_id),
		 required=COALESCE($7,required), days_offset=COALESCE($8,days_offset), sort_order=COALESCE($9,sort_order),
		 estimated_minutes=COALESCE($10,estimated_minutes)
		 WHERE id=$1
		 RETURNING id, template_id, title, description, category, responsible_type, responsible_user_id, required, days_offset, sort_order, estimated_minutes`,
		id, req.Title, req.Description, req.Category, req.ResponsibleType, req.ResponsibleUserID, required, req.DaysOffset, req.SortOrder, req.EstimatedMinutes,
	).Scan(&t.ID, &t.TemplateID, &t.Title, &t.Description, &t.Category, &t.ResponsibleType, &t.ResponsibleUserID, &t.Required, &t.DaysOffset, &t.SortOrder, &t.EstimatedMinutes)
	return t, err
}

func (r *Repository) DeleteTemplateTask(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM onboarding_template_tasks WHERE id=$1`, id)
	return err
}

// Processes

func (r *Repository) CreateProcess(ctx context.Context, p *OnboardingProcess) error {
	if p.CompletionPolicy == "" {
		p.CompletionPolicy = "STRICT"
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_processes (company_id, employee_id, template_id, start_date, target_completion_date, status, completion_policy, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, progress_percentage, created_at, updated_at`,
		p.CompanyID, p.EmployeeID, p.TemplateID, p.StartDate, p.TargetCompletionDate, "NOT_STARTED", p.CompletionPolicy, p.CreatedBy,
	).Scan(&p.ID, &p.ProgressPercentage, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) GetProcess(ctx context.Context, companyID, id string) (*OnboardingProcess, error) {
	p := &OnboardingProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, template_id, start_date, target_completion_date, status, progress_percentage, completion_policy, completed_at, cancelled_at, cancellation_reason, created_by, created_at, updated_at
		 FROM onboarding_processes WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.TemplateID, &p.StartDate, &p.TargetCompletionDate, &p.Status, &p.ProgressPercentage, &p.CompletionPolicy, &p.CompletedAt, &p.CancelledAt, &p.CancellationReason, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetProcessByEmployee(ctx context.Context, companyID, employeeID string) (*OnboardingProcess, error) {
	p := &OnboardingProcess{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, template_id, start_date, target_completion_date, status, progress_percentage, completion_policy, completed_at, cancelled_at, cancellation_reason, created_by, created_at, updated_at
		 FROM onboarding_processes WHERE company_id=$1 AND employee_id=$2 ORDER BY created_at DESC LIMIT 1`, companyID, employeeID,
	).Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.TemplateID, &p.StartDate, &p.TargetCompletionDate, &p.Status, &p.ProgressPercentage, &p.CompletionPolicy, &p.CompletedAt, &p.CancelledAt, &p.CancellationReason, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ListProcesses(ctx context.Context, companyID string, filters OnboardingFilters) ([]OnboardingProcess, error) {
	query := `SELECT id, company_id, employee_id, template_id, start_date, target_completion_date, status, progress_percentage, completion_policy, completed_at, cancelled_at, cancellation_reason, created_by, created_at, updated_at
		 FROM onboarding_processes WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ps []OnboardingProcess
	for rows.Next() {
		var p OnboardingProcess
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.TemplateID, &p.StartDate, &p.TargetCompletionDate, &p.Status, &p.ProgressPercentage, &p.CompletionPolicy, &p.CompletedAt, &p.CancelledAt, &p.CancellationReason, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (r *Repository) UpdateProcessStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) UpdateProcessProgress(ctx context.Context, companyID, id string, progress int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET progress_percentage=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, progress)
	return err
}

func (r *Repository) CompleteProcess(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status='COMPLETED', completed_at=NOW(), progress_percentage=100, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) CancelProcess(ctx context.Context, companyID, id, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET status='CANCELLED', cancelled_at=NOW(), cancellation_reason=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, reason)
	return err
}

func (r *Repository) UpdateProcess(ctx context.Context, companyID, id string, req *UpdateOnboardingRequest) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_processes SET
		 template_id=COALESCE($3,template_id), completion_policy=COALESCE($4,completion_policy), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, req.TemplateID, req.CompletionPolicy)
	return err
}

// Tasks

func (r *Repository) CreateTask(ctx context.Context, t *OnboardingTask) error {
	dueDate := t.DueDate
	if dueDate.IsZero() {
		dueDate = time.Now().AddDate(0, 0, 7)
	}
	req := t.Required
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_tasks (onboarding_id, company_id, employee_id, title, description, category, responsible_type, responsible_id, due_date, status, required, sort_order, estimated_minutes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING id, created_at, updated_at`,
		t.OnboardingID, t.CompanyID, t.EmployeeID, t.Title, t.Description, t.Category, t.ResponsibleType, t.ResponsibleID, dueDate, "PENDING", req, t.SortOrder, t.EstimatedMinutes,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) GetTask(ctx context.Context, companyID, taskID string) (*OnboardingTask, error) {
	t := &OnboardingTask{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, title, description, category, responsible_type, responsible_id, due_date, status, required, sort_order, estimated_minutes, completed_at, started_at, blocked_reason, created_at, updated_at
		 FROM onboarding_tasks WHERE company_id=$1 AND id=$2`, companyID, taskID,
	).Scan(&t.ID, &t.OnboardingID, &t.CompanyID, &t.EmployeeID, &t.Title, &t.Description, &t.Category, &t.ResponsibleType, &t.ResponsibleID, &t.DueDate, &t.Status, &t.Required, &t.SortOrder, &t.EstimatedMinutes, &t.CompletedAt, &t.StartedAt, &t.BlockedReason, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *Repository) ListTasks(ctx context.Context, onboardingID string) ([]OnboardingTask, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, title, description, category, responsible_type, responsible_id, due_date, status, required, sort_order, estimated_minutes, completed_at, started_at, blocked_reason, created_at, updated_at
		 FROM onboarding_tasks WHERE onboarding_id=$1 ORDER BY sort_order`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []OnboardingTask
	for rows.Next() {
		var t OnboardingTask
		if err := rows.Scan(&t.ID, &t.OnboardingID, &t.CompanyID, &t.EmployeeID, &t.Title, &t.Description, &t.Category, &t.ResponsibleType, &t.ResponsibleID, &t.DueDate, &t.Status, &t.Required, &t.SortOrder, &t.EstimatedMinutes, &t.CompletedAt, &t.StartedAt, &t.BlockedReason, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func (r *Repository) UpdateTask(ctx context.Context, companyID, id string, req *UpdateTaskRequest) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_tasks SET
		 title=COALESCE($3,title), description=COALESCE($4,description), category=COALESCE($5,category),
		 responsible_type=COALESCE($6,responsible_type), responsible_id=COALESCE($7,responsible_id),
		 required=COALESCE($8,required), sort_order=COALESCE($9,sort_order),
		 estimated_minutes=COALESCE($10,estimated_minutes), status=COALESCE($11,status),
		 blocked_reason=COALESCE($12,blocked_reason), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, req.Title, req.Description, req.Category, req.ResponsibleType, req.ResponsibleID, req.Required, req.SortOrder, req.EstimatedMinutes, req.Status, req.BlockedReason)
	return err
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_tasks SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) CompleteTask(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_tasks SET status='COMPLETED', completed_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) StartTask(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_tasks SET status='IN_PROGRESS', started_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) BlockTask(ctx context.Context, companyID, id, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_tasks SET status='BLOCKED', blocked_reason=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, reason)
	return err
}

func (r *Repository) GetRequiredTaskCount(ctx context.Context, onboardingID string) (total, completed int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(*) FILTER (WHERE status='COMPLETED') FROM onboarding_tasks WHERE onboarding_id=$1 AND required=true`,
		onboardingID).Scan(&total, &completed)
	return
}

func (r *Repository) GetTaskCountByStatus(ctx context.Context, companyID string) (dueToday int, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_tasks WHERE company_id=$1 AND due_date=CURRENT_DATE AND status NOT IN ('COMPLETED','CANCELLED')`,
		companyID).Scan(&dueToday)
	return
}

// Documents

func (r *Repository) CreateDocument(ctx context.Context, d *OnboardingDocument) error {
	req := d.Required
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_documents (onboarding_id, company_id, employee_id, document_type, status, required)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		d.OnboardingID, d.CompanyID, d.EmployeeID, d.DocumentType, "REQUIRED", req,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *Repository) GetDocument(ctx context.Context, companyID, id string) (*OnboardingDocument, error) {
	d := &OnboardingDocument{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, document_type, file_name, mime_type, size_bytes, checksum, storage_provider, storage_key, status, required, rejection_reason, uploaded_at, reviewed_at, reviewed_by, created_at, updated_at
		 FROM onboarding_documents WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&d.ID, &d.OnboardingID, &d.CompanyID, &d.EmployeeID, &d.DocumentType, &d.FileName, &d.MimeType, &d.SizeBytes, &d.Checksum, &d.StorageProvider, &d.StorageKey, &d.Status, &d.Required, &d.RejectionReason, &d.UploadedAt, &d.ReviewedAt, &d.ReviewedBy, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

func (r *Repository) ListDocuments(ctx context.Context, onboardingID string) ([]OnboardingDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, document_type, file_name, mime_type, size_bytes, checksum, storage_provider, storage_key, status, required, rejection_reason, uploaded_at, reviewed_at, reviewed_by, created_at, updated_at
		 FROM onboarding_documents WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ds []OnboardingDocument
	for rows.Next() {
		var d OnboardingDocument
		if err := rows.Scan(&d.ID, &d.OnboardingID, &d.CompanyID, &d.EmployeeID, &d.DocumentType, &d.FileName, &d.MimeType, &d.SizeBytes, &d.Checksum, &d.StorageProvider, &d.StorageKey, &d.Status, &d.Required, &d.RejectionReason, &d.UploadedAt, &d.ReviewedAt, &d.ReviewedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *Repository) UpdateDocumentStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) UpdateDocumentAfterUpload(ctx context.Context, companyID, id, fileName, mimeType, storageProvider, storageKey, checksum string, size int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET file_name=$3, mime_type=$4, size_bytes=$5, storage_provider=$6, storage_key=$7, checksum=$8, status='UPLOADED', uploaded_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, fileName, mimeType, size, storageProvider, storageKey, checksum)
	return err
}

func (r *Repository) ApproveDocument(ctx context.Context, companyID, id, reviewedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET status='APPROVED', reviewed_by=$3, reviewed_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, reviewedBy)
	return err
}

func (r *Repository) RejectDocument(ctx context.Context, companyID, id, reviewedBy, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_documents SET status='REJECTED', reviewed_by=$3, rejection_reason=$4, reviewed_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, reviewedBy, reason)
	return err
}

func (r *Repository) GetPendingReviewDocumentCount(ctx context.Context, companyID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_documents WHERE company_id=$1 AND status='UPLOADED'`, companyID).Scan(&count)
	return count, err
}

// Assets

func (r *Repository) CreateAsset(ctx context.Context, a *OnboardingAsset) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_assets (onboarding_id, company_id, employee_id, asset_type, description, serial_number, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`,
		a.OnboardingID, a.CompanyID, a.EmployeeID, a.AssetType, a.Description, a.SerialNumber, "REQUESTED", a.Notes,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *Repository) GetAsset(ctx context.Context, companyID, id string) (*OnboardingAsset, error) {
	a := &OnboardingAsset{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, asset_type, description, serial_number, status, assigned_by, assigned_at, delivered_at, returned_at, notes, created_at, updated_at
		 FROM onboarding_assets WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&a.ID, &a.OnboardingID, &a.CompanyID, &a.EmployeeID, &a.AssetType, &a.Description, &a.SerialNumber, &a.Status, &a.AssignedBy, &a.AssignedAt, &a.DeliveredAt, &a.ReturnedAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *Repository) ListAssets(ctx context.Context, onboardingID string) ([]OnboardingAsset, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, asset_type, description, serial_number, status, assigned_by, assigned_at, delivered_at, returned_at, notes, created_at, updated_at
		 FROM onboarding_assets WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []OnboardingAsset
	for rows.Next() {
		var a OnboardingAsset
		if err := rows.Scan(&a.ID, &a.OnboardingID, &a.CompanyID, &a.EmployeeID, &a.AssetType, &a.Description, &a.SerialNumber, &a.Status, &a.AssignedBy, &a.AssignedAt, &a.DeliveredAt, &a.ReturnedAt, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func (r *Repository) UpdateAssetStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_assets SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) AssignAsset(ctx context.Context, companyID, id, assignedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_assets SET status='ASSIGNED', assigned_by=$3, assigned_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, assignedBy)
	return err
}

func (r *Repository) DeliverAsset(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_assets SET status='DELIVERED', delivered_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) ReturnAsset(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_assets SET status='RETURNED', returned_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

// Access Requests

func (r *Repository) CreateAccessRequest(ctx context.Context, ar *AccessRequest) error {
	if ar.AccessType == "" {
		ar.AccessType = "STANDARD"
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO access_requests (onboarding_id, company_id, employee_id, system_name, access_type, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, requested_at, created_at, updated_at`,
		ar.OnboardingID, ar.CompanyID, ar.EmployeeID, ar.SystemName, ar.AccessType, "REQUESTED", ar.Notes,
	).Scan(&ar.ID, &ar.RequestedAt, &ar.CreatedAt, &ar.UpdatedAt)
}

func (r *Repository) GetAccessRequest(ctx context.Context, companyID, id string) (*AccessRequest, error) {
	ar := &AccessRequest{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, system_name, access_type, status, requested_at, approved_at, approved_by, activated_at, revoked_at, notes, created_at, updated_at
		 FROM access_requests WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&ar.ID, &ar.OnboardingID, &ar.CompanyID, &ar.EmployeeID, &ar.SystemName, &ar.AccessType, &ar.Status, &ar.RequestedAt, &ar.ApprovedAt, &ar.ApprovedBy, &ar.ActivatedAt, &ar.RevokedAt, &ar.Notes, &ar.CreatedAt, &ar.UpdatedAt)
	return ar, err
}

func (r *Repository) ListAccessRequests(ctx context.Context, onboardingID string) ([]AccessRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, system_name, access_type, status, requested_at, approved_at, approved_by, activated_at, revoked_at, notes, created_at, updated_at
		 FROM access_requests WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ars []AccessRequest
	for rows.Next() {
		var ar AccessRequest
		if err := rows.Scan(&ar.ID, &ar.OnboardingID, &ar.CompanyID, &ar.EmployeeID, &ar.SystemName, &ar.AccessType, &ar.Status, &ar.RequestedAt, &ar.ApprovedAt, &ar.ApprovedBy, &ar.ActivatedAt, &ar.RevokedAt, &ar.Notes, &ar.CreatedAt, &ar.UpdatedAt); err != nil {
			return nil, err
		}
		ars = append(ars, ar)
	}
	return ars, nil
}

func (r *Repository) ApproveAccess(ctx context.Context, companyID, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE access_requests SET status='APPROVED', approved_by=$3, approved_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id, approvedBy)
	return err
}

func (r *Repository) RejectAccess(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE access_requests SET status='REJECTED', updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) ActivateAccess(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE access_requests SET status='ACTIVATED', activated_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) RevokeAccess(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE access_requests SET status='REVOKED', revoked_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

// Milestones

func (r *Repository) CreateMilestone(ctx context.Context, m *OnboardingMilestone) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_milestones (onboarding_id, company_id, employee_id, milestone_type, title, description, days_offset, due_date, responsible_type, responsible_id, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`,
		m.OnboardingID, m.CompanyID, m.EmployeeID, m.MilestoneType, m.Title, m.Description, m.DaysOffset, m.DueDate, m.ResponsibleType, m.ResponsibleID, "PENDING",
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *Repository) ListMilestones(ctx context.Context, onboardingID string) ([]OnboardingMilestone, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, milestone_type, title, description, days_offset, due_date, responsible_type, responsible_id, status, completed_at, created_at, updated_at
		 FROM onboarding_milestones WHERE onboarding_id=$1 ORDER BY days_offset`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []OnboardingMilestone
	for rows.Next() {
		var m OnboardingMilestone
		if err := rows.Scan(&m.ID, &m.OnboardingID, &m.CompanyID, &m.EmployeeID, &m.MilestoneType, &m.Title, &m.Description, &m.DaysOffset, &m.DueDate, &m.ResponsibleType, &m.ResponsibleID, &m.Status, &m.CompletedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

func (r *Repository) UpdateMilestone(ctx context.Context, companyID, id string, req *UpdateMilestoneRequest) (*OnboardingMilestone, error) {
	m := &OnboardingMilestone{}
	err := r.pool.QueryRow(ctx,
		`UPDATE onboarding_milestones SET
		 title=COALESCE($3,title), description=COALESCE($4,description), days_offset=COALESCE($5,days_offset),
		 responsible_type=COALESCE($6,responsible_type), responsible_id=COALESCE($7,responsible_id),
		 status=COALESCE($8,status), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, onboarding_id, company_id, employee_id, milestone_type, title, description, days_offset, due_date, responsible_type, responsible_id, status, completed_at, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.DaysOffset, req.ResponsibleType, req.ResponsibleID, req.Status,
	).Scan(&m.ID, &m.OnboardingID, &m.CompanyID, &m.EmployeeID, &m.MilestoneType, &m.Title, &m.Description, &m.DaysOffset, &m.DueDate, &m.ResponsibleType, &m.ResponsibleID, &m.Status, &m.CompletedAt, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *Repository) CompleteMilestone(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_milestones SET status='COMPLETED', completed_at=NOW(), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) GetNextMilestone(ctx context.Context, onboardingID string) (*OnboardingMilestone, error) {
	m := &OnboardingMilestone{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, milestone_type, title, description, days_offset, due_date, responsible_type, responsible_id, status, completed_at, created_at, updated_at
		 FROM onboarding_milestones WHERE onboarding_id=$1 AND status='PENDING' ORDER BY due_date ASC LIMIT 1`, onboardingID,
	).Scan(&m.ID, &m.OnboardingID, &m.CompanyID, &m.EmployeeID, &m.MilestoneType, &m.Title, &m.Description, &m.DaysOffset, &m.DueDate, &m.ResponsibleType, &m.ResponsibleID, &m.Status, &m.CompletedAt, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

// Feedback

func (r *Repository) CreateFeedback(ctx context.Context, f *OnboardingFeedback) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_feedback (onboarding_id, company_id, employee_id, feedback_type, submitted_by, adaptation_score, team_score, knowledge_score, communication_score, overall_score, comments)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, submitted_at, created_at`,
		f.OnboardingID, f.CompanyID, f.EmployeeID, f.FeedbackType, f.SubmittedBy, f.AdaptationScore, f.TeamScore, f.KnowledgeScore, f.CommunicationScore, f.OverallScore, f.Comments,
	).Scan(&f.ID, &f.SubmittedAt, &f.CreatedAt)
}

func (r *Repository) ListFeedback(ctx context.Context, onboardingID string) ([]OnboardingFeedback, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, feedback_type, submitted_by, adaptation_score, team_score, knowledge_score, communication_score, overall_score, comments, submitted_at, created_at
		 FROM onboarding_feedback WHERE onboarding_id=$1 ORDER BY created_at DESC`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fs []OnboardingFeedback
	for rows.Next() {
		var f OnboardingFeedback
		if err := rows.Scan(&f.ID, &f.OnboardingID, &f.CompanyID, &f.EmployeeID, &f.FeedbackType, &f.SubmittedBy, &f.AdaptationScore, &f.TeamScore, &f.KnowledgeScore, &f.CommunicationScore, &f.OverallScore, &f.Comments, &f.SubmittedAt, &f.CreatedAt); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}

// Buddies

func (r *Repository) AssignBuddy(ctx context.Context, b *OnboardingBuddy) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_buddies (onboarding_id, company_id, employee_id, buddy_employee_id, start_date, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at, updated_at`,
		b.OnboardingID, b.CompanyID, b.EmployeeID, b.BuddyEmployeeID, b.StartDate, "ACTIVE", b.Notes,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *Repository) GetBuddy(ctx context.Context, onboardingID string) (*OnboardingBuddy, error) {
	b := &OnboardingBuddy{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, buddy_employee_id, start_date, end_date, status, notes, created_at, updated_at
		 FROM onboarding_buddies WHERE onboarding_id=$1 ORDER BY created_at DESC LIMIT 1`, onboardingID,
	).Scan(&b.ID, &b.OnboardingID, &b.CompanyID, &b.EmployeeID, &b.BuddyEmployeeID, &b.StartDate, &b.EndDate, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *Repository) ListBuddies(ctx context.Context, onboardingID string) ([]OnboardingBuddy, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, buddy_employee_id, start_date, end_date, status, notes, created_at, updated_at
		 FROM onboarding_buddies WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bs []OnboardingBuddy
	for rows.Next() {
		var b OnboardingBuddy
		if err := rows.Scan(&b.ID, &b.OnboardingID, &b.CompanyID, &b.EmployeeID, &b.BuddyEmployeeID, &b.StartDate, &b.EndDate, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

func (r *Repository) CompleteBuddy(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE onboarding_buddies SET status='COMPLETED', end_date=CURRENT_DATE, updated_at=NOW()
		 WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

// Exceptions

func (r *Repository) CreateException(ctx context.Context, e *OnboardingException) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO onboarding_exceptions (onboarding_id, company_id, entity_type, entity_id, reason, created_by, expires_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at`,
		e.OnboardingID, e.CompanyID, e.EntityType, e.EntityID, e.Reason, e.CreatedBy, e.ExpiresAt,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *Repository) ListExceptions(ctx context.Context, onboardingID string) ([]OnboardingException, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, entity_type, entity_id, reason, created_by, expires_at, created_at
		 FROM onboarding_exceptions WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var es []OnboardingException
	for rows.Next() {
		var e OnboardingException
		if err := rows.Scan(&e.ID, &e.OnboardingID, &e.CompanyID, &e.EntityType, &e.EntityID, &e.Reason, &e.CreatedBy, &e.ExpiresAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

// Training Assignments

func (r *Repository) CreateTraining(ctx context.Context, t *TrainingAssignment) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO training_assignments (onboarding_id, company_id, employee_id, course_name, description, training_type, status, due_date, external_provider, external_course_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, created_at, updated_at`,
		t.OnboardingID, t.CompanyID, t.EmployeeID, t.CourseName, t.Description, t.TrainingType, "ASSIGNED", t.DueDate, t.ExternalProvider, t.ExternalCourseID,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) ListTraining(ctx context.Context, onboardingID string) ([]TrainingAssignment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, onboarding_id, company_id, employee_id, course_name, description, training_type, status, due_date, completed_at, external_provider, external_course_id, created_at, updated_at
		 FROM training_assignments WHERE onboarding_id=$1 ORDER BY created_at`, onboardingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []TrainingAssignment
	for rows.Next() {
		var t TrainingAssignment
		if err := rows.Scan(&t.ID, &t.OnboardingID, &t.CompanyID, &t.EmployeeID, &t.CourseName, &t.Description, &t.TrainingType, &t.Status, &t.DueDate, &t.CompletedAt, &t.ExternalProvider, &t.ExternalCourseID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

// Domain Events

func (r *Repository) CreateEvent(ctx context.Context, e *DomainEvent) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO domain_events (event_type, company_id, aggregate_id, aggregate_type, payload)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`,
		e.EventType, e.CompanyID, e.AggregateID, e.AggregateType, e.Payload,
	).Scan(&e.ID, &e.CreatedAt)
}

// Notifications

func (r *Repository) CreateNotification(ctx context.Context, n *Notification) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO notifications (company_id, user_id, title, body, notification_type, channel, reference_type, reference_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at`,
		n.CompanyID, n.UserID, n.Title, n.Body, n.NotificationType, n.Channel, n.ReferenceType, n.ReferenceID,
	).Scan(&n.ID, &n.CreatedAt)
}

// Audit Log

func (r *Repository) CreateAuditLog(ctx context.Context, a *OnboardingAuditLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO onboarding_audit_log (company_id, user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		a.CompanyID, a.UserID, a.Action, a.EntityType, a.EntityID, a.OldValue, a.NewValue, a.IPAddress, a.UserAgent)
	return err
}

// Dashboard

func (r *Repository) GetDashboard(ctx context.Context, companyID string) (*OnboardingDashboard, error) {
	d := &OnboardingDashboard{}

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='IN_PROGRESS'`, companyID).Scan(&d.Active)
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='NOT_STARTED'`, companyID).Scan(&d.Pending)
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status='COMPLETED'`, companyID).Scan(&d.Completed)
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_processes WHERE company_id=$1 AND status IN ('IN_PROGRESS','NOT_STARTED') AND target_completion_date < CURRENT_DATE`, companyID).Scan(&d.Overdue)
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(progress_percentage),0) FROM onboarding_processes WHERE company_id=$1 AND status IN ('IN_PROGRESS','NOT_STARTED')`, companyID).Scan(&d.AverageProgress)

	d.TasksDueToday, _ = r.GetTaskCountByStatus(ctx, companyID)
	d.DocumentsPendingReview, _ = r.GetPendingReviewDocumentCount(ctx, companyID)

	return d, nil
}

func (r *Repository) GetEmployeeDashboard(ctx context.Context, companyID, employeeID string) (*EmployeeDashboard, error) {
	d := &EmployeeDashboard{}

	p, err := r.GetProcessByEmployee(ctx, companyID, employeeID)
	if err != nil {
		return nil, err
	}

	d.Status = p.Status
	d.Progress = p.ProgressPercentage

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_tasks WHERE onboarding_id=$1`, p.ID).Scan(&d.TasksTotal)
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_tasks WHERE onboarding_id=$1 AND status='COMPLETED'`, p.ID).Scan(&d.TasksCompleted)
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM onboarding_tasks WHERE onboarding_id=$1 AND status NOT IN ('COMPLETED','CANCELLED')`, p.ID).Scan(&d.PendingTasks)

	m, err := r.GetNextMilestone(ctx, p.ID)
	if err == nil {
		d.NextMilestone = &MilestoneRef{
			Name: m.Title,
			Date: m.DueDate,
		}
	}

	return d, nil
}
