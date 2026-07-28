package performance

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

	logger.Info("Performance worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Performance worker stopped")
			return
		case <-ticker.C:
			w.detectOverdueObjectives(ctx)
			w.detectUpcomingDeadlines(ctx)
		}
	}
}

func (w *Worker) detectOverdueObjectives(ctx context.Context) {
	_, err := w.repo.pool.Exec(ctx,
		`UPDATE performance_objectives SET status='AT_RISK', updated_at=NOW()
		 WHERE status='ACTIVE' AND due_date IS NOT NULL AND due_date < CURRENT_DATE`)
	if err != nil {
		logger.Errorf("Failed to detect overdue objectives", zap.Error(err))
	}
}

func (w *Worker) detectUpcomingDeadlines(ctx context.Context) {
	rows, err := w.repo.pool.Query(ctx,
		`SELECT id, company_id, name, evaluation_deadline
		 FROM performance_cycles
		 WHERE status IN ('OPEN','IN_PROGRESS')
		 AND evaluation_deadline IS NOT NULL
		 AND evaluation_deadline BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'`)
	if err != nil {
		logger.Errorf("Failed to detect upcoming deadlines", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, companyID, name string
		var deadline time.Time
		rows.Scan(&id, &companyID, &name, &deadline)
		logger.Info("Performance cycle deadline approaching",
			zap.String("cycle", name),
			zap.String("deadline", deadline.Format("2006-01-02")),
		)
	}
}
