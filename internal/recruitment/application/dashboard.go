package application

import (
	"context"

	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepo
}

func NewDashboardService(dashboardRepo *repository.DashboardRepo) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

func (s *DashboardService) GetDashboard(ctx context.Context, companyID string) (*repository.DashboardStats, error) {
	const op = "GetDashboard"
	return s.dashboardRepo.GetDashboardStats(ctx, companyID)
}

func (s *DashboardService) GetFunnel(ctx context.Context, companyID string) (*repository.FunnelData, error) {
	const op = "GetFunnel"
	return s.dashboardRepo.GetFunnelData(ctx, companyID)
}

func (s *DashboardService) GetTimeToHire(ctx context.Context, companyID string) (*repository.TimeToHire, error) {
	const op = "GetTimeToHire"
	return s.dashboardRepo.GetTimeToHire(ctx, companyID)
}
