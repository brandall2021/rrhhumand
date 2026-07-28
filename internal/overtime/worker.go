package overtime

import (
	"context"
	"time"

	"github.com/rrhhumand/api/pkg/logger"
)

type Worker struct {
	service *Service
}

func NewWorker(service *Service) *Worker {
	return &Worker{service: service}
}

func (w *Worker) DetectDailyOvertime(ctx context.Context, companyID string) {
	yesterday := time.Now().AddDate(0, 0, -1)
	from := yesterday.Format("2006-01-02")
	to := from

	policy, _ := w.service.repo.GetActivePolicy(ctx, companyID)
	if policy == nil {
		return
	}

	records, count, err := w.service.DetectOvertime(ctx, companyID, from, to)
	if err != nil {
		logger.Errorf("Failed to detect overtime", logger.Err(err))
		return
	}
	logger.Info("Overtime detection completed",
		logger.String("company_id", companyID),
		logger.String("date", from),
		logger.Int("records", count),
	)
	_ = records
}

func (w *Worker) ProcessPendingApprovals(ctx context.Context, companyID string) {
	requests, err := w.service.ListRequests(ctx, companyID, OvertimeFilters{Status: "PENDING"})
	if err != nil {
		logger.Errorf("Failed to list pending overtime requests", logger.Err(err))
		return
	}
	logger.Info("Pending overtime requests",
		logger.String("company_id", companyID),
		logger.Int("count", len(requests)),
	)
}

func (w *Worker) Run(ctx context.Context, companyID string) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.DetectDailyOvertime(ctx, companyID)
			w.ProcessPendingApprovals(ctx, companyID)
		}
	}
}
