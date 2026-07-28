package positions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type PositionRepository struct {
	pool *pgxpool.Pool
}

func NewPositionRepository(pool *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{pool: pool}
}

func (r *PositionRepository) Create(ctx context.Context, pos *models.Position) error {
	query := `
		INSERT INTO positions (id, company_id, name, code, description, department_id, level, min_salary, max_salary, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		pos.ID, pos.CompanyID, pos.Name, pos.Code,
		pos.Description, pos.DepartmentID, pos.Level,
		pos.MinSalary, pos.MaxSalary, pos.Active,
	).Scan(&pos.CreatedAt, &pos.UpdatedAt)
}

func (r *PositionRepository) FindByID(ctx context.Context, id, companyID string) (*models.Position, error) {
	query := `
		SELECT id, company_id, name, code, description, department_id, level, min_salary, max_salary, active, created_at, updated_at
		FROM positions WHERE id = $1 AND company_id = $2`
	pos := &models.Position{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&pos.ID, &pos.CompanyID, &pos.Name, &pos.Code,
		&pos.Description, &pos.DepartmentID, &pos.Level,
		&pos.MinSalary, &pos.MaxSalary, &pos.Active,
		&pos.CreatedAt, &pos.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pos, nil
}

func (r *PositionRepository) List(ctx context.Context, companyID string, offset, limit int, search string, departmentID string) ([]models.Position, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM positions WHERE company_id = $1`
	args := []interface{}{companyID}
	argIdx := 2

	if search != "" {
		countQuery += ` AND (name ILIKE $` + string(rune('0'+argIdx)) + ` OR code ILIKE $` + string(rune('0'+argIdx)) + `)`
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if departmentID != "" {
		countQuery += ` AND department_id = $` + string(rune('0'+argIdx))
		args = append(args, departmentID)
		argIdx++
	}

	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, company_id, name, code, description, department_id, level, min_salary, max_salary, active, created_at, updated_at
		FROM positions WHERE company_id = $1`
	listArgs := []interface{}{companyID}
	paramIdx := 2

	if search != "" {
		query += ` AND (name ILIKE $` + string(rune('0'+paramIdx)) + ` OR code ILIKE $` + string(rune('0'+paramIdx)) + `)`
		listArgs = append(listArgs, "%"+search+"%")
		paramIdx++
	}
	if departmentID != "" {
		query += ` AND department_id = $` + string(rune('0'+paramIdx))
		listArgs = append(listArgs, departmentID)
		paramIdx++
	}

	query += ` ORDER BY level ASC, name ASC LIMIT $` + string(rune('0'+paramIdx)) + ` OFFSET $` + string(rune('0'+paramIdx+1))
	listArgs = append(listArgs, limit, offset)

	rows, err := r.pool.Query(ctx, query, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var positions []models.Position
	for rows.Next() {
		var p models.Position
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.Name, &p.Code,
			&p.Description, &p.DepartmentID, &p.Level,
			&p.MinSalary, &p.MaxSalary, &p.Active,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		positions = append(positions, p)
	}
	return positions, total, nil
}

func (r *PositionRepository) Update(ctx context.Context, pos *models.Position) error {
	query := `
		UPDATE positions
		SET name=$1, code=$2, description=$3, department_id=$4, level=$5, min_salary=$6, max_salary=$7, active=$8, updated_at=NOW()
		WHERE id=$9 AND company_id=$10`
	_, err := r.pool.Exec(ctx, query,
		pos.Name, pos.Code, pos.Description, pos.DepartmentID,
		pos.Level, pos.MinSalary, pos.MaxSalary, pos.Active,
		pos.ID, pos.CompanyID,
	)
	return err
}

func (r *PositionRepository) Delete(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM positions WHERE id=$1 AND company_id=$2`, id, companyID)
	return err
}
