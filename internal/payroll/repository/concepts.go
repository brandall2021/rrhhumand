package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) CreateConcept(ctx context.Context, c *domain.PayrollConcept) error {
	q := `INSERT INTO payroll_concepts (id,company_id,code,name,description,concept_type,taxability,
		calculation_type,base_concept_id,active,effective_from,effective_to,sort_order,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Code, c.Name, c.Description, c.ConceptType, c.Taxability,
		c.CalculationType, c.BaseConceptID, c.Active, c.EffectiveFrom, c.EffectiveTo, c.SortOrder, c.CreatedBy)
	return repoErr("CreateConcept", err)
}

func (r *Repository) UpdateConcept(ctx context.Context, c *domain.PayrollConcept) error {
	q := `UPDATE payroll_concepts SET name=$1,description=$2,concept_type=$3,taxability=$4,
		calculation_type=$5,base_concept_id=$6,active=$7,sort_order=$8,updated_at=NOW() WHERE id=$9 AND company_id=$10`
	_, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.ConceptType, c.Taxability,
		c.CalculationType, c.BaseConceptID, c.Active, c.SortOrder, c.ID, c.CompanyID)
	return repoErr("UpdateConcept", err)
}

func (r *Repository) GetConcept(ctx context.Context, companyID, id string) (*domain.PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var c domain.PayrollConcept
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
		&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetConcept", err)
	}
	return &c, nil
}

func (r *Repository) GetConceptByCode(ctx context.Context, companyID, code string) (*domain.PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE company_id=$1 AND code=$2`
	row := r.pool.QueryRow(ctx, q, companyID, code)
	var c domain.PayrollConcept
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
		&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetConceptByCode", err)
	}
	return &c, nil
}

func (r *Repository) ListConcepts(ctx context.Context, companyID string, conceptType, taxability *string, active *bool) ([]domain.PayrollConcept, error) {
	q := `SELECT id,company_id,code,name,description,concept_type,taxability,calculation_type,
		base_concept_id,active,effective_from,effective_to,sort_order,created_by,created_at,updated_at
		FROM payroll_concepts WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if conceptType != nil {
		q += fmt.Sprintf(" AND concept_type=$%d", n)
		args = append(args, *conceptType)
		n++
	}
	if taxability != nil {
		q += fmt.Sprintf(" AND taxability=$%d", n)
		args = append(args, *taxability)
		n++
	}
	if active != nil {
		q += fmt.Sprintf(" AND active=$%d", n)
		args = append(args, *active)
		n++
	}
	q += " ORDER BY sort_order, code"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListConcepts", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollConcept, error) {
		var c domain.PayrollConcept
		err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ConceptType, &c.Taxability, &c.CalculationType,
			&c.BaseConceptID, &c.Active, &c.EffectiveFrom, &c.EffectiveTo, &c.SortOrder, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		return c, err
	})
}

func (r *Repository) CreateRule(ctx context.Context, rule *domain.PayrollRule) error {
	params, _ := json.Marshal(rule.Parameters)
	q := `INSERT INTO payroll_rules (id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,
		rule_type,formula,parameters,priority,effective_from,effective_to,version,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, rule.ID, rule.CompanyID, rule.Country, rule.Jurisdiction, rule.AgreementID,
		rule.CategoryID, rule.ConceptID, rule.RuleType, rule.Formula, params, rule.Priority,
		rule.EffectiveFrom, rule.EffectiveTo, rule.Version, rule.CreatedBy)
	return repoErr("CreateRule", err)
}

func (r *Repository) UpdateRule(ctx context.Context, rule *domain.PayrollRule) error {
	params, _ := json.Marshal(rule.Parameters)
	q := `UPDATE payroll_rules SET rule_type=$1,formula=$2,parameters=$3,priority=$4,active=$5,effective_to=$6,version=version+1,updated_at=NOW()
		WHERE id=$7 AND company_id=$8`
	_, err := r.pool.Exec(ctx, q, rule.RuleType, rule.Formula, params, rule.Priority, rule.Active, rule.EffectiveTo, rule.ID, rule.CompanyID)
	return repoErr("UpdateRule", err)
}

func (r *Repository) GetRule(ctx context.Context, companyID, id string) (*domain.PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at FROM payroll_rules WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var rule domain.PayrollRule
	var params []byte
	err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
		&rule.ConceptID, &rule.RuleType, &rule.Formula, &params, &rule.Priority, &rule.EffectiveFrom,
		&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetRule", err)
	}
	if params != nil {
		json.Unmarshal(params, &rule.Parameters)
	}
	return &rule, nil
}

func (r *Repository) ListRules(ctx context.Context, companyID string) ([]domain.PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at
		FROM payroll_rules WHERE company_id=$1 ORDER BY priority, created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollRule, error) {
		var rule domain.PayrollRule
		var params []byte
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
			&rule.ConceptID, &rule.RuleType, &rule.Formula, &params, &rule.Priority, &rule.EffectiveFrom,
			&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		if err == nil && params != nil {
			json.Unmarshal(params, &rule.Parameters)
		}
		return rule, err
	})
}

func (r *Repository) GetActiveRulesByConceptIDs(ctx context.Context, companyID string, conceptIDs []string, date time.Time) ([]domain.PayrollRule, error) {
	q := `SELECT id,company_id,country,jurisdiction,agreement_id,category_id,concept_id,rule_type,formula,parameters,
		priority,effective_from,effective_to,version,active,created_by,created_at,updated_at
		FROM payroll_rules WHERE company_id=$1 AND concept_id=ANY($2) AND active=true AND effective_from<=$3 AND (effective_to IS NULL OR effective_to>=$3)
		ORDER BY concept_id, priority DESC`
	rows, err := r.pool.Query(ctx, q, companyID, conceptIDs, date)
	if err != nil {
		return nil, repoErr("GetActiveRulesByConceptIDs", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollRule, error) {
		var rule domain.PayrollRule
		var params []byte
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.Country, &rule.Jurisdiction, &rule.AgreementID, &rule.CategoryID,
			&rule.ConceptID, &rule.RuleType, &rule.Formula, &params, &rule.Priority, &rule.EffectiveFrom,
			&rule.EffectiveTo, &rule.Version, &rule.Active, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		if err == nil && params != nil {
			json.Unmarshal(params, &rule.Parameters)
		}
		return rule, err
	})
}
