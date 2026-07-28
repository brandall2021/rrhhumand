package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/shopspring/decimal"
)

type BankRepo struct {
	pool *pgxpool.Pool
}

func NewBankRepo(pool *pgxpool.Pool) *BankRepo {
	return &BankRepo{pool: pool}
}

func (r *BankRepo) CreateBatch(ctx context.Context, b *domain.BankBatch) error {
	q := `INSERT INTO payroll_bank_batches (id,company_id,run_id,batch_number,bank_code,bank_name,
		payment_type,total_amount,total_employees,currency,payment_date,status,file_format,
		file_name,storage_path,file_content,sent_at,processed_at,error_message,notes,generated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.RunID, b.BatchNumber, b.BankCode, b.BankName,
		b.PaymentType, b.TotalAmount, b.TotalEmployees, b.Currency, b.PaymentDate, b.Status, b.FileFormat,
		b.FileName, b.StoragePath, b.FileContent, b.SentAt, b.ProcessedAt, b.ErrorMessage, b.Notes, b.GeneratedBy)
	return repoErr("CreateBatch", err)
}

func (r *BankRepo) GetBatch(ctx context.Context, companyID, id uuid.UUID) (*domain.BankBatch, error) {
	q := `SELECT id,company_id,run_id,batch_number,bank_code,bank_name,payment_type,total_amount,
		total_employees,currency,payment_date,status,file_format,file_name,storage_path,file_content,
		sent_at,processed_at,error_message,notes,generated_by,created_at,updated_at
		FROM payroll_bank_batches WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var b domain.BankBatch
	err := row.Scan(&b.ID, &b.CompanyID, &b.RunID, &b.BatchNumber, &b.BankCode, &b.BankName, &b.PaymentType, &b.TotalAmount,
		&b.TotalEmployees, &b.Currency, &b.PaymentDate, &b.Status, &b.FileFormat, &b.FileName, &b.StoragePath, &b.FileContent,
		&b.SentAt, &b.ProcessedAt, &b.ErrorMessage, &b.Notes, &b.GeneratedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBatch", err)
	}
	return &b, nil
}

func (r *BankRepo) ListBatches(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, limit, offset int) ([]domain.BankBatch, error) {
	q := `SELECT id,company_id,run_id,batch_number,bank_code,bank_name,payment_type,total_amount,
		total_employees,currency,payment_date,status,file_format,file_name,storage_path,file_content,
		sent_at,processed_at,error_message,notes,generated_by,created_at,updated_at
		FROM payroll_bank_batches WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if runID != nil {
		q += fmt.Sprintf(" AND run_id=$%d", n)
		args = append(args, *runID)
		n++
	}
	q += " ORDER BY created_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBatches", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BankBatch, error) {
		var b domain.BankBatch
		err := row.Scan(&b.ID, &b.CompanyID, &b.RunID, &b.BatchNumber, &b.BankCode, &b.BankName, &b.PaymentType, &b.TotalAmount,
			&b.TotalEmployees, &b.Currency, &b.PaymentDate, &b.Status, &b.FileFormat, &b.FileName, &b.StoragePath, &b.FileContent,
			&b.SentAt, &b.ProcessedAt, &b.ErrorMessage, &b.Notes, &b.GeneratedBy, &b.CreatedAt, &b.UpdatedAt)
		return b, err
	})
}

func (r *BankRepo) UpdateBatchStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_bank_batches SET status=$1`
	args := []any{status}
	n := 2
	for k, v := range fields {
		q += fmt.Sprintf(",%s=$%d", k, n)
		args = append(args, v)
		n++
	}
	q += fmt.Sprintf(",updated_at=NOW() WHERE id=$%d", n)
	args = append(args, id)
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("UpdateBatchStatus", err)
}

func (r *BankRepo) CreateBatchItem(ctx context.Context, it *domain.BankBatchItem) error {
	q := `INSERT INTO payroll_bank_batch_items (id,batch_id,employee_id,run_employee_id,cuil,surname,name,
		bank_code,bank_name,branch_code,account_type,account_number,cbu,alias,amount,currency,concept,status,
		error_message,payment_date,transaction_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`
	_, err := r.pool.Exec(ctx, q, it.ID, it.BatchID, it.EmployeeID, it.RunEmployeeID, it.CUIL, it.Surname, it.Name,
		it.BankCode, it.BankName, it.BranchCode, it.AccountType, it.AccountNumber, it.CBU, it.Alias,
		it.Amount, it.Currency, it.Concept, it.Status, it.ErrorMessage, it.PaymentDate, it.TransactionID)
	return repoErr("CreateBatchItem", err)
}

func (r *BankRepo) BulkCreateBatchItems(ctx context.Context, items []domain.BankBatchItem) error {
	if len(items) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_bank_batch_items (id,batch_id,employee_id,run_employee_id,cuil,surname,name,
		bank_code,bank_name,branch_code,account_type,account_number,cbu,alias,amount,currency,concept,status,
		error_message,payment_date,transaction_id) VALUES `
	args := []any{}
	n := 1
	for _, it := range items {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14, n+15, n+16, n+17, n+18, n+19, n+20)
		args = append(args, it.ID, it.BatchID, it.EmployeeID, it.RunEmployeeID, it.CUIL, it.Surname, it.Name,
			it.BankCode, it.BankName, it.BranchCode, it.AccountType, it.AccountNumber, it.CBU, it.Alias,
			it.Amount, it.Currency, it.Concept, it.Status, it.ErrorMessage, it.PaymentDate, it.TransactionID)
		n += 21
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateBatchItems", err)
}

func (r *BankRepo) ListBatchItems(ctx context.Context, batchID uuid.UUID) ([]domain.BankBatchItem, error) {
	q := `SELECT id,batch_id,employee_id,run_employee_id,cuil,surname,name,
		bank_code,bank_name,branch_code,account_type,account_number,cbu,alias,amount,currency,concept,status,
		error_message,payment_date,transaction_id,created_at,updated_at
		FROM payroll_bank_batch_items WHERE batch_id=$1 ORDER BY surname, name`
	rows, err := r.pool.Query(ctx, q, batchID)
	if err != nil {
		return nil, repoErr("ListBatchItems", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BankBatchItem, error) {
		var it domain.BankBatchItem
		err := row.Scan(&it.ID, &it.BatchID, &it.EmployeeID, &it.RunEmployeeID, &it.CUIL, &it.Surname, &it.Name,
			&it.BankCode, &it.BankName, &it.BranchCode, &it.AccountType, &it.AccountNumber, &it.CBU, &it.Alias,
			&it.Amount, &it.Currency, &it.Concept, &it.Status, &it.ErrorMessage, &it.PaymentDate, &it.TransactionID,
			&it.CreatedAt, &it.UpdatedAt)
		return it, err
	})
}

func (r *BankRepo) ListBatchItemsByEmployee(ctx context.Context, companyID, employeeID uuid.UUID, limit, offset int) ([]domain.BankBatchItem, error) {
	q := `SELECT i.id,i.batch_id,i.employee_id,i.run_employee_id,i.cuil,i.surname,i.name,
		i.bank_code,i.bank_name,i.branch_code,i.account_type,i.account_number,i.cbu,i.alias,
		i.amount,i.currency,i.concept,i.status,i.error_message,i.payment_date,i.transaction_id,i.created_at,i.updated_at
		FROM payroll_bank_batch_items i
		JOIN payroll_bank_batches b ON b.id=i.batch_id
		WHERE b.company_id=$1 AND i.employee_id=$2 ORDER BY i.created_at DESC`
	args := []any{companyID, employeeID}
	n := 3
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBatchItemsByEmployee", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BankBatchItem, error) {
		var it domain.BankBatchItem
		err := row.Scan(&it.ID, &it.BatchID, &it.EmployeeID, &it.RunEmployeeID, &it.CUIL, &it.Surname, &it.Name,
			&it.BankCode, &it.BankName, &it.BranchCode, &it.AccountType, &it.AccountNumber, &it.CBU, &it.Alias,
			&it.Amount, &it.Currency, &it.Concept, &it.Status, &it.ErrorMessage, &it.PaymentDate, &it.TransactionID,
			&it.CreatedAt, &it.UpdatedAt)
		return it, err
	})
}

func (r *BankRepo) UpdateBatchItemStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_bank_batch_items SET status=$1`
	args := []any{status}
	n := 2
	for k, v := range fields {
		q += fmt.Sprintf(",%s=$%d", k, n)
		args = append(args, v)
		n++
	}
	q += fmt.Sprintf(",updated_at=NOW() WHERE id=$%d", n)
	args = append(args, id)
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("UpdateBatchItemStatus", err)
}

type BatchSummary struct {
	TotalItems    int             `json:"total_items"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	ProcessedItems int            `json:"processed_items"`
	FailedItems   int             `json:"failed_items"`
}

func (r *BankRepo) GetBatchSummary(ctx context.Context, batchID uuid.UUID) (*BatchSummary, error) {
	q := `SELECT COUNT(*),COALESCE(SUM(amount),0),
		COUNT(*) FILTER (WHERE status='PROCESSED' OR status='SENT'),
		COUNT(*) FILTER (WHERE status='FAILED')
		FROM payroll_bank_batch_items WHERE batch_id=$1`
	row := r.pool.QueryRow(ctx, q, batchID)
	var s BatchSummary
	err := row.Scan(&s.TotalItems, &s.TotalAmount, &s.ProcessedItems, &s.FailedItems)
	if err != nil {
		return nil, repoErr("GetBatchSummary", err)
	}
	return &s, nil
}
