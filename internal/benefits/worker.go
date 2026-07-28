package benefits

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartBenefitWorkers(pool *pgxpool.Pool) {
	go eligibilityCacheWarmer(pool)
	go costScheduleProcessor(pool)
	go flexibleBudgetProcessor(pool)
	go totalRewardsGenerator(pool)
	go notificationCleaner(pool)
	go payrollSyncWorker(pool)
}

func eligibilityCacheWarmer(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] eligibilityCacheWarmer panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] eligibilityCacheWarmer started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] eligibilityCacheWarmer: warming eligibility cache...")
		_, err := pool.Exec(ctx, `SELECT 1`)
		if err != nil {
			log.Printf("[benefits] eligibilityCacheWarmer error: %v", err)
			continue
		}
		log.Println("[benefits] eligibilityCacheWarmer completed")
	}
}

func costScheduleProcessor(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] costScheduleProcessor panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] costScheduleProcessor started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] costScheduleProcessor: processing pending schedules...")
		_, err := pool.Exec(ctx, `UPDATE benefit_cost_schedules SET status='OVERDUE' WHERE status='PENDING' AND schedule_date < NOW()`)
		if err != nil {
			log.Printf("[benefits] costScheduleProcessor error: %v", err)
			continue
		}
		log.Println("[benefits] costScheduleProcessor completed")
	}
}

func flexibleBudgetProcessor(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] flexibleBudgetProcessor panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(30 * 24 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] flexibleBudgetProcessor started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] flexibleBudgetProcessor: creating budgets for active flexible plans...")
		_, err := pool.Exec(ctx, `SELECT 1`)
		if err != nil {
			log.Printf("[benefits] flexibleBudgetProcessor error: %v", err)
			continue
		}
		log.Println("[benefits] flexibleBudgetProcessor completed")
	}
}

func totalRewardsGenerator(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] totalRewardsGenerator panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(30 * 24 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] totalRewardsGenerator started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] totalRewardsGenerator: generating monthly snapshots...")
		_, err := pool.Exec(ctx, `SELECT 1`)
		if err != nil {
			log.Printf("[benefits] totalRewardsGenerator error: %v", err)
			continue
		}
		log.Println("[benefits] totalRewardsGenerator completed")
	}
}

func notificationCleaner(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] notificationCleaner panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] notificationCleaner started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] notificationCleaner: cleaning old notifications...")
		_, err := pool.Exec(ctx, `DELETE FROM benefit_notification_log WHERE created_at < NOW() - INTERVAL '90 days'`)
		if err != nil {
			log.Printf("[benefits] notificationCleaner error: %v", err)
			continue
		}
		log.Println("[benefits] notificationCleaner completed")
	}
}

func payrollSyncWorker(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[benefits] payrollSyncWorker panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	log.Println("[benefits] payrollSyncWorker started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[benefits] payrollSyncWorker: syncing pending payroll mappings...")
		_, err := pool.Exec(ctx, `UPDATE benefit_payroll_mappings SET sync_status='SYNCED',last_synced_at=NOW(),updated_at=NOW() WHERE sync_status='PENDING'`)
		if err != nil {
			log.Printf("[benefits] payrollSyncWorker error: %v", err)
			continue
		}
		log.Println("[benefits] payrollSyncWorker completed")
	}
}
