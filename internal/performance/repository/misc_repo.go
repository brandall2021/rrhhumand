package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/performance/domain"
)

type evidenceRepo struct {
	pool *pgxpool.Pool
}

func (r *evidenceRepo) Create(ctx context.Context, e *domain.PerformanceEvidence) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_evidence (company_id, evaluation_id, objective_id, feedback_id, title, description, evidence_type, storage_key, file_name, mime_type, size_bytes, url, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id`,
		e.CompanyID, e.EvaluationID, e.ObjectiveID, e.FeedbackID, e.Title, e.Description, e.EvidenceType,
		e.StorageKey, e.FileName, e.MimeType, e.SizeBytes, e.URL, e.CreatedBy, now,
	).Scan(&e.ID)
}

func (r *evidenceRepo) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceEvidence, error) {
	e := &domain.PerformanceEvidence{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, evaluation_id, objective_id, feedback_id, title, description, evidence_type,
		 storage_key, file_name, mime_type, size_bytes, url, created_by, created_at
		 FROM performance_evidence WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&e.ID, &e.CompanyID, &e.EvaluationID, &e.ObjectiveID, &e.FeedbackID, &e.Title, &e.Description, &e.EvidenceType,
		&e.StorageKey, &e.FileName, &e.MimeType, &e.SizeBytes, &e.URL, &e.CreatedBy, &e.CreatedAt)
	return e, err
}

func (r *evidenceRepo) ListByEvaluation(ctx context.Context, evaluationID string) ([]domain.PerformanceEvidence, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, evaluation_id, objective_id, feedback_id, title, description, evidence_type,
		 storage_key, file_name, mime_type, size_bytes, url, created_by, created_at
		 FROM performance_evidence WHERE evaluation_id=$1 ORDER BY created_at`, evaluationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evidence []domain.PerformanceEvidence
	for rows.Next() {
		var e domain.PerformanceEvidence
		rows.Scan(&e.ID, &e.CompanyID, &e.EvaluationID, &e.ObjectiveID, &e.FeedbackID, &e.Title, &e.Description, &e.EvidenceType,
			&e.StorageKey, &e.FileName, &e.MimeType, &e.SizeBytes, &e.URL, &e.CreatedBy, &e.CreatedAt)
		evidence = append(evidence, e)
	}
	return evidence, nil
}

func (r *evidenceRepo) ListByObjective(ctx context.Context, objectiveID string) ([]domain.PerformanceEvidence, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, evaluation_id, objective_id, feedback_id, title, description, evidence_type,
		 storage_key, file_name, mime_type, size_bytes, url, created_by, created_at
		 FROM performance_evidence WHERE objective_id=$1 ORDER BY created_at`, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evidence []domain.PerformanceEvidence
	for rows.Next() {
		var e domain.PerformanceEvidence
		rows.Scan(&e.ID, &e.CompanyID, &e.EvaluationID, &e.ObjectiveID, &e.FeedbackID, &e.Title, &e.Description, &e.EvidenceType,
			&e.StorageKey, &e.FileName, &e.MimeType, &e.SizeBytes, &e.URL, &e.CreatedBy, &e.CreatedAt)
		evidence = append(evidence, e)
	}
	return evidence, nil
}

func (r *evidenceRepo) ListByFeedback(ctx context.Context, feedbackID string) ([]domain.PerformanceEvidence, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, evaluation_id, objective_id, feedback_id, title, description, evidence_type,
		 storage_key, file_name, mime_type, size_bytes, url, created_by, created_at
		 FROM performance_evidence WHERE feedback_id=$1 ORDER BY created_at`, feedbackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evidence []domain.PerformanceEvidence
	for rows.Next() {
		var e domain.PerformanceEvidence
		rows.Scan(&e.ID, &e.CompanyID, &e.EvaluationID, &e.ObjectiveID, &e.FeedbackID, &e.Title, &e.Description, &e.EvidenceType,
			&e.StorageKey, &e.FileName, &e.MimeType, &e.SizeBytes, &e.URL, &e.CreatedBy, &e.CreatedAt)
		evidence = append(evidence, e)
	}
	return evidence, nil
}

func (r *evidenceRepo) Delete(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM performance_evidence WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

// ResultRepository

type resultRepo struct {
	pool *pgxpool.Pool
}

func (r *resultRepo) Upsert(ctx context.Context, res *domain.PerformanceResult) error {
	now := time.Now()
	err := r.pool.QueryRow(ctx,
		`INSERT INTO performance_results (company_id, cycle_id, employee_id, objective_score, competency_score,
		 self_score, manager_score, peer_score, hr_score, final_score, final_rating, final_rating_label, strengths, improvement_areas, summary, calculated_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		 ON CONFLICT (cycle_id, employee_id) DO UPDATE SET
		 objective_score=$4, competency_score=$5, self_score=$6, manager_score=$7, peer_score=$8, hr_score=$9,
		 final_score=$10, final_rating=$11, final_rating_label=$12, strengths=$13, improvement_areas=$14, summary=$15,
		 calculated_at=NOW(), updated_at=NOW()
		 RETURNING id`,
		res.CompanyID, res.CycleID, res.EmployeeID, res.ObjectiveScore, res.CompetencyScore,
		res.SelfScore, res.ManagerScore, res.PeerScore, res.HRScore, res.FinalScore, res.FinalRating, res.FinalRatingLabel,
		res.Strengths, res.ImprovementAreas, res.Summary, now, now, now,
	).Scan(&res.ID)
	return err
}

func (r *resultRepo) GetByCycleEmployee(ctx context.Context, companyID, cycleID, employeeID string) (*domain.PerformanceResult, error) {
	res := &domain.PerformanceResult{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cycle_id, employee_id, objective_score, competency_score,
		 self_score, manager_score, peer_score, hr_score, final_score, final_rating, final_rating_label,
		 strengths, improvement_areas, summary, calculated_at, created_at, updated_at
		 FROM performance_results WHERE company_id=$1 AND cycle_id=$2 AND employee_id=$3`,
		companyID, cycleID, employeeID,
	).Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.EmployeeID, &res.ObjectiveScore, &res.CompetencyScore,
		&res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore, &res.FinalScore, &res.FinalRating, &res.FinalRatingLabel,
		&res.Strengths, &res.ImprovementAreas, &res.Summary, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
	return res, err
}

func (r *resultRepo) ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, objective_score, competency_score,
		 self_score, manager_score, peer_score, hr_score, final_score, final_rating, final_rating_label,
		 strengths, improvement_areas, summary, calculated_at, created_at, updated_at
		 FROM performance_results WHERE company_id=$1 AND cycle_id=$2 ORDER BY final_score DESC NULLS LAST`,
		companyID, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.PerformanceResult
	for rows.Next() {
		var res domain.PerformanceResult
		rows.Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.EmployeeID, &res.ObjectiveScore, &res.CompetencyScore,
			&res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore, &res.FinalScore, &res.FinalRating, &res.FinalRatingLabel,
			&res.Strengths, &res.ImprovementAreas, &res.Summary, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
		results = append(results, res)
	}
	return results, nil
}

func (r *resultRepo) ListByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, cycle_id, employee_id, objective_score, competency_score,
		 self_score, manager_score, peer_score, hr_score, final_score, final_rating, final_rating_label,
		 strengths, improvement_areas, summary, calculated_at, created_at, updated_at
		 FROM performance_results WHERE company_id=$1 AND employee_id=$2 ORDER BY calculated_at DESC`,
		companyID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.PerformanceResult
	for rows.Next() {
		var res domain.PerformanceResult
		rows.Scan(&res.ID, &res.CompanyID, &res.CycleID, &res.EmployeeID, &res.ObjectiveScore, &res.CompetencyScore,
			&res.SelfScore, &res.ManagerScore, &res.PeerScore, &res.HRScore, &res.FinalScore, &res.FinalRating, &res.FinalRatingLabel,
			&res.Strengths, &res.ImprovementAreas, &res.Summary, &res.CalculatedAt, &res.CreatedAt, &res.UpdatedAt)
		results = append(results, res)
	}
	return results, nil
}

func (r *resultRepo) Delete(ctx context.Context, companyID, cycleID, employeeID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM performance_results WHERE company_id=$1 AND cycle_id=$2 AND employee_id=$3`,
		companyID, cycleID, employeeID)
	return err
}

// DashboardRepository

type dashboardRepo struct {
	pool *pgxpool.Pool
}

func (r *dashboardRepo) GetDashboard(ctx context.Context, companyID string) (*domain.PerformanceDashboard, error) {
	dash := &domain.PerformanceDashboard{}

	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_cycles WHERE company_id=$1`, companyID).Scan(&dash.TotalCycles)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_cycles WHERE company_id=$1 AND status IN ('OPEN','IN_PROGRESS','REVIEW','CALIBRATION')`, companyID).Scan(&dash.ActiveCycles)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1`, companyID).Scan(&dash.TotalEvaluations)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1 AND status='SUBMITTED'`, companyID).Scan(&dash.CompletedEvaluations)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_evaluations WHERE company_id=$1 AND status IN ('DRAFT','PENDING')`, companyID).Scan(&dash.PendingEvaluations)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(final_score),0) FROM performance_results WHERE company_id=$1`, companyID).Scan(&dash.AverageScore)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_objectives WHERE company_id=$1`, companyID).Scan(&dash.TotalObjectives)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_objectives WHERE company_id=$1 AND status='COMPLETED'`, companyID).Scan(&dash.CompletedObjectives)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_feedback WHERE company_id=$1`, companyID).Scan(&dash.TotalFeedback)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_improvement_plans WHERE company_id=$1`, companyID).Scan(&dash.TotalImprovementPlans)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM performance_development_plans WHERE company_id=$1`, companyID).Scan(&dash.TotalDevelopmentPlans)

	rows, err := r.pool.Query(ctx,
		`SELECT COALESCE(final_rating,'UNRATED'), COUNT(*) FROM performance_results WHERE company_id=$1 GROUP BY final_rating ORDER BY COUNT(*) DESC`,
		companyID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rc domain.RatingCount
			rows.Scan(&rc.Rating, &rc.Count)
			dash.RatingDistribution = append(dash.RatingDistribution, rc)
		}
	}

	return dash, nil
}

// AuditLogRepository

type auditLogRepo struct {
	pool *pgxpool.Pool
}

func (r *auditLogRepo) Create(ctx context.Context, l *domain.PerformanceAuditLog) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO performance_audit_log (company_id, user_id, entity_type, entity_id, action, old_values, new_values, ip_address, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		l.CompanyID, l.UserID, l.EntityType, l.EntityID, l.Action, l.OldValues, l.NewValues, l.IPAddress, now)
	return err
}

func (r *auditLogRepo) ListByEntity(ctx context.Context, companyID, entityType, entityID string) ([]domain.PerformanceAuditLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, user_id, entity_type, entity_id, action, old_values, new_values, ip_address, created_at
		 FROM performance_audit_log WHERE company_id=$1 AND entity_type=$2 AND entity_id=$3 ORDER BY created_at DESC`,
		companyID, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.PerformanceAuditLog
	for rows.Next() {
		var l domain.PerformanceAuditLog
		rows.Scan(&l.ID, &l.CompanyID, &l.UserID, &l.EntityType, &l.EntityID, &l.Action, &l.OldValues, &l.NewValues, &l.IPAddress, &l.CreatedAt)
		logs = append(logs, l)
	}
	return logs, nil
}

// OutboxRepository

type outboxRepo struct {
	pool *pgxpool.Pool
}

func (r *outboxRepo) Create(ctx context.Context, e *domain.OutboxEvent) error {
	now := time.Now()
	return r.pool.QueryRow(ctx,
		`INSERT INTO performance_outbox (company_id, event_type, aggregate_type, aggregate_id, payload, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		e.CompanyID, e.EventType, e.AggregateType, e.AggregateID, e.Payload, e.Status, now,
	).Scan(&e.ID)
}

func (r *outboxRepo) ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, event_type, aggregate_type, aggregate_id, payload, status, retry_count, last_error, created_at, processed_at
		 FROM performance_outbox WHERE status='PENDING' ORDER BY created_at LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var e domain.OutboxEvent
		rows.Scan(&e.ID, &e.CompanyID, &e.EventType, &e.AggregateType, &e.AggregateID, &e.Payload, &e.Status, &e.RetryCount, &e.LastError, &e.CreatedAt, &e.ProcessedAt)
		events = append(events, e)
	}
	return events, nil
}

func (r *outboxRepo) MarkProcessed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE performance_outbox SET status='PROCESSED', processed_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *outboxRepo) MarkFailed(ctx context.Context, id string, errMsg error) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE performance_outbox SET status='FAILED', retry_count=retry_count+1, last_error=$2 WHERE id=$1`,
		id, errMsg.Error())
	return err
}

func NewEvidenceRepository(pool *pgxpool.Pool) EvidenceRepository {
	return &evidenceRepo{pool: pool}
}

func NewResultRepository(pool *pgxpool.Pool) ResultRepository {
	return &resultRepo{pool: pool}
}

func NewDashboardRepository(pool *pgxpool.Pool) DashboardRepository {
	return &dashboardRepo{pool: pool}
}

func NewAuditLogRepository(pool *pgxpool.Pool) AuditLogRepository {
	return &auditLogRepo{pool: pool}
}

func NewOutboxRepository(pool *pgxpool.Pool) OutboxRepository {
	return &outboxRepo{pool: pool}
}
