package payroll

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	svc    *Service
	repo   *Repository
	stopCh chan struct{}
	doneCh chan struct{}
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

	w.checkPendingRuns()
	w.checkOpenPeriods()

	for {
		select {
		case <-w.stopCh:
			log.Println("[payroll-worker] stopping")
			return
		case <-ticker.C:
			w.checkPendingRuns()
			w.checkOpenPeriods()
		}
	}
}

func (w *Worker) checkPendingRuns() {
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Future: process queue-based calculation jobs
	// For now, placeholder for async calculation logic
	log.Println("[payroll-worker] checking pending runs")
}

func (w *Worker) checkOpenPeriods() {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Future: auto-close periods past end_date without activity
	log.Println("[payroll-worker] checking open periods")
}
