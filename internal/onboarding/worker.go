package onboarding

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

	logger.Info("Onboarding worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Onboarding worker stopped")
			return
		case <-ticker.C:
			w.processOverdueTasks(ctx)
			w.processUpcomingMilestones(ctx)
		}
	}
}

func (w *Worker) processOverdueTasks(ctx context.Context) {
	_, err := w.repo.pool.Exec(ctx,
		`INSERT INTO notifications (company_id, user_id, title, body, notification_type, channel, reference_type, reference_id)
		 SELECT ot.company_id, COALESCE(ot.responsible_id, ot.employee_id),
		 'Tarea vencida', 'La tarea "' || ot.title || '" está vencida desde ' || ot.due_date::text,
		 'TASK_OVERDUE', 'IN_APP', 'onboarding_task', ot.id
		 FROM onboarding_tasks ot
		 WHERE ot.due_date < CURRENT_DATE AND ot.status NOT IN ('COMPLETED','CANCELLED')`)
	if err != nil {
		logger.Errorf("Failed to process overdue tasks", zap.Error(err))
	}
}

func (w *Worker) processUpcomingMilestones(ctx context.Context) {
	_, err := w.repo.pool.Exec(ctx,
		`INSERT INTO notifications (company_id, user_id, title, body, notification_type, channel, reference_type, reference_id)
		 SELECT om.company_id, om.responsible_id,
		 'Hito próximo', 'El hito "' || om.title || '" vence el ' || om.due_date::text,
		 'MILESTONE_DUE', 'IN_APP', 'onboarding_milestone', om.id
		 FROM onboarding_milestones om
		 WHERE om.due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
		 AND om.status = 'PENDING'`)
	if err != nil {
		logger.Errorf("Failed to process upcoming milestones", zap.Error(err))
	}
}
