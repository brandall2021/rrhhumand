package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type RequisitionRepo struct {
	pool *pgxpool.Pool
}

func NewRequisitionRepo(pool *pgxpool.Pool) *RequisitionRepo {
	return &RequisitionRepo{pool: pool}
}

func (r *RequisitionRepo) Create(ctx context.Context, companyID, requestedBy string, req *domain.Requisition) (*domain.Requisition, error) {
	rec := &domain.Requisition{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_requisitions (company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, justification, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, urgency, reason)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, justification, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, urgency, reason, status, approved_at, opened_at, closed_at, closed_reason, created_at, updated_at`,
		companyID, req.PositionID, req.DepartmentID, requestedBy, req.HiringManagerID,
		req.Title, req.Description, req.Justification, req.Vacancies, req.EmploymentType, req.WorkMode,
		req.Location, req.SalaryMin, req.SalaryMax, req.Currency, req.Urgency, req.Reason,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Justification, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode,
		&rec.Location, &rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Urgency, &rec.Reason,
		&rec.Status, &rec.ApprovedAt, &rec.OpenedAt, &rec.ClosedAt, &rec.ClosedReason, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *RequisitionRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Requisition, error) {
	rec := &domain.Requisition{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, justification, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, urgency, reason, status, approved_at, opened_at, closed_at, closed_reason, created_at, updated_at
		 FROM job_requisitions WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Justification, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode,
		&rec.Location, &rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Urgency, &rec.Reason,
		&rec.Status, &rec.ApprovedAt, &rec.OpenedAt, &rec.ClosedAt, &rec.ClosedReason, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *RequisitionRepo) List(ctx context.Context, companyID string, status string) ([]domain.Requisition, error) {
	query := `SELECT id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, justification, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, urgency, reason, status, approved_at, opened_at, closed_at, closed_reason, created_at, updated_at
		 FROM job_requisitions WHERE company_id=$1`
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

	var recs []domain.Requisition
	for rows.Next() {
		var rec domain.Requisition
		rows.Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
			&rec.Title, &rec.Description, &rec.Justification, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode,
			&rec.Location, &rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Urgency, &rec.Reason,
			&rec.Status, &rec.ApprovedAt, &rec.OpenedAt, &rec.ClosedAt, &rec.ClosedReason, &rec.CreatedAt, &rec.UpdatedAt)
		recs = append(recs, rec)
	}
	return recs, nil
}

func (r *RequisitionRepo) Update(ctx context.Context, companyID, id string, req *domain.Requisition) (*domain.Requisition, error) {
	rec := &domain.Requisition{}
	err := r.pool.QueryRow(ctx,
		`UPDATE job_requisitions SET
		 title=COALESCE($3,title), description=COALESCE($4,description), justification=COALESCE($5,justification),
		 vacancies=COALESCE($6,vacancies), employment_type=COALESCE($7,employment_type), work_mode=COALESCE($8,work_mode),
		 location=COALESCE($9,location), salary_min=COALESCE($10,salary_min), salary_max=COALESCE($11,salary_max),
		 currency=COALESCE($12,currency), reason=COALESCE($13,reason), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, position_id, department_id, requested_by, hiring_manager_id, title, description, justification, vacancies, employment_type, work_mode, location, salary_min, salary_max, currency, urgency, reason, status, approved_at, opened_at, closed_at, closed_reason, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.Justification, req.Vacancies, req.EmploymentType, req.WorkMode,
		req.Location, req.SalaryMin, req.SalaryMax, req.Currency, req.Reason,
	).Scan(&rec.ID, &rec.CompanyID, &rec.PositionID, &rec.DepartmentID, &rec.RequestedBy, &rec.HiringManagerID,
		&rec.Title, &rec.Description, &rec.Justification, &rec.Vacancies, &rec.EmploymentType, &rec.WorkMode,
		&rec.Location, &rec.SalaryMin, &rec.SalaryMax, &rec.Currency, &rec.Urgency, &rec.Reason,
		&rec.Status, &rec.ApprovedAt, &rec.OpenedAt, &rec.ClosedAt, &rec.ClosedReason, &rec.CreatedAt, &rec.UpdatedAt)
	return rec, err
}

func (r *RequisitionRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_requisitions SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *RequisitionRepo) AddSkill(ctx context.Context, req *domain.RequisitionSkill) (*domain.RequisitionSkill, error) {
	s := &domain.RequisitionSkill{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_requisition_skills (requisition_id, skill, category, required, min_years)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, requisition_id, skill, category, required, min_years`,
		req.RequisitionID, req.Skill, req.Category, req.Required, req.MinYears,
	).Scan(&s.ID, &s.RequisitionID, &s.Skill, &s.Category, &s.Required, &s.MinYears)
	return s, err
}

func (r *RequisitionRepo) RemoveSkill(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM job_requisition_skills WHERE id=$1`, id)
	return err
}

func (r *RequisitionRepo) ListSkills(ctx context.Context, requisitionID string) ([]domain.RequisitionSkill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, requisition_id, skill, category, required, min_years
		 FROM job_requisition_skills WHERE requisition_id=$1`, requisitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.RequisitionSkill
	for rows.Next() {
		var s domain.RequisitionSkill
		rows.Scan(&s.ID, &s.RequisitionID, &s.Skill, &s.Category, &s.Required, &s.MinYears)
		skills = append(skills, s)
	}
	return skills, nil
}

func (r *RequisitionRepo) AddApproval(ctx context.Context, req *domain.RequisitionApproval) (*domain.RequisitionApproval, error) {
	a := &domain.RequisitionApproval{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_requisition_approvals (requisition_id, approver_id, step_order, status)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, requisition_id, approver_id, step_order, status, comment, decided_at, created_at`,
		req.RequisitionID, req.ApproverID, req.StepOrder, req.Status,
	).Scan(&a.ID, &a.RequisitionID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
	return a, err
}

func (r *RequisitionRepo) UpdateApproval(ctx context.Context, id string, status string, comment *string, decidedAt time.Time) (*domain.RequisitionApproval, error) {
	a := &domain.RequisitionApproval{}
	err := r.pool.QueryRow(ctx,
		`UPDATE job_requisition_approvals SET status=$2, comment=COALESCE($3,comment), decided_at=$4 WHERE id=$1
		 RETURNING id, requisition_id, approver_id, step_order, status, comment, decided_at, created_at`,
		id, status, comment, decidedAt,
	).Scan(&a.ID, &a.RequisitionID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
	return a, err
}

func (r *RequisitionRepo) ListApprovals(ctx context.Context, requisitionID string) ([]domain.RequisitionApproval, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, requisition_id, approver_id, step_order, status, comment, decided_at, created_at
		 FROM job_requisition_approvals WHERE requisition_id=$1 ORDER BY step_order`, requisitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []domain.RequisitionApproval
	for rows.Next() {
		var a domain.RequisitionApproval
		rows.Scan(&a.ID, &a.RequisitionID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
		approvals = append(approvals, a)
	}
	return approvals, nil
}
