package training

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	svc *Service
	log *zap.Logger
}

func NewWorker(svc *Service, log *zap.Logger) *Worker {
	return &Worker{svc: svc, log: log}
}

func (w *Worker) Start(ctx context.Context, companyID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	w.log.Info("training worker started", zap.String("company", companyID), zap.Duration("interval", interval))

	for {
		select {
		case <-ctx.Done():
			w.log.Info("training worker stopped")
			return
		case <-ticker.C:
			w.run(ctx, companyID)
		}
	}
}

func (w *Worker) run(ctx context.Context, companyID string) {
	w.log.Debug("training worker cycle", zap.String("company", companyID))

	if companyID == "" {
		w.log.Warn("training worker skipping cycle: no company context provided")
		return
	}

	// Process pending events (notifications)
	if err := w.svc.ProcessPendingEvents(ctx, companyID); err != nil {
		w.log.Error("process pending events", zap.Error(err))
	}

	// Check overdue enrollments
	w.svc.CheckOverdueEnrollments(ctx, companyID)

	// Check certificate expirations
	w.svc.CheckCertificateExpirations(ctx, companyID)
}
