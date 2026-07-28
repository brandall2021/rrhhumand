package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type ScoringRepo struct {
	pool *pgxpool.Pool
}

func NewScoringRepo(pool *pgxpool.Pool) *ScoringRepo {
	return &ScoringRepo{pool: pool}
}

func (r *ScoringRepo) CreateScoringModel(ctx context.Context, companyID string, req *domain.ScoringModel) (*domain.ScoringModel, error) {
	m := &domain.ScoringModel{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO scoring_models (company_id, name, description, config, is_default, active)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, name, description, config, is_default, active, created_at`,
		companyID, req.Name, req.Description, req.Config, req.IsDefault, req.Active,
	).Scan(&m.ID, &m.CompanyID, &m.Name, &m.Description, &m.Config, &m.IsDefault, &m.Active, &m.CreatedAt)
	return m, err
}

func (r *ScoringRepo) GetScoringModel(ctx context.Context, companyID, id string) (*domain.ScoringModel, error) {
	m := &domain.ScoringModel{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, config, is_default, active, created_at
		 FROM scoring_models WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&m.ID, &m.CompanyID, &m.Name, &m.Description, &m.Config, &m.IsDefault, &m.Active, &m.CreatedAt)
	return m, err
}

func (r *ScoringRepo) ListScoringModels(ctx context.Context, companyID string) ([]domain.ScoringModel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, config, is_default, active, created_at
		 FROM scoring_models WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []domain.ScoringModel
	for rows.Next() {
		var m domain.ScoringModel
		rows.Scan(&m.ID, &m.CompanyID, &m.Name, &m.Description, &m.Config, &m.IsDefault, &m.Active, &m.CreatedAt)
		models = append(models, m)
	}
	return models, nil
}

func (r *ScoringRepo) UpdateScoringModel(ctx context.Context, companyID, id string, req *domain.ScoringModel) (*domain.ScoringModel, error) {
	m := &domain.ScoringModel{}
	err := r.pool.QueryRow(ctx,
		`UPDATE scoring_models SET
		 name=COALESCE($3,name), description=COALESCE($4,description),
		 config=COALESCE($5,config), is_default=COALESCE($6,is_default),
		 active=COALESCE($7,active) WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, description, config, is_default, active, created_at`,
		companyID, id, req.Name, req.Description, req.Config, req.IsDefault, req.Active,
	).Scan(&m.ID, &m.CompanyID, &m.Name, &m.Description, &m.Config, &m.IsDefault, &m.Active, &m.CreatedAt)
	return m, err
}

func (r *ScoringRepo) DeleteScoringModel(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM scoring_models WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *ScoringRepo) AddCriterion(ctx context.Context, req *domain.ScoringCriterion) (*domain.ScoringCriterion, error) {
	c := &domain.ScoringCriterion{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO scoring_criteria (model_id, name, field, weight, scoring_type, config, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, model_id, name, field, weight, scoring_type, config, sort_order`,
		req.ModelID, req.Name, req.Field, req.Weight, req.ScoringType, req.Config, req.SortOrder,
	).Scan(&c.ID, &c.ModelID, &c.Name, &c.Field, &c.Weight, &c.ScoringType, &c.Config, &c.SortOrder)
	return c, err
}

func (r *ScoringRepo) UpdateCriterion(ctx context.Context, id string, req *domain.ScoringCriterion) (*domain.ScoringCriterion, error) {
	c := &domain.ScoringCriterion{}
	err := r.pool.QueryRow(ctx,
		`UPDATE scoring_criteria SET name=COALESCE($2,name), field=COALESCE($3,field),
		 weight=COALESCE($4,weight), scoring_type=COALESCE($5,scoring_type),
		 config=COALESCE($6,config), sort_order=COALESCE($7,sort_order) WHERE id=$1
		 RETURNING id, model_id, name, field, weight, scoring_type, config, sort_order`,
		id, req.Name, req.Field, req.Weight, req.ScoringType, req.Config, req.SortOrder,
	).Scan(&c.ID, &c.ModelID, &c.Name, &c.Field, &c.Weight, &c.ScoringType, &c.Config, &c.SortOrder)
	return c, err
}

func (r *ScoringRepo) DeleteCriterion(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM scoring_criteria WHERE id=$1`, id)
	return err
}

func (r *ScoringRepo) ListCriteria(ctx context.Context, modelID string) ([]domain.ScoringCriterion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, model_id, name, field, weight, scoring_type, config, sort_order
		 FROM scoring_criteria WHERE model_id=$1 ORDER BY sort_order`, modelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var criteria []domain.ScoringCriterion
	for rows.Next() {
		var c domain.ScoringCriterion
		rows.Scan(&c.ID, &c.ModelID, &c.Name, &c.Field, &c.Weight, &c.ScoringType, &c.Config, &c.SortOrder)
		criteria = append(criteria, c)
	}
	return criteria, nil
}

func (r *ScoringRepo) SaveMatchingResult(ctx context.Context, req *domain.MatchingResult) (*domain.MatchingResult, error) {
	mr := &domain.MatchingResult{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO matching_results (candidate_id, position_id, overall_score, skill_score, experience_score, education_score, culture_score, details)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT (candidate_id, position_id) DO UPDATE SET
		 overall_score=$3, skill_score=$4, experience_score=$5, education_score=$6, culture_score=$7, details=$8, matched_at=NOW()
		 RETURNING id, candidate_id, position_id, overall_score, skill_score, experience_score, education_score, culture_score, details, matched_at`,
		req.CandidateID, req.PositionID, req.OverallScore, req.SkillScore, req.ExperienceScore,
		req.EducationScore, req.CultureScore, req.Details,
	).Scan(&mr.ID, &mr.CandidateID, &mr.PositionID, &mr.OverallScore, &mr.SkillScore,
		&mr.ExperienceScore, &mr.EducationScore, &mr.CultureScore, &mr.Details, &mr.MatchedAt)
	return mr, err
}

func (r *ScoringRepo) GetMatchingResult(ctx context.Context, candidateID, positionID string) (*domain.MatchingResult, error) {
	mr := &domain.MatchingResult{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, candidate_id, position_id, overall_score, skill_score, experience_score, education_score, culture_score, details, matched_at
		 FROM matching_results WHERE candidate_id=$1 AND position_id=$2`, candidateID, positionID,
	).Scan(&mr.ID, &mr.CandidateID, &mr.PositionID, &mr.OverallScore, &mr.SkillScore,
		&mr.ExperienceScore, &mr.EducationScore, &mr.CultureScore, &mr.Details, &mr.MatchedAt)
	return mr, err
}

func (r *ScoringRepo) ListMatchingResults(ctx context.Context, positionID string) ([]domain.MatchingResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, position_id, overall_score, skill_score, experience_score, education_score, culture_score, details, matched_at
		 FROM matching_results WHERE position_id=$1 ORDER BY overall_score DESC NULLS LAST`, positionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.MatchingResult
	for rows.Next() {
		var mr domain.MatchingResult
		rows.Scan(&mr.ID, &mr.CandidateID, &mr.PositionID, &mr.OverallScore, &mr.SkillScore,
			&mr.ExperienceScore, &mr.EducationScore, &mr.CultureScore, &mr.Details, &mr.MatchedAt)
		results = append(results, mr)
	}
	return results, nil
}

func (r *ScoringRepo) DeleteMatchingResults(ctx context.Context, positionID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM matching_results WHERE position_id=$1`, positionID)
	return err
}
