package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type TalentPoolRepo struct {
	pool *pgxpool.Pool
}

func NewTalentPoolRepo(pool *pgxpool.Pool) *TalentPoolRepo {
	return &TalentPoolRepo{pool: pool}
}

func (r *TalentPoolRepo) CreatePool(ctx context.Context, companyID string, req *domain.TalentPool) (*domain.TalentPool, error) {
	tp := &domain.TalentPool{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO talent_pools (company_id, name, description, criteria, is_auto)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, company_id, name, description, criteria, is_auto, created_at`,
		companyID, req.Name, req.Description, req.Criteria, req.IsAuto,
	).Scan(&tp.ID, &tp.CompanyID, &tp.Name, &tp.Description, &tp.Criteria, &tp.IsAuto, &tp.CreatedAt)
	return tp, err
}

func (r *TalentPoolRepo) GetPool(ctx context.Context, companyID, id string) (*domain.TalentPool, error) {
	tp := &domain.TalentPool{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, criteria, is_auto, created_at
		 FROM talent_pools WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&tp.ID, &tp.CompanyID, &tp.Name, &tp.Description, &tp.Criteria, &tp.IsAuto, &tp.CreatedAt)
	return tp, err
}

func (r *TalentPoolRepo) ListPools(ctx context.Context, companyID string) ([]domain.TalentPool, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, criteria, is_auto, created_at
		 FROM talent_pools WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []domain.TalentPool
	for rows.Next() {
		var tp domain.TalentPool
		rows.Scan(&tp.ID, &tp.CompanyID, &tp.Name, &tp.Description, &tp.Criteria, &tp.IsAuto, &tp.CreatedAt)
		pools = append(pools, tp)
	}
	return pools, nil
}

func (r *TalentPoolRepo) UpdatePool(ctx context.Context, companyID, id string, req *domain.TalentPool) (*domain.TalentPool, error) {
	tp := &domain.TalentPool{}
	err := r.pool.QueryRow(ctx,
		`UPDATE talent_pools SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 criteria=COALESCE($5,criteria), is_auto=COALESCE($6,is_auto)
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, criteria, is_auto, created_at`,
		companyID, id, req.Name, req.Description, req.Criteria, req.IsAuto,
	).Scan(&tp.ID, &tp.CompanyID, &tp.Name, &tp.Description, &tp.Criteria, &tp.IsAuto, &tp.CreatedAt)
	return tp, err
}

func (r *TalentPoolRepo) DeletePool(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM talent_pools WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *TalentPoolRepo) AddCandidate(ctx context.Context, req *domain.TalentPoolCandidate) (*domain.TalentPoolCandidate, error) {
	tpc := &domain.TalentPoolCandidate{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO talent_pool_candidates (pool_id, candidate_id, added_by, added_reason)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (pool_id, candidate_id) DO NOTHING
		 RETURNING id, pool_id, candidate_id, added_by, added_reason, created_at`,
		req.PoolID, req.CandidateID, req.AddedBy, req.AddedReason,
	).Scan(&tpc.ID, &tpc.PoolID, &tpc.CandidateID, &tpc.AddedBy, &tpc.AddedReason, &tpc.CreatedAt)
	return tpc, err
}

func (r *TalentPoolRepo) RemoveCandidate(ctx context.Context, poolID, candidateID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM talent_pool_candidates WHERE pool_id=$1 AND candidate_id=$2`,
		poolID, candidateID)
	return err
}

func (r *TalentPoolRepo) ListPoolCandidates(ctx context.Context, poolID string) ([]domain.TalentPoolCandidate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, pool_id, candidate_id, added_by, added_reason, created_at
		 FROM talent_pool_candidates WHERE pool_id=$1 ORDER BY created_at DESC`, poolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.TalentPoolCandidate
	for rows.Next() {
		var tpc domain.TalentPoolCandidate
		rows.Scan(&tpc.ID, &tpc.PoolID, &tpc.CandidateID, &tpc.AddedBy, &tpc.AddedReason, &tpc.CreatedAt)
		candidates = append(candidates, tpc)
	}
	return candidates, nil
}
