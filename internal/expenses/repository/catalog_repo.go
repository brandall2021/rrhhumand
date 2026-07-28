package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

func repoErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("expenses_repo.%s: %w", op, err)
}

type CatalogRepo struct {
	pool *pgxpool.Pool
}

func NewCatalogRepo(pool *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{pool: pool}
}

func (r *CatalogRepo) CreateCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	q := `INSERT INTO expense_categories (id,company_id,name,description,parent_id,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.ParentID, c.IsActive, c.CreatedBy)
	return repoErr("CreateCategory", err)
}

func (r *CatalogRepo) GetCategory(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseCategory, error) {
	q := `SELECT id,company_id,name,description,parent_id,is_active,created_by,created_at,updated_at
		FROM expense_categories WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.ExpenseCategory
	err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.ParentID, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCategory", err)
	}
	return &c, nil
}

func (r *CatalogRepo) ListCategories(ctx context.Context, companyID uuid.UUID, parentID *uuid.UUID) ([]domain.ExpenseCategory, error) {
	q := `SELECT id,company_id,name,description,parent_id,is_active,created_by,created_at,updated_at
		FROM expense_categories WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if parentID != nil {
		q += fmt.Sprintf(" AND parent_id=$%d", n)
		args = append(args, *parentID)
		n++
	}
	q += " ORDER BY name"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListCategories", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpenseCategory, error) {
		var c domain.ExpenseCategory
		err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.ParentID, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *CatalogRepo) UpdateCategory(ctx context.Context, c *domain.ExpenseCategory) error {
	q := `UPDATE expense_categories SET name=$1,description=$2,parent_id=$3,is_active=$4,updated_at=NOW()
		WHERE id=$5 AND company_id=$6`
	_, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.ParentID, c.IsActive, c.ID, c.CompanyID)
	return repoErr("UpdateCategory", err)
}

func (r *CatalogRepo) DeleteCategory(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_categories WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteCategory", err)
}

func (r *CatalogRepo) CreatePaymentMethod(ctx context.Context, m *domain.ExpensePaymentMethod) error {
	q := `INSERT INTO expense_payment_methods (id,company_id,name,description,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, m.ID, m.CompanyID, m.Name, m.Description, m.IsActive, m.CreatedBy)
	return repoErr("CreatePaymentMethod", err)
}

func (r *CatalogRepo) ListPaymentMethods(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePaymentMethod, error) {
	q := `SELECT id,company_id,name,description,is_active,created_by,created_at,updated_at
		FROM expense_payment_methods WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListPaymentMethods", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePaymentMethod, error) {
		var m domain.ExpensePaymentMethod
		err := row.Scan(&m.ID, &m.CompanyID, &m.Name, &m.Description, &m.IsActive, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
		return m, err
	})
}

func (r *CatalogRepo) UpdatePaymentMethod(ctx context.Context, m *domain.ExpensePaymentMethod) error {
	q := `UPDATE expense_payment_methods SET name=$1,description=$2,is_active=$3,updated_at=NOW()
		WHERE id=$4 AND company_id=$5`
	_, err := r.pool.Exec(ctx, q, m.Name, m.Description, m.IsActive, m.ID, m.CompanyID)
	return repoErr("UpdatePaymentMethod", err)
}
