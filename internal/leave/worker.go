package leave

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

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	w.processReminders(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processReminders(ctx)
		}
	}
}

func (w *Worker) processReminders(ctx context.Context) {
	logger.Info("Processing leave reminders")

	requests, err := w.service.repo.GetExpiringRequests(ctx, 7)
	if err != nil {
		logger.Info("Failed to get expiring requests", logger.String("error", err.Error()))
		return
	}

	for _, req := range requests {
		daysUntil := int(time.Until(req.StartDate).Hours() / 24)
		if daysUntil <= 1 {
			logger.Info("Leave starts tomorrow",
				logger.String("employee_id", req.EmployeeID),
				logger.String("start_date", req.StartDate.Format("2006-01-02")),
			)
		} else if daysUntil <= 7 {
			logger.Info("Leave starts soon",
				logger.String("employee_id", req.EmployeeID),
				logger.String("start_date", req.StartDate.Format("2006-01-02")),
			)
		}
	}
}
