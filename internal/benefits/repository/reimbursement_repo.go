package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type ReimbursementRepo struct {
	pool *pgxpool.Pool
}

func NewReimbursementRepo(pool *pgxpool.Pool) *ReimbursementRepo {
	return &ReimbursementRepo{pool: pool}
}

func scanReimbursement(row pgx.CollectableRow) (domain.BenefitReimbursement, error) {
	var r domain.BenefitReimbursement
	err := row.Scan(&r.ID, &r.CompanyID, &r.EmployeeID, &r.BenefitID, &r.FlexiblePlanID, &r.WalletID,
		&r.RequestID, &r.Category, &r.Description, &r.Amount, &r.ApprovedAmount, &r.Currency,
		&r.ReceiptDate, &r.ExpenseDate, &r.Status, &r.RejectionReason, &r.PaymentMethod, &r.PaidAt,
		&r.PaymentReference, &r.SubmittedBy, &r.ReviewedBy, &r.ReviewedAt, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (r *ReimbursementRepo) Create(ctx context.Context, reimb *domain.BenefitReimbursement) error {
	q := `INSERT INTO benefit_reimbursements (id,company_id,employee_id,benefit_id,flexible_plan_id,wallet_id,
		request_id,category,description,amount,approved_amount,currency,receipt_date,expense_date,status,
		rejection_reason,payment_method,paid_at,payment_reference,submitted_by,reviewed_by,reviewed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`
	_, err := r.pool.Exec(ctx, q, reimb.ID, reimb.CompanyID, reimb.EmployeeID, reimb.BenefitID, reimb.FlexiblePlanID, reimb.WalletID,
		reimb.RequestID, reimb.Category, reimb.Description, reimb.Amount, reimb.ApprovedAmount, reimb.Currency,
		reimb.ReceiptDate, reimb.ExpenseDate, reimb.Status, reimb.RejectionReason, reimb.PaymentMethod, reimb.PaidAt,
		reimb.PaymentReference, reimb.SubmittedBy, reimb.ReviewedBy, reimb.ReviewedAt)
	return repoErr("Create", err)
}

func (r *ReimbursementRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitReimbursement, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,flexible_plan_id,wallet_id,request_id,category,description,
		amount,approved_amount,currency,receipt_date,expense_date,status,rejection_reason,payment_method,paid_at,
		payment_reference,submitted_by,reviewed_by,reviewed_at,created_at,updated_at
		FROM benefit_reimbursements WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	reimb, err := scanReimbursement(row)
	if err != nil {
		return nil, repoErr("Get", err)
	}
	return &reimb, nil
}

func (r *ReimbursementRepo) List(ctx context.Context, companyID, employeeID, benefitID *uuid.UUID, status *string) ([]domain.BenefitReimbursement, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,flexible_plan_id,wallet_id,request_id,category,description,
		amount,approved_amount,currency,receipt_date,expense_date,status,rejection_reason,payment_method,paid_at,
		payment_reference,submitted_by,reviewed_by,reviewed_at,created_at,updated_at
		FROM benefit_reimbursements WHERE 1=1`
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
	return pgx.CollectRows(rows, scanReimbursement)
}

func (r *ReimbursementRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reviewedBy *uuid.UUID) error {
	q := `UPDATE benefit_reimbursements SET status=$1,reviewed_by=$2,reviewed_at=NOW(),updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, status, reviewedBy, id)
	return repoErr("UpdateStatus", err)
}

func (r *ReimbursementRepo) CreateDocument(ctx context.Context, d *domain.BenefitReimbursementDocument) error {
	q := `INSERT INTO benefit_reimbursement_documents (id,reimbursement_id,file_name,file_type,file_size,storage_path,ocr_text,is_verified,uploaded_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, d.ID, d.ReimbursementID, d.FileName, d.FileType, d.FileSize, d.StoragePath, d.OCRText, d.IsVerified, d.UploadedBy)
	return repoErr("CreateDocument", err)
}

func (r *ReimbursementRepo) ListDocuments(ctx context.Context, reimbursementID uuid.UUID) ([]domain.BenefitReimbursementDocument, error) {
	q := `SELECT id,reimbursement_id,file_name,file_type,file_size,storage_path,ocr_text,is_verified,uploaded_by,uploaded_at
		FROM benefit_reimbursement_documents WHERE reimbursement_id=$1 ORDER BY uploaded_at`
	rows, err := r.pool.Query(ctx, q, reimbursementID)
	if err != nil {
		return nil, repoErr("ListDocuments", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitReimbursementDocument, error) {
		var d domain.BenefitReimbursementDocument
		err := row.Scan(&d.ID, &d.ReimbursementID, &d.FileName, &d.FileType, &d.FileSize, &d.StoragePath, &d.OCRText, &d.IsVerified, &d.UploadedBy, &d.UploadedAt)
		return d, err
	})
}
