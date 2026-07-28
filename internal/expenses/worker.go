package expenses

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartExpenseWorkers(pool *pgxpool.Pool) {
	go policyWorker(pool)
	go reminderWorker(pool)
	go budgetWorker(pool)
	go auditCleaner(pool)
	go duplicateCheckWorker(pool)
}

func policyWorker(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expenses] policyWorker panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	log.Println("[expenses] policyWorker started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[expenses] policyWorker: checking pending expenses policy compliance...")
		_, err := pool.Exec(ctx, `
			UPDATE expenses SET is_policy_compliant = TRUE, policy_status = 'COMPLIANT', updated_at = NOW()
			WHERE status = 'SUBMITTED' AND is_policy_compliant = FALSE
			  AND category_id IN (
				SELECT DISTINCT pr.category_id FROM expense_policy_rules pr
				JOIN expense_policies p ON p.id = pr.policy_id
				WHERE p.is_active = TRUE AND pr.is_active = TRUE AND pr.max_amount IS NULL
			  )`)
		if err != nil {
			log.Printf("[expenses] policyWorker error: %v", err)
			continue
		}
		log.Println("[expenses] policyWorker completed")
	}
}

func reminderWorker(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expenses] reminderWorker panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	log.Println("[expenses] reminderWorker started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[expenses] reminderWorker: sending pending reminders...")

		_, err := pool.Exec(ctx, `
			INSERT INTO expense_notification_log (id, company_id, employee_id, notification_type, channel, title, body, sent_at, created_at)
			SELECT gen_random_uuid(), t.company_id, t.employee_id, 'TRAVEL_REPORT_REMINDER', 'in_app',
				   'Expense Report Reminder', 'Please submit your expense report for the completed travel.',
				   NOW(), NOW()
			FROM travels t
			WHERE t.status = 'COMPLETED'
			  AND t.return_date < NOW() - INTERVAL '3 days'
			  AND NOT EXISTS (
				SELECT 1 FROM expense_reports er WHERE er.travel_id = t.id AND er.status != 'CANCELLED'
			  )`)
		if err != nil {
			log.Printf("[expenses] reminderWorker travel reminder error: %v", err)
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO expense_notification_log (id, company_id, employee_id, notification_type, channel, title, body, sent_at, created_at)
			SELECT gen_random_uuid(), ea.company_id, ea.employee_id, 'ADVANCE_REMINDER', 'in_app',
				   'Advance Settlement Reminder', 'You have an unsettled advance that requires attention.',
				   NOW(), NOW()
			FROM expense_advances ea
			WHERE ea.status = 'PAID'
			  AND ea.paid_date < NOW() - INTERVAL '15 days'
			  AND NOT EXISTS (
				SELECT 1 FROM expense_reports er WHERE er.advance_id = ea.id AND er.status IN ('SUBMITTED','APPROVED','PAID')
			  )`)
		if err != nil {
			log.Printf("[expenses] reminderWorker advance reminder error: %v", err)
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO expense_notification_log (id, company_id, employee_id, notification_type, channel, title, body, sent_at, created_at)
			SELECT gen_random_uuid(), e.company_id, e.employee_id, 'UNSUBMITTED_EXPENSE_REMINDER', 'in_app',
				   'Expense Submission Reminder', 'You have unsubmitted expenses that need attention.',
				   NOW(), NOW()
			FROM expenses e
			WHERE e.status = 'DRAFT'
			  AND e.created_at < NOW() - INTERVAL '7 days'
			  AND NOT EXISTS (
				SELECT 1 FROM expense_notification_log n
				WHERE n.employee_id = e.employee_id AND n.notification_type = 'UNSUBMITTED_EXPENSE_REMINDER'
				  AND n.created_at > NOW() - INTERVAL '24 hours'
			  )`)
		if err != nil {
			log.Printf("[expenses] reminderWorker expense reminder error: %v", err)
		}

		log.Println("[expenses] reminderWorker completed")
	}
}

func budgetWorker(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expenses] budgetWorker panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	log.Println("[expenses] budgetWorker started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[expenses] budgetWorker: monitoring budget thresholds...")

		rows, err := pool.Query(ctx, `
			SELECT id, company_id, name, total_amount, used_amount,
			       CASE
				 WHEN total_amount > 0 THEN (used_amount::numeric / total_amount::numeric) * 100
				 ELSE 0
			       END as usage_pct
			FROM expense_budgets
			WHERE is_active = TRUE
			  AND period_end >= CURRENT_DATE`)
		if err != nil {
			log.Printf("[expenses] budgetWorker query error: %v", err)
			continue
		}

		type budgetUsage struct {
			ID         string
			CompanyID  string
			Name       string
			Total      float64
			Used       float64
			UsagePct   float64
		}

		for rows.Next() {
			var b budgetUsage
			if err := rows.Scan(&b.ID, &b.CompanyID, &b.Name, &b.Total, &b.Used, &b.UsagePct); err != nil {
				log.Printf("[expenses] budgetWorker scan error: %v", err)
				continue
			}
			thresholds := []float64{50, 75, 90, 100}
			for _, t := range thresholds {
				if b.UsagePct >= t {
					log.Printf("[expenses] budgetWorker: budget %s (%s) at %.1f%% usage (threshold: %.0f%%)",
						b.Name, b.ID, b.UsagePct, t)
				}
			}
		}
		rows.Close()

		log.Println("[expenses] budgetWorker completed")
	}
}

func auditCleaner(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expenses] auditCleaner panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	log.Println("[expenses] auditCleaner started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[expenses] auditCleaner: cleaning old audit logs...")
		result, err := pool.Exec(ctx, `DELETE FROM expense_audit_logs WHERE created_at < NOW() - INTERVAL '90 days'`)
		if err != nil {
			log.Printf("[expenses] auditCleaner error: %v", err)
			continue
		}
		log.Printf("[expenses] auditCleaner removed %d records", result.RowsAffected())
	}
}

func duplicateCheckWorker(pool *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expenses] duplicateCheckWorker panicked: %v", r)
		}
	}()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	log.Println("[expenses] duplicateCheckWorker started")
	for range ticker.C {
		ctx := context.Background()
		log.Println("[expenses] duplicateCheckWorker: scanning for potential duplicates...")

		_, err := pool.Exec(ctx, `
			INSERT INTO expense_duplicate_checks (id, company_id, expense_id, duplicate_expense_id, match_reason, match_score, status, created_at)
			SELECT gen_random_uuid(), e1.company_id, e1.id, e2.id, 'MERCHANT_TAX_ID_AMOUNT', 0.95, 'PENDING', NOW()
			FROM expenses e1
			JOIN expenses e2 ON e1.company_id = e2.company_id
				AND e1.id != e2.id
				AND e1.merchant_tax_id IS NOT NULL
				AND e1.merchant_tax_id = e2.merchant_tax_id
				AND e1.base_amount = e2.base_amount
			WHERE e1.status = 'SUBMITTED'
			  AND e2.status IN ('SUBMITTED','APPROVED','PAID')
			  AND e1.created_at > NOW() - INTERVAL '24 hours'
			  AND NOT EXISTS (
				SELECT 1 FROM expense_duplicate_checks dc
				WHERE dc.expense_id = e1.id AND dc.duplicate_expense_id = e2.id
			  )`)
		if err != nil {
			log.Printf("[expenses] duplicateCheckWorker error: %v", err)
			continue
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO expense_duplicate_checks (id, company_id, expense_id, duplicate_expense_id, match_reason, match_score, status, created_at)
			SELECT gen_random_uuid(), e1.company_id, e1.id, e2.id, 'RECEIPT_NUMBER', 0.90, 'PENDING', NOW()
			FROM expenses e1
			JOIN expenses e2 ON e1.company_id = e2.company_id
				AND e1.id != e2.id
				AND e1.receipt_number IS NOT NULL
				AND e1.receipt_number = e2.receipt_number
			WHERE e1.status = 'SUBMITTED'
			  AND e2.status IN ('SUBMITTED','APPROVED','PAID')
			  AND e1.created_at > NOW() - INTERVAL '24 hours'
			  AND NOT EXISTS (
				SELECT 1 FROM expense_duplicate_checks dc
				WHERE dc.expense_id = e1.id AND dc.duplicate_expense_id = e2.id
			  )`)
		if err != nil {
			log.Printf("[expenses] duplicateCheckWorker error: %v", err)
			continue
		}

		log.Println("[expenses] duplicateCheckWorker completed")
	}
}
