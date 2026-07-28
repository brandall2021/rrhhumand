package recruitment

import (
	"context"
	"time"

	"github.com/rrhhumand/api/pkg/logger"
	"go.uber.org/zap"
)

type Worker struct {
	repo *Repository
}

func NewWorker(repo *Repository) *Worker {
	return &Worker{repo: repo}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	logger.Info("Recruitment worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Recruitment worker stopped")
			return
		case <-ticker.C:
			w.closeExpiredPostings(ctx)
			w.closeExpiredOffers(ctx)
		}
	}
}

func (w *Worker) closeExpiredPostings(ctx context.Context) {
	_, err := w.repo.pool.Exec(ctx,
		`UPDATE job_postings SET status='CLOSED', closing_at=NOW()
		 WHERE status='PUBLISHED' AND closing_at IS NOT NULL AND closing_at < NOW()`)
	if err != nil {
		logger.Errorf("Failed to close expired postings", zap.Error(err))
	}
}

func (w *Worker) closeExpiredOffers(ctx context.Context) {
	_, err := w.repo.pool.Exec(ctx,
		`UPDATE job_offers SET status='EXPIRED', updated_at=NOW()
		 WHERE status='SENT' AND response_deadline IS NOT NULL AND response_deadline < CURRENT_DATE`)
	if err != nil {
		logger.Errorf("Failed to close expired offers", zap.Error(err))
	}
}
