package compensation

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	svc     *Service
	repo    *Repository
	stopCh  chan struct{}
	doneCh  chan struct{}
}

func NewWorker(svc *Service, repo *Repository) *Worker {
	return &Worker{
		svc:    svc,
		repo:   repo,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (w *Worker) Start(interval time.Duration) {
	go w.run(interval)
}

func (w *Worker) Stop() {
	close(w.stopCh)
	<-w.doneCh
}

func (w *Worker) run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(w.doneCh)

	if err := w.checkBenefitExpirations(); err != nil {
		log.Printf("[compensation-worker] initial benefit check error: %v", err)
	}
	if err := w.checkBudgetThresholds(); err != nil {
		log.Printf("[compensation-worker] initial budget check error: %v", err)
	}

	for {
		select {
		case <-w.stopCh:
			log.Println("[compensation-worker] stopping")
			return
		case <-ticker.C:
			if err := w.checkBenefitExpirations(); err != nil {
				log.Printf("[compensation-worker] benefit check error: %v", err)
			}
			if err := w.checkBudgetThresholds(); err != nil {
				log.Printf("[compensation-worker] budget check error: %v", err)
			}
			if err := w.processDomainEvents(); err != nil {
				log.Printf("[compensation-worker] event processing error: %v", err)
			}
		}
	}
}

func (w *Worker) checkBenefitExpirations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	expiring, err := w.repo.GetExpiringBenefits(ctx)
	if err != nil {
		return err
	}

	for _, eb := range expiring {
		if err := w.repo.NotifyBenefitExpiration(ctx, eb); err != nil {
			log.Printf("[compensation-worker] failed to notify benefit expiration %s: %v", eb.ID, err)
		}
	}
	return nil
}

func (w *Worker) checkBudgetThresholds() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	budgets, err := w.repo.GetBudgetsNearThreshold(ctx)
	if err != nil {
		return err
	}

	for _, b := range budgets {
		if err := w.repo.NotifyBudgetAlert(ctx, b); err != nil {
			log.Printf("[compensation-worker] failed to notify budget alert %s: %v", b.ID, err)
		}
	}
	return nil
}

func (w *Worker) processDomainEvents() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	events, err := w.repo.GetUnprocessedDomainEvents(ctx, "compensation")
	if err != nil {
		return err
	}

	for _, evt := range events {
		switch evt.EventType {
		case "review.opened":
			log.Printf("[compensation-worker] review opened event: %s", evt.ID)
		case "adjustment.applied":
			log.Printf("[compensation-worker] adjustment applied event: %s", evt.ID)
		case "bonus.approved":
			log.Printf("[compensation-worker] bonus approved event: %s", evt.ID)
		default:
			log.Printf("[compensation-worker] unknown event type: %s", evt.EventType)
		}
		if err := w.repo.MarkDomainEventProcessed(ctx, evt.ID); err != nil {
			log.Printf("[compensation-worker] failed to mark event %s as processed: %v", evt.ID, err)
		}
	}
	return nil
}
