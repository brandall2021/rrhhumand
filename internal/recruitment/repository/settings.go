package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type SettingsRepo struct {
	pool *pgxpool.Pool
}

func NewSettingsRepo(pool *pgxpool.Pool) *SettingsRepo {
	return &SettingsRepo{pool: pool}
}

func (r *SettingsRepo) CreateSource(ctx context.Context, companyID string, req *domain.RecruitmentSource) (*domain.RecruitmentSource, error) {
	s := &domain.RecruitmentSource{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_sources (company_id, name, type, config, active)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, company_id, name, type, config, active, created_at`,
		companyID, req.Name, req.Type, req.Config, req.Active,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Type, &s.Config, &s.Active, &s.CreatedAt)
	return s, err
}

func (r *SettingsRepo) ListSources(ctx context.Context, companyID string) ([]domain.RecruitmentSource, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, type, config, active, created_at
		 FROM recruitment_sources WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []domain.RecruitmentSource
	for rows.Next() {
		var s domain.RecruitmentSource
		rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.Type, &s.Config, &s.Active, &s.CreatedAt)
		sources = append(sources, s)
	}
	return sources, nil
}

func (r *SettingsRepo) UpdateSource(ctx context.Context, companyID, id string, req *domain.RecruitmentSource) (*domain.RecruitmentSource, error) {
	s := &domain.RecruitmentSource{}
	err := r.pool.QueryRow(ctx,
		`UPDATE recruitment_sources SET
		 name=COALESCE($3,name), type=COALESCE($4,type), config=COALESCE($5,config),
		 active=COALESCE($6,active) WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, type, config, active, created_at`,
		companyID, id, req.Name, req.Type, req.Config, req.Active,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Type, &s.Config, &s.Active, &s.CreatedAt)
	return s, err
}

func (r *SettingsRepo) CreateStage(ctx context.Context, companyID string, req *domain.RecruitmentStage) (*domain.RecruitmentStage, error) {
	s := &domain.RecruitmentStage{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_stages (company_id, name, category, sort_order, color, active)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, name, category, sort_order, color, active, created_at`,
		companyID, req.Name, req.Category, req.SortOrder, req.Color, req.Active,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Category, &s.SortOrder, &s.Color, &s.Active, &s.CreatedAt)
	return s, err
}

func (r *SettingsRepo) ListStages(ctx context.Context, companyID string) ([]domain.RecruitmentStage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, category, sort_order, color, active, created_at
		 FROM recruitment_stages WHERE company_id=$1 ORDER BY sort_order`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []domain.RecruitmentStage
	for rows.Next() {
		var s domain.RecruitmentStage
		rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.Category, &s.SortOrder, &s.Color, &s.Active, &s.CreatedAt)
		stages = append(stages, s)
	}
	return stages, nil
}

func (r *SettingsRepo) UpdateStage(ctx context.Context, companyID, id string, req *domain.RecruitmentStage) (*domain.RecruitmentStage, error) {
	s := &domain.RecruitmentStage{}
	err := r.pool.QueryRow(ctx,
		`UPDATE recruitment_stages SET
		 name=COALESCE($3,name), category=COALESCE($4,category),
		 sort_order=COALESCE($5,sort_order), color=COALESCE($6,color),
		 active=COALESCE($7,active) WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, category, sort_order, color, active, created_at`,
		companyID, id, req.Name, req.Category, req.SortOrder, req.Color, req.Active,
	).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Category, &s.SortOrder, &s.Color, &s.Active, &s.CreatedAt)
	return s, err
}

func (r *SettingsRepo) ReorderStages(ctx context.Context, stageIDs []string) error {
	for i, id := range stageIDs {
		_, err := r.pool.Exec(ctx,
			`UPDATE recruitment_stages SET sort_order=$2 WHERE id=$1`,
			id, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SettingsRepo) CreateTransition(ctx context.Context, req *domain.StageTransition) (*domain.StageTransition, error) {
	t := &domain.StageTransition{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recruitment_stage_transitions (company_id, from_stage_id, to_stage_id, required_actions)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, from_stage_id, to_stage_id, required_actions, created_at`,
		req.CompanyID, req.FromStageID, req.ToStageID, req.RequiredActions,
	).Scan(&t.ID, &t.CompanyID, &t.FromStageID, &t.ToStageID, &t.RequiredActions, &t.CreatedAt)
	return t, err
}

func (r *SettingsRepo) ListTransitions(ctx context.Context, companyID string) ([]domain.StageTransition, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, from_stage_id, to_stage_id, required_actions, created_at
		 FROM recruitment_stage_transitions WHERE company_id=$1`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []domain.StageTransition
	for rows.Next() {
		var t domain.StageTransition
		rows.Scan(&t.ID, &t.CompanyID, &t.FromStageID, &t.ToStageID, &t.RequiredActions, &t.CreatedAt)
		transitions = append(transitions, t)
	}
	return transitions, nil
}

func (r *SettingsRepo) DeleteTransition(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM recruitment_stage_transitions WHERE id=$1 AND company_id=$2`,
		id, companyID)
	return err
}

func (r *SettingsRepo) CreateRejectionReason(ctx context.Context, companyID string, req *domain.RejectionReason) (*domain.RejectionReason, error) {
	rr := &domain.RejectionReason{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO rejection_reasons (company_id, name, category, active, sort_order)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, company_id, name, category, active, sort_order, created_at`,
		companyID, req.Name, req.Category, req.Active, req.SortOrder,
	).Scan(&rr.ID, &rr.CompanyID, &rr.Name, &rr.Category, &rr.Active, &rr.SortOrder, &rr.CreatedAt)
	return rr, err
}

func (r *SettingsRepo) ListRejectionReasons(ctx context.Context, companyID string) ([]domain.RejectionReason, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, category, active, sort_order, created_at
		 FROM rejection_reasons WHERE company_id=$1 ORDER BY sort_order`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []domain.RejectionReason
	for rows.Next() {
		var rr domain.RejectionReason
		rows.Scan(&rr.ID, &rr.CompanyID, &rr.Name, &rr.Category, &rr.Active, &rr.SortOrder, &rr.CreatedAt)
		reasons = append(reasons, rr)
	}
	return reasons, nil
}

func (r *SettingsRepo) UpdateRejectionReason(ctx context.Context, companyID, id string, req *domain.RejectionReason) (*domain.RejectionReason, error) {
	rr := &domain.RejectionReason{}
	err := r.pool.QueryRow(ctx,
		`UPDATE rejection_reasons SET
		 name=COALESCE($3,name), category=COALESCE($4,category),
		 active=COALESCE($5,active), sort_order=COALESCE($6,sort_order)
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, category, active, sort_order, created_at`,
		companyID, id, req.Name, req.Category, req.Active, req.SortOrder,
	).Scan(&rr.ID, &rr.CompanyID, &rr.Name, &rr.Category, &rr.Active, &rr.SortOrder, &rr.CreatedAt)
	return rr, err
}
