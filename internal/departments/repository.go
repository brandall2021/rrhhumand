package departments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type DepartmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{pool: pool}
}

func (r *DepartmentRepository) Create(ctx context.Context, dept *models.Department) error {
	query := `
		INSERT INTO departments (id, company_id, name, code, description, branch_id, parent_id, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		dept.ID, dept.CompanyID, dept.Name, dept.Code,
		dept.Description, dept.BranchID, dept.ParentID, dept.Active,
	).Scan(&dept.CreatedAt, &dept.UpdatedAt)
}

func (r *DepartmentRepository) FindByID(ctx context.Context, id, companyID string) (*models.Department, error) {
	query := `
		SELECT id, company_id, name, code, description, branch_id, parent_id, active, created_at, updated_at
		FROM departments WHERE id = $1 AND company_id = $2`
	dept := &models.Department{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&dept.ID, &dept.CompanyID, &dept.Name, &dept.Code,
		&dept.Description, &dept.BranchID, &dept.ParentID, &dept.Active,
		&dept.CreatedAt, &dept.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (r *DepartmentRepository) List(ctx context.Context, companyID string, offset, limit int, search string) ([]models.Department, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM departments WHERE company_id = $1`
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
		SELECT id, company_id, name, code, description, branch_id, parent_id, active, created_at, updated_at
		FROM departments WHERE company_id = $1`

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

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.Name, &d.Code,
			&d.Description, &d.BranchID, &d.ParentID, &d.Active,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		depts = append(depts, d)
	}
	return depts, total, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, dept *models.Department) error {
	query := `
		UPDATE departments
		SET name=$1, code=$2, description=$3, branch_id=$4, parent_id=$5, active=$6, updated_at=NOW()
		WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, query,
		dept.Name, dept.Code, dept.Description, dept.BranchID,
		dept.ParentID, dept.Active, dept.ID, dept.CompanyID,
	)
	return err
}

func (r *DepartmentRepository) Delete(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM departments WHERE id=$1 AND company_id=$2`, id, companyID)
	return err
}
