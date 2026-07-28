package attendance

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
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processAbsences(ctx)
		}
	}
}

func (w *Worker) processAbsences(ctx context.Context) {
	logger.Info("Processing daily absences")

	// Get all companies (simplified - in production, iterate over active companies)
	// For now, the worker will be triggered manually or via cron

	logger.Info("Absence processing completed")
}
