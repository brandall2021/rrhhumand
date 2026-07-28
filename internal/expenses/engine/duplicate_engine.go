package engine

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/rrhhumand/api/internal/expenses/repository"
)

type DuplicateEngine struct {
	duplicateRepo *repository.DuplicateRepo
	pool          *pgxpool.Pool
}

func NewDuplicateEngine(dr *repository.DuplicateRepo, pool *pgxpool.Pool) *DuplicateEngine {
	return &DuplicateEngine{duplicateRepo: dr, pool: pool}
}

type matchCandidate struct {
	ID     uuid.UUID
	Reason string
	Score  *float64
}

func (e *DuplicateEngine) CheckDuplicate(ctx context.Context, expense *domain.Expense) ([]domain.ExpenseDuplicateCheck, error) {
	now := time.Now()
	var results []domain.ExpenseDuplicateCheck

	candidates, err := e.findCandidates(ctx, expense)
	if err != nil {
		return nil, engErr("CheckDuplicate.findCandidates", err)
	}

	for _, m := range candidates {
		check := domain.ExpenseDuplicateCheck{
			ID:                uuid.New(),
			ExpenseID:         expense.ID,
			DuplicateExpenseID: &m.ID,
			MatchReason:       m.Reason,
			MatchScore:        m.Score,
			Status:            "PENDING",
			CreatedAt:         now,
		}

		if err := e.duplicateRepo.Create(ctx, &check); err != nil {
			return nil, engErr("CheckDuplicate.create", err)
		}

		results = append(results, check)
	}

	return results, nil
}

func (e *DuplicateEngine) findCandidates(ctx context.Context, expense *domain.Expense) ([]matchCandidate, error) {
	var candidates []matchCandidate

	if expense.MerchantTaxID != nil && *expense.MerchantTaxID != "" {
		rows, err := e.pool.Query(ctx, `
			SELECT id FROM expenses WHERE company_id=$1 AND merchant_tax_id=$2 AND base_amount=$3
			AND id!=$4 AND status IN ('SUBMITTED','APPROVED','PAID')`,
			expense.CompanyID, *expense.MerchantTaxID, expense.BaseAmount, expense.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			score := 0.95
			candidates = append(candidates, matchCandidate{ID: id, Reason: "MERCHANT_TAX_ID_AMOUNT", Score: &score})
		}
		rows.Close()
	}

	if expense.ReceiptNumber != nil && *expense.ReceiptNumber != "" {
		rows, err := e.pool.Query(ctx, `
			SELECT id FROM expenses WHERE company_id=$1 AND receipt_number=$2
			AND id!=$3 AND status IN ('SUBMITTED','APPROVED','PAID')`,
			expense.CompanyID, *expense.ReceiptNumber, expense.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			score := 0.90
			candidates = append(candidates, matchCandidate{ID: id, Reason: "RECEIPT_NUMBER", Score: &score})
		}
		rows.Close()
	}

	return candidates, nil
}

func (e *DuplicateEngine) MarkDuplicate(ctx context.Context, checkID uuid.UUID, resolvedBy uuid.UUID, status string) error {
	return e.duplicateRepo.UpdateStatus(ctx, checkID, status, &resolvedBy)
}
