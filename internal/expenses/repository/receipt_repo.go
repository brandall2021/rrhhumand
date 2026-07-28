package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ReceiptRepo struct {
	pool *pgxpool.Pool
}

func NewReceiptRepo(pool *pgxpool.Pool) *ReceiptRepo {
	return &ReceiptRepo{pool: pool}
}

func (r *ReceiptRepo) Create(ctx context.Context, rec *domain.ExpenseReceipt) error {
	q := `INSERT INTO expense_receipts (id,expense_id,file_name,file_path,file_size,mime_type,uploaded_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, rec.ID, rec.ExpenseID, rec.FileName, rec.FilePath, rec.FileSize, rec.MimeType, rec.UploadedBy)
	return repoErr("ReceiptRepo.Create", err)
}

func (r *ReceiptRepo) Get(ctx context.Context, id uuid.UUID) (*domain.ExpenseReceipt, error) {
	q := `SELECT id,expense_id,file_name,file_path,file_size,mime_type,uploaded_by,created_at
		FROM expense_receipts WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var rec domain.ExpenseReceipt
	err := row.Scan(&rec.ID, &rec.ExpenseID, &rec.FileName, &rec.FilePath, &rec.FileSize, &rec.MimeType, &rec.UploadedBy, &rec.CreatedAt)
	if err != nil {
		return nil, repoErr("ReceiptRepo.Get", err)
	}
	return &rec, nil
}

func (r *ReceiptRepo) ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error) {
	q := `SELECT id,expense_id,file_name,file_path,file_size,mime_type,uploaded_by,created_at
		FROM expense_receipts WHERE expense_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, expenseID)
	if err != nil {
		return nil, repoErr("ReceiptRepo.ListByExpense", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseReceipt, error) {
		var rec domain.ExpenseReceipt
		err := row.Scan(&rec.ID, &rec.ExpenseID, &rec.FileName, &rec.FilePath, &rec.FileSize, &rec.MimeType, &rec.UploadedBy, &rec.CreatedAt)
		return rec, err
	})
}

func (r *ReceiptRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_receipts WHERE id=$1`, id)
	return repoErr("ReceiptRepo.Delete", err)
}
