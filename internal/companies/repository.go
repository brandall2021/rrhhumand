package companies

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}

func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	query := `
		INSERT INTO companies (id, name, slug, logo_url, plan, settings, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		company.ID, company.Name, company.Slug, company.LogoURL,
		company.Plan, company.Settings, company.Active,
	).Scan(&company.CreatedAt, &company.UpdatedAt)
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*models.Company, error) {
	query := `
		SELECT id, name, slug, logo_url, plan, settings, active, created_at, updated_at
		FROM companies WHERE id = $1`
	company := &models.Company{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&company.ID, &company.Name, &company.Slug, &company.LogoURL,
		&company.Plan, &company.Settings, &company.Active,
		&company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (r *CompanyRepository) FindBySlug(ctx context.Context, slug string) (*models.Company, error) {
	query := `
		SELECT id, name, slug, logo_url, plan, settings, active, created_at, updated_at
		FROM companies WHERE slug = $1`
	company := &models.Company{}
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&company.ID, &company.Name, &company.Slug, &company.LogoURL,
		&company.Plan, &company.Settings, &company.Active,
		&company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	query := `
		UPDATE companies
		SET name = $1, slug = $2, logo_url = $3, plan = $4, settings = $5, active = $6, updated_at = NOW()
		WHERE id = $7`
	_, err := r.pool.Exec(ctx, query,
		company.Name, company.Slug, company.LogoURL,
		company.Plan, company.Settings, company.Active, company.ID,
	)
	return err
}

func (r *CompanyRepository) List(ctx context.Context, offset, limit int) ([]models.Company, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM companies`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, slug, logo_url, plan, settings, active, created_at, updated_at
		FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.LogoURL,
			&c.Plan, &c.Settings, &c.Active,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, c)
	}
	return companies, total, nil
}
