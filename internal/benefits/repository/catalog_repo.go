package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func repoErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("benefits_repo.%s: %w", op, err)
}

type CatalogRepo struct {
	pool *pgxpool.Pool
}

func NewCatalogRepo(pool *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{pool: pool}
}

func (r *CatalogRepo) CreateCategory(ctx context.Context, c *domain.BenefitCategory) error {
	q := `INSERT INTO benefit_categories (id,company_id,name,description,icon,color,sort_order,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.Icon, c.Color, c.SortOrder, c.IsActive, c.CreatedBy)
	return repoErr("CreateCategory", err)
}

func (r *CatalogRepo) GetCategory(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitCategory, error) {
	q := `SELECT id,company_id,name,description,icon,color,sort_order,is_active,created_by,created_at,updated_at
		FROM benefit_categories WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.BenefitCategory
	err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Icon, &c.Color, &c.SortOrder, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCategory", err)
	}
	return &c, nil
}

func (r *CatalogRepo) ListCategories(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitCategory, error) {
	q := `SELECT id,company_id,name,description,icon,color,sort_order,is_active,created_by,created_at,updated_at
		FROM benefit_categories WHERE company_id=$1 ORDER BY sort_order,name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCategories", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitCategory, error) {
		var c domain.BenefitCategory
		err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.Icon, &c.Color, &c.SortOrder, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *CatalogRepo) UpdateCategory(ctx context.Context, c *domain.BenefitCategory) error {
	q := `UPDATE benefit_categories SET name=$1,description=$2,icon=$3,color=$4,sort_order=$5,is_active=$6,updated_at=NOW()
		WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.Icon, c.Color, c.SortOrder, c.IsActive, c.ID, c.CompanyID)
	return repoErr("UpdateCategory", err)
}

func (r *CatalogRepo) DeleteCategory(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_categories WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteCategory", err)
}

func (r *CatalogRepo) CreateType(ctx context.Context, t *domain.BenefitType) error {
	q := `INSERT INTO benefit_types (id,company_id,category_id,name,description,code,nature,tax_treatment,
		requires_approval,is_reimbursable,is_flexible,has_wallet,sort_order,is_active,config_schema,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CompanyID, t.CategoryID, t.Name, t.Description, t.Code, t.Nature, t.TaxTreatment,
		t.RequiresApproval, t.IsReimbursable, t.IsFlexible, t.HasWallet, t.SortOrder, t.IsActive, t.ConfigSchema, t.CreatedBy)
	return repoErr("CreateType", err)
}

func (r *CatalogRepo) GetType(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitType, error) {
	q := `SELECT id,company_id,category_id,name,description,code,nature,tax_treatment,
		requires_approval,is_reimbursable,is_flexible,has_wallet,sort_order,is_active,config_schema,created_by,created_at,updated_at
		FROM benefit_types WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var t domain.BenefitType
	err := row.Scan(&t.ID, &t.CompanyID, &t.CategoryID, &t.Name, &t.Description, &t.Code, &t.Nature, &t.TaxTreatment,
		&t.RequiresApproval, &t.IsReimbursable, &t.IsFlexible, &t.HasWallet, &t.SortOrder, &t.IsActive, &t.ConfigSchema, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetType", err)
	}
	return &t, nil
}

func (r *CatalogRepo) ListTypes(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitType, error) {
	q := `SELECT id,company_id,category_id,name,description,code,nature,tax_treatment,
		requires_approval,is_reimbursable,is_flexible,has_wallet,sort_order,is_active,config_schema,created_by,created_at,updated_at
		FROM benefit_types WHERE company_id=$1 ORDER BY sort_order,name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListTypes", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitType, error) {
		var t domain.BenefitType
		err := row.Scan(&t.ID, &t.CompanyID, &t.CategoryID, &t.Name, &t.Description, &t.Code, &t.Nature, &t.TaxTreatment,
			&t.RequiresApproval, &t.IsReimbursable, &t.IsFlexible, &t.HasWallet, &t.SortOrder, &t.IsActive, &t.ConfigSchema, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *CatalogRepo) UpdateType(ctx context.Context, t *domain.BenefitType) error {
	q := `UPDATE benefit_types SET category_id=$1,name=$2,description=$3,code=$4,nature=$5,tax_treatment=$6,
		requires_approval=$7,is_reimbursable=$8,is_flexible=$9,has_wallet=$10,sort_order=$11,is_active=$12,config_schema=$13,updated_at=NOW()
		WHERE id=$14 AND company_id=$15`
	_, err := r.pool.Exec(ctx, q, t.CategoryID, t.Name, t.Description, t.Code, t.Nature, t.TaxTreatment,
		t.RequiresApproval, t.IsReimbursable, t.IsFlexible, t.HasWallet, t.SortOrder, t.IsActive, t.ConfigSchema, t.ID, t.CompanyID)
	return repoErr("UpdateType", err)
}

func (r *CatalogRepo) DeleteType(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_types WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteType", err)
}

func (r *CatalogRepo) CreateProvider(ctx context.Context, p *domain.BenefitProvider) error {
	q := `INSERT INTO benefit_providers (id,company_id,name,legal_name,tax_id,provider_type,contact_name,
		contact_email,contact_phone,website,address,service_region,contract_start,contract_end,
		contract_file_path,billing_cycle,billing_currency,rating,notes,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Name, p.LegalName, p.TaxID, p.ProviderType,
		p.ContactName, p.ContactEmail, p.ContactPhone, p.Website, p.Address, p.ServiceRegion,
		p.ContractStart, p.ContractEnd, p.ContractFilePath, p.BillingCycle, p.BillingCurrency,
		p.Rating, p.Notes, p.IsActive, p.CreatedBy)
	return repoErr("CreateProvider", err)
}

func (r *CatalogRepo) GetProvider(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitProvider, error) {
	q := `SELECT id,company_id,name,legal_name,tax_id,provider_type,contact_name,contact_email,contact_phone,
		website,address,service_region,contract_start,contract_end,contract_file_path,billing_cycle,
		billing_currency,rating,notes,is_active,created_by,created_at,updated_at
		FROM benefit_providers WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p domain.BenefitProvider
	err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.LegalName, &p.TaxID, &p.ProviderType,
		&p.ContactName, &p.ContactEmail, &p.ContactPhone, &p.Website, &p.Address, &p.ServiceRegion,
		&p.ContractStart, &p.ContractEnd, &p.ContractFilePath, &p.BillingCycle,
		&p.BillingCurrency, &p.Rating, &p.Notes, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetProvider", err)
	}
	return &p, nil
}

func (r *CatalogRepo) ListProviders(ctx context.Context, companyID uuid.UUID) ([]domain.BenefitProvider, error) {
	q := `SELECT id,company_id,name,legal_name,tax_id,provider_type,contact_name,contact_email,contact_phone,
		website,address,service_region,contract_start,contract_end,contract_file_path,billing_cycle,
		billing_currency,rating,notes,is_active,created_by,created_at,updated_at
		FROM benefit_providers WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListProviders", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitProvider, error) {
		var p domain.BenefitProvider
		err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.LegalName, &p.TaxID, &p.ProviderType,
			&p.ContactName, &p.ContactEmail, &p.ContactPhone, &p.Website, &p.Address, &p.ServiceRegion,
			&p.ContractStart, &p.ContractEnd, &p.ContractFilePath, &p.BillingCycle,
			&p.BillingCurrency, &p.Rating, &p.Notes, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *CatalogRepo) UpdateProvider(ctx context.Context, p *domain.BenefitProvider) error {
	q := `UPDATE benefit_providers SET name=$1,legal_name=$2,tax_id=$3,provider_type=$4,contact_name=$5,
		contact_email=$6,contact_phone=$7,website=$8,address=$9,service_region=$10,contract_start=$11,
		contract_end=$12,contract_file_path=$13,billing_cycle=$14,billing_currency=$15,rating=$16,notes=$17,
		is_active=$18,updated_at=NOW() WHERE id=$19 AND company_id=$20`
	_, err := r.pool.Exec(ctx, q, p.Name, p.LegalName, p.TaxID, p.ProviderType,
		p.ContactName, p.ContactEmail, p.ContactPhone, p.Website, p.Address, p.ServiceRegion,
		p.ContractStart, p.ContractEnd, p.ContractFilePath, p.BillingCycle, p.BillingCurrency,
		p.Rating, p.Notes, p.IsActive, p.ID, p.CompanyID)
	return repoErr("UpdateProvider", err)
}

func (r *CatalogRepo) DeleteProvider(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_providers WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteProvider", err)
}
