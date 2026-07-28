package payroll

import (
	"context"

	"github.com/rrhhumand/api/pkg/logger"
)

type Worker struct {
	service *Service
}

func NewWorker(service *Service) *Worker {
	return &Worker{service: service}
}

func (w *Worker) ValidatePeriods(ctx context.Context, companyID string) {
	periods, err := w.service.ListPeriods(ctx, companyID)
	if err != nil {
		logger.Errorf("Failed to list periods for validation", logger.Err(err))
		return
	}

	for _, period := range periods {
		if period.Status == "OPEN" {
			logger.Info("Open payroll period detected",
				logger.String("period_id", period.ID),
				logger.String("period_name", period.Name),
			)
		}
	}
}

func (w *Worker) CheckPendingBonuses(ctx context.Context, companyID string) {
	bonuses, err := w.service.ListBonuses(ctx, companyID, PayrollFilters{Status: "PENDING"})
	if err != nil {
		logger.Errorf("Failed to check pending bonuses", logger.Err(err))
		return
	}
	if len(bonuses) > 0 {
		logger.Info("Pending bonuses detected",
			logger.String("company_id", companyID),
			logger.Int("count", len(bonuses)),
		)
	}
}

func (w *Worker) CheckPendingAdvances(ctx context.Context, companyID string) {
	advances, err := w.service.ListAdvances(ctx, companyID, PayrollFilters{Status: "PENDING"})
	if err != nil {
		logger.Errorf("Failed to check pending advances", logger.Err(err))
		return
	}
	if len(advances) > 0 {
		logger.Info("Pending advances detected",
			logger.String("company_id", companyID),
			logger.Int("count", len(advances)),
		)
	}
}
