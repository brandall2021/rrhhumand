package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type BenefitRepo struct {
	pool *pgxpool.Pool
}

func NewBenefitRepo(pool *pgxpool.Pool) *BenefitRepo {
	return &BenefitRepo{pool: pool}
}

func (r *BenefitRepo) Create(ctx context.Context, b *domain.Benefit) error {
	q := `INSERT INTO benefits (id,company_id,type_id,provider_id,code,name,description,short_description,
		coverage_details,eligibility_summary,logo_url,banner_url,website_url,terms_url,documentation_url,
		provider_reference,start_date,end_date,max_beneficiaries,current_beneficiaries,waiting_period_days,
		minimum_service_months,deductible_amount,deductible_period,copay_percentage,max_coverage_amount,
		max_coverage_period,auto_enroll,enrollment_deadline_days,requires_evidence,evidence_description,
		status,visibility,sort_order,metadata,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.TypeID, b.ProviderID, b.Code, b.Name, b.Description, b.ShortDescription,
		b.CoverageDetails, b.EligibilitySummary, b.LogoURL, b.BannerURL, b.WebsiteURL, b.TermsURL, b.DocumentationURL,
		b.ProviderReference, b.StartDate, b.EndDate, b.MaxBeneficiaries, b.CurrentBeneficiaries, b.WaitingPeriodDays,
		b.MinimumServiceMonths, b.DeductibleAmount, b.DeductiblePeriod, b.CopayPercentage, b.MaxCoverageAmount,
		b.MaxCoveragePeriod, b.AutoEnroll, b.EnrollmentDeadlineDays, b.RequiresEvidence, b.EvidenceDescription,
		b.Status, b.Visibility, b.SortOrder, b.Metadata, b.CreatedBy)
	return repoErr("Create", err)
}

func scanBenefit(row pgx.Row) (domain.Benefit, error) {
	var b domain.Benefit
	err := row.Scan(&b.ID, &b.CompanyID, &b.TypeID, &b.ProviderID, &b.Code, &b.Name, &b.Description, &b.ShortDescription,
		&b.CoverageDetails, &b.EligibilitySummary, &b.LogoURL, &b.BannerURL, &b.WebsiteURL, &b.TermsURL, &b.DocumentationURL,
		&b.ProviderReference, &b.StartDate, &b.EndDate, &b.MaxBeneficiaries, &b.CurrentBeneficiaries, &b.WaitingPeriodDays,
		&b.MinimumServiceMonths, &b.DeductibleAmount, &b.DeductiblePeriod, &b.CopayPercentage, &b.MaxCoverageAmount,
		&b.MaxCoveragePeriod, &b.AutoEnroll, &b.EnrollmentDeadlineDays, &b.RequiresEvidence, &b.EvidenceDescription,
		&b.Status, &b.Visibility, &b.SortOrder, &b.Metadata, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *BenefitRepo) Get(ctx context.Context, companyID, id uuid.UUID) (*domain.Benefit, error) {
	q := `SELECT id,company_id,type_id,provider_id,code,name,description,short_description,coverage_details,
		eligibility_summary,logo_url,banner_url,website_url,terms_url,documentation_url,provider_reference,
		start_date,end_date,max_beneficiaries,current_beneficiaries,waiting_period_days,minimum_service_months,
		deductible_amount,deductible_period,copay_percentage,max_coverage_amount,max_coverage_period,auto_enroll,
		enrollment_deadline_days,requires_evidence,evidence_description,status,visibility,sort_order,metadata,
		created_by,created_at,updated_at
		FROM benefits WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	b, err := scanBenefit(row)
	if err != nil {
		return nil, repoErr("Get", err)
	}
	return &b, nil
}

func (r *BenefitRepo) List(ctx context.Context, companyID uuid.UUID, status, typeID, visibility *string, limit, offset int) ([]domain.Benefit, error) {
	q := `SELECT id,company_id,type_id,provider_id,code,name,description,short_description,coverage_details,
		eligibility_summary,logo_url,banner_url,website_url,terms_url,documentation_url,provider_reference,
		start_date,end_date,max_beneficiaries,current_beneficiaries,waiting_period_days,minimum_service_months,
		deductible_amount,deductible_period,copay_percentage,max_coverage_amount,max_coverage_period,auto_enroll,
		enrollment_deadline_days,requires_evidence,evidence_description,status,visibility,sort_order,metadata,
		created_by,created_at,updated_at
		FROM benefits WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	if typeID != nil {
		q += fmt.Sprintf(" AND type_id=$%d", n)
		args = append(args, *typeID)
		n++
	}
	if visibility != nil {
		q += fmt.Sprintf(" AND visibility=$%d", n)
		args = append(args, *visibility)
		n++
	}
	q += " ORDER BY sort_order,name"
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
		return nil, repoErr("List", err)
	}
	defer rows.Close()
	var result []domain.Benefit
	for rows.Next() {
		var b domain.Benefit
		err := rows.Scan(&b.ID, &b.CompanyID, &b.TypeID, &b.ProviderID, &b.Code, &b.Name, &b.Description, &b.ShortDescription,
			&b.CoverageDetails, &b.EligibilitySummary, &b.LogoURL, &b.BannerURL, &b.WebsiteURL, &b.TermsURL, &b.DocumentationURL,
			&b.ProviderReference, &b.StartDate, &b.EndDate, &b.MaxBeneficiaries, &b.CurrentBeneficiaries, &b.WaitingPeriodDays,
			&b.MinimumServiceMonths, &b.DeductibleAmount, &b.DeductiblePeriod, &b.CopayPercentage, &b.MaxCoverageAmount,
			&b.MaxCoveragePeriod, &b.AutoEnroll, &b.EnrollmentDeadlineDays, &b.RequiresEvidence, &b.EvidenceDescription,
			&b.Status, &b.Visibility, &b.SortOrder, &b.Metadata, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, repoErr("List", err)
		}
		result = append(result, b)
	}
	return result, nil
}

func (r *BenefitRepo) Update(ctx context.Context, b *domain.Benefit) error {
	q := `UPDATE benefits SET type_id=$1,provider_id=$2,code=$3,name=$4,description=$5,short_description=$6,
		coverage_details=$7,eligibility_summary=$8,logo_url=$9,banner_url=$10,website_url=$11,terms_url=$12,
		documentation_url=$13,provider_reference=$14,start_date=$15,end_date=$16,max_beneficiaries=$17,
		waiting_period_days=$18,minimum_service_months=$19,deductible_amount=$20,deductible_period=$21,
		copay_percentage=$22,max_coverage_amount=$23,max_coverage_period=$24,auto_enroll=$25,
		enrollment_deadline_days=$26,requires_evidence=$27,evidence_description=$28,status=$29,visibility=$30,
		sort_order=$31,metadata=$32,updated_at=NOW() WHERE id=$33 AND company_id=$34`
	_, err := r.pool.Exec(ctx, q, b.TypeID, b.ProviderID, b.Code, b.Name, b.Description, b.ShortDescription,
		b.CoverageDetails, b.EligibilitySummary, b.LogoURL, b.BannerURL, b.WebsiteURL, b.TermsURL,
		b.DocumentationURL, b.ProviderReference, b.StartDate, b.EndDate, b.MaxBeneficiaries,
		b.WaitingPeriodDays, b.MinimumServiceMonths, b.DeductibleAmount, b.DeductiblePeriod,
		b.CopayPercentage, b.MaxCoverageAmount, b.MaxCoveragePeriod, b.AutoEnroll,
		b.EnrollmentDeadlineDays, b.RequiresEvidence, b.EvidenceDescription, b.Status, b.Visibility,
		b.SortOrder, b.Metadata, b.ID, b.CompanyID)
	return repoErr("Update", err)
}

func (r *BenefitRepo) Delete(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefits WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("Delete", err)
}

func (r *BenefitRepo) IncrementBeneficiaries(ctx context.Context, benefitID uuid.UUID, delta int) error {
	q := `UPDATE benefits SET current_beneficiaries=current_beneficiaries+$1,updated_at=NOW() WHERE id=$2`
	_, err := r.pool.Exec(ctx, q, delta, benefitID)
	return repoErr("IncrementBeneficiaries", err)
}

func (r *BenefitRepo) SearchBenefits(ctx context.Context, companyID uuid.UUID, query string) ([]domain.Benefit, error) {
	q := `SELECT id,company_id,type_id,provider_id,code,name,description,short_description,coverage_details,
		eligibility_summary,logo_url,banner_url,website_url,terms_url,documentation_url,provider_reference,
		start_date,end_date,max_beneficiaries,current_beneficiaries,waiting_period_days,minimum_service_months,
		deductible_amount,deductible_period,copay_percentage,max_coverage_amount,max_coverage_period,auto_enroll,
		enrollment_deadline_days,requires_evidence,evidence_description,status,visibility,sort_order,metadata,
		created_by,created_at,updated_at
		FROM benefits WHERE company_id=$1 AND (name ILIKE $2 OR code ILIKE $2) ORDER BY sort_order,name`
	pattern := "%" + query + "%"
	rows, err := r.pool.Query(ctx, q, companyID, pattern)
	if err != nil {
		return nil, repoErr("SearchBenefits", err)
	}
	defer rows.Close()
	var result []domain.Benefit
	for rows.Next() {
		var b domain.Benefit
		err := rows.Scan(&b.ID, &b.CompanyID, &b.TypeID, &b.ProviderID, &b.Code, &b.Name, &b.Description, &b.ShortDescription,
			&b.CoverageDetails, &b.EligibilitySummary, &b.LogoURL, &b.BannerURL, &b.WebsiteURL, &b.TermsURL, &b.DocumentationURL,
			&b.ProviderReference, &b.StartDate, &b.EndDate, &b.MaxBeneficiaries, &b.CurrentBeneficiaries, &b.WaitingPeriodDays,
			&b.MinimumServiceMonths, &b.DeductibleAmount, &b.DeductiblePeriod, &b.CopayPercentage, &b.MaxCoverageAmount,
			&b.MaxCoveragePeriod, &b.AutoEnroll, &b.EnrollmentDeadlineDays, &b.RequiresEvidence, &b.EvidenceDescription,
			&b.Status, &b.Visibility, &b.SortOrder, &b.Metadata, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, repoErr("SearchBenefits", err)
		}
		result = append(result, b)
	}
	return result, nil
}
