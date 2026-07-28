package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type OfferRepo struct {
	pool *pgxpool.Pool
}

func NewOfferRepo(pool *pgxpool.Pool) *OfferRepo {
	return &OfferRepo{pool: pool}
}

func (r *OfferRepo) Create(ctx context.Context, companyID string, req *domain.Offer) (*domain.Offer, error) {
	o := &domain.Offer{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO offers (company_id, application_id, position_title, department_id, offer_type, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, variable_compensation, benefits_summary, equity_terms, conditions, notes, response_deadline, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		 RETURNING id, company_id, application_id, position_title, department_id, offer_type, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, variable_compensation, benefits_summary, equity_terms, conditions, notes, response_deadline, status, sent_at, accepted_at, rejected_at, rejection_reason, expired_at, created_by, created_at, updated_at`,
		companyID, req.ApplicationID, req.PositionTitle, req.DepartmentID, req.OfferType,
		req.StartDate, req.EmploymentType, req.WorkMode, req.SalaryAmount, req.SalaryCurrency,
		req.SalaryPeriod, req.VariableComp, req.BenefitsSummary, req.EquityTerms, req.Conditions,
		req.Notes, req.ResponseDeadline, req.CreatedBy,
	).Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.PositionTitle, &o.DepartmentID, &o.OfferType,
		&o.StartDate, &o.EmploymentType, &o.WorkMode, &o.SalaryAmount, &o.SalaryCurrency,
		&o.SalaryPeriod, &o.VariableComp, &o.BenefitsSummary, &o.EquityTerms, &o.Conditions,
		&o.Notes, &o.ResponseDeadline, &o.Status, &o.SentAt, &o.AcceptedAt, &o.RejectedAt,
		&o.RejectionReason, &o.ExpiredAt, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *OfferRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Offer, error) {
	o := &domain.Offer{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, application_id, position_title, department_id, offer_type, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, variable_compensation, benefits_summary, equity_terms, conditions, notes, response_deadline, status, sent_at, accepted_at, rejected_at, rejection_reason, expired_at, created_by, created_at, updated_at
		 FROM offers WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.PositionTitle, &o.DepartmentID, &o.OfferType,
		&o.StartDate, &o.EmploymentType, &o.WorkMode, &o.SalaryAmount, &o.SalaryCurrency,
		&o.SalaryPeriod, &o.VariableComp, &o.BenefitsSummary, &o.EquityTerms, &o.Conditions,
		&o.Notes, &o.ResponseDeadline, &o.Status, &o.SentAt, &o.AcceptedAt, &o.RejectedAt,
		&o.RejectionReason, &o.ExpiredAt, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *OfferRepo) List(ctx context.Context, companyID string, applicationID, status string) ([]domain.Offer, error) {
	query := `SELECT id, company_id, application_id, position_title, department_id, offer_type, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, variable_compensation, benefits_summary, equity_terms, conditions, notes, response_deadline, status, sent_at, accepted_at, rejected_at, rejection_reason, expired_at, created_by, created_at, updated_at
		 FROM offers WHERE company_id=$1`
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

	var offers []domain.Offer
	for rows.Next() {
		var o domain.Offer
		rows.Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.PositionTitle, &o.DepartmentID, &o.OfferType,
			&o.StartDate, &o.EmploymentType, &o.WorkMode, &o.SalaryAmount, &o.SalaryCurrency,
			&o.SalaryPeriod, &o.VariableComp, &o.BenefitsSummary, &o.EquityTerms, &o.Conditions,
			&o.Notes, &o.ResponseDeadline, &o.Status, &o.SentAt, &o.AcceptedAt, &o.RejectedAt,
			&o.RejectionReason, &o.ExpiredAt, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
		offers = append(offers, o)
	}
	return offers, nil
}

func (r *OfferRepo) Update(ctx context.Context, companyID, id string, req *domain.Offer) (*domain.Offer, error) {
	o := &domain.Offer{}
	err := r.pool.QueryRow(ctx,
		`UPDATE offers SET
		 position_title=COALESCE($3,position_title), offer_type=COALESCE($4,offer_type),
		 start_date=COALESCE($5,start_date), employment_type=COALESCE($6,employment_type),
		 work_mode=COALESCE($7,work_mode), salary_amount=COALESCE($8,salary_amount),
		 salary_currency=COALESCE($9,salary_currency), salary_period=COALESCE($10,salary_period),
		 variable_compensation=COALESCE($11,variable_compensation),
		 benefits_summary=COALESCE($12,benefits_summary), equity_terms=COALESCE($13,equity_terms),
		 conditions=COALESCE($14,conditions), notes=COALESCE($15,notes),
		 response_deadline=COALESCE($16,response_deadline), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, application_id, position_title, department_id, offer_type, start_date, employment_type, work_mode, salary_amount, salary_currency, salary_period, variable_compensation, benefits_summary, equity_terms, conditions, notes, response_deadline, status, sent_at, accepted_at, rejected_at, rejection_reason, expired_at, created_by, created_at, updated_at`,
		companyID, id, req.PositionTitle, req.OfferType, req.StartDate, req.EmploymentType,
		req.WorkMode, req.SalaryAmount, req.SalaryCurrency, req.SalaryPeriod, req.VariableComp,
		req.BenefitsSummary, req.EquityTerms, req.Conditions, req.Notes, req.ResponseDeadline,
	).Scan(&o.ID, &o.CompanyID, &o.ApplicationID, &o.PositionTitle, &o.DepartmentID, &o.OfferType,
		&o.StartDate, &o.EmploymentType, &o.WorkMode, &o.SalaryAmount, &o.SalaryCurrency,
		&o.SalaryPeriod, &o.VariableComp, &o.BenefitsSummary, &o.EquityTerms, &o.Conditions,
		&o.Notes, &o.ResponseDeadline, &o.Status, &o.SentAt, &o.AcceptedAt, &o.RejectedAt,
		&o.RejectionReason, &o.ExpiredAt, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *OfferRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE offers SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *OfferRepo) AddApproval(ctx context.Context, req *domain.OfferApproval) (*domain.OfferApproval, error) {
	a := &domain.OfferApproval{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO offer_approvals (offer_id, approver_id, step_order, status)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, offer_id, approver_id, step_order, status, comment, decided_at, created_at`,
		req.OfferID, req.ApproverID, req.StepOrder, req.Status,
	).Scan(&a.ID, &a.OfferID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
	return a, err
}

func (r *OfferRepo) UpdateApproval(ctx context.Context, id, status string, comment *string) (*domain.OfferApproval, error) {
	a := &domain.OfferApproval{}
	err := r.pool.QueryRow(ctx,
		`UPDATE offer_approvals SET status=$2, comment=COALESCE($3,comment), decided_at=NOW() WHERE id=$1
		 RETURNING id, offer_id, approver_id, step_order, status, comment, decided_at, created_at`,
		id, status, comment,
	).Scan(&a.ID, &a.OfferID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
	return a, err
}

func (r *OfferRepo) ListApprovals(ctx context.Context, offerID string) ([]domain.OfferApproval, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, offer_id, approver_id, step_order, status, comment, decided_at, created_at
		 FROM offer_approvals WHERE offer_id=$1 ORDER BY step_order`, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []domain.OfferApproval
	for rows.Next() {
		var a domain.OfferApproval
		rows.Scan(&a.ID, &a.OfferID, &a.ApproverID, &a.StepOrder, &a.Status, &a.Comment, &a.DecidedAt, &a.CreatedAt)
		approvals = append(approvals, a)
	}
	return approvals, nil
}

func (r *OfferRepo) AddNegotiation(ctx context.Context, req *domain.OfferNegotiation) (*domain.OfferNegotiation, error) {
	n := &domain.OfferNegotiation{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO offer_negotiations (offer_id, requested_by, field, original_value, requested_value, counter_value, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, offer_id, requested_by, field, original_value, requested_value, counter_value, status, notes, created_at, resolved_at`,
		req.OfferID, req.RequestedBy, req.Field, req.OriginalValue, req.RequestedValue,
		req.CounterValue, req.Status, req.Notes,
	).Scan(&n.ID, &n.OfferID, &n.RequestedBy, &n.Field, &n.OriginalValue, &n.RequestedValue,
		&n.CounterValue, &n.Status, &n.Notes, &n.CreatedAt, &n.ResolvedAt)
	return n, err
}

func (r *OfferRepo) UpdateNegotiation(ctx context.Context, id string, req *domain.OfferNegotiation) (*domain.OfferNegotiation, error) {
	n := &domain.OfferNegotiation{}
	err := r.pool.QueryRow(ctx,
		`UPDATE offer_negotiations SET counter_value=COALESCE($2,counter_value), status=COALESCE($3,status),
		 notes=COALESCE($4,notes), resolved_at=CASE WHEN $3 IN ('ACCEPTED','REJECTED') THEN NOW() ELSE resolved_at END
		 WHERE id=$1
		 RETURNING id, offer_id, requested_by, field, original_value, requested_value, counter_value, status, notes, created_at, resolved_at`,
		id, req.CounterValue, req.Status, req.Notes,
	).Scan(&n.ID, &n.OfferID, &n.RequestedBy, &n.Field, &n.OriginalValue, &n.RequestedValue,
		&n.CounterValue, &n.Status, &n.Notes, &n.CreatedAt, &n.ResolvedAt)
	return n, err
}

func (r *OfferRepo) ListNegotiations(ctx context.Context, offerID string) ([]domain.OfferNegotiation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, offer_id, requested_by, field, original_value, requested_value, counter_value, status, notes, created_at, resolved_at
		 FROM offer_negotiations WHERE offer_id=$1 ORDER BY created_at`, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var negotiations []domain.OfferNegotiation
	for rows.Next() {
		var n domain.OfferNegotiation
		rows.Scan(&n.ID, &n.OfferID, &n.RequestedBy, &n.Field, &n.OriginalValue, &n.RequestedValue,
			&n.CounterValue, &n.Status, &n.Notes, &n.CreatedAt, &n.ResolvedAt)
		negotiations = append(negotiations, n)
	}
	return negotiations, nil
}

func (r *OfferRepo) AddDocument(ctx context.Context, req *domain.OfferDocument) (*domain.OfferDocument, error) {
	d := &domain.OfferDocument{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO offer_documents (offer_id, document_type, file_name, storage_key)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, offer_id, document_type, file_name, storage_key, signed_at, created_at`,
		req.OfferID, req.DocumentType, req.FileName, req.StorageKey,
	).Scan(&d.ID, &d.OfferID, &d.DocumentType, &d.FileName, &d.StorageKey, &d.SignedAt, &d.CreatedAt)
	return d, err
}

func (r *OfferRepo) ListDocuments(ctx context.Context, offerID string) ([]domain.OfferDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, offer_id, document_type, file_name, storage_key, signed_at, created_at
		 FROM offer_documents WHERE offer_id=$1 ORDER BY created_at`, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.OfferDocument
	for rows.Next() {
		var d domain.OfferDocument
		rows.Scan(&d.ID, &d.OfferID, &d.DocumentType, &d.FileName, &d.StorageKey, &d.SignedAt, &d.CreatedAt)
		docs = append(docs, d)
	}
	return docs, nil
}
