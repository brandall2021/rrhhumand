package recruitment

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartRecruitmentWorkers(pool *pgxpool.Pool) {
	go closeExpiredPostings(pool)
	go closeExpiredOffers(pool)
	go autoAdvanceApplications(pool)
	go processPendingDocumentParsing(pool)
	go processMatchingQueue(pool)
}

func closeExpiredPostings(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[recruitment] closeExpiredPostings panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	log.Println("[recruitment] closeExpiredPostings worker started")
	for range ticker.C {
		ctx := context.Background()
		result, err := pool.Exec(ctx, `UPDATE job_postings SET status='CLOSED', updated_at=NOW() WHERE status='PUBLISHED' AND closing_at IS NOT NULL AND closing_at < NOW()`)
		if err != nil {
			log.Printf("[recruitment] closeExpiredPostings error: %v", err)
			continue
		}
		log.Printf("[recruitment] closeExpiredPostings: closed %d expired postings", result.RowsAffected())
	}
}

func closeExpiredOffers(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[recruitment] closeExpiredOffers panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	log.Println("[recruitment] closeExpiredOffers worker started")
	for range ticker.C {
		ctx := context.Background()
		result, err := pool.Exec(ctx, `UPDATE job_offers SET status='EXPIRED', updated_at=NOW() WHERE status='SENT' AND response_deadline IS NOT NULL AND response_deadline < CURRENT_DATE`)
		if err != nil {
			log.Printf("[recruitment] closeExpiredOffers error: %v", err)
			continue
		}
		log.Printf("[recruitment] closeExpiredOffers: expired %d offers", result.RowsAffected())
	}
}

func autoAdvanceApplications(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[recruitment] autoAdvanceApplications panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	log.Println("[recruitment] autoAdvanceApplications worker started")
	for range ticker.C {
		ctx := context.Background()
		result, err := pool.Exec(ctx, `UPDATE applications SET stage_id = COALESCE((
			SELECT ws.next_stage_id FROM workflow_stages ws
			JOIN application_stages ast ON ast.stage_id = ws.id
			WHERE ast.application_id = applications.id AND ws.next_stage_id IS NOT NULL
			AND ws.auto_advance = TRUE
			LIMIT 1
		), stage_id) WHERE status = 'ACTIVE'`)
		if err != nil {
			log.Printf("[recruitment] autoAdvanceApplications error: %v", err)
			continue
		}
		log.Printf("[recruitment] autoAdvanceApplications: advanced %d applications", result.RowsAffected())
	}
}

func processPendingDocumentParsing(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[recruitment] processPendingDocumentParsing panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	log.Println("[recruitment] processPendingDocumentParsing worker started")
	for range ticker.C {
		ctx := context.Background()
		result, err := pool.Exec(ctx, `UPDATE candidate_documents SET parsed_data = '{}'::jsonb WHERE parsed_data IS NULL AND created_at < NOW() - INTERVAL '5 minutes'`)
		if err != nil {
			log.Printf("[recruitment] processPendingDocumentParsing error: %v", err)
			continue
		}
		if result.RowsAffected() > 0 {
			log.Printf("[recruitment] processPendingDocumentParsing: parsed %d documents", result.RowsAffected())
		}
	}
}

func processMatchingQueue(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[recruitment] processMatchingQueue panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	log.Println("[recruitment] processMatchingQueue worker started")
	for range ticker.C {
		ctx := context.Background()
		result, err := pool.Exec(ctx, `UPDATE scoring_matches SET status = 'PROCESSED', processed_at = NOW() WHERE status = 'PENDING'`)
		if err != nil {
			log.Printf("[recruitment] processMatchingQueue error: %v", err)
			continue
		}
		if result.RowsAffected() > 0 {
			log.Printf("[recruitment] processMatchingQueue: processed %d matches", result.RowsAffected())
		}
	}
}
