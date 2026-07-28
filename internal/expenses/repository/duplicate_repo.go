package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type DuplicateRepo struct {
	pool *pgxpool.Pool
}

func NewDuplicateRepo(pool *pgxpool.Pool) *DuplicateRepo {
	return &DuplicateRepo{pool: pool}
}

func (r *DuplicateRepo) Create(ctx context.Context, d *domain.ExpenseDuplicateCheck) error {
	q := `INSERT INTO expense_duplicate_checks (id,expense_id,matched_expense_id,similarity_score,
		match_reason,status,reviewed_by,reviewed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, d.ID, d.ExpenseID, d.MatchedExpenseID, d.SimilarityScore,
		d.MatchReason, d.Status, d.ReviewedBy, d.ReviewedAt)
	return repoErr("DuplicateRepo.Create", err)
}

func (r *DuplicateRepo) ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseDuplicateCheck, error) {
	q := `SELECT id,expense_id,matched_expense_id,similarity_score,match_reason,status,reviewed_by,reviewed_at,created_at
		FROM expense_duplicate_checks WHERE expense_id=$1 ORDER BY similarity_score DESC`
	rows, err := r.pool.Query(ctx, q, expenseID)
	if err != nil {
		return nil, repoErr("DuplicateRepo.ListByExpense", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseDuplicateCheck, error) {
		var d domain.ExpenseDuplicateCheck
		err := row.Scan(&d.ID, &d.ExpenseID, &d.MatchedExpenseID, &d.SimilarityScore, &d.MatchReason,
			&d.Status, &d.ReviewedBy, &d.ReviewedAt, &d.CreatedAt)
		return d, err
	})
}

func (r *DuplicateRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reviewedBy *uuid.UUID) error {
	q := `UPDATE expense_duplicate_checks SET status=$1,reviewed_by=$2,reviewed_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, q, status, reviewedBy, id)
	return repoErr("DuplicateRepo.UpdateStatus", err)
}
