package branches

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type BranchRepository struct {
	pool *pgxpool.Pool
}

func NewBranchRepository(pool *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{pool: pool}
}

func (r *BranchRepository) Create(ctx context.Context, branch *models.Branch) error {
	query := `
		INSERT INTO branches (id, company_id, name, code, address, city, state, country, phone, email, timezone, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		branch.ID, branch.CompanyID, branch.Name, branch.Code,
		branch.Address, branch.City, branch.State, branch.Country,
		branch.Phone, branch.Email, branch.Timezone, branch.Active,
	).Scan(&branch.CreatedAt, &branch.UpdatedAt)
}

func (r *BranchRepository) FindByID(ctx context.Context, id, companyID string) (*models.Branch, error) {
	query := `
		SELECT id, company_id, name, code, address, city, state, country, phone, email, timezone, active, created_at, updated_at
		FROM branches WHERE id = $1 AND company_id = $2`
	branch := &models.Branch{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&branch.ID, &branch.CompanyID, &branch.Name, &branch.Code,
		&branch.Address, &branch.City, &branch.State, &branch.Country,
		&branch.Phone, &branch.Email, &branch.Timezone, &branch.Active,
		&branch.CreatedAt, &branch.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (r *BranchRepository) List(ctx context.Context, companyID string, offset, limit int, search string) ([]models.Branch, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM branches WHERE company_id = $1`
	args := []interface{}{companyID}

	if search != "" {
		countQuery += ` AND (name ILIKE $2 OR code ILIKE $2)`
		args = append(args, "%"+search+"%")
	}

	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, company_id, name, code, address, city, state, country, phone, email, timezone, active, created_at, updated_at
		FROM branches WHERE company_id = $1`

	listArgs := []interface{}{companyID}
	paramIdx := 2
	if search != "" {
		query += ` AND (name ILIKE $` + string(rune('0'+paramIdx)) + ` OR code ILIKE $` + string(rune('0'+paramIdx)) + `)`
		listArgs = append(listArgs, "%"+search+"%")
		paramIdx++
	}

	query += ` ORDER BY name ASC LIMIT $` + string(rune('0'+paramIdx)) + ` OFFSET $` + string(rune('0'+paramIdx+1))
	listArgs = append(listArgs, limit, offset)

	rows, err := r.pool.Query(ctx, query, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		var b models.Branch
		if err := rows.Scan(
			&b.ID, &b.CompanyID, &b.Name, &b.Code,
			&b.Address, &b.City, &b.State, &b.Country,
			&b.Phone, &b.Email, &b.Timezone, &b.Active,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		branches = append(branches, b)
	}
	return branches, total, nil
}

func (r *BranchRepository) Update(ctx context.Context, branch *models.Branch) error {
	query := `
		UPDATE branches
		SET name=$1, code=$2, address=$3, city=$4, state=$5, country=$6, phone=$7, email=$8, timezone=$9, active=$10, updated_at=NOW()
		WHERE id=$11 AND company_id=$12`
	_, err := r.pool.Exec(ctx, query,
		branch.Name, branch.Code, branch.Address, branch.City,
		branch.State, branch.Country, branch.Phone, branch.Email,
		branch.Timezone, branch.Active, branch.ID, branch.CompanyID,
	)
	return err
}

func (r *BranchRepository) Delete(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM branches WHERE id=$1 AND company_id=$2`, id, companyID)
	return err
}
