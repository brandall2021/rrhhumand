package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardStats struct {
	OpenRequisitions   int `json:"open_requisitions"`
	TotalCandidates    int `json:"total_candidates"`
	ApplicationsWeek   int `json:"applications_this_week"`
	PendingOffers      int `json:"pending_offers"`
	HiresThisMonth     int `json:"hires_this_month"`
	TotalInterviews    int `json:"total_interviews"`
	AvgTimeToHire      float64 `json:"avg_time_to_hire"`
}

type FunnelStage struct {
	Stage string `json:"stage"`
	Count int    `json:"count"`
}

type FunnelData struct {
	Stages []FunnelStage `json:"stages"`
}

type TimeToHire struct {
	AvgDays float64 `json:"avg_days"`
	Median  float64 `json:"median_days"`
	MinDays int     `json:"min_days"`
	MaxDays int     `json:"max_days"`
}

type DashboardCache struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	CacheKey  string    `json:"cache_key"`
	Data      string    `json:"cache_data"`
	CachedAt  time.Time `json:"cached_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type DashboardRepo struct {
	pool *pgxpool.Pool
}

func NewDashboardRepo(pool *pgxpool.Pool) *DashboardRepo {
	return &DashboardRepo{pool: pool}
}

func (r *DashboardRepo) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	ds := &DashboardStats{}

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM job_requisitions WHERE company_id=$1 AND status IN ('OPEN','APPROVED')`,
		companyID).Scan(&ds.OpenRequisitions)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM candidates WHERE company_id=$1`,
		companyID).Scan(&ds.TotalCandidates)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM applications WHERE company_id=$1 AND applied_at >= NOW() - INTERVAL '7 days'`,
		companyID).Scan(&ds.ApplicationsWeek)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM offers WHERE company_id=$1 AND status IN ('DRAFT','PENDING_APPROVAL','SENT')`,
		companyID).Scan(&ds.PendingOffers)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM applications WHERE company_id=$1 AND status='HIRED' AND hired_at >= date_trunc('month', NOW())`,
		companyID).Scan(&ds.HiresThisMonth)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM interviews WHERE company_id=$1`,
		companyID).Scan(&ds.TotalInterviews)

	r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (hired_at - applied_at))/86400),0) FROM applications WHERE company_id=$1 AND status='HIRED' AND hired_at IS NOT NULL`,
		companyID).Scan(&ds.AvgTimeToHire)

	return ds, nil
}

func (r *DashboardRepo) GetFunnelData(ctx context.Context, companyID string) (*FunnelData, error) {
	fd := &FunnelData{}

	rows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM applications WHERE company_id=$1 GROUP BY status ORDER BY COUNT(*) DESC`,
		companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fs FunnelStage
		rows.Scan(&fs.Stage, &fs.Count)
		fd.Stages = append(fd.Stages, fs)
	}
	return fd, nil
}

func (r *DashboardRepo) GetTimeToHire(ctx context.Context, companyID string) (*TimeToHire, error) {
	tth := &TimeToHire{}

	r.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (hired_at - applied_at))/86400),0),
		        COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (hired_at - applied_at))/86400),0),
		        COALESCE(MIN(EXTRACT(EPOCH FROM (hired_at - applied_at))/86400)::INT,0),
		        COALESCE(MAX(EXTRACT(EPOCH FROM (hired_at - applied_at))/86400)::INT,0)
		 FROM applications WHERE company_id=$1 AND status='HIRED' AND hired_at IS NOT NULL`,
		companyID).Scan(&tth.AvgDays, &tth.Median, &tth.MinDays, &tth.MaxDays)

	return tth, nil
}

func (r *DashboardRepo) CacheDashboard(ctx context.Context, companyID, cacheKey, data string, expiresAt *time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO recruitment_dashboard_cache (company_id, cache_key, cache_data, expires_at)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (company_id, cache_key) DO UPDATE SET cache_data=$3, cached_at=NOW(), expires_at=$4`,
		companyID, cacheKey, data, expiresAt)
	return err
}

func (r *DashboardRepo) GetCachedDashboard(ctx context.Context, companyID, cacheKey string) (*DashboardCache, error) {
	dc := &DashboardCache{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, cache_key, cache_data, cached_at, expires_at
		 FROM recruitment_dashboard_cache WHERE company_id=$1 AND cache_key=$2
		 AND (expires_at IS NULL OR expires_at > NOW())`,
		companyID, cacheKey,
	).Scan(&dc.ID, &dc.CompanyID, &dc.CacheKey, &dc.Data, &dc.CachedAt, &dc.ExpiresAt)
	return dc, err
}
