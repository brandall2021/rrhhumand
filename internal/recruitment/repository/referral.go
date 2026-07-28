package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type ReferralRepo struct {
	pool *pgxpool.Pool
}

func NewReferralRepo(pool *pgxpool.Pool) *ReferralRepo {
	return &ReferralRepo{pool: pool}
}

func (r *ReferralRepo) CreateReward(ctx context.Context, companyID string, req *domain.ReferralReward) (*domain.ReferralReward, error) {
	rr := &domain.ReferralReward{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO referral_rewards (company_id, referral_id, reward_type, amount, currency, status)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, referral_id, reward_type, amount, currency, status, paid_at, created_at`,
		companyID, req.ReferralID, req.RewardType, req.Amount, req.Currency, req.Status,
	).Scan(&rr.ID, &rr.CompanyID, &rr.ReferralID, &rr.RewardType, &rr.Amount, &rr.Currency, &rr.Status, &rr.PaidAt, &rr.CreatedAt)
	return rr, err
}

func (r *ReferralRepo) GetReward(ctx context.Context, companyID, id string) (*domain.ReferralReward, error) {
	rr := &domain.ReferralReward{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, referral_id, reward_type, amount, currency, status, paid_at, created_at
		 FROM referral_rewards WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&rr.ID, &rr.CompanyID, &rr.ReferralID, &rr.RewardType, &rr.Amount, &rr.Currency, &rr.Status, &rr.PaidAt, &rr.CreatedAt)
	return rr, err
}

func (r *ReferralRepo) ListRewards(ctx context.Context, companyID string, referralID string) ([]domain.ReferralReward, error) {
	query := `SELECT id, company_id, referral_id, reward_type, amount, currency, status, paid_at, created_at
		 FROM referral_rewards WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if referralID != "" {
		query += fmt.Sprintf(" AND referral_id=$%d", argIdx)
		args = append(args, referralID)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []domain.ReferralReward
	for rows.Next() {
		var rr domain.ReferralReward
		rows.Scan(&rr.ID, &rr.CompanyID, &rr.ReferralID, &rr.RewardType, &rr.Amount, &rr.Currency, &rr.Status, &rr.PaidAt, &rr.CreatedAt)
		rewards = append(rewards, rr)
	}
	return rewards, nil
}

func (r *ReferralRepo) UpdateReward(ctx context.Context, id, status string) (*domain.ReferralReward, error) {
	rr := &domain.ReferralReward{}
	err := r.pool.QueryRow(ctx,
		`UPDATE referral_rewards SET status=$2, paid_at=CASE WHEN $2='PAID' THEN NOW() ELSE paid_at END WHERE id=$1
		 RETURNING id, company_id, referral_id, reward_type, amount, currency, status, paid_at, created_at`,
		id, status,
	).Scan(&rr.ID, &rr.CompanyID, &rr.ReferralID, &rr.RewardType, &rr.Amount, &rr.Currency, &rr.Status, &rr.PaidAt, &rr.CreatedAt)
	return rr, err
}
