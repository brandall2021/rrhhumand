package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type AssignmentRepo struct {
	pool *pgxpool.Pool
}

func NewAssignmentRepo(pool *pgxpool.Pool) *AssignmentRepo {
	return &AssignmentRepo{pool: pool}
}

func scanEmployeeBenefit(row pgx.Row) (domain.EmployeeBenefit, error) {
	var eb domain.EmployeeBenefit
	err := row.Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitID, &eb.PlanID, &eb.ProviderID,
		&eb.Status, &eb.EnrollmentDate, &eb.StartDate, &eb.EndDate, &eb.CancellationDate, &eb.CancellationReason,
		&eb.AutoRenew, &eb.RenewalDate, &eb.CoverageLevel, &eb.Dependents, &eb.EmergencyContact,
		&eb.BeneficiaryInfo, &eb.EmployeeCost, &eb.EmployerCost, &eb.Currency,
		&eb.PayrollDeductionEnabled, &eb.PayrollDeductionAmount, &eb.ExternalMemberID,
		&eb.ExternalPolicyNumber, &eb.ExternalGroupNumber, &eb.Documents, &eb.Notes,
		&eb.Source, &eb.EnrolledBy, &eb.EnrolledAt, &eb.CreatedAt, &eb.UpdatedAt)
	return eb, err
}

func (r *AssignmentRepo) Create(ctx context.Context, eb *domain.EmployeeBenefit) error {
	q := `INSERT INTO employee_benefits (id,company_id,employee_id,benefit_id,plan_id,provider_id,status,
		enrollment_date,start_date,end_date,cancellation_date,cancellation_reason,auto_renew,renewal_date,
		coverage_level,dependents,emergency_contact,beneficiary_info,employee_cost,employer_cost,currency,
		payroll_deduction_enabled,payroll_deduction_amount,external_member_id,external_policy_number,
		external_group_number,documents,notes,source,enrolled_by,enrolled_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)`
	_, err := r.pool.Exec(ctx, q, eb.ID, eb.CompanyID, eb.EmployeeID, eb.BenefitID, eb.PlanID, eb.ProviderID, eb.Status,
		eb.EnrollmentDate, eb.StartDate, eb.EndDate, eb.CancellationDate, eb.CancellationReason, eb.AutoRenew, eb.RenewalDate,
		eb.CoverageLevel, eb.Dependents, eb.EmergencyContact, eb.BeneficiaryInfo, eb.EmployeeCost, eb.EmployerCost, eb.Currency,
		eb.PayrollDeductionEnabled, eb.PayrollDeductionAmount, eb.ExternalMemberID, eb.ExternalPolicyNumber,
		eb.ExternalGroupNumber, eb.Documents, eb.Notes, eb.Source, eb.EnrolledBy, eb.EnrolledAt)
	return repoErr("Create", err)
}

func (r *AssignmentRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeBenefit, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,plan_id,provider_id,status,enrollment_date,start_date,
		end_date,cancellation_date,cancellation_reason,auto_renew,renewal_date,coverage_level,dependents,
		emergency_contact,beneficiary_info,employee_cost,employer_cost,currency,payroll_deduction_enabled,
		payroll_deduction_amount,external_member_id,external_policy_number,external_group_number,documents,notes,
		source,enrolled_by,enrolled_at,created_at,updated_at
		FROM employee_benefits WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	eb, err := scanEmployeeBenefit(row)
	if err != nil {
		return nil, repoErr("Get", err)
	}
	return &eb, nil
}

func (r *AssignmentRepo) List(ctx context.Context, companyID, employeeID, benefitID *uuid.UUID, status *string) ([]domain.EmployeeBenefit, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,plan_id,provider_id,status,enrollment_date,start_date,
		end_date,cancellation_date,cancellation_reason,auto_renew,renewal_date,coverage_level,dependents,
		emergency_contact,beneficiary_info,employee_cost,employer_cost,currency,payroll_deduction_enabled,
		payroll_deduction_amount,external_member_id,external_policy_number,external_group_number,documents,notes,
		source,enrolled_by,enrolled_at,created_at,updated_at
		FROM employee_benefits WHERE 1=1`
	args := []any{}
	n := 1
	if companyID != nil {
		q += fmt.Sprintf(" AND company_id=$%d", n)
		args = append(args, *companyID)
		n++
	}
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if benefitID != nil {
		q += fmt.Sprintf(" AND benefit_id=$%d", n)
		args = append(args, *benefitID)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("List", err)
	}
	defer rows.Close()
	var result []domain.EmployeeBenefit
	for rows.Next() {
		var eb domain.EmployeeBenefit
		err := rows.Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitID, &eb.PlanID, &eb.ProviderID,
			&eb.Status, &eb.EnrollmentDate, &eb.StartDate, &eb.EndDate, &eb.CancellationDate, &eb.CancellationReason,
			&eb.AutoRenew, &eb.RenewalDate, &eb.CoverageLevel, &eb.Dependents, &eb.EmergencyContact,
			&eb.BeneficiaryInfo, &eb.EmployeeCost, &eb.EmployerCost, &eb.Currency,
			&eb.PayrollDeductionEnabled, &eb.PayrollDeductionAmount, &eb.ExternalMemberID,
			&eb.ExternalPolicyNumber, &eb.ExternalGroupNumber, &eb.Documents, &eb.Notes,
			&eb.Source, &eb.EnrolledBy, &eb.EnrolledAt, &eb.CreatedAt, &eb.UpdatedAt)
		if err != nil {
			return nil, repoErr("List", err)
		}
		result = append(result, eb)
	}
	return result, nil
}

func (r *AssignmentRepo) Update(ctx context.Context, eb *domain.EmployeeBenefit) error {
	q := `UPDATE employee_benefits SET plan_id=$1,provider_id=$2,status=$3,enrollment_date=$4,start_date=$5,
		end_date=$6,cancellation_date=$7,cancellation_reason=$8,auto_renew=$9,renewal_date=$10,
		coverage_level=$11,dependents=$12,emergency_contact=$13,beneficiary_info=$14,employee_cost=$15,
		employer_cost=$16,currency=$17,payroll_deduction_enabled=$18,payroll_deduction_amount=$19,
		external_member_id=$20,external_policy_number=$21,external_group_number=$22,documents=$23,notes=$24,
		updated_at=NOW() WHERE id=$25 AND company_id=$26`
	_, err := r.pool.Exec(ctx, q, eb.PlanID, eb.ProviderID, eb.Status, eb.EnrollmentDate, eb.StartDate,
		eb.EndDate, eb.CancellationDate, eb.CancellationReason, eb.AutoRenew, eb.RenewalDate,
		eb.CoverageLevel, eb.Dependents, eb.EmergencyContact, eb.BeneficiaryInfo, eb.EmployeeCost,
		eb.EmployerCost, eb.Currency, eb.PayrollDeductionEnabled, eb.PayrollDeductionAmount,
		eb.ExternalMemberID, eb.ExternalPolicyNumber, eb.ExternalGroupNumber, eb.Documents, eb.Notes,
		eb.ID, eb.CompanyID)
	return repoErr("Update", err)
}

func (r *AssignmentRepo) Delete(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM employee_benefits WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("Delete", err)
}

func (r *AssignmentRepo) CreateHistory(ctx context.Context, h *domain.EmployeeBenefitHistory) error {
	q := `INSERT INTO employee_benefit_history (id,employee_benefit_id,employee_id,benefit_id,action,
		previous_value,new_value,change_reason,changed_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, h.ID, h.EmployeeBenefitID, h.EmployeeID, h.BenefitID, h.Action,
		h.PreviousValue, h.NewValue, h.ChangeReason, h.ChangedBy)
	return repoErr("CreateHistory", err)
}

func (r *AssignmentRepo) ListHistory(ctx context.Context, employeeBenefitID, employeeID *uuid.UUID) ([]domain.EmployeeBenefitHistory, error) {
	q := `SELECT id,employee_benefit_id,employee_id,benefit_id,action,previous_value,new_value,change_reason,changed_by,changed_at
		FROM employee_benefit_history WHERE 1=1`
	args := []any{}
	n := 1
	if employeeBenefitID != nil {
		q += fmt.Sprintf(" AND employee_benefit_id=$%d", n)
		args = append(args, *employeeBenefitID)
		n++
	}
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	q += " ORDER BY changed_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListHistory", err)
	}
	defer rows.Close()
	var result []domain.EmployeeBenefitHistory
	for rows.Next() {
		var h domain.EmployeeBenefitHistory
		err := rows.Scan(&h.ID, &h.EmployeeBenefitID, &h.EmployeeID, &h.BenefitID, &h.Action,
			&h.PreviousValue, &h.NewValue, &h.ChangeReason, &h.ChangedBy, &h.ChangedAt)
		if err != nil {
			return nil, repoErr("ListHistory", err)
		}
		result = append(result, h)
	}
	return result, nil
}

func scanBenefitRequest(row pgx.Row) (domain.BenefitRequest, error) {
	var br domain.BenefitRequest
	err := row.Scan(&br.ID, &br.CompanyID, &br.EmployeeID, &br.BenefitID, &br.PlanID, &br.EmployeeBenefitID,
		&br.RequestType, &br.Status, &br.RequestData, &br.Justification, &br.CoverageLevel, &br.Dependents,
		&br.EffectiveDate, &br.Notes, &br.SubmittedBy, &br.SubmittedAt, &br.ResolvedBy, &br.ResolvedAt,
		&br.ResolutionNotes, &br.CreatedAt, &br.UpdatedAt)
	return br, err
}

func (r *AssignmentRepo) CreateRequest(ctx context.Context, br *domain.BenefitRequest) error {
	q := `INSERT INTO benefit_requests (id,company_id,employee_id,benefit_id,plan_id,employee_benefit_id,
		request_type,status,request_data,justification,coverage_level,dependents,effective_date,notes,
		submitted_by,submitted_at,resolved_by,resolved_at,resolution_notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`
	_, err := r.pool.Exec(ctx, q, br.ID, br.CompanyID, br.EmployeeID, br.BenefitID, br.PlanID, br.EmployeeBenefitID,
		br.RequestType, br.Status, br.RequestData, br.Justification, br.CoverageLevel, br.Dependents,
		br.EffectiveDate, br.Notes, br.SubmittedBy, br.SubmittedAt, br.ResolvedBy, br.ResolvedAt,
		br.ResolutionNotes)
	return repoErr("CreateRequest", err)
}

func (r *AssignmentRepo) GetRequest(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitRequest, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,plan_id,employee_benefit_id,request_type,status,
		request_data,justification,coverage_level,dependents,effective_date,notes,submitted_by,submitted_at,
		resolved_by,resolved_at,resolution_notes,created_at,updated_at
		FROM benefit_requests WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	br, err := scanBenefitRequest(row)
	if err != nil {
		return nil, repoErr("GetRequest", err)
	}
	return &br, nil
}

func (r *AssignmentRepo) ListRequests(ctx context.Context, companyID, employeeID, benefitID *uuid.UUID, status *string) ([]domain.BenefitRequest, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,plan_id,employee_benefit_id,request_type,status,
		request_data,justification,coverage_level,dependents,effective_date,notes,submitted_by,submitted_at,
		resolved_by,resolved_at,resolution_notes,created_at,updated_at
		FROM benefit_requests WHERE 1=1`
	args := []any{}
	n := 1
	if companyID != nil {
		q += fmt.Sprintf(" AND company_id=$%d", n)
		args = append(args, *companyID)
		n++
	}
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if benefitID != nil {
		q += fmt.Sprintf(" AND benefit_id=$%d", n)
		args = append(args, *benefitID)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListRequests", err)
	}
	defer rows.Close()
	var result []domain.BenefitRequest
	for rows.Next() {
		var br domain.BenefitRequest
		err := rows.Scan(&br.ID, &br.CompanyID, &br.EmployeeID, &br.BenefitID, &br.PlanID, &br.EmployeeBenefitID,
			&br.RequestType, &br.Status, &br.RequestData, &br.Justification, &br.CoverageLevel, &br.Dependents,
			&br.EffectiveDate, &br.Notes, &br.SubmittedBy, &br.SubmittedAt, &br.ResolvedBy, &br.ResolvedAt,
			&br.ResolutionNotes, &br.CreatedAt, &br.UpdatedAt)
		if err != nil {
			return nil, repoErr("ListRequests", err)
		}
		result = append(result, br)
	}
	return result, nil
}

func (r *AssignmentRepo) UpdateRequest(ctx context.Context, br *domain.BenefitRequest) error {
	q := `UPDATE benefit_requests SET plan_id=$1,employee_benefit_id=$2,request_type=$3,status=$4,
		request_data=$5,justification=$6,coverage_level=$7,dependents=$8,effective_date=$9,notes=$10,
		submitted_by=$11,submitted_at=$12,resolved_by=$13,resolved_at=$14,resolution_notes=$15,
		updated_at=NOW() WHERE id=$16 AND company_id=$17`
	_, err := r.pool.Exec(ctx, q, br.PlanID, br.EmployeeBenefitID, br.RequestType, br.Status,
		br.RequestData, br.Justification, br.CoverageLevel, br.Dependents, br.EffectiveDate, br.Notes,
		br.SubmittedBy, br.SubmittedAt, br.ResolvedBy, br.ResolvedAt, br.ResolutionNotes,
		br.ID, br.CompanyID)
	return repoErr("UpdateRequest", err)
}

func (r *AssignmentRepo) CreateReview(ctx context.Context, rev *domain.BenefitRequestReview) error {
	q := `INSERT INTO benefit_request_reviews (id,request_id,step_id,reviewer_id,review_type,comment)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, rev.ID, rev.RequestID, rev.StepID, rev.ReviewerID, rev.ReviewType, rev.Comment)
	return repoErr("CreateReview", err)
}

func (r *AssignmentRepo) ListReviews(ctx context.Context, requestID uuid.UUID) ([]domain.BenefitRequestReview, error) {
	q := `SELECT id,request_id,step_id,reviewer_id,review_type,comment,reviewed_at
		FROM benefit_request_reviews WHERE request_id=$1 ORDER BY reviewed_at`
	rows, err := r.pool.Query(ctx, q, requestID)
	if err != nil {
		return nil, repoErr("ListReviews", err)
	}
	defer rows.Close()
	var result []domain.BenefitRequestReview
	for rows.Next() {
		var rev domain.BenefitRequestReview
		err := rows.Scan(&rev.ID, &rev.RequestID, &rev.StepID, &rev.ReviewerID, &rev.ReviewType, &rev.Comment, &rev.ReviewedAt)
		if err != nil {
			return nil, repoErr("ListReviews", err)
		}
		result = append(result, rev)
	}
	return result, nil
}
