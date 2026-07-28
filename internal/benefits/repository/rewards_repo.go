package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type RewardsRepo struct {
	pool *pgxpool.Pool
}

func NewRewardsRepo(pool *pgxpool.Pool) *RewardsRepo {
	return &RewardsRepo{pool: pool}
}

func (r *RewardsRepo) CreateItem(ctx context.Context, item *domain.TotalRewardsItem) error {
	q := `INSERT INTO total_rewards_items (id,company_id,name,category,description,amount_type,amount_value,
		amount_percentage,currency,frequency,display_order,is_monetary,include_in_statement,icon,color,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, item.ID, item.CompanyID, item.Name, item.Category, item.Description,
		item.AmountType, item.AmountValue, item.AmountPercentage, item.Currency, item.Frequency,
		item.DisplayOrder, item.IsMonetary, item.IncludeInStatement, item.Icon, item.Color, item.IsActive, item.CreatedBy)
	return repoErr("CreateItem", err)
}

func (r *RewardsRepo) GetItem(ctx context.Context, companyID, id uuid.UUID) (*domain.TotalRewardsItem, error) {
	q := `SELECT id,company_id,name,category,description,amount_type,amount_value,amount_percentage,currency,
		frequency,display_order,is_monetary,include_in_statement,icon,color,is_active,created_by,created_at,updated_at
		FROM total_rewards_items WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var item domain.TotalRewardsItem
	err := row.Scan(&item.ID, &item.CompanyID, &item.Name, &item.Category, &item.Description,
		&item.AmountType, &item.AmountValue, &item.AmountPercentage, &item.Currency, &item.Frequency,
		&item.DisplayOrder, &item.IsMonetary, &item.IncludeInStatement, &item.Icon, &item.Color, &item.IsActive, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetItem", err)
	}
	return &item, nil
}

func (r *RewardsRepo) ListItems(ctx context.Context, companyID uuid.UUID) ([]domain.TotalRewardsItem, error) {
	q := `SELECT id,company_id,name,category,description,amount_type,amount_value,amount_percentage,currency,
		frequency,display_order,is_monetary,include_in_statement,icon,color,is_active,created_by,created_at,updated_at
		FROM total_rewards_items WHERE company_id=$1 ORDER BY display_order,category,name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListItems", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.TotalRewardsItem, error) {
		var item domain.TotalRewardsItem
		err := row.Scan(&item.ID, &item.CompanyID, &item.Name, &item.Category, &item.Description,
			&item.AmountType, &item.AmountValue, &item.AmountPercentage, &item.Currency, &item.Frequency,
			&item.DisplayOrder, &item.IsMonetary, &item.IncludeInStatement, &item.Icon, &item.Color, &item.IsActive, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt)
		return item, err
	})
}

func (r *RewardsRepo) UpdateItem(ctx context.Context, item *domain.TotalRewardsItem) error {
	q := `UPDATE total_rewards_items SET name=$1,category=$2,description=$3,amount_type=$4,amount_value=$5,
		amount_percentage=$6,currency=$7,frequency=$8,display_order=$9,is_monetary=$10,
		include_in_statement=$11,icon=$12,color=$13,is_active=$14,updated_at=NOW()
		WHERE id=$15 AND company_id=$16`
	_, err := r.pool.Exec(ctx, q, item.Name, item.Category, item.Description,
		item.AmountType, item.AmountValue, item.AmountPercentage, item.Currency, item.Frequency,
		item.DisplayOrder, item.IsMonetary, item.IncludeInStatement, item.Icon, item.Color, item.IsActive,
		item.ID, item.CompanyID)
	return repoErr("UpdateItem", err)
}

func scanSnapshot(row pgx.CollectableRow) (domain.TotalRewardsSnapshot, error) {
	var s domain.TotalRewardsSnapshot
	err := row.Scan(&s.ID, &s.CompanyID, &s.EmployeeID, &s.SnapshotDate, &s.FiscalYear, &s.PeriodName,
		&s.BaseSalary, &s.VariablePay, &s.BonusesTotal, &s.IncentivesTotal, &s.BenefitsTotal,
		&s.EmployerContributions, &s.FlexibleSpending, &s.InsuranceValue, &s.DevelopmentValue,
		&s.WellnessValue, &s.RecognitionValue, &s.PerksValue, &s.TotalRewards, &s.Currency,
		&s.Items, &s.Metadata, &s.GeneratedBy, &s.GeneratedAt, &s.CreatedAt)
	return s, err
}

func (r *RewardsRepo) CreateSnapshot(ctx context.Context, s *domain.TotalRewardsSnapshot) error {
	q := `INSERT INTO total_rewards_snapshots (id,company_id,employee_id,snapshot_date,fiscal_year,period_name,
		base_salary,variable_pay,bonuses_total,incentives_total,benefits_total,employer_contributions,
		flexible_spending,insurance_value,development_value,wellness_value,recognition_value,perks_value,
		total_rewards,currency,items,metadata,generated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.EmployeeID, s.SnapshotDate, s.FiscalYear, s.PeriodName,
		s.BaseSalary, s.VariablePay, s.BonusesTotal, s.IncentivesTotal, s.BenefitsTotal,
		s.EmployerContributions, s.FlexibleSpending, s.InsuranceValue, s.DevelopmentValue,
		s.WellnessValue, s.RecognitionValue, s.PerksValue, s.TotalRewards, s.Currency,
		s.Items, s.Metadata, s.GeneratedBy)
	return repoErr("CreateSnapshot", err)
}

func (r *RewardsRepo) GetSnapshotByEmployeeAndYear(ctx context.Context, employeeID uuid.UUID, fiscalYear int) (*domain.TotalRewardsSnapshot, error) {
	q := `SELECT id,company_id,employee_id,snapshot_date,fiscal_year,period_name,base_salary,variable_pay,
		bonuses_total,incentives_total,benefits_total,employer_contributions,flexible_spending,insurance_value,
		development_value,wellness_value,recognition_value,perks_value,total_rewards,currency,items,metadata,
		generated_by,generated_at,created_at
		FROM total_rewards_snapshots WHERE employee_id=$1 AND fiscal_year=$2 ORDER BY snapshot_date DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, employeeID, fiscalYear)
	s, err := scanSnapshot(row)
	if err != nil {
		return nil, repoErr("GetSnapshotByEmployeeAndYear", err)
	}
	return &s, nil
}

func (r *RewardsRepo) GetLatestSnapshot(ctx context.Context, employeeID uuid.UUID) (*domain.TotalRewardsSnapshot, error) {
	q := `SELECT id,company_id,employee_id,snapshot_date,fiscal_year,period_name,base_salary,variable_pay,
		bonuses_total,incentives_total,benefits_total,employer_contributions,flexible_spending,insurance_value,
		development_value,wellness_value,recognition_value,perks_value,total_rewards,currency,items,metadata,
		generated_by,generated_at,created_at
		FROM total_rewards_snapshots WHERE employee_id=$1 ORDER BY snapshot_date DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, employeeID)
	s, err := scanSnapshot(row)
	if err != nil {
		return nil, repoErr("GetLatestSnapshot", err)
	}
	return &s, nil
}

func (r *RewardsRepo) ListSnapshots(ctx context.Context, companyID uuid.UUID, fiscalYear int) ([]domain.TotalRewardsSnapshot, error) {
	q := `SELECT id,company_id,employee_id,snapshot_date,fiscal_year,period_name,base_salary,variable_pay,
		bonuses_total,incentives_total,benefits_total,employer_contributions,flexible_spending,insurance_value,
		development_value,wellness_value,recognition_value,perks_value,total_rewards,currency,items,metadata,
		generated_by,generated_at,created_at
		FROM total_rewards_snapshots WHERE company_id=$1 AND fiscal_year=$2 ORDER BY employee_id,snapshot_date DESC`
	rows, err := r.pool.Query(ctx, q, companyID, fiscalYear)
	if err != nil {
		return nil, repoErr("ListSnapshots", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, scanSnapshot)
}

func (r *RewardsRepo) CreateNotification(ctx context.Context, n *domain.BenefitNotificationLog) error {
	q := `INSERT INTO benefit_notification_log (id,company_id,employee_id,notification_type,channel,title,body,metadata)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, n.ID, n.CompanyID, n.EmployeeID, n.NotificationType, n.Channel, n.Title, n.Body, n.Metadata)
	return repoErr("CreateNotification", err)
}

func (r *RewardsRepo) ListNotifications(ctx context.Context, employeeID, notificationType *uuid.UUID, limit, offset int) ([]domain.BenefitNotificationLog, error) {
	q := `SELECT id,company_id,employee_id,notification_type,channel,title,body,metadata,read_at,sent_at,created_at
		FROM benefit_notification_log WHERE 1=1`
	args := []any{}
	n := 1
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if notificationType != nil {
		q += fmt.Sprintf(" AND notification_type=$%d", n)
		args = append(args, *notificationType)
		n++
	}
	q += " ORDER BY sent_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListNotifications", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitNotificationLog, error) {
		var n domain.BenefitNotificationLog
		err := row.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.NotificationType, &n.Channel, &n.Title, &n.Body, &n.Metadata, &n.ReadAt, &n.SentAt, &n.CreatedAt)
		return n, err
	})
}

func (r *RewardsRepo) MarkRead(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE benefit_notification_log SET read_at=NOW() WHERE id=$1 AND read_at IS NULL`, id)
	return repoErr("MarkRead", err)
}

func (r *RewardsRepo) CreateReportDefinition(ctx context.Context, d *domain.BenefitReportDefinition) error {
	q := `INSERT INTO benefit_report_definitions (id,company_id,name,description,report_type,config,schedule_cron,recipients,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, d.ID, d.CompanyID, d.Name, d.Description, d.ReportType, d.Config, d.ScheduleCron, d.Recipients, d.IsActive, d.CreatedBy)
	return repoErr("CreateReportDefinition", err)
}

func (r *RewardsRepo) GetReportDefinition(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitReportDefinition, error) {
	q := `SELECT id,company_id,name,description,report_type,config,schedule_cron,recipients,is_active,created_by,created_at,updated_at
		FROM benefit_report_definitions WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var d domain.BenefitReportDefinition
	err := row.Scan(&d.ID, &d.CompanyID, &d.Name, &d.Description, &d.ReportType, &d.Config, &d.ScheduleCron, &d.Recipients, &d.IsActive, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetReportDefinition", err)
	}
	return &d, nil
}

func (r *RewardsRepo) ListReportDefinitions(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitReportDefinition, error) {
	q := `SELECT id,company_id,name,description,report_type,config,schedule_cron,recipients,is_active,created_by,created_at,updated_at
		FROM benefit_report_definitions WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListReportDefinitions", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitReportDefinition, error) {
		var d domain.BenefitReportDefinition
		err := row.Scan(&d.ID, &d.CompanyID, &d.Name, &d.Description, &d.ReportType, &d.Config, &d.ScheduleCron, &d.Recipients, &d.IsActive, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt)
		return d, err
	})
}

func (r *RewardsRepo) UpdateReportDefinition(ctx context.Context, d *domain.BenefitReportDefinition) error {
	q := `UPDATE benefit_report_definitions SET name=$1,description=$2,report_type=$3,config=$4,schedule_cron=$5,
		recipients=$6,is_active=$7,updated_at=NOW() WHERE id=$8 AND company_id=$9`
	_, err := r.pool.Exec(ctx, q, d.Name, d.Description, d.ReportType, d.Config, d.ScheduleCron, d.Recipients, d.IsActive, d.ID, d.CompanyID)
	return repoErr("UpdateReportDefinition", err)
}

func (r *RewardsRepo) CreateReportResult(ctx context.Context, res *domain.BenefitReportResult) error {
	q := `INSERT INTO benefit_report_results (id,definition_id,company_id,report_type,period_start,period_end,
		file_name,file_content,storage_path,file_size,format,status,error_message,generated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, res.ID, res.DefinitionID, res.CompanyID, res.ReportType, res.PeriodStart, res.PeriodEnd,
		res.FileName, res.FileContent, res.StoragePath, res.FileSize, res.Format, res.Status, res.ErrorMessage, res.GeneratedBy)
	return repoErr("CreateReportResult", err)
}

func (r *RewardsRepo) GetReportResult(ctx context.Context, id uuid.UUID) (*domain.BenefitReportResult, error) {
	q := `SELECT id,definition_id,company_id,report_type,period_start,period_end,file_name,file_content,storage_path,
		file_size,format,status,error_message,generated_by,generated_at,created_at
		FROM benefit_report_results WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var res domain.BenefitReportResult
	err := row.Scan(&res.ID, &res.DefinitionID, &res.CompanyID, &res.ReportType, &res.PeriodStart, &res.PeriodEnd,
		&res.FileName, &res.FileContent, &res.StoragePath, &res.FileSize, &res.Format, &res.Status, &res.ErrorMessage, &res.GeneratedBy, &res.GeneratedAt, &res.CreatedAt)
	if err != nil {
		return nil, repoErr("GetReportResult", err)
	}
	return &res, nil
}

func (r *RewardsRepo) ListReportResults(ctx context.Context, definitionID, companyID *uuid.UUID) ([]domain.BenefitReportResult, error) {
	q := `SELECT id,definition_id,company_id,report_type,period_start,period_end,file_name,file_content,storage_path,
		file_size,format,status,error_message,generated_by,generated_at,created_at
		FROM benefit_report_results WHERE 1=1`
	args := []any{}
	n := 1
	if definitionID != nil {
		q += fmt.Sprintf(" AND definition_id=$%d", n)
		args = append(args, *definitionID)
		n++
	}
	if companyID != nil {
		q += fmt.Sprintf(" AND company_id=$%d", n)
		args = append(args, *companyID)
		n++
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListReportResults", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitReportResult, error) {
		var res domain.BenefitReportResult
		err := row.Scan(&res.ID, &res.DefinitionID, &res.CompanyID, &res.ReportType, &res.PeriodStart, &res.PeriodEnd,
			&res.FileName, &res.FileContent, &res.StoragePath, &res.FileSize, &res.Format, &res.Status, &res.ErrorMessage, &res.GeneratedBy, &res.GeneratedAt, &res.CreatedAt)
		return res, err
	})
}
