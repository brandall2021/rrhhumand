package roles

import (
	"context"

	"github.com/rrhhumand/api/internal/models"
)

type RoleService struct {
	repo *RoleRepository
}

func NewRoleService(repo *RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) FindByName(ctx context.Context, name string) (*models.Role, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *RoleService) GetAll(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetAll(ctx)
}

func (s *RoleService) GetPermissionsByUserRole(ctx context.Context, userID, companyID string) ([]models.Permission, error) {
	return s.repo.GetPermissionsByUserRole(ctx, userID, companyID)
}
