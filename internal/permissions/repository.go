package permissions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionChecker struct {
	pool *pgxpool.Pool
}

func NewPermissionChecker(pool *pgxpool.Pool) *PermissionChecker {
	return &PermissionChecker{pool: pool}
}

func (pc *PermissionChecker) HasPermission(ctx context.Context, userID, companyID, resource, action string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			JOIN role_permissions rp ON rp.permission_id = p.id
			JOIN user_roles ur ON ur.role_id = rp.role_id
			WHERE ur.user_id = $1
			AND ur.company_id = $2
			AND p.resource = $3
			AND p.action = $4
		)`

	var exists bool
	err := pc.pool.QueryRow(ctx, query, userID, companyID, resource, action).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
